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

package resources_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestTransformCached(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]	
-- assets/css/main.css --
body {
	  background: #fff;
}
-- content/p1.md --
---
title: "P1"
---
P1.
-- content/p2.md --
---
title: "P2"
---
P2.
-- layouts/_default/list.html --
List.
-- layouts/_default/single.html --
{{ $css := resources.Get "css/main.css" | resources.Minify  }}
CSS: {{ $css.Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "CSS: body{background:#fff}")
}
