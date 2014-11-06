package target

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
)

const ALIAS = "<!DOCTYPE html><html><head><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0;url={{ .Permalink }}\" /></head></html>"
const ALIAS_XHTML = "<!DOCTYPE html><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0;url={{ .Permalink }}\" /></head></html>"

var DefaultAliasTemplates *template.Template

func init() {
	DefaultAliasTemplates = template.New("")
	template.Must(DefaultAliasTemplates.New("alias").Parse(ALIAS))
	template.Must(DefaultAliasTemplates.New("alias-xhtml").Parse(ALIAS_XHTML))
}

type AliasPublisher interface {
	Translator
	Publish(string, template.HTML) error
}

type HTMLRedirectAlias struct {
	PublishDir string
	Templates  *template.Template
}

func (h *HTMLRedirectAlias) Translate(alias string) (aliasPath string, err error) {
	if len(alias) <= 0 {
		return
	}

	if strings.HasSuffix(alias, "/") {
		alias = alias + "index.html"
	} else if !strings.HasSuffix(alias, ".html") {
		alias = alias + "/index.html"
	}
	return filepath.Join(h.PublishDir, helpers.MakePath(alias)), nil
}

type AliasNode struct {
	Permalink template.HTML
}

func (h *HTMLRedirectAlias) Publish(path string, permalink template.HTML) (err error) {
	if path, err = h.Translate(path); err != nil {
		return
	}

	t := "alias"
	if strings.HasSuffix(path, ".xhtml") {
		t = "alias-xhtml"
	}

	template := DefaultAliasTemplates
	if h.Templates != nil {
		template = h.Templates
	}

	buffer := new(bytes.Buffer)
	err = template.ExecuteTemplate(buffer, t, &AliasNode{permalink})
	if err != nil {
		return
	}

	return helpers.WriteToDisk(path, buffer, hugofs.DestinationFS)
}
