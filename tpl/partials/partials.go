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
	"reflect"
	"strings"
	"sync"
	"time"

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
	variant any
}

func (k partialCacheKey) templateName() string {
	if !strings.HasPrefix(k.name, "partials/") {
		return "partials/" + k.name
	}
	return k.name
}

// partialCache represents a cache of partials protected by a mutex.
type partialCache struct {
	sync.RWMutex
	p map[partialCacheKey]any
}

func (p *partialCache) clear() {
	p.Lock()
	defer p.Unlock()
	p.p = make(map[partialCacheKey]any)
}

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	cache := &partialCache{p: make(map[partialCacheKey]any)}
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
	Arg    any
	Result any
}

// Set sets the return value and returns an empty string.
func (c *contextWrapper) Set(in any) string {
	c.Result = in
	return ""
}

// Include executes the named partial.
// If the partial contains a return statement, that value will be returned.
// Else, the rendered output will be returned:
// A string if the partial is a text/template, or template.HTML when html/template.
// Note that ctx is provided by Hugo, not the end user.
func (ns *Namespace) Include(ctx context.Context, name string, contextList ...any) (any, error) {
	name, result, err := ns.include(ctx, name, contextList...)
	if err != nil {
		return result, err
	}

	if ns.deps.Metrics != nil {
		ns.deps.Metrics.TrackValue(name, result, false)
	}

	return result, nil
}

// include is a helper function that lookups and executes the named partial.
// Returns the final template name and the rendered output.
func (ns *Namespace) include(ctx context.Context, name string, dataList ...any) (string, any, error) {
	var data any
	if len(dataList) > 0 {
		data = dataList[0]
	}

	var n string
	if strings.HasPrefix(name, "partials/") {
		n = name
	} else {
		n = "partials/" + name
	}

	templ, found := ns.deps.Tmpl().Lookup(n)
	if !found {
		// For legacy reasons.
		templ, found = ns.deps.Tmpl().Lookup(n + ".html")
	}

	if !found {
		return "", "", fmt.Errorf("partial %q not found", name)
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
		w = io.Discard
	} else {
		b := bp.GetBuffer()
		defer bp.PutBuffer(b)
		w = b
	}

	if err := ns.deps.Tmpl().ExecuteWithContext(ctx, templ, w, data); err != nil {
		return "", nil, err
	}

	var result any

	if ctx, ok := data.(*contextWrapper); ok {
		result = ctx.Result
	} else if _, ok := templ.(*texttemplate.Template); ok {
		result = w.(fmt.Stringer).String()
	} else {
		result = template.HTML(w.(fmt.Stringer).String())
	}

	return templ.Name(), result, nil
}

// IncludeCached executes and caches partial templates.  The cache is created with name+variants as the key.
// Note that ctx is provided by Hugo, not the end user.
func (ns *Namespace) IncludeCached(ctx context.Context, name string, context any, variants ...any) (any, error) {
	key, err := createKey(name, variants...)
	if err != nil {
		return nil, err
	}

	result, err := ns.getOrCreate(ctx, key, context)
	if err == errUnHashable {
		// Try one more
		key.variant = helpers.HashString(key.variant)
		result, err = ns.getOrCreate(ctx, key, context)
	}

	return result, err
}

func createKey(name string, variants ...any) (partialCacheKey, error) {
	var variant any

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

func (ns *Namespace) getOrCreate(ctx context.Context, key partialCacheKey, context any) (result any, err error) {
	start := time.Now()
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
		if ns.deps.Metrics != nil {
			ns.deps.Metrics.TrackValue(key.templateName(), p, true)
			// The templates that gets executed is measured in Execute.
			// We need to track the time spent in the cache to
			// get the totals correct.
			ns.deps.Metrics.MeasureSince(key.templateName(), start)

		}
		return p, nil
	}

	// This needs to be done outside the lock.
	// See #9588
	_, p, err = ns.include(ctx, key.name, context)
	if err != nil {
		return nil, err
	}

	ns.cachedPartials.Lock()
	defer ns.cachedPartials.Unlock()
	// Double-check.
	if p2, ok := ns.cachedPartials.p[key]; ok {
		if ns.deps.Metrics != nil {
			ns.deps.Metrics.TrackValue(key.templateName(), p, true)
			ns.deps.Metrics.MeasureSince(key.templateName(), start)
		}
		return p2, nil

	}
	if ns.deps.Metrics != nil {
		ns.deps.Metrics.TrackValue(key.templateName(), p, false)
	}

	ns.cachedPartials.p[key] = p

	return p, nil
}
