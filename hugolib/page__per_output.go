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
	"fmt"
	"html/template"
	"strings"
	"sync"
	"unicode/utf8"

	"errors"

	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/lazy"

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

var pageContentOutputDependenciesID = identity.KeyValueIdentity{Key: "pageOutput", Value: "dependencies"}

func newPageContentOutput(p *pageState, po *pageOutput) (*pageContentOutput, error) {
	parent := p.init

	var dependencyTracker identity.Manager
	if p.s.running() {
		dependencyTracker = identity.NewManager(pageContentOutputDependenciesID)
	}

	cp := &pageContentOutput{
		dependencyTracker: dependencyTracker,
		p:                 p,
		f:                 po.f,
		renderHooks:       &renderHooks{},
	}

	initToC := func(ctx context.Context) (err error) {
		if p.cmap == nil {
			// Nothing to do.
			return nil
		}

		if err := po.cp.initRenderHooks(); err != nil {
			return err
		}

		f := po.f
		cp.contentPlaceholders, err = p.shortcodeState.prepareShortcodesForPage(ctx, p, f)
		if err != nil {
			return err
		}

		var hasVariants bool
		cp.workContent, hasVariants, err = p.contentToRender(ctx, p.source.parsed, p.cmap, cp.contentPlaceholders)
		if err != nil {
			return err
		}
		if hasVariants {
			p.pageOutputTemplateVariationsState.Store(2)
		}

		isHTML := cp.p.m.markup == "html"

		if !isHTML {
			createAndSetToC := func(tocProvider converter.TableOfContentsProvider) {
				cfg := p.s.ContentSpec.Converters.GetMarkupConfig()
				cp.tableOfContents = tocProvider.TableOfContents()
				cp.tableOfContentsHTML = template.HTML(
					cp.tableOfContents.ToHTML(
						cfg.TableOfContents.StartLevel,
						cfg.TableOfContents.EndLevel,
						cfg.TableOfContents.Ordered,
					),
				)
			}
			// If the converter supports doing the parsing separately, we do that.
			parseResult, ok, err := po.contentRenderer.ParseContent(ctx, cp.workContent)
			if err != nil {
				return err
			}
			if ok {
				// This is Goldmark.
				// Store away the parse result for later use.
				createAndSetToC(parseResult)
				cp.astDoc = parseResult.Doc()

				return nil
			}

			// This is Asciidoctor etc.
			r, err := po.contentRenderer.ParseAndRenderContent(ctx, cp.workContent, true)
			if err != nil {
				return err
			}

			cp.workContent = r.Bytes()

			if tocProvider, ok := r.(converter.TableOfContentsProvider); ok {
				createAndSetToC(tocProvider)
			} else {
				tmpContent, tmpTableOfContents := helpers.ExtractTOC(cp.workContent)
				cp.tableOfContentsHTML = helpers.BytesToHTML(tmpTableOfContents)
				cp.tableOfContents = tableofcontents.Empty
				cp.workContent = tmpContent
			}
		}

		return nil

	}

	initContent := func(ctx context.Context) (err error) {

		p.s.h.IncrContentRender()

		if p.cmap == nil {
			// Nothing to do.
			return nil
		}

		if cp.astDoc != nil {
			// The content is parsed, but not rendered.
			r, ok, err := po.contentRenderer.RenderContent(ctx, cp.workContent, cp.astDoc)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("invalid state: astDoc is set but RenderContent returned false")
			}

			cp.workContent = r.Bytes()
		}

		if p.cmap.hasNonMarkdownShortcode || cp.placeholdersEnabled {
			// There are one or more replacement tokens to be replaced.
			var hasShortcodeVariants bool
			tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
				if token == tocShortcodePlaceholder {
					// The Page's TableOfContents was accessed in a shortcode.
					if cp.tableOfContentsHTML == "" {
						cp.p.s.initInit(ctx, cp.initToC, cp.p)
					}
					return []byte(cp.tableOfContentsHTML), nil
				}
				renderer, found := cp.contentPlaceholders[token]
				if found {
					repl, more, err := renderer.renderShortcode(ctx)
					if err != nil {
						return nil, err
					}
					hasShortcodeVariants = hasShortcodeVariants || more
					return repl, nil
				}
				// This should never happen.
				return nil, fmt.Errorf("unknown shortcode token %q", token)
			}

			cp.workContent, err = expandShortcodeTokens(ctx, cp.workContent, tokenHandler)
			if err != nil {
				return err
			}
			if hasShortcodeVariants {
				p.pageOutputTemplateVariationsState.Store(2)
			}
		}

		if cp.p.source.hasSummaryDivider {
			isHTML := cp.p.m.markup == "html"
			if isHTML {
				src := p.source.parsed.Input()

				// Use the summary sections as they are provided by the user.
				if p.source.posSummaryEnd != -1 {
					cp.summary = helpers.BytesToHTML(src[p.source.posMainContent:p.source.posSummaryEnd])
				}

				if cp.p.source.posBodyStart != -1 {
					cp.workContent = src[cp.p.source.posBodyStart:]
				}

			} else {
				summary, content, err := splitUserDefinedSummaryAndContent(cp.p.m.markup, cp.workContent)
				if err != nil {
					cp.p.s.Log.Errorf("Failed to set user defined summary for page %q: %s", cp.p.pathOrTitle(), err)
				} else {
					cp.workContent = content
					cp.summary = helpers.BytesToHTML(summary)
				}
			}
		} else if cp.p.m.summary != "" {
			b, err := po.contentRenderer.ParseAndRenderContent(ctx, []byte(cp.p.m.summary), false)
			if err != nil {
				return err
			}
			html := cp.p.s.ContentSpec.TrimShortHTML(b.Bytes())
			cp.summary = helpers.BytesToHTML(html)
		}

		cp.content = helpers.BytesToHTML(cp.workContent)

		return nil
	}

	cp.initToC = parent.Branch(func(ctx context.Context) (any, error) {
		return nil, initToC(ctx)
	})

	// There may be recursive loops in shortcodes and render hooks.
	cp.initMain = cp.initToC.BranchWithTimeout(p.s.siteCfg.timeout, func(ctx context.Context) (any, error) {
		return nil, initContent(ctx)
	})

	cp.initPlain = cp.initMain.Branch(func(context.Context) (any, error) {
		cp.plain = tpl.StripHTML(string(cp.content))
		cp.plainWords = strings.Fields(cp.plain)
		cp.setWordCounts(p.m.isCJKLanguage)

		if err := cp.setAutoSummary(); err != nil {
			return err, nil
		}

		return nil, nil
	})

	return cp, nil
}

type renderHooks struct {
	getRenderer hooks.GetRendererFunc
	init        sync.Once
}

// pageContentOutput represents the Page content for a given output format.
type pageContentOutput struct {
	f output.Format

	p *pageState

	// Lazy load dependencies
	initToC   *lazy.Init
	initMain  *lazy.Init
	initPlain *lazy.Init

	placeholdersEnabled     bool
	placeholdersEnabledInit sync.Once

	// Renders Markdown hooks.
	renderHooks *renderHooks

	workContent       []byte
	dependencyTracker identity.Manager // Set in server mode.

	// Temporary storage of placeholders mapped to their content.
	// These are shortcodes etc. Some of these will need to be replaced
	// after any markup is rendered, so they share a common prefix.
	contentPlaceholders map[string]shortcodeRenderer

	// Content sections
	content             template.HTML
	summary             template.HTML
	tableOfContents     *tableofcontents.Fragments
	tableOfContentsHTML template.HTML
	// For Goldmark we split Parse and Render.
	astDoc any

	truncated bool

	plainWords     []string
	plain          string
	fuzzyWordCount int
	wordCount      int
	readingTime    int
}

func (p *pageContentOutput) trackDependency(id identity.Provider) {
	if p.dependencyTracker != nil {
		p.dependencyTracker.Add(id)
	}

}

func (p *pageContentOutput) Reset() {
	if p.dependencyTracker != nil {
		p.dependencyTracker.Reset()
	}
	p.initToC.Reset()
	p.initMain.Reset()
	p.initPlain.Reset()
	p.renderHooks = &renderHooks{}
}

func (p *pageContentOutput) Content(ctx context.Context) (any, error) {
	p.p.s.initInit(ctx, p.initMain, p.p)
	return p.content, nil
}

func (p *pageContentOutput) FuzzyWordCount(ctx context.Context) int {
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.fuzzyWordCount
}

func (p *pageContentOutput) Len(ctx context.Context) int {
	p.p.s.initInit(ctx, p.initMain, p.p)
	return len(p.content)
}

func (p *pageContentOutput) Plain(ctx context.Context) string {
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.plain
}

func (p *pageContentOutput) PlainWords(ctx context.Context) []string {
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.plainWords
}

func (p *pageContentOutput) ReadingTime(ctx context.Context) int {
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.readingTime
}

func (p *pageContentOutput) Summary(ctx context.Context) template.HTML {
	p.p.s.initInit(ctx, p.initMain, p.p)
	if !p.p.source.hasSummaryDivider {
		p.p.s.initInit(ctx, p.initPlain, p.p)
	}
	return p.summary
}

func (p *pageContentOutput) TableOfContents(ctx context.Context) template.HTML {
	p.p.s.initInit(ctx, p.initMain, p.p)
	return p.tableOfContentsHTML
}

func (p *pageContentOutput) Truncated(ctx context.Context) bool {
	if p.p.truncated {
		return true
	}
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.truncated
}

func (p *pageContentOutput) WordCount(ctx context.Context) int {
	p.p.s.initInit(ctx, p.initPlain, p.p)
	return p.wordCount
}

func (p *pageContentOutput) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
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

	if err = p.initRenderHooks(); err != nil {
		return "", err
	}

	conv := p.p.getContentConverter()
	if opts.Markup != "" && opts.Markup != p.p.m.markup {
		var err error
		// TODO(bep) consider cache
		conv, err = p.p.m.newContentConverter(p.p, opts.Markup)
		if err != nil {
			return "", p.p.wrapError(err)
		}
	}

	var rendered []byte

	if pageparser.HasShortcode(contentToRender) {
		// String contains a shortcode.
		parsed, err := pageparser.ParseMain(strings.NewReader(contentToRender), pageparser.Config{})
		if err != nil {
			return "", err
		}
		pm := &pageContentMap{
			items: make([]any, 0, 20),
		}
		s := newShortcodeHandler(p.p, p.p.s)

		if err := p.p.mapContentForResult(
			parsed,
			s,
			pm,
			opts.Markup,
			nil,
		); err != nil {
			return "", err
		}

		placeholders, err := s.prepareShortcodesForPage(ctx, p.p, p.f)
		if err != nil {
			return "", err
		}

		contentToRender, hasVariants, err := p.p.contentToRender(ctx, parsed, pm, placeholders)
		if err != nil {
			return "", err
		}
		if hasVariants {
			p.p.pageOutputTemplateVariationsState.Store(2)
		}
		b, err := p.renderContentWithConverter(ctx, conv, contentToRender, false)
		if err != nil {
			return "", p.p.wrapError(err)
		}
		rendered = b.Bytes()

		if pm.hasNonMarkdownShortcode || p.placeholdersEnabled {
			var hasShortcodeVariants bool

			tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
				if token == tocShortcodePlaceholder {
					// The Page's TableOfContents was accessed in a shortcode.
					if p.tableOfContentsHTML == "" {
						p.p.s.initInit(ctx, p.initToC, p.p)
					}
					return []byte(p.tableOfContentsHTML), nil
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
				p.p.pageOutputTemplateVariationsState.Store(2)
			}
		}

		// We need a consolidated view in $page.HasShortcode
		p.p.shortcodeState.transferNames(s)

	} else {
		c, err := p.renderContentWithConverter(ctx, conv, []byte(contentToRender), false)
		if err != nil {
			return "", p.p.wrapError(err)
		}

		rendered = c.Bytes()
	}

	if opts.Display == "inline" {
		// We may have to rethink this in the future when we get other
		// renderers.
		rendered = p.p.s.ContentSpec.TrimShortHTML(rendered)
	}

	return template.HTML(string(rendered)), nil
}

func (p *pageContentOutput) RenderWithTemplateInfo(ctx context.Context, info tpl.Info, layout ...string) (template.HTML, error) {
	p.p.addDependency(info)
	return p.Render(ctx, layout...)
}

func (p *pageContentOutput) Render(ctx context.Context, layout ...string) (template.HTML, error) {
	templ, found, err := p.p.resolveTemplate(layout...)
	if err != nil {
		return "", p.p.wrapError(err)
	}

	if !found {
		return "", nil
	}

	p.p.addDependency(templ.(tpl.Info))

	// Make sure to send the *pageState and not the *pageContentOutput to the template.
	res, err := executeToString(ctx, p.p.s.Tmpl(), templ, p.p)
	if err != nil {
		return "", p.p.wrapError(fmt.Errorf("failed to execute template %s: %w", templ.Name(), err))
	}
	return template.HTML(res), nil
}

func (p *pageContentOutput) initRenderHooks() error {
	if p == nil {
		return nil
	}

	p.renderHooks.init.Do(func() {
		if p.p.pageOutputTemplateVariationsState.Load() == 0 {
			p.p.pageOutputTemplateVariationsState.Store(1)
		}

		type cacheKey struct {
			tp hooks.RendererType
			id any
			f  output.Format
		}

		renderCache := make(map[cacheKey]any)
		var renderCacheMu sync.Mutex

		resolvePosition := func(ctx any) text.Position {
			var offset int

			switch v := ctx.(type) {
			case hooks.CodeblockContext:
				offset = bytes.Index(p.p.source.parsed.Input(), []byte(v.Inner()))
			}

			pos := p.p.posFromInput(p.p.source.parsed.Input(), offset)

			if pos.LineNumber > 0 {
				// Move up to the code fence delimiter.
				// This is in line with how we report on shortcodes.
				pos.LineNumber = pos.LineNumber - 1
			}

			return pos
		}

		p.renderHooks.getRenderer = func(tp hooks.RendererType, id any) any {
			renderCacheMu.Lock()
			defer renderCacheMu.Unlock()

			key := cacheKey{tp: tp, id: id, f: p.f}
			if r, ok := renderCache[key]; ok {
				return r
			}

			layoutDescriptor := p.p.getLayoutDescriptor()
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
				templ, found, err := p.p.s.Tmpl().LookupLayout(layoutDescriptor, f)
				if err != nil {
					panic(err)
				}
				return templ, found
			}

			templ, found1 := getHookTemplate(p.f)

			if p.p.reusePageOutputContent() {
				// Check if some of the other output formats would give a different template.
				for _, f := range p.p.s.renderFormats {
					if f.Name == p.f.Name {
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
							p.p.pageOutputTemplateVariationsState.Store(2)
							break
						}
					}
				}
			}
			if !found1 {
				if tp == hooks.CodeBlockRendererType {
					// No user provided tempplate for code blocks, so we use the native Go code version -- which is also faster.
					r := p.p.s.ContentSpec.Converters.GetHighlighter()
					renderCache[key] = r
					return r
				}
				return nil
			}

			r := hookRendererTemplate{
				templateHandler: p.p.s.Tmpl(),
				SearchProvider:  templ.(identity.SearchProvider),
				templ:           templ,
				resolvePosition: resolvePosition,
			}
			renderCache[key] = r
			return r
		}
	})

	return nil
}

func (p *pageContentOutput) setAutoSummary() error {
	if p.p.source.hasSummaryDivider || p.p.m.summary != "" {
		return nil
	}

	var summary string
	var truncated bool

	if p.p.m.isCJKLanguage {
		summary, truncated = p.p.s.ContentSpec.TruncateWordsByRune(p.plainWords)
	} else {
		summary, truncated = p.p.s.ContentSpec.TruncateWordsToWholeSentence(p.plain)
	}
	p.summary = template.HTML(summary)

	p.truncated = truncated

	return nil
}

func (cp *pageContentOutput) getContentConverter() (converter.Converter, error) {
	if err := cp.initRenderHooks(); err != nil {
		return nil, err
	}
	return cp.p.getContentConverter(), nil
}

func (cp *pageContentOutput) ParseAndRenderContent(ctx context.Context, content []byte, renderTOC bool) (converter.ResultRender, error) {
	c, err := cp.getContentConverter()
	if err != nil {
		return nil, err
	}
	return cp.renderContentWithConverter(ctx, c, content, renderTOC)
}

func (cp *pageContentOutput) ParseContent(ctx context.Context, content []byte) (converter.ResultParse, bool, error) {
	c, err := cp.getContentConverter()
	if err != nil {
		return nil, false, err
	}
	p, ok := c.(converter.ParseRenderer)
	if !ok {
		return nil, ok, nil
	}
	rctx := converter.RenderContext{
		Src:         content,
		RenderTOC:   true,
		GetRenderer: cp.renderHooks.getRenderer,
	}
	r, err := p.Parse(rctx)
	return r, ok, err

}
func (cp *pageContentOutput) RenderContent(ctx context.Context, content []byte, doc any) (converter.ResultRender, bool, error) {
	c, err := cp.getContentConverter()
	if err != nil {
		return nil, false, err
	}
	p, ok := c.(converter.ParseRenderer)
	if !ok {
		return nil, ok, nil
	}
	rctx := converter.RenderContext{
		Src:         content,
		RenderTOC:   true,
		GetRenderer: cp.renderHooks.getRenderer,
	}
	r, err := p.Render(rctx, doc)
	if err == nil {
		if ids, ok := r.(identity.IdentitiesProvider); ok {
			for _, v := range ids.GetIdentities() {
				cp.trackDependency(v)
			}
		}
	}

	return r, ok, err
}

func (cp *pageContentOutput) renderContentWithConverter(ctx context.Context, c converter.Converter, content []byte, renderTOC bool) (converter.ResultRender, error) {
	r, err := c.Convert(
		converter.RenderContext{
			Src:         content,
			RenderTOC:   renderTOC,
			GetRenderer: cp.renderHooks.getRenderer,
		})

	if err == nil {
		if ids, ok := r.(identity.IdentitiesProvider); ok {
			for _, v := range ids.GetIdentities() {
				cp.trackDependency(v)
			}
		}
	}

	return r, err
}

func (p *pageContentOutput) setWordCounts(isCJKLanguage bool) {
	if isCJKLanguage {
		p.wordCount = 0
		for _, word := range p.plainWords {
			runeCount := utf8.RuneCountInString(word)
			if len(word) == runeCount {
				p.wordCount++
			} else {
				p.wordCount += runeCount
			}
		}
	} else {
		p.wordCount = helpers.TotalWords(p.plain)
	}

	// TODO(bep) is set in a test. Fix that.
	if p.fuzzyWordCount == 0 {
		p.fuzzyWordCount = (p.wordCount + 100) / 100 * 100
	}

	if isCJKLanguage {
		p.readingTime = (p.wordCount + 500) / 501
	} else {
		p.readingTime = (p.wordCount + 212) / 213
	}
}

// A callback to signal that we have inserted a placeholder into the rendered
// content. This avoids doing extra replacement work.
func (p *pageContentOutput) enablePlaceholders() {
	p.placeholdersEnabledInit.Do(func() {
		p.placeholdersEnabled = true
	})
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
