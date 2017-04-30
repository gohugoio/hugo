// Copyright 2017 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/hugo/tpl"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/spf13/hugo/helpers"
)

const (
	alias      = "<!DOCTYPE html><html><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"
	aliasXHtml = "<!DOCTYPE html><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"
)

var defaultAliasTemplates *template.Template

func init() {
	//TODO(bep) consolidate
	defaultAliasTemplates = template.New("")
	template.Must(defaultAliasTemplates.New("alias").Parse(alias))
	template.Must(defaultAliasTemplates.New("alias-xhtml").Parse(aliasXHtml))
}

type aliasHandler struct {
	t         tpl.TemplateFinder
	log       *jww.Notepad
	allowRoot bool
}

func newAliasHandler(t tpl.TemplateFinder, l *jww.Notepad, allowRoot bool) aliasHandler {
	return aliasHandler{t, l, allowRoot}
}

func (a aliasHandler) renderAlias(isXHTML bool, permalink string, page *Page) (io.Reader, error) {
	t := "alias"
	if isXHTML {
		t = "alias-xhtml"
	}

	var templ *tpl.TemplateAdapter

	if a.t != nil {
		templ = a.t.Lookup("alias.html")
	}

	if templ == nil {
		def := defaultAliasTemplates.Lookup(t)
		if def != nil {
			templ = &tpl.TemplateAdapter{Template: def}
		}

	}
	data := struct {
		Permalink string
		Page      *Page
	}{
		permalink,
		page,
	}

	buffer := new(bytes.Buffer)
	err := templ.Execute(buffer, data)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *Site) writeDestAlias(path, permalink string, p *Page) (err error) {
	return s.publishDestAlias(false, path, permalink, p)
}

func (s *Site) publishDestAlias(allowRoot bool, path, permalink string, p *Page) (err error) {
	handler := newAliasHandler(s.Tmpl, s.Log, allowRoot)

	isXHTML := strings.HasSuffix(path, ".xhtml")

	s.Log.DEBUG.Println("creating alias:", path, "redirecting to", permalink)

	targetPath, err := handler.targetPathAlias(path)
	if err != nil {
		return err
	}

	aliasContent, err := handler.renderAlias(isXHTML, permalink, p)
	if err != nil {
		return err
	}

	return s.publish(targetPath, aliasContent)

}

func (a aliasHandler) targetPathAlias(src string) (string, error) {
	originalAlias := src
	if len(src) <= 0 {
		return "", fmt.Errorf("Alias \"\" is an empty string")
	}

	alias := filepath.Clean(src)
	components := strings.Split(alias, helpers.FilePathSeparator)

	if !a.allowRoot && alias == helpers.FilePathSeparator {
		return "", fmt.Errorf("Alias \"%s\" resolves to website root directory", originalAlias)
	}

	// Validate against directory traversal
	if components[0] == ".." {
		return "", fmt.Errorf("Alias \"%s\" traverses outside the website root directory", originalAlias)
	}

	// Handle Windows file and directory naming restrictions
	// See "Naming Files, Paths, and Namespaces" on MSDN
	// https://msdn.microsoft.com/en-us/library/aa365247%28v=VS.85%29.aspx?f=255&MSPPError=-2147217396
	msgs := []string{}
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM0", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT0", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	if strings.ContainsAny(alias, ":*?\"<>|") {
		msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains invalid characters on Windows: : * ? \" < > |", originalAlias))
	}
	for _, ch := range alias {
		if ch < ' ' {
			msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains ASCII control code (0x00 to 0x1F), invalid on Windows: : * ? \" < > |", originalAlias))
			continue
		}
	}
	for _, comp := range components {
		if strings.HasSuffix(comp, " ") || strings.HasSuffix(comp, ".") {
			msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains component with a trailing space or period, problematic on Windows", originalAlias))
		}
		for _, r := range reservedNames {
			if comp == r {
				msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains component with reserved name \"%s\" on Windows", originalAlias, r))
			}
		}
	}
	if len(msgs) > 0 {
		if runtime.GOOS == "windows" {
			for _, m := range msgs {
				a.log.ERROR.Println(m)
			}
			return "", fmt.Errorf("Cannot create \"%s\": Windows filename restriction", originalAlias)
		}
		for _, m := range msgs {
			a.log.WARN.Println(m)
		}
	}

	// Add the final touch
	alias = strings.TrimPrefix(alias, helpers.FilePathSeparator)
	if strings.HasSuffix(alias, helpers.FilePathSeparator) {
		alias = alias + "index.html"
	} else if !strings.HasSuffix(alias, ".html") {
		alias = alias + helpers.FilePathSeparator + "index.html"
	}
	if originalAlias != alias {
		a.log.INFO.Printf("Alias \"%s\" translated to \"%s\"\n", originalAlias, alias)
	}

	return alias, nil
}
