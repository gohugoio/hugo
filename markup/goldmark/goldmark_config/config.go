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

// Package goldmark_config holds Goldmark related configuration.
package goldmark_config

const (
	AutoHeadingIDTypeGitHub      = "github"
	AutoHeadingIDTypeGitHubAscii = "github-ascii"
	AutoHeadingIDTypeBlackfriday = "blackfriday"
)

// Default holds the default Goldmark configuration.
var Default = Config{
	Extensions: Extensions{
		Typographer: Typographer{
			Disable:          false,
			LeftSingleQuote:  "&lsquo;",
			RightSingleQuote: "&rsquo;",
			LeftDoubleQuote:  "&ldquo;",
			RightDoubleQuote: "&rdquo;",
			EnDash:           "&ndash;",
			EmDash:           "&mdash;",
			Ellipsis:         "&hellip;",
			LeftAngleQuote:   "&laquo;",
			RightAngleQuote:  "&raquo;",
			Apostrophe:       "&rsquo;",
		},
		Footnote:        true,
		DefinitionList:  true,
		Table:           true,
		Strikethrough:   true,
		Linkify:         true,
		LinkifyProtocol: "https",
		TaskList:        true,
		CJK: CJK{
			Enable:                   false,
			EastAsianLineBreaks:      false,
			EastAsianLineBreaksStyle: "simple",
			EscapedSpace:             false,
		},
		Extras: Extras{
			Delete: Delete{
				Enable: false,
			},
			Insert: Insert{
				Enable: false,
			},
			Mark: Mark{
				Enable: false,
			},
			Subscript: Subscript{
				Enable: false,
			},
			Superscript: Superscript{
				Enable: false,
			},
		},
		Passthrough: Passthrough{
			Enable: false,
			Delimiters: DelimitersConfig{
				Inline: [][]string{},
				Block:  [][]string{},
			},
		},
	},
	Renderer: Renderer{
		Unsafe: false,
	},
	Parser: Parser{
		AutoHeadingID:                      true,
		AutoHeadingIDType:                  AutoHeadingIDTypeGitHub,
		WrapStandAloneImageWithinParagraph: true,
		Attribute: ParserAttribute{
			Title: true,
			Block: false,
		},
	},
}

// Config configures Goldmark.
type Config struct {
	Renderer               Renderer
	Parser                 Parser
	Extensions             Extensions
	DuplicateResourceFiles bool
	RenderHooks            RenderHooks
}

// RenderHooks contains configuration for Goldmark render hooks.
type RenderHooks struct {
	Image ImageRenderHook
	Link  LinkRenderHook
}

// ImageRenderHook contains configuration for the image render hook.
type ImageRenderHook struct {
	// Enable the default image render hook.
	// We need to know if it is set or not, hence the pointer.
	EnableDefault *bool
}

func (h ImageRenderHook) IsEnableDefault() bool {
	return h.EnableDefault != nil && *h.EnableDefault
}

// LinkRenderHook contains configuration for the link render hook.
type LinkRenderHook struct {
	// Disable the default image render hook.
	// We need to know if it is set or not, hence the pointer.
	EnableDefault *bool
}

func (h LinkRenderHook) IsEnableDefault() bool {
	return h.EnableDefault != nil && *h.EnableDefault
}

type Extensions struct {
	Typographer    Typographer
	Footnote       bool
	DefinitionList bool
	Extras         Extras
	Passthrough    Passthrough

	// GitHub flavored markdown
	Table           bool
	Strikethrough   bool
	Linkify         bool
	LinkifyProtocol string
	TaskList        bool
	CJK             CJK
}

// Typographer holds typographer configuration.
type Typographer struct {
	// Whether to disable typographer.
	Disable bool

	// Value used for left single quote.
	LeftSingleQuote string
	// Value used for right single quote.
	RightSingleQuote string
	// Value used for left double quote.
	LeftDoubleQuote string
	// Value used for right double quote.
	RightDoubleQuote string
	// Value used for en dash.
	EnDash string
	// Value used for em dash.
	EmDash string
	// Value used for ellipsis.
	Ellipsis string
	// Value used for left angle quote.
	LeftAngleQuote string
	// Value used for right angle quote.
	RightAngleQuote string
	// Value used for apostrophe.
	Apostrophe string
}

// Extras holds extras configuration.
// github.com/hugoio/hugo-goldmark-extensions/extras
type Extras struct {
	Delete      Delete
	Insert      Insert
	Mark        Mark
	Subscript   Subscript
	Superscript Superscript
}

type Delete struct {
	Enable bool
}

type Insert struct {
	Enable bool
}

type Mark struct {
	Enable bool
}

type Subscript struct {
	Enable bool
}

type Superscript struct {
	Enable bool
}

// Passthrough holds passthrough configuration.
// github.com/hugoio/hugo-goldmark-extensions/passthrough
type Passthrough struct {
	// Whether to enable the extension
	Enable bool

	// The delimiters to use for inline and block passthroughs.
	Delimiters DelimitersConfig
}

type DelimitersConfig struct {
	// The delimiters to use for inline passthroughs. Each entry in the list
	// is a size-2 list of strings, where the first string is the opening delimiter
	// and the second string is the closing delimiter, e.g.,
	//
	// [["$", "$"], ["\\(", "\\)"]]
	Inline [][]string

	// The delimiters to use for block passthroughs. Same format as Inline.
	Block [][]string
}

type CJK struct {
	// Whether to enable CJK support.
	Enable bool

	// Whether softline breaks between east asian wide characters should be ignored.
	EastAsianLineBreaks bool

	// Styles of Line Breaking of EastAsianLineBreaks: "simple" or "css3draft"
	EastAsianLineBreaksStyle string

	// Whether a '\' escaped half-space(0x20) should not be rendered.
	EscapedSpace bool
}

type Renderer struct {
	// Whether softline breaks should be rendered as '<br>'
	HardWraps bool

	// XHTML instead of HTML5.
	XHTML bool

	// Allow raw HTML etc.
	Unsafe bool
}

type Parser struct {
	// Enables custom heading ids and
	// auto generated heading ids.
	AutoHeadingID bool

	// The strategy to use when generating heading IDs.
	// Available options are "github", "github-ascii".
	// Default is "github", which will create GitHub-compatible anchor names.
	AutoHeadingIDType string

	// Enables custom attributes.
	Attribute ParserAttribute

	// Whether to wrap stand-alone images within a paragraph or not.
	WrapStandAloneImageWithinParagraph bool
}

type ParserAttribute struct {
	// Enables custom attributes for titles.
	Title bool
	// Enables custom attributes for blocks.
	Block bool
}
