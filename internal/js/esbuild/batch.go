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
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

var _ Batcher = (*batcher)(nil)

const (
	NsBatch = "_hugo-js-batch"

	propsKeyImportContext = "importContext"
	propsResoure          = "resource"
)

//go:embed batch-esm-runner.gotmpl
var runnerTemplateStr string

var _ BatchPackage = (*Package)(nil)

var _ buildToucher = (*optsHolder[scriptOptions])(nil)

var (
	_ buildToucher             = (*scriptGroup)(nil)
	_ isBuiltOrTouchedProvider = (*scriptGroup)(nil)
)

func NewBatcherClient(deps *deps.Deps) (*BatcherClient, error) {
	c := &BatcherClient{
		d:            deps,
		buildClient:  NewBuildClient(deps.BaseFs.Assets, deps.ResourceSpec),
		createClient: create.New(deps.ResourceSpec),
		bundlesCache: maps.NewCache[string, BatchPackage](),
	}

	deps.BuildEndListeners.Add(func(...any) bool {
		c.bundlesCache.Reset()
		return false
	})

	return c, nil
}

func (o optionsMap[K, C]) ByKey() optionsGetSetters[K, C] {
	var values []optionsGetSetter[K, C]
	for _, v := range o {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Key().String() < values[j].Key().String()
	})

	return values
}

func (o *opts[K, C]) Compiled() C {
	o.h.checkCompileErr()
	return o.h.compiled
}

func (os optionsGetSetters[K, C]) Filter(predicate func(K) bool) optionsGetSetters[K, C] {
	var a optionsGetSetters[K, C]
	for _, v := range os {
		if predicate(v.Key()) {
			a = append(a, v)
		}
	}
	return a
}

func (o *optsHolder[C]) IdentifierBase() string {
	return o.optionsID
}

func (o *opts[K, C]) Key() K {
	return o.key
}

func (o *opts[K, C]) Reset() {
	mu := o.once.ResetWithLock()
	defer mu.Unlock()
	o.h.resetCounter++
}

func (o *opts[K, C]) Get(id uint32) OptionsSetter {
	var b *optsHolder[C]
	o.once.Do(func() {
		b = o.h
		b.setBuilt(id)
	})
	return b
}

func (o *opts[K, C]) GetIdentity() identity.Identity {
	return o.h
}

func (o *optsHolder[C]) SetOptions(m map[string]any) string {
	o.optsSetCounter++
	o.optsPrev = o.optsCurr
	o.optsCurr = m
	o.compiledPrev = o.compiled
	o.compiled, o.compileErr = o.compiled.compileOptions(m, o.defaults)
	o.checkCompileErr()
	return ""
}

// ValidateBatchID validates the given ID according to some very
func ValidateBatchID(id string, isTopLevel bool) error {
	if id == "" {
		return fmt.Errorf("id must be set")
	}
	// No Windows slashes.
	if strings.Contains(id, "\\") {
		return fmt.Errorf("id must not contain backslashes")
	}

	// Allow forward slashes in top level IDs only.
	if !isTopLevel && strings.Contains(id, "/") {
		return fmt.Errorf("id must not contain forward slashes")
	}

	return nil
}

func newIsBuiltOrTouched() isBuiltOrTouched {
	return isBuiltOrTouched{
		built:   make(buildIDs),
		touched: make(buildIDs),
	}
}

func newOpts[K any, C optionsCompiler[C]](key K, optionsID string, defaults defaultOptionValues) *opts[K, C] {
	return &opts[K, C]{
		key: key,
		h: &optsHolder[C]{
			optionsID:        optionsID,
			defaults:         defaults,
			isBuiltOrTouched: newIsBuiltOrTouched(),
		},
	}
}

// BatchPackage holds a group of JavaScript resources.
type BatchPackage interface {
	Groups() map[string]resource.Resources
}

// Batcher is used to build JavaScript packages.
type Batcher interface {
	Build(context.Context) (BatchPackage, error)
	Config(ctx context.Context) OptionsSetter
	Group(ctx context.Context, id string) BatcherGroup
}

// BatcherClient is a client for building JavaScript packages.
type BatcherClient struct {
	d *deps.Deps

	once           sync.Once
	runnerTemplate tpl.Template

	createClient *create.Client
	buildClient  *BuildClient

	bundlesCache *maps.Cache[string, BatchPackage]
}

// New creates a new Batcher with the given ID.
// This will be typically created once and reused across rebuilds.
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

	dependencyManager := c.d.Conf.NewIdentityManager("jsbatch_" + id)
	configID := "config_" + id

	b := &batcher{
		id:                id,
		scriptGroups:      make(map[string]*scriptGroup),
		dependencyManager: dependencyManager,
		client:            c,
		configOptions: newOpts[scriptID, configOptions](
			scriptID(configID),
			configID,
			defaultOptionValues{},
		),
	}

	c.d.BuildEndListeners.Add(func(...any) bool {
		b.reset()
		return false
	})

	idFinder := identity.NewFinder(identity.FinderConfig{})

	c.d.OnChangeListeners.Add(func(ids ...identity.Identity) bool {
		for _, id := range ids {
			if r := idFinder.Contains(id, b.dependencyManager, 50); r > 0 {
				b.staleVersion.Add(1)
				return false
			}

			sp, ok := id.(identity.DependencyManagerScopedProvider)
			if !ok {
				continue
			}
			idms := sp.GetDependencyManagerForScopesAll()

			for _, g := range b.scriptGroups {
				g.forEachIdentity(func(id2 identity.Identity) bool {
					bt, ok := id2.(buildToucher)
					if !ok {
						return false
					}
					for _, id3 := range idms {
						// This handles the removal of the only source for a script group (e.g. all shortcodes in a contnt page).
						// Note the very shallow search.
						if r := idFinder.Contains(id2, id3, 0); r > 0 {
							bt.setTouched(b.buildCount)
							return false
						}
					}
					return false
				})
			}
		}

		return false
	})

	return b, nil
}

func (c *BatcherClient) buildBatchGroup(ctx context.Context, t *batchGroupTemplateContext) (resource.Resource, string, error) {
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

// BatcherGroup is a group of scripts and instances.
type BatcherGroup interface {
	Instance(sid, iid string) OptionsSetter
	Runner(id string) OptionsSetter
	Script(id string) OptionsSetter
}

// OptionsSetter is used to set options for a batch, script or instance.
type OptionsSetter interface {
	SetOptions(map[string]any) string
}

// Package holds a group of JavaScript resources.
type Package struct {
	id string
	b  *batcher

	groups map[string]resource.Resources
}

func (p *Package) Groups() map[string]resource.Resources {
	return p.groups
}

type batchGroupTemplateContext struct {
	keyPath string
	ID      string
	Runners []scriptRunnerTemplateContext
	Scripts []scriptBatchTemplateContext
}

type batcher struct {
	mu           sync.Mutex
	id           string
	buildCount   uint32
	staleVersion atomic.Uint32
	scriptGroups scriptGroups

	client            *BatcherClient
	dependencyManager identity.Manager

	configOptions optionsGetSetter[scriptID, configOptions]

	// The last successfully built package.
	// If this is non-nil and not stale, we can reuse it (e.g. on server rebuilds)
	prevBuild *Package
}

// Build builds the batch if not already built or if it's stale.
func (b *batcher) Build(ctx context.Context) (BatchPackage, error) {
	key := dynacache.CleanKey(b.id + ".js")
	p, err := b.client.bundlesCache.GetOrCreate(key, func() (BatchPackage, error) {
		return b.build(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build JS batch %q: %w", b.id, err)
	}
	return p, nil
}

func (b *batcher) Config(ctx context.Context) OptionsSetter {
	return b.configOptions.Get(b.buildCount)
}

func (b *batcher) Group(ctx context.Context, id string) BatcherGroup {
	if err := ValidateBatchID(id, false); err != nil {
		panic(err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	group, found := b.scriptGroups[id]
	if !found {
		idm := b.client.d.Conf.NewIdentityManager("jsbatch_" + id)
		b.dependencyManager.AddIdentity(idm)

		group = &scriptGroup{
			id: id, b: b,
			isBuiltOrTouched:  newIsBuiltOrTouched(),
			dependencyManager: idm,
			scriptsOptions:    make(optionsMap[scriptID, scriptOptions]),
			instancesOptions:  make(optionsMap[instanceID, paramsOptions]),
			runnersOptions:    make(optionsMap[scriptID, scriptOptions]),
		}
		b.scriptGroups[id] = group
	}

	group.setBuilt(b.buildCount)

	return group
}

func (b *batcher) isStale() bool {
	if b.staleVersion.Load() > 0 {
		return true
	}

	if b.removeNotSet() {
		return true
	}

	if b.configOptions.isStale() {
		return true
	}

	for _, v := range b.scriptGroups {
		if v.isStale() {
			return true
		}
	}

	return false
}

func (b *batcher) build(ctx context.Context) (BatchPackage, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	defer func() {
		b.staleVersion.Store(0)
		b.buildCount++
	}()

	if b.prevBuild != nil {
		if !b.isStale() {
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

func (b *batcher) doBuild(ctx context.Context) (*Package, error) {
	type importContext struct {
		name           string
		resourceGetter resource.ResourceGetter
		scriptOptions  scriptOptions
		dm             identity.Manager
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
		state.pathGroup.Set(paths.TrimExt(pth), group)
		state.importResource.Set(pth, r)
		if isResult {
			state.resultResource.Set(pth, r)
		}
		entryPoints = append(entryPoints, pth)
	}

	for _, g := range b.scriptGroups.Sorted() {
		keyPath := g.id
		var runners []scriptRunnerTemplateContext
		for _, vv := range g.runnersOptions.ByKey() {
			runnerKeyPath := keyPath + "_" + vv.Key().String()
			runnerImpPath := paths.AddLeadingSlash(runnerKeyPath + "_runner" + vv.Compiled().Resource.MediaType().FirstSuffix.FullSuffix)
			runners = append(runners, scriptRunnerTemplateContext{opts: vv, Import: runnerImpPath})
			addResource(g.id, runnerImpPath, vv.Compiled().Resource, false)
		}

		t := &batchGroupTemplateContext{
			keyPath: keyPath,
			ID:      g.id,
			Runners: runners,
		}

		instances := g.instancesOptions.ByKey()

		for _, vv := range g.scriptsOptions.ByKey() {
			keyPath := keyPath + "_" + vv.Key().String()
			opts := vv.Compiled()
			impPath := path.Join(PrefixHugoVirtual, opts.Dir(), keyPath+opts.Resource.MediaType().FirstSuffix.FullSuffix)
			impCtx := opts.ImportContext

			state.importerImportContext.Set(impPath, importContext{
				name:           keyPath,
				resourceGetter: impCtx,
				scriptOptions:  opts,
				dm:             g.dependencyManager,
			})

			bt := scriptBatchTemplateContext{
				opts:   vv,
				Import: impPath,
			}
			state.importResource.Set(bt.Import, vv.Compiled().Resource)
			predicate := func(k instanceID) bool {
				return k.scriptID == vv.Key()
			}
			for _, vvv := range instances.Filter(predicate) {
				bt.Instances = append(bt.Instances, scriptInstanceBatchTemplateContext{opts: vvv})
			}

			t.Scripts = append(t.Scripts, bt)
		}

		r, s, err := b.client.buildBatchGroup(ctx, t)
		if err != nil {
			return nil, fmt.Errorf("failed to build JS batch: %w", err)
		}

		state.importerImportContext.Set(s, importContext{
			name:           s,
			resourceGetter: nil,
			dm:             g.dependencyManager,
		})

		addResource(g.id, s, r, true)
	}

	mediaTypes := b.client.d.ResourceSpec.MediaTypes()

	externalOptions := b.configOptions.Compiled().Options
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
			Splitting:         true,
			ImportOnResolveFunc: func(imp string, args api.OnResolveArgs) string {
				var importContextPath string
				if args.Kind == api.ResolveEntryPoint {
					importContextPath = args.Path
				} else {
					importContextPath = args.Importer
				}
				importContext, importContextFound := state.importerImportContext.Get(importContextPath)

				// We want to track the dependencies closest to where they're used.
				dm := b.dependencyManager
				if importContextFound {
					dm = importContext.dm
				}

				if r, found := state.importResource.Get(imp); found {
					dm.AddIdentity(identity.FirstIdentity(r))
					return imp
				}

				if importContext.resourceGetter != nil {
					resolved := ResolveResource(imp, importContext.resourceGetter)
					if resolved != nil {
						resolvePath := resources.InternalResourceTargetPath(resolved)
						dm.AddIdentity(identity.FirstIdentity(resolved))
						imp := PrefixHugoVirtual + resolvePath
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
					content, err := r.(resource.ContentProvider).Content(ctx)
					if err != nil {
						panic(err)
					}
					return cast.ToString(content)
				}
				return ""
			},
			ImportParamsOnLoadFunc: func(args api.OnLoadArgs) json.RawMessage {
				if importContext, found := state.importerImportContext.Get(args.Path); found {
					if !importContext.scriptOptions.IsZero() {
						return importContext.scriptOptions.Params
					}
				}
				return nil
			},
			ErrorMessageResolveFunc: func(args api.Message) *ErrorMessageResolved {
				if loc := args.Location; loc != nil {
					path := strings.TrimPrefix(loc.File, NsHugoImportResolveFunc+":")
					if r, found := state.importResource.Get(path); found {
						sourcePath := resources.InternalResourceSourcePathBestEffort(r)

						var contentr hugio.ReadSeekCloser
						if cp, ok := r.(hugio.ReadSeekCloserProvider); ok {
							contentr, _ = cp.ReadSeekCloser()
						}
						return &ErrorMessageResolved{
							Content: contentr,
							Path:    sourcePath,
							Message: args.Text,
						}

					}

				}
				return nil
			},
			ResolveSourceMapSource: func(s string) string {
				if r, found := state.importResource.Get(s); found {
					if ss := resources.InternalResourceSourcePath(r); ss != "" {
						return ss
					}
					return PrefixHugoMemory + s
				}
				return ""
			},
			EntryPoints: entryPoints,
		},
	}

	result, err := b.client.buildClient.Build(jsOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build JS bundle: %w", err)
	}

	groups := make(map[string]resource.Resources)

	createAndAddResource := func(targetPath, group string, o api.OutputFile, mt media.Type) error {
		var sourceFilename string
		if r, found := state.importResource.Get(targetPath); found {
			sourceFilename = resources.InternalResourceSourcePathBestEffort(r)
		}
		targetPath = path.Join(b.id, targetPath)

		rd := resources.ResourceSourceDescriptor{
			LazyPublish: true,
			OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
				return hugio.NewReadSeekerNoOpCloserFromBytes(o.Contents), nil
			},
			MediaType:            mt,
			TargetPath:           targetPath,
			SourceFilenameOrPath: sourceFilename,
		}
		r, err := b.client.d.ResourceSpec.NewResource(rd)
		if err != nil {
			return err
		}

		groups[group] = append(groups[group], r)

		return nil
	}

	outDir := b.client.d.AbsPublishDir

	createAndAddResources := func(o api.OutputFile) (bool, error) {
		p := paths.ToSlashPreserveLeading(strings.TrimPrefix(o.Path, outDir))
		ext := path.Ext(p)
		mt, _, found := mediaTypes.GetBySuffix(ext)
		if !found {
			return false, nil
		}

		group, found := state.pathGroup.Get(paths.TrimExt(p))

		if !found {
			return false, nil
		}

		if err := createAndAddResource(p, group, o, mt); err != nil {
			return false, err
		}

		return true, nil
	}

	for _, o := range result.OutputFiles {
		handled, err := createAndAddResources(o)
		if err != nil {
			return nil, err
		}

		if !handled {
			//  Copy to destination.
			p := strings.TrimPrefix(o.Path, outDir)
			targetFilename := filepath.Join(b.id, p)
			fs := b.client.d.BaseFs.PublishFs
			if err := fs.MkdirAll(filepath.Dir(targetFilename), 0o777); err != nil {
				return nil, fmt.Errorf("failed to create dir %q: %w", targetFilename, err)
			}

			if err := afero.WriteFile(fs, targetFilename, o.Contents, 0o666); err != nil {
				return nil, fmt.Errorf("failed to write to %q: %w", targetFilename, err)
			}
		}
	}

	p := &Package{
		id:     path.Join(NsBatch, b.id),
		b:      b,
		groups: groups,
	}

	return p, nil
}

func (b *batcher) removeNotSet() bool {
	// We already have the lock.
	var removed bool
	currentBuildID := b.buildCount
	for k, v := range b.scriptGroups {
		if !v.isBuilt(currentBuildID) && v.isTouched(currentBuildID) {
			// Remove entire group.
			removed = true
			delete(b.scriptGroups, k)
			continue
		}
		if v.removeTouchedButNotSet() {
			removed = true
		}
		if v.removeNotSet() {
			removed = true
		}
	}

	return removed
}

func (b *batcher) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.configOptions.Reset()
	for _, v := range b.scriptGroups {
		v.Reset()
	}
}

type buildIDs map[uint32]bool

func (b buildIDs) Has(buildID uint32) bool {
	return b[buildID]
}

func (b buildIDs) Set(buildID uint32) {
	b[buildID] = true
}

type buildToucher interface {
	setTouched(buildID uint32)
}

type configOptions struct {
	Options ExternalOptions
}

func (s configOptions) isStaleCompiled(prev configOptions) bool {
	return false
}

func (s configOptions) compileOptions(m map[string]any, defaults defaultOptionValues) (configOptions, error) {
	config, err := DecodeExternalOptions(m)
	if err != nil {
		return configOptions{}, err
	}

	return configOptions{
		Options: config,
	}, nil
}

type defaultOptionValues struct {
	defaultExport string
}

type instanceID struct {
	scriptID   scriptID
	instanceID string
}

func (i instanceID) String() string {
	return i.scriptID.String() + "_" + i.instanceID
}

type isBuiltOrTouched struct {
	built   buildIDs
	touched buildIDs
}

func (i isBuiltOrTouched) setBuilt(id uint32) {
	i.built.Set(id)
}

func (i isBuiltOrTouched) isBuilt(id uint32) bool {
	return i.built.Has(id)
}

func (i isBuiltOrTouched) setTouched(id uint32) {
	i.touched.Set(id)
}

func (i isBuiltOrTouched) isTouched(id uint32) bool {
	return i.touched.Has(id)
}

type isBuiltOrTouchedProvider interface {
	isBuilt(uint32) bool
	isTouched(uint32) bool
}

type key interface {
	comparable
	fmt.Stringer
}

type optionsCompiler[C any] interface {
	isStaleCompiled(C) bool
	compileOptions(map[string]any, defaultOptionValues) (C, error)
}

type optionsGetSetter[K, C any] interface {
	isBuiltOrTouchedProvider
	identity.IdentityProvider
	// resource.StaleInfo

	Compiled() C
	Key() K
	Reset()

	Get(uint32) OptionsSetter
	isStale() bool
	currPrev() (map[string]any, map[string]any)
}

type optionsGetSetters[K key, C any] []optionsGetSetter[K, C]

type optionsMap[K key, C any] map[K]optionsGetSetter[K, C]

type opts[K any, C optionsCompiler[C]] struct {
	key  K
	h    *optsHolder[C]
	once lazy.OnceMore
}

type optsHolder[C optionsCompiler[C]] struct {
	optionsID string

	defaults defaultOptionValues

	// Keep track of one generation so we can detect changes.
	// Note that most of this tracking is performed on the options/map level.
	compiled     C
	compiledPrev C
	compileErr   error

	resetCounter   uint32
	optsSetCounter uint32
	optsCurr       map[string]any
	optsPrev       map[string]any

	isBuiltOrTouched
}

type paramsOptions struct {
	Params json.RawMessage
}

func (s paramsOptions) isStaleCompiled(prev paramsOptions) bool {
	return false
}

func (s paramsOptions) compileOptions(m map[string]any, defaults defaultOptionValues) (paramsOptions, error) {
	v := struct {
		Params map[string]any
	}{}

	if err := mapstructure.WeakDecode(m, &v); err != nil {
		return paramsOptions{}, err
	}

	paramsJSON, err := json.Marshal(v.Params)
	if err != nil {
		return paramsOptions{}, err
	}

	return paramsOptions{
		Params: paramsJSON,
	}, nil
}

type scriptBatchTemplateContext struct {
	opts      optionsGetSetter[scriptID, scriptOptions]
	Import    string
	Instances []scriptInstanceBatchTemplateContext
}

func (s *scriptBatchTemplateContext) Export() string {
	return s.opts.Compiled().Export
}

func (c scriptBatchTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID        string                               `json:"id"`
		Instances []scriptInstanceBatchTemplateContext `json:"instances"`
	}{
		ID:        c.opts.Key().String(),
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
		b.opts.Key().String(),
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
	b  *batcher
	isBuiltOrTouched
	dependencyManager identity.Manager

	scriptsOptions   optionsMap[scriptID, scriptOptions]
	instancesOptions optionsMap[instanceID, paramsOptions]
	runnersOptions   optionsMap[scriptID, scriptOptions]
}

// For internal use only.
func (b *scriptGroup) GetDependencyManager() identity.Manager {
	return b.dependencyManager
}

// For internal use only.
func (b *scriptGroup) IdentifierBase() string {
	return b.id
}

func (s *scriptGroup) Instance(sid, id string) OptionsSetter {
	if err := ValidateBatchID(sid, false); err != nil {
		panic(err)
	}
	if err := ValidateBatchID(id, false); err != nil {
		panic(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	iid := instanceID{scriptID: scriptID(sid), instanceID: id}
	if v, found := s.instancesOptions[iid]; found {
		return v.Get(s.b.buildCount)
	}

	fullID := "instance_" + s.key() + "_" + iid.String()

	s.instancesOptions[iid] = newOpts[instanceID, paramsOptions](
		iid,
		fullID,
		defaultOptionValues{},
	)

	return s.instancesOptions[iid].Get(s.b.buildCount)
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
	if err := ValidateBatchID(id, false); err != nil {
		panic(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	sid := scriptID(id)
	if v, found := s.runnersOptions[sid]; found {
		return v.Get(s.b.buildCount)
	}

	runnerIdentity := "runner_" + s.key() + "_" + id

	// A typical signature for a runner would be:
	//     export default function Run(scripts) {}
	// The user can override the default export in the templates.

	s.runnersOptions[sid] = newOpts[scriptID, scriptOptions](
		sid,
		runnerIdentity,
		defaultOptionValues{
			defaultExport: "default",
		},
	)

	return s.runnersOptions[sid].Get(s.b.buildCount)
}

func (s *scriptGroup) Script(id string) OptionsSetter {
	if err := ValidateBatchID(id, false); err != nil {
		panic(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	sid := scriptID(id)
	if v, found := s.scriptsOptions[sid]; found {
		return v.Get(s.b.buildCount)
	}

	scriptIdentity := "script_" + s.key() + "_" + id

	s.scriptsOptions[sid] = newOpts[scriptID, scriptOptions](
		sid,
		scriptIdentity,
		defaultOptionValues{
			defaultExport: "*",
		},
	)

	return s.scriptsOptions[sid].Get(s.b.buildCount)
}

func (s *scriptGroup) isStale() bool {
	for _, v := range s.scriptsOptions {
		if v.isStale() {
			return true
		}
	}

	for _, v := range s.instancesOptions {
		if v.isStale() {
			return true
		}
	}

	for _, v := range s.runnersOptions {
		if v.isStale() {
			return true
		}
	}

	return false
}

func (v *scriptGroup) forEachIdentity(
	f func(id identity.Identity) bool,
) bool {
	if f(v) {
		return true
	}
	for _, vv := range v.instancesOptions {
		if f(vv.GetIdentity()) {
			return true
		}
	}

	for _, vv := range v.scriptsOptions {
		if f(vv.GetIdentity()) {
			return true
		}
	}

	for _, vv := range v.runnersOptions {
		if f(vv.GetIdentity()) {
			return true
		}
	}

	return false
}

func (s *scriptGroup) key() string {
	return s.b.id + "_" + s.id
}

func (g *scriptGroup) removeNotSet() bool {
	currentBuildID := g.b.buildCount
	if !g.isBuilt(currentBuildID) {
		// This group was never accessed in this build.
		return false
	}
	var removed bool

	if g.instancesOptions.isBuilt(currentBuildID) {
		// A new instance has been set in this group for this build.
		// Remove any instance that has not been set in this build.
		for k, v := range g.instancesOptions {
			if v.isBuilt(currentBuildID) {
				continue
			}
			delete(g.instancesOptions, k)
			removed = true
		}
	}

	if g.runnersOptions.isBuilt(currentBuildID) {
		// A new runner has been set in this group for this build.
		// Remove any runner that has not been set in this build.
		for k, v := range g.runnersOptions {
			if v.isBuilt(currentBuildID) {
				continue
			}
			delete(g.runnersOptions, k)
			removed = true
		}
	}

	if g.scriptsOptions.isBuilt(currentBuildID) {
		// A new script has been set in this group for this build.
		// Remove any script that has not been set in this build.
		for k, v := range g.scriptsOptions {
			if v.isBuilt(currentBuildID) {
				continue
			}
			delete(g.scriptsOptions, k)

			// Also remove any instance with this ID.
			for kk := range g.instancesOptions {
				if kk.scriptID == k {
					delete(g.instancesOptions, kk)
				}
			}
			removed = true
		}
	}

	return removed
}

func (g *scriptGroup) removeTouchedButNotSet() bool {
	currentBuildID := g.b.buildCount
	var removed bool
	for k, v := range g.instancesOptions {
		if v.isBuilt(currentBuildID) {
			continue
		}
		if v.isTouched(currentBuildID) {
			delete(g.instancesOptions, k)
			removed = true
		}
	}
	for k, v := range g.runnersOptions {
		if v.isBuilt(currentBuildID) {
			continue
		}
		if v.isTouched(currentBuildID) {
			delete(g.runnersOptions, k)
			removed = true
		}
	}
	for k, v := range g.scriptsOptions {
		if v.isBuilt(currentBuildID) {
			continue
		}
		if v.isTouched(currentBuildID) {
			delete(g.scriptsOptions, k)
			removed = true

			// Also remove any instance with this ID.
			for kk := range g.instancesOptions {
				if kk.scriptID == k {
					delete(g.instancesOptions, kk)
				}
			}
		}

	}
	return removed
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

type scriptID string

func (s scriptID) String() string {
	return string(s)
}

type scriptInstanceBatchTemplateContext struct {
	opts optionsGetSetter[instanceID, paramsOptions]
}

func (c scriptInstanceBatchTemplateContext) ID() string {
	return c.opts.Key().instanceID
}

func (c scriptInstanceBatchTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID     string          `json:"id"`
		Params json.RawMessage `json:"params"`
	}{
		ID:     c.opts.Key().instanceID,
		Params: c.opts.Compiled().Params,
	})
}

type scriptOptions struct {
	// The script to build.
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

func (o *scriptOptions) Dir() string {
	return path.Dir(resources.InternalResourceTargetPath(o.Resource))
}

func (s scriptOptions) IsZero() bool {
	return s.Resource == nil
}

func (s scriptOptions) isStaleCompiled(prev scriptOptions) bool {
	if prev.IsZero() {
		return false
	}

	// All but the ImportContext are checked at the options/map level.
	i1nil, i2nil := prev.ImportContext == nil, s.ImportContext == nil
	if i1nil && i2nil {
		return false
	}
	if i1nil || i2nil {
		return true
	}
	// On its own this check would have too many false positives, but combined with the other checks it should be fine.
	// We cannot do equality checking here.
	if !prev.ImportContext.(resource.IsProbablySameResourceGetter).IsProbablySameResourceGetter(s.ImportContext) {
		return true
	}

	return false
}

func (s scriptOptions) compileOptions(m map[string]any, defaults defaultOptionValues) (scriptOptions, error) {
	v := struct {
		Resource      resource.Resource
		ImportContext any
		Export        string
		Params        map[string]any
	}{}

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
		v.Export = defaults.defaultExport
	}

	compiled := scriptOptions{
		Resource:      v.Resource,
		Export:        v.Export,
		ImportContext: resource.NewCachedResourceGetter(v.ImportContext),
		Params:        paramsJSON,
	}

	if compiled.Resource == nil {
		return scriptOptions{}, fmt.Errorf("resource not set")
	}

	return compiled, nil
}

type scriptRunnerTemplateContext struct {
	opts   optionsGetSetter[scriptID, scriptOptions]
	Import string
}

func (s *scriptRunnerTemplateContext) Export() string {
	return s.opts.Compiled().Export
}

func (c scriptRunnerTemplateContext) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		ID string `json:"id"`
	}{
		ID: c.opts.Key().String(),
	})
}

func (o optionsMap[K, C]) isBuilt(id uint32) bool {
	for _, v := range o {
		if v.isBuilt(id) {
			return true
		}
	}

	return false
}

func (o *opts[K, C]) isBuilt(id uint32) bool {
	return o.h.isBuilt(id)
}

func (o *opts[K, C]) isStale() bool {
	if o.h.isStaleOpts() {
		return true
	}
	if o.h.compiled.isStaleCompiled(o.h.compiledPrev) {
		return true
	}
	return false
}

func (o *optsHolder[C]) isStaleOpts() bool {
	if o.optsSetCounter == 1 && o.resetCounter > 0 {
		return false
	}
	isStale := func() bool {
		if len(o.optsCurr) != len(o.optsPrev) {
			return true
		}
		for k, v := range o.optsPrev {
			vv, found := o.optsCurr[k]
			if !found {
				return true
			}
			if strings.EqualFold(k, propsKeyImportContext) {
				// This is checked later.
			} else if si, ok := vv.(resource.StaleInfo); ok {
				if si.StaleVersion() > 0 {
					return true
				}
			} else {
				if !reflect.DeepEqual(v, vv) {
					return true
				}
			}
		}
		return false
	}()

	return isStale
}

func (o *opts[K, C]) isTouched(id uint32) bool {
	return o.h.isTouched(id)
}

func (o *optsHolder[C]) checkCompileErr() {
	if o.compileErr != nil {
		panic(o.compileErr)
	}
}

func (o *opts[K, C]) currPrev() (map[string]any, map[string]any) {
	return o.h.optsCurr, o.h.optsPrev
}

func init() {
	// We don't want any dependencies/change tracking on the top level Package,
	// we want finer grained control via Package.Group.
	var p any = &Package{}
	if _, ok := p.(identity.Identity); ok {
		panic("esbuid.Package should not implement identity.Identity")
	}
	if _, ok := p.(identity.DependencyManagerProvider); ok {
		panic("esbuid.Package should not implement identity.DependencyManagerProvider")
	}
}
