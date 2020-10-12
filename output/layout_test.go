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

	qt "github.com/frankban/quicktest"
	"github.com/kylelemons/godebug/diff"
)

func TestLayout(t *testing.T) {
	c := qt.New(t)

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
		name             string
		layoutDescriptor LayoutDescriptor
		layoutOverride   string
		format           Format
		expect           []string
	}{
		{
			"Home", LayoutDescriptor{Kind: "home"}, "", ampType,
			[]string{
				"index.amp.html",
				"home.amp.html",
				"list.amp.html",
				"index.html",
				"home.html",
				"list.html",
				"_default/index.amp.html",
				"_default/home.amp.html",
				"_default/list.amp.html",
				"_default/index.html",
				"_default/home.html",
				"_default/list.html",
			},
		},
		{
			"Home baseof", LayoutDescriptor{Kind: "home", Baseof: true}, "", ampType,
			[]string{
				"index-baseof.amp.html",
				"home-baseof.amp.html",
				"list-baseof.amp.html",
				"baseof.amp.html",
				"index-baseof.html",
				"home-baseof.html",
				"list-baseof.html",
				"baseof.html",
				"_default/index-baseof.amp.html",
				"_default/home-baseof.amp.html",
				"_default/list-baseof.amp.html",
				"_default/baseof.amp.html",
				"_default/index-baseof.html",
				"_default/home-baseof.html",
				"_default/list-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Home, HTML", LayoutDescriptor{Kind: "home"}, "", htmlFormat,
			// We will eventually get to index.html. This looks stuttery, but makes the lookup logic easy to understand.
			[]string{
				"index.html.html",
				"home.html.html",
				"list.html.html",
				"index.html",
				"home.html",
				"list.html",
				"_default/index.html.html",
				"_default/home.html.html",
				"_default/list.html.html",
				"_default/index.html",
				"_default/home.html",
				"_default/list.html",
			},
		},
		{
			"Home, HTML, baseof", LayoutDescriptor{Kind: "home", Baseof: true}, "", htmlFormat,
			[]string{
				"index-baseof.html.html",
				"home-baseof.html.html",
				"list-baseof.html.html",
				"baseof.html.html",
				"index-baseof.html",
				"home-baseof.html",
				"list-baseof.html",
				"baseof.html",
				"_default/index-baseof.html.html",
				"_default/home-baseof.html.html",
				"_default/list-baseof.html.html",
				"_default/baseof.html.html",
				"_default/index-baseof.html",
				"_default/home-baseof.html",
				"_default/list-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Home, french language", LayoutDescriptor{Kind: "home", Lang: "fr"}, "", ampType,
			[]string{
				"index.fr.amp.html",
				"home.fr.amp.html",
				"list.fr.amp.html",
				"index.amp.html",
				"home.amp.html",
				"list.amp.html",
				"index.fr.html",
				"home.fr.html",
				"list.fr.html",
				"index.html",
				"home.html",
				"list.html",
				"_default/index.fr.amp.html",
				"_default/home.fr.amp.html",
				"_default/list.fr.amp.html",
				"_default/index.amp.html",
				"_default/home.amp.html",
				"_default/list.amp.html",
				"_default/index.fr.html",
				"_default/home.fr.html",
				"_default/list.fr.html",
				"_default/index.html",
				"_default/home.html",
				"_default/list.html",
			},
		},
		{
			"Home, no ext or delim", LayoutDescriptor{Kind: "home"}, "", noExtDelimFormat,
			[]string{
				"index.nem",
				"home.nem",
				"list.nem",
				"_default/index.nem",
				"_default/home.nem",
				"_default/list.nem",
			},
		},
		{
			"Home, no ext", LayoutDescriptor{Kind: "home"}, "", noExt,
			[]string{
				"index.nex",
				"home.nex",
				"list.nex",
				"_default/index.nex",
				"_default/home.nex",
				"_default/list.nex",
			},
		},
		{
			"Page, no ext or delim", LayoutDescriptor{Kind: "page"}, "", noExtDelimFormat,
			[]string{"_default/single.nem"},
		},
		{
			"Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, "", ampType,
			[]string{
				"sect1/sect1.amp.html",
				"sect1/section.amp.html",
				"sect1/list.amp.html",
				"sect1/sect1.html",
				"sect1/section.html",
				"sect1/list.html",
				"section/sect1.amp.html",
				"section/section.amp.html",
				"section/list.amp.html",
				"section/sect1.html",
				"section/section.html",
				"section/list.html",
				"_default/sect1.amp.html",
				"_default/section.amp.html",
				"_default/list.amp.html",
				"_default/sect1.html",
				"_default/section.html",
				"_default/list.html",
			},
		},
		{
			"Section, baseof", LayoutDescriptor{Kind: "section", Section: "sect1", Baseof: true}, "", ampType,
			[]string{
				"sect1/sect1-baseof.amp.html",
				"sect1/section-baseof.amp.html",
				"sect1/list-baseof.amp.html",
				"sect1/baseof.amp.html",
				"sect1/sect1-baseof.html",
				"sect1/section-baseof.html",
				"sect1/list-baseof.html",
				"sect1/baseof.html",
				"section/sect1-baseof.amp.html",
				"section/section-baseof.amp.html",
				"section/list-baseof.amp.html",
				"section/baseof.amp.html",
				"section/sect1-baseof.html",
				"section/section-baseof.html",
				"section/list-baseof.html",
				"section/baseof.html",
				"_default/sect1-baseof.amp.html",
				"_default/section-baseof.amp.html",
				"_default/list-baseof.amp.html",
				"_default/baseof.amp.html",
				"_default/sect1-baseof.html",
				"_default/section-baseof.html",
				"_default/list-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Section with layout", LayoutDescriptor{Kind: "section", Section: "sect1", Layout: "mylayout"}, "", ampType,
			[]string{
				"sect1/mylayout.amp.html",
				"sect1/sect1.amp.html",
				"sect1/section.amp.html",
				"sect1/list.amp.html",
				"sect1/mylayout.html",
				"sect1/sect1.html",
				"sect1/section.html",
				"sect1/list.html",
				"section/mylayout.amp.html",
				"section/sect1.amp.html",
				"section/section.amp.html",
				"section/list.amp.html",
				"section/mylayout.html",
				"section/sect1.html",
				"section/section.html",
				"section/list.html",
				"_default/mylayout.amp.html",
				"_default/sect1.amp.html",
				"_default/section.amp.html",
				"_default/list.amp.html",
				"_default/mylayout.html",
				"_default/sect1.html",
				"_default/section.html",
				"_default/list.html",
			},
		},
		{
			"Term", LayoutDescriptor{Kind: "term", Section: "tags"}, "", ampType,
			[]string{
				"term/term.amp.html",
				"term/tags.amp.html",
				"term/taxonomy.amp.html",
				"term/list.amp.html",
				"term/term.html",
				"term/tags.html",
				"term/taxonomy.html",
				"term/list.html",
				"taxonomy/term.amp.html",
				"taxonomy/tags.amp.html",
				"taxonomy/taxonomy.amp.html",
				"taxonomy/list.amp.html",
				"taxonomy/term.html",
				"taxonomy/tags.html",
				"taxonomy/taxonomy.html",
				"taxonomy/list.html",
				"tags/term.amp.html",
				"tags/tags.amp.html",
				"tags/taxonomy.amp.html",
				"tags/list.amp.html",
				"tags/term.html",
				"tags/tags.html",
				"tags/taxonomy.html",
				"tags/list.html",
				"_default/term.amp.html",
				"_default/tags.amp.html",
				"_default/taxonomy.amp.html",
				"_default/list.amp.html",
				"_default/term.html",
				"_default/tags.html",
				"_default/taxonomy.html",
				"_default/list.html",
			},
		},
		{
			"Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "categories"}, "", ampType,
			[]string{
				"categories/categories.terms.amp.html",
				"categories/terms.amp.html",
				"categories/taxonomy.amp.html",
				"categories/list.amp.html",
				"categories/categories.terms.html",
				"categories/terms.html",
				"categories/taxonomy.html",
				"categories/list.html",
				"taxonomy/categories.terms.amp.html",
				"taxonomy/terms.amp.html",
				"taxonomy/taxonomy.amp.html",
				"taxonomy/list.amp.html",
				"taxonomy/categories.terms.html",
				"taxonomy/terms.html",
				"taxonomy/taxonomy.html",
				"taxonomy/list.html",
				"_default/categories.terms.amp.html",
				"_default/terms.amp.html",
				"_default/taxonomy.amp.html",
				"_default/list.amp.html",
				"_default/categories.terms.html",
				"_default/terms.html",
				"_default/taxonomy.html",
				"_default/list.html",
			},
		},
		{
			"Page", LayoutDescriptor{Kind: "page"}, "", ampType,
			[]string{
				"_default/single.amp.html",
				"_default/single.html",
			},
		},
		{
			"Page, baseof", LayoutDescriptor{Kind: "page", Baseof: true}, "", ampType,
			[]string{
				"_default/single-baseof.amp.html",
				"_default/baseof.amp.html",
				"_default/single-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Page with layout", LayoutDescriptor{Kind: "page", Layout: "mylayout"}, "", ampType,
			[]string{
				"_default/mylayout.amp.html",
				"_default/single.amp.html",
				"_default/mylayout.html",
				"_default/single.html",
			},
		},
		{
			"Page with layout, baseof", LayoutDescriptor{Kind: "page", Layout: "mylayout", Baseof: true}, "", ampType,
			[]string{
				"_default/mylayout-baseof.amp.html",
				"_default/single-baseof.amp.html",
				"_default/baseof.amp.html",
				"_default/mylayout-baseof.html",
				"_default/single-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Page with layout and type", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype"}, "", ampType,
			[]string{
				"myttype/mylayout.amp.html",
				"myttype/single.amp.html",
				"myttype/mylayout.html",
				"myttype/single.html",
				"_default/mylayout.amp.html",
				"_default/single.amp.html",
				"_default/mylayout.html",
				"_default/single.html",
			},
		},
		{
			"Page with layout and type with subtype", LayoutDescriptor{Kind: "page", Layout: "mylayout", Type: "myttype/mysubtype"}, "", ampType,
			[]string{
				"myttype/mysubtype/mylayout.amp.html",
				"myttype/mysubtype/single.amp.html",
				"myttype/mysubtype/mylayout.html",
				"myttype/mysubtype/single.html",
				"_default/mylayout.amp.html",
				"_default/single.amp.html",
				"_default/mylayout.html",
				"_default/single.html",
			},
		},
		// RSS
		{
			"RSS Home", LayoutDescriptor{Kind: "home"}, "", RSSFormat,
			[]string{
				"index.rss.xml",
				"home.rss.xml",
				"rss.xml",
				"list.rss.xml",
				"index.xml",
				"home.xml",
				"list.xml",
				"_default/index.rss.xml",
				"_default/home.rss.xml",
				"_default/rss.xml",
				"_default/list.rss.xml",
				"_default/index.xml",
				"_default/home.xml",
				"_default/list.xml",
				"_internal/_default/rss.xml",
			},
		},
		{
			"RSS Home, baseof", LayoutDescriptor{Kind: "home", Baseof: true}, "", RSSFormat,
			[]string{
				"index-baseof.rss.xml",
				"home-baseof.rss.xml",
				"list-baseof.rss.xml",
				"baseof.rss.xml",
				"index-baseof.xml",
				"home-baseof.xml",
				"list-baseof.xml",
				"baseof.xml",
				"_default/index-baseof.rss.xml",
				"_default/home-baseof.rss.xml",
				"_default/list-baseof.rss.xml",
				"_default/baseof.rss.xml",
				"_default/index-baseof.xml",
				"_default/home-baseof.xml",
				"_default/list-baseof.xml",
				"_default/baseof.xml",
			},
		},
		{
			"RSS Section", LayoutDescriptor{Kind: "section", Section: "sect1"}, "", RSSFormat,
			[]string{
				"sect1/sect1.rss.xml",
				"sect1/section.rss.xml",
				"sect1/rss.xml",
				"sect1/list.rss.xml",
				"sect1/sect1.xml",
				"sect1/section.xml",
				"sect1/list.xml",
				"section/sect1.rss.xml",
				"section/section.rss.xml",
				"section/rss.xml",
				"section/list.rss.xml",
				"section/sect1.xml",
				"section/section.xml",
				"section/list.xml",
				"_default/sect1.rss.xml",
				"_default/section.rss.xml",
				"_default/rss.xml",
				"_default/list.rss.xml",
				"_default/sect1.xml",
				"_default/section.xml",
				"_default/list.xml",
				"_internal/_default/rss.xml",
			},
		},
		{
			"RSS Term", LayoutDescriptor{Kind: "term", Section: "tag"}, "", RSSFormat,
			[]string{
				"term/term.rss.xml",
				"term/tag.rss.xml",
				"term/taxonomy.rss.xml",
				"term/rss.xml",
				"term/list.rss.xml",
				"term/term.xml",
				"term/tag.xml",
				"term/taxonomy.xml",
				"term/list.xml",
				"taxonomy/term.rss.xml",
				"taxonomy/tag.rss.xml",
				"taxonomy/taxonomy.rss.xml",
				"taxonomy/rss.xml",
				"taxonomy/list.rss.xml",
				"taxonomy/term.xml",
				"taxonomy/tag.xml",
				"taxonomy/taxonomy.xml",
				"taxonomy/list.xml",
				"tag/term.rss.xml",
				"tag/tag.rss.xml",
				"tag/taxonomy.rss.xml",
				"tag/rss.xml",
				"tag/list.rss.xml",
				"tag/term.xml",
				"tag/tag.xml",
				"tag/taxonomy.xml",
				"tag/list.xml",
				"_default/term.rss.xml",
				"_default/tag.rss.xml",
				"_default/taxonomy.rss.xml",
				"_default/rss.xml",
				"_default/list.rss.xml",
				"_default/term.xml",
				"_default/tag.xml",
				"_default/taxonomy.xml",
				"_default/list.xml",
				"_internal/_default/rss.xml",
			},
		},
		{
			"RSS Taxonomy", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, "", RSSFormat,
			[]string{
				"tag/tag.terms.rss.xml",
				"tag/terms.rss.xml",
				"tag/taxonomy.rss.xml",
				"tag/rss.xml",
				"tag/list.rss.xml",
				"tag/tag.terms.xml",
				"tag/terms.xml",
				"tag/taxonomy.xml",
				"tag/list.xml",
				"taxonomy/tag.terms.rss.xml",
				"taxonomy/terms.rss.xml",
				"taxonomy/taxonomy.rss.xml",
				"taxonomy/rss.xml",
				"taxonomy/list.rss.xml",
				"taxonomy/tag.terms.xml",
				"taxonomy/terms.xml",
				"taxonomy/taxonomy.xml",
				"taxonomy/list.xml",
				"_default/tag.terms.rss.xml",
				"_default/terms.rss.xml",
				"_default/taxonomy.rss.xml",
				"_default/rss.xml",
				"_default/list.rss.xml",
				"_default/tag.terms.xml",
				"_default/terms.xml",
				"_default/taxonomy.xml",
				"_default/list.xml",
				"_internal/_default/rss.xml",
			},
		},
		{
			"Home plain text", LayoutDescriptor{Kind: "home"}, "", JSONFormat,
			[]string{
				"index.json.json",
				"home.json.json",
				"list.json.json",
				"index.json",
				"home.json",
				"list.json",
				"_default/index.json.json",
				"_default/home.json.json",
				"_default/list.json.json",
				"_default/index.json",
				"_default/home.json",
				"_default/list.json",
			},
		},
		{
			"Page plain text", LayoutDescriptor{Kind: "page"}, "", JSONFormat,
			[]string{
				"_default/single.json.json",
				"_default/single.json",
			},
		},
		{
			"Reserved section, shortcodes", LayoutDescriptor{Kind: "section", Section: "shortcodes", Type: "shortcodes"}, "", ampType,
			[]string{
				"section/shortcodes.amp.html",
				"section/section.amp.html",
				"section/list.amp.html",
				"section/shortcodes.html",
				"section/section.html",
				"section/list.html",
				"_default/shortcodes.amp.html",
				"_default/section.amp.html",
				"_default/list.amp.html",
				"_default/shortcodes.html",
				"_default/section.html",
				"_default/list.html",
			},
		},
		{
			"Reserved section, partials", LayoutDescriptor{Kind: "section", Section: "partials", Type: "partials"}, "", ampType,
			[]string{
				"section/partials.amp.html",
				"section/section.amp.html",
				"section/list.amp.html",
				"section/partials.html",
				"section/section.html",
				"section/list.html",
				"_default/partials.amp.html",
				"_default/section.amp.html",
				"_default/list.amp.html",
				"_default/partials.html",
				"_default/section.html",
				"_default/list.html",
			},
		},
		// This is currently always HTML only
		{
			"404, HTML", LayoutDescriptor{Kind: "404"}, "", htmlFormat,
			[]string{
				"404.html.html",
				"404.html",
			},
		},
		{
			"404, HTML baseof", LayoutDescriptor{Kind: "404", Baseof: true}, "", htmlFormat,
			[]string{
				"404-baseof.html.html",
				"baseof.html.html",
				"404-baseof.html",
				"baseof.html",
				"_default/404-baseof.html.html",
				"_default/baseof.html.html",
				"_default/404-baseof.html",
				"_default/baseof.html",
			},
		},
		{
			"Content hook", LayoutDescriptor{Kind: "render-link", RenderingHook: true, Layout: "mylayout", Section: "blog"}, "", ampType,
			[]string{
				"blog/_markup/render-link.amp.html",
				"blog/_markup/render-link.html",
				"_default/_markup/render-link.amp.html",
				"_default/_markup/render-link.html",
			},
		},
	} {
		c.Run(this.name, func(c *qt.C) {
			l := NewLayoutHandler()

			layouts, err := l.For(this.layoutDescriptor, this.format)

			c.Assert(err, qt.IsNil)
			c.Assert(layouts, qt.Not(qt.IsNil), qt.Commentf(this.layoutDescriptor.Kind))

			if !reflect.DeepEqual(layouts, this.expect) {
				r := strings.NewReplacer(
					"[", "\t\"",
					"]", "\",",
					" ", "\",\n\t\"",
				)
				fmtGot := r.Replace(fmt.Sprintf("%v", layouts))
				fmtExp := r.Replace(fmt.Sprintf("%v", this.expect))

				c.Fatalf("got %d items, expected %d:\nGot:\n\t%v\nExpected:\n\t%v\nDiff:\n%s", len(layouts), len(this.expect), layouts, this.expect, diff.Diff(fmtExp, fmtGot))

			}
		})
	}
}

func BenchmarkLayout(b *testing.B) {
	descriptor := LayoutDescriptor{Kind: "taxonomy", Section: "categories"}
	l := NewLayoutHandler()

	for i := 0; i < b.N; i++ {
		_, err := l.For(descriptor, HTMLFormat)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkLayoutUncached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		descriptor := LayoutDescriptor{Kind: "taxonomy", Section: "categories"}
		l := NewLayoutHandler()

		_, err := l.For(descriptor, HTMLFormat)
		if err != nil {
			panic(err)
		}
	}
}
