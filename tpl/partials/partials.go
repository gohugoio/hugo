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
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/bep/lazycache"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/identity"

	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/tpl"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
)

type partialCacheKey struct {
	Name     string
	Variants []any
}
type includeResult struct {
	name     string
	result   any
	mangager identity.Manager
	err      error
}

func (k partialCacheKey) Key() string {
	if k.Variants == nil {
		return k.Name
	}
	return hashing.HashString(append([]any{k.Name}, k.Variants...)...)
}

func (k partialCacheKey) templateName() string {
	if !strings.HasPrefix(k.Name, "partials/") {
		return "partials/" + k.Name
	}
	return k.Name
}

// partialCache represents a LRU cache of partials.
type partialCache struct {
	cache *lazycache.Cache[string, includeResult]
}

func (p *partialCache) clear() {
	p.cache.DeleteFunc(func(s string, r includeResult) bool {
		return true
	})
}

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	// This lazycache was introduced in Hugo 0.111.0.
	// We're going to expand and consolidate all memory caches in Hugo using this,
	// so just set a high limit for now.
	lru := lazycache.New(lazycache.Options[string, includeResult]{MaxEntries: 1000})

	cache := &partialCache{cache: lru}
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
	res := ns.includWithTimeout(ctx, name, contextList...)
	if res.err != nil {
		return nil, res.err
	}

	if ns.deps.Metrics != nil {
		ns.deps.Metrics.TrackValue(res.name, res.result, false)
	}

	return res.result, nil
}

func (ns *Namespace) includWithTimeout(ctx context.Context, name string, dataList ...any) includeResult {
	// Create a new context with a timeout not connected to the incoming context.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), ns.deps.Conf.Timeout())
	defer cancel()

	res := make(chan includeResult, 1)

	go func() {
		res <- ns.include(ctx, name, dataList...)
	}()

	select {
	case r := <-res:
		return r
	case <-timeoutCtx.Done():
		err := timeoutCtx.Err()
		if err == context.DeadlineExceeded {
			//lint:ignore ST1005 end user message.
			err = fmt.Errorf("partial %q timed out after %s. This is most likely due to infinite recursion. If this is just a slow template, you can try to increase the 'timeout' config setting.", name, ns.deps.Conf.Timeout())
		}
		return includeResult{err: err}
	}
}

// include is a helper function that lookups and executes the named partial.
// Returns the final template name and the rendered output.
func (ns *Namespace) include(ctx context.Context, name string, dataList ...any) includeResult {
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
		return includeResult{err: fmt.Errorf("partial %q not found", name)}
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
		return includeResult{err: err}
	}

	var result any

	if ctx, ok := data.(*contextWrapper); ok {
		result = ctx.Result
	} else if _, ok := templ.(*texttemplate.Template); ok {
		result = w.(fmt.Stringer).String()
	} else {
		result = template.HTML(w.(fmt.Stringer).String())
	}

	return includeResult{
		name:   templ.Name(),
		result: result,
	}
}

// IncludeCached executes and caches partial templates.  The cache is created with name+variants as the key.
// Note that ctx is provided by Hugo, not the end user.
func (ns *Namespace) IncludeCached(ctx context.Context, name string, context any, variants ...any) (any, error) {
	start := time.Now()
	key := partialCacheKey{
		Name:     name,
		Variants: variants,
	}
	depsManagerIn := tpl.Context.GetDependencyManagerInCurrentScope(ctx)

	r, found, err := ns.cachedPartials.cache.GetOrCreate(key.Key(), func(string) (includeResult, error) {
		var depsManagerShared identity.Manager
		if ns.deps.Conf.Watching() {
			// We need to create a shared dependency manager to pass downwards
			// and add those same dependencies to any cached invocation of this partial.
			depsManagerShared = identity.NewManager("partials")
			ctx = tpl.Context.DependencyManagerScopedProvider.Set(ctx, depsManagerShared.(identity.DependencyManagerScopedProvider))
		}
		r := ns.includWithTimeout(ctx, key.Name, context)
		if ns.deps.Conf.Watching() {
			r.mangager = depsManagerShared
		}
		return r, r.err
	})
	if err != nil {
		return nil, err
	}

	if ns.deps.Metrics != nil {
		if found {
			// The templates that gets executed is measured in Execute.
			// We need to track the time spent in the cache to
			// get the totals correct.
			ns.deps.Metrics.MeasureSince(key.templateName(), start)
		}
		ns.deps.Metrics.TrackValue(key.templateName(), r.result, found)
	}

	if r.mangager != nil && depsManagerIn != nil {
		depsManagerIn.AddIdentity(r.mangager)
	}

	return r.result, nil
}
