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
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/mitchellh/hashstructure"

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
	return &transformedResource{
		Resource:                    r,
		transformation:              t,
		transformedResourceMetadata: transformedResourceMetadata{MetaData: make(map[string]interface{})},
		cache:                       s.ResourceCache}, nil
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

// ResourceTransformationKey are provided by the different transformation implementations.
// It identifies the transformation (name) and its configuration (elements).
// We combine this in a chain with the rest of the transformations
// with the target filename and a content hash of the origin to use as cache key.
type ResourceTransformationKey struct {
	name     string
	elements []interface{}
}

// NewResourceTransformationKey creates a new ResourceTransformationKey from the transformation
// name and elements. We will create a 64 bit FNV hash from the elements, which when combined
// with the other key elements should be unique for all practical applications.
func NewResourceTransformationKey(name string, elements ...interface{}) ResourceTransformationKey {
	return ResourceTransformationKey{name: name, elements: elements}
}

// Do not change this without good reasons.
func (k ResourceTransformationKey) key() string {
	if len(k.elements) == 0 {
		return k.name
	}

	sb := bp.GetBuffer()
	defer bp.PutBuffer(sb)

	sb.WriteString(k.name)
	for _, element := range k.elements {
		hash, err := hashstructure.Hash(element, nil)
		if err != nil {
			panic(err)
		}
		sb.WriteString("_")
		sb.WriteString(strconv.FormatUint(hash, 10))
	}

	return sb.String()
}

// ResourceTransformation is the interface that a resource transformation step
// needs to implement.
type ResourceTransformation interface {
	Key() ResourceTransformationKey
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

func (r *transformedResource) openPublishFileForWriting(relTargetPath string) (io.WriteCloser, error) {
	return helpers.OpenFilesForWriting(r.cache.rs.PublishFs, r.linker.relTargetPathsFor(relTargetPath)...)
}

func (r *transformedResource) transform(setContent, publish bool) (err error) {

	// This can be the last resource in a chain.
	// Rewind and create a processing chain.
	var chain []resource.Resource
	current := r
	for {
		rr := current.Resource
		chain = append(chain[:0], append([]resource.Resource{rr}, chain[0:]...)...)
		if tr, ok := rr.(*transformedResource); ok {
			current = tr
		} else {
			break
		}
	}

	// Append the current transformer at the end
	chain = append(chain, r)

	first := chain[0]

	// Files with a suffix will be stored in cache (both on disk and in memory)
	// partitioned by their suffix. There will be other files below /other.
	// This partition is also how we determine what to delete on server reloads.
	var key, base string
	for _, element := range chain {
		switch v := element.(type) {
		case *transformedResource:
			key = key + "_" + v.transformation.Key().key()
		case permalinker:
			r.linker = v
			p := v.targetPath()
			if p == "" {
				panic("target path needed for key creation")
			}
			partition := ResourceKeyPartition(p)
			base = partition + "/" + p
		default:
			return fmt.Errorf("transformation not supported for type %T", element)
		}
	}

	key = r.cache.cleanKey(base + "_" + helpers.MD5String(key))

	cached, found := r.cache.get(key)
	if found {
		r.transferTransformedValues(cached.(*transformedResource))
		return
	}

	// Acquire a write lock for the named transformation.
	r.cache.nlocker.Lock(key)
	// Check the cache again.
	cached, found = r.cache.get(key)
	if found {
		r.transferTransformedValues(cached.(*transformedResource))
		r.cache.nlocker.Unlock(key)
		return
	}

	defer r.cache.nlocker.Unlock(key)
	defer r.cache.set(key, r)

	b1 := bp.GetBuffer()
	b2 := bp.GetBuffer()
	defer bp.PutBuffer(b1)
	defer bp.PutBuffer(b2)

	tctx := &ResourceTransformationCtx{
		Data:                  r.transformedResourceMetadata.MetaData,
		OpenResourcePublisher: r.openPublishFileForWriting,
	}

	tctx.InMediaType = first.MediaType()
	tctx.OutMediaType = first.MediaType()

	contentrc, err := contentReadSeekerCloser(first)
	if err != nil {
		return err
	}
	defer contentrc.Close()

	tctx.From = contentrc
	tctx.To = b1

	if r.linker != nil {
		tctx.InPath = r.linker.targetPath()
		tctx.SourcePath = tctx.InPath
	}

	counter := 0

	var transformedContentr io.Reader

	for _, element := range chain {
		tr, ok := element.(*transformedResource)
		if !ok {
			continue
		}
		counter++
		if counter != 1 {
			tctx.InMediaType = tctx.OutMediaType
		}
		if counter%2 == 0 {
			tctx.From = b1
			b2.Reset()
			tctx.To = b2
		} else {
			if counter != 1 {
				// The first reader is the file.
				tctx.From = b2
			}
			b1.Reset()
			tctx.To = b1
		}

		if err := tr.transformation.Transform(tctx); err != nil {
			if err == herrors.ErrFeatureNotAvailable {
				// This transformation is not available in this
				// Hugo installation (scss not compiled in, PostCSS not available etc.)
				// If a prepared bundle for this transformation chain is available, use that.
				f := r.tryTransformedFileCache(key)
				if f == nil {
					errMsg := err.Error()
					if tr.transformation.Key().name == "postcss" {
						errMsg = "PostCSS not found; install with \"npm install postcss-cli\". See https://gohugo.io/hugo-pipes/postcss/"
					}
					return fmt.Errorf("%s: failed to transform %q (%s): %s", strings.ToUpper(tr.transformation.Key().name), tctx.InPath, tctx.InMediaType.Type(), errMsg)
				}
				transformedContentr = f
				defer f.Close()

				// The reader above is all we need.
				break
			}

			// Abort.
			return err
		}

		if tctx.OutPath != "" {
			tctx.InPath = tctx.OutPath
			tctx.OutPath = ""
		}
	}

	if transformedContentr == nil {
		r.Target = tctx.InPath
		r.MediaTypeV = tctx.OutMediaType.Type()
	}

	var publishwriters []io.WriteCloser

	if publish {
		publicw, err := r.openPublishFileForWriting(r.Target)
		if err != nil {
			r.transformErr = err
			return err
		}
		defer publicw.Close()

		publishwriters = append(publishwriters, publicw)
	}

	if transformedContentr == nil {
		// Also write it to the cache
		fi, metaw, err := r.cache.writeMeta(key, r.transformedResourceMetadata)
		if err != nil {
			return err
		}
		r.sourceFilename = fi.Name

		publishwriters = append(publishwriters, metaw)

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
	_, r.transformErr = io.Copy(publishw, transformedContentr)
	publishw.Close()

	if setContent {
		r.contentInit.Do(func() {
			r.content = contentmemw.String()
		})
	}

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
