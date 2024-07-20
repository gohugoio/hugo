// Copyright 2024 The Hugo Authors. All rights reserved.
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

package media

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecodeTypes(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name        string
		m           map[string]any
		shouldError bool
		assert      func(t *testing.T, name string, tt Types)
	}{
		{
			"Redefine JSON",
			map[string]any{
				"application/json": map[string]any{
					"suffixes": []string{"jasn"},
				},
			},

			false,
			func(t *testing.T, name string, tt Types) {
				for _, ttt := range tt {
					if _, ok := DefaultTypes.GetByType(ttt.Type); !ok {
						fmt.Println(ttt.Type, "not found in default types")
					}
				}

				c.Assert(len(tt), qt.Equals, len(DefaultTypes))
				json, si, found := tt.GetBySuffix("jasn")
				c.Assert(found, qt.Equals, true)
				c.Assert(json.String(), qt.Equals, "application/json")
				c.Assert(si.FullSuffix, qt.Equals, ".jasn")
			},
		},
		{
			"MIME suffix in key, multiple file suffixes, custom delimiter",
			map[string]any{
				"application/hugo+hg": map[string]any{
					"suffixes":  []string{"hg1", "hG2"},
					"Delimiter": "_",
				},
			},
			false,
			func(t *testing.T, name string, tt Types) {
				c.Assert(len(tt), qt.Equals, len(DefaultTypes)+1)
				hg, si, found := tt.GetBySuffix("hg2")
				c.Assert(found, qt.Equals, true)
				c.Assert(hg.FirstSuffix.Suffix, qt.Equals, "hg1")
				c.Assert(hg.FirstSuffix.FullSuffix, qt.Equals, "_hg1")
				c.Assert(si.Suffix, qt.Equals, "hg2")
				c.Assert(si.FullSuffix, qt.Equals, "_hg2")
				c.Assert(hg.String(), qt.Equals, "application/hugo+hg")

				_, found = tt.GetByType("application/hugo+hg")
				c.Assert(found, qt.Equals, true)
			},
		},
		{
			"Add custom media type",
			map[string]any{
				"text/hugo+hgo": map[string]any{
					"Suffixes": []string{"hgo2"},
				},
			},
			false,
			func(t *testing.T, name string, tp Types) {
				c.Assert(len(tp), qt.Equals, len(DefaultTypes)+1)
				// Make sure we have not broken the default config.

				_, _, found := tp.GetBySuffix("json")
				c.Assert(found, qt.Equals, true)

				hugo, _, found := tp.GetBySuffix("hgo2")
				c.Assert(found, qt.Equals, true)
				c.Assert(hugo.String(), qt.Equals, "text/hugo+hgo")
			},
		},
	}

	for _, test := range tests {
		result, err := DecodeTypes(test.m)
		if test.shouldError {
			c.Assert(err, qt.Not(qt.IsNil))
		} else {
			c.Assert(err, qt.IsNil)
			test.assert(t, test.name, result.Config)
		}
	}
}

func TestDefaultTypes(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		tp               Type
		expectedMainType string
		expectedSubType  string
		expectedSuffixes string
		expectedType     string
		expectedString   string
	}{
		{Builtin.CalendarType, "text", "calendar", "ics", "text/calendar", "text/calendar"},
		{Builtin.CSSType, "text", "css", "css", "text/css", "text/css"},
		{Builtin.SCSSType, "text", "x-scss", "scss", "text/x-scss", "text/x-scss"},
		{Builtin.CSVType, "text", "csv", "csv", "text/csv", "text/csv"},
		{Builtin.HTMLType, "text", "html", "html,htm", "text/html", "text/html"},
		{Builtin.MarkdownType, "text", "markdown", "md,mdown,markdown", "text/markdown", "text/markdown"},
		{Builtin.EmacsOrgModeType, "text", "org", "org", "text/org", "text/org"},
		{Builtin.PandocType, "text", "pandoc", "pandoc,pdc", "text/pandoc", "text/pandoc"},
		{Builtin.ReStructuredTextType, "text", "rst", "rst", "text/rst", "text/rst"},
		{Builtin.AsciiDocType, "text", "asciidoc", "adoc,asciidoc,ad", "text/asciidoc", "text/asciidoc"},
		{Builtin.JavascriptType, "text", "javascript", "js,jsm,mjs", "text/javascript", "text/javascript"},
		{Builtin.TypeScriptType, "text", "typescript", "ts", "text/typescript", "text/typescript"},
		{Builtin.TSXType, "text", "tsx", "tsx", "text/tsx", "text/tsx"},
		{Builtin.JSXType, "text", "jsx", "jsx", "text/jsx", "text/jsx"},
		{Builtin.JSONType, "application", "json", "json", "application/json", "application/json"},
		{Builtin.RSSType, "application", "rss", "xml,rss", "application/rss+xml", "application/rss+xml"},
		{Builtin.SVGType, "image", "svg", "svg", "image/svg+xml", "image/svg+xml"},
		{Builtin.TextType, "text", "plain", "txt", "text/plain", "text/plain"},
		{Builtin.XMLType, "application", "xml", "xml", "application/xml", "application/xml"},
		{Builtin.TOMLType, "application", "toml", "toml", "application/toml", "application/toml"},
		{Builtin.YAMLType, "application", "yaml", "yaml,yml", "application/yaml", "application/yaml"},
		{Builtin.PDFType, "application", "pdf", "pdf", "application/pdf", "application/pdf"},
		{Builtin.TrueTypeFontType, "font", "ttf", "ttf", "font/ttf", "font/ttf"},
		{Builtin.OpenTypeFontType, "font", "otf", "otf", "font/otf", "font/otf"},
	} {
		c.Assert(test.tp.MainType, qt.Equals, test.expectedMainType)
		c.Assert(test.tp.SubType, qt.Equals, test.expectedSubType)
		c.Assert(test.tp.SuffixesCSV, qt.Equals, test.expectedSuffixes)
		c.Assert(test.tp.Type, qt.Equals, test.expectedType)
		c.Assert(test.tp.String(), qt.Equals, test.expectedString)

	}

	c.Assert(len(DefaultTypes), qt.Equals, 40)
}
