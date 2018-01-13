// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package output

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {

	noExtNoDelimMediaType := media.TextType
	noExtNoDelimMediaType.Suffix = ""
	noExtNoDelimMediaType.Delimiter = ""

	noExtMediaType := media.TextType
	noExtMediaType.Suffix = ""

	var (
		ampType = Format{
			Name:      "AMP",
			MediaType: media.HTMLType,
			BaseName:  "index",
		}

		htmlFormat = HTMLFormat

		noExtDelimFormat = Format{
			Name:      "NEM",
			MediaType: noExtNoDelimMediaType,
			BaseName:  "_redirects",
		}
		noExt = Format{
			Name:      "NEX",
			MediaType: noExtMediaType,
			BaseName:  "next",
		}
	)

	for _, this := range []struct {
		name           string
		d              LayoutDescriptor
		hasTheme       bool
		layoutOverride string
		tp             Format
		expect         []string
		expectCount    int
	}{
		{"Home", LayoutDescriptor{Kind: "home"}, true, "", ampType,
			[]string{"index.amp.html", "theme/index.amp.html", "home.amp.html", "theme/home.amp.html", "list.amp.html", "theme/list.amp.html", "index.html", "theme/index.html", "home.html", "theme/home.html", "list.html", "theme/list.html", "_default/index.amp.html"}, 24},
		{"Home, HTML", LayoutDescriptor{Kind: "home"}, true, "", htmlFormat,
			// We will eventually get to index.html. This looks stuttery, but makes the lookup logic easy to understand.
			[]string{"index.html.html", "theme/index.html.html", "home.html.html"}, 24},
		{"Home, french language", LayoutDescriptor{Kind: "home", Lang: "fr"}, true, "", ampType,
			[]string{"index.fr.amp.html", "theme/index.fr.amp.html"},
			48},
		{"Home, no ext or delim", LayoutDescriptor{Kind: "home"}, true, "", noExtDelimFormat,
			[]string{"index.nem", "theme/index.nem", "home.nem", "theme/home.nem", "list.nem"}, 12},
		{"Home, no ext", LayoutDescriptor{Kind: "home"}, true, "", noExt,
			[]string{"index.nex", "theme/index.nex", "home.nex", "theme/home.nex", "list.nex"}, 12},
		{"Page, no ext or delim", LayoutDescriptor{Kind: "page"}, true, "", noExtDelimFormat,
			[]string{"_default/single.nem", "theme/_default/single.nem"}, 2},
		{"Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, false, "", ampType,
			[]string{"sect1/sect1.amp.html", "sect1/section.amp.html", "sect1/list.amp.html", "sect1/sect1.html", "sect1/section.html", "sect1/list.html", "section/sect1.amp.html", "section/section.amp.html"}, 18},
		{"Section with layout", LayoutDescriptor{Kind: "section", Section: "sect1", Layout: "mylayout"}, false, "", ampType,
			[]string{"sect1/mylayout.amp.html", "sect1/sect1.amp.html", "sect1/section.amp.html", "sect1/list.amp.html", "sect1/mylayout.html", "sect1/sect1.html"}, 24},
		{"Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, false, "", ampType,
			[]string{"taxonomy/tag.amp.html", "taxonomy/taxonomy.amp.html", "taxonomy/list.amp.html", "taxonomy/tag.html", "taxonomy/taxonomy.html"}, 18},
		{"Taxonomy term", LayoutDescriptor{Kind: "taxonomyTerm", Section: "categories"}, false, "", ampType,
			[]string{"taxonomy/categories.terms.amp.html", "taxonomy/terms.amp.html", "taxonomy/list.amp.html", "taxonomy/categories.terms.html", "taxonomy/terms.html"}, 18},
		{"Page", LayoutDescriptor{Kind: "page"}, true, "", ampType,
			[]string{"_default/single.amp.html", "theme/_default/single.amp.html", "_default/single.html", "theme/_default/single.html"}, 4},
		{"Page with layout", LayoutDescriptor{Kind: "page", Layout: "mylayout"}, false, "", ampType,
			[]string{"_default/mylayout.amp.html", "_default/single.amp.html", "_default/mylayout.html", "_default/single.html"}, 4},
		{"Page with layout and type", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype"}, false, "", ampType,
			[]string{"myttype/mylayout.amp.html", "myttype/single.amp.html", "myttype/mylayout.html"}, 8},
		{"Page with layout and type with subtype", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype/mysubtype"}, false, "", ampType,
			[]string{"myttype/mysubtype/mylayout.amp.html", "myttype/mysubtype/single.amp.html", "myttype/mysubtype/mylayout.html"}, 8},
		// RSS
		{"RSS Home with theme", LayoutDescriptor{Kind: "home"}, true, "", RSSFormat,
			[]string{"index.rss.xml", "theme/index.rss.xml", "home.rss.xml", "theme/home.rss.xml", "rss.xml"}, 29},
		{"RSS Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, false, "", RSSFormat,
			[]string{"sect1/sect1.rss.xml", "sect1/section.rss.xml", "sect1/rss.xml", "sect1/list.rss.xml", "sect1/sect1.xml", "sect1/section.xml"}, 22},
		{"RSS Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, false, "", RSSFormat,
			[]string{"taxonomy/tag.rss.xml", "taxonomy/taxonomy.rss.xml", "taxonomy/rss.xml", "taxonomy/list.rss.xml", "taxonomy/tag.xml", "taxonomy/taxonomy.xml"}, 22},
		{"RSS Taxonomy term", LayoutDescriptor{Kind: "taxonomyTerm", Section: "tag"}, false, "", RSSFormat,
			[]string{"taxonomy/tag.terms.rss.xml", "taxonomy/terms.rss.xml", "taxonomy/rss.xml", "taxonomy/list.rss.xml", "taxonomy/tag.terms.xml"}, 22},
		{"Home plain text", LayoutDescriptor{Kind: "home"}, true, "", JSONFormat,
			[]string{"_text/index.json.json", "_text/theme/index.json.json", "_text/home.json.json", "_text/theme/home.json.json"}, 24},
		{"Page plain text", LayoutDescriptor{Kind: "page"}, true, "", JSONFormat,
			[]string{"_text/_default/single.json.json", "_text/theme/_default/single.json.json", "_text/_default/single.json", "_text/theme/_default/single.json"}, 4},
		{"Reserved section, shortcodes", LayoutDescriptor{Kind: "section", Section: "shortcodes", Type: "shortcodes"}, true, "", ampType,
			[]string{"section/shortcodes.amp.html", "theme/section/shortcodes.amp.html"}, 24},
		{"Reserved section, partials", LayoutDescriptor{Kind: "section", Section: "partials", Type: "partials"}, true, "", ampType,
			[]string{"section/partials.amp.html", "theme/section/partials.amp.html"}, 24},
	} {
		t.Run(this.name, func(t *testing.T) {
			l := NewLayoutHandler(this.hasTheme)

			layouts, err := l.For(this.d, this.tp)

			require.NoError(t, err)
			require.NotNil(t, layouts)
			require.True(t, len(layouts) >= len(this.expect), fmt.Sprint(layouts))
			// Not checking the complete list for now ...
			got := layouts[:len(this.expect)]
			if len(layouts) != this.expectCount || !reflect.DeepEqual(got, this.expect) {
				formatted := strings.Replace(fmt.Sprintf("%v", layouts), "[", "\"", 1)
				formatted = strings.Replace(formatted, "]", "\"", 1)
				formatted = strings.Replace(formatted, " ", "\", \"", -1)

				t.Fatalf("Got %d/%d:\n%v\nExpected:\n%v\nAll:\n%v\nFormatted:\n%s", len(layouts), this.expectCount, got, this.expect, layouts, formatted)

			}

			if !this.hasTheme {
				for _, layout := range layouts {
					require.NotContains(t, layout, "theme")
				}
			}
		})
	}

}

func BenchmarkLayout(b *testing.B) {
	descriptor := LayoutDescriptor{Kind: "taxonomyTerm", Section: "categories"}
	l := NewLayoutHandler(false)

	for i := 0; i < b.N; i++ {
		layouts, err := l.For(descriptor, HTMLFormat)
		require.NoError(b, err)
		require.NotEmpty(b, layouts)
	}
}
