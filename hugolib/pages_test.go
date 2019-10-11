package hugolib

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
)

func newPagesPrevNextTestSite(t testing.TB, numPages int) *sitesBuilder {
	pageTemplate := `
---
title: "Page %d"
weight: %d
---

`
	b := newTestSitesBuilder(t)

	for i := 1; i <= numPages; i++ {
		b.WithContent(fmt.Sprintf("page%d.md", i), fmt.Sprintf(pageTemplate, i, rand.Intn(numPages)))
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
		Variant{".Next", nil, func(p page.Page, pages page.Pages) { p.Next() }},
		Variant{".Prev", nil, func(p page.Page, pages page.Pages) { p.Prev() }},
		Variant{"Pages.Next", nil, func(p page.Page, pages page.Pages) { pages.Next(p) }},
		Variant{"Pages.Prev", nil, func(p page.Page, pages page.Pages) { pages.Prev(p) }},
		Variant{"Pages.Shuffled.Next", shufflePages, func(p page.Page, pages page.Pages) { pages.Next(p) }},
		Variant{"Pages.Shuffled.Prev", shufflePages, func(p page.Page, pages page.Pages) { pages.Prev(p) }},
		Variant{"Pages.ByTitle.Next", func(pages page.Pages) page.Pages { return pages.ByTitle() }, func(p page.Page, pages page.Pages) { pages.Next(p) }},
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
