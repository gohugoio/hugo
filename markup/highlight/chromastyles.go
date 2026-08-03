// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/alecthomas/chroma/v3"
	"github.com/alecthomas/chroma/v3/formatters/html"
	"github.com/alecthomas/chroma/v3/styles"
)

// ChromaStylesOptions holds options for ChromaStylesCSS.
type ChromaStylesOptions struct {
	// Highlighter style, e.g. "monokai".
	Style string

	// Style mode, "light" or "dark".
	Mode string

	// Scope selectors under a top level mode class, e.g. ".dark .chroma".
	ModeSelector bool

	// Foreground and background colors for highlighted lines, e.g. "#fff000 bg:#000fff".
	HighlightStyle string

	// Foreground and background colors for inline line numbers.
	LineNumbersInlineStyle string

	// Foreground and background colors for table line numbers.
	LineNumbersTableStyle string

	// Omit CSS class comment prefixes in the generated CSS.
	OmitClassComments bool
}

// ChromaStylesCSS generates a CSS stylesheet for the Chroma code highlighter.
// This stylesheet is needed if markup.highlight.noClasses is disabled in config.
func ChromaStylesCSS(opts ChromaStylesOptions) (string, error) {
	name := strings.ToLower(opts.Style)
	if !slices.Contains(styles.Names(), name) {
		return "", fmt.Errorf("invalid style: %s", name)
	}
	var chromaStyle *chroma.Style
	if opts.Mode != "" {
		var mode chroma.Mode
		switch opts.Mode {
		case "light":
			mode = chroma.Light
		case "dark":
			mode = chroma.Dark
		default:
			return "", fmt.Errorf("invalid mode: %s", opts.Mode)
		}

		chromaStyle = styles.GetForMode(name, mode)
		if chromaStyle.Mode() != mode {
			return "", fmt.Errorf("style %q does not have a %q mode", name, opts.Mode)
		}
	} else {
		chromaStyle = styles.Get(name)
	}
	builder := chromaStyle.Builder()
	if opts.HighlightStyle != "" {
		builder.Add(chroma.LineHighlight, opts.HighlightStyle)
	}
	if opts.LineNumbersInlineStyle != "" {
		builder.Add(chroma.LineNumbers, opts.LineNumbersInlineStyle)
	}
	if opts.LineNumbersTableStyle != "" {
		builder.Add(chroma.LineNumbersTable, opts.LineNumbersTableStyle)
	}
	style, err := builder.Build()
	if err != nil {
		return "", err
	}

	options := []html.Option{
		html.WithCSSComments(!opts.OmitClassComments),
		html.WithCustomCSS(chromaCSSOverrides(style)),
	}

	var buf bytes.Buffer
	if err := html.New(options...).WriteCSS(&buf, style); err != nil {
		return "", err
	}
	css := buf.String()
	if opts.ModeSelector {
		// Scope every selector under a top level mode class, e.g. ".dark .chroma".
		// This allows generating both light and dark stylesheets and toggling
		// them with a parent class on the page.
		// There's no upstream option for this (I think), so do string replacements for now.
		// TODO(bep) upstream option for this.
		var prefix string
		switch style.Mode() {
		case chroma.Light:
			prefix = ".light "
		case chroma.Dark:
			prefix = ".dark "
		default:
			return "", fmt.Errorf("style %q does not have a light or dark mode", style.Name)
		}
		replacer := strings.NewReplacer(
			".bg {", prefix+".bg {",
			".chroma ", prefix+".chroma ",
		)
		css = replacer.Replace(css)
	}

	return css, nil
}

// chromaCSSOverrides re-emits leaf token colors that Chroma's minifier drops
// because they equal the style's default foreground (the .chroma color).
// In a paired light/dark setup, the other sheet's explicit rule (e.g. .chroma .nx)
// would otherwise leak into this mode, since an explicit declaration beats
// inheritance from .chroma.
// TODO(bep): replace with an upstream option, see issue 15161.
func chromaCSSOverrides(style *chroma.Style) map[chroma.TokenType]string {
	bg := style.Get(chroma.Background)
	m := make(map[chroma.TokenType]string)
	for tt := range chroma.StandardTypes {
		if tt == chroma.Background || !style.Has(tt) || !chromaLeafToken(tt) {
			continue
		}
		entry := style.Get(tt)
		if !entry.Sub(bg).IsZero() || !entry.Colour.IsSet() {
			continue
		}
		if css := html.StyleEntryToCSS(chroma.StyleEntry{Colour: entry.Colour}); css != "" {
			m[tt] = css
		}
	}
	return m
}

func chromaLeafToken(tt chroma.TokenType) bool {
	for other := range chroma.StandardTypes {
		if other != tt && (other.Category() == tt || other.SubCategory() == tt) {
			return false
		}
	}
	return true
}
