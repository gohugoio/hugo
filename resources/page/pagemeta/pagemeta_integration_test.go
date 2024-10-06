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
	"strings"
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

func TestDateValidation(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- content/_index.md --
FRONT_MATTER
-- layouts/index.html --
{{ .Date.UTC.Format "2006-01-02" }}
--
`
	errorMsg := `ERROR the "date" front matter field is not a parsable date`

	// TOML: unquoted date/time (valid)
	f := strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = 2024-10-01
+++
	`)
	b := hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "2024-10-01")

	// TOML: string (valid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = "2024-10-01"
+++
	`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "2024-10-01")

	// TOML: empty string (valid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = ""
+++
	`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "0001-01-01")

	// TOML: int (valid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = 0
+++
	`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "1970-01-01")

	// TOML: string (invalid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = "2024-42-42"
+++
	`)
	b, _ = hugolib.TestE(t, f)
	b.AssertLogContains(errorMsg)

	// TOML: bool (invalid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = true
+++
	`)
	b, _ = hugolib.TestE(t, f)
	b.AssertLogContains(errorMsg)

	// TOML: float (invalid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
+++
date = 6.7
+++
	`)
	b, _ = hugolib.TestE(t, f)
	b.AssertLogContains(errorMsg)

	// JSON: null (valid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
{
  "date": null
}
	`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "0001-01-01")

	// YAML: null (valid)
	f = strings.ReplaceAll(files, "FRONT_MATTER", `
---
date:
---
	`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "0001-01-01")
}
