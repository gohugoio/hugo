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
	"testing"

	"fmt"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

const (
	menuPageTemplate = `---
title: %q
weight: %d
menu:
  %s:
    weight: %d
---
# Doc Menu
`
)

func TestSectionPagesMenu(t *testing.T) {
	t.Parallel()

	siteConfig := `
baseurl = "http://example.com/"
title = "Section Menu"
sectionPagesMenu = "sect"
`

	th, h := newTestSitesFromConfig(t, afero.NewMemMapFs(), siteConfig,
		"layouts/partials/menu.html", `{{- $p := .page -}}
{{- $m := .menu -}}
{{ range (index $p.Site.Menus $m) -}}
{{- .URL }}|{{ .Name }}|{{ .Weight -}}|
{{- if $p.IsMenuCurrent $m . }}IsMenuCurrent{{ else }}-{{ end -}}|
{{- if $p.HasMenuCurrent $m . }}HasMenuCurrent{{ else }}-{{ end -}}|
{{- end -}}
`,
		"layouts/_default/single.html",
		`Single|{{ .Title }}
Menu Sect:  {{ partial "menu.html" (dict "page" . "menu" "sect") }}
Menu Main:  {{ partial "menu.html" (dict "page" . "menu" "main") }}`,
		"layouts/_default/list.html", "List|{{ .Title }}|{{ .Content }}",
	)
	require.Len(t, h.Sites, 1)

	fs := th.Fs

	writeSource(t, fs, "content/sect1/p1.md", fmt.Sprintf(menuPageTemplate, "p1", 1, "main", 40))
	writeSource(t, fs, "content/sect1/p2.md", fmt.Sprintf(menuPageTemplate, "p2", 2, "main", 30))
	writeSource(t, fs, "content/sect2/p3.md", fmt.Sprintf(menuPageTemplate, "p3", 3, "main", 20))
	writeSource(t, fs, "content/sect2/p4.md", fmt.Sprintf(menuPageTemplate, "p4", 4, "main", 10))
	writeSource(t, fs, "content/sect3/p5.md", fmt.Sprintf(menuPageTemplate, "p5", 5, "main", 5))

	writeNewContentFile(t, fs, "Section One", "2017-01-01", "content/sect1/_index.md", 100)
	writeNewContentFile(t, fs, "Section Five", "2017-01-01", "content/sect5/_index.md", 10)

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	s := h.Sites[0]

	require.Len(t, s.Menus, 2)

	p1 := s.RegularPages[0].Menus()

	// There is only one menu in the page, but it is "member of" 2
	require.Len(t, p1, 1)

	th.assertFileContent("public/sect1/p1/index.html", "Single",
		"Menu Sect:  /sect5/|Section Five|10|-|-|/sect1/|Section One|100|-|HasMenuCurrent|/sect2/|Sect2s|0|-|-|/sect3/|Sect3s|0|-|-|",
		"Menu Main:  /sect3/p5/|p5|5|-|-|/sect2/p4/|p4|10|-|-|/sect2/p3/|p3|20|-|-|/sect1/p2/|p2|30|-|-|/sect1/p1/|p1|40|IsMenuCurrent|-|",
	)

	th.assertFileContent("public/sect2/p3/index.html", "Single",
		"Menu Sect:  /sect5/|Section Five|10|-|-|/sect1/|Section One|100|-|-|/sect2/|Sect2s|0|-|HasMenuCurrent|/sect3/|Sect3s|0|-|-|")

}
