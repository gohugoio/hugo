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

package output

import (
	"testing"

	"github.com/spf13/hugo/media"
	"github.com/stretchr/testify/require"
)

func TestDefaultTypes(t *testing.T) {
	require.Equal(t, "HTML", HTMLType.Name)
	require.Equal(t, media.HTMLType, HTMLType.MediaType)
	require.Empty(t, HTMLType.Path)
	require.False(t, HTMLType.IsPlainText)

	require.Equal(t, "RSS", RSSType.Name)
	require.Equal(t, media.RSSType, RSSType.MediaType)
	require.Empty(t, RSSType.Path)
	require.False(t, RSSType.IsPlainText)
	require.True(t, RSSType.NoUgly)
}

func TestGetType(t *testing.T) {
	tp, _ := GetFormat("html")
	require.Equal(t, HTMLType, tp)
	tp, _ = GetFormat("HTML")
	require.Equal(t, HTMLType, tp)
	_, found := GetFormat("FOO")
	require.False(t, found)
}
