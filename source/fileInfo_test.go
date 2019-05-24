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

package source

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFileInfo(t *testing.T) {
	assert := require.New(t)

	s := newTestSourceSpec()

	for _, this := range []struct {
		base     string
		filename string
		assert   func(f *FileInfo)
	}{
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/page.md"), func(f *FileInfo) {
			assert.Equal(filepath.FromSlash("/a/b/page.md"), f.Filename())
			assert.Equal(filepath.FromSlash("b/"), f.Dir())
			assert.Equal(filepath.FromSlash("b/page.md"), f.Path())
			assert.Equal("b", f.Section())
			assert.Equal(filepath.FromSlash("page"), f.TranslationBaseName())
			assert.Equal(filepath.FromSlash("page"), f.BaseFileName())

		}},
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/c/d/page.md"), func(f *FileInfo) {
			assert.Equal("b", f.Section())

		}},
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/page.en.MD"), func(f *FileInfo) {
			assert.Equal("b", f.Section())
			assert.Equal(filepath.FromSlash("b/page.en.MD"), f.Path())
			assert.Equal(filepath.FromSlash("page"), f.TranslationBaseName())
			assert.Equal(filepath.FromSlash("page.en"), f.BaseFileName())

		}},
	} {
		f := s.NewFileInfo(this.base, this.filename, false, nil)
		this.assert(f)
	}

}

func TestFileInfoLanguage(t *testing.T) {
	assert := require.New(t)
	langs := map[string]bool{
		"sv": true,
		"en": true,
	}

	m := afero.NewMemMapFs()
	lfs := hugofs.NewLanguageFs("sv", langs, m)
	v := newTestConfig()

	fs := hugofs.NewFrom(m, v)

	ps, err := helpers.NewPathSpec(fs, v)
	assert.NoError(err)
	s := SourceSpec{SourceFs: lfs, PathSpec: ps}
	s.Languages = map[string]interface{}{
		"en": true,
	}

	err = afero.WriteFile(lfs, "page.md", []byte("abc"), 0777)
	assert.NoError(err)
	err = afero.WriteFile(lfs, "page.en.md", []byte("abc"), 0777)
	assert.NoError(err)

	sv, _ := lfs.Stat("page.md")
	en, _ := lfs.Stat("page.en.md")

	fiSv := s.NewFileInfo("", "page.md", false, sv)
	fiEn := s.NewFileInfo("", "page.en.md", false, en)

	assert.Equal("sv", fiSv.Lang())
	assert.Equal("en", fiEn.Lang())

	// test contentBaseName implementation
	fi := s.NewFileInfo("", "2018-10-01-contentbasename.md", false, nil)
	assert.Equal("2018-10-01-contentbasename", fi.ContentBaseName())

	fi = s.NewFileInfo("", "2018-10-01-contentbasename.en.md", false, nil)
	assert.Equal("2018-10-01-contentbasename", fi.ContentBaseName())

	fi = s.NewFileInfo("", filepath.Join("2018-10-01-contentbasename", "index.en.md"), true, nil)
	assert.Equal("2018-10-01-contentbasename", fi.ContentBaseName())

	fi = s.NewFileInfo("", filepath.Join("2018-10-01-contentbasename", "_index.en.md"), false, nil)
	assert.Equal("_index", fi.ContentBaseName())
}
