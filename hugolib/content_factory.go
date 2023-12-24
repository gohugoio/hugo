// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/spf13/afero"
)

// ContentFactory creates content files from archetype templates.
type ContentFactory struct {
	h *HugoSites

	// We parse the archetype templates as Go templates, so we need
	// to replace any shortcode with a temporary placeholder.
	shortcodeReplacerPre  *strings.Replacer
	shortcodeReplacerPost *strings.Replacer
}

// ApplyArchetypeFilename archetypeFilename to w as a template using the given Page p as the foundation for the data context.
func (f ContentFactory) ApplyArchetypeFi(w io.Writer, p page.Page, archetypeKind string, fi hugofs.FileMetaInfo) error {
	if fi.IsDir() {
		return fmt.Errorf("archetype directory (%q) not supported", fi.Meta().Filename)
	}

	templateSource, err := fi.Meta().ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read archetype file %q: %s: %w", fi.Meta().Filename, err, err)
	}

	return f.ApplyArchetypeTemplate(w, p, archetypeKind, string(templateSource))
}

// ApplyArchetypeTemplate templateSource to w as a template using the given Page p as the foundation for the data context.
func (f ContentFactory) ApplyArchetypeTemplate(w io.Writer, p page.Page, archetypeKind, templateSource string) error {
	ps := p.(*pageState)
	if archetypeKind == "" {
		archetypeKind = p.Type()
	}

	d := &archetypeFileData{
		Type: archetypeKind,
		Date: htime.Now().Format(time.RFC3339),
		Page: p,
		File: p.File(),
	}

	templateSource = f.shortcodeReplacerPre.Replace(templateSource)

	templ, err := ps.s.TextTmpl().Parse("archetype.md", string(templateSource))
	if err != nil {
		return fmt.Errorf("failed to parse archetype template: %s: %w", err, err)
	}

	result, err := executeToString(context.Background(), ps.s.Tmpl(), templ, d)
	if err != nil {
		return fmt.Errorf("failed to execute archetype template: %s: %w", err, err)
	}

	_, err = io.WriteString(w, f.shortcodeReplacerPost.Replace(result))

	return err
}

func (f ContentFactory) SectionFromFilename(filename string) (string, error) {
	filename = filepath.Clean(filename)
	rel, _, err := f.h.AbsProjectContentDir(filename)
	if err != nil {
		return "", err
	}

	parts := strings.Split(paths.ToSlashTrimLeading(rel), "/")
	if len(parts) < 2 {
		return "", nil
	}
	return parts[0], nil
}

// CreateContentPlaceHolder creates a content placeholder file inside the
// best matching content directory.
func (f ContentFactory) CreateContentPlaceHolder(filename string, force bool) (string, error) {
	filename = filepath.Clean(filename)
	_, abs, err := f.h.AbsProjectContentDir(filename)
	if err != nil {
		return "", err
	}

	// This will be overwritten later, just write a placeholder to get
	// the paths correct.
	placeholder := `---
title: "Content Placeholder"
_build:
  render: never
  list: never
  publishResources: false
---

`

	if force {
		return abs, afero.WriteReader(f.h.Fs.Source, abs, strings.NewReader(placeholder))
	}
	return abs, afero.SafeWriteReader(f.h.Fs.Source, abs, strings.NewReader(placeholder))
}

// NewContentFactory creates a new ContentFactory for h.
func NewContentFactory(h *HugoSites) ContentFactory {
	return ContentFactory{
		h: h,
		shortcodeReplacerPre: strings.NewReplacer(
			"{{<", "{x{<",
			"{{%", "{x{%",
			">}}", ">}x}",
			"%}}", "%}x}"),
		shortcodeReplacerPost: strings.NewReplacer(
			"{x{<", "{{<",
			"{x{%", "{{%",
			">}x}", ">}}",
			"%}x}", "%}}"),
	}
}

// archetypeFileData represents the data available to an archetype template.
type archetypeFileData struct {
	// The archetype content type, either given as --kind option or extracted
	// from the target path's section, i.e. "blog/mypost.md" will resolve to
	// "blog".
	Type string

	// The current date and time as a RFC3339 formatted string, suitable for use in front matter.
	Date string

	// The temporary page. Note that only the file path information is relevant at this stage.
	Page page.Page

	// File is the same as Page.File, embedded here for historic reasons.
	// TODO(bep) make this a method.
	*source.File
}

func (f *archetypeFileData) Site() page.Site {
	return f.Page.Site()
}

func (f *archetypeFileData) Name() string {
	return f.Page.File().ContentBaseName()
}
