// Copyright 2016 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eknkc/amber"
	"github.com/spf13/afero"
	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/yosssi/ace"
)

var localTemplates *template.Template
var tmpl Template

// TODO(bep) an interface with hundreds of methods ... remove it.
// And unexport most of these methods.
type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	Lookup(name string) *template.Template
	Templates() []*template.Template
	New(name string) *template.Template
	GetClone() *template.Template
	LoadTemplates(absPath string)
	LoadTemplatesWithPrefix(absPath, prefix string)
	MarkReady()
	AddTemplate(name, tpl string) error
	AddTemplateFileWithMaster(name, overlayFilename, masterFilename string) error
	AddAceTemplate(name, basePath, innerPath string, baseContent, innerContent []byte) error
	AddInternalTemplate(prefix, name, tpl string) error
	AddInternalShortcode(name, tpl string) error
	PrintErrors()
}

type templateErr struct {
	name string
	err  error
}

type GoHTMLTemplate struct {
	template.Template
	clone *template.Template

	// a separate storage for the overlays created from cloned master templates.
	// note: No mutex protection, so we add these in one Go routine, then just read.
	overlays map[string]*template.Template

	errors []*templateErr
}

// T is the "global" template system
func T() Template {
	if tmpl == nil {
		tmpl = New()
	}

	return tmpl
}

// InitializeT resets the internal template state to its initial state
func InitializeT() Template {
	tmpl = New()
	return tmpl
}

// New returns a new Hugo Template System
// with all the additional features, templates & functions
func New() Template {
	var templates = &GoHTMLTemplate{
		Template: *template.New(""),
		overlays: make(map[string]*template.Template),
		errors:   make([]*templateErr, 0),
	}

	localTemplates = &templates.Template

	// The URL funcs in the funcMap is somewhat language dependent,
	// so we need to wait until the language and site config is loaded.
	initFuncMap()

	for k, v := range funcMap {
		amber.FuncMap[k] = v
	}
	templates.Funcs(funcMap)
	templates.LoadEmbedded()
	return templates
}

func partial(name string, contextList ...interface{}) template.HTML {
	if strings.HasPrefix("partials/", name) {
		name = name[8:]
	}
	var context interface{}

	if len(contextList) == 0 {
		context = nil
	} else {
		context = contextList[0]
	}
	return ExecuteTemplateToHTML(context, "partials/"+name, "theme/partials/"+name)
}

func executeTemplate(context interface{}, w io.Writer, layouts ...string) {
	var worked bool
	for _, layout := range layouts {
		templ := Lookup(layout)
		if templ == nil {
			layout += ".html"
			templ = Lookup(layout)
		}

		if templ != nil {
			if err := templ.Execute(w, context); err != nil {
				// Printing the err is spammy, see https://github.com/golang/go/issues/17414
				helpers.DistinctErrorLog.Println(layout, "is an incomplete or empty template")
			}
			worked = true
			break
		}
	}
	if !worked {
		jww.ERROR.Println("Unable to render", layouts)
		jww.ERROR.Println("Expecting to find a template in either the theme/layouts or /layouts in one of the following relative locations", layouts)
	}
}

func ExecuteTemplateToHTML(context interface{}, layouts ...string) template.HTML {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	executeTemplate(context, b, layouts...)
	return template.HTML(b.String())
}

func Lookup(name string) *template.Template {
	return (tmpl.(*GoHTMLTemplate)).Lookup(name)
}

func (t *GoHTMLTemplate) Lookup(name string) *template.Template {

	if templ := localTemplates.Lookup(name); templ != nil {
		return templ
	}

	if t.overlays != nil {
		if templ, ok := t.overlays[name]; ok {
			return templ
		}
	}

	if t.clone != nil {
		if templ := t.clone.Lookup(name); templ != nil {
			return templ
		}
	}

	return nil

}

func (t *GoHTMLTemplate) GetClone() *template.Template {
	return t.clone
}

func (t *GoHTMLTemplate) LoadEmbedded() {
	t.EmbedShortcodes()
	t.EmbedTemplates()
}

// MarkReady marks the template as "ready for execution". No changes allowed
// after this is set.
func (t *GoHTMLTemplate) MarkReady() {
	if t.clone == nil {
		t.clone = template.Must(t.Template.Clone())
	}
}

func (t *GoHTMLTemplate) checkState() {
	if t.clone != nil {
		panic("template is cloned and cannot be modfified")
	}
}

func (t *GoHTMLTemplate) AddInternalTemplate(prefix, name, tpl string) error {
	if prefix != "" {
		return t.AddTemplate("_internal/"+prefix+"/"+name, tpl)
	}
	return t.AddTemplate("_internal/"+name, tpl)
}

func (t *GoHTMLTemplate) AddInternalShortcode(name, content string) error {
	return t.AddInternalTemplate("shortcodes", name, content)
}

func (t *GoHTMLTemplate) AddTemplate(name, tpl string) error {
	t.checkState()
	templ, err := t.New(name).Parse(tpl)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
		return err
	}
	if err := applyTemplateTransformers(templ); err != nil {
		return err
	}

	return nil
}

func (t *GoHTMLTemplate) AddTemplateFileWithMaster(name, overlayFilename, masterFilename string) error {

	// There is currently no known way to associate a cloned template with an existing one.
	// This funky master/overlay design will hopefully improve in a future version of Go.
	//
	// Simplicity is hard.
	//
	// Until then we'll have to live with this hackery.
	//
	// See https://github.com/golang/go/issues/14285
	//
	// So, to do minimum amount of changes to get this to work:
	//
	// 1. Lookup or Parse the master
	// 2. Parse and store the overlay in a separate map

	masterTpl := t.Lookup(masterFilename)

	if masterTpl == nil {
		b, err := afero.ReadFile(hugofs.Source(), masterFilename)
		if err != nil {
			return err
		}
		masterTpl, err = t.New(masterFilename).Parse(string(b))

		if err != nil {
			// TODO(bep) Add a method that does this
			t.errors = append(t.errors, &templateErr{name: name, err: err})
			return err
		}
	}

	b, err := afero.ReadFile(hugofs.Source(), overlayFilename)
	if err != nil {
		return err
	}

	overlayTpl, err := template.Must(masterTpl.Clone()).Parse(string(b))
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	} else {
		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/spf13/hugo/issues/2549
		overlayTpl = overlayTpl.Lookup(overlayTpl.Name())
		if err := applyTemplateTransformers(overlayTpl); err != nil {
			return err
		}
		t.overlays[name] = overlayTpl
	}

	return err
}

func (t *GoHTMLTemplate) AddAceTemplate(name, basePath, innerPath string, baseContent, innerContent []byte) error {
	t.checkState()
	var base, inner *ace.File
	name = name[:len(name)-len(filepath.Ext(innerPath))] + ".html"

	// Fixes issue #1178
	basePath = strings.Replace(basePath, "\\", "/", -1)
	innerPath = strings.Replace(innerPath, "\\", "/", -1)

	if basePath != "" {
		base = ace.NewFile(basePath, baseContent)
		inner = ace.NewFile(innerPath, innerContent)
	} else {
		base = ace.NewFile(innerPath, innerContent)
		inner = ace.NewFile("", []byte{})
	}
	parsed, err := ace.ParseSource(ace.NewSource(base, inner, []*ace.File{}), nil)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
		return err
	}
	templ, err := ace.CompileResultWithTemplate(t.New(name), parsed, nil)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
		return err
	}
	return applyTemplateTransformers(templ)
}

func (t *GoHTMLTemplate) AddTemplateFile(name, baseTemplatePath, path string) error {
	t.checkState()
	// get the suffix and switch on that
	ext := filepath.Ext(path)
	switch ext {
	case ".amber":
		templateName := strings.TrimSuffix(name, filepath.Ext(name)) + ".html"
		compiler := amber.New()
		b, err := afero.ReadFile(hugofs.Source(), path)

		if err != nil {
			return err
		}

		// Parse the input data
		if err := compiler.ParseData(b, path); err != nil {
			return err
		}

		templ, err := compiler.CompileWithTemplate(t.New(templateName))
		if err != nil {
			return err
		}

		return applyTemplateTransformers(templ)
	case ".ace":
		var innerContent, baseContent []byte
		innerContent, err := afero.ReadFile(hugofs.Source(), path)

		if err != nil {
			return err
		}

		if baseTemplatePath != "" {
			baseContent, err = afero.ReadFile(hugofs.Source(), baseTemplatePath)
			if err != nil {
				return err
			}
		}

		return t.AddAceTemplate(name, baseTemplatePath, path, baseContent, innerContent)
	default:

		if baseTemplatePath != "" {
			return t.AddTemplateFileWithMaster(name, path, baseTemplatePath)
		}

		b, err := afero.ReadFile(hugofs.Source(), path)

		if err != nil {
			return err
		}

		jww.DEBUG.Printf("Add template file from path %s", path)

		return t.AddTemplate(name, string(b))
	}

}

func (t *GoHTMLTemplate) GenerateTemplateNameFrom(base, path string) string {
	name, _ := filepath.Rel(base, path)
	return filepath.ToSlash(name)
}

func isDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}

func isBackupFile(path string) bool {
	return path[len(path)-1] == '~'
}

const baseFileBase = "baseof"

var aceTemplateInnerMarkers = [][]byte{[]byte("= content")}
var goTemplateInnerMarkers = [][]byte{[]byte("{{define"), []byte("{{ define")}

func isBaseTemplate(path string) bool {
	return strings.Contains(path, baseFileBase)
}

func (t *GoHTMLTemplate) loadTemplates(absPath string, prefix string) {
	jww.DEBUG.Printf("Load templates from path %q prefix %q", absPath, prefix)
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		jww.DEBUG.Println("Template path", path)
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := filepath.EvalSymlinks(absPath)
			if err != nil {
				jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", absPath, err)
				return nil
			}
			linkfi, err := hugofs.Source().Stat(link)
			if err != nil {
				jww.ERROR.Printf("Cannot stat '%s', error was: %s", link, err)
				return nil
			}
			if !linkfi.Mode().IsRegular() {
				jww.ERROR.Printf("Symbolic links for directories not supported, skipping '%s'", absPath)
			}
			return nil
		}

		if !fi.IsDir() {
			if isDotFile(path) || isBackupFile(path) || isBaseTemplate(path) {
				return nil
			}

			tplName := t.GenerateTemplateNameFrom(absPath, path)

			if prefix != "" {
				tplName = strings.Trim(prefix, "/") + "/" + tplName
			}

			var baseTemplatePath string

			// Ace and Go templates may have both a base and inner template.
			pathDir := filepath.Dir(path)
			if filepath.Ext(path) != ".amber" && !strings.HasSuffix(pathDir, "partials") && !strings.HasSuffix(pathDir, "shortcodes") {

				innerMarkers := goTemplateInnerMarkers
				baseFileName := fmt.Sprintf("%s.html", baseFileBase)

				if filepath.Ext(path) == ".ace" {
					innerMarkers = aceTemplateInnerMarkers
					baseFileName = fmt.Sprintf("%s.ace", baseFileBase)
				}

				// This may be a view that shouldn't have base template
				// Have to look inside it to make sure
				needsBase, err := helpers.FileContainsAny(path, innerMarkers, hugofs.Source())
				if err != nil {
					return err
				}
				if needsBase {

					layoutDir := helpers.GetLayoutDirPath()
					currBaseFilename := fmt.Sprintf("%s-%s", helpers.Filename(path), baseFileName)
					templateDir := filepath.Dir(path)
					themeDir := filepath.Join(helpers.GetThemeDir())
					relativeThemeLayoutsDir := filepath.Join(helpers.GetRelativeThemeDir(), "layouts")

					var baseTemplatedDir string

					if strings.HasPrefix(templateDir, relativeThemeLayoutsDir) {
						baseTemplatedDir = strings.TrimPrefix(templateDir, relativeThemeLayoutsDir)
					} else {
						baseTemplatedDir = strings.TrimPrefix(templateDir, layoutDir)
					}

					baseTemplatedDir = strings.TrimPrefix(baseTemplatedDir, helpers.FilePathSeparator)

					// Look for base template in the follwing order:
					//   1. <current-path>/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
					//   2. <current-path>/baseof.<suffix>
					//   3. _default/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
					//   4. _default/baseof.<suffix>
					// For each of the steps above, it will first look in the project, then, if theme is set,
					// in the theme's layouts folder.

					pairsToCheck := [][]string{
						[]string{baseTemplatedDir, currBaseFilename},
						[]string{baseTemplatedDir, baseFileName},
						[]string{"_default", currBaseFilename},
						[]string{"_default", baseFileName},
					}

				Loop:
					for _, pair := range pairsToCheck {
						pathsToCheck := basePathsToCheck(pair, layoutDir, themeDir)
						for _, pathToCheck := range pathsToCheck {
							if ok, err := helpers.Exists(pathToCheck, hugofs.Source()); err == nil && ok {
								baseTemplatePath = pathToCheck
								break Loop
							}
						}
					}
				}
			}

			if err := t.AddTemplateFile(tplName, baseTemplatePath, path); err != nil {
				jww.ERROR.Printf("Failed to add template %s in path %s: %s", tplName, path, err)
			}

		}
		return nil
	}
	if err := helpers.SymbolicWalk(hugofs.Source(), absPath, walker); err != nil {
		jww.ERROR.Printf("Failed to load templates: %s", err)
	}
}

func basePathsToCheck(path []string, layoutDir, themeDir string) []string {
	// Always look in the project.
	pathsToCheck := []string{filepath.Join((append([]string{layoutDir}, path...))...)}

	// May have a theme
	if themeDir != "" {
		pathsToCheck = append(pathsToCheck, filepath.Join((append([]string{themeDir, "layouts"}, path...))...))
	}

	return pathsToCheck

}

func (t *GoHTMLTemplate) LoadTemplatesWithPrefix(absPath string, prefix string) {
	t.loadTemplates(absPath, prefix)
}

func (t *GoHTMLTemplate) LoadTemplates(absPath string) {
	t.loadTemplates(absPath, "")
}

func (t *GoHTMLTemplate) PrintErrors() {
	for _, e := range t.errors {
		jww.ERROR.Println(e.err)
	}
}
