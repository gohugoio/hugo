package tplimpl

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/kinds"
)

func TestTemplateDescriptorCompare(t *testing.T) {
	c := qt.New(t)

	dh := descriptorHandler{
		opts: StoreOptions{
			OutputFormats:       output.DefaultFormats,
			DefaultOutputFormat: "html",
		},
	}

	less := func(category Category, this, other1, other2 TemplateDescriptor) {
		c.Helper()
		result1 := dh.compareDescriptors(category, this, other1)
		result2 := dh.compareDescriptors(category, this, other2)
		c.Assert(result1.w1 < result2.w1, qt.IsTrue, qt.Commentf("%d < %d", result1, result2))
	}

	check := func(category Category, this, other TemplateDescriptor, less bool) {
		c.Helper()
		result := dh.compareDescriptors(category, this, other)
		if less {
			c.Assert(result.w1 < 0, qt.IsTrue, qt.Commentf("%d", result))
		} else {
			c.Assert(result.w1 >= 0, qt.IsTrue, qt.Commentf("%d", result))
		}
	}

	check(

		CategoryBaseof,
		TemplateDescriptor{Kind: "", Layout: "", Lang: "", OutputFormat: "404", MediaType: "text/html"},
		TemplateDescriptor{Kind: "", Layout: "", Lang: "", OutputFormat: "html", MediaType: "text/html"},
		false,
	)

	check(
		CategoryLayout,
		TemplateDescriptor{Kind: "", Lang: "en", OutputFormat: "404", MediaType: "text/html"},
		TemplateDescriptor{Kind: "", Layout: "", Lang: "", OutputFormat: "alias", MediaType: "text/html"},
		true,
	)

	less(
		CategoryLayout,
		TemplateDescriptor{Kind: kinds.KindHome, Layout: "list", OutputFormat: "html"},
		TemplateDescriptor{Layout: "list", OutputFormat: "html"},
		TemplateDescriptor{Kind: kinds.KindHome, OutputFormat: "html"},
	)

	check(
		CategoryLayout,
		TemplateDescriptor{Kind: kinds.KindHome, Layout: "list", OutputFormat: "html", MediaType: "text/html"},
		TemplateDescriptor{Kind: kinds.KindHome, Layout: "list", OutputFormat: "myformat", MediaType: "text/html"},
		false,
	)
}

// INFO  timer:  name resolveTemplate count 779 duration 5.482274ms average 7.037µs median 4µs
func BenchmarkCompareDescriptors(b *testing.B) {
	dh := descriptorHandler{
		opts: StoreOptions{
			OutputFormats:       output.DefaultFormats,
			DefaultOutputFormat: "html",
		},
	}

	pairs := []struct {
		d1, d2 TemplateDescriptor
	}{
		{
			TemplateDescriptor{Kind: "", Layout: "", OutputFormat: "404", MediaType: "text/html", Lang: "en", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
			TemplateDescriptor{Kind: "", Layout: "", OutputFormat: "rss", MediaType: "application/rss+xml", Lang: "", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
		},
		{
			TemplateDescriptor{Kind: "page", Layout: "single", OutputFormat: "html", MediaType: "text/html", Lang: "en", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
			TemplateDescriptor{Kind: "", Layout: "list", OutputFormat: "", MediaType: "application/rss+xml", Lang: "", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
		},
		{
			TemplateDescriptor{Kind: "page", Layout: "single", OutputFormat: "html", MediaType: "text/html", Lang: "en", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
			TemplateDescriptor{Kind: "", Layout: "", OutputFormat: "alias", MediaType: "text/html", Lang: "", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
		},
		{
			TemplateDescriptor{Kind: "page", Layout: "single", OutputFormat: "rss", MediaType: "application/rss+xml", Lang: "en", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
			TemplateDescriptor{Kind: "", Layout: "single", OutputFormat: "rss", MediaType: "application/rss+xml", Lang: "nn", Variant1: "", Variant2: "", LayoutMustMatch: false, IsPlainText: false},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pair := range pairs {
			_ = dh.compareDescriptors(CategoryLayout, pair.d1, pair.d2)
		}
	}
}
