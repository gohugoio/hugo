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

		require.Equal(t, test.expectedType, test.tp.Type())
		require.Equal(t, test.expectedString, test.tp.String())

	}

}
