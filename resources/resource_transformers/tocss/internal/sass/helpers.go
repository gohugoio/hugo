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

package sass

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/types/css"
)

const (
	HugoVarsNamespace = "hugo:vars"
)

func CreateVarsStyleSheet(vars map[string]any) string {
	if vars == nil {
		return ""
	}
	var varsStylesheet string

	var varsSlice []string
	for k, v := range vars {
		var prefix string
		if !strings.HasPrefix(k, "$") {
			prefix = "$"
		}

		switch v.(type) {
		case css.QuotedString:
			// Marked by the user as a string that needs to be quoted.
			varsSlice = append(varsSlice, fmt.Sprintf("%s%s: %q;", prefix, k, v))
		default:
			if isTypedCSSValue(v) {
				// E.g. 24px, 1.5rem, 10%, hsl(0, 0%, 100%), calc(24px + 36px), #fff, #ffffff.
				varsSlice = append(varsSlice, fmt.Sprintf("%s%s: %v;", prefix, k, v))
			} else {
				// unquote will preserve quotes around URLs etc. if needed.
				varsSlice = append(varsSlice, fmt.Sprintf("%s%s: unquote(%q);", prefix, k, v))
			}
		}
	}
	sort.Strings(varsSlice)
	varsStylesheet = strings.Join(varsSlice, "\n")
	return varsStylesheet
}

var (
	isCSSColor = regexp.MustCompile(`^#[0-9a-fA-F]{3,6}$`)
	isCSSFunc  = regexp.MustCompile(`^([a-zA-Z-]+)\(`)
	isCSSUnit  = regexp.MustCompile(`^([0-9]+)(\.[0-9]+)?([a-zA-Z-%]+)$`)
)

// isTypedCSSValue returns true if the given string is a CSS value that
// we should preserve the type of, as in: Not wrap it in quotes.
func isTypedCSSValue(v any) bool {
	switch s := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, css.UnquotedString:
		return true
	case string:
		if isCSSColor.MatchString(s) {
			return true
		}
		if isCSSFunc.MatchString(s) {
			return true
		}
		if isCSSUnit.MatchString(s) {
			return true
		}

	}

	return false
}
