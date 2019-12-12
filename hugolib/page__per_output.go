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

func newPageContentOutput(p *pageState) func(f output.Format) (*pageContentOutput, error) {

	parent := p.init

	return func(f output.Format) (*pageContentOutput, error) {
		cp := &pageContentOutput{
			p: p,
			f: f,
		}

		initContent := func() (err error) {
			if p.cmap == nil {
				// Nothing to do.
				return nil
			}
			defer func() {
				// See https://github.com/gohugoio/hugo/issues/6210
				if r := recover(); r != nil {
					err = fmt.Errorf("%s", r)
					p.s.Log.ERROR.Printf("[BUG] Got panic:\n%s\n%s", r, string(debug.Stack()))
				}
			}()

			var hasVariants bool

			cp.contentPlaceholders, hasVariants, err = p.shortcodeState.renderShortcodesForPage(p, f)
			if err != nil {
				return err
			}

			if p.render && !hasVariants {
				// We can reuse this for the other output formats
				cp.enableReuse()
			}

			cp.workContent = p.contentToRender(cp.contentPlaceholders)

			isHTML := cp.p.m.markup == "html"

			if p.renderable {
				if !isHTML {
					r, err := cp.renderContent(cp.workContent)
					if err != nil {
						return err
					}
					cp.convertedResult = r
					cp.workContent = r.Bytes()

					if _, ok := r.(converter.TableOfContentsProvider); !ok {
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
							cp.p.s.Log.ERROR.Printf("Failed to set user defined summary for page %q: %s", cp.p.pathOrTitle(), err)
						} else {
							cp.workContent = content
							cp.summary = helpers.BytesToHTML(summary)
						}
					}
				} else if cp.p.m.summary != "" {
					b, err := cp.p.getContentConverter().Convert(
						converter.RenderContext{
							Src: []byte(cp.p.m.summary),
						},
					)

					if err != nil {
						return err
					}
					html := cp.p.s.ContentSpec.TrimShortHTML(b.Bytes())
					cp.summary = helpers.BytesToHTML(html)
				}
			}

			cp.content = helpers.BytesToHTML(cp.workContent)

			if !p.renderable {
				err := cp.addSelfTemplate()
				return err
			}

			return nil

		}

		// Recursive loops can only happen in content files with template code (shortcodes etc.)
		// Avoid creating new goroutines if we don't have to.
		needTimeout := !p.renderable || p.shortcodeState.hasShortcodes()

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

}

// pageContentOutput represents the Page content for a given output format.
type pageContentOutput struct {
	f output.Format

	// If we can safely reuse this for other output formats.
	reuse     bool
	reuseInit sync.Once

	p *pageState

	// Lazy load dependencies
	initMain  *lazy.Init
	initPlain *lazy.Init

	placeholdersEnabled     bool
	placeholdersEnabledInit sync.Once

	// Content state

	workContent     []byte
	convertedResult converter.Result

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
	if tocProvider, ok := p.convertedResult.(converter.TableOfContentsProvider); ok {
		cfg := p.p.s.ContentSpec.Converters.GetMarkupConfig()
		return template.HTML(tocProvider.TableOfContents().ToHTML(cfg.TableOfContents.StartLevel, cfg.TableOfContents.EndLevel, cfg.TableOfContents.Ordered))
	}
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

func (cp *pageContentOutput) renderContent(content []byte) (converter.Result, error) {
	return cp.p.getContentConverter().Convert(
		converter.RenderContext{
			Src:       content,
			RenderTOC: true,
		})
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

func (p *pageContentOutput) addSelfTemplate() error {
	self := p.p.selfLayoutForOutput(p.f)
	err := p.p.s.TemplateHandler().AddLateTemplate(self, string(p.content))
	if err != nil {
		return err
	}
	return nil
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
	page.ContentProvider
	page.PaginatorProvider
	page.TableOfContentsProvider
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
