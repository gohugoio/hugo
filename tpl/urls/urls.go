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

// Package urls provides template functions to deal with URLs.
package urls

import (
	"errors"
	"fmt"

	"html/template"

	"net/url"

	"github.com/gohugoio/hugo/common/urls"
	"github.com/gohugoio/hugo/deps"
	_errors "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// New returns a new instance of the urls-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps:      deps,
		multihost: deps.Cfg.GetBool("multihost"),
	}
}

// Namespace provides template functions for the "urls" namespace.
type Namespace struct {
	deps      *deps.Deps
	multihost bool
}

// AbsURL takes a given string and converts it to an absolute URL.
func (ns *Namespace) AbsURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}

	return template.HTML(ns.deps.PathSpec.AbsURL(s, false)), nil
}

// Parse parses rawurl into a URL structure. The rawurl may be relative or
// absolute.
func (ns *Namespace) Parse(rawurl interface{}) (*url.URL, error) {
	s, err := cast.ToStringE(rawurl)
	if err != nil {
		return nil, _errors.Wrap(err, "Error in Parse")
	}

	return url.Parse(s)
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

// URLize returns the given argument formatted as URL.
func (ns *Namespace) URLize(a interface{}) (string, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}
	return ns.deps.PathSpec.URLize(s), nil
}

// Anchorize creates sanitized anchor names that are compatible with Blackfriday.
func (ns *Namespace) Anchorize(a interface{}) (string, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", nil
	}
	return ns.deps.ContentSpec.SanitizeAnchorName(s), nil
}

// Ref returns the absolute URL path to a given content item.
func (ns *Namespace) Ref(in interface{}, args interface{}) (template.HTML, error) {
	p, ok := in.(urls.RefLinker)
	if !ok {
		return "", errors.New("invalid Page received in Ref")
	}
	argsm, err := ns.refArgsToMap(args)
	if err != nil {
		return "", err
	}
	s, err := p.Ref(argsm)
	return template.HTML(s), err
}

// RelRef returns the relative URL path to a given content item.
func (ns *Namespace) RelRef(in interface{}, args interface{}) (template.HTML, error) {
	p, ok := in.(urls.RefLinker)
	if !ok {
		return "", errors.New("invalid Page received in RelRef")
	}
	argsm, err := ns.refArgsToMap(args)
	if err != nil {
		return "", err
	}

	s, err := p.RelRef(argsm)
	return template.HTML(s), err
}

func (ns *Namespace) refArgsToMap(args interface{}) (map[string]interface{}, error) {
	var (
		s  string
		of string
	)

	v := args
	if _, ok := v.([]interface{}); ok {
		v = cast.ToStringSlice(v)
	}

	switch v := v.(type) {
	case map[string]interface{}:
		return v, nil
	case map[string]string:
		m := make(map[string]interface{})
		for k, v := range v {
			m[k] = v
		}
		return m, nil
	case []string:
		if len(v) == 0 || len(v) > 2 {
			return nil, fmt.Errorf("invalid numer of arguments to ref")
		}
		// These where the options before we introduced the map type:
		s = v[0]
		if len(v) == 2 {
			of = v[1]
		}
	default:
		var err error
		s, err = cast.ToStringE(args)
		if err != nil {
			return nil, err
		}

	}

	return map[string]interface{}{
		"path":         s,
		"outputFormat": of,
	}, nil
}

// RelLangURL takes a given string and prepends the relative path according to a
// page's position in the project directory structure and the current language.
func (ns *Namespace) RelLangURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	return template.HTML(ns.deps.PathSpec.RelURL(s, !ns.multihost)), nil
}

// AbsLangURL takes a given string and converts it to an absolute URL according
// to a page's position in the project directory structure and the current
// language.
func (ns *Namespace) AbsLangURL(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	return template.HTML(ns.deps.PathSpec.AbsURL(s, !ns.multihost)), nil
}
