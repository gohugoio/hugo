package attributes

import (
	"strings"

	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// This extension is based on/inspired by https://github.com/mdigger/goldmark-attributes
// MIT License
// Copyright (c) 2019 Dmitry Sedykh

var (
	kindAttributesBlock = ast.NewNodeKind("AttributesBlock")
	attrNameID          = []byte("id")

	defaultParser = new(attrParser)
)

func New(cfg goldmark_config.Parser) goldmark.Extender {
	return &attrExtension{cfg: cfg}
}

type attrExtension struct {
	cfg goldmark_config.Parser
}

func (a *attrExtension) Extend(m goldmark.Markdown) {
	if a.cfg.Attribute.Block {
		m.Parser().AddOptions(
			parser.WithBlockParsers(
				util.Prioritized(defaultParser, 100)),
		)
	}
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{cfg: a.cfg}, 100),
		),
	)
}

type attrParser struct{}

func (a *attrParser) CanAcceptIndentedLine() bool {
	return false
}

func (a *attrParser) CanInterruptParagraph() bool {
	return true
}

func (a *attrParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
}

func (a *attrParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (a *attrParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	if attrs, ok := parser.ParseAttributes(reader); ok {
		// add attributes
		node := &attributesBlock{
			BaseBlock: ast.BaseBlock{},
		}
		for _, attr := range attrs {
			node.SetAttribute(attr.Name, attr.Value)
		}
		return node, parser.NoChildren
	}
	return nil, parser.RequireParagraph
}

func (a *attrParser) Trigger() []byte {
	return []byte{'{'}
}

type attributesBlock struct {
	ast.BaseBlock
}

func (a *attributesBlock) Dump(source []byte, level int) {
	attrs := a.Attributes()
	list := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		var (
			name  = util.BytesToReadOnlyString(attr.Name)
			value = util.BytesToReadOnlyString(util.EscapeHTML(attr.Value.([]byte)))
		)
		list[name] = value
	}
	ast.DumpHelper(a, source, level, list, nil)
}

func (a *attributesBlock) Kind() ast.NodeKind {
	return kindAttributesBlock
}

type transformer struct {
	cfg goldmark_config.Parser
}

func (a *transformer) isFragmentNode(n ast.Node) bool {
	switch n.Kind() {
	case east.KindDefinitionTerm, ast.KindHeading:
		return true
	default:
		return false
	}
}

func (a *transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	var attributes []ast.Node
	if a.cfg.Attribute.Block {
		attributes = make([]ast.Node, 0, 500)
	}
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if a.isFragmentNode(node) {
			if id, found := node.Attribute(attrNameID); !found {
				a.generateAutoID(node, reader, pc)
			} else {
				pc.IDs().Put(id.([]byte))
			}
		}

		if a.cfg.Attribute.Block && node.Kind() == kindAttributesBlock {
			// Attributes for fenced code blocks are handled in their own extension,
			// but note that we currently only support code block attributes when
			// CodeFences=true.
			if node.PreviousSibling() != nil && node.PreviousSibling().Kind() != ast.KindFencedCodeBlock && !node.HasBlankPreviousLines() {
				attributes = append(attributes, node)
				return ast.WalkSkipChildren, nil
			} else {
				// remove attributes node
				node.Parent().RemoveChild(node.Parent(), node)
			}
		}

		return ast.WalkContinue, nil
	})

	for _, attr := range attributes {
		if prev := attr.PreviousSibling(); prev != nil &&
			prev.Type() == ast.TypeBlock {
			for _, attr := range attr.Attributes() {
				if _, found := prev.Attribute(attr.Name); !found {
					prev.SetAttribute(attr.Name, attr.Value)
				}
			}
		}
		// remove attributes node
		attr.Parent().RemoveChild(attr.Parent(), attr)
	}
}

func (a *transformer) generateAutoID(n ast.Node, reader text.Reader, pc parser.Context) {
	var text []byte
	switch n := n.(type) {
	case *ast.Heading:
		if a.cfg.AutoHeadingID {
			text = textHeadingID(n, reader)
		}
	case *east.DefinitionTerm:
		if a.cfg.AutoDefinitionTermID {
			text = []byte(render.TextPlain(n, reader.Source()))
		}
	}

	if len(text) > 0 {
		headingID := pc.IDs().Generate(text, n.Kind())
		n.SetAttribute(attrNameID, headingID)
	}
}

// Markdown settext headers can have multiple lines, use the last line for the ID.
func textHeadingID(n *ast.Heading, reader text.Reader) []byte {
	text := render.TextPlain(n, reader.Source())
	if n.Lines().Len() > 1 {

		// For multiline headings, Goldmark's extension for headings returns the last line.
		// We have a slightly different approach, but in most cases the end result should be the same.
		// Instead of looking at the text segments in Lines (see #13405 for issues with that),
		// we split the text above and use the last line.
		parts := strings.Split(text, "\n")
		text = parts[len(parts)-1]
	}

	return []byte(text)
}
