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

package page

import "strings"

const (
	KindPage = "page"

	// The rest are node types; home page, sections etc.

	KindHome         = "home"
	KindSection      = "section"
	KindTaxonomy     = "taxonomy"
	KindTaxonomyTerm = "taxonomyTerm"
)

var kindMap = map[string]string{
	strings.ToLower(KindPage):         KindPage,
	strings.ToLower(KindHome):         KindHome,
	strings.ToLower(KindSection):      KindSection,
	strings.ToLower(KindTaxonomy):     KindTaxonomy,
	strings.ToLower(KindTaxonomyTerm): KindTaxonomyTerm,
}

// GetKind gets the page kind given a string, empty if not found.
func GetKind(s string) string {
	return kindMap[strings.ToLower(s)]
}
