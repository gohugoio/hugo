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
	"bytes"
	"sync"

	"github.com/kyokomi/emoji"
)

var (
	emojiInit sync.Once

	emojis = make(map[string][]byte)

	emojiDelim     = []byte(":")
	emojiWordDelim = []byte(" ")
	emojiMaxSize   int
)

// Emojify "emojifies" the input source.
// Note that the input byte slice will be modified if needed.
// See http://www.emoji-cheat-sheet.com/
func Emojify(source []byte) []byte {
	emojiInit.Do(initEmoji)

	start := 0
	k := bytes.Index(source[start:], emojiDelim)

	for k != -1 {

		j := start + k

		upper := j + emojiMaxSize

		if upper > len(source) {
			upper = len(source)
		}

		endEmoji := bytes.Index(source[j+1:upper], emojiDelim)
		nextWordDelim := bytes.Index(source[j:upper], emojiWordDelim)

		if endEmoji < 0 {
			start++
		} else if endEmoji == 0 || (nextWordDelim != -1 && nextWordDelim < endEmoji) {
			start += endEmoji + 1
		} else {
			endKey := endEmoji + j + 2
			emojiKey := source[j:endKey]

			if emoji, ok := emojis[string(emojiKey)]; ok {
				source = append(source[:j], append(emoji, source[endKey:]...)...)
			}

			start += endEmoji
		}

		if start >= len(source) {
			break
		}

		k = bytes.Index(source[start:], emojiDelim)
	}

	return source
}

func initEmoji() {
	emojiMap := emoji.CodeMap()

	for k, v := range emojiMap {
		emojis[k] = []byte(v)

		if len(k) > emojiMaxSize {
			emojiMaxSize = len(k)
		}
	}

}
