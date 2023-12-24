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
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type FilenameFilter struct {
	shouldInclude func(filename string) bool
	inclusions    []glob.Glob
	dirInclusions []glob.Glob
	exclusions    []glob.Glob
	isWindows     bool

	nested []*FilenameFilter
}

func normalizeFilenameGlobPattern(s string) string {
	// Use Unix separators even on Windows.
	s = filepath.ToSlash(s)
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return s
}

// NewFilenameFilter creates a new Glob where the Match method will
// return true if the file should be included.
// Note that the inclusions will be checked first.
func NewFilenameFilter(inclusions, exclusions []string) (*FilenameFilter, error) {
	if inclusions == nil && exclusions == nil {
		return nil, nil
	}
	filter := &FilenameFilter{isWindows: isWindows}

	for _, include := range inclusions {
		include = normalizeFilenameGlobPattern(include)
		g, err := GetGlob(include)
		if err != nil {
			return nil, err
		}
		filter.inclusions = append(filter.inclusions, g)

		// For mounts that do directory walking (e.g. content) we
		// must make sure that all directories up to this inclusion also
		// gets included.
		dir := path.Dir(include)
		parts := strings.Split(dir, "/")
		for i := range parts {
			pattern := "/" + filepath.Join(parts[:i+1]...)
			g, err := GetGlob(pattern)
			if err != nil {
				return nil, err
			}
			filter.dirInclusions = append(filter.dirInclusions, g)
		}
	}

	for _, exclude := range exclusions {
		exclude = normalizeFilenameGlobPattern(exclude)
		g, err := GetGlob(exclude)
		if err != nil {
			return nil, err
		}
		filter.exclusions = append(filter.exclusions, g)
	}

	return filter, nil
}

// MustNewFilenameFilter invokes NewFilenameFilter and panics on error.
func MustNewFilenameFilter(inclusions, exclusions []string) *FilenameFilter {
	filter, err := NewFilenameFilter(inclusions, exclusions)
	if err != nil {
		panic(err)
	}
	return filter
}

// NewFilenameFilterForInclusionFunc create a new filter using the provided inclusion func.
func NewFilenameFilterForInclusionFunc(shouldInclude func(filename string) bool) *FilenameFilter {
	return &FilenameFilter{shouldInclude: shouldInclude, isWindows: isWindows}
}

// Match returns whether filename should be included.
func (f *FilenameFilter) Match(filename string, isDir bool) bool {
	if f == nil {
		return true
	}
	if !f.doMatch(filename, isDir) {
		return false
	}

	for _, nested := range f.nested {
		if !nested.Match(filename, isDir) {
			return false
		}
	}

	return true
}

// Append appends a filter to the chain. The receiver will be copied if needed.
func (f *FilenameFilter) Append(other *FilenameFilter) *FilenameFilter {
	if f == nil {
		return other
	}

	clone := *f
	nested := make([]*FilenameFilter, len(clone.nested)+1)
	copy(nested, clone.nested)
	nested[len(nested)-1] = other
	clone.nested = nested

	return &clone
}

func (f *FilenameFilter) doMatch(filename string, isDir bool) bool {
	if f == nil {
		return true
	}

	if !strings.HasPrefix(filename, filepathSeparator) {
		filename = filepathSeparator + filename
	}

	if f.shouldInclude != nil {
		if f.shouldInclude(filename) {
			return true
		}
		if f.isWindows {
			// The Glob matchers below handles this by themselves,
			// for the shouldInclude we need to take some extra steps
			// to make this robust.
			winFilename := filepath.FromSlash(filename)
			if filename != winFilename {
				if f.shouldInclude(winFilename) {
					return true
				}
			}
		}

	}

	for _, inclusion := range f.inclusions {
		if inclusion.Match(filename) {
			return true
		}
	}

	if isDir && f.inclusions != nil {
		for _, inclusion := range f.dirInclusions {
			if inclusion.Match(filename) {
				return true
			}
		}
	}

	for _, exclusion := range f.exclusions {
		if exclusion.Match(filename) {
			return false
		}
	}

	return f.inclusions == nil && f.shouldInclude == nil
}
