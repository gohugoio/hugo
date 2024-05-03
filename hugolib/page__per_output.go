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

	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/hugocontext"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/markup/converter"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/helpers"
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
		otherOutputs: make(map[uint64]*pageContentOutput),
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
	otherOutputs map[uint64]*pageContentOutput

	contentRenderedVersion uint32 // Incremented on reset.
	contentRendered        bool   // Set on content render.

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
	pco.contentRendered = false
	pco.renderHooks = &renderHooks{}
}

func (pco *pageContentOutput) Fragments(ctx context.Context) *tableofcontents.Fragments {
	return pco.po.p.m.content.mustContentToC(ctx, pco).tableOfContents
}

func (pco *pageContentOutput) RenderShortcodes(ctx context.Context) (template.HTML, error) {
	content := pco.po.p.m.content
	source, err := content.pi.contentSource(content)
	if err != nil {
		return "", err
	}
	ct, err := content.contentToC(ctx, pco)
	if err != nil {
		return "", err
	}

	var insertPlaceholders bool
	var hasVariants bool
	cb := setGetContentCallbackInContext.Get(ctx)
	if cb != nil {
		insertPlaceholders = true
	}
	c := make([]byte, 0, len(source)+(len(source)/10))
	for _, it := range content.pi.itemsStep2 {
		switch v := it.(type) {
		case pageparser.Item:
			c = append(c, source[v.Pos():v.Pos()+len(v.Val(source))]...)
		case pageContentReplacement:
			// Ignore.
		case *shortcode:
			if !insertPlaceholders || !v.insertPlaceholder() {
				// Insert the rendered shortcode.
				renderedShortcode, found := ct.contentPlaceholders[v.placeholder]
				if !found {
					// This should never happen.
					panic(fmt.Sprintf("rendered shortcode %q not found", v.placeholder))
				}

				b, more, err := renderedShortcode.renderShortcode(ctx)
				if err != nil {
					return "", fmt.Errorf("failed to render shortcode: %w", err)
				}
				hasVariants = hasVariants || more
				c = append(c, []byte(b)...)

			} else {
				// Insert the placeholder so we can insert the content after
				// markdown processing.
				c = append(c, []byte(v.placeholder)...)
			}
		default:
			panic(fmt.Sprintf("unknown item type %T", it))
		}
	}

	if hasVariants {
		pco.po.p.pageOutputTemplateVariationsState.Add(1)
	}

	if cb != nil {
		cb(pco, ct)
	}

	if tpl.Context.IsInGoldmark.Get(ctx) {
		// This content will be parsed and rendered by Goldmark.
		// Wrap it in a special Hugo markup to assign the correct Page from
		// the stack.
		return template.HTML(hugocontext.Wrap(c, pco.po.p.pid)), nil
	}

	return helpers.BytesToHTML(c), nil
}

func (pco *pageContentOutput) Content(ctx context.Context) (any, error) {
	r, err := pco.po.p.m.content.contentRendered(ctx, pco)
	return r.content, err
}

func (pco *pageContentOutput) TableOfContents(ctx context.Context) template.HTML {
	return pco.po.p.m.content.mustContentToC(ctx, pco).tableOfContentsHTML
}

func (p *pageContentOutput) Len(ctx context.Context) int {
	return len(p.mustContentRendered(ctx).content)
}

func (pco *pageContentOutput) mustContentRendered(ctx context.Context) contentSummary {
	r, err := pco.po.p.m.content.contentRendered(ctx, pco)
	if err != nil {
		pco.fail(err)
	}
	return r
}

func (pco *pageContentOutput) mustContentPlain(ctx context.Context) contentPlainPlainWords {
	r, err := pco.po.p.m.content.contentPlain(ctx, pco)
	if err != nil {
		pco.fail(err)
	}
	return r
}

func (pco *pageContentOutput) fail(err error) {
	pco.po.p.s.h.FatalError(pco.po.p.wrapError(err))
}

func (pco *pageContentOutput) Plain(ctx context.Context) string {
	return pco.mustContentPlain(ctx).plain
}

func (pco *pageContentOutput) PlainWords(ctx context.Context) []string {
	return pco.mustContentPlain(ctx).plainWords
}

func (pco *pageContentOutput) ReadingTime(ctx context.Context) int {
	return pco.mustContentPlain(ctx).readingTime
}

func (pco *pageContentOutput) WordCount(ctx context.Context) int {
	return pco.mustContentPlain(ctx).wordCount
}

func (pco *pageContentOutput) FuzzyWordCount(ctx context.Context) int {
	return pco.mustContentPlain(ctx).fuzzyWordCount
}

func (pco *pageContentOutput) Summary(ctx context.Context) template.HTML {
	return pco.mustContentPlain(ctx).summary
}

func (pco *pageContentOutput) Truncated(ctx context.Context) bool {
	return pco.mustContentPlain(ctx).summaryTruncated
}

func (pco *pageContentOutput) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	if len(args) < 1 || len(args) > 2 {
		return "", errors.New("want 1 or 2 arguments")
	}

	var contentToRender string
	opts := defaultRenderStringOpts
	sidx := 1

	if len(args) == 1 {
		sidx = 0
	} else {
		m, ok := args[0].(map[string]any)
		if !ok {
			return "", errors.New("first argument must be a map")
		}

		if err := mapstructure.WeakDecode(m, &opts); err != nil {
			return "", fmt.Errorf("failed to decode options: %w", err)
		}
	}

	contentToRenderv := args[sidx]

	if _, ok := contentToRenderv.(hstring.RenderedString); ok {
		// This content is already rendered, this is potentially
		// a infinite recursion.
		return "", errors.New("text is already rendered, repeating it may cause infinite recursion")
	}

	var err error
	contentToRender, err = cast.ToStringE(contentToRenderv)
	if err != nil {
		return "", err
	}

	if err = pco.initRenderHooks(); err != nil {
		return "", err
	}

	conv := pco.po.p.getContentConverter()
	if opts.Markup != "" && opts.Markup != pco.po.p.m.pageConfig.Markup {
		var err error
		conv, err = pco.po.p.m.newContentConverter(pco.po.p, opts.Markup)
		if err != nil {
			return "", pco.po.p.wrapError(err)
		}
	}

	var rendered []byte

	parseInfo := &contentParseInfo{
		h:   pco.po.p.s.h,
		pid: pco.po.p.pid,
	}

	if pageparser.HasShortcode(contentToRender) {
		contentToRenderb := []byte(contentToRender)
		// String contains a shortcode.
		parseInfo.itemsStep1, err = pageparser.ParseBytesMain(contentToRenderb, pageparser.Config{})
		if err != nil {
			return "", err
		}

		s := newShortcodeHandler(pco.po.p.pathOrTitle(), pco.po.p.s)
		if err := parseInfo.mapItemsAfterFrontMatter(contentToRenderb, s); err != nil {
			return "", err
		}

		placeholders, err := s.prepareShortcodesForPage(ctx, pco.po.p, pco.po.f, true)
		if err != nil {
			return "", err
		}

		contentToRender, hasVariants, err := parseInfo.contentToRender(ctx, contentToRenderb, placeholders)
		if err != nil {
			return "", err
		}
		if hasVariants {
			pco.po.p.pageOutputTemplateVariationsState.Add(1)
		}
		b, err := pco.renderContentWithConverter(ctx, conv, contentToRender, false)
		if err != nil {
			return "", pco.po.p.wrapError(err)
		}
		rendered = b.Bytes()

		if parseInfo.hasNonMarkdownShortcode {
			var hasShortcodeVariants bool

			tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
				if token == tocShortcodePlaceholder {
					toc, err := pco.po.p.m.content.contentToC(ctx, pco)
					if err != nil {
						return nil, err
					}
					// The Page's TableOfContents was accessed in a shortcode.
					return []byte(toc.tableOfContentsHTML), nil
				}
				renderer, found := placeholders[token]
				if found {
					repl, more, err := renderer.renderShortcode(ctx)
					if err != nil {
						return nil, err
					}
					hasShortcodeVariants = hasShortcodeVariants || more
					return repl, nil
				}
				// This should not happen.
				return nil, fmt.Errorf("unknown shortcode token %q", token)
			}

			rendered, err = expandShortcodeTokens(ctx, rendered, tokenHandler)
			if err != nil {
				return "", err
			}
			if hasShortcodeVariants {
				pco.po.p.pageOutputTemplateVariationsState.Add(1)
			}
		}

		// We need a consolidated view in $page.HasShortcode
		pco.po.p.m.content.shortcodeState.transferNames(s)

	} else {
		c, err := pco.renderContentWithConverter(ctx, conv, []byte(contentToRender), false)
		if err != nil {
			return "", pco.po.p.wrapError(err)
		}

		rendered = c.Bytes()
	}

	if opts.Display == "inline" {
		markup := pco.po.p.m.pageConfig.Markup
		if opts.Markup != "" {
			markup = pco.po.p.s.ContentSpec.ResolveMarkup(opts.Markup)
		}
		rendered = pco.po.p.s.ContentSpec.TrimShortHTML(rendered, markup)
	}

	return template.HTML(string(rendered)), nil
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
			case hooks.CodeblockContext:
				offset = bytes.Index(source, []byte(v.Inner()))
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

			if pco.po.p.reusePageOutputContent() {
				// Check if some of the other output formats would give a different template.
				for _, f := range pco.po.p.s.renderFormats {
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

func splitUserDefinedSummaryAndContent(markup string, c []byte) (summary []byte, content []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("summary split failed: %s", r)
		}
	}()

	startDivider := bytes.Index(c, internalSummaryDividerBaseBytes)

	if startDivider == -1 {
		return
	}

	startTag := "p"
	switch markup {
	case "asciidocext":
		startTag = "div"
	}

	// Walk back and forward to the surrounding tags.
	start := bytes.LastIndex(c[:startDivider], []byte("<"+startTag))
	end := bytes.Index(c[startDivider:], []byte("</"+startTag))

	if start == -1 {
		start = startDivider
	} else {
		start = startDivider - (startDivider - start)
	}

	if end == -1 {
		end = startDivider + len(internalSummaryDividerBase)
	} else {
		end = startDivider + end + len(startTag) + 3
	}

	var addDiv bool

	switch markup {
	case "rst":
		addDiv = true
	}

	withoutDivider := append(c[:start], bytes.Trim(c[end:], "\n")...)

	if len(withoutDivider) > 0 {
		summary = bytes.TrimSpace(withoutDivider[:start])
	}

	if addDiv {
		// For the rst
		summary = append(append([]byte(nil), summary...), []byte("</div>")...)
	}

	if err != nil {
		return
	}

	content = bytes.TrimSpace(withoutDivider)

	return
}
