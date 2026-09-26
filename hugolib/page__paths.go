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
	"path"
	"strings"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

func newPagePaths(ps *pageState) (pagePaths, error) {
	s := ps.s
	pm := ps.m

	targetPathDescriptor, err := createTargetPathDescriptor(ps)
	if err != nil {
		return pagePaths{}, err
	}

	var outputFormats output.Formats

	if ps.m.isStandalone() {
		outputFormats = output.Formats{ps.m.standaloneOutputFormat}
	} else {
		outputFormats = ps.outputFormats()
		if len(outputFormats) == 0 {
			return pagePaths{}, nil
		}

		if pm.noRender() {
			outputFormats = outputFormats[:1]
		}
	}

	pageOutputFormats := make(page.OutputFormats, len(outputFormats))
	targets := make(map[string]targetPathsHolder)

	for i, f := range outputFormats {
		desc := targetPathDescriptor
		desc.Type = f
		paths := page.CreateTargetPaths(desc)

		var relPermalink, permalink string

		// If a page is headless or bundled in another,
		// it will not get published on its own and it will have no links.
		// We also check the build options if it's set to not render or have
		// a link.
		if !pm.noLink() && !pm.bundled {
			relPermalink = paths.RelPermalink(s.PathSpec)
			permalink = paths.PermalinkForOutputFormat(s.PathSpec, f)
		}

		pageOutputFormats[i] = page.NewOutputFormat(relPermalink, permalink, len(outputFormats) == 1, f)

		// Use the main format for permalinks, usually HTML.
		permalinksIndex := 0
		if f.Permalinkable {
			// Unless it's permalinkable.
			permalinksIndex = i
		}

		targets[f.Name] = targetPathsHolder{
			relURL:       relPermalink,
			paths:        paths,
			OutputFormat: pageOutputFormats[permalinksIndex],
		}

	}

	var out page.OutputFormats
	if !pm.noLink() {
		out = pageOutputFormats
	}

	return pagePaths{
		outputFormats:        out,
		firstOutputFormat:    pageOutputFormats[0],
		targetPaths:          targets,
		targetPathDescriptor: targetPathDescriptor,
	}, nil
}

// targetPathDir returns the directory part of p's target path, with any
// slug set in the section, taxonomy or term pages above it applied.
func targetPathDir(p *pageState) string {
	unnormalized := p.s.Conf.DisablePathToLower()
	pi := p.PathInfo()
	if unnormalized {
		pi = pi.Unnormalized()
	}

	var sb strings.Builder
	sb.Grow(len(pi.Base()))

	if p.IsBranch() {
		writeBranchDir(&sb, p, unnormalized)
	} else {
		section := p.CurrentSection()
		writeBranchDir(&sb, section, unnormalized)
		sb.WriteString(strings.TrimPrefix(pi.ContainerDir(), pathBase(section, unnormalized)))
	}

	return sb.String()
}

func writeBranchDir(sb *strings.Builder, p page.Page, unnormalized bool) {
	if p.IsHome() {
		return
	}
	parent := p.Parent()
	writeBranchDir(sb, parent, unnormalized)
	rel := strings.TrimPrefix(pathBase(p, unnormalized), pathBase(parent, unnormalized))
	if slug := p.Slug(); slug != "" {
		sb.WriteString(rel[:strings.LastIndexByte(rel, '/')+1])
		sb.WriteString(slug)
		return
	}
	sb.WriteString(rel)
}

func pathBase(p page.Page, unnormalized bool) string {
	pi := p.PathInfo()
	if unnormalized {
		pi = pi.Unnormalized()
	}
	return pi.Base()
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

func createTargetPathDescriptor(p *pageState) (page.TargetPathDescriptor, error) {
	s := p.s
	d := s.Deps
	pm := p.m
	alwaysInSubDir := p.Kind() == kinds.KindSitemap
	pageInfoPage := p.PathInfo()
	if p.s.Conf.DisablePathToLower() {
		pageInfoPage = pageInfoPage.Unnormalized()
	}

	desc := page.TargetPathDescriptor{
		PathSpec:    d.PathSpec,
		Kind:        p.Kind(),
		Path:        pageInfoPage,
		Dir:         targetPathDir(p),
		UglyURLs:    s.h.Conf.IsUglyURLs(p.Section()),
		ForcePrefix: s.h.Conf.IsMultihost() || alwaysInSubDir,
		URL:         pm.pageConfig.URL,
	}

	if pm.Slug() != "" {
		desc.BaseName = pm.Slug()
	} else if pm.isStandalone() && pm.standaloneOutputFormat.BaseName != "" {
		desc.BaseName = pm.standaloneOutputFormat.BaseName
	} else {
		desc.BaseName = pageInfoPage.BaseNameNoIdentifier()
	}

	addPrefix := func(filePath, link string) {
		if filePath != "" {
			if desc.PrefixFilePath != "" {
				desc.PrefixFilePath += "/" + filePath
			} else {
				desc.PrefixFilePath += filePath
			}
		}

		if link != "" {
			if desc.PrefixLink != "" {
				desc.PrefixLink += "/" + link
			} else {
				desc.PrefixLink += link
			}
		}
	}

	// Add path prefixes.
	// Add any role first, as that is a natural candidate for external ACL checks.
	// For multihost sites, the file path needs to start with the language code.
	if s.h.Conf.IsMultihost() {
		addPrefix(s.getLanguageTargetPathLang(true), "")
	}
	rolePrefix := s.getPrefixRole()
	addPrefix(rolePrefix, rolePrefix)
	versionPrefix := s.getPrefixVersion()
	addPrefix(versionPrefix, versionPrefix)
	if s.h.Conf.IsMultihost() {
		addPrefix("", s.getLanguagePermalinkLang(alwaysInSubDir))
	} else {
		addPrefix(s.getLanguageTargetPathLang(alwaysInSubDir), s.getLanguagePermalinkLang(alwaysInSubDir))
	}

	if desc.URL != "" && strings.IndexByte(desc.URL, ':') >= 0 {
		// Attempt to parse and expand an url
		opath, err := d.ResourceSpec.Permalinks.ExpandPattern(desc.URL, p)
		if err != nil {
			return desc, err
		}

		if opath != "" {
			opath, _ = url.QueryUnescape(opath)
			desc.URL = opath
		}
	}

	opath, err := d.ResourceSpec.Permalinks.Expand(p)
	if err != nil {
		return desc, err
	}

	if opath != "" {
		opath, _ = url.QueryUnescape(opath)
		opath = path.Clean(opath)

		desc.ExpandedPermalink = opath

		if p.File() != nil {
			s.Log.Debugf("Set expanded permalink path for %s %s to %#v", p.Kind(), p.File().Path(), opath)
		} else {
			s.Log.Debugf("Set expanded permalink path for %s %s to %#v", p.Kind(), p.Path(), opath)
		}
	}

	return desc, nil
}
