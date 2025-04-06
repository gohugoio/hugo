// Copyright 2025 The Hugo Authors. All rights reserved.
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

package tplimpl

import (
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resources/kinds"
)

type layoutLegacyMapping struct {
	sourcePath string
	target     layoutLegacyMappingTarget
}

type layoutLegacyMappingTarget struct {
	targetPath     string
	targetDesc     TemplateDescriptor
	targetCategory Category
}

var (
	ltermPlural = layoutLegacyMappingTarget{
		targetPath:     "/PLURAL",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindTerm},
		targetCategory: CategoryLayout,
	}
	ltermBase = layoutLegacyMappingTarget{
		targetPath:     "",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindTerm},
		targetCategory: CategoryLayout,
	}

	ltaxPlural = layoutLegacyMappingTarget{
		targetPath:     "/PLURAL",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindTaxonomy},
		targetCategory: CategoryLayout,
	}
	ltaxBase = layoutLegacyMappingTarget{
		targetPath:     "",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindTaxonomy},
		targetCategory: CategoryLayout,
	}

	lsectBase = layoutLegacyMappingTarget{
		targetPath:     "",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindSection},
		targetCategory: CategoryLayout,
	}
	lsectTheSection = layoutLegacyMappingTarget{
		targetPath:     "/THESECTION",
		targetDesc:     TemplateDescriptor{Kind: kinds.KindSection},
		targetCategory: CategoryLayout,
	}
)

type legacyTargetPathIdentifiers struct {
	targetPath     string
	targetCategory Category
	kind           string
	lang           string
	outputFormat   string
	ext            string
}

type legacyOrdinalMapping struct {
	ordinal int
	mapping layoutLegacyMappingTarget
}

type legacyOrdinalMappingFi struct {
	m  legacyOrdinalMapping
	fi hugofs.FileMetaInfo
}

var legacyTermMappings = []layoutLegacyMapping{
	{sourcePath: "/PLURAL/term", target: ltermPlural},
	{sourcePath: "/PLURAL/SINGULAR", target: ltermPlural},
	{sourcePath: "/term/term", target: ltermBase},
	{sourcePath: "/term/SINGULAR", target: ltermPlural},
	{sourcePath: "/term/taxonomy", target: ltermPlural},
	{sourcePath: "/term/list", target: ltermBase},
	{sourcePath: "/taxonomy/term", target: ltermBase},
	{sourcePath: "/taxonomy/SINGULAR", target: ltermPlural},
	{sourcePath: "/SINGULAR/term", target: ltermPlural},
	{sourcePath: "/SINGULAR/SINGULAR", target: ltermPlural},
	{sourcePath: "/_default/SINGULAR", target: ltermPlural},
	{sourcePath: "/_default/taxonomy", target: ltermBase},
}

var legacyTaxonomyMappings = []layoutLegacyMapping{
	{sourcePath: "/PLURAL/SINGULAR.terms", target: ltaxPlural},
	{sourcePath: "/PLURAL/terms", target: ltaxPlural},
	{sourcePath: "/PLURAL/taxonomy", target: ltaxPlural},
	{sourcePath: "/PLURAL/list", target: ltaxPlural},
	{sourcePath: "/SINGULAR/SINGULAR.terms", target: ltaxPlural},
	{sourcePath: "/SINGULAR/terms", target: ltaxPlural},
	{sourcePath: "/SINGULAR/taxonomy", target: ltaxPlural},
	{sourcePath: "/SINGULAR/list", target: ltaxPlural},
	{sourcePath: "/taxonomy/SINGULAR.terms", target: ltaxPlural},
	{sourcePath: "/taxonomy/terms", target: ltaxBase},
	{sourcePath: "/taxonomy/taxonomy", target: ltaxBase},
	{sourcePath: "/taxonomy/list", target: ltaxBase},
	{sourcePath: "/_default/SINGULAR.terms", target: ltaxBase},
	{sourcePath: "/_default/terms", target: ltaxBase},
	{sourcePath: "/_default/taxonomy", target: ltaxBase},
}

var legacySectionMappings = []layoutLegacyMapping{
	// E.g. /mysection/mysection.html
	{sourcePath: "/THESECTION/THESECTION", target: lsectTheSection},
	// E.g. /section/mysection.html
	{sourcePath: "/SECTIONKIND/THESECTION", target: lsectTheSection},
	// E.g. /section/section.html
	{sourcePath: "/SECTIONKIND/SECTIONKIND", target: lsectBase},
	// E.g. /section/list.html
	{sourcePath: "/SECTIONKIND/list", target: lsectBase},
	// E.g. /_default/mysection.html
	{sourcePath: "/_default/THESECTION", target: lsectTheSection},
}
