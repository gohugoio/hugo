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

package hooks

import (
	"context"
	"io"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/internal/attributes"
)

var _ AttributesOptionsSliceProvider = (*attributes.AttributesHolder)(nil)

type AttributesProvider interface {
	// Attributes passed in from Markdown (e.g. { attrName1=attrValue1 attrName2="attr Value 2" }).
	Attributes() map[string]any
}

// LinkContext is the context passed to a link render hook.
type LinkContext interface {
	PageProvider

	// The link URL.
	Destination() string

	// The link title attribute.
	Title() string

	// The rendered (HTML) text.
	Text() hstring.HTML

	// The plain variant of Text.
	PlainText() string
}

// ImageLinkContext is the context passed to a image link render hook.
type ImageLinkContext interface {
	LinkContext

	// Returns true if this is a standalone image and the config option
	// markup.goldmark.parser.wrapStandAloneImageWithinParagraph is disabled.
	IsBlock() bool

	// Zero-based ordinal for all the images in the current document.
	Ordinal() int
}

// CodeblockContext is the context passed to a code block render hook.
type CodeblockContext interface {
	BaseContext
	AttributesProvider

	// Chroma highlighting processing options. This will only be filled if Type is a known Chroma Lexer.
	Options() map[string]any

	// The type of code block. This will be the programming language, e.g. bash, when doing code highlighting.
	Type() string

	// The text between the code fences.
	Inner() string
}

// TableContext is the context passed to a table render hook.
type TableContext interface {
	BaseContext
	AttributesProvider

	THead() []TableRow
	TBody() []TableRow
}

// BaseContext is the base context used in most render hooks.
type BaseContext interface {
	text.Positioner
	PageProvider

	// Zero-based ordinal for all elements of this kind in the current document.
	Ordinal() int
}

// BlockquoteContext is the context passed to a blockquote render hook.
type BlockquoteContext interface {
	BaseContext

	AttributesProvider

	// The blockquote text.
	// If type is "alert", this will be the alert text.
	Text() hstring.HTML

	/// Returns the blockquote type, one of "regular" and "alert".
	// Type "alert" indicates that this is a GitHub type alert.
	Type() string

	// The GitHub alert type converted to lowercase, e.g. "note".
	// Only set if Type is "alert".
	AlertType() string

	// The alert title.
	// Currently only relevant for Obsidian alerts.
	// GitHub does not suport alert titles and will not render alerts with titles.
	AlertTitle() hstring.HTML

	// The alert sign,  "+" or "-" or "" used to indicate folding.
	// Currently only relevant for Obsidian alerts.
	// GitHub does not suport alert signs and will not render alerts with signs.
	AlertSign() string
}

type PositionerSourceTargetProvider interface {
	// For internal use.
	PositionerSourceTarget() []byte
}

// PassThroughContext is the context passed to a passthrough render hook.
type PassthroughContext interface {
	BaseContext
	AttributesProvider

	// Currently one of "inline" or "block".
	Type() string

	// The inner content of the passthrough element, excluding the delimiters.
	Inner() string
}

type AttributesOptionsSliceProvider interface {
	AttributesSlice() []attributes.Attribute
	OptionsSlice() []attributes.Attribute
}

type LinkRenderer interface {
	RenderLink(cctx context.Context, w io.Writer, ctx LinkContext) error
}

type CodeBlockRenderer interface {
	RenderCodeblock(cctx context.Context, w hugio.FlexiWriter, ctx CodeblockContext) error
}

type BlockquoteRenderer interface {
	RenderBlockquote(cctx context.Context, w hugio.FlexiWriter, ctx BlockquoteContext) error
}

type TableRenderer interface {
	RenderTable(cctx context.Context, w hugio.FlexiWriter, ctx TableContext) error
}

type PassthroughRenderer interface {
	RenderPassthrough(cctx context.Context, w io.Writer, ctx PassthroughContext) error
}

type IsDefaultCodeBlockRendererProvider interface {
	IsDefaultCodeBlockRenderer() bool
}

// HeadingContext contains accessors to all attributes that a HeadingRenderer
// can use to render a heading.
type HeadingContext interface {
	PageProvider
	// Level is the level of the header (i.e. 1 for top-level, 2 for sub-level, etc.).
	Level() int
	// Anchor is the HTML id assigned to the heading.
	Anchor() string
	// Text is the rendered (HTML) heading text, excluding the heading marker.
	Text() hstring.HTML
	// PlainText is the unrendered version of Text.
	PlainText() string

	// Attributes (e.g. CSS classes)
	AttributesProvider
}

type PageProvider interface {
	// Page is the page being rendered.
	Page() any

	// PageInner may be different than Page when .RenderShortcodes is in play.
	// The main use case for this is to include other pages' markdown into the current page
	// but resolve resources and pages relative to the original.
	PageInner() any
}

// HeadingRenderer describes a uniquely identifiable rendering hook.
type HeadingRenderer interface {
	// RenderHeading writes the rendered content to w using the data in w.
	RenderHeading(cctx context.Context, w io.Writer, ctx HeadingContext) error
}

// ElementPositionResolver provides a way to resolve the start Position
// of a markdown element in the original source document.
// This may be both slow and approximate, so should only be
// used for error logging.
type ElementPositionResolver interface {
	ResolvePosition(ctx any) text.Position
}

type RendererType int

const (
	LinkRendererType RendererType = iota + 1
	ImageRendererType
	HeadingRendererType
	CodeBlockRendererType
	PassthroughRendererType
	BlockquoteRendererType
	TableRendererType
)

type GetRendererFunc func(t RendererType, id any) any

type TableCell struct {
	Text      hstring.HTML
	Alignment string // left, center, or right
}

type TableRow []TableCell

type Table struct {
	THead []TableRow
	TBody []TableRow
}
