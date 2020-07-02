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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
)

func TestJS_Build(t *testing.T) {
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

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-babel")
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
	{{ $options := dict "minify" true }}
	{{ $transpiled := resources.Get "js/main.js" | js.Build $options }}
	Built: {{ $transpiled.Content | safeJS }}
	`)

	jsDir := filepath.Join(workDir, "assets", "js")
	b.Assert(os.MkdirAll(jsDir, 0777), qt.IsNil)
	b.WithSourceFile("assets/js/main.js", mainJS)
	b.WithSourceFile("assets/js/included.js", includedJS)

	_, err = captureStdout(func() error {
		return b.BuildE(BuildCfg{})
	})
	b.Assert(err, qt.IsNil)

	b.AssertFileContent("public/index.html", `
  Built: (()=&gt;{console.log(&#34;included&#34;);console.log(&#34;main&#34;);})();
	`)

}
