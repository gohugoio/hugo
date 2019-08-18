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
	"path"
	"strings"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/resource"

	"fmt"
	"io"
	"sync"

	"github.com/gohugoio/hugo/media"

	bp "github.com/gohugoio/hugo/bufferpool"
)

var (
	_ resource.ContentResource        = (*transformedResource)(nil)
	_ resource.ReadSeekCloserResource = (*transformedResource)(nil)
	_ collections.Slicer              = (*transformedResource)(nil)
	_ resource.Identifier             = (*transformedResource)(nil)
)

func (s *Spec) Transform(r resource.Resource, t ResourceTransformation) (resource.Resource, error) {
	if r == nil {
		return nil, errors.New("got nil Resource in transformation. Make sure you check with 'with' or 'if' when you get a resource, e.g. with resources.Get.")
	}

	if transformer, ok := r.(Transformer); ok {
		return transformer.Transform(t)
	}

	return nil, errors.Errorf("transform not supported for type %T", r)

	// TODO1
	// Something ala: return r.(ResourceTransformer).Apply(foo)
	// Which clones r
	tr := &transformedResource{
		Resource:                    r,
		transformation:              t,
		transformedResourceMetadata: transformedResourceMetadata{MetaData: make(map[string]interface{})},
		cache:                       s.ResourceCache}

	return tr, nil

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

func (ctx *ResourceTransformationCtx) addPathIdentifier(inPath, identifier string) string {
	dir, file := path.Split(inPath)
	base, ext := helpers.PathAndExt(file)
	return path.Join(dir, (base + identifier + ext))
}

// ReplaceOutPathExtension transforming InPath to OutPath replacing the file
// extension, e.g. ".scss"
func (ctx *ResourceTransformationCtx) ReplaceOutPathExtension(newExt string) {
	dir, file := path.Split(ctx.InPath)
	base, _ := helpers.PathAndExt(file)
	ctx.OutPath = path.Join(dir, (base + newExt))
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

// ResourceTransformation is the interface that a resource transformation step
// needs to implement.
type ResourceTransformation interface {
	Key() internal.ResourceTransformationKey
	Transform(ctx *ResourceTransformationCtx) error
}

// We will persist this information to disk.
type transformedResourceMetadata struct {
	Target     string                 `json:"Target"`
	MediaTypeV string                 `json:"MediaType"`
	MetaData   map[string]interface{} `json:"Data"`
}

type transformedResource struct {
	commonResource

	cache *ResourceCache

	// This is the filename inside resources/_gen/assets
	sourceFilename string

	linker permalinker

	// The transformation to apply.
	transformation ResourceTransformation

	// We apply the tranformations lazily.
	transformInit sync.Once
	transformErr  error

	// We delay publishing until either .RelPermalink or .Permalink
	// is invoked.
	publishInit sync.Once
	published   bool

	// The transformed values
	content     string
	contentInit sync.Once
	transformedResourceMetadata

	// The source
	resource.Resource
}

func (r *transformedResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	if err := r.initContent(); err != nil {
		return nil, err
	}
	return hugio.NewReadSeekerNoOpCloserFromString(r.content), nil
}

func (r *transformedResource) transferTransformedValues(another *transformedResource) {
	if another.content != "" {
		r.contentInit.Do(func() {
			r.content = another.content
		})
	}
	r.transformedResourceMetadata = another.transformedResourceMetadata
}

func (r *transformedResource) tryTransformedFileCache(key string) io.ReadCloser {
	fi, f, meta, found := r.cache.getFromFile(key)
	if !found {
		return nil
	}
	r.transformedResourceMetadata = meta
	r.sourceFilename = fi.Name

	return f
}

func (r *transformedResource) Content() (interface{}, error) {
	if err := r.initTransform(true, false); err != nil {
		return nil, err
	}
	if err := r.initContent(); err != nil {
		return "", err
	}
	return r.content, nil
}

func (r *transformedResource) Data() interface{} {
	if err := r.initTransform(false, false); err != nil {
		return noData
	}
	return r.MetaData
}

func (r *transformedResource) MediaType() media.Type {
	if err := r.initTransform(false, false); err != nil {
		return media.Type{}
	}
	m, _ := r.cache.rs.MediaTypes.GetByType(r.MediaTypeV)
	return m
}

func (r *transformedResource) Key() string {
	if err := r.initTransform(false, false); err != nil {
		return ""
	}
	return r.linker.relPermalinkFor(r.Target)
}

func (r *transformedResource) Permalink() string {
	if err := r.initTransform(false, true); err != nil {
		return ""
	}
	return r.linker.permalinkFor(r.Target)
}

func (r *transformedResource) RelPermalink() string {
	if err := r.initTransform(false, true); err != nil {
		return ""
	}
	return r.linker.relPermalinkFor(r.Target)
}

func (r *transformedResource) initContent() error {
	var err error
	r.contentInit.Do(func() {
		var b []byte
		_, b, err = r.cache.fileCache.GetBytes(r.sourceFilename)
		if err != nil {
			return
		}
		r.content = string(b)
	})
	return err
}

// TODO1
func (r *genericResource) openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error) {
	return helpers.OpenFilesForWriting(r.spec.BaseFs.PublishFs, r.relTargetPathsFor(relTargetPath)...)
}

func (r *transformedResource) openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error) {
	return helpers.OpenFilesForWriting(r.cache.rs.PublishFs, r.linker.relTargetPathsFor(relTargetPath)...)
}

type resourceTransformation struct {
	init            sync.Once
	transformations []ResourceTransformation
}

func (t *resourceTransformation) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RT: %d", len(t.transformations))
}

func (t *resourceTransformation) Add(transformation ResourceTransformation) {
	t.transformations = append(t.transformations, transformation)
}

func (t resourceTransformation) Clone() *resourceTransformation {
	transformations := make([]ResourceTransformation, len(t.transformations))
	copy(transformations, t.transformations)
	t.transformations = transformations

	return &t
}

func (t *resourceTransformation) Apply(r *genericResource) error {
	var transformErr error
	t.init.Do(func() {

		setContent := true // TODO1
		publish := true
		cache := r.spec.ResourceCache

		// Files with a suffix will be stored in cache (both on disk and in memory)
		// partitioned by their suffix.
		var key string
		for _, tr := range t.transformations {
			key = key + "_" + tr.Key().Value()
		}

		base := ResourceCacheKey(r.TargetPath())

		key = cache.cleanKey(base) + "_" + helpers.MD5String(key)

		// Acquire a write lock for the named transformation.
		cache.nlocker.Lock(key)

		defer cache.nlocker.Unlock(key)
		defer cache.set(key, r)

		b1 := bp.GetBuffer()
		b2 := bp.GetBuffer()
		defer bp.PutBuffer(b1)
		defer bp.PutBuffer(b2)

		tctx := &ResourceTransformationCtx{
			Data:                  make(map[string]interface{}),
			OpenResourcePublisher: nil, // TODO1
		}

		tctx.InMediaType = r.MediaType()
		tctx.OutMediaType = r.MediaType()

		var contentrc hugio.ReadSeekCloser
		contentrc, transformErr = contentReadSeekerCloser(r)
		if transformErr != nil {
			return
		}
		defer contentrc.Close()

		tctx.From = contentrc
		tctx.To = b1

		tctx.InPath = r.TargetPath()
		tctx.SourcePath = tctx.InPath

		counter := 0

		var transformedContentr io.Reader

		for i, tr := range t.transformations {
			if i != 0 {
				tctx.InMediaType = tctx.OutMediaType
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

			if transformErr = tr.Transform(tctx); transformErr != nil {

				if transformErr == herrors.ErrFeatureNotAvailable {
					// This transformation is not available in this
					// Hugo installation (scss not compiled in, PostCSS not available etc.)
					// If a prepared bundle for this transformation chain is available, use that.
					// TODO1
					f := r.tryTransformedFileCache(key)
					if f == nil {
						errMsg := transformErr.Error()
						if tr.Key().Name == "postcss" {
							errMsg = "PostCSS not found; install with \"npm install postcss-cli\". See https://gohugo.io/hugo-pipes/postcss/"
						}
						transformErr = fmt.Errorf("%s: failed to transform %q (%s): %s", strings.ToUpper(tr.Key().Name), tctx.InPath, tctx.InMediaType.Type(), errMsg)
						return
					}
					transformedContentr = f
					defer f.Close()

					// The reader above is all we need.
					break
				}

				// Abort.
				return
			}

			if tctx.OutPath != "" {
				tctx.InPath = tctx.OutPath
				tctx.OutPath = ""
			}
		}

		// TODO1 transformedContentr?
		// TODO1 optimize for no writes + test
		if transformedContentr == nil {
			r.setTransformedValues(tctx)
			//r.Target = tctx.InPath
			//r.MediaTypeV = tctx.OutMediaType.Type()
		}

		var publishwriters []io.WriteCloser

		if publish {
			var publicw io.WriteCloser
			publicw, transformErr = r.openPublishFileForWriting(r.TargetPath())
			if transformErr != nil {
				return
			}
			defer publicw.Close()

			publishwriters = append(publishwriters, publicw)
		}

		if transformedContentr == nil {
			// Also write it to the cache
			// TODO1
			/*fi, metaw, err := cache.writeMeta(key, r.transformedResourceMetadata)
			if err != nil {
				return err
			}*/
			//			r.sourceFilename = fi.Name

			//	publishwriters = append(publishwriters, metaw)

			if counter > 0 {
				transformedContentr = tctx.To.(*bytes.Buffer)
			} else {
				transformedContentr = contentrc
			}
		}

		// Also write it to memory
		var contentmemw *bytes.Buffer

		if setContent {
			contentmemw = bp.GetBuffer()
			defer bp.PutBuffer(contentmemw)
			publishwriters = append(publishwriters, hugio.ToWriteCloser(contentmemw))
		}

		publishw := hugio.NewMultiWriteCloser(publishwriters...)
		_, transformErr = io.Copy(publishw, transformedContentr)
		publishw.Close()

		// TODO1
		if setContent {
			r.contentInit.Do(func() {
				r.content = contentmemw.String()
			})
		}

	})
	return transformErr
}

func (r *transformedResource) transform(setContent, publish bool) (err error) {

	return nil
}

func (r *transformedResource) initTransform(setContent, publish bool) error {
	r.transformInit.Do(func() {
		r.published = publish
		if err := r.transform(setContent, publish); err != nil {
			r.transformErr = err
			r.cache.rs.Logger.ERROR.Println("error: failed to transform resource:", err)
		}

	})

	if !publish {
		return r.transformErr
	}

	r.publishInit.Do(func() {
		if r.published {
			return
		}

		r.published = true

		// Copy the file from cache to /public
		_, src, err := r.cache.fileCache.Get(r.sourceFilename)
		if src == nil {
			panic(fmt.Sprintf("[BUG] resource cache file not found: %q", r.sourceFilename))
		}

		if err == nil {
			defer src.Close()

			var dst io.WriteCloser
			dst, err = r.openPublishFileForWriting(r.Target)
			if err == nil {
				defer dst.Close()
				io.Copy(dst, src)
			}
		}

		if err != nil {
			r.transformErr = err
			r.cache.rs.Logger.ERROR.Println("error: failed to publish resource:", err)
			return
		}

	})

	return r.transformErr
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
