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
	"regexp"
	"sync"

	"github.com/spf13/cast"
)

// FindRE returns a list of strings that match the regular expression. By default all matches
// will be included. The number of matches can be limited with an optional third parameter.
func (ns *Namespace) FindRE(expr string, content interface{}, limit ...interface{}) ([]string, error) {
	re, err := reCache.Get(expr)
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

// ReplaceRE returns a copy of s, replacing all matches of the regular
// expression pattern with the replacement text repl.
func (ns *Namespace) ReplaceRE(pattern, repl, s interface{}) (_ string, err error) {
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

	re, err := reCache.Get(sp)
	if err != nil {
		return "", err
	}

	return re.ReplaceAllString(ss, sr), nil
}

// regexpCache represents a cache of regexp objects protected by a mutex.
type regexpCache struct {
	mu sync.RWMutex
	re map[string]*regexp.Regexp
}

// Get retrieves a regexp object from the cache based upon the pattern.
// If the pattern is not found in the cache, create one
func (rc *regexpCache) Get(pattern string) (re *regexp.Regexp, err error) {
	var ok bool

	if re, ok = rc.get(pattern); !ok {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		rc.set(pattern, re)
	}

	return re, nil
}

func (rc *regexpCache) get(key string) (re *regexp.Regexp, ok bool) {
	rc.mu.RLock()
	re, ok = rc.re[key]
	rc.mu.RUnlock()
	return
}

func (rc *regexpCache) set(key string, re *regexp.Regexp) {
	rc.mu.Lock()
	rc.re[key] = re
	rc.mu.Unlock()
}

var reCache = regexpCache{re: make(map[string]*regexp.Regexp)}
