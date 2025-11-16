// Copyright 2024 The Hugo Authors. All rights reserved.
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

package pageparser

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_lexIntroSection(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	for i, tt := range []struct {
		input                string
		expectItemType       ItemType
		expectSummaryDivider []byte
	}{
		{"{\"title\": \"JSON\"}\n", TypeFrontMatterJSON, summaryDivider},
		{"#+TITLE: ORG\n", TypeFrontMatterORG, summaryDividerOrg},
		{"+++\ntitle = \"TOML\"\n+++\n", TypeFrontMatterTOML, summaryDivider},
		{"---\ntitle: YAML\n---\n", TypeFrontMatterYAML, summaryDivider},
		// Issue 13152
		{"# ATX Header Level 1\n", tText, summaryDivider},
	} {
		errMsg := qt.Commentf("[%d] %v", i, tt.input)

		l := newPageLexer([]byte(tt.input), lexIntroSection, Config{})
		l.run()

		c.Assert(l.items[0].Type, qt.Equals, tt.expectItemType, errMsg)
		c.Assert(l.summaryDivider, qt.DeepEquals, tt.expectSummaryDivider, errMsg)

	}
}
