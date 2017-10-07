// Copyright 2017 The Hugo Authors. All rights reserved.
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

package metrics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimilarPercentage(t *testing.T) {
	assert := require.New(t)

	sentence := "this is some words about nothing, Hugo!"
	words := strings.Fields(sentence)
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	sentenceReversed := strings.Join(words, " ")

	assert.Equal(100, howSimilar("Hugo Rules", "Hugo Rules"))
	assert.Equal(50, howSimilar("Hugo Rules", "Hugo Rocks"))
	assert.Equal(66, howSimilar("The Hugo Rules", "The Hugo Rocks"))
	assert.Equal(66, howSimilar("The Hugo Rules", "The Hugo"))
	assert.Equal(66, howSimilar("The Hugo", "The Hugo Rules"))
	assert.Equal(0, howSimilar("Totally different", "Not Same"))
	assert.Equal(14, howSimilar(sentence, sentenceReversed))

}

func BenchmarkHowSimilar(b *testing.B) {
	s1 := "Hugo is cool and " + strings.Repeat("fun ", 10) + "!"
	s2 := "Hugo is cool and " + strings.Repeat("cool ", 10) + "!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		howSimilar(s1, s2)
	}
}
