// Copyright 2015 The Hugo Authors. All rights reserved.
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

func TestPrev(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	assert.Equal(t, pages.Prev(pages[0]), pages[4])
	assert.Equal(t, pages.Prev(pages[1]), pages[0])
	assert.Equal(t, pages.Prev(pages[4]), pages[3])
}

func TestNext(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	assert.Equal(t, pages.Next(pages[0]), pages[1])
	assert.Equal(t, pages.Next(pages[1]), pages[2])
	assert.Equal(t, pages.Next(pages[4]), pages[0])
}

func prepareWeightedPagesPrevNext(t *testing.T) WeightedPages {
	s := newTestSite(t)
	w := WeightedPages{}

	for _, src := range pagePNTestSources {
		p, err := s.NewPage(src.path)
		if err != nil {
			t.Fatalf("failed to prepare test page %s", src.path)
		}
		p.Weight = src.weight
		p.Date = cast.ToTime(src.date)
		p.PublishDate = cast.ToTime(src.date)
		w = append(w, WeightedPage{p.Weight, p})
	}

	w.Sort()
	return w
}

func TestWeightedPagesPrev(t *testing.T) {
	t.Parallel()
	w := prepareWeightedPagesPrevNext(t)
	assert.Equal(t, w.Prev(w[0].Page), w[4].Page)
	assert.Equal(t, w.Prev(w[1].Page), w[0].Page)
	assert.Equal(t, w.Prev(w[4].Page), w[3].Page)
}

func TestWeightedPagesNext(t *testing.T) {
	t.Parallel()
	w := prepareWeightedPagesPrevNext(t)
	assert.Equal(t, w.Next(w[0].Page), w[1].Page)
	assert.Equal(t, w.Next(w[1].Page), w[2].Page)
	assert.Equal(t, w.Next(w[4].Page), w[0].Page)
}
