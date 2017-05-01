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

package partials

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
	texttemplate "text/template"

	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/deps"
)

var TestTemplateProvider deps.ResourceProvider

// partialCache represents a cache of partials protected by a mutex.
type partialCache struct {
	sync.RWMutex
	p map[string]interface{}
}

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps:           deps,
		cachedPartials: partialCache{p: make(map[string]interface{})},
	}
}

// Namespace provides template functions for the "templates" namespace.
type Namespace struct {
	deps           *deps.Deps
	cachedPartials partialCache
}

// Include executes the named partial and returns either a string,
// when the partial is a text/template, or template.HTML when html/template.
func (ns *Namespace) Include(name string, contextList ...interface{}) (interface{}, error) {
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
		templ := ns.deps.Tmpl.Lookup(n)
		if templ == nil {
			// For legacy reasons.
			templ = ns.deps.Tmpl.Lookup(n + ".html")
		}
		if templ != nil {
			b := bp.GetBuffer()
			defer bp.PutBuffer(b)

			if err := templ.Execute(b, context); err != nil {
				return "", err
			}

			if _, ok := templ.Template.(*texttemplate.Template); ok {
				return b.String(), nil
			}

			return template.HTML(b.String()), nil

		}
	}

	return "", fmt.Errorf("Partial %q not found", name)
}

// getCached executes and caches partial templates.  An optional variant
// string parameter (a string slice actually, but be only use a variadic
// argument to make it optional) can be passed so that a given partial can have
// multiple uses. The cache is created with name+variant as the key.
func (ns *Namespace) getCached(name string, context interface{}, variant ...string) (interface{}, error) {
	key := name
	if len(variant) > 0 {
		for i := 0; i < len(variant); i++ {
			key += variant[i]
		}
	}
	return ns.getOrCreate(key, name, context)
}

func (ns *Namespace) getOrCreate(key, name string, context interface{}) (p interface{}, err error) {
	var ok bool

	ns.cachedPartials.RLock()
	p, ok = ns.cachedPartials.p[key]
	ns.cachedPartials.RUnlock()

	if ok {
		return
	}

	ns.cachedPartials.Lock()
	if p, ok = ns.cachedPartials.p[key]; !ok {
		ns.cachedPartials.Unlock()
		p, err = ns.Include(name, context)

		ns.cachedPartials.Lock()
		ns.cachedPartials.p[key] = p

	}
	ns.cachedPartials.Unlock()

	return
}
