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
	"unicode/utf8"

	"github.com/gohugoio/hugo/lazy"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/page"
)

var (
	nopTargetPath    = targetPathString("")
	nopPagePerOutput = struct {
		page.ContentProvider
		page.PageRenderProvider
		page.AlternativeOutputFormatsProvider
		targetPather
	}{
		page.NopPage,
		page.NopPage,
		page.NopPage,
		nopTargetPath,
	}
)

func newPageContentProvider(p *pageState) func(f output.Format) (*pageContentProvider, error) {

	parent := p.lateInit

	return func(f output.Format) (*pageContentProvider, error) {
		cp := &pageContentProvider{
			p: p,
			f: f,
		}

		init := parent.Branch(func() (interface{}, error) {
			// TODO(bep) page fast path when no shortcodes
			// Each page output format will get its own copy, if needed.
			if !p.renderable {
				// No markdown or similar to render, but it may still
				// contain shortcodes.
				cp.workContent = make([]byte, len(p.workContent))
				copy(cp.workContent, p.workContent)
				return nil, nil
			}

			cp.workContent = cp.renderContent(p, p.workContent)
			tmpContent, tmpTableOfContents := helpers.ExtractTOC(cp.workContent)
			cp.tableOfContents = helpers.BytesToHTML(tmpTableOfContents)
			cp.workContent = tmpContent

			return nil, nil
		})

		renderedContent := init.BranchdWithTimeout(p.s.siteConfigHolder.timeout, func(ctx context.Context) (interface{}, error) {

			c, err := cp.handleShortcodes(p, f, cp.workContent)
			if err != nil {
				return nil, err
			}

			if cp.p.source.hasSummaryDivider && cp.p.m.markup != "html" {
				summary, content, err := splitUserDefinedSummaryAndContent(cp.p.m.markup, c)
				if err != nil {
					cp.p.s.Log.ERROR.Printf("Failed to set user defined summary for page %q: %s", cp.p.pathOrTitle(), err)
				} else {
					c = content
					cp.summary = helpers.BytesToHTML(summary)
				}
			}

			cp.content = helpers.BytesToHTML(c)

			if !p.renderable {
				err := cp.addSelfTemplate()
				return nil, err
			}

			return nil, nil
		})

		plainInit := renderedContent.Branch(func() (interface{}, error) {
			cp.plain = helpers.StripHTML(string(cp.content))
			cp.plainWords = strings.Fields(cp.plain)
			cp.setWordCounts(p.m.isCJKLanguage)

			if err := cp.setAutoSummary(); err != nil {
				return err, nil
			}

			return nil, nil
		})

		cp.mainInit = renderedContent
		cp.plainInit = plainInit

		return cp, nil

	}

	// TODO(bep) page consider/remove page shifter logic

}

type pageContentProvider struct {
	f output.Format

	p *pageState

	// Lazy load dependencies
	mainInit  *lazy.Init
	plainInit *lazy.Init

	// Content state

	workContent []byte

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

// TODO(bep) page vs format
func (p *pageContentProvider) addSelfTemplate() error {
	self := p.p.selfLayoutForOutput(p.f)
	err := p.p.s.TemplateHandler().AddLateTemplate(self, string(p.content))
	if err != nil {
		return err
	}
	return nil
}

func (p *pageContentProvider) AlternativeOutputFormats() page.OutputFormats {
	var o page.OutputFormats
	for _, of := range p.p.OutputFormats() {
		if of.Format.NotAlternative || of.Format.Name == p.f.Name {
			continue
		}

		o = append(o, of)
	}
	return o
}

func (p *pageContentProvider) Content() (interface{}, error) {
	p.p.s.initInit(p.mainInit)
	return p.content, nil
}

func (p *pageContentProvider) FuzzyWordCount() int {
	p.p.s.initInit(p.plainInit)
	return p.fuzzyWordCount
}

func (p *pageContentProvider) Len() int {
	p.p.s.initInit(p.mainInit)
	return len(p.content)
}

func (p *pageContentProvider) Plain() string {
	p.p.s.initInit(p.plainInit)
	return p.plain
}

func (p *pageContentProvider) PlainWords() []string {
	p.p.s.initInit(p.plainInit)
	return p.plainWords
}

func (p *pageContentProvider) ReadingTime() int {
	p.p.s.initInit(p.plainInit)
	return p.readingTime
}

func (p *pageContentProvider) Render(layout ...string) template.HTML {
	l, err := p.p.getLayouts(p.f, layout...)
	if err != nil {
		p.p.s.DistinctErrorLog.Printf(".Render: Failed to resolve layout %q for page %q", layout, p.p.Path())
		return ""
	}

	for _, layout := range l {
		templ, found := p.p.s.Tmpl.Lookup(layout)
		if !found {
			// This is legacy from when we had only one output format and
			// HTML templates only. Some have references to layouts without suffix.
			// We default to good old HTML.
			templ, found = p.p.s.Tmpl.Lookup(layout + ".html")
		}
		if templ != nil {
			// Note that we're passing the full Page to the template.
			res, err := executeToString(templ, p.p)
			if err != nil {
				p.p.s.DistinctErrorLog.Printf(".Render: Failed to execute template %q: %s", layout, err)
				return template.HTML("")
			}
			return template.HTML(res)
		}
	}

	return ""

}

func (p *pageContentProvider) Summary() template.HTML {
	p.p.s.initInit(p.mainInit)
	if !p.p.source.hasSummaryDivider {
		p.p.s.initInit(p.plainInit)
	}
	return p.summary
}

func (p *pageContentProvider) TableOfContents() template.HTML {
	p.p.s.initInit(p.mainInit)
	return p.tableOfContents
}

func (p *pageContentProvider) Truncated() bool {
	return p.p.truncated || p.truncated
}

func (p *pageContentProvider) WordCount() int {
	p.p.s.initInit(p.plainInit)
	return p.wordCount
}

func (cp *pageContentProvider) handleShortcodes(p *pageState, f output.Format, rawContentCopy []byte) ([]byte, error) {
	if p.shortcodeState.getContentShortcodes().Len() == 0 {
		return rawContentCopy, nil
	}

	rendered, err := p.shortcodeState.executeShortcodesForOuputFormat(p, f)
	if err != nil {
		return rawContentCopy, err
	}

	rawContentCopy, err = replaceShortcodeTokens(rawContentCopy, shortcodePlaceholderPrefix, rendered)
	if err != nil {
		return nil, err
	}

	return rawContentCopy, nil
}

func (cp *pageContentProvider) renderContent(p page.Page, content []byte) []byte {
	return cp.p.s.ContentSpec.RenderBytes(&helpers.RenderingContext{
		Content: content, RenderTOC: true, PageFmt: cp.p.m.markup,
		Cfg:        p.Language(),
		DocumentID: p.File().UniqueID(), DocumentName: p.File().Path(),
		Config: cp.p.m.getRenderingConfig()})
}

func (p *pageContentProvider) setAutoSummary() error {
	if p.p.source.hasSummaryDivider {
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

func (p *pageContentProvider) setWordCounts(isCJKLanguage bool) {
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

// these will be shifted out when rendering a given output format.
type pagePerOutputProviders interface {
	targetPather
	page.ContentProvider
	page.PageRenderProvider
	page.AlternativeOutputFormatsProvider
}

type targetPathString string

func (s targetPathString) targetPath() string {
	return string(s)
}

type targetPather interface {
	targetPath() string
}

func executeToString(templ tpl.Template, data interface{}) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := templ.Execute(b, data); err != nil {
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
	case "asciidoc":
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
