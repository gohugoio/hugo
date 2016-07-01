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
	"regexp"
	"sync"
)

var (
	emojiInit sync.Once
	regEmoji  *regexp.Regexp
	emojiMap  = make(map[string]string)
)

// Emojify "emojifies" the input source.
// See http://www.emoji-cheat-sheet.com/
func Emojify(source []byte) []byte {
	emojiInit.Do(initEmoji)

	return []byte(regEmoji.ReplaceAllStringFunc(string(source), EmojiTranslateFunc))
}

func EmojiTranslateFunc(str string) string {
	emoji, ok := emojiMap[str]
	if ok {
		return emoji
	}
	return str
}

func initEmoji() {
	emojiMap = emoji.CodeMap()
	regEmoji = regexp.MustCompile(":[a-zA-Z0-9_-]+:")
}
