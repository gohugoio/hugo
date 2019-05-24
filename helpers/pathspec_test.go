// Copyright 2018 The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/langs"
	"github.com/stretchr/testify/require"
)

func TestNewPathSpecFromConfig(t *testing.T) {
	v := newTestCfg()
	l := langs.NewLanguage("no", v)
	v.Set("disablePathToLower", true)
	v.Set("removePathAccents", true)
	v.Set("uglyURLs", true)
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
	require.True(t, p.CanonifyURLs)
	require.True(t, p.DisablePathToLower)
	require.True(t, p.RemovePathAccents)
	require.True(t, p.UglyURLs)
	require.Equal(t, "no", p.Language.Lang)
	require.Equal(t, "side", p.PaginatePath)

	require.Equal(t, "http://base.com", p.BaseURL.String())
	require.Equal(t, "thethemes", p.ThemesDir)
	require.Equal(t, "thework", p.WorkingDir)
	require.Equal(t, []string{"thetheme"}, p.Themes())
}
