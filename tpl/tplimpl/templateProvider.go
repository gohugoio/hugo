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
	"github.com/gohugoio/hugo/deps"
)

// TemplateProvider manages templates.
type TemplateProvider struct{}

// DefaultTemplateProvider is a globally available TemplateProvider.
var DefaultTemplateProvider *TemplateProvider

// Update updates the Hugo Template System in the provided Deps
// with all the additional features, templates & functions.
func (*TemplateProvider) Update(d *deps.Deps) error {
	tmpl, err := newTemplateExec(d)
	if err != nil {
		return err
	}
	return tmpl.postTransform()
}

// Clone clones.
func (*TemplateProvider) Clone(d *deps.Deps) error {
	t := d.Tmpl().(*templateExec)
	d.SetTmpl(t.Clone(d))
	return nil
}
