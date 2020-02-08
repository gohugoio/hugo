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

package resources

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/gohugoio/hugo/resources/internal"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
)

var (
	_ resource.ContentResource         = (*genericResource)(nil)
	_ resource.ReadSeekCloserResource  = (*genericResource)(nil)
	_ resource.Resource                = (*genericResource)(nil)
	_ resource.Source                  = (*genericResource)(nil)
	_ resource.Cloner                  = (*genericResource)(nil)
	_ resource.ResourcesLanguageMerger = (*resource.Resources)(nil)
	_ permalinker                      = (*genericResource)(nil)
	_ resource.Identifier              = (*genericResource)(nil)
	_ fileInfo                         = (*genericResource)(nil)
)

type ResourceSourceDescriptor struct {
	// TargetPaths is a callback to fetch paths's relative to its owner.
	TargetPaths func() page.TargetPaths

	// Need one of these to load the resource content.
	SourceFile         source.File
	OpenReadSeekCloser resource.OpenReadSeekCloser

	FileInfo os.FileInfo

	// If OpenReadSeekerCloser is not set, we use this to open the file.
	SourceFilename string

	Fs afero.Fs

	// The relative target filename without any language code.
	RelTargetFilename string

	// Any base paths prepended to the target path. This will also typically be the
	// language code, but setting it here means that it should not have any effect on
	// the permalink.
	// This may be several values. In multihost mode we may publish the same resources to
	// multiple targets.
	TargetBasePaths []string

	// Delay publishing until either Permalink or RelPermalink is called. Maybe never.
	LazyPublish bool
}

func (r ResourceSourceDescriptor) Filename() string {
	if r.SourceFile != nil {
		return r.SourceFile.Filename()
	}
	return r.SourceFilename
}

type ResourceTransformer interface {
	resource.Resource
	Transformer
}

type Transformer interface {
	Transform(...ResourceTransformation) (ResourceTransformer, error)
}

func NewFeatureNotAvailableTransformer(key string, elements ...interface{}) ResourceTransformation {
	return transformerNotAvailable{
		key: internal.NewResourceTransformationKey(key, elements...),
	}
}

type transformerNotAvailable struct {
	key internal.ResourceTransformationKey
}

func (t transformerNotAvailable) Transform(ctx *ResourceTransformationCtx) error {
	return herrors.ErrFeatureNotAvailable
}

func (t transformerNotAvailable) Key() internal.ResourceTransformationKey {
	return t.key
}

type baseResourceResource interface {
	resource.Cloner
	resource.ContentProvider
	resource.Resource
	resource.Identifier
}

type baseResourceInternal interface {
	resource.Source

	fileInfo
	metaAssigner
	targetPather

	ReadSeekCloser() (hugio.ReadSeekCloser, error)

	// Internal
	cloneWithUpdates(*transformationUpdate) (baseResource, error)
	tryTransformedFileCache(key string, u *transformationUpdate) io.ReadCloser

	specProvider
	getResourcePaths() *resourcePathDescriptor
	getTargetFilenames() []string
	openDestinationsForWriting() (io.WriteCloser, error)
	openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error)

	relTargetPathForRel(rel string, addBaseTargetPath, isAbs, isURL bool) string
}

type specProvider interface {
	getSpec() *Spec
}

type baseResource interface {
	baseResourceResource
	baseResourceInternal
}

type commonResource struct {
}

// Slice is not meant to be used externally. It's a bridge function
// for the template functions. See collections.Slice.
func (commonResource) Slice(in interface{}) (interface{}, error) {
	switch items := in.(type) {
	case resource.Resources:
		return items, nil
	case []interface{}:
		groups := make(resource.Resources, len(items))
		for i, v := range items {
			g, ok := v.(resource.Resource)
			if !ok {
				return nil, fmt.Errorf("type %T is not a Resource", v)
			}
			groups[i] = g
			{
			}
		}
		return groups, nil
	default:
		return nil, fmt.Errorf("invalid slice type %T", items)
	}
}

type dirFile struct {
	// This is the directory component with Unix-style slashes.
	dir string
	// This is the file component.
	file string
}

func (d dirFile) path() string {
	return path.Join(d.dir, d.file)
}

type fileInfo interface {
	getSourceFilename() string
	setSourceFilename(string)
	setSourceFs(afero.Fs)
	getFileInfo() hugofs.FileMetaInfo
	hash() (string, error)
	size() int
}

// genericResource represents a generic linkable resource.
type genericResource struct {
	*resourcePathDescriptor
	*resourceFileInfo
	*resourceContent

	spec *Spec

	title  string
	name   string
	params map[string]interface{}
	data   map[string]interface{}

	resourceType string
	mediaType    media.Type
}

func (l *genericResource) Clone() resource.Resource {
	return l.clone()
}

func (l *genericResource) Content() (interface{}, error) {
	if err := l.initContent(); err != nil {
		return nil, err
	}

	return l.content, nil
}

func (l *genericResource) Data() interface{} {
	return l.data
}

func (l *genericResource) Key() string {
	return l.RelPermalink()
}

func (l *genericResource) MediaType() media.Type {
	return l.mediaType
}

func (l *genericResource) setMediaType(mediaType media.Type) {
	l.mediaType = mediaType
}

func (l *genericResource) Name() string {
	return l.name
}

func (l *genericResource) Params() maps.Params {
	return l.params
}

func (l *genericResource) Permalink() string {
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(l.relTargetDirFile.path(), true), l.spec.BaseURL.HostURL())
}

func (l *genericResource) Publish() error {
	var err error
	l.publishInit.Do(func() {
		var fr hugio.ReadSeekCloser
		fr, err = l.ReadSeekCloser()
		if err != nil {
			return
		}
		defer fr.Close()

		var fw io.WriteCloser
		fw, err = helpers.OpenFilesForWriting(l.spec.BaseFs.PublishFs, l.getTargetFilenames()...)
		if err != nil {
			return
		}
		defer fw.Close()

		_, err = io.Copy(fw, fr)

	})

	return err
}

func (l *genericResource) RelPermalink() string {
	return l.relPermalinkFor(l.relTargetDirFile.path())
}

func (l *genericResource) ResourceType() string {
	return l.resourceType
}

func (l *genericResource) String() string {
	return fmt.Sprintf("Resource(%s: %s)", l.resourceType, l.name)
}

// Path is stored with Unix style slashes.
func (l *genericResource) TargetPath() string {
	return l.relTargetDirFile.path()
}

func (l *genericResource) Title() string {
	return l.title
}

func (l *genericResource) createBasePath(rel string, isURL bool) string {
	if l.targetPathBuilder == nil {
		return rel
	}
	tp := l.targetPathBuilder()

	if isURL {
		return path.Join(tp.SubResourceBaseLink, rel)
	}

	// TODO(bep) path
	return path.Join(filepath.ToSlash(tp.SubResourceBaseTarget), rel)
}

func (l *genericResource) initContent() error {
	var err error
	l.contentInit.Do(func() {
		var r hugio.ReadSeekCloser
		r, err = l.ReadSeekCloser()
		if err != nil {
			return
		}
		defer r.Close()

		var b []byte
		b, err = ioutil.ReadAll(r)
		if err != nil {
			return
		}

		l.content = string(b)
	})

	return err
}

func (l *genericResource) setName(name string) {
	l.name = name
}

func (l *genericResource) getResourcePaths() *resourcePathDescriptor {
	return l.resourcePathDescriptor
}

func (l *genericResource) getSpec() *Spec {
	return l.spec
}

func (l *genericResource) getTargetFilenames() []string {
	paths := l.relTargetPaths()
	for i, p := range paths {
		paths[i] = filepath.Clean(p)
	}
	return paths
}

func (l *genericResource) setTitle(title string) {
	l.title = title
}

func (r *genericResource) tryTransformedFileCache(key string, u *transformationUpdate) io.ReadCloser {
	fi, f, meta, found := r.spec.ResourceCache.getFromFile(key)
	if !found {
		return nil
	}
	u.sourceFilename = &fi.Name
	mt, _ := r.spec.MediaTypes.GetByType(meta.MediaTypeV)
	u.mediaType = mt
	u.data = meta.MetaData
	u.targetPath = meta.Target
	return f
}

func (r *genericResource) mergeData(in map[string]interface{}) {
	if len(in) == 0 {
		return
	}
	if r.data == nil {
		r.data = make(map[string]interface{})
	}
	for k, v := range in {
		if _, found := r.data[k]; !found {
			r.data[k] = v
		}
	}
}

func (rc *genericResource) cloneWithUpdates(u *transformationUpdate) (baseResource, error) {
	r := rc.clone()

	if u.content != nil {
		r.contentInit.Do(func() {
			r.content = *u.content
			r.openReadSeekerCloser = func() (hugio.ReadSeekCloser, error) {
				return hugio.NewReadSeekerNoOpCloserFromString(r.content), nil
			}
		})
	}

	r.mediaType = u.mediaType

	if u.sourceFilename != nil {
		r.setSourceFilename(*u.sourceFilename)
	}

	if u.sourceFs != nil {
		r.setSourceFs(u.sourceFs)
	}

	if u.targetPath == "" {
		return nil, errors.New("missing targetPath")
	}

	fpath, fname := path.Split(u.targetPath)
	r.resourcePathDescriptor.relTargetDirFile = dirFile{dir: fpath, file: fname}

	r.mergeData(u.data)

	return r, nil
}

func (l genericResource) clone() *genericResource {
	gi := *l.resourceFileInfo
	rp := *l.resourcePathDescriptor
	l.resourceFileInfo = &gi
	l.resourcePathDescriptor = &rp
	l.resourceContent = &resourceContent{}
	return &l
}

// returns an opened file or nil if nothing to write (it may already be published).
func (l *genericResource) openDestinationsForWriting() (w io.WriteCloser, err error) {

	l.publishInit.Do(func() {
		targetFilenames := l.getTargetFilenames()
		var changedFilenames []string

		// Fast path:
		// This is a processed version of the original;
		// check if it already existis at the destination.
		for _, targetFilename := range targetFilenames {
			if _, err := l.getSpec().BaseFs.PublishFs.Stat(targetFilename); err == nil {
				continue
			}

			changedFilenames = append(changedFilenames, targetFilename)
		}

		if len(changedFilenames) == 0 {
			return
		}

		w, err = helpers.OpenFilesForWriting(l.getSpec().BaseFs.PublishFs, changedFilenames...)

	})

	return

}

func (r *genericResource) openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error) {
	return helpers.OpenFilesForWriting(r.spec.BaseFs.PublishFs, r.relTargetPathsFor(relTargetPath)...)
}

func (l *genericResource) permalinkFor(target string) string {
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(target, true), l.spec.BaseURL.HostURL())
}

func (l *genericResource) relPermalinkFor(target string) string {
	return l.relPermalinkForRel(target, false)
}

func (l *genericResource) relPermalinkForRel(rel string, isAbs bool) string {
	return l.spec.PathSpec.URLizeFilename(l.relTargetPathForRel(rel, false, isAbs, true))
}

func (l *genericResource) relTargetPathForRel(rel string, addBaseTargetPath, isAbs, isURL bool) string {
	if addBaseTargetPath && len(l.baseTargetPathDirs) > 1 {
		panic("multiple baseTargetPathDirs")
	}
	var basePath string
	if addBaseTargetPath && len(l.baseTargetPathDirs) > 0 {
		basePath = l.baseTargetPathDirs[0]
	}

	return l.relTargetPathForRelAndBasePath(rel, basePath, isAbs, isURL)
}

func (l *genericResource) relTargetPathForRelAndBasePath(rel, basePath string, isAbs, isURL bool) string {
	rel = l.createBasePath(rel, isURL)

	if basePath != "" {
		rel = path.Join(basePath, rel)
	}

	if l.baseOffset != "" {
		rel = path.Join(l.baseOffset, rel)
	}

	if isURL {
		bp := l.spec.PathSpec.GetBasePath(!isAbs)
		if bp != "" {
			rel = path.Join(bp, rel)
		}
	}

	if len(rel) == 0 || rel[0] != '/' {
		rel = "/" + rel
	}

	return rel
}

func (l *genericResource) relTargetPaths() []string {
	return l.relTargetPathsForRel(l.TargetPath())
}

func (l *genericResource) relTargetPathsFor(target string) []string {
	return l.relTargetPathsForRel(target)
}

func (l *genericResource) relTargetPathsForRel(rel string) []string {
	if len(l.baseTargetPathDirs) == 0 {
		return []string{l.relTargetPathForRelAndBasePath(rel, "", false, false)}
	}

	targetPaths := make([]string, len(l.baseTargetPathDirs))
	for i, dir := range l.baseTargetPathDirs {
		targetPaths[i] = l.relTargetPathForRelAndBasePath(rel, dir, false, false)
	}
	return targetPaths
}

func (l *genericResource) updateParams(params map[string]interface{}) {
	if l.params == nil {
		l.params = params
		return
	}

	// Sets the params not already set
	for k, v := range params {
		if _, found := l.params[k]; !found {
			l.params[k] = v
		}
	}
}

type targetPather interface {
	TargetPath() string
}

type permalinker interface {
	targetPather
	permalinkFor(target string) string
	relPermalinkFor(target string) string
	relTargetPaths() []string
	relTargetPathsFor(target string) []string
}

type resourceContent struct {
	content     string
	contentInit sync.Once

	publishInit sync.Once
}

type resourceFileInfo struct {
	// Will be set if this resource is backed by something other than a file.
	openReadSeekerCloser resource.OpenReadSeekCloser

	// This may be set to tell us to look in another filesystem for this resource.
	// We, by default, use the sourceFs filesystem in the spec below.
	sourceFs afero.Fs

	// Absolute filename to the source, including any content folder path.
	// Note that this is absolute in relation to the filesystem it is stored in.
	// It can be a base path filesystem, and then this filename will not match
	// the path to the file on the real filesystem.
	sourceFilename string

	fi hugofs.FileMetaInfo

	// A hash of the source content. Is only calculated in caching situations.
	h *resourceHash
}

func (fi *resourceFileInfo) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	if fi.openReadSeekerCloser != nil {
		return fi.openReadSeekerCloser()
	}

	f, err := fi.getSourceFs().Open(fi.getSourceFilename())
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fi *resourceFileInfo) getFileInfo() hugofs.FileMetaInfo {
	return fi.fi
}

func (fi *resourceFileInfo) getSourceFilename() string {
	return fi.sourceFilename
}

func (fi *resourceFileInfo) setSourceFilename(s string) {
	// Make sure it's always loaded by sourceFilename.
	fi.openReadSeekerCloser = nil
	fi.sourceFilename = s
}

func (fi *resourceFileInfo) getSourceFs() afero.Fs {
	return fi.sourceFs
}

func (fi *resourceFileInfo) setSourceFs(fs afero.Fs) {
	fi.sourceFs = fs
}

func (fi *resourceFileInfo) hash() (string, error) {
	var err error
	fi.h.init.Do(func() {
		var hash string
		var f hugio.ReadSeekCloser
		f, err = fi.ReadSeekCloser()
		if err != nil {
			err = errors.Wrap(err, "failed to open source file")
			return
		}
		defer f.Close()

		hash, err = helpers.MD5FromFileFast(f)
		if err != nil {
			return
		}
		fi.h.value = hash
	})

	return fi.h.value, err
}

func (fi *resourceFileInfo) size() int {
	if fi.fi == nil {
		return 0
	}

	return int(fi.fi.Size())
}

type resourceHash struct {
	value string
	init  sync.Once
}

type resourcePathDescriptor struct {
	// The relative target directory and filename.
	relTargetDirFile dirFile

	// Callback used to construct a target path relative to its owner.
	targetPathBuilder func() page.TargetPaths

	// This will normally be the same as above, but this will only apply to publishing
	// of resources. It may be mulltiple values when in multihost mode.
	baseTargetPathDirs []string

	// baseOffset is set when the output format's path has a offset, e.g. for AMP.
	baseOffset string
}
