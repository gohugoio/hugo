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

package herrors

import (
	"regexp"
	"strconv"
)

var lineNumberExtractors = []lineNumberExtractor{
	// Template/shortcode parse errors
	newLineNumberErrHandlerFromRegexp(`:(\d+):(\d*):`),
	newLineNumberErrHandlerFromRegexp(`:(\d+):`),

	// YAML parse errors
	newLineNumberErrHandlerFromRegexp(`line (\d+):`),

	// i18n bundle errors
	newLineNumberErrHandlerFromRegexp(`\((\d+),\s(\d*)`),
}

type lineNumberExtractor func(e error) (int, int)

func newLineNumberErrHandlerFromRegexp(expression string) lineNumberExtractor {
	re := regexp.MustCompile(expression)
	return extractLineNo(re)
}

func extractLineNo(re *regexp.Regexp) lineNumberExtractor {
	return func(e error) (int, int) {
		if e == nil {
			panic("no error")
		}
		col := 1
		s := e.Error()
		m := re.FindStringSubmatch(s)
		if len(m) >= 2 {
			lno, _ := strconv.Atoi(m[1])
			if len(m) > 2 {
				col, _ = strconv.Atoi(m[2])
			}

			if col <= 0 {
				col = 1
			}

			return lno, col
		}

		return 0, col
	}
}
