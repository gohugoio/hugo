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
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"
)

var (
	_ Resource     = (*genericResource)(nil)
	_ metaAssigner = (*genericResource)(nil)
	_ Source       = (*genericResource)(nil)
	_ Cloner       = (*genericResource)(nil)
)

const DefaultResourceType = "unknown"

// Source is an internal template and not meant for use in the templates. It
// may change without notice.
type Source interface {
	AbsSourceFilename() string
	Publish() error
}

// Cloner is an internal template and not meant for use in the templates. It
// may change without notice.
type Cloner interface {
	WithNewBase(base string) Resource
}

type metaAssigner interface {
	setTitle(title string)
	setName(name string)
	updateParams(params map[string]interface{})
}

// Resource represents a linkable resource, i.e. a content page, image etc.
type Resource interface {
	// Permalink represents the absolute link to this resource.
	Permalink() string

	// RelPermalink represents the host relative link to this resource.
	RelPermalink() string

	// ResourceType is the resource type. For most file types, this is the main
	// part of the MIME type, e.g. "image", "application", "text" etc.
	// For content pages, this value is "page".
	ResourceType() string

	// Name is the logical name of this resource. This can be set in the front matter
	// metadata for this resource. If not set, Hugo will assign a value.
	// This will in most cases be the base filename.
	// So, for the image "/some/path/sunset.jpg" this will be "sunset.jpg".
	// The value returned by this method will be used in the GetByPrefix and ByPrefix methods
	// on Resources.
	Name() string

	// Title returns the title if set in front matter. For content pages, this will be the expected value.
	Title() string

	// Params set in front matter for this resource.
	Params() map[string]interface{}
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

const prefixDeprecatedMsg = `We have added the more flexible Resources.GetMatch (find one) and Resources.Match (many) to replace the "prefix" methods. 

These matches by a given globbing pattern, e.g. "*.jpg".

Some examples:

* To find all resources by its prefix in the root dir of the bundle: .Match image*
* To find one resource by its prefix in the root dir of the bundle: .GetMatch image*
* To find all JPEG images anywhere in the bundle: .Match **.jpg`

// GetByPrefix gets the first resource matching the given filename prefix, e.g
// "logo" will match logo.png. It returns nil of none found.
// In potential ambiguous situations, combine it with ByType.
func (r Resources) GetByPrefix(prefix string) Resource {
	helpers.Deprecated("Resources", "GetByPrefix", prefixDeprecatedMsg, false)
	prefix = strings.ToLower(prefix)
	for _, resource := range r {
		if matchesPrefix(resource, prefix) {
			return resource
		}
	}
	return nil
}

// ByPrefix gets all resources matching the given base filename prefix, e.g
// "logo" will match logo.png.
func (r Resources) ByPrefix(prefix string) Resources {
	helpers.Deprecated("Resources", "ByPrefix", prefixDeprecatedMsg, false)
	var matches Resources
	prefix = strings.ToLower(prefix)
	for _, resource := range r {
		if matchesPrefix(resource, prefix) {
			matches = append(matches, resource)
		}
	}
	return matches
}

// GetMatch finds the first Resource matching the given pattern, or nil if none found.
// See Match for a more complete explanation about the rules used.
func (r Resources) GetMatch(pattern string) Resource {
	g, err := getGlob(pattern)
	if err != nil {
		return nil
	}

	for _, resource := range r {
		if g.Match(strings.ToLower(resource.Name())) {
			return resource
		}
	}

	return nil
}

// Match gets all resources matching the given base filename prefix, e.g
// "*.png" will match all png files. The "*" does not match path delimiters (/),
// so if you organize your resources in sub-folders, you need to be explicit about it, e.g.:
// "images/*.png". To match any PNG image anywhere in the bundle you can do "**.png", and
// to match all PNG images below the images folder, use "images/**.jpg".
// The matching is case insensitive.
// Match matches by using the value of Resource.Name, which, by default, is a filename with
// path relative to the bundle root with Unix style slashes (/) and no leading slash, e.g. "images/logo.png".
// See https://github.com/gobwas/glob for the full rules set.
func (r Resources) Match(pattern string) Resources {
	g, err := getGlob(pattern)
	if err != nil {
		return nil
	}

	var matches Resources
	for _, resource := range r {
		if g.Match(strings.ToLower(resource.Name())) {
			matches = append(matches, resource)
		}
	}
	return matches
}

func matchesPrefix(r Resource, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(r.Name()), prefix)
}

var (
	globCache = make(map[string]glob.Glob)
	globMu    sync.RWMutex
)

func getGlob(pattern string) (glob.Glob, error) {
	var g glob.Glob

	globMu.RLock()
	g, found := globCache[pattern]
	globMu.RUnlock()
	if !found {
		var err error
		g, err = glob.Compile(strings.ToLower(pattern), '/')
		if err != nil {
			return nil, err
		}

		globMu.Lock()
		globCache[pattern] = g
		globMu.Unlock()
	}

	return g, nil

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

	title  string
	name   string
	params map[string]interface{}

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

const counterPlaceHolder = ":counter"

// AssignMetadata assigns the given metadata to those resources that supports updates
// and matching by wildcard given in `src` using `filepath.Match` with lower cased values.
// This assignment is additive, but the most specific match needs to be first.
// The `name` and `title` metadata field support shell-matched collection it got a match in.
// See https://golang.org/pkg/path/#Match
func AssignMetadata(metadata []map[string]interface{}, resources ...Resource) error {

	counters := make(map[string]int)

	for _, r := range resources {
		if _, ok := r.(metaAssigner); !ok {
			continue
		}

		var (
			nameSet, titleSet                   bool
			nameCounter, titleCounter           = 0, 0
			nameCounterFound, titleCounterFound bool
			resourceSrcKey                      = strings.ToLower(r.Name())
		)

		ma := r.(metaAssigner)
		for _, meta := range metadata {
			src, found := meta["src"]
			if !found {
				return fmt.Errorf("missing 'src' in metadata for resource")
			}

			srcKey := strings.ToLower(cast.ToString(src))

			glob, err := getGlob(srcKey)
			if err != nil {
				return fmt.Errorf("failed to match resource with metadata: %s", err)
			}

			match := glob.Match(resourceSrcKey)

			if match {
				if !nameSet {
					name, found := meta["name"]
					if found {
						name := cast.ToString(name)
						if !nameCounterFound {
							nameCounterFound = strings.Contains(name, counterPlaceHolder)
						}
						if nameCounterFound && nameCounter == 0 {
							counterKey := "name_" + srcKey
							nameCounter = counters[counterKey] + 1
							counters[counterKey] = nameCounter
						}

						ma.setName(replaceResourcePlaceholders(name, nameCounter))
						nameSet = true
					}
				}

				if !titleSet {
					title, found := meta["title"]
					if found {
						title := cast.ToString(title)
						if !titleCounterFound {
							titleCounterFound = strings.Contains(title, counterPlaceHolder)
						}
						if titleCounterFound && titleCounter == 0 {
							counterKey := "title_" + srcKey
							titleCounter = counters[counterKey] + 1
							counters[counterKey] = titleCounter
						}
						ma.setTitle((replaceResourcePlaceholders(title, titleCounter)))
						titleSet = true
					}
				}

				params, found := meta["params"]
				if found {
					m := cast.ToStringMap(params)
					// Needed for case insensitive fetching of params values
					helpers.ToLowerMap(m)
					ma.updateParams(m)
				}
			}
		}
	}

	return nil
}

func replaceResourcePlaceholders(in string, counter int) string {
	return strings.Replace(in, counterPlaceHolder, strconv.Itoa(counter), -1)
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
		params:            make(map[string]interface{}),
		name:              baseFilename,
		title:             baseFilename,
	}
}
