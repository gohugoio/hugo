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

package pagemeta_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestLastModEq(t *testing.T) {
	files := `
-- hugo.toml --
timeZone = "Europe/London"
-- content/p1.md --
---
title: p1
date: 2024-03-13T06:00:00
---
-- layouts/_default/single.html --
Date: {{ .Date }}
Lastmod: {{ .Lastmod }}
Eq: {{ eq .Date .Lastmod }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
Date: 2024-03-13 06:00:00 &#43;0000 GMT
Lastmod: 2024-03-13 06:00:00 &#43;0000 GMT
Eq: true
`)
}
