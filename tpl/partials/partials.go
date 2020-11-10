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

// Package partials provides template functions for working with reusable
// templates.
package partials

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/identity"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/tpl"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
)

// TestTemplateProvider is global deps.ResourceProvider.
// NOTE: It's currently unused.
var TestTemplateProvider deps.ResourceProvider

type partialCacheKey struct {
	name    string
	variant interface{}
}

type partialCacheEntry struct {
	templateIdentity identity.Identity
	v                interface{}
}

// partialCache represents a cache of partials protected by a mutex.
type partialCache struct {
	sync.RWMutex
	p map[partialCacheKey]partialCacheEntry
}

func (p *partialCache) clear() {
	p.Lock()
	defer p.Unlock()
	p.p = make(map[partialCacheKey]partialCacheEntry)
}

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	cache := &partialCache{p: make(map[partialCacheKey]partialCacheEntry)}
	deps.BuildStartListeners.Add(
		func() {
			cache.clear()
		})

	return &Namespace{
		deps:           deps,
		cachedPartials: cache,
	}
}

// Namespace provides template functions for the "templates" namespace.
type Namespace struct {
	deps           *deps.Deps
	cachedPartials *partialCache
}

// contextWrapper makes room for a return value in a partial invocation.
type contextWrapper struct {
	Arg    interface{}
	Result interface{}
}

// Set sets the return value and returns an empty string.
func (c *contextWrapper) Set(in interface{}) string {
	c.Result = in
	return ""
}

// Include executes the named partial.
// If the partial contains a return statement, that value will be returned.
// Else, the rendered output will be returned:
// A string if the partial is a text/template, or template.HTML when html/template.
// Note that ctx is provided by Hugo, not the end user. TODO1 ctx used?
func (ns *Namespace) Include(ctx context.Context, name string, dataList ...interface{}) (interface{}, error) {
	v, id, err := ns.include(ctx, name, dataList...)
	if err != nil {
		return nil, err
	}
	if ns.deps.Running {
		// Track the usage of this partial so we know when to re-render pages using it.
		tpl.AddIdentiesToDataContext(ctx, id)
	}
	return v, nil
}

func (ns *Namespace) include(ctx context.Context, name string, dataList ...interface{}) (interface{}, identity.Identity, error) {
	name = strings.TrimPrefix(name, "partials/")

	var data interface{}
	if len(dataList) > 0 {
		data = dataList[0]
	}

	n := "partials/" + name
	templ, found := ns.deps.Tmpl().Lookup(n)

	if !found {
		// For legacy reasons.
		templ, found = ns.deps.Tmpl().Lookup(n + ".html")
	}

	if !found {
		return "", nil, fmt.Errorf("partial %q not found", name)
	}

	var info tpl.ParseInfo
	if ip, ok := templ.(tpl.Info); ok {
		info = ip.ParseInfo()
	}

	var w io.Writer

	if info.HasReturn {
		// Wrap the context sent to the template to capture the return value.
		// Note that the template is rewritten to make sure that the dot (".")
		// and the $ variable points to Arg.
		data = &contextWrapper{
			Arg: data,
		}

		// We don't care about any template output.
		w = ioutil.Discard
	} else {
		b := bp.GetBuffer()
		defer bp.PutBuffer(b)
		w = b
	}

	if err := ns.deps.Tmpl().ExecuteWithContext(ctx, templ, w, data); err != nil {
		return "", nil, err
	}

	var result interface{}

	if ctx, ok := data.(*contextWrapper); ok {
		result = ctx.Result
	} else if _, ok := templ.(*texttemplate.Template); ok {
		result = w.(fmt.Stringer).String()
	} else {
		result = template.HTML(w.(fmt.Stringer).String())
	}

	if ns.deps.Metrics != nil {
		ns.deps.Metrics.TrackValue(templ.Name(), result)
	}

	return result, templ.(identity.Identity), nil
}

// IncludeCached executes and caches partial templates.  The cache is created with name+variants as the key.
// Note that ctx is provided by Hugo and not the end user.
func (ns *Namespace) IncludeCached(ctx context.Context, name string, data interface{}, variants ...interface{}) (interface{}, error) {
	key, err := createKey(name, variants...)
	if err != nil {
		return nil, err
	}

	result, err := ns.getOrCreate(ctx, key, data)
	if err == errUnHashable {
		// Try one more
		key.variant = helpers.HashString(key.variant)
		result, err = ns.getOrCreate(ctx, key, data)
	}

	if ns.deps.Running {
		// Track the usage of this partial so we know when to re-render pages using it.
		tpl.AddIdentiesToDataContext(ctx, result.templateIdentity)
	}

	return result.v, err
}

func createKey(name string, variants ...interface{}) (partialCacheKey, error) {
	var variant interface{}

	if len(variants) > 1 {
		variant = helpers.HashString(variants...)
	} else if len(variants) == 1 {
		variant = variants[0]
		t := reflect.TypeOf(variant)
		switch t.Kind() {
		// This isn't an exhaustive list of unhashable types.
		// There may be structs with slices,
		// but that should be very rare. We do recover from that situation
		// below.
		case reflect.Slice, reflect.Array, reflect.Map:
			variant = helpers.HashString(variant)
		}
	}

	return partialCacheKey{name: name, variant: variant}, nil
}

var errUnHashable = errors.New("unhashable")

func (ns *Namespace) getOrCreate(ctx context.Context, key partialCacheKey, dot interface{}) (pe partialCacheEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			if strings.Contains(err.Error(), "unhashable type") {
				ns.cachedPartials.RUnlock()
				err = errUnHashable
			}
		}
	}()

	ns.cachedPartials.RLock()
	p, ok := ns.cachedPartials.p[key]
	ns.cachedPartials.RUnlock()

	if ok {
		return p, nil
	}

	v, id, err := ns.include(ctx, key.name, dot)
	if err != nil {
		return
	}

	ns.cachedPartials.Lock()
	defer ns.cachedPartials.Unlock()
	// Double-check.
	if p2, ok := ns.cachedPartials.p[key]; ok {
		return p2, nil
	}

	pe = partialCacheEntry{
		templateIdentity: id,
		v:                v,
	}
	ns.cachedPartials.p[key] = pe

	return
}
