package output

import (
	"strings"

	//	"fmt"

	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/output/layouts"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() docshelper.DocProvider {
		return docshelper.DocProvider{
			"output": map[string]any{
				"layouts": createLayoutExamples(),
			},
		}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}

func createLayoutExamples() any {
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
		d    layouts.LayoutDescriptor
	}{
		// Taxonomy layouts.LayoutDescriptor={categories category taxonomy en  false Type Section
		{"Single page in \"posts\" section", layouts.LayoutDescriptor{Kind: "page", Type: "posts", OutputFormatName: "html", Suffix: "html"}},
		{"Base template for single page in \"posts\" section", layouts.LayoutDescriptor{Baseof: true, Kind: "page", Type: "posts", OutputFormatName: "html", Suffix: "html"}},
		{"Single page in \"posts\" section with layout set to \"demolayout\"", layouts.LayoutDescriptor{Kind: "page", Type: "posts", Layout: demoLayout, OutputFormatName: "html", Suffix: "html"}},
		{"Base template for single page in \"posts\" section with layout set to \"demolayout\"", layouts.LayoutDescriptor{Baseof: true, Kind: "page", Type: "posts", Layout: demoLayout, OutputFormatName: "html", Suffix: "html"}},
		{"AMP single page in \"posts\" section", layouts.LayoutDescriptor{Kind: "page", Type: "posts", OutputFormatName: "amp", Suffix: "html"}},
		{"AMP single page in \"posts\" section, French language", layouts.LayoutDescriptor{Kind: "page", Type: "posts", Lang: "fr", OutputFormatName: "amp", Suffix: "html"}},
		// Typeless pages get "page" as type
		{"Home page", layouts.LayoutDescriptor{Kind: "home", Type: "page", OutputFormatName: "html", Suffix: "html"}},
		{"Base template for home page", layouts.LayoutDescriptor{Baseof: true, Kind: "home", Type: "page", OutputFormatName: "html", Suffix: "html"}},
		{"Home page with type set to \"demotype\"", layouts.LayoutDescriptor{Kind: "home", Type: demoType, OutputFormatName: "html", Suffix: "html"}},
		{"Base template for home page with type set to \"demotype\"", layouts.LayoutDescriptor{Baseof: true, Kind: "home", Type: demoType, OutputFormatName: "html", Suffix: "html"}},
		{"Home page with layout set to \"demolayout\"", layouts.LayoutDescriptor{Kind: "home", Type: "page", Layout: demoLayout, OutputFormatName: "html", Suffix: "html"}},
		{"AMP home, French language", layouts.LayoutDescriptor{Kind: "home", Type: "page", Lang: "fr", OutputFormatName: "amp", Suffix: "html"}},
		{"JSON home", layouts.LayoutDescriptor{Kind: "home", Type: "page", OutputFormatName: "json", Suffix: "json"}},
		{"RSS home", layouts.LayoutDescriptor{Kind: "home", Type: "page", OutputFormatName: "rss", Suffix: "xml"}},

		{"Section list for \"posts\"", layouts.LayoutDescriptor{Kind: "section", Type: "posts", Section: "posts", OutputFormatName: "html", Suffix: "html"}},
		{"Section list for \"posts\" with type set to \"blog\"", layouts.LayoutDescriptor{Kind: "section", Type: "blog", Section: "posts", OutputFormatName: "html", Suffix: "html"}},
		{"Section list for \"posts\" with layout set to \"demolayout\"", layouts.LayoutDescriptor{Kind: "section", Layout: demoLayout, Section: "posts", OutputFormatName: "html", Suffix: "html"}},
		{"Section list for \"posts\"", layouts.LayoutDescriptor{Kind: "section", Type: "posts", OutputFormatName: "rss", Suffix: "xml"}},

		{"Taxonomy list for \"categories\"", layouts.LayoutDescriptor{Kind: "taxonomy", Type: "categories", Section: "category", OutputFormatName: "html", Suffix: "html"}},
		{"Taxonomy list for \"categories\"", layouts.LayoutDescriptor{Kind: "taxonomy", Type: "categories", Section: "category", OutputFormatName: "rss", Suffix: "xml"}},

		{"Term list for \"categories\"", layouts.LayoutDescriptor{Kind: "term", Type: "categories", Section: "category", OutputFormatName: "html", Suffix: "html"}},
		{"Term list for \"categories\"", layouts.LayoutDescriptor{Kind: "term", Type: "categories", Section: "category", OutputFormatName: "rss", Suffix: "xml"}},
	} {

		l := layouts.NewLayoutHandler()
		layouts, _ := l.For(example.d)

		basicExamples = append(basicExamples, Example{
			Example:      example.name,
			Kind:         example.d.Kind,
			OutputFormat: example.d.OutputFormatName,
			Suffix:       example.d.Suffix,
			Layouts:      makeLayoutsPresentable(layouts),
		})
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
