// Copyright 2016 The Hugo Authors. All rights reserved.
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
package helpers

import (
	"github.com/kyokomi/emoji"
	"github.com/spf13/hugo/bufferpool"
	"reflect"
	"strings"
	"testing"
)

func TestEmojiCustom(t *testing.T) {
	for i, this := range []struct {
		input  string
		expect []byte
	}{
		{"A :smile: a day", []byte(emoji.Sprint("A :smile: a day"))},
		{"A few :smile:s a day", []byte(emoji.Sprint("A few :smile:s a day"))},
		{"A :smile: and a :beer: makes the day for sure.", []byte(emoji.Sprint("A :smile: and a :beer: makes the day for sure."))},
		{"A :smile: and: a :beer:", []byte(emoji.Sprint("A :smile: and: a :beer:"))},
		{"A :diamond_shape_with_a_dot_inside: and then some.", []byte(emoji.Sprint("A :diamond_shape_with_a_dot_inside: and then some."))},
		{":smile:", []byte(emoji.Sprint(":smile:"))},
		{":smi", []byte(":smi")},
		{"A :smile:", []byte(emoji.Sprint("A :smile:"))},
		{":beer:!", []byte(emoji.Sprint(":beer:!"))},
		{"::smile:", []byte(emoji.Sprint("::smile:"))},
		{":beer::", []byte(emoji.Sprint(":beer::"))},
		{" :beer: :", []byte(emoji.Sprint(" :beer: :"))},
		{":beer: and :smile: and another :beer:!", []byte(emoji.Sprint(":beer: and :smile: and another :beer:!"))},
		{" :beer: : ", []byte(emoji.Sprint(" :beer: : "))},
		{"No smilies for you!", []byte("No smilies for you!")},
		{" The motto: no smiles! ", []byte(" The motto: no smiles! ")},
		{":hugo_is_the_best_static_gen:", []byte(":hugo_is_the_best_static_gen:")},
		{"은행 :smile: 은행", []byte(emoji.Sprint("은행 :smile: 은행"))},
	} {
		result := Emojify([]byte(this.input))

		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] got '%q' but expected %q", i, result, this.expect)
		}

	}
}

// The Emoji benchmarks below are heavily skewed in Hugo's direction:
//
// Hugo have a byte slice, wants a byte slice and doesn't mind if the original is modified.

func BenchmarkEmojiKyokomiFprint(b *testing.B) {

	f := func(in []byte) []byte {
		buff := bufferpool.GetBuffer()
		defer bufferpool.PutBuffer(buff)
		emoji.Fprint(buff, string(in))

		bc := make([]byte, buff.Len(), buff.Len())
		copy(bc, buff.Bytes())
		return bc
	}

	doBenchmarkEmoji(b, f)
}

func BenchmarkEmojiKyokomiSprint(b *testing.B) {

	f := func(in []byte) []byte {
		return []byte(emoji.Sprint(string(in)))
	}

	doBenchmarkEmoji(b, f)
}

func BenchmarkHugoEmoji(b *testing.B) {
	doBenchmarkEmoji(b, Emojify)
}

func doBenchmarkEmoji(b *testing.B, f func(in []byte) []byte) {

	type input struct {
		in     []byte
		expect []byte
	}

	data := []struct {
		input  string
		expect string
	}{
		{"A :smile: a day", emoji.Sprint("A :smile: a day")},
		{"A :smile: and a :beer: day keeps the doctor away", emoji.Sprint("A :smile: and a :beer: day keeps the doctor away")},
		{"A :smile: a day and 10 " + strings.Repeat(":beer: ", 10), emoji.Sprint("A :smile: a day and 10 " + strings.Repeat(":beer: ", 10))},
		{"No smiles today.", "No smiles today."},
		{"No smiles for you or " + strings.Repeat("you ", 1000), "No smiles for you or " + strings.Repeat("you ", 1000)},
	}

	var in []input = make([]input, b.N*len(data))
	var cnt = 0
	for i := 0; i < b.N; i++ {
		for _, this := range data {
			in[cnt] = input{[]byte(this.input), []byte(this.expect)}
			cnt++
		}
	}

	b.ResetTimer()
	cnt = 0
	for i := 0; i < b.N; i++ {
		for j := range data {
			currIn := in[cnt]
			cnt++
			result := f(currIn.in)
			if len(result) != len(currIn.expect) {
				b.Fatalf("[%d] emoji std, got \n%q but expected \n%q", j, result, currIn.expect)
			}
		}

	}
}
