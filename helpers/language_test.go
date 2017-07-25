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

package helpers

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetGlobalOnlySetting(t *testing.T) {
	v := viper.New()
	lang := NewDefaultLanguage(v)
	lang.SetParam("defaultContentLanguageInSubdir", false)
	lang.SetParam("paginatePath", "side")
	v.Set("defaultContentLanguageInSubdir", true)
	v.Set("paginatePath", "page")

	require.True(t, lang.GetBool("defaultContentLanguageInSubdir"))
	require.Equal(t, "side", lang.GetString("paginatePath"))
}
