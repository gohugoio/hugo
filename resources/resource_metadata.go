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
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"
)

var (
	_ mediaTypeAssigner             = (*genericResource)(nil)
	_ mediaTypeAssigner             = (*imageResource)(nil)
	_ resource.Staler               = (*genericResource)(nil)
	_ resource.NameOriginalProvider = (*genericResource)(nil)
)

// metaAssigner allows updating metadata in resources that supports it.
type metaAssigner interface {
	setTitle(title string)
	setName(name string)
	updateParams(params map[string]any)
}

// metaAssigner allows updating the media type in resources that supports it.
type mediaTypeAssigner interface {
	setMediaType(mediaType media.Type)
}

const counterPlaceHolder = ":counter"

var _ metaAssigner = (*metaResource)(nil)

// metaResource is a resource with metadata that can be updated.
type metaResource struct {
	changed bool
	title   string
	name    string
	params  maps.Params
}

func (r *metaResource) Name() string {
	return r.name
}

func (r *metaResource) Title() string {
	return r.title
}

func (r *metaResource) Params() maps.Params {
	return r.params
}

func (r *metaResource) setTitle(title string) {
	r.title = title
	r.changed = true
}

func (r *metaResource) setName(name string) {
	r.name = name
	r.changed = true
}

func (r *metaResource) updateParams(params map[string]any) {
	if r.params == nil {
		r.params = make(map[string]interface{})
	}
	for k, v := range params {
		r.params[k] = v
	}
	r.changed = true
}

func CloneWithMetadataIfNeeded(m []map[string]any, r resource.Resource) resource.Resource {
	wmp, ok := r.(resource.WithResourceMetaProvider)
	if !ok {
		return r
	}

	wrapped := &metaResource{
		name:   r.Name(),
		title:  r.Title(),
		params: r.Params(),
	}

	assignMetadata(m, wrapped)
	if !wrapped.changed {
		return r
	}

	return wmp.WithResourceMeta(wrapped)
}

// AssignMetadata assigns the given metadata to those resources that supports updates
// and matching by wildcard given in `src` using `filepath.Match` with lower cased values.
// This assignment is additive, but the most specific match needs to be first.
// The `name` and `title` metadata field support shell-matched collection it got a match in.
// See https://golang.org/pkg/path/#Match
func assignMetadata(metadata []map[string]any, ma *metaResource) error {
	counters := make(map[string]int)

	var (
		nameSet, titleSet                   bool
		nameCounter, titleCounter           = 0, 0
		nameCounterFound, titleCounterFound bool
		resourceSrcKey                      = strings.ToLower(ma.Name())
	)

	for _, meta := range metadata {
		src, found := meta["src"]
		if !found {
			return fmt.Errorf("missing 'src' in metadata for resource")
		}

		srcKey := strings.ToLower(cast.ToString(src))

		glob, err := glob.GetGlob(srcKey)
		if err != nil {
			return fmt.Errorf("failed to match resource with metadata: %w", err)
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
				m := maps.ToStringMap(params)
				// Needed for case insensitive fetching of params values
				maps.PrepareParams(m)
				ma.updateParams(m)
			}
		}
	}

	return nil
}

func replaceResourcePlaceholders(in string, counter int) string {
	return strings.Replace(in, counterPlaceHolder, strconv.Itoa(counter), -1)
}
