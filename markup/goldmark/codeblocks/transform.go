package codeblocks

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Kind is the kind of an Hugo code block.
var KindCodeBlock = ast.NewNodeKind("HugoCodeBlock")

// Its raw contents are the plain text of the code block.
type codeBlock struct {
	ast.BaseBlock
	ordinal int
	b       *ast.FencedCodeBlock
}

func (*codeBlock) Kind() ast.NodeKind { return KindCodeBlock }

func (*codeBlock) IsRaw() bool { return true }

func (b *codeBlock) Dump(src []byte, level int) {
}

type Transformer struct{}

// Transform transforms the provided Markdown AST.
func (*Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	var codeBlocks []*ast.FencedCodeBlock

	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}

		cb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		codeBlocks = append(codeBlocks, cb)

		return ast.WalkContinue, nil
	})

	for i, cb := range codeBlocks {
		b := &codeBlock{b: cb, ordinal: i}
		parent := cb.Parent()
		if parent != nil {
			parent.ReplaceChild(parent, cb, b)
		}
	}
}
