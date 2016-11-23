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

package target

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

const alias = "<!DOCTYPE html><html><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"
const aliasXHtml = "<!DOCTYPE html><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"

var defaultAliasTemplates *template.Template

func init() {
	defaultAliasTemplates = template.New("")
	template.Must(defaultAliasTemplates.New("alias").Parse(alias))
	template.Must(defaultAliasTemplates.New("alias-xhtml").Parse(aliasXHtml))
}

type AliasPublisher interface {
	Translator
	Publish(path string, permalink string, page interface{}) error
}

type HTMLRedirectAlias struct {
	PublishDir string
	Templates  *template.Template
	AllowRoot  bool // for the language redirects
}

func (h *HTMLRedirectAlias) Translate(alias string) (aliasPath string, err error) {
	originalAlias := alias
	if len(alias) <= 0 {
		return "", fmt.Errorf("Alias \"\" is an empty string")
	}

	alias = filepath.Clean(alias)
	components := strings.Split(alias, helpers.FilePathSeparator)

	if !h.AllowRoot && alias == helpers.FilePathSeparator {
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
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

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
				jww.ERROR.Println(m)
			}
			return "", fmt.Errorf("Cannot create \"%s\": Windows filename restriction", originalAlias)
		}
		for _, m := range msgs {
			jww.WARN.Println(m)
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
		jww.INFO.Printf("Alias \"%s\" translated to \"%s\"\n", originalAlias, alias)
	}

	return filepath.Join(h.PublishDir, alias), nil
}

type AliasNode struct {
	Permalink string
	Page      interface{}
}

func (h *HTMLRedirectAlias) Publish(path string, permalink string, page interface{}) (err error) {
	if path, err = h.Translate(path); err != nil {
		jww.ERROR.Printf("%s, skipping.", err)
		return nil
	}

	t := "alias"
	if strings.HasSuffix(path, ".xhtml") {
		t = "alias-xhtml"
	}

	template := defaultAliasTemplates
	if h.Templates != nil {
		template = h.Templates
		t = "alias.html"
	}

	buffer := new(bytes.Buffer)
	err = template.ExecuteTemplate(buffer, t, &AliasNode{permalink, page})
	if err != nil {
		return
	}

	return helpers.WriteToDisk(path, buffer, hugofs.Destination())
}
