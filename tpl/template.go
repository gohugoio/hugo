// Copyright 2019 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"reflect"

	"io"
	"regexp"

	"github.com/gohugoio/hugo/output"

	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
)

// TemplateManager manages the collection of templates.
type TemplateManager interface {
	TemplateHandler
	TemplateFuncGetter
	AddTemplate(name, tpl string) error
	MarkReady() error
}

// TemplateVariants describes the possible variants of a template.
// All of these may be empty.
type TemplateVariants struct {
	Language     string
	OutputFormat output.Format
}

// TemplateFinder finds templates.
type TemplateFinder interface {
	TemplateLookup
	TemplateLookupVariant
}

// TemplateHandler finds and executes templates.
type TemplateHandler interface {
	TemplateFinder
	Execute(t Template, wr io.Writer, data interface{}) error
	LookupLayout(d output.LayoutDescriptor, f output.Format) (Template, bool, error)
	HasTemplate(name string) bool
}

type TemplateLookup interface {
	Lookup(name string) (Template, bool)
}

type TemplateLookupVariant interface {
	// TODO(bep) this currently only works for shortcodes.
	// We may unify and expand this variant pattern to the
	// other templates, but we need this now for the shortcodes to
	// quickly determine if a shortcode has a template for a given
	// output format.
	// It returns the template, if it was found or not and if there are
	// alternative representations (output format, language).
	// We are currently only interested in output formats, so we should improve
	// this for speed.
	LookupVariant(name string, variants TemplateVariants) (Template, bool, bool)
}

// Template is the common interface between text/template and html/template.
type Template interface {
	Name() string
	Prepare() (*texttemplate.Template, error)
}

// TemplateParser is used to parse ad-hoc templates, e.g. in the Resource chain.
type TemplateParser interface {
	Parse(name, tpl string) (Template, error)
}

// TemplateParseFinder provides both parsing and finding.
type TemplateParseFinder interface {
	TemplateParser
	TemplateFinder
}

// TemplateDebugger prints some debug info to stdoud.
type TemplateDebugger interface {
	Debug()
}

// templateInfo wraps a Template with some additional information.
type templateInfo struct {
	Template
	Info
}

// templateInfo wraps a Template with some additional information.
type templateInfoManager struct {
	Template
	InfoManager
}

// TemplatesProvider as implemented by deps.Deps.
type TemplatesProvider interface {
	Tmpl() TemplateHandler
	TextTmpl() TemplateParseFinder
}

// WithInfo wraps the info in a template.
func WithInfo(templ Template, info Info) Template {
	if manager, ok := info.(InfoManager); ok {
		return &templateInfoManager{
			Template:    templ,
			InfoManager: manager,
		}
	}

	return &templateInfo{
		Template: templ,
		Info:     info,
	}
}

var baseOfRe = regexp.MustCompile("template: (.*?):")

func extractBaseOf(err string) string {
	m := baseOfRe.FindStringSubmatch(err)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

// TemplateFuncGetter allows to find a template func by name.
type TemplateFuncGetter interface {
	GetFunc(name string) (reflect.Value, bool)
}
