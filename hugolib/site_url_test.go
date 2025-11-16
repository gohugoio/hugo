// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"testing"
)

func TestSectionsEntries(t *testing.T) {
	files := `
-- hugo.toml --
-- content/withfile/_index.md --
-- content/withoutfile/p1.md --
-- layouts/_default/list.html --
SectionsEntries: {{ .SectionsEntries }}


`

	b := Test(t, files)

	b.AssertFileContent("public/withfile/index.html", "SectionsEntries: [withfile]")
	b.AssertFileContent("public/withoutfile/index.html", "SectionsEntries: [withoutfile]")
}
