// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package source_test

import (@@ -247,11 +247,38 @@ i18n: {{ T "hugo" }}
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/source"
)

func TestFileInfo(t *testing.T) {func TestNewContentForce(t *testing.T) {

	mm := afero.NewMemMapFs()

	c := qt.New(t)

	archetypeDir := filepath.Join("archetypes", "my-bundle")

	c.Assert(mm.MkdirAll(archetypeDir, 0o755), qt.IsNil)

	c.Assert(afero.WriteFile(mm, filepath.Join(archetypeDir, "index.md"), []byte(""), 0o755), qt.IsNil)

	c.Assert(afero.WriteFile(mm, filepath.Join(archetypeDir, "index.nn.md"), []byte(""), 0o755), qt.IsNil)

	c.Assert(initFs(mm), qt.IsNil)

	cfg, fs := newTestCfg(c, mm)

	h, err := hugolib.NewHugoSites(deps.DepsCfg{Cfg: cfg, Fs: fs})

	c.Assert(err, qt.IsNil)

	c.Assert(len(h.Sites), qt.Equals, 2)

	// from file

	c.Assert(create.NewContent(h, "post", "post/my-post.md", false), qt.IsNil)

	c.Assert(create.NewContent(h, "post", "post/my-post.md", false), qt.IsNotNil)

	c.Assert(create.NewContent(h, "post", "post/my-post.md", true), qt.IsNil)

	// from dir

	c.Assert(create.NewContent(h, "my-bundle", "post/my-post", false), qt.IsNil)

	c.Assert(create.NewContent(h, "my-bundle", "post/my-post", false), qt.IsNotNil)

	c.Assert(create.NewContent(h, "my-bundle", "post/my-post", true), qt.IsNil)

}
	c := qt.New(t)

	s := newTestSourceSpec()

	for _, this := range []struct {
		base     string
		filename string
		assert   func(f *source.FileInfo)
	}{
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/page.md"), func(f *source.FileInfo) {
			c.Assert(f.Filename(), qt.Equals, filepath.FromSlash("/a/b/page.md"))
			c.Assert(f.Dir(), qt.Equals, filepath.FromSlash("b/"))
			c.Assert(f.Path(), qt.Equals, filepath.FromSlash("b/page.md"))
			c.Assert(f.Section(), qt.Equals, "b")
			c.Assert(f.TranslationBaseName(), qt.Equals, filepath.FromSlash("page"))
			c.Assert(f.BaseFileName(), qt.Equals, filepath.FromSlash("page"))
		}},
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/c/d/page.md"), func(f *source.FileInfo) {
			c.Assert(f.Section(), qt.Equals, "b")
		}},
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/page.en.MD"), func(f *source.FileInfo) {
			c.Assert(f.Section(), qt.Equals, "b")
			c.Assert(f.Path(), qt.Equals, filepath.FromSlash("b/page.en.MD"))
			c.Assert(f.TranslationBaseName(), qt.Equals, filepath.FromSlash("page"))
			c.Assert(f.BaseFileName(), qt.Equals, filepath.FromSlash("page.en"))
		}},
	} {
		path := strings.TrimPrefix(this.filename, this.base)
		f, err := s.NewFileInfoFrom(path, this.filename)
		c.Assert(err, qt.IsNil)
		this.assert(f)
	}
}
