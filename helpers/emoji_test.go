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
	"reflect"
	"testing"
)

func TestEmojiCustom(t *testing.T) {
	for i, this := range []struct {
		input  string
		expect []byte
	}{
		{"A :smile: a day", []byte("A 😄 a day")},
		{"A few :smile:s a day", []byte("A few 😄s a day")},
		{"A :smile: and a :beer: makes the day for sure.", []byte("A 😄 and a 🍺 makes the day for sure.")},
		{"A :smile: and: a :beer:", []byte("A 😄 and: a 🍺")},
		{"A :diamond_shape_with_a_dot_inside: and then some.", []byte("A 💠 and then some.")},
		{":smile:", []byte("😄")},
		{":smi", []byte(":smi")},
		{"A :smile:", []byte("A 😄")},
		{":beer:!", []byte("🍺!")},
		{"::smile:", []byte(":😄")},
		{":beer::", []byte("🍺:")},
		{" :beer: :", []byte(" 🍺 :")},
		{":beer: and :smile: and another :beer:!", []byte("🍺 and 😄 and another 🍺!")},
		{" :beer: : ", []byte(" 🍺 : ")},
		{"No smiles for you!", []byte("No smiles for you!")},
		{" The motto: no smiles! ", []byte(" The motto: no smiles! ")},
		{":hugo_is_the_best_static_gen:", []byte(":hugo_is_the_best_static_gen:")},
		{"은행 :smile: 은행", []byte("은행 😄 은행")},
		// #2198
		{"See: A :beer:!", []byte("See: A 🍺!")},
		{`Aaaaaaaaaa: aaaaaaaaaa aaaaaaaaaa aaaaaaaaaa.

:beer:`, []byte(`Aaaaaaaaaa: aaaaaaaaaa aaaaaaaaaa aaaaaaaaaa.

🍺`)},
		{"test :\n```bash\nthis is a test\n```\n\ntest\n\n:cool::blush:::pizza:\\:blush : : blush: :pizza:", []byte("test :\n```bash\nthis is a test\n```\n\ntest\n\n🆒😊:🍕\\:blush : : blush: 🍕")},
		{
			// 2391
			"[a](http://gohugo.io) :smile: [r](http://gohugo.io/introduction/overview/) :beer:",
			[]byte(`[a](http://gohugo.io) 😄 [r](http://gohugo.io/introduction/overview/) 🍺`),
		},
	} {

		result := Emojify([]byte(this.input))

		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] got %q but expected %q", i, result, this.expect)
		}

	}
}
