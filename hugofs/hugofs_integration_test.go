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

package hugofs_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestMountRestrictTheme(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
theme = "mytheme"
[[module.mounts]]
source = '../file2.txt'
target = 'assets/file2.txt'
-- themes/mytheme/hugo.toml --
[[module.mounts]]
source = '../../file1.txt'
target = 'assets/file1.txt'
-- file1.txt --
file1
-- file2.txt --
file2
-- layouts/all.html --
All.
`
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "mount source must be a local path for modules/themes")
}

// Issue 14089.
func TestMountNodeMoudulesFromTheme(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
theme = "mytheme"
-- node_modules/bootstrap/foo.txt --
foo project.
-- layouts/all.html --
{{ $foo := resources.Get "vendor/bootstrap/foo.txt" }}
Foo: {{ with $foo }}{{ .Content }}{{ else }}Fail{{ end }}
-- themes/mytheme/hugo.toml --
[[module.mounts]]
source = 'NODE_MODULES_SOURCE' # tries first in theme, then in project root
target = 'assets/vendor/bootstrap'

`
	runFiles := func(files string) *hugolib.IntegrationTestBuilder {
		return hugolib.Test(t, files, hugolib.TestOptOsFs())
	}
	files := strings.ReplaceAll(filesTemplate, "NODE_MODULES_SOURCE", "node_modules/bootstrap")
	b := runFiles(files)
	b.AssertFileContent("public/index.html", "Foo: foo project.")

	// This is for backwards compatibility. ../../node_modules/bootstrap works exactly the same as node_modules/bootstrap.
	files = strings.ReplaceAll(filesTemplate, "NODE_MODULES_SOURCE", "../../node_modules/bootstrap")
	b = runFiles(files)
	b.AssertFileContent("public/index.html", "Foo: foo project.")

	files = strings.ReplaceAll(filesTemplate, "NODE_MODULES_SOURCE", "node_modules/bootstrap")
	files += `
-- themes/mytheme/node_modules/bootstrap/foo.txt --
foo theme.
`

	b = runFiles(files)
	b.AssertFileContent("public/index.html", "Foo: foo theme.")
}
