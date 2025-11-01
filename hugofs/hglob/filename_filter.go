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

package hglob

import (
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gobwas/glob"
)

type FilenameFilter struct {
	shouldInclude func(filename string) bool

	entries []globFilenameFilterEntry

	isWindows bool

	nested []*FilenameFilter
}

type globFilenameFilterEntryType int

const (
	globFilenameFilterEntryTypeInclusion globFilenameFilterEntryType = iota
	globFilenameFilterEntryTypeInclusionDir
	globFilenameFilterEntryTypeExclusion
)

// NegationPrefix is the prefix that makes a pattern an exclusion.
const NegationPrefix = "! "

type globFilenameFilterEntry struct {
	g glob.Glob
	t globFilenameFilterEntryType
}

func normalizeFilenameGlobPattern(s string) string {
	// Use Unix separators even on Windows.
	s = filepath.ToSlash(s)
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return s
}

func NewFilenameFilterV2(patterns []string) (*FilenameFilter, error) {
	if len(patterns) == 0 {
		return nil, nil
	}
	filter := &FilenameFilter{isWindows: isWindows}
	for _, p := range patterns {
		var t globFilenameFilterEntryType
		if strings.HasPrefix(p, NegationPrefix) {
			t = globFilenameFilterEntryTypeExclusion
			p = strings.TrimPrefix(p, NegationPrefix)
		} else {
			t = globFilenameFilterEntryTypeInclusion
		}
		p = normalizeFilenameGlobPattern(p)
		g, err := GetGlob(p)
		if err != nil {
			return nil, err
		}
		filter.entries = append(filter.entries, globFilenameFilterEntry{t: t, g: g})
		if t == globFilenameFilterEntryTypeInclusion {
			// For mounts that do directory walking (e.g. content) we
			// must make sure that all directories up to this inclusion also
			// gets included.
			dir := path.Dir(p)
			parts := strings.Split(dir, "/")
			for i := range parts {
				pattern := "/" + filepath.Join(parts[:i+1]...)
				g, err := GetGlob(pattern)
				if err != nil {
					return nil, err
				}
				filter.entries = append(filter.entries, globFilenameFilterEntry{t: globFilenameFilterEntryTypeInclusionDir, g: g})
			}
		}

	}

	return filter, nil
}

// NewFilenameFilter creates a new Glob where the Match method will
// return true if the file should be included.
// Note that the exclusions will be checked first.
// Deprecated: Use NewFilenameFilterV2.
func NewFilenameFilter(inclusions, exclusions []string) (*FilenameFilter, error) {
	for i, p := range exclusions {
		if !strings.HasPrefix(p, NegationPrefix) {
			exclusions[i] = NegationPrefix + p
		}
	}
	all := slices.Concat(inclusions, exclusions)
	return NewFilenameFilterV2(all)
}

// MustNewFilenameFilter invokes NewFilenameFilter and panics on error.
// Deprecated: Use NewFilenameFilterV2.
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

	var hasInclude bool
	for _, entry := range f.entries {
		switch entry.t {
		case globFilenameFilterEntryTypeExclusion:
			if entry.g.Match(filename) {
				return false
			}
		case globFilenameFilterEntryTypeInclusion:
			if entry.g.Match(filename) {
				return true
			}
			hasInclude = true
		case globFilenameFilterEntryTypeInclusionDir:
			if isDir && entry.g.Match(filename) {
				return true
			}
		}
	}

	return !hasInclude && f.shouldInclude == nil
}
