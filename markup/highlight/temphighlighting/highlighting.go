// package highlighting is a extension for the goldmark(http://github.com/yuin/goldmark).
//
// This extension adds syntax-highlighting to the fenced code blocks using
// chroma(https://github.com/alecthomas/chroma).
//
// TODO(bep) this is a very temporary fork based on https://github.com/yuin/goldmark-highlighting/pull/10
// MIT Licensed, Copyright Yusuke Inuzuka
package temphighlighting

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// ImmutableAttributes is a read-only interface for ast.Attributes.
type ImmutableAttributes interface {
	// Get returns (value, true) if an attribute associated with given
	// name exists, otherwise (nil, false)
	Get(name []byte) (interface{}, bool)

	// GetString returns (value, true) if an attribute associated with given
	// name exists, otherwise (nil, false)
	GetString(name string) (interface{}, bool)

	// All returns all attributes.
	All() []ast.Attribute
}

type immutableAttributes struct {
	n ast.Node
}

func (a *immutableAttributes) Get(name []byte) (interface{}, bool) {
	return a.n.Attribute(name)
}

func (a *immutableAttributes) GetString(name string) (interface{}, bool) {
	return a.n.AttributeString(name)
}

func (a *immutableAttributes) All() []ast.Attribute {
	if a.n.Attributes() == nil {
		return []ast.Attribute{}
	}
	return a.n.Attributes()
}

// CodeBlockContext holds contextual information of code highlighting.
type CodeBlockContext interface {
	// Language returns (language, true) if specified, otherwise (nil, false).
	Language() ([]byte, bool)

	// Highlighted returns true if this code block can be highlighted, otherwise false.
	Highlighted() bool

	// Attributes return attributes of the code block.
	Attributes() ImmutableAttributes
}

type codeBlockContext struct {
	language    []byte
	highlighted bool
	attributes  ImmutableAttributes
}

func newCodeBlockContext(language []byte, highlighted bool, attrs ImmutableAttributes) CodeBlockContext {
	return &codeBlockContext{
		language:    language,
		highlighted: highlighted,
		attributes:  attrs,
	}
}

func (c *codeBlockContext) Language() ([]byte, bool) {
	if c.language != nil {
		return c.language, true
	}
	return nil, false
}

func (c *codeBlockContext) Highlighted() bool {
	return c.highlighted
}

func (c *codeBlockContext) Attributes() ImmutableAttributes {
	return c.attributes
}

// WrapperRenderer renders wrapper elements like div, pre, etc.
type WrapperRenderer func(w util.BufWriter, context CodeBlockContext, entering bool)

// CodeBlockOptions creates Chroma options per code block.
type CodeBlockOptions func(ctx CodeBlockContext) []chromahtml.Option

// Config struct holds options for the extension.
type Config struct {
	html.Config

	// Style is a highlighting style.
	// Supported styles are defined under https://github.com/alecthomas/chroma/tree/master/formatters.
	Style string

	// FormatOptions is a option related to output formats.
	// See https://github.com/alecthomas/chroma#the-html-formatter for details.
	FormatOptions []chromahtml.Option

	// CSSWriter is an io.Writer that will be used as CSS data output buffer.
	// If WithClasses() is enabled, you can get CSS data corresponds to the style.
	CSSWriter io.Writer

	// CodeBlockOptions allows set Chroma options per code block.
	CodeBlockOptions CodeBlockOptions

	// WrapperRendererCodeBlockOptions allows you to change wrapper elements.
	WrapperRenderer WrapperRenderer
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{
		Config:           html.NewConfig(),
		Style:            "github",
		FormatOptions:    []chromahtml.Option{},
		CSSWriter:        nil,
		WrapperRenderer:  nil,
		CodeBlockOptions: nil,
	}
}

// SetOption implements renderer.SetOptioner.
func (c *Config) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optStyle:
		c.Style = value.(string)
	case optFormatOptions:
		if value != nil {
			c.FormatOptions = value.([]chromahtml.Option)
		}
	case optCSSWriter:
		c.CSSWriter = value.(io.Writer)
	case optWrapperRenderer:
		c.WrapperRenderer = value.(WrapperRenderer)
	case optCodeBlockOptions:
		c.CodeBlockOptions = value.(CodeBlockOptions)
	default:
		c.Config.SetOption(name, value)
	}
}

// Option interface is a functional option interface for the extension.
type Option interface {
	renderer.Option
	// SetHighlightingOption sets given option to the extension.
	SetHighlightingOption(*Config)
}

type withHTMLOptions struct {
	value []html.Option
}

func (o *withHTMLOptions) SetConfig(c *renderer.Config) {
	if o.value != nil {
		for _, v := range o.value {
			v.(renderer.Option).SetConfig(c)
		}
	}
}

func (o *withHTMLOptions) SetHighlightingOption(c *Config) {
	if o.value != nil {
		for _, v := range o.value {
			v.SetHTMLOption(&c.Config)
		}
	}
}

// WithHTMLOptions is functional option that wraps goldmark HTMLRenderer options.
func WithHTMLOptions(opts ...html.Option) Option {
	return &withHTMLOptions{opts}
}

const optStyle renderer.OptionName = "HighlightingStyle"

var highlightLinesAttrName = []byte("hl_lines")

var styleAttrName = []byte("hl_style")
var nohlAttrName = []byte("nohl")
var linenosAttrName = []byte("linenos")
var linenosTableAttrValue = []byte("table")
var linenosInlineAttrValue = []byte("inline")
var linenostartAttrName = []byte("linenostart")

type withStyle struct {
	value string
}

func (o *withStyle) SetConfig(c *renderer.Config) {
	c.Options[optStyle] = o.value
}

func (o *withStyle) SetHighlightingOption(c *Config) {
	c.Style = o.value
}

// WithStyle is a functional option that changes highlighting style.
func WithStyle(style string) Option {
	return &withStyle{style}
}

const optCSSWriter renderer.OptionName = "HighlightingCSSWriter"

type withCSSWriter struct {
	value io.Writer
}

func (o *withCSSWriter) SetConfig(c *renderer.Config) {
	c.Options[optCSSWriter] = o.value
}

func (o *withCSSWriter) SetHighlightingOption(c *Config) {
	c.CSSWriter = o.value
}

// WithCSSWriter is a functional option that sets io.Writer for CSS data.
func WithCSSWriter(w io.Writer) Option {
	return &withCSSWriter{w}
}

const optWrapperRenderer renderer.OptionName = "HighlightingWrapperRenderer"

type withWrapperRenderer struct {
	value WrapperRenderer
}

func (o *withWrapperRenderer) SetConfig(c *renderer.Config) {
	c.Options[optWrapperRenderer] = o.value
}

func (o *withWrapperRenderer) SetHighlightingOption(c *Config) {
	c.WrapperRenderer = o.value
}

// WithWrapperRenderer is a functional option that sets WrapperRenderer that
// renders wrapper elements like div, pre, etc.
func WithWrapperRenderer(w WrapperRenderer) Option {
	return &withWrapperRenderer{w}
}

const optCodeBlockOptions renderer.OptionName = "HighlightingCodeBlockOptions"

type withCodeBlockOptions struct {
	value CodeBlockOptions
}

func (o *withCodeBlockOptions) SetConfig(c *renderer.Config) {
	c.Options[optWrapperRenderer] = o.value
}

func (o *withCodeBlockOptions) SetHighlightingOption(c *Config) {
	c.CodeBlockOptions = o.value
}

// WithCodeBlockOptions is a functional option that sets CodeBlockOptions that
// allows setting Chroma options per code block.
func WithCodeBlockOptions(c CodeBlockOptions) Option {
	return &withCodeBlockOptions{value: c}
}

const optFormatOptions renderer.OptionName = "HighlightingFormatOptions"

type withFormatOptions struct {
	value []chromahtml.Option
}

func (o *withFormatOptions) SetConfig(c *renderer.Config) {
	if _, ok := c.Options[optFormatOptions]; !ok {
		c.Options[optFormatOptions] = []chromahtml.Option{}
	}
	c.Options[optStyle] = append(c.Options[optFormatOptions].([]chromahtml.Option), o.value...)
}

func (o *withFormatOptions) SetHighlightingOption(c *Config) {
	c.FormatOptions = append(c.FormatOptions, o.value...)
}

// WithFormatOptions is a functional option that wraps chroma HTML formatter options.
func WithFormatOptions(opts ...chromahtml.Option) Option {
	return &withFormatOptions{opts}
}

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct {
	Config
}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer(opts ...Option) renderer.NodeRenderer {
	r := &HTMLRenderer{
		Config: NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHighlightingOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
}

func getAttributes(node *ast.FencedCodeBlock, infostr []byte) ImmutableAttributes {
	if node.Attributes() != nil {
		return &immutableAttributes{node}
	}
	if infostr != nil {
		attrStartIdx := -1

		for idx, char := range infostr {
			if char == '{' {
				attrStartIdx = idx
				break
			}
		}
		if attrStartIdx > 0 {
			n := ast.NewTextBlock() // dummy node for storing attributes
			attrStr := infostr[attrStartIdx:]
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				for _, attr := range attrs {
					n.SetAttribute(attr.Name, attr.Value)
				}
				return &immutableAttributes{n}
			}
		}
	}
	return nil
}

func (r *HTMLRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if !entering {
		return ast.WalkContinue, nil
	}
	language := n.Language(source)

	chromaFormatterOptions := make([]chromahtml.Option, len(r.FormatOptions))
	copy(chromaFormatterOptions, r.FormatOptions)
	style := styles.Get(r.Style)
	nohl := false

	var info []byte
	if n.Info != nil {
		info = n.Info.Segment.Value(source)
	}
	attrs := getAttributes(n, info)
	if attrs != nil {
		baseLineNumber := 1
		if linenostartAttr, ok := attrs.Get(linenostartAttrName); ok {
			baseLineNumber = int(linenostartAttr.(float64))
			chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.BaseLineNumber(baseLineNumber))
		}
		if linesAttr, hasLinesAttr := attrs.Get(highlightLinesAttrName); hasLinesAttr {
			if lines, ok := linesAttr.([]interface{}); ok {
				var hlRanges [][2]int
				for _, l := range lines {
					if ln, ok := l.(float64); ok {
						hlRanges = append(hlRanges, [2]int{int(ln) + baseLineNumber - 1, int(ln) + baseLineNumber - 1})
					}
					if rng, ok := l.([]uint8); ok {
						slices := strings.Split(string([]byte(rng)), "-")
						lhs, err := strconv.Atoi(slices[0])
						if err != nil {
							continue
						}
						rhs := lhs
						if len(slices) > 1 {
							rhs, err = strconv.Atoi(slices[1])
							if err != nil {
								continue
							}
						}
						hlRanges = append(hlRanges, [2]int{lhs + baseLineNumber - 1, rhs + baseLineNumber - 1})
					}
				}
				chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.HighlightLines(hlRanges))
			}
		}
		if styleAttr, hasStyleAttr := attrs.Get(styleAttrName); hasStyleAttr {
			styleStr := string([]byte(styleAttr.([]uint8)))
			style = styles.Get(styleStr)
		}
		if _, hasNohlAttr := attrs.Get(nohlAttrName); hasNohlAttr {
			nohl = true
		}

		if linenosAttr, ok := attrs.Get(linenosAttrName); ok {
			switch v := linenosAttr.(type) {
			case bool:
				chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.WithLineNumbers(v))
			case []uint8:
				if v != nil {
					chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.WithLineNumbers(true))
				}
				if bytes.Equal(v, linenosTableAttrValue) {
					chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.LineNumbersInTable(true))
				} else if bytes.Equal(v, linenosInlineAttrValue) {
					chromaFormatterOptions = append(chromaFormatterOptions, chromahtml.LineNumbersInTable(false))
				}
			}
		}
	}

	var lexer chroma.Lexer
	if language != nil {
		lexer = lexers.Get(string(language))
	}
	if !nohl && lexer != nil {
		if style == nil {
			style = styles.Fallback
		}
		var buffer bytes.Buffer
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			buffer.Write(line.Value(source))
		}
		iterator, err := lexer.Tokenise(nil, buffer.String())
		if err == nil {
			c := newCodeBlockContext(language, true, attrs)

			if r.CodeBlockOptions != nil {
				chromaFormatterOptions = append(chromaFormatterOptions, r.CodeBlockOptions(c)...)
			}
			formatter := chromahtml.New(chromaFormatterOptions...)
			if r.WrapperRenderer != nil {
				r.WrapperRenderer(w, c, true)
			}
			_ = formatter.Format(w, style, iterator) == nil
			if r.WrapperRenderer != nil {
				r.WrapperRenderer(w, c, false)
			}
			if r.CSSWriter != nil {
				_ = formatter.WriteCSS(r.CSSWriter, style)
			}
			return ast.WalkContinue, nil
		}
	}

	var c CodeBlockContext
	if r.WrapperRenderer != nil {
		c = newCodeBlockContext(language, false, attrs)
		r.WrapperRenderer(w, c, true)
	} else {
		_, _ = w.WriteString("<pre><code")
		language := n.Language(source)
		if language != nil {
			_, _ = w.WriteString(" class=\"language-")
			r.Writer.Write(w, language)
			_, _ = w.WriteString("\"")
		}
		_ = w.WriteByte('>')
	}
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
	if r.WrapperRenderer != nil {
		r.WrapperRenderer(w, c, false)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

type highlighting struct {
	options []Option
}

// Highlighting is a goldmark.Extender implementation.
var Highlighting = &highlighting{
	options: []Option{},
}

// NewHighlighting returns a new extension with given options.
func NewHighlighting(opts ...Option) goldmark.Extender {
	return &highlighting{
		options: opts,
	}
}

// Extend implements goldmark.Extender.
func (e *highlighting) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(e.options...), 200),
	))
}
