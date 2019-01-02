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
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/tpl"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"
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
)

var noData = make(map[string]interface{})

type permalinker interface {
	relPermalinkFor(target string) string
	permalinkFor(target string) string
	relTargetPathsFor(target string) []string
	relTargetPaths() []string
	targetPath() string
}

type Spec struct {
	*helpers.PathSpec

	MediaTypes    media.Types
	OutputFormats output.Formats

	Logger *loggers.Logger

	TextTemplates tpl.TemplateParseFinder

	// Holds default filter settings etc.
	imaging *Imaging

	imageCache    *imageCache
	ResourceCache *ResourceCache
	FileCaches    filecache.Caches
}

func NewSpec(
	s *helpers.PathSpec,
	fileCaches filecache.Caches,
	logger *loggers.Logger,
	outputFormats output.Formats,
	mimeTypes media.Types) (*Spec, error) {

	imaging, err := decodeImaging(s.Cfg.GetStringMap("imaging"))
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = loggers.NewErrorLogger()
	}

	rs := &Spec{PathSpec: s,
		Logger:        logger,
		imaging:       &imaging,
		MediaTypes:    mimeTypes,
		OutputFormats: outputFormats,
		FileCaches:    fileCaches,
		imageCache: newImageCache(
			fileCaches.ImageCache(),

			s,
		)}

	rs.ResourceCache = newResourceCache(rs)

	return rs, nil

}

type ResourceSourceDescriptor struct {
	// TargetPathBuilder is a callback to create target paths's relative to its owner.
	TargetPathBuilder func(base string) string

	// Need one of these to load the resource content.
	SourceFile         source.File
	OpenReadSeekCloser resource.OpenReadSeekCloser

	// If OpenReadSeekerCloser is not set, we use this to open the file.
	SourceFilename string

	// The relative target filename without any language code.
	RelTargetFilename string

	// Any base path prepeneded to the permalink.
	// Typically the language code if this resource should be published to its sub-folder.
	URLBase string

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

func (r *Spec) sourceFs() afero.Fs {
	return r.PathSpec.BaseFs.Content.Fs
}

func (r *Spec) New(fd ResourceSourceDescriptor) (resource.Resource, error) {
	return r.newResourceForFs(r.sourceFs(), fd)
}

func (r *Spec) NewForFs(sourceFs afero.Fs, fd ResourceSourceDescriptor) (resource.Resource, error) {
	return r.newResourceForFs(sourceFs, fd)
}

func (r *Spec) newResourceForFs(sourceFs afero.Fs, fd ResourceSourceDescriptor) (resource.Resource, error) {
	if fd.OpenReadSeekCloser == nil {
		if fd.SourceFile != nil && fd.SourceFilename != "" {
			return nil, errors.New("both SourceFile and AbsSourceFilename provided")
		} else if fd.SourceFile == nil && fd.SourceFilename == "" {
			return nil, errors.New("either SourceFile or AbsSourceFilename must be provided")
		}
	}

	if fd.RelTargetFilename == "" {
		fd.RelTargetFilename = fd.Filename()
	}

	if len(fd.TargetBasePaths) == 0 {
		// If not set, we publish the same resource to all hosts.
		fd.TargetBasePaths = r.MultihostTargetBasePaths
	}

	return r.newResource(sourceFs, fd)
}

func (r *Spec) newResource(sourceFs afero.Fs, fd ResourceSourceDescriptor) (resource.Resource, error) {
	var fi os.FileInfo
	var sourceFilename string

	if fd.OpenReadSeekCloser != nil {

	} else if fd.SourceFilename != "" {
		var err error
		fi, err = sourceFs.Stat(fd.SourceFilename)
		if err != nil {
			return nil, err
		}
		sourceFilename = fd.SourceFilename
	} else {
		fi = fd.SourceFile.FileInfo()
		sourceFilename = fd.SourceFile.Filename()
	}

	if fd.RelTargetFilename == "" {
		fd.RelTargetFilename = sourceFilename
	}

	ext := filepath.Ext(fd.RelTargetFilename)
	mimeType, found := r.MediaTypes.GetFirstBySuffix(strings.TrimPrefix(ext, "."))
	// TODO(bep) we need to handle these ambigous types better, but in this context
	// we most likely want the application/xml type.
	if mimeType.Suffix() == "xml" && mimeType.SubType == "rss" {
		mimeType, found = r.MediaTypes.GetByType("application/xml")
	}

	if !found {
		mimeStr := mime.TypeByExtension(ext)
		if mimeStr != "" {
			mimeType, _ = media.FromStringAndExt(mimeStr, ext)
		}
	}

	gr := r.newGenericResourceWithBase(
		sourceFs,
		fd.LazyPublish,
		fd.OpenReadSeekCloser,
		fd.URLBase,
		fd.TargetBasePaths,
		fd.TargetPathBuilder,
		fi,
		sourceFilename,
		fd.RelTargetFilename,
		mimeType)

	if mimeType.MainType == "image" {
		ext := strings.ToLower(helpers.Ext(sourceFilename))

		imgFormat, ok := imageFormats[ext]
		if !ok {
			// This allows SVG etc. to be used as resources. They will not have the methods of the Image, but
			// that would not (currently) have worked.
			return gr, nil
		}

		if err := gr.initHash(); err != nil {
			return nil, err
		}

		return &Image{
			format:          imgFormat,
			imaging:         r.imaging,
			genericResource: gr}, nil
	}
	return gr, nil

}

// TODO(bep) unify
func (r *Spec) IsInImageCache(key string) bool {
	// This is used for cache pruning. We currently only have images, but we could
	// imagine expanding on this.
	return r.imageCache.isInCache(key)
}

func (r *Spec) DeleteCacheByPrefix(prefix string) {
	r.imageCache.deleteByPrefix(prefix)
}

func (r *Spec) ClearCaches() {
	r.imageCache.clear()
	r.ResourceCache.clear()
}

func (r *Spec) CacheStats() string {
	r.imageCache.mu.RLock()
	defer r.imageCache.mu.RUnlock()

	s := fmt.Sprintf("Cache entries: %d", len(r.imageCache.store))

	count := 0
	for k := range r.imageCache.store {
		if count > 5 {
			break
		}
		s += "\n" + k
		count++
	}

	return s
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

type resourcePathDescriptor struct {
	// The relative target directory and filename.
	relTargetDirFile dirFile

	// Callback used to construct a target path relative to its owner.
	targetPathBuilder func(rel string) string

	// baseURLDir is the fixed sub-folder for a resource in permalinks. This will typically
	// be the language code if we publish to the language's sub-folder.
	baseURLDir string

	// This will normally be the same as above, but this will only apply to publishing
	// of resources. It may be mulltiple values when in multihost mode.
	baseTargetPathDirs []string

	// baseOffset is set when the output format's path has a offset, e.g. for AMP.
	baseOffset string
}

type resourceContent struct {
	content     string
	contentInit sync.Once
}

type resourceHash struct {
	hash     string
	hashInit sync.Once
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

// genericResource represents a generic linkable resource.
type genericResource struct {
	commonResource
	resourcePathDescriptor

	title  string
	name   string
	params map[string]interface{}

	// Absolute filename to the source, including any content folder path.
	// Note that this is absolute in relation to the filesystem it is stored in.
	// It can be a base path filesystem, and then this filename will not match
	// the path to the file on the real filesystem.
	sourceFilename string

	// Will be set if this resource is backed by something other than a file.
	openReadSeekerCloser resource.OpenReadSeekCloser

	// A hash of the source content. Is only calculated in caching situations.
	*resourceHash

	// This may be set to tell us to look in another filesystem for this resource.
	// We, by default, use the sourceFs filesystem in the spec below.
	overriddenSourceFs afero.Fs

	spec *Spec

	resourceType string
	mediaType    media.Type

	osFileInfo os.FileInfo

	// We create copies of this struct, so this needs to be a pointer.
	*resourceContent

	// May be set to signal lazy/delayed publishing.
	*publishOnce
}

type commonResource struct {
}

func (l *genericResource) Data() interface{} {
	return noData
}

func (l *genericResource) Content() (interface{}, error) {
	if err := l.initContent(); err != nil {
		return nil, err
	}

	return l.content, nil
}

func (l *genericResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	if l.openReadSeekerCloser != nil {
		return l.openReadSeekerCloser()
	}
	f, err := l.sourceFs().Open(l.sourceFilename)
	if err != nil {
		return nil, err
	}
	return f, nil

}

func (l *genericResource) MediaType() media.Type {
	return l.mediaType
}

// Implement the Cloner interface.
func (l genericResource) WithNewBase(base string) resource.Resource {
	l.baseOffset = base
	l.resourceContent = &resourceContent{}
	return &l
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
		}
		return groups, nil
	default:
		return nil, fmt.Errorf("invalid slice type %T", items)
	}
}

func (l *genericResource) initHash() error {
	var err error
	l.hashInit.Do(func() {
		var hash string
		var f hugio.ReadSeekCloser
		f, err = l.ReadSeekCloser()
		if err != nil {
			err = errors.Wrap(err, "failed to open source file")
			return
		}
		defer f.Close()

		hash, err = helpers.MD5FromFileFast(f)
		if err != nil {
			return
		}
		l.hash = hash

	})

	return err
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

func (l *genericResource) sourceFs() afero.Fs {
	if l.overriddenSourceFs != nil {
		return l.overriddenSourceFs
	}
	return l.spec.sourceFs()
}

func (l *genericResource) publishIfNeeded() {
	if l.publishOnce != nil {
		l.publishOnce.publish(l)
	}
}

func (l *genericResource) Permalink() string {
	l.publishIfNeeded()
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(l.relTargetDirFile.path(), true), l.spec.BaseURL.HostURL())
}

func (l *genericResource) RelPermalink() string {
	l.publishIfNeeded()
	return l.relPermalinkFor(l.relTargetDirFile.path())
}

func (l *genericResource) Key() string {
	return l.relTargetDirFile.path()
}

func (l *genericResource) relPermalinkFor(target string) string {
	return l.relPermalinkForRel(target, false)

}
func (l *genericResource) permalinkFor(target string) string {
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(target, true), l.spec.BaseURL.HostURL())

}
func (l *genericResource) relTargetPathsFor(target string) []string {
	return l.relTargetPathsForRel(target)
}

func (l *genericResource) relTargetPaths() []string {
	return l.relTargetPathsForRel(l.targetPath())
}

func (l *genericResource) Name() string {
	return l.name
}

func (l *genericResource) Title() string {
	return l.title
}

func (l *genericResource) Params() map[string]interface{} {
	return l.params
}

func (l *genericResource) setTitle(title string) {
	l.title = title
}

func (l *genericResource) setName(name string) {
	l.name = name
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

func (l *genericResource) relPermalinkForRel(rel string, isAbs bool) string {
	return l.spec.PathSpec.URLizeFilename(l.relTargetPathForRel(rel, false, isAbs, true))
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
	if l.targetPathBuilder != nil {
		rel = l.targetPathBuilder(rel)
	}

	if isURL && l.baseURLDir != "" {
		rel = path.Join(l.baseURLDir, rel)
	}

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

func (l *genericResource) ResourceType() string {
	return l.resourceType
}

func (l *genericResource) String() string {
	return fmt.Sprintf("Resource(%s: %s)", l.resourceType, l.name)
}

func (l *genericResource) Publish() error {
	fr, err := l.ReadSeekCloser()
	if err != nil {
		return err
	}
	defer fr.Close()
	fw, err := helpers.OpenFilesForWriting(l.spec.BaseFs.PublishFs, l.targetFilenames()...)
	if err != nil {
		return err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	return err
}

// Path is stored with Unix style slashes.
func (l *genericResource) targetPath() string {
	return l.relTargetDirFile.path()
}

func (l *genericResource) targetFilenames() []string {
	paths := l.relTargetPaths()
	for i, p := range paths {
		paths[i] = filepath.Clean(p)
	}
	return paths
}

// TODO(bep) clean up below
func (r *Spec) newGenericResource(sourceFs afero.Fs,
	targetPathBuilder func(base string) string,
	osFileInfo os.FileInfo,
	sourceFilename,
	baseFilename string,
	mediaType media.Type) *genericResource {
	return r.newGenericResourceWithBase(
		sourceFs,
		false,
		nil,
		"",
		nil,
		targetPathBuilder,
		osFileInfo,
		sourceFilename,
		baseFilename,
		mediaType,
	)

}

func (r *Spec) newGenericResourceWithBase(
	sourceFs afero.Fs,
	lazyPublish bool,
	openReadSeekerCloser resource.OpenReadSeekCloser,
	urlBaseDir string,
	targetPathBaseDirs []string,
	targetPathBuilder func(base string) string,
	osFileInfo os.FileInfo,
	sourceFilename,
	baseFilename string,
	mediaType media.Type) *genericResource {

	// This value is used both to construct URLs and file paths, but start
	// with a Unix-styled path.
	baseFilename = helpers.ToSlashTrimLeading(baseFilename)
	fpath, fname := path.Split(baseFilename)

	var resourceType string
	if mediaType.MainType == "image" {
		resourceType = mediaType.MainType
	} else {
		resourceType = mediaType.SubType
	}

	pathDescriptor := resourcePathDescriptor{
		baseURLDir:         urlBaseDir,
		baseTargetPathDirs: targetPathBaseDirs,
		targetPathBuilder:  targetPathBuilder,
		relTargetDirFile:   dirFile{dir: fpath, file: fname},
	}

	var po *publishOnce
	if lazyPublish {
		po = &publishOnce{logger: r.Logger}
	}

	return &genericResource{
		openReadSeekerCloser:   openReadSeekerCloser,
		publishOnce:            po,
		resourcePathDescriptor: pathDescriptor,
		overriddenSourceFs:     sourceFs,
		osFileInfo:             osFileInfo,
		sourceFilename:         sourceFilename,
		mediaType:              mediaType,
		resourceType:           resourceType,
		spec:                   r,
		params:                 make(map[string]interface{}),
		name:                   baseFilename,
		title:                  baseFilename,
		resourceContent:        &resourceContent{},
		resourceHash:           &resourceHash{},
	}
}
