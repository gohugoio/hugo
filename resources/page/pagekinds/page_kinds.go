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

package pagekinds

import "strings"

const (
	Page = "page"

	// Branch nodes.
	Home     = "home"
	Section  = "section"
	Taxonomy = "taxonomy"
	Term     = "term"

	// Special purpose page kinds.
	Sitemap   = "sitemap"
	RobotsTXT = "robotsTXT"
	Status404 = "404"
)

var kindMap = map[string]string{
	strings.ToLower(Page):     Page,
	strings.ToLower(Home):     Home,
	strings.ToLower(Section):  Section,
	strings.ToLower(Taxonomy): Taxonomy,
	strings.ToLower(Term):     Term,

	// Legacy.
	"taxonomyterm": Taxonomy,
	"rss":          "RSS",
}

// Get gets the page kind given a string, empty if not found.
func Get(s string) string {
	return kindMap[strings.ToLower(s)]
}

// IsBranch determines whether s represents a branch node (e.g. a section).
func IsBranch(s string) bool {
	return s == Home || s == Section || s == Taxonomy || s == Term
}
