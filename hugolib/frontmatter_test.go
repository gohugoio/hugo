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

package hugolib

import "testing"

// Issue 10624
func TestFrontmatterPreserveDatatypesForSlices(t *testing.T) {
	t.Parallel()

	files := `
-- content/post/one.md --
---
ints: [1, 2, 3]
mixed: ["1", 2, 3]
strings: ["1", "2","3"]
---	
-- layouts/_default/single.html --
Ints: {{ printf "%T" .Params.ints }} {{ range .Params.ints }}Int: {{ fmt.Printf "%[1]v (%[1]T)" . }}|{{ end }}
Mixed: {{ printf "%T" .Params.mixed }} {{ range .Params.mixed }}Mixed: {{ fmt.Printf "%[1]v (%[1]T)" . }}|{{ end }}
Strings: {{ printf "%T" .Params.strings }} {{ range .Params.strings }}Strings: {{ fmt.Printf "%[1]v (%[1]T)" . }}|{{ end }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	b.Build()

	b.AssertFileContent("public/post/one/index.html", "Ints: []interface {} Int: 1 (int)|Int: 2 (int)|Int: 3 (int)|")
	b.AssertFileContent("public/post/one/index.html", "Mixed: []interface {} Mixed: 1 (string)|Mixed: 2 (int)|Mixed: 3 (int)|")
	b.AssertFileContent("public/post/one/index.html", "Strings: []string Strings: 1 (string)|Strings: 2 (string)|Strings: 3 (string)|")
}
