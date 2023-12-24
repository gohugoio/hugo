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
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/markup/internal/attributes"
)

// Markdown attributes used by the Chroma highlighter.
var chromaHighlightProcessingAttributes = map[string]bool{
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
	for k, v := range chromaHighlightProcessingAttributes {
		chromaHighlightProcessingAttributes[strings.ToLower(k)] = v
	}
}

func New(cfg Config) Highlighter {
	return chromaHighlighter{
		cfg: cfg,
	}
}

type Highlighter interface {
	Highlight(code, lang string, opts any) (string, error)
	HighlightCodeBlock(ctx hooks.CodeblockContext, opts any) (HighlightResult, error)
	hooks.CodeBlockRenderer
	hooks.IsDefaultCodeBlockRendererProvider
}

type chromaHighlighter struct {
	cfg Config
}

func (h chromaHighlighter) Highlight(code, lang string, opts any) (string, error) {
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

func (h chromaHighlighter) HighlightCodeBlock(ctx hooks.CodeblockContext, opts any) (HighlightResult, error) {
	cfg := h.cfg

	var b strings.Builder

	attributes := ctx.(hooks.AttributesOptionsSliceProvider).AttributesSlice()

	options := ctx.Options()

	if err := applyOptionsFromMap(options, &cfg); err != nil {
		return HighlightResult{}, err
	}

	// Apply these last so the user can override them.
	if err := applyOptions(opts, &cfg); err != nil {
		return HighlightResult{}, err
	}

	if err := applyOptionsFromCodeBlockContext(ctx, &cfg); err != nil {
		return HighlightResult{}, err
	}

	low, high, err := highlight(&b, ctx.Inner(), ctx.Type(), attributes, cfg)
	if err != nil {
		return HighlightResult{}, err
	}

	highlighted := b.String()
	if high == 0 {
		high = len(highlighted)
	}

	return HighlightResult{
		highlighted: template.HTML(highlighted),
		innerLow:    low,
		innerHigh:   high,
	}, nil
}

func (h chromaHighlighter) RenderCodeblock(cctx context.Context, w hugio.FlexiWriter, ctx hooks.CodeblockContext) error {
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

// HighlightResult holds the result of an highlighting operation.
type HighlightResult struct {
	innerLow    int
	innerHigh   int
	highlighted template.HTML
}

// Wrapped returns the highlighted code wrapped in a <div>, <pre> and <code> tag.
func (h HighlightResult) Wrapped() template.HTML {
	return h.highlighted
}

// Inner returns the highlighted code without the wrapping <div>, <pre> and <code> tag, suitable for inline use.
func (h HighlightResult) Inner() template.HTML {
	return h.highlighted[h.innerLow:h.innerHigh]
}

func highlight(fw hugio.FlexiWriter, code, lang string, attributes []attributes.Attribute, cfg Config) (int, int, error) {
	var lexer chroma.Lexer
	if lang != "" {
		lexer = chromalexers.Get(lang)
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
		if cfg.Hl_inline {
			fmt.Fprintf(w, "<code%s>%s</code>", inlineCodeAttrs(lang), gohtml.EscapeString(code))
		} else {
			preWrapper := getPreWrapper(lang, w)
			fmt.Fprint(w, preWrapper.Start(true, ""))
			fmt.Fprint(w, gohtml.EscapeString(code))
			fmt.Fprint(w, preWrapper.End(true))
		}
		return 0, 0, nil
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

	if !cfg.Hl_inline {
		writeDivStart(w, attributes)
	}

	options := cfg.toHTMLOptions()
	var wrapper html.PreWrapper

	if cfg.Hl_inline {
		wrapper = startEnd{
			start: func(code bool, styleAttr string) string {
				if code {
					return fmt.Sprintf(`<code%s>`, inlineCodeAttrs(lang))
				}
				return ``
			},
			end: func(code bool) string {
				if code {
					return `</code>`
				}

				return ``
			},
		}
	} else {
		wrapper = getPreWrapper(lang, w)
	}

	options = append(options, html.WithPreWrapper(wrapper))

	formatter := html.New(options...)

	if err := formatter.Format(w, style, iterator); err != nil {
		return 0, 0, err
	}

	if !cfg.Hl_inline {
		writeDivEnd(w)
	}

	if p, ok := wrapper.(*preWrapper); ok {
		return p.low, p.high, nil
	}

	return 0, 0, nil
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

func inlineCodeAttrs(lang string) string {
	return fmt.Sprintf(` class="code-inline language-%s"`, lang)
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

type startEnd struct {
	start func(code bool, styleAttr string) string
	end   func(code bool) string
}

func (s startEnd) Start(code bool, styleAttr string) string {
	return s.start(code, styleAttr)
}

func (s startEnd) End(code bool) string {
	return s.end(code)
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
