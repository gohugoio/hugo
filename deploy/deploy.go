// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nodeploy
// +build !nodeploy

package deploy

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deploy/deployconfig"
	"github.com/gohugoio/hugo/media"
	"github.com/spf13/afero"
	"golang.org/x/text/unicode/norm"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob" // import
	_ "gocloud.dev/blob/gcsblob"  // import
	_ "gocloud.dev/blob/s3blob"   // import
	"gocloud.dev/gcerrors"
)

// Deployer supports deploying the site to target cloud providers.
type Deployer struct {
	localFs afero.Fs
	bucket  *blob.Bucket

	mediaTypes media.Types // Hugo's MediaType to guess ContentType
	quiet      bool        // true reduces STDOUT // TODO(bep) remove, this is a global feature.

	cfg    deployconfig.DeployConfig
	logger loggers.Logger

	target *deployconfig.Target // the target to deploy to

	// For tests...
	summary deploySummary // summary of latest Deploy results
}

type deploySummary struct {
	NumLocal, NumRemote, NumUploads, NumDeletes int
}

const metaMD5Hash = "md5chksum" // the meta key to store md5hash in

// New constructs a new *Deployer.
func New(cfg config.AllProvider, logger loggers.Logger, localFs afero.Fs) (*Deployer, error) {
	dcfg := cfg.GetConfigSection(deployconfig.DeploymentConfigKey).(deployconfig.DeployConfig)
	targetName := dcfg.Target

	if len(dcfg.Targets) == 0 {
		return nil, errors.New("no deployment targets found")
	}
	mediaTypes := cfg.GetConfigSection("mediaTypes").(media.Types)

	// Find the target to deploy to.
	var tgt *deployconfig.Target
	if targetName == "" {
		// Default to the first target.
		tgt = dcfg.Targets[0]
	} else {
		for _, t := range dcfg.Targets {
			if t.Name == targetName {
				tgt = t
			}
		}
		if tgt == nil {
			return nil, fmt.Errorf("deployment target %q not found", targetName)
		}
	}

	return &Deployer{
		localFs:    localFs,
		target:     tgt,
		quiet:      cfg.BuildExpired(),
		mediaTypes: mediaTypes,
		cfg:        dcfg,
	}, nil
}

func (d *Deployer) openBucket(ctx context.Context) (*blob.Bucket, error) {
	if d.bucket != nil {
		return d.bucket, nil
	}
	d.logger.Printf("Deploying to target %q (%s)\n", d.target.Name, d.target.URL)
	return blob.OpenBucket(ctx, d.target.URL)
}

// Deploy deploys the site to a target.
func (d *Deployer) Deploy(ctx context.Context) error {
	if d.logger == nil {
		d.logger = loggers.NewDefault()
	}

	bucket, err := d.openBucket(ctx)
	if err != nil {
		return err
	}

	if d.cfg.Workers <= 0 {
		d.cfg.Workers = 10
	}

	// Load local files from the source directory.
	var include, exclude glob.Glob
	var mappath func(string) string
	if d.target != nil {
		include, exclude = d.target.IncludeGlob, d.target.ExcludeGlob
		if d.target.StripIndexHTML {
			mappath = stripIndexHTML
		}
	}
	local, err := d.walkLocal(d.localFs, d.cfg.Matchers, include, exclude, d.mediaTypes, mappath)
	if err != nil {
		return err
	}
	d.logger.Infof("Found %d local files.\n", len(local))
	d.summary.NumLocal = len(local)

	// Load remote files from the target.
	remote, err := d.walkRemote(ctx, bucket, include, exclude)
	if err != nil {
		return err
	}
	d.logger.Infof("Found %d remote files.\n", len(remote))
	d.summary.NumRemote = len(remote)

	// Diff local vs remote to see what changes need to be applied.
	uploads, deletes := d.findDiffs(local, remote, d.cfg.Force)
	d.summary.NumUploads = len(uploads)
	d.summary.NumDeletes = len(deletes)
	if len(uploads)+len(deletes) == 0 {
		if !d.quiet {
			d.logger.Println("No changes required.")
		}
		return nil
	}
	if !d.quiet {
		d.logger.Println(summarizeChanges(uploads, deletes))
	}

	// Ask for confirmation before proceeding.
	if d.cfg.Confirm && !d.cfg.DryRun {
		fmt.Printf("Continue? (Y/n) ")
		var confirm string
		if _, err := fmt.Scanln(&confirm); err != nil {
			return err
		}
		if confirm != "" && confirm[0] != 'y' && confirm[0] != 'Y' {
			return errors.New("aborted")
		}
	}

	// Order the uploads. They are organized in groups; all uploads in a group
	// must be complete before moving on to the next group.
	uploadGroups := applyOrdering(d.cfg.Ordering, uploads)

	nParallel := d.cfg.Workers
	var errs []error
	var errMu sync.Mutex // protects errs

	for _, uploads := range uploadGroups {
		// Short-circuit for an empty group.
		if len(uploads) == 0 {
			continue
		}

		// Within the group, apply uploads in parallel.
		sem := make(chan struct{}, nParallel)
		for _, upload := range uploads {
			if d.cfg.DryRun {
				if !d.quiet {
					d.logger.Printf("[DRY RUN] Would upload: %v\n", upload)
				}
				continue
			}

			sem <- struct{}{}
			go func(upload *fileToUpload) {
				if err := d.doSingleUpload(ctx, bucket, upload); err != nil {
					errMu.Lock()
					defer errMu.Unlock()
					errs = append(errs, err)
				}
				<-sem
			}(upload)
		}
		// Wait for all uploads in the group to finish.
		for n := nParallel; n > 0; n-- {
			sem <- struct{}{}
		}
	}

	if d.cfg.MaxDeletes != -1 && len(deletes) > d.cfg.MaxDeletes {
		d.logger.Warnf("Skipping %d deletes because it is more than --maxDeletes (%d). If this is expected, set --maxDeletes to a larger number, or -1 to disable this check.\n", len(deletes), d.cfg.MaxDeletes)
		d.summary.NumDeletes = 0
	} else {
		// Apply deletes in parallel.
		sort.Slice(deletes, func(i, j int) bool { return deletes[i] < deletes[j] })
		sem := make(chan struct{}, nParallel)
		for _, del := range deletes {
			if d.cfg.DryRun {
				if !d.quiet {
					d.logger.Printf("[DRY RUN] Would delete %s\n", del)
				}
				continue
			}
			sem <- struct{}{}
			go func(del string) {
				d.logger.Infof("Deleting %s...\n", del)
				if err := bucket.Delete(ctx, del); err != nil {
					if gcerrors.Code(err) == gcerrors.NotFound {
						d.logger.Warnf("Failed to delete %q because it wasn't found: %v", del, err)
					} else {
						errMu.Lock()
						defer errMu.Unlock()
						errs = append(errs, err)
					}
				}
				<-sem
			}(del)
		}
		// Wait for all deletes to finish.
		for n := nParallel; n > 0; n-- {
			sem <- struct{}{}
		}
	}

	if len(errs) > 0 {
		if !d.quiet {
			d.logger.Printf("Encountered %d errors.\n", len(errs))
		}
		return errs[0]
	}
	if !d.quiet {
		d.logger.Println("Success!")
	}

	if d.cfg.InvalidateCDN {
		if d.target.CloudFrontDistributionID != "" {
			if d.cfg.DryRun {
				if !d.quiet {
					d.logger.Printf("[DRY RUN] Would invalidate CloudFront CDN with ID %s\n", d.target.CloudFrontDistributionID)
				}
			} else {
				d.logger.Println("Invalidating CloudFront CDN...")
				if err := InvalidateCloudFront(ctx, d.target); err != nil {
					d.logger.Printf("Failed to invalidate CloudFront CDN: %v\n", err)
					return err
				}
			}
		}
		if d.target.GoogleCloudCDNOrigin != "" {
			if d.cfg.DryRun {
				if !d.quiet {
					d.logger.Printf("[DRY RUN] Would invalidate Google Cloud CDN with origin %s\n", d.target.GoogleCloudCDNOrigin)
				}
			} else {
				d.logger.Println("Invalidating Google Cloud CDN...")
				if err := InvalidateGoogleCloudCDN(ctx, d.target.GoogleCloudCDNOrigin); err != nil {
					d.logger.Printf("Failed to invalidate Google Cloud CDN: %v\n", err)
					return err
				}
			}
		}
		d.logger.Println("Success!")
	}
	return nil
}

// summarizeChanges creates a text description of the proposed changes.
func summarizeChanges(uploads []*fileToUpload, deletes []string) string {
	uploadSize := int64(0)
	for _, u := range uploads {
		uploadSize += u.Local.UploadSize
	}
	return fmt.Sprintf("Identified %d file(s) to upload, totaling %s, and %d file(s) to delete.", len(uploads), humanize.Bytes(uint64(uploadSize)), len(deletes))
}

// doSingleUpload executes a single file upload.
func (d *Deployer) doSingleUpload(ctx context.Context, bucket *blob.Bucket, upload *fileToUpload) error {
	d.logger.Infof("Uploading %v...\n", upload)
	opts := &blob.WriterOptions{
		CacheControl:    upload.Local.CacheControl(),
		ContentEncoding: upload.Local.ContentEncoding(),
		ContentType:     upload.Local.ContentType(),
		Metadata:        map[string]string{metaMD5Hash: hex.EncodeToString(upload.Local.MD5())},
	}
	w, err := bucket.NewWriter(ctx, upload.Local.SlashPath, opts)
	if err != nil {
		return err
	}
	r, err := upload.Local.Reader()
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// localFile represents a local file from the source. Use newLocalFile to
// construct one.
type localFile struct {
	// NativePath is the native path to the file (using file.Separator).
	NativePath string
	// SlashPath is NativePath converted to use /.
	SlashPath string
	// UploadSize is the size of the content to be uploaded. It may not
	// be the same as the local file size if the content will be
	// gzipped before upload.
	UploadSize int64

	fs         afero.Fs
	matcher    *deployconfig.Matcher
	md5        []byte       // cache
	gzipped    bytes.Buffer // cached of gzipped contents if gzipping
	mediaTypes media.Types
}

// newLocalFile initializes a *localFile.
func newLocalFile(fs afero.Fs, nativePath, slashpath string, m *deployconfig.Matcher, mt media.Types) (*localFile, error) {
	f, err := fs.Open(nativePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lf := &localFile{
		NativePath: nativePath,
		SlashPath:  slashpath,
		fs:         fs,
		matcher:    m,
		mediaTypes: mt,
	}
	if m != nil && m.Gzip {
		// We're going to gzip the content. Do it once now, and cache the result
		// in gzipped. The UploadSize is the size of the gzipped content.
		gz := gzip.NewWriter(&lf.gzipped)
		if _, err := io.Copy(gz, f); err != nil {
			return nil, err
		}
		if err := gz.Close(); err != nil {
			return nil, err
		}
		lf.UploadSize = int64(lf.gzipped.Len())
	} else {
		// Raw content. Just get the UploadSize.
		info, err := f.Stat()
		if err != nil {
			return nil, err
		}
		lf.UploadSize = info.Size()
	}
	return lf, nil
}

// Reader returns an io.ReadCloser for reading the content to be uploaded.
// The caller must call Close on the returned ReaderCloser.
// The reader content may not be the same as the local file content due to
// gzipping.
func (lf *localFile) Reader() (io.ReadCloser, error) {
	if lf.matcher != nil && lf.matcher.Gzip {
		// We've got the gzipped contents cached in gzipped.
		// Note: we can't use lf.gzipped directly as a Reader, since we it discards
		// data after it is read, and we may read it more than once.
		return io.NopCloser(bytes.NewReader(lf.gzipped.Bytes())), nil
	}
	// Not expected to fail since we did it successfully earlier in newLocalFile,
	// but could happen due to changes in the underlying filesystem.
	return lf.fs.Open(lf.NativePath)
}

// CacheControl returns the Cache-Control header to use for lf, based on the
// first matching matcher (if any).
func (lf *localFile) CacheControl() string {
	if lf.matcher == nil {
		return ""
	}
	return lf.matcher.CacheControl
}

// ContentEncoding returns the Content-Encoding header to use for lf, based
// on the matcher's Content-Encoding and Gzip fields.
func (lf *localFile) ContentEncoding() string {
	if lf.matcher == nil {
		return ""
	}
	if lf.matcher.Gzip {
		return "gzip"
	}
	return lf.matcher.ContentEncoding
}

// ContentType returns the Content-Type header to use for lf.
// It first checks if there's a Content-Type header configured via a matching
// matcher; if not, it tries to generate one based on the filename extension.
// If this fails, the Content-Type will be the empty string. In this case, Go
// Cloud will automatically try to infer a Content-Type based on the file
// content.
func (lf *localFile) ContentType() string {
	if lf.matcher != nil && lf.matcher.ContentType != "" {
		return lf.matcher.ContentType
	}

	ext := filepath.Ext(lf.NativePath)
	if mimeType, _, found := lf.mediaTypes.GetFirstBySuffix(strings.TrimPrefix(ext, ".")); found {
		return mimeType.Type
	}

	return mime.TypeByExtension(ext)
}

// Force returns true if the file should be forced to re-upload based on the
// matching matcher.
func (lf *localFile) Force() bool {
	return lf.matcher != nil && lf.matcher.Force
}

// MD5 returns an MD5 hash of the content to be uploaded.
func (lf *localFile) MD5() []byte {
	if len(lf.md5) > 0 {
		return lf.md5
	}
	h := md5.New()
	r, err := lf.Reader()
	if err != nil {
		return nil
	}
	defer r.Close()
	if _, err := io.Copy(h, r); err != nil {
		return nil
	}
	lf.md5 = h.Sum(nil)
	return lf.md5
}

// knownHiddenDirectory checks if the specified name is a well known
// hidden directory.
func knownHiddenDirectory(name string) bool {
	knownDirectories := []string{
		".well-known",
	}

	for _, dir := range knownDirectories {
		if name == dir {
			return true
		}
	}
	return false
}

// walkLocal walks the source directory and returns a flat list of files,
// using localFile.SlashPath as the map keys.
func (d *Deployer) walkLocal(fs afero.Fs, matchers []*deployconfig.Matcher, include, exclude glob.Glob, mediaTypes media.Types, mappath func(string) string) (map[string]*localFile, error) {
	retval := map[string]*localFile{}
	err := afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden directories.
			if path != "" && strings.HasPrefix(info.Name(), ".") {
				// Except for specific hidden directories
				if !knownHiddenDirectory(info.Name()) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// .DS_Store is an internal MacOS attribute file; skip it.
		if info.Name() == ".DS_Store" {
			return nil
		}

		// When a file system is HFS+, its filepath is in NFD form.
		if runtime.GOOS == "darwin" {
			path = norm.NFC.String(path)
		}

		// Check include/exclude matchers.
		slashpath := filepath.ToSlash(path)
		if include != nil && !include.Match(slashpath) {
			d.logger.Infof("  dropping %q due to include\n", slashpath)
			return nil
		}
		if exclude != nil && exclude.Match(slashpath) {
			d.logger.Infof("  dropping %q due to exclude\n", slashpath)
			return nil
		}

		// Find the first matching matcher (if any).
		var m *deployconfig.Matcher
		for _, cur := range matchers {
			if cur.Matches(slashpath) {
				m = cur
				break
			}
		}
		// Apply any additional modifications to the local path, to map it to
		// the remote path.
		if mappath != nil {
			slashpath = mappath(slashpath)
		}
		lf, err := newLocalFile(fs, path, slashpath, m, mediaTypes)
		if err != nil {
			return err
		}
		retval[lf.SlashPath] = lf
		return nil
	})
	if err != nil {
		return nil, err
	}
	return retval, nil
}

// stripIndexHTML remaps keys matching "<dir>/index.html" to "<dir>/".
func stripIndexHTML(slashpath string) string {
	const suffix = "/index.html"
	if strings.HasSuffix(slashpath, suffix) {
		return slashpath[:len(slashpath)-len(suffix)+1]
	}
	return slashpath
}

// walkRemote walks the target bucket and returns a flat list.
func (d *Deployer) walkRemote(ctx context.Context, bucket *blob.Bucket, include, exclude glob.Glob) (map[string]*blob.ListObject, error) {
	retval := map[string]*blob.ListObject{}
	iter := bucket.List(nil)
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// Check include/exclude matchers.
		if include != nil && !include.Match(obj.Key) {
			d.logger.Infof("  remote dropping %q due to include\n", obj.Key)
			continue
		}
		if exclude != nil && exclude.Match(obj.Key) {
			d.logger.Infof("  remote dropping %q due to exclude\n", obj.Key)
			continue
		}
		// If the remote didn't give us an MD5, use remote attributes MD5, if that doesn't exist compute one.
		// This can happen for some providers (e.g., fileblob, which uses the
		// local filesystem), but not for the most common Cloud providers
		// (S3, GCS, Azure). Although, it can happen for S3 if the blob was uploaded
		// via a multi-part upload.
		// Although it's unfortunate to have to read the file, it's likely better
		// than assuming a delta and re-uploading it.
		if len(obj.MD5) == 0 {
			var attrMD5 []byte
			attrs, err := bucket.Attributes(ctx, obj.Key)
			if err == nil {
				md5String, exists := attrs.Metadata[metaMD5Hash]
				if exists {
					attrMD5, _ = hex.DecodeString(md5String)
				}
			}
			if len(attrMD5) == 0 {
				r, err := bucket.NewReader(ctx, obj.Key, nil)
				if err == nil {
					h := md5.New()
					if _, err := io.Copy(h, r); err == nil {
						obj.MD5 = h.Sum(nil)
					}
					r.Close()
				}
			} else {
				obj.MD5 = attrMD5
			}
		}
		retval[obj.Key] = obj
	}
	return retval, nil
}

// uploadReason is an enum of reasons why a file must be uploaded.
type uploadReason string

const (
	reasonUnknown    uploadReason = "unknown"
	reasonNotFound   uploadReason = "not found at target"
	reasonForce      uploadReason = "--force"
	reasonSize       uploadReason = "size differs"
	reasonMD5Differs uploadReason = "md5 differs"
	reasonMD5Missing uploadReason = "remote md5 missing"
)

// fileToUpload represents a single local file that should be uploaded to
// the target.
type fileToUpload struct {
	Local  *localFile
	Reason uploadReason
}

func (u *fileToUpload) String() string {
	details := []string{humanize.Bytes(uint64(u.Local.UploadSize))}
	if s := u.Local.CacheControl(); s != "" {
		details = append(details, fmt.Sprintf("Cache-Control: %q", s))
	}
	if s := u.Local.ContentEncoding(); s != "" {
		details = append(details, fmt.Sprintf("Content-Encoding: %q", s))
	}
	if s := u.Local.ContentType(); s != "" {
		details = append(details, fmt.Sprintf("Content-Type: %q", s))
	}
	return fmt.Sprintf("%s (%s): %v", u.Local.SlashPath, strings.Join(details, ", "), u.Reason)
}

// findDiffs diffs localFiles vs remoteFiles to see what changes should be
// applied to the remote target. It returns a slice of *fileToUpload and a
// slice of paths for files to delete.
func (d *Deployer) findDiffs(localFiles map[string]*localFile, remoteFiles map[string]*blob.ListObject, force bool) ([]*fileToUpload, []string) {
	var uploads []*fileToUpload
	var deletes []string

	found := map[string]bool{}
	for path, lf := range localFiles {
		upload := false
		reason := reasonUnknown

		if remoteFile, ok := remoteFiles[path]; ok {
			// The file exists in remote. Let's see if we need to upload it anyway.

			// TODO: We don't register a diff if the metadata (e.g., Content-Type
			// header) has changed. This would be difficult/expensive to detect; some
			// providers return metadata along with their "List" result, but others
			// (notably AWS S3) do not, so gocloud.dev's blob.Bucket doesn't expose
			// it in the list result. It would require a separate request per blob
			// to fetch. At least for now, we work around this by documenting it and
			// providing a "force" flag (to re-upload everything) and a "force" bool
			// per matcher (to re-upload all files in a matcher whose headers may have
			// changed).
			// Idea: extract a sample set of 1 file per extension + 1 file per matcher
			// and check those files?
			if force {
				upload = true
				reason = reasonForce
			} else if lf.Force() {
				upload = true
				reason = reasonForce
			} else if lf.UploadSize != remoteFile.Size {
				upload = true
				reason = reasonSize
			} else if len(remoteFile.MD5) == 0 {
				// This shouldn't happen unless the remote didn't give us an MD5 hash
				// from List, AND we failed to compute one by reading the remote file.
				// Default to considering the files different.
				upload = true
				reason = reasonMD5Missing
			} else if !bytes.Equal(lf.MD5(), remoteFile.MD5) {
				upload = true
				reason = reasonMD5Differs
			}
			found[path] = true
		} else {
			// The file doesn't exist in remote.
			upload = true
			reason = reasonNotFound
		}
		if upload {
			d.logger.Debugf("%s needs to be uploaded: %v\n", path, reason)
			uploads = append(uploads, &fileToUpload{lf, reason})
		} else {
			d.logger.Debugf("%s exists at target and does not need to be uploaded", path)
		}
	}

	// Remote files that weren't found locally should be deleted.
	for path := range remoteFiles {
		if !found[path] {
			deletes = append(deletes, path)
		}
	}
	return uploads, deletes
}

// applyOrdering returns an ordered slice of slices of uploads.
//
// The returned slice will have length len(ordering)+1.
//
// The subslice at index i, for i = 0 ... len(ordering)-1, will have all of the
// uploads whose Local.SlashPath matched the regex at ordering[i] (but not any
// previous ordering regex).
// The subslice at index len(ordering) will have the remaining uploads that
// didn't match any ordering regex.
//
// The subslices are sorted by Local.SlashPath.
func applyOrdering(ordering []*regexp.Regexp, uploads []*fileToUpload) [][]*fileToUpload {
	// Sort the whole slice by Local.SlashPath first.
	sort.Slice(uploads, func(i, j int) bool { return uploads[i].Local.SlashPath < uploads[j].Local.SlashPath })

	retval := make([][]*fileToUpload, len(ordering)+1)
	for _, u := range uploads {
		matched := false
		for i, re := range ordering {
			if re.MatchString(u.Local.SlashPath) {
				retval[i] = append(retval[i], u)
				matched = true
				break
			}
		}
		if !matched {
			retval[len(ordering)] = append(retval[len(ordering)], u)
		}
	}
	return retval
}
