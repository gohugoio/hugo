// Copyright 2017 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	"fmt"

	"github.com/gohugoio/hugo/output"
)

func TestPageTargetPath(t *testing.T) {

	pathSpec := newTestDefaultPathSpec(t)

	noExtNoDelimMediaType := media.TextType
	noExtNoDelimMediaType.Suffixes = []string{}
	noExtNoDelimMediaType.Delimiter = ""

	// Netlify style _redirects
	noExtDelimFormat := output.Format{
		Name:      "NER",
		MediaType: noExtNoDelimMediaType,
		BaseName:  "_redirects",
	}

	for _, multiHost := range []bool{false, true} {
		for _, langPrefix := range []string{"", "no"} {
			for _, uglyURLs := range []bool{false, true} {
				t.Run(fmt.Sprintf("multihost=%t,langPrefix=%q,uglyURLs=%t", multiHost, langPrefix, uglyURLs),
					func(t *testing.T) {

						tests := []struct {
							name     string
							d        targetPathDescriptor
							expected string
						}{
							{"JSON home", targetPathDescriptor{Kind: KindHome, Type: output.JSONFormat}, "/index.json"},
							{"AMP home", targetPathDescriptor{Kind: KindHome, Type: output.AMPFormat}, "/amp/index.html"},
							{"HTML home", targetPathDescriptor{Kind: KindHome, BaseName: "_index", Type: output.HTMLFormat}, "/index.html"},
							{"Netlify redirects", targetPathDescriptor{Kind: KindHome, BaseName: "_index", Type: noExtDelimFormat}, "/_redirects"},
							{"HTML section list", targetPathDescriptor{
								Kind:     KindSection,
								Sections: []string{"sect1"},
								BaseName: "_index",
								Type:     output.HTMLFormat}, "/sect1/index.html"},
							{"HTML taxonomy list", targetPathDescriptor{
								Kind:     KindTaxonomy,
								Sections: []string{"tags", "hugo"},
								BaseName: "_index",
								Type:     output.HTMLFormat}, "/tags/hugo/index.html"},
							{"HTML taxonomy term", targetPathDescriptor{
								Kind:     KindTaxonomy,
								Sections: []string{"tags"},
								BaseName: "_index",
								Type:     output.HTMLFormat}, "/tags/index.html"},
							{
								"HTML page", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/a/b",
									BaseName: "mypage",
									Sections: []string{"a"},
									Type:     output.HTMLFormat}, "/a/b/mypage/index.html"},

							{
								"HTML page with index as base", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/a/b",
									BaseName: "index",
									Sections: []string{"a"},
									Type:     output.HTMLFormat}, "/a/b/index.html"},

							{
								"HTML page with special chars", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/a/b",
									BaseName: "My Page!",
									Type:     output.HTMLFormat}, "/a/b/My-Page/index.html"},
							{"RSS home", targetPathDescriptor{Kind: kindRSS, Type: output.RSSFormat}, "/index.xml"},
							{"RSS section list", targetPathDescriptor{
								Kind:     kindRSS,
								Sections: []string{"sect1"},
								Type:     output.RSSFormat}, "/sect1/index.xml"},
							{
								"AMP page", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/a/b/c",
									BaseName: "myamp",
									Type:     output.AMPFormat}, "/amp/a/b/c/myamp/index.html"},
							{
								"AMP page with URL with suffix", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/sect/",
									BaseName: "mypage",
									URL:      "/some/other/url.xhtml",
									Type:     output.HTMLFormat}, "/some/other/url.xhtml"},
							{
								"JSON page with URL without suffix", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/sect/",
									BaseName: "mypage",
									URL:      "/some/other/path/",
									Type:     output.JSONFormat}, "/some/other/path/index.json"},
							{
								"JSON page with URL without suffix and no trailing slash", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/sect/",
									BaseName: "mypage",
									URL:      "/some/other/path",
									Type:     output.JSONFormat}, "/some/other/path/index.json"},
							{
								"HTML page with expanded permalink", targetPathDescriptor{
									Kind:              KindPage,
									Dir:               "/a/b",
									BaseName:          "mypage",
									ExpandedPermalink: "/2017/10/my-title",
									Type:              output.HTMLFormat}, "/2017/10/my-title/index.html"},
							{
								"Paginated HTML home", targetPathDescriptor{
									Kind:     KindHome,
									BaseName: "_index",
									Type:     output.HTMLFormat,
									Addends:  "page/3"}, "/page/3/index.html"},
							{
								"Paginated Taxonomy list", targetPathDescriptor{
									Kind:     KindTaxonomy,
									BaseName: "_index",
									Sections: []string{"tags", "hugo"},
									Type:     output.HTMLFormat,
									Addends:  "page/3"}, "/tags/hugo/page/3/index.html"},
							{
								"Regular page with addend", targetPathDescriptor{
									Kind:     KindPage,
									Dir:      "/a/b",
									BaseName: "mypage",
									Addends:  "c/d/e",
									Type:     output.HTMLFormat}, "/a/b/mypage/c/d/e/index.html"},
						}

						for i, test := range tests {
							test.d.PathSpec = pathSpec
							test.d.UglyURLs = uglyURLs
							test.d.LangPrefix = langPrefix
							test.d.IsMultihost = multiHost
							test.d.Dir = filepath.FromSlash(test.d.Dir)
							isUgly := uglyURLs && !test.d.Type.NoUgly

							expected := test.expected

							// TODO(bep) simplify
							if test.d.Kind == KindPage && test.d.BaseName == test.d.Type.BaseName {

							} else if test.d.Kind == KindHome && test.d.Type.Path != "" {
							} else if (!strings.HasPrefix(expected, "/index") || test.d.Addends != "") && test.d.URL == "" && isUgly {
								expected = strings.Replace(expected,
									"/"+test.d.Type.BaseName+"."+test.d.Type.MediaType.Suffix(),
									"."+test.d.Type.MediaType.Suffix(), -1)
							}

							if test.d.LangPrefix != "" && !(test.d.Kind == KindPage && test.d.URL != "") {
								expected = "/" + test.d.LangPrefix + expected
							} else if multiHost && test.d.LangPrefix != "" && test.d.URL != "" {
								expected = "/" + test.d.LangPrefix + expected
							}

							expected = filepath.FromSlash(expected)

							pagePath := createTargetPath(test.d)

							if pagePath != expected {
								t.Fatalf("[%d] [%s] targetPath expected %q, got: %q", i, test.name, expected, pagePath)
							}
						}
					})
			}
		}
	}
}
