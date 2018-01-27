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

		}},
		{filepath.FromSlash("/a/"), filepath.FromSlash("/a/b/c/d/page.md"), func(f *FileInfo) {
			assert.Equal("b", f.Section())

		}},
	} {
		f := s.NewFileInfo(this.base, this.filename, false, nil)
		this.assert(f)
	}

}
