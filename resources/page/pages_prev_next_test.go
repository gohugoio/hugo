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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cast"
)

type pagePNTestObject struct {
	path   string
	weight int
	date   string
}

var pagePNTestSources = []pagePNTestObject{
	{"/section1/testpage1.md", 5, "2012-04-06"},
	{"/section1/testpage2.md", 4, "2012-01-01"},
	{"/section1/testpage3.md", 3, "2012-04-06"},
	{"/section2/testpage4.md", 2, "2012-03-02"},
	{"/section2/testpage5.md", 1, "2012-04-06"},
}

func TestPrev(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := preparePageGroupTestPages(t)
	c.Assert(pages[4], qt.Equals, pages.Prev(pages[0]))
	c.Assert(pages[0], qt.Equals, pages.Prev(pages[1]))
	c.Assert(pages[3], qt.Equals, pages.Prev(pages[4]))
}

func TestNext(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := preparePageGroupTestPages(t)
	c.Assert(pages[1], qt.Equals, pages.Next(pages[0]))
	c.Assert(pages[2], qt.Equals, pages.Next(pages[1]))
	c.Assert(pages[0], qt.Equals, pages.Next(pages[4]))
}

func prepareWeightedPagesPrevNext(t *testing.T) WeightedPages {
	w := WeightedPages{}

	for _, src := range pagePNTestSources {
		p := newTestPage()
		p.path = src.path
		p.weight = src.weight
		p.date = cast.ToTime(src.date)
		p.pubDate = cast.ToTime(src.date)
		w = append(w, WeightedPage{Weight: p.weight, Page: p})
	}

	w.Sort()
	return w
}

func TestWeightedPagesPrev(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	w := prepareWeightedPagesPrevNext(t)
	c.Assert(w[4].Page, qt.Equals, w.Prev(w[0].Page))
	c.Assert(w[0].Page, qt.Equals, w.Prev(w[1].Page))
	c.Assert(w[3].Page, qt.Equals, w.Prev(w[4].Page))
}

func TestWeightedPagesNext(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	w := prepareWeightedPagesPrevNext(t)
	c.Assert(w[1].Page, qt.Equals, w.Next(w[0].Page))
	c.Assert(w[2].Page, qt.Equals, w.Next(w[1].Page))
	c.Assert(w[0].Page, qt.Equals, w.Next(w[4].Page))
}
