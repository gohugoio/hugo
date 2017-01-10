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

	"github.com/spf13/hugo/hugofs"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewPathSpecFromConfig(t *testing.T) {
	viper.Set("disablePathToLower", true)
	viper.Set("removePathAccents", true)
	viper.Set("uglyURLs", true)
	viper.Set("multilingual", true)
	viper.Set("defaultContentLanguageInSubdir", true)
	viper.Set("defaultContentLanguage", "no")
	viper.Set("currentContentLanguage", NewLanguage("no"))
	viper.Set("canonifyURLs", true)
	viper.Set("paginatePath", "side")

	p := NewPathSpec(hugofs.NewMem(), viper.GetViper())

	require.True(t, p.canonifyURLs)
	require.True(t, p.defaultContentLanguageInSubdir)
	require.True(t, p.disablePathToLower)
	require.True(t, p.multilingual)
	require.True(t, p.removePathAccents)
	require.True(t, p.uglyURLs)
	require.Equal(t, "no", p.defaultContentLanguage)
	require.Equal(t, "no", p.currentContentLanguage.Lang)
	require.Equal(t, "side", p.paginatePath)
}
