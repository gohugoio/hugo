// Copyright 2021 The Hugo Authors. All rights reserved.
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

package babel_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
)

func TestTransformBabel(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	files := `
-- assets/js/main.js --
/* A Car */
class Car {
	constructor(brand) {
	this.carname = brand;
	}
}
-- assets/js/main2.js --
/* A Car2 */
class Car2 {
	constructor(brand) {
	this.carname = brand;
	}
}
-- babel.config.js --
console.error("Hugo Environment:", process.env.HUGO_ENVIRONMENT );

module.exports = {
	presets: ["@babel/preset-env"],
};
-- hugo.toml --
disablekinds = ['taxonomy', 'term', 'page']
[security]
	[security.exec]
	allow = ['^node$', '^babel$']
-- layouts/home.html --
{{ $options := dict "noComments" true }}
{{ $transpiled := resources.Get "js/main.js" | babel -}}
Transpiled: {{ $transpiled.Content | safeJS }}

{{ $transpiled := resources.Get "js/main2.js" | babel (dict "sourceMap" "inline") -}}
Transpiled2: {{ $transpiled.Content | safeJS }}

{{ $transpiled := resources.Get "js/main2.js" | babel (dict "sourceMap" "external") -}}
Transpiled3: {{ $transpiled.Permalink }}
-- package.json --
{
	"scripts": {},

	"devDependencies": {
	"@babel/cli": "7.28.6",
	"@babel/core": "7.28.6",
	"@babel/preset-env": "7.28.6"
	}
}

	`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall(), hugolib.TestOptInfo())

	b.AssertLogContains("babel: Hugo Environment: production")
	b.AssertFileContent("public/index.html", `var Car2 =`)
	b.AssertFileContent("public/js/main2.js", `var Car2 =`)
	b.AssertFileContent("public/js/main2.js.map", `{"version":3,`)
	b.AssertFileContent("public/index.html", `
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozL`)
}

// See Issue 15043.
// See Issue 15040.
func TestTransformBabelConfigResolution(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[[module.imports]]
path = "github.com/bep/hugo-mod-nop"
[security]
  [security.exec]
  allow = ['^go$', '^node$', '^babel$']
-- assets/js/main.js --
/* A Car */
class Car {
  constructor(brand) {
    this.carname = brand;
  }
}
-- layouts/home.html --
{{ $js := resources.Get "js/main.js" | babel }}
RelPermalink: {{ $js.RelPermalink }}|HasBody: {{ gt (len $js.Content) 0 }}|
-- package.json --
{
  "devDependencies": {
    "@babel/cli": "7.28.6",
    "@babel/core": "7.28.6"
  }
}
-- go.mod --
module github.com/example/project

go 1.26

replace github.com/bep/hugo-mod-nop => ../external-module
-- ../external-module/go.mod --
module github.com/bep/hugo-mod-nop

go 1.26
-- ../external-module/CONFIG_FILE_NAME --
CONFIG_FILE_CONTENT
	`

	tests := []struct {
		name              string
		configFileName    string
		configFileContent string
	}{
		{
			name:              "mjs in module",
			configFileName:    "babel.config.mjs",
			configFileContent: "export default {};\n",
		},
		{
			name:              "cjs in module",
			configFileName:    "babel.config.cjs",
			configFileContent: "module.exports = {};\n",
		},
		{
			name:              "js in module",
			configFileName:    "babel.config.js",
			configFileContent: "module.exports = {};\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			rootDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
			c.Assert(err, qt.IsNil)
			c.Cleanup(clean)

			projectDir := filepath.Join(rootDir, "project")
			moduleDir := filepath.Join(rootDir, "external-module")
			c.Assert(os.MkdirAll(projectDir, 0o755), qt.IsNil)
			c.Assert(os.MkdirAll(moduleDir, 0o755), qt.IsNil)

			f := strings.ReplaceAll(files, "CONFIG_FILE_NAME", tt.configFileName)
			f = strings.ReplaceAll(f, "CONFIG_FILE_CONTENT", tt.configFileContent)

			b := hugolib.Test(c, f,
				hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
					cfg.WorkingDir = projectDir
					cfg.NeedsOsFS = true
					cfg.NeedsNpmInstall = true
				}),
				hugolib.TestOptInfo(),
			)

			b.AssertFileContent("public/index.html",
				"RelPermalink: /js/main.js|HasBody: true|",
			)
			b.AssertLogContains(tt.configFileName)
		})
	}
}
