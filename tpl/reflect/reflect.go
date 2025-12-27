// Copyright 2018 The Hugo Authors. All rights reserved.
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

package reflect

import (
	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

// New returns a new instance of the reflect-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "reflect" namespace.
type Namespace struct{}

// IsMap reports whether v is a map.
func (ns *Namespace) IsMap(v any) bool {
	return hreflect.IsMap(v)
}

// IsSlice reports whether v is a slice.
func (ns *Namespace) IsSlice(v any) bool {
	return hreflect.IsSlice(v)
}

// IsPage reports whether v is a Hugo Page.
func (ns *Namespace) IsPage(v any) bool {
	_, ok := v.(page.Page)
	return ok
}

// IsResource reports whether v is a Hugo Resource.
func (ns *Namespace) IsResource(v any) bool {
	_, ok := v.(resource.Resource)
	return ok
}

// IsSite reports whether v is a Hugo Site.
func (ns *Namespace) IsSite(v any) bool {
	_, ok := v.(page.Site)
	return ok
}

// IsImageResource reports whether v is an Hugo Image Resource.
// If this returns true, you may process it and get information about its width, height, etc.
func (ns *Namespace) IsImageResource(v any) bool {
	return resources.IsImage(v)
}
