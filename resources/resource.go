// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/helpers"
)

var (
	_ resource.ContentResource           = (*genericResource)(nil)
	_ resource.ReadSeekCloserResource    = (*genericResource)(nil)
	_ resource.Resource                  = (*genericResource)(nil)
	_ resource.Source                    = (*genericResource)(nil)
	_ resource.Cloner                    = (*genericResource)(nil)
	_ resource.ResourcesLanguageMerger   = (*resource.Resources)(nil)
	_ resource.Identifier                = (*genericResource)(nil)
	_ identity.IdentityGroupProvider     = (*genericResource)(nil)
	_ identity.DependencyManagerProvider = (*genericResource)(nil)
	_ identity.Identity                  = (*genericResource)(nil)
	_ fileInfo                           = (*genericResource)(nil)
)

type ResourceSourceDescriptor struct {
	// The source content.
	OpenReadSeekCloser hugio.OpenReadSeekCloser

	// The canonical source path.
	Path *paths.Path

	// The name of the resource.
	Name string

	// The name of the resource as it was read from the source.
	NameOriginal string

	// Any base paths prepended to the target path. This will also typically be the
	// language code, but setting it here means that it should not have any effect on
	// the permalink.
	// This may be several values. In multihost mode we may publish the same resources to
	// multiple targets.
	TargetBasePaths []string

	TargetPath           string
	BasePathRelPermalink string
	BasePathTargetPath   string

	// The Data to associate with this resource.
	Data map[string]any

	// Delay publishing until either Permalink or RelPermalink is called. Maybe never.
	LazyPublish bool

	// Set when its known up front, else it's resolved from the target filename.
	MediaType media.Type

	// Used to track dependencies (e.g. imports). May be nil if that's of no concern.
	DependencyManager identity.Manager

	// A shared identity for this resource and all its clones.
	// If this is not set, it's set to Anonymous.
	GroupIdentity identity.Identity
}

func (fd *ResourceSourceDescriptor) init(r *Spec) error {
	if len(fd.TargetBasePaths) == 0 {
		// If not set, we publish the same resource to all hosts.
		fd.TargetBasePaths = r.MultihostTargetBasePaths
	}

	if fd.OpenReadSeekCloser == nil {
		panic(errors.New("OpenReadSeekCloser is nil"))
	}

	if fd.TargetPath == "" {
		panic(errors.New("RelPath is empty"))
	}

	if fd.Path == nil {
		fd.Path = paths.Parse("", fd.TargetPath)
	}

	if fd.TargetPath == "" {
		fd.TargetPath = fd.Path.Path()
	} else {
		fd.TargetPath = paths.ToSlashPreserveLeading(fd.TargetPath)
	}

	fd.BasePathRelPermalink = paths.ToSlashPreserveLeading(fd.BasePathRelPermalink)
	if fd.BasePathRelPermalink == "/" {
		fd.BasePathRelPermalink = ""
	}
	fd.BasePathTargetPath = paths.ToSlashPreserveLeading(fd.BasePathTargetPath)
	if fd.BasePathTargetPath == "/" {
		fd.BasePathTargetPath = ""
	}

	fd.TargetPath = paths.ToSlashPreserveLeading(fd.TargetPath)
	for i, base := range fd.TargetBasePaths {
		dir := paths.ToSlashPreserveLeading(base)
		if dir == "/" {
			dir = ""
		}
		fd.TargetBasePaths[i] = dir
	}

	if fd.Name == "" {
		fd.Name = fd.TargetPath
	}

	if fd.NameOriginal == "" {
		fd.NameOriginal = fd.Name
	}

	mediaType := fd.MediaType
	if mediaType.IsZero() {
		ext := fd.Path.Ext()
		var (
			found      bool
			suffixInfo media.SuffixInfo
		)
		mediaType, suffixInfo, found = r.MediaTypes().GetFirstBySuffix(ext)
		// TODO(bep) we need to handle these ambiguous types better, but in this context
		// we most likely want the application/xml type.
		if suffixInfo.Suffix == "xml" && mediaType.SubType == "rss" {
			mediaType, found = r.MediaTypes().GetByType("application/xml")
		}

		if !found {
			// A fallback. Note that mime.TypeByExtension is slow by Hugo standards,
			// so we should configure media types to avoid this lookup for most
			// situations.
			mimeStr := mime.TypeByExtension("." + ext)
			if mimeStr != "" {
				mediaType, _ = media.FromStringAndExt(mimeStr, ext)
			}
		}
	}

	fd.MediaType = mediaType

	if fd.DependencyManager == nil {
		if r.Cfg.Watching() {
			fd.DependencyManager = identity.NewManager("resource")
		} else {
			fd.DependencyManager = identity.NopManager
		}
	}

	if fd.GroupIdentity == nil {
		fd.GroupIdentity = identity.Anonymous
	}

	return nil
}

type ResourceTransformer interface {
	resource.Resource
	Transformer
}

type Transformer interface {
	Transform(...ResourceTransformation) (ResourceTransformer, error)
	TransformWithContext(context.Context, ...ResourceTransformation) (ResourceTransformer, error)
}

func NewFeatureNotAvailableTransformer(key string, elements ...any) ResourceTransformation {
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

// resourceCopier is for internal use.
type resourceCopier interface {
	cloneTo(targetPath string) resource.Resource
}

// Copy copies r to the targetPath given.
func Copy(r resource.Resource, targetPath string) resource.Resource {
	if r.Err() != nil {
		panic(fmt.Sprintf("Resource has an .Err: %s", r.Err()))
	}
	return r.(resourceCopier).cloneTo(targetPath)
}

type baseResourceResource interface {
	resource.Cloner
	resourceCopier
	resource.ContentProvider
	resource.Resource
	resource.Identifier
}

type baseResourceInternal interface {
	resource.Source
	resource.NameOriginalProvider

	fileInfo
	mediaTypeAssigner
	targetPather

	ReadSeekCloser() (hugio.ReadSeekCloser, error)

	identity.IdentityGroupProvider
	identity.DependencyManagerProvider

	// For internal use.
	cloneWithUpdates(*transformationUpdate) (baseResource, error)
	tryTransformedFileCache(key string, u *transformationUpdate) io.ReadCloser

	getResourcePaths() internal.ResourcePaths

	specProvider
	openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error)
}

type specProvider interface {
	getSpec() *Spec
}

type baseResource interface {
	baseResourceResource
	baseResourceInternal
	resource.Staler
}

type commonResource struct{}

// Slice is for internal use.
// for the template functions. See collections.Slice.
func (commonResource) Slice(in any) (any, error) {
	switch items := in.(type) {
	case resource.Resources:
		return items, nil
	case []any:
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

type fileInfo interface {
	setOpenSource(hugio.OpenReadSeekCloser)
	setSourceFilenameIsHash(bool)
	setTargetPath(internal.ResourcePaths)
	size() int64
	hashProvider
}

type hashProvider interface {
	hash() string
}

type StaleValue[V any] struct {
	// The value.
	Value V

	// IsStaleFunc reports whether the value is stale.
	IsStaleFunc func() bool
}

func (s *StaleValue[V]) IsStale() bool {
	return s.IsStaleFunc()
}

type AtomicStaler struct {
	stale uint32
}

func (s *AtomicStaler) MarkStale() {
	atomic.StoreUint32(&s.stale, 1)
}

func (s *AtomicStaler) IsStale() bool {
	return atomic.LoadUint32(&(s.stale)) > 0
}

// For internal use.
type GenericResourceTestInfo struct {
	Paths internal.ResourcePaths
}

// For internal use.
func GetTestInfoForResource(r resource.Resource) GenericResourceTestInfo {
	var gr *genericResource
	switch v := r.(type) {
	case *genericResource:
		gr = v
	case *resourceAdapter:
		gr = v.target.(*genericResource)
	default:
		panic(fmt.Sprintf("unknown resource type: %T", r))
	}
	return GenericResourceTestInfo{
		Paths: gr.paths,
	}
}

// genericResource represents a generic linkable resource.
type genericResource struct {
	publishInit *sync.Once

	sd    ResourceSourceDescriptor
	paths internal.ResourcePaths

	sourceFilenameIsHash bool

	h *resourceHash // A hash of the source content. Is only calculated in caching situations.

	resource.Staler

	title  string
	name   string
	params map[string]any

	spec *Spec
}

func (l *genericResource) IdentifierBase() string {
	return l.sd.Path.IdentifierBase()
}

func (l *genericResource) GetIdentityGroup() identity.Identity {
	return l.sd.GroupIdentity
}

func (l *genericResource) GetDependencyManager() identity.Manager {
	return l.sd.DependencyManager
}

func (l *genericResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	return l.sd.OpenReadSeekCloser()
}

func (l *genericResource) Clone() resource.Resource {
	return l.clone()
}

func (l *genericResource) size() int64 {
	l.hash()
	return l.h.size
}

func (l *genericResource) hash() string {
	if err := l.h.init(l); err != nil {
		panic(err)
	}
	return l.h.value
}

func (l *genericResource) setOpenSource(openSource hugio.OpenReadSeekCloser) {
	l.sd.OpenReadSeekCloser = openSource
}

func (l *genericResource) setSourceFilenameIsHash(b bool) {
	l.sourceFilenameIsHash = b
}

func (l *genericResource) setTargetPath(d internal.ResourcePaths) {
	l.paths = d
}

func (l *genericResource) cloneTo(targetPath string) resource.Resource {
	c := l.clone()
	c.paths = c.paths.FromTargetPath(targetPath)
	return c
}

func (l *genericResource) Content(context.Context) (any, error) {
	r, err := l.ReadSeekCloser()
	if err != nil {
		return "", err
	}
	defer r.Close()

	return hugio.ReadString(r)
}

func (r *genericResource) Err() resource.ResourceError {
	return nil
}

func (l *genericResource) Data() any {
	return l.sd.Data
}

func (l *genericResource) Key() string {
	basePath := l.spec.Cfg.BaseURL().BasePathNoTrailingSlash
	var key string
	if basePath == "" {
		key = l.RelPermalink()
	} else {
		key = strings.TrimPrefix(l.RelPermalink(), basePath)
	}

	if l.spec.Cfg.IsMultihost() {
		key = l.spec.Lang() + key
	}

	return key
}

func (l *genericResource) MediaType() media.Type {
	return l.sd.MediaType
}

func (l *genericResource) setMediaType(mediaType media.Type) {
	l.sd.MediaType = mediaType
}

func (l *genericResource) Name() string {
	return l.name
}

func (l *genericResource) NameOriginal() string {
	return l.sd.NameOriginal
}

func (l *genericResource) Params() maps.Params {
	return l.params
}

func (l *genericResource) Publish() error {
	var err error
	l.publishInit.Do(func() {
		targetFilenames := l.getResourcePaths().TargetFilenames()

		if l.sourceFilenameIsHash {
			// This is a processed image. We want to avoid copying it if it hasn't changed.
			var changedFilenames []string
			for _, targetFilename := range targetFilenames {
				if _, err := l.getSpec().BaseFs.PublishFs.Stat(targetFilename); err == nil {
					continue
				}
				changedFilenames = append(changedFilenames, targetFilename)
			}
			if len(changedFilenames) == 0 {
				return
			}
			targetFilenames = changedFilenames
		}
		var fr hugio.ReadSeekCloser
		fr, err = l.ReadSeekCloser()
		if err != nil {
			return
		}
		defer fr.Close()

		var fw io.WriteCloser
		fw, err = helpers.OpenFilesForWriting(l.spec.BaseFs.PublishFs, targetFilenames...)
		if err != nil {
			return
		}
		defer fw.Close()

		_, err = io.Copy(fw, fr)
	})

	return err
}

func (l *genericResource) RelPermalink() string {
	return l.spec.PathSpec.GetBasePath(false) + paths.PathEscape(l.paths.TargetLink())
}

func (l *genericResource) Permalink() string {
	return l.spec.Cfg.BaseURL().WithPathNoTrailingSlash + paths.PathEscape(l.paths.TargetPath())
}

func (l *genericResource) ResourceType() string {
	return l.MediaType().MainType
}

func (l *genericResource) String() string {
	return fmt.Sprintf("Resource(%s: %s)", l.ResourceType(), l.name)
}

// Path is stored with Unix style slashes.
func (l *genericResource) TargetPath() string {
	return l.paths.TargetPath()
}

func (l *genericResource) Title() string {
	return l.title
}

func (l *genericResource) getSpec() *Spec {
	return l.spec
}

func (l *genericResource) getResourcePaths() internal.ResourcePaths {
	return l.paths
}

func (r *genericResource) tryTransformedFileCache(key string, u *transformationUpdate) io.ReadCloser {
	fi, f, meta, found := r.spec.ResourceCache.getFromFile(key)
	if !found {
		return nil
	}
	u.sourceFilename = &fi.Name
	mt, _ := r.spec.MediaTypes().GetByType(meta.MediaTypeV)
	u.mediaType = mt
	u.data = meta.MetaData
	u.targetPath = meta.Target
	return f
}

func (r *genericResource) mergeData(in map[string]any) {
	if len(in) == 0 {
		return
	}
	if r.sd.Data == nil {
		r.sd.Data = make(map[string]any)
	}
	for k, v := range in {
		if _, found := r.sd.Data[k]; !found {
			r.sd.Data[k] = v
		}
	}
}

func (rc *genericResource) cloneWithUpdates(u *transformationUpdate) (baseResource, error) {
	r := rc.clone()

	if u.content != nil {
		r.sd.OpenReadSeekCloser = func() (hugio.ReadSeekCloser, error) {
			return hugio.NewReadSeekerNoOpCloserFromString(*u.content), nil
		}
	}

	r.sd.MediaType = u.mediaType

	if u.sourceFilename != nil {
		if u.sourceFs == nil {
			return nil, errors.New("sourceFs is nil")
		}
		r.setOpenSource(func() (hugio.ReadSeekCloser, error) {
			return u.sourceFs.Open(*u.sourceFilename)
		})
	} else if u.sourceFs != nil {
		return nil, errors.New("sourceFs is set without sourceFilename")
	}

	if u.targetPath == "" {
		return nil, errors.New("missing targetPath")
	}

	r.setTargetPath(r.paths.FromTargetPath(u.targetPath))
	r.mergeData(u.data)

	return r, nil
}

func (l genericResource) clone() *genericResource {
	l.publishInit = &sync.Once{}
	return &l
}

func (r *genericResource) openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error) {
	filenames := r.paths.FromTargetPath(relTargetPath).TargetFilenames()
	return helpers.OpenFilesForWriting(r.spec.BaseFs.PublishFs, filenames...)
}

type targetPather interface {
	TargetPath() string
}

type resourceHash struct {
	value    string
	size     int64
	initOnce sync.Once
}

func (r *resourceHash) init(l hugio.ReadSeekCloserProvider) error {
	var initErr error
	r.initOnce.Do(func() {
		var hash string
		var size int64
		f, err := l.ReadSeekCloser()
		if err != nil {
			initErr = fmt.Errorf("failed to open source: %w", err)
			return
		}
		defer f.Close()
		hash, size, err = helpers.MD5FromReaderFast(f)
		if err != nil {
			initErr = fmt.Errorf("failed to calculate hash: %w", err)
			return
		}
		r.value = hash
		r.size = size
	})

	return initErr
}
