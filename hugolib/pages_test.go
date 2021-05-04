package hugolib

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
)

func newPagesPrevNextTestSite(t testing.TB, numPages int) *sitesBuilder {
	categories := []string{"blue", "green", "red", "orange", "indigo", "amber", "lime"}
	cat1, cat2 := categories[rand.Intn(len(categories))], categories[rand.Intn(len(categories))]
	categoriesSlice := fmt.Sprintf("[%q,%q]", cat1, cat2)
	pageTemplate := `
---
title: "Page %d"
weight: %d
categories: %s
---

`
	b := newTestSitesBuilder(t)

	for i := 1; i <= numPages; i++ {
		b.WithContent(fmt.Sprintf("page%d.md", i), fmt.Sprintf(pageTemplate, i, rand.Intn(numPages), categoriesSlice))
	}

	return b
}

func TestPagesPrevNext(t *testing.T) {
	b := newPagesPrevNextTestSite(t, 100)
	b.Build(BuildCfg{SkipRender: true})

	pages := b.H.Sites[0].RegularPages()

	b.Assert(pages, qt.HasLen, 100)

	for _, p := range pages {
		msg := qt.Commentf("w=%d", p.Weight())
		b.Assert(pages.Next(p), qt.Equals, p.Next(), msg)
		b.Assert(pages.Prev(p), qt.Equals, p.Prev(), msg)
	}
}

func BenchmarkPagesPrevNext(b *testing.B) {
	type Variant struct {
		name         string
		preparePages func(pages page.Pages) page.Pages
		run          func(p page.Page, pages page.Pages)
	}

	shufflePages := func(pages page.Pages) page.Pages {
		rand.Shuffle(len(pages), func(i, j int) { pages[i], pages[j] = pages[j], pages[i] })
		return pages
	}

	for _, variant := range []Variant{
		{".Next", nil, func(p page.Page, pages page.Pages) { p.Next() }},
		{".Prev", nil, func(p page.Page, pages page.Pages) { p.Prev() }},
		{"Pages.Next", nil, func(p page.Page, pages page.Pages) { pages.Next(p) }},
		{"Pages.Prev", nil, func(p page.Page, pages page.Pages) { pages.Prev(p) }},
		{"Pages.Shuffled.Next", shufflePages, func(p page.Page, pages page.Pages) { pages.Next(p) }},
		{"Pages.Shuffled.Prev", shufflePages, func(p page.Page, pages page.Pages) { pages.Prev(p) }},
		{"Pages.ByTitle.Next", func(pages page.Pages) page.Pages { return pages.ByTitle() }, func(p page.Page, pages page.Pages) { pages.Next(p) }},
	} {
		for _, numPages := range []int{300, 5000} {
			b.Run(fmt.Sprintf("%s-pages-%d", variant.name, numPages), func(b *testing.B) {
				b.StopTimer()
				builder := newPagesPrevNextTestSite(b, numPages)
				builder.Build(BuildCfg{SkipRender: true})
				pages := builder.H.Sites[0].RegularPages()
				if variant.preparePages != nil {
					pages = variant.preparePages(pages)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					p := pages[rand.Intn(len(pages))]
					variant.run(p, pages)
				}
			})
		}
	}
}

func BenchmarkPagePageCollections(b *testing.B) {
	type Variant struct {
		name string
		run  func(p page.Page)
	}

	for _, variant := range []Variant{
		{".Pages", func(p page.Page) { p.Pages() }},
		{".RegularPages", func(p page.Page) { p.RegularPages() }},
		{".RegularPagesRecursive", func(p page.Page) { p.RegularPagesRecursive() }},
	} {
		for _, numPages := range []int{300, 5000} {
			b.Run(fmt.Sprintf("%s-%d", variant.name, numPages), func(b *testing.B) {
				b.StopTimer()
				builder := newPagesPrevNextTestSite(b, numPages)
				builder.Build(BuildCfg{SkipRender: true})
				var pages page.Pages
				for _, p := range builder.H.Sites[0].Pages() {
					if !p.IsPage() {
						pages = append(pages, p)
					}
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					p := pages[rand.Intn(len(pages))]
					variant.run(p)
				}
			})
		}
	}
}
