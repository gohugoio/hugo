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
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
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
	Store() *maps.Scratch

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

func (p *pagesFromDataTemplateContext) toPathMap(v any) (string, map[string]any, error) {
	m, err := maps.ToStringMapE(v)
	if err != nil {
		return "", nil, err
	}
	pathv, ok := m["path"]
	if !ok {
		return "", nil, fmt.Errorf("path not set")
	}
	path, err := cast.ToStringE(pathv)
	if err != nil || path == "" {
		return "", nil, fmt.Errorf("invalid path %q", path)
	}
	return path, m, nil
}

func (p *pagesFromDataTemplateContext) AddPage(v any) (string, error) {
	path, m, err := p.toPathMap(v)
	if err != nil {
		return "", err
	}

	if !p.p.buildState.checkHasChangedAndSetSourceInfo(path, m) {
		return "", nil
	}

	pd := pagemeta.DefaultPageConfig
	pd.IsFromContentAdapter = true

	if err := mapstructure.WeakDecode(m, &pd); err != nil {
		return "", fmt.Errorf("failed to decode page map: %w", err)
	}

	p.p.buildState.NumPagesAdded++

	if err := pd.Validate(true); err != nil {
		return "", err
	}

	return "", p.p.HandlePage(p.p, &pd)
}

func (p *pagesFromDataTemplateContext) AddResource(v any) (string, error) {
	path, m, err := p.toPathMap(v)
	if err != nil {
		return "", err
	}

	if !p.p.buildState.checkHasChangedAndSetSourceInfo(path, m) {
		return "", nil
	}

	var rd pagemeta.ResourceConfig
	if err := mapstructure.WeakDecode(m, &rd); err != nil {
		return "", err
	}

	p.p.buildState.NumResourcesAdded++

	if err := rd.Validate(); err != nil {
		return "", err
	}

	return "", p.p.HandleResource(p.p, &rd)
}

func (p *pagesFromDataTemplateContext) Site() page.Site {
	return p.p.Site
}

func (p *pagesFromDataTemplateContext) Store() *maps.Scratch {
	return p.p.store
}

func (p *pagesFromDataTemplateContext) EnableAllLanguages() string {
	p.p.buildState.EnableAllLanguages = true
	return ""
}

func NewPagesFromTemplate(opts PagesFromTemplateOptions) *PagesFromTemplate {
	return &PagesFromTemplate{
		PagesFromTemplateOptions: opts,
		PagesFromTemplateDeps:    opts.DepsFromSite(opts.Site),
		buildState: &BuildState{
			sourceInfosCurrent: maps.NewCache[string, *sourceInfo](),
		},
		store: maps.NewScratch(),
	}
}

type PagesFromTemplateOptions struct {
	Site         page.Site
	DepsFromSite func(page.Site) PagesFromTemplateDeps

	DependencyManager identity.Manager

	Watching bool

	HandlePage     func(pt *PagesFromTemplate, p *pagemeta.PageConfig) error
	HandleResource func(pt *PagesFromTemplate, p *pagemeta.ResourceConfig) error

	GoTmplFi hugofs.FileMetaInfo
}

type PagesFromTemplateDeps struct {
	TmplFinder tpl.TemplateParseFinder
	TmplExec   tpl.TemplateExecutor
}

var _ resource.Staler = (*PagesFromTemplate)(nil)

type PagesFromTemplate struct {
	PagesFromTemplateOptions
	PagesFromTemplateDeps
	buildState *BuildState
	store      *maps.Scratch
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
	NumPagesAdded      uint64
	NumResourcesAdded  uint64
	EnableAllLanguages bool
	ChangedIdentities  []identity.Identity
	DeletedPaths       []string
	Path               *paths.Path
}

type BuildState struct {
	StaleVersion uint32

	EnableAllLanguages bool

	// Paths deleted in the current build.
	DeletedPaths []string

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

func (b *BuildState) checkHasChangedAndSetSourceInfo(changedPath string, v any) bool {
	h := b.hash(v)
	si, found := b.sourceInfosPrevious.Get(changedPath)
	if found {
		b.sourceInfosCurrent.Set(changedPath, si)
		if si.hash == h {
			return false
		}
	} else {
		si = &sourceInfo{}
		b.sourceInfosCurrent.Set(changedPath, si)
	}
	si.hash = h
	return true
}

func (b *BuildState) resolveDeletedPaths() {
	if b.sourceInfosPrevious == nil {
		b.DeletedPaths = nil
		return
	}
	var paths []string
	b.sourceInfosPrevious.ForEeach(func(k string, _ *sourceInfo) {
		if _, found := b.sourceInfosCurrent.Get(k); !found {
			paths = append(paths, k)
		}
	})

	b.DeletedPaths = paths
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

type sourceInfo struct {
	hash uint64
}

func (p PagesFromTemplate) CloneForSite(s page.Site) *PagesFromTemplate {
	// We deliberately make them share the same DepenencyManager and Store.
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

func (p *PagesFromTemplate) Execute(ctx context.Context) (BuildInfo, error) {
	defer func() {
		p.buildState.PrepareNextBuild()
	}()

	f, err := p.GoTmplFi.Meta().Open()
	if err != nil {
		return BuildInfo{}, err
	}
	defer f.Close()

	tmpl, err := p.TmplFinder.Parse(filepath.ToSlash(p.GoTmplFi.Meta().Filename), helpers.ReaderToString(f))
	if err != nil {
		return BuildInfo{}, err
	}

	data := &pagesFromDataTemplateContext{
		p: p,
	}

	ctx = tpl.Context.DependencyManagerScopedProvider.Set(ctx, p)

	if err := p.TmplExec.ExecuteWithContext(ctx, tmpl, io.Discard, data); err != nil {
		return BuildInfo{}, err
	}

	if p.Watching {
		p.buildState.resolveDeletedPaths()
	}

	bi := BuildInfo{
		NumPagesAdded:      p.buildState.NumPagesAdded,
		NumResourcesAdded:  p.buildState.NumResourcesAdded,
		EnableAllLanguages: p.buildState.EnableAllLanguages,
		ChangedIdentities:  p.buildState.ChangedIdentities,
		DeletedPaths:       p.buildState.DeletedPaths,
		Path:               p.GoTmplFi.Meta().PathInfo,
	}

	return bi, nil
}

//////////////
