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

package page

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	"fmt"

	"github.com/gohugoio/hugo/output"
)

func TestPageTargetPath(t *testing.T) {

	pathSpec := newTestPathSpec()

	noExtNoDelimMediaType := media.TextType
	noExtNoDelimMediaType.Suffixes = []string{}
	noExtNoDelimMediaType.Delimiter = ""

	// Netlify style _redirects
	noExtDelimFormat := output.Format{
		Name:      "NER",
		MediaType: noExtNoDelimMediaType,
		BaseName:  "_redirects",
	}

	for _, langPrefixPath := range []string{"", "no"} {
		for _, langPrefixLink := range []string{"", "no"} {
			for _, uglyURLs := range []bool{false, true} {

				tests := []struct {
					name     string
					d        TargetPathDescriptor
					expected TargetPaths
				}{
					{"JSON home", TargetPathDescriptor{Kind: KindHome, Type: output.JSONFormat}, TargetPaths{TargetFilename: "/index.json", SubResourceBaseTarget: "", Link: "/index.json"}},
					{"AMP home", TargetPathDescriptor{Kind: KindHome, Type: output.AMPFormat}, TargetPaths{TargetFilename: "/amp/index.html", SubResourceBaseTarget: "/amp", Link: "/amp/"}},
					{"HTML home", TargetPathDescriptor{Kind: KindHome, BaseName: "_index", Type: output.HTMLFormat}, TargetPaths{TargetFilename: "/index.html", SubResourceBaseTarget: "", Link: "/"}},
					{"Netlify redirects", TargetPathDescriptor{Kind: KindHome, BaseName: "_index", Type: noExtDelimFormat}, TargetPaths{TargetFilename: "/_redirects", SubResourceBaseTarget: "", Link: "/_redirects"}},
					{"HTML section list", TargetPathDescriptor{
						Kind:     KindSection,
						Sections: []string{"sect1"},
						BaseName: "_index",
						Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/sect1/index.html", SubResourceBaseTarget: "/sect1", Link: "/sect1/"}},
					{"HTML taxonomy list", TargetPathDescriptor{
						Kind:     KindTaxonomy,
						Sections: []string{"tags", "hugo"},
						BaseName: "_index",
						Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/tags/hugo/index.html", SubResourceBaseTarget: "/tags/hugo", Link: "/tags/hugo/"}},
					{"HTML taxonomy term", TargetPathDescriptor{
						Kind:     KindTaxonomy,
						Sections: []string{"tags"},
						BaseName: "_index",
						Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/tags/index.html", SubResourceBaseTarget: "/tags", Link: "/tags/"}},
					{
						"HTML page", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/a/b",
							BaseName: "mypage",
							Sections: []string{"a"},
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/a/b/mypage/index.html", SubResourceBaseTarget: "/a/b/mypage", Link: "/a/b/mypage/"}},

					{
						"HTML page with index as base", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/a/b",
							BaseName: "index",
							Sections: []string{"a"},
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/a/b/index.html", SubResourceBaseTarget: "/a/b", Link: "/a/b/"}},

					{
						"HTML page with special chars", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/a/b",
							BaseName: "My Page!",
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/a/b/my-page/index.html", SubResourceBaseTarget: "/a/b/my-page", Link: "/a/b/my-page/"}},
					{"RSS home", TargetPathDescriptor{Kind: "rss", Type: output.RSSFormat}, TargetPaths{TargetFilename: "/index.xml", SubResourceBaseTarget: "", Link: "/index.xml"}},
					{"RSS section list", TargetPathDescriptor{
						Kind:     "rss",
						Sections: []string{"sect1"},
						Type:     output.RSSFormat}, TargetPaths{TargetFilename: "/sect1/index.xml", SubResourceBaseTarget: "/sect1", Link: "/sect1/index.xml"}},
					{
						"AMP page", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/a/b/c",
							BaseName: "myamp",
							Type:     output.AMPFormat}, TargetPaths{TargetFilename: "/amp/a/b/c/myamp/index.html", SubResourceBaseTarget: "/amp/a/b/c/myamp", Link: "/amp/a/b/c/myamp/"}},
					{
						"AMP page with URL with suffix", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/url.xhtml",
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/some/other/url.xhtml", SubResourceBaseTarget: "/some/other", Link: "/some/other/url.xhtml"}},
					{
						"JSON page with URL without suffix", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path/",
							Type:     output.JSONFormat}, TargetPaths{TargetFilename: "/some/other/path/index.json", SubResourceBaseTarget: "/some/other/path", Link: "/some/other/path/index.json"}},
					{
						"JSON page with URL without suffix and no trailing slash", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path",
							Type:     output.JSONFormat}, TargetPaths{TargetFilename: "/some/other/path/index.json", SubResourceBaseTarget: "/some/other/path", Link: "/some/other/path/index.json"}},
					{
						"HTML page with URL without suffix and no trailing slash", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path",
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/some/other/path/index.html", SubResourceBaseTarget: "/some/other/path", Link: "/some/other/path/"}},
					{
						"HTML page with expanded permalink", TargetPathDescriptor{
							Kind:              KindPage,
							Dir:               "/a/b",
							BaseName:          "mypage",
							ExpandedPermalink: "/2017/10/my-title/",
							Type:              output.HTMLFormat}, TargetPaths{TargetFilename: "/2017/10/my-title/index.html", SubResourceBaseTarget: "/2017/10/my-title", Link: "/2017/10/my-title/"}},
					{
						"Paginated HTML home", TargetPathDescriptor{
							Kind:     KindHome,
							BaseName: "_index",
							Type:     output.HTMLFormat,
							Addends:  "page/3"}, TargetPaths{TargetFilename: "/page/3/index.html", SubResourceBaseTarget: "/page/3", Link: "/page/3/"}},
					{
						"Paginated Taxonomy list", TargetPathDescriptor{
							Kind:     KindTaxonomy,
							BaseName: "_index",
							Sections: []string{"tags", "hugo"},
							Type:     output.HTMLFormat,
							Addends:  "page/3"}, TargetPaths{TargetFilename: "/tags/hugo/page/3/index.html", SubResourceBaseTarget: "/tags/hugo/page/3", Link: "/tags/hugo/page/3/"}},
					{
						"Regular page with addend", TargetPathDescriptor{
							Kind:     KindPage,
							Dir:      "/a/b",
							BaseName: "mypage",
							Addends:  "c/d/e",
							Type:     output.HTMLFormat}, TargetPaths{TargetFilename: "/a/b/mypage/c/d/e/index.html", SubResourceBaseTarget: "/a/b/mypage/c/d/e", Link: "/a/b/mypage/c/d/e/"}},
				}

				for i, test := range tests {
					t.Run(fmt.Sprintf("langPrefixPath=%s,langPrefixLink=%s,uglyURLs=%t,name=%s", langPrefixPath, langPrefixLink, uglyURLs, test.name),
						func(t *testing.T) {

							test.d.ForcePrefix = true
							test.d.PathSpec = pathSpec
							test.d.UglyURLs = uglyURLs
							test.d.PrefixFilePath = langPrefixPath
							test.d.PrefixLink = langPrefixLink
							test.d.Dir = filepath.FromSlash(test.d.Dir)
							isUgly := uglyURLs && !test.d.Type.NoUgly

							expected := test.expected

							// TODO(bep) simplify
							if test.d.Kind == KindPage && test.d.BaseName == test.d.Type.BaseName {
							} else if test.d.Kind == KindHome && test.d.Type.Path != "" {
							} else if test.d.Type.MediaType.Suffix() != "" && (!strings.HasPrefix(expected.TargetFilename, "/index") || test.d.Addends != "") && test.d.URL == "" && isUgly {
								expected.TargetFilename = strings.Replace(expected.TargetFilename,
									"/"+test.d.Type.BaseName+"."+test.d.Type.MediaType.Suffix(),
									"."+test.d.Type.MediaType.Suffix(), 1)
								expected.Link = strings.TrimSuffix(expected.Link, "/") + "." + test.d.Type.MediaType.Suffix()

							}

							if test.d.PrefixFilePath != "" && !strings.HasPrefix(test.d.URL, "/"+test.d.PrefixFilePath) {
								expected.TargetFilename = "/" + test.d.PrefixFilePath + expected.TargetFilename
								expected.SubResourceBaseTarget = "/" + test.d.PrefixFilePath + expected.SubResourceBaseTarget
							}

							if test.d.PrefixLink != "" && !strings.HasPrefix(test.d.URL, "/"+test.d.PrefixLink) {
								expected.Link = "/" + test.d.PrefixLink + expected.Link
							}

							expected.TargetFilename = filepath.FromSlash(expected.TargetFilename)
							expected.SubResourceBaseTarget = filepath.FromSlash(expected.SubResourceBaseTarget)

							pagePath := CreateTargetPaths(test.d)

							if !eqTargetPaths(pagePath, expected) {
								t.Fatalf("[%d] [%s] targetPath expected\n%#v, got:\n%#v", i, test.name, expected, pagePath)

							}
						})
				}
			}

		}
	}
}

func TestPageTargetPathPrefix(t *testing.T) {
	pathSpec := newTestPathSpec()
	tests := []struct {
		name     string
		d        TargetPathDescriptor
		expected TargetPaths
	}{
		{"URL set, prefix both, no force", TargetPathDescriptor{Kind: KindPage, Type: output.JSONFormat, URL: "/mydir/my.json", ForcePrefix: false, PrefixFilePath: "pf", PrefixLink: "pl"},
			TargetPaths{TargetFilename: "/mydir/my.json", SubResourceBaseTarget: "/mydir", SubResourceBaseLink: "/mydir", Link: "/mydir/my.json"}},
		{"URL set, prefix both, force", TargetPathDescriptor{Kind: KindPage, Type: output.JSONFormat, URL: "/mydir/my.json", ForcePrefix: true, PrefixFilePath: "pf", PrefixLink: "pl"},
			TargetPaths{TargetFilename: "/pf/mydir/my.json", SubResourceBaseTarget: "/pf/mydir", SubResourceBaseLink: "/pl/mydir", Link: "/pl/mydir/my.json"}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf(test.name),
			func(t *testing.T) {
				test.d.PathSpec = pathSpec
				expected := test.expected
				expected.TargetFilename = filepath.FromSlash(expected.TargetFilename)
				expected.SubResourceBaseTarget = filepath.FromSlash(expected.SubResourceBaseTarget)

				pagePath := CreateTargetPaths(test.d)

				if pagePath != expected {
					t.Fatalf("[%d] [%s] targetPath expected\n%#v, got:\n%#v", i, test.name, expected, pagePath)
				}
			})
	}

}

func eqTargetPaths(p1, p2 TargetPaths) bool {

	if p1.Link != p2.Link {
		return false
	}

	if p1.SubResourceBaseTarget != p2.SubResourceBaseTarget {
		return false
	}

	if p1.TargetFilename != p2.TargetFilename {
		return false
	}

	return true
}
