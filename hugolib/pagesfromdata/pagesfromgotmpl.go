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

package pagesfromdata

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hstore"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/tpl/tplimpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

type PagesFromDataTemplateContext interface {
	// AddPage adds a new page to the site.
	// The first return value will always be an empty string.
	AddPage(any) (string, error)

	// AddResource adds a new resource to the site.
	// The first return value will always be an empty string.
	AddResource(any) (string, error)

	// The site to which the pages will be added.
	Site() page.Site

	// The same template may be executed multiple times for multiple languages.
	// The Store can be used to store state between these invocations.
	Store() *hstore.Scratch

	// By default, the template will be executed for the language
	// defined by the _content.gotmpl file (e.g. its mount definition).
	// This method can be used to activate the template for all languages.
	// The return value will always be an empty string.
	EnableAllLanguages() string
}

var _ PagesFromDataTemplateContext = (*pagesFromDataTemplateContext)(nil)

type pagesFromDataTemplateContext struct {
	p *PagesFromTemplate
}

func (p *pagesFromDataTemplateContext) toPathSitesMap(v any) (string, map[string]any, map[string]any, error) {
	m, err := maps.ToStringMapE(v)
	if err != nil {
		return "", nil, nil, err
	}

	path, err := cast.ToStringE(m["path"])
	if err != nil {
		return "", nil, nil, fmt.Errorf("invalid path %q", path)
	}

	sites := maps.ToStringMap(m["sites"])

	return path, sites, m, nil
}

func (p *pagesFromDataTemplateContext) AddPage(v any) (string, error) {
	path, sites, m, err := p.toPathSitesMap(v)
	if err != nil {
		return "", err
	}

	hash, hasChanged := p.p.buildState.checkHasChangedAndSetSourceInfo(path, sites, m)
	if !hasChanged {
		return "", nil
	}

	pe := &pagemeta.PageConfigEarly{
		IsFromContentAdapter: true,
		Frontmatter:          m,
		SourceEntryHash:      hash,
	}

	// The rest will be handled after the cascade is calculated and applied.
	if err := mapstructure.WeakDecode(pe.Frontmatter, pe); err != nil {
		err = fmt.Errorf("failed to decode page map: %w", err)
		return "", err
	}

	if err := pe.Init(true); err != nil {
		return "", err
	}

	p.p.buildState.NumPagesAdded++

	return "", p.p.HandlePage(p.p, pe)
}

func (p *pagesFromDataTemplateContext) AddResource(v any) (string, error) {
	path, sites, m, err := p.toPathSitesMap(v)
	if err != nil {
		return "", err
	}

	hash, hasChanged := p.p.buildState.checkHasChangedAndSetSourceInfo(path, sites, m)
	if !hasChanged {
		return "", nil
	}

	var rd pagemeta.ResourceConfig
	if err := mapstructure.WeakDecode(m, &rd); err != nil {
		return "", err
	}
	rd.ContentAdapterSourceEntryHash = hash

	p.p.buildState.NumResourcesAdded++

	if err := rd.Validate(); err != nil {
		return "", err
	}

	return "", p.p.HandleResource(p.p, &rd)
}

func (p *pagesFromDataTemplateContext) Site() page.Site {
	return p.p.Site
}

func (p *pagesFromDataTemplateContext) Store() *hstore.Scratch {
	return p.p.store
}

func (p *pagesFromDataTemplateContext) EnableAllLanguages() string {
	p.p.buildState.EnableAllLanguages = true
	return ""
}

func (p *pagesFromDataTemplateContext) EnableAllDimensions() string {
	p.p.buildState.EnableAllDimensions = true
	return ""
}

func NewPagesFromTemplate(opts PagesFromTemplateOptions) *PagesFromTemplate {
	return &PagesFromTemplate{
		PagesFromTemplateOptions: opts,
		PagesFromTemplateDeps:    opts.DepsFromSite(opts.Site),
		buildState: &BuildState{
			sourceInfosCurrent: maps.NewCache[string, *sourceInfo](),
		},
		store: hstore.NewScratch(),
	}
}

type PagesFromTemplateOptions struct {
	Site         page.Site
	DepsFromSite func(page.Site) PagesFromTemplateDeps

	DependencyManager identity.Manager

	Watching bool

	HandlePage     func(pt *PagesFromTemplate, p *pagemeta.PageConfigEarly) error
	HandleResource func(pt *PagesFromTemplate, p *pagemeta.ResourceConfig) error

	GoTmplFi hugofs.FileMetaInfo
}

type PagesFromTemplateDeps struct {
	TemplateStore *tplimpl.TemplateStore
}

var _ resource.Staler = (*PagesFromTemplate)(nil)

type PagesFromTemplate struct {
	PagesFromTemplateOptions
	PagesFromTemplateDeps
	buildState *BuildState
	store      *hstore.Scratch
}

func (b *PagesFromTemplate) AddChange(id identity.Identity) {
	b.buildState.ChangedIdentities = append(b.buildState.ChangedIdentities, id)
}

func (b *PagesFromTemplate) MarkStale() {
	b.buildState.StaleVersion++
}

func (b *PagesFromTemplate) StaleVersion() uint32 {
	return b.buildState.StaleVersion
}

type BuildInfo struct {
	NumPagesAdded       uint64
	NumResourcesAdded   uint64
	EnableAllLanguages  bool
	EnableAllDimensions bool
	ChangedIdentities   []identity.Identity
	DeletedPaths        []PathHashes
	Path                *paths.Path
}

type BuildState struct {
	StaleVersion uint32

	EnableAllLanguages  bool
	EnableAllDimensions bool

	// PathHashes deleted in the current build.
	DeletedPaths []PathHashes

	// Changed identities in the current build.
	ChangedIdentities []identity.Identity

	NumPagesAdded     uint64
	NumResourcesAdded uint64

	sourceInfosCurrent  *maps.Cache[string, *sourceInfo]
	sourceInfosPrevious *maps.Cache[string, *sourceInfo]
}

func (b *BuildState) hash(v any) uint64 {
	return hashing.HashUint64(v)
}

type sourceInfo struct {
	siteHashes map[uint64]uint64
}

func (b *BuildState) checkHasChangedAndSetSourceInfo(changedPath string, sites map[string]any, v any) (uint64, bool) {
	hv := b.hash(v)
	hsites := b.hash(sites)

	si, _ := b.sourceInfosCurrent.GetOrCreate(changedPath, func() (*sourceInfo, error) {
		return &sourceInfo{
			siteHashes: make(map[uint64]uint64),
		}, nil
	})

	if h, found := si.siteHashes[hsites]; found && h == hv {
		return hv, false
	}

	if psi, found := b.sourceInfosPrevious.Get(changedPath); found {
		if h, found := psi.siteHashes[hsites]; found && h == hv {
			// Not changed.
			si.siteHashes[hsites] = hv
			return hv, false
		}
	}

	// It has changed.
	si.siteHashes[hsites] = hv
	return hv, true
}

type PathHashes struct {
	Path   string
	Hashes map[uint64]struct{}
}

func (b *BuildState) resolveDeletedPaths() {
	if b.sourceInfosPrevious == nil {
		b.DeletedPaths = nil
		return
	}
	var pathsHashes []PathHashes
	b.sourceInfosPrevious.ForEeach(func(k string, pv *sourceInfo) bool {
		if cv, found := b.sourceInfosCurrent.Get(k); !found {
			pathsHashes = append(pathsHashes, PathHashes{Path: k, Hashes: map[uint64]struct{}{}})
		} else {
			deleted := map[uint64]struct{}{}
			for k, ph := range pv.siteHashes {
				ch, found := cv.siteHashes[k]
				if !found || ch != ph {
					deleted[ph] = struct{}{}
				}
			}
			if len(deleted) > 0 {
				pathsHashes = append(pathsHashes, PathHashes{Path: k, Hashes: deleted})
			}
		}
		return true
	})

	b.DeletedPaths = pathsHashes
}

func (b *BuildState) PrepareNextBuild() {
	b.sourceInfosPrevious = b.sourceInfosCurrent
	b.sourceInfosCurrent = maps.NewCache[string, *sourceInfo]()
	b.StaleVersion = 0
	b.DeletedPaths = nil
	b.ChangedIdentities = nil
	b.NumPagesAdded = 0
	b.NumResourcesAdded = 0
}

func (p PagesFromTemplate) CloneForSite(s page.Site) *PagesFromTemplate {
	// We deliberately make them share the same DependencyManager and Store.
	p.PagesFromTemplateOptions.Site = s
	p.PagesFromTemplateDeps = p.PagesFromTemplateOptions.DepsFromSite(s)
	p.buildState = &BuildState{
		sourceInfosCurrent: maps.NewCache[string, *sourceInfo](),
	}
	return &p
}

func (p PagesFromTemplate) CloneForGoTmpl(fi hugofs.FileMetaInfo) *PagesFromTemplate {
	p.PagesFromTemplateOptions.GoTmplFi = fi
	return &p
}

func (p *PagesFromTemplate) GetDependencyManagerForScope(scope int) identity.Manager {
	return p.DependencyManager
}

func (p *PagesFromTemplate) GetDependencyManagerForScopesAll() []identity.Manager {
	return []identity.Manager{p.DependencyManager}
}

func (p *PagesFromTemplate) Execute(ctx context.Context) (BuildInfo, error) {
	defer func() {
		p.buildState.PrepareNextBuild()
	}()

	f, err := p.GoTmplFi.Meta().Open()
	if err != nil {
		return BuildInfo{}, err
	}
	defer f.Close()

	tmpl, err := p.TemplateStore.TextParse(filepath.ToSlash(p.GoTmplFi.Meta().Filename), helpers.ReaderToString(f))
	if err != nil {
		return BuildInfo{}, err
	}

	data := &pagesFromDataTemplateContext{
		p: p,
	}

	ctx = tpl.Context.DependencyManagerScopedProvider.Set(ctx, p)

	if err := p.TemplateStore.ExecuteWithContext(ctx, tmpl, io.Discard, data); err != nil {
		return BuildInfo{}, err
	}

	if p.Watching {
		p.buildState.resolveDeletedPaths()
	}

	bi := BuildInfo{
		NumPagesAdded:       p.buildState.NumPagesAdded,
		NumResourcesAdded:   p.buildState.NumResourcesAdded,
		EnableAllLanguages:  p.buildState.EnableAllLanguages,
		EnableAllDimensions: p.buildState.EnableAllDimensions,
		ChangedIdentities:   p.buildState.ChangedIdentities,
		DeletedPaths:        p.buildState.DeletedPaths,
		Path:                p.GoTmplFi.Meta().PathInfo,
	}

	return bi, nil
}
