// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	"strings"

	"github.com/gohugoio/hugo/deps"
)

func TestAllTemplateEngines(t *testing.T) {
	t.Parallel()
	noOp := func(s string) string {
		return s
	}

	amberFixer := func(s string) string {
		fixed := strings.Replace(s, "{{ .Title", "{{ Title", -1)
		fixed = strings.Replace(fixed, ".Content", "Content", -1)
		fixed = strings.Replace(fixed, ".IsNamedParams", "IsNamedParams", -1)
		fixed = strings.Replace(fixed, "{{", "#{", -1)
		fixed = strings.Replace(fixed, "}}", "}", -1)
		fixed = strings.Replace(fixed, `title "hello world"`, `title("hello world")`, -1)

		return fixed
	}

	for _, config := range []struct {
		suffix        string
		templateFixer func(s string) string
	}{
		{"amber", amberFixer},
		{"html", noOp},
		{"ace", noOp},
	} {
		t.Run(config.suffix,
			func(t *testing.T) {
				doTestTemplateEngine(t, config.suffix, config.templateFixer)
			})
	}

}

func doTestTemplateEngine(t *testing.T, suffix string, templateFixer func(s string) string) {

	cfg, fs := newTestCfg()

	t.Log("Testing", suffix)

	templTemplate := `
p
	|
	| Page Title: {{ .Title }}
	br
	| Page Content: {{ .Content }}
	br
	| {{ title "hello world" }}

`

	templShortcodeTemplate := `
p
	|
	| Shortcode: {{ .IsNamedParams }}
`

	templ := templateFixer(templTemplate)
	shortcodeTempl := templateFixer(templShortcodeTemplate)

	writeSource(t, fs, filepath.Join("content", "p.md"), `
---
title: My Title 
---
My Content

Shortcode: {{< myShort >}}

`)

	writeSource(t, fs, filepath.Join("layouts", "_default", fmt.Sprintf("single.%s", suffix)), templ)
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", fmt.Sprintf("myShort.%s", suffix)), shortcodeTempl)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileContent(filepath.Join("public", "p", "index.html"),
		"Page Title: My Title",
		"My Content",
		"Hello World",
		"Shortcode: false",
	)

}
