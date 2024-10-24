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

// Package esbuild provides functions for building JavaScript resources.
package esbuild

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
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

var _ Batcher = (*batcher)(nil)

var _ resource.StaleInfo = (*options)(nil)

var _ identity.Identity = (*Package)(nil)

const (
	NsBatch = "__hugo-js-batch"
)

func (g *getOnce[T]) Get() T {
	var v T
	g.once.Do(func() {
		v = g.v
	})
	return v
}

func newOptions() *options {
	return &options{getOnce[*optionsSetter]{
		v: &optionsSetter{},
	}}
}

type Batcher interface {
	Build(context.Context) (*Package, error)
	Config() OptionsSetter
	Group(id string) BatcherGroup
}

//go:embed batch-esm-runner.gotmpl
var runnerTemplateStr string

func NewBatcherClient(deps *deps.Deps) (*BatcherClient, error) {
	return &BatcherClient{
		d:            deps,
		buildClient:  NewBuildClient(deps.BaseFs.Assets, deps.ResourceSpec),
		createClient: create.New(deps.ResourceSpec),
		bundlesCache: dynacache.GetOrCreatePartition[string, *Package](
			deps.MemCache,
			"/jsb1",
			// Mark it to clear on rebuild, but each package will evaluate itself for changes.
			dynacache.OptionsPartition{ClearWhen: dynacache.ClearOnRebuild, Weight: 10},
		),
	}, nil
}

type BatcherClient struct {
	d *deps.Deps

	once           sync.Once
	runnerTemplate tpl.Template

	createClient *create.Client
	buildClient  *BuildClient

	bundlesCache *dynacache.Partition[string, *Package]
}

// New creates a new Batcher with the given ID.
func (c *BatcherClient) New(id string) (Batcher, error) {
	var initErr error
	c.once.Do(func() {
		// We should fix the initialization order here (or use the Go template package directly), but we need to wait
		// for the Hugo templates to be ready.
		tmpl, err := c.d.TextTmpl().Parse("batch-esm-runner", runnerTemplateStr)
		if err != nil {
			initErr = err
			return
		}
		c.runnerTemplate = tmpl
	})

	if initErr != nil {
		return nil, initErr
	}

	return &batcher{
		id:                id,
		scriptGroups:      make(map[string]*scriptGroup),
		dependencyManager: c.d.Conf.NewIdentityManager("jsbatch" + id),
		client:            c,
		configOptions:     newOptions(),
	}, nil
}

func (c *BatcherClient) buildBatch(ctx context.Context, t *batchTemplateContext) (resource.Resource, string, error) {
	var buf bytes.Buffer

	if err := c.d.Tmpl().ExecuteWithContext(ctx, c.runnerTemplate, &buf, t); err != nil {
		return nil, "", err
	}

	s := paths.AddLeadingSlash(t.keyPath + ".js")
	r, err := c.createClient.FromString(s, buf.String())
	if err != nil {
		return nil, "", err
	}

	return r, s, nil
}

type BatcherGroup interface {
	Instance(sid, iid string) OptionsSetter
	Runner(id string) OptionsSetter
	Script(id string) OptionsSetter
}

type OptionsSetter interface {
	SetOptions(map[string]any) string
}

// TODO1 names.
type Package struct {
	origin       *batcher
	outDir       string
	id           string
	staleVersion uint32
	b            *batcher
	Groups       map[string]resource.Resources
}

func (b *Package) GetDependencyManager() identity.Manager {
	return b.origin.dependencyManager
}

func (p *Package) IdentifierBase() string {
	return p.id
}

func (p *Package) MarkStale() {
	p.origin.reset()
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

			for _, vv := range v.runnersOptions {
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

type ParamsOptions struct {
	Params json.RawMessage
}

type ScriptOptions struct {
	// The script to build.
	// TODO1 handle stale.
	Resource resource.Resource

	// The import context to use.
	// Note that we will always fall back to the resource's own import context.
	ImportContext resource.ResourceGetter

	// The export name to use for this script's group's runners (if any).
	// If not set, the default export will be used.
	Export string

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

type ScriptOptionsGetSetter interface {
	GetOptions() *ScriptOptions
	SetOptions(map[string]any) string
}

type batchTemplateContext struct {
	keyPath string
	ID      string
	Runners []scriptRunnerTemplateContext
	Scripts []scriptBatchTemplateContext
}

type batcher struct {
	mu           sync.Mutex
	id           string
	scriptGroups scriptGroups

	client            *BatcherClient
	dependencyManager identity.Manager

	configOptions *options

	// The last successfully built package.
	// If this is non-nil and not stale, we can reuse it (e.g. on server rebuilds)
	prevBuild *Package

	// Compiled.
	config ExternalOptions
}

func (b *batcher) Build(ctx context.Context) (*Package, error) {
	key := dynacache.CleanKey(b.id + ".js")
	p, err := b.client.bundlesCache.GetOrCreateWitTimeout(key, b.client.d.Conf.Timeout(), func(string) (*Package, error) {
		return b.build(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build Batch %q: %w", b.id, err)
	}

	if p.b != b {
		panic("bundler mismatch")
	}

	return p, nil
}

func (b *batcher) Config() OptionsSetter {
	return b.configOptions.Get()
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
			runnersOptions:   make(map[string]*options),
		}
		b.scriptGroups[id] = group
	}

	return group
}

func (b *batcher) build(ctx context.Context) (*Package, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Use the unexported calculateStaleVersion
	if b.prevBuild != nil {
		staleVersion := b.prevBuild.calculateStaleVersion()
		if staleVersion == 0 {
			return b.prevBuild, nil
		}
	}

	p, err := b.doBuild(ctx)
	if err != nil {
		return nil, err
	}
	b.prevBuild = p
	return p, nil
}

func (b *batcher) compile() error {
	var err error
	b.config, err = DecodeExternalOptions(b.configOptions.commit().opts)
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

func (b *batcher) doBuild(ctx context.Context) (*Package, error) {
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

	// Entry points passed to ESBuid.
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
		var runners []scriptRunnerTemplateContext
		for _, vv := range v.runners.Sorted() {
			runnerKeyPath := keyPath + "_" + vv.ID
			runnerImpPath := paths.AddLeadingSlash(runnerKeyPath + "_runner" + vv.Resource.MediaType().FirstSuffix.FullSuffix)
			runners = append(runners, scriptRunnerTemplateContext{script: vv, Import: runnerImpPath})
			addResource(k, runnerImpPath, vv.Resource, false)
		}

		t := &batchTemplateContext{
			keyPath: keyPath,
			ID:      v.id,
			Runners: runners,
		}

		instances := v.instances.Sorted()

		for _, vv := range v.scripts.Sorted() {
			keyPath := keyPath + "_" + vv.ID
			opts := vv.ScriptOptions
			impPath := path.Join(PrefixHugoVirtual, opts.Dir(), keyPath+opts.Resource.MediaType().FirstSuffix.FullSuffix)
			impCtx := opts.ImportContext

			state.importerImportContext.Set(impPath, importContext{
				name:           keyPath,
				resourceGetter: impCtx,
				scriptOptions:  opts,
			})

			bt := scriptBatchTemplateContext{
				script: vv,
				Import: impPath,
			}

			state.importResource.Set(bt.Import, vv.Resource)
			for _, vvv := range instances.ByScriptID(vv.ID) {
				bt.Instances = append(bt.Instances, scriptInstanceBatchTemplateContext{instance: vvv})
			}

			t.Scripts = append(t.Scripts, bt)
		}

		r, s, err := b.client.buildBatch(ctx, t)
		if err != nil {
			return nil, fmt.Errorf("failed to build batch: %w", err)
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

	outDir, err := b.client.d.MkdirTemp("hugo-jsbatch")
	if err != nil {
		return nil, err
	}

	externalOptions := b.config
	if externalOptions.Format == "" {
		externalOptions.Format = "esm"
	}
	if externalOptions.Format != "esm" {
		return nil, fmt.Errorf("only esm format is currently supported")
	}

	jsOpts := Options{
		ExternalOptions: externalOptions,
		InternalOptions: InternalOptions{
			DependencyManager: b.dependencyManager,
			OutDir:            outDir,
			Write:             true,
			AllowOverwrite:    true,
			Splitting:         true,
			ImportOnResolveFunc: func(depsManager identity.Manager, imp string, args api.OnResolveArgs) string {
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
						depsManager.AddIdentity(identity.FirstIdentity(resolved))
						imp := PrefixHugoVirtual + resolved.(resource.PathProvider).Path()
						state.importResource.Set(imp, resolved)
						state.importerImportContext.Set(imp, importContext)
						return imp

					}
				}
				return ""
			},
			ImportOnLoadFunc: func(args api.OnLoadArgs) string {
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
			ErrorMessageResolveFunc: func(args api.Message) *ErrorMessageResolved {
				if loc := args.Location; loc != nil {
					path := strings.TrimPrefix(loc.File, NsHugoImportResolveFunc+":")
					if r, found := state.importResource.Get(path); found {
						path = strings.TrimPrefix(path, PrefixHugoVirtual)
						var contentr hugio.ReadSeekCloser
						if cp, ok := r.(hugio.ReadSeekCloserProvider); ok {
							contentr, _ = cp.ReadSeekCloser()
						}
						return &ErrorMessageResolved{
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

	result, err := b.client.buildClient.Build(jsOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build bundle: %w", err)
	}

	m := fromJSONToESBuildResultMeta(b.client.d.Conf.WorkingDir(), result.Metafile)

	groups := make(map[string]resource.Resources)

	createAndAddResource := func(filename, targetPath, group string, mt media.Type) error {
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

	createAndAddResources := func(o esBuildResultMetaOutput) (bool, error) {
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

		if err := createAndAddResource(o.filename, p, group, mt); err != nil {
			return false, err
		}

		if o.CSSBundle != "" {
			p := filepath.ToSlash(strings.TrimPrefix(o.CSSBundle, outDir))
			if err := createAndAddResource(o.CSSBundle, p, group, cssMt); err != nil {
				return false, err
			}
		}

		return true, nil
	}

	for _, o := range m.Outputs {
		handled, err := createAndAddResources(o)
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
		id:     path.Join(NsBatch, b.id),
		Groups: groups,
	}, nil
}

func (b *batcher) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, v := range b.scriptGroups {
		// TODO1 check if this is complete.
		v.Reset()
	}
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

type getOnce[T any] struct {
	v    T
	once lazy.OnceMore
}

type instance struct {
	instanceID
	*ParamsOptions
}

type instanceID struct {
	scriptID   string
	instanceID string
}

type (
	instanceMap map[instanceID]*ParamsOptions
	instances   []*instance
)

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

func (i instances) ByScriptID(id string) instances {
	var a instances
	for _, v := range i {
		if v.instanceID.scriptID == id {
			a = append(a, v)
		}
	}
	return a
}

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

type script struct {
	ID string
	*ScriptOptions
}

type scriptBatchTemplateContext struct {
	*script
	Import    string
	Instances []scriptInstanceBatchTemplateContext
}

func (c scriptBatchTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID        string                               `json:"id"`
		Instances []scriptInstanceBatchTemplateContext `json:"instances"`
	}{
		ID:        c.ID,
		Instances: c.Instances,
	})
}

func (b scriptBatchTemplateContext) RunnerJSON(i int) string {
	script := fmt.Sprintf("Script%d", i)

	v := struct {
		ID string `json:"id"`

		// Read-only live JavaScript binding.
		Binding   string                               `json:"binding"`
		Instances []scriptInstanceBatchTemplateContext `json:"instances"`
	}{
		b.ID,
		script,
		b.Instances,
	}

	bb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s := string(bb)

	// Remove the quotes to make it a valid JS object.
	s = strings.ReplaceAll(s, fmt.Sprintf("%q", script), script)

	return s
}

type scriptGroup struct {
	mu sync.Mutex

	id string

	client *BatcherClient

	scriptsOptions   map[string]*options
	instancesOptions map[instanceID]*options
	runnersOptions   map[string]*options

	// Compiled.
	scripts   scriptMap
	instances instanceMap
	runners   scriptMap
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

func (g *scriptGroup) Reset() {
	for _, v := range g.scriptsOptions {
		v.Reset()
	}
	for _, v := range g.instancesOptions {
		v.Reset()
	}
	for _, v := range g.runnersOptions {
		v.Reset()
	}
}

func (s *scriptGroup) Runner(id string) OptionsSetter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, found := s.runnersOptions[id]; found {
		return v.Get()
	}
	s.runnersOptions[id] = newOptions()
	return s.runnersOptions[id].Get()
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

func (s *scriptGroup) errFailedToCompile(what, id string, err error) error {
	return fmt.Errorf("failed to compile %s for %q > %q: %w", what, s.id, id, err)
}

func (s *scriptGroup) compile() error {
	// TODO1 lock?
	s.scripts = make(map[string]*ScriptOptions)
	s.instances = make(map[instanceID]*ParamsOptions)
	s.runners = make(map[string]*ScriptOptions)

	for k, v := range s.scriptsOptions {
		compiled, err := compileScriptOptions(v, "*")
		if err != nil {
			return s.errFailedToCompile("Script", k, err)
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

	for k, v := range s.runnersOptions {
		compiled, err := compileScriptOptions(v, "default")
		if err != nil {
			return s.errFailedToCompile("Runner", k, err)
		}
		s.runners[k] = compiled
	}

	return nil
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

type scriptInstanceBatchTemplateContext struct {
	*instance
}

func (c scriptInstanceBatchTemplateContext) ID() string {
	return c.instanceID.instanceID
}

func (c scriptInstanceBatchTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID     string          `json:"id"`
		Params json.RawMessage `json:"params"`
	}{
		ID:     c.instanceID.instanceID,
		Params: c.Params,
	})
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

type scriptRunnerTemplateContext struct {
	*script
	Import string
}

func (c scriptRunnerTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID string `json:"id"`
	}{
		ID: c.ID,
	})
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

func compileScriptOptions(o *options, defaultExport string) (*ScriptOptions, error) {
	v := struct {
		Resource      resource.Resource
		ImportContext any
		Export        string
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

	if v.Export == "" {
		v.Export = defaultExport
	}

	compiled := &ScriptOptions{
		Resource:      v.Resource,
		Export:        v.Export,
		ImportContext: resource.NewResourceGetter(v.ImportContext),
		Params:        paramsJSON,
	}

	if compiled.Resource == nil {
		return nil, fmt.Errorf("resource not set")
	}

	return compiled, nil
}

// TODO1 remove.
func deb(what string, v ...any) {
	fmt.Println(what, v)
}

func fromJSONToESBuildResultMeta(cwd, jsons string) esBuildResultMeta {
	var m esBuildResultMeta
	if err := json.Unmarshal([]byte(jsons), &m); err != nil {
		panic(err)
	}
	if err := m.Compile(cwd); err != nil {
		panic(err)
	}
	return m
}

func logTime(name string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("%s in %s\n", name, elapsed)
}
