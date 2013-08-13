package hugolib

import (
	"strings"
	"testing"
)

func TestSitePossibleIndexes(t *testing.T) {
	site := new(Site)
	page, _ := ReadFrom(strings.NewReader(PAGE_YAML_WITH_INDEXES_A), "path/to/page")
	site.Pages = append(site.Pages, page)
	indexes := site.possibleIndexes()
	if !compareStringSlice(indexes, []string{"tags", "categories"}) {
		t.Fatalf("possible indexes do not match [tags categories].  Got: %s", indexes)
	}
}
