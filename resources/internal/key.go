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

package internal

import "github.com/gohugoio/hugo/helpers"

// ResourceTransformationKey are provided by the different transformation implementations.
// It identifies the transformation (name) and its configuration (elements).
// We combine this in a chain with the rest of the transformations
// with the target filename and a content hash of the origin to use as cache key.
type ResourceTransformationKey struct {
	Name     string
	elements []interface{}
}

// NewResourceTransformationKey creates a new ResourceTransformationKey from the transformation
// name and elements. We will create a 64 bit FNV hash from the elements, which when combined
// with the other key elements should be unique for all practical applications.
func NewResourceTransformationKey(name string, elements ...interface{}) ResourceTransformationKey {
	return ResourceTransformationKey{Name: name, elements: elements}
}

// Value returns the Key as a string.
// Do not change this without good reasons.
func (k ResourceTransformationKey) Value() string {
	if len(k.elements) == 0 {
		return k.Name
	}

	return k.Name + "_" + helpers.HashString(k.elements...)

}
