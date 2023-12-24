// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
)

func TestMountFilters(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t)
	workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-mountfilters")
	b.Assert(err, qt.IsNil)
	defer clean()

	for _, component := range files.ComponentFolders {
		b.Assert(os.MkdirAll(filepath.Join(workingDir, component), 0o777), qt.IsNil)
	}
	b.WithWorkingDir(workingDir).WithLogger(loggers.NewDefault())
	b.WithConfigFile("toml", fmt.Sprintf(`
workingDir = %q

[module]
[[module.mounts]]
source = 'content'
target = 'content'
excludeFiles = "/a/c/**"
[[module.mounts]]
source = 'static'
target = 'static'
[[module.mounts]]
source = 'layouts'
target = 'layouts'
excludeFiles = "/**/foo.html"
[[module.mounts]]
source = 'data'
target = 'data'
includeFiles = "/mydata/**"
[[module.mounts]]
source = 'assets'
target = 'assets'
excludeFiles = ["/**exclude.*", "/moooo.*"]
[[module.mounts]]
source = 'i18n'
target = 'i18n'
[[module.mounts]]
source = 'archetypes'
target = 'archetypes'

	
`, workingDir))

	b.WithContent("/a/b/p1.md", "---\ntitle: Include\n---")
	b.WithContent("/a/c/p2.md", "---\ntitle: Exclude\n---")

	b.WithSourceFile(
		"data/mydata/b.toml", `b1='bval'`,
		"data/nodata/c.toml", `c1='bval'`,
		"layouts/partials/foo.html", `foo`,
		"assets/exclude.txt", `foo`,
		"assets/js/exclude.js", `foo`,
		"assets/js/include.js", `foo`,
		"assets/js/exclude.js", `foo`,
	)

	b.WithTemplatesAdded("index.html", `

Data: {{ site.Data }}:END

Template: {{ templates.Exists "partials/foo.html" }}:END
Resource1: {{ resources.Get "js/include.js" }}:END
Resource2: {{ resources.Get "js/exclude.js" }}:END
Resource3: {{ resources.Get "exclude.txt" }}:END
Resources: {{ resources.Match "**.js" }}
`)

	b.Build(BuildCfg{})

	assertExists := func(name string, shouldExist bool) {
		b.Helper()
		b.Assert(b.CheckExists(name), qt.Equals, shouldExist)
	}

	assertExists("public/a/b/p1/index.html", true)
	assertExists("public/a/c/p2/index.html", false)

	b.AssertFileContent(filepath.Join("public", "index.html"), `
Data: map[mydata:map[b:map[b1:bval]]]:END	
Template: false
Resource1: /js/include.js:END
Resource2: :END
Resource3: :END
Resources: [/js/include.js]
`)
}
