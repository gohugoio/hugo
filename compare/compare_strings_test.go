// Copyright 2019 The Hugo Authors. All rights reserved.
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

package compare

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	assert := require.New(t)
	for i, test := range []struct {
		a string
		b string
	}{
		{"a", "a"},
		{"A", "a"},
		{"Ab", "Ac"},
		{"az", "Za"},
		{"C", "D"},
		{"B", "a"},
		{"C", ""},
		{"", ""},
		{"αβδC", "ΑΒΔD"},
		{"αβδC", "ΑΒΔ"},
		{"αβδ", "ΑΒΔD"},
		{"αβδ", "ΑΒΔ"},
		{"β", "δ"},
		{"好", strings.ToLower("好")},
	} {

		expect := strings.Compare(strings.ToLower(test.a), strings.ToLower(test.b))
		got := compareFold(test.a, test.b)

		assert.Equal(expect, got, fmt.Sprintf("test %d: %d", i, expect))

	}
}

func TestLexicographicSort(t *testing.T) {
	assert := require.New(t)

	s := []string{"b", "Bz", "ba", "A", "Ba", "ba"}

	sort.Slice(s, func(i, j int) bool {
		return LessStrings(s[i], s[j])
	})

	assert.Equal([]string{"A", "b", "Ba", "ba", "ba", "Bz"}, s)

}
