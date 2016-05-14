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

	"github.com/spf13/viper"
)

func TestByCountOrderOfTaxonomies(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	viper.Set("taxonomies", taxonomies)

	site := new(Site)
	page, _ := NewPageFrom(strings.NewReader(pageYamlWithTaxonomiesA), "path/to/page")
	site.Pages = append(site.Pages, page)
	site.assembleTaxonomies()

	st := make([]string, 0)
	for _, t := range site.Taxonomies["tags"].ByCount() {
		st = append(st, t.Name)
	}

	if !compareStringSlice(st, []string{"a", "b", "c"}) {
		t.Fatalf("ordered taxonomies do not match [a, b, c].  Got: %s", st)
	}
}
