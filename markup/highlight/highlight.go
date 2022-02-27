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
	"html/template"
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/text"
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
	hooks.IsDefaultCodeBlockRendererProvider
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

	if _, _, err := highlight(&b, code, lang, nil, cfg); err != nil {
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

	if err := applyOptionsFromCodeBlockContext(ctx, &cfg); err != nil {
		return HightlightResult{}, err
	}

	low, high, err := highlight(&b, ctx.Inner(), ctx.Type(), attributes, cfg)
	if err != nil {
		return HightlightResult{}, err
	}

	return HightlightResult{
		highlighted: template.HTML(b.String()),
		innerLow:    low,
		innerHigh:   high,
	}, nil
}

func (h chromaHighlighter) RenderCodeblock(w hugio.FlexiWriter, ctx hooks.CodeblockContext) error {
	cfg := h.cfg
	attributes := ctx.(hooks.AttributesOptionsSliceProvider).AttributesSlice()

	if err := applyOptionsFromMap(ctx.Options(), &cfg); err != nil {
		return err
	}

	if err := applyOptionsFromCodeBlockContext(ctx, &cfg); err != nil {
		return err
	}

	code := text.Puts(ctx.Inner())

	_, _, err := highlight(w, code, ctx.Type(), attributes, cfg)
	return err
}

func (h chromaHighlighter) IsDefaultCodeBlockRenderer() bool {
	return true
}

var id = identity.NewPathIdentity("chroma", "highlight")

func (h chromaHighlighter) GetIdentity() identity.Identity {
	return id
}

type HightlightResult struct {
	innerLow    int
	innerHigh   int
	highlighted template.HTML
}

func (h HightlightResult) Wrapped() template.HTML {
	return h.highlighted
}

func (h HightlightResult) Inner() template.HTML {
	return h.highlighted[h.innerLow:h.innerHigh]
}

func highlight(fw hugio.FlexiWriter, code, lang string, attributes []attributes.Attribute, cfg Config) (int, int, error) {
	var low, high int

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

	w := &byteCountFlexiWriter{delegate: fw}

	if lexer == nil {
		wrapper := getPreWrapper(lang, w)
		fmt.Fprint(w, wrapper.Start(true, ""))
		fmt.Fprint(w, gohtml.EscapeString(code))
		fmt.Fprint(w, wrapper.End(true))
		return low, high, nil
	}

	style := styles.Get(cfg.Style)
	if style == nil {
		style = styles.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return 0, 0, err
	}

	options := cfg.ToHTMLOptions()
	preWrapper := getPreWrapper(lang, w)
	options = append(options, html.WithPreWrapper(preWrapper))

	formatter := html.New(options...)

	writeDivStart(w, attributes)

	if err := formatter.Format(w, style, iterator); err != nil {
		return 0, 0, err
	}
	writeDivEnd(w)

	return preWrapper.low, preWrapper.high, nil
}

func getPreWrapper(language string, writeCounter *byteCountFlexiWriter) *preWrapper {
	return &preWrapper{language: language, writeCounter: writeCounter}
}

type preWrapper struct {
	low          int
	high         int
	writeCounter *byteCountFlexiWriter
	language     string
}

func (p *preWrapper) Start(code bool, styleAttr string) string {
	var language string
	if code {
		language = p.language
	}
	w := &strings.Builder{}
	WritePreStart(w, language, styleAttr)
	p.low = p.writeCounter.counter + w.Len()
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

func (p *preWrapper) End(code bool) string {
	p.high = p.writeCounter.counter
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

type byteCountFlexiWriter struct {
	delegate hugio.FlexiWriter
	counter  int
}

func (w *byteCountFlexiWriter) Write(p []byte) (int, error) {
	n, err := w.delegate.Write(p)
	w.counter += n
	return n, err
}

func (w *byteCountFlexiWriter) WriteByte(c byte) error {
	w.counter++
	return w.delegate.WriteByte(c)
}

func (w *byteCountFlexiWriter) WriteString(s string) (int, error) {
	n, err := w.delegate.WriteString(s)
	w.counter += n
	return n, err
}

func (w *byteCountFlexiWriter) WriteRune(r rune) (int, error) {
	n, err := w.delegate.WriteRune(r)
	w.counter += n
	return n, err
}
