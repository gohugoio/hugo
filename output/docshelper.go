package output

import (
	"strings"

	//	"fmt"

	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})

		docs["formats"] = DefaultFormats
		docs["layouts"] = createLayoutExamples()
		return docs
	}

	docshelper.AddDocProvider("output", docsProvider)
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
		name     string
		d        LayoutDescriptor
		hasTheme bool
		f        Format
	}{
		// Taxonomy output.LayoutDescriptor={categories category taxonomy en  false Type Section
		{"Single page in \"posts\" section", LayoutDescriptor{Kind: "page", Type: "posts"}, false, HTMLFormat},
		{"Single page in \"posts\" section with layout set", LayoutDescriptor{Kind: "page", Type: "posts", Layout: demoLayout}, false, HTMLFormat},
		{"Single page in \"posts\" section with theme", LayoutDescriptor{Kind: "page", Type: "posts"}, true, HTMLFormat},
		{"AMP single page", LayoutDescriptor{Kind: "page", Type: "posts"}, false, AMPFormat},
		{"AMP single page, French language", LayoutDescriptor{Kind: "page", Type: "posts", Lang: "fr"}, false, AMPFormat},
		// All section or typeless pages gets "page" as type
		{"Home page", LayoutDescriptor{Kind: "home", Type: "page"}, false, HTMLFormat},
		{"Home page with type set", LayoutDescriptor{Kind: "home", Type: demoType}, false, HTMLFormat},
		{"Home page with layout set", LayoutDescriptor{Kind: "home", Type: "page", Layout: demoLayout}, false, HTMLFormat},
		{`Home page with theme`, LayoutDescriptor{Kind: "home", Type: "page"}, true, HTMLFormat},
		{`AMP home, French language"`, LayoutDescriptor{Kind: "home", Type: "page", Lang: "fr"}, false, AMPFormat},
		{"JSON home", LayoutDescriptor{Kind: "home", Type: "page"}, false, JSONFormat},
		{"RSS home", LayoutDescriptor{Kind: "home", Type: "page"}, false, RSSFormat},

		{"Section list for \"posts\" section", LayoutDescriptor{Kind: "section", Type: "posts", Section: "posts"}, false, HTMLFormat},
		{"Section list for \"posts\" section with type set to \"blog\"", LayoutDescriptor{Kind: "section", Type: "blog", Section: "posts"}, false, HTMLFormat},
		{"Section list for \"posts\" section with layout set to \"demoLayout\"", LayoutDescriptor{Kind: "section", Layout: demoLayout, Section: "posts"}, false, HTMLFormat},

		{"Taxonomy list in categories", LayoutDescriptor{Kind: "taxonomy", Type: "categories", Section: "category"}, false, HTMLFormat},
		{"Taxonomy term in categories", LayoutDescriptor{Kind: "taxonomyTerm", Type: "categories", Section: "category"}, false, HTMLFormat},
	} {

		l := NewLayoutHandler(example.hasTheme)
		layouts, _ := l.For(example.d, example.f)

		basicExamples = append(basicExamples, Example{
			Example:      example.name,
			Kind:         example.d.Kind,
			OutputFormat: example.f.Name,
			Suffix:       example.f.MediaType.Suffix,
			Layouts:      makeLayoutsPresentable(layouts)})
	}

	return basicExamples

}

func makeLayoutsPresentable(l []string) []string {
	var filtered []string
	for _, ll := range l {
		ll = strings.TrimPrefix(ll, "_text/")
		if strings.Contains(ll, "theme/") {
			ll = strings.Replace(ll, "theme/", "demoTheme/layouts/", -1)
		} else {
			ll = "layouts/" + ll
		}
		if !strings.Contains(ll, "indexes") {
			filtered = append(filtered, ll)
		}
	}

	return filtered
}
