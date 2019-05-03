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
	"fmt"
	"html/template"
	"strings"
	texttemplate "text/template"
	"text/template/parse"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/tpl/tplimpl/embedded"
	"github.com/pkg/errors"

	"github.com/eknkc/amber"

	"os"

	"github.com/gohugoio/hugo/output"

	"path/filepath"
	"sync"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/afero"
)

const (
	textTmplNamePrefix = "_text/"
)

var (
	_ tpl.TemplateHandler       = (*templateHandler)(nil)
	_ tpl.TemplateDebugger      = (*templateHandler)(nil)
	_ tpl.TemplateFuncsGetter   = (*templateHandler)(nil)
	_ tpl.TemplateTestMocker    = (*templateHandler)(nil)
	_ tpl.TemplateFinder        = (*htmlTemplates)(nil)
	_ tpl.TemplateFinder        = (*textTemplates)(nil)
	_ templateLoader            = (*htmlTemplates)(nil)
	_ templateLoader            = (*textTemplates)(nil)
	_ templateFuncsterTemplater = (*htmlTemplates)(nil)
	_ templateFuncsterTemplater = (*textTemplates)(nil)
)

// Protecting  global map access (Amber)
var amberMu sync.Mutex

type templateErr struct {
	name string
	err  error
}

type templateLoader interface {
	handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error
	addTemplate(name, tpl string) (*templateContext, error)
	addLateTemplate(name, tpl string) error
}

type templateFuncsterTemplater interface {
	templateFuncsterSetter
	tpl.TemplateFinder
	setFuncs(funcMap map[string]interface{})
}

type templateFuncsterSetter interface {
	setTemplateFuncster(f *templateFuncster)
}

// templateHandler holds the templates in play.
// It implements the templateLoader and tpl.TemplateHandler interfaces.
type templateHandler struct {
	mu sync.Mutex

	// shortcodes maps shortcode name to template variants
	// (language, output format etc.) of that shortcode.
	shortcodes map[string]*shortcodeTemplates

	// templateInfo maps template name to some additional information about that template.
	// Note that for shortcodes that same information is embedded in the
	// shortcodeTemplates type.
	templateInfo map[string]tpl.Info

	// text holds all the pure text templates.
	text *textTemplates
	html *htmlTemplates

	extTextTemplates []*textTemplate

	amberFuncMap template.FuncMap

	errors []*templateErr

	// This is the filesystem to load the templates from. All the templates are
	// stored in the root of this filesystem.
	layoutsFs afero.Fs

	*deps.Deps
}

const (
	shortcodesPathPrefix = "shortcodes/"
	internalPathPrefix   = "_internal/"
)

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

func (t *templateHandler) addShortcodeVariant(name string, info tpl.Info, templ tpl.Template) {
	base := templateBaseName(templateShortcode, name)

	shortcodename, variants := templateNameAndVariants(base)

	templs, found := t.shortcodes[shortcodename]
	if !found {
		templs = &shortcodeTemplates{}
		t.shortcodes[shortcodename] = templs
	}

	sv := shortcodeVariant{variants: variants, info: info, templ: templ}

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

// NewTextTemplate provides a text template parser that has all the Hugo
// template funcs etc. built-in.
func (t *templateHandler) NewTextTemplate() tpl.TemplateParseFinder {
	t.mu.Lock()
	defer t.mu.Unlock()

	tt := &textTemplate{t: texttemplate.New("")}
	t.extTextTemplates = append(t.extTextTemplates, tt)

	return struct {
		tpl.TemplateParser
		tpl.TemplateLookup
		tpl.TemplateLookupVariant
	}{
		tt,
		tt,
		new(nopLookupVariant),
	}

}

type nopLookupVariant int

func (l nopLookupVariant) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return nil, false, false
}

func (t *templateHandler) Debug() {
	fmt.Println("HTML templates:\n", t.html.t.DefinedTemplates())
	fmt.Println("\n\nText templates:\n", t.text.t.DefinedTemplates())
}

// Lookup tries to find a template with the given name in both template
// collections: First HTML, then the plain text template collection.
func (t *templateHandler) Lookup(name string) (tpl.Template, bool) {

	if strings.HasPrefix(name, textTmplNamePrefix) {
		// The caller has explicitly asked for a text template, so only look
		// in the text template collection.
		// The templates are stored without the prefix identificator.
		name = strings.TrimPrefix(name, textTmplNamePrefix)

		return t.applyTemplateInfo(t.text.Lookup(name))
	}

	// Look in both
	if te, found := t.html.Lookup(name); found {
		return t.applyTemplateInfo(te, true)
	}

	return t.applyTemplateInfo(t.text.Lookup(name))

}

func (t *templateHandler) applyTemplateInfo(templ tpl.Template, found bool) (tpl.Template, bool) {
	if adapter, ok := templ.(*tpl.TemplateAdapter); ok {
		if adapter.Info.IsZero() {
			if info, found := t.templateInfo[templ.Name()]; found {
				adapter.Info = info
			}
		}
	}

	return templ, found
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

	return &tpl.TemplateAdapter{
		Template:             sv.templ,
		Info:                 sv.info,
		Metrics:              t.Deps.Metrics,
		Fs:                   t.layoutsFs,
		NameBaseTemplateName: t.html.nameBaseTemplateName}, true, more

}

func (t *textTemplates) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return t.handler.LookupVariant(name, variants)
}

func (t *htmlTemplates) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return t.handler.LookupVariant(name, variants)
}

func (t *templateHandler) lookupTemplate(in interface{}) tpl.Template {
	switch templ := in.(type) {
	case *texttemplate.Template:
		return t.text.lookup(templ.Name())
	case *template.Template:
		return t.html.lookup(templ.Name())
	}

	panic(fmt.Sprintf("%T is not a template", in))
}

func (t *templateHandler) setFuncMapInTemplate(in interface{}, funcs map[string]interface{}) {
	switch templ := in.(type) {
	case *texttemplate.Template:
		templ.Funcs(funcs)
		return
	case *template.Template:
		templ.Funcs(funcs)
		return
	}

	panic(fmt.Sprintf("%T is not a template", in))
}

func (t *templateHandler) clone(d *deps.Deps) *templateHandler {
	c := &templateHandler{
		Deps:         d,
		layoutsFs:    d.BaseFs.Layouts.Fs,
		shortcodes:   make(map[string]*shortcodeTemplates),
		templateInfo: t.templateInfo,
		html:         &htmlTemplates{t: template.Must(t.html.t.Clone()), overlays: make(map[string]*template.Template), templatesCommon: t.html.templatesCommon},
		text:         &textTemplates{textTemplate: &textTemplate{t: texttemplate.Must(t.text.t.Clone())}, overlays: make(map[string]*texttemplate.Template), templatesCommon: t.text.templatesCommon},
		errors:       make([]*templateErr, 0),
	}

	for k, v := range t.shortcodes {
		other := *v
		variantsc := make([]shortcodeVariant, len(v.variants))
		for i, variant := range v.variants {
			variantsc[i] = shortcodeVariant{
				info:     variant.info,
				variants: variant.variants,
				templ:    c.lookupTemplate(variant.templ),
			}
		}
		other.variants = variantsc
		c.shortcodes[k] = &other
	}

	d.Tmpl = c

	c.initFuncs()

	for k, v := range t.html.overlays {
		vc := template.Must(v.Clone())
		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/gohugoio/hugo/issues/2549
		vc = vc.Lookup(vc.Name())
		vc.Funcs(c.html.funcster.funcMap)
		c.html.overlays[k] = vc
	}

	for k, v := range t.text.overlays {
		vc := texttemplate.Must(v.Clone())
		vc = vc.Lookup(vc.Name())
		vc.Funcs(texttemplate.FuncMap(c.text.funcster.funcMap))
		c.text.overlays[k] = vc
	}

	return c

}

func newTemplateAdapter(deps *deps.Deps) *templateHandler {
	common := &templatesCommon{
		nameBaseTemplateName: make(map[string]string),
		transformNotFound:    make(map[string]bool),
	}

	htmlT := &htmlTemplates{
		t:               template.New(""),
		overlays:        make(map[string]*template.Template),
		templatesCommon: common,
	}
	textT := &textTemplates{
		textTemplate:    &textTemplate{t: texttemplate.New("")},
		overlays:        make(map[string]*texttemplate.Template),
		templatesCommon: common,
	}
	h := &templateHandler{
		Deps:         deps,
		layoutsFs:    deps.BaseFs.Layouts.Fs,
		shortcodes:   make(map[string]*shortcodeTemplates),
		templateInfo: make(map[string]tpl.Info),
		html:         htmlT,
		text:         textT,
		errors:       make([]*templateErr, 0),
	}

	common.handler = h

	return h

}

// Shared by both HTML and text templates.
type templatesCommon struct {
	handler  *templateHandler
	funcster *templateFuncster

	// Used to get proper filenames in errors
	nameBaseTemplateName map[string]string

	// Holds names of the templates not found during the first AST transformation
	// pass.
	transformNotFound map[string]bool
}
type htmlTemplates struct {
	mu sync.RWMutex

	*templatesCommon

	t *template.Template

	// This looks, and is, strange.
	// The clone is used by non-renderable content pages, and these need to be
	// re-parsed on content change, and to avoid the
	// "cannot Parse after Execute" error, we need to re-clone it from the original clone.
	clone      *template.Template
	cloneClone *template.Template

	// a separate storage for the overlays created from cloned master templates.
	// note: No mutex protection, so we add these in one Go routine, then just read.
	overlays map[string]*template.Template
}

func (t *htmlTemplates) setTemplateFuncster(f *templateFuncster) {
	t.funcster = f
}

func (t *htmlTemplates) Lookup(name string) (tpl.Template, bool) {
	templ := t.lookup(name)
	if templ == nil {
		return nil, false
	}

	return &tpl.TemplateAdapter{Template: templ, Metrics: t.funcster.Deps.Metrics, Fs: t.handler.layoutsFs, NameBaseTemplateName: t.nameBaseTemplateName}, true
}

func (t *htmlTemplates) lookup(name string) *template.Template {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Need to check in the overlay registry first as it will also be found below.
	if t.overlays != nil {
		if templ, ok := t.overlays[name]; ok {
			return templ
		}
	}

	if templ := t.t.Lookup(name); templ != nil {
		return templ
	}

	if t.clone != nil {
		return t.clone.Lookup(name)
	}

	return nil
}

func (t *textTemplates) setTemplateFuncster(f *templateFuncster) {
	t.funcster = f
}

type textTemplates struct {
	*templatesCommon
	*textTemplate
	clone      *texttemplate.Template
	cloneClone *texttemplate.Template

	overlays map[string]*texttemplate.Template
}

func (t *textTemplates) Lookup(name string) (tpl.Template, bool) {
	templ := t.lookup(name)
	if templ == nil {
		return nil, false
	}
	return &tpl.TemplateAdapter{Template: templ, Metrics: t.funcster.Deps.Metrics, Fs: t.handler.layoutsFs, NameBaseTemplateName: t.nameBaseTemplateName}, true
}

func (t *textTemplates) lookup(name string) *texttemplate.Template {

	// Need to check in the overlay registry first as it will also be found below.
	if t.overlays != nil {
		if templ, ok := t.overlays[name]; ok {
			return templ
		}
	}

	if templ := t.t.Lookup(name); templ != nil {
		return templ
	}

	if t.clone != nil {
		return t.clone.Lookup(name)
	}

	return nil
}

func (t *templateHandler) setFuncs(funcMap map[string]interface{}) {
	t.html.setFuncs(funcMap)
	t.text.setFuncs(funcMap)
}

// SetFuncs replaces the funcs in the func maps with new definitions.
// This is only used in tests.
func (t *templateHandler) SetFuncs(funcMap map[string]interface{}) {
	t.setFuncs(funcMap)
}

func (t *templateHandler) GetFuncs() map[string]interface{} {
	return t.html.funcster.funcMap
}

func (t *htmlTemplates) setFuncs(funcMap map[string]interface{}) {
	t.t.Funcs(funcMap)
}

func (t *textTemplates) setFuncs(funcMap map[string]interface{}) {
	t.t.Funcs(funcMap)
}

// LoadTemplates loads the templates from the layouts filesystem.
// A prefix can be given to indicate a template namespace to load the templates
// into, i.e. "_internal" etc.
func (t *templateHandler) LoadTemplates(prefix string) error {
	return t.loadTemplates(prefix)

}

func (t *htmlTemplates) addTemplateIn(tt *template.Template, name, tpl string) (*templateContext, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	templ, err := tt.New(name).Parse(tpl)
	if err != nil {
		return nil, err
	}

	typ := resolveTemplateType(name)

	c, err := applyTemplateTransformersToHMLTTemplate(typ, templ)
	if err != nil {
		return nil, err
	}

	for k, _ := range c.notFound {
		t.transformNotFound[k] = true
	}

	if typ == templateShortcode {
		t.handler.addShortcodeVariant(name, c.Info, templ)
	} else {
		t.handler.templateInfo[name] = c.Info
	}

	return c, nil
}

func (t *htmlTemplates) addTemplate(name, tpl string) (*templateContext, error) {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *htmlTemplates) addLateTemplate(name, tpl string) error {
	_, err := t.addTemplateIn(t.clone, name, tpl)
	return err
}

type textTemplate struct {
	mu sync.RWMutex
	t  *texttemplate.Template
}

func (t *textTemplate) Parse(name, tpl string) (tpl.Template, error) {
	return t.parseIn(t.t, name, tpl)
}

func (t *textTemplate) Lookup(name string) (tpl.Template, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tpl := t.t.Lookup(name)
	return tpl, tpl != nil
}

func (t *textTemplate) parseIn(tt *texttemplate.Template, name, tpl string) (*texttemplate.Template, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	templ, err := tt.New(name).Parse(tpl)
	if err != nil {
		return nil, err
	}

	if _, err := applyTemplateTransformersToTextTemplate(templateUndefined, templ); err != nil {
		return nil, err
	}
	return templ, nil
}

func (t *textTemplates) addTemplateIn(tt *texttemplate.Template, name, tpl string) (*templateContext, error) {
	name = strings.TrimPrefix(name, textTmplNamePrefix)
	templ, err := t.parseIn(tt, name, tpl)
	if err != nil {
		return nil, err
	}

	typ := resolveTemplateType(name)

	c, err := applyTemplateTransformersToTextTemplate(typ, templ)
	if err != nil {
		return nil, err
	}

	for k, _ := range c.notFound {
		t.transformNotFound[k] = true
	}

	if typ == templateShortcode {
		t.handler.addShortcodeVariant(name, c.Info, templ)
	} else {
		t.handler.templateInfo[name] = c.Info
	}

	return c, nil
}

func (t *textTemplates) addTemplate(name, tpl string) (*templateContext, error) {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *textTemplates) addLateTemplate(name, tpl string) error {
	_, err := t.addTemplateIn(t.clone, name, tpl)
	return err
}

func (t *templateHandler) addTemplate(name, tpl string) error {
	return t.AddTemplate(name, tpl)
}

func (t *templateHandler) postTransform() error {
	if len(t.html.transformNotFound) == 0 && len(t.text.transformNotFound) == 0 {
		return nil
	}

	defer func() {
		t.text.transformNotFound = make(map[string]bool)
		t.html.transformNotFound = make(map[string]bool)
	}()

	for _, s := range []struct {
		lookup            func(name string) *parse.Tree
		transformNotFound map[string]bool
	}{
		// html templates
		{func(name string) *parse.Tree {
			templ := t.html.lookup(name)
			if templ == nil {
				return nil
			}
			return templ.Tree
		}, t.html.transformNotFound},
		// text templates
		{func(name string) *parse.Tree {
			templT := t.text.lookup(name)
			if templT == nil {
				return nil
			}
			return templT.Tree
		}, t.text.transformNotFound},
	} {
		for name, _ := range s.transformNotFound {
			templ := s.lookup(name)
			if templ != nil {
				_, err := applyTemplateTransformers(templateUndefined, templ, s.lookup)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *templateHandler) addLateTemplate(name, tpl string) error {
	return t.AddLateTemplate(name, tpl)
}

// AddLateTemplate is used to add a template late, i.e. after the
// regular templates have started its execution.
func (t *templateHandler) AddLateTemplate(name, tpl string) error {
	h := t.getTemplateHandler(name)
	if err := h.addLateTemplate(name, tpl); err != nil {
		return err
	}
	return nil
}

// AddTemplate parses and adds a template to the collection.
// Templates with name prefixed with "_text" will be handled as plain
// text templates.
// TODO(bep) clean up these addTemplate variants
func (t *templateHandler) AddTemplate(name, tpl string) error {
	h := t.getTemplateHandler(name)
	_, err := h.addTemplate(name, tpl)
	if err != nil {
		return err
	}
	return nil
}

// MarkReady marks the templates as "ready for execution". No changes allowed
// after this is set.
// TODO(bep) if this proves to be resource heavy, we could detect
// earlier if we really need this, or make it lazy.
func (t *templateHandler) MarkReady() error {
	if err := t.postTransform(); err != nil {
		return err
	}

	if t.html.clone == nil {
		t.html.clone = template.Must(t.html.t.Clone())
		t.html.cloneClone = template.Must(t.html.clone.Clone())
	}
	if t.text.clone == nil {
		t.text.clone = texttemplate.Must(t.text.t.Clone())
		t.text.cloneClone = texttemplate.Must(t.text.clone.Clone())
	}

	return nil
}

// RebuildClone rebuilds the cloned templates. Used for live-reloads.
func (t *templateHandler) RebuildClone() {
	if t.html != nil && t.html.cloneClone != nil {
		t.html.clone = template.Must(t.html.cloneClone.Clone())
	}
	if t.text != nil && t.text.cloneClone != nil {
		t.text.clone = texttemplate.Must(t.text.cloneClone.Clone())
	}
}

func (t *templateHandler) loadTemplates(prefix string) error {

	walker := func(path string, fi hugofs.FileMetaInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		if isDotFile(path) || isBackupFile(path) || isBaseTemplate(path) {
			return nil
		}

		workingDir := t.PathSpec.WorkingDir

		descriptor := output.TemplateLookupDescriptor{
			WorkingDir:    workingDir,
			RelPath:       path,
			Prefix:        prefix,
			OutputFormats: t.OutputFormatsConfig,
			FileExists: func(filename string) (bool, error) {
				return helpers.Exists(filename, t.Layouts.Fs)
			},
			ContainsAny: func(filename string, subslices [][]byte) (bool, error) {
				return helpers.FileContainsAny(filename, subslices, t.Layouts.Fs)
			},
		}

		tplID, err := output.CreateTemplateNames(descriptor)
		if err != nil {
			t.Log.ERROR.Printf("Failed to resolve template in path %q: %s", path, err)
			return nil
		}

		if err := t.addTemplateFile(tplID.Name, tplID.MasterFilename, tplID.OverlayFilename); err != nil {
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

func (t *templateHandler) initFuncs() {

	// Both template types will get their own funcster instance, which
	// in the current case contains the same set of funcs.
	funcMap := createFuncMap(t.Deps)
	for _, funcsterHolder := range []templateFuncsterSetter{t.html, t.text} {
		funcster := newTemplateFuncster(t.Deps)

		// The URL funcs in the funcMap is somewhat language dependent,
		// so we need to wait until the language and site config is loaded.
		funcster.initFuncMap(funcMap)

		funcsterHolder.setTemplateFuncster(funcster)

	}

	for _, v := range t.shortcodes {
		for _, variant := range v.variants {
			t.setFuncMapInTemplate(variant.templ, funcMap)
		}
	}

	for _, extText := range t.extTextTemplates {
		extText.t.Funcs(funcMap)
	}

	// Amber is HTML only.
	t.amberFuncMap = template.FuncMap{}

	amberMu.Lock()
	for k, v := range amber.FuncMap {
		t.amberFuncMap[k] = v
	}

	for k, v := range t.html.funcster.funcMap {
		t.amberFuncMap[k] = v
		// Hacky, but we need to make sure that the func names are in the global map.
		amber.FuncMap[k] = func() string {
			panic("should never be invoked")
		}
	}
	amberMu.Unlock()

}

func (t *templateHandler) getTemplateHandler(name string) templateLoader {
	if strings.HasPrefix(name, textTmplNamePrefix) {
		return t.text
	}
	return t.html
}

func (t *templateHandler) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error {
	h := t.getTemplateHandler(name)
	return h.handleMaster(name, overlayFilename, masterFilename, onMissing)
}

func (t *htmlTemplates) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error {

	masterTpl := t.lookup(masterFilename)

	if masterTpl == nil {
		templ, err := onMissing(masterFilename)
		if err != nil {
			return err
		}

		masterTpl, err = t.t.New(overlayFilename).Parse(templ.template)
		if err != nil {
			return templ.errWithFileContext("parse master failed", err)
		}
	}

	templ, err := onMissing(overlayFilename)
	if err != nil {
		return err
	}

	overlayTpl, err := template.Must(masterTpl.Clone()).Parse(templ.template)
	if err != nil {
		return templ.errWithFileContext("parse failed", err)
	}

	// The extra lookup is a workaround, see
	// * https://github.com/golang/go/issues/16101
	// * https://github.com/gohugoio/hugo/issues/2549
	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())
	if _, err := applyTemplateTransformersToHMLTTemplate(templateUndefined, overlayTpl); err != nil {
		return err
	}

	t.overlays[name] = overlayTpl
	t.nameBaseTemplateName[name] = masterFilename

	return err

}

func (t *textTemplates) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error {

	name = strings.TrimPrefix(name, textTmplNamePrefix)
	masterTpl := t.lookup(masterFilename)

	if masterTpl == nil {
		templ, err := onMissing(masterFilename)
		if err != nil {
			return err
		}

		masterTpl, err = t.t.New(masterFilename).Parse(templ.template)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %q:", templ.filename)
		}
		t.nameBaseTemplateName[masterFilename] = templ.filename
	}

	templ, err := onMissing(overlayFilename)
	if err != nil {
		return err
	}

	overlayTpl, err := texttemplate.Must(masterTpl.Clone()).Parse(templ.template)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %q:", templ.filename)
	}

	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())
	if _, err := applyTemplateTransformersToTextTemplate(templateUndefined, overlayTpl); err != nil {
		return err
	}
	t.overlays[name] = overlayTpl
	t.nameBaseTemplateName[name] = templ.filename

	return err

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

func (t *templateHandler) addTemplateFile(name, baseTemplatePath, path string) error {
	t.checkState()

	t.Log.DEBUG.Printf("Add template file: name %q, baseTemplatePath %q, path %q", name, baseTemplatePath, path)

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
				realFilename = fim.Meta().Filename()
			}
		}

		return templateInfo{template: s, filename: filename, realFilename: realFilename, fs: fs}, nil
	}

	// get the suffix and switch on that
	ext := filepath.Ext(path)
	switch ext {
	case ".amber":
		//	Only HTML support for Amber
		withoutExt := strings.TrimSuffix(name, filepath.Ext(name))
		templateName := withoutExt + ".html"
		b, err := afero.ReadFile(t.Layouts.Fs, path)

		if err != nil {
			return err
		}

		amberMu.Lock()
		templ, err := t.compileAmberWithTemplate(b, path, t.html.t.New(templateName))
		amberMu.Unlock()
		if err != nil {
			return err
		}

		typ := resolveTemplateType(name)

		c, err := applyTemplateTransformersToHMLTTemplate(typ, templ)
		if err != nil {
			return err
		}

		if typ == templateShortcode {
			t.addShortcodeVariant(templateName, c.Info, templ)
		} else {
			t.templateInfo[name] = c.Info
		}

		return nil

	case ".ace":
		//	Only HTML support for Ace
		var innerContent, baseContent []byte
		innerContent, err := afero.ReadFile(t.Layouts.Fs, path)

		if err != nil {
			return err
		}

		if baseTemplatePath != "" {
			baseContent, err = afero.ReadFile(t.Layouts.Fs, baseTemplatePath)
			if err != nil {
				return err
			}
		}

		return t.addAceTemplate(name, baseTemplatePath, path, baseContent, innerContent)
	default:

		if baseTemplatePath != "" {
			return t.handleMaster(name, path, baseTemplatePath, getTemplate)
		}

		templ, err := getTemplate(path)

		if err != nil {
			return err
		}

		err = t.AddTemplate(name, templ.template)
		if err != nil {
			return templ.errWithFileContext("parse failed", err)
		}
		return nil
	}
}

var embeddedTemplatesAliases = map[string][]string{
	"shortcodes/twitter.html": {"shortcodes/tweet.html"},
}

func (t *templateHandler) loadEmbedded() error {
	for _, kv := range embedded.EmbeddedTemplates {
		name, templ := kv[0], kv[1]
		if err := t.addInternalTemplate(name, templ); err != nil {
			return err
		}
		if aliases, found := embeddedTemplatesAliases[name]; found {
			for _, alias := range aliases {
				if err := t.addInternalTemplate(alias, templ); err != nil {
					return err
				}
			}

		}
	}

	return nil

}

func (t *templateHandler) addInternalTemplate(name, tpl string) error {
	return t.AddTemplate("_internal/"+name, tpl)
}

func (t *templateHandler) checkState() {
	if t.html.clone != nil || t.text.clone != nil {
		panic("template is cloned and cannot be modfified")
	}
}

func isDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}

func isBackupFile(path string) bool {
	return path[len(path)-1] == '~'
}

const baseFileBase = "baseof"

func isBaseTemplate(path string) bool {
	return strings.Contains(filepath.Base(path), baseFileBase)
}
