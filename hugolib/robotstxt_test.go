// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"fmt"
	"strings"
	"testing"
)

func TestRobotsTXTOutput(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://auth/bub/"
enableRobotsTXT = true
-- layouts/robots.txt --
User-agent: Googlebot
  {{ range .Data.Pages }}
	Disallow: {{.RelPermalink}}
	{{ end }}
`
	b := Test(t, files)

	b.AssertFileContent("public/robots.txt", "User-agent: Googlebot", "! Sitemap:")
}

func TestRobotsTXTDefaultTemplate(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "http://auth/bub/"
enableRobotsTXT = true
`

	b := Test(t, files)

	b.AssertFileContentEquals("public/robots.txt", "User-agent: *\nSitemap: http://auth/bub/sitemap.xml\n")
}

// See issue 4678.
func TestRobotsTXTSitemap(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name        string
		config      string
		contentPage string
		robots      map[string]string
		sitemapKind string
	}{
		{
			name:        "monolingual",
			robots:      map[string]string{"public/robots.txt": "https://example.org/docs/"},
			sitemapKind: "urlset",
		},
		{
			name:        "sitemap name collision",
			contentPage: "_sitemap",
			robots:      map[string]string{"public/robots.txt": ""},
		},
		{
			name:        "sitemap index name collision",
			contentPage: "_sitemapindex",
			robots:      map[string]string{"public/robots.txt": "https://example.org/docs/"},
			sitemapKind: "urlset",
		},
		{
			name:        "monolingual subdir",
			config:      "defaultContentLanguageInSubdir = true",
			robots:      map[string]string{"public/robots.txt": "https://example.org/docs/"},
			sitemapKind: "sitemapindex",
		},
		{
			name: "multilingual",
			config: `
[languages.en]
weight = 1
[languages.fr]
weight = 2
`,
			robots:      map[string]string{"public/robots.txt": "https://example.org/docs/"},
			sitemapKind: "sitemapindex",
		},
		{
			name: "multilingual subdir",
			config: `
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.fr]
weight = 2
`,
			robots:      map[string]string{"public/robots.txt": "https://example.org/docs/"},
			sitemapKind: "sitemapindex",
		},
		{
			name: "multihost",
			config: `
[languages.en]
weight = 1
baseURL = 'https://en.example.org/docs/'
[languages.fr]
weight = 2
baseURL = 'https://fr.example.org/guide/'
`,
			robots: map[string]string{
				"public/en/robots.txt": "https://en.example.org/docs/",
				"public/fr/robots.txt": "https://fr.example.org/guide/",
			},
			sitemapKind: "urlset",
		},
		{
			name: "multihost disabled language",
			config: `
[languages.en]
weight = 1
baseURL = 'https://en.example.org/docs/'
[languages.fr]
weight = 2
baseURL = 'https://fr.example.org/guide/'
disableKinds = ['sitemap']
`,
			robots: map[string]string{
				"public/en/robots.txt": "https://en.example.org/docs/",
				"public/fr/robots.txt": "",
			},
			sitemapKind: "urlset",
		},
	} {
		for _, filename := range []string{"sitemap.xml", "custom.xml"} {
			for _, disabled := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/disabled=%t", test.name, filename, disabled), func(t *testing.T) {
					t.Parallel()
					disableSitemap := ""
					if disabled {
						disableSitemap = ", 'sitemap'"
					}
					files := fmt.Sprintf(`
-- hugo.toml --
baseURL = 'https://example.org/docs/'
enableRobotsTXT = true
relativeURLs = true
disableKinds = ['taxonomy', 'term', 'rss'%s]
%s
[sitemap]
filename = '%s'
-- layouts/home.html --
{{ .Title }}
-- layouts/page.html --
{{ .Title }}
`, disableSitemap, test.config, filename)
					if test.contentPage != "" {
						files += fmt.Sprintf("\n-- content/%s.md --\n---\ntitle: Example\n---\n", test.contentPage)
					}

					b := Test(t, files)
					if test.contentPage != "" {
						b.AssertFileExists("public/"+test.contentPage+"/index.html", true)
					}
					for robotsFilename, baseURL := range test.robots {
						enabled := !disabled && baseURL != ""
						want := "User-agent: *\n"
						if enabled {
							want += "Sitemap: " + baseURL + filename + "\n"
						}
						b.AssertFileContentEquals(robotsFilename, want)
						sitemapFilename := strings.TrimSuffix(robotsFilename, "robots.txt") + filename
						b.AssertFileExists(sitemapFilename, enabled)
						if enabled {
							b.AssertFileContent(sitemapFilename, "<"+test.sitemapKind)
						}
					}
				})
			}
		}
	}
}
