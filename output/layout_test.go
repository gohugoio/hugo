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

	"github.com/stretchr/testify/require"
)

type testLayoutIdentifier struct {
	pageKind    string
	pageSection string
	pageLayout  string
	pageType    string
}

func (l testLayoutIdentifier) PageKind() string {
	return l.pageKind
}

func (l testLayoutIdentifier) PageLayout() string {
	return l.pageLayout
}

func (l testLayoutIdentifier) PageType() string {
	return l.pageType
}

func (l testLayoutIdentifier) PageSection() string {
	return l.pageSection
}

func TestLayout(t *testing.T) {

	for i, this := range []struct {
		li             testLayoutIdentifier
		hasTheme       bool
		layoutOverride string
		tp             Type
		expect         []string
	}{
		{testLayoutIdentifier{"home", "", "", ""}, true, "", HTMLType,
			[]string{"index.html", "_default/list.html", "theme/index.html", "theme/_default/list.html"}},
		{testLayoutIdentifier{"section", "sect1", "", ""}, false, "", HTMLType,
			[]string{"section/sect1.html", "sect1/list.html"}},
		{testLayoutIdentifier{"taxonomy", "tag", "", ""}, false, "", HTMLType,
			[]string{"taxonomy/tag.html", "indexes/tag.html"}},
		{testLayoutIdentifier{"taxonomyTerm", "categories", "", ""}, false, "", HTMLType,
			[]string{"taxonomy/categories.terms.html", "_default/terms.html"}},
		{testLayoutIdentifier{"page", "", "", ""}, true, "", HTMLType,
			[]string{"_default/single.html", "theme/_default/single.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", ""}, false, "", HTMLType,
			[]string{"_default/mylayout.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", "myttype"}, false, "", HTMLType,
			[]string{"myttype/mylayout.html", "_default/mylayout.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", "myttype/mysubtype"}, false, "", HTMLType,
			[]string{"myttype/mysubtype/mylayout.html", "myttype/mylayout.html", "_default/mylayout.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", "myttype"}, false, "myotherlayout", HTMLType,
			[]string{"myttype/myotherlayout.html", "_default/myotherlayout.html"}},
	} {
		l := NewLayoutHandler(this.hasTheme)
		logMsg := fmt.Sprintf("Test %d", i)
		layouts := l.For(this.li, this.layoutOverride, this.tp)
		require.NotNil(t, layouts, logMsg)
		require.True(t, len(layouts) >= len(this.expect), logMsg)
		// Not checking the complete list for now ...
		require.Equal(t, this.expect, layouts[:len(this.expect)], logMsg)

		if !this.hasTheme {
			for _, layout := range layouts {
				require.NotContains(t, layout, "theme", logMsg)
			}
		}
	}
}
