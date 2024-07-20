// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/source"
)

type fileInfo struct {
	*source.File

	overriddenLang string
}

func (fi *fileInfo) Open() (afero.File, error) {
	f, err := fi.FileInfo().Meta().Open()
	if err != nil {
		err = fmt.Errorf("fileInfo: %w", err)
	}

	return f, err
}

func (fi *fileInfo) Lang() string {
	if fi.overriddenLang != "" {
		return fi.overriddenLang
	}
	return fi.File.Lang()
}

func (fi *fileInfo) String() string {
	if fi == nil || fi.File == nil {
		return ""
	}
	return fi.Path()
}
