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
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_transformers/js"
	template "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// TODO1 move/consolidate.
const (
	hugoVirtualNS = "@hugo-virtual"
	paramsNS      = "@params"
)

type Batcher interface {
	Config() OptionsSetter
	Group(id string) BatcherGroup
	Build() (*Package, error)
}

type BatcherGroup interface {
	Callback(id string) OptionsSetter
	Script(id string) OptionsSetter
	Instance(sid, iid string) OptionsSetter
}

type ScriptOptions struct {
	// The script to build.
	// TODO1 handle stale.
	Resource resource.Resource

	// The import context to use.
	// Note that we will always fall back to the resource's own import context.
	ImportContext resource.ResourceGetter

	// Params marshaled to JSON.
	Params json.RawMessage
}

func (o ScriptOptions) Compile(m map[string]any) (*ScriptOptions, error) {
	var s optionsGetSet // TODO1 type.
	if err := mapstructure.WeakDecode(m, &s); err != nil {
		return nil, err
	}

	paramsJSON, err := json.Marshal(s.Params)
	if err != nil {
		panic(err)
	}

	return &ScriptOptions{
		Resource:      s.Resource,
		ImportContext: resource.NewResourceGetter(s.ImportContext),
		Params:        paramsJSON,
	}, nil
}

func (o *ScriptOptions) Dir() string {
	return path.Dir(o.Resource.(resource.PathProvider).Path())
}

type ParamsOptions struct {
	Params json.RawMessage
}

type ScriptOptionsGetSetter interface {
	GetOptions() *ScriptOptions
	SetOptions(map[string]any) string
}

type OptionsSetter interface {
	SetOptions(map[string]any) string
}

func (ns *Namespace) Batch(id string, store *maps.Scratch) (Batcher, error) {
	key := path.Join(nsBundle, id)
	b := store.GetOrCreate(key, func() any {
		return &batcher{
			id:            id,
			scriptGroups:  make(map[string]*scriptGroup),
			client:        ns,
			configOptions: newOptions(),
		}
	})
	return b.(*batcher), nil
}

func (b *batcher) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, v := range b.scriptGroups {
		// TODO1 check if this is complete.
		v.Reset()
	}
}

func (b *batcher) Group(id string) BatcherGroup {
	b.mu.Lock()
	defer b.mu.Unlock()

	group, found := b.scriptGroups[id]
	if !found {
		group = &scriptGroup{
			id: id, client: b.client,
			scriptsOptions:   make(map[string]*options),
			instancesOptions: make(map[instanceID]*options),
			callbacksOptions: make(map[string]*options),
		}
		b.scriptGroups[id] = group
	}

	return group
}

var _ Batcher = (*batcher)(nil)

type batchTemplateContext struct {
	keyPath   string
	ID        string
	Callbacks []string
	Modules   []batchTemplateExecutionsContext
}

type batchTemplateExecutionsContext struct {
	ID         string                   `json:"id"`
	ImportPath string                   `json:"-"`
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
	ID     string          `json:"id"`
	Params json.RawMessage `json:"params"`
}

type scriptGroup struct {
	mu sync.Mutex

	id string

	client *Namespace

	scriptsOptions   map[string]*options
	instancesOptions map[instanceID]*options
	callbacksOptions map[string]*options

	// Compiled.
	scripts   scriptMap
	instances instanceMap
	callbacks scriptMap
}

func (g *scriptGroup) Reset() {
	for _, v := range g.scriptsOptions {
		v.Reset()
	}
	for _, v := range g.instancesOptions {
		v.Reset()
	}
	for _, v := range g.callbacksOptions {
		v.Reset()
	}
}

type script struct {
	ID string
	*ScriptOptions
}

type instance struct {
	instanceID
	*ParamsOptions
}

type (
	instanceMap map[instanceID]*ParamsOptions
	instances   []*instance
)

func (i instances) ByScriptID(id string) instances {
	var a instances
	for _, v := range i {
		if v.instanceID.scriptID == id {
			a = append(a, v)
		}
	}
	return a
}

func (p instanceMap) Sorted() instances {
	var a []*instance
	for k, v := range p {
		a = append(a, &instance{instanceID: k, ParamsOptions: v})
	}
	sort.Slice(a, func(i, j int) bool {
		ai := a[i]
		aj := a[j]
		if ai.instanceID.scriptID != aj.instanceID.scriptID {
			return ai.instanceID.scriptID < aj.instanceID.scriptID
		}
		return ai.instanceID.instanceID < aj.instanceID.instanceID
	})
	return a
}

type scriptMap map[string]*ScriptOptions

func (s scriptMap) Sorted() []*script {
	var a []*script
	for k, v := range s {
		a = append(a, &script{ID: k, ScriptOptions: v})
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID < a[j].ID
	})
	return a
}

func (s *scriptGroup) compile() error {
	// TODO1 lock?
	s.scripts = make(map[string]*ScriptOptions)
	s.instances = make(map[instanceID]*ParamsOptions)
	s.callbacks = make(map[string]*ScriptOptions)

	for k, v := range s.scriptsOptions {
		compiled, err := compileScriptOptions(v)
		if err != nil {
			return err
		}
		s.scripts[k] = compiled
	}

	for k, v := range s.instancesOptions {
		compiled, err := compileParamsOptions(v)
		if err != nil {
			return err
		}
		s.instances[k] = compiled
	}

	for k, v := range s.callbacksOptions {
		compiled, err := compileScriptOptions(v)
		if err != nil {
			return err
		}
		s.callbacks[k] = compiled
	}

	return nil
}

type instanceID struct {
	scriptID   string
	instanceID string
}

func (s *scriptGroup) Script(id string) OptionsSetter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, found := s.scriptsOptions[id]; found {
		return v.Get()
	}
	s.scriptsOptions[id] = newOptions()
	return s.scriptsOptions[id].Get()
}

func (s *scriptGroup) Instance(sid, iid string) OptionsSetter {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := instanceID{scriptID: sid, instanceID: iid}
	if v, found := s.instancesOptions[id]; found {
		return v.Get()
	}
	s.instancesOptions[id] = newOptions()
	return s.instancesOptions[id].Get()
}

func (s *scriptGroup) Callback(id string) OptionsSetter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, found := s.callbacksOptions[id]; found {
		return v.Get()
	}
	s.callbacksOptions[id] = newOptions()
	return s.callbacksOptions[id].Get()
}

type scriptGroups map[string]*scriptGroup

func (s scriptGroups) Sorted() []*scriptGroup {
	var a []*scriptGroup
	for _, v := range s {
		a = append(a, v)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].id < a[j].id
	})
	return a
}

type batcher struct {
	mu           sync.Mutex
	id           string
	scriptGroups scriptGroups

	client *Namespace

	configOptions *options

	// The last successfully built package.
	// If this is non-nil and not stale, we can reuse it (e.g. on server rebuilds)
	prevBuild *Package

	// Compiled.
	config js.ExternalOptions
}

func (b *batcher) compile() error {
	var err error
	b.config, err = js.DecodeExternalOptions(b.configOptions.commit().opts)
	if err != nil {
		return err
	}

	for _, v := range b.scriptGroups {
		if err := v.compile(); err != nil {
			return err
		}
	}
	return nil
}

func (b *batcher) Config() OptionsSetter {
	return b.configOptions.Get()
}

var _ resource.StaleInfo = (*options)(nil)

type options struct {
	getOnce[*optionsSetter]
}

func (o *options) Reset() {
	mu := o.once.ResetWithLock()
	o.v.staleVersion.Store(0)
	mu.Unlock()
}

func (o *options) StaleVersion() uint32 {
	return o.v.staleVersion.Load()
}

// func (o ScriptOptions) Compile(m map[string]any) (*ScriptOptions, error) {

func newOptions() *options {
	return &options{getOnce[*optionsSetter]{
		v: &optionsSetter{},
	}}
}

type optionsSetter struct {
	staleVersion atomic.Uint32
	opts         map[string]any
}

// TODO1 try to avoid stale page resources when changing the head.
func (o *optionsSetter) SetOptions(m map[string]any) string {
	if o.opts != nil {
		if reflect.DeepEqual(o.opts, m) {
			return ""
		}
		var isStale bool
		for k, v := range m {
			vv, found := o.opts[k]
			if !found {
				isStale = true
			} else {
				if si, ok := vv.(resource.StaleInfo); ok {
					isStale = si.StaleVersion() > 0
				} else {
					isStale = !reflect.DeepEqual(v, vv)
				}
			}

			if isStale {
				break
			}
		}

		if !isStale {
			return ""
		}

		o.staleVersion.Add(1)
	}

	o.opts = m

	return ""
}

type getOnce[T any] struct {
	v    T
	once lazy.OnceMore
}

func (g *getOnce[T]) Get() T {
	var v T
	g.once.Do(func() {
		v = g.v
	})
	return v
}

func (g *getOnce[T]) commit() T {
	g.once.Do(func() {})
	return g.v
}

func compileParamsOptions(o *options) (*ParamsOptions, error) {
	v := struct {
		Params map[string]any
	}{}

	m := o.commit().opts

	if err := mapstructure.WeakDecode(m, &v); err != nil {
		return nil, err
	}

	paramsJSON, err := json.Marshal(v.Params)
	if err != nil {
		return nil, err
	}

	return &ParamsOptions{
		Params: paramsJSON,
	}, nil
}

func compileScriptOptions(o *options) (*ScriptOptions, error) {
	v := struct {
		Resource      resource.Resource
		ImportContext any
		Params        map[string]any
	}{}

	m := o.commit().opts

	if err := mapstructure.WeakDecode(m, &v); err != nil {
		panic(err)
	}

	var paramsJSON []byte
	if v.Params != nil {
		var err error
		paramsJSON, err = json.Marshal(v.Params)
		if err != nil {
			panic(err)
		}
	}

	compiled := &ScriptOptions{
		Resource:      v.Resource,
		ImportContext: resource.NewResourceGetter(v.ImportContext),
		Params:        paramsJSON,
	}

	return compiled, nil
}

type optionsGetSet struct {
	Resource      resource.Resource
	ImportContext any
	Params        map[string]any

	// Compiled values
	compiled *ScriptOptions
}

func (s *optionsGetSet) GetOptions() *ScriptOptions {
	return s.compiled
}

func (s *optionsGetSet) SetOptions(m map[string]any) string {
	if err := mapstructure.WeakDecode(m, &s); err != nil {
		panic(err)
	}

	paramsJSON, err := json.Marshal(s.Params)
	if err != nil {
		panic(err)
	}

	s.compiled = &ScriptOptions{
		Resource:      s.Resource,
		ImportContext: resource.NewResourceGetter(s.ImportContext),
		Params:        paramsJSON,
	}

	return ""
}

var (
	_ resource.StaleInfo                    = (*Package)(nil)
	_ identity.IsProbablyDependencyProvider = (*Package)(nil)
	_ identity.Identity                     = (*Package)(nil)
)

// TODO1 names.
type Package struct {
	origin       *batcher
	outDir       string
	id           string
	staleVersion uint32
	b            *batcher
	Groups       map[string]resource.Resources
}

// TODO1 add a group.Subscribe (identity) method?

func (p *Package) IdentifierBase() string {
	return p.id
}

// TODO1 need this? We cannot cache this value for long.
func (p *Package) StaleVersion() uint32 {
	p.b.mu.Lock()
	defer p.b.mu.Unlock()
	if p.staleVersion == 0 {
		p.staleVersion = p.calculateStaleVersion()
	}
	return p.staleVersion
}

func (p *Package) MarkStale() {
	p.origin.reset()
}

func (p *Package) IsProbablyDependency(other identity.Identity) bool {
	depsFinder := identity.NewFinder(identity.FinderConfig{})
	var b bool
	p.forEeachStaleInfo(func(rr resource.StaleInfo) bool {
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

// You should not depend on the invocation order when calling this.
// TODO1 check that this does not get called on first build.
func (p *Package) forEeachStaleInfo(f func(si resource.StaleInfo) bool) {
	check := func(v any) bool {
		if si, ok := v.(resource.StaleInfo); ok {
			return f(si)
		}
		return false
	}
	for _, v := range p.b.scriptGroups {
		if b := func() bool {
			v.mu.Lock()
			defer v.mu.Unlock()

			for _, vv := range v.instancesOptions {
				if check(vv) {
					return true
				}
			}

			for _, vv := range v.scriptsOptions {
				if check(vv) {
					return true
				}
			}

			for _, vv := range v.callbacksOptions {
				if check(vv) {
					return true
				}
			}

			return false
		}(); b {
			return
		}
	}
}

func (p *Package) calculateStaleVersion() uint32 {
	// Return the first  non-zero value found.
	var i uint32
	p.forEeachStaleInfo(func(si resource.StaleInfo) bool {
		if i = si.StaleVersion(); i > 0 {
			return true
		}
		return false
	})

	return i
}

func logTime(name string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("%s in %s\n", name, elapsed)
}

func (b *batcher) Build() (*Package, error) {
	key := dynacache.CleanKey(b.id + ".js")

	p, err := b.client.bundlesCache.GetOrCreate(key, func(string) (*Package, error) {
		return b.build()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build batch: %w", err)
	}

	if p.b != b {
		panic("bundler mismatch")
	}

	return p, nil
}

// TODO1 remove.
func deb(what string, v ...any) {
	fmt.Println(what, v)
}

func (b *batcher) build() (*Package, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Use the unexported calculateStaleVersion
	if b.prevBuild != nil && b.prevBuild.calculateStaleVersion() == 0 {
		return b.prevBuild, nil
	}

	p, err := b.doBuild()
	if err != nil {
		return nil, err
	}
	b.prevBuild = p
	return p, nil
}

func (b *batcher) doBuild() (*Package, error) {
	start := time.Now()
	defer logTime("build", start) // TODO1 remove.

	keyPath := b.id

	type importContext struct {
		name           string
		resourceGetter resource.ResourceGetter
		scriptOptions  *ScriptOptions // TODO1 remove resourceGetter?
	}

	state := struct {
		importResource        *maps.Cache[string, resource.Resource]
		resultResource        *maps.Cache[string, resource.Resource]
		importerImportContext *maps.Cache[string, importContext]
		pathGroup             *maps.Cache[string, string]
	}{
		importResource:        maps.NewCache[string, resource.Resource](),
		resultResource:        maps.NewCache[string, resource.Resource](),
		importerImportContext: maps.NewCache[string, importContext](),
		pathGroup:             maps.NewCache[string, string](),
	}

	var entryPoints []string
	addResource := func(group, pth string, r resource.Resource, isResult bool) {
		state.pathGroup.Set(pth, group)
		state.importResource.Set(pth, r)
		if isResult {
			state.resultResource.Set(pth, r)
		}
		entryPoints = append(entryPoints, pth)
	}

	if err := b.compile(); err != nil {
		return nil, err
	}

	for k, v := range b.scriptGroups {
		keyPath := keyPath + "_" + k

		var callbacks []string
		for _, vv := range v.callbacks.Sorted() {
			callbackKeyPath := keyPath + "_" + vv.ID
			callbackImpPath := paths.AddLeadingSlash(callbackKeyPath + "_callback" + vv.Resource.MediaType().FirstSuffix.FullSuffix)
			callbacks = append(callbacks, callbackImpPath)
			addResource(k, callbackImpPath, vv.Resource, false)
		}

		t := &batchTemplateContext{
			keyPath:   keyPath,
			ID:        v.id,
			Callbacks: callbacks,
		}

		instances := v.instances.Sorted()

		for _, vv := range v.scripts.Sorted() {
			if vv.Resource == nil {
				// TODO1 others, init.
				return nil, fmt.Errorf("resource not set for %q", vv.ID)
			}
			keyPath := keyPath + "_" + vv.ID
			opts := vv.ScriptOptions
			impPath := path.Join(hugoVirtualNS, opts.Dir(), keyPath+opts.Resource.MediaType().FirstSuffix.FullSuffix)
			impCtx := opts.ImportContext

			state.importerImportContext.Set(impPath, importContext{
				name:           keyPath,
				resourceGetter: impCtx,
				scriptOptions:  opts,
			})

			bt := batchTemplateExecutionsContext{
				ID:         vv.ID,
				r:          vv.Resource,
				ImportPath: impPath,
			}
			state.importResource.Set(bt.ImportPath, vv.Resource)
			for _, vvv := range instances.ByScriptID(vv.ID) {
				bt.Instances = append(bt.Instances, batchTemplateExecution{ID: vvv.instanceID.instanceID, Params: vvv.Params})
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

		state.importerImportContext.Set(s, importContext{
			name:           s,
			resourceGetter: nil,
			scriptOptions:  nil,
		})

		addResource(v.id, s, r, true)
	}

	absPublishDir := b.client.d.AbsPublishDir
	mediaTypes := b.client.d.ResourceSpec.MediaTypes()
	cssMt, _, _ := mediaTypes.GetFirstBySuffix("css")

	cacheDir := filepath.Join(b.client.d.SourceSpec.Cfg.Dirs().CacheDir, "_jsbatch")
	if err := os.Mkdir(cacheDir, 0o777); err != nil && !herrors.IsExist(err) {
		return nil, err
	}
	outDir, err := os.MkdirTemp(cacheDir, "jsbatch")
	if err != nil {
		return nil, err
	}

	var importResulveMu sync.Mutex

	externalOptions := b.config
	externalOptions.Format = "esm" // Maybe allow other formats for simple 1 script setups. Also consider splitting below.

	jopts := js.Options{
		ExternalOptions: externalOptions,
		InternalOptions: js.InternalOptions{
			OutDir:         outDir,
			Write:          true,
			AllowOverwrite: true,
			Splitting:      true,
			ImportOnResolveFunc: func(imp string, args api.OnResolveArgs) string {
				importResulveMu.Lock()
				defer importResulveMu.Unlock()

				if _, found := state.importResource.Get(imp); found {
					return imp
				}

				var importContextPath string
				if args.Kind == api.ResolveEntryPoint {
					importContextPath = args.Path
				} else {
					importContextPath = args.Importer
				}
				importContext, _ := state.importerImportContext.Get(importContextPath)

				if importContext.resourceGetter != nil {
					resolved := importContext.resourceGetter.Get(imp)

					if resolved != nil {
						imp := hugoVirtualNS + resolved.(resource.PathProvider).Path()
						state.importResource.Set(imp, resolved)
						state.importerImportContext.Set(imp, importContext) // TODO1 test case (see headlessui)
						return imp

					}
				}
				return ""
			},
			ImportOnLoadFunc: func(args api.OnLoadArgs) string {
				importResulveMu.Lock()
				defer importResulveMu.Unlock()

				imp := args.Path

				if r, found := state.importResource.Get(imp); found {
					content, err := r.(resource.ContentProvider).Content(context.Background()) // TODO1
					if err != nil {
						panic(err)
					}
					return cast.ToString(content)
				}
				return ""
			},
			ImportParamsOnLoadFunc: func(args api.OnLoadArgs) json.RawMessage {
				if importContext, found := state.importerImportContext.Get(args.Path); found {
					if importContext.scriptOptions != nil {
						return importContext.scriptOptions.Params
					}
				}
				return nil
			},
			ErrorMessageResolveFunc: func(args api.Message) *js.ErrorMessageResolved {
				if loc := args.Location; loc != nil {
					path := strings.TrimPrefix(loc.File, "ns-hugo:") // TODO1
					if r, found := state.importResource.Get(path); found {
						path = strings.TrimPrefix(path, hugoVirtualNS)
						var contentr hugio.ReadSeekCloser
						if cp, ok := r.(hugio.ReadSeekCloserProvider); ok {
							contentr, _ = cp.ReadSeekCloser()
						}
						return &js.ErrorMessageResolved{
							Content: contentr,
							Path:    path,
							Message: args.Text,
						}

					}

				}
				return nil
			},
			EntryPoints: entryPoints,
		},
	}

	result, err := b.client.client.BuildBundle(jopts)
	if err != nil {
		return nil, fmt.Errorf("failed to build bundle: %w", err)
	}

	cwd, err := os.Getwd() // TODO1
	if err != nil {
		return nil, err
	}
	m := fromJSONToMeta(cwd, result.Metafile)

	groups := make(map[string]resource.Resources)

	// TODO1
	addFoo := func(filename, targetPath, group string, mt media.Type) error {
		rd := resources.ResourceSourceDescriptor{
			LazyPublish: true,
			OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
				return os.Open(filename)
			},
			MediaType:  mt,
			TargetPath: targetPath,
			// DependencyManager:  idm, TODO1
		}
		r, err := b.client.d.ResourceSpec.NewResource(rd)
		if err != nil {
			return err
		}

		groups[group] = append(groups[group], r)

		return nil
	}

	createAndAddResource3 := func(o esBuildResultMetaOutput) (bool, error) {
		p := filepath.ToSlash(strings.TrimPrefix(o.filename, outDir))
		ext := path.Ext(p)
		mt, _, found := mediaTypes.GetBySuffix(ext)
		if !found {
			return false, nil
		}
		groupPath := p
		group, found := state.pathGroup.Get(groupPath)

		if !found {
			return false, nil
		}

		if err := addFoo(o.filename, p, group, mt); err != nil {
			// TODO1
			return false, err
		}

		if o.CSSBundle != "" {
			p := filepath.ToSlash(strings.TrimPrefix(o.CSSBundle, outDir))
			if err := addFoo(o.CSSBundle, p, group, cssMt); err != nil {
				return false, err
			}
		}

		return true, nil
	}

	for _, o := range m.Outputs {
		handled, err := createAndAddResource3(o)
		if err != nil {
			return nil, err
		}
		if !handled {
			//  Copy to destination.
			p := strings.TrimPrefix(o.filename, outDir)
			if err := hugio.CopyFile(hugofs.Os, o.filename, filepath.Join(absPublishDir, p)); err != nil {
				return nil, fmt.Errorf("failed to copy %q to %q: %w", o.filename, absPublishDir, err)
			}
		}
	}

	return &Package{
		origin: b,
		outDir: outDir,
		b:      b,
		id:     path.Join(nsBundle, b.id),
		Groups: groups,
	}, nil
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

func fromJSONToMeta(cwd, s string) esBuildResultMeta {
	var m esBuildResultMeta
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	if err := m.Compile(cwd); err != nil {
		panic(err)
	}

	return m
}

// https://esbuild.github.io/api/#metafile
type esBuildResultMeta struct {
	Outputs map[string]esBuildResultMetaOutput

	// Compiled values.
	cssBundleEntryPoint map[string]esBuildResultMetaOutput
}

func (e *esBuildResultMeta) Compile(cwd string) error {
	// Rewrite the paths to be absolute.
	// See https://github.com/evanw/esbuild/issues/338
	outputs := make(map[string]esBuildResultMetaOutput)
	for k, v := range e.Outputs {
		filename := filepath.Join(cwd, k)
		if err := v.Compile(filename); err != nil {
			return err
		}
		deb("Output.Exports", k, v.Exports)
		if v.CSSBundle != "" {
			v.CSSBundle = filepath.Join(cwd, v.CSSBundle)
		}
		outputs[filename] = v
	}
	e.Outputs = outputs

	e.cssBundleEntryPoint = make(map[string]esBuildResultMetaOutput)
	for _, v := range e.Outputs {
		if v.CSSBundle != "" {
			e.cssBundleEntryPoint[v.CSSBundle] = v
		}
	}
	return nil
}

type esBuildResultMetaOutput struct {
	Bytes      int64
	Exports    []string
	Imports    []esBuildResultMetaOutputImport
	EntryPoint string
	CSSBundle  string

	// compiled values.
	filename string
}

func (e *esBuildResultMetaOutput) Compile(filename string) error {
	e.filename = filename
	return nil
}

type esBuildResultMetaOutputImport struct {
	Path string
	Kind string
}

// TODO1 remove the now superflous Close harness in the template package.
