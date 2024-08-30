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
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/identity"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/markup/converter"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/tpl"

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
	res, err := executeToString(ctx, pco.po.p.s.Tmpl(), templ, pco.po.p)
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

			layoutDescriptor := pco.po.p.getLayoutDescriptor()
			layoutDescriptor.RenderingHook = true
			layoutDescriptor.LayoutOverride = false
			layoutDescriptor.Layout = ""

			switch tp {
			case hooks.LinkRendererType:
				layoutDescriptor.Kind = "render-link"
			case hooks.ImageRendererType:
				layoutDescriptor.Kind = "render-image"
			case hooks.HeadingRendererType:
				layoutDescriptor.Kind = "render-heading"
			case hooks.PassthroughRendererType:
				layoutDescriptor.Kind = "render-passthrough"
				if id != nil {
					layoutDescriptor.KindVariants = id.(string)
				}
			case hooks.BlockquoteRendererType:
				layoutDescriptor.Kind = "render-blockquote"
				if id != nil {
					layoutDescriptor.KindVariants = id.(string)
				}
			case hooks.TableRendererType:
				layoutDescriptor.Kind = "render-table"
			case hooks.CodeBlockRendererType:
				layoutDescriptor.Kind = "render-codeblock"
				if id != nil {
					lang := id.(string)
					lexer := chromalexers.Get(lang)
					if lexer != nil {
						layoutDescriptor.KindVariants = strings.Join(lexer.Config().Aliases, ",")
					} else {
						layoutDescriptor.KindVariants = lang
					}
				}
			}

			getHookTemplate := func(f output.Format) (tpl.Template, bool) {
				templ, found, err := pco.po.p.s.Tmpl().LookupLayout(layoutDescriptor, f)
				if err != nil {
					panic(err)
				}
				if found {
					if isitp, ok := templ.(tpl.IsInternalTemplateProvider); ok && isitp.IsInternalTemplate() {
						renderHookConfig := pco.po.p.s.conf.Markup.Goldmark.RenderHooks
						switch templ.Name() {
						case "_default/_markup/render-link.html":
							if !renderHookConfig.Link.IsEnableDefault() {
								return nil, false
							}
						case "_default/_markup/render-image.html":
							if !renderHookConfig.Image.IsEnableDefault() {
								return nil, false
							}
						}
					}
				}
				return templ, found
			}

			templ, found1 := getHookTemplate(pco.po.f)

			if !found1 || pco.po.p.reusePageOutputContent() {
				// Some hooks may only be available in HTML, and if
				// this site is configured to not have HTML output, we need to
				// make sure we have a fallback. This should be very rare.
				candidates := pco.po.p.s.renderFormats
				if pco.po.f.MediaType.FirstSuffix.Suffix != "html" {
					if _, found := candidates.GetBySuffix("html"); !found {
						candidates = append(candidates, output.HTMLFormat)
					}
				}
				// Check if some of the other output formats would give a different template.
				for _, f := range candidates {
					if f.Name == pco.po.f.Name {
						continue
					}
					templ2, found2 := getHookTemplate(f)

					if found2 {
						if !found1 {
							templ = templ2
							found1 = true
							break
						}

						if templ != templ2 {
							pco.po.p.pageOutputTemplateVariationsState.Add(1)
							break
						}
					}
				}
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
				templateHandler: pco.po.p.s.Tmpl(),
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
}

type targetPathsHolder struct {
	paths page.TargetPaths
	page.OutputFormat
}

func (t targetPathsHolder) targetPaths() page.TargetPaths {
	return t.paths
}

func executeToString(ctx context.Context, h tpl.TemplateHandler, templ tpl.Template, data any) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := h.ExecuteWithContext(ctx, templ, b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
