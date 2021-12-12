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
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/config"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/htesting"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
)

func TestResourceChainBabel(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	packageJSON := `{
  "scripts": {},

  "devDependencies": {
    "@babel/cli": "7.8.4",
    "@babel/core": "7.9.0",	
    "@babel/preset-env": "7.9.5"
  }
}
`

	babelConfig := `
console.error("Hugo Environment:", process.env.HUGO_ENVIRONMENT );

module.exports = {
  presets: ["@babel/preset-env"],
};

`

	js := `
/* A Car */
class Car {
  constructor(brand) {
    this.carname = brand;
  }
}
`

	js2 := `
/* A Car2 */
class Car2 {
  constructor(brand) {
    this.carname = brand;
  }
}
`

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-babel")
	c.Assert(err, qt.IsNil)
	defer clean()

	var logBuf bytes.Buffer
	logger := loggers.NewBasicLoggerForWriter(jww.LevelInfo, &logBuf)

	v := config.New()
	v.Set("workingDir", workDir)
	v.Set("disableKinds", []string{"taxonomy", "term", "page"})
	v.Set("security", map[string]interface{}{
		"exec": map[string]interface{}{
			"allow": []string{"^npx$", "^babel$"},
		},
	})

	b := newTestSitesBuilder(t).WithLogger(logger)

	// Need to use OS fs for this.
	b.Fs = hugofs.NewDefault(v)
	b.WithWorkingDir(workDir)
	b.WithViper(v)
	b.WithContent("p1.md", "")

	b.WithTemplates("index.html", `
{{ $options := dict "noComments" true }}
{{ $transpiled := resources.Get "js/main.js" | babel -}}
Transpiled: {{ $transpiled.Content | safeJS }}

{{ $transpiled := resources.Get "js/main2.js" | babel (dict "sourceMap" "inline") -}}
Transpiled2: {{ $transpiled.Content | safeJS }}

{{ $transpiled := resources.Get "js/main2.js" | babel (dict "sourceMap" "external") -}}
Transpiled3: {{ $transpiled.Permalink }}

`)

	jsDir := filepath.Join(workDir, "assets", "js")
	b.Assert(os.MkdirAll(jsDir, 0777), qt.IsNil)
	b.WithSourceFile("assets/js/main.js", js)
	b.WithSourceFile("assets/js/main2.js", js2)
	b.WithSourceFile("package.json", packageJSON)
	b.WithSourceFile("babel.config.js", babelConfig)

	b.Assert(os.Chdir(workDir), qt.IsNil)
	cmd := b.NpmInstall()
	err = cmd.Run()
	b.Assert(err, qt.IsNil)

	b.Build(BuildCfg{})

	// Make sure Node sees this.
	b.Assert(logBuf.String(), qt.Contains, "babel: Hugo Environment: production")
	b.Assert(err, qt.IsNil)

	b.AssertFileContent("public/index.html", `var Car =`)
	b.AssertFileContent("public/index.html", `var Car2 =`)
	b.AssertFileContent("public/js/main2.js", `var Car2 =`)
	b.AssertFileContent("public/js/main2.js.map", `{"version":3,`)
	b.AssertFileContent("public/index.html", `
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozL`)
}
