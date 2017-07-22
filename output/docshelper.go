package output

import (
	"strings"

	"fmt"

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
		name           string
		d              LayoutDescriptor
		hasTheme       bool
		layoutOverride string
		f              Format
	}{
		{`AMP home, with theme "demoTheme".`, LayoutDescriptor{Kind: "home"}, true, "", AMPFormat},
		{`AMP home, French language".`, LayoutDescriptor{Kind: "home", Lang: "fr"}, false, "", AMPFormat},
		{"RSS home, no theme.", LayoutDescriptor{Kind: "home"}, false, "", RSSFormat},
		{"JSON home, no theme.", LayoutDescriptor{Kind: "home"}, false, "", JSONFormat},
		{fmt.Sprintf(`CSV regular, "layout: %s" in front matter.`, demoLayout), LayoutDescriptor{Kind: "page", Layout: demoLayout}, false, "", CSVFormat},
		{fmt.Sprintf(`JSON regular, "type: %s" in front matter.`, demoType), LayoutDescriptor{Kind: "page", Type: demoType}, false, "", JSONFormat},
		{"HTML regular.", LayoutDescriptor{Kind: "page"}, false, "", HTMLFormat},
		{"AMP regular.", LayoutDescriptor{Kind: "page"}, false, "", AMPFormat},
		{"Calendar blog section.", LayoutDescriptor{Kind: "section", Section: "blog"}, false, "", CalendarFormat},
		{"Calendar taxonomy list.", LayoutDescriptor{Kind: "taxonomy", Section: "tag"}, false, "", CalendarFormat},
		{"Calendar taxonomy term.", LayoutDescriptor{Kind: "taxonomyTerm", Section: "tag"}, false, "", CalendarFormat},
	} {

		l := NewLayoutHandler(example.hasTheme)
		layouts, _ := l.For(example.d, example.layoutOverride, example.f)

		basicExamples = append(basicExamples, Example{
			Example:      example.name,
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
