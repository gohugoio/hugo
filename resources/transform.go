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
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/images/exif"
	"github.com/spf13/afero"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/media"
)

var (
	_ resource.ContentResource        = (*resourceAdapter)(nil)
	_ resource.ReadSeekCloserResource = (*resourceAdapter)(nil)
	_ resource.Resource               = (*resourceAdapter)(nil)
	_ resource.Source                 = (*resourceAdapter)(nil)
	_ resource.Identifier             = (*resourceAdapter)(nil)
	_ resource.ResourceMetaProvider   = (*resourceAdapter)(nil)
)

// These are transformations that need special support in Hugo that may not
// be available when building the theme/site so we write the transformation
// result to disk and reuse if needed for these,
var transformationsToCacheOnDisk = map[string]bool{
	"postcss": true,
	"tocss":   true,
}

func newResourceAdapter(spec *Spec, lazyPublish bool, target transformableResource) *resourceAdapter {
	var po *publishOnce
	if lazyPublish {
		po = &publishOnce{}
	}
	return &resourceAdapter{
		resourceTransformations: &resourceTransformations{},
		resourceAdapterInner: &resourceAdapterInner{
			spec:        spec,
			publishOnce: po,
			target:      target,
		},
	}
}

// ResourceTransformation is the interface that a resource transformation step
// needs to implement.
type ResourceTransformation interface {
	Key() internal.ResourceTransformationKey
	Transform(ctx *ResourceTransformationCtx) error
}

type ResourceTransformationCtx struct {
	// The content to transform.
	From io.Reader

	// The target of content transformation.
	// The current implementation requires that r is written to w
	// even if no transformation is performed.
	To io.Writer

	// This is the relative path to the original source. Unix styled slashes.
	SourcePath string

	// This is the relative target path to the resource. Unix styled slashes.
	InPath string

	// The relative target path to the transformed resource. Unix styled slashes.
	OutPath string

	// The input media type
	InMediaType media.Type

	// The media type of the transformed resource.
	OutMediaType media.Type

	// Data data can be set on the transformed Resource. Not that this need
	// to be simple types, as it needs to be serialized to JSON and back.
	Data map[string]interface{}

	// This is used to publis additional artifacts, e.g. source maps.
	// We may improve this.
	OpenResourcePublisher func(relTargetPath string) (io.WriteCloser, error)
}

// AddOutPathIdentifier transforming InPath to OutPath adding an identifier,
// eg '.min' before any extension.
func (ctx *ResourceTransformationCtx) AddOutPathIdentifier(identifier string) {
	ctx.OutPath = ctx.addPathIdentifier(ctx.InPath, identifier)
}

// PublishSourceMap writes the content to the target folder of the main resource
// with the ".map" extension added.
func (ctx *ResourceTransformationCtx) PublishSourceMap(content string) error {
	target := ctx.OutPath + ".map"
	f, err := ctx.OpenResourcePublisher(target)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	return err
}

// ReplaceOutPathExtension transforming InPath to OutPath replacing the file
// extension, e.g. ".scss"
func (ctx *ResourceTransformationCtx) ReplaceOutPathExtension(newExt string) {
	dir, file := path.Split(ctx.InPath)
	base, _ := helpers.PathAndExt(file)
	ctx.OutPath = path.Join(dir, (base + newExt))
}

func (ctx *ResourceTransformationCtx) addPathIdentifier(inPath, identifier string) string {
	dir, file := path.Split(inPath)
	base, ext := helpers.PathAndExt(file)
	return path.Join(dir, (base + identifier + ext))
}

type publishOnce struct {
	publisherInit sync.Once
	publisherErr  error
}

type resourceAdapter struct {
	commonResource
	*resourceTransformations
	*resourceAdapterInner
}

func (r *resourceAdapter) Content() (interface{}, error) {
	r.init(false, true)
	if r.transformationsErr != nil {
		return nil, r.transformationsErr
	}
	return r.target.Content()
}

func (r *resourceAdapter) Data() interface{} {
	r.init(false, false)
	return r.target.Data()
}

func (r *resourceAdapter) Fill(spec string) (resource.Image, error) {
	return r.getImageOps().Fill(spec)
}

func (r *resourceAdapter) Fit(spec string) (resource.Image, error) {
	return r.getImageOps().Fit(spec)
}

func (r *resourceAdapter) Filter(filters ...interface{}) (resource.Image, error) {
	return r.getImageOps().Filter(filters...)
}

func (r *resourceAdapter) Height() int {
	return r.getImageOps().Height()
}

func (r *resourceAdapter) Exif() (*exif.Exif, error) {
	return r.getImageOps().Exif()
}

func (r *resourceAdapter) Key() string {
	r.init(false, false)
	return r.target.(resource.Identifier).Key()
}

func (r *resourceAdapter) MediaType() media.Type {
	r.init(false, false)
	return r.target.MediaType()
}

func (r *resourceAdapter) Name() string {
	r.init(false, false)
	return r.target.Name()
}

func (r *resourceAdapter) Params() maps.Params {
	r.init(false, false)
	return r.target.Params()
}

func (r *resourceAdapter) Permalink() string {
	r.init(true, false)
	return r.target.Permalink()
}

func (r *resourceAdapter) Publish() error {
	r.init(false, false)

	return r.target.Publish()
}

func (r *resourceAdapter) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	r.init(false, false)
	return r.target.ReadSeekCloser()
}

func (r *resourceAdapter) RelPermalink() string {
	r.init(true, false)
	return r.target.RelPermalink()
}

func (r *resourceAdapter) Resize(spec string) (resource.Image, error) {
	return r.getImageOps().Resize(spec)
}

func (r *resourceAdapter) ResourceType() string {
	r.init(false, false)
	return r.target.ResourceType()
}

func (r *resourceAdapter) String() string {
	return r.Name()
}

func (r *resourceAdapter) Title() string {
	r.init(false, false)
	return r.target.Title()
}

func (r resourceAdapter) Transform(t ...ResourceTransformation) (ResourceTransformer, error) {
	r.resourceTransformations = &resourceTransformations{
		transformations: append(r.transformations, t...),
	}

	r.resourceAdapterInner = &resourceAdapterInner{
		spec:        r.spec,
		publishOnce: &publishOnce{},
		target:      r.target,
	}

	return &r, nil
}

func (r *resourceAdapter) Width() int {
	return r.getImageOps().Width()
}

func (r *resourceAdapter) getImageOps() resource.ImageOps {
	img, ok := r.target.(resource.ImageOps)
	if !ok {
		panic(fmt.Sprintf("%T is not an image", r.target))
	}
	r.init(false, false)
	return img
}

func (r *resourceAdapter) getMetaAssigner() metaAssigner {
	return r.target
}

func (r *resourceAdapter) getSpec() *Spec {
	return r.spec
}

func (r *resourceAdapter) publish() {
	if r.publishOnce == nil {
		return
	}

	r.publisherInit.Do(func() {
		r.publisherErr = r.target.Publish()

		if r.publisherErr != nil {
			r.spec.Logger.ERROR.Printf("Failed to publish Resource: %s", r.publisherErr)
		}
	})

}

func (r *resourceAdapter) TransformationKey() string {
	// Files with a suffix will be stored in cache (both on disk and in memory)
	// partitioned by their suffix.
	var key string
	for _, tr := range r.transformations {
		key = key + "_" + tr.Key().Value()
	}

	base := ResourceCacheKey(r.target.Key())
	return r.spec.ResourceCache.cleanKey(base) + "_" + helpers.MD5String(key)
}

func (r *resourceAdapter) transform(publish, setContent bool) error {
	cache := r.spec.ResourceCache

	key := r.TransformationKey()

	cached, found := cache.get(key)

	if found {
		r.resourceAdapterInner = cached.(*resourceAdapterInner)
		return nil
	}

	// Acquire a write lock for the named transformation.
	cache.nlocker.Lock(key)
	// Check the cache again.
	cached, found = cache.get(key)
	if found {
		r.resourceAdapterInner = cached.(*resourceAdapterInner)
		cache.nlocker.Unlock(key)
		return nil
	}

	defer cache.nlocker.Unlock(key)
	defer cache.set(key, r.resourceAdapterInner)

	b1 := bp.GetBuffer()
	b2 := bp.GetBuffer()
	defer bp.PutBuffer(b1)
	defer bp.PutBuffer(b2)

	tctx := &ResourceTransformationCtx{
		Data:                  make(map[string]interface{}),
		OpenResourcePublisher: r.target.openPublishFileForWriting,
	}

	tctx.InMediaType = r.target.MediaType()
	tctx.OutMediaType = r.target.MediaType()

	startCtx := *tctx
	updates := &transformationUpdate{startCtx: startCtx}

	var contentrc hugio.ReadSeekCloser

	contentrc, err := contentReadSeekerCloser(r.target)
	if err != nil {
		return err
	}

	defer contentrc.Close()

	tctx.From = contentrc
	tctx.To = b1

	tctx.InPath = r.target.TargetPath()
	tctx.SourcePath = tctx.InPath

	counter := 0
	writeToFileCache := false

	var transformedContentr io.Reader

	for i, tr := range r.transformations {
		if i != 0 {
			tctx.InMediaType = tctx.OutMediaType
		}

		mayBeCachedOnDisk := transformationsToCacheOnDisk[tr.Key().Name]
		if !writeToFileCache {
			writeToFileCache = mayBeCachedOnDisk
		}

		if i > 0 {
			hasWrites := tctx.To.(*bytes.Buffer).Len() > 0
			if hasWrites {
				counter++
				// Switch the buffers
				if counter%2 == 0 {
					tctx.From = b2
					b1.Reset()
					tctx.To = b1
				} else {
					tctx.From = b1
					b2.Reset()
					tctx.To = b2
				}
			}
		}

		newErr := func(err error) error {

			msg := fmt.Sprintf("%s: failed to transform %q (%s)", strings.ToUpper(tr.Key().Name), tctx.InPath, tctx.InMediaType.Type())

			if err == herrors.ErrFeatureNotAvailable {
				var errMsg string
				if tr.Key().Name == "postcss" {
					// This transformation is not available in this
					// Most likely because PostCSS is not installed.
					errMsg = ". Check your PostCSS installation; install with \"npm install postcss-cli\". See https://gohugo.io/hugo-pipes/postcss/"
				} else if tr.Key().Name == "tocss" {
					errMsg = ". Check your Hugo installation; you need the extended version to build SCSS/SASS."
				}

				return errors.New(msg + errMsg)
			}

			return errors.Wrap(err, msg)

		}

		var tryFileCache bool

		if mayBeCachedOnDisk && r.spec.BuildConfig.UseResourceCache(nil) {
			tryFileCache = true
		} else {
			err = tr.Transform(tctx)
			if err != nil && err != herrors.ErrFeatureNotAvailable {
				return newErr(err)
			}

			if mayBeCachedOnDisk {
				tryFileCache = r.spec.BuildConfig.UseResourceCache(err)
			}
			if err != nil && !tryFileCache {
				return newErr(err)
			}
		}

		if tryFileCache {
			f := r.target.tryTransformedFileCache(key, updates)
			if f == nil {
				return newErr(errors.Errorf("resource %q not found in file cache", key))
			}
			transformedContentr = f
			updates.sourceFs = cache.fileCache.Fs
			defer f.Close()

			// The reader above is all we need.
			break
		}

		if tctx.OutPath != "" {
			tctx.InPath = tctx.OutPath
			tctx.OutPath = ""
		}
	}

	if transformedContentr == nil {
		updates.updateFromCtx(tctx)
	}

	var publishwriters []io.WriteCloser

	if publish {
		publicw, err := r.target.openPublishFileForWriting(updates.targetPath)
		if err != nil {
			return err
		}
		publishwriters = append(publishwriters, publicw)
	}

	if transformedContentr == nil {
		if writeToFileCache {
			// Also write it to the cache
			fi, metaw, err := cache.writeMeta(key, updates.toTransformedResourceMetadata())
			if err != nil {
				return err
			}
			updates.sourceFilename = &fi.Name
			updates.sourceFs = cache.fileCache.Fs
			publishwriters = append(publishwriters, metaw)
		}

		// Any transofrmations reading from From must also write to To.
		// This means that if the target buffer is empty, we can just reuse
		// the original reader.
		if b, ok := tctx.To.(*bytes.Buffer); ok && b.Len() > 0 {
			transformedContentr = tctx.To.(*bytes.Buffer)
		} else {
			transformedContentr = contentrc
		}
	}

	// Also write it to memory
	var contentmemw *bytes.Buffer

	setContent = setContent || !writeToFileCache

	if setContent {
		contentmemw = bp.GetBuffer()
		defer bp.PutBuffer(contentmemw)
		publishwriters = append(publishwriters, hugio.ToWriteCloser(contentmemw))
	}

	publishw := hugio.NewMultiWriteCloser(publishwriters...)
	_, err = io.Copy(publishw, transformedContentr)
	if err != nil {
		return err
	}
	publishw.Close()

	if setContent {
		s := contentmemw.String()
		updates.content = &s
	}

	newTarget, err := r.target.cloneWithUpdates(updates)
	if err != nil {
		return err
	}
	r.target = newTarget

	return nil
}

func (r *resourceAdapter) init(publish, setContent bool) {
	r.initTransform(publish, setContent)
}

func (r *resourceAdapter) initTransform(publish, setContent bool) {
	r.transformationsInit.Do(func() {
		if len(r.transformations) == 0 {
			// Nothing to do.
			return
		}

		if publish {
			// The transformation will write the content directly to
			// the destination.
			r.publishOnce = nil
		}

		r.transformationsErr = r.transform(publish, setContent)
		if r.transformationsErr != nil {
			if r.spec.ErrorSender != nil {
				r.spec.ErrorSender.SendError(r.transformationsErr)
			} else {
				r.spec.Logger.ERROR.Printf("Transformation failed: %s", r.transformationsErr)
			}
		}
	})

	if publish && r.publishOnce != nil {
		r.publish()
	}
}

type resourceAdapterInner struct {
	target transformableResource

	spec *Spec

	// Handles publishing (to /public) if needed.
	*publishOnce
}

type resourceTransformations struct {
	transformationsInit sync.Once
	transformationsErr  error
	transformations     []ResourceTransformation
}

type transformableResource interface {
	baseResourceInternal

	resource.ContentProvider
	resource.Resource
	resource.Identifier
}

type transformationUpdate struct {
	content        *string
	sourceFilename *string
	sourceFs       afero.Fs
	targetPath     string
	mediaType      media.Type
	data           map[string]interface{}

	startCtx ResourceTransformationCtx
}

func (u *transformationUpdate) isContenChanged() bool {
	return u.content != nil || u.sourceFilename != nil
}

func (u *transformationUpdate) toTransformedResourceMetadata() transformedResourceMetadata {
	return transformedResourceMetadata{
		MediaTypeV: u.mediaType.Type(),
		Target:     u.targetPath,
		MetaData:   u.data,
	}
}

func (u *transformationUpdate) updateFromCtx(ctx *ResourceTransformationCtx) {
	u.targetPath = ctx.OutPath
	u.mediaType = ctx.OutMediaType
	u.data = ctx.Data
	u.targetPath = ctx.InPath
}

// We will persist this information to disk.
type transformedResourceMetadata struct {
	Target     string                 `json:"Target"`
	MediaTypeV string                 `json:"MediaType"`
	MetaData   map[string]interface{} `json:"Data"`
}

// contentReadSeekerCloser returns a ReadSeekerCloser if possible for a given Resource.
func contentReadSeekerCloser(r resource.Resource) (hugio.ReadSeekCloser, error) {
	switch rr := r.(type) {
	case resource.ReadSeekCloserResource:
		rc, err := rr.ReadSeekCloser()
		if err != nil {
			return nil, err
		}
		return rc, nil
	default:
		return nil, fmt.Errorf("cannot transform content of Resource of type %T", r)

	}
}
