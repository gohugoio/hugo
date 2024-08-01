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
	"context"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cast"
)

type pageGroupTestObject struct {
	path   string
	weight int
	date   string
	param  string
}

var pageGroupTestSources = []pageGroupTestObject{
	{"/section1/testpage1.md", 3, "2012-04-06", "foo"},
	{"/section1/testpage2.md", 3, "2012-01-01", "bar"},
	{"/section1/testpage3.md", 2, "2012-04-06", "foo"},
	{"/section2/testpage4.md", 1, "2012-03-02", "bar"},
	// date might also be a full datetime:
	{"/section2/testpage5.md", 1, "2012-04-06T00:00:00Z", "baz"},
}

func preparePageGroupTestPages(t *testing.T) Pages {
	var pages Pages
	for _, src := range pageGroupTestSources {
		p := newTestPage()
		p.path = src.path
		if p.path != "" {
			p.section = strings.Split(strings.TrimPrefix(p.path, "/"), "/")[0]
		}
		p.weight = src.weight
		p.date = cast.ToTime(src.date)
		p.pubDate = cast.ToTime(src.date)
		p.expiryDate = cast.ToTime(src.date)
		p.lastMod = cast.ToTime(src.date).AddDate(3, 0, 0)
		p.params["custom_param"] = src.param
		p.params["custom_date"] = cast.ToTime(src.date)
		p.params["custom_string_date"] = src.date
		p.params["custom_object"] = map[string]any{
			"param":       src.param,
			"date":        cast.ToTime(src.date),
			"string_date": src.date,
		}
		pages = append(pages, p)
	}
	return pages
}

var comparePageGroup = qt.CmpEquals(cmp.Comparer(func(a, b Page) bool {
	return a == b
}))

func TestGroupByWithFieldNameArg(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: 1, Pages: Pages{pages[3], pages[4]}},
		{Key: 2, Pages: Pages{pages[2]}},
		{Key: 3, Pages: Pages{pages[0], pages[1]}},
	}

	groups, err := pages.GroupBy(context.Background(), "Weight")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByWithMethodNameArg(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "section1", Pages: Pages{pages[0], pages[1], pages[2]}},
		{Key: "section2", Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy(context.Background(), "Type")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByWithSectionArg(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "section1", Pages: Pages{pages[0], pages[1], pages[2]}},
		{Key: "section2", Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy(context.Background(), "Section")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: 3, Pages: Pages{pages[0], pages[1]}},
		{Key: 2, Pages: Pages{pages[2]}},
		{Key: 1, Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy(context.Background(), "Weight", "desc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByCalledWithEmptyPages(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	var pages Pages
	groups, err := pages.GroupBy(context.Background(), "Weight")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, qt.IsNil)
}

func TestReverse(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)

	groups1, err := pages.GroupBy(context.Background(), "Weight", "desc")
	c.Assert(err, qt.IsNil)

	groups2, err := pages.GroupBy(context.Background(), "Weight")
	c.Assert(err, qt.IsNil)

	groups2 = groups2.Reverse()
	c.Assert(groups2, comparePageGroup, groups1)
}

func TestGroupByParam(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "bar", Pages: Pages{pages[1], pages[3]}},
		{Key: "baz", Pages: Pages{pages[4]}},
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
	}

	groups, err := pages.GroupByParam("custom_param")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
		{Key: "baz", Pages: Pages{pages[4]}},
		{Key: "bar", Pages: Pages{pages[1], pages[3]}},
	}

	groups, err := pages.GroupByParam("custom_param", "desc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamCalledWithCapitalLetterString(t *testing.T) {
	c := qt.New(t)
	testStr := "TestString"
	p := newTestPage()
	p.params["custom_param"] = testStr
	pages := Pages{p}

	groups, err := pages.GroupByParam("custom_param")

	c.Assert(err, qt.IsNil)
	c.Assert(groups[0].Key, qt.DeepEquals, testStr)
}

func TestGroupByParamCalledWithSomeUnavailableParams(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	delete(pages[1].Params(), "custom_param")
	delete(pages[3].Params(), "custom_param")
	delete(pages[4].Params(), "custom_param")

	expect := PagesGroup{
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
	}

	groups, err := pages.GroupByParam("custom_param")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamCalledWithEmptyPages(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByParam("custom_param")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, qt.IsNil)
}

func TestGroupByParamCalledWithUnavailableParam(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	_, err := pages.GroupByParam("unavailable_param")
	c.Assert(err, qt.IsNil)
}

func TestGroupByParamNested(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)

	expect := PagesGroup{
		{Key: "bar", Pages: Pages{pages[1], pages[3]}},
		{Key: "baz", Pages: Pages{pages[4]}},
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
	}

	groups, err := pages.GroupByParam("custom_object.param")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByDate(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByDate("2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByDateInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByDate("2006-01", "asc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByPublishDate(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByPublishDate("2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByPublishDateInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByDate("2006-01", "asc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByPublishDateWithEmptyPages(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByPublishDate("2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, qt.IsNil)
}

func TestGroupByExpiryDate(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByExpiryDate("2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamDate(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByParamDate("custom_date", "2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamDateNested(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByParamDate("custom_object.date", "2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

// https://github.com/gohugoio/hugo/issues/3983
func TestGroupByParamDateWithStringParams(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByParamDate("custom_string_date", "2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamDateNestedWithStringParams(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByParamDate("custom_object.string_date", "2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByLastmod(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2015-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2015-03", Pages: Pages{pages[3]}},
		{Key: "2015-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByLastmod("2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByLastmodInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2015-01", Pages: Pages{pages[1]}},
		{Key: "2015-03", Pages: Pages{pages[3]}},
		{Key: "2015-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByLastmod("2006-01", "asc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamDateInReverseOrder(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByParamDate("custom_date", "2006-01", "asc")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, comparePageGroup, expect)
}

func TestGroupByParamDateWithEmptyPages(t *testing.T) {
	c := qt.New(t)

	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByParamDate("custom_date", "2006-01")
	c.Assert(err, qt.IsNil)
	c.Assert(groups, qt.IsNil)
}
