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
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/media"
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

func TestGetFormatByName(t *testing.T) {
	formats := Formats{AMPFormat, CalendarFormat}
	tp, _ := formats.GetByName("AMp")
	require.Equal(t, AMPFormat, tp)
	_, found := formats.GetByName("HTML")
	require.False(t, found)
	_, found = formats.GetByName("FOO")
	require.False(t, found)
}

func TestGetFormatByExt(t *testing.T) {
	formats1 := Formats{AMPFormat, CalendarFormat}
	formats2 := Formats{AMPFormat, HTMLFormat, CalendarFormat}
	tp, _ := formats1.GetBySuffix("html")
	require.Equal(t, AMPFormat, tp)
	tp, _ = formats1.GetBySuffix("ics")
	require.Equal(t, CalendarFormat, tp)
	_, found := formats1.GetBySuffix("not")
	require.False(t, found)

	// ambiguous
	_, found = formats2.GetBySuffix("html")
	require.False(t, found)
}

func TestGetFormatByFilename(t *testing.T) {
	noExtNoDelimMediaType := media.TextType
	noExtNoDelimMediaType.Suffix = ""
	noExtNoDelimMediaType.Delimiter = ""

	noExtMediaType := media.TextType
	noExtMediaType.Suffix = ""

	var (
		noExtDelimFormat = Format{
			Name:      "NEM",
			MediaType: noExtNoDelimMediaType,
			BaseName:  "_redirects",
		}
		noExt = Format{
			Name:      "NEX",
			MediaType: noExtMediaType,
			BaseName:  "next",
		}
	)

	formats := Formats{AMPFormat, HTMLFormat, noExtDelimFormat, noExt, CalendarFormat}
	f, found := formats.FromFilename("my.amp.html")
	require.True(t, found)
	require.Equal(t, AMPFormat, f)
	f, found = formats.FromFilename("my.ics")
	require.True(t, found)
	f, found = formats.FromFilename("my.html")
	require.True(t, found)
	require.Equal(t, HTMLFormat, f)
	f, found = formats.FromFilename("my.nem")
	require.True(t, found)
	require.Equal(t, noExtDelimFormat, f)
	f, found = formats.FromFilename("my.nex")
	require.True(t, found)
	require.Equal(t, noExt, f)
	_, found = formats.FromFilename("my.css")
	require.False(t, found)

}

func TestDecodeFormats(t *testing.T) {

	mediaTypes := media.Types{media.JSONType, media.XMLType}

	var tests = []struct {
		name        string
		maps        []map[string]interface{}
		shouldError bool
		assert      func(t *testing.T, name string, f Formats)
	}{
		{
			"Redefine JSON",
			[]map[string]interface{}{
				{
					"JsON": map[string]interface{}{
						"baseName":    "myindex",
						"isPlainText": "false"}}},
			false,
			func(t *testing.T, name string, f Formats) {
				require.Len(t, f, len(DefaultFormats), name)
				json, _ := f.GetByName("JSON")
				require.Equal(t, "myindex", json.BaseName)
				require.Equal(t, media.JSONType, json.MediaType)
				require.False(t, json.IsPlainText)

			}},
		{
			"Add XML format with string as mediatype",
			[]map[string]interface{}{
				{
					"MYXMLFORMAT": map[string]interface{}{
						"baseName":  "myxml",
						"mediaType": "application/xml",
					}}},
			false,
			func(t *testing.T, name string, f Formats) {
				require.Len(t, f, len(DefaultFormats)+1, name)
				xml, found := f.GetByName("MYXMLFORMAT")
				require.True(t, found)
				require.Equal(t, "myxml", xml.BaseName, fmt.Sprint(xml))
				require.Equal(t, media.XMLType, xml.MediaType)

				// Verify that we haven't changed the DefaultFormats slice.
				json, _ := f.GetByName("JSON")
				require.Equal(t, "index", json.BaseName, name)

			}},
		{
			"Add format unknown mediatype",
			[]map[string]interface{}{
				{
					"MYINVALID": map[string]interface{}{
						"baseName":  "mymy",
						"mediaType": "application/hugo",
					}}},
			true,
			func(t *testing.T, name string, f Formats) {

			}},
		{
			"Add and redefine XML format",
			[]map[string]interface{}{
				{
					"MYOTHERXMLFORMAT": map[string]interface{}{
						"baseName":  "myotherxml",
						"mediaType": media.XMLType,
					}},
				{
					"MYOTHERXMLFORMAT": map[string]interface{}{
						"baseName": "myredefined",
					}},
			},
			false,
			func(t *testing.T, name string, f Formats) {
				require.Len(t, f, len(DefaultFormats)+1, name)
				xml, found := f.GetByName("MYOTHERXMLFORMAT")
				require.True(t, found)
				require.Equal(t, "myredefined", xml.BaseName, fmt.Sprint(xml))
				require.Equal(t, media.XMLType, xml.MediaType)
			}},
	}

	for _, test := range tests {
		result, err := DecodeFormats(mediaTypes, test.maps...)
		if test.shouldError {
			require.Error(t, err, test.name)
		} else {
			require.NoError(t, err, test.name)
			test.assert(t, test.name, result)
		}
	}
}
