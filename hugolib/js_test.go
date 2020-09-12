// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
)

func TestJSBuildWithNPM(t *testing.T) {
	if !isCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	if runtime.GOOS == "windows" {
		t.Skip("skip NPM test on Windows")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	mainJS := `
	import "./included";
	import { toCamelCase } from "to-camel-case";
	
	console.log("main");
	console.log("To camel:", toCamelCase("space case"));
`
	includedJS := `
	console.log("included");
	
	`

	jsxContent := `
import * as React from 'react'
import * as ReactDOM from 'react-dom'

 ReactDOM.render(
   <h1>Hello, world!</h1>,
   document.getElementById('root')
 );
`

	tsContent := `function greeter(person: string) {
    return "Hello, " + person;
}

let user = [0, 1, 2];

document.body.textContent = greeter(user);`

	packageJSON := `{
  "scripts": {},

  "dependencies": {
		"to-camel-case": "1.0.0",
		"react": "^16",
		"react-dom": "^16"
  }
}
`

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-js-npm")
	c.Assert(err, qt.IsNil)
	defer clean()

	v := viper.New()
	v.Set("workingDir", workDir)
	v.Set("disableKinds", []string{"taxonomy", "term", "page"})
	b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())

	// Need to use OS fs for this.
	b.Fs = hugofs.NewDefault(v)
	b.WithWorkingDir(workDir)
	b.WithViper(v)
	b.WithContent("p1.md", "")

	b.WithTemplates("index.html", `
{{ $options := dict "minify" false "externals" (slice "react" "react-dom") }}
{{ $js := resources.Get "js/main.js" | js.Build $options }}
JS:  {{ template "print" $js }}
{{ $jsx := resources.Get "js/myjsx.jsx" | js.Build $options }}
JSX: {{ template "print" $jsx }}
{{ $ts := resources.Get "js/myts.ts" | js.Build }}
TS: {{ template "print" $ts }}

{{ define "print" }}RelPermalink: {{.RelPermalink}}|MIME: {{ .MediaType }}|Content: {{ .Content | safeJS }}{{ end }}

`)

	jsDir := filepath.Join(workDir, "assets", "js")
	b.Assert(os.MkdirAll(jsDir, 0777), qt.IsNil)
	b.Assert(os.Chdir(workDir), qt.IsNil)
	b.WithSourceFile("package.json", packageJSON)
	b.WithSourceFile("assets/js/main.js", mainJS)
	b.WithSourceFile("assets/js/myjsx.jsx", jsxContent)
	b.WithSourceFile("assets/js/myts.ts", tsContent)

	b.WithSourceFile("assets/js/included.js", includedJS)

	out, err := exec.Command("npm", "install").CombinedOutput()
	b.Assert(err, qt.IsNil, qt.Commentf(string(out)))

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
console.log(&#34;included&#34;);
if (hasSpace.test(string))
const React = __toModule(require(&#34;react&#34;));
function greeter(person) {
`)

}

func TestJSBuild(t *testing.T) {
	if !isCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	mainJS := `
	import "./included";
	
	console.log("main");

`
	includedJS := `
	console.log("included");
	
	`

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-js")
	c.Assert(err, qt.IsNil)
	defer clean()

	v := viper.New()
	v.Set("workingDir", workDir)
	v.Set("disableKinds", []string{"taxonomy", "term", "page"})
	b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())

	b.Fs = hugofs.NewDefault(v)
	b.WithWorkingDir(workDir)
	b.WithViper(v)
	b.WithContent("p1.md", "")

	b.WithTemplates("index.html", `
{{ $js := resources.Get "js/main.js" | js.Build }}
JS:  {{ template "print" $js }}


{{ define "print" }}RelPermalink: {{.RelPermalink}}|MIME: {{ .MediaType }}|Content: {{ .Content | safeJS }}{{ end }}

`)

	jsDir := filepath.Join(workDir, "assets", "js")
	b.Assert(os.MkdirAll(jsDir, 0777), qt.IsNil)
	b.Assert(os.Chdir(workDir), qt.IsNil)
	b.WithSourceFile("assets/js/main.js", mainJS)
	b.WithSourceFile("assets/js/included.js", includedJS)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
console.log(&#34;included&#34;);

`)

}

func TestJSBuildGlobals(t *testing.T) {
	if !isCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-js")
	c.Assert(err, qt.IsNil)
	defer clean()

	v := viper.New()
	v.Set("workingDir", workDir)
	v.Set("disableKinds", []string{"taxonomy", "term", "page"})
	b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())

	b.Fs = hugofs.NewDefault(v)
	b.WithWorkingDir(workDir)
	b.WithViper(v)
	b.WithContent("p1.md", "")

	jsDir := filepath.Join(workDir, "assets", "js")
	b.Assert(os.MkdirAll(jsDir, 0777), qt.IsNil)
	b.Assert(os.Chdir(workDir), qt.IsNil)

	b.WithTemplates("index.html", `
{{- $js := resources.Get "js/main-project.js" | js.Build -}}
{{ template "print" (dict "js" $js "name" "root") }}

{{- define "print" -}}
{{ printf "rellink-%s-%s" .name .js.RelPermalink | safeHTML }}
{{ printf "mime-%s-%s" .name .js.MediaType | safeHTML }}
{{ printf "content-%s-%s" .name .js.Content | safeHTML }}
{{- end -}}
`)

	b.WithSourceFile("assets/js/normal.js", `
const name = "root-normal";
export default name;
`)
	b.WithSourceFile("assets/js/main-project.js", `
import normal from "@js/normal";
window.normal = normal; // make sure not to tree-shake
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
const name = "root-normal";
`)
}

func TestJSBuildOverride(t *testing.T) {
	if !isCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-js2")
	c.Assert(err, qt.IsNil)
	defer clean()
	// workDir := "/tmp/hugo-test-js2"
	c.Assert(os.Chdir(workDir), qt.IsNil)

	cfg := viper.New()
	cfg.Set("workingDir", workDir)
	fs := hugofs.NewFrom(afero.NewOsFs(), cfg)

	b := newTestSitesBuilder(t)
	b.Fs = fs
	b.WithLogger(loggers.NewWarningLogger())

	realWrite := func(name string, content string) {
		realLocation := filepath.Join(workDir, name)
		realDir := filepath.Dir(realLocation)
		if _, err := os.Stat(realDir); err != nil {
			os.MkdirAll(realDir, 0777)
		}
		bytesContent := []byte(content)
		// c.Assert(ioutil.WriteFile(realLocation, bytesContent, 0777), qt.IsNil)
		c.Assert(afero.WriteFile(b.Fs.Source, realLocation, bytesContent, 0777), qt.IsNil)
	}

	realWrite("config.toml", `
baseURL="https://example.org"

[module]
[[module.imports]]
path="mod2"
[[module.imports.mounts]]
source="assets"
target="assets"
[[module.imports.mounts]]
source="layouts"
target="layouts"
[[module.imports]]
path="mod1"
[[module.imports.mounts]]
source="assets"
target="assets"
[[module.imports.mounts]]
source="layouts"
target="layouts"
`)

	realWrite("content/p1.md", `---
layout: sample
---
`)
	realWrite("themes/mod1/layouts/_default/sample.html", `
{{- $js := resources.Get "js/main-project.js" | js.Build -}}
{{ template "print" (dict "js" $js "name" "root") }}

{{- $js = resources.Get "js/main-mod1.js" | js.Build -}}
{{ template "print" (dict "js" $js "name" "mod1") }}

{{- $js = resources.Get "js/main-mod2.js" | js.Build (dict "data" .Site.Params) -}}
{{ template "print" (dict "js" $js "name" "mod2") }}

{{- $js = resources.Get "js/main-mod2.js" | js.Build (dict "data" .Site.Params "sourceMap" "inline" "targetPath" "js/main-mod2-inline.js") -}}
{{ template "print" (dict "js" $js "name" "mod2") }}

{{- $js = resources.Get "js/main-mod2.js" | js.Build (dict "data" .Site.Params "sourceMap" "external" "targetPath" "js/main-mod2-external.js") -}}
{{ template "print" (dict "js" $js "name" "mod2") }}

{{- define "print" -}}
{{ printf "rellink-%s-%s" .name .js.RelPermalink | safeHTML }}
{{ printf "mime-%s-%s" .name .js.MediaType | safeHTML }}
{{ printf "content-%s-%s" .name .js.Content | safeHTML }}
{{- end -}}
`)

	// Override project included file
	// This file will override the one in mod1 and mod2
	realWrite("assets/js/override.js", `
const name = "root-override";
export default name;
`)

	// Add empty theme mod config files
	realWrite("themes/mod1/config.yml", ``)
	realWrite("themes/mod2/config.yml", ``)

	// This is the main project js file.
	// try to include @js/override which is overridden inside of project
	// try to include @js/override-mod which is overridden in mod2
	realWrite("assets/js/main-project.js", `
import override from "@js/override";
import overrideMod from "@js/override-mod";
window.override = override; // make sure to prevent tree-shake
window.overrideMod  = overrideMod; // make sure to prevent tree-shake
`)
	// This is the mod1 js file
	// try to include @js/override which is overridden inside of the project
	// try to include @js/override-mod which is overridden in mod2
	realWrite("themes/mod1/assets/js/main-mod1.js", `
import override from "@js/override";
import overrideMod from "@js/override-mod";
window.mod = "mod1";
window.override = override; // make sure to prevent tree-shake
window.overrideMod  = overrideMod; // make sure to prevent tree-shake
`)
	// This is the mod1 js file (overridden in mod2)
	// try to include @js/override which is overridden inside of the project
	// try to include @js/override-mod which is overridden in mod2
	realWrite("themes/mod2/assets/js/main-mod1.js", `
import override from "@js/override";
import overrideMod from "@js/override-mod";
window.mod = "mod2";
window.override = override; // make sure to prevent tree-shake
window.overrideMod  = overrideMod; // make sure to prevent tree-shake
`)
	// This is mod2 js file
	// try to include @js/override which is overridden inside of the project
	// try to include @js/override-mod which is overridden in mod2
	// try to include @config which is declared in a local jsconfig.json file
	// try to include @data which was passed as "data" into js.Build
	realWrite("themes/mod2/assets/js/main-mod2.js", `
import override from "@js/override";
import overrideMod from "@js/override-mod";
import config from "@config";
import data from "@data";
window.data = data;
window.override = override; // make sure to prevent tree-shake
window.overrideMod  = overrideMod; // make sure to prevent tree-shake
window.config = config;
`)
	realWrite("themes/mod2/assets/js/jsconfig.json", `
{
	"compilerOptions": {
		"baseUrl": ".",
		"paths": {
			"@config": ["./config.json"]
		}
	}
}
`)
	realWrite("themes/mod2/assets/js/config.json", `
{
	"data": {
		"sample": "sample"
	}
}
`)
	realWrite("themes/mod1/assets/js/override.js", `
const name = "mod1-override";
export default name;
`)
	realWrite("themes/mod2/assets/js/override.js", `
const name = "mod2-override";
export default name;
`)
	realWrite("themes/mod1/assets/js/override-mod.js", `
const nameMod = "mod1-override";
export default nameMod;
`)
	realWrite("themes/mod2/assets/js/override-mod.js", `
const nameMod = "mod2-override";
export default nameMod;
`)
	b.WithConfigFile("toml", `
baseURL="https://example.org"
themesDir="./themes"
[module]
[[module.imports]]
path="mod2"
[[module.imports.mounts]]
source="assets"
target="assets"
[[module.imports.mounts]]
source="layouts"
target="layouts"
[[module.imports]]
path="mod1"
[[module.imports.mounts]]
source="assets"
target="assets"
[[module.imports.mounts]]
source="layouts"
target="layouts"
`)

	b.WithWorkingDir(workDir)
	b.LoadConfig()

	b.Build(BuildCfg{})

	b.AssertFileContent("public/js/main-mod1.js", `
name = "root-override";
nameMod = "mod2-override";
window.mod = "mod2";
`)
	b.AssertFileContent("public/js/main-mod2.js", `
name = "root-override";
nameMod = "mod2-override";
sample: "sample"
"sect"
`)
	b.AssertFileContent("public/js/main-project.js", `
name = "root-override";
nameMod = "mod2-override";
`)
	b.AssertFileContent("public/js/main-mod2-external.js.map", `
const nameMod = \"mod2-override\";\nexport default nameMod;\n
"\nimport override from \"@js/override\";\nimport overrideMod from \"@js/override-mod\";\nimport config from \"@config\";\nimport data from \"@data\";\nwindow.data = data;\nwindow.override = override; // make sure to prevent tree-shake\nwindow.overrideMod  = overrideMod; // make sure to prevent tree-shake\nwindow.config = config;\n"
`)
	b.AssertFileContent("public/js/main-mod2-inline.js", `
	sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiYXNzZXRzL2pzL292ZXJyaWRlLmpzIiwgInRoZW
`)
}
