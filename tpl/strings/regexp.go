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

package strings

import (
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/spf13/cast"
)

// FindRE returns a list of strings that match the regular expression. By default all matches
// will be included. The number of matches can be limited with an optional third parameter.
func (ns *Namespace) FindRE(expr string, content any, limit ...any) ([]string, error) {
	re, err := hstrings.GetOrCompileRegexp(expr)
	if err != nil {
		return nil, err
	}

	conv, err := cast.ToStringE(content)
	if err != nil {
		return nil, err
	}

	if len(limit) == 0 {
		return re.FindAllString(conv, -1), nil
	}

	lim, err := cast.ToIntE(limit[0])
	if err != nil {
		return nil, err
	}

	return re.FindAllString(conv, lim), nil
}

// FindRESubmatch returns a slice of all successive matches of the regular
// expression in content. Each element is a slice of strings holding the text
// of the leftmost match of the regular expression and the matches, if any, of
// its subexpressions.
//
// By default all matches will be included. The number of matches can be
// limited with the optional limit parameter. A return value of nil indicates
// no match.
func (ns *Namespace) FindRESubmatch(expr string, content any, limit ...any) ([][]string, error) {
	re, err := hstrings.GetOrCompileRegexp(expr)
	if err != nil {
		return nil, err
	}

	conv, err := cast.ToStringE(content)
	if err != nil {
		return nil, err
	}
	n := -1
	if len(limit) > 0 {
		n, err = cast.ToIntE(limit[0])
		if err != nil {
			return nil, err
		}
	}

	return re.FindAllStringSubmatch(conv, n), nil
}

// ReplaceRE returns a copy of s, replacing all matches of the regular
// expression pattern with the replacement text repl. The number of replacements
// can be limited with an optional fourth parameter.
func (ns *Namespace) ReplaceRE(pattern, repl, s any, n ...any) (_ string, err error) {
	sp, err := cast.ToStringE(pattern)
	if err != nil {
		return
	}

	sr, err := cast.ToStringE(repl)
	if err != nil {
		return
	}

	ss, err := cast.ToStringE(s)
	if err != nil {
		return
	}

	nn := -1
	if len(n) > 0 {
		nn, err = cast.ToIntE(n[0])
		if err != nil {
			return
		}
	}

	re, err := hstrings.GetOrCompileRegexp(sp)
	if err != nil {
		return "", err
	}

	return re.ReplaceAllStringFunc(ss, func(str string) string {
		if nn == 0 {
			return str
		}

		nn -= 1
		return re.ReplaceAllString(str, sr)
	}), nil
}
