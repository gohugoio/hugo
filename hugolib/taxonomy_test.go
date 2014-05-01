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
