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
	"reflect"
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
	{"/section2/testpage5.md", 1, "2012-04-06", "baz"},
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
		p.params["custom_param"] = src.param
		p.params["custom_date"] = cast.ToTime(src.date)
		pages = append(pages, p)
	}
	return pages
}

func TestGroupByWithFieldNameArg(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: 1, Pages: Pages{pages[3], pages[4]}},
		{Key: 2, Pages: Pages{pages[2]}},
		{Key: 3, Pages: Pages{pages[0], pages[1]}},
	}

	groups, err := pages.GroupBy("Weight")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByWithMethodNameArg(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "section1", Pages: Pages{pages[0], pages[1], pages[2]}},
		{Key: "section2", Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy("Type")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByWithSectionArg(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "section1", Pages: Pages{pages[0], pages[1], pages[2]}},
		{Key: "section2", Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy("Section")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be\n%#v, got\n%#v", expect, groups)
	}
}

func TestGroupByInReverseOrder(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: 3, Pages: Pages{pages[0], pages[1]}},
		{Key: 2, Pages: Pages{pages[2]}},
		{Key: 1, Pages: Pages{pages[3], pages[4]}},
	}

	groups, err := pages.GroupBy("Weight", "desc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByCalledWithEmptyPages(t *testing.T) {
	t.Parallel()
	var pages Pages
	groups, err := pages.GroupBy("Weight")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if groups != nil {
		t.Errorf("PagesGroup isn't empty. It should be %#v, got %#v", nil, groups)
	}
}

func TestGroupByParamCalledWithUnavailableKey(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	_, err := pages.GroupByParam("UnavailableKey")
	if err == nil {
		t.Errorf("GroupByParam should return an error but didn't")
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)

	groups1, err := pages.GroupBy("Weight", "desc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}

	groups2, err := pages.GroupBy("Weight")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	groups2 = groups2.Reverse()

	if !reflect.DeepEqual(groups2, groups1) {
		t.Errorf("PagesGroup is sorted in unexpected order. It should be %#v, got %#v", groups2, groups1)
	}
}

func TestGroupByParam(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "bar", Pages: Pages{pages[1], pages[3]}},
		{Key: "baz", Pages: Pages{pages[4]}},
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
	}

	groups, err := pages.GroupByParam("custom_param")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamInReverseOrder(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
		{Key: "baz", Pages: Pages{pages[4]}},
		{Key: "bar", Pages: Pages{pages[1], pages[3]}},
	}

	groups, err := pages.GroupByParam("custom_param", "desc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamCalledWithCapitalLetterString(t *testing.T) {
	c := qt.New(t)
	testStr := "TestString"
	p := newTestPage()
	p.params["custom_param"] = testStr
	pages := Pages{p}

	groups, err := pages.GroupByParam("custom_param")

	c.Assert(err, qt.IsNil)
	c.Assert(groups[0].Key, qt.Equals, testStr)

}

func TestGroupByParamCalledWithSomeUnavailableParams(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	delete(pages[1].Params(), "custom_param")
	delete(pages[3].Params(), "custom_param")
	delete(pages[4].Params(), "custom_param")

	expect := PagesGroup{
		{Key: "foo", Pages: Pages{pages[0], pages[2]}},
	}

	groups, err := pages.GroupByParam("custom_param")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamCalledWithEmptyPages(t *testing.T) {
	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByParam("custom_param")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if groups != nil {
		t.Errorf("PagesGroup isn't empty. It should be %#v, got %#v", nil, groups)
	}
}

func TestGroupByParamCalledWithUnavailableParam(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	_, err := pages.GroupByParam("unavailable_param")
	if err == nil {
		t.Errorf("GroupByParam should return an error but didn't")
	}
}

func TestGroupByDate(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByDate("2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByDateInReverseOrder(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByDate("2006-01", "asc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByPublishDate(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByPublishDate("2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByPublishDateInReverseOrder(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByDate("2006-01", "asc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByPublishDateWithEmptyPages(t *testing.T) {
	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByPublishDate("2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if groups != nil {
		t.Errorf("PagesGroup isn't empty. It should be %#v, got %#v", nil, groups)
	}
}

func TestGroupByExpiryDate(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByExpiryDate("2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamDate(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-04", Pages: Pages{pages[4], pages[2], pages[0]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-01", Pages: Pages{pages[1]}},
	}

	groups, err := pages.GroupByParamDate("custom_date", "2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamDateInReverseOrder(t *testing.T) {
	t.Parallel()
	pages := preparePageGroupTestPages(t)
	expect := PagesGroup{
		{Key: "2012-01", Pages: Pages{pages[1]}},
		{Key: "2012-03", Pages: Pages{pages[3]}},
		{Key: "2012-04", Pages: Pages{pages[0], pages[2], pages[4]}},
	}

	groups, err := pages.GroupByParamDate("custom_date", "2006-01", "asc")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if !reflect.DeepEqual(groups, expect) {
		t.Errorf("PagesGroup has unexpected groups. It should be %#v, got %#v", expect, groups)
	}
}

func TestGroupByParamDateWithEmptyPages(t *testing.T) {
	t.Parallel()
	var pages Pages
	groups, err := pages.GroupByParamDate("custom_date", "2006-01")
	if err != nil {
		t.Fatalf("Unable to make PagesGroup array: %s", err)
	}
	if groups != nil {
		t.Errorf("PagesGroup isn't empty. It should be %#v, got %#v", nil, groups)
	}
}
