// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
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

func preparePagePNTestPages(t *testing.T) Pages {
	var pages Pages
	for _, s := range pagePNTestSources {
		p, err := NewPage(s.path)
		if err != nil {
			t.Fatalf("failed to prepare test page %s", s.path)
		}
		p.Weight = s.weight
		p.Date = cast.ToTime(s.date)
		p.PublishDate = cast.ToTime(s.date)
		pages = append(pages, p)
	}
	return pages
}

func TestPrev(t *testing.T) {
	pages := preparePageGroupTestPages(t)
	assert.Equal(t, pages.Prev(pages[0]), pages[4])
	assert.Equal(t, pages.Prev(pages[1]), pages[0])
	assert.Equal(t, pages.Prev(pages[4]), pages[3])
}

func TestNext(t *testing.T) {
	pages := preparePageGroupTestPages(t)
	assert.Equal(t, pages.Next(pages[0]), pages[1])
	assert.Equal(t, pages.Next(pages[1]), pages[2])
	assert.Equal(t, pages.Next(pages[4]), pages[0])
}

func prepareWeightedPagesPrevNext(t *testing.T) WeightedPages {
	w := WeightedPages{}

	for _, s := range pagePNTestSources {
		p, err := NewPage(s.path)
		if err != nil {
			t.Fatalf("failed to prepare test page %s", s.path)
		}
		p.Weight = s.weight
		p.Date = cast.ToTime(s.date)
		p.PublishDate = cast.ToTime(s.date)
		w = append(w, WeightedPage{p.Weight, p})
	}

	w.Sort()
	return w
}

func TestWeightedPagesPrev(t *testing.T) {
	w := prepareWeightedPagesPrevNext(t)
	assert.Equal(t, w.Prev(w[0].Page), w[4].Page)
	assert.Equal(t, w.Prev(w[1].Page), w[0].Page)
	assert.Equal(t, w.Prev(w[4].Page), w[3].Page)
}

func TestWeightedPagesNext(t *testing.T) {
	w := prepareWeightedPagesPrevNext(t)
	assert.Equal(t, w.Next(w[0].Page), w[1].Page)
	assert.Equal(t, w.Next(w[1].Page), w[2].Page)
	assert.Equal(t, w.Next(w[4].Page), w[0].Page)
}
