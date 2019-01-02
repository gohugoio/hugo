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

package resource

import (
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/common/hugio"
)

// Cloner is an internal template and not meant for use in the templates. It
// may change without notice.
type Cloner interface {
	WithNewBase(base string) Resource
}

// Resource represents a linkable resource, i.e. a content page, image etc.
type Resource interface {
	resourceBase

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

	// Resource specific data set by Hugo.
	// One example would be.Data.Digest for fingerprinted resources.
	Data() interface{}

	// Params set in front matter for this resource.
	Params() map[string]interface{}
}

// resourceBase pulls out the minimal set of operations to define a Resource,
// to simplify testing etc.
type resourceBase interface {
	// MediaType is this resource's MIME type.
	MediaType() media.Type
}

// ResourcesLanguageMerger describes an interface for merging resources from a
// different language.
type ResourcesLanguageMerger interface {
	MergeByLanguage(other Resources) Resources
	// Needed for integration with the tpl package.
	MergeByLanguageInterface(other interface{}) (interface{}, error)
}

// Identifier identifies a resource.
type Identifier interface {
	Key() string
}

// ContentResource represents a Resource that provides a way to get to its content.
// Most Resource types in Hugo implements this interface, including Page.
// This should be used with care, as it will read the file content into memory, but it
// should be cached as effectively as possible by the implementation.
type ContentResource interface {
	resourceBase

	// Content returns this resource's content. It will be equivalent to reading the content
	// that RelPermalink points to in the published folder.
	// The return type will be contextual, and should be what you would expect:
	// * Page: template.HTML
	// * JSON: String
	// * Etc.
	Content() (interface{}, error)
}

// OpenReadSeekCloser allows setting some other way (than reading from a filesystem)
// to open or create a ReadSeekCloser.
type OpenReadSeekCloser func() (hugio.ReadSeekCloser, error)

// ReadSeekCloserResource is a Resource that supports loading its content.
type ReadSeekCloserResource interface {
	resourceBase
	ReadSeekCloser() (hugio.ReadSeekCloser, error)
}
