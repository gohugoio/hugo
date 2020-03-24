// Copyright 2019 The Hugo Authors. All rights reserved.
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

package files

import (
	"os"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestIsContentFile(t *testing.T) {
	c := qt.New(t)

	c.Assert(IsContentFile(filepath.FromSlash("my/file.md")), qt.Equals, true)
	c.Assert(IsContentFile(filepath.FromSlash("my/file.ad")), qt.Equals, true)
	c.Assert(IsContentFile(filepath.FromSlash("textfile.txt")), qt.Equals, false)
	c.Assert(IsContentExt("md"), qt.Equals, true)
	c.Assert(IsContentExt("json"), qt.Equals, false)
}

// newFs returns a filesystem with the specified files at the specified filepaths.
func newFs(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()

	for filename, content := range files {
		err := afero.WriteFile(fs, filename, []byte(content), os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create file %q: %v", filename, err)
		}
	}

	return fs
}

func TestIsHTMLContent(t *testing.T) {
	c := qt.New(t)

	opener := func(t *testing.T, s string) func() (afero.File, error) {
		fs := newFs(t, map[string]string{"temp": s})
		return func() (afero.File, error) {
			return fs.Open("/temp")
		}
	}

	c.Assert(hasFrontMatter(opener(t, "   <html>")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, "   <!--\n---")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, "   <!--")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, "   ---<")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, " foo  <")), qt.Equals, false)

	c.Assert(hasFrontMatter(opener(t, "---\ntitle: dummy\n---")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, "+++\ntitle= \"dummy\"\n+++")), qt.Equals, false)
	c.Assert(hasFrontMatter(opener(t, "{\n\"title\": \"dummy\"\n}")), qt.Equals, false)
}

func TestComponentFolders(t *testing.T) {
	c := qt.New(t)

	// It's important that these are absolutely right and not changed.
	c.Assert(len(componentFoldersSet), qt.Equals, len(ComponentFolders))
	c.Assert(IsComponentFolder("archetypes"), qt.Equals, true)
	c.Assert(IsComponentFolder("layouts"), qt.Equals, true)
	c.Assert(IsComponentFolder("data"), qt.Equals, true)
	c.Assert(IsComponentFolder("i18n"), qt.Equals, true)
	c.Assert(IsComponentFolder("assets"), qt.Equals, true)
	c.Assert(IsComponentFolder("resources"), qt.Equals, false)
	c.Assert(IsComponentFolder("static"), qt.Equals, true)
	c.Assert(IsComponentFolder("content"), qt.Equals, true)
	c.Assert(IsComponentFolder("foo"), qt.Equals, false)
	c.Assert(IsComponentFolder(""), qt.Equals, false)

}

func TestClassifyContentFile(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		filename string
		content  string
		class    ContentClass
	}{
		{
			filename: "_index.md",
			content:  "---\ntitle:Top\n---",
			class:    ContentClassBranch,
		}, {
			filename: "sub/_index.md",
			content:  "---\ntitle:Other\n---",
			class:    ContentClassBranch,
		}, {
			filename: "blog/post/index.md",
			content:  "---\ntitle:Today\n---\n<html></html>",
			class:    ContentClassLeaf,
		}, {
			filename: "blog/post/yesterday.md",
			content:  "---\ntitle:Yesterday\n---",
			class:    ContentClassContent,
		}, {
			filename: "blog/post/static/index.html",
			content:  "<html></html>",
			class:    ContentClassFile,
		}, {
			filename: "blog/post/static/image.png",
			content:  "",
			class:    ContentClassFile,
		},
	}

	fsMap := make(map[string]string, len(tests))
	for _, test := range tests {
		fsMap[test.filename] = test.content
	}
	fs := newFs(t, fsMap)

	opener := func(s string) func() (afero.File, error) {
		return func() (afero.File, error) { return fs.Open(s) }
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			got := ClassifyContentFile(filepath.Base(test.filename), opener(test.filename))
			c.Check(got, qt.Equals, test.class, qt.Commentf("%q", test.filename))
		})
	}
}
