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

	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewPathSpecFromConfig(t *testing.T) {
	v := viper.New()
	l := NewLanguage("no", v)
	v.Set("disablePathToLower", true)
	v.Set("removePathAccents", true)
	v.Set("uglyURLs", true)
	v.Set("multilingual", true)
	v.Set("defaultContentLanguageInSubdir", true)
	v.Set("defaultContentLanguage", "no")
	v.Set("canonifyURLs", true)
	v.Set("paginatePath", "side")
	v.Set("baseURL", "http://base.com")
	v.Set("themesDir", "thethemes")
	v.Set("layoutDir", "thelayouts")
	v.Set("workingDir", "thework")
	v.Set("staticDir", "thestatic")
	v.Set("theme", "thetheme")

	p, err := NewPathSpec(hugofs.NewMem(v), l)

	require.NoError(t, err)
	require.True(t, p.canonifyURLs)
	require.True(t, p.defaultContentLanguageInSubdir)
	require.True(t, p.disablePathToLower)
	require.True(t, p.multilingual)
	require.True(t, p.removePathAccents)
	require.True(t, p.uglyURLs)
	require.Equal(t, "no", p.defaultContentLanguage)
	require.Equal(t, "no", p.language.Lang)
	require.Equal(t, "side", p.paginatePath)

	require.Equal(t, "http://base.com", p.BaseURL.String())
	require.Equal(t, "thethemes", p.themesDir)
	require.Equal(t, "thelayouts", p.layoutDir)
	require.Equal(t, "thework", p.workingDir)
	require.Equal(t, "thetheme", p.theme)
}
