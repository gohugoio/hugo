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
	"fmt"
	"html/template"
	"strings"

	bp "github.com/spf13/hugo/bufferpool"

	"image"

	"github.com/spf13/hugo/deps"
)

// Some of the template funcs are'nt entirely stateless.
type templateFuncster struct {
	funcMap        template.FuncMap
	cachedPartials partialCache
	image          *imageHandler

	// Make sure each funcster gets its own TemplateFinder to get
	// proper text and HTML template separation.
	Tmpl templateFuncsterTemplater

	*deps.Deps
}

func newTemplateFuncster(deps *deps.Deps, t templateFuncsterTemplater) *templateFuncster {
	return &templateFuncster{
		Deps:           deps,
		Tmpl:           t,
		cachedPartials: partialCache{p: make(map[string]interface{})},
		image:          &imageHandler{fs: deps.Fs, imageConfigCache: map[string]image.Config{}},
	}
}

// Partial executes the named partial and returns either a string,
// when called from text/template, for or a template.HTML.
func (t *templateFuncster) partial(name string, contextList ...interface{}) (interface{}, error) {
	if strings.HasPrefix("partials/", name) {
		name = name[8:]
	}
	var context interface{}

	if len(contextList) == 0 {
		context = nil
	} else {
		context = contextList[0]
	}

	for _, n := range []string{"partials/" + name, "theme/partials/" + name} {
		templ := t.Tmpl.Lookup(n)
		if templ == nil {
			// For legacy reasons.
			templ = t.Tmpl.Lookup(n + ".html")
		}
		if templ != nil {
			b := bp.GetBuffer()
			defer bp.PutBuffer(b)

			if err := templ.Execute(b, context); err != nil {
				return "", err
			}

			switch t.Tmpl.(type) {
			case *htmlTemplates:
				return template.HTML(b.String()), nil
			case *textTemplates:
				return b.String(), nil
			default:
				panic("Unknown type")
			}
		}
	}

	return "", fmt.Errorf("Partial %q not found", name)
}
