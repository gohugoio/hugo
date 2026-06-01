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
// See the License for the specific language gxoverning permissions and
// limitations under the License.

package npm_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/modules/npm"
	"github.com/spf13/afero"
)

func getPackageBuilderTestFiles() string {
	files := `
-- hugo.toml --
-- package.json --
{
  "workspaces": [
    "packages/*"
  ]
}
-- packages/hugoautogen/package.json --
-- packages/a/package.json --
PACKAGE_CONTENT
-- packages/b/package.json --
PACKAGE_CONTENT
-- packages/c/package.json --
PACKAGE_CONTENT
-- packages/d/package.json --
PACKAGE_CONTENT
-- packages/e/package.json --
PACKAGE_CONTENT
`
	packageContent := `{
"name": "foo",
"version": "0.1.1",
"dependencies": {
	"react-dom": "1.1.1",
	"tailwindcss": "1.2.0",
	"@babel/cli": "7.8.4",
	"@babel/core": "7.9.0",
	"@babel/preset-env": "7.9.5"
},
"devDependencies": {
	"postcss-cli": "7.1.0",
	"tailwindcss": "1.2.0",
	"@babel/cli": "7.8.4",
	"@babel/core": "7.9.0",
	"@babel/preset-env": "7.9.5"
}
}`
	files = strings.ReplaceAll(files, "PACKAGE_CONTENT", packageContent)
	return files
}

func TestPackageBuilder(t *testing.T) {
	files := getPackageBuilderTestFiles()
	b := hugolib.Test(t, files)
	fs := b.H.Fs.WorkingDirReadOnly
	sourceFs := b.H.BaseFs.ProjectSourceFs
	assetsFs := b.H.BaseFs.AssetsWithDuplicatesPreserved.Fs
	mods := b.H.Configs.Modules

	sum := npm.PackageFilesSum(fs, b.H.AllModules())
	b.Assert(sum, qt.Equals, "528bb6507993888c")

	b.Assert(npm.Pack(sourceFs, assetsFs, mods), qt.IsNil)
	b.Assert(npm.NpmPackNeedsUpdate(sourceFs, mods), qt.IsFalse)

	meta1, err := afero.ReadFile(sourceFs, "packages/hugoautogen/hugo_packagemeta.json")
	b.Assert(err, qt.IsNil)

	b.Assert(npm.Pack(sourceFs, assetsFs, mods), qt.IsNil)
	b.Assert(npm.NpmPackNeedsUpdate(sourceFs, mods), qt.IsFalse)

	meta2, err := afero.ReadFile(sourceFs, "packages/hugoautogen/hugo_packagemeta.json")
	b.Assert(err, qt.IsNil)
	b.Assert(string(meta2), qt.Equals, string(meta1))
}

func BenchmarkPackageFilesSum(b *testing.B) {
	files := getPackageBuilderTestFiles()
	bb := hugolib.Test(b, files)
	fs := bb.H.Fs.WorkingDirReadOnly
	b.ResetTimer()

	for b.Loop() {
		sum := npm.PackageFilesSum(fs, modules.Modules{})
		bb.Assert(sum, qt.Equals, "528bb6507993888c")
	}
}
