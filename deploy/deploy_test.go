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
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deploy/deployconfig"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/memblob"
)

func TestFindDiffs(t *testing.T) {
	hash1 := []byte("hash 1")
	hash2 := []byte("hash 2")
	makeLocal := func(path string, size int64, hash []byte) *localFile {
		return &localFile{NativePath: path, SlashPath: filepath.ToSlash(path), UploadSize: size, md5: hash}
	}
	makeRemote := func(path string, size int64, hash []byte) *blob.ListObject {
		return &blob.ListObject{Key: path, Size: size, MD5: hash}
	}

	tests := []struct {
		Description string
		Local       []*localFile
		Remote      []*blob.ListObject
		Force       bool
		WantUpdates []*fileToUpload
		WantDeletes []string
	}{
		{
			Description: "empty -> no diffs",
		},
		{
			Description: "local == remote -> no diffs",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 3, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
				makeRemote("ccc", 3, hash2),
			},
		},
		{
			Description: "local w/ separators == remote -> no diffs",
			Local: []*localFile{
				makeLocal(filepath.Join("aaa", "aaa"), 1, hash1),
				makeLocal(filepath.Join("bbb", "bbb"), 2, hash1),
				makeLocal(filepath.Join("ccc", "ccc"), 3, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa/aaa", 1, hash1),
				makeRemote("bbb/bbb", 2, hash1),
				makeRemote("ccc/ccc", 3, hash2),
			},
		},
		{
			Description: "local == remote with force flag true -> diffs",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 3, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
				makeRemote("ccc", 3, hash2),
			},
			Force: true,
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonForce},
				{makeLocal("bbb", 2, nil), reasonForce},
				{makeLocal("ccc", 3, nil), reasonForce},
			},
		},
		{
			Description: "local == remote with route.Force true -> diffs",
			Local: []*localFile{
				{NativePath: "aaa", SlashPath: "aaa", UploadSize: 1, matcher: &deployconfig.Matcher{Force: true}, md5: hash1},
				makeLocal("bbb", 2, hash1),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonForce},
			},
		},
		{
			Description: "extra local file -> upload",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("bbb", 2, nil), reasonNotFound},
			},
		},
		{
			Description: "extra remote file -> delete",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash2),
			},
			WantDeletes: []string{"bbb"},
		},
		{
			Description: "diffs in size or md5 -> upload",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 1, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, nil),
				makeRemote("bbb", 1, hash1),
				makeRemote("ccc", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonMD5Missing},
				{makeLocal("bbb", 2, nil), reasonSize},
				{makeLocal("ccc", 1, nil), reasonMD5Differs},
			},
		},
		{
			Description: "mix of updates and deletes",
			Local: []*localFile{
				makeLocal("same", 1, hash1),
				makeLocal("updated", 2, hash1),
				makeLocal("updated2", 1, hash2),
				makeLocal("new", 1, hash1),
				makeLocal("new2", 2, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("same", 1, hash1),
				makeRemote("updated", 1, hash1),
				makeRemote("updated2", 1, hash1),
				makeRemote("stale", 1, hash1),
				makeRemote("stale2", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("new", 1, nil), reasonNotFound},
				{makeLocal("new2", 2, nil), reasonNotFound},
				{makeLocal("updated", 2, nil), reasonSize},
				{makeLocal("updated2", 1, nil), reasonMD5Differs},
			},
			WantDeletes: []string{"stale", "stale2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			local := map[string]*localFile{}
			for _, l := range tc.Local {
				local[l.SlashPath] = l
			}
			remote := map[string]*blob.ListObject{}
			for _, r := range tc.Remote {
				remote[r.Key] = r
			}
			d := newDeployer()
			gotUpdates, gotDeletes := d.findDiffs(local, remote, tc.Force)
			gotUpdates = applyOrdering(nil, gotUpdates)[0]
			sort.Slice(gotDeletes, func(i, j int) bool { return gotDeletes[i] < gotDeletes[j] })
			if diff := cmp.Diff(gotUpdates, tc.WantUpdates, cmpopts.IgnoreUnexported(localFile{})); diff != "" {
				t.Errorf("updates differ:\n%s", diff)
			}
			if diff := cmp.Diff(gotDeletes, tc.WantDeletes); diff != "" {
				t.Errorf("deletes differ:\n%s", diff)
			}
		})
	}
}

func TestWalkLocal(t *testing.T) {
	tests := map[string]struct {
		Given   []string
		Expect  []string
		MapPath func(string) string
	}{
		"Empty": {
			Given:  []string{},
			Expect: []string{},
		},
		"Normal": {
			Given:  []string{"file.txt", "normal_dir/file.txt"},
			Expect: []string{"file.txt", "normal_dir/file.txt"},
		},
		"Hidden": {
			Given:  []string{"file.txt", ".hidden_dir/file.txt", "normal_dir/file.txt"},
			Expect: []string{"file.txt", "normal_dir/file.txt"},
		},
		"Well Known": {
			Given:  []string{"file.txt", ".hidden_dir/file.txt", ".well-known/file.txt"},
			Expect: []string{"file.txt", ".well-known/file.txt"},
		},
		"StripIndexHTML": {
			Given:   []string{"index.html", "file.txt", "dir/index.html", "dir/file.txt"},
			Expect:  []string{"index.html", "file.txt", "dir/", "dir/file.txt"},
			MapPath: stripIndexHTML,
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for _, name := range tc.Given {
				dir, _ := path.Split(name)
				if dir != "" {
					if err := fs.MkdirAll(dir, 0o755); err != nil {
						t.Fatal(err)
					}
				}
				if fd, err := fs.Create(name); err != nil {
					t.Fatal(err)
				} else {
					fd.Close()
				}
			}
			d := newDeployer()
			if got, err := d.walkLocal(fs, nil, nil, nil, media.DefaultTypes, tc.MapPath); err != nil {
				t.Fatal(err)
			} else {
				expect := map[string]any{}
				for _, path := range tc.Expect {
					if _, ok := got[path]; !ok {
						t.Errorf("expected %q in results, but was not found", path)
					}
					expect[path] = nil
				}
				for path := range got {
					if _, ok := expect[path]; !ok {
						t.Errorf("got %q in results unexpectedly", path)
					}
				}
			}
		})
	}
}

func TestStripIndexHTML(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Output string
	}{
		"Unmapped": {Input: "normal_file.txt", Output: "normal_file.txt"},
		"Stripped": {Input: "directory/index.html", Output: "directory/"},
		"NoSlash":  {Input: "prefix_index.html", Output: "prefix_index.html"},
		"Root":     {Input: "index.html", Output: "index.html"},
	}
	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			got := stripIndexHTML(tc.Input)
			if got != tc.Output {
				t.Errorf("got %q, expect %q", got, tc.Output)
			}
		})
	}
}

func TestStripIndexHTMLMatcher(t *testing.T) {
	// StripIndexHTML should not affect matchers.
	fs := afero.NewMemMapFs()
	if err := fs.Mkdir("dir", 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.html", "dir/index.html", "file.txt"} {
		if fd, err := fs.Create(name); err != nil {
			t.Fatal(err)
		} else {
			fd.Close()
		}
	}
	d := newDeployer()
	const pattern = `\.html$`
	matcher := &deployconfig.Matcher{Pattern: pattern, Gzip: true, Re: regexp.MustCompile(pattern)}
	if got, err := d.walkLocal(fs, []*deployconfig.Matcher{matcher}, nil, nil, media.DefaultTypes, stripIndexHTML); err != nil {
		t.Fatal(err)
	} else {
		for _, name := range []string{"index.html", "dir/"} {
			lf := got[name]
			if lf == nil {
				t.Errorf("missing file %q", name)
			} else if lf.matcher == nil {
				t.Errorf("file %q has nil matcher, expect %q", name, pattern)
			}
		}
		const name = "file.txt"
		lf := got[name]
		if lf == nil {
			t.Errorf("missing file %q", name)
		} else if lf.matcher != nil {
			t.Errorf("file %q has matcher %q, expect nil", name, lf.matcher.Pattern)
		}
	}
}

func TestLocalFile(t *testing.T) {
	const (
		content = "hello world!"
	)
	contentBytes := []byte(content)
	contentLen := int64(len(contentBytes))
	contentMD5 := md5.Sum(contentBytes)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(contentBytes); err != nil {
		t.Fatal(err)
	}
	gz.Close()
	gzBytes := buf.Bytes()
	gzLen := int64(len(gzBytes))
	gzMD5 := md5.Sum(gzBytes)

	tests := []struct {
		Description         string
		Path                string
		Matcher             *deployconfig.Matcher
		MediaTypesConfig    map[string]any
		WantContent         []byte
		WantSize            int64
		WantMD5             []byte
		WantContentType     string // empty string is always OK, since content type detection is OS-specific
		WantCacheControl    string
		WantContentEncoding string
	}{
		{
			Description: "file with no suffix",
			Path:        "foo",
			WantContent: contentBytes,
			WantSize:    contentLen,
			WantMD5:     contentMD5[:],
		},
		{
			Description: "file with .txt suffix",
			Path:        "foo.txt",
			WantContent: contentBytes,
			WantSize:    contentLen,
			WantMD5:     contentMD5[:],
		},
		{
			Description:      "CacheControl from matcher",
			Path:             "foo.txt",
			Matcher:          &deployconfig.Matcher{CacheControl: "max-age=630720000"},
			WantContent:      contentBytes,
			WantSize:         contentLen,
			WantMD5:          contentMD5[:],
			WantCacheControl: "max-age=630720000",
		},
		{
			Description:         "ContentEncoding from matcher",
			Path:                "foo.txt",
			Matcher:             &deployconfig.Matcher{ContentEncoding: "foobar"},
			WantContent:         contentBytes,
			WantSize:            contentLen,
			WantMD5:             contentMD5[:],
			WantContentEncoding: "foobar",
		},
		{
			Description:     "ContentType from matcher",
			Path:            "foo.txt",
			Matcher:         &deployconfig.Matcher{ContentType: "foo/bar"},
			WantContent:     contentBytes,
			WantSize:        contentLen,
			WantMD5:         contentMD5[:],
			WantContentType: "foo/bar",
		},
		{
			Description:         "gzipped content",
			Path:                "foo.txt",
			Matcher:             &deployconfig.Matcher{Gzip: true},
			WantContent:         gzBytes,
			WantSize:            gzLen,
			WantMD5:             gzMD5[:],
			WantContentEncoding: "gzip",
		},
		{
			Description: "Custom MediaType",
			Path:        "foo.hugo",
			MediaTypesConfig: map[string]any{
				"hugo/custom": map[string]any{
					"suffixes": []string{"hugo"},
				},
			},
			WantContent:     contentBytes,
			WantSize:        contentLen,
			WantMD5:         contentMD5[:],
			WantContentType: "hugo/custom",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			fs := new(afero.MemMapFs)
			if err := afero.WriteFile(fs, tc.Path, []byte(content), os.ModePerm); err != nil {
				t.Fatal(err)
			}
			mediaTypes := media.DefaultTypes
			if len(tc.MediaTypesConfig) > 0 {
				mt, err := media.DecodeTypes(tc.MediaTypesConfig)
				if err != nil {
					t.Fatal(err)
				}
				mediaTypes = mt.Config
			}
			lf, err := newLocalFile(fs, tc.Path, filepath.ToSlash(tc.Path), tc.Matcher, mediaTypes)
			if err != nil {
				t.Fatal(err)
			}
			if got := lf.UploadSize; got != tc.WantSize {
				t.Errorf("got size %d want %d", got, tc.WantSize)
			}
			if got := lf.MD5(); !bytes.Equal(got, tc.WantMD5) {
				t.Errorf("got MD5 %x want %x", got, tc.WantMD5)
			}
			if got := lf.CacheControl(); got != tc.WantCacheControl {
				t.Errorf("got CacheControl %q want %q", got, tc.WantCacheControl)
			}
			if got := lf.ContentEncoding(); got != tc.WantContentEncoding {
				t.Errorf("got ContentEncoding %q want %q", got, tc.WantContentEncoding)
			}
			if tc.WantContentType != "" {
				if got := lf.ContentType(); got != tc.WantContentType {
					t.Errorf("got ContentType %q want %q", got, tc.WantContentType)
				}
			}
			// Verify the reader last to ensure the previous operations don't
			// interfere with it.
			r, err := lf.Reader()
			if err != nil {
				t.Fatal(err)
			}
			gotContent, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(gotContent, tc.WantContent) {
				t.Errorf("got content %q want %q", string(gotContent), string(tc.WantContent))
			}
			r.Close()
			// Verify we can read again.
			r, err = lf.Reader()
			if err != nil {
				t.Fatal(err)
			}
			gotContent, err = io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			r.Close()
			if !bytes.Equal(gotContent, tc.WantContent) {
				t.Errorf("got content %q want %q", string(gotContent), string(tc.WantContent))
			}
		})
	}
}

func TestOrdering(t *testing.T) {
	tests := []struct {
		Description string
		Uploads     []string
		Ordering    []*regexp.Regexp
		Want        [][]string
	}{
		{
			Description: "empty",
			Want:        [][]string{nil},
		},
		{
			Description: "no ordering",
			Uploads:     []string{"c", "b", "a", "d"},
			Want:        [][]string{{"a", "b", "c", "d"}},
		},
		{
			Description: "one ordering",
			Uploads:     []string{"db", "c", "b", "a", "da"},
			Ordering:    []*regexp.Regexp{regexp.MustCompile("^d")},
			Want:        [][]string{{"da", "db"}, {"a", "b", "c"}},
		},
		{
			Description: "two orderings",
			Uploads:     []string{"db", "c", "b", "a", "da"},
			Ordering: []*regexp.Regexp{
				regexp.MustCompile("^d"),
				regexp.MustCompile("^b"),
			},
			Want: [][]string{{"da", "db"}, {"b"}, {"a", "c"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			uploads := make([]*fileToUpload, len(tc.Uploads))
			for i, u := range tc.Uploads {
				uploads[i] = &fileToUpload{Local: &localFile{SlashPath: u}}
			}
			gotUploads := applyOrdering(tc.Ordering, uploads)
			var got [][]string
			for _, subslice := range gotUploads {
				var gotsubslice []string
				for _, u := range subslice {
					gotsubslice = append(gotsubslice, u.Local.SlashPath)
				}
				got = append(got, gotsubslice)
			}
			if diff := cmp.Diff(got, tc.Want); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type fileData struct {
	Name     string // name of the file
	Contents string // contents of the file
}

// initLocalFs initializes fs with some test files.
func initLocalFs(ctx context.Context, fs afero.Fs) ([]*fileData, error) {
	// The initial local filesystem.
	local := []*fileData{
		{"aaa", "aaa"},
		{"bbb", "bbb"},
		{"subdir/aaa", "subdir-aaa"},
		{"subdir/nested/aaa", "subdir-nested-aaa"},
		{"subdir2/bbb", "subdir2-bbb"},
	}
	if err := writeFiles(fs, local); err != nil {
		return nil, err
	}
	return local, nil
}

// fsTest represents an (afero.FS, Go CDK blob.Bucket) against which end-to-end
// tests can be run.
type fsTest struct {
	name   string
	fs     afero.Fs
	bucket *blob.Bucket
}

// initFsTests initializes a pair of tests for end-to-end test:
// 1. An in-memory afero.Fs paired with an in-memory Go CDK bucket.
// 2. A filesystem-based afero.Fs paired with an filesystem-based Go CDK bucket.
// It returns the pair of tests and a cleanup function.
func initFsTests(t *testing.T) []*fsTest {
	t.Helper()

	tmpfsdir := t.TempDir()
	tmpbucketdir := t.TempDir()

	memfs := afero.NewMemMapFs()
	membucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { membucket.Close() })

	filefs := hugofs.NewBasePathFs(afero.NewOsFs(), tmpfsdir)
	filebucket, err := fileblob.OpenBucket(tmpbucketdir, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { filebucket.Close() })

	tests := []*fsTest{
		{"mem", memfs, membucket},
		{"file", filefs, filebucket},
	}
	return tests
}

// TestEndToEndSync verifies that basic adds, updates, and deletes are working
// correctly.
func TestEndToEndSync(t *testing.T) {
	ctx := context.Background()
	tests := initFsTests(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			local, err := initLocalFs(ctx, test.fs)
			if err != nil {
				t.Fatal(err)
			}
			deployer := &Deployer{
				localFs:    test.fs,
				bucket:     test.bucket,
				mediaTypes: media.DefaultTypes,
				cfg:        deployconfig.DeployConfig{MaxDeletes: -1},
			}

			// Initial deployment should sync remote with local.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("initial deploy: failed: %v", err)
			}
			wantSummary := deploySummary{NumLocal: 5, NumRemote: 0, NumUploads: 5, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("initial deploy: got %v, want %v", deployer.summary, wantSummary)
			}
			if diff, err := verifyRemote(ctx, deployer.bucket, local); err != nil {
				t.Errorf("initial deploy: failed to verify remote: %v", err)
			} else if diff != "" {
				t.Errorf("initial deploy: remote snapshot doesn't match expected:\n%v", diff)
			}

			// A repeat deployment shouldn't change anything.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("no-op deploy: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 0, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("no-op deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// Make some changes to the local filesystem:
			// 1. Modify file [0].
			// 2. Delete file [1].
			// 3. Add a new file (sorted last).
			updatefd := local[0]
			updatefd.Contents = "new contents"
			deletefd := local[1]
			local = append(local[:1], local[2:]...) // removing deleted [1]
			newfd := &fileData{"zzz", "zzz"}
			local = append(local, newfd)
			if err := writeFiles(test.fs, []*fileData{updatefd, newfd}); err != nil {
				t.Fatal(err)
			}
			if err := test.fs.Remove(deletefd.Name); err != nil {
				t.Fatal(err)
			}

			// A deployment should apply those 3 changes.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy after changes: failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 2, NumDeletes: 1}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy after changes: got %v, want %v", deployer.summary, wantSummary)
			}
			if diff, err := verifyRemote(ctx, deployer.bucket, local); err != nil {
				t.Errorf("deploy after changes: failed to verify remote: %v", err)
			} else if diff != "" {
				t.Errorf("deploy after changes: remote snapshot doesn't match expected:\n%v", diff)
			}

			// Again, a repeat deployment shouldn't change anything.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("no-op deploy: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 0, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("no-op deploy: got %v, want %v", deployer.summary, wantSummary)
			}
		})
	}
}

// TestMaxDeletes verifies that the "maxDeletes" flag is working correctly.
func TestMaxDeletes(t *testing.T) {
	ctx := context.Background()
	tests := initFsTests(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			local, err := initLocalFs(ctx, test.fs)
			if err != nil {
				t.Fatal(err)
			}
			deployer := &Deployer{
				localFs:    test.fs,
				bucket:     test.bucket,
				mediaTypes: media.DefaultTypes,
				cfg:        deployconfig.DeployConfig{MaxDeletes: -1},
			}

			// Sync remote with local.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("initial deploy: failed: %v", err)
			}
			wantSummary := deploySummary{NumLocal: 5, NumRemote: 0, NumUploads: 5, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("initial deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// Delete two files, [1] and [2].
			if err := test.fs.Remove(local[1].Name); err != nil {
				t.Fatal(err)
			}
			if err := test.fs.Remove(local[2].Name); err != nil {
				t.Fatal(err)
			}

			// A deployment with maxDeletes=0 shouldn't change anything.
			deployer.cfg.MaxDeletes = 0
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 3, NumRemote: 5, NumUploads: 0, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// A deployment with maxDeletes=1 shouldn't change anything either.
			deployer.cfg.MaxDeletes = 1
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 3, NumRemote: 5, NumUploads: 0, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// A deployment with maxDeletes=2 should make the changes.
			deployer.cfg.MaxDeletes = 2
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 3, NumRemote: 5, NumUploads: 0, NumDeletes: 2}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// Delete two more files, [0] and [3].
			if err := test.fs.Remove(local[0].Name); err != nil {
				t.Fatal(err)
			}
			if err := test.fs.Remove(local[3].Name); err != nil {
				t.Fatal(err)
			}

			// A deployment with maxDeletes=-1 should make the changes.
			deployer.cfg.MaxDeletes = -1
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 1, NumRemote: 3, NumUploads: 0, NumDeletes: 2}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, wantSummary)
			}
		})
	}
}

// TestIncludeExclude verifies that the include/exclude options for targets work.
func TestIncludeExclude(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		Include string
		Exclude string
		Want    deploySummary
	}{
		{
			Want: deploySummary{NumLocal: 5, NumUploads: 5},
		},
		{
			Include: "**aaa",
			Want:    deploySummary{NumLocal: 3, NumUploads: 3},
		},
		{
			Include: "**bbb",
			Want:    deploySummary{NumLocal: 2, NumUploads: 2},
		},
		{
			Include: "aaa",
			Want:    deploySummary{NumLocal: 1, NumUploads: 1},
		},
		{
			Exclude: "**aaa",
			Want:    deploySummary{NumLocal: 2, NumUploads: 2},
		},
		{
			Exclude: "**bbb",
			Want:    deploySummary{NumLocal: 3, NumUploads: 3},
		},
		{
			Exclude: "aaa",
			Want:    deploySummary{NumLocal: 4, NumUploads: 4},
		},
		{
			Include: "**aaa",
			Exclude: "**nested**",
			Want:    deploySummary{NumLocal: 2, NumUploads: 2},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("include %q exclude %q", test.Include, test.Exclude), func(t *testing.T) {
			fsTests := initFsTests(t)
			fsTest := fsTests[1] // just do file-based test

			_, err := initLocalFs(ctx, fsTest.fs)
			if err != nil {
				t.Fatal(err)
			}
			tgt := &deployconfig.Target{
				Include: test.Include,
				Exclude: test.Exclude,
			}
			if err := tgt.ParseIncludeExclude(); err != nil {
				t.Error(err)
			}
			deployer := &Deployer{
				localFs: fsTest.fs,
				cfg:     deployconfig.DeployConfig{MaxDeletes: -1}, bucket: fsTest.bucket,
				target:     tgt,
				mediaTypes: media.DefaultTypes,
			}

			// Sync remote with local.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy: failed: %v", err)
			}
			if !cmp.Equal(deployer.summary, test.Want) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, test.Want)
			}
		})
	}
}

// TestIncludeExcludeRemoteDelete verifies deleted local files that don't match include/exclude patterns
// are not deleted on the remote.
func TestIncludeExcludeRemoteDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		Include string
		Exclude string
		Want    deploySummary
	}{
		{
			Want: deploySummary{NumLocal: 3, NumRemote: 5, NumUploads: 0, NumDeletes: 2},
		},
		{
			Include: "**aaa",
			Want:    deploySummary{NumLocal: 2, NumRemote: 3, NumUploads: 0, NumDeletes: 1},
		},
		{
			Include: "subdir/**",
			Want:    deploySummary{NumLocal: 1, NumRemote: 2, NumUploads: 0, NumDeletes: 1},
		},
		{
			Exclude: "**bbb",
			Want:    deploySummary{NumLocal: 2, NumRemote: 3, NumUploads: 0, NumDeletes: 1},
		},
		{
			Exclude: "bbb",
			Want:    deploySummary{NumLocal: 3, NumRemote: 4, NumUploads: 0, NumDeletes: 1},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("include %q exclude %q", test.Include, test.Exclude), func(t *testing.T) {
			fsTests := initFsTests(t)
			fsTest := fsTests[1] // just do file-based test

			local, err := initLocalFs(ctx, fsTest.fs)
			if err != nil {
				t.Fatal(err)
			}
			deployer := &Deployer{
				localFs: fsTest.fs,
				cfg:     deployconfig.DeployConfig{MaxDeletes: -1}, bucket: fsTest.bucket,
				mediaTypes: media.DefaultTypes,
			}

			// Initial sync to get the files on the remote
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy: failed: %v", err)
			}

			// Delete two files, [1] and [2].
			if err := fsTest.fs.Remove(local[1].Name); err != nil {
				t.Fatal(err)
			}
			if err := fsTest.fs.Remove(local[2].Name); err != nil {
				t.Fatal(err)
			}

			// Second sync
			tgt := &deployconfig.Target{
				Include: test.Include,
				Exclude: test.Exclude,
			}
			if err := tgt.ParseIncludeExclude(); err != nil {
				t.Error(err)
			}
			deployer.target = tgt
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy: failed: %v", err)
			}

			if !cmp.Equal(deployer.summary, test.Want) {
				t.Errorf("deploy: got %v, want %v", deployer.summary, test.Want)
			}
		})
	}
}

// TestCompression verifies that gzip compression works correctly.
// In particular, MD5 hashes must be of the compressed content.
func TestCompression(t *testing.T) {
	ctx := context.Background()

	tests := initFsTests(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			local, err := initLocalFs(ctx, test.fs)
			if err != nil {
				t.Fatal(err)
			}
			deployer := &Deployer{
				localFs:    test.fs,
				bucket:     test.bucket,
				cfg:        deployconfig.DeployConfig{MaxDeletes: -1, Matchers: []*deployconfig.Matcher{{Pattern: ".*", Gzip: true, Re: regexp.MustCompile(".*")}}},
				mediaTypes: media.DefaultTypes,
			}

			// Initial deployment should sync remote with local.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("initial deploy: failed: %v", err)
			}
			wantSummary := deploySummary{NumLocal: 5, NumRemote: 0, NumUploads: 5, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("initial deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// A repeat deployment shouldn't change anything.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("no-op deploy: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 0, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("no-op deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// Make an update to the local filesystem, on [1].
			updatefd := local[1]
			updatefd.Contents = "new contents"
			if err := writeFiles(test.fs, []*fileData{updatefd}); err != nil {
				t.Fatal(err)
			}

			// A deployment should apply the changes.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("deploy after changes: failed: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 1, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("deploy after changes: got %v, want %v", deployer.summary, wantSummary)
			}
		})
	}
}

// TestMatching verifies that matchers match correctly, and that the Force
// attribute for matcher works.
func TestMatching(t *testing.T) {
	ctx := context.Background()
	tests := initFsTests(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := initLocalFs(ctx, test.fs)
			if err != nil {
				t.Fatal(err)
			}
			deployer := &Deployer{
				localFs:    test.fs,
				bucket:     test.bucket,
				cfg:        deployconfig.DeployConfig{MaxDeletes: -1, Matchers: []*deployconfig.Matcher{{Pattern: "^subdir/aaa$", Force: true, Re: regexp.MustCompile("^subdir/aaa$")}}},
				mediaTypes: media.DefaultTypes,
			}

			// Initial deployment to sync remote with local.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("initial deploy: failed: %v", err)
			}
			wantSummary := deploySummary{NumLocal: 5, NumRemote: 0, NumUploads: 5, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("initial deploy: got %v, want %v", deployer.summary, wantSummary)
			}

			// A repeat deployment should upload a single file, the one that matched the Force matcher.
			// Note that matching happens based on the ToSlash form, so this matches
			// even on Windows.
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("no-op deploy with single force matcher: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 1, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("no-op deploy with single force matcher: got %v, want %v", deployer.summary, wantSummary)
			}

			// Repeat with a matcher that should now match 3 files.
			deployer.cfg.Matchers = []*deployconfig.Matcher{{Pattern: "aaa", Force: true, Re: regexp.MustCompile("aaa")}}
			if err := deployer.Deploy(ctx); err != nil {
				t.Errorf("no-op deploy with triple force matcher: %v", err)
			}
			wantSummary = deploySummary{NumLocal: 5, NumRemote: 5, NumUploads: 3, NumDeletes: 0}
			if !cmp.Equal(deployer.summary, wantSummary) {
				t.Errorf("no-op deploy with triple force matcher: got %v, want %v", deployer.summary, wantSummary)
			}
		})
	}
}

// writeFiles writes the files in fds to fd.
func writeFiles(fs afero.Fs, fds []*fileData) error {
	for _, fd := range fds {
		dir := path.Dir(fd.Name)
		if dir != "." {
			err := fs.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		f, err := fs.Create(fd.Name)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(fd.Contents)
		if err != nil {
			return err
		}
	}
	return nil
}

// verifyRemote that the current contents of bucket matches local.
// It returns an empty string if the contents matched, and a non-empty string
// capturing the diff if they didn't.
func verifyRemote(ctx context.Context, bucket *blob.Bucket, local []*fileData) (string, error) {
	var cur []*fileData
	iter := bucket.List(nil)
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		contents, err := bucket.ReadAll(ctx, obj.Key)
		if err != nil {
			return "", err
		}
		cur = append(cur, &fileData{obj.Key, string(contents)})
	}
	if cmp.Equal(cur, local) {
		return "", nil
	}
	diff := "got: \n"
	for _, f := range cur {
		diff += fmt.Sprintf("  %s: %s\n", f.Name, f.Contents)
	}
	diff += "want: \n"
	for _, f := range local {
		diff += fmt.Sprintf("  %s: %s\n", f.Name, f.Contents)
	}
	return diff, nil
}

func newDeployer() *Deployer {
	return &Deployer{
		logger: loggers.NewDefault(),
	}
}
