// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package tplimpl

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/metrics"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl"
	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"
	"github.com/spf13/afero"
)

const (
	CategoryLayout Category = iota + 1
	CategoryBaseof
	CategoryMarkup
	CategoryShortcode
	CategoryPartial
	// Internal categories
	CategoryServer
	CategoryHugo
)

const (
	SubCategoryMain     SubCategory = iota
	SubCategoryEmbedded             // Internal Hugo templates
	SubCategoryInline               // Inline partials
)

const (
	containerMarkup          = "_markup"
	containerShortcodes      = "_shortcodes"
	shortcodesPathIdentifier = "/_shortcodes/"
	containerPartials        = "_partials"
)

const (
	layoutAll    = "all"
	layoutList   = "list"
	layoutSingle = "single"
)

var (
	_ identity.IdentityProvider             = (*TemplInfo)(nil)
	_ identity.IsProbablyDependentProvider  = (*TemplInfo)(nil)
	_ identity.IsProbablyDependencyProvider = (*TemplInfo)(nil)
)

const (
	processingStateInitial processingState = iota
	processingStateTransformed
)

// The identifiers may be truncated in the log, e.g.
// "executing "main" at <$scaled.SRelPermalin...>: can't evaluate field SRelPermalink in type *resource.Image"
// We need this to identify position in templates with base templates applied.
var identifiersRe = regexp.MustCompile(`at \<(.*?)(\.{3})?\>:`)

var weightNoMatch = weight{w1: -1}

//
//go:embed all:embedded/templates/*
var embeddedTemplatesFs embed.FS

func NewStore(opts StoreOptions, siteOpts SiteOptions) (*TemplateStore, error) {
	html, ok := opts.OutputFormats.GetByName("html")
	if !ok {
		panic("HTML output format not found")
	}
	s := &TemplateStore{
		opts:                 opts,
		siteOpts:             siteOpts,
		optsOrig:             opts,
		siteOptsOrig:         siteOpts,
		htmlFormat:           html,
		storeSite:            configureSiteStorage(siteOpts, opts.Watching),
		treeMain:             doctree.NewSimpleTree[map[nodeKey]*TemplInfo](),
		treeShortcodes:       doctree.NewSimpleTree[map[string]map[TemplateDescriptor]*TemplInfo](),
		templatesByPath:      maps.NewCache[string, *TemplInfo](),
		shortcodesByName:     maps.NewCache[string, *TemplInfo](),
		cacheLookupPartials:  maps.NewCache[string, *TemplInfo](),
		templatesSnapshotSet: maps.NewCache[*parse.Tree, struct{}](),

		// Note that the funcs passed below is just for name validation.
		tns: newTemplateNamespace(siteOpts.TemplateFuncs),

		dh: descriptorHandler{
			opts: opts,
		},
	}

	if err := s.init(); err != nil {
		return nil, err
	}
	if err := s.insertTemplates(nil, false); err != nil {
		return nil, err
	}
	if err := s.insertEmbedded(); err != nil {
		return nil, err
	}
	if err := s.parseTemplates(false); err != nil {
		return nil, err
	}
	if err := s.extractInlinePartials(false); err != nil {
		return nil, err
	}
	if err := s.transformTemplates(); err != nil {
		return nil, err
	}
	if err := s.tns.createPrototypes(true); err != nil {
		return nil, err
	}
	if err := s.prepareTemplates(); err != nil {
		return nil, err
	}
	return s, nil
}

//go:generate stringer -type Category

type Category int

type SiteOptions struct {
	Site          page.Site
	TemplateFuncs map[string]any
}

type StoreOptions struct {
	// The filesystem to use.
	Fs afero.Fs

	// The logger to use.
	Log loggers.Logger

	// The path parser to use.
	PathParser *paths.PathParser

	// Set when --enableTemplateMetrics is set.
	Metrics metrics.Provider

	// All configured output formats.
	OutputFormats output.Formats

	// All configured media types.
	MediaTypes media.Types

	// The default content language.
	DefaultContentLanguage string

	// The default output format.
	DefaultOutputFormat string

	// Taxonomy config.
	TaxonomySingularPlural map[string]string

	// Whether we are in watch or server mode.
	Watching bool

	// compiled.
	legacyMappingTaxonomy map[string]legacyOrdinalMapping
	legacyMappingTerm     map[string]legacyOrdinalMapping
	legacyMappingSection  map[string]legacyOrdinalMapping
}

//go:generate stringer -type SubCategory

type SubCategory int

type TemplInfo struct {
	// The category of this template.
	category Category

	subCategory SubCategory

	// PathInfo info.
	PathInfo *paths.Path

	// Set when backed by a file.
	Fi hugofs.FileMetaInfo

	// The template content with any leading BOM removed.
	content string

	// The parsed template.
	// Note that any baseof template will be applied later.
	Template tpl.Template

	// If no baseof is needed, this will be set to true.
	// E.g. shortcode templates do not need a baseof.
	noBaseOf bool

	// If NoBaseOf is false, we will look for the final template in this tree.
	baseVariants *doctree.SimpleTree[map[TemplateDescriptor]*TemplWithBaseApplied]

	// The template variants that are based on this template.
	overlays []*TemplInfo

	// The base template used, if any.
	base *TemplInfo

	// The descriptior that this template represents.
	D TemplateDescriptor

	// Parser state.
	ParseInfo ParseInfo

	// The execution counter for this template.
	executionCounter atomic.Uint64

	// processing state.
	state          processingState
	isLegacyMapped bool
}

func (ti *TemplInfo) SubCategory() SubCategory {
	return ti.subCategory
}

func (ti *TemplInfo) BaseVariantsSeq() iter.Seq[*TemplWithBaseApplied] {
	return func(yield func(*TemplWithBaseApplied) bool) {
		ti.baseVariants.Walk(func(key string, v map[TemplateDescriptor]*TemplWithBaseApplied) (bool, error) {
			for _, vv := range v {
				if !yield(vv) {
					return true, nil
				}
			}
			return false, nil
		})
	}
}

func (t *TemplInfo) IdentifierBase() string {
	if t.PathInfo == nil {
		return t.Name()
	}
	return t.PathInfo.IdentifierBase()
}

func (t *TemplInfo) GetIdentity() identity.Identity {
	return t
}

func (ti *TemplInfo) Name() string {
	if ti.Template == nil {
		if ti.PathInfo != nil {
			return ti.PathInfo.PathNoLeadingSlash()
		}
	}
	return ti.Template.Name()
}

func (ti *TemplInfo) Prepare() (*texttemplate.Template, error) {
	return ti.Template.Prepare()
}

func (t *TemplInfo) IsProbablyDependency(other identity.Identity) bool {
	return t.isProbablyTheSameIDAs(other)
}

func (t *TemplInfo) IsProbablyDependent(other identity.Identity) bool {
	for _, overlay := range t.overlays {
		if overlay.isProbablyTheSameIDAs(other) {
			return true
		}
	}
	return t.isProbablyTheSameIDAs(other)
}

func (ti *TemplInfo) String() string {
	if ti == nil {
		return "<nil>"
	}
	return ti.PathInfo.String()
}

func (ti *TemplInfo) findBestMatchBaseof(s *TemplateStore, d1 TemplateDescriptor, k1 string, slashCountK1 int, best *bestMatch) {
	if ti.baseVariants == nil {
		return
	}

	ti.baseVariants.WalkPath(k1, func(k2 string, v map[TemplateDescriptor]*TemplWithBaseApplied) (bool, error) {
		if !s.inPath(k1, k2) {
			return false, nil
		}
		slashCountK2 := strings.Count(k2, "/")
		distance := slashCountK1 - slashCountK2

		for d2, vv := range v {
			weight := s.dh.compareDescriptors(CategoryBaseof, false, d1, d2)
			weight.distance = distance
			if best.isBetter(weight, vv.Template) {
				best.updateValues(weight, k2, d2, vv.Template)
			}
		}
		return false, nil
	})
}

func (t *TemplInfo) isProbablyTheSameIDAs(other identity.Identity) bool {
	if t.IdentifierBase() == other.IdentifierBase() {
		return true
	}

	if t.Fi != nil && t.Fi.Meta().PathInfo != t.PathInfo {
		return other.IdentifierBase() == t.Fi.Meta().PathInfo.IdentifierBase()
	}

	return false
}

// Implements the additional methods in tpl.CurrentTemplateInfoOps.
func (ti *TemplInfo) Base() tpl.CurrentTemplateInfoCommonOps {
	return ti.base
}

func (ti *TemplInfo) Filename() string {
	if ti.Fi == nil {
		return ""
	}
	return ti.Fi.Meta().Filename
}

type TemplWithBaseApplied struct {
	// The template that's overlaid on top of the base template.
	Overlay *TemplInfo
	// The base template.
	Base *TemplInfo
	// This is the final template that can be used to render a page.
	Template *TemplInfo
}

// TemplateQuery is used in LookupPagesLayout to find the best matching template.
type TemplateQuery struct {
	// The path to walk down to.
	Path string

	// The name to look for. Used for shortcode queries.
	Name string

	// The category to look in.
	Category Category

	// The template descriptor to match against.
	Desc TemplateDescriptor

	// Whether to even consider this candidate.
	Consider func(candidate *TemplInfo) bool
}

func (q *TemplateQuery) init() {
	if q.Desc.Kind == kinds.KindTemporary {
		q.Desc.Kind = ""
	} else if kinds.GetKindMain(q.Desc.Kind) == "" {
		q.Desc.Kind = ""
	}
	if q.Desc.LayoutFromTemplate == "" && q.Desc.Kind != "" {
		if q.Desc.Kind == kinds.KindPage {
			q.Desc.LayoutFromTemplate = layoutSingle
		} else {
			q.Desc.LayoutFromTemplate = layoutList
		}
	}

	if q.Consider == nil {
		q.Consider = func(match *TemplInfo) bool {
			return true
		}
	}

	q.Name = strings.ToLower(q.Name)

	if q.Category == 0 {
		panic("category not set")
	}
}

type TemplateStore struct {
	opts       StoreOptions
	siteOpts   SiteOptions
	htmlFormat output.Format

	treeMain             *doctree.SimpleTree[map[nodeKey]*TemplInfo]
	treeShortcodes       *doctree.SimpleTree[map[string]map[TemplateDescriptor]*TemplInfo]
	templatesByPath      *maps.Cache[string, *TemplInfo]
	shortcodesByName     *maps.Cache[string, *TemplInfo]
	templatesSnapshotSet *maps.Cache[*parse.Tree, struct{}]

	dh descriptorHandler

	// The template namespace.
	tns *templateNamespace

	// Site specific state.
	// All above this is reused.
	storeSite *storeSite

	// For testing benchmarking.
	optsOrig     StoreOptions
	siteOptsOrig SiteOptions

	// caches. These need to be refreshed when the templates are refreshed.
	cacheLookupPartials *maps.Cache[string, *TemplInfo]
}

// NewFromOpts creates a new store with the same configuration as the original.
// Used for testing/benchmarking.
func (s *TemplateStore) NewFromOpts() (*TemplateStore, error) {
	return NewStore(s.optsOrig, s.siteOptsOrig)
}

// In the previous implementation of base templates in Hugo, we parsed and applied these base templates on
// request, e.g. in the middle of rendering. The idea was that we coulnd't know upfront which layoyt/base template
// combination that would be used.
// This, however, added a lot of complexity involving a careful dance of template cloning and parsing
// (Go HTML tenplates cannot be parsed after any of the templates in the tree have been executed).
// FindAllBaseTemplateCandidates finds all base template candidates for the given descriptor so we can apply them upfront.
// In this setup we may end up with unused base templates, but not having to do the cloning should more than make up for that.
func (s *TemplateStore) FindAllBaseTemplateCandidates(overlayKey string, desc TemplateDescriptor) []keyTemplateInfo {
	var result []keyTemplateInfo
	descBaseof := desc
	s.treeMain.Walk(func(k string, v map[nodeKey]*TemplInfo) (bool, error) {
		for _, vv := range v {
			if vv.category != CategoryBaseof {
				continue
			}

			if vv.D.isKindInLayout(desc.LayoutFromTemplate) && s.dh.compareDescriptors(CategoryBaseof, false, descBaseof, vv.D).w1 > 0 {
				result = append(result, keyTemplateInfo{Key: k, Info: vv})
			}
		}
		return false, nil
	})

	return result
}

func (t *TemplateStore) ExecuteWithContext(ctx context.Context, ti *TemplInfo, wr io.Writer, data any) error {
	defer func() {
		ti.executionCounter.Add(1)
		if ti.base != nil {
			ti.base.executionCounter.Add(1)
		}
	}()

	templ := ti.Template

	parent := tpl.Context.CurrentTemplate.Get(ctx)
	var level int
	if parent != nil {
		level = parent.Level + 1
	}
	currentTi := &tpl.CurrentTemplateInfo{
		Parent:                 parent,
		Level:                  level,
		CurrentTemplateInfoOps: ti,
	}

	ctx = tpl.Context.CurrentTemplate.Set(ctx, currentTi)

	const levelThreshold = 999
	if level > levelThreshold {
		return fmt.Errorf("maximum template call stack size exceeded in %q", ti.Filename())
	}

	if t.opts.Metrics != nil {
		defer t.opts.Metrics.MeasureSince(templ.Name(), time.Now())
	}

	execErr := t.storeSite.executer.ExecuteWithContext(ctx, ti, wr, data)
	if execErr != nil {
		return t.addFileContext(ti, "execute of template failed", execErr)
	}
	return nil
}

func (t *TemplateStore) GetFunc(name string) (reflect.Value, bool) {
	v, found := t.storeSite.execHelper.funcs[name]
	return v, found
}

func (s *TemplateStore) GetIdentity(p string) identity.Identity {
	p = paths.AddLeadingSlash(p)
	v, found := s.templatesByPath.Get(p)
	if !found {
		return nil
	}
	return v.GetIdentity()
}

func (t *TemplateStore) LookupByPath(templatePath string) *TemplInfo {
	v, _ := t.templatesByPath.Get(templatePath)
	return v
}

var bestPool = sync.Pool{
	New: func() any {
		return &bestMatch{}
	},
}

func (s *TemplateStore) getBest() *bestMatch {
	v := bestPool.Get()
	b := v.(*bestMatch)
	b.defaultOutputformat = s.opts.DefaultOutputFormat
	return b
}

func (s *TemplateStore) putBest(b *bestMatch) {
	b.reset()
	bestPool.Put(b)
}

func (s *TemplateStore) LookupPagesLayout(q TemplateQuery) *TemplInfo {
	q.init()
	key := s.key(q.Path)

	slashCountKey := strings.Count(key, "/")
	best1 := s.getBest()
	defer s.putBest(best1)
	s.findBestMatchWalkPath(q, key, slashCountKey, best1)
	if best1.w.w1 <= 0 {
		return nil
	}
	m := best1.templ
	if m.noBaseOf {
		return m
	}
	best1.reset()
	m.findBestMatchBaseof(s, q.Desc, key, slashCountKey, best1)
	if best1.w.w1 <= 0 {
		return nil
	}
	return best1.templ
}

func (s *TemplateStore) LookupPartial(pth string) *TemplInfo {
	ti, _ := s.cacheLookupPartials.GetOrCreate(pth, func() (*TemplInfo, error) {
		pi := s.opts.PathParser.Parse(files.ComponentFolderLayouts, pth).ForType(paths.TypePartial)
		k1, _, _, desc, err := s.toKeyCategoryAndDescriptor(pi)
		if err != nil {
			return nil, err
		}
		if desc.OutputFormat == "" && desc.MediaType == "" {
			// Assume HTML.
			desc.OutputFormat = s.htmlFormat.Name
			desc.MediaType = s.htmlFormat.MediaType.Type
			desc.IsPlainText = s.htmlFormat.IsPlainText
		}

		best := s.getBest()
		defer s.putBest(best)
		s.findBestMatchGet(s.key(path.Join(containerPartials, k1)), CategoryPartial, nil, desc, best)
		return best.templ, nil
	})

	return ti
}

func (s *TemplateStore) LookupShortcodeByName(name string) *TemplInfo {
	name = strings.ToLower(name)
	ti, _ := s.shortcodesByName.Get(name)
	if ti == nil {
		return nil
	}
	return ti
}

func (s *TemplateStore) LookupShortcode(q TemplateQuery) (*TemplInfo, error) {
	q.init()
	k1 := s.key(q.Path)

	slashCountK1 := strings.Count(k1, "/")

	best := s.getBest()
	defer s.putBest(best)

	s.treeShortcodes.WalkPath(k1, func(k2 string, m map[string]map[TemplateDescriptor]*TemplInfo) (bool, error) {
		if !s.inPath(k1, k2) {
			return false, nil
		}
		slashCountK2 := strings.Count(k2, "/")
		distance := slashCountK1 - slashCountK2

		v, found := m[q.Name]
		if !found {
			return false, nil
		}

		for k, vv := range v {
			best.candidates = append(best.candidates, vv)
			if !q.Consider(vv) {
				continue
			}

			weight := s.dh.compareDescriptors(q.Category, vv.subCategory == SubCategoryEmbedded, q.Desc, k)
			weight.distance = distance
			isBetter := best.isBetter(weight, vv)
			if isBetter {
				best.updateValues(weight, k2, k, vv)
			}
		}

		return false, nil
	})

	if best.w.w1 <= 0 {
		var err error
		if s := best.candidatesAsStringSlice(); s != nil {
			msg := fmt.Sprintf("no compatible template found for shortcode %q in %s", q.Name, s)
			if !q.Desc.IsPlainText {
				msg += "; note that to use plain text template shortcodes in HTML you need to use the shortcode {{% delimiter"
			}
			err = errors.New(msg)
		} else {
			err = fmt.Errorf("no template found for shortcode %q", q.Name)
		}
		return nil, err
	}

	return best.templ, nil
}

// PrintDebug is for testing/debugging only.
func (s *TemplateStore) PrintDebug(prefix string, category Category, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	printOne := func(key string, vv *TemplInfo) {
		level := strings.Count(key, "/")
		if category != vv.category {
			return
		}
		s := strings.ReplaceAll(strings.TrimSpace(vv.content), "\n", " ")
		ts := fmt.Sprintf("kind: %q layout: %q lang: %q content: %.30s", vv.D.Kind, vv.D.LayoutFromTemplate, vv.D.Lang, s)
		fmt.Fprintf(w, "%s%s %s\n", strings.Repeat(" ", level), key, ts)
	}
	s.treeMain.WalkPrefix(prefix, func(key string, v map[nodeKey]*TemplInfo) (bool, error) {
		for _, vv := range v {
			printOne(key, vv)
		}
		return false, nil
	})
	s.treeShortcodes.WalkPrefix(prefix, func(key string, v map[string]map[TemplateDescriptor]*TemplInfo) (bool, error) {
		for _, vv := range v {
			for _, vv2 := range vv {
				printOne(key, vv2)
			}
		}
		return false, nil
	})
}

func (s *TemplateStore) clearCaches() {
	s.cacheLookupPartials.Reset()
}

// RefreshFiles refreshes this store for the files matching the given predicate.
func (s *TemplateStore) RefreshFiles(include func(fi hugofs.FileMetaInfo) bool) error {
	s.clearCaches()

	if err := s.tns.createPrototypesParse(); err != nil {
		return err
	}
	if err := s.insertTemplates(include, true); err != nil {
		return err
	}
	if err := s.createTemplatesSnapshot(); err != nil {
		return err
	}
	if err := s.parseTemplates(true); err != nil {
		return err
	}
	if err := s.extractInlinePartials(true); err != nil {
		return err
	}

	if err := s.transformTemplates(); err != nil {
		return err
	}
	if err := s.tns.createPrototypes(false); err != nil {
		return err
	}
	if err := s.prepareTemplates(); err != nil {
		return err
	}
	return nil
}

func (s *TemplateStore) HasTemplate(templatePath string) bool {
	templatePath = strings.ToLower(templatePath)
	templatePath = paths.AddLeadingSlash(templatePath)
	return s.templatesByPath.Contains(templatePath)
}

func (t *TemplateStore) TextLookup(name string) *TemplInfo {
	templ := t.tns.standaloneText.Lookup(name)
	if templ == nil {
		return nil
	}
	return &TemplInfo{
		Template: templ,
	}
}

func (t *TemplateStore) TextParse(name, tpl string) (*TemplInfo, error) {
	templ, err := t.tns.standaloneText.New(name).Parse(tpl)
	if err != nil {
		return nil, err
	}
	return &TemplInfo{
		Template: templ,
	}, nil
}

func (t *TemplateStore) UnusedTemplates() []*TemplInfo {
	var unused []*TemplInfo

	for vv := range t.templates() {
		if vv.subCategory != SubCategoryMain || vv.isLegacyMapped {
			// Skip inline partials and internal templates.
			continue
		}
		if vv.executionCounter.Load() == 0 {
			unused = append(unused, vv)
		}
	}

	sort.Sort(byPath(unused))
	return unused
}

// WithSiteOpts creates a new store with the given site options.
// This is used to create per site template store, all sharing the same templates,
// but with a different template function execution context.
func (s TemplateStore) WithSiteOpts(opts SiteOptions) *TemplateStore {
	s.siteOpts = opts
	s.storeSite = configureSiteStorage(opts, s.opts.Watching)
	return &s
}

func (s *TemplateStore) findBestMatchGet(key string, category Category, consider func(candidate *TemplInfo) bool, desc TemplateDescriptor, best *bestMatch) {
	key = strings.ToLower(key)

	v := s.treeMain.Get(key)
	if v == nil {
		return
	}

	for k, vv := range v {
		if vv.category != category {
			continue
		}

		if consider != nil && !consider(vv) {
			continue
		}

		weight := s.dh.compareDescriptors(category, vv.subCategory == SubCategoryEmbedded, desc, k.d)
		if best.isBetter(weight, vv) {
			best.updateValues(weight, key, k.d, vv)
		}
	}
}

func (s *TemplateStore) inPath(k1, k2 string) bool {
	if k1 != k2 && !strings.HasPrefix(k1, k2+"/") {
		return false
	}
	return true
}

func (s *TemplateStore) findBestMatchWalkPath(q TemplateQuery, k1 string, slashCountK1 int, best *bestMatch) {
	s.treeMain.WalkPath(k1, func(k2 string, v map[nodeKey]*TemplInfo) (bool, error) {
		if !s.inPath(k1, k2) {
			return false, nil
		}
		slashCountK2 := strings.Count(k2, "/")
		distance := slashCountK1 - slashCountK2

		for k, vv := range v {
			if vv.category != q.Category {
				continue
			}

			if !q.Consider(vv) {
				continue
			}

			weight := s.dh.compareDescriptors(q.Category, vv.subCategory == SubCategoryEmbedded, q.Desc, k.d)

			weight.distance = distance
			isBetter := best.isBetter(weight, vv)

			if isBetter {
				best.updateValues(weight, k2, k.d, vv)
			}
		}

		return false, nil
	})
}

func (t *TemplateStore) addDeferredTemplate(owner *TemplInfo, name string, n *parse.ListNode) error {
	if _, found := t.templatesByPath.Get(name); found {
		return nil
	}

	var templ tpl.Template

	if owner.D.IsPlainText {
		prototype := t.tns.parseText
		tt, err := prototype.New(name).Parse("")
		if err != nil {
			return fmt.Errorf("failed to parse empty text template %q: %w", name, err)
		}
		tt.Tree.Root = n
		templ = tt
	} else {
		prototype := t.tns.parseHTML
		tt, err := prototype.New(name).Parse("")
		if err != nil {
			return fmt.Errorf("failed to parse empty HTML template %q: %w", name, err)
		}
		tt.Tree.Root = n
		templ = tt
	}

	t.templatesByPath.Set(name, &TemplInfo{
		Fi:       owner.Fi,
		PathInfo: owner.PathInfo,
		D:        owner.D,
		Template: templ,
	})

	return nil
}

func (s *TemplateStore) addFileContext(ti *TemplInfo, what string, inerr error) error {
	if ti.Fi == nil {
		return inerr
	}

	identifiers := s.extractIdentifiers(inerr.Error())

	checkFilename := func(fi hugofs.FileMetaInfo, inErr error) (error, bool) {
		lineMatcher := func(m herrors.LineMatcher) int {
			if m.Position.LineNumber != m.LineNumber {
				return -1
			}

			for _, id := range identifiers {
				if strings.Contains(m.Line, id) {
					// We found the line, but return a 0 to signal to
					// use the column from the error message.
					return 0
				}
			}
			return -1
		}

		f, err := fi.Meta().Open()
		if err != nil {
			return inErr, false
		}
		defer f.Close()

		fe := herrors.NewFileErrorFromName(inErr, fi.Meta().Filename)
		fe.UpdateContent(f, lineMatcher)

		return fe, fe.ErrorContext().Position.IsValid()
	}

	inerr = fmt.Errorf("%s: %w", what, inerr)

	var (
		currentErr error
		ok         bool
	)

	if currentErr, ok = checkFilename(ti.Fi, inerr); ok {
		return currentErr
	}

	if ti.base != nil {
		if currentErr, ok = checkFilename(ti.base.Fi, inerr); ok {
			return currentErr
		}
	}

	return currentErr
}

func (s *TemplateStore) extractIdentifiers(line string) []string {
	m := identifiersRe.FindAllStringSubmatch(line, -1)
	identifiers := make([]string, len(m))
	for i := range m {
		identifiers[i] = m[i][1]
	}
	return identifiers
}

func (s *TemplateStore) extractInlinePartials(rebuild bool) error {
	isPartialName := func(s string) bool {
		return strings.HasPrefix(s, "partials/") || strings.HasPrefix(s, "_partials/")
	}

	// We may find both inline and external partials in the current template namespaces,
	// so only add the ones we have not seen before.
	for templ := range s.allRawTemplates() {
		if templ.Name() == "" || !isPartialName(templ.Name()) {
			continue
		}
		if rebuild && s.templatesSnapshotSet.Contains(getParseTree(templ)) {
			// This partial was not created during this build.
			continue
		}
		name := templ.Name()
		if !paths.HasExt(name) {
			// Assume HTML. This in line with how the lookup works.
			name = name + s.htmlFormat.MediaType.FirstSuffix.FullSuffix
		}
		if !strings.HasPrefix(name, "_") {
			name = "_" + name
		}
		pi := s.opts.PathParser.Parse(files.ComponentFolderLayouts, name)
		ti, err := s.insertTemplate(pi, nil, SubCategoryInline, false, s.treeMain)
		if err != nil {
			return err
		}

		if ti != nil {
			ti.Template = templ
			ti.noBaseOf = true
			ti.subCategory = SubCategoryInline
			ti.D.IsPlainText = isText(templ)
		}
	}

	return nil
}

func (s *TemplateStore) allRawTemplates() iter.Seq[tpl.Template] {
	p := s.tns
	return func(yield func(tpl.Template) bool) {
		for t := range p.templatesIn(p.parseHTML) {
			if !yield(t) {
				return
			}
		}
		for t := range p.templatesIn(p.parseText) {
			if !yield(t) {
				return
			}
		}

		for _, tt := range p.baseofHtmlClones {
			for t := range p.templatesIn(tt) {
				if !yield(t) {
					return
				}
			}
		}
		for _, tt := range p.baseofTextClones {
			for t := range p.templatesIn(tt) {
				if !yield(t) {
					return
				}
			}
		}
	}
}

func (s *TemplateStore) insertEmbedded() error {
	return fs.WalkDir(embeddedTemplatesFs, ".", func(tpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d == nil || d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		templb, err := embeddedTemplatesFs.ReadFile(tpath)
		if err != nil {
			return err
		}

		// Get the newlines on Windows in line with how we had it back when we used Go Generate
		// to write the templates to Go files.
		templ := string(bytes.ReplaceAll(templb, []byte("\r\n"), []byte("\n")))
		name := strings.TrimPrefix(filepath.ToSlash(tpath), "embedded/templates/")

		insertOne := func(name, content string) error {
			pi := s.opts.PathParser.Parse(files.ComponentFolderLayouts, name)
			var (
				ti  *TemplInfo
				err error
			)
			if pi.Section() == containerShortcodes {
				ti, err = s.insertShortcode(pi, nil, false, s.treeShortcodes)
				if err != nil {
					return err
				}
			} else {
				ti, err = s.insertTemplate(pi, nil, SubCategoryEmbedded, false, s.treeMain)
				if err != nil {
					return err
				}
			}

			if ti != nil {
				// Currently none of the embedded templates need a baseof template.
				ti.noBaseOf = true
				ti.content = content
				ti.subCategory = SubCategoryEmbedded
			}

			return nil
		}

		// Copy the embedded HTML table render hook to each output format.
		// See https://github.com/gohugoio/hugo/issues/13351.
		if name == path.Join(containerMarkup, "render-table.html") {
			for _, of := range s.opts.OutputFormats {
				path := paths.TrimExt(name) + "." + of.Name + of.MediaType.FirstSuffix.FullSuffix
				if err := insertOne(path, templ); err != nil {
					return err
				}
			}

			return nil
		}

		if err := insertOne(name, templ); err != nil {
			return err
		}

		if aliases, found := embeddedTemplatesAliases[name]; found {
			for _, alias := range aliases {
				if err := insertOne(alias, templ); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *TemplateStore) setTemplateByPath(p string, ti *TemplInfo) {
	s.templatesByPath.Set(p, ti)
}

func (s *TemplateStore) insertShortcode(pi *paths.Path, fi hugofs.FileMetaInfo, replace bool, tree doctree.Tree[map[string]map[TemplateDescriptor]*TemplInfo]) (*TemplInfo, error) {
	k1, k2, _, d, err := s.toKeyCategoryAndDescriptor(pi)
	if err != nil {
		return nil, err
	}
	m := tree.Get(k1)
	if m == nil {
		m = make(map[string]map[TemplateDescriptor]*TemplInfo)
		tree.Insert(k1, m)
	}

	m1, found := m[k2]
	if found {
		if _, found := m1[d]; found {
			if !replace {
				return nil, nil
			}
		}
	} else {
		m1 = make(map[TemplateDescriptor]*TemplInfo)
		m[k2] = m1
	}

	ti := &TemplInfo{
		PathInfo: pi,
		Fi:       fi,
		D:        d,
		category: CategoryShortcode,
		noBaseOf: true,
	}

	m1[d] = ti

	s.shortcodesByName.Set(k2, ti)
	s.setTemplateByPath(pi.Path(), ti)

	if fi != nil {
		if pi2 := fi.Meta().PathInfo; pi2 != pi {
			s.setTemplateByPath(pi2.Path(), ti)
		}
	}

	return ti, nil
}

func (s *TemplateStore) insertTemplate(pi *paths.Path, fi hugofs.FileMetaInfo, subCategory SubCategory, replace bool, tree doctree.Tree[map[nodeKey]*TemplInfo]) (*TemplInfo, error) {
	key, _, category, d, err := s.toKeyCategoryAndDescriptor(pi)
	// See #13577. Warn for now.
	if err != nil {
		var loc string
		if fi != nil {
			loc = fmt.Sprintf("file %q", fi.Meta().Filename)
		} else {
			loc = fmt.Sprintf("path %q", pi.Path())
		}
		s.opts.Log.Warnf("skipping template %s: %s", loc, err)
		return nil, nil
	}

	return s.insertTemplate2(pi, fi, key, category, subCategory, d, replace, false, tree)
}

func (s *TemplateStore) insertTemplate2(
	pi *paths.Path,
	fi hugofs.FileMetaInfo,
	key string,
	category Category,
	subCategory SubCategory,
	d TemplateDescriptor,
	replace, isLegacyMapped bool,
	tree doctree.Tree[map[nodeKey]*TemplInfo],
) (*TemplInfo, error) {
	if category == 0 {
		panic("category not set")
	}

	if category == CategoryPartial && d.OutputFormat == "" && d.MediaType == "" {
		// See issue #13601.
		d.OutputFormat = s.htmlFormat.Name
		d.MediaType = s.htmlFormat.MediaType.Type
	}

	m := tree.Get(key)
	nk := nodeKey{c: category, d: d}

	if m == nil {
		m = make(map[nodeKey]*TemplInfo)
		tree.Insert(key, m)
	}

	nkExisting, existingFound := m[nk]
	if !replace && existingFound && fi != nil && nkExisting.Fi != nil {
		// See issue #13715.
		// We do the merge on the file system level, but from Hugo v0.146.0 we have a situation where
		// the project may well have a different layouts layout compared to the theme(s) it uses.
		// We could possibly have fixed that on a lower (file system) level, but since this is just
		// a temporary situation (until all projects are updated),
		// do a replace here if the file comes from higher up in the module chain.
		replace = fi.Meta().ModuleOrdinal < nkExisting.Fi.Meta().ModuleOrdinal
	}

	if !replace && existingFound {
		// Always replace inline partials to allow for reloading.
		replace = subCategory == SubCategoryInline && nkExisting.subCategory == SubCategoryInline
	}

	if !replace && existingFound {
		if len(pi.Identifiers()) >= len(nkExisting.PathInfo.Identifiers()) {
			// e.g. /pages/home.foo.html and  /pages/home.html where foo may be a valid language name in another site.
			return nil, nil
		}
	}

	ti := &TemplInfo{
		PathInfo:       pi,
		Fi:             fi,
		D:              d,
		category:       category,
		noBaseOf:       category > CategoryLayout,
		isLegacyMapped: isLegacyMapped,
	}

	m[nk] = ti

	if !isLegacyMapped {
		s.setTemplateByPath(pi.Path(), ti)
		if fi != nil {
			if pi2 := fi.Meta().PathInfo; pi2 != pi {
				s.setTemplateByPath(pi2.Path(), ti)
			}
		}
	}

	return ti, nil
}

func (s *TemplateStore) insertTemplates(include func(fi hugofs.FileMetaInfo) bool, partialRebuild bool) error {
	if include == nil {
		include = func(fi hugofs.FileMetaInfo) bool {
			return true
		}
	}

	// Set if we need to reset the base variants.
	var (
		resetBaseVariants bool
	)

	legacyOrdinalMappings := map[legacyTargetPathIdentifiers]legacyOrdinalMappingFi{}

	walker := func(pth string, fi hugofs.FileMetaInfo) error {
		if fi.IsDir() {
			return nil
		}

		if isDotFile(pth) || isBackupFile(pth) {
			return nil
		}

		if !include(fi) {
			return nil
		}

		piOrig := fi.Meta().PathInfo

		// Convert any legacy value to new format.
		fromLegacyPath := func(pi *paths.Path) *paths.Path {
			p := pi.Path()
			p = strings.TrimPrefix(p, "/_default")
			if strings.HasPrefix(p, "/shortcodes") || strings.HasPrefix(p, "/partials") {
				// Insert an underscore so it becomes /_shortcodes or /_partials.
				p = "/_" + p[1:]
			}

			if strings.Contains(p, "-"+baseNameBaseof) {
				// Before Hugo 0.146.0 we prepended one identifier (layout, type or kind) in front of the baseof keyword,
				// and then separated with a hyphen before the baseof keyword.
				// This identifier needs to be moved right after the baseof keyword and the hyphen removed, e.g.
				// /docs/list-baseof.html => /docs/baseof.list.html.
				dir, name := path.Split(p)
				hyphenIdx := strings.Index(name, "-")
				if hyphenIdx > 0 {
					id := name[:hyphenIdx]
					name = name[hyphenIdx+1+len(baseNameBaseof):]
					if !strings.HasPrefix(name, ".") {
						name = "." + name
					}
					p = path.Join(dir, baseNameBaseof+"."+id+name)
				}
			}
			if p == pi.Path() {
				return pi
			}
			return s.opts.PathParser.Parse(files.ComponentFolderLayouts, p)
		}

		pi := piOrig
		var applyLegacyMapping bool
		switch pi.Section() {
		case containerPartials, containerShortcodes, containerMarkup:
			// OK.
		default:
			pi = fromLegacyPath(pi)
			applyLegacyMapping = strings.Count(pi.Path(), "/") <= 2
		}

		if applyLegacyMapping {
			handleMapping := func(m1 legacyOrdinalMapping) {
				key := legacyTargetPathIdentifiers{
					targetPath:     m1.mapping.targetPath,
					targetCategory: m1.mapping.targetCategory,
					kind:           m1.mapping.targetDesc.Kind,
					lang:           pi.Lang(),
					ext:            pi.Ext(),
					outputFormat:   pi.OutputFormat(),
				}

				if m2, ok := legacyOrdinalMappings[key]; ok {
					if m1.ordinal < m2.m.ordinal {
						// Higher up == better match.
						legacyOrdinalMappings[key] = legacyOrdinalMappingFi{m1, fi}
					}
				} else {
					legacyOrdinalMappings[key] = legacyOrdinalMappingFi{m1, fi}
				}
			}

			if m1, ok := s.opts.legacyMappingTaxonomy[piOrig.PathBeforeLangAndOutputFormatAndExt()]; ok {
				handleMapping(m1)
			}

			if m1, ok := s.opts.legacyMappingTerm[piOrig.PathBeforeLangAndOutputFormatAndExt()]; ok {
				handleMapping(m1)
			}

			const (
				sectionKindToken = "SECTIONKIND"
				sectionToken     = "THESECTION"
			)

			base := piOrig.PathBeforeLangAndOutputFormatAndExt()
			identifiers := []string{}
			if pi.Layout() != "" {
				identifiers = append(identifiers, pi.Layout())
			}
			if pi.Kind() != "" {
				identifiers = append(identifiers, pi.Kind())
			}

			shouldIncludeSection := func(section string) bool {
				switch section {
				case containerShortcodes, containerPartials, containerMarkup:
					return false
				case "taxonomy", "":
					return false
				default:
					for k, v := range s.opts.TaxonomySingularPlural {
						if k == section || v == section {
							return false
						}
					}
					return true
				}
			}
			if shouldIncludeSection(pi.Section()) {
				identifiers = append(identifiers, pi.Section())
			}

			identifiers = helpers.UniqueStrings(identifiers)

			// Tokens on e.g. form /SECTIONKIND/THESECTION
			insertSectionTokens := func(section string) []string {
				kindOnly := isLayoutStandard(section)
				var ss []string
				s1 := base
				if !kindOnly {
					s1 = strings.ReplaceAll(s1, section, sectionToken)
				}
				s1 = strings.ReplaceAll(s1, kinds.KindSection, sectionKindToken)
				if s1 != base {
					ss = append(ss, s1)
				}
				s1 = strings.ReplaceAll(base, kinds.KindSection, sectionKindToken)
				if !kindOnly {
					s1 = strings.ReplaceAll(s1, section, sectionToken)
				}
				if s1 != base {
					ss = append(ss, s1)
				}

				helpers.UniqueStringsReuse(ss)

				return ss
			}

			for _, id := range identifiers {
				if id == "" {
					continue
				}

				p := insertSectionTokens(id)
				for _, ss := range p {
					if m1, ok := s.opts.legacyMappingSection[ss]; ok {
						targetPath := m1.mapping.targetPath

						if targetPath != "" {
							targetPath = strings.ReplaceAll(targetPath, sectionToken, id)
							targetPath = strings.ReplaceAll(targetPath, sectionKindToken, id)
							targetPath = strings.ReplaceAll(targetPath, "//", "/")
						}
						m1.mapping.targetPath = targetPath
						handleMapping(m1)
					}
				}
			}

		}

		if partialRebuild && pi.NameNoIdentifier() == baseNameBaseof {
			// A baseof file has changed.
			resetBaseVariants = true
		}

		var ti *TemplInfo
		var err error
		if pi.Type() == paths.TypeShortcode {
			ti, err = s.insertShortcode(pi, fi, partialRebuild, s.treeShortcodes)
			if err != nil || ti == nil {
				return err
			}
		} else {
			ti, err = s.insertTemplate(pi, fi, SubCategoryMain, partialRebuild, s.treeMain)
			if err != nil || ti == nil {
				return err
			}
		}

		if err := s.tns.readTemplateInto(ti); err != nil {
			return err
		}

		return nil
	}

	if err := helpers.Walk(s.opts.Fs, "", walker); err != nil {
		if !herrors.IsNotExist(err) {
			return err
		}
		return nil
	}

	for k, v := range legacyOrdinalMappings {
		targetPath := k.targetPath
		m := v.m.mapping
		fi := v.fi
		pi := fi.Meta().PathInfo
		outputFormat, mediaType := s.resolveOutputFormatAndOrMediaType(k.outputFormat, k.ext)
		category := m.targetCategory
		desc := m.targetDesc
		desc.Kind = k.kind
		desc.Lang = k.lang
		desc.OutputFormat = outputFormat.Name
		desc.IsPlainText = outputFormat.IsPlainText
		desc.MediaType = mediaType.Type

		ti, err := s.insertTemplate2(pi, fi, targetPath, category, SubCategoryMain, desc, true, true, s.treeMain)
		if err != nil {
			return err
		}
		if ti == nil {
			continue
		}
		ti.isLegacyMapped = true
		if err := s.tns.readTemplateInto(ti); err != nil {
			return err
		}

	}

	if resetBaseVariants {
		s.tns.baseofHtmlClones = nil
		s.tns.baseofTextClones = nil
		s.treeMain.Walk(func(key string, v map[nodeKey]*TemplInfo) (bool, error) {
			for _, vv := range v {
				if !vv.noBaseOf {
					vv.state = processingStateInitial
				}
			}
			return false, nil
		})
	}

	return nil
}

func (s *TemplateStore) key(dir string) string {
	dir = paths.AddLeadingSlash(dir)
	if dir == "/" {
		return ""
	}
	return paths.TrimTrailing(dir)
}

func (s *TemplateStore) createTemplatesSnapshot() error {
	s.templatesSnapshotSet.Reset()
	for t := range s.allRawTemplates() {
		s.templatesSnapshotSet.Set(getParseTree(t), struct{}{})
	}
	return nil
}

func (s *TemplateStore) parseTemplates(replace bool) error {
	if err := func() error {
		// Read and parse all templates.
		for _, v := range s.treeMain.All() {
			for _, vv := range v {
				if vv.state == processingStateTransformed {
					continue
				}
				if err := s.parseTemplate(vv, replace); err != nil {
					return err
				}
			}
		}

		// Lookup and apply base templates where needed.
		for key, v := range s.treeMain.All() {
			for _, vv := range v {
				if vv.state == processingStateTransformed {
					continue
				}
				if !vv.noBaseOf {
					d := vv.D
					// Find all compatible base templates.
					baseTemplates := s.FindAllBaseTemplateCandidates(key, d)
					if len(baseTemplates) == 0 {
						// The regular expression used to detect if a template needs a base template has some
						// rare false positives. Assume we don't need one.
						vv.noBaseOf = true
						if err := s.parseTemplate(vv, replace); err != nil {
							return err
						}
						continue
					}
					vv.baseVariants = doctree.NewSimpleTree[map[TemplateDescriptor]*TemplWithBaseApplied]()

					for _, base := range baseTemplates {
						if err := s.tns.applyBaseTemplate(vv, base); err != nil {
							return err
						}
					}

				}
			}
		}

		return nil
	}(); err != nil {
		return err
	}

	// Prese shortcodes.
	for _, v := range s.treeShortcodes.All() {
		for _, vv := range v {
			for _, vvv := range vv {
				if vvv.state == processingStateTransformed {
					continue
				}
				if err := s.parseTemplate(vvv, replace); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// prepareTemplates prepares all templates for execution.
func (s *TemplateStore) prepareTemplates() error {
	for t := range s.templates() {
		if t.category == CategoryBaseof {
			continue
		}
		if _, err := t.Prepare(); err != nil {
			return err
		}
	}
	return nil
}

type PathTemplateDescriptor struct {
	Path string
	Desc TemplateDescriptor
}

// resolveOutputFormatAndOrMediaType resolves the output format and/or media type
// based on the given output format suffix and media type suffix.
// Either of the suffixes can be empty, and the function will try to find a match
// based on the other suffix. If both are empty, the function will return zero values.
func (s *TemplateStore) resolveOutputFormatAndOrMediaType(ofs, mns string) (output.Format, media.Type) {
	var outputFormat output.Format
	var mediaType media.Type

	if ofs != "" {
		if of, found := s.opts.OutputFormats.GetByName(ofs); found {
			outputFormat = of
			mediaType = of.MediaType
		}
	}

	if mns != "" && mediaType.IsZero() {
		if of, found := s.opts.OutputFormats.GetBySuffix(mns); found {
			outputFormat = of
			mediaType = of.MediaType
		} else {
			if mt, _, found := s.opts.MediaTypes.GetFirstBySuffix(mns); found {
				mediaType = mt
				if outputFormat.IsZero() {
					// For e.g. index.xml we will in the default confg now have the application/rss+xml  media type.
					// Try a last time to find the output format using the SubType as the name.
					// As to template resolution, this value is currently only used to
					// decide if this is a text or HTML template.
					outputFormat, _ = s.opts.OutputFormats.GetByName(mt.SubType)
				}
			}
		}
	}

	return outputFormat, mediaType
}

// templates iterates over all templates in the store.
// Note that for templates with one or more base templates applied,
// we will yield the variants, e.g. the templates that's actually in use.
func (s *TemplateStore) templates() iter.Seq[*TemplInfo] {
	return func(yield func(*TemplInfo) bool) {
		for _, v := range s.treeMain.All() {
			for _, vv := range v {
				if !vv.noBaseOf {
					for vvv := range vv.BaseVariantsSeq() {
						if !yield(vvv.Template) {
							return
						}
					}
				} else {
					if !yield(vv) {
						return
					}
				}
			}
		}
		for _, v := range s.treeShortcodes.All() {
			for _, vv := range v {
				for _, vvv := range vv {
					if !yield(vvv) {
						return
					}
				}
			}
		}
	}
}

func (s *TemplateStore) toKeyCategoryAndDescriptor(p *paths.Path) (string, string, Category, TemplateDescriptor, error) {
	k1 := p.Dir()
	k2 := ""

	outputFormat, mediaType := s.resolveOutputFormatAndOrMediaType(p.OutputFormat(), p.Ext())
	nameNoIdentifier := p.NameNoIdentifier()

	d := TemplateDescriptor{
		Lang:               p.Lang(),
		OutputFormat:       p.OutputFormat(),
		MediaType:          mediaType.Type,
		Kind:               p.Kind(),
		LayoutFromTemplate: p.Layout(),
		IsPlainText:        outputFormat.IsPlainText,
	}

	d.normalizeFromFile()

	section := p.Section()

	var category Category
	switch p.Type() {
	case paths.TypeShortcode:
		category = CategoryShortcode
	case paths.TypePartial:
		category = CategoryPartial
	case paths.TypeMarkup:
		category = CategoryMarkup
	}

	if category == 0 {
		if nameNoIdentifier == baseNameBaseof {
			category = CategoryBaseof
		} else {
			switch section {
			case "_hugo":
				category = CategoryHugo
			case "_server":
				category = CategoryServer
			default:
				category = CategoryLayout
			}
		}
	}

	if category == CategoryPartial {
		d.LayoutFromTemplate = ""
		k1 = p.PathNoIdentifier()
	}

	if category == CategoryShortcode {
		k1 = p.PathNoIdentifier()

		parts := strings.Split(k1, "/"+containerShortcodes+"/")
		k1 = parts[0]
		if len(parts) > 1 {
			k2 = parts[1]
		}
		k1 = s.key(k1)
	}

	// Legacy layout for home page.
	if d.LayoutFromTemplate == "index" {
		if d.Kind == "" {
			d.Kind = kinds.KindHome
		}
		d.LayoutFromTemplate = ""
	}

	if d.LayoutFromTemplate == d.Kind {
		d.LayoutFromTemplate = ""
	}

	k1 = strings.TrimPrefix(k1, "/_default")
	if k1 == "/" {
		k1 = ""
	}

	if category == CategoryMarkup {
		// We store all template nodes for a given directory on the same level.
		k1 = strings.TrimSuffix(k1, "/_markup")
		parts := strings.Split(d.LayoutFromTemplate, "-")
		if len(parts) < 2 {
			return "", "", 0, TemplateDescriptor{}, fmt.Errorf("unrecognized render hook template")
		}
		// Either 2 or 3 parts, e.g. render-codeblock-go.
		d.Variant1 = parts[1]
		if len(parts) > 2 {
			d.Variant2 = parts[2]
		}
		d.LayoutFromTemplate = "" // This allows using page layout as part of the key for lookups.
	}

	return k1, k2, category, d, nil
}

func (s *TemplateStore) transformTemplates() error {
	lookup := func(name string, in *TemplInfo) *TemplInfo {
		if in.D.IsPlainText {
			templ := in.Template.(*texttemplate.Template).Lookup(name)
			if templ != nil {
				return &TemplInfo{
					Template: templ,
				}
			}
		} else {
			templ := in.Template.(*htmltemplate.Template).Lookup(name)
			if templ != nil {
				return &TemplInfo{
					Template: templ,
				}
			}
		}

		return nil
	}

	for vv := range s.templates() {
		if vv.state == processingStateTransformed {
			continue
		}
		vv.state = processingStateTransformed
		if vv.category == CategoryBaseof {
			continue
		}
		tctx, err := applyTemplateTransformers(vv, lookup)
		if err != nil {
			return err
		}
		for name, node := range tctx.deferNodes {
			if err := s.addDeferredTemplate(vv, name, node); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TemplateStore) init() error {
	// Before Hugo 0.146 we had a very elaborate template lookup system, especially for
	// terms and taxonomies. This is a way of preserving backwards compatibility
	// by mapping old paths into the new tree.
	s.opts.legacyMappingTaxonomy = make(map[string]legacyOrdinalMapping)
	s.opts.legacyMappingTerm = make(map[string]legacyOrdinalMapping)
	s.opts.legacyMappingSection = make(map[string]legacyOrdinalMapping)

	// Placeholders.
	const singular = "SINGULAR"
	const plural = "PLURAL"

	replaceTokens := func(s, singularv, pluralv string) string {
		s = strings.Replace(s, singular, singularv, -1)
		s = strings.Replace(s, plural, pluralv, -1)
		return s
	}

	hasSingularOrPlural := func(s string) bool {
		return strings.Contains(s, singular) || strings.Contains(s, plural)
	}

	expand := func(v layoutLegacyMapping) []layoutLegacyMapping {
		var result []layoutLegacyMapping

		if hasSingularOrPlural(v.sourcePath) || hasSingularOrPlural(v.target.targetPath) {
			for s, p := range s.opts.TaxonomySingularPlural {
				target := v.target
				target.targetPath = replaceTokens(target.targetPath, s, p)
				vv := replaceTokens(v.sourcePath, s, p)
				result = append(result, layoutLegacyMapping{sourcePath: vv, target: target})
			}
		} else {
			result = append(result, v)
		}
		return result
	}

	expandSections := func(v layoutLegacyMapping) []layoutLegacyMapping {
		var result []layoutLegacyMapping
		result = append(result, v)
		baseofVariant := v
		baseofVariant.sourcePath += "-" + baseNameBaseof
		baseofVariant.target.targetCategory = CategoryBaseof
		result = append(result, baseofVariant)
		return result
	}

	var terms []layoutLegacyMapping
	for _, v := range legacyTermMappings {
		terms = append(terms, expand(v)...)
	}
	var taxonomies []layoutLegacyMapping
	for _, v := range legacyTaxonomyMappings {
		taxonomies = append(taxonomies, expand(v)...)
	}
	var sections []layoutLegacyMapping
	for _, v := range legacySectionMappings {
		sections = append(sections, expandSections(v)...)
	}

	for i, m := range terms {
		s.opts.legacyMappingTerm[m.sourcePath] = legacyOrdinalMapping{ordinal: i, mapping: m.target}
	}
	for i, m := range taxonomies {
		s.opts.legacyMappingTaxonomy[m.sourcePath] = legacyOrdinalMapping{ordinal: i, mapping: m.target}
	}
	for i, m := range sections {
		s.opts.legacyMappingSection[m.sourcePath] = legacyOrdinalMapping{ordinal: i, mapping: m.target}
	}

	return nil
}

type TemplateStoreProvider interface {
	GetTemplateStore() *TemplateStore
}

type TextTemplatHandler interface {
	ExecuteWithContext(ctx context.Context, ti *TemplInfo, wr io.Writer, data any) error
	TextLookup(name string) *TemplInfo
	TextParse(name, tpl string) (*TemplInfo, error)
}

type bestMatch struct {
	templ      *TemplInfo
	desc       TemplateDescriptor
	w          weight
	key        string
	candidates []*TemplInfo

	// settings.
	defaultOutputformat string
}

func (best *bestMatch) reset() {
	best.templ = nil
	best.w = weight{}
	best.desc = TemplateDescriptor{}
	best.key = ""
	best.candidates = nil
}

func (best *bestMatch) candidatesAsStringSlice() []string {
	if len(best.candidates) == 0 {
		return nil
	}
	candidates := make([]string, len(best.candidates))
	for i, v := range best.candidates {
		candidates[i] = v.PathInfo.Path()
	}
	return candidates
}

func (best *bestMatch) isBetter(w weight, ti *TemplInfo) bool {
	if best.templ == nil {
		// Anything is better than nothing.
		return true
	}

	if w.w1 <= 0 {
		if best.w.w1 <= 0 {
			return ti.PathInfo.Path() < best.templ.PathInfo.Path()
		}
		return false
	}

	// Note that for render hook templates, we need to make
	// the embedded render hook template wih if they're a better match,
	// e.g. render-codeblock-goat.html.
	if best.templ.category != CategoryMarkup && best.w.w1 > 0 {
		currentBestIsEmbedded := best.templ.subCategory == SubCategoryEmbedded
		if currentBestIsEmbedded {
			if ti.subCategory != SubCategoryEmbedded {
				return true
			}
		} else {
			if ti.subCategory == SubCategoryEmbedded {
				// Prefer user provided template.
				return false
			}
		}
	}

	if w.distance < best.w.distance {
		if w.w2 < best.w.w2 {
			return false
		}
		if w.w3 < best.w.w3 {
			return false
		}
	} else {
		if w.w1 < best.w.w1 {
			return false
		}
	}

	if w.isEqualWeights(best.w) {
		// Tie breakers.
		if w.distance < best.w.distance {
			return true
		}

		return ti.PathInfo.Path() < best.templ.PathInfo.Path()
	}

	return true
}

func (best *bestMatch) updateValues(w weight, key string, k TemplateDescriptor, vv *TemplInfo) {
	best.w = w
	best.templ = vv
	best.desc = k
	best.key = key
}

type byPath []*TemplInfo

func (a byPath) Len() int { return len(a) }
func (a byPath) Less(i, j int) bool {
	return a[i].PathInfo.Path() < a[j].PathInfo.Path()
}

func (a byPath) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type keyTemplateInfo struct {
	Key  string
	Info *TemplInfo
}

type nodeKey struct {
	c Category
	d TemplateDescriptor
}

type processingState int

// the parts of a template store that's set per site.
type storeSite struct {
	opts       SiteOptions
	execHelper *templateExecHelper
	executer   texttemplate.Executer
}

type weight struct {
	w1       int
	w2       int
	w3       int
	distance int
}

func isLayoutStandard(s string) bool {
	switch s {
	case layoutAll, layoutList, layoutSingle:
		return true
	default:
		return false
	}
}

func (w weight) isEqualWeights(other weight) bool {
	return w.w1 == other.w1 && w.w2 == other.w2 && w.w3 == other.w3
}

func configureSiteStorage(opts SiteOptions, watching bool) *storeSite {
	funcsv := make(map[string]reflect.Value)

	for k, v := range opts.TemplateFuncs {
		vv := reflect.ValueOf(v)
		funcsv[k] = vv
	}

	// Duplicate Go's internal funcs here for faster lookups.
	for k, v := range htmltemplate.GoFuncs {
		if _, exists := funcsv[k]; !exists {
			vv, ok := v.(reflect.Value)
			if !ok {
				vv = reflect.ValueOf(v)
			}
			funcsv[k] = vv
		}
	}

	for k, v := range texttemplate.GoFuncs {
		if _, exists := funcsv[k]; !exists {
			funcsv[k] = v
		}
	}

	s := &storeSite{
		opts: opts,
		execHelper: &templateExecHelper{
			watching:   watching,
			funcs:      funcsv,
			site:       reflect.ValueOf(opts.Site),
			siteParams: reflect.ValueOf(opts.Site.Params()),
		},
	}

	s.executer = texttemplate.NewExecuter(s.execHelper)

	return s
}

func isBackupFile(path string) bool {
	return path[len(path)-1] == '~'
}

func isDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}
