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
	"html/template"
	"io"
)

const (
	alias      = "<!DOCTYPE html><html><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"
	aliasXHtml = "<!DOCTYPE html><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><title>{{ .Permalink }}</title><link rel=\"canonical\" href=\"{{ .Permalink }}\"/><meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\" /><meta http-equiv=\"refresh\" content=\"0; url={{ .Permalink }}\" /></head></html>"
)

var defaultAliasTemplates *template.Template

func init() {
	defaultAliasTemplates = template.New("")
	template.Must(defaultAliasTemplates.New("alias").Parse(alias))
	template.Must(defaultAliasTemplates.New("alias-xhtml").Parse(aliasXHtml))
}

type aliasHandler struct {
	Templates *template.Template
}

func newAliasHandler(t *template.Template) aliasHandler {
	return aliasHandler{t}
}

func (a aliasHandler) renderAlias(isXHTML bool, permalink string, page *Page) (io.Reader, error) {
	t := "alias"
	if isXHTML {
		t = "alias-xhtml"
	}

	template := defaultAliasTemplates
	if a.Templates != nil {
		template = a.Templates
		t = "alias.html"
	}

	data := struct {
		Permalink string
		Page      *Page
	}{
		permalink,
		page,
	}

	buffer := new(bytes.Buffer)
	err := template.ExecuteTemplate(buffer, t, data)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
