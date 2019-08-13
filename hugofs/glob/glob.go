// Copyright 2019 The Hugo Authors. All rights reserved.
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

package glob

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobwas/glob"
	"github.com/gobwas/glob/syntax"
)

var (
	globCache = make(map[string]glob.Glob)
	globMu    sync.RWMutex
)

func GetGlob(pattern string) (glob.Glob, error) {
	var g glob.Glob

	globMu.RLock()
	g, found := globCache[pattern]
	globMu.RUnlock()
	if !found {
		var err error
		g, err = glob.Compile(strings.ToLower(pattern), '/')
		if err != nil {
			return nil, err
		}

		globMu.Lock()
		globCache[pattern] = g
		globMu.Unlock()
	}

	return g, nil

}

func NormalizePath(p string) string {
	return strings.Trim(path.Clean(filepath.ToSlash(strings.ToLower(p))), "/.")
}

// ResolveRootDir takes a normalized path on the form "assets/**.json" and
// determines any root dir, i.e. any start path without any wildcards.
func ResolveRootDir(p string) string {
	parts := strings.Split(path.Dir(p), "/")
	var roots []string
	for _, part := range parts {
		if HasGlobChar(part) {
			break
		}
		roots = append(roots, part)
	}

	if len(roots) == 0 {
		return ""
	}

	return strings.Join(roots, "/")
}

// FilterGlobParts removes any string with glob wildcard.
func FilterGlobParts(a []string) []string {
	b := a[:0]
	for _, x := range a {
		if !HasGlobChar(x) {
			b = append(b, x)
		}
	}
	return b
}

// HasGlobChar returns whether s contains any glob wildcards.
func HasGlobChar(s string) bool {
	for i := 0; i < len(s); i++ {
		if syntax.Special(s[i]) {
			return true
		}
	}
	return false

}
