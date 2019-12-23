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
	"io"
	"reflect"
	"regexp"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/common/herrors"

	"strings"

	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"

	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/tpl/tplimpl/embedded"
	"github.com/pkg/errors"

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
	_ tpl.TemplateManager    = (*templateHandler)(nil)
	_ tpl.TemplateHandler    = (*templateHandler)(nil)
	_ tpl.TemplateDebugger   = (*templateHandler)(nil)
	_ tpl.TemplateFuncGetter = (*templateHandler)(nil)
	_ tpl.TemplateFinder     = (*htmlTemplates)(nil)
	_ tpl.TemplateFinder     = (*textTemplates)(nil)
	_ templateLoader         = (*htmlTemplates)(nil)
	_ templateLoader         = (*textTemplates)(nil)
)

const (
	shortcodesPathPrefix = "shortcodes/"
	internalPathPrefix   = "_internal/"
)

// The identifiers may be truncated in the log, e.g.
// "executing "main" at <$scaled.SRelPermalin...>: can't evaluate field SRelPermalink in type *resource.Image"
var identifiersRe = regexp.MustCompile(`at \<(.*?)(\.{3})?\>:`)

var embeddedTemplatesAliases = map[string][]string{
	"shortcodes/twitter.html": {"shortcodes/tweet.html"},
}

const baseFileBase = "baseof"

func newTemplateAdapter(deps *deps.Deps) *templateHandler {

	common := &templatesCommon{
		nameBaseTemplateName: make(map[string]string),
		transformNotFound:    make(map[string]bool),
		identityNotFound:     make(map[string][]identity.Manager),
	}

	htmlT := &htmlTemplates{
		t:               template.New(""),
		overlays:        make(map[string]*template.Template),
		templatesCommon: common,
	}

	textT := &textTemplates{
		textTemplate:    &textTemplate{t: texttemplate.New("")},
		standalone:      &textTemplate{t: texttemplate.New("")},
		overlays:        make(map[string]*texttemplate.Template),
		templatesCommon: common,
	}

	h := &templateHandler{
		Deps:      deps,
		layoutsFs: deps.BaseFs.Layouts.Fs,
		templateHandlerCommon: &templateHandlerCommon{
			shortcodes:       make(map[string]*shortcodeTemplates),
			templateInfo:     make(map[string]tpl.Info),
			templateInfoTree: make(map[string]*templateInfoTree),
			html:             htmlT,
			text:             textT,
		},
	}

	textT.textTemplate.templates = textT
	textT.standalone.templates = textT
	common.handler = h

	return h

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

func (t *htmlTemplates) Lookup(name string) (tpl.Template, bool) {
	templ := t.lookup(name)
	if templ == nil {
		return nil, false
	}

	return templ, true
}

func (t *htmlTemplates) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return t.handler.LookupVariant(name, variants)
}

func (t *htmlTemplates) addLateTemplate(name, tpl string) error {
	_, err := t.addTemplateIn(t.clone, name, tpl)
	return err
}

func (t *htmlTemplates) addTemplate(name, tpl string) (*templateContext, error) {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *htmlTemplates) addTemplateIn(tt *template.Template, name, templstr string) (*templateContext, error) {
	templ, err := tt.New(name).Parse(templstr)
	if err != nil {
		return nil, err
	}

	typ := resolveTemplateType(name)

	c, err := t.handler.applyTemplateTransformersToHMLTTemplate(typ, templ)
	if err != nil {
		return nil, err
	}

	for k := range c.templateNotFound {
		t.transformNotFound[k] = true
		t.identityNotFound[k] = append(t.identityNotFound[k], c.id)
	}

	for k := range c.identityNotFound {
		t.identityNotFound[k] = append(t.identityNotFound[k], c.id)
	}

	return c, nil
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
	if _, err := t.handler.applyTemplateTransformersToHMLTTemplate(templateUndefined, overlayTpl); err != nil {
		return err
	}

	t.overlays[name] = overlayTpl
	t.nameBaseTemplateName[name] = masterFilename

	return err

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

func (t htmlTemplates) withNewHandler(h *templateHandler) *htmlTemplates {
	t.templatesCommon = t.templatesCommon.withNewHandler(h)
	return &t
}

type nopLookupVariant int

func (l nopLookupVariant) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return nil, false, false
}

// templateHandler holds the templates in play.
// It implements the templateLoader and tpl.TemplateHandler interfaces.
// There is one templateHandler created per Site.
type templateHandler struct {
	ready bool

	executor texttemplate.Executer
	funcs    map[string]reflect.Value

	// This is the filesystem to load the templates from. All the templates are
	// stored in the root of this filesystem.
	layoutsFs afero.Fs

	*deps.Deps

	*templateHandlerCommon
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

func (t *templateHandler) Debug() {
	fmt.Println("HTML templates:\n", t.html.t.DefinedTemplates())
	fmt.Println("\n\nText templates:\n", t.text.t.DefinedTemplates())
}

func (t *templateHandler) Execute(templ tpl.Template, wr io.Writer, data interface{}) error {
	if t.Metrics != nil {
		defer t.Metrics.MeasureSince(templ.Name(), time.Now())
	}

	execErr := t.executor.Execute(templ, wr, data)
	if execErr != nil {
		execErr = t.addFileContext(templ.Name(), execErr)
	}

	return execErr

}

func (t *templateHandler) GetFunc(name string) (reflect.Value, bool) {
	v, found := t.funcs[name]
	return v, found

}

// LoadTemplates loads the templates from the layouts filesystem.
// A prefix can be given to indicate a template namespace to load the templates
// into, i.e. "_internal" etc.
func (t *templateHandler) LoadTemplates(prefix string) error {
	return t.loadTemplates(prefix)

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

// This currently only applies to shortcodes and what we get here is the
// shortcode name.
func (t *templateHandler) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	if !t.ready {
		panic("handler not ready")
	}
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

	return tpl.WithInfo(sv.templ, sv.info), true, more

}

// markReady marks the templates as "ready for execution". No changes allowed
// after this is set.
func (t *templateHandler) markReady() error {
	defer func() {
		t.ready = true
	}()

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

func (h *templateHandler) initTemplateExecuter() {
	exec, funcs := newTemplateExecuter(h.Deps)
	h.executor = exec
	h.funcs = funcs
	funcMap := make(map[string]interface{})
	for k, v := range funcs {
		funcMap[k] = v.Interface()
	}

	// Note that these funcs are not the ones getting called
	// on execution, but they are needed at parse time.
	h.text.textTemplate.t.Funcs(funcMap)
	h.text.standalone.t.Funcs(funcMap)
	h.html.t.Funcs(funcMap)
}

func (t *templateHandler) getTemplateHandler(name string) templateLoader {
	if strings.HasPrefix(name, textTmplNamePrefix) {
		return t.text
	}
	return t.html
}

func (t *templateHandler) addFileContext(name string, inerr error) error {
	if strings.HasPrefix(name, "_internal") {
		return inerr
	}

	f, realFilename, err := t.fileAndFilename(name)
	if err != nil {
		return inerr

	}
	defer f.Close()

	master, hasMaster := t.html.nameBaseTemplateName[name]

	ferr := errors.Wrap(inerr, "execute of template failed")

	// Since this can be a composite of multiple template files (single.html + baseof.html etc.)
	// we potentially need to look in both -- and cannot rely on line number alone.
	lineMatcher := func(m herrors.LineMatcher) bool {
		if m.Position.LineNumber != m.LineNumber {
			return false
		}
		if !hasMaster {
			return true
		}

		identifiers := t.extractIdentifiers(m.Error.Error())

		for _, id := range identifiers {
			if strings.Contains(m.Line, id) {
				return true
			}
		}
		return false
	}

	fe, ok := herrors.WithFileContext(ferr, realFilename, f, lineMatcher)
	if ok || !hasMaster {
		return fe
	}

	// Try the base template if relevant
	f, realFilename, err = t.fileAndFilename(master)
	if err != nil {
		return err
	}
	defer f.Close()

	fe, ok = herrors.WithFileContext(ferr, realFilename, f, lineMatcher)

	if !ok {
		// Return the most specific.
		return ferr

	}
	return fe

}

func (t *templateHandler) addInternalTemplate(name, tpl string) error {
	return t.AddTemplate("_internal/"+name, tpl)
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
		helpers.Deprecated("Amber templates are no longer supported.", "Use Go templates or a Hugo version <= 0.60.", true)
		return nil
	case ".ace":
		helpers.Deprecated("ACE templates are no longer supported.", "Use Go templates or a Hugo version <= 0.60.", true)
		return nil
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

func (t *templateHandler) applyTemplateInfo(templ tpl.Template, found bool) (tpl.Template, bool) {
	if templ != nil {
		if info, found := t.templateInfo[templ.Name()]; found {
			return tpl.WithInfo(templ, info), true
		}
	}

	return templ, found
}

func (t *templateHandler) checkState() {
	if t.html.clone != nil || t.text.clone != nil {
		panic("template is cloned and cannot be modfified")
	}
}

func (t *templateHandler) clone(d *deps.Deps) *templateHandler {
	if !t.ready {
		panic("invalid state")
	}
	c := &templateHandler{
		ready:     true,
		Deps:      d,
		layoutsFs: d.BaseFs.Layouts.Fs,
	}

	c.templateHandlerCommon = t.templateHandlerCommon.withNewHandler(c)
	d.Tmpl = c
	d.TextTmpl = c.wrapTextTemplate(c.text.standalone)
	c.executor, c.funcs = newTemplateExecuter(d)

	return c

}

func (t *templateHandler) extractIdentifiers(line string) []string {
	m := identifiersRe.FindAllStringSubmatch(line, -1)
	identifiers := make([]string, len(m))
	for i := 0; i < len(m); i++ {
		identifiers[i] = m[i][1]
	}
	return identifiers
}

func (t *templateHandler) fileAndFilename(name string) (afero.File, string, error) {
	fs := t.layoutsFs
	filename := filepath.FromSlash(name)

	fi, err := fs.Stat(filename)
	if err != nil {
		return nil, "", err
	}
	fim := fi.(hugofs.FileMetaInfo)
	meta := fim.Meta()

	f, err := meta.Open()
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to open template file %q:", filename)
	}

	return f, meta.Filename(), nil
}

func (t *templateHandler) handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error {
	h := t.getTemplateHandler(name)
	return h.handleMaster(name, overlayFilename, masterFilename, onMissing)
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

func (t *templateHandler) getOrCreateTemplateInfo(name string) (identity.Manager, tpl.ParseInfo) {
	info, found := t.templateInfo[name]
	if found {
		return info.(identity.Manager), info.ParseInfo()
	}
	return identity.NewManager(identity.NewPathIdentity(files.ComponentFolderLayouts, name)), tpl.DefaultParseInfo
}

func (t *templateHandler) createTemplateInfo(name string) (identity.Manager, tpl.ParseInfo) {
	_, found := t.templateInfo[name]
	if found {
		panic("already created: " + name)
	}

	return identity.NewManager(identity.NewPathIdentity(files.ComponentFolderLayouts, name)), tpl.DefaultParseInfo
}

func (t *templateHandler) postTransform() error {
	for k, v := range t.templateInfoTree {
		if v.id != nil {
			info := tpl.NewInfo(
				v.id,
				v.info,
			)
			t.templateInfo[k] = info

			if v.typ == templateShortcode {
				t.addShortcodeVariant(k, info, v.templ)
			}
		}
	}

	for _, s := range []struct {
		lookup            func(name string) *templateInfoTree
		transformNotFound map[string]bool
		identityNotFound  map[string][]identity.Manager
	}{
		// html templates
		{func(name string) *templateInfoTree {
			templ := t.html.lookup(name)
			if templ == nil {
				return nil
			}
			id, info := t.getOrCreateTemplateInfo(name)
			return &templateInfoTree{
				id:   id,
				info: info,
				tree: templ.Tree,
			}
		}, t.html.transformNotFound, t.html.identityNotFound},
		// text templates
		{func(name string) *templateInfoTree {
			templT := t.text.lookup(name)
			if templT == nil {
				return nil
			}
			id, info := t.getOrCreateTemplateInfo(name)
			return &templateInfoTree{
				id:   id,
				info: info,
				tree: templT.Tree,
			}
		}, t.text.transformNotFound, t.text.identityNotFound},
	} {
		for name := range s.transformNotFound {
			templ := s.lookup(name)
			if templ != nil {
				_, err := applyTemplateTransformers(templateUndefined, templ, s.lookup)
				if err != nil {
					return err
				}
			}
		}

		for k, v := range s.identityNotFound {
			tmpl := s.lookup(k)
			if tmpl != nil {
				for _, im := range v {
					im.Add(tmpl.id)
				}
			}
		}
	}

	return nil
}

func (t *templateHandler) wrapTextTemplate(tt *textTemplate) tpl.TemplateParseFinder {
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

type templateHandlerCommon struct {
	// shortcodes maps shortcode name to template variants
	// (language, output format etc.) of that shortcode.
	shortcodes map[string]*shortcodeTemplates

	// templateInfo maps template name to some additional information about that template.
	// Note that for shortcodes that same information is embedded in the
	// shortcodeTemplates type.
	templateInfo map[string]tpl.Info

	// Used to track templates during the AST transformations.
	templateInfoTree map[string]*templateInfoTree

	// text holds all the pure text templates.
	text *textTemplates
	html *htmlTemplates
}

func (t templateHandlerCommon) withNewHandler(h *templateHandler) *templateHandlerCommon {
	t.text = t.text.withNewHandler(h)
	t.html = t.html.withNewHandler(h)
	return &t
}

type templateLoader interface {
	addLateTemplate(name, tpl string) error
	addTemplate(name, tpl string) (*templateContext, error)
	handleMaster(name, overlayFilename, masterFilename string, onMissing func(filename string) (templateInfo, error)) error
}

// Shared by both HTML and text templates.
type templatesCommon struct {
	handler *templateHandler

	// Used to get proper filenames in errors
	nameBaseTemplateName map[string]string

	// Holds names of the template definitions not found during the first AST transformation
	// pass.
	transformNotFound map[string]bool

	// Holds identities of templates not found during first pass.
	identityNotFound map[string][]identity.Manager
}

func (t templatesCommon) withNewHandler(h *templateHandler) *templatesCommon {
	t.handler = h
	return &t
}

type textTemplate struct {
	mu        sync.RWMutex
	t         *texttemplate.Template
	templates *textTemplates
}

func (t *textTemplate) Lookup(name string) (tpl.Template, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tpl := t.t.Lookup(name)
	return tpl, tpl != nil
}

func (t *textTemplate) Parse(name, tpl string) (tpl.Template, error) {
	return t.parseIn(t.t, name, tpl)
}

func (t *textTemplate) parseIn(tt *texttemplate.Template, name, tpl string) (*texttemplate.Template, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	templ, err := tt.New(name).Parse(tpl)
	if err != nil {
		return nil, err
	}

	if _, err := t.templates.handler.applyTemplateTransformersToTextTemplate(templateUndefined, templ); err != nil {
		return nil, err
	}
	return templ, nil
}

type textTemplates struct {
	*templatesCommon
	*textTemplate
	standalone *textTemplate
	clone      *texttemplate.Template
	cloneClone *texttemplate.Template

	overlays map[string]*texttemplate.Template
}

func (t *textTemplates) Lookup(name string) (tpl.Template, bool) {
	templ := t.lookup(name)
	if templ == nil {
		return nil, false
	}
	return templ, true
}

func (t *textTemplates) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return t.handler.LookupVariant(name, variants)
}

func (t *textTemplates) addLateTemplate(name, tpl string) error {
	_, err := t.addTemplateIn(t.clone, name, tpl)
	return err
}

func (t *textTemplates) addTemplate(name, tpl string) (*templateContext, error) {
	return t.addTemplateIn(t.t, name, tpl)
}

func (t *textTemplates) addTemplateIn(tt *texttemplate.Template, name, tplstr string) (*templateContext, error) {
	name = strings.TrimPrefix(name, textTmplNamePrefix)
	templ, err := t.parseIn(tt, name, tplstr)
	if err != nil {
		return nil, err
	}

	typ := resolveTemplateType(name)

	c, err := t.handler.applyTemplateTransformersToTextTemplate(typ, templ)
	if err != nil {
		return nil, err
	}

	for k := range c.templateNotFound {
		t.transformNotFound[k] = true
	}

	return c, nil
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
	if _, err := t.handler.applyTemplateTransformersToTextTemplate(templateUndefined, overlayTpl); err != nil {
		return err
	}
	t.overlays[name] = overlayTpl
	t.nameBaseTemplateName[name] = templ.filename

	return err

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

func (t textTemplates) withNewHandler(h *templateHandler) *textTemplates {
	t.templatesCommon = t.templatesCommon.withNewHandler(h)
	return &t
}

func isBackupFile(path string) bool {
	return path[len(path)-1] == '~'
}

func isBaseTemplate(path string) bool {
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
