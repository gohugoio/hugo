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
	"context"
	"fmt"
	gohtml "html"
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/internal/attributes"
)

// Markdown attributes used by the Chroma hightlighter.
var chromaHightlightProcessingAttributes = map[string]bool{
	"anchorLineNos":      true,
	"guessSyntax":        true,
	"hl_Lines":           true,
	"lineAnchors":        true,
	"lineNos":            true,
	"lineNoStart":        true,
	"lineNumbersInTable": true,
	"noClasses":          true,
	"style":              true,
	"tabWidth":           true,
}

func init() {
	for k, v := range chromaHightlightProcessingAttributes {
		chromaHightlightProcessingAttributes[strings.ToLower(k)] = v
	}
}

func New(cfg Config) Highlighter {
	return chromaHighlighter{
		cfg: cfg,
	}
}

type Highlighter interface {
	Highlight(code, lang string, opts interface{}) (string, error)
	HighlightCodeBlock(ctx hooks.CodeblockContext, opts interface{}) (HightlightResult, error)
	hooks.CodeBlockRenderer
}

type chromaHighlighter struct {
	cfg Config
}

func (h chromaHighlighter) Highlight(code, lang string, opts interface{}) (string, error) {
	cfg := h.cfg
	if err := applyOptions(opts, &cfg); err != nil {
		return "", err
	}
	var b strings.Builder

	if err := highlight(&b, code, lang, nil, cfg); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (h chromaHighlighter) HighlightCodeBlock(ctx hooks.CodeblockContext, opts interface{}) (HightlightResult, error) {
	cfg := h.cfg

	var b strings.Builder

	attributes := ctx.(hooks.AttributesOptionsSliceProvider).AttributesSlice()
	options := ctx.Options()

	if err := applyOptionsFromMap(options, &cfg); err != nil {
		return HightlightResult{}, err
	}

	// Apply these last so the user can override them.
	if err := applyOptions(opts, &cfg); err != nil {
		return HightlightResult{}, err
	}

	err := highlight(&b, ctx.Code(), ctx.Lang(), attributes, cfg)
	if err != nil {
		return HightlightResult{}, err
	}

	return HightlightResult{
		Body: template.HTML(b.String()),
	}, nil
}

func (h chromaHighlighter) RenderCodeblock(ctx context.Context, w hugio.FlexiWriter, ctxt hooks.CodeblockContext) error {
	cfg := h.cfg
	attributes := ctxt.(hooks.AttributesOptionsSliceProvider).AttributesSlice()

	if err := applyOptionsFromMap(ctxt.Options(), &cfg); err != nil {
		return err
	}

	return highlight(w, ctxt.Code(), ctxt.Lang(), attributes, cfg)
}

var id = identity.NewPathIdentity("chroma", "highlight")

func (h chromaHighlighter) GetIdentity() identity.Identity {
	return id
}

type HightlightResult struct {
	Body template.HTML
}

func (h HightlightResult) Highlighted() template.HTML {
	return h.Body
}

func (h chromaHighlighter) toHighlightOptionsAttributes(ctx hooks.CodeblockContext) (map[string]interface{}, map[string]interface{}) {
	attributes := ctx.Attributes()
	if attributes == nil || len(attributes) == 0 {
		return nil, nil
	}

	options := make(map[string]interface{})
	attrs := make(map[string]interface{})

	for k, v := range attributes {
		klow := strings.ToLower(k)
		if chromaHightlightProcessingAttributes[klow] {
			options[klow] = v
		} else {
			attrs[k] = v
		}
	}
	const lineanchorsKey = "lineanchors"
	if _, found := options[lineanchorsKey]; !found {
		// Set it to the ordinal.
		options[lineanchorsKey] = strconv.Itoa(ctx.Ordinal())
	}
	return options, attrs
}

func highlight(w hugio.FlexiWriter, code, lang string, attributes []attributes.Attribute, cfg Config) error {
	var lexer chroma.Lexer
	if lang != "" {
		lexer = lexers.Get(lang)
	}

	if lexer == nil && (cfg.GuessSyntax && !cfg.NoHl) {
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
		return nil
	}

	style := styles.Get(cfg.Style)
	if style == nil {
		style = styles.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return err
	}

	options := cfg.ToHTMLOptions()
	options = append(options, getHtmlPreWrapper(lang))

	formatter := html.New(options...)

	writeDivStart(w, attributes)
	if err := formatter.Format(w, style, iterator); err != nil {
		return err
	}
	writeDivEnd(w)

	return nil
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
	var language string
	if code {
		language = p.language
	}
	w := &strings.Builder{}
	WritePreStart(w, language, styleAttr)
	return w.String()
}

func WritePreStart(w io.Writer, language, styleAttr string) {
	fmt.Fprintf(w, `<pre tabindex="0"%s>`, styleAttr)
	fmt.Fprint(w, "<code")
	if language != "" {
		fmt.Fprint(w, ` class="language-`+language+`"`)
		fmt.Fprint(w, ` data-lang="`+language+`"`)
	}
	fmt.Fprint(w, ">")
}

const preEnd = "</code></pre>"

func (p preWrapper) End(code bool) string {
	return preEnd
}

func WritePreEnd(w io.Writer) {
	fmt.Fprint(w, preEnd)
}

func writeDivStart(w hugio.FlexiWriter, attrs []attributes.Attribute) {
	w.WriteString(`<div class="highlight`)
	if attrs != nil {
		for _, attr := range attrs {
			if attr.Name == "class" {
				w.WriteString(" " + attr.ValueString())
				break
			}
		}
		_, _ = w.WriteString("\"")
		attributes.RenderAttributes(w, true, attrs...)
	} else {
		_, _ = w.WriteString("\"")
	}

	w.WriteString(">")
}

func writeDivEnd(w hugio.FlexiWriter) {
	w.WriteString("</div>")
}
