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

package cssjs_test

import (
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

func TestTailwindV4Basic(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	files := `
-- hugo.toml --
-- package.json --
{
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/bep/hugo-starter-tailwind-basic.git"
  },
  "devDependencies": {
    "@tailwindcss/cli": "^4.0.1",
    "tailwindcss": "^4.0.1"
  },
  "name": "hugo-starter-tailwind-basic",
  "version": "0.1.0"
}
-- assets/css/styles.css --
@import "tailwindcss";

@theme {
  --font-family-display: "Satoshi", "sans-serif";

  --breakpoint-3xl: 1920px;

  --color-neon-pink: oklch(71.7% 0.25 360);
  --color-neon-lime: oklch(91.5% 0.258 129);
  --color-neon-cyan: oklch(91.3% 0.139 195.8);
}
-- layouts/index.html --
{{ $css := resources.Get "css/styles.css" | css.TailwindCSS }}
CSS: {{ $css.Content | safeCSS }}|
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               t,
			TxtarString:     files,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			LogLevel:        logg.LevelInfo,
		}).Build()

	b.AssertFileContent("public/index.html", "/*! tailwindcss v4.")
}
