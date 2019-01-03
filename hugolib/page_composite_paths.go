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
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/resources/page"
)

func newPagePaths(
	s *Site,
	p page.Page,
	pm *pageMeta) (pagePaths, error) {

	d := s.Deps

	targetPathDescriptor, err := createTargetPathDescriptor(s, p, pm)
	if err != nil {
		return pagePaths{}, err
	}

	outputFormats := pm.outputFormats()
	if len(outputFormats) == 0 {
		// TODO(bep)
		outputFormats = pm.s.outputFormats[pm.Kind()]
	}

	if len(outputFormats) == 0 {
		return pagePaths{}, nil
	}

	pageOutputFormats := make(page.OutputFormats, len(outputFormats))
	targetPaths := make(map[string]targetPathString)

	for i, f := range outputFormats {
		desc := targetPathDescriptor
		desc.Type = f
		targetPath := page.CreateTargetPath(desc)
		rel := targetPath

		// For /index.json etc. we must  use the full path.
		if f.MediaType.FullSuffix() == ".html" && filepath.Base(rel) == "index.html" {
			rel = strings.TrimSuffix(rel, f.BaseFilename())
		}

		rel = d.PathSpec.URLizeFilename(filepath.ToSlash(rel))
		perm, err := permalinkForOutputFormat(d.PathSpec, rel, f)
		if err != nil {
			return pagePaths{}, err
		}

		// TODO(bep) page
		pageOutputFormats[i] = page.NewOutputFormat(s.PathSpec.PrependBasePath(rel, false), perm, len(outputFormats) == 1, f)
		targetPaths[f.Name] = targetPathString(targetPath)

	}

	f := outputFormats[0]
	target := targetPaths[f.Name]

	relTargetPathBase := strings.TrimSuffix(string(target), f.BaseFilename())
	relTargetPathBase = strings.Trim(strings.TrimSuffix(string(relTargetPathBase), f.MediaType.FullSuffix()), helpers.FilePathSeparator)
	if prefix := s.GetLanguagePrefix(); prefix != "" {
		// Any language code in the path will be added later.
		relTargetPathBase = strings.TrimPrefix(relTargetPathBase, prefix+helpers.FilePathSeparator)
	}

	return pagePaths{
		outputFormats:        pageOutputFormats,
		targetPaths:          targetPaths,
		targetPathDescriptor: targetPathDescriptor,
		relTargetPathBase:    relTargetPathBase,
	}, nil

	/*	target := filepath.ToSlash(p.createRelativeTargetPath())
		rel := d.PathSpec.URLizeFilename(target)

		var err error
		f := asdf[0]
		p.permalink, err = p.s.permalinkForOutputFormat(rel, f)
		if err != nil {
			return err
		}

		p.relTargetPathBase = strings.TrimPrefix(strings.TrimSuffix(target, f.MediaType.FullSuffix()), "/")
		if prefix := p.s.GetLanguagePrefix(); prefix != "" {
			// Any language code in the path will be added later.
			p.relTargetPathBase = strings.TrimPrefix(p.relTargetPathBase, prefix+"/")
		}
		p.relPermalink = p.s.PathSpec.PrependBasePath(rel, false)
		p.layoutDescriptor = p.createLayoutDescriptor()
	*/

}

type pagePaths struct {
	outputFormats page.OutputFormats

	targetPaths          map[string]targetPathString
	targetPathDescriptor page.TargetPathDescriptor

	// relative target path without extension and any base path element
	// from the baseURL or the language code.
	// This is used to construct paths in the page resources.
	relTargetPathBase string
}

func (l pagePaths) OutputFormats() page.OutputFormats {
	return l.outputFormats
}

func (l pagePaths) Permalink() string {
	return l.outputFormats[0].Permalink()
}

func (l pagePaths) RelPermalink() string {
	return l.outputFormats[0].RelPermalink()
}

func createTargetPathDescriptor(s *Site, p page.Page, pm *pageMeta) (page.TargetPathDescriptor, error) {
	var (
		dir      string
		baseName string
	)

	d := s.Deps

	if p.File() != nil {
		dir = p.File().Dir()
		baseName = p.File().TranslationBaseName()
	}

	desc := page.TargetPathDescriptor{
		PathSpec:    d.PathSpec,
		Kind:        p.Kind(),
		Sections:    p.SectionsEntries(),
		UglyURLs:    s.Info.uglyURLs(p),
		Dir:         dir,
		URL:         pm.urlPaths.URL,
		IsMultihost: s.h.IsMultihost(),
	}

	if pm.Slug() != "" {
		desc.BaseName = pm.Slug()
	} else {
		desc.BaseName = baseName
	}

	desc.LangPrefix = s.getLanguageTargetPathLang()

	// Expand only page.KindPage and page.KindTaxonomy; don't expand other Kinds of Pages
	// like page.KindSection or page.KindTaxonomyTerm because they are "shallower" and
	// the permalink configuration values are likely to be redundant, e.g.
	// naively expanding /category/:slug/ would give /category/categories/ for
	// the "categories" page.KindTaxonomyTerm.
	if p.Kind() == page.KindPage || p.Kind() == page.KindTaxonomy {
		opath, err := d.ResourceSpec.Permalinks.Expand(p.Section(), p)
		if err != nil {
			return desc, err
		}

		if opath != "" {
			opath, _ = url.QueryUnescape(opath)
			opath = filepath.FromSlash(opath)
			desc.ExpandedPermalink = opath
		}

	}

	return desc, nil

}
