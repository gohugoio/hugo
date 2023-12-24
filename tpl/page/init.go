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

// Package page provides template functions for accessing the current Page object,
// the entry level context for the current template.
package page

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "page"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ns := &internal.TemplateFuncsNamespace{
			Name: name,
			Context: func(ctx context.Context, args ...interface{}) (interface{}, error) {
				v := tpl.Context.Page.Get(ctx)
				if v == nil {
					// The multilingual sitemap does not have a page as its context.
					return nil, nil
				}

				return v.(page.Page), nil
			},
		}

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
