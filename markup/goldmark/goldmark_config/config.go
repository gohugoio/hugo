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

// DefaultConfig holds the default Goldmark configuration.
var Default = Config{
	Extensions: Extensions{
		Typographer:    true,
		Footnote:       true,
		DefinitionList: true,
		Table:          true,
		Strikethrough:  true,
		Linkify:        true,
		TaskList:       true,
	},
	Renderer: Renderer{
		Unsafe: false,
	},
	Parser: Parser{
		AutoHeadingID:     true,
		AutoHeadingIDType: AutoHeadingIDTypeGitHub,
		Attribute:         true,
	},
}

// Config configures Goldmark.
type Config struct {
	Renderer   Renderer
	Parser     Parser
	Extensions Extensions
}

type Extensions struct {
	Typographer    bool
	Footnote       bool
	DefinitionList bool

	// GitHub flavored markdown
	Table         bool
	Strikethrough bool
	Linkify       bool
	TaskList      bool
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
	Attribute bool
}
