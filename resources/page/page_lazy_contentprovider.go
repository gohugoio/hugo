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

package page

import (
	"html/template"

	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/markup/converter"
)

// OutputFormatContentProvider represents the method set that is "outputFormat aware" and that we
// provide lazy initialization for in case they get invoked outside of their normal rendering context, e.g. via .Translations.
// Note that this set is currently not complete, but should cover the most common use cases.
// For the others, the implementation will be from the page.NoopPage.
type OutputFormatContentProvider interface {
	OutputFormatPageContentProvider

	// for internal use.
	ContentRenderer
}

// OutputFormatPageContentProvider holds the exported methods from Page that are "outputFormat aware".
type OutputFormatPageContentProvider interface {
	ContentProvider
	TableOfContentsProvider
	PageRenderProvider
}

// LazyContentProvider initializes itself when read. Each method of the
// ContentProvider interface initializes a content provider and shares it
// with other methods.
//
// Used in cases where we cannot guarantee whether the content provider
// will be needed. Must create via NewLazyContentProvider.
type LazyContentProvider struct {
	init *lazy.Init
	cp   OutputFormatContentProvider
}

// NewLazyContentProvider returns a LazyContentProvider initialized with
// function f. The resulting LazyContentProvider calls f in order to
// retrieve a ContentProvider
func NewLazyContentProvider(f func() (OutputFormatContentProvider, error)) *LazyContentProvider {
	lcp := LazyContentProvider{
		init: lazy.New(),
		cp:   NopCPageContentRenderer,
	}
	lcp.init.Add(func() (any, error) {
		cp, err := f()
		if err != nil {
			return nil, err
		}
		lcp.cp = cp
		return nil, nil
	})
	return &lcp
}

func (lcp *LazyContentProvider) Reset() {
	lcp.init.Reset()
}

func (lcp *LazyContentProvider) Content() (any, error) {
	lcp.init.Do()
	return lcp.cp.Content()
}

func (lcp *LazyContentProvider) Plain() string {
	lcp.init.Do()
	return lcp.cp.Plain()
}

func (lcp *LazyContentProvider) PlainWords() []string {
	lcp.init.Do()
	return lcp.cp.PlainWords()
}

func (lcp *LazyContentProvider) Summary() template.HTML {
	lcp.init.Do()
	return lcp.cp.Summary()
}

func (lcp *LazyContentProvider) Truncated() bool {
	lcp.init.Do()
	return lcp.cp.Truncated()
}

func (lcp *LazyContentProvider) FuzzyWordCount() int {
	lcp.init.Do()
	return lcp.cp.FuzzyWordCount()
}

func (lcp *LazyContentProvider) WordCount() int {
	lcp.init.Do()
	return lcp.cp.WordCount()
}

func (lcp *LazyContentProvider) ReadingTime() int {
	lcp.init.Do()
	return lcp.cp.ReadingTime()
}

func (lcp *LazyContentProvider) Len() int {
	lcp.init.Do()
	return lcp.cp.Len()
}

func (lcp *LazyContentProvider) Render(layout ...string) (template.HTML, error) {
	lcp.init.Do()
	return lcp.cp.Render(layout...)
}

func (lcp *LazyContentProvider) RenderString(args ...any) (template.HTML, error) {
	lcp.init.Do()
	return lcp.cp.RenderString(args...)
}

func (lcp *LazyContentProvider) TableOfContents() template.HTML {
	lcp.init.Do()
	return lcp.cp.TableOfContents()
}

func (lcp *LazyContentProvider) RenderContent(content []byte, renderTOC bool) (converter.Result, error) {
	lcp.init.Do()
	return lcp.cp.RenderContent(content, renderTOC)
}
