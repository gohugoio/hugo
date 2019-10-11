// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestSearchPage(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := createSortTestPages(10)
	for i, p := range pages {
		p.(*testPage).title = fmt.Sprintf("Title %d", i%2)
	}

	for _, pages := range []Pages{pages.ByTitle(), pages.ByTitle().Reverse()} {
		less := isPagesProbablySorted(pages, lessPageTitle)
		c.Assert(less, qt.Not(qt.IsNil))
		for i, p := range pages {
			idx := searchPageBinary(p, pages, less)
			c.Assert(idx, qt.Equals, i)
		}
	}

}

func BenchmarkSearchPage(b *testing.B) {
	type Variant struct {
		name         string
		preparePages func(pages Pages) Pages
		search       func(p Page, pages Pages) int
	}

	shufflePages := func(pages Pages) Pages {
		rand.Shuffle(len(pages), func(i, j int) { pages[i], pages[j] = pages[j], pages[i] })
		return pages
	}

	linearSearch := func(p Page, pages Pages) int {
		return searchPageLinear(p, pages, 0)
	}

	createPages := func(num int) Pages {
		pages := createSortTestPages(num)
		for _, p := range pages {
			tp := p.(*testPage)
			tp.weight = rand.Intn(len(pages))
			tp.title = fmt.Sprintf("Title %d", rand.Intn(len(pages)))

			tp.pubDate = time.Now().Add(time.Duration(rand.Intn(len(pages)/5)) * time.Hour)
			tp.date = time.Now().Add(time.Duration(rand.Intn(len(pages)/5)) * time.Hour)
		}

		return pages
	}

	for _, variant := range []Variant{
		Variant{"Shuffled", shufflePages, searchPage},
		Variant{"ByWeight", func(pages Pages) Pages {
			return pages.ByWeight()
		}, searchPage},
		Variant{"ByWeight.Reverse", func(pages Pages) Pages {
			return pages.ByWeight().Reverse()
		}, searchPage},
		Variant{"ByDate", func(pages Pages) Pages {
			return pages.ByDate()
		}, searchPage},
		Variant{"ByPublishDate", func(pages Pages) Pages {
			return pages.ByPublishDate()
		}, searchPage},
		Variant{"ByTitle", func(pages Pages) Pages {
			return pages.ByTitle()
		}, searchPage},
		Variant{"ByTitle Linear", func(pages Pages) Pages {
			return pages.ByTitle()
		}, linearSearch},
	} {
		for _, numPages := range []int{100, 500, 1000, 5000} {
			b.Run(fmt.Sprintf("%s-%d", variant.name, numPages), func(b *testing.B) {
				b.StopTimer()
				pages := createPages(numPages)
				if variant.preparePages != nil {
					pages = variant.preparePages(pages)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					j := rand.Intn(numPages)
					k := variant.search(pages[j], pages)
					if k != j {
						b.Fatalf("%d != %d", k, j)
					}
				}
			})
		}
	}
}

func TestIsPagesProbablySorted(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(isPagesProbablySorted(createSortTestPages(6).ByWeight(), DefaultPageSort), qt.Not(qt.IsNil))
	c.Assert(isPagesProbablySorted(createSortTestPages(300).ByWeight(), DefaultPageSort), qt.Not(qt.IsNil))
	c.Assert(isPagesProbablySorted(createSortTestPages(6), DefaultPageSort), qt.IsNil)
	c.Assert(isPagesProbablySorted(createSortTestPages(300).ByTitle(), pageLessFunctions...), qt.Not(qt.IsNil))

}
