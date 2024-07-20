// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gobwas/glob"
	"github.com/gobwas/glob/syntax"
)

const filepathSeparator = string(os.PathSeparator)

var (
	isWindows        = runtime.GOOS == "windows"
	defaultGlobCache = &globCache{
		isWindows: isWindows,
		cache:     make(map[string]globErr),
	}
)

type globErr struct {
	glob glob.Glob
	err  error
}

type globCache struct {
	// Config
	isWindows bool

	// Cache
	sync.RWMutex
	cache map[string]globErr
}

func (gc *globCache) GetGlob(pattern string) (glob.Glob, error) {
	var eg globErr

	gc.RLock()
	var found bool
	eg, found = gc.cache[pattern]
	gc.RUnlock()
	if found {
		return eg.glob, eg.err
	}

	var g glob.Glob
	var err error

	pattern = filepath.ToSlash(pattern)
	g, err = glob.Compile(strings.ToLower(pattern), '/')

	eg = globErr{
		globDecorator{
			g:         g,
			isWindows: gc.isWindows,
		},
		err,
	}

	gc.Lock()
	gc.cache[pattern] = eg
	gc.Unlock()

	return eg.glob, eg.err
}

// Or creates a new Glob from the given globs.
func Or(globs ...glob.Glob) glob.Glob {
	return globSlice{globs: globs}
}

// MatchesFunc is a convenience type to create a glob.Glob from a function.
type MatchesFunc func(s string) bool

func (m MatchesFunc) Match(s string) bool {
	return m(s)
}

type globSlice struct {
	globs []glob.Glob
}

func (g globSlice) Match(s string) bool {
	for _, g := range g.globs {
		if g.Match(s) {
			return true
		}
	}
	return false
}

type globDecorator struct {
	// On Windows we may get filenames with Windows slashes to match,
	// which we need to normalize.
	isWindows bool

	g glob.Glob
}

func (g globDecorator) Match(s string) bool {
	if g.isWindows {
		s = filepath.ToSlash(s)
	}
	s = strings.ToLower(s)
	return g.g.Match(s)
}

func GetGlob(pattern string) (glob.Glob, error) {
	return defaultGlobCache.GetGlob(pattern)
}

func NormalizePath(p string) string {
	return strings.ToLower(NormalizePathNoLower(p))
}

func NormalizePathNoLower(p string) string {
	return strings.Trim(path.Clean(filepath.ToSlash(p)), "/.")
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
