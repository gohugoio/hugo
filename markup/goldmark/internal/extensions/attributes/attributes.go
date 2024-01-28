package attributes

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// This extension is based on/inspired by https://github.com/mdigger/goldmark-attributes
// MIT License
// Copyright (c) 2019 Dmitry Sedykh

var (
	kindAttributesBlock = ast.NewNodeKind("AttributesBlock")

	defaultParser                        = new(attrParser)
	defaultTransformer                   = new(transformer)
	attributes         goldmark.Extender = new(attrExtension)
)

func New() goldmark.Extender {
	return attributes
}

type attrExtension struct{}

func (a *attrExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(defaultParser, 100)),
		parser.WithASTTransformers(
			util.Prioritized(defaultTransformer, 100),
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

type transformer struct{}

func (a *transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	attributes := make([]ast.Node, 0, 500)
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == kindAttributesBlock {
			// Attributes for fenced code blocks are handled in their own extension,
			// but note that we currently only support code block attributes when
			// CodeFences=true.
			if node.PreviousSibling() != nil && node.PreviousSibling().Kind() != ast.KindFencedCodeBlock && !node.HasBlankPreviousLines() {
				attributes = append(attributes, node)
				return ast.WalkSkipChildren, nil
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
