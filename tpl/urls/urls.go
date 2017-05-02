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

package urls

import (
	"errors"
	"html/template"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/deps"
)

// New returns a new instance of the urls-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "urls" namespace.
type Namespace struct {
	deps *deps.Deps
}

// AbsURL takes a given string and converts it to an absolute URL.
func (ns *Namespace) AbsURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}

	return template.HTML(ns.deps.PathSpec.AbsURL(s, false)), nil
}

// RelURL takes a given string and prepends the relative path according to a
// page's position in the project directory structure.
func (ns *Namespace) RelURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}

	return template.HTML(ns.deps.PathSpec.RelURL(s, false)), nil
}

func (ns *Namespace) URLize(a interface{}) (string, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}
	return ns.deps.PathSpec.URLize(s), nil
}

type reflinker interface {
	Ref(refs ...string) (string, error)
	RelRef(refs ...string) (string, error)
}

// Ref returns the absolute URL path to a given content item.
func (ns *Namespace) Ref(in interface{}, refs ...string) (template.HTML, error) {
	p, ok := in.(reflinker)
	if !ok {
		return "", errors.New("invalid Page received in Ref")
	}
	s, err := p.Ref(refs...)
	return template.HTML(s), err
}

// RelRef returns the relative URL path to a given content item.
func (ns *Namespace) RelRef(in interface{}, refs ...string) (template.HTML, error) {
	p, ok := in.(reflinker)
	if !ok {
		return "", errors.New("invalid Page received in RelRef")
	}
	s, err := p.RelRef(refs...)
	return template.HTML(s), err
}

// RelLangURL takes a given string and prepends the relative path according to a
// page's position in the project directory structure and the current language.
func (ns *Namespace) RelLangURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	return template.HTML(ns.deps.PathSpec.RelURL(s, true)), nil
}

// AbsLangURL takes a given string and converts it to an absolute URL according
// to a page's position in the project directory structure and the current
// language.
func (ns *Namespace) AbsLangURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	return template.HTML(ns.deps.PathSpec.AbsURL(s, true)), nil
}
