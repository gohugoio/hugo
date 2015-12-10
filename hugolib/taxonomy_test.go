// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"testing"
)

func TestSitePossibleTaxonomies(t *testing.T) {
	site := new(Site)
	page, _ := NewPageFrom(strings.NewReader(PAGE_YAML_WITH_TAXONOMIES_A), "path/to/page")
	site.Pages = append(site.Pages, page)
	taxonomies := site.possibleTaxonomies()
	if !compareStringSlice(taxonomies, []string{"tags", "categories"}) {
		if !compareStringSlice(taxonomies, []string{"categories", "tags"}) {
			t.Fatalf("possible taxonomies do not match [tags categories].  Got: %s", taxonomies)
		}
	}
}
