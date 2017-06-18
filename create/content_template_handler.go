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

package create

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/afero"
)

const (
	archetypeTemplateTemplate = `+++
title = "{{ replace .BaseFileName "-" " " | title }}"
date = {{ .Date }}
draft = true
+++`
)

func executeArcheTypeAsTemplate(s *hugolib.Site, kind, targetPath, archetypeFilename string) ([]byte, error) {

	var (
		archetypeContent  []byte
		archetypeTemplate []byte
		err               error
	)

	sp := source.NewSourceSpec(s.Deps.Cfg, s.Deps.Fs)
	f := sp.NewFile(targetPath)

	data := struct {
		Type string
		Date string
		*source.File
	}{
		Type: kind,
		Date: time.Now().Format(time.RFC3339),
		File: f,
	}

	if archetypeFilename == "" {
		// TODO(bep) archetype revive the issue about wrong tpl funcs arg order
		archetypeTemplate = []byte(archetypeTemplateTemplate)
	} else {
		archetypeTemplate, err = afero.ReadFile(s.Fs.Source, archetypeFilename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read archetype file %q: %s", archetypeFilename, err)
		}

	}

	// Reuse the Hugo template setup to get the template funcs properly set up.
	templateHandler := s.Deps.Tmpl.(tpl.TemplateHandler)
	if err := templateHandler.AddTemplate("_text/archetype", string(archetypeTemplate)); err != nil {
		return nil, fmt.Errorf("Failed to parse archetype file %q: %s", archetypeFilename, err)
	}

	templ := templateHandler.Lookup("_text/archetype")

	var buff bytes.Buffer
	if err := templ.Execute(&buff, data); err != nil {
		return nil, fmt.Errorf("Failed to process archetype file %q: %s", archetypeFilename, err)
	}

	archetypeContent = buff.Bytes()

	if !bytes.Contains(archetypeContent, []byte("date")) || !bytes.Contains(archetypeContent, []byte("title")) {
		// TODO(bep) remove some time in the future.
		s.Log.FEEDBACK.Println(fmt.Sprintf(`WARNING: date and/or title missing from archetype file %q. 
From Hugo 0.24 this must be provided in the archetype file itself, if needed. Example:
%s
`, archetypeFilename, archetypeTemplateTemplate))

	}

	return archetypeContent, nil

}
