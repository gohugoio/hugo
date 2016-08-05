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

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadGlobalConfig(t *testing.T) {
	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContent := `
	PaginatePath = "side"
	`

	writeSource(t, "hugo.toml", configContent)

	require.NoError(t, LoadGlobalConfig("", "hugo.toml"))
	assert.Equal(t, "side", viper.GetString("PaginatePath"))
	// default
	assert.Equal(t, "layouts", viper.GetString("LayoutDir"))
}
