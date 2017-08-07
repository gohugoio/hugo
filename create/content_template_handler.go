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
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/afero"
)

// ArchetypeFileData represents the data available to an archetype template.
type ArchetypeFileData struct {
	// The archetype content type, either given as --kind option or extracted
	// from the target path's section, i.e. "blog/mypost.md" will resolve to
	// "blog".
	Type string

	// The current date and time as a RFC3339 formatted string, suitable for use in front matter.
	Date string

	// The Site, fully equipped with all the pages etc. Note: This will only be set if it is actually
	// used in the archetype template. Also, if this is a multilingual setup,
	// this site is the site that best matches the target content file, based
	// on the presence of language code in the filename.
	Site *hugolib.Site

	// The target content file. Note that the .Content will be empty, as that
	// has not been created yet.
	*source.File
}

const (
	// ArchetypeTemplateTemplate is used as initial template when adding an archetype template.
	ArchetypeTemplateTemplate = `---
title: "{{ replace .TranslationBaseName "-" " " | title }}"
date: {{ .Date }}
draft: true
---

`
)

var (
	archetypeShortcodeReplacementsPre = strings.NewReplacer(
		"{{<", "{x{<",
		"{{%", "{x{%",
		">}}", ">}x}",
		"%}}", "%}x}")

	archetypeShortcodeReplacementsPost = strings.NewReplacer(
		"{x{<", "{{<",
		"{x{%", "{{%",
		">}x}", ">}}",
		"%}x}", "%}}")
)

func executeArcheTypeAsTemplate(s *hugolib.Site, kind, targetPath, archetypeFilename string) ([]byte, error) {

	var (
		archetypeContent  []byte
		archetypeTemplate []byte
		err               error
	)

	sp := source.NewSourceSpec(s.Deps.Cfg, s.Deps.Fs)
	f := sp.NewFile(targetPath)

	data := ArchetypeFileData{
		Type: kind,
		Date: time.Now().Format(time.RFC3339),
		File: f,
		Site: s,
	}

	if archetypeFilename == "" {
		// TODO(bep) archetype revive the issue about wrong tpl funcs arg order
		archetypeTemplate = []byte(ArchetypeTemplateTemplate)
	} else {
		archetypeTemplate, err = afero.ReadFile(s.Fs.Source, archetypeFilename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read archetype file %q: %s", archetypeFilename, err)
		}

	}

	// The archetype template may contain shortcodes, and these does not play well
	// with the Go templates. Need to set some temporary delimiters.
	archetypeTemplate = []byte(archetypeShortcodeReplacementsPre.Replace(string(archetypeTemplate)))

	// Reuse the Hugo template setup to get the template funcs properly set up.
	templateHandler := s.Deps.Tmpl.(tpl.TemplateHandler)
	templateName := "_text/" + helpers.Filename(archetypeFilename)
	if err := templateHandler.AddTemplate(templateName, string(archetypeTemplate)); err != nil {
		return nil, fmt.Errorf("Failed to parse archetype file %q: %s", archetypeFilename, err)
	}

	templ := templateHandler.Lookup(templateName)

	var buff bytes.Buffer
	if err := templ.Execute(&buff, data); err != nil {
		return nil, fmt.Errorf("Failed to process archetype file %q: %s", archetypeFilename, err)
	}

	archetypeContent = []byte(archetypeShortcodeReplacementsPost.Replace(buff.String()))

	if !bytes.Contains(archetypeContent, []byte("date")) || !bytes.Contains(archetypeContent, []byte("title")) {
		// TODO(bep) remove some time in the future.
		s.Log.FEEDBACK.Println(fmt.Sprintf(`WARNING: date and/or title missing from archetype file %q.
From Hugo 0.24 this must be provided in the archetype file itself, if needed. Example:
%s
`, archetypeFilename, ArchetypeTemplateTemplate))

	}

	return archetypeContent, nil

}
