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
// limitatio	ns under the License.

package herrors

import (
	"fmt"
	"regexp"
	"strconv"
)

var lineNumberExtractors = []lineNumberExtractor{
	// Template/shortcode parse errors
	newLineNumberErrHandlerFromRegexp("(.*?:)(\\d+)(:.*)"),

	// TOML parse errors
	newLineNumberErrHandlerFromRegexp("(.*Near line )(\\d+)(\\s.*)"),

	// YAML parse errors
	newLineNumberErrHandlerFromRegexp("(line )(\\d+)(:)"),
}

type lineNumberExtractor func(e error, offset int) (int, string)

func newLineNumberErrHandlerFromRegexp(expression string) lineNumberExtractor {
	re := regexp.MustCompile(expression)
	return extractLineNo(re)
}

func extractLineNo(re *regexp.Regexp) lineNumberExtractor {
	return func(e error, offset int) (int, string) {
		if e == nil {
			panic("no error")
		}
		s := e.Error()
		m := re.FindStringSubmatch(s)
		if len(m) == 4 {
			i, _ := strconv.Atoi(m[2])
			msg := e.Error()
			if offset != 0 {
				i = i + offset
				msg = re.ReplaceAllString(s, fmt.Sprintf("${1}%d${3}", i))
			}
			return i, msg
		}

		return -1, ""
	}
}
