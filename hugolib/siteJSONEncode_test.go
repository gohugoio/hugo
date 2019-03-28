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

package hugolib

import (
	"testing"
)

// Issue #1123
// Testing prevention of cyclic refs in JSON encoding
// May be smart to run with: -timeout 4000ms
func TestEncodePage(t *testing.T) {
	t.Parallel()

	templ := `Page: |{{ index .Site.RegularPages 0 | jsonify }}|
Site: {{ site | jsonify }}
`

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded("index.html", templ)
	b.WithContent("page.md", `---
title: "Page"
date: 2019-02-28
---

Content.

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `"Date":"2019-02-28T00:00:00Z"`)

}
