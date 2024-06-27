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

type Batcher interface {
	UseScript(id string) BatcherScript
	UseScriptGroup(id string) BatcherScriptMany
	Build() (*Package, error)
}

type BatcherScript interface {
	BatcherScriptOps
	types.Closer
}

type BatcherScriptOps interface {
	ResourceGetSetter
	AddInstance(id string, opts any) string
}

type BatcherScriptMany interface {
	BatcherScriptManyOps
	types.Closer
}

type BatcherScriptManyOps interface {
	CallbackGetSetter
	UseScript(id string) BatcherScript
}

type CallbackGetSetter interface {
	GetCallback() resource.Resource
	SetCallback(r resource.Resource) string
}

type ResourceGetSetter interface {
	GetResource() resource.Resource
	SetResource(r resource.Resource) string
}

func (ns *Namespace) Batch(id string, store *maps.Scratch) (Batcher, error) {
	key := path.Join(nsBundle, id)
	b := store.GetOrCreate(key, func() any {
		return &batcher{id: id, scriptOnes: make(map[string]*scriptOne), scriptManys: make(map[string]*scriptMany), client: ns}
	})
	return b.(*batcher), nil
}

func (b *batcher) UseScript(id string) BatcherScript {
	b.mu.Lock()

	one, found := b.scriptOnes[id]
	if !found {
		one = &scriptOne{
			id:        id,
			instances: make(map[string]scriptInstance),
			client:    b.client,
		}
		b.scriptOnes[id] = one
	}

	b.mu.Unlock()
	one.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BatcherScriptOps
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

func (b *batcher) UseScriptGroup(id string) BatcherScriptMany {
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
		BatcherScriptManyOps
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

type close func() error

func (c close) Close() error {
	return c()
}

var (
	_ Batcher              = (*batcher)(nil)
	_ BatcherScriptOps     = (*scriptOne)(nil)
	_ BatcherScriptManyOps = (*scriptMany)(nil)
)

func (b *scriptOne) AddInstance(id string, opts any) string {
	if b.r == nil {
		panic("resource not set")
	}
	if id == "" {
		panic("id not set")
	}

	b.instances[id] = decodeScriptInstance(opts)
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
	if r == nil {
		// TODO1 apply this to all Setters.
		panic("resource not set")
	}
	b.callback = r
	return ""
}

func (b *scriptMany) UseScript(id string) BatcherScript {
	item, found := b.items[id]
	if !found {
		item = &scriptManyItem{
			id:        id,
			instances: make(map[string]scriptInstance),
			client:    b.client,
		}
		b.items[id] = item
	}

	item.mu.Lock()

	// This will be auto closed if used in a with statement.
	// But the caller may also call Close, so make sure we only do it once.
	var closeOnce sync.Once

	return struct {
		BatcherScriptOps
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
	instances map[string]scriptInstance

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

type batcher struct {
	mu          sync.Mutex
	id          string
	scriptOnes  map[string]*scriptOne
	scriptManys map[string]*scriptMany

	client *Namespace
}

type resourceGetSet struct {
	r resource.Resource
}

func (r *resourceGetSet) Dir() string {
	return path.Dir(r.r.(resource.PathProvider).Path())
}

func (r *resourceGetSet) GetResource() resource.Resource {
	if resource.StaleVersion(r.r) > 0 {
		// Allow the client to set a new resource.
		return nil
	}
	return r.r
}

func (r *resourceGetSet) SetResource(res resource.Resource) string {
	if res == nil {
		panic("resource not set")
	}
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

func (b *batcher) build() (*Package, error) {
	defer herrors.Recover() // TODO1
	b.mu.Lock()
	defer b.mu.Unlock()

	keyPath := b.id

	importResource := make(map[string]resource.Resource)
	resultResource := make(map[string]resource.Resource)
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

	if len(b.scriptOnes) > 0 {
		for k, v := range b.scriptOnes {
			if v.r == nil {
				return nil, fmt.Errorf("resource not set for %q", k)
			}
			keyPath := keyPath + "_" + k
			resourcePath := paths.AddLeadingSlash(keyPath + v.r.MediaType().FirstSuffix.FullSuffix)
			addResource(k, resourcePath, v.r, true)

		}
	}

	if len(b.scriptManys) > 0 {
		for k, v := range b.scriptManys {
			keyPath := keyPath + "_" + k

			bopts := batchBuildOpts{
				Callback: v.callback,
			}
			var callbackImpPath string
			if bopts.Callback != nil {
				callbackImpPath = paths.AddLeadingSlash(keyPath + "_callback" + bopts.Callback.MediaType().FirstSuffix.FullSuffix)
				addResource(k, callbackImpPath, bopts.Callback, false)
			}

			t := &batchTemplateContext{
				keyPath:            keyPath,
				ID:                 v.id,
				CallbackImportPath: callbackImpPath,
			}

			for kk, vv := range v.items {
				if vv.r == nil {
					// TODO1 others.
					return nil, fmt.Errorf("resource not set for %q", kk)
				}
				keyPath := keyPath + "_" + kk
				const namespace = "@hugo-virtual"
				impPath := path.Join(namespace, vv.Dir(), keyPath+vv.r.MediaType().FirstSuffix.FullSuffix)
				bt := batchTemplateExecutionsContext{
					ID:         kk,
					r:          vv.r,
					ImportPath: impPath,
				}
				importResource[bt.ImportPath] = vv.r
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
			addResource(v.id, s, r, true)
		}
	}

	target := "es2018"

	conf := b.client.d.Conf
	absPublishDir := b.client.d.AbsPublishDir
	mediaTypes := b.client.d.ResourceSpec.MediaTypes()
	cssMt, _, _ := mediaTypes.GetFirstBySuffix("css")

	// TODO1 remove on close?
	fmt.Println(conf.Dirs().CacheDir)
	cacheDir := filepath.Join(b.client.d.SourceSpec.Cfg.Dirs().CacheDir, "_jsbatch")
	if err := os.Mkdir(cacheDir, 0o777); err != nil && !herrors.IsExist(err) {
		return nil, err
	}
	outDir, err := os.MkdirTemp(cacheDir, "jsbatch")
	if err != nil {
		return nil, err
	}

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
			ImportOnResolveFunc: func(imp string) string {
				if _, found := importResource[imp]; found {
					return imp
				}
				return ""
			},
			ImportOnLoadFunc: func(imp string) string {
				if r, found := importResource[imp]; found {
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
		if v.Exports != nil {
			// TODO1
			fmt.Println("   ", v.EntryPoint, " exports", v.Exports)
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
