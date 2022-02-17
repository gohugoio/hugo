// Copyright 2019 The Hugo Authors. All rights reserved.
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

package tplimpl

import (
	"bytes"
	"context"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/pkg/errors"

	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/tpl"
)

const (
	textTmplNamePrefix = "_text/"

	shortcodesPathPrefix = "shortcodes/"
	internalPathPrefix   = "_internal/"
	baseFileBase         = "baseof"
)

// The identifiers may be truncated in the log, e.g.
// "executing "main" at <$scaled.SRelPermalin...>: can't evaluate field SRelPermalink in type *resource.Image"
var identifiersRe = regexp.MustCompile(`at \<(.*?)(\.{3})?\>:`)

var embeddedTemplatesAliases = map[string][]string{
	"shortcodes/twitter.html": {"shortcodes/tweet.html"},
}

var (
	_ tpl.TemplateManager         = (*templateExec)(nil)
	_ tpl.TemplateHandler         = (*templateExec)(nil)
	_ tpl.TemplateFuncGetter      = (*templateExec)(nil)
	_ tpl.TemplateFinder          = (*templateExec)(nil)
	_ tpl.UnusedTemplatesProvider = (*templateExec)(nil)

	_ tpl.Template = (*templateState)(nil)
	_ tpl.Info     = (*templateState)(nil)
)

var baseTemplateDefineRe = regexp.MustCompile(`^{{-?\s*define`)

// needsBaseTemplate returns true if the first non-comment template block is a
// define block.
// If a base template does not exist, we will handle that when it's used.
func needsBaseTemplate(templ string) bool {
	idx := -1
	inComment := false
	for i := 0; i < len(templ); {
		if !inComment && strings.HasPrefix(templ[i:], "{{/*") {
			inComment = true
			i += 4
		} else if inComment && strings.HasPrefix(templ[i:], "*/}}") {
			inComment = false
			i += 4
		} else {
			r, size := utf8.DecodeRuneInString(templ[i:])
			if !inComment {
				if strings.HasPrefix(templ[i:], "{{") {
					idx = i
					break
				} else if !unicode.IsSpace(r) {
					break
				}
			}
			i += size
		}
	}

	if idx == -1 {
		return false
	}

	return baseTemplateDefineRe.MatchString(templ[idx:])
}

func newIdentity(name string) identity.Manager {
	return identity.NewManager(identity.NewPathIdentity(files.ComponentFolderLayouts, name))
}

func newStandaloneTextTemplate(funcs map[string]interface{}) tpl.TemplateParseFinder {
	return &textTemplateWrapperWithLock{
		RWMutex:  &sync.RWMutex{},
		Template: texttemplate.New("").Funcs(funcs),
	}
}

func newTemplateExec(d *deps.Deps) (*templateExec, error) {
	exec, funcs := newTemplateExecuter(d)
	funcMap := make(map[string]interface{})
	for k, v := range funcs {
		funcMap[k] = v.Interface()
	}

	var templateUsageTracker map[string]templateInfo
	if d.Cfg.GetBool("printUnusedTemplates") {
		templateUsageTracker = make(map[string]templateInfo)
	}

	h := &templateHandler{
		nameBaseTemplateName: make(map[string]string),
		transformNotFound:    make(map[string]*templateState),
		identityNotFound:     make(map[string][]identity.Manager),

		shortcodes:   make(map[string]*shortcodeTemplates),
		templateInfo: make(map[string]tpl.Info),
		baseof:       make(map[string]templateInfo),
		needsBaseof:  make(map[string]templateInfo),

		main: newTemplateNamespace(funcMap),

		Deps:                d,
		layoutHandler:       output.NewLayoutHandler(),
		layoutsFs:           d.BaseFs.Layouts.Fs,
		layoutTemplateCache: make(map[layoutCacheKey]tpl.Template),

		templateUsageTracker: templateUsageTracker,
	}

	if err := h.loadEmbedded(); err != nil {
		return nil, err
	}

	if err := h.loadTemplates(); err != nil {
		return nil, err
	}

	e := &templateExec{
		d:               d,
		executor:        exec,
		funcs:           funcs,
		templateHandler: h,
	}

	d.SetTmpl(e)
	d.SetTextTmpl(newStandaloneTextTemplate(funcMap))

	if d.WithTemplate != nil {
		if err := d.WithTemplate(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func newTemplateNamespace(funcs map[string]interface{}) *templateNamespace {
	return &templateNamespace{
		prototypeHTML: htmltemplate.New("").Funcs(funcs),
		prototypeText: texttemplate.New("").Funcs(funcs),
		templateStateMap: &templateStateMap{
			templates: make(map[string]*templateState),
		},
	}
}

func newTemplateState(templ tpl.Template, info templateInfo) *templateState {
	return &templateState{
		info:      info,
		typ:       info.resolveType(),
		Template:  templ,
		Manager:   newIdentity(info.name),
		parseInfo: tpl.DefaultParseInfo,
	}
}

type layoutCacheKey struct {
	d output.LayoutDescriptor
	f string
}

type templateExec struct {
	d        *deps.Deps
	executor texttemplate.Executer
	funcs    map[string]reflect.Value

	*templateHandler
}

func (t templateExec) Clone(d *deps.Deps) *templateExec {
	exec, funcs := newTemplateExecuter(d)
	t.executor = exec
	t.funcs = funcs
	t.d = d
	return &t
}

func (t *templateExec) Execute(templ tpl.Template, wr io.Writer, data interface{}) error {
	return t.ExecuteWithContext(context.Background(), templ, wr, data)
}

func (t *templateExec) ExecuteWithContext(ctx context.Context, templ tpl.Template, wr io.Writer, data interface{}) error {
	if rlocker, ok := templ.(types.RLocker); ok {
		rlocker.RLock()
		defer rlocker.RUnlock()
	}
	if t.Metrics != nil {
		defer t.Metrics.MeasureSince(templ.Name(), time.Now())
	}

	if t.templateUsageTracker != nil {
		if ts, ok := templ.(*templateState); ok {
			t.templateUsageTrackerMu.Lock()
			if _, found := t.templateUsageTracker[ts.Name()]; !found {
				t.templateUsageTracker[ts.Name()] = ts.info
			}

			if !ts.baseInfo.IsZero() {
				if _, found := t.templateUsageTracker[ts.baseInfo.name]; !found {
					t.templateUsageTracker[ts.baseInfo.name] = ts.baseInfo
				}
			}
			t.templateUsageTrackerMu.Unlock()
		}
	}

	execErr := t.executor.ExecuteWithContext(ctx, templ, wr, data)
	if execErr != nil {
		execErr = t.addFileContext(templ, execErr)
	}
	return execErr
}

func (t *templateExec) UnusedTemplates() []tpl.FileInfo {
	if t.templateUsageTracker == nil {
		return nil
	}
	var unused []tpl.FileInfo

	for _, ti := range t.needsBaseof {
		if _, found := t.templateUsageTracker[ti.name]; !found {
			unused = append(unused, ti)
		}
	}

	for _, ti := range t.baseof {
		if _, found := t.templateUsageTracker[ti.name]; !found {
			unused = append(unused, ti)
		}
	}

	for _, ts := range t.main.templates {
		ti := ts.info
		if strings.HasPrefix(ti.name, "_internal/") || ti.realFilename == "" {
			continue
		}

		if _, found := t.templateUsageTracker[ti.name]; !found {
			unused = append(unused, ti)
		}
	}

	sort.Slice(unused, func(i, j int) bool {
		return unused[i].Name() < unused[j].Name()
	})

	return unused
}

func (t *templateExec) GetFunc(name string) (reflect.Value, bool) {
	v, found := t.funcs[name]
	return v, found
}

func (t *templateExec) MarkReady() error {
	var err error
	t.readyInit.Do(func() {
		// We only need the clones if base templates are in use.
		if len(t.needsBaseof) > 0 {
			err = t.main.createPrototypes()
		}
	})

	return err
}

type templateHandler struct {
	main        *templateNamespace
	needsBaseof map[string]templateInfo
	baseof      map[string]templateInfo

	readyInit sync.Once

	// This is the filesystem to load the templates from. All the templates are
	// stored in the root of this filesystem.
	layoutsFs afero.Fs

	layoutHandler *output.LayoutHandler

	layoutTemplateCache   map[layoutCacheKey]tpl.Template
	layoutTemplateCacheMu sync.RWMutex

	*deps.Deps

	// Used to get proper filenames in errors
	nameBaseTemplateName map[string]string

	// Holds name and source of template definitions not found during the first
	// AST transformation pass.
	transformNotFound map[string]*templateState

	// Holds identities of templates not found during first pass.
	identityNotFound map[string][]identity.Manager

	// shortcodes maps shortcode name to template variants
	// (language, output format etc.) of that shortcode.
	shortcodes map[string]*shortcodeTemplates

	// templateInfo maps template name to some additional information about that template.
	// Note that for shortcodes that same information is embedded in the
	// shortcodeTemplates type.
	templateInfo map[string]tpl.Info

	// May be nil.
	templateUsageTracker   map[string]templateInfo
	templateUsageTrackerMu sync.Mutex
}

// AddTemplate parses and adds a template to the collection.
// Templates with name prefixed with "_text" will be handled as plain
// text templates.
func (t *templateHandler) AddTemplate(name, tpl string) error {
	templ, err := t.addTemplateTo(t.newTemplateInfo(name, tpl), t.main)
	if err == nil {
		t.applyTemplateTransformers(t.main, templ)
	}
	return err
}

func (t *templateHandler) Lookup(name string) (tpl.Template, bool) {
	templ, found := t.main.Lookup(name)
	if found {
		return templ, true
	}

	return nil, false
}

func (t *templateHandler) LookupLayout(d output.LayoutDescriptor, f output.Format) (tpl.Template, bool, error) {
	key := layoutCacheKey{d, f.Name}
	t.layoutTemplateCacheMu.RLock()
	if cacheVal, found := t.layoutTemplateCache[key]; found {
		t.layoutTemplateCacheMu.RUnlock()
		return cacheVal, true, nil
	}
	t.layoutTemplateCacheMu.RUnlock()

	t.layoutTemplateCacheMu.Lock()
	defer t.layoutTemplateCacheMu.Unlock()

	templ, found, err := t.findLayout(d, f)
	if err == nil && found {
		t.layoutTemplateCache[key] = templ
		return templ, true, nil
	}

	return nil, false, err
}

// This currently only applies to shortcodes and what we get here is the
// shortcode name.
func (t *templateHandler) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	name = templateBaseName(templateShortcode, name)
	s, found := t.shortcodes[name]
	if !found {
		return nil, false, false
	}

	sv, found := s.fromVariants(variants)
	if !found {
		return nil, false, false
	}

	more := len(s.variants) > 1

	return sv.ts, true, more
}

// LookupVariants returns all variants of name, nil if none found.
func (t *templateHandler) LookupVariants(name string) []tpl.Template {
	name = templateBaseName(templateShortcode, name)
	s, found := t.shortcodes[name]
	if !found {
		return nil
	}

	variants := make([]tpl.Template, len(s.variants))
	for i := 0; i < len(variants); i++ {
		variants[i] = s.variants[i].ts
	}

	return variants
}

func (t *templateHandler) HasTemplate(name string) bool {
	if _, found := t.baseof[name]; found {
		return true
	}

	if _, found := t.needsBaseof[name]; found {
		return true
	}

	_, found := t.Lookup(name)
	return found
}

func (t *templateHandler) findLayout(d output.LayoutDescriptor, f output.Format) (tpl.Template, bool, error) {
	layouts, _ := t.layoutHandler.For(d, f)
	for _, name := range layouts {
		templ, found := t.main.Lookup(name)
		if found {
			return templ, true, nil
		}

		overlay, found := t.needsBaseof[name]

		if !found {
			continue
		}

		d.Baseof = true
		baseLayouts, _ := t.layoutHandler.For(d, f)
		var base templateInfo
		found = false
		for _, l := range baseLayouts {
			base, found = t.baseof[l]
			if found {
				break
			}
		}

		templ, err := t.applyBaseTemplate(overlay, base)
		if err != nil {
			return nil, false, err
		}

		ts := newTemplateState(templ, overlay)

		if found {
			ts.baseInfo = base

			// Add the base identity to detect changes
			ts.Add(identity.NewPathIdentity(files.ComponentFolderLayouts, base.name))
		}

		t.applyTemplateTransformers(t.main, ts)

		if err := t.extractPartials(ts.Template); err != nil {
			return nil, false, err
		}

		return ts, true, nil

	}

	return nil, false, nil
}

func (t *templateHandler) findTemplate(name string) *templateState {
	if templ, found := t.Lookup(name); found {
		return templ.(*templateState)
	}
	return nil
}

func (t *templateHandler) newTemplateInfo(name, tpl string) templateInfo {
	var isText bool
	name, isText = t.nameIsText(name)
	return templateInfo{
		name:     name,
		isText:   isText,
		template: tpl,
	}
}

func (t *templateHandler) addFileContext(templ tpl.Template, inerr error) error {
	if strings.HasPrefix(templ.Name(), "_internal") {
		return inerr
	}

	ts, ok := templ.(*templateState)
	if !ok {
		return inerr
	}

	//lint:ignore ST1008 the error is the main result
	checkFilename := func(info templateInfo, inErr error) (error, bool) {
		if info.filename == "" {
			return inErr, false
		}

		lineMatcher := func(m herrors.LineMatcher) bool {
			if m.Position.LineNumber != m.LineNumber {
				return false
			}

			identifiers := t.extractIdentifiers(m.Error.Error())

			for _, id := range identifiers {
				if strings.Contains(m.Line, id) {
					return true
				}
			}
			return false
		}

		f, err := t.layoutsFs.Open(info.filename)
		if err != nil {
			return inErr, false
		}
		defer f.Close()

		fe, ok := herrors.WithFileContext(inErr, info.realFilename, f, lineMatcher)
		if ok {
			return fe, true
		}
		return inErr, false
	}

	inerr = errors.Wrap(inerr, "execute of template failed")

	if err, ok := checkFilename(ts.info, inerr); ok {
		return err
	}

	err, _ := checkFilename(ts.baseInfo, inerr)

	return err
}

func (t *templateHandler) addShortcodeVariant(ts *templateState) {
	name := ts.Name()
	base := templateBaseName(templateShortcode, name)

	shortcodename, variants := templateNameAndVariants(base)

	templs, found := t.shortcodes[shortcodename]
	if !found {
		templs = &shortcodeTemplates{}
		t.shortcodes[shortcodename] = templs
	}

	sv := shortcodeVariant{variants: variants, ts: ts}

	i := templs.indexOf(variants)

	if i != -1 {
		// Only replace if it's an override of an internal template.
		if !isInternal(name) {
			templs.variants[i] = sv
		}
	} else {
		templs.variants = append(templs.variants, sv)
	}
}

func (t *templateHandler) addTemplateFile(name, path string) error {
	getTemplate := func(filename string) (templateInfo, error) {
		fs := t.Layouts.Fs
		b, err := afero.ReadFile(fs, filename)
		if err != nil {
			return templateInfo{filename: filename, fs: fs}, err
		}

		s := removeLeadingBOM(string(b))

		realFilename := filename
		if fi, err := fs.Stat(filename); err == nil {
			if fim, ok := fi.(hugofs.FileMetaInfo); ok {
				realFilename = fim.Meta().Filename
			}
		}

		var isText bool
		name, isText = t.nameIsText(name)

		return templateInfo{
			name:         name,
			isText:       isText,
			template:     s,
			filename:     filename,
			realFilename: realFilename,
			fs:           fs,
		}, nil
	}

	tinfo, err := getTemplate(path)
	if err != nil {
		return err
	}

	if isBaseTemplatePath(name) {
		// Store it for later.
		t.baseof[name] = tinfo
		return nil
	}

	needsBaseof := !t.noBaseNeeded(name) && needsBaseTemplate(tinfo.template)
	if needsBaseof {
		t.needsBaseof[name] = tinfo
		return nil
	}

	templ, err := t.addTemplateTo(tinfo, t.main)
	if err != nil {
		return tinfo.errWithFileContext("parse failed", err)
	}
	t.applyTemplateTransformers(t.main, templ)

	return nil
}

func (t *templateHandler) addTemplateTo(info templateInfo, to *templateNamespace) (*templateState, error) {
	return to.parse(info)
}

func (t *templateHandler) applyBaseTemplate(overlay, base templateInfo) (tpl.Template, error) {
	if overlay.isText {
		var (
			templ = t.main.prototypeTextClone.New(overlay.name)
			err   error
		)

		if !base.IsZero() {
			templ, err = templ.Parse(base.template)
			if err != nil {
				return nil, base.errWithFileContext("parse failed", err)
			}
		}

		templ, err = texttemplate.Must(templ.Clone()).Parse(overlay.template)
		if err != nil {
			return nil, overlay.errWithFileContext("parse failed", err)
		}

		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/gohugoio/hugo/issues/2549
		// templ = templ.Lookup(templ.Name())

		return templ, nil
	}

	var (
		templ = t.main.prototypeHTMLClone.New(overlay.name)
		err   error
	)

	if !base.IsZero() {
		templ, err = templ.Parse(base.template)
		if err != nil {
			return nil, base.errWithFileContext("parse failed", err)
		}
	}

	templ, err = htmltemplate.Must(templ.Clone()).Parse(overlay.template)
	if err != nil {
		return nil, overlay.errWithFileContext("parse failed", err)
	}

	// The extra lookup is a workaround, see
	// * https://github.com/golang/go/issues/16101
	// * https://github.com/gohugoio/hugo/issues/2549
	templ = templ.Lookup(templ.Name())

	return templ, err
}

func (t *templateHandler) applyTemplateTransformers(ns *templateNamespace, ts *templateState) (*templateContext, error) {
	c, err := applyTemplateTransformers(ts, ns.newTemplateLookup(ts))
	if err != nil {
		return nil, err
	}

	for k := range c.templateNotFound {
		t.transformNotFound[k] = ts
		t.identityNotFound[k] = append(t.identityNotFound[k], c.t)
	}

	for k := range c.identityNotFound {
		t.identityNotFound[k] = append(t.identityNotFound[k], c.t)
	}

	return c, err
}

func (t *templateHandler) extractIdentifiers(line string) []string {
	m := identifiersRe.FindAllStringSubmatch(line, -1)
	identifiers := make([]string, len(m))
	for i := 0; i < len(m); i++ {
		identifiers[i] = m[i][1]
	}
	return identifiers
}

//go:embed embedded/templates/*
//go:embed embedded/templates/_default/*
var embededTemplatesFs embed.FS

func (t *templateHandler) loadEmbedded() error {
	return fs.WalkDir(embededTemplatesFs, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		templb, err := embededTemplatesFs.ReadFile(path)
		if err != nil {
			return err
		}

		// Get the newlines on Windows in line with how we had it back when we used Go Generate
		// to write the templates to Go files.
		templ := string(bytes.ReplaceAll(templb, []byte("\r\n"), []byte("\n")))
		name := strings.TrimPrefix(filepath.ToSlash(path), "embedded/templates/")
		templateName := name

		// For the render hooks it does not make sense to preseve the
		// double _indternal double book-keeping,
		// just add it if its now provided by the user.
		if !strings.Contains(path, "_default/_markup") {
			templateName = internalPathPrefix + name
		}

		if _, found := t.Lookup(templateName); !found {
			if err := t.AddTemplate(templateName, templ); err != nil {
				return err
			}
		}

		if aliases, found := embeddedTemplatesAliases[name]; found {
			// TODO(bep) avoid reparsing these aliases
			for _, alias := range aliases {
				alias = internalPathPrefix + alias
				if err := t.AddTemplate(alias, templ); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (t *templateHandler) loadTemplates() error {
	walker := func(path string, fi hugofs.FileMetaInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		if isDotFile(path) || isBackupFile(path) {
			return nil
		}

		name := strings.TrimPrefix(filepath.ToSlash(path), "/")
		filename := filepath.Base(path)
		outputFormat, found := t.OutputFormatsConfig.FromFilename(filename)

		if found && outputFormat.IsPlainText {
			name = textTmplNamePrefix + name
		}

		if err := t.addTemplateFile(name, path); err != nil {
			return err
		}

		return nil
	}

	if err := helpers.SymbolicWalk(t.Layouts.Fs, "", walker); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	return nil
}

func (t *templateHandler) nameIsText(name string) (string, bool) {
	isText := strings.HasPrefix(name, textTmplNamePrefix)
	if isText {
		name = strings.TrimPrefix(name, textTmplNamePrefix)
	}
	return name, isText
}

func (t *templateHandler) noBaseNeeded(name string) bool {
	if strings.HasPrefix(name, "shortcodes/") || strings.HasPrefix(name, "partials/") {
		return true
	}
	return strings.Contains(name, "_markup/")
}

func (t *templateHandler) extractPartials(templ tpl.Template) error {
	templs := templates(templ)
	for _, templ := range templs {
		if templ.Name() == "" || !strings.HasPrefix(templ.Name(), "partials/") {
			continue
		}

		ts := newTemplateState(templ, templateInfo{name: templ.Name()})
		ts.typ = templatePartial

		t.main.mu.RLock()
		_, found := t.main.templates[templ.Name()]
		t.main.mu.RUnlock()

		if !found {
			t.main.mu.Lock()
			// This is a template defined inline.
			_, err := applyTemplateTransformers(ts, t.main.newTemplateLookup(ts))
			if err != nil {
				t.main.mu.Unlock()
				return err
			}
			t.main.templates[templ.Name()] = ts
			t.main.mu.Unlock()

		}
	}

	return nil
}

func (t *templateHandler) postTransform() error {
	defineCheckedHTML := false
	defineCheckedText := false

	for _, v := range t.main.templates {
		if v.typ == templateShortcode {
			t.addShortcodeVariant(v)
		}

		if defineCheckedHTML && defineCheckedText {
			continue
		}

		isText := isText(v.Template)
		if isText {
			if defineCheckedText {
				continue
			}
			defineCheckedText = true
		} else {
			if defineCheckedHTML {
				continue
			}
			defineCheckedHTML = true
		}

		if err := t.extractPartials(v.Template); err != nil {
			return err
		}
	}

	for name, source := range t.transformNotFound {
		lookup := t.main.newTemplateLookup(source)
		templ := lookup(name)
		if templ != nil {
			_, err := applyTemplateTransformers(templ, lookup)
			if err != nil {
				return err
			}
		}
	}

	for k, v := range t.identityNotFound {
		ts := t.findTemplate(k)
		if ts != nil {
			for _, im := range v {
				im.Add(ts)
			}
		}
	}

	for _, v := range t.shortcodes {
		sort.Slice(v.variants, func(i, j int) bool {
			v1, v2 := v.variants[i], v.variants[j]
			name1, name2 := v1.ts.Name(), v2.ts.Name()
			isHTMl1, isHTML2 := strings.HasSuffix(name1, "html"), strings.HasSuffix(name2, "html")

			// There will be a weighted selection later, but make
			// sure these are sorted to get a stable selection for
			// output formats missing specific templates.
			// Prefer HTML.
			if isHTMl1 || isHTML2 && !(isHTMl1 && isHTML2) {
				return isHTMl1
			}

			return name1 < name2
		})
	}

	return nil
}

type templateNamespace struct {
	prototypeText      *texttemplate.Template
	prototypeHTML      *htmltemplate.Template
	prototypeTextClone *texttemplate.Template
	prototypeHTMLClone *htmltemplate.Template

	*templateStateMap
}

func (t templateNamespace) Clone() *templateNamespace {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.templateStateMap = &templateStateMap{
		templates: make(map[string]*templateState),
	}

	t.prototypeText = texttemplate.Must(t.prototypeText.Clone())
	t.prototypeHTML = htmltemplate.Must(t.prototypeHTML.Clone())

	return &t
}

func (t *templateNamespace) Lookup(name string) (tpl.Template, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	templ, found := t.templates[name]
	if !found {
		return nil, false
	}

	return templ, found
}

func (t *templateNamespace) createPrototypes() error {
	t.prototypeTextClone = texttemplate.Must(t.prototypeText.Clone())
	t.prototypeHTMLClone = htmltemplate.Must(t.prototypeHTML.Clone())

	return nil
}

func (t *templateNamespace) newTemplateLookup(in *templateState) func(name string) *templateState {
	return func(name string) *templateState {
		if templ, found := t.templates[name]; found {
			if templ.isText() != in.isText() {
				return nil
			}
			return templ
		}
		if templ, found := findTemplateIn(name, in); found {
			return newTemplateState(templ, templateInfo{name: templ.Name()})
		}
		return nil
	}
}

func (t *templateNamespace) parse(info templateInfo) (*templateState, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if info.isText {
		prototype := t.prototypeText

		templ, err := prototype.New(info.name).Parse(info.template)
		if err != nil {
			return nil, err
		}

		ts := newTemplateState(templ, info)

		t.templates[info.name] = ts

		return ts, nil
	}

	prototype := t.prototypeHTML

	templ, err := prototype.New(info.name).Parse(info.template)
	if err != nil {
		return nil, err
	}

	ts := newTemplateState(templ, info)

	t.templates[info.name] = ts

	return ts, nil
}

type templateState struct {
	tpl.Template

	typ       templateType
	parseInfo tpl.ParseInfo
	identity.Manager

	info     templateInfo
	baseInfo templateInfo // Set when a base template is used.
}

func (t *templateState) ParseInfo() tpl.ParseInfo {
	return t.parseInfo
}

func (t *templateState) isText() bool {
	return isText(t.Template)
}

func isText(templ tpl.Template) bool {
	_, isText := templ.(*texttemplate.Template)
	return isText
}

type templateStateMap struct {
	mu        sync.RWMutex
	templates map[string]*templateState
}

type templateWrapperWithLock struct {
	*sync.RWMutex
	tpl.Template
}

type textTemplateWrapperWithLock struct {
	*sync.RWMutex
	*texttemplate.Template
}

func (t *textTemplateWrapperWithLock) Lookup(name string) (tpl.Template, bool) {
	t.RLock()
	templ := t.Template.Lookup(name)
	t.RUnlock()
	if templ == nil {
		return nil, false
	}
	return &textTemplateWrapperWithLock{
		RWMutex:  t.RWMutex,
		Template: templ,
	}, true
}

func (t *textTemplateWrapperWithLock) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	panic("not supported")
}

func (t *textTemplateWrapperWithLock) LookupVariants(name string) []tpl.Template {
	panic("not supported")
}

func (t *textTemplateWrapperWithLock) Parse(name, tpl string) (tpl.Template, error) {
	t.Lock()
	defer t.Unlock()
	return t.Template.New(name).Parse(tpl)
}

func isBackupFile(path string) bool {
	return path[len(path)-1] == '~'
}

func isBaseTemplatePath(path string) bool {
	return strings.Contains(filepath.Base(path), baseFileBase)
}

func isDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}

func removeLeadingBOM(s string) string {
	const bom = '\ufeff'

	for i, r := range s {
		if i == 0 && r != bom {
			return s
		}
		if i > 0 {
			return s[i:]
		}
	}

	return s
}

// resolves _internal/shortcodes/param.html => param.html etc.
func templateBaseName(typ templateType, name string) string {
	name = strings.TrimPrefix(name, internalPathPrefix)
	switch typ {
	case templateShortcode:
		return strings.TrimPrefix(name, shortcodesPathPrefix)
	default:
		panic("not implemented")
	}
}

func unwrap(templ tpl.Template) tpl.Template {
	if ts, ok := templ.(*templateState); ok {
		return ts.Template
	}
	return templ
}

func templates(in tpl.Template) []tpl.Template {
	var templs []tpl.Template
	in = unwrap(in)
	if textt, ok := in.(*texttemplate.Template); ok {
		for _, t := range textt.Templates() {
			templs = append(templs, t)
		}
	}

	if htmlt, ok := in.(*htmltemplate.Template); ok {
		for _, t := range htmlt.Templates() {
			templs = append(templs, t)
		}
	}

	return templs
}
