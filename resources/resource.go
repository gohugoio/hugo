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

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/loggers"
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
	_ collections.Slicer               = (*genericResource)(nil)
	_ resource.Identifier              = (*genericResource)(nil)
	_ fileInfo                         = (*genericResource)(nil)
	_ Transformer                      = (*genericResource)(nil)
)

var noData = make(map[string]interface{})

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

type Transformer interface {
	Transform(ResourceTransformation) (resource.Resource, error)
}

type baseResource interface {
	fileInfo
	metaAssigner

	resource.Cloner
	resource.ContentProvider

	resource.Resource
	resource.Source
	ReadSeekCloser() (hugio.ReadSeekCloser, error)

	// Internal

	getResourcePaths() *resourcePathDescriptor
	getSpec() *Spec
	getTargetFilenames() []string
	openDestinationsForWriting() (io.WriteCloser, error)

	addTransformation(t ResourceTransformation)

	relTargetPathForRel(rel string, addBaseTargetPath, isAbs, isURL bool) string
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
	setOpenReadSeekerCloser(r resource.OpenReadSeekCloser)
	getSourceFs() afero.Fs
	setSourceFs(afero.Fs)
	hash() (string, error)
	size() int
}

// genericResource represents a generic linkable resource.
type genericResource struct {
	commonResource

	transformation *resourceTransformation

	*resourcePathDescriptor
	*resourceFileInfo
	*resourceContent
	*publishOnce

	spec *Spec

	title  string
	name   string
	params map[string]interface{}

	resourceType string
	mediaType    media.Type
}

func (l *genericResource) Clone() resource.Resource {
	return l.clone()
}

// Implement the Cloner interface.
func (l genericResource) CloneWithNewBase(base string) resource.Resource {
	lc := l.clone()
	lc.baseOffset = base
	return lc
}

func (l *genericResource) Content() (interface{}, error) {
	if err := l.initContent(); err != nil {
		return nil, err
	}

	return l.content, nil
}

func (l *genericResource) Data() interface{} {
	return noData
}

func (l *genericResource) Key() string {
	return l.relTargetDirFile.path()
}

func (l *genericResource) MediaType() media.Type {
	return l.mediaType
}

func (l *genericResource) Name() string {
	return l.name
}

func (l *genericResource) Params() map[string]interface{} {
	return l.params
}

func (l *genericResource) Permalink() string {
	l.publishIfNeeded()
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(l.relTargetDirFile.path(), true), l.spec.BaseURL.HostURL())
}

func (l *genericResource) Publish() error {
	fr, err := l.ReadSeekCloser()
	if err != nil {
		return err
	}
	defer fr.Close()

	fw, err := helpers.OpenFilesForWriting(l.spec.BaseFs.PublishFs, l.getTargetFilenames()...)
	if err != nil {
		return err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	return err
}

func (l *genericResource) RelPermalink() string {
	l.publishIfNeeded()
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

func (r *genericResource) Transform(t ResourceTransformation) (resource.Resource, error) {
	rt := r.clone()
	rt.addTransformation(t)
	return rt, nil
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

func (l *genericResource) setMediaType(m media.Type) {
	l.mediaType = m
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

func (r *genericResource) initTransform() {
	if r.transformation == nil {
		// Nothing to do
		return
	}

	if err := r.transformation.Apply(r); err != nil {
		r.logger.ERROR.Println(err)
	}
}

func (r *genericResource) tryTransformedFileCache(key string) io.ReadCloser {
	fi, f, _, found := r.spec.ResourceCache.getFromFile(key)
	if !found {
		return nil
	}
	// TODO1
	//r.transformedResourceMetadata = meta
	r.sourceFilename = fi.Name

	return f
}

func (r *genericResource) setTransformedValues(ctx *ResourceTransformationCtx) {
	fpath, fname := path.Split(ctx.InPath)
	r.mediaType = ctx.OutMediaType
	r.resourcePathDescriptor.relTargetDirFile = dirFile{dir: fpath, file: fname}
}

func (r *genericResource) addTransformation(t ResourceTransformation) {
	if r.transformation == nil {
		r.transformation = &resourceTransformation{}
	}
	r.transformation.Add(t)
}

func (l genericResource) clone() *genericResource {
	gi := *l.resourceFileInfo
	rp := *l.resourcePathDescriptor
	l.resourceFileInfo = &gi
	l.resourcePathDescriptor = &rp
	l.resourceContent = &resourceContent{}
	// TODO1
	if l.publishOnce != nil {
		l.publishOnce = &publishOnce{logger: l.publishOnce.logger}
	}
	if l.transformation != nil {
		l.transformation = l.transformation.Clone()
	}
	return &l
}

// returns an opened file or nil if nothing to write.
func (l *genericResource) openDestinationsForWriting() (io.WriteCloser, error) {
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
		return nil, nil
	}

	return helpers.OpenFilesForWriting(l.getSpec().BaseFs.PublishFs, changedFilenames...)

}

func (l *genericResource) permalinkFor(target string) string {
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(target, true), l.spec.BaseURL.HostURL())

}

func (l *genericResource) publishIfNeeded() {
	l.initTransform()
	if l.publishOnce != nil {
		l.publishOnce.publish(l)
	}
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

	var targetPaths = make([]string, len(l.baseTargetPathDirs))
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

type permalinker interface {
	TargetPath() string
	permalinkFor(target string) string
	relPermalinkFor(target string) string
	relTargetPaths() []string
	relTargetPathsFor(target string) []string
}

type publishOnce struct {
	publisherInit sync.Once
	publisherErr  error
	logger        *loggers.Logger
}

func (l *publishOnce) publish(s resource.Source) error {
	l.publisherInit.Do(func() {
		l.publisherErr = s.Publish()
		if l.publisherErr != nil {
			l.logger.ERROR.Printf("failed to publish Resource: %s", l.publisherErr)
		}
	})
	return l.publisherErr
}

type resourceContent struct {
	content     string
	contentInit sync.Once
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

	fi os.FileInfo

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

func (fi *resourceFileInfo) getSourceFilename() string {
	return fi.sourceFilename
}

func (fi *resourceFileInfo) setSourceFilename(s string) {
	// Make sure it's always loaded by sourceFilename.
	fi.openReadSeekerCloser = nil
	fi.sourceFilename = s
}

func (fi *resourceFileInfo) setOpenReadSeekerCloser(r resource.OpenReadSeekCloser) {
	fi.openReadSeekerCloser = r
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
