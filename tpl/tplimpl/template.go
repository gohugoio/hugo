// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"path"
	"strings"
	texttemplate "text/template"

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
	_ templateLoader            = (*templateHandler)(nil)
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
	addTemplate(name, tpl string) error
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

// NewTextTemplate provides a text template parser that has all the Hugo
// template funcs etc. built-in.
func (t *templateHandler) NewTextTemplate() tpl.TemplateParseFinder {
	t.mu.Lock()
	t.mu.Unlock()

	tt := &textTemplate{t: texttemplate.New("")}
	t.extTextTemplates = append(t.extTextTemplates, tt)

	return tt

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

		return t.text.Lookup(name)
	}

	// Look in both
	if te, found := t.html.Lookup(name); found {
		return te, true
	}

	return t.text.Lookup(name)

}

func (t *templateHandler) clone(d *deps.Deps) *templateHandler {
	c := &templateHandler{
		Deps:      d,
		layoutsFs: d.BaseFs.Layouts.Fs,
		html:      &htmlTemplates{t: template.Must(t.html.t.Clone()), overlays: make(map[string]*template.Template), templatesCommon: t.html.templatesCommon},
		text:      &textTemplates{textTemplate: &textTemplate{t: texttemplate.Must(t.text.t.Clone())}, overlays: make(map[string]*texttemplate.Template), templatesCommon: t.text.templatesCommon},
		errors:    make([]*templateErr, 0),
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
		Deps:      deps,
		layoutsFs: deps.BaseFs.Layouts.Fs,
		html:      htmlT,
		text:      textT,
		errors:    make([]*templateErr, 0),
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
}
type htmlTemplates struct {
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

func (t *htmlTemplates) addTemplateIn(tt *template.Template, name, tpl string) error {
	templ, err := tt.New(name).Parse(tpl)
	if err != nil {
		return err
	}

	if err := applyTemplateTransformersToHMLTTemplate(templ); err != nil {
		return err
	}

	if strings.Contains(name, "shortcodes") {
		// We need to keep track of one ot the output format's shortcode template
		// without knowing the rendering context.
		withoutExt := strings.TrimSuffix(name, path.Ext(name))
		clone := template.Must(templ.Clone())
		tt.AddParseTree(withoutExt, clone.Tree)
	}

	return nil
}

func (t *htmlTemplates) addTemplate(name, tpl string) error {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *htmlTemplates) addLateTemplate(name, tpl string) error {
	return t.addTemplateIn(t.clone, name, tpl)
}

type textTemplate struct {
	t *texttemplate.Template
}

func (t *textTemplate) Parse(name, tpl string) (tpl.Template, error) {
	return t.parSeIn(t.t, name, tpl)
}

func (t *textTemplate) Lookup(name string) (tpl.Template, bool) {
	tpl := t.t.Lookup(name)
	return tpl, tpl != nil
}

func (t *textTemplate) parSeIn(tt *texttemplate.Template, name, tpl string) (*texttemplate.Template, error) {
	templ, err := tt.New(name).Parse(tpl)
	if err != nil {
		return nil, err
	}

	if err := applyTemplateTransformersToTextTemplate(templ); err != nil {
		return nil, err
	}
	return templ, nil
}

func (t *textTemplates) addTemplateIn(tt *texttemplate.Template, name, tpl string) error {
	name = strings.TrimPrefix(name, textTmplNamePrefix)
	templ, err := t.parSeIn(tt, name, tpl)
	if err != nil {
		return err
	}

	if err := applyTemplateTransformersToTextTemplate(templ); err != nil {
		return err
	}

	if strings.Contains(name, "shortcodes") {
		// We need to keep track of one ot the output format's shortcode template
		// without knowing the rendering context.
		withoutExt := strings.TrimSuffix(name, path.Ext(name))
		clone := texttemplate.Must(templ.Clone())
		tt.AddParseTree(withoutExt, clone.Tree)
	}

	return nil
}

func (t *textTemplates) addTemplate(name, tpl string) error {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *textTemplates) addLateTemplate(name, tpl string) error {
	return t.addTemplateIn(t.clone, name, tpl)
}

func (t *templateHandler) addTemplate(name, tpl string) error {
	return t.AddTemplate(name, tpl)
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
func (t *templateHandler) AddTemplate(name, tpl string) error {
	h := t.getTemplateHandler(name)
	if err := h.addTemplate(name, tpl); err != nil {
		return err
	}
	return nil
}

// MarkReady marks the templates as "ready for execution". No changes allowed
// after this is set.
// TODO(bep) if this proves to be resource heavy, we could detect
// earlier if we really need this, or make it lazy.
func (t *templateHandler) MarkReady() {
	if t.html.clone == nil {
		t.html.clone = template.Must(t.html.t.Clone())
		t.html.cloneClone = template.Must(t.html.clone.Clone())
	}
	if t.text.clone == nil {
		t.text.clone = texttemplate.Must(t.text.t.Clone())
		t.text.cloneClone = texttemplate.Must(t.text.clone.Clone())
	}
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

	walker := func(path string, fi os.FileInfo, err error) error {
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
	if err := applyTemplateTransformersToHMLTTemplate(overlayTpl); err != nil {
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
	if err := applyTemplateTransformersToTextTemplate(overlayTpl); err != nil {
		return err
	}
	t.overlays[name] = overlayTpl
	t.nameBaseTemplateName[name] = templ.filename

	return err

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
		s := string(b)

		realFilename := filename
		if fi, err := fs.Stat(filename); err == nil {
			if fir, ok := fi.(hugofs.RealFilenameInfo); ok {
				realFilename = fir.RealFilename()
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

		if err := applyTemplateTransformersToHMLTTemplate(templ); err != nil {
			return err
		}

		if strings.Contains(templateName, "shortcodes") {
			// We need to keep track of one ot the output format's shortcode template
			// without knowing the rendering context.
			clone := template.Must(templ.Clone())
			t.html.t.AddParseTree(withoutExt, clone.Tree)
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
