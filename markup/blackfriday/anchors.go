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

// Package blackfriday holds some compability functions for the old Blackfriday v1 Markdown engine.
package blackfriday

import "unicode"

// SanitizedAnchorName is how Blackfriday sanitizes anchor names.
// Implementation borrowed from https://github.com/russross/blackfriday/blob/a477dd1646916742841ed20379f941cfa6c5bb6f/block.go#L1464
// Note that Hugo removed its Blackfriday support in v0.100.0, but you can still use this strategy for
// auto ID generation.
func SanitizedAnchorName(text string) string {
	var anchorName []rune
	futureDash := false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}
