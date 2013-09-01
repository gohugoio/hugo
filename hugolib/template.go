package hugolib

import (
	"io/ioutil"
	"github.com/eknkc/amber"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// HTML encapsulates a known safe HTML document fragment.
// It should not be used for HTML from a third-party, or HTML with
// unclosed tags or comments. The outputs of a sound HTML sanitizer
// and a template escaped by this package are fine for use with HTML.
type HTML template.HTML

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	Lookup(name string) *template.Template
	Templates() []*template.Template
	New(name string) *template.Template
	LoadTemplates(absPath string)
	AddTemplate(name, tpl string) error
}

type URL template.URL

type templateErr struct {
	name string
	err  error
}

type GoHtmlTemplate struct {
	template.Template
	errors []*templateErr
}

func NewTemplate() Template {
	var templates = &GoHtmlTemplate{
		Template: *template.New(""),
		errors:   make([]*templateErr, 0),
	}

	funcMap := template.FuncMap{
		"urlize":    Urlize,
		"gt":        Gt,
		"isset":     IsSet,
		"echoParam": ReturnWhenSet,
	}

	templates.Funcs(funcMap)
	templates.primeTemplates()
	return templates
}

func (t *GoHtmlTemplate) AddTemplate(name, tpl string) error {
	_, err := t.New(name).Parse(tpl)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	}
	return err
}

func (t *GoHtmlTemplate) AddTemplateFile(name, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	_, err = t.New(name).Parse(s)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	}
	return err
}

func (t *GoHtmlTemplate) generateTemplateNameFrom(base, path string) string {
	return filepath.ToSlash(path[len(base)+1:])
}

func (t *GoHtmlTemplate) primeTemplates() {
	alias := "<!DOCTYPE html>\n <html>\n <head>\n <link rel=\"canonical\" href=\"{{ .Permalink }}\"/>\n <meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" />\n <meta http-equiv=\"refresh\" content=\"0;url={{ .Permalink }}\" />\n </head>\n </html>"
	alias_xhtml := "<!DOCTYPE html>\n <html xmlns=\"http://www.w3.org/1999/xhtml\">\n <head>\n <link rel=\"canonical\" href=\"{{ .Permalink }}\"/>\n <meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" />\n <meta http-equiv=\"refresh\" content=\"0;url={{ .Permalink }}\" />\n </head>\n </html>"

	t.AddTemplate("alias", alias)
	t.AddTemplate("alias-xhtml", alias_xhtml)
}

func (t *GoHtmlTemplate) LoadTemplates(absPath string) {
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			PrintErr("Walker: ", err)
			return nil
		}

		if !fi.IsDir() {
			if ignoreDotFile(path) {
				return nil
			}

			tplName := t.generateTemplateNameFrom(absPath, path)

			if strings.HasSuffix(path, ".amber") {
				compiler := amber.New()
				// Parse the input file
				if err := compiler.ParseFile(path); err != nil {
					return nil
				}

				// note t.New(tplName)
				if _, err := compiler.CompileWithTemplate(t.New(tplName)); err != nil {
					PrintErr("Could not compile amber file: "+path, err)
					return err
				}

			} else {
				t.AddTemplateFile(tplName, path)
			}
		}
		return nil
	}

	filepath.Walk(absPath, walker)
}
