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
	"github.com/gohugoio/hugo/tpl"
)

// TemplateProvider manages templates.
type TemplateProvider struct{}

// DefaultTemplateProvider is a globally available TemplateProvider.
var DefaultTemplateProvider *TemplateProvider

// Update updates the Hugo Template System in the provided Deps
// with all the additional features, templates & functions.
func (*TemplateProvider) NewResource(dst *deps.Deps) error {
	handlers, err := newTemplateHandlers(dst)
	if err != nil {
		return err
	}
	dst.SetTempl(handlers)
	return nil
}

// Clone clones.
func (*TemplateProvider) CloneResource(dst, src *deps.Deps) error {
	t := src.Tmpl().(*templateExec)
	c := t.Clone(dst)
	funcMap := make(map[string]any)
	for k, v := range c.funcs {
		funcMap[k] = v.Interface()
	}
	dst.SetTempl(&tpl.TemplateHandlers{
		Tmpl:    c,
		TxtTmpl: newStandaloneTextTemplate(funcMap),
	})
	return nil
}
