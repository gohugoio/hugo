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
package js_test

import (
	"fmt"
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/hugolib"
)

// TODO1 fix shims vs headlessui.

func TestBatch(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
disableLiveReload = true
baseURL = "https://example.com"

# TOOD1
[build]
    [[build.cachebusters]]
        source = '.*'
        target = '.*'

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
-- content/mybundle/bundlestyles.css --
import './foo.css'
import './bar.css'
.bundlestyles {
	background-color: blue;
}
-- content/mybundle/bundlereact.jsx --
import * as React from "react";
import './foo.css'
import './bundlestyles.css'
window.React1 = React;

let text = 'Click me, too!'

export default function MyBundleButton() {
    return (
        <button>${text}</button>
    )
}

-- assets/js/reactcallback.js --
import * as ReactDOM from 'react-dom/client';
import * as React from 'react';

export default function Callback(modules) {
	for (const module of modules) {
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
import './bundlestyles.css'
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

export default function MyOtherButton() {
    return (
        <button>${text}</button>
    )
}
-- assets/js/main1.js --
import * as React from "react";
import * as params from '@params';

console.log('main1.React', React)
console.log('main1.params.id', params.id)

// TODO1 make it work without this.
export default function Main1() {};

-- assets/js/main2.js --
import * as React from "react";
import * as params from '@params';

console.log('main2.React', React)
console.log('main2.params.id', params.id)

export default function Main2() {};

-- assets/js/main3.js --
import * as React from "react";
import * as params from '@params';

console.log('main3.params.id', params.id)

export default function Main3() {};

-- layouts/_default/single.html --
Single.
{{ $r := .Resources.GetMatch "*.jsx" }}
{{ $batch := (js.Batch "mybundle" site.Home.Store) }}
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
{{ $batch := (js.Batch "mybundle" site.Home.Store) }}
{{ range $k, $v := $batch.Build.Groups }}
 {{ $k }}:
 {{ range . }}
	{{ .RelPermalink }}
  {{ end }}
 {{ end }}
{{ end }}
{{ $myContentBundle := site.GetPage "mybundle" }}
{{ $batch := (js.Batch "mybundle" site.Home.Store) }}
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
{{ with .Instance "main1" "m1i1" }}{{ .SetOptions (dict "title" "Main1 Instance 1")}}{{ end }}
{{ with .Instance "main1" "m1i2" }}{{ .SetOptions (dict "title" "Main1 Instance 2")}}{{ end }}
{{ end }}
{{ with $batch.Group "reactbatch" }}
 	{{ with .Callback "reactcallback" }}
		{{ .SetOptions ( dict "resource"  (resources.Get "js/reactcallback.js") )}}
	{{ end }}
	{{ with .Script "r1" }}
		{{ .SetOptions (dict
			"resource" (resources.Get "js/react1.jsx")
			"importContext" (slice $myContentBundle $otherCSS)
			"params" (dict "id" "r1")
			)
		}}
	{{ end }}
	{{ with .Instance "r1" "i1" }}{{ .SetOptions (dict "title" "Instance 1")}}{{ end }}
	{{ with .Instance "r1" "i2" }}{{ .SetOptions (dict "title" "Instance 2")}}{{ end }}
	{{ with .Script "r2" }}
		{{ .SetOptions (dict
			"resource" (resources.Get "js/react2.jsx")
			"importContext" $otherCSS
			"params" (dict "id" "r2")
			)
		}}
	{{ end }}
	{{ with .Instance "r2" "i1" }}{{ .SetOptions (dict "title" "Instance 2-1")}}{{ end }}
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

	fmt.Println(b.LogString())

	b.AssertFileContent("public/mybundle_reactbatch.css",
		".bar {",
	)

	// Verify params resolution.
	b.AssertFileContent("public/mybundle_mains.js",
		`
var id = "main1";
console.log("main1.params.id", id);
var id2 = "main2";
console.log("main2.params.id", id2);

# Params from top level config.
var id3 = "config";
console.log("main3.params.id", id3);
`)

	b.AssertFileContent("public/mybundle_reactbatch.js", `
"mod": MyButton, "id": "r1"
"mod": MyOtherButton, "id": "r2"
"mod": MyBundleButton, "id": "r3"

	`)

	b.EditFileReplaceAll("content/mybundle/bundlestyles.css", ".bundlestyles", ".bundlestyles-edit").Build()
	b.AssertFileContent("public/mybundle_reactbatch.css", ".bundlestyles-edit {")

	b.EditFileReplaceAll("assets/other/bar.css", ".bar {", ".bar-edit {").Build()
	b.AssertFileContent("public/mybundle_reactbatch.css", ".bar-edit {")

	b.EditFileReplaceAll("assets/other/bar.css", ".bar-edit {", ".bar-edit2 {").Build()
	b.AssertFileContent("public/mybundle_reactbatch.css", ".bar-edit2 {")
}

func TestEsBuildResolvePageBundle(t *testing.T) {
	files := `
-- hugo.toml --
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/mystyles1.css --
body {
	background-color: blue;
}
-- content/mybundle/mystyles2.css --
button {
	background-color: red;
}
-- content/mybundle/myscript.js --
import "./mystyles1.css";
import "./mystyles2.css";
console.log('Hello, world!');

// TODO1 make it work without this.
export default {};
-- layouts/_default/single.html --
Single.
TODO1 directory structure vs ID:
{{ $batch := (js.Batch "myjsbundle" .Store) }}
{{ $js := .Resources.GetMatch "*.js" }}
{{ with $batch.UseScriptGroup "g1" }}
	{{ with .Script "s1" }}
	 	{{ if not .GetImportContext }}
			{{ .SetImportContext $ }}
		{{ end }}
	 	{{ if not .GetResource }}
		  {{ .SetResource $js }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance s1-1") }}
	{{ end }}
{{ end }}
{{ range $batch.Build.Groups }}
 {{ range $i, $e := . }}
	{{ $i }}: {{ $e.RelPermalink }}|
 {{ end }}
{{ end }}

`

	// TODO1 check what happens without AddInstance.

	b := hugolib.Test(t, files, hugolib.TestOptWithOSFs())

	b.AssertFileContent("public/myjsbundle_g1.css", "body", "button")
	b.AssertFileContent("public/myjsbundle_g1.js", `Hello, world!`)
}

// TODO1  executing "_default/single.html" at <$batch.Build.Groups>: error calling Build: Could not resolve "./mystyles.css"` error file source.

// TODO1 move this.
func TestResourcesGet(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/text/txt1.txt --
Text 1.
-- assets/text/txt2.txt --
Text 2.
-- assets/text/sub/txt3.txt --
Text 3.
-- assets/text/sub/txt4.txt --
Text 4.
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/txt1.txt --
Text 1.
-- content/mybundle/sub/txt2.txt --
Text 1.
-- layouts/index.html --
{{ $mybundle := site.GetPage "mybundle" }}
{{ $subResources := resources.Match "text/sub/*.*"  }}
 {{ $subResourcesMount :=  $subResources.Mount "newroot" }}
resources:text/txt1.txt:{{ with resources.Get "text/txt1.txt" }}{{ .Name }}{{ end }}|
resources:text/txt2.txt:{{ with resources.Get "text/txt2.txt" }}{{ .Name }}{{ end }}|
resources:text/sub/txt3.txt:{{ with resources.Get "text/sub/txt3.txt" }}{{ .Name }}{{ end }}|
subResources.range:{{ range $subResources }}{{ .Name }}|{{ end }}|
subResources:"text/sub/txt3.txt:{{ with $subResources.Get "text/sub/txt3.txt" }}{{ .Name }}{{ end }}|
subResourcesMount:newroot/txt3.txt:{{ with $subResourcesMount.Get "newroot/txt3.txt" }}{{ .Name }}{{ end }}|
page:txt1.txt:{{ with $mybundle.Resources.Get "txt1.txt" }}{{ .Name }}{{ end }}|
page:./txt1.txt:{{ with $mybundle.Resources.Get "./txt1.txt" }}{{ .Name }}{{ end }}|
page:sub/txt2.txt:{{ with $mybundle.Resources.Get "sub/txt2.txt" }}{{ .Name }}{{ end }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
	asdf
resources:text/txt1.txt:/text/txt1.txt|
resources:text/txt2.txt:/text/txt2.txt|
resources:text/sub/txt3.txt:/text/sub/txt3.txt|
subResources:"text/sub/txt3.txt:/text/sub/txt3.txt|
subResourcesMount:newroot/txt3.txt:/text/sub/txt3.txt|
page:txt1.txt:txt1.txt|
page:./txt1.txt:txt1.txt|
page:sub/txt2.txt:sub/txt2.txt|
`)
}

// TODO1 check .Name in bundles on renames.
// TODO1 https://esbuild.github.io/content-types/#local-css
