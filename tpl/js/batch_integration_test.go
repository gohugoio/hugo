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
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestBatch(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "page"]
baseURL = "https://example.com"
-- package.json --
{
  "devDependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  }
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
import './bar.css'
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
-- assets/js/react1.jsx --
import * as React from "react";
import './button.css'
import './foo.css'


window.React1 = React;

let text = 'Click me'

export default function MyButton() {
    return (
        <button>${text}</button>
    )
}
	
-- assets/js/react2.jsx --
import * as React from "react";

window.React2 = React;

let text = 'Click me, too!'

export default function MyOtherButton() {
    return (
        <button>${text}</button>
    )
}
-- assets/js/main1.js --
import * as React from "react";

console.log('main1.React', React)

-- assets/js/main2.js --
import * as React from "react";

console.log('main2.React', React)

-- layouts/index.html --
Home.
{{ $bundle := (js.Batch "mybundle" .Store) }}
{{ $otherCSS := (resources.Match "other/*.css").Mount "." }}
{{ with $bundle.UseScriptGroup "reactbatch" }}
 	{{ if not .GetCallback }}
		{{ .SetCallback (resources.Get "js/reactcallback.js") }}
	{{ end }}
	{{ with .UseScript "r1" }}
	 	{{ if not .GetImportContext }}
			{{ .SetImportContext $otherCSS }}
		{{ end }}
		{{ if not .GetResource }}
		  {{ .SetResource (resources.Get "js/react1.jsx") }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance 1") }}
		{{ .AddInstance "i2" (dict "title" "Instance 2") }}
	{{ end }}
	 {{ with .UseScript "r2" }}
		{{ if not .GetResource }}
		  {{ .SetResource (resources.Get "js/react2.jsx") }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance 2-1") }}
	{{ end }}
{{ end }}
{{ range $k, $v := $bundle.Build.Groups }}
 {{ range . }}
	{{ $k }}: {{ .RelPermalink }}
{{ end }}
{{ end }}}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               t,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			TxtarString:     files,
			// PrintAndKeepTempDir: true,
		}).Build()

	// b.AssertPublishDir("sadf")

	b.AssertFileContent("public/mybundle_reactbatch.js", `
"mod": MyButton, "id": "r1"
"mod": MyOtherButton, "id": "r2"


	`)
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
{{ $bundle := (js.Batch "myjsbundle" .Store) }}
{{ $js := .Resources.GetMatch "*.js" }}
{{ with $bundle.UseScriptGroup "g1" }}
	{{ with .UseScript "s1" }}
	 	{{ if not .GetImportContext }}
			{{ .SetImportContext $ }}
		{{ end }}
	 	{{ if not .GetResource }}
		  {{ .SetResource $js }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance s1-1") }}
	{{ end }}
{{ end }}
{{ range $bundle.Build.Groups }}
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

// TODO1  executing "_default/single.html" at <$bundle.Build.Groups>: error calling Build: Could not resolve "./mystyles.css"` error file source.

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
