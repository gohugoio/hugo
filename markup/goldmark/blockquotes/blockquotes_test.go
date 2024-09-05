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

package blockquotes

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestResolveBlockQuoteAlert(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	tests := []struct {
		input    string
		expected blockQuoteAlert
	}{
		{
			input:    "[!NOTE]",
			expected: blockQuoteAlert{typ: "note"},
		},
		{
			input:    "[!FaQ]",
			expected: blockQuoteAlert{typ: "faq"},
		},
		{
			input:    "[!NOTE]+",
			expected: blockQuoteAlert{typ: "note", sign: "+"},
		},
		{
			input:    "[!NOTE]-",
			expected: blockQuoteAlert{typ: "note", sign: "-"},
		},
		{
			input:    "[!NOTE] This is a note",
			expected: blockQuoteAlert{typ: "note", title: "This is a note"},
		},
		{
			input:    "[!NOTE]+ This is a note",
			expected: blockQuoteAlert{typ: "note", sign: "+", title: "This is a note"},
		},
		{
			input:    "[!NOTE]+ This is a title\nThis is not.",
			expected: blockQuoteAlert{typ: "note", sign: "+", title: "This is a title"},
		},
		{
			input:    "[!NOTE]\nThis is not.",
			expected: blockQuoteAlert{typ: "note"},
		},
	}

	for i, test := range tests {
		c.Assert(resolveBlockQuoteAlert("<p>"+test.input+"</p>"), qt.Equals, test.expected, qt.Commentf("Test %d", i))
	}
}
