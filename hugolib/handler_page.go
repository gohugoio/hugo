// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
)

func init() {
	RegisterHandler(new(markdownHandler))
	RegisterHandler(new(htmlHandler))
	RegisterHandler(new(asciidocHandler))
	RegisterHandler(new(rstHandler))
	RegisterHandler(new(mmarkHandler))
}

type basicPageHandler Handle

func (b basicPageHandler) Read(f *source.File, s *Site) HandledResult {
	page, err := NewPage(f.Path())
	if err != nil {
		return HandledResult{file: f, err: err}
	}

	if _, err := page.ReadFrom(f.Contents); err != nil {
		return HandledResult{file: f, err: err}
	}

	page.Site = &s.Info
	page.Tmpl = s.Tmpl

	return HandledResult{file: f, page: page, err: err}
}

func (b basicPageHandler) FileConvert(*source.File, *Site) HandledResult {
	return HandledResult{}
}

type markdownHandler struct {
	basicPageHandler
}

func (h markdownHandler) Extensions() []string { return []string{"mdown", "markdown", "md"} }
func (h markdownHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	p.ProcessShortcodes(t)

	tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.renderContent(helpers.RemoveSummaryDivider(p.rawContent)))

	if len(p.contentShortCodes) > 0 {
		replaced, err := replaceShortcodeTokensInsources(shortcodePlaceholderPrefix, true, p.contentShortCodes,
			tmpContent, tmpTableOfContents)
		if err != nil {
			jww.FATAL.Printf("Fail to replace shortcode tokens in %s:\n%s", p.BaseFileName(), err.Error())
			return HandledResult{err: err}
		}
		tmpContent = replaced[0]
		tmpTableOfContents = replaced[1]
	}

	p.Content = helpers.BytesToHTML(tmpContent)
	p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)

	return HandledResult{err: nil}
}

type htmlHandler struct {
	basicPageHandler
}

func (h htmlHandler) Extensions() []string { return []string{"html", "htm"} }
func (h htmlHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	p.ProcessShortcodes(t)
	var content []byte
	var err error

	if len(p.contentShortCodes) > 0 {
		content, err = replaceShortcodeTokens(p.rawContent, shortcodePlaceholderPrefix, true, p.contentShortCodes)

		if err != nil {
			jww.FATAL.Printf("Fail to replace shortcode tokens in %s:\n%s", p.BaseFileName(), err.Error())
			return HandledResult{err: err}
		}
	} else {
		content = p.rawContent
	}

	p.Content = helpers.BytesToHTML(content)
	return HandledResult{err: nil}
}

type asciidocHandler struct {
	basicPageHandler
}

func (h asciidocHandler) Extensions() []string { return []string{"asciidoc", "adoc", "ad"} }
func (h asciidocHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	p.ProcessShortcodes(t)

	// TODO(spf13) Add/Refactor AsciiDoc Logic here
	tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.renderContent(helpers.RemoveSummaryDivider(p.rawContent)))

	if len(p.contentShortCodes) > 0 {
		replaced, err := replaceShortcodeTokensInsources(shortcodePlaceholderPrefix, true, p.contentShortCodes,
			tmpContent, tmpTableOfContents)
		if err != nil {
			jww.FATAL.Printf("Fail to replace shortcode tokens in %s:\n%s", p.BaseFileName(), err.Error())
			return HandledResult{err: err}
		}
		tmpContent = replaced[0]
		tmpTableOfContents = replaced[1]
	}

	p.Content = helpers.BytesToHTML(tmpContent)
	p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)

	//err := p.Convert()
	return HandledResult{page: p, err: nil}
}

type rstHandler struct {
	basicPageHandler
}

func (h rstHandler) Extensions() []string { return []string{"rest", "rst"} }
func (h rstHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	p.ProcessShortcodes(t)

	tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.renderContent(helpers.RemoveSummaryDivider(p.rawContent)))

	if len(p.contentShortCodes) > 0 {
		replaced, err := replaceShortcodeTokensInsources(shortcodePlaceholderPrefix, true, p.contentShortCodes,
			tmpContent, tmpTableOfContents)
		if err != nil {
			jww.FATAL.Printf("Fail to replace shortcode tokens in %s:\n%s", p.BaseFileName(), err.Error())
			return HandledResult{err: err}
		}
		tmpContent = replaced[0]
		tmpTableOfContents = replaced[1]
	}

	p.Content = helpers.BytesToHTML(tmpContent)
	p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)

	return HandledResult{err: nil}
}

type mmarkHandler struct {
	basicPageHandler
}

func (h mmarkHandler) Extensions() []string { return []string{"mmark"} }
func (h mmarkHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	p.ProcessShortcodes(t)

	tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.renderContent(helpers.RemoveSummaryDivider(p.rawContent)))

	if len(p.contentShortCodes) > 0 {
		tmpContentWithTokensReplaced, err := replaceShortcodeTokens(tmpContent, shortcodePlaceholderPrefix, true, p.contentShortCodes)

		if err != nil {
			jww.FATAL.Printf("Fail to replace short code tokens in %s:\n%s", p.BaseFileName(), err.Error())
			return HandledResult{err: err}
		} else {
			tmpContent = tmpContentWithTokensReplaced
		}
	}

	p.Content = helpers.BytesToHTML(tmpContent)
	p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)

	return HandledResult{err: nil}
}
