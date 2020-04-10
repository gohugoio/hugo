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
	"strings"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/resources/page"
)

func newPagePaths(
	s *Site,
	p page.Page,
	pm *pageMeta) (pagePaths, error) {

	targetPathDescriptor, err := createTargetPathDescriptor(s, p, pm)
	if err != nil {
		return pagePaths{}, err
	}

	outputFormats := pm.outputFormats()
	if len(outputFormats) == 0 {
		return pagePaths{}, nil
	}

	if pm.noRender() {
		outputFormats = outputFormats[:1]
	}

	pageOutputFormats := make(page.OutputFormats, len(outputFormats))
	targets := make(map[string]targetPathsHolder)

	for i, f := range outputFormats {
		desc := targetPathDescriptor
		desc.Type = f
		paths := page.CreateTargetPaths(desc)

		var relPermalink, permalink string

		// If a page is headless or marked as "no render", or bundled in another,
		// it will not get published on its own and it will have no links.
		if !pm.noRender() && !pm.bundled {
			relPermalink = paths.RelPermalink(s.PathSpec)
			permalink = paths.PermalinkForOutputFormat(s.PathSpec, f)
		}

		pageOutputFormats[i] = page.NewOutputFormat(relPermalink, permalink, len(outputFormats) == 1, f)

		// Use the main format for permalinks, usually HTML.
		permalinksIndex := 0
		if f.Permalinkable {
			// Unless it's permalinkable
			permalinksIndex = i
		}

		targets[f.Name] = targetPathsHolder{
			paths:        paths,
			OutputFormat: pageOutputFormats[permalinksIndex]}

	}

	var out page.OutputFormats
	if !pm.noRender() {
		out = pageOutputFormats
	}

	return pagePaths{
		outputFormats:        out,
		firstOutputFormat:    pageOutputFormats[0],
		targetPaths:          targets,
		targetPathDescriptor: targetPathDescriptor,
	}, nil

}

type pagePaths struct {
	outputFormats     page.OutputFormats
	firstOutputFormat page.OutputFormat

	targetPaths          map[string]targetPathsHolder
	targetPathDescriptor page.TargetPathDescriptor
}

func (l pagePaths) OutputFormats() page.OutputFormats {
	return l.outputFormats
}

func createTargetPathDescriptor(s *Site, p page.Page, pm *pageMeta) (page.TargetPathDescriptor, error) {
	var (
		dir             string
		baseName        string
		contentBaseName string
	)

	d := s.Deps

	if !p.File().IsZero() {
		dir = p.File().Dir()
		baseName = p.File().TranslationBaseName()
		contentBaseName = p.File().ContentBaseName()
	}

	if baseName != contentBaseName {
		// See https://github.com/gohugoio/hugo/issues/4870
		// A leaf bundle
		dir = strings.TrimSuffix(dir, contentBaseName+helpers.FilePathSeparator)
		baseName = contentBaseName
	}

	alwaysInSubDir := p.Kind() == kindSitemap

	desc := page.TargetPathDescriptor{
		PathSpec:    d.PathSpec,
		Kind:        p.Kind(),
		Sections:    p.SectionsEntries(),
		UglyURLs:    s.Info.uglyURLs(p),
		ForcePrefix: s.h.IsMultihost() || alwaysInSubDir,
		Dir:         dir,
		URL:         pm.urlPaths.URL,
	}

	if pm.Slug() != "" {
		desc.BaseName = pm.Slug()
	} else {
		desc.BaseName = baseName
	}

	desc.PrefixFilePath = s.getLanguageTargetPathLang(alwaysInSubDir)
	desc.PrefixLink = s.getLanguagePermalinkLang(alwaysInSubDir)

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
			desc.ExpandedPermalink = opath
		}

	}

	return desc, nil

}
