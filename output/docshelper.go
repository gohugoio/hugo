package output

import (
	"strings"

	//	"fmt"

	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() docshelper.DocProvider {
		return docshelper.DocProvider{
			"output": map[string]interface{}{
				"formats": DefaultFormats,
				"layouts": createLayoutExamples(),
			},
		}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}

func createLayoutExamples() interface{} {

	type Example struct {
		Example      string
		Kind         string
		OutputFormat string
		Suffix       string
		Layouts      []string `json:"Template Lookup Order"`
	}

	var (
		basicExamples []Example
		demoLayout    = "demolayout"
		demoType      = "demotype"
	)

	for _, example := range []struct {
		name string
		d    LayoutDescriptor
		f    Format
	}{
		// Taxonomy output.LayoutDescriptor={categories category taxonomy en  false Type Section
		{"Single page in \"posts\" section", LayoutDescriptor{Kind: "page", Type: "posts"}, HTMLFormat},
		{"Base template for single page in \"posts\" section", LayoutDescriptor{Baseof: true, Kind: "page", Type: "posts"}, HTMLFormat},
		{"Single page in \"posts\" section with layout set", LayoutDescriptor{Kind: "page", Type: "posts", Layout: demoLayout}, HTMLFormat},
		{"Base template for single page in \"posts\" section with layout set", LayoutDescriptor{Baseof: true, Kind: "page", Type: "posts", Layout: demoLayout}, HTMLFormat},
		{"AMP single page", LayoutDescriptor{Kind: "page", Type: "posts"}, AMPFormat},
		{"AMP single page, French language", LayoutDescriptor{Kind: "page", Type: "posts", Lang: "fr"}, AMPFormat},
		// All section or typeless pages gets "page" as type
		{"Home page", LayoutDescriptor{Kind: "home", Type: "page"}, HTMLFormat},
		{"Base template for home page", LayoutDescriptor{Baseof: true, Kind: "home", Type: "page"}, HTMLFormat},
		{"Home page with type set", LayoutDescriptor{Kind: "home", Type: demoType}, HTMLFormat},
		{"Base template for home page with type set", LayoutDescriptor{Baseof: true, Kind: "home", Type: demoType}, HTMLFormat},
		{"Home page with layout set", LayoutDescriptor{Kind: "home", Type: "page", Layout: demoLayout}, HTMLFormat},
		{`AMP home, French language"`, LayoutDescriptor{Kind: "home", Type: "page", Lang: "fr"}, AMPFormat},
		{"JSON home", LayoutDescriptor{Kind: "home", Type: "page"}, JSONFormat},
		{"RSS home", LayoutDescriptor{Kind: "home", Type: "page"}, RSSFormat},
		{"RSS section posts", LayoutDescriptor{Kind: "section", Type: "posts"}, RSSFormat},
		{"Taxonomy list in categories", LayoutDescriptor{Kind: "taxonomy", Type: "categories", Section: "category"}, RSSFormat},
		{"Taxonomy terms in categories", LayoutDescriptor{Kind: "taxonomyTerm", Type: "categories", Section: "category"}, RSSFormat},
		{"Section list for \"posts\" section", LayoutDescriptor{Kind: "section", Type: "posts", Section: "posts"}, HTMLFormat},
		{"Section list for \"posts\" section with type set to \"blog\"", LayoutDescriptor{Kind: "section", Type: "blog", Section: "posts"}, HTMLFormat},
		{"Section list for \"posts\" section with layout set to \"demoLayout\"", LayoutDescriptor{Kind: "section", Layout: demoLayout, Section: "posts"}, HTMLFormat},

		{"Taxonomy list in categories", LayoutDescriptor{Kind: "taxonomy", Type: "categories", Section: "category"}, HTMLFormat},
		{"Taxonomy term in categories", LayoutDescriptor{Kind: "taxonomyTerm", Type: "categories", Section: "category"}, HTMLFormat},
	} {

		l := NewLayoutHandler()
		layouts, _ := l.For(example.d, example.f)

		basicExamples = append(basicExamples, Example{
			Example:      example.name,
			Kind:         example.d.Kind,
			OutputFormat: example.f.Name,
			Suffix:       example.f.MediaType.Suffix(),
			Layouts:      makeLayoutsPresentable(layouts)})
	}

	return basicExamples

}

func makeLayoutsPresentable(l []string) []string {
	var filtered []string
	for _, ll := range l {
		if strings.Contains(ll, "page/") {
			// This is a valid lookup, but it's more confusing than useful.
			continue
		}
		ll = "layouts/" + strings.TrimPrefix(ll, "_text/")

		if !strings.Contains(ll, "indexes") {
			filtered = append(filtered, ll)
		}
	}

	return filtered
}
