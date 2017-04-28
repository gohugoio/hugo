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

package widget

import (
	"fmt"
	"html/template"

	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/deps"
)

var TestTemplateProvider deps.ResourceProvider

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "templates" namespace.
type Namespace struct {
	deps *deps.Deps
}

// IncludeWidgetArea retrieves and display a widget area using the /widgets/ shortcode
func (ns *Namespace) Widgets(name string, context interface{}) (interface{}, error) {
	// Add (_wa: name) index/value to context to access it inside
	// the embedded template (as Widget Area)
	outcontext := make(map[string]interface{})
	outcontext["c"] = context
	outcontext["_wa"] = name

	// See in template_embedded for widgets.html
	templ := ns.deps.Tmpl.Lookup("_internal/widgets.html")
	if templ != nil {
		b := bp.GetBuffer()
		defer bp.PutBuffer(b)

		if err := templ.Execute(b, outcontext); err != nil {
			return "", err
		}

		return template.HTML(b.String()), nil
	}
	return "", fmt.Errorf("Widget area %q not found", name)
}
