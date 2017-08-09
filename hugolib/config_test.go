// Copyright 2016-present The Hugo Authors. All rights reserved.
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

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContent := `
	PaginatePath = "side"
	`

	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	cfg, err := LoadConfig(mm, "", "hugo.toml")
	require.NoError(t, err)

	assert.Equal(t, "side", cfg.GetString("paginatePath"))
	// default
	assert.Equal(t, "layouts", cfg.GetString("layoutDir"))
}
func TestLoadMultiConfig(t *testing.T) {
	t.Parallel()

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContentBase := `
	DontChange = "same"
	PaginatePath = "side"
	`
	configContentSub := `
	PaginatePath = "top"
	`
	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "base.toml", configContentBase)

	writeToFs(t, mm, "override.toml", configContentSub)

	cfg, err := LoadConfig(mm, "", "base.toml,override.toml")
	require.NoError(t, err)

	assert.Equal(t, "top", cfg.GetString("paginatePath"))
	assert.Equal(t, "same", cfg.GetString("DontChange"))
}
