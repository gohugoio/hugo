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
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hcontext"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/goldmark/hugocontext"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

const (
	internalSummaryDividerBase = "HUGOMORE42"
)

var (
	internalSummaryDividerPreString = "\n\n" + internalSummaryDividerBase + "\n\n"
	internalSummaryDividerPre       = []byte(internalSummaryDividerPreString)
)

type pageContentReplacement struct {
	val []byte

	source pageparser.Item
}

func (m *pageMeta) parseFrontMatter(h *HugoSites, pid uint64) (*contentParseInfo, error) {
	var (
		sourceKey            string
		openSource           hugio.OpenReadSeekCloser
		isFromContentAdapter = m.pageConfig.IsFromContentAdapter
	)

	if m.f != nil && !isFromContentAdapter {
		sourceKey = filepath.ToSlash(m.f.Filename())
		if !isFromContentAdapter {
			meta := m.f.FileInfo().Meta()
			openSource = func() (hugio.ReadSeekCloser, error) {
				r, err := meta.Open()
				if err != nil {
					return nil, fmt.Errorf("failed to open file %q: %w", meta.Filename, err)
				}
				return r, nil
			}
		}
	} else if isFromContentAdapter {
		openSource = m.pageConfig.Content.ValueAsOpenReadSeekCloser()
	}

	if sourceKey == "" {
		sourceKey = strconv.FormatUint(pid, 10)
	}

	pi := &contentParseInfo{
		h:          h,
		pid:        pid,
		sourceKey:  sourceKey,
		openSource: openSource,
	}

	source, err := pi.contentSource(m)
	if err != nil {
		return nil, err
	}

	items, err := pageparser.ParseBytes(
		source,
		pageparser.Config{
			NoFrontMatter: isFromContentAdapter,
		},
	)
	if err != nil {
		return nil, err
	}

	pi.itemsStep1 = items

	if isFromContentAdapter {
		// No front matter.
		return pi, nil
	}

	if err := pi.mapFrontMatter(source); err != nil {
		return nil, err
	}

	return pi, nil
}

func (m *pageMeta) newCachedContent(h *HugoSites, pi *contentParseInfo) (*cachedContent, error) {
	var filename string
	if m.f != nil {
		filename = m.f.Filename()
	}

	c := &cachedContent{
		pm:             m.s.pageMap,
		StaleInfo:      m,
		shortcodeState: newShortcodeHandler(filename, m.s),
		pi:             pi,
		enableEmoji:    m.s.conf.EnableEmoji,
		scopes:         maps.NewCache[string, *cachedContentScope](),
	}

	source, err := c.pi.contentSource(m)
	if err != nil {
		return nil, err
	}

	if err := c.parseContentFile(source); err != nil {
		return nil, err
	}

	return c, nil
}

type cachedContent struct {
	pm *pageMap

	resource.StaleInfo

	shortcodeState *shortcodeHandler

	// Parsed content.
	pi *contentParseInfo

	enableEmoji bool

	scopes *maps.Cache[string, *cachedContentScope]
}

func (c *cachedContent) getOrCreateScope(scope string, pco *pageContentOutput) *cachedContentScope {
	key := scope + pco.po.f.Name
	cs, _ := c.scopes.GetOrCreate(key, func() (*cachedContentScope, error) {
		return &cachedContentScope{
			cachedContent: c,
			pco:           pco,
			scope:         scope,
		}, nil
	})
	return cs
}

type contentParseInfo struct {
	h *HugoSites

	pid       uint64
	sourceKey string

	// The source bytes.
	openSource hugio.OpenReadSeekCloser

	frontMatter map[string]any

	// Whether the parsed content contains a summary separator.
	hasSummaryDivider bool

	// Returns the position in bytes after any front matter.
	posMainContent int

	// Indicates whether we must do placeholder replacements.
	hasNonMarkdownShortcode bool

	// Items from the page parser.
	// These maps directly to the source
	itemsStep1 pageparser.Items

	//  *shortcode, pageContentReplacement or pageparser.Item
	itemsStep2 []any
}

func (p *contentParseInfo) AddBytes(item pageparser.Item) {
	p.itemsStep2 = append(p.itemsStep2, item)
}

func (p *contentParseInfo) AddReplacement(val []byte, source pageparser.Item) {
	p.itemsStep2 = append(p.itemsStep2, pageContentReplacement{val: val, source: source})
}

func (p *contentParseInfo) AddShortcode(s *shortcode) {
	p.itemsStep2 = append(p.itemsStep2, s)
	if s.insertPlaceholder() {
		p.hasNonMarkdownShortcode = true
	}
}

// contentToRenderForItems returns the content to be processed by Goldmark or similar.
func (pi *contentParseInfo) contentToRender(ctx context.Context, source []byte, renderedShortcodes map[string]shortcodeRenderer) ([]byte, bool, error) {
	var hasVariants bool
	c := make([]byte, 0, len(source)+(len(source)/10))

	for _, it := range pi.itemsStep2 {
		switch v := it.(type) {
		case pageparser.Item:
			c = append(c, source[v.Pos():v.Pos()+len(v.Val(source))]...)
		case pageContentReplacement:
			c = append(c, v.val...)
		case *shortcode:
			if !v.insertPlaceholder() {
				// Insert the rendered shortcode.
				renderedShortcode, found := renderedShortcodes[v.placeholder]
				if !found {
					// This should never happen.
					panic(fmt.Sprintf("rendered shortcode %q not found", v.placeholder))
				}

				b, more, err := renderedShortcode.renderShortcode(ctx)
				if err != nil {
					return nil, false, fmt.Errorf("failed to render shortcode: %w", err)
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

	return c, hasVariants, nil
}

func (c *cachedContent) IsZero() bool {
	return len(c.pi.itemsStep2) == 0
}

func (c *cachedContent) parseContentFile(source []byte) error {
	if source == nil || c.pi.openSource == nil {
		return nil
	}

	return c.pi.mapItemsAfterFrontMatter(source, c.shortcodeState)
}

func (c *contentParseInfo) parseFrontMatter(it pageparser.Item, iter *pageparser.Iterator, source []byte) error {
	if c.frontMatter != nil {
		return nil
	}

	f := pageparser.FormatFromFrontMatterType(it.Type)
	var err error
	c.frontMatter, err = metadecoders.Default.UnmarshalToMap(it.Val(source), f)
	if err != nil {
		if fe, ok := err.(herrors.FileError); ok {
			pos := fe.Position()

			// Offset the starting position of front matter.
			offset := iter.LineNumber(source) - 1
			if f == metadecoders.YAML {
				offset -= 1
			}
			pos.LineNumber += offset

			fe.UpdatePosition(pos)
			fe.SetFilename("") // It will be set later.

			return fe
		} else {
			return err
		}
	}

	return nil
}

func (rn *contentParseInfo) failMap(source []byte, err error, i pageparser.Item) error {
	if fe, ok := err.(herrors.FileError); ok {
		return fe
	}

	pos := posFromInput("", source, i.Pos())

	return herrors.NewFileErrorFromPos(err, pos)
}

func (rn *contentParseInfo) mapFrontMatter(source []byte) error {
	if len(rn.itemsStep1) == 0 {
		return nil
	}
	iter := pageparser.NewIterator(rn.itemsStep1)

Loop:
	for {
		it := iter.Next()
		switch {
		case it.IsFrontMatter():
			if err := rn.parseFrontMatter(it, iter, source); err != nil {
				return err
			}
			next := iter.Peek()
			if !next.IsDone() {
				rn.posMainContent = next.Pos()
			}
			// Done.
			break Loop
		case it.IsEOF():
			break Loop
		case it.IsError():
			return rn.failMap(source, it.Err, it)
		default:

		}
	}

	return nil
}

func (rn *contentParseInfo) mapItemsAfterFrontMatter(
	source []byte,
	s *shortcodeHandler,
) error {
	if len(rn.itemsStep1) == 0 {
		return nil
	}

	fail := func(err error, i pageparser.Item) error {
		if fe, ok := err.(herrors.FileError); ok {
			return fe
		}

		pos := posFromInput("", source, i.Pos())

		return herrors.NewFileErrorFromPos(err, pos)
	}

	iter := pageparser.NewIterator(rn.itemsStep1)

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var ordinal int

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.IsFrontMatter():
			// Ignore.
		case it.Type == pageparser.TypeLeadSummaryDivider:
			posBody := -1
			f := func(item pageparser.Item) bool {
				if posBody == -1 && !item.IsDone() {
					posBody = item.Pos()
				}

				if item.IsNonWhitespace(source) {
					// Done
					return false
				}
				return true
			}
			iter.PeekWalk(f)

			rn.hasSummaryDivider = true

			// The content may be rendered by Goldmark or similar,
			// and we need to track the summary.
			rn.AddReplacement(internalSummaryDividerPre, it)

		// Handle shortcode
		case it.IsLeftShortcodeDelim():
			// let extractShortcode handle left delim (will do so recursively)
			iter.Backup()

			currShortcode, err := s.extractShortcode(ordinal, 0, source, iter)
			if err != nil {
				return fail(err, it)
			}

			currShortcode.pos = it.Pos()
			currShortcode.length = iter.Current().Pos() - it.Pos()
			if currShortcode.placeholder == "" {
				currShortcode.placeholder = createShortcodePlaceholder("s", rn.pid, currShortcode.ordinal)
			}

			if currShortcode.name != "" {
				s.addName(currShortcode.name)
			}

			if currShortcode.params == nil {
				var s []string
				currShortcode.params = s
			}

			currShortcode.placeholder = createShortcodePlaceholder("s", rn.pid, ordinal)
			ordinal++
			s.shortcodes = append(s.shortcodes, currShortcode)

			rn.AddShortcode(currShortcode)

		case it.IsEOF():
			break Loop
		case it.IsError():
			return fail(it.Err, it)
		default:
			rn.AddBytes(it)
		}
	}

	return nil
}

func (c *cachedContent) mustSource() []byte {
	source, err := c.pi.contentSource(c)
	if err != nil {
		panic(err)
	}
	return source
}

func (c *contentParseInfo) contentSource(s resource.StaleInfo) ([]byte, error) {
	key := c.sourceKey
	versionv := s.StaleVersion()

	v, err := c.h.cacheContentSource.GetOrCreate(key, func(string) (*resources.StaleValue[[]byte], error) {
		b, err := c.readSourceAll()
		if err != nil {
			return nil, err
		}

		return &resources.StaleValue[[]byte]{
			Value: b,
			StaleVersionFunc: func() uint32 {
				return s.StaleVersion() - versionv
			},
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return v.Value, nil
}

func (c *contentParseInfo) readSourceAll() ([]byte, error) {
	if c.openSource == nil {
		return []byte{}, nil
	}
	r, err := c.openSource()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

type contentTableOfContents struct {
	// For Goldmark we split Parse and Render.
	astDoc any

	tableOfContents     *tableofcontents.Fragments
	tableOfContentsHTML template.HTML

	// Temporary storage of placeholders mapped to their content.
	// These are shortcodes etc. Some of these will need to be replaced
	// after any markup is rendered, so they share a common prefix.
	contentPlaceholders map[string]shortcodeRenderer

	contentToRender []byte
}

type contentSummary struct {
	content               template.HTML
	contentWithoutSummary template.HTML
	summary               page.Summary
}

type contentPlainPlainWords struct {
	plain      string
	plainWords []string

	wordCount      int
	fuzzyWordCount int
	readingTime    int
}

func (c *cachedContentScope) keyScope(ctx context.Context) string {
	return hugo.GetMarkupScope(ctx) + c.pco.po.f.Name
}

func (c *cachedContentScope) contentRendered(ctx context.Context) (contentSummary, error) {
	cp := c.pco
	ctx = tpl.Context.DependencyScope.Set(ctx, pageDependencyScopeGlobal)
	key := c.pi.sourceKey + "/" + c.keyScope(ctx)
	versionv := c.version(cp)

	v, err := c.pm.cacheContentRendered.GetOrCreate(key, func(string) (*resources.StaleValue[contentSummary], error) {
		cp.po.p.s.Log.Trace(logg.StringFunc(func() string {
			return fmt.Sprintln("contentRendered", key)
		}))

		cp.po.p.s.h.contentRenderCounter.Add(1)
		cp.contentRendered.Store(true)
		po := cp.po

		ct, err := c.contentToC(ctx)
		if err != nil {
			return nil, err
		}

		rs, err := func() (*resources.StaleValue[contentSummary], error) {
			rs := &resources.StaleValue[contentSummary]{
				StaleVersionFunc: func() uint32 {
					return c.version(cp) - versionv
				},
			}

			if len(c.pi.itemsStep2) == 0 {
				// Nothing to do.
				return rs, nil
			}

			var b []byte

			if ct.astDoc != nil {
				// The content is parsed, but not rendered.
				r, ok, err := po.contentRenderer.RenderContent(ctx, ct.contentToRender, ct.astDoc)
				if err != nil {
					return nil, err
				}

				if !ok {
					return nil, errors.New("invalid state: astDoc is set but RenderContent returned false")
				}

				b = r.Bytes()

			} else {
				// Copy the content to be rendered.
				b = make([]byte, len(ct.contentToRender))
				copy(b, ct.contentToRender)
			}

			// There are one or more replacement tokens to be replaced.
			var hasShortcodeVariants bool
			tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
				if token == tocShortcodePlaceholder {
					return []byte(ct.tableOfContentsHTML), nil
				}
				renderer, found := ct.contentPlaceholders[token]
				if found {
					repl, more, err := renderer.renderShortcode(ctx)
					if err != nil {
						return nil, err
					}
					hasShortcodeVariants = hasShortcodeVariants || more
					return repl, nil
				}
				// This should never happen.
				panic(fmt.Errorf("unknown shortcode token %q (number of tokens: %d)", token, len(ct.contentPlaceholders)))
			}

			b, err = expandShortcodeTokens(ctx, b, tokenHandler)
			if err != nil {
				return nil, err
			}
			if hasShortcodeVariants {
				cp.po.p.pageOutputTemplateVariationsState.Add(1)
			}

			var result contentSummary
			if c.pi.hasSummaryDivider {
				s := string(b)
				summarized := page.ExtractSummaryFromHTMLWithDivider(cp.po.p.m.pageConfig.ContentMediaType, s, internalSummaryDividerBase)
				result.summary = page.Summary{
					Text:      template.HTML(summarized.Summary()),
					Type:      page.SummaryTypeManual,
					Truncated: summarized.Truncated(),
				}
				result.contentWithoutSummary = template.HTML(summarized.ContentWithoutSummary())
				result.content = template.HTML(summarized.Content())
			} else {
				result.content = template.HTML(string(b))
			}

			if !c.pi.hasSummaryDivider && cp.po.p.m.pageConfig.Summary == "" {
				numWords := cp.po.p.s.conf.SummaryLength
				isCJKLanguage := cp.po.p.m.pageConfig.IsCJKLanguage
				summary := page.ExtractSummaryFromHTML(cp.po.p.m.pageConfig.ContentMediaType, string(result.content), numWords, isCJKLanguage)
				result.summary = page.Summary{
					Text:      template.HTML(summary.Summary()),
					Type:      page.SummaryTypeAuto,
					Truncated: summary.Truncated(),
				}
				result.contentWithoutSummary = template.HTML(summary.ContentWithoutSummary())
			}
			rs.Value = result

			return rs, nil
		}()
		if err != nil {
			return rs, cp.po.p.wrapError(err)
		}

		if rs.Value.summary.IsZero() {
			b, err := cp.po.contentRenderer.ParseAndRenderContent(ctx, []byte(cp.po.p.m.pageConfig.Summary), false)
			if err != nil {
				return nil, err
			}
			html := cp.po.p.s.ContentSpec.TrimShortHTML(b.Bytes(), cp.po.p.m.pageConfig.Content.Markup)
			rs.Value.summary = page.Summary{
				Text: helpers.BytesToHTML(html),
				Type: page.SummaryTypeFrontMatter,
			}
			rs.Value.contentWithoutSummary = rs.Value.content
		}

		return rs, err
	})
	if err != nil {
		return contentSummary{}, cp.po.p.wrapError(err)
	}

	return v.Value, nil
}

func (c *cachedContentScope) mustContentToC(ctx context.Context) contentTableOfContents {
	ct, err := c.contentToC(ctx)
	if err != nil {
		panic(err)
	}
	return ct
}

var setGetContentCallbackInContext = hcontext.NewContextDispatcher[func(*pageContentOutput, contentTableOfContents)]("contentCallback")

func (c *cachedContentScope) contentToC(ctx context.Context) (contentTableOfContents, error) {
	cp := c.pco
	key := c.pi.sourceKey + "/" + c.keyScope(ctx)
	versionv := c.version(cp)

	v, err := c.pm.contentTableOfContents.GetOrCreate(key, func(string) (*resources.StaleValue[contentTableOfContents], error) {
		source, err := c.pi.contentSource(c)
		if err != nil {
			return nil, err
		}

		var ct contentTableOfContents
		if err := cp.initRenderHooks(); err != nil {
			return nil, err
		}
		f := cp.po.f
		po := cp.po
		p := po.p
		ct.contentPlaceholders, err = c.shortcodeState.prepareShortcodesForPage(ctx, p, f, false)
		if err != nil {
			return nil, err
		}

		// Callback called from below (e.g. in .RenderString)
		ctxCallback := func(cp2 *pageContentOutput, ct2 contentTableOfContents) {
			cp.otherOutputs.Set(cp2.po.p.pid, cp2)

			// Merge content placeholders
			for k, v := range ct2.contentPlaceholders {
				ct.contentPlaceholders[k] = v
			}

			if p.s.conf.Internal.Watch {
				for _, s := range cp2.po.p.m.content.shortcodeState.shortcodes {
					for _, templ := range s.templs {
						cp.trackDependency(templ.(identity.IdentityProvider))
					}
				}
			}

			// Transfer shortcode names so HasShortcode works for shortcodes from included pages.
			cp.po.p.m.content.shortcodeState.transferNames(cp2.po.p.m.content.shortcodeState)
			if cp2.po.p.pageOutputTemplateVariationsState.Load() > 0 {
				cp.po.p.pageOutputTemplateVariationsState.Add(1)
			}
		}

		ctx = setGetContentCallbackInContext.Set(ctx, ctxCallback)

		var hasVariants bool
		ct.contentToRender, hasVariants, err = c.pi.contentToRender(ctx, source, ct.contentPlaceholders)
		if err != nil {
			return nil, err
		}

		if hasVariants {
			p.pageOutputTemplateVariationsState.Add(1)
		}

		isHTML := cp.po.p.m.pageConfig.ContentMediaType.IsHTML()

		if !isHTML {
			createAndSetToC := func(tocProvider converter.TableOfContentsProvider) {
				cfg := p.s.ContentSpec.Converters.GetMarkupConfig()
				ct.tableOfContents = tocProvider.TableOfContents()
				ct.tableOfContentsHTML = template.HTML(
					ct.tableOfContents.ToHTML(
						cfg.TableOfContents.StartLevel,
						cfg.TableOfContents.EndLevel,
						cfg.TableOfContents.Ordered,
					),
				)
			}

			// If the converter supports doing the parsing separately, we do that.
			parseResult, ok, err := po.contentRenderer.ParseContent(ctx, ct.contentToRender)
			if err != nil {
				return nil, err
			}
			if ok {
				// This is Goldmark.
				// Store away the parse result for later use.
				createAndSetToC(parseResult)

				ct.astDoc = parseResult.Doc()

			} else {

				// This is Asciidoctor etc.
				r, err := po.contentRenderer.ParseAndRenderContent(ctx, ct.contentToRender, true)
				if err != nil {
					return nil, err
				}

				ct.contentToRender = r.Bytes()

				if tocProvider, ok := r.(converter.TableOfContentsProvider); ok {
					createAndSetToC(tocProvider)
				} else {
					tmpContent, tmpTableOfContents := helpers.ExtractTOC(ct.contentToRender)
					ct.tableOfContentsHTML = helpers.BytesToHTML(tmpTableOfContents)
					ct.tableOfContents = tableofcontents.Empty
					ct.contentToRender = tmpContent
				}
			}
		}

		return &resources.StaleValue[contentTableOfContents]{
			Value: ct,
			StaleVersionFunc: func() uint32 {
				return c.version(cp) - versionv
			},
		}, nil
	})
	if err != nil {
		return contentTableOfContents{}, err
	}

	return v.Value, nil
}

func (c *cachedContent) version(cp *pageContentOutput) uint32 {
	// Both of these gets incremented on change.
	return c.StaleVersion() + cp.contentRenderedVersion
}

func (c *cachedContentScope) contentPlain(ctx context.Context) (contentPlainPlainWords, error) {
	cp := c.pco
	key := c.pi.sourceKey + "/" + c.keyScope(ctx)

	versionv := c.version(cp)

	v, err := c.pm.cacheContentPlain.GetOrCreateWitTimeout(key, cp.po.p.s.Conf.Timeout(), func(string) (*resources.StaleValue[contentPlainPlainWords], error) {
		var result contentPlainPlainWords
		rs := &resources.StaleValue[contentPlainPlainWords]{
			StaleVersionFunc: func() uint32 {
				return c.version(cp) - versionv
			},
		}

		rendered, err := c.contentRendered(ctx)
		if err != nil {
			return nil, err
		}

		result.plain = tpl.StripHTML(string(rendered.content))
		result.plainWords = strings.Fields(result.plain)

		isCJKLanguage := cp.po.p.m.pageConfig.IsCJKLanguage

		if isCJKLanguage {
			result.wordCount = 0
			for _, word := range result.plainWords {
				runeCount := utf8.RuneCountInString(word)
				if len(word) == runeCount {
					result.wordCount++
				} else {
					result.wordCount += runeCount
				}
			}
		} else {
			result.wordCount = helpers.TotalWords(result.plain)
		}

		// TODO(bep) is set in a test. Fix that.
		if result.fuzzyWordCount == 0 {
			result.fuzzyWordCount = (result.wordCount + 100) / 100 * 100
		}

		if isCJKLanguage {
			result.readingTime = (result.wordCount + 500) / 501
		} else {
			result.readingTime = (result.wordCount + 212) / 213
		}

		rs.Value = result

		return rs, nil
	})
	if err != nil {
		if herrors.IsTimeoutError(err) {
			err = fmt.Errorf("timed out rendering the page content. You may have a circular loop in a shortcode, or your site may have resources that take longer to build than the `timeout` limit in your Hugo config file: %w", err)
		}
		return contentPlainPlainWords{}, err
	}
	return v.Value, nil
}

type cachedContentScope struct {
	*cachedContent
	pco   *pageContentOutput
	scope string
}

func (c *cachedContentScope) prepareContext(ctx context.Context) context.Context {
	// The markup scope is recursive, so if already set to a non zero value, preserve that value.
	if s := hugo.GetMarkupScope(ctx); s != "" || s == c.scope {
		return ctx
	}
	return hugo.SetMarkupScope(ctx, c.scope)
}

func (c *cachedContentScope) Render(ctx context.Context) (page.Content, error) {
	return c, nil
}

func (c *cachedContentScope) Content(ctx context.Context) (template.HTML, error) {
	ctx = c.prepareContext(ctx)
	cr, err := c.contentRendered(ctx)
	if err != nil {
		return "", err
	}
	return cr.content, nil
}

func (c *cachedContentScope) ContentWithoutSummary(ctx context.Context) (template.HTML, error) {
	ctx = c.prepareContext(ctx)
	cr, err := c.contentRendered(ctx)
	if err != nil {
		return "", err
	}
	return cr.contentWithoutSummary, nil
}

func (c *cachedContentScope) Summary(ctx context.Context) (page.Summary, error) {
	ctx = c.prepareContext(ctx)
	rendered, err := c.contentRendered(ctx)
	return rendered.summary, err
}

func (c *cachedContentScope) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	ctx = c.prepareContext(ctx)

	if len(args) < 1 || len(args) > 2 {
		return "", errors.New("want 1 or 2 arguments")
	}

	pco := c.pco

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
		if opts.Markup != "" {
			opts.Markup = markup.ResolveMarkup(opts.Markup)
		}
	}

	contentToRenderv := args[sidx]

	if _, ok := contentToRenderv.(hstring.HTML); ok {
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

	if opts.Markup != "" && opts.Markup != pco.po.p.m.pageConfig.ContentMediaType.SubType {
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
		parseInfo.itemsStep1, err = pageparser.ParseBytes(contentToRenderb, pageparser.Config{
			NoFrontMatter:    true,
			NoSummaryDivider: true,
		})
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
					toc, err := c.contentToC(ctx)
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
		markup := pco.po.p.m.pageConfig.Content.Markup
		if opts.Markup != "" {
			markup = pco.po.p.s.ContentSpec.ResolveMarkup(opts.Markup)
		}
		rendered = pco.po.p.s.ContentSpec.TrimShortHTML(rendered, markup)
	}

	return template.HTML(string(rendered)), nil
}

func (c *cachedContentScope) RenderShortcodes(ctx context.Context) (template.HTML, error) {
	ctx = c.prepareContext(ctx)

	pco := c.pco
	content := pco.po.p.m.content

	source, err := content.pi.contentSource(content)
	if err != nil {
		return "", err
	}
	ct, err := c.contentToC(ctx)
	if err != nil {
		return "", err
	}

	var insertPlaceholders bool
	var hasVariants bool
	cb := setGetContentCallbackInContext.Get(ctx)
	if cb != nil {
		insertPlaceholders = true
	}
	cc := make([]byte, 0, len(source)+(len(source)/10))
	for _, it := range content.pi.itemsStep2 {
		switch v := it.(type) {
		case pageparser.Item:
			cc = append(cc, source[v.Pos():v.Pos()+len(v.Val(source))]...)
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
				cc = append(cc, []byte(b)...)

			} else {
				// Insert the placeholder so we can insert the content after
				// markdown processing.
				cc = append(cc, []byte(v.placeholder)...)
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
		return template.HTML(hugocontext.Wrap(cc, pco.po.p.pid)), nil
	}

	return helpers.BytesToHTML(cc), nil
}

func (c *cachedContentScope) Plain(ctx context.Context) string {
	ctx = c.prepareContext(ctx)
	return c.mustContentPlain(ctx).plain
}

func (c *cachedContentScope) PlainWords(ctx context.Context) []string {
	ctx = c.prepareContext(ctx)
	return c.mustContentPlain(ctx).plainWords
}

func (c *cachedContentScope) WordCount(ctx context.Context) int {
	ctx = c.prepareContext(ctx)
	return c.mustContentPlain(ctx).wordCount
}

func (c *cachedContentScope) FuzzyWordCount(ctx context.Context) int {
	ctx = c.prepareContext(ctx)
	return c.mustContentPlain(ctx).fuzzyWordCount
}

func (c *cachedContentScope) ReadingTime(ctx context.Context) int {
	ctx = c.prepareContext(ctx)
	return c.mustContentPlain(ctx).readingTime
}

func (c *cachedContentScope) Len(ctx context.Context) int {
	ctx = c.prepareContext(ctx)
	return len(c.mustContentRendered(ctx).content)
}

func (c *cachedContentScope) Fragments(ctx context.Context) *tableofcontents.Fragments {
	ctx = c.prepareContext(ctx)
	toc := c.mustContentToC(ctx).tableOfContents
	if toc == nil {
		return nil
	}
	return toc
}

func (c *cachedContentScope) fragmentsHTML(ctx context.Context) template.HTML {
	ctx = c.prepareContext(ctx)
	return c.mustContentToC(ctx).tableOfContentsHTML
}

func (c *cachedContentScope) mustContentPlain(ctx context.Context) contentPlainPlainWords {
	r, err := c.contentPlain(ctx)
	if err != nil {
		c.pco.fail(err)
	}
	return r
}

func (c *cachedContentScope) mustContentRendered(ctx context.Context) contentSummary {
	r, err := c.contentRendered(ctx)
	if err != nil {
		c.pco.fail(err)
	}
	return r
}
