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

// Package js provides functions for building JavaScript resources
package esbuild_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/internal/js/esbuild"
)

// Used to test misc. error situations etc.
const jsBatchFilesTemplate = `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "section"]
disableLiveReload = true
-- assets/js/styles.css --
body {
	background-color: red;
}
-- assets/js/main.js --
import './styles.css';
import * as params from '@params';
import * as foo from 'mylib';
console.log("Hello, Main!");
console.log("params.p1", params.p1);
export default function Main() {};
-- assets/js/runner.js --
console.log("Hello, Runner!");
-- node_modules/mylib/index.js --
console.log("Hello, My Lib!");
-- layouts/shortcodes/hdx.html --
{{ $path := .Get "r" }}
{{ $r := or (.Page.Resources.Get $path) (resources.Get $path) }}
{{ $batch := (js.Batch "mybatch") }}
{{ $scriptID := $path | anchorize }}
{{ $instanceID :=  .Ordinal | string }}
{{ $group := .Page.RelPermalink | anchorize }}
{{ $params := .Params | default dict }}
{{ $export := .Get "export" | default "default" }}
{{ with $batch.Group $group }}
	{{ with .Runner "create-elements" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/runner.js")) }}
	{{ end }}
	{{ with .Script $scriptID }}
		{{ .SetOptions (dict
			"resource" $r
			"export" $export
			"importContext" (slice $.Page)
		)
		}}
		{{ end }}
	{{ with .Instance $scriptID $instanceID }}
		{{ .SetOptions (dict "params" $params) }}
	{{ end }}
{{ end }}
hdx-instance: {{ $scriptID }}: {{ $instanceID }}|
-- layouts/_default/baseof.html --
Base.
{{ $batch := (js.Batch "mybatch") }}
 {{ with $batch.Config }}
	{{ .SetOptions (dict 
		"params" (dict "id" "config")
		"sourceMap" ""
		)
	}}
{{ end }}
{{ with (templates.Defer (dict "key" "global")) }}
Defer:
{{ $batch := (js.Batch "mybatch") }}
{{ range $k, $v := $batch.Build.Groups }}
	{{ range $kk, $vv := . -}}
		{{ $k }}: {{ .RelPermalink }}
	{{ end }}
{{ end -}}
{{ end }}
{{ block "main" . }}Main{{ end }}
End.
-- layouts/_default/single.html --
{{ define "main" }}
==> Single Template Content: {{ .Content }}$
{{ $batch := (js.Batch "mybatch") }}
{{ with $batch.Group "mygroup" }}
 	{{ with .Runner "run" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/runner.js")) }}
	{{ end }}
	{{ with .Script "main" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/main.js") "params" (dict "p1" "param-p1-main" )) }}
	{{ end }}
	{{ with .Instance "main" "i1" }}
		{{ .SetOptions (dict "params" (dict "title" "Instance 1")) }}
	{{ end }}
{{ end }}
{{ end }}
-- layouts/index.html --
{{ define "main" }}
Home.
{{ end }}
-- content/p1/index.md --
---
title: "P1"
---

Some content.

{{< hdx r="p1script.js" myparam="p1-param-1" >}}
{{< hdx r="p1script.js" myparam="p1-param-2" >}}

-- content/p1/p1script.js --
console.log("P1 Script");


`

// Just to verify that the above file setup works.
func TestBatchTemplateOKBuild(t *testing.T) {
	b := hugolib.Test(t, jsBatchFilesTemplate, hugolib.TestOptWithOSFs())
	b.AssertPublishDir("mybatch/mygroup.js", "mybatch/mygroup.css")
}

func TestBatchRemoveAllInGroup(t *testing.T) {
	files := jsBatchFilesTemplate
	b := hugolib.TestRunning(t, files, hugolib.TestOptWithOSFs())

	b.AssertFileContent("public/p1/index.html", "p1: /mybatch/p1.js")

	b.EditFiles("content/p1/index.md", `
---
title: "P1"
---
Empty.
`)
	b.Build()

	b.AssertFileContent("public/p1/index.html", "! p1: /mybatch/p1.js")

	// Add one script back.
	b.EditFiles("content/p1/index.md", `
---
title: "P1"
---

{{< hdx r="p1script.js" myparam="p1-param-1-new" >}}
`)
	b.Build()

	b.AssertFileContent("public/mybatch/p1.js",
		"p1-param-1-new",
		"p1script.js")
}

func TestBatchEditInstance(t *testing.T) {
	files := jsBatchFilesTemplate
	b := hugolib.TestRunning(t, files, hugolib.TestOptWithOSFs())
	b.AssertFileContent("public/mybatch/mygroup.js", "Instance 1")
	b.EditFileReplaceAll("layouts/_default/single.html", "Instance 1", "Instance 1 Edit").Build()
	b.AssertFileContent("public/mybatch/mygroup.js", "Instance 1 Edit")
}

func TestBatchEditScriptParam(t *testing.T) {
	files := jsBatchFilesTemplate
	b := hugolib.TestRunning(t, files, hugolib.TestOptWithOSFs())
	b.AssertFileContent("public/mybatch/mygroup.js", "param-p1-main")
	b.EditFileReplaceAll("layouts/_default/single.html", "param-p1-main", "param-p1-main-edited").Build()
	b.AssertFileContent("public/mybatch/mygroup.js", "param-p1-main-edited")
}

func TestBatchRenameBundledScript(t *testing.T) {
	files := jsBatchFilesTemplate
	b := hugolib.TestRunning(t, files, hugolib.TestOptWithOSFs())
	b.AssertFileContent("public/mybatch/p1.js", "P1 Script")
	b.RenameFile("content/p1/p1script.js", "content/p1/p1script2.js")
	_, err := b.BuildE()
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "resource not set")

	// Rename it back.
	b.RenameFile("content/p1/p1script2.js", "content/p1/p1script.js")
	b.Build()
}

func TestBatchErrorScriptResourceNotSet(t *testing.T) {
	files := strings.Replace(jsBatchFilesTemplate, `(resources.Get "js/main.js")`, `(resources.Get "js/doesnotexist.js")`, 1)
	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `error calling SetOptions: resource not set`)
}

func TestBatchSlashInBatchID(t *testing.T) {
	files := strings.ReplaceAll(jsBatchFilesTemplate, `"mybatch"`, `"my/batch"`)
	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNil)
	b.AssertPublishDir("my/batch/mygroup.js")
}

func TestBatchSourceMaps(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "section"]
disableLiveReload = true
-- assets/js/styles.css --
body {
	background-color: red;
}
-- assets/js/main.js --
import * as foo from 'mylib';
console.log("Hello, Main!");
-- assets/js/runner.js --
console.log("Hello, Runner!");
-- node_modules/mylib/index.js --
console.log("Hello, My Lib!");
-- layouts/shortcodes/hdx.html --
{{ $path := .Get "r" }}
{{ $r := or (.Page.Resources.Get $path) (resources.Get $path) }}
{{ $batch := (js.Batch "mybatch") }}
{{ $scriptID := $path | anchorize }}
{{ $instanceID :=  .Ordinal | string }}
{{ $group := .Page.RelPermalink | anchorize }}
{{ $params := .Params | default dict }}
{{ $export := .Get "export" | default "default" }}
{{ with $batch.Group $group }}
	{{ with .Runner "create-elements" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/runner.js")) }}
	{{ end }}
	{{ with .Script $scriptID }}
		{{ .SetOptions (dict
			"resource" $r
			"export" $export
			"importContext" (slice $.Page)
		)
		}}
		{{ end }}
	{{ with .Instance $scriptID $instanceID }}
		{{ .SetOptions (dict "params" $params) }}
	{{ end }}
{{ end }}
hdx-instance: {{ $scriptID }}: {{ $instanceID }}|
-- layouts/_default/baseof.html --
Base.
{{ $batch := (js.Batch "mybatch") }}
 {{ with $batch.Config }}
	{{ .SetOptions (dict 
		"params" (dict "id" "config")
		"sourceMap" ""
		)
	}}
{{ end }}
{{ with (templates.Defer (dict "key" "global")) }}
Defer:
{{ $batch := (js.Batch "mybatch") }}
{{ range $k, $v := $batch.Build.Groups }}
	{{ range $kk, $vv := . -}}
		{{ $k }}: {{ .RelPermalink }}
	{{ end }}
{{ end -}}
{{ end }}
{{ block "main" . }}Main{{ end }}
End.
-- layouts/_default/single.html --
{{ define "main" }}
==> Single Template Content: {{ .Content }}$
{{ $batch := (js.Batch "mybatch") }}
{{ with $batch.Group "mygroup" }}
 	{{ with .Runner "run" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/runner.js")) }}
	{{ end }}
	{{ with .Script "main" }}
		{{ .SetOptions (dict "resource" (resources.Get "js/main.js") "params" (dict "p1" "param-p1-main" )) }}
	{{ end }}
	{{ with .Instance "main" "i1" }}
		{{ .SetOptions (dict "params" (dict "title" "Instance 1")) }}
	{{ end }}
{{ end }}
{{ end }}
-- layouts/index.html --
{{ define "main" }}
Home.
{{ end }}
-- content/p1/index.md --
---
title: "P1"
---

Some content.

{{< hdx r="p1script.js" myparam="p1-param-1" >}}
{{< hdx r="p1script.js" myparam="p1-param-2" >}}

-- content/p1/p1script.js --
import * as foo from 'mylib';
console.lg("Foo", foo);
console.log("P1 Script");
export default function P1Script() {};


`
	files := strings.Replace(filesTemplate, `"sourceMap" ""`, `"sourceMap" "linked"`, 1)
	b := hugolib.TestRunning(t, files, hugolib.TestOptWithOSFs())
	b.AssertFileContent("public/mybatch/mygroup.js.map", "main.js", "! ns-hugo")
	b.AssertFileContent("public/mybatch/mygroup.js", "sourceMappingURL=mygroup.js.map")
	b.AssertFileContent("public/mybatch/p1.js", "sourceMappingURL=p1.js.map")
	b.AssertFileContent("public/mybatch/mygroup_run_runner.js", "sourceMappingURL=mygroup_run_runner.js.map")
	b.AssertFileContent("public/mybatch/chunk-UQKPPNA6.js", "sourceMappingURL=chunk-UQKPPNA6.js.map")

	checkMap := func(p string, expectLen int) {
		s := b.FileContent(p)
		sources := esbuild.SourcesFromSourceMap(s)
		b.Assert(sources, qt.HasLen, expectLen)

		// Check that all source files exist.
		for _, src := range sources {
			filename, ok := paths.UrlStringToFilename(src)
			b.Assert(ok, qt.IsTrue)
			_, err := os.Stat(filename)
			b.Assert(err, qt.IsNil)
		}
	}

	checkMap("public/mybatch/mygroup.js.map", 1)
	checkMap("public/mybatch/p1.js.map", 1)
	checkMap("public/mybatch/mygroup_run_runner.js.map", 0)
	checkMap("public/mybatch/chunk-UQKPPNA6.js.map", 1)
}

func TestBatchErrorRunnerResourceNotSet(t *testing.T) {
	files := strings.Replace(jsBatchFilesTemplate, `(resources.Get "js/runner.js")`, `(resources.Get "js/doesnotexist.js")`, 1)
	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `resource not set`)
}

func TestBatchErrorScriptResourceInAssetsSyntaxError(t *testing.T) {
	// Introduce JS syntax error in assets/js/main.js
	files := strings.Replace(jsBatchFilesTemplate, `console.log("Hello, Main!");`, `console.log("Hello, Main!"`, 1)
	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`assets/js/main.js:5:0": Expected ")" but found "console"`))
}

func TestBatchErrorScriptResourceInBundleSyntaxError(t *testing.T) {
	// Introduce JS syntax error in content/p1/p1script.js
	files := strings.Replace(jsBatchFilesTemplate, `console.log("P1 Script");`, `console.log("P1 Script"`, 1)
	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`/content/p1/p1script.js:3:0": Expected ")" but found end of file`))
}

func TestBatch(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
disableLiveReload = true
baseURL = "https://example.com"
-- package.json --
{
  "devDependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  }
}
-- assets/js/shims/react.js --
-- assets/js/shims/react-dom.js --
module.exports = window.ReactDOM;
module.exports = window.React;
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/mybundlestyles.css --
@import './foo.css';
@import './bar.css';
@import './otherbundlestyles.css';

.mybundlestyles {
	background-color: blue;
}
-- content/mybundle/bundlereact.jsx --
import * as React from "react";
import './foo.css';
import './mybundlestyles.css';
window.React1 = React;

let text = 'Click me, too!'

export default function MyBundleButton() {
    return (
        <button>${text}</button>
    )
}

-- assets/js/reactrunner.js --
import * as ReactDOM from 'react-dom/client';
import * as React from 'react';

export default function Run(group) {
	for (const module of group.scripts) {
		for (const instance of module.instances) {
			/* This is a convention in this project. */
			let elId = §§${module.id}-${instance.id}§§;
			let el = document.getElementById(elId);
			if (!el) {
				console.warn(§§Element with id ${elId} not found§§);
				continue;
			}
			const root = ReactDOM.createRoot(el);
			const reactEl = React.createElement(module.mod, instance.params);
			root.render(reactEl);
		}
	}
}
-- assets/other/otherbundlestyles.css --
.otherbundlestyles {
	background-color: red;
}
-- assets/other/foo.css --
@import './bar.css';

.foo {
	background-color: blue;
}
-- assets/other/bar.css --
.bar {
	background-color: red;
}
-- assets/js/button.css --
button {
	background-color: red;
}
-- assets/js/bar.css --
.bar-assets {
	background-color: red;
}
-- assets/js/helper.js --
import './bar.css'

export function helper() {
	console.log('helper');
}	

-- assets/js/react1styles_nested.css --
.react1styles_nested {
	background-color: red;
}
-- assets/js/react1styles.css --
@import './react1styles_nested.css';
.react1styles {
	background-color: red;
}
-- assets/js/react1.jsx --
import * as React from "react";
import './button.css'
import './foo.css'
import './react1styles.css'

window.React1 = React;

let text = 'Click me'

export default function MyButton() {
    return (
        <button>${text}</button>
    )
}
	
-- assets/js/react2.jsx --
import * as React from "react";
import { helper } from './helper.js'
import './foo.css'

window.React2 = React;

let text = 'Click me, too!'

export function MyOtherButton() {
    return (
        <button>${text}</button>
    )
}
-- assets/js/main1.js --
import * as React from "react";
import * as params from '@params';

console.log('main1.React', React)
console.log('main1.params.id', params.id)

-- assets/js/main2.js --
import * as React from "react";
import * as params from '@params';

console.log('main2.React', React)
console.log('main2.params.id', params.id)

export default function Main2() {};

-- assets/js/main3.js --
import * as React from "react";
import * as params from '@params';
import * as config from '@params/config';

console.log('main3.params.id', params.id)
console.log('config.params.id', config.id)

export default function Main3() {};

-- layouts/_default/single.html --
Single.

{{ $r := .Resources.GetMatch "*.jsx" }}
{{ $batch := (js.Batch "mybundle") }}
{{ $otherCSS := (resources.Match "/other/*.css").Mount "/other" "." }}
 {{ with $batch.Config }}
	  {{ $shims := dict "react" "js/shims/react.js"  "react-dom/client" "js/shims/react-dom.js" }}
      {{ .SetOptions (dict 
           "target" "es2018"
           "params" (dict "id" "config")
		   "shims" $shims
	     )
      }}
{{ end }}
{{ with $batch.Group "reactbatch" }}
	{{ with .Script "r3" }}
		 	{{ .SetOptions (dict
			 	"resource" $r
 				"importContext" (slice $ $otherCSS)
				"params" (dict "id" "r3")
               )
			}}
	{{ end }}
	{{ with .Instance "r3" "r2i1" }}
	 	{{ .SetOptions  (dict "title" "r2 instance 1")}}
	{{ end }}
{{ end }}
-- layouts/index.html --
Home.
{{ with (templates.Defer (dict "key" "global")) }}
{{ $batch := (js.Batch "mybundle") }}
{{ range $k, $v := $batch.Build.Groups }}
 {{ range $kk, $vv := . }}
	 {{ $k }}: {{ $kk }}: {{ .RelPermalink }}
  {{ end }}
 {{ end }}
{{ end }}
{{ $myContentBundle := site.GetPage "mybundle" }}
{{ $batch := (js.Batch "mybundle") }}
{{ $otherCSS := (resources.Match "/other/*.css").Mount "/other" "." }}
{{ with $batch.Group "mains" }}
  {{ with .Script "main1" }}
	{{ .SetOptions (dict
			 	"resource" (resources.Get "js/main1.js")
				"params" (dict "id" "main1")
            )
	}}
  {{ end }}
  {{ with .Script "main2" }}
   {{ .SetOptions (dict
			 	"resource" (resources.Get "js/main2.js")
				"params" (dict "id" "main2")
            )
  }}
  {{ end }}
 {{ with .Script "main3" }}
   {{ .SetOptions (dict
			 	"resource" (resources.Get "js/main3.js")
            )
  }}
  {{ end }}
{{ with .Instance "main1" "m1i1" }}{{ .SetOptions (dict "params" (dict "title" "Main1 Instance 1"))}}{{ end }}
{{ with .Instance "main1" "m1i2" }}{{ .SetOptions (dict "params" (dict "title" "Main1 Instance 2"))}}{{ end }}
{{ end }}
{{ with $batch.Group "reactbatch" }}
 	{{ with .Runner "reactrunner" }}
		{{ .SetOptions ( dict "resource"  (resources.Get "js/reactrunner.js") )}}
	{{ end }}
	{{ with .Script "r1" }}
		{{ .SetOptions (dict
			"resource" (resources.Get "js/react1.jsx")
			"importContext" (slice $myContentBundle $otherCSS)
			"params" (dict "id" "r1")
			)
		}}
	{{ end }}
	{{ with .Instance "r1" "i1" }}{{ .SetOptions (dict "params" (dict "title" "Instance 1"))}}{{ end }}
	{{ with .Instance "r1" "i2" }}{{ .SetOptions (dict "params" (dict "title" "Instance 2"))}}{{ end }}
	{{ with .Script "r2" }}
		{{ .SetOptions (dict
			"resource" (resources.Get "js/react2.jsx")
			"export" "MyOtherButton"
			"importContext" $otherCSS
			"params" (dict "id" "r2")
			)
		}}
	{{ end }}
	{{ with .Instance "r2" "i1" }}{{ .SetOptions (dict "params" (dict "title" "Instance 2-1"))}}{{ end }}
{{ end }}
 
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               t,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			TxtarString:     files,
			Running:         true,
			LogLevel:        logg.LevelWarn,
			// PrintAndKeepTempDir: true,
		}).Build()

	b.AssertFileContent("public/index.html",
		"mains: 0: /mybundle/mains.js",
		"reactbatch: 2: /mybundle/reactbatch.css",
	)

	b.AssertFileContent("public/mybundle/reactbatch.css",
		".bar {",
	)

	// Verify params resolution.
	b.AssertFileContent("public/mybundle/mains.js",
		`
var id = "main1";
console.log("main1.params.id", id);
var id2 = "main2";
console.log("main2.params.id", id2);


# Params from top level config.
var id3 = "config";
console.log("main3.params.id", void 0);
console.log("config.params.id", id3);
`)

	b.EditFileReplaceAll("content/mybundle/mybundlestyles.css", ".mybundlestyles", ".mybundlestyles-edit").Build()
	b.AssertFileContent("public/mybundle/reactbatch.css", ".mybundlestyles-edit {")

	b.EditFileReplaceAll("assets/other/bar.css", ".bar {", ".bar-edit {").Build()
	b.AssertFileContent("public/mybundle/reactbatch.css", ".bar-edit {")

	b.EditFileReplaceAll("assets/other/bar.css", ".bar-edit {", ".bar-edit2 {").Build()
	b.AssertFileContent("public/mybundle/reactbatch.css", ".bar-edit2 {")
}

func TestEditBaseofManyTimes(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
disableKinds = ["taxonomy", "term"]
-- layouts/_default/baseof.html --
Baseof.
{{ block "main" . }}{{ end }}
{{ with (templates.Defer (dict "key" "global")) }}
Now. {{ now }}
{{ end }}
-- layouts/_default/single.html --
{{ define "main" }}
Single.
{{ end }}
--
-- layouts/_default/list.html --
{{ define "main" }}
List.
{{ end }}
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/_index.md --
---
title: "Home"
---
`

	b := hugolib.TestRunning(t, files)
	b.AssertFileContent("public/index.html", "Baseof.")

	for i := 0; i < 100; i++ {
		b.EditFileReplaceAll("layouts/_default/baseof.html", "Now", "Now.").Build()
		b.AssertFileContent("public/index.html", "Now..")
	}
}
