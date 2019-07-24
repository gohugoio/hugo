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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsContentFile(t *testing.T) {
	assert := require.New(t)

	assert.True(IsContentFile(filepath.FromSlash("my/file.md")))
	assert.True(IsContentFile(filepath.FromSlash("my/file.ad")))
	assert.False(IsContentFile(filepath.FromSlash("textfile.txt")))
	assert.True(IsContentExt("md"))
	assert.False(IsContentExt("json"))
}

func TestComponentFolders(t *testing.T) {
	assert := require.New(t)

	// It's important that these are absolutely right and not changed.
	assert.Equal(len(ComponentFolders), len(componentFoldersSet))
	assert.True(IsComponentFolder("archetypes"))
	assert.True(IsComponentFolder("layouts"))
	assert.True(IsComponentFolder("data"))
	assert.True(IsComponentFolder("i18n"))
	assert.True(IsComponentFolder("assets"))
	assert.False(IsComponentFolder("resources"))
	assert.True(IsComponentFolder("static"))
	assert.True(IsComponentFolder("content"))
	assert.False(IsComponentFolder("foo"))
	assert.False(IsComponentFolder(""))
}
