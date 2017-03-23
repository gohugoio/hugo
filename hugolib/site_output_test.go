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

package hugolib

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"fmt"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/output"
	"github.com/spf13/viper"
)

func TestDefaultOutputFormats(t *testing.T) {
	t.Parallel()
	defs, err := createDefaultOutputFormats(viper.New())

	require.NoError(t, err)

	tests := []struct {
		name string
		kind string
		want output.Formats
	}{
		{"RSS not for regular pages", KindPage, output.Formats{output.HTMLType}},
		{"Home Sweet Home", KindHome, output.Formats{output.HTMLType, output.RSSType}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defs[tt.kind]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDefaultOutputFormats(%v) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}

func TestSiteWithPageOutputs(t *testing.T) {
	for _, outputs := range [][]string{{"html", "json", "calendar"}, {"json"}} {
		t.Run(fmt.Sprintf("%v", outputs), func(t *testing.T) {
			doTestSiteWithPageOutputs(t, outputs)
		})
	}
}

// TODO(bep) output add test for site outputs config
func doTestSiteWithPageOutputs(t *testing.T, outputs []string) {
	t.Parallel()

	outputsStr := strings.Replace(fmt.Sprintf("%q", outputs), " ", ", ", -1)

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "en"

disableKinds = ["page", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"]

[Taxonomies]
tag = "tags"
category = "categories"
`

	pageTemplate := `---
title: "%s"
outputs: %s
---
# Doc
`

	th, h := newTestSitesFromConfig(t, siteConfig,
		"layouts/_default/list.json", `List JSON|{{ .Title }}|{{ .Content }}|Alt formats: {{ len .AlternativeOutputFormats -}}|
{{- range .AlternativeOutputFormats -}}
Alt Output: {{ .Name -}}|
{{- end -}}|
{{- range .OutputFormats -}}
Output/Rel: {{ .Name -}}/{{ .Rel }}|
{{- end -}}
`,
	)
	require.Len(t, h.Sites, 1)

	fs := th.Fs

	writeSource(t, fs, "content/_index.md", fmt.Sprintf(pageTemplate, "JSON Home", outputsStr))

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	s := h.Sites[0]
	home := s.getPage(KindHome)

	require.NotNil(t, home)

	lenOut := len(outputs)

	require.Len(t, home.outputFormats, lenOut)

	// TODO(bep) output assert template/text
	// There is currently always a JSON output to make it simpler ...
	altFormats := lenOut - 1
	hasHTML := helpers.InStringArray(outputs, "html")
	th.assertFileContent("public/index.json",
		"List JSON",
		fmt.Sprintf("Alt formats: %d", altFormats),
	)

	if hasHTML {
		th.assertFileContent("public/index.json",
			"Alt Output: HTML",
			"Output/Rel: JSON/alternate|",
			"Output/Rel: HTML/canonical|",
		)
	} else {
		th.assertFileContent("public/index.json",
			"Output/Rel: JSON/canonical|",
		)
	}

	of := home.OutputFormats()
	require.Len(t, of, lenOut)
	require.Nil(t, of.Get("Hugo"))
	require.NotNil(t, of.Get("json"))
	json := of.Get("JSON")
	_, err = home.AlternativeOutputFormats()
	require.Error(t, err)
	require.NotNil(t, json)
	require.Equal(t, "/blog/index.json", json.RelPermalink())
	require.Equal(t, "http://example.com/blog/index.json", json.Permalink())

	if helpers.InStringArray(outputs, "cal") {
		// TODO(bep) output have do some protocil handling for the default too if set.
		cal := of.Get("calendar")
		require.NotNil(t, cal)
		require.Equal(t, "/blog/index.ics", cal.RelPermalink())
		require.Equal(t, "webcal://example.com/blog/index.ics", cal.Permalink())
	}

}
