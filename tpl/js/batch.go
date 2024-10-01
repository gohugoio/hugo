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
	"sort"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
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
	UseScriptGroup(id string) BatcherScriptGroup
	Build() (*Package, error)
}

type BatcherScriptGroup interface {
	BatcherScriptGroupOps
	types.Closer
}

type BatcherScriptGroupOps interface {
	CallbackOptionsGetSetter
	Script(id string) BatcherScript
}

type BatcherScript interface {
	ScriptOptionsGetSetter
	AddInstance(id string, opts any) string
}

type ScriptOptions struct {
	// The script to build.
	// TODO1 handle stale.
	Resource resource.Resource

	// The import context to use.
	// Note that we will always fall back to the resource's own import context.
	ImportContext resource.ResourceGetter

	// Params marshaled to JSON.
	Params string
}

func (o *ScriptOptions) Dir() string {
	return path.Dir(o.Resource.(resource.PathProvider).Path())
}

type ScriptOptionsGetSetter interface {
	GetOptions() *ScriptOptions
	SetOptions(map[string]any) string
}

type CallbackOptionsGetSetter interface {
	GetCallbackOptions() *ScriptOptions
	SetCallbackOptions(map[string]any) string
}

func (ns *Namespace) Batch(id string, store *maps.Scratch) (Batcher, error) {
	key := path.Join(nsBundle, id)
	b := store.GetOrCreate(key, func() any {
		return &batcher{id: id, scriptGroups: make(map[string]*scriptGroup), client: ns}
	})
	return b.(*batcher), nil
}

func (b *batcher) UseScriptGroup(id string) BatcherScriptGroup {
	b.mu.Lock()

	group, found := b.scriptGroups[id]
	if !found {
		group = &scriptGroup{id: id, client: b.client, items: make(map[string]*scriptGroupItem)}
		b.scriptGroups[id] = group
	}

	// Aquire a more fine grained lock.
	b.mu.Unlock()
	group.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BatcherScriptGroupOps
		types.Closer
	}{
		group,
		close(func() error {
			closeOnce.Do(func() {
				group.mu.Unlock()
			})
			return nil
		}),
	}
}

type close func() error

func (c close) Close() error {
	return c()
}

var (
	_ Batcher               = (*batcher)(nil)
	_ BatcherScriptGroupOps = (*scriptGroup)(nil)
)

func decodeScriptInstance(opts any) scriptInstance {
	var inst scriptInstance
	if err := mapstructure.WeakDecode(opts, &inst); err != nil {
		panic(err)
	}
	return inst
}

func (b *scriptGroupItem) AddInstance(id string, opts any) string {
	b.instances[id] = decodeScriptInstance(opts)
	return ""
}

func (b *scriptGroup) Script(id string) BatcherScript {
	// The script group is locked at this point.
	item, found := b.items[id]
	if !found {
		item = &scriptGroupItem{
			id:        id,
			instances: make(map[string]scriptInstance),
			client:    b.client,
		}
		b.items[id] = item
	}

	return item
}

type batchTemplateContext struct {
	keyPath            string
	ID                 string
	CallbackImportPath string
	Modules            []batchTemplateExecutionsContext
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
	ID     string `json:"id"`
	Params any    `json:"params"`
}

type batchBuildOpts struct {
	Callback           *ScriptOptions
	js.ExternalOptions `mapstructure:",squash"`
}

type scriptsMany struct {
	mu sync.Mutex

	id string

	batches map[string]*scriptGroup

	client *Namespace
}

type scriptInstance struct {
	Params map[string]any
}

type scriptGroupItem struct {
	id string
	optionsGetSet
	instances map[string]scriptInstance

	client *Namespace
}

type scriptGroup struct {
	mu              sync.Mutex
	id              string
	callBackOptions optionsGetSet
	items           scripts

	client *Namespace
}

func (s *scriptGroup) GetCallbackOptions() *ScriptOptions {
	return s.callBackOptions.GetOptions()
}

func (s *scriptGroup) SetCallbackOptions(m map[string]any) string {
	return s.callBackOptions.SetOptions(m)
}

type scripts map[string]*scriptGroupItem

func (s scripts) Sorted() []*scriptGroupItem {
	var a []*scriptGroupItem
	for _, v := range s {
		a = append(a, v)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].id < a[j].id
	})
	return a
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
		Params:        string(paramsJSON),
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
	outDir       string
	id           string
	staleVersion uint32
	b            *batcher
	Groups       map[string]resource.Resources
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
	for _, v := range p.b.scriptGroups.Sorted() {
		if b := func() bool {
			v.mu.Lock()
			defer v.mu.Unlock()
			callbackOptions := v.GetCallbackOptions() // TODO1 validate.
			if callbackOptions != nil {
				if f(callbackOptions.Resource) { // TODO1 options.
					return true
				}
			}
			for _, vv := range v.items.Sorted() {
				if f(vv.Resource) {
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

func (b *batcher) Build() (*Package, error) {
	key := dynacache.CleanKey(b.id + ".js")
	p, err := b.client.bundlesCache.GetOrCreate(key, func(string) (*Package, error) {
		return b.build()
	})
	if err != nil {
		return nil, err
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

	keyPath := b.id

	type importContext struct {
		name           string
		resourceGetter resource.ResourceGetter
		scriptOptions  *ScriptOptions // TODO1 remove resourceGetter?
	}

	importResource := make(map[string]resource.Resource)
	resultResource := make(map[string]resource.Resource)
	importerImportContext := make(map[string]importContext)
	// importerImportContextStack := collections.NewStack[map[string]resource.ResourceGetter]()
	pathGroup := make(map[string]string)

	var entryPoints []string
	addResource := func(group, pth string, r resource.Resource, isResult bool) {
		pathGroup[pth] = group
		importResource[pth] = r
		if isResult {
			resultResource[pth] = r
		}
		entryPoints = append(entryPoints, pth)
	}

	if len(b.scriptGroups) > 0 {
		for k, v := range b.scriptGroups {
			keyPath := keyPath + "_" + k

			bopts := batchBuildOpts{
				Callback: v.GetCallbackOptions(),
			}
			var callbackImpPath string
			if bopts.Callback != nil {
				callbackImpPath = paths.AddLeadingSlash(keyPath + "_callback" + bopts.Callback.Resource.MediaType().FirstSuffix.FullSuffix)
				addResource(k, callbackImpPath, bopts.Callback.Resource, false)
			}

			t := &batchTemplateContext{
				keyPath:            keyPath,
				ID:                 v.id,
				CallbackImportPath: callbackImpPath,
			}

			for _, vv := range v.items.Sorted() {
				if vv.Resource == nil {
					// TODO1 others, init.
					return nil, fmt.Errorf("resource not set for %q", vv.id)
				}
				//  /mybundle_mains.js
				keyPath := keyPath + "_" + vv.id
				opts := vv.GetOptions()
				impPath := path.Join(hugoVirtualNS, opts.Dir(), keyPath+opts.Resource.MediaType().FirstSuffix.FullSuffix)
				impCtx := opts.ImportContext
				importerImportContext[impPath] = importContext{
					name:           keyPath,
					resourceGetter: impCtx,
					scriptOptions:  opts,
				}

				bt := batchTemplateExecutionsContext{
					ID:         vv.id,
					r:          vv.Resource,
					ImportPath: impPath,
				}
				importResource[bt.ImportPath] = vv.Resource
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
			fmt.Println("add batch B", v.id, s)
			importerImportContext[s] = importContext{
				name:           s,
				resourceGetter: nil,
				scriptOptions:  nil,
			}
			addResource(v.id, s, r, true)
		}
	}

	target := "es2018"

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

	var currentImportContext importContext

	var importResulveMu sync.Mutex

	jopts := js.Options{
		ExternalOptions: js.ExternalOptions{
			Format:  "esm",
			Target:  target,
			Defines: map[string]any{
				//"process.env.NODE_ENV": `"development"`,
			},
		},
		InternalOptions: js.InternalOptions{
			OutDir:         outDir,
			Write:          true,
			AllowOverwrite: true,
			Splitting:      true,
			ImportOnResolveFunc: func(imp string, args api.OnResolveArgs) string {
				importResulveMu.Lock()
				defer importResulveMu.Unlock()

				var found bool
				currentImportContext, found = importerImportContext[imp]
				if !found {
					currentImportContext = importerImportContext[args.Importer]
				}

				fmt.Println("ImportOnResolveFunc", "importer.base", path.Base(args.Importer), "path", args.Path)

				if _, found := importResource[imp]; found {
					return imp
				}

				if imp == "./foo.css" {
					fmt.Println("foo.css:", currentImportContext.name)
				}

				if currentImportContext.resourceGetter != nil {
					resolved := currentImportContext.resourceGetter.Get(imp)

					if resolved != nil {
						imp := hugoVirtualNS + resolved.(resource.PathProvider).Path()
						// TODO1 mu
						importResource[imp] = resolved
						return imp

					}
				}
				return ""
			},
			ImportOnLoadFunc: func(args api.OnLoadArgs) string {
				importResulveMu.Lock()
				defer importResulveMu.Unlock()

				imp := args.Path

				if r, found := importResource[imp]; found {
					content, err := r.(resource.ContentProvider).Content(context.Background()) // TODO1
					if err != nil {
						panic(err)
					}
					return cast.ToString(content)
				}
				return ""
			},
			ImportParamsOnLoadFunc: func(args api.OnLoadArgs) string {
				if currentImportContext.scriptOptions != nil {
					return currentImportContext.scriptOptions.Params
				}
				return ""
			},
			EntryPoints: entryPoints,
		},
	}

	result, err := b.client.client.BuildBundle(jopts)
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
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
		group, found := pathGroup[groupPath]

		if !found {
			return false, nil
		}

		if err := addFoo(o.filename, p, group, mt); err != nil {
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
		outDir: outDir,
		b:      b,
		id:     path.Join(nsBundle, b.id),
		Groups: groups,
	}, nil
}

type bundleResource string

func (b bundleResource) Name() string {
	return path.Base(string(b))
}

func (b bundleResource) Title() string {
	return b.Name()
}

func (b bundleResource) RelPermalink() string {
	return "/js/bundles/mybundle" + string(b)
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
