// Copyright 2024 The Hugo Authors. All rights reserved.
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

package transform

import (
	"os"

	"github.com/gohugoio/hugo/deps"

	"github.com/go-openapi/loads"
)

// OpenAPIDocFromJSON loads an OpenAPI document from a file path.
func (h *helpers) OpenAPIDocFromJSON(filename string) (*loads.Document, error) {
	r, err := h.s.rs.Get(filename) // <-- MODIFIED LINE 1
	if err != nil {               // <-- MODIFIED LINE 2
		return nil, err           // <-- MODIFIED LINE 3
	}                             // <-- MODIFIED LINE 4
	data, err := r.Content()      // <-- MODIFIED LINE 5

	if err != nil {
		return nil, err
	}

	doc, err := loads.Analyzed(data, "")
	if err != nil {
		return nil, err
	}

	return doc, nil
}
