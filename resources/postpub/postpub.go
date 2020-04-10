// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postpub

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"
)

type PostPublishedResource interface {
	resource.ResourceTypeProvider
	resource.ResourceLinksProvider
	resource.ResourceMetaProvider
	resource.ResourceParamsProvider
	resource.ResourceDataProvider
	resource.OriginProvider

	MediaType() map[string]interface{}
}

const (
	PostProcessPrefix = "__h_pp_l1"
	PostProcessSuffix = "__e"
)

func NewPostPublishResource(id int, r resource.Resource) PostPublishedResource {
	return &PostPublishResource{
		prefix:   PostProcessPrefix + "_" + strconv.Itoa(id) + "_",
		delegate: r,
	}
}

// postPublishResource holds a Resource to be transformed post publishing.
type PostPublishResource struct {
	prefix   string
	delegate resource.Resource
}

func (r *PostPublishResource) field(name string) string {
	return r.prefix + name + PostProcessSuffix
}

func (r *PostPublishResource) Permalink() string {
	return r.field("Permalink")
}

func (r *PostPublishResource) RelPermalink() string {
	return r.field("RelPermalink")
}

func (r *PostPublishResource) Origin() resource.Resource {
	return r.delegate
}

func (r *PostPublishResource) GetFieldString(pattern string) (string, bool) {
	if r == nil {
		panic("resource is nil")
	}
	prefixIdx := strings.Index(pattern, r.prefix)
	if prefixIdx == -1 {
		// Not a method on this resource.
		return "", false
	}

	fieldAccessor := pattern[prefixIdx+len(r.prefix) : strings.Index(pattern, PostProcessSuffix)]

	d := r.delegate
	switch {
	case fieldAccessor == "RelPermalink":
		return d.RelPermalink(), true
	case fieldAccessor == "Permalink":
		return d.Permalink(), true
	case fieldAccessor == "Name":
		return d.Name(), true
	case fieldAccessor == "Title":
		return d.Title(), true
	case fieldAccessor == "ResourceType":
		return d.ResourceType(), true
	case fieldAccessor == "Content":
		content, err := d.(resource.ContentProvider).Content()
		if err != nil {
			return "", true
		}
		return cast.ToString(content), true
	case strings.HasPrefix(fieldAccessor, "MediaType"):
		return r.fieldToString(d.MediaType(), fieldAccessor), true
	case fieldAccessor == "Data.Integrity":
		return cast.ToString((d.Data().(map[string]interface{})["Integrity"])), true
	default:
		panic(fmt.Sprintf("unknown field accessor %q", fieldAccessor))
	}

}

func (r *PostPublishResource) fieldToString(receiver interface{}, path string) string {
	fieldname := strings.Split(path, ".")[1]

	receiverv := reflect.ValueOf(receiver)
	switch receiverv.Kind() {
	case reflect.Map:
		v := receiverv.MapIndex(reflect.ValueOf(fieldname))
		return cast.ToString(v.Interface())
	default:
		v := receiverv.FieldByName(fieldname)
		if !v.IsValid() {
			method := receiverv.MethodByName(fieldname)
			if method.IsValid() {
				vals := method.Call(nil)
				if len(vals) > 0 {
					v = vals[0]
				}

			}
		}

		if v.IsValid() {
			return cast.ToString(v.Interface())
		}
		return ""
	}
}

func (r *PostPublishResource) Data() interface{} {
	m := map[string]interface{}{
		"Integrity": "",
	}
	insertFieldPlaceholders("Data", m, r.field)
	return m
}

func (r *PostPublishResource) MediaType() map[string]interface{} {
	m := structToMapWithPlaceholders("MediaType", media.Type{}, r.field)
	return m
}

func (r *PostPublishResource) ResourceType() string {
	return r.field("ResourceType")
}

func (r *PostPublishResource) Name() string {
	return r.field("Name")
}

func (r *PostPublishResource) Title() string {
	return r.field("Title")
}

func (r *PostPublishResource) Params() maps.Params {
	panic(r.fieldNotSupported("Params"))
}

func (r *PostPublishResource) Content() (interface{}, error) {
	return r.field("Content"), nil
}

func (r *PostPublishResource) fieldNotSupported(name string) string {
	return fmt.Sprintf("method .%s is currently not supported in post-publish transformations.", name)
}
