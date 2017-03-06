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
	l := &Layout{}

	for _, this := range []struct {
		li     testLayoutIdentifier
		tp     Type
		expect []string
	}{
		{testLayoutIdentifier{"home", "", "", ""}, HTMLType, []string{"index.html", "_default/list.html", "theme/index.html", "theme/_default/list.html"}},
		{testLayoutIdentifier{"section", "sect1", "", ""}, HTMLType, []string{"section/sect1.html", "sect1/list.html"}},
		{testLayoutIdentifier{"taxonomy", "tag", "", ""}, HTMLType, []string{"taxonomy/tag.html", "indexes/tag.html"}},
		{testLayoutIdentifier{"taxonomyTerm", "categories", "", ""}, HTMLType, []string{"taxonomy/categories.terms.html", "_default/terms.html"}},
		{testLayoutIdentifier{"page", "", "", ""}, HTMLType, []string{"_default/single.html", "theme/_default/single.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", ""}, HTMLType, []string{"_default/mylayout.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", "myttype"}, HTMLType, []string{"myttype/mylayout.html", "_default/mylayout.html"}},
		{testLayoutIdentifier{"page", "", "mylayout", "myttype/mysubtype"}, HTMLType, []string{"myttype/mysubtype/mylayout.html", "myttype/mylayout.html", "_default/mylayout.html"}},
	} {
		layouts := l.For(this.li, this.tp)
		require.NotNil(t, layouts)
		require.True(t, len(layouts) >= len(this.expect))
		// Not checking the complete list for now ...
		require.Equal(t, this.expect, layouts[:len(this.expect)])
	}
}
