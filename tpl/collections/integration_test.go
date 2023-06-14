// Copyright 2022 The Hugo Authors. All rights reserved.
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

package collections_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue 9585
func TestApplyWithContext(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
{{ apply (seq 3) "partial" "foo.html"}}
-- layouts/partials/foo.html --
{{ return "foo"}}
  `

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
	[foo foo foo]
`)
}

// Issue 9865
func TestSortStable(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- layouts/index.html --
{{ $values := slice (dict "a" 1 "b" 2) (dict "a" 3 "b" 1) (dict "a" 2 "b" 0) (dict "a" 1 "b" 0) (dict "a" 3 "b" 1) (dict "a" 2 "b" 2) (dict "a" 2 "b" 1) (dict "a" 0 "b" 3) (dict "a" 3 "b" 3) (dict "a" 0 "b" 0) (dict "a" 0 "b" 0) (dict "a" 2 "b" 0) (dict "a" 1 "b" 2) (dict "a" 1 "b" 1) (dict "a" 3 "b" 0) (dict "a" 2 "b" 0) (dict "a" 3 "b" 0) (dict "a" 3 "b" 0) (dict "a" 3 "b" 0) (dict "a" 3 "b" 1) }}
Asc:  {{ sort (sort $values "b" "asc") "a" "asc" }}
Desc: {{ sort (sort $values "b" "desc") "a" "desc" }}

  `

	for i := 0; i < 4; i++ {

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.AssertFileContent("public/index.html", `
Asc:  [map[a:0 b:0] map[a:0 b:0] map[a:0 b:3] map[a:1 b:0] map[a:1 b:1] map[a:1 b:2] map[a:1 b:2] map[a:2 b:0] map[a:2 b:0] map[a:2 b:0] map[a:2 b:1] map[a:2 b:2] map[a:3 b:0] map[a:3 b:0] map[a:3 b:0] map[a:3 b:0] map[a:3 b:1] map[a:3 b:1] map[a:3 b:1] map[a:3 b:3]]
Desc: [map[a:3 b:3] map[a:3 b:1] map[a:3 b:1] map[a:3 b:1] map[a:3 b:0] map[a:3 b:0] map[a:3 b:0] map[a:3 b:0] map[a:2 b:2] map[a:2 b:1] map[a:2 b:0] map[a:2 b:0] map[a:2 b:0] map[a:1 b:2] map[a:1 b:2] map[a:1 b:1] map[a:1 b:0] map[a:0 b:3] map[a:0 b:0] map[a:0 b:0]]
`)

	}
}

// Issue #11004.
func TestAppendSliceToASliceOfSlices(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/index.html --
{{ $obj := slice (slice "a") }}
{{ $obj = $obj | append (slice "b") }}
{{ $obj = $obj | append (slice "c") }}

{{ $obj }}


  `

	for i := 0; i < 4; i++ {

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.AssertFileContent("public/index.html", "[[a] [b] [c]]")

	}

}
