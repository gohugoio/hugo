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

package collections

import (
	"errors"
	"fmt"
	"html/template"

	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/maps"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type tstNoStringer struct{}

func TestAfter(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		index  interface{}
		seq    interface{}
		expect interface{}
	}{
		{int(2), []string{"a", "b", "c", "d"}, []string{"c", "d"}},
		{int32(3), []string{"a", "b"}, []string{}},
		{int64(2), []int{100, 200, 300}, []int{300}},
		{100, []int{100, 200}, []int{}},
		{"1", []int{100, 200, 300}, []int{200, 300}},
		{0, []int{100, 200, 300, 400, 500}, []int{100, 200, 300, 400, 500}},
		{0, []string{"a", "b", "c", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{2, []string{}, []string{}},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.After(test.index, test.seq)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.DeepEquals, test.expect, errMsg)
	}
}

type tstGrouper struct {
}

type tstGroupers []*tstGrouper

func (g tstGrouper) Group(key interface{}, items interface{}) (interface{}, error) {
	ilen := reflect.ValueOf(items).Len()
	return fmt.Sprintf("%v(%d)", key, ilen), nil
}

type tstGrouper2 struct {
}

func (g *tstGrouper2) Group(key interface{}, items interface{}) (interface{}, error) {
	ilen := reflect.ValueOf(items).Len()
	return fmt.Sprintf("%v(%d)", key, ilen), nil
}

func TestGroup(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	for i, test := range []struct {
		key    interface{}
		items  interface{}
		expect interface{}
	}{
		{"a", []*tstGrouper{{}, {}}, "a(2)"},
		{"b", tstGroupers{&tstGrouper{}, &tstGrouper{}}, "b(2)"},
		{"a", []tstGrouper{{}, {}}, "a(2)"},
		{"a", []*tstGrouper2{{}, {}}, "a(2)"},
		{"b", []tstGrouper2{{}, {}}, "b(2)"},
		{"a", []*tstGrouper{}, "a(0)"},
		{"a", []string{"a", "b"}, false},
		{"a", "asdf", false},
		{"a", nil, false},
		{nil, []*tstGrouper{{}, {}}, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Group(test.key, test.items)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestDelimit(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		seq       interface{}
		delimiter interface{}
		last      interface{}
		expect    template.HTML
	}{
		{[]string{"class1", "class2", "class3"}, " ", nil, "class1 class2 class3"},
		{[]int{1, 2, 3, 4, 5}, ",", nil, "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", nil, "1, 2, 3, 4, 5"},
		{[]string{"class1", "class2", "class3"}, " ", " and ", "class1 class2 and class3"},
		{[]int{1, 2, 3, 4, 5}, ",", ",", "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", ", and ", "1, 2, 3, 4, and 5"},
		// test maps with and without sorting required
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", nil, "10--20--30--40--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", nil, "10--20--30--40--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", nil, "50--40--10--30--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", nil, "10--20--30--40--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", nil, "30--20--10--40--50"},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, "--", nil, "30--20--10--40--50"},
		// test maps with a last delimiter
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", "--and--", "50--40--10--30--and--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[float64]string{3.5: "10", 2.5: "20", 1.5: "30", 4.5: "40", 5.5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		var result template.HTML
		var err error

		if test.last == nil {
			result, err = ns.Delimit(test.seq, test.delimiter)
		} else {
			result, err = ns.Delimit(test.seq, test.delimiter, test.last)
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestDictionary(t *testing.T) {
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		values []interface{}
		expect interface{}
	}{
		{[]interface{}{"a", "b"}, map[string]interface{}{"a": "b"}},
		{[]interface{}{[]string{"a", "b"}, "c"}, map[string]interface{}{"a": map[string]interface{}{"b": "c"}}},
		{[]interface{}{[]string{"a", "b"}, "c", []string{"a", "b2"}, "c2", "b", "c"},
			map[string]interface{}{"a": map[string]interface{}{"b": "c", "b2": "c2"}, "b": "c"}},
		{[]interface{}{"a", 12, "b", []int{4}}, map[string]interface{}{"a": 12, "b": []int{4}}},
		// errors
		{[]interface{}{5, "b"}, false},
		{[]interface{}{"a", "b", "c"}, false},
	} {
		i := i
		test := test
		c.Run(fmt.Sprint(i), func(c *qt.C) {
			c.Parallel()
			errMsg := qt.Commentf("[%d] %v", i, test.values)

			result, err := ns.Dictionary(test.values...)

			if b, ok := test.expect.(bool); ok && !b {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}

			c.Assert(err, qt.IsNil, errMsg)
			c.Assert(result, qt.DeepEquals, test.expect, qt.Commentf(fmt.Sprint(result)))
		})
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	s := []string{"a", "b", "c"}
	reversed, err := ns.Reverse(s)
	c.Assert(err, qt.IsNil)
	c.Assert(reversed, qt.DeepEquals, []string{"c", "b", "a"}, qt.Commentf(fmt.Sprint(reversed)))
	c.Assert(s, qt.DeepEquals, []string{"a", "b", "c"})

	reversed, err = ns.Reverse(nil)
	c.Assert(err, qt.IsNil)
	c.Assert(reversed, qt.IsNil)
	_, err = ns.Reverse(43)
	c.Assert(err, qt.Not(qt.IsNil))

}

func TestEchoParam(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		a      interface{}
		key    interface{}
		expect interface{}
	}{
		{[]int{1, 2, 3}, 1, int64(2)},
		{[]uint{1, 2, 3}, 1, uint64(2)},
		{[]float64{1.1, 2.2, 3.3}, 1, float64(2.2)},
		{[]string{"foo", "bar", "baz"}, 1, "bar"},
		{[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}}, 1, ""},
		{map[string]int{"foo": 1, "bar": 2, "baz": 3}, "bar", int64(2)},
		{map[string]uint{"foo": 1, "bar": 2, "baz": 3}, "bar", uint64(2)},
		{map[string]float64{"foo": 1.1, "bar": 2.2, "baz": 3.3}, "bar", float64(2.2)},
		{map[string]string{"foo": "FOO", "bar": "BAR", "baz": "BAZ"}, "bar", "BAR"},
		{map[string]TstX{"foo": {A: "a", B: "b"}, "bar": {A: "c", B: "d"}, "baz": {A: "e", B: "f"}}, "bar", ""},
		{map[string]interface{}{"foo": nil}, "foo", ""},
		{(*[]string)(nil), "bar", ""},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result := ns.EchoParam(test.a, test.key)

		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		limit  interface{}
		seq    interface{}
		expect interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"a", "b"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{100, 200}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{100}},
		{0, []string{"h", "u", "g", "o"}, []string{}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.First(test.limit, test.seq)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, test.expect, errMsg)
	}
}

func TestIn(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		l1     interface{}
		l2     interface{}
		expect bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]interface{}{"a", "b", "c"}, "b", true},
		{[]interface{}{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "12", "c"}, 12, false},
		{[]string{"a", "b", "c"}, nil, false},
		{[]int{1, 2, 4}, 2, true},
		{[]interface{}{1, 2, 4}, 2, true},
		{[]interface{}{1, 2, 4}, nil, false},
		{[]interface{}{nil}, nil, false},
		{[]int{1, 2, 4}, 3, false},
		{[]float64{1.23, 2.45, 4.67}, 1.23, true},
		{[]float64{1.234567, 2.45, 4.67}, 1.234568, false},
		{[]float64{1, 2, 3}, 1, true},
		{[]float32{1, 2, 3}, 1, true},
		{"this substring should be found", "substring", true},
		{"this substring should not be found", "subseastring", false},
		{nil, "foo", false},
		// Pointers
		{pagesPtr{p1, p2, p3, p2}, p2, true},
		{pagesPtr{p1, p2, p3, p2}, p4, false},
		// Structs
		{pagesVals{p3v, p2v, p3v, p2v}, p2v, true},
		{pagesVals{p3v, p2v, p3v, p2v}, p4v, false},
		// template.HTML
		{template.HTML("this substring should be found"), "substring", true},
		{template.HTML("this substring should not be found"), "subseastring", false},
		// Uncomparable, use hashstructure
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{[][]string{{"a", "b"}}, []string{"a", "b"}, true},
	} {

		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.In(test.l1, test.l2)
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

type testPage struct {
	Title string
}

func (p testPage) String() string {
	return "p-" + p.Title
}

type pagesPtr []*testPage
type pagesVals []testPage

var (
	p1 = &testPage{"A"}
	p2 = &testPage{"B"}
	p3 = &testPage{"C"}
	p4 = &testPage{"D"}

	p1v = testPage{"A"}
	p2v = testPage{"B"}
	p3v = testPage{"C"}
	p4v = testPage{"D"}
)

func TestIntersect(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		l1, l2 interface{}
		expect interface{}
	}{
		{[]string{"a", "b", "c", "c"}, []string{"a", "b", "b"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, []string{}},
		{[]string{}, []string{}, []string{}},
		{[]string{"a", "b"}, nil, []interface{}{}},
		{nil, []string{"a", "b"}, []interface{}{}},
		{nil, nil, []interface{}{}},
		{[]string{"1", "2"}, []int{1, 2}, []string{}},
		{[]int{1, 2}, []string{"1", "2"}, []int{}},
		{[]int{1, 2, 4}, []int{2, 4}, []int{2, 4}},
		{[]int{2, 4}, []int{1, 2, 4}, []int{2, 4}},
		{[]int{1, 2, 4}, []int{3, 6}, []int{}},
		{[]float64{2.2, 4.4}, []float64{1.1, 2.2, 4.4}, []float64{2.2, 4.4}},

		// []interface{} ∩ []interface{}
		{[]interface{}{"a", "b", "c"}, []interface{}{"a", "b", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{1, 2, 3}, []interface{}{1, 2, 2}, []interface{}{1, 2}},
		{[]interface{}{int8(1), int8(2), int8(3)}, []interface{}{int8(1), int8(2), int8(2)}, []interface{}{int8(1), int8(2)}},
		{[]interface{}{int16(1), int16(2), int16(3)}, []interface{}{int16(1), int16(2), int16(2)}, []interface{}{int16(1), int16(2)}},
		{[]interface{}{int32(1), int32(2), int32(3)}, []interface{}{int32(1), int32(2), int32(2)}, []interface{}{int32(1), int32(2)}},
		{[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{int64(1), int64(2), int64(2)}, []interface{}{int64(1), int64(2)}},
		{[]interface{}{float32(1), float32(2), float32(3)}, []interface{}{float32(1), float32(2), float32(2)}, []interface{}{float32(1), float32(2)}},
		{[]interface{}{float64(1), float64(2), float64(3)}, []interface{}{float64(1), float64(2), float64(2)}, []interface{}{float64(1), float64(2)}},

		// []interface{} ∩ []T
		{[]interface{}{"a", "b", "c"}, []string{"a", "b", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{1, 2, 3}, []int{1, 2, 2}, []interface{}{1, 2}},
		{[]interface{}{int8(1), int8(2), int8(3)}, []int8{1, 2, 2}, []interface{}{int8(1), int8(2)}},
		{[]interface{}{int16(1), int16(2), int16(3)}, []int16{1, 2, 2}, []interface{}{int16(1), int16(2)}},
		{[]interface{}{int32(1), int32(2), int32(3)}, []int32{1, 2, 2}, []interface{}{int32(1), int32(2)}},
		{[]interface{}{int64(1), int64(2), int64(3)}, []int64{1, 2, 2}, []interface{}{int64(1), int64(2)}},
		{[]interface{}{uint(1), uint(2), uint(3)}, []uint{1, 2, 2}, []interface{}{uint(1), uint(2)}},
		{[]interface{}{float32(1), float32(2), float32(3)}, []float32{1, 2, 2}, []interface{}{float32(1), float32(2)}},
		{[]interface{}{float64(1), float64(2), float64(3)}, []float64{1, 2, 2}, []interface{}{float64(1), float64(2)}},

		// []T ∩ []interface{}
		{[]string{"a", "b", "c"}, []interface{}{"a", "b", "b"}, []string{"a", "b"}},
		{[]int{1, 2, 3}, []interface{}{1, 2, 2}, []int{1, 2}},
		{[]int8{1, 2, 3}, []interface{}{int8(1), int8(2), int8(2)}, []int8{1, 2}},
		{[]int16{1, 2, 3}, []interface{}{int16(1), int16(2), int16(2)}, []int16{1, 2}},
		{[]int32{1, 2, 3}, []interface{}{int32(1), int32(2), int32(2)}, []int32{1, 2}},
		{[]int64{1, 2, 3}, []interface{}{int64(1), int64(2), int64(2)}, []int64{1, 2}},
		{[]float32{1, 2, 3}, []interface{}{float32(1), float32(2), float32(2)}, []float32{1, 2}},
		{[]float64{1, 2, 3}, []interface{}{float64(1), float64(2), float64(2)}, []float64{1, 2}},

		// Structs
		{pagesPtr{p1, p4, p2, p3}, pagesPtr{p4, p2, p2}, pagesPtr{p4, p2}},
		{pagesVals{p1v, p4v, p2v, p3v}, pagesVals{p1v, p3v, p3v}, pagesVals{p1v, p3v}},
		{[]interface{}{p1, p4, p2, p3}, []interface{}{p4, p2, p2}, []interface{}{p4, p2}},
		{[]interface{}{p1v, p4v, p2v, p3v}, []interface{}{p1v, p3v, p3v}, []interface{}{p1v, p3v}},
		{pagesPtr{p1, p4, p2, p3}, pagesPtr{}, pagesPtr{}},
		{pagesVals{}, pagesVals{p1v, p3v, p3v}, pagesVals{}},
		{[]interface{}{p1, p4, p2, p3}, []interface{}{}, []interface{}{}},
		{[]interface{}{}, []interface{}{p1v, p3v, p3v}, []interface{}{}},

		// errors
		{"not array or slice", []string{"a"}, false},
		{[]string{"a"}, "not array or slice", false},

		// uncomparable types - #3820
		{[]map[int]int{{1: 1}, {2: 2}}, []map[int]int{{2: 2}, {3: 3}}, false},
		{[][]int{{1, 1}, {1, 2}}, [][]int{{1, 2}, {1, 2}, {1, 3}}, false},
		{[]int{1, 1}, [][]int{{1, 2}, {1, 2}, {1, 3}}, false},
	} {

		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Intersect(test.l1, test.l2)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		if !reflect.DeepEqual(result, test.expect) {
			t.Fatalf("[%d] Got\n%v expected\n%v", i, result, test.expect)
		}
	}
}

func TestIsSet(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newTestNs()

	for i, test := range []struct {
		a      interface{}
		key    interface{}
		expect bool
		isErr  bool
	}{
		{[]interface{}{1, 2, 3, 5}, 2, true, false},
		{[]interface{}{1, 2, 3, 5}, "2", true, false},
		{[]interface{}{1, 2, 3, 5}, 2.0, true, false},

		{[]interface{}{1, 2, 3, 5}, 22, false, false},

		{map[string]interface{}{"a": 1, "b": 2}, "b", true, false},
		{map[string]interface{}{"a": 1, "b": 2}, "bc", false, false},

		{time.Now(), "Day", false, false},
		{nil, "nil", false, false},
		{[]interface{}{1, 2, 3, 5}, TstX{}, false, true},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.IsSet(test.a, test.key)
		if test.isErr {
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestLast(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		limit  interface{}
		seq    interface{}
		expect interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"b", "c"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{200, 300}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{300}},
		{"0", []int{100, 200, 300}, []int{}},
		{"0", []string{"a", "b", "c"}, []string{}},
		// errors
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
		{1, (*string)(nil), false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Last(test.limit, test.seq)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.DeepEquals, test.expect, errMsg)
	}
}

func TestQuerify(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	for i, test := range []struct {
		params []interface{}
		expect interface{}
	}{
		{[]interface{}{"a", "b"}, "a=b"},
		{[]interface{}{"a", "b", "c", "d", "f", " &"}, `a=b&c=d&f=+%26`},
		// errors
		{[]interface{}{5, "b"}, false},
		{[]interface{}{"a", "b", "c"}, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.params)

		result, err := ns.Querify(test.params...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestSeq(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	for i, test := range []struct {
		args   []interface{}
		expect interface{}
	}{
		{[]interface{}{-2, 5}, []int{-2, -1, 0, 1, 2, 3, 4, 5}},
		{[]interface{}{1, 2, 4}, []int{1, 3}},
		{[]interface{}{1}, []int{1}},
		{[]interface{}{3}, []int{1, 2, 3}},
		{[]interface{}{3.2}, []int{1, 2, 3}},
		{[]interface{}{0}, []int{}},
		{[]interface{}{-1}, []int{-1}},
		{[]interface{}{-3}, []int{-1, -2, -3}},
		{[]interface{}{3, -2}, []int{3, 2, 1, 0, -1, -2}},
		{[]interface{}{6, -2, 2}, []int{6, 4, 2}},
		// errors
		{[]interface{}{1, 0, 2}, false},
		{[]interface{}{1, -1, 2}, false},
		{[]interface{}{2, 1, 1}, false},
		{[]interface{}{2, 1, 1, 1}, false},
		{[]interface{}{2001}, false},
		{[]interface{}{}, false},
		{[]interface{}{0, -1000000}, false},
		{[]interface{}{tstNoStringer{}}, false},
		{nil, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Seq(test.args...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.DeepEquals, test.expect, errMsg)
	}
}

func TestShuffle(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	for i, test := range []struct {
		seq     interface{}
		success bool
	}{
		{[]string{"a", "b", "c", "d"}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200}, true},
		{[]string{"a", "b"}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100, 200, 300}, true},
		{[]int{100}, true},
		// errors
		{nil, false},
		{t, false},
		{(*string)(nil), false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Shuffle(test.seq)

		if !test.success {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)

		resultv := reflect.ValueOf(result)
		seqv := reflect.ValueOf(test.seq)

		c.Assert(seqv.Len(), qt.Equals, resultv.Len(), errMsg)
	}
}

func TestShuffleRandomising(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	// Note that this test can fail with false negative result if the shuffle
	// of the sequence happens to be the same as the original sequence. However
	// the propability of the event is 10^-158 which is negligible.
	seqLen := 100
	rand.Seed(time.Now().UTC().UnixNano())

	for _, test := range []struct {
		seq []int
	}{
		{rand.Perm(seqLen)},
	} {
		result, err := ns.Shuffle(test.seq)
		resultv := reflect.ValueOf(result)

		c.Assert(err, qt.IsNil)

		allSame := true
		for i, v := range test.seq {
			allSame = allSame && (resultv.Index(i).Interface() == v)
		}

		c.Assert(allSame, qt.Equals, false)
	}
}

// Also see tests in commons/collection.
func TestSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})

	for i, test := range []struct {
		args     []interface{}
		expected interface{}
	}{
		{[]interface{}{"a", "b"}, []string{"a", "b"}},
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{nil}, []interface{}{nil}},
		{[]interface{}{5, "b"}, []interface{}{5, "b"}},
		{[]interface{}{tstNoStringer{}}, []tstNoStringer{{}}},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.args)

		result := ns.Slice(test.args...)

		c.Assert(result, qt.DeepEquals, test.expected, errMsg)
	}

}

func TestUnion(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		l1     interface{}
		l2     interface{}
		expect interface{}
		isErr  bool
	}{
		{nil, nil, []interface{}{}, false},
		{nil, []string{"a", "b"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, nil, []string{"a", "b"}, false},

		// []A ∪ []B
		{[]string{"1", "2"}, []int{3}, []string{}, false},
		{[]int{1, 2}, []string{"1", "2"}, []int{}, false},

		// []T ∪ []T
		{[]string{"a", "b", "c", "c"}, []string{"a", "b", "b"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"a", "b", "c", "d", "e"}, false},
		{[]string{}, []string{}, []string{}, false},
		{[]int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2, 3, 4, 5}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2, 4}, []int{2, 4}, []int{1, 2, 4}, false},
		{[]int{2, 4}, []int{1, 2, 4}, []int{2, 4, 1}, false},
		{[]int{1, 2, 4}, []int{3, 6}, []int{1, 2, 4, 3, 6}, false},
		{[]float64{2.2, 4.4}, []float64{1.1, 2.2, 4.4}, []float64{2.2, 4.4, 1.1}, false},
		{[]interface{}{"a", "b", "c", "c"}, []interface{}{"a", "b", "b"}, []interface{}{"a", "b", "c"}, false},

		// []T ∪ []interface{}
		{[]string{"1", "2"}, []interface{}{"9"}, []string{"1", "2", "9"}, false},
		{[]int{2, 4}, []interface{}{1, 2, 4}, []int{2, 4, 1}, false},
		{[]int8{2, 4}, []interface{}{int8(1), int8(2), int8(4)}, []int8{2, 4, 1}, false},
		{[]int8{2, 4}, []interface{}{1, 2, 4}, []int8{2, 4, 1}, false},
		{[]int16{2, 4}, []interface{}{1, 2, 4}, []int16{2, 4, 1}, false},
		{[]int32{2, 4}, []interface{}{1, 2, 4}, []int32{2, 4, 1}, false},
		{[]int64{2, 4}, []interface{}{1, 2, 4}, []int64{2, 4, 1}, false},

		{[]float64{2.2, 4.4}, []interface{}{1.1, 2.2, 4.4}, []float64{2.2, 4.4, 1.1}, false},
		{[]float32{2.2, 4.4}, []interface{}{1.1, 2.2, 4.4}, []float32{2.2, 4.4, 1.1}, false},

		// []interface{} ∪ []T
		{[]interface{}{"a", "b", "c", "c"}, []string{"a", "b", "d"}, []interface{}{"a", "b", "c", "d"}, false},
		{[]interface{}{}, []string{}, []interface{}{}, false},
		{[]interface{}{1, 2}, []int{2, 3}, []interface{}{1, 2, 3}, false},
		{[]interface{}{1, 2}, []int8{2, 3}, []interface{}{1, 2, 3}, false}, // 28
		{[]interface{}{uint(1), uint(2)}, []uint{2, 3}, []interface{}{uint(1), uint(2), uint(3)}, false},
		{[]interface{}{1.1, 2.2}, []float64{2.2, 3.3}, []interface{}{1.1, 2.2, 3.3}, false},

		// Structs
		{pagesPtr{p1, p4}, pagesPtr{p4, p2, p2}, pagesPtr{p1, p4, p2}, false},
		{pagesVals{p1v}, pagesVals{p3v, p3v}, pagesVals{p1v, p3v}, false},
		{[]interface{}{p1, p4}, []interface{}{p4, p2, p2}, []interface{}{p1, p4, p2}, false},
		{[]interface{}{p1v}, []interface{}{p3v, p3v}, []interface{}{p1v, p3v}, false},
		// #3686
		{[]interface{}{p1v}, []interface{}{}, []interface{}{p1v}, false},
		{[]interface{}{}, []interface{}{p1v}, []interface{}{p1v}, false},
		{pagesPtr{p1}, pagesPtr{}, pagesPtr{p1}, false},
		{pagesVals{p1v}, pagesVals{}, pagesVals{p1v}, false},
		{pagesPtr{}, pagesPtr{p1}, pagesPtr{p1}, false},
		{pagesVals{}, pagesVals{p1v}, pagesVals{p1v}, false},

		// errors
		{"not array or slice", []string{"a"}, false, true},
		{[]string{"a"}, "not array or slice", false, true},

		// uncomparable types - #3820
		{[]map[string]int{{"K1": 1}}, []map[string]int{{"K2": 2}, {"K2": 2}}, false, true},
		{[][]int{{1, 1}, {1, 2}}, [][]int{{2, 1}, {2, 2}}, false, true},
	} {

		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Union(test.l1, test.l2)
		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		if !reflect.DeepEqual(result, test.expect) {
			t.Fatalf("[%d] Got\n%v expected\n%v", i, result, test.expect)
		}
	}
}

func TestUniq(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(&deps.Deps{})
	for i, test := range []struct {
		l      interface{}
		expect interface{}
		isErr  bool
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b", "c", "c"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b", "b", "c"}, []string{"a", "b", "c"}, false},
		{[]string{"a", "b", "c", "b"}, []string{"a", "b", "c"}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2, 3, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2, 3, 2}, []int{1, 2, 3}, false},
		{[4]int{1, 2, 3, 2}, []int{1, 2, 3}, false},
		{nil, make([]interface{}, 0), false},
		// Pointers
		{pagesPtr{p1, p2, p3, p2}, pagesPtr{p1, p2, p3}, false},
		{pagesPtr{}, pagesPtr{}, false},
		// Structs
		{pagesVals{p3v, p2v, p3v, p2v}, pagesVals{p3v, p2v}, false},

		// not Comparable(), use hashstruscture
		{[]map[string]int{
			{"K1": 1}, {"K2": 2}, {"K1": 1}, {"K2": 1},
		}, []map[string]int{
			{"K1": 1}, {"K2": 2}, {"K2": 1},
		}, false},

		// should fail
		{1, 1, true},
		{"foo", "fo", true},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Uniq(test.l)
		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.DeepEquals, test.expect, errMsg)
	}
}

func (x *TstX) TstRp() string {
	return "r" + x.A
}

func (x TstX) TstRv() string {
	return "r" + x.B
}

func (x TstX) TstRv2() string {
	return "r" + x.B
}

func (x TstX) unexportedMethod() string {
	return x.unexported
}

func (x TstX) MethodWithArg(s string) string {
	return s
}

func (x TstX) MethodReturnNothing() {}

func (x TstX) MethodReturnErrorOnly() error {
	return errors.New("some error occurred")
}

func (x TstX) MethodReturnTwoValues() (string, string) {
	return "foo", "bar"
}

func (x TstX) MethodReturnValueWithError() (string, error) {
	return "", errors.New("some error occurred")
}

func (x TstX) String() string {
	return fmt.Sprintf("A: %s, B: %s", x.A, x.B)
}

type TstX struct {
	A, B       string
	unexported string
}

type TstParams struct {
	params maps.Params
}

func (x TstParams) Params() maps.Params {
	return x.params

}

type TstXIHolder struct {
	XI TstXI
}

// Partially implemented by the TstX struct.
type TstXI interface {
	TstRv2() string
}

func ToTstXIs(slice interface{}) []TstXI {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	tis := make([]TstXI, s.Len())

	for i := 0; i < s.Len(); i++ {
		tsti, ok := s.Index(i).Interface().(TstXI)
		if !ok {
			return nil
		}
		tis[i] = tsti
	}

	return tis
}

func newDeps(cfg config.Provider) *deps.Deps {
	l := langs.NewLanguage("en", cfg)
	l.Set("i18nDir", "i18n")
	cs, err := helpers.NewContentSpec(l, loggers.NewErrorLogger(), afero.NewMemMapFs())
	if err != nil {
		panic(err)
	}
	return &deps.Deps{
		Cfg:         cfg,
		Fs:          hugofs.NewMem(l),
		ContentSpec: cs,
		Log:         loggers.NewErrorLogger(),
	}
}

func newTestNs() *Namespace {
	v := viper.New()
	v.Set("contentDir", "content")
	return New(newDeps(v))
}
