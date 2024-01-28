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
	"io"
	"strings"
	"unicode/utf8"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hcontext"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
)

const (
	internalSummaryDividerBase = "HUGOMORE42"
)

var (
	internalSummaryDividerBaseBytes = []byte(internalSummaryDividerBase)
	internalSummaryDividerPre       = []byte("\n\n" + internalSummaryDividerBase + "\n\n")
)

type pageContentReplacement struct {
	val []byte

	source pageparser.Item
}

func newCachedContent(m *pageMeta, pid uint64) (*cachedContent, error) {
	var openSource hugio.OpenReadSeekCloser
	var filename string
	if m.f != nil {
		meta := m.f.FileInfo().Meta()
		openSource = func() (hugio.ReadSeekCloser, error) {
			r, err := meta.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file %q: %w", meta.Filename, err)
			}
			return r, nil
		}
		filename = m.f.Filename()
	}

	c := &cachedContent{
		pm:             m.s.pageMap,
		StaleInfo:      m,
		shortcodeState: newShortcodeHandler(filename, m.s),
		parseInfo: &contentParseInfo{
			pid: pid,
		},
		cacheBaseKey: m.pathInfo.PathNoLang(),
		openSource:   openSource,
		enableEmoji:  m.s.conf.EnableEmoji,
	}

	source, err := c.contentSource()
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

	cacheBaseKey string

	// The source bytes.
	openSource hugio.OpenReadSeekCloser

	resource.StaleInfo

	shortcodeState *shortcodeHandler

	// Parsed content.
	parseInfo *contentParseInfo

	enableEmoji bool
}

type contentParseInfo struct {
	pid         uint64
	frontMatter map[string]any

	// Whether the parsed content contains a summary separator.
	hasSummaryDivider bool

	// Whether there are more content after the summary divider.
	summaryTruncated bool

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
	return len(c.parseInfo.itemsStep2) == 0
}

func (c *cachedContent) parseContentFile(source []byte) error {
	if source == nil || c.openSource == nil {
		return nil
	}

	items, err := pageparser.ParseBytes(
		source,
		pageparser.Config{},
	)
	if err != nil {
		return err
	}

	c.parseInfo.itemsStep1 = items

	return c.parseInfo.mapItems(source, c.shortcodeState)
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

func (rn *contentParseInfo) mapItems(
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
			if err := rn.parseFrontMatter(it, iter, source); err != nil {
				return err
			}
			next := iter.Peek()
			if !next.IsDone() {
				rn.posMainContent = next.Pos()
			}
		case it.Type == pageparser.TypeLeadSummaryDivider:
			posBody := -1
			f := func(item pageparser.Item) bool {
				if posBody == -1 && !item.IsDone() {
					posBody = item.Pos()
				}

				if item.IsNonWhitespace(source) {
					rn.summaryTruncated = true

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
	source, err := c.contentSource()
	if err != nil {
		panic(err)
	}
	return source
}

func (c *cachedContent) contentSource() ([]byte, error) {
	key := c.cacheBaseKey
	v, err := c.pm.cacheContentSource.GetOrCreate(key, func(string) (*resources.StaleValue[[]byte], error) {
		b, err := c.readSourceAll()
		if err != nil {
			return nil, err
		}

		return &resources.StaleValue[[]byte]{
			Value: b,
			IsStaleFunc: func() bool {
				return c.IsStale()
			},
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return v.Value, nil
}

func (c *cachedContent) readSourceAll() ([]byte, error) {
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
	content          template.HTML
	summary          template.HTML
	summaryTruncated bool
}

type contentPlainPlainWords struct {
	plain      string
	plainWords []string

	summary          template.HTML
	summaryTruncated bool

	wordCount      int
	fuzzyWordCount int
	readingTime    int
}

func (c *cachedContent) contentRendered(ctx context.Context, cp *pageContentOutput) (contentSummary, error) {
	ctx = tpl.Context.DependencyScope.Set(ctx, pageDependencyScopeGlobal)
	key := c.cacheBaseKey + "/" + cp.po.f.Name
	versionv := cp.contentRenderedVersion

	v, err := c.pm.cacheContentRendered.GetOrCreate(key, func(string) (*resources.StaleValue[contentSummary], error) {
		cp.po.p.s.Log.Trace(logg.StringFunc(func() string {
			return fmt.Sprintln("contentRendered", key)
		}))

		cp.po.p.s.h.contentRenderCounter.Add(1)
		cp.contentRendered = true
		po := cp.po

		ct, err := c.contentToC(ctx, cp)
		if err != nil {
			return nil, err
		}

		rs := &resources.StaleValue[contentSummary]{
			IsStaleFunc: func() bool {
				return c.IsStale() || cp.contentRenderedVersion != versionv
			},
		}

		if len(c.parseInfo.itemsStep2) == 0 {
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

		var result contentSummary // hasVariants bool

		if c.parseInfo.hasSummaryDivider {
			isHTML := cp.po.p.m.markup == "html"
			if isHTML {
				// Use the summary sections as provided by the user.
				i := bytes.Index(b, internalSummaryDividerPre)
				result.summary = helpers.BytesToHTML(b[:i])
				b = b[i+len(internalSummaryDividerPre):]

			} else {
				summary, content, err := splitUserDefinedSummaryAndContent(cp.po.p.m.markup, b)
				if err != nil {
					cp.po.p.s.Log.Errorf("Failed to set user defined summary for page %q: %s", cp.po.p.pathOrTitle(), err)
				} else {
					b = content
					result.summary = helpers.BytesToHTML(summary)
				}
			}
			result.summaryTruncated = c.parseInfo.summaryTruncated
		}
		result.content = helpers.BytesToHTML(b)
		rs.Value = result

		return rs, nil
	})
	if err != nil {
		return contentSummary{}, cp.po.p.wrapError(err)
	}

	return v.Value, nil
}

func (c *cachedContent) mustContentToC(ctx context.Context, cp *pageContentOutput) contentTableOfContents {
	ct, err := c.contentToC(ctx, cp)
	if err != nil {
		panic(err)
	}
	return ct
}

var setGetContentCallbackInContext = hcontext.NewContextDispatcher[func(*pageContentOutput, contentTableOfContents)]("contentCallback")

func (c *cachedContent) contentToC(ctx context.Context, cp *pageContentOutput) (contentTableOfContents, error) {
	key := c.cacheBaseKey + "/" + cp.po.f.Name
	versionv := cp.contentRenderedVersion

	v, err := c.pm.contentTableOfContents.GetOrCreate(key, func(string) (*resources.StaleValue[contentTableOfContents], error) {
		source, err := c.contentSource()
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

		// Callback called from above (e.g. in .RenderString)
		ctxCallback := func(cp2 *pageContentOutput, ct2 contentTableOfContents) {
			// Merge content placeholders
			for k, v := range ct2.contentPlaceholders {
				ct.contentPlaceholders[k] = v
			}

			if p.s.conf.Internal.Watch {
				for _, s := range cp2.po.p.content.shortcodeState.shortcodes {
					for _, templ := range s.templs {
						cp.trackDependency(templ.(identity.IdentityProvider))
					}
				}
			}

			// Transfer shortcode names so HasShortcode works for shortcodes from included pages.
			cp.po.p.content.shortcodeState.transferNames(cp2.po.p.content.shortcodeState)
			if cp2.po.p.pageOutputTemplateVariationsState.Load() > 0 {
				cp.po.p.pageOutputTemplateVariationsState.Add(1)
			}
		}

		ctx = setGetContentCallbackInContext.Set(ctx, ctxCallback)

		var hasVariants bool
		ct.contentToRender, hasVariants, err = c.parseInfo.contentToRender(ctx, source, ct.contentPlaceholders)
		if err != nil {
			return nil, err
		}

		if hasVariants {
			p.pageOutputTemplateVariationsState.Add(1)
		}

		isHTML := cp.po.p.m.markup == "html"

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
			IsStaleFunc: func() bool {
				return c.IsStale() || cp.contentRenderedVersion != versionv
			},
		}, nil
	})
	if err != nil {
		return contentTableOfContents{}, err
	}

	return v.Value, nil
}

func (c *cachedContent) contentPlain(ctx context.Context, cp *pageContentOutput) (contentPlainPlainWords, error) {
	key := c.cacheBaseKey + "/" + cp.po.f.Name

	versionv := cp.contentRenderedVersion

	v, err := c.pm.cacheContentPlain.GetOrCreateWitTimeout(key, cp.po.p.s.Conf.Timeout(), func(string) (*resources.StaleValue[contentPlainPlainWords], error) {
		var result contentPlainPlainWords
		rs := &resources.StaleValue[contentPlainPlainWords]{
			IsStaleFunc: func() bool {
				return c.IsStale() || cp.contentRenderedVersion != versionv
			},
		}

		rendered, err := c.contentRendered(ctx, cp)
		if err != nil {
			return nil, err
		}

		result.plain = tpl.StripHTML(string(rendered.content))
		result.plainWords = strings.Fields(result.plain)

		isCJKLanguage := cp.po.p.m.isCJKLanguage

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

		if rendered.summary != "" {
			result.summary = rendered.summary
			result.summaryTruncated = rendered.summaryTruncated
		} else if cp.po.p.m.summary != "" {
			b, err := cp.po.contentRenderer.ParseAndRenderContent(ctx, []byte(cp.po.p.m.summary), false)
			if err != nil {
				return nil, err
			}
			html := cp.po.p.s.ContentSpec.TrimShortHTML(b.Bytes())
			result.summary = helpers.BytesToHTML(html)
		} else {
			var summary string
			var truncated bool
			if isCJKLanguage {
				summary, truncated = cp.po.p.s.ContentSpec.TruncateWordsByRune(result.plainWords)
			} else {
				summary, truncated = cp.po.p.s.ContentSpec.TruncateWordsToWholeSentence(result.plain)
			}
			result.summary = template.HTML(summary)
			result.summaryTruncated = truncated
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
