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
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/types/css"
)

const (
	HugoVarsNamespace = "hugo:vars"
	// Transpiler implementation can be controlled from the client by
	// setting the 'transpiler' option.
	// Default is currently 'libsass', but that may change.
	TranspilerDart    = "dartsass"
	TranspilerLibSass = "libsass"
)

// HugoVarsSubPath returns the slash-separated sub-path of a "hugo:vars" URL,
// e.g. "hugo:vars" -> "" and "hugo:vars/mobile" -> "mobile". The second return
// is false if url is not in the "hugo:vars" namespace.
func HugoVarsSubPath(url string) (string, bool) {
	if url == HugoVarsNamespace {
		return "", true
	}
	if rest, ok := strings.CutPrefix(url, HugoVarsNamespace+"/"); ok {
		return rest, true
	}
	return "", false
}

// PrepareVars lowercases all keys for any map value recursively and returns a clone if modified.
func PrepareVars(vars map[string]any) map[string]any {
	if vars == nil {
		return nil
	}

	// Lowercase all keys for map values recursively, so that they can be accessed case-insensitively from the stylesheet.
	var isCloned bool
	for k, v := range vars {
		if hstrings.HasUppercase(k) && hreflect.IsMap(v) {
			if !isCloned {
				vars = maps.Clone(vars)
			}
			delete(vars, k)
			vars[strings.ToLower(k)] = PrepareVars(hmaps.ToStringMap(v))
			isCloned = true
		}
	}
	return vars
}

// ResolveVars returns the entries of vars at the given slash-separated path.
// Nested map entries are excluded from the result, so only scalar/typed values remain.
// An empty path returns the top-level scalars.
func ResolveVars(vars map[string]any, path string) map[string]any {
	if vars == nil {
		return nil
	}
	if path == "" || path == "/" {
		return removeMaps(vars)
	}
	vv, err := hmaps.GetNestedParam(path, "/", vars)
	if err != nil {
		return nil
	}

	return removeMaps(hmaps.ToStringMap(vv))
}

func removeMaps(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		if hreflect.IsMap(v) {
			continue
		}
		res[k] = v
	}
	return res
}

func CreateVarsStyleSheet(transpiler string, vars map[string]any) string {
	if vars == nil {
		return ""
	}
	var varsStylesheet string

	var varsSlice []string
	for k, v := range vars {
		if hreflect.IsMap(v) {
			// Nested vars are exposed via "hugo:vars/<name>" namespaces, skip here.
			continue
		}
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
				if transpiler == TranspilerDart {
					varsSlice = append(varsSlice, fmt.Sprintf("%s%s: string.unquote(%q);", prefix, k, v))
				} else {
					varsSlice = append(varsSlice, fmt.Sprintf("%s%s: unquote(%q);", prefix, k, v))
				}
			}
		}
	}
	sort.Strings(varsSlice)

	if transpiler == TranspilerDart {
		varsStylesheet = `@use "sass:string";` + "\n" + strings.Join(varsSlice, "\n")
	} else {
		varsStylesheet = strings.Join(varsSlice, "\n")
	}

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
