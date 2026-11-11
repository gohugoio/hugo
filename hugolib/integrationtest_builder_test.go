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

package hugolib_test

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/hugolib"
)

// See issue 15568.
func TestNpmInstallPreservesWorkingDirectory(t *testing.T) {
	if !hexec.InPath("npm") {
		t.Skip("npm not installed")
	}

	wd, err := os.Getwd()
	qt.Assert(t, err, qt.IsNil)

	files := `
-- hugo.toml --
disableKinds = ['rss', 'sitemap', 'taxonomy', 'term']
-- layouts/home.html --
Home.
-- package.json --
{
  "name": "hugo-npm-working-directory-test",
  "private": true,
  "scripts": {
    "install": "node -e \"require('fs').writeFileSync('npm-installed.txt', 'installed')\""
  }
}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall())
	b.AssertFileContent("npm-installed.txt", "installed")
	b.AssertFileContent("public/index.html", "Home.")

	got, err := os.Getwd()
	b.Assert(err, qt.IsNil)
	b.Assert(got, qt.Equals, wd)
}
