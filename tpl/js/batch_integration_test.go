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
-- assets/js/react1.jsx --
import * as React from "react";

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
{{ $bundle := (js.Bundle "mybundle" .Store) }}
{{ with $bundle.UseScriptOne "main1" }}
	{{ if not .GetResource }}
		{{ .SetResource (resources.Get "js/main1.js") }}
	{{ end }}
	{{ .SetInstance (dict "title" "Main1 Instance") }}
{{ end }}
 {{ with $bundle.UseScriptOne "main2" }}
	{{ if not .GetResource }}
		{{ .SetResource (resources.Get "js/main2.js") }}
	{{ end }}
	{{ .SetInstance (dict "title" "Main2 Instance") }}
{{ end }}
{{ with $bundle.UseScriptMany "reactbatch" }}
 	{{ if not .GetCallback }}
		{{ .SetCallback (resources.Get "js/reactcallback.js") }}
	{{ end }}
	{{ with .UseItem "r1" }}
		{{ if not .GetResource }}
		  {{ .SetResource (resources.Get "js/react1.jsx") }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance 1") }}
		{{ .AddInstance "i2" (dict "title" "Instance 2") }}
	{{ end }}
	 {{ with .UseItem "r2" }}
		{{ if not .GetResource }}
		  {{ .SetResource (resources.Get "js/react2.jsx") }}
		{{ end }}
		{{ .AddInstance "i1" (dict "title" "Instance 2-1") }}
	{{ end }}
{{ end }}
{{ range $k, $v := $bundle.Build.Scripts }}
{{ $k }}: {{ .RelPermalink }}
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

	b.AssertFileContent("public/index.html", `
main1: /js/bundles/mybundle/main1.js
main2: /js/bundles/mybundle/main2.js
reactbatch: /js/bundles/mybundle/reactbatch.js


	`)
}

// TODO1 make instance into a map with params as only key (for now)
