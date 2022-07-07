// Copyright 2018 The Hugo Authors. All rights reserved.
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

func TestMain(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	mainTests := []lexerTest{
		{"emoji #1", "Some text with :emoji:", []typeText{nti(tText, "Some text with "), nti(TypeEmoji, ":emoji:"), tstEOF}},
		{"emoji #2", "Some text with :emoji: and some text.", []typeText{nti(tText, "Some text with "), nti(TypeEmoji, ":emoji:"), nti(tText, " and some text."), tstEOF}},
		{"looks like an emoji #1", "Some text and then :emoji", []typeText{nti(tText, "Some text and then "), nti(tText, ":"), nti(tText, "emoji"), tstEOF}},
		{"looks like an emoji #2", "Some text and then ::", []typeText{nti(tText, "Some text and then "), nti(tText, ":"), nti(tText, ":"), tstEOF}},
		{"looks like an emoji #3", ":Some :text", []typeText{nti(tText, ":"), nti(tText, "Some "), nti(tText, ":"), nti(tText, "text"), tstEOF}},
	}

	for i, test := range mainTests {
		items := collectWithConfig([]byte(test.input), false, lexMainSection, Config{EnableEmoji: true})
		if !equal(test.input, items, test.items) {
			got := itemsToString(items, []byte(test.input))
			expected := testItemsToString(test.items)
			c.Assert(got, qt.Equals, expected, qt.Commentf("Test %d: %s", i, test.name))
		}
	}
}
