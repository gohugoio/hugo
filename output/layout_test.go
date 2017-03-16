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

var ampType = Format{
	Name:      "AMP",
	MediaType: media.HTMLType,
	BaseName:  "index",
}

func TestLayout(t *testing.T) {

	for _, this := range []struct {
		name           string
		li             testLayoutIdentifier
		hasTheme       bool
		layoutOverride string
		tp             Format
		expect         []string
	}{
		{"Home", testLayoutIdentifier{"home", "", "", ""}, true, "", ampType,
			[]string{"index.amp.html", "index.html", "_default/list.amp.html", "_default/list.html", "theme/index.amp.html", "theme/index.html"}},
		{"Section", testLayoutIdentifier{"section", "sect1", "", ""}, false, "", ampType,
			[]string{"section/sect1.amp.html", "section/sect1.html"}},
		{"Taxonomy", testLayoutIdentifier{"taxonomy", "tag", "", ""}, false, "", ampType,
			[]string{"taxonomy/tag.amp.html", "taxonomy/tag.html"}},
		{"Taxonomy term", testLayoutIdentifier{"taxonomyTerm", "categories", "", ""}, false, "", ampType,
			[]string{"taxonomy/categories.terms.amp.html", "taxonomy/categories.terms.html", "_default/terms.amp.html"}},
		{"Page", testLayoutIdentifier{"page", "", "", ""}, true, "", ampType,
			[]string{"_default/single.amp.html", "_default/single.html", "theme/_default/single.amp.html"}},
		{"Page with layout", testLayoutIdentifier{"page", "", "mylayout", ""}, false, "", ampType,
			[]string{"_default/mylayout.amp.html", "_default/mylayout.html"}},
		{"Page with layout and type", testLayoutIdentifier{"page", "", "mylayout", "myttype"}, false, "", ampType,
			[]string{"myttype/mylayout.amp.html", "myttype/mylayout.html", "_default/mylayout.amp.html"}},
		{"Page with layout and type with subtype", testLayoutIdentifier{"page", "", "mylayout", "myttype/mysubtype"}, false, "", ampType,
			[]string{"myttype/mysubtype/mylayout.amp.html", "myttype/mysubtype/mylayout.html", "myttype/mylayout.amp.html"}},
		{"Page with overridden layout", testLayoutIdentifier{"page", "", "mylayout", "myttype"}, false, "myotherlayout", ampType,
			[]string{"myttype/myotherlayout.amp.html", "myttype/myotherlayout.html"}},
	} {
		t.Run(this.name, func(t *testing.T) {
			l := NewLayoutHandler(this.hasTheme)

			layouts := l.For(this.li, this.layoutOverride, this.tp)

			require.NotNil(t, layouts)
			require.True(t, len(layouts) >= len(this.expect))
			// Not checking the complete list for now ...
			require.Equal(t, this.expect, layouts[:len(this.expect)])

			if !this.hasTheme {
				for _, layout := range layouts {
					require.NotContains(t, layout, "theme")
				}
			}
		})
	}
}
