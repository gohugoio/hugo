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
package page_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestNextPrevConfig(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
-- content/mysection/_index.md --
-- content/mysection/p1.md --
---
title: "Page 1"
weight: 10
---
-- content/mysection/p2.md --
---
title: "Page 2"
weight: 20
---
-- content/mysection/p3.md --
---
title: "Page 3"
weight: 30
---
-- layouts/_default/single.html --
{{ .Title }}|Next: {{ with .Next}}{{ .Title}}{{ end }}|Prev: {{ with .Prev}}{{ .Title}}{{ end }}|NextInSection: {{ with .NextInSection}}{{ .Title}}{{ end }}|PrevInSection: {{ with .PrevInSection}}{{ .Title}}{{ end }}|

`
	b := hugolib.Test(t, filesTemplate)

	b.AssertFileContent("public/mysection/p1/index.html", "Page 1|Next: |Prev: Page 2|NextInSection: |PrevInSection: Page 2|")
	b.AssertFileContent("public/mysection/p2/index.html", "Page 2|Next: Page 1|Prev: Page 3|NextInSection: Page 1|PrevInSection: Page 3|")
	b.AssertFileContent("public/mysection/p3/index.html", "Page 3|Next: Page 2|Prev: |NextInSection: Page 2|PrevInSection: |")

	files := strings.ReplaceAll(filesTemplate, "-- hugo.toml --", `-- hugo.toml --
[page]
nextPrevSortOrder="aSc"
nextPrevInSectionSortOrder="asC"
`)

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/mysection/p1/index.html", "Page 1|Next: Page 2|Prev: |NextInSection: Page 2|PrevInSection: |")
	b.AssertFileContent("public/mysection/p2/index.html", "Page 2|Next: Page 3|Prev: Page 1|NextInSection: Page 3|PrevInSection: Page 1|")
	b.AssertFileContent("public/mysection/p3/index.html", "Page 3|Next: |Prev: Page 2|NextInSection: |PrevInSection: Page 2|")

	files = strings.ReplaceAll(filesTemplate, "-- hugo.toml --", `-- hugo.toml --
[page]
nextPrevSortOrder="aSc"
`)

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/mysection/p1/index.html", "Page 1|Next: Page 2|Prev: |NextInSection: |PrevInSection: Page 2|")
	b.AssertFileContent("public/mysection/p2/index.html", "Page 2|Next: Page 3|Prev: Page 1|NextInSection: Page 1|PrevInSection: Page 3|")
	b.AssertFileContent("public/mysection/p3/index.html", "Page 3|Next: |Prev: Page 2|NextInSection: Page 2|PrevInSection: |")

	files = strings.ReplaceAll(filesTemplate, "-- hugo.toml --", `-- hugo.toml --
[page]
nextPrevInSectionSortOrder="aSc"
`)

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/mysection/p1/index.html", "Page 1|Next: |Prev: Page 2|NextInSection: Page 2|PrevInSection: |")
}
