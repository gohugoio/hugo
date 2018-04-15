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

package path

import (
	"fmt"
	_path "path"

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

// DirFile holds the result from path.Split.
type DirFile struct {
	Dir  string
	File string
}

// Used in test.
func (df DirFile) String() string {
	return fmt.Sprintf("%s|%s", df.Dir, df.File)
}

// Split splits path immediately following the final slash,
// separating it into a directory and file name component.
// If there is no slash in path, Split returns an empty dir and
// file set to path.
// The returned values have the property that path = dir+file.
func (ns *Namespace) Split(path interface{}) (DirFile, error) {
	spath, err := cast.ToStringE(path)
	if err != nil {
		return DirFile{}, err
	}
	dir, file := _path.Split(spath)

	return DirFile{Dir: dir, File: file}, nil
}
