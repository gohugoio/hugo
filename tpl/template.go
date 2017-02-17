package tpl

import (
	"html/template"
	"io"
)

// TODO(bep) make smaller
type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	ExecuteTemplateToHTML(context interface{}, layouts ...string) template.HTML
	Lookup(name string) *template.Template
	Templates() []*template.Template
	New(name string) *template.Template
	GetClone() *template.Template
	LoadTemplates(absPath string)
	LoadTemplatesWithPrefix(absPath, prefix string)
	AddTemplate(name, tpl string) error
	AddTemplateFileWithMaster(name, overlayFilename, masterFilename string) error
	AddAceTemplate(name, basePath, innerPath string, baseContent, innerContent []byte) error
	AddInternalTemplate(prefix, name, tpl string) error
	AddInternalShortcode(name, tpl string) error
	Partial(name string, contextList ...interface{}) template.HTML
	PrintErrors()
	Funcs(funcMap template.FuncMap)
	MarkReady()
}
