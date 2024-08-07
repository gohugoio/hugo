// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package js provides functions for building JavaScript resources
package js

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_transformers/js"
	template "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

func (ns *Namespace) Bundle(id string, store *maps.Scratch) (Bundler, error) {
	key := path.Join(nsBundle, id)
	b := store.GetOrCreate(key, func() any {
		return &bundler{id: id, scriptOnes: make(map[string]*scriptOne), scriptManys: make(map[string]*scriptMany), client: ns}
	})
	return b.(*bundler), nil
}

func (b *bundler) UseScriptOne(id string) BundleScriptOne {
	b.mu.Lock()

	one, found := b.scriptOnes[id]
	if !found {
		one = &scriptOne{id: id, client: b.client}
		b.scriptOnes[id] = one
	}

	b.mu.Unlock()
	one.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BundleScriptOneOps
		types.Closer
	}{
		one,
		close(func() error {
			closeOnce.Do(func() {
				one.mu.Unlock()
			})
			return nil
		}),
	}
}

func (b *bundler) UseScriptMany(id string) BundleScriptMany {
	b.mu.Lock()

	many, found := b.scriptManys[id]
	if !found {
		many = &scriptMany{id: id, client: b.client, items: make(map[string]*scriptManyItem)}
		b.scriptManys[id] = many
	}

	b.mu.Unlock()
	many.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BundleScriptManyOps
		types.Closer
	}{
		many,
		close(func() error {
			closeOnce.Do(func() {
				many.mu.Unlock()
			})
			return nil
		}),
	}
}

type Bundler interface {
	UseScriptOne(id string) BundleScriptOne
	UseScriptMany(id string) BundleScriptMany
	Build() (*Package, error)
}

type BundleScriptOne interface {
	BundleScriptOneOps
	types.Closer
}

type BundleScriptOneOps interface {
	ResourceGetSetter
	SetInstance(opts any) string
}

type BundleScriptMany interface {
	BundleScriptManyOps
	types.Closer
}

type BundleScriptItem interface {
	BundleScriptItemOps
	types.Closer
}

type BundleScriptItemOps interface {
	ResourceGetSetter
	AddInstance(id string, opts any) string
}

type BundleScriptManyOps interface {
	GetCallback() resource.Resource
	SetCallback(r resource.Resource) string
	UseItem(id string) BundleScriptItem
}

type ScriptItem interface{}

type BundleCommonScriptOps interface {
	ResourceGetSetter
}

type ResourceGetSetter interface {
	GetResource() resource.Resource
	SetResource(r resource.Resource) string
}

type close func() error

func (c close) Close() error {
	return c()
}

var (
	_ Bundler             = (*bundler)(nil)
	_ BundleScriptOneOps  = (*scriptOne)(nil)
	_ BundleScriptManyOps = (*scriptMany)(nil)
)

func (b *scriptOne) SetInstance(opts any) string {
	/*if b.r == nil {
		panic("resource not set")
	}
	if id == "" {
		panic("id not set")
	}

	paramsm := cast.ToStringMap(params)
	b.instances[id] = &batchInstance{params: paramsm}
	*/ // TODO1

	return ""
}

func decodeScriptInstance(opts any) scriptInstance {
	var inst scriptInstance
	if err := mapstructure.WeakDecode(opts, &inst); err != nil {
		panic(err)
	}
	return inst
}

func (b *scriptManyItem) AddInstance(id string, opts any) string {
	b.instances[id] = decodeScriptInstance(opts)
	return ""
}

func (b *scriptMany) GetCallback() resource.Resource {
	if resource.StaleVersion(b.callback) > 0 {
		// Allow the client to set a new resource.
		return nil
	}
	return b.callback
}

func (b *scriptMany) SetCallback(r resource.Resource) string {
	b.callback = r
	return ""
}

func (b *scriptMany) UseItem(id string) BundleScriptItem {
	item, found := b.items[id]
	if !found {
		item = &scriptManyItem{id: id, instances: make(map[string]scriptInstance), client: b.client}
		b.items[id] = item
	}

	item.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BundleScriptItemOps
		types.Closer
	}{
		item,
		close(func() error {
			closeOnce.Do(func() {
				item.mu.Unlock()
			})
			return nil
		}),
	}
}

type batchTemplateContext struct {
	keyPath            string
	ID                 string
	CallbackImportPath string
	Modules            []batchTemplateExecutionsContext
}

type batchTemplateExecutionsContext struct {
	ID         string                   `json:"id"`
	ImportPath string                   `json:"importPath"`
	Instances  []batchTemplateExecution `json:"instances"`

	r resource.Resource
}

func (b batchTemplateExecutionsContext) CallbackJSON(i int) string {
	mod := fmt.Sprintf("Mod%d", i)

	v := struct {
		Mod string `json:"mod"`
		batchTemplateExecutionsContext
	}{
		mod,
		b,
	}

	bb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s := string(bb)

	s = strings.ReplaceAll(s, fmt.Sprintf("%q", mod), mod)

	return s
}

type batchTemplateExecution struct {
	ID     string `json:"id"`
	Params any    `json:"params"`
}

type batchBuildOpts struct {
	Callback           resource.Resource
	js.ExternalOptions `mapstructure:",squash"`
}

type scriptsOne struct {
	mu sync.Mutex

	id      string
	batches map[string]*scriptOne

	client *Namespace
}

type scriptsMany struct {
	mu sync.Mutex

	id      string
	batches map[string]*scriptMany

	client *Namespace
}

type scriptOne struct {
	mu sync.Mutex
	id string

	resourceGetSet

	client *Namespace
}

type scriptInstance struct {
	Params map[string]any
}

type scriptManyItem struct {
	mu sync.Mutex
	id string

	resourceGetSet
	instances map[string]scriptInstance

	client *Namespace
}

type scriptMany struct {
	mu       sync.Mutex
	id       string
	callback resource.Resource

	items map[string]*scriptManyItem

	client *Namespace
}

type bundler struct {
	mu          sync.Mutex
	id          string
	scriptOnes  map[string]*scriptOne
	scriptManys map[string]*scriptMany

	client *Namespace
}

type resourceGetSet struct {
	r resource.Resource
}

func (r *resourceGetSet) GetResource() resource.Resource {
	if resource.StaleVersion(r.r) > 0 {
		// Allow the client to set a new resource.
		return nil
	}
	return r.r
}

func (r *resourceGetSet) SetResource(res resource.Resource) string {
	r.r = res
	return ""
}

var (
	_ resource.StaleInfo                    = (*Package)(nil)
	_ identity.IsProbablyDependencyProvider = (*Package)(nil)
	_ identity.Identity                     = (*Package)(nil)
)

// TODO1 names.
type Package struct {
	id           string
	staleVersion uint32
	b            *bundler
	Scripts      map[string]resource.Resource
}

func (p *Package) IdentifierBase() string {
	return p.id
}

func (p *Package) StaleVersion() uint32 {
	p.b.mu.Lock()
	defer p.b.mu.Unlock()
	if p.staleVersion == 0 {
		p.staleVersion = p.calculateStaleVersion()
	}
	return p.staleVersion
}

func (p *Package) IsProbablyDependency(other identity.Identity) bool {
	depsFinder := identity.NewFinder(identity.FinderConfig{})
	var b bool
	p.forEeachResource(func(rr resource.Resource) bool {
		identity.WalkIdentitiesShallow(other, func(level int, left identity.Identity) bool {
			identity.WalkIdentitiesShallow(rr, func(level int, right identity.Identity) bool {
				if i := depsFinder.Contains(left, right, -1); i > 0 {
					b = true
				}
				return b
			})
			return b
		})
		return b
	})

	// TODO1 why is this called twice on change?

	return b
}

func (p *Package) forEeachResource(f func(r resource.Resource) bool) {
	for _, v := range p.b.scriptManys {
		if b := func() bool {
			v.mu.Lock()
			defer v.mu.Unlock()
			if v.callback != nil {
				if f(v.callback) {
					return true
				}
			}
			for _, vv := range v.items {
				vv.mu.Lock()
				defer vv.mu.Unlock()
				if f(vv.r) {
					return true
				}
			}
			return false
		}(); b {
			return
		}
	}

	for _, v := range p.b.scriptOnes {
		if b := func() bool {
			v.mu.Lock()
			defer v.mu.Unlock()
			if f(v.r) {
				return true
			}
			return false
		}(); b {
			return
		}
	}
}

func (p *Package) calculateStaleVersion() uint32 {
	// Return the first 0 zero value of the resources in this bundle.
	var i uint32
	p.forEeachResource(func(r resource.Resource) bool {
		if i = resource.StaleVersion(r); i > 0 {
			return true
		}
		return false
	})

	return i
}

func (b *bundler) Build() (*Package, error) {
	key := dynacache.CleanKey(b.id + ".js")
	p, err := b.client.bundlesCache.GetOrCreate(key, func(string) (*Package, error) {
		return b.build()
	})

	if p.b != b {
		panic("bundler mismatch")
	}

	return p, err
}

func (b *bundler) build() (*Package, error) {
	defer herrors.Recover() // TODO1
	b.mu.Lock()
	defer b.mu.Unlock()

	keyPath := b.id

	idResource := make(map[string]resource.Resource)
	impMap := make(map[string]resource.Resource)
	var entryPoints []string
	addEntryPoint := func(id, s string, r resource.Resource) {
		impMap[s] = r
		entryPoints = append(entryPoints, s)
		idResource[id] = bundleResource(s)
	}

	if len(b.scriptOnes) > 0 {
		for k, v := range b.scriptOnes {
			if v.r == nil {
				return nil, fmt.Errorf("resource not set for %q", k)
			}
			keyPath := path.Join(keyPath, k)
			resourcePath := paths.AddLeadingSlash(keyPath + v.r.MediaType().FirstSuffix.FullSuffix)
			addEntryPoint(k, resourcePath, v.r)

		}
	}

	if len(b.scriptManys) > 0 {
		for k, v := range b.scriptManys {
			keyPath := path.Join(keyPath, k)
			bopts := batchBuildOpts{
				Callback: v.callback,
			}
			var callbackImpPath string
			if bopts.Callback != nil {
				callbackImpPath = paths.AddLeadingSlash(keyPath + "_callback" + bopts.Callback.MediaType().FirstSuffix.FullSuffix)
				addEntryPoint(k, callbackImpPath, bopts.Callback)
			}

			t := &batchTemplateContext{
				keyPath:            keyPath,
				ID:                 v.id,
				CallbackImportPath: callbackImpPath,
			}

			for kk, vv := range v.items {
				keyPath := path.Join(keyPath, kk)
				bt := batchTemplateExecutionsContext{
					ID:         kk,
					r:          vv.r,
					ImportPath: keyPath + vv.r.MediaType().FirstSuffix.FullSuffix,
				}
				impMap[bt.ImportPath] = vv.r
				for kkk, vvv := range vv.instances {
					bt.Instances = append(bt.Instances, batchTemplateExecution{ID: kkk, Params: vvv.Params})
					sort.Slice(bt.Instances, func(i, j int) bool {
						return bt.Instances[i].ID < bt.Instances[j].ID
					})
				}
				t.Modules = append(t.Modules, bt)
			}
			sort.Slice(t.Modules, func(i, j int) bool {
				return t.Modules[i].ID < t.Modules[j].ID
			})

			r, s, err := b.client.buildBatch(t)
			if err != nil {
				return nil, err
			}
			addEntryPoint(v.id, s, r)
		}
	}

	target := "es2018"

	outDir := filepath.Join(b.client.d.Paths.AbsPublishDir, "js", "bundles", b.id)

	jopts := js.Options{
		ExternalOptions: js.ExternalOptions{
			Format:  "esm",
			Target:  target,
			Defines: map[string]any{
				//"process.env.NODE_ENV": `"development"`,
			},
		},
		InternalOptions: js.InternalOptions{
			OutDir: outDir,
			// TODO1 maybe not.
			Write:          true,
			AllowOverwrite: true,
			Splitting:      true,
			ImportOnResolveFunc: func(imp string) string {
				if _, found := impMap[imp]; found {
					return imp
				}
				return ""
			},
			ImportOnLoadFunc: func(imp string) string {
				if r, found := impMap[imp]; found {
					content, err := r.(resource.ContentProvider).Content(context.Background()) // TODO1
					if err != nil {
						panic(err)
					}
					return cast.ToString(content)
				}

				return ""
			},
			EntryPoints: entryPoints,
		},
	}

	if err := b.client.client.BuildBundle(jopts); err != nil {
		return nil, err
	}

	return &Package{b: b, id: path.Join(nsBundle, b.id), Scripts: idResource}, nil
}

type bundleResource string

func (b bundleResource) Name() string {
	return path.Base(string(b))
}

func (b bundleResource) Title() string {
	return b.Name()
}

func (b bundleResource) RelPermalink() string {
	return "/js/bundles" + string(b)
}

func (b bundleResource) Permalink() string {
	panic("not implemented")
}

func (b bundleResource) ResourceType() string {
	panic("not implemented")
}

func (b bundleResource) MediaType() media.Type {
	panic("not implemented")
}

func (b bundleResource) Data() any {
	panic("not implemented")
}

func (b bundleResource) Err() resource.ResourceError {
	return nil
}

func (b bundleResource) Params() maps.Params {
	panic("not implemented")
}

const nsBundle = "__hugo-js-bundle"

func (ns *Namespace) buildBatch(t *batchTemplateContext) (resource.Resource, string, error) {
	var buf bytes.Buffer
	if err := batchEsmCallbackTemplate.Execute(&buf, t); err != nil {
		return nil, "", err
	}

	s := paths.AddLeadingSlash(t.keyPath + ".js")
	r, err := ns.createClient.FromString(s, buf.String())
	if err != nil {
		return nil, "", err
	}

	return r, s, nil
}

//go:embed batch-esm-callback.gotmpl
var batchEsmCallbackTemplateString string
var batchEsmCallbackTemplate *template.Template

func init() {
	batchEsmCallbackTemplate = template.Must(template.New("batch-esm-callback").Parse(batchEsmCallbackTemplateString))
}
