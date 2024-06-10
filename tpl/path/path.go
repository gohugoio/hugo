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
// limitations under the License.

// Package path provides template functions for manipulating paths.
package path

import (
	_path "path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

// New returns a new instance of the path-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "os" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
// The input path is passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
func (ns *Namespace) Ext(path any) (string, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}
	spath = filepath.ToSlash(spath)
	return _path.Ext(spath), nil
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element using Split, the path is Cleaned and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of slashes followed by non-slash bytes, Dir
// returns a single slash. In any other case, the returned path does not end in a
// slash.
// The input path is passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
func (ns *Namespace) Dir(path any) (string, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}
	spath = filepath.ToSlash(spath)
	return _path.Dir(spath), nil
}

// Base returns the last element of path.
// Trailing slashes are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of slashes, Base returns "/".
// The input path is passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
func (ns *Namespace) Base(path any) (string, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}
	spath = filepath.ToSlash(spath)
	return _path.Base(spath), nil
}

// BaseName returns the last element of path, removing the extension if present.
// Trailing slashes are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of slashes, Base returns "/".
// The input path is passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
func (ns *Namespace) BaseName(path any) (string, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}
	spath = filepath.ToSlash(spath)
	return strings.TrimSuffix(_path.Base(spath), _path.Ext(spath)), nil
}

// Split splits path immediately following the final slash,
// separating it into a directory and file name component.
// If there is no slash in path, Split returns an empty dir and
// file set to path.
// The input path is passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
// The returned values have the property that path = dir+file.
func (ns *Namespace) Split(path any) (paths.DirFile, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return paths.DirFile{}, err
	}
	spath = filepath.ToSlash(spath)
	dir, file := _path.Split(spath)

	return paths.DirFile{Dir: dir, File: file}, nil
}

// Join joins any number of path elements into a single path, adding a
// separating slash if necessary. All the input
// path elements are passed into filepath.ToSlash converting any Windows slashes
// to forward slashes.
// The result is Cleaned; in particular,
// all empty strings are ignored.
func (ns *Namespace) Join(elements ...any) (string, error) {
	var pathElements []string
	for _, elem := range elements {
		switch v := elem.(type) {
		case []string:
			for _, e := range v {
				pathElements = append(pathElements, filepath.ToSlash(e))
			}
		case []any:
			for _, e := range v {
				elemStr, err := cast.ToStringE(e)
				if err != nil {
					return "", err
				}
				pathElements = append(pathElements, filepath.ToSlash(elemStr))
			}
		default:
			elemStr, err := cast.ToStringE(elem)
			if err != nil {
				return "", err
			}
			pathElements = append(pathElements, filepath.ToSlash(elemStr))
		}
	}
	return _path.Join(pathElements...), nil
}

// Clean replaces the separators used with standard slashes and then
// extraneous slashes are removed.
func (ns *Namespace) Clean(path any) (string, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return "", err
	}
	spath = filepath.ToSlash(spath)
	return _path.Clean(spath), nil
}
