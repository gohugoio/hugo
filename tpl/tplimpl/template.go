// Copyright 2017-present The Hugo Authors. All rights reserved.
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

	"github.com/eknkc/amber"

	"os"

	"github.com/spf13/hugo/output"

	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/tpl"
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

// Protecting global map access (Amber)
var amberMu sync.Mutex

type templateErr struct {
	name string
	err  error
}

type templateLoader interface {
	handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (string, error)) error
	addTemplate(name, tpl string) error
	addLateTemplate(name, tpl string) error
}

type templateFuncsterTemplater interface {
	tpl.TemplateFinder
	setFuncs(funcMap map[string]interface{})
	setTemplateFuncster(f *templateFuncster)
}

// templateHandler holds the templates in play.
// It implements the templateLoader and tpl.TemplateHandler interfaces.
type templateHandler struct {
	// text holds all the pure text templates.
	text *textTemplates
	html *htmlTemplates

	amberFuncMap template.FuncMap

	errors []*templateErr

	*deps.Deps
}

func (t *templateHandler) addError(name string, err error) {
	t.errors = append(t.errors, &templateErr{name, err})
}

func (t *templateHandler) Debug() {
	fmt.Println("HTML templates:\n", t.html.t.DefinedTemplates())
	fmt.Println("\n\nText templates:\n", t.text.t.DefinedTemplates())
}

// PrintErrors prints the accumulated errors as ERROR to the log.
func (t *templateHandler) PrintErrors() {
	for _, e := range t.errors {
		t.Log.ERROR.Println(e.name, ":", e.err)
	}
}

// Lookup tries to find a template with the given name in both template
// collections: First HTML, then the plain text template collection.
func (t *templateHandler) Lookup(name string) *tpl.TemplateAdapter {

	if strings.HasPrefix(name, textTmplNamePrefix) {
		// The caller has explicitly asked for a text template, so only look
		// in the text template collection.
		// The templates are stored without the prefix identificator.
		name = strings.TrimPrefix(name, textTmplNamePrefix)
		return t.text.Lookup(name)
	}

	// Look in both
	if te := t.html.Lookup(name); te != nil {
		return te
	}

	return t.text.Lookup(name)
}

func (t *templateHandler) clone(d *deps.Deps) *templateHandler {
	c := &templateHandler{
		Deps:   d,
		html:   &htmlTemplates{t: template.Must(t.html.t.Clone()), overlays: make(map[string]*template.Template)},
		text:   &textTemplates{t: texttemplate.Must(t.text.t.Clone()), overlays: make(map[string]*texttemplate.Template)},
		errors: make([]*templateErr, 0),
	}

	d.Tmpl = c

	c.initFuncs()

	for k, v := range t.html.overlays {
		vc := template.Must(v.Clone())
		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/spf13/hugo/issues/2549
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
	htmlT := &htmlTemplates{
		t:        template.New(""),
		overlays: make(map[string]*template.Template),
	}
	textT := &textTemplates{
		t:        texttemplate.New(""),
		overlays: make(map[string]*texttemplate.Template),
	}
	return &templateHandler{
		Deps:   deps,
		html:   htmlT,
		text:   textT,
		errors: make([]*templateErr, 0),
	}

}

type htmlTemplates struct {
	funcster *templateFuncster

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

func (t *htmlTemplates) Lookup(name string) *tpl.TemplateAdapter {
	templ := t.lookup(name)
	if templ == nil {
		return nil
	}
	return &tpl.TemplateAdapter{Template: templ}
}

func (t *htmlTemplates) lookup(name string) *template.Template {
	if templ := t.t.Lookup(name); templ != nil {
		return templ
	}
	if t.overlays != nil {
		if templ, ok := t.overlays[name]; ok {
			return templ
		}
	}

	if t.clone != nil {
		return t.clone.Lookup(name)
	}

	return nil
}

type textTemplates struct {
	funcster *templateFuncster

	t *texttemplate.Template

	clone      *texttemplate.Template
	cloneClone *texttemplate.Template

	overlays map[string]*texttemplate.Template
}

func (t *textTemplates) setTemplateFuncster(f *templateFuncster) {
	t.funcster = f
}

func (t *textTemplates) Lookup(name string) *tpl.TemplateAdapter {
	templ := t.lookup(name)
	if templ == nil {
		return nil
	}
	return &tpl.TemplateAdapter{Template: templ}
}

func (t *textTemplates) lookup(name string) *texttemplate.Template {
	if templ := t.t.Lookup(name); templ != nil {
		return templ
	}
	if t.overlays != nil {
		if templ, ok := t.overlays[name]; ok {
			return templ
		}
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

// LoadTemplates loads the templates, starting from the given absolute path.
// A prefix can be given to indicate a template namespace to load the templates
// into, i.e. "_internal" etc.
func (t *templateHandler) LoadTemplates(absPath, prefix string) {
	t.loadTemplates(absPath, prefix)

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
		tt.AddParseTree(withoutExt, templ.Tree)
	}

	return nil
}

func (t *htmlTemplates) addTemplate(name, tpl string) error {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *htmlTemplates) addLateTemplate(name, tpl string) error {
	return t.addTemplateIn(t.clone, name, tpl)
}

func (t *textTemplates) addTemplateIn(tt *texttemplate.Template, name, tpl string) error {
	name = strings.TrimPrefix(name, textTmplNamePrefix)
	templ, err := tt.New(name).Parse(tpl)
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
		tt.AddParseTree(withoutExt, templ.Tree)
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
		t.addError(name, err)
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
		t.addError(name, err)
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
	t.html.clone = template.Must(t.html.cloneClone.Clone())
	t.text.clone = texttemplate.Must(t.text.cloneClone.Clone())
}

func (t *templateHandler) loadTemplates(absPath string, prefix string) {
	t.Log.DEBUG.Printf("Load templates from path %q prefix %q", absPath, prefix)
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		t.Log.DEBUG.Println("Template path", path)
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := filepath.EvalSymlinks(absPath)
			if err != nil {
				t.Log.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", absPath, err)
				return nil
			}

			linkfi, err := t.Fs.Source.Stat(link)
			if err != nil {
				t.Log.ERROR.Printf("Cannot stat '%s', error was: %s", link, err)
				return nil
			}

			if !linkfi.Mode().IsRegular() {
				t.Log.ERROR.Printf("Symbolic links for directories not supported, skipping '%s'", absPath)
			}
			return nil
		}

		if !fi.IsDir() {
			if isDotFile(path) || isBackupFile(path) || isBaseTemplate(path) {
				return nil
			}

			var (
				workingDir = t.PathSpec.WorkingDir()
				themeDir   = t.PathSpec.GetThemeDir()
				layoutDir  = t.PathSpec.LayoutDir()
			)

			if themeDir != "" && strings.HasPrefix(absPath, themeDir) {
				workingDir = themeDir
				layoutDir = "layouts"
			}

			li := strings.LastIndex(path, layoutDir) + len(layoutDir) + 1
			relPath := path[li:]
			templateDir := path[:li-len(layoutDir)-1]

			descriptor := output.TemplateLookupDescriptor{
				TemplateDir:   templateDir,
				WorkingDir:    workingDir,
				LayoutDir:     layoutDir,
				RelPath:       relPath,
				Prefix:        prefix,
				ThemeDir:      themeDir,
				OutputFormats: t.OutputFormatsConfig,
				FileExists: func(filename string) (bool, error) {
					return helpers.Exists(filename, t.Fs.Source)
				},
				ContainsAny: func(filename string, subslices [][]byte) (bool, error) {
					return helpers.FileContainsAny(filename, subslices, t.Fs.Source)
				},
			}

			tplID, err := output.CreateTemplateNames(descriptor)
			if err != nil {
				t.Log.ERROR.Printf("Failed to resolve template in path %q: %s", path, err)

				return nil
			}

			if err := t.addTemplateFile(tplID.Name, tplID.MasterFilename, tplID.OverlayFilename); err != nil {
				t.Log.ERROR.Printf("Failed to add template %q in path %q: %s", tplID.Name, path, err)
			}

		}
		return nil
	}
	if err := helpers.SymbolicWalk(t.Fs.Source, absPath, walker); err != nil {
		t.Log.ERROR.Printf("Failed to load templates: %s", err)
	}
}

func (t *templateHandler) initFuncs() {

	// Both template types will get their own funcster instance, which
	// in the current case contains the same set of funcs.
	for _, funcsterHolder := range []templateFuncsterTemplater{t.html, t.text} {
		funcster := newTemplateFuncster(t.Deps)

		// The URL funcs in the funcMap is somewhat language dependent,
		// so we need to wait until the language and site config is loaded.
		funcster.initFuncMap()

		funcsterHolder.setTemplateFuncster(funcster)

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

func (t *templateHandler) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (string, error)) error {
	h := t.getTemplateHandler(name)
	return h.handleMaster(name, overlayFilename, masterFilename, onMissing)
}

func (t *htmlTemplates) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (string, error)) error {
	masterTpl := t.lookup(masterFilename)

	if masterTpl == nil {
		templ, err := onMissing(masterFilename)
		if err != nil {
			return err
		}

		masterTpl, err = t.t.New(overlayFilename).Parse(templ)
		if err != nil {
			return err
		}
	}

	templ, err := onMissing(overlayFilename)
	if err != nil {
		return err
	}

	overlayTpl, err := template.Must(masterTpl.Clone()).Parse(templ)
	if err != nil {
		return err
	}

	// The extra lookup is a workaround, see
	// * https://github.com/golang/go/issues/16101
	// * https://github.com/spf13/hugo/issues/2549
	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())
	if err := applyTemplateTransformersToHMLTTemplate(overlayTpl); err != nil {
		return err
	}
	t.overlays[name] = overlayTpl

	return err

}

func (t *textTemplates) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (string, error)) error {
	name = strings.TrimPrefix(name, textTmplNamePrefix)
	masterTpl := t.lookup(masterFilename)

	if masterTpl == nil {
		templ, err := onMissing(masterFilename)
		if err != nil {
			return err
		}

		masterTpl, err = t.t.New(overlayFilename).Parse(templ)
		if err != nil {
			return err
		}
	}

	templ, err := onMissing(overlayFilename)
	if err != nil {
		return err
	}

	overlayTpl, err := texttemplate.Must(masterTpl.Clone()).Parse(templ)
	if err != nil {
		return err
	}

	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())
	if err := applyTemplateTransformersToTextTemplate(overlayTpl); err != nil {
		return err
	}
	t.overlays[name] = overlayTpl

	return err

}

func (t *templateHandler) addTemplateFile(name, baseTemplatePath, path string) error {
	t.checkState()

	getTemplate := func(filename string) (string, error) {
		b, err := afero.ReadFile(t.Fs.Source, filename)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	// get the suffix and switch on that
	ext := filepath.Ext(path)
	switch ext {
	case ".amber":
		//	Only HTML support for Amber
		templateName := strings.TrimSuffix(name, filepath.Ext(name)) + ".html"
		b, err := afero.ReadFile(t.Fs.Source, path)

		if err != nil {
			return err
		}

		amberMu.Lock()
		templ, err := t.compileAmberWithTemplate(b, path, t.html.t.New(templateName))
		amberMu.Unlock()
		if err != nil {
			return err
		}

		return applyTemplateTransformersToHMLTTemplate(templ)
	case ".ace":
		//	Only HTML support for Ace
		var innerContent, baseContent []byte
		innerContent, err := afero.ReadFile(t.Fs.Source, path)

		if err != nil {
			return err
		}

		if baseTemplatePath != "" {
			baseContent, err = afero.ReadFile(t.Fs.Source, baseTemplatePath)
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

		t.Log.DEBUG.Printf("Add template file from path %s", path)

		return t.AddTemplate(name, templ)
	}

}

func (t *templateHandler) loadEmbedded() {
	t.embedShortcodes()
	t.embedTemplates()
}

func (t *templateHandler) addInternalTemplate(prefix, name, tpl string) error {
	if prefix != "" {
		return t.AddTemplate("_internal/"+prefix+"/"+name, tpl)
	}
	return t.AddTemplate("_internal/"+name, tpl)
}

func (t *templateHandler) addInternalShortcode(name, content string) error {
	return t.addInternalTemplate("shortcodes", name, content)
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
	return strings.Contains(path, baseFileBase)
}
