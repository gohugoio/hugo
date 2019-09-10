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
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

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
	Site *hugolib.SiteInfo

	// Name will in most cases be the same as TranslationBaseName, e.g. "my-post".
	// But if that value is "index" (bundles), the Name is instead the owning folder.
	// This is the value you in most cases would want to use to construct the title in your
	// archetype template.
	Name string

	// The target content file. Note that the .Content will be empty, as that
	// has not been created yet.
	source.File
}

const (
	// ArchetypeTemplateTemplate is used as initial template when adding an archetype template.
	ArchetypeTemplateTemplate = `---
title: "{{ replace .Name "-" " " | title }}"
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

func executeArcheTypeAsTemplate(s *hugolib.Site, name, kind, targetPath, archetypeFilename string) ([]byte, error) {

	var (
		archetypeContent  []byte
		archetypeTemplate []byte
		err               error
	)

	f, err := s.SourceSpec.NewFileInfoFrom(targetPath, targetPath)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = f.TranslationBaseName()

		if name == "index" || name == "_index" {
			// Page bundles; the directory name will hopefully have a better name.
			dir := strings.TrimSuffix(f.Dir(), helpers.FilePathSeparator)
			_, name = filepath.Split(dir)
		}
	}

	data := ArchetypeFileData{
		Type: kind,
		Date: time.Now().Format(time.RFC3339),
		Name: name,
		File: f,
		Site: s.Info,
	}

	if archetypeFilename == "" {
		// TODO(bep) archetype revive the issue about wrong tpl funcs arg order
		archetypeTemplate = []byte(ArchetypeTemplateTemplate)
	} else {
		archetypeTemplate, err = afero.ReadFile(s.BaseFs.Archetypes.Fs, archetypeFilename)
		if err != nil {
			return nil, fmt.Errorf("failed to read archetype file %s", err)
		}

	}

	// The archetype template may contain shortcodes, and these does not play well
	// with the Go templates. Need to set some temporary delimiters.
	archetypeTemplate = []byte(archetypeShortcodeReplacementsPre.Replace(string(archetypeTemplate)))

	// Reuse the Hugo template setup to get the template funcs properly set up.
	templateHandler := s.Deps.Tmpl().(tpl.TemplateManager)
	templateName := helpers.Filename(archetypeFilename)
	if err := templateHandler.AddTemplate("_text/"+templateName, string(archetypeTemplate)); err != nil {
		return nil, errors.Wrapf(err, "Failed to parse archetype file %q:", archetypeFilename)
	}

	templ, _ := templateHandler.Lookup(templateName)

	var buff bytes.Buffer
	if err := templateHandler.Execute(templ, &buff, data); err != nil {
		return nil, errors.Wrapf(err, "Failed to process archetype file %q:", archetypeFilename)
	}

	archetypeContent = []byte(archetypeShortcodeReplacementsPost.Replace(buff.String()))

	return archetypeContent, nil

}
