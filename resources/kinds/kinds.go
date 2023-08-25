// Copyright 2023 The Hugo Authors. All rights reserved.
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

package kinds

import (
	"sort"
	"strings"
)

const (
	KindPage = "page"

	// The rest are node types; home page, sections etc.

	KindHome    = "home"
	KindSection = "section"

	// Note that before Hugo 0.73 these were confusingly named
	// taxonomy (now: term)
	// taxonomyTerm (now: taxonomy)
	KindTaxonomy = "taxonomy"
	KindTerm     = "term"

	// The following are (currently) temporary nodes,
	// i.e. nodes we create just to render in isolation.
	KindRSS       = "rss"
	KindSitemap   = "sitemap"
	KindRobotsTXT = "robotstxt"
	Kind404       = "404"
)

var (
	// This is all the kinds we can expect to find in .Site.Pages.
	AllKindsInPages []string
	// This is all the kinds, including the temporary ones.
	AllKinds []string
)

func init() {
	for k := range kindMapMain {
		AllKindsInPages = append(AllKindsInPages, k)
		AllKinds = append(AllKinds, k)
	}

	for k := range kindMapTemporary {
		AllKinds = append(AllKinds, k)
	}

	// Sort the slices for determinism.
	sort.Strings(AllKindsInPages)
	sort.Strings(AllKinds)
}

var kindMapMain = map[string]string{
	KindPage:     KindPage,
	KindHome:     KindHome,
	KindSection:  KindSection,
	KindTaxonomy: KindTaxonomy,
	KindTerm:     KindTerm,

	// Legacy, pre v0.53.0.
	"taxonomyterm": KindTaxonomy,
}

var kindMapTemporary = map[string]string{
	KindRSS:       KindRSS,
	KindSitemap:   KindSitemap,
	KindRobotsTXT: KindRobotsTXT,
	Kind404:       Kind404,
}

// GetKindMain gets the page kind given a string, empty if not found.
// Note that this will not return any temporary kinds (e.g. robotstxt).
func GetKindMain(s string) string {
	return kindMapMain[strings.ToLower(s)]
}

// GetKindAny gets the page kind given a string, empty if not found.
func GetKindAny(s string) string {
	if pkind := GetKindMain(s); pkind != "" {
		return pkind
	}
	return kindMapTemporary[strings.ToLower(s)]
}

// IsDeprecatedAndReplacedWith returns the new kind if the given kind is deprecated.
func IsDeprecatedAndReplacedWith(s string) string {
	s = strings.ToLower(s)

	switch s {
	case "taxonomyterm":
		return KindTaxonomy
	default:
		return ""
	}
}
