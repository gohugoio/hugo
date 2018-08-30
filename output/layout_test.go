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
	noExtNoDelimMediaType.Suffixes = nil
	noExtNoDelimMediaType.Delimiter = ""

	noExtMediaType := media.TextType
	noExtMediaType.Suffixes = nil

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
		layoutOverride string
		tp             Format
		expect         []string
		expectCount    int
	}{
		{"Home", LayoutDescriptor{Kind: "home"}, "", ampType,
			[]string{"index.amp.html", "home.amp.html", "list.amp.html", "index.html", "home.html", "list.html", "_default/index.amp.html"}, 12},
		{"Home, HTML", LayoutDescriptor{Kind: "home"}, "", htmlFormat,
			// We will eventually get to index.html. This looks stuttery, but makes the lookup logic easy to understand.
			[]string{"index.html.html", "home.html.html"}, 12},
		{"Home, french language", LayoutDescriptor{Kind: "home", Lang: "fr"}, "", ampType,
			[]string{"index.fr.amp.html"},
			24},
		{"Home, no ext or delim", LayoutDescriptor{Kind: "home"}, "", noExtDelimFormat,
			[]string{"index.nem", "home.nem", "list.nem"}, 6},
		{"Home, no ext", LayoutDescriptor{Kind: "home"}, "", noExt,
			[]string{"index.nex", "home.nex", "list.nex"}, 6},
		{"Page, no ext or delim", LayoutDescriptor{Kind: "page"}, "", noExtDelimFormat,
			[]string{"_default/single.nem"}, 1},
		{"Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, "", ampType,
			[]string{"sect1/sect1.amp.html", "sect1/section.amp.html", "sect1/list.amp.html", "sect1/sect1.html", "sect1/section.html", "sect1/list.html", "section/sect1.amp.html", "section/section.amp.html"}, 18},
		{"Section with layout", LayoutDescriptor{Kind: "section", Section: "sect1", Layout: "mylayout"}, "", ampType,
			[]string{"sect1/mylayout.amp.html", "sect1/sect1.amp.html", "sect1/section.amp.html", "sect1/list.amp.html", "sect1/mylayout.html", "sect1/sect1.html"}, 24},
		{"Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, "", ampType,
			[]string{"taxonomy/tag.amp.html", "taxonomy/taxonomy.amp.html", "taxonomy/list.amp.html", "taxonomy/tag.html", "taxonomy/taxonomy.html"}, 18},
		{"Taxonomy term", LayoutDescriptor{Kind: "taxonomyTerm", Section: "categories"}, "", ampType,
			[]string{"taxonomy/categories.terms.amp.html", "taxonomy/terms.amp.html", "taxonomy/list.amp.html", "taxonomy/categories.terms.html", "taxonomy/terms.html"}, 18},
		{"Page", LayoutDescriptor{Kind: "page"}, "", ampType,
			[]string{"_default/single.amp.html", "_default/single.html"}, 2},
		{"Page with layout", LayoutDescriptor{Kind: "page", Layout: "mylayout"}, "", ampType,
			[]string{"_default/mylayout.amp.html", "_default/single.amp.html", "_default/mylayout.html", "_default/single.html"}, 4},
		{"Page with layout and type", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype"}, "", ampType,
			[]string{"myttype/mylayout.amp.html", "myttype/single.amp.html", "myttype/mylayout.html"}, 8},
		{"Page with layout and type with subtype", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype/mysubtype"}, "", ampType,
			[]string{"myttype/mysubtype/mylayout.amp.html", "myttype/mysubtype/single.amp.html", "myttype/mysubtype/mylayout.html"}, 8},
		// RSS
		{"RSS Home", LayoutDescriptor{Kind: "home"}, "", RSSFormat,
			[]string{"index.rss.xml", "home.rss.xml", "rss.xml"}, 15},
		{"RSS Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, "", RSSFormat,
			[]string{"sect1/sect1.rss.xml", "sect1/section.rss.xml", "sect1/rss.xml", "sect1/list.rss.xml", "sect1/sect1.xml", "sect1/section.xml"}, 22},
		{"RSS Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, "", RSSFormat,
			[]string{"taxonomy/tag.rss.xml", "taxonomy/taxonomy.rss.xml", "taxonomy/rss.xml", "taxonomy/list.rss.xml", "taxonomy/tag.xml", "taxonomy/taxonomy.xml"}, 22},
		{"RSS Taxonomy term", LayoutDescriptor{Kind: "taxonomyTerm", Section: "tag"}, "", RSSFormat,
			[]string{"taxonomy/tag.terms.rss.xml", "taxonomy/terms.rss.xml", "taxonomy/rss.xml", "taxonomy/list.rss.xml", "taxonomy/tag.terms.xml"}, 22},
		{"Home plain text", LayoutDescriptor{Kind: "home"}, "", JSONFormat,
			[]string{"_text/index.json.json", "_text/home.json.json"}, 12},
		{"Page plain text", LayoutDescriptor{Kind: "page"}, "", JSONFormat,
			[]string{"_text/_default/single.json.json", "_text/_default/single.json"}, 2},
		{"Reserved section, shortcodes", LayoutDescriptor{Kind: "section", Section: "shortcodes", Type: "shortcodes"}, "", ampType,
			[]string{"section/shortcodes.amp.html"}, 12},
		{"Reserved section, partials", LayoutDescriptor{Kind: "section", Section: "partials", Type: "partials"}, "", ampType,
			[]string{"section/partials.amp.html"}, 12},
	} {
		t.Run(this.name, func(t *testing.T) {
			l := NewLayoutHandler()

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

		})
	}

}

func BenchmarkLayout(b *testing.B) {
	descriptor := LayoutDescriptor{Kind: "taxonomyTerm", Section: "categories"}
	l := NewLayoutHandler()

	for i := 0; i < b.N; i++ {
		layouts, err := l.For(descriptor, HTMLFormat)
		require.NoError(b, err)
		require.NotEmpty(b, layouts)
	}
}
