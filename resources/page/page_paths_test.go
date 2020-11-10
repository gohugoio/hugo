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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/resources/page/pagekinds"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/output"
)

func TestPageTargetPath(t *testing.T) {
	pathSpec := newTestPathSpec()

	noExtNoDelimMediaType := media.WithDelimiterAndSuffixes(media.TextType, "", "")
	noExtNoDelimMediaType.Delimiter = ""

	// Netlify style _redirects
	noExtDelimFormat := output.Format{
		Name:      "NER",
		MediaType: noExtNoDelimMediaType,
		BaseName:  "_redirects",
	}

	htmlCustomBaseName := output.HTMLFormat
	htmlCustomBaseName.BaseName = "cindex"

	type variant struct {
		langPrefixPath string
		langPrefixLink string
		isUgly         bool
	}

	applyPathPrefixes := func(v variant, tp *TargetPaths) {
		if v.langPrefixLink != "" {
			tp.Link = fmt.Sprintf("/%s%s", v.langPrefixLink, tp.Link)
			if tp.SubResourceBaseLink != "" {
				tp.SubResourceBaseLink = fmt.Sprintf("/%s%s", v.langPrefixLink, tp.SubResourceBaseLink)
			}
		}
		if v.langPrefixPath != "" {
			tp.TargetFilename = fmt.Sprintf("/%s%s", v.langPrefixPath, tp.TargetFilename)
			if tp.SubResourceBaseTarget != "" {
				tp.SubResourceBaseTarget = fmt.Sprintf("/%s%s", v.langPrefixPath, tp.SubResourceBaseTarget)
			}
		}
	}

	for _, langPrefixPath := range []string{"", "no"} {
		for _, langPrefixLink := range []string{"", "no"} {
			for _, uglyURLs := range []bool{false, true} {
				tests := []struct {
					name         string
					d            TargetPathDescriptor
					expectedFunc func(v variant) (TargetPaths, bool)
				}{
					{
						"JSON home",
						TargetPathDescriptor{Kind: pagekinds.Home, Type: output.JSONFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/index.json", Link: "/index.json"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"AMP home",
						TargetPathDescriptor{Kind: pagekinds.Home, Type: output.AMPFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/amp/index.html", SubResourceBaseTarget: "/amp", Link: "/amp/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML home",
						TargetPathDescriptor{Kind: pagekinds.Home, BaseName: "_index", Type: output.HTMLFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/index.html", Link: "/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"Netlify redirects",
						TargetPathDescriptor{Kind: pagekinds.Home, BaseName: "_index", Type: noExtDelimFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/_redirects", Link: "/_redirects"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML section list", TargetPathDescriptor{
							Kind:     pagekinds.Section,
							Sections: []string{"sect1"},
							BaseName: "_index",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/sect1.html", SubResourceBaseTarget: "/sect1", Link: "/sect1.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/sect1/index.html", SubResourceBaseTarget: "/sect1", Link: "/sect1/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML taxonomy term", TargetPathDescriptor{
							Kind:     pagekinds.Term,
							Sections: []string{"tags", "hugo"},
							BaseName: "_index",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/tags/hugo.html", SubResourceBaseTarget: "/tags/hugo", Link: "/tags/hugo.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/tags/hugo/index.html", SubResourceBaseTarget: "/tags/hugo", Link: "/tags/hugo/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML taxonomy", TargetPathDescriptor{
							Kind:     pagekinds.Taxonomy,
							Sections: []string{"tags"},
							BaseName: "_index",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/tags.html", SubResourceBaseTarget: "/tags", Link: "/tags.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/tags/index.html", SubResourceBaseTarget: "/tags", Link: "/tags/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b",
							BaseName: "mypage",
							Sections: []string{"a"},
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/a/b/mypage.html", SubResourceBaseTarget: "/a/b/mypage", Link: "/a/b/mypage.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/a/b/mypage/index.html", SubResourceBaseTarget: "/a/b/mypage", Link: "/a/b/mypage/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page, custom base", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b/mypage",
							Sections: []string{"a"},
							Type:     htmlCustomBaseName,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/a/b/mypage.html", SubResourceBaseTarget: "/a/b/mypage", Link: "/a/b/mypage.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/a/b/mypage/cindex.html", SubResourceBaseTarget: "/a/b/mypage", Link: "/a/b/mypage/cindex.html"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page with index as base", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b",
							BaseName: "index",
							Sections: []string{"a"},
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/a/b/index.html", SubResourceBaseTarget: "/a/b", Link: "/a/b/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},

					{
						"HTML page with special chars", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b",
							BaseName: "My Page!",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/a/b/my-page.html", SubResourceBaseTarget: "/a/b/my-page", Link: "/a/b/my-page.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/a/b/my-page/index.html", SubResourceBaseTarget: "/a/b/my-page", Link: "/a/b/my-page/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"RSS home", TargetPathDescriptor{Kind: "rss", Type: output.RSSFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/index.xml", SubResourceBaseTarget: "", Link: "/index.xml"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"RSS section list", TargetPathDescriptor{
							Kind:     "rss",
							Sections: []string{"sect1"},
							Type:     output.RSSFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/sect1/index.xml", SubResourceBaseTarget: "/sect1", Link: "/sect1/index.xml"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"AMP page", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b/c",
							BaseName: "myamp",
							Type:     output.AMPFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/amp/a/b/c/myamp.html", SubResourceBaseTarget: "/amp/a/b/c/myamp", Link: "/amp/a/b/c/myamp.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/amp/a/b/c/myamp/index.html", SubResourceBaseTarget: "/amp/a/b/c/myamp", Link: "/amp/a/b/c/myamp/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"AMP page with URL with suffix", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/url.xhtml",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/some/other/url.xhtml", SubResourceBaseTarget: "/some/other/url", Link: "/some/other/url.xhtml"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"JSON page with URL without suffix", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path/",
							Type:     output.JSONFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/some/other/path/index.json", Link: "/some/other/path/index.json"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"JSON page with URL without suffix and no trailing slash", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path",
							Type:     output.JSONFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/some/other/path/index.json", Link: "/some/other/path/index.json"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page with URL without suffix and no trailing slash", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other/path",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/some/other/path/index.html", SubResourceBaseTarget: "/some/other/path", Link: "/some/other/path/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page with URL containing double hyphen", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/sect/",
							BaseName: "mypage",
							URL:      "/some/other--url/",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/some/other--url/index.html", SubResourceBaseTarget: "/some/other--url", Link: "/some/other--url/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page with URL with lots of dots", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							BaseName: "mypage",
							URL:      "../../../../../myblog/p2/",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/myblog/p2/index.html", SubResourceBaseTarget: "/myblog/p2", Link: "/myblog/p2/"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"HTML page with expanded permalink", TargetPathDescriptor{
							Kind:              pagekinds.Page,
							Dir:               "/a/b",
							BaseName:          "mypage",
							ExpandedPermalink: "/2017/10/my-title/",
							Type:              output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/2017/10/my-title.html", SubResourceBaseTarget: "/2017/10/my-title", SubResourceBaseLink: "/2017/10/my-title", Link: "/2017/10/my-title.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/2017/10/my-title/index.html", SubResourceBaseTarget: "/2017/10/my-title", SubResourceBaseLink: "/2017/10/my-title", Link: "/2017/10/my-title/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"Paginated HTML home", TargetPathDescriptor{
							Kind:     pagekinds.Home,
							BaseName: "_index",
							Type:     output.HTMLFormat,
							Addends:  "page/3",
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/page/3.html", SubResourceBaseTarget: "/page/3", Link: "/page/3.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/page/3/index.html", SubResourceBaseTarget: "/page/3", Link: "/page/3/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"Paginated Taxonomy terms list", TargetPathDescriptor{
							Kind:     pagekinds.Term,
							BaseName: "_index",
							Sections: []string{"tags", "hugo"},
							Type:     output.HTMLFormat,
							Addends:  "page/3",
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/tags/hugo/page/3.html", Link: "/tags/hugo/page/3.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/tags/hugo/page/3/index.html", Link: "/tags/hugo/page/3/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"Regular page with addend", TargetPathDescriptor{
							Kind:     pagekinds.Page,
							Dir:      "/a/b",
							BaseName: "mypage",
							Addends:  "c/d/e",
							Type:     output.HTMLFormat,
						},
						func(v variant) (expected TargetPaths, skip bool) {
							if v.isUgly {
								expected = TargetPaths{TargetFilename: "/a/b/mypage/c/d/e.html", SubResourceBaseTarget: "/a/b/mypage/c/d/e", Link: "/a/b/mypage/c/d/e.html"}
							} else {
								expected = TargetPaths{TargetFilename: "/a/b/mypage/c/d/e/index.html", SubResourceBaseTarget: "/a/b/mypage/c/d/e", Link: "/a/b/mypage/c/d/e/"}
							}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{
						"404", TargetPathDescriptor{Kind: pagekinds.Status404, Type: output.HTTPStatusHTMLFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/404.html", SubResourceBaseTarget: "", Link: "/404.html"}
							applyPathPrefixes(v, &expected)
							return
						},
					},
					{"robots.txt", TargetPathDescriptor{Kind: pagekinds.RobotsTXT, Type: output.RobotsTxtFormat},
						func(v variant) (expected TargetPaths, skip bool) {
							expected = TargetPaths{TargetFilename: "/robots.txt", SubResourceBaseTarget: "", Link: "/robots.txt"}
							return
						},
					},
				}

				for i, test := range tests {
					t.Run(fmt.Sprintf("langPrefixPath=%s,langPrefixLink=%s,uglyURLs=%t,name=%s", langPrefixPath, langPrefixLink, uglyURLs, test.name),
						func(t *testing.T) {
							test.d.ForcePrefix = true
							test.d.PathSpec = pathSpec
							test.d.UglyURLs = uglyURLs
							if !test.d.Type.Root {
								test.d.PrefixFilePath = langPrefixPath
								test.d.PrefixLink = langPrefixLink
							}
							test.d.Dir = filepath.FromSlash(test.d.Dir)
							isUgly := test.d.Type.Ugly || (uglyURLs && !test.d.Type.NoUgly)

							v := variant{
								langPrefixLink: langPrefixLink,
								langPrefixPath: langPrefixPath,
								isUgly:         isUgly,
							}

							expected, skip := test.expectedFunc(v)
							if skip {
								return
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
		{
			"URL set, prefix both, no force",
			TargetPathDescriptor{Kind: pagekinds.Page, Type: output.JSONFormat, URL: "/mydir/my.json", ForcePrefix: false, PrefixFilePath: "pf", PrefixLink: "pl"},
			TargetPaths{TargetFilename: "/mydir/my.json", SubResourceBaseTarget: "/mydir/my", SubResourceBaseLink: "/mydir/my", Link: "/mydir/my.json"},
		},
		{
			"URL set, prefix both, force",
			TargetPathDescriptor{Kind: pagekinds.Page, Type: output.JSONFormat, URL: "/mydir/my.json", ForcePrefix: true, PrefixFilePath: "pf", PrefixLink: "pl"},
			TargetPaths{TargetFilename: "/pf/mydir/my.json", SubResourceBaseTarget: "/pf/mydir/my", SubResourceBaseLink: "/pl/mydir/my", Link: "/pl/mydir/my.json"},
		},
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

func BenchmarkCreateTargetPaths(b *testing.B) {
	pathSpec := newTestPathSpec()
	descriptors := []TargetPathDescriptor{
		{Kind: pagekinds.Home, Type: output.JSONFormat, PathSpec: pathSpec},
		{Kind: pagekinds.Home, Type: output.HTMLFormat, PathSpec: pathSpec},
		{Kind: pagekinds.Section, Type: output.HTMLFormat, Sections: []string{"a", "b", "c"}, PathSpec: pathSpec},
		{Kind: pagekinds.Page, Dir: "/sect/", Type: output.HTMLFormat, PathSpec: pathSpec},
		{Kind: pagekinds.Page, ExpandedPermalink: "/foo/bar/", UglyURLs: true, Type: output.HTMLFormat, PathSpec: pathSpec},
		{Kind: pagekinds.Page, URL: "/sect/foo.html", Type: output.HTMLFormat, PathSpec: pathSpec},
		{Kind: pagekinds.Status404, Type: output.HTTPStatusHTMLFormat, PathSpec: pathSpec},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, d := range descriptors {
			_ = CreateTargetPaths(d)
		}
	}
}

func eqTargetPaths(got, expected TargetPaths) bool {
	if got.Link != expected.Link {
		return false
	}

	// Be a little lenient with these sub resource paths as it's not filled in in all cases.
	if expected.SubResourceBaseLink != "" && got.SubResourceBaseLink != expected.SubResourceBaseLink {
		return false
	}

	if expected.SubResourceBaseTarget != "" && got.SubResourceBaseTarget != expected.SubResourceBaseTarget {
		return false
	}

	if got.TargetFilename != expected.TargetFilename {
		return false
	}

	return true
}
