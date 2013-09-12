package target

import (
	helpers "github.com/spf13/hugo/template"
	"path"
	"strings"
)

type HTMLRedirectAlias struct {
	PublishDir string
}

func (h *HTMLRedirectAlias) Translate(alias string) (aliasPath string, err error) {
	if strings.HasSuffix(alias, "/") {
		alias = alias + "index.html"
	}
	return path.Join(h.PublishDir, helpers.Urlize(alias)), nil
}
