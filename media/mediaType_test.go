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
	for _, test := range []struct {
		tp               Type
		expectedMainType string
		expectedSubType  string
		expectedSuffix   string
		expectedType     string
		expectedString   string
	}{
		{CalendarType, "text", "calendar", "ics", "text/calendar", "text/calendar"},
		{CSSType, "text", "css", "css", "text/css", "text/css"},
		{SCSSType, "text", "x-scss", "scss", "text/x-scss", "text/x-scss"},
		{CSVType, "text", "csv", "csv", "text/csv", "text/csv"},
		{HTMLType, "text", "html", "html", "text/html", "text/html"},
		{JavascriptType, "application", "javascript", "js", "application/javascript", "application/javascript"},
		{JSONType, "application", "json", "json", "application/json", "application/json"},
		{RSSType, "application", "rss", "xml", "application/rss+xml", "application/rss+xml"},
		{SVGType, "image", "svg", "svg", "image/svg+xml", "image/svg+xml"},
		{TextType, "text", "plain", "txt", "text/plain", "text/plain"},
		{XMLType, "application", "xml", "xml", "application/xml", "application/xml"},
		{TOMLType, "application", "toml", "toml", "application/toml", "application/toml"},
		{YAMLType, "application", "yaml", "yaml", "application/yaml", "application/yaml"},
	} {
		require.Equal(t, test.expectedMainType, test.tp.MainType)
		require.Equal(t, test.expectedSubType, test.tp.SubType)
		require.Equal(t, test.expectedSuffix, test.tp.Suffix(), test.tp.String())
		require.Equal(t, defaultDelimiter, test.tp.Delimiter)

		require.Equal(t, test.expectedType, test.tp.Type())
		require.Equal(t, test.expectedString, test.tp.String())

	}

	require.Equal(t, 15, len(DefaultTypes))

}

func TestGetByType(t *testing.T) {
	types := Types{HTMLType, RSSType}

	mt, found := types.GetByType("text/HTML")
	require.True(t, found)
	require.Equal(t, mt, HTMLType)

	_, found = types.GetByType("text/nono")
	require.False(t, found)

	mt, found = types.GetByType("application/rss+xml")
	require.True(t, found)
	require.Equal(t, mt, RSSType)

	mt, found = types.GetByType("application/rss")
	require.True(t, found)
	require.Equal(t, mt, RSSType)
}

func TestGetByMainSubType(t *testing.T) {
	assert := require.New(t)
	f, found := DefaultTypes.GetByMainSubType("text", "plain")
	assert.True(found)
	assert.Equal(f, TextType)
	_, found = DefaultTypes.GetByMainSubType("foo", "plain")
	assert.False(found)
}

func TestBySuffix(t *testing.T) {
	assert := require.New(t)
	formats := DefaultTypes.BySuffix("xml")
	assert.Equal(2, len(formats))
	assert.Equal("rss", formats[0].SubType)
	assert.Equal("xml", formats[1].SubType)
}

func TestGetFirstBySuffix(t *testing.T) {
	assert := require.New(t)
	f, found := DefaultTypes.GetFirstBySuffix("xml")
	assert.True(found)
	assert.Equal(Type{MainType: "application", SubType: "rss", mimeSuffix: "xml", Delimiter: ".", Suffixes: []string{"xml"}, fileSuffix: "xml"}, f)
}

func TestFromTypeString(t *testing.T) {
	f, err := fromString("text/html")
	require.NoError(t, err)
	require.Equal(t, HTMLType.Type(), f.Type())

	f, err = fromString("application/custom")
	require.NoError(t, err)
	require.Equal(t, Type{MainType: "application", SubType: "custom", mimeSuffix: "", fileSuffix: ""}, f)

	f, err = fromString("application/custom+sfx")
	require.NoError(t, err)
	require.Equal(t, Type{MainType: "application", SubType: "custom", mimeSuffix: "sfx"}, f)

	_, err = fromString("noslash")
	require.Error(t, err)

	f, err = fromString("text/xml; charset=utf-8")
	require.NoError(t, err)
	require.Equal(t, Type{MainType: "text", SubType: "xml", mimeSuffix: ""}, f)
	require.Equal(t, "", f.Suffix())
}

// Add a test for the SVG case
// https://github.com/gohugoio/hugo/issues/4920
func TestFromExtensionMultipleSuffixes(t *testing.T) {
	assert := require.New(t)
	tp, found := DefaultTypes.GetBySuffix("svg")
	assert.True(found)
	assert.Equal("image/svg+xml", tp.String())
	assert.Equal("svg", tp.fileSuffix)
	assert.Equal(".svg", tp.FullSuffix())
	tp, found = DefaultTypes.GetByType("image/svg+xml")
	assert.True(found)
	assert.Equal("image/svg+xml", tp.String())
	assert.True(found)
	assert.Equal(".svg", tp.FullSuffix())

}

func TestDecodeTypes(t *testing.T) {

	var tests = []struct {
		name        string
		maps        []map[string]interface{}
		shouldError bool
		assert      func(t *testing.T, name string, tt Types)
	}{
		{
			"Redefine JSON",
			[]map[string]interface{}{
				{
					"application/json": map[string]interface{}{
						"suffixes": []string{"jasn"}}}},
			false,
			func(t *testing.T, name string, tt Types) {
				require.Len(t, tt, len(DefaultTypes))
				json, found := tt.GetBySuffix("jasn")
				require.True(t, found)
				require.Equal(t, "application/json", json.String(), name)
				require.Equal(t, ".jasn", json.FullSuffix())
			}},
		{
			"MIME suffix in key, multiple file suffixes, custom delimiter",
			[]map[string]interface{}{
				{
					"application/hugo+hg": map[string]interface{}{
						"suffixes":  []string{"hg1", "hg2"},
						"Delimiter": "_",
					}}},
			false,
			func(t *testing.T, name string, tt Types) {
				require.Len(t, tt, len(DefaultTypes)+1)
				hg, found := tt.GetBySuffix("hg2")
				require.True(t, found)
				require.Equal(t, "hg", hg.mimeSuffix)
				require.Equal(t, "hg2", hg.Suffix())
				require.Equal(t, "_hg2", hg.FullSuffix())
				require.Equal(t, "application/hugo+hg", hg.String(), name)

				hg, found = tt.GetByType("application/hugo+hg")
				require.True(t, found)

			}},
		{
			"Add custom media type",
			[]map[string]interface{}{
				{
					"text/hugo+hgo": map[string]interface{}{
						"Suffixes": []string{"hgo2"}}}},
			false,
			func(t *testing.T, name string, tt Types) {
				require.Len(t, tt, len(DefaultTypes)+1)
				// Make sure we have not broken the default config.

				_, found := tt.GetBySuffix("json")
				require.True(t, found)

				hugo, found := tt.GetBySuffix("hgo2")
				require.True(t, found)
				require.Equal(t, "text/hugo+hgo", hugo.String(), name)
			}},
	}

	for _, test := range tests {
		result, err := DecodeTypes(test.maps...)
		if test.shouldError {
			require.Error(t, err, test.name)
		} else {
			require.NoError(t, err, test.name)
			test.assert(t, test.name, result)
		}
	}
}
