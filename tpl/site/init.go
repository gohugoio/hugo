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

// Package site provides template functions for accessing the Site object.
package site

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "site"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		s := page.WrapSite(d.Site)
		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return s, nil },
		}

		// We just add the Site as the namespace here. No method mappings.

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
