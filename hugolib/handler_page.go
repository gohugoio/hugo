// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"fmt"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"
)

func init() {
	RegisterHandler(new(markdownHandler))
	RegisterHandler(new(htmlHandler))
	RegisterHandler(new(asciidocHandler))
	RegisterHandler(new(rstHandler))
	RegisterHandler(new(mmarkHandler))
	RegisterHandler(new(orgHandler))
}

type basicPageHandler Handle

func (b basicPageHandler) Read(f *source.File, s *Site) HandledResult {
	page, err := s.NewPage(f.Path())

	if err != nil {
		return HandledResult{file: f, err: err}
	}

	if _, err := page.ReadFrom(f.Contents); err != nil {
		return HandledResult{file: f, err: err}
	}

	// In a multilanguage setup, we use the first site to
	// do the initial processing.
	// That site may be different than where the page will end up,
	// so we do the assignment here.
	// We should clean up this, but that will have to wait.
	s.assignSiteByLanguage(page)

	return HandledResult{file: f, page: page, err: err}
}

func (b basicPageHandler) FileConvert(*source.File, *Site) HandledResult {
	return HandledResult{}
}

type markdownHandler struct {
	basicPageHandler
}

func (h markdownHandler) Extensions() []string { return []string{"mdown", "markdown", "md"} }
func (h markdownHandler) PageConvert(p *Page) HandledResult {
	return commonConvert(p)
}

type htmlHandler struct {
	basicPageHandler
}

func (h htmlHandler) Extensions() []string { return []string{"html", "htm"} }

func (h htmlHandler) PageConvert(p *Page) HandledResult {
	if p.rendered {
		panic(fmt.Sprintf("Page %q already rendered, does not need conversion", p.BaseFileName()))
	}

	// Work on a copy of the raw content from now on.
	p.createWorkContentCopy()

	if err := p.processShortcodes(); err != nil {
		p.s.Log.ERROR.Println(err)
	}

	return HandledResult{err: nil}
}

type asciidocHandler struct {
	basicPageHandler
}

func (h asciidocHandler) Extensions() []string { return []string{"asciidoc", "adoc", "ad"} }
func (h asciidocHandler) PageConvert(p *Page) HandledResult {
	return commonConvert(p)
}

type rstHandler struct {
	basicPageHandler
}

func (h rstHandler) Extensions() []string { return []string{"rest", "rst"} }
func (h rstHandler) PageConvert(p *Page) HandledResult {
	return commonConvert(p)
}

type mmarkHandler struct {
	basicPageHandler
}

func (h mmarkHandler) Extensions() []string { return []string{"mmark"} }
func (h mmarkHandler) PageConvert(p *Page) HandledResult {
	return commonConvert(p)
}

type orgHandler struct {
	basicPageHandler
}

func (h orgHandler) Extensions() []string { return []string{"org"} }
func (h orgHandler) PageConvert(p *Page) HandledResult {
	return commonConvert(p)
}

func commonConvert(p *Page) HandledResult {
	if p.rendered {
		panic(fmt.Sprintf("Page %q already rendered, does not need conversion", p.BaseFileName()))
	}

	// Work on a copy of the raw content from now on.
	p.createWorkContentCopy()

	if err := p.processShortcodes(); err != nil {
		p.s.Log.ERROR.Println(err)
	}

	// TODO(bep) these page handlers need to be re-evaluated, as it is hard to
	// process a page in isolation. See the new preRender func.
	if p.s.Cfg.GetBool("enableEmoji") {
		p.workContent = helpers.Emojify(p.workContent)
	}

	p.workContent = p.replaceDivider(p.workContent)
	p.workContent = p.renderContent(p.workContent)

	return HandledResult{err: nil}
}
