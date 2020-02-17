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

package highlight

import (
	"fmt"
	gohtml "html"
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	hl "github.com/yuin/goldmark-highlighting"
)

func New(cfg Config) Highlighter {
	return Highlighter{
		cfg: cfg,
	}
}

type Highlighter struct {
	cfg Config
}

func (h Highlighter) Highlight(code, lang, optsStr string) (string, error) {
	cfg := h.cfg
	if optsStr != "" {
		if err := applyOptionsFromString(optsStr, &cfg); err != nil {
			return "", err
		}
	}
	return highlight(code, lang, cfg)
}

func highlight(code, lang string, cfg Config) (string, error) {
	w := &strings.Builder{}
	var lexer chroma.Lexer
	if lang != "" {
		lexer = lexers.Get(lang)
	}

	if lexer == nil && cfg.GuessSyntax {
		lexer = lexers.Analyse(code)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lang = strings.ToLower(lexer.Config().Name)
	}

	if lexer == nil {
		wrapper := getPreWrapper(lang)
		fmt.Fprint(w, wrapper.Start(true, ""))
		fmt.Fprint(w, gohtml.EscapeString(code))
		fmt.Fprint(w, wrapper.End(true))
		return w.String(), nil
	}

	style := styles.Get(cfg.Style)
	if style == nil {
		style = styles.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	options := cfg.ToHTMLOptions()
	options = append(options, getHtmlPreWrapper(lang))

	formatter := html.New(options...)

	fmt.Fprint(w, `<div class="highlight">`)
	if err := formatter.Format(w, style, iterator); err != nil {
		return "", err
	}
	fmt.Fprint(w, `</div>`)

	return w.String(), nil
}

func GetCodeBlockOptions() func(ctx hl.CodeBlockContext) []html.Option {
	return func(ctx hl.CodeBlockContext) []html.Option {
		var language string
		if l, ok := ctx.Language(); ok {
			language = string(l)
		}
		return []html.Option{
			getHtmlPreWrapper(language),
		}
	}
}

func getPreWrapper(language string) preWrapper {
	return preWrapper{language: language}
}

func getHtmlPreWrapper(language string) html.Option {
	return html.WithPreWrapper(getPreWrapper(language))
}

type preWrapper struct {
	language string
}

func (p preWrapper) Start(code bool, styleAttr string) string {
	w := &strings.Builder{}
	fmt.Fprintf(w, "<pre%s>", styleAttr)
	var language string
	if code {
		language = p.language
	}
	WriteCodeTag(w, language)
	return w.String()
}

func WriteCodeTag(w io.Writer, language string) {
	fmt.Fprint(w, "<code")
	if language != "" {
		fmt.Fprint(w, ` class="language-`+language+`"`)
		fmt.Fprint(w, ` data-lang="`+language+`"`)
	}
	fmt.Fprint(w, ">")
}

func (p preWrapper) End(code bool) string {
	return "</code></pre>"
}
