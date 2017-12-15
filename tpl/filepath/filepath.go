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

package filepath

import (
	"path/filepath"

	"github.com/spf13/cast"
)

// New returns a new instance of the filepath-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "filepath" namespace.
type Namespace struct{}

// Base returns the last element of path. Trailing path separators are removed
// before extracting the last element. If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func (ns *Namespace) Base(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.Base(s), nil
}

// Clean returns the shortest path name equivalent to path
// by purely lexical processing. It applies the following rules
// iteratively until no further processing can be done:
//
//	1. Replace multiple Separator elements with a single one.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path,
//	   assuming Separator is '/'.
//
// The returned path ends in a slash only if it represents a root directory,
// such as "/" on Unix or `C:\` on Windows.
//
// Finally, any occurrences of slash are replaced by Separator.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
func (ns *Namespace) Clean(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(s), nil
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed. If the path is empty, Dir returns ".". If the path
// consists entirely of separators, Dir returns a single separator. The returned
// path does not end in a separator unless it is the root directory.
func (ns *Namespace) Dir(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.Dir(s), nil
}

// Ext returns the file name extension used by path. The extension is the suffix
// beginning at the final dot in the final element of path; it is empty if there
// is no dot.
func (ns *Namespace) Ext(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.Ext(s), nil
}

// FromSlash returns the result of replacing each slash ('/') character in path
// with a separator character. Multiple slashes are replaced by multiple
// separators.
func (ns *Namespace) FromSlash(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.FromSlash(s), nil
}

// Separator returns the OS-specific path separator.
func (ns *Namespace) Separator() string {
	return string(filepath.Separator)
}

// Split splits path immediately following the final Separator, separating it
// into a directory and file name component. If there is no Separator in path,
// Split returns an empty dir and file set to path. The returned values have the
// property that path = dir+file.
func (ns *Namespace) Split(path interface{}) (string, string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", "", err
	}

	dir, file := filepath.Split(s)
	return dir, file, nil
}

// ToSlash returns the result of replacing each separator character in path with
// a slash ('/') character. Multiple separators are replaced by multiple slashes.
func (ns *Namespace) ToSlash(path interface{}) (string, error) {
	s, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(s), nil
}
