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

// Package page contains the core interfaces and types for the Page resource,
// a core component in Hugo.
package page

import (
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources/resource"
)

// TODO(bep) page there is language and stuff going on. There will be
// page sources that does not care about that, so a "DefaultLanguagePage" wrapper...

type Page interface {
	resource.Resource
	resource.ContentProvider
	resource.LanguageProvider
	resource.Dated

	Kind() string

	Param(key interface{}) (interface{}, error)

	Weight() int
	LinkTitle() string

	Resources() resource.Resources

	// Make it indexable as a related.Document
	SearchKeywords(cfg related.IndexConfig) ([]related.Keyword, error)
}

// TranslationProvider provides translated versions of a Page.
type TranslationProvider interface {
}
