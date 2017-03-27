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

package tplimpl

import (
	"path/filepath"

	"strings"

	"github.com/yosssi/ace"
)

func (t *templateHandler) addAceTemplate(name, basePath, innerPath string, baseContent, innerContent []byte) error {
	t.checkState()
	var base, inner *ace.File
	name = name[:len(name)-len(filepath.Ext(innerPath))] + ".html"

	// Fixes issue #1178
	basePath = strings.Replace(basePath, "\\", "/", -1)
	innerPath = strings.Replace(innerPath, "\\", "/", -1)

	if basePath != "" {
		base = ace.NewFile(basePath, baseContent)
		inner = ace.NewFile(innerPath, innerContent)
	} else {
		base = ace.NewFile(innerPath, innerContent)
		inner = ace.NewFile("", []byte{})
	}
	parsed, err := ace.ParseSource(ace.NewSource(base, inner, []*ace.File{}), nil)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
		return err
	}
	templ, err := ace.CompileResultWithTemplate(t.html.t.New(name), parsed, nil)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
		return err
	}
	return applyTemplateTransformersToHMLTTemplate(templ)
}
