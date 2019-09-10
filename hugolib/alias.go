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

package hugolib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/publisher"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl"
)

type aliasHandler struct {
	t         tpl.TemplateHandler
	log       *loggers.Logger
	allowRoot bool
}

func newAliasHandler(t tpl.TemplateHandler, l *loggers.Logger, allowRoot bool) aliasHandler {
	return aliasHandler{t, l, allowRoot}
}

type aliasPage struct {
	Permalink string
	page.Page
}

func (a aliasHandler) renderAlias(permalink string, p page.Page) (io.Reader, error) {

	var templ tpl.Template
	var found bool

	templ, found = a.t.Lookup("alias.html")
	if !found {
		// TODO(bep) consolidate
		templ, found = a.t.Lookup("_internal/alias.html")
		if !found {
			return nil, errors.New("no alias template found")
		}
	}

	data := aliasPage{
		permalink,
		p,
	}

	buffer := new(bytes.Buffer)
	err := a.t.Execute(templ, buffer, data)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *Site) writeDestAlias(path, permalink string, outputFormat output.Format, p page.Page) (err error) {
	return s.publishDestAlias(false, path, permalink, outputFormat, p)
}

func (s *Site) publishDestAlias(allowRoot bool, path, permalink string, outputFormat output.Format, p page.Page) (err error) {
	handler := newAliasHandler(s.Tmpl(), s.Log, allowRoot)

	s.Log.DEBUG.Println("creating alias:", path, "redirecting to", permalink)

	targetPath, err := handler.targetPathAlias(path)
	if err != nil {
		return err
	}

	aliasContent, err := handler.renderAlias(permalink, p)
	if err != nil {
		return err
	}

	pd := publisher.Descriptor{
		Src:          aliasContent,
		TargetPath:   targetPath,
		StatCounter:  &s.PathSpec.ProcessingStats.Aliases,
		OutputFormat: outputFormat,
	}

	return s.publisher.Publish(pd)

}

func (a aliasHandler) targetPathAlias(src string) (string, error) {
	originalAlias := src
	if len(src) <= 0 {
		return "", fmt.Errorf("alias \"\" is an empty string")
	}

	alias := path.Clean(filepath.ToSlash(src))

	if !a.allowRoot && alias == "/" {
		return "", fmt.Errorf("alias \"%s\" resolves to website root directory", originalAlias)
	}

	components := strings.Split(alias, "/")

	// Validate against directory traversal
	if components[0] == ".." {
		return "", fmt.Errorf("alias \"%s\" traverses outside the website root directory", originalAlias)
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
			return "", fmt.Errorf("cannot create \"%s\": Windows filename restriction", originalAlias)
		}
		for _, m := range msgs {
			a.log.INFO.Println(m)
		}
	}

	// Add the final touch
	alias = strings.TrimPrefix(alias, "/")
	if strings.HasSuffix(alias, "/") {
		alias = alias + "index.html"
	} else if !strings.HasSuffix(alias, ".html") {
		alias = alias + "/" + "index.html"
	}

	return filepath.FromSlash(alias), nil
}
