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

package hugolib

import (
	"testing"

	"path/filepath"

	"github.com/gohugoio/hugo/source"
	"github.com/stretchr/testify/require"
)

func TestBundleFileInfo(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	cfg, fs := newTestBundleSourcesMultilingual(t)
	sourceSpec := source.NewSourceSpec(cfg, fs)

	for _, this := range []struct {
		filename string
		check    func(f *fileInfo)
	}{
		{"/path/to/file.md", func(fi *fileInfo) {
			assert.Equal("md", fi.Ext())
			assert.Equal("en", fi.Lang())
			assert.False(fi.isOwner())
			assert.True(fi.isContentFile())
		}},
		{"/path/to/file.JPG", func(fi *fileInfo) {
			assert.Equal("jpg", fi.Ext())
			assert.False(fi.isContentFile())
		}},
		{"/path/to/file.nn.png", func(fi *fileInfo) {
			assert.Equal("png", fi.Ext())
			assert.Equal("nn", fi.Lang())
			assert.Equal("file", fi.TranslationBaseName())
			assert.False(fi.isContentFile())
		}},
	} {
		fi := newFileInfo(
			sourceSpec,
			filepath.FromSlash("/work/base"),
			filepath.FromSlash(this.filename),
			nil, bundleNot)
		this.check(fi)
	}

}
