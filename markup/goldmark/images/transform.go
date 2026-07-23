package images

import (
	"strconv"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

type (
	imagesExtension struct {
		wrapStandAloneImageWithinParagraph bool
	}
)

const (
	// Used to signal to the rendering step that an image is used in a block context.
	// Dont's change this; the prefix must match the internalAttrPrefix in the root goldmark package.
	AttrIsBlock = "_h__isBlock"
	AttrOrdinal = "_h__ordinal"
)

// New returns a goldmark v2 parser extension for images.
func New(wrapStandAloneImageWithinParagraph bool) parser.Extension {
	return &imagesExtension{wrapStandAloneImageWithinParagraph: wrapStandAloneImageWithinParagraph}
}

func (e *imagesExtension) ParserOptions(*parser.Config) []parser.Option {
	return []parser.Option{
		parser.WithASTTransformers(
			util.Prioritized[parser.ASTTransformer](&Transformer{wrapStandAloneImageWithinParagraph: e.wrapStandAloneImageWithinParagraph}, 300),
		),
	}
}

type Transformer struct {
	wrapStandAloneImageWithinParagraph bool
}

// Transform transforms the provided Markdown AST.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	var ordinal int
	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}

		if n, ok := node.(*ast.Image); ok {
			parent := n.Parent()
			// GOLDMARK-V2: node attribute values are text.MultilineValue now, so
			// these internal signal attributes must be encoded as strings and
			// decoded in the render hook.
			n.SetAttribute(AttrOrdinal, text.NewStringMultilineValue(strconv.Itoa(ordinal)))
			ordinal++

			if !t.wrapStandAloneImageWithinParagraph {
				isBlock := parent.ChildCount() == 1
				if isBlock {
					n.SetAttribute(AttrIsBlock, text.NewStringMultilineValue("true"))
				}

				if isBlock && parent.Kind() == ast.KindParagraph {
					for _, attr := range parent.Attributes() {
						// Transfer any attribute set down to the image.
						// Image elements does not support attributes on its own,
						// so it's safe to just set without checking first.
						n.SetAttribute(attr.Name, attr.Value)
					}
					grandParent := parent.Parent()
					// GOLDMARK-V2: ReplaceChild no longer takes the receiver as its
					// first argument.
					grandParent.ReplaceChild(parent, n)
				}
			}

		}

		return ast.WalkContinue, nil
	})
}
