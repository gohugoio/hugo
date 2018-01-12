// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package resource

import (
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/helpers"
)

var (
	_ Resource = (*genericResource)(nil)
	_ Source   = (*genericResource)(nil)
	_ Cloner   = (*genericResource)(nil)
)

const DefaultResourceType = "unknown"

type Source interface {
	AbsSourceFilename() string
	Publish() error
}

type Cloner interface {
	WithNewBase(base string) Resource
}

// Resource represents a linkable resource, i.e. a content page, image etc.
type Resource interface {
	Permalink() string
	RelPermalink() string
	ResourceType() string
}

// Resources represents a slice of resources, which can be a mix of different types.
// I.e. both pages and images etc.
type Resources []Resource

func (r Resources) ByType(tp string) Resources {
	var filtered Resources

	for _, resource := range r {
		if resource.ResourceType() == tp {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// GetBySuffix gets the first resource matching the given filename prefix, e.g
// "logo" will match logo.png. It returns nil of none found.
// In potential ambiguous situations, combine it with ByType.
func (r Resources) GetByPrefix(prefix string) Resource {
	prefix = strings.ToLower(prefix)
	for _, resource := range r {
		var name string
		f, ok := resource.(source.File)
		if ok {
			name = f.BaseFileName()
		} else {
			_, name = filepath.Split(resource.RelPermalink())
		}
		name = strings.ToLower(name)

		if strings.HasPrefix(name, prefix) {
			return resource
		}
	}
	return nil
}

type Spec struct {
	*helpers.PathSpec
	mimeTypes media.Types

	// Holds default filter settings etc.
	imaging *Imaging

	imageCache *imageCache

	AbsGenImagePath string
}

func NewSpec(s *helpers.PathSpec, mimeTypes media.Types) (*Spec, error) {

	imaging, err := decodeImaging(s.Cfg.GetStringMap("imaging"))
	if err != nil {
		return nil, err
	}
	s.GetLayoutDirPath()

	genImagePath := s.AbsPathify(filepath.Join(s.Cfg.GetString("resourceDir"), "_gen", "images"))

	return &Spec{AbsGenImagePath: genImagePath, PathSpec: s, imaging: &imaging, mimeTypes: mimeTypes, imageCache: newImageCache(
		s,
		// We're going to write a cache pruning routine later, so make it extremely
		// unlikely that the user shoots him or herself in the foot
		// and this is set to a value that represents data he/she
		// cares about. This should be set in stone once released.
		genImagePath,
		s.AbsPathify(s.Cfg.GetString("publishDir")))}, nil
}

func (r *Spec) NewResourceFromFile(
	targetPathBuilder func(base string) string,
	absPublishDir string,
	file source.File, relTargetFilename string) (Resource, error) {

	return r.newResource(targetPathBuilder, absPublishDir, file.Filename(), file.FileInfo(), relTargetFilename)
}

func (r *Spec) NewResourceFromFilename(
	targetPathBuilder func(base string) string,
	absPublishDir,
	absSourceFilename, relTargetFilename string) (Resource, error) {

	fi, err := r.Fs.Source.Stat(absSourceFilename)
	if err != nil {
		return nil, err
	}
	return r.newResource(targetPathBuilder, absPublishDir, absSourceFilename, fi, relTargetFilename)
}

func (r *Spec) newResource(
	targetPathBuilder func(base string) string,
	absPublishDir,
	absSourceFilename string, fi os.FileInfo, relTargetFilename string) (Resource, error) {

	var mimeType string
	ext := filepath.Ext(relTargetFilename)
	m, found := r.mimeTypes.GetBySuffix(strings.TrimPrefix(ext, "."))
	if found {
		mimeType = m.SubType
	} else {
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = DefaultResourceType
		} else {
			mimeType = mimeType[:strings.Index(mimeType, "/")]
		}
	}

	gr := r.newGenericResource(targetPathBuilder, fi, absPublishDir, absSourceFilename, filepath.ToSlash(relTargetFilename), mimeType)

	if mimeType == "image" {
		f, err := r.Fs.Source.Open(absSourceFilename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		hash, err := helpers.MD5FromFileFast(f)
		if err != nil {
			return nil, err
		}

		return &Image{
			hash:            hash,
			imaging:         r.imaging,
			genericResource: gr}, nil
	}
	return gr, nil
}

func (r *Spec) IsInCache(key string) bool {
	// This is used for cache pruning. We currently only have images, but we could
	// imagine expanding on this.
	return r.imageCache.isInCache(key)
}

func (r *Spec) DeleteCacheByPrefix(prefix string) {
	r.imageCache.deleteByPrefix(prefix)
}

func (r *Spec) CacheStats() string {
	r.imageCache.mu.RLock()
	defer r.imageCache.mu.RUnlock()

	s := fmt.Sprintf("Cache entries: %d", len(r.imageCache.store))

	count := 0
	for k, _ := range r.imageCache.store {
		if count > 5 {
			break
		}
		s += "\n" + k
		count++
	}

	return s
}

// genericResource represents a generic linkable resource.
type genericResource struct {
	// The relative path to this resource.
	relTargetPath string

	// Base is set when the output format's path has a offset, e.g. for AMP.
	base string

	// Absolute filename to the source, including any content folder path.
	absSourceFilename string
	absPublishDir     string
	resourceType      string
	osFileInfo        os.FileInfo

	spec              *Spec
	targetPathBuilder func(rel string) string
}

func (l *genericResource) Permalink() string {
	return l.spec.PermalinkForBaseURL(l.relPermalinkForRel(l.relTargetPath, false), l.spec.BaseURL.String())
}

func (l *genericResource) RelPermalink() string {
	return l.relPermalinkForRel(l.relTargetPath, true)
}

// Implement the Cloner interface.
func (l genericResource) WithNewBase(base string) Resource {
	l.base = base
	return &l
}

func (l *genericResource) relPermalinkForRel(rel string, addBasePath bool) string {
	return l.spec.PathSpec.URLizeFilename(l.relTargetPathForRel(rel, addBasePath))
}

func (l *genericResource) relTargetPathForRel(rel string, addBasePath bool) string {
	if l.targetPathBuilder != nil {
		rel = l.targetPathBuilder(rel)
	}

	if l.base != "" {
		rel = path.Join(l.base, rel)
	}

	if addBasePath && l.spec.PathSpec.BasePath != "" {
		rel = path.Join(l.spec.PathSpec.BasePath, rel)
	}

	if rel[0] != '/' {
		rel = "/" + rel
	}

	return rel
}

func (l *genericResource) ResourceType() string {
	return l.resourceType
}

func (l *genericResource) AbsSourceFilename() string {
	return l.absSourceFilename
}

func (l *genericResource) Publish() error {
	f, err := l.spec.Fs.Source.Open(l.AbsSourceFilename())
	if err != nil {
		return err
	}
	defer f.Close()

	target := filepath.Join(l.absPublishDir, l.target())

	return helpers.WriteToDisk(target, f, l.spec.Fs.Destination)
}

func (l *genericResource) target() string {
	target := l.relTargetPathForRel(l.relTargetPath, false)
	if l.spec.PathSpec.Languages.IsMultihost() {
		target = path.Join(l.spec.PathSpec.Language.Lang, target)
	}
	return target
}

func (r *Spec) newGenericResource(
	targetPathBuilder func(base string) string,
	osFileInfo os.FileInfo,
	absPublishDir,
	absSourceFilename,
	baseFilename,
	resourceType string) *genericResource {

	return &genericResource{
		targetPathBuilder: targetPathBuilder,
		osFileInfo:        osFileInfo,
		absPublishDir:     absPublishDir,
		absSourceFilename: absSourceFilename,
		relTargetPath:     baseFilename,
		resourceType:      resourceType,
		spec:              r,
	}
}
