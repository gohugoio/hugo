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

package media

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultTypes(t *testing.T) {
	require.Equal(t, "text", HTMLType.MainType)
	require.Equal(t, "html", HTMLType.SubType)
	require.Equal(t, "html", HTMLType.Suffix)

	require.Equal(t, "text/html", HTMLType.MainType())
	require.Equal(t, "text/html+html", HTMLType.String())

	require.Equal(t, "application", RSSType.MainType)
	require.Equal(t, "rss", RSSType.SubType)
	require.Equal(t, "xml", RSSType.Suffix)

	require.Equal(t, "application/rss", RSSType.MainType())
	require.Equal(t, "application/rss+xml", RSSType.String())

}
