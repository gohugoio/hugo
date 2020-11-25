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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/htesting"

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
    "to-camel-case": "1.0.0"
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
var React = __toModule(require(&#34;react&#34;));
function greeter(person) {
`)

}

func TestJSBuild(t *testing.T) {
	if !isCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	if runtime.GOOS == "windows" {
		// TODO(bep) we really need to get this working on Travis.
		t.Skip("skip npm test on Windows")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-js-mod")
	c.Assert(err, qt.IsNil)
	defer clean()

	config := fmt.Sprintf(`
baseURL = "https://example.org"
workingDir = %q

disableKinds = ["page", "section", "term", "taxonomy"]

[module]
[[module.imports]]
path="github.com/gohugoio/hugoTestProjectJSModImports"



`, workDir)

	b := newTestSitesBuilder(t)
	b.Fs = hugofs.NewDefault(viper.New())
	b.WithWorkingDir(workDir).WithConfigFile("toml", config).WithLogger(loggers.NewInfoLogger())
	b.WithSourceFile("go.mod", `module github.com/gohugoio/tests/testHugoModules
        
go 1.15
        
require github.com/gohugoio/hugoTestProjectJSModImports v0.5.0 // indirect

`)

	b.WithContent("p1.md", "").WithNothingAdded()

	b.WithSourceFile("package.json", `{
 "dependencies": {
  "date-fns": "^2.16.1"
 }
}`)

	b.Assert(os.Chdir(workDir), qt.IsNil)
	_, err = exec.Command("npm", "install").CombinedOutput()
	b.Assert(err, qt.IsNil)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/js/main.js", `
greeting: "greeting configured in mod2"
Hello1 from mod1: $
return "Hello2 from mod1";
var Hugo = "Rocks!";
Hello3 from mod2. Date from date-fns: ${today}
Hello from lib in the main project
Hello5 from mod2.
var myparam = "Hugo Rocks!";`)

}
