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
	"bytes"
	"fmt"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
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
	return commonConvert(p, t)
}

type htmlHandler struct {
	basicPageHandler
}

func (h htmlHandler) Extensions() []string { return []string{"html", "htm"} }
func (h htmlHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	if p.rendered {
		panic(fmt.Sprintf("Page %q already rendered, does not need conversion", p.BaseFileName()))
	}

	// Work on a copy of the raw content from now on.
	p.createWorkContentCopy()

	p.ProcessShortcodes(t)

	return HandledResult{err: nil}
}

type asciidocHandler struct {
	basicPageHandler
}

func (h asciidocHandler) Extensions() []string { return []string{"asciidoc", "adoc", "ad"} }
func (h asciidocHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	return commonConvert(p, t)
}

type rstHandler struct {
	basicPageHandler
}

func (h rstHandler) Extensions() []string { return []string{"rest", "rst"} }
func (h rstHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	return commonConvert(p, t)
}

type mmarkHandler struct {
	basicPageHandler
}

func (h mmarkHandler) Extensions() []string { return []string{"mmark"} }
func (h mmarkHandler) PageConvert(p *Page, t tpl.Template) HandledResult {
	return commonConvert(p, t)
}

func commonConvert(p *Page, t tpl.Template) HandledResult {
	if p.rendered {
		panic(fmt.Sprintf("Page %q already rendered, does not need conversion", p.BaseFileName()))
	}

	// Work on a copy of the raw content from now on.
	p.createWorkContentCopy()

	p.ProcessShortcodes(t)

	// TODO(bep) these page handlers need to be re-evaluated, as it is hard to
	// process a page in isolation. See the new preRender func.
	if viper.GetBool("enableEmoji") {
		p.workContent = helpers.Emojify(p.workContent)
	}

	// We have to replace the <!--more--> with something that survives all the
	// rendering engines.
	// TODO(bep) inline replace
	p.workContent = bytes.Replace(p.workContent, []byte(helpers.SummaryDivider), internalSummaryDivider, 1)
	p.workContent = p.renderContent(p.workContent)

	return HandledResult{err: nil}
}
