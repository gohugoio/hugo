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
	"runtime/debug"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/gohugoio/hugo/identity"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/markup/converter/hooks"

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

	initContent := func() (err error) {
		p.s.h.IncrContentRender()

		if p.cmap == nil {
			// Nothing to do.
			return nil
		}
		defer func() {
			// See https://github.com/gohugoio/hugo/issues/6210
			if r := recover(); r != nil {
				err = fmt.Errorf("%s", r)
				p.s.Log.Errorf("[BUG] Got panic:\n%s\n%s", r, string(debug.Stack()))
			}
		}()

		if err := po.cp.initRenderHooks(); err != nil {
			return err
		}

		var hasShortcodeVariants bool

		f := po.f
		cp.contentPlaceholders, hasShortcodeVariants, err = p.shortcodeState.renderShortcodesForPage(p, f)
		if err != nil {
			return err
		}

		enableReuse := !(hasShortcodeVariants || cp.renderHooksHaveVariants)

		if enableReuse {
			// Reuse this for the other output formats.
			// We may improve on this, but we really want to avoid re-rendering the content
			// to all output formats.
			// The current rule is that if you need output format-aware shortcodes or
			// content rendering hooks, create a output format-specific template, e.g.
			// myshortcode.amp.html.
			cp.enableReuse()
		}

		cp.workContent = p.contentToRender(cp.contentPlaceholders)

		isHTML := cp.p.m.markup == "html"

		if !isHTML {
			r, err := cp.renderContent(cp.workContent, true)
			if err != nil {
				return err
			}

			cp.workContent = r.Bytes()

			if tocProvider, ok := r.(converter.TableOfContentsProvider); ok {
				cfg := p.s.ContentSpec.Converters.GetMarkupConfig()
				cp.tableOfContents = template.HTML(
					tocProvider.TableOfContents().ToHTML(
						cfg.TableOfContents.StartLevel,
						cfg.TableOfContents.EndLevel,
						cfg.TableOfContents.Ordered,
					),
				)
			} else {
				tmpContent, tmpTableOfContents := helpers.ExtractTOC(cp.workContent)
				cp.tableOfContents = helpers.BytesToHTML(tmpTableOfContents)
				cp.workContent = tmpContent
			}
		}

		if cp.placeholdersEnabled {
			// ToC was accessed via .Page.TableOfContents in the shortcode,
			// at a time when the ToC wasn't ready.
			cp.contentPlaceholders[tocShortcodePlaceholder] = string(cp.tableOfContents)
		}

		if p.cmap.hasNonMarkdownShortcode || cp.placeholdersEnabled {
			// There are one or more replacement tokens to be replaced.
			cp.workContent, err = replaceShortcodeTokens(cp.workContent, cp.contentPlaceholders)
			if err != nil {
				return err
			}
		}

		if cp.p.source.hasSummaryDivider {
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
			b, err := cp.renderContent([]byte(cp.p.m.summary), false)
			if err != nil {
				return err
			}
			html := cp.p.s.ContentSpec.TrimShortHTML(b.Bytes())
			cp.summary = helpers.BytesToHTML(html)
		}

		cp.content = helpers.BytesToHTML(cp.workContent)

		return nil
	}

	// Recursive loops can only happen in content files with template code (shortcodes etc.)
	// Avoid creating new goroutines if we don't have to.
	needTimeout := p.shortcodeState.hasShortcodes() || cp.renderHooks != nil

	if needTimeout {
		cp.initMain = parent.BranchWithTimeout(p.s.siteCfg.timeout, func(ctx context.Context) (interface{}, error) {
			return nil, initContent()
		})
	} else {
		cp.initMain = parent.Branch(func() (interface{}, error) {
			return nil, initContent()
		})
	}

	cp.initPlain = cp.initMain.Branch(func() (interface{}, error) {
		cp.plain = helpers.StripHTML(string(cp.content))
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
	hooks hooks.Renderers
	init  sync.Once
}

// pageContentOutput represents the Page content for a given output format.
type pageContentOutput struct {
	f output.Format

	// If we can reuse this for other output formats.
	reuse     bool
	reuseInit sync.Once

	p *pageState

	// Lazy load dependencies
	initMain  *lazy.Init
	initPlain *lazy.Init

	placeholdersEnabled     bool
	placeholdersEnabledInit sync.Once

	renderHooks *renderHooks

	// Set if there are more than one output format variant
	renderHooksHaveVariants bool // TODO(bep) reimplement this in another way, consolidate with shortcodes

	// Content state

	workContent       []byte
	dependencyTracker identity.Manager // Set in server mode.

	// Temporary storage of placeholders mapped to their content.
	// These are shortcodes etc. Some of these will need to be replaced
	// after any markup is rendered, so they share a common prefix.
	contentPlaceholders map[string]string

	// Content sections
	content         template.HTML
	summary         template.HTML
	tableOfContents template.HTML

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
	p.initMain.Reset()
	p.initPlain.Reset()
	p.renderHooks = &renderHooks{}
}

func (p *pageContentOutput) Content() (interface{}, error) {
	if p.p.s.initInit(p.initMain, p.p) {
		return p.content, nil
	}
	return nil, nil
}

func (p *pageContentOutput) FuzzyWordCount() int {
	p.p.s.initInit(p.initPlain, p.p)
	return p.fuzzyWordCount
}

func (p *pageContentOutput) Len() int {
	p.p.s.initInit(p.initMain, p.p)
	return len(p.content)
}

func (p *pageContentOutput) Plain() string {
	p.p.s.initInit(p.initPlain, p.p)
	return p.plain
}

func (p *pageContentOutput) PlainWords() []string {
	p.p.s.initInit(p.initPlain, p.p)
	return p.plainWords
}

func (p *pageContentOutput) ReadingTime() int {
	p.p.s.initInit(p.initPlain, p.p)
	return p.readingTime
}

func (p *pageContentOutput) Summary() template.HTML {
	p.p.s.initInit(p.initMain, p.p)
	if !p.p.source.hasSummaryDivider {
		p.p.s.initInit(p.initPlain, p.p)
	}
	return p.summary
}

func (p *pageContentOutput) TableOfContents() template.HTML {
	p.p.s.initInit(p.initMain, p.p)
	return p.tableOfContents
}

func (p *pageContentOutput) Truncated() bool {
	if p.p.truncated {
		return true
	}
	p.p.s.initInit(p.initPlain, p.p)
	return p.truncated
}

func (p *pageContentOutput) WordCount() int {
	p.p.s.initInit(p.initPlain, p.p)
	return p.wordCount
}

func (p *pageContentOutput) RenderString(args ...interface{}) (template.HTML, error) {
	if len(args) < 1 || len(args) > 2 {
		return "", errors.New("want 1 or 2 arguments")
	}

	var s string
	opts := defaultRenderStringOpts
	sidx := 1

	if len(args) == 1 {
		sidx = 0
	} else {
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return "", errors.New("first argument must be a map")
		}

		if err := mapstructure.WeakDecode(m, &opts); err != nil {
			return "", errors.WithMessage(err, "failed to decode options")
		}
	}

	var err error
	s, err = cast.ToStringE(args[sidx])
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
		conv, err = p.p.m.newContentConverter(p.p, opts.Markup, nil)
		if err != nil {
			return "", p.p.wrapError(err)
		}
	}

	c, err := p.renderContentWithConverter(conv, []byte(s), false)
	if err != nil {
		return "", p.p.wrapError(err)
	}

	b := c.Bytes()

	if opts.Display == "inline" {
		// We may have to rethink this in the future when we get other
		// renderers.
		b = p.p.s.ContentSpec.TrimShortHTML(b)
	}

	return template.HTML(string(b)), nil
}

func (p *pageContentOutput) RenderWithTemplateInfo(info tpl.Info, layout ...string) (template.HTML, error) {
	p.p.addDependency(info)
	return p.Render(layout...)
}

func (p *pageContentOutput) Render(layout ...string) (template.HTML, error) {
	templ, found, err := p.p.resolveTemplate(layout...)
	if err != nil {
		return "", p.p.wrapError(err)
	}

	if !found {
		return "", nil
	}

	p.p.addDependency(templ.(tpl.Info))

	// Make sure to send the *pageState and not the *pageContentOutput to the template.
	res, err := executeToString(p.p.s.Tmpl(), templ, p.p)
	if err != nil {
		return "", p.p.wrapError(errors.Wrapf(err, "failed to execute template %q v", layout))
	}
	return template.HTML(res), nil
}

func (p *pageContentOutput) initRenderHooks() error {
	if p == nil {
		return nil
	}

	var initErr error

	p.renderHooks.init.Do(func() {
		ps := p.p

		c := ps.getContentConverter()
		if c == nil || !c.Supports(converter.FeatureRenderHooks) {
			return
		}

		h, err := ps.createRenderHooks(p.f)
		if err != nil {
			initErr = err
			return
		}
		p.renderHooks.hooks = h

		if !p.renderHooksHaveVariants || h.IsZero() {
			// Check if there is a different render hooks template
			// for any of the other page output formats.
			// If not, we can reuse this.
			for _, po := range ps.pageOutputs {
				if po.f.Name != p.f.Name {
					h2, err := ps.createRenderHooks(po.f)
					if err != nil {
						initErr = err
						return
					}

					if h2.IsZero() {
						continue
					}

					if p.renderHooks.hooks.IsZero() {
						p.renderHooks.hooks = h2
					}

					p.renderHooksHaveVariants = !h2.Eq(p.renderHooks.hooks)

					if p.renderHooksHaveVariants {
						break
					}

				}
			}
		}
	})

	return initErr
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

func (cp *pageContentOutput) renderContent(content []byte, renderTOC bool) (converter.Result, error) {
	c := cp.p.getContentConverter()
	return cp.renderContentWithConverter(c, content, renderTOC)
}

func (cp *pageContentOutput) renderContentWithConverter(c converter.Converter, content []byte, renderTOC bool) (converter.Result, error) {
	r, err := c.Convert(
		converter.RenderContext{
			Src:         content,
			RenderTOC:   renderTOC,
			RenderHooks: cp.renderHooks.hooks,
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

func (p *pageContentOutput) enableReuse() {
	p.reuseInit.Do(func() {
		p.reuse = true
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

func executeToString(h tpl.TemplateHandler, templ tpl.Template, data interface{}) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := h.Execute(templ, b, data); err != nil {
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
