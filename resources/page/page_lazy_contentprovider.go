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
	"context"
	"html/template"

	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"
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
	lcp.init.Add(func(context.Context) (any, error) {
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

func (lcp *LazyContentProvider) TableOfContents(ctx context.Context) template.HTML {
	lcp.init.Do(ctx)
	return lcp.cp.TableOfContents(ctx)
}

func (lcp *LazyContentProvider) Fragments(ctx context.Context) *tableofcontents.Fragments {
	lcp.init.Do(ctx)
	return lcp.cp.Fragments(ctx)
}

func (lcp *LazyContentProvider) Content(ctx context.Context) (any, error) {
	lcp.init.Do(ctx)
	return lcp.cp.Content(ctx)
}

func (lcp *LazyContentProvider) Plain(ctx context.Context) string {
	lcp.init.Do(ctx)
	return lcp.cp.Plain(ctx)
}

func (lcp *LazyContentProvider) PlainWords(ctx context.Context) []string {
	lcp.init.Do(ctx)
	return lcp.cp.PlainWords(ctx)
}

func (lcp *LazyContentProvider) Summary(ctx context.Context) template.HTML {
	lcp.init.Do(ctx)
	return lcp.cp.Summary(ctx)
}

func (lcp *LazyContentProvider) Truncated(ctx context.Context) bool {
	lcp.init.Do(ctx)
	return lcp.cp.Truncated(ctx)
}

func (lcp *LazyContentProvider) FuzzyWordCount(ctx context.Context) int {
	lcp.init.Do(ctx)
	return lcp.cp.FuzzyWordCount(ctx)
}

func (lcp *LazyContentProvider) WordCount(ctx context.Context) int {
	lcp.init.Do(ctx)
	return lcp.cp.WordCount(ctx)
}

func (lcp *LazyContentProvider) ReadingTime(ctx context.Context) int {
	lcp.init.Do(ctx)
	return lcp.cp.ReadingTime(ctx)
}

func (lcp *LazyContentProvider) Len(ctx context.Context) int {
	lcp.init.Do(ctx)
	return lcp.cp.Len(ctx)
}

func (lcp *LazyContentProvider) Render(ctx context.Context, layout ...string) (template.HTML, error) {
	lcp.init.Do(ctx)
	return lcp.cp.Render(ctx, layout...)
}

func (lcp *LazyContentProvider) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	lcp.init.Do(ctx)
	return lcp.cp.RenderString(ctx, args...)
}

func (lcp *LazyContentProvider) ParseAndRenderContent(ctx context.Context, content []byte, renderTOC bool) (converter.ResultRender, error) {
	lcp.init.Do(ctx)
	return lcp.cp.ParseAndRenderContent(ctx, content, renderTOC)
}

func (lcp *LazyContentProvider) ParseContent(ctx context.Context, content []byte) (converter.ResultParse, bool, error) {
	lcp.init.Do(ctx)
	return lcp.cp.ParseContent(ctx, content)
}

func (lcp *LazyContentProvider) RenderContent(ctx context.Context, content []byte, doc any) (converter.ResultRender, bool, error) {
	lcp.init.Do(ctx)
	return lcp.cp.RenderContent(ctx, content, doc)
}
