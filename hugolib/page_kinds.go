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
	"strings"

	"github.com/gohugoio/hugo/resources/page"
)

var (

	// This is all the kinds we can expect to find in .Site.Pages.
	allKindsInPages = []string{page.KindPage, page.KindHome, page.KindSection, page.KindTaxonomy, page.KindTaxonomyTerm}
)

const (

	// Temporary state.
	kindUnknown = "unknown"

	// The following are (currently) temporary nodes,
	// i.e. nodes we create just to render in isolation.
	kindRSS       = "RSS"
	kindSitemap   = "sitemap"
	kindRobotsTXT = "robotsTXT"
	kind404       = "404"

	pageResourceType = "page"
)

var kindMap = map[string]string{
	strings.ToLower(kindRSS):       kindRSS,
	strings.ToLower(kindSitemap):   kindSitemap,
	strings.ToLower(kindRobotsTXT): kindRobotsTXT,
	strings.ToLower(kind404):       kind404,
}

func getKind(s string) string {
	if pkind := page.GetKind(s); pkind != "" {
		return pkind
	}
	return kindMap[strings.ToLower(s)]
}
