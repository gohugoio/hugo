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

package hugolib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/tpl/tplimpl"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/markup/converter"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	nopTargetPath    = targetPathsHolder{}
	nopPagePerOutput = struct {
		resource.ResourceLinksProvider
		page.ContentProvider
		page.PageRenderProvider
		page.PaginatorProvider
		page.TableOfContentsProvider
		page.AlternativeOutputFormatsProvider

		targetPather
	}{
		page.NopPage,
		page.NopPage,
		page.NopPage,
		page.NopPage,
		page.NopPage,
		page.NopPage,
		nopTargetPath,
	}
)

func newPageContentOutput(po *pageOutput) (*pageContentOutput, error) {
	cp := &pageContentOutput{
		po:           po,
		renderHooks:  &renderHooks{},
		otherOutputs: maps.NewCache[uint64, *pageContentOutput](),
	}
	return cp, nil
}

type renderHooks struct {
	getRenderer hooks.GetRendererFunc
	init        sync.Once
}

// pageContentOutput represents the Page content for a given output format.
type pageContentOutput struct {
	po *pageOutput

	// Other pages involved in rendering of this page,
	// typically included with .RenderShortcodes.
	otherOutputs *maps.Cache[uint64, *pageContentOutput]

	contentRenderedVersion uint32      // Incremented on reset.
	contentRendered        atomic.Bool // Set on content render.

	// Renders Markdown hooks.
	renderHooks *renderHooks
}

func (pco *pageContentOutput) trackDependency(idp identity.IdentityProvider) {
	pco.po.p.dependencyManagerOutput.AddIdentity(idp.GetIdentity())
}

func (pco *pageContentOutput) Reset() {
	if pco == nil {
		return
	}
	pco.contentRenderedVersion++
	pco.contentRendered.Store(false)
	pco.renderHooks = &renderHooks{}
}

func (pco *pageContentOutput) Render(ctx context.Context, layout ...string) (template.HTML, error) {
	if len(layout) == 0 {
		return "", errors.New("no layout given")
	}
	templ, found, err := pco.po.p.resolveTemplate(layout...)
	if err != nil {
		return "", pco.po.p.wrapError(err)
	}

	if !found {
		return "", nil
	}

	// Make sure to send the *pageState and not the *pageContentOutput to the template.
	res, err := executeToString(ctx, pco.po.p.s.GetTemplateStore(), templ, pco.po.p)
	if err != nil {
		return "", pco.po.p.wrapError(fmt.Errorf("failed to execute template %s: %w", templ.Name(), err))
	}
	return template.HTML(res), nil
}

func (pco *pageContentOutput) Fragments(ctx context.Context) *tableofcontents.Fragments {
	return pco.c().Fragments(ctx)
}

func (pco *pageContentOutput) RenderShortcodes(ctx context.Context) (template.HTML, error) {
	return pco.c().RenderShortcodes(ctx)
}

func (pco *pageContentOutput) Markup(opts ...any) page.Markup {
	if len(opts) > 1 {
		panic("too many arguments, expected 0 or 1")
	}
	var scope string
	if len(opts) == 1 {
		scope = cast.ToString(opts[0])
	}
	return pco.po.p.m.content.getOrCreateScope(scope, pco)
}

func (pco *pageContentOutput) c() page.Markup {
	return pco.po.p.m.content.getOrCreateScope("", pco)
}

func (pco *pageContentOutput) Content(ctx context.Context) (any, error) {
	r, err := pco.c().Render(ctx)
	if err != nil {
		return nil, err
	}
	return r.Content(ctx)
}

func (pco *pageContentOutput) ContentWithoutSummary(ctx context.Context) (template.HTML, error) {
	r, err := pco.c().Render(ctx)
	if err != nil {
		return "", err
	}
	return r.ContentWithoutSummary(ctx)
}

func (pco *pageContentOutput) TableOfContents(ctx context.Context) template.HTML {
	return pco.c().(*cachedContentScope).fragmentsHTML(ctx)
}

func (pco *pageContentOutput) Len(ctx context.Context) int {
	return pco.mustRender(ctx).Len(ctx)
}

func (pco *pageContentOutput) mustRender(ctx context.Context) page.Content {
	c, err := pco.c().Render(ctx)
	if err != nil {
		pco.fail(err)
	}
	return c
}

func (pco *pageContentOutput) fail(err error) {
	pco.po.p.s.h.FatalError(pco.po.p.wrapError(err))
}

func (pco *pageContentOutput) Plain(ctx context.Context) string {
	return pco.mustRender(ctx).Plain(ctx)
}

func (pco *pageContentOutput) PlainWords(ctx context.Context) []string {
	return pco.mustRender(ctx).PlainWords(ctx)
}

func (pco *pageContentOutput) ReadingTime(ctx context.Context) int {
	return pco.mustRender(ctx).ReadingTime(ctx)
}

func (pco *pageContentOutput) WordCount(ctx context.Context) int {
	return pco.mustRender(ctx).WordCount(ctx)
}

func (pco *pageContentOutput) FuzzyWordCount(ctx context.Context) int {
	return pco.mustRender(ctx).FuzzyWordCount(ctx)
}

func (pco *pageContentOutput) Summary(ctx context.Context) template.HTML {
	summary, err := pco.mustRender(ctx).Summary(ctx)
	if err != nil {
		pco.fail(err)
	}
	return summary.Text
}

func (pco *pageContentOutput) Truncated(ctx context.Context) bool {
	summary, err := pco.mustRender(ctx).Summary(ctx)
	if err != nil {
		pco.fail(err)
	}
	return summary.Truncated
}

func (pco *pageContentOutput) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	return pco.c().RenderString(ctx, args...)
}

func (pco *pageContentOutput) initRenderHooks() error {
	if pco == nil {
		return nil
	}

	pco.renderHooks.init.Do(func() {
		if pco.po.p.pageOutputTemplateVariationsState.Load() == 0 {
			pco.po.p.pageOutputTemplateVariationsState.Store(1)
		}

		type cacheKey struct {
			tp hooks.RendererType
			id any
			f  output.Format
		}

		renderCache := make(map[cacheKey]any)
		var renderCacheMu sync.Mutex

		resolvePosition := func(ctx any) text.Position {
			source := pco.po.p.m.content.mustSource()
			var offset int

			switch v := ctx.(type) {
			case hooks.PositionerSourceTargetProvider:
				offset = bytes.Index(source, v.PositionerSourceTarget())
			}

			pos := pco.po.p.posFromInput(source, offset)

			if pos.LineNumber > 0 {
				// Move up to the code fence delimiter.
				// This is in line with how we report on shortcodes.
				pos.LineNumber = pos.LineNumber - 1
			}

			return pos
		}

		pco.renderHooks.getRenderer = func(tp hooks.RendererType, id any) any {
			renderCacheMu.Lock()
			defer renderCacheMu.Unlock()

			key := cacheKey{tp: tp, id: id, f: pco.po.f}
			if r, ok := renderCache[key]; ok {
				return r
			}

			// Inherit the descriptor from the page/current output format.
			// This allows for fine-grained control of the template used for
			// rendering of e.g. links.
			base, layoutDescriptor := pco.po.p.GetInternalTemplateBasePathAndDescriptor()

			switch tp {
			case hooks.LinkRendererType:
				layoutDescriptor.Variant1 = "link"
			case hooks.ImageRendererType:
				layoutDescriptor.Variant1 = "image"
			case hooks.HeadingRendererType:
				layoutDescriptor.Variant1 = "heading"
			case hooks.PassthroughRendererType:
				layoutDescriptor.Variant1 = "passthrough"
				if id != nil {
					layoutDescriptor.Variant2 = id.(string)
				}
			case hooks.BlockquoteRendererType:
				layoutDescriptor.Variant1 = "blockquote"
				if id != nil {
					layoutDescriptor.Variant2 = id.(string)
				}
			case hooks.TableRendererType:
				layoutDescriptor.Variant1 = "table"
			case hooks.CodeBlockRendererType:
				layoutDescriptor.Variant1 = "codeblock"
				if id != nil {
					layoutDescriptor.Variant2 = id.(string)
				}
			}

			renderHookConfig := pco.po.p.s.conf.Markup.Goldmark.RenderHooks
			var ignoreInternal bool
			switch layoutDescriptor.Variant1 {
			case "link":
				ignoreInternal = !renderHookConfig.Link.IsEnableDefault()
			case "image":
				ignoreInternal = !renderHookConfig.Image.IsEnableDefault()
			}

			candidates := pco.po.p.s.renderFormats
			var numCandidatesFound int
			consider := func(candidate *tplimpl.TemplInfo) bool {
				if layoutDescriptor.Variant1 != candidate.D.Variant1 {
					return false
				}

				if layoutDescriptor.Variant2 != "" && candidate.D.Variant2 != "" && layoutDescriptor.Variant2 != candidate.D.Variant2 {
					return false
				}

				if ignoreInternal && candidate.SubCategory() == tplimpl.SubCategoryEmbedded {
					// Don't consider the internal hook templates.
					return false
				}

				if pco.po.p.pageOutputTemplateVariationsState.Load() > 1 {
					return true
				}

				if candidate.D.OutputFormat == "" {
					numCandidatesFound++
				} else if _, found := candidates.GetByName(candidate.D.OutputFormat); found {
					numCandidatesFound++
				}

				return true
			}

			getHookTemplate := func() (*tplimpl.TemplInfo, bool) {
				q := tplimpl.TemplateQuery{
					Path:     base,
					Category: tplimpl.CategoryMarkup,
					Desc:     layoutDescriptor,
					Consider: consider,
				}

				v := pco.po.p.s.TemplateStore.LookupPagesLayout(q)
				return v, v != nil
			}

			templ, found1 := getHookTemplate()
			if found1 && templ == nil {
				panic("found1 is true, but templ is nil")
			}

			if !found1 && layoutDescriptor.OutputFormat == pco.po.p.s.conf.DefaultOutputFormat {
				numCandidatesFound++
			}

			if numCandidatesFound > 1 {
				// More than one output format candidate found for this hook temoplate,
				// so we cannot reuse the same rendered content.
				pco.po.p.incrPageOutputTemplateVariation()
			}

			if !found1 {
				if tp == hooks.CodeBlockRendererType {
					// No user provided template for code blocks, so we use the native Go version -- which is also faster.
					r := pco.po.p.s.ContentSpec.Converters.GetHighlighter()
					renderCache[key] = r
					return r
				}
				return nil
			}

			r := hookRendererTemplate{
				templateHandler: pco.po.p.s.GetTemplateStore(),
				templ:           templ,
				resolvePosition: resolvePosition,
			}
			renderCache[key] = r
			return r
		}
	})

	return nil
}

func (pco *pageContentOutput) getContentConverter() (converter.Converter, error) {
	if err := pco.initRenderHooks(); err != nil {
		return nil, err
	}
	return pco.po.p.getContentConverter(), nil
}

func (cp *pageContentOutput) ParseAndRenderContent(ctx context.Context, content []byte, renderTOC bool) (converter.ResultRender, error) {
	c, err := cp.getContentConverter()
	if err != nil {
		return nil, err
	}
	return cp.renderContentWithConverter(ctx, c, content, renderTOC)
}

func (pco *pageContentOutput) ParseContent(ctx context.Context, content []byte) (converter.ResultParse, bool, error) {
	c, err := pco.getContentConverter()
	if err != nil {
		return nil, false, err
	}
	p, ok := c.(converter.ParseRenderer)
	if !ok {
		return nil, ok, nil
	}
	rctx := converter.RenderContext{
		Ctx:         ctx,
		Src:         content,
		RenderTOC:   true,
		GetRenderer: pco.renderHooks.getRenderer,
	}
	r, err := p.Parse(rctx)
	return r, ok, err
}

func (pco *pageContentOutput) RenderContent(ctx context.Context, content []byte, doc any) (converter.ResultRender, bool, error) {
	c, err := pco.getContentConverter()
	if err != nil {
		return nil, false, err
	}
	p, ok := c.(converter.ParseRenderer)
	if !ok {
		return nil, ok, nil
	}
	rctx := converter.RenderContext{
		Ctx:         ctx,
		Src:         content,
		RenderTOC:   true,
		GetRenderer: pco.renderHooks.getRenderer,
	}
	r, err := p.Render(rctx, doc)
	return r, ok, err
}

func (pco *pageContentOutput) renderContentWithConverter(ctx context.Context, c converter.Converter, content []byte, renderTOC bool) (converter.ResultRender, error) {
	r, err := c.Convert(
		converter.RenderContext{
			Ctx:         ctx,
			Src:         content,
			RenderTOC:   renderTOC,
			GetRenderer: pco.renderHooks.getRenderer,
		})
	return r, err
}

// these will be shifted out when rendering a given output format.
type pagePerOutputProviders interface {
	targetPather
	page.PaginatorProvider
	resource.ResourceLinksProvider
}

type targetPather interface {
	targetPaths() page.TargetPaths
	getRelURL() string
}

type targetPathsHolder struct {
	// relURL is usually the same as OutputFormat.RelPermalink, but can be different
	// for non-permalinkable output formats. These shares RelPermalink with the main (first) output format.
	relURL string
	paths  page.TargetPaths
	page.OutputFormat
}

func (t targetPathsHolder) getRelURL() string {
	return t.relURL
}

func (t targetPathsHolder) targetPaths() page.TargetPaths {
	return t.paths
}

func executeToString(ctx context.Context, h *tplimpl.TemplateStore, templ *tplimpl.TemplInfo, data any) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := h.ExecuteWithContext(ctx, templ, b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
