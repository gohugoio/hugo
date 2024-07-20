// Copyright 2019 The Hugo Authors. All rights reserved.
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

package tplimpl_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"

	"github.com/gohugoio/hugo/tpl/internal"
)

func TestTemplateFuncsExamples(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds=["home", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
ignoreErrors = ["my-err-id"]
[outputs]
home=["HTML"]
-- layouts/partials/header.html --
<title>Hugo Rocks!</title>
-- files/README.txt --
Hugo Rocks!
-- content/blog/hugo-rocks.md --
---
title: "**BatMan**"
---
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	d := b.H.Sites[0].Deps

	var (
		templates []string
		expected  []string
	)

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(d)
		for _, mm := range ns.MethodMappings {
			for _, example := range mm.Examples {
				// These will fail the build, so skip.
				if strings.Contains(example[0], "errorf") ||
					strings.Contains(example[0], "transform.XMLEscape") ||
					strings.Contains(example[0], "math.Rand") {
					continue
				}
				templates = append(templates, example[0])
				expected = append(expected, example[1])
			}
		}
	}

	files += fmt.Sprintf("-- layouts/_default/single.html --\n%s\n", strings.Join(templates, "\n"))
	b = hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/blog/hugo-rocks/index.html", expected...)
}
