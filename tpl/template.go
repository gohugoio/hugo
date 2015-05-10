// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"bytes"
	"fmt"
	"github.com/eknkc/amber"
	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/yosssi/ace"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var localTemplates *template.Template
var tmpl Template

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	Lookup(name string) *template.Template
	Templates() []*template.Template
	New(name string) *template.Template
	LoadTemplates(absPath string)
	LoadTemplatesWithPrefix(absPath, prefix string)
	AddTemplate(name, tpl string) error
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
	errors []*templateErr
}

// The "Global" Template System
func T() Template {
	if tmpl == nil {
		tmpl = New()
	}

	return tmpl
}

// Resets the internal template state to it's initial state
func InitializeT() Template {
	tmpl = New()
	return tmpl
}

// Return a new Hugo Template System
// With all the additional features, templates & functions
func New() Template {
	var templates = &GoHTMLTemplate{
		Template: *template.New(""),
		errors:   make([]*templateErr, 0),
	}

	localTemplates = &templates.Template

	templates.Funcs(funcMap)
	templates.LoadEmbedded()
	return templates
}

func Partial(name string, context_list ...interface{}) template.HTML {
	if strings.HasPrefix("partials/", name) {
		name = name[8:]
	}
	var context interface{}

	if len(context_list) == 0 {
		context = nil
	} else {
		context = context_list[0]
	}
	return ExecuteTemplateToHTML(context, "partials/"+name, "theme/partials/"+name)
}

func ExecuteTemplate(context interface{}, buffer *bytes.Buffer, layouts ...string) {
	worked := false
	for _, layout := range layouts {

		name := layout

		if localTemplates.Lookup(name) == nil {
			name = layout + ".html"
		}

		if localTemplates.Lookup(name) != nil {
			err := localTemplates.ExecuteTemplate(buffer, name, context)
			if err != nil {
				jww.ERROR.Println(err, "in", name)
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
	ExecuteTemplate(context, b, layouts...)
	return template.HTML(b.String())
}

func (t *GoHTMLTemplate) LoadEmbedded() {
	t.EmbedShortcodes()
	t.EmbedTemplates()
}

func (t *GoHTMLTemplate) AddInternalTemplate(prefix, name, tpl string) error {
	if prefix != "" {
		return t.AddTemplate("_internal/"+prefix+"/"+name, tpl)
	} else {
		return t.AddTemplate("_internal/"+name, tpl)
	}
}

func (t *GoHTMLTemplate) AddInternalShortcode(name, content string) error {
	return t.AddInternalTemplate("shortcodes", name, content)
}

func (t *GoHTMLTemplate) AddTemplate(name, tpl string) error {
	_, err := t.New(name).Parse(tpl)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	}
	return err
}

func (t *GoHTMLTemplate) AddTemplateFile(name, baseTemplatePath, path string) error {
	// get the suffix and switch on that
	ext := filepath.Ext(path)
	switch ext {
	case ".amber":
		compiler := amber.New()
		// Parse the input file
		if err := compiler.ParseFile(path); err != nil {
			return nil
		}

		if _, err := compiler.CompileWithTemplate(t.New(name)); err != nil {
			return err
		}
	case ".ace":
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var base, inner *ace.File

		name = name[:len(name)-len(ext)] + ".html"
		if baseTemplatePath != "" {
			b2, err := ioutil.ReadFile(baseTemplatePath)
			if err != nil {
				return err
			}
			base = ace.NewFile(baseTemplatePath, b2)
			inner = ace.NewFile(path, b)
		} else {
			base = ace.NewFile(path, b)
			inner = ace.NewFile("", []byte{})
		}
		rslt, err := ace.ParseSource(ace.NewSource(base, inner, []*ace.File{}), nil)
		if err != nil {
			t.errors = append(t.errors, &templateErr{name: name, err: err})
			return err
		}
		_, err = ace.CompileResultWithTemplate(t.New(name), rslt, nil)
		if err != nil {
			t.errors = append(t.errors, &templateErr{name: name, err: err})
		}
		return err
	default:
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return t.AddTemplate(name, string(b))
	}

	return nil

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

const baseAceFilename = "baseof.ace"

var aceTemplateInnerMarker = []byte("= content")

func isBaseTemplate(path string) bool {
	return strings.HasSuffix(path, baseAceFilename)
}

func (t *GoHTMLTemplate) loadTemplates(absPath string, prefix string) {
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := filepath.EvalSymlinks(absPath)
			if err != nil {
				jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", absPath, err)
				return nil
			}
			linkfi, err := os.Stat(link)
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

			// ACE templates may have both a base and inner template.
			if filepath.Ext(path) == ".ace" && !strings.HasSuffix(filepath.Dir(path), "partials") {
				// This may be a view that shouldn't have base template
				// Have to look inside it to make sure
				needsBase, err := helpers.FileContains(path, aceTemplateInnerMarker, hugofs.OsFs)
				if err != nil {
					return err
				}
				if needsBase {

					// Look for base template in the follwing order:
					//   1. <current-path>/<template-name>-baseof.ace, e.g. list-baseof.ace.
					//   2. <current-path>/baseof.ace
					//   3. _default/<template-name>-baseof.ace, e.g. list-baseof.ace.
					//   4. _default/baseof.ace

					currBaseAceFilename := fmt.Sprintf("%s-%s", helpers.Filename(path), baseAceFilename)
					templateDir := filepath.Dir(path)

					pathsToCheck := []string{
						filepath.Join(templateDir, currBaseAceFilename),
						filepath.Join(templateDir, baseAceFilename),
						filepath.Join(absPath, "_default", currBaseAceFilename),
						filepath.Join(absPath, "_default", baseAceFilename)}

					for _, pathToCheck := range pathsToCheck {
						if ok, err := helpers.Exists(pathToCheck, hugofs.OsFs); err == nil && ok {
							baseTemplatePath = pathToCheck
							break
						}
					}
				}
			}

			t.AddTemplateFile(tplName, baseTemplatePath, path)

		}
		return nil
	}

	filepath.Walk(absPath, walker)
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
