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

package js_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/internal/js/esbuild"
)

func TestBuildVariants(t *testing.T) {
	c := qt.New(t)

	mainWithImport := `
-- config.toml --
disableKinds=["page", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
disableLiveReload = true
-- assets/js/main.js --
import { hello1, hello2 } from './util1';
hello1();
hello2();
-- assets/js/util1.js --
import { hello3 } from './util2';
export function hello1() {
	return 'abcd';
}
export function hello2() {
	return hello3();
}
-- assets/js/util2.js --
export function hello3() {
	return 'efgh';
}
-- layouts/index.html --
{{ $js := resources.Get "js/main.js" | js.Build }}
JS Content:{{ $js.Content }}:End:

			`

	c.Run("Basic", func(c *qt.C) {
		b := hugolib.NewIntegrationTestBuilder(hugolib.IntegrationTestConfig{T: c, NeedsOsFS: true, TxtarString: mainWithImport}).Build()

		b.AssertFileContent("public/index.html", `abcd`)
	})

	c.Run("Edit Import", func(c *qt.C) {
		b := hugolib.NewIntegrationTestBuilder(hugolib.IntegrationTestConfig{T: c, Running: true, NeedsOsFS: true, TxtarString: mainWithImport}).Build()

		b.AssertFileContent("public/index.html", `abcd`)
		b.EditFileReplaceFunc("assets/js/util1.js", func(s string) string { return strings.ReplaceAll(s, "abcd", "1234") }).Build()
		b.AssertFileContent("public/index.html", `1234`)
	})

	c.Run("Edit Import Nested", func(c *qt.C) {
		b := hugolib.NewIntegrationTestBuilder(hugolib.IntegrationTestConfig{T: c, Running: true, NeedsOsFS: true, TxtarString: mainWithImport}).Build()

		b.AssertFileContent("public/index.html", `efgh`)
		b.EditFileReplaceFunc("assets/js/util2.js", func(s string) string { return strings.ReplaceAll(s, "efgh", "1234") }).Build()
		b.AssertFileContent("public/index.html", `1234`)
	})
}

func TestBuildWithModAndNpm(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	c := qt.New(t)

	files := `
-- config.toml --
baseURL = "https://example.org"
disableKinds=["page", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
[module]
[[module.imports]]
path="github.com/gohugoio/hugoTestProjectJSModImports"
-- go.mod --
module github.com/gohugoio/tests/testHugoModules

go 1.16

require github.com/gohugoio/hugoTestProjectJSModImports v0.10.0 // indirect
-- package.json --
{
	"dependencies": {
	"date-fns": "^2.16.1"
	}
}

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			TxtarString:     files,
			Verbose:         true,
		}).Build()

	b.AssertFileContent("public/js/main.js", `
greeting: "greeting configured in mod2"
Hello1 from mod1: $
return "Hello2 from mod1";
var Hugo = "Rocks!";
Hello3 from mod2. Date from date-fns: ${today}
Hello from lib in the main project
Hello5 from mod2.
var myparam = "Hugo Rocks!";
shim cwd
`)

	// React JSX, verify the shimming.
	b.AssertFileContent("public/js/like.js", filepath.FromSlash(`@v0.10.0/assets/js/shims/react.js
module.exports = window.ReactDOM;
`))
}

func TestBuildWithNpm(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	c := qt.New(t)

	files := `
-- assets/js/included.js --
console.log("included");
-- assets/js/main.js --
import "./included";
	import { toCamelCase } from "to-camel-case";

	console.log("main");
	console.log("To camel:", toCamelCase("space case"));
-- assets/js/myjsx.jsx --
import * as React from 'react'
import * as ReactDOM from 'react-dom'

	ReactDOM.render(
	<h1>Hello, world!</h1>,
	document.getElementById('root')
	);
-- assets/js/myts.ts --
function greeter(person: string) {
	return "Hello, " + person;
}
let user = [0, 1, 2];
document.body.textContent = greeter(user);
-- config.toml --
disablekinds = ['taxonomy', 'term', 'page']
-- content/p1.md --
Content.
-- data/hugo.toml --
slogan = "Hugo Rocks!"
-- i18n/en.yaml --
hello:
   other: "Hello"
-- i18n/fr.yaml --
hello:
   other: "Bonjour"
-- layouts/index.html --
{{ $options := dict "minify" false "externals" (slice "react" "react-dom")  "sourcemap" "linked" }}
{{ $js := resources.Get "js/main.js" | js.Build $options }}
JS:  {{ template "print" $js }}
{{ $jsx := resources.Get "js/myjsx.jsx" | js.Build $options }}
JSX: {{ template "print" $jsx }}
{{ $ts := resources.Get "js/myts.ts" | js.Build (dict "sourcemap" "inline")}}
TS: {{ template "print" $ts }}
{{ $ts2 := resources.Get "js/myts.ts" | js.Build (dict "sourcemap" "external" "TargetPath" "js/myts2.js")}}
TS2: {{ template "print" $ts2 }}
{{ define "print" }}RelPermalink: {{.RelPermalink}}|MIME: {{ .MediaType }}|Content: {{ .Content | safeJS }}{{ end }}
-- package.json --
{
	"scripts": {},

	"dependencies": {
	"to-camel-case": "1.0.0"
	}
}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			TxtarString:     files,
		}).Build()

	b.AssertFileContent("public/js/main.js", `//# sourceMappingURL=main.js.map`)
	b.AssertFileContent("public/js/main.js.map", `"version":3`, "! ns-hugo")                                   // linked
	b.AssertFileContent("public/js/myts.js", `//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJz`) // inline
	b.AssertFileContent("public/index.html", `
		console.log(&#34;included&#34;);
		if (hasSpace.test(string))
		var React = __toESM(__require(&#34;react&#34;));
		function greeter(person) {
`)

	checkMap := func(p string, expectLen int) {
		s := b.FileContent(p)
		sources := esbuild.SourcesFromSourceMap(s)
		b.Assert(sources, qt.HasLen, expectLen)

		// Check that all source files exist.
		for _, src := range sources {
			filename, ok := paths.UrlStringToFilename(src)
			b.Assert(ok, qt.IsTrue)
			_, err := os.Stat(filename)
			b.Assert(err, qt.IsNil, qt.Commentf("src: %q", src))
		}
	}

	checkMap("public/js/main.js.map", 4)
}

func TestBuildError(t *testing.T) {
	c := qt.New(t)

	filesTemplate := `
-- config.toml --
disableKinds=["page", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
-- assets/js/main.js --
// A comment.
import { hello1, hello2 } from './util1';
hello1();
hello2();
-- assets/js/util1.js --
/* Some
comments.
*/
import { hello3 } from './util2';
export function hello1() {
	return 'abcd';
}
export function hello2() {
	return hello3();
}
-- assets/js/util2.js --
export function hello3() {
	return 'efgh';
}
-- layouts/index.html --
{{ $js := resources.Get "js/main.js" | js.Build }}
JS Content:{{ $js.Content }}:End:

			`

	c.Run("Import from main not found", func(c *qt.C) {
		c.Parallel()
		files := strings.Replace(filesTemplate, "import { hello1, hello2 }", "import { hello1, hello2, FOOBAR }", 1)
		b, err := hugolib.NewIntegrationTestBuilder(hugolib.IntegrationTestConfig{T: c, NeedsOsFS: true, TxtarString: files}).BuildE()
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, `main.js:2:25": No matching export`)
	})

	c.Run("Import from import not found", func(c *qt.C) {
		c.Parallel()
		files := strings.Replace(filesTemplate, "import { hello3 } from './util2';", "import { hello3, FOOBAR } from './util2';", 1)
		b, err := hugolib.NewIntegrationTestBuilder(hugolib.IntegrationTestConfig{T: c, NeedsOsFS: true, TxtarString: files}).BuildE()
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, `util1.js:4:17": No matching export in`)
	})
}

// See issue 10527.
func TestImportHugoVsESBuild(t *testing.T) {
	c := qt.New(t)

	for _, importSrcDir := range []string{"node_modules", "assets"} {
		c.Run(importSrcDir, func(c *qt.C) {
			files := `
-- IMPORT_SRC_DIR/imp1/index.js --
console.log("IMPORT_SRC_DIR:imp1/index.js");
-- IMPORT_SRC_DIR/imp2/index.ts --
console.log("IMPORT_SRC_DIR:imp2/index.ts");
-- IMPORT_SRC_DIR/imp3/foo.ts --
console.log("IMPORT_SRC_DIR:imp3/foo.ts");
-- assets/js/main.js --
import 'imp1/index.js';
import 'imp2/index.js';
import 'imp3/foo.js';
-- layouts/index.html --
{{ $js := resources.Get "js/main.js" | js.Build }}
{{ $js.RelPermalink }}
			`

			files = strings.ReplaceAll(files, "IMPORT_SRC_DIR", importSrcDir)

			b := hugolib.NewIntegrationTestBuilder(
				hugolib.IntegrationTestConfig{
					T:           c,
					NeedsOsFS:   true,
					TxtarString: files,
				}).Build()

			expected := `
IMPORT_SRC_DIR:imp1/index.js
IMPORT_SRC_DIR:imp2/index.ts
IMPORT_SRC_DIR:imp3/foo.ts
`
			expected = strings.ReplaceAll(expected, "IMPORT_SRC_DIR", importSrcDir)

			b.AssertFileContent("public/js/main.js", expected)
		})
	}
}

// See https://github.com/evanw/esbuild/issues/2745
func TestPreserveLegalComments(t *testing.T) {
	t.Parallel()

	files := `
-- assets/js/main.js --
/* @license
 * Main license.
 */
import * as foo from 'js/utils';
console.log("Hello Main");
-- assets/js/utils/index.js --
export * from './util1';
export * from './util2';
-- assets/js/utils/util1.js --
/*! License util1  */
console.log("Hello 1");
-- assets/js/utils/util2.js --
//! License util2  */
console.log("Hello 2");
-- layouts/index.html --
{{ $js := resources.Get "js/main.js" | js.Build (dict "minify" false) }}
{{ $js.RelPermalink }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			NeedsOsFS:   true,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/js/main.js", `
License util1
License util2
Main license

	`)
}

// Issue #11232
func TestTypeScriptExperimentalDecorators(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
-- tsconfig.json --
{
  "compilerOptions": {
    "experimentalDecorators": true,
  }
}
-- assets/ts/main.ts --
function addFoo(target: any) {target.prototype.foo = 'bar'}
@addFoo
class A {}
-- layouts/index.html --
{{ $opts := dict "target" "es2020" "targetPath" "js/main.js" }}
{{ (resources.Get "ts/main.ts" | js.Build $opts).Publish }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			NeedsOsFS:   true,
			TxtarString: files,
		}).Build()
	b.AssertFileContent("public/js/main.js", "__decorateClass")
}

// Issue 13183.
func TestExternalsInAssets(t *testing.T) {
	files := `
-- assets/js/util1.js --
export function hello1() {
	return 'abcd';
}
-- assets/js/util2.js --
export function hello2() {
	return 'efgh';
}
-- assets/js/main.js --
import { hello1 } from './util1.js';
import { hello2 } from './util2.js';

hello1();
hello2();
-- layouts/index.html --
Home.
{{ $js := resources.Get "js/main.js" | js.Build (dict "externals" (slice "./util1.js")) }}
{{ $js.Publish }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())

	b.AssertFileContent("public/js/main.js", "efgh")
	b.AssertFileContent("public/js/main.js", "! abcd")
}
