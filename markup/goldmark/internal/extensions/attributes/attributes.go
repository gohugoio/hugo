package attributes

import (
	"strings"

	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/yuin/goldmark/v2/ast"
	east "github.com/yuin/goldmark/v2/extension/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// This extension is based on/inspired by https://github.com/mdigger/goldmark-attributes
// MIT License
// Copyright (c) 2019 Dmitry Sedykh

var (
	kindAttributesBlock = ast.NewNodeKind("AttributesBlock")
	attrNameID          = "id"

	defaultParser = new(attrParser)
)

// New returns a goldmark v2 parser extension for block attributes and auto IDs.
func New(cfg goldmark_config.Parser) parser.Extension {
	return &attrExtension{cfg: cfg}
}

type attrExtension struct {
	cfg goldmark_config.Parser
}

func (a *attrExtension) ParserOptions(*parser.Config) []parser.Option {
	var opts []parser.Option
	if a.cfg.Attribute.Block {
		opts = append(opts,
			parser.WithBlockParsers(
				util.Prioritized[parser.BlockParser](defaultParser, 100)),
		)
	}
	opts = append(opts,
		parser.WithASTTransformers(
			util.Prioritized[parser.ASTTransformer](&transformer{cfg: a.cfg}, 100),
		),
	)
	return opts
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

// Dump implements Node.Dump.
// GOLDMARK-V2: Dump now returns *ast.NodeDump instead of taking a level and
// calling ast.DumpHelper (which was removed).
func (a *attributesBlock) Dump(source []byte) *ast.NodeDump {
	attrs := a.Attributes()
	list := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		list[attr.Name] = string(util.EscapeHTML(attr.Value.Bytes(nil)))
	}
	return ast.NewNodeDump(a, list)
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
	var solitaryAttributeNodes []ast.Node
	if a.cfg.Attribute.Block {
		attributes = make([]ast.Node, 0, 100)
	}
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if a.isFragmentNode(node) {
			if id, found := node.Attribute(attrNameID); !found {
				a.generateAutoID(node, reader, pc)
			} else {
				pc.IDs().Put(id.Bytes(reader.Source()))
			}
		}

		if a.cfg.Attribute.Block && node.Kind() == kindAttributesBlock {
			// Attributes for fenced code blocks are handled in their own extension,
			// but note that we currently only support code block attributes when
			// CodeFences=true.
			// GOLDMARK-V2: KindFencedCodeBlock was merged into KindCodeBlock;
			// HasBlankPreviousLines moved to the ast.BlockNode interface.
			bn, isBlock := node.(ast.BlockNode)
			if node.PreviousSibling() != nil && node.PreviousSibling().Kind() != ast.KindCodeBlock && isBlock && !bn.HasBlankPreviousLines() {
				attributes = append(attributes, node)
				return ast.WalkSkipChildren, nil
			} else {
				solitaryAttributeNodes = append(solitaryAttributeNodes, node)
			}
		}

		return ast.WalkContinue, nil
	})

	for _, attr := range attributes {
		// GOLDMARK-V2: Node.Type()/ast.TypeBlock were removed; use ast.BlockNode.
		if prev := attr.PreviousSibling(); prev != nil && isBlockNode(prev) {
			for _, attr := range attr.Attributes() {
				if _, found := prev.Attribute(attr.Name); !found {
					prev.SetAttribute(attr.Name, attr.Value)
				}
			}
		}
		// remove attributes node
		// GOLDMARK-V2: RemoveChild no longer takes the receiver as first arg.
		attr.Parent().RemoveChild(attr)
	}

	// Remove any solitary attribute nodes.
	for _, n := range solitaryAttributeNodes {
		n.Parent().RemoveChild(n)
	}
}

func (a *transformer) generateAutoID(n ast.Node, reader text.Reader, pc parser.Context) {
	var idText []byte
	switch n := n.(type) {
	case *ast.Heading:
		if a.cfg.AutoHeadingID {
			idText = textHeadingID(n, reader)
		}
	case *east.DefinitionTerm:
		if a.cfg.AutoDefinitionTermID {
			idText = []byte(render.TextPlain(n, reader.Source()))
		}
	}

	if len(idText) > 0 {
		headingID := pc.IDs().Generate(idText, n.Kind())
		n.SetAttribute(attrNameID, text.NewStringMultilineValue(string(headingID)))
	}
}

func isBlockNode(n ast.Node) bool {
	_, ok := n.(ast.BlockNode)
	return ok
}

// Markdown settext headers can have multiple lines, use the last line for the ID.
func textHeadingID(n *ast.Heading, reader text.Reader) []byte {
	text := render.TextPlain(n, reader.Source())
	// GOLDMARK-V2: block nodes no longer expose Lines(), so we can't cheaply
	// detect multi-line setext headings. We always take the last line of the
	// plain text, which matches the previous behavior in the common cases.
	// See #13405 for the caveats with the segment-based approach.
	parts := strings.Split(text, "\n")
	text = parts[len(parts)-1]

	return []byte(text)
}
