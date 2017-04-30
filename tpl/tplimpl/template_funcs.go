// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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
	"sync"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/tpl/internal"

	// Init the namespaces
	_ "github.com/spf13/hugo/tpl/collections"
	_ "github.com/spf13/hugo/tpl/compare"
	_ "github.com/spf13/hugo/tpl/crypto"
	_ "github.com/spf13/hugo/tpl/data"
	_ "github.com/spf13/hugo/tpl/encoding"
	_ "github.com/spf13/hugo/tpl/lang"
	_ "github.com/spf13/hugo/tpl/math"
	_ "github.com/spf13/hugo/tpl/strings"
)

// Get retrieves partial output from the cache based upon the partial name.
// If the partial is not found in the cache, the partial is rendered and added
// to the cache.
func (t *templateFuncster) Get(key, name string, context interface{}) (p interface{}, err error) {
	var ok bool

	t.cachedPartials.RLock()
	p, ok = t.cachedPartials.p[key]
	t.cachedPartials.RUnlock()

	if ok {
		return
	}

	t.cachedPartials.Lock()
	if p, ok = t.cachedPartials.p[key]; !ok {
		t.cachedPartials.Unlock()
		p, err = t.partial(name, context)

		t.cachedPartials.Lock()
		t.cachedPartials.p[key] = p

	}
	t.cachedPartials.Unlock()

	return
}

// partialCache represents a cache of partials protected by a mutex.
type partialCache struct {
	sync.RWMutex
	p map[string]interface{}
}

// partialCached executes and caches partial templates.  An optional variant
// string parameter (a string slice actually, but be only use a variadic
// argument to make it optional) can be passed so that a given partial can have
// multiple uses.  The cache is created with name+variant as the key.
func (t *templateFuncster) partialCached(name string, context interface{}, variant ...string) (interface{}, error) {
	key := name
	if len(variant) > 0 {
		for i := 0; i < len(variant); i++ {
			key += variant[i]
		}
	}
	return t.Get(key, name, context)
}

func (t *templateFuncster) initFuncMap() {
	funcMap := template.FuncMap{
		// Namespaces
		"images":  t.images.Namespace,
		"inflect": t.inflect.Namespace,
		"os":      t.os.Namespace,
		"safe":    t.safe.Namespace,
		//"time":        t.time.Namespace,
		"transform": t.transform.Namespace,
		"urls":      t.urls.Namespace,

		"absURL":        t.urls.AbsURL,
		"absLangURL":    t.urls.AbsLangURL,
		"dateFormat":    t.time.Format,
		"emojify":       t.transform.Emojify,
		"getenv":        t.os.Getenv,
		"highlight":     t.transform.Highlight,
		"htmlEscape":    t.transform.HTMLEscape,
		"htmlUnescape":  t.transform.HTMLUnescape,
		"humanize":      t.inflect.Humanize,
		"imageConfig":   t.images.Config,
		"int":           func(v interface{}) (int, error) { return cast.ToIntE(v) },
		"markdownify":   t.transform.Markdownify,
		"now":           t.time.Now,
		"partial":       t.partial,
		"partialCached": t.partialCached,
		"plainify":      t.transform.Plainify,
		"pluralize":     t.inflect.Pluralize,
		"print":         fmt.Sprint,
		"printf":        fmt.Sprintf,
		"println":       fmt.Sprintln,
		"readDir":       t.os.ReadDir,
		"readFile":      t.os.ReadFile,
		"ref":           t.urls.Ref,
		"relURL":        t.urls.RelURL,
		"relLangURL":    t.urls.RelLangURL,
		"relref":        t.urls.RelRef,
		"safeCSS":       t.safe.CSS,
		"safeHTML":      t.safe.HTML,
		"safeHTMLAttr":  t.safe.HTMLAttr,
		"safeJS":        t.safe.JS,
		"safeJSStr":     t.safe.JSStr,
		"safeURL":       t.safe.URL,
		"sanitizeURL":   t.safe.SanitizeURL,
		"sanitizeurl":   t.safe.SanitizeURL,
		"singularize":   t.inflect.Singularize,
		"string":        func(v interface{}) (string, error) { return cast.ToStringE(v) },
		"time":          t.time.AsTime,
		"urlize":        t.PathSpec.URLize,
	}

	// Merge the namespace funcs
	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(t.Deps)
		// TODO(bep) namespace ns.Context is a dummy func just to make this work.
		// Consider if we can add this context to the rendering context in an easy
		// way to make this cleaner. Maybe.
		funcMap[ns.Name] = ns.Context
		for k, v := range ns.Aliases {
			funcMap[k] = v
		}
	}

	t.funcMap = funcMap
	t.Tmpl.(*templateHandler).setFuncs(funcMap)
}
