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
	require.Equal(t, "Calendar", CalendarFormat.Name)
	require.Equal(t, media.CalendarType, CalendarFormat.MediaType)
	require.Equal(t, "webcal://", CalendarFormat.Protocol)
	require.Empty(t, CalendarFormat.Path)
	require.True(t, CalendarFormat.IsPlainText)
	require.False(t, CalendarFormat.IsHTML)

	require.Equal(t, "CSS", CSSFormat.Name)
	require.Equal(t, media.CSSType, CSSFormat.MediaType)
	require.Empty(t, CSSFormat.Path)
	require.Empty(t, CSSFormat.Protocol) // Will inherit the BaseURL protocol.
	require.True(t, CSSFormat.IsPlainText)
	require.False(t, CSSFormat.IsHTML)

	require.Equal(t, "CSV", CSVFormat.Name)
	require.Equal(t, media.CSVType, CSVFormat.MediaType)
	require.Empty(t, CSVFormat.Path)
	require.Empty(t, CSVFormat.Protocol)
	require.True(t, CSVFormat.IsPlainText)
	require.False(t, CSVFormat.IsHTML)

	require.Equal(t, "HTML", HTMLFormat.Name)
	require.Equal(t, media.HTMLType, HTMLFormat.MediaType)
	require.Empty(t, HTMLFormat.Path)
	require.Empty(t, HTMLFormat.Protocol)
	require.False(t, HTMLFormat.IsPlainText)
	require.True(t, HTMLFormat.IsHTML)

	require.Equal(t, "AMP", AMPFormat.Name)
	require.Equal(t, media.HTMLType, AMPFormat.MediaType)
	require.Equal(t, "amp", AMPFormat.Path)
	require.Empty(t, AMPFormat.Protocol)
	require.False(t, AMPFormat.IsPlainText)
	require.True(t, AMPFormat.IsHTML)

	require.Equal(t, "RSS", RSSFormat.Name)
	require.Equal(t, media.RSSType, RSSFormat.MediaType)
	require.Empty(t, RSSFormat.Path)
	require.False(t, RSSFormat.IsPlainText)
	require.True(t, RSSFormat.NoUgly)
	require.False(t, CalendarFormat.IsHTML)

}

func TestGetFormat(t *testing.T) {
	tp, _ := GetFormat("html")
	require.Equal(t, HTMLFormat, tp)
	tp, _ = GetFormat("HTML")
	require.Equal(t, HTMLFormat, tp)
	_, found := GetFormat("FOO")
	require.False(t, found)
}

func TestGeGetFormatByName(t *testing.T) {
	formats := Formats{AMPFormat, CalendarFormat}
	tp, _ := formats.GetByName("AMP")
	require.Equal(t, AMPFormat, tp)
	_, found := formats.GetByName("HTML")
	require.False(t, found)
	_, found = formats.GetByName("FOO")
	require.False(t, found)
}

func TestGeGetFormatByExt(t *testing.T) {
	formats1 := Formats{AMPFormat, CalendarFormat}
	formats2 := Formats{AMPFormat, HTMLFormat, CalendarFormat}
	tp, _ := formats1.GetBySuffix("html")
	require.Equal(t, AMPFormat, tp)
	tp, _ = formats1.GetBySuffix("ics")
	require.Equal(t, CalendarFormat, tp)
	_, found := formats1.GetBySuffix("not")
	require.False(t, found)

	// ambiguous
	_, found = formats2.GetByName("html")
	require.False(t, found)
}
