// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hexec_test

import (
	"testing"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

func TestNPMGlobalInstalls(t *testing.T) {
	if !htesting.IsRealCI() {
		t.Skip("We only ever want to run this in CI.")
	}
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[security.exec]
allow = ['^(babel|node|postcss|tailwindcss)$']

-- package.json --
{}
-- hugo_stats.json --
-- assets/js/main.js --
console.log("Hello, world!");
-- assets/css/main1.css --
body { color: red }
-- assets/css/main2.css --
@import "tailwindcss";
@plugin "@tailwindcss/typography";
@source "hugo_stats.json";
body { color: blue }
-- layouts/home.html --
{{ with resources.Get "css/main1.css"  }}
	{{ with . | css.PostCSS  }}
		 CSS1 size: {{ .Content | len }}|{{ .RelPermalink }}|
	{{ end }}
{{ end }}
 {{ with resources.Get "css/main2.css"  }}
	{{ with . | css.TailwindCSS  }}
		 CSS2 size: {{ .Content | len }}|{{ .RelPermalink }}|
	{{ end }}
{{ end }}
{{ with resources.Get "js/main.js"  }}
	{{ with . | js.Babel  }}
		 JS size: {{ .Content | len }}|{{ .RelPermalink }}|
	{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstallGlobal(
		"postcss", "postcss-cli",
		"@babel/core", "@babel/cli",
		"tailwindcss", "@tailwindcss/cli", "@tailwindcss/typography",
	))

	b.AssertFileContent("public/index.html",
		"CSS1 size: 233|/css/main1.css|",
		"CSS2 size: 4557|/css/main2.css|",
		"JS size: 31|/js/main.js|",
	)
}
