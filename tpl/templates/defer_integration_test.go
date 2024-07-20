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

package templates_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

const deferFilesCommon = `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["taxonomy", "term", "rss", "sitemap", "robotsTXT", "404", "section"]
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- i18n/en.toml --
[hello]
other = "Hello"
-- i18n/nn.toml --
[hello]
other = "Hei"
-- content/_index.en.md --
---
title: "Home"
outputs: ["html", "amp"]
---
-- content/_index.nn.md --
---
title: "Heim"
outputs: ["html", "amp"]
---
-- assets/mytext.txt --
Hello.
-- layouts/baseof.html --
HTML|{{ block "main" . }}{{ end }}$
-- layouts/index.html --
{{ define "main" }}
EDIT_COUNTER_OUTSIDE_0
{{ .Store.Set "hello" "Hello" }}
{{ $data := dict "page" . }}
{{ with (templates.Defer (dict "data" $data) ) }}
{{ $mytext := resources.Get "mytext.txt" }}
REPLACE_ME|Title: {{ .page.Title }}|{{ .page.RelPermalink }}|Hello: {{ T "hello" }}|Hello Store: {{ .page.Store.Get "hello" }}|Mytext: {{ $mytext.Content }}|
EDIT_COUNTER_DEFER_0
{{ end }}$
{{ end }}
-- layouts/index.amp.html --
AMP.
{{ $data := dict "page" . }}
{{ with (templates.Defer (dict "data" $data) ) }}Title AMP: {{ .page.Title }}|{{ .page.RelPermalink }}|Hello: {{ T "hello" }}{{ end }}$

`

func TestDeferBasic(t *testing.T) {
	t.Parallel()

	b := hugolib.Test(t, deferFilesCommon)

	b.AssertFileContent("public/index.html", "Title: Home|/|Hello: Hello|Hello Store: Hello|Mytext: Hello.|")
	b.AssertFileContent("public/amp/index.html", "Title AMP: Home|/amp/|Hello: Hello")
	b.AssertFileContent("public/nn/index.html", "Title: Heim|/nn/|Hello: Hei")
	b.AssertFileContent("public/nn/amp/index.html", "Title AMP: Heim|/nn/amp/|Hello: Hei")
}

func TestDeferRepeatedBuildsEditOutside(t *testing.T) {
	t.Parallel()

	b := hugolib.TestRunning(t, deferFilesCommon)

	for i := 0; i < 5; i++ {
		old := fmt.Sprintf("EDIT_COUNTER_OUTSIDE_%d", i)
		new := fmt.Sprintf("EDIT_COUNTER_OUTSIDE_%d", i+1)
		b.EditFileReplaceAll("layouts/index.html", old, new).Build()
		b.AssertFileContent("public/index.html", new)
	}
}

func TestDeferRepeatedBuildsEditDefer(t *testing.T) {
	t.Parallel()

	b := hugolib.TestRunning(t, deferFilesCommon)

	for i := 0; i < 8; i++ {
		old := fmt.Sprintf("EDIT_COUNTER_DEFER_%d", i)
		new := fmt.Sprintf("EDIT_COUNTER_DEFER_%d", i+1)
		b.EditFileReplaceAll("layouts/index.html", old, new).Build()
		b.AssertFileContent("public/index.html", new)
	}
}

func TestDeferErrorParse(t *testing.T) {
	t.Parallel()

	b, err := hugolib.TestE(t, strings.ReplaceAll(deferFilesCommon, "Title AMP: {{ .page.Title }}", "{{ .page.Title }"))

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, `index.amp.html:3: unexpected "}" in operand`)
}

func TestDeferErrorRuntime(t *testing.T) {
	t.Parallel()

	b, err := hugolib.TestE(t, strings.ReplaceAll(deferFilesCommon, "Title AMP: {{ .page.Title }}", "{{ .page.Titles }}"))

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`/layouts/index.amp.html:3:57`))
	b.Assert(err.Error(), qt.Contains, `execute of template failed: template: index.amp.html:3:57: executing at <.page.Titles>: can't evaluate field Titles`)
}

func TestDeferEditDeferBlock(t *testing.T) {
	t.Parallel()

	b := hugolib.TestRunning(t, deferFilesCommon)
	b.AssertRenderCountPage(4)
	b.EditFileReplaceAll("layouts/index.html", "REPLACE_ME", "Edited.").Build()
	b.AssertFileContent("public/index.html", "Edited.")
	b.AssertRenderCountPage(2)
}

//

func TestDeferEditResourceUsedInDeferBlock(t *testing.T) {
	t.Parallel()

	b := hugolib.TestRunning(t, deferFilesCommon)
	b.AssertRenderCountPage(4)
	b.EditFiles("assets/mytext.txt", "Mytext Hello Edited.").Build()
	b.AssertFileContent("public/index.html", "Mytext Hello Edited.")
	b.AssertRenderCountPage(2)
}

func TestDeferMountPublic(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[module]
[[module.mounts]]
source = "content"
target = "content"
[[module.mounts]]
source = "layouts"
target = "layouts"
[[module.mounts]]
source = 'public'
target = 'assets/public'
disableWatch = true
-- layouts/index.html --
Home.
{{ $mydata := dict "v1" "v1value" }}
{{ $json := resources.FromString "mydata/data.json" ($mydata | jsonify ) }}
{{ $nop := $json.RelPermalink }}
{{ with (templates.Defer (dict "key" "foo")) }}
	  {{ $jsonFilePublic := resources.Get "public/mydata/data.json" }}
	  {{ with  $jsonFilePublic }}
      {{ $m := $jsonFilePublic | transform.Unmarshal }}
	  v1: {{ $m.v1 }}
	  {{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "v1: v1value")
}

func TestDeferFromContentAdapterShouldFail(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_content.gotmpl --
{{ with (templates.Defer (dict "key" "foo")) }}
 Foo.
{{ end }}
`

	b, err := hugolib.TestE(t, files)

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "error calling Defer: this method cannot be called before the site is fully initialized")
}

func TestDeferPostProcessShouldThrowAnError(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/mytext.txt --
ABCD.
-- layouts/index.html --
Home
{{ with (templates.Defer (dict "key" "foo")) }}
{{ $mytext := resources.Get "mytext.txt" | minify | resources.PostProcess }}
{{ end }}

`
	b, err := hugolib.TestE(t, files)

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "resources.PostProcess cannot be used in a deferred template")
}
