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
		{CalendarType, "text", "calendar", "ics", "text/calendar", "text/calendar+ics"},
		{CSSType, "text", "css", "css", "text/css", "text/css+css"},
		{CSVType, "text", "csv", "csv", "text/csv", "text/csv+csv"},
		{HTMLType, "text", "html", "html", "text/html", "text/html+html"},
		{JavascriptType, "application", "javascript", "js", "application/javascript", "application/javascript+js"},
		{JSONType, "application", "json", "json", "application/json", "application/json+json"},
		{RSSType, "application", "rss", "xml", "application/rss", "application/rss+xml"},
		{TextType, "text", "plain", "txt", "text/plain", "text/plain+txt"},
	} {
		require.Equal(t, test.expectedMainType, test.tp.MainType)
		require.Equal(t, test.expectedSubType, test.tp.SubType)
		require.Equal(t, test.expectedSuffix, test.tp.Suffix)
		require.Equal(t, defaultDelimiter, test.tp.Delimiter)

		require.Equal(t, test.expectedType, test.tp.Type())
		require.Equal(t, test.expectedString, test.tp.String())

	}

}

func TestGetByType(t *testing.T) {
	types := Types{HTMLType, RSSType}

	mt, found := types.GetByType("text/HTML")
	require.True(t, found)
	require.Equal(t, mt, HTMLType)

	_, found = types.GetByType("text/nono")
	require.False(t, found)
}

func TestFromTypeString(t *testing.T) {
	f, err := FromString("text/html")
	require.NoError(t, err)
	require.Equal(t, HTMLType, f)

	f, err = FromString("application/custom")
	require.NoError(t, err)
	require.Equal(t, Type{MainType: "application", SubType: "custom", Suffix: "custom", Delimiter: defaultDelimiter}, f)

	f, err = FromString("application/custom+pdf")
	require.NoError(t, err)
	require.Equal(t, Type{MainType: "application", SubType: "custom", Suffix: "pdf", Delimiter: defaultDelimiter}, f)

	_, err = FromString("noslash")
	require.Error(t, err)

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
						"suffix": "jsn"}}},
			false,
			func(t *testing.T, name string, tt Types) {
				require.Len(t, tt, len(DefaultTypes))
				json, found := tt.GetBySuffix("jsn")
				require.True(t, found)
				require.Equal(t, "application/json+jsn", json.String(), name)
			}},
		{
			"Add custom media type",
			[]map[string]interface{}{
				{
					"text/hugo": map[string]interface{}{
						"suffix": "hgo"}}},
			false,
			func(t *testing.T, name string, tt Types) {
				require.Len(t, tt, len(DefaultTypes)+1)
				// Make sure we have not broken the default config.
				_, found := tt.GetBySuffix("json")
				require.True(t, found)

				hugo, found := tt.GetBySuffix("hgo")
				require.True(t, found)
				require.Equal(t, "text/hugo+hgo", hugo.String(), name)
			}},
		{
			"Add media type invalid key",
			[]map[string]interface{}{
				{
					"text/hugo+hgo": map[string]interface{}{}}},
			true,
			func(t *testing.T, name string, tt Types) {

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
