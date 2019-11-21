// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"fmt"
	"reflect"
	"testing"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/deps"
)

type stringsSlice []string

func TestSort(t *testing.T) {
	t.Parallel()

	ns := New(&deps.Deps{})

	type ts struct {
		MyInt    int
		MyFloat  float64
		MyString string
	}
	type mid struct {
		Tst TstX
	}

	for i, test := range []struct {
		seq         interface{}
		sortByField interface{}
		sortAsc     string
		expect      interface{}
	}{
		{[]string{"class1", "class2", "class3"}, nil, "asc", []string{"class1", "class2", "class3"}},
		{[]string{"class3", "class1", "class2"}, nil, "asc", []string{"class1", "class2", "class3"}},
		{[]string{"CLASS3", "class1", "class2"}, nil, "asc", []string{"class1", "class2", "CLASS3"}},
		// Issue 6023
		{stringsSlice{"class3", "class1", "class2"}, nil, "asc", stringsSlice{"class1", "class2", "class3"}},

		{[]int{1, 2, 3, 4, 5}, nil, "asc", []int{1, 2, 3, 4, 5}},
		{[]int{5, 4, 3, 1, 2}, nil, "asc", []int{1, 2, 3, 4, 5}},
		// test sort key parameter is focibly set empty
		{[]string{"class3", "class1", "class2"}, map[int]string{1: "a"}, "asc", []string{"class1", "class2", "class3"}},
		// test map sorting by keys
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, nil, "asc", []int{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, nil, "asc", []int{30, 20, 10, 40, 50}},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, nil, "asc", []string{"10", "20", "30", "40", "50"}},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, nil, "asc", []string{"50", "40", "10", "30", "20"}},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, nil, "asc", []string{"10", "20", "30", "40", "50"}},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, nil, "asc", []string{"30", "20", "10", "40", "50"}},
		// test map sorting by value
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "value", "asc", []int{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "value", "asc", []int{10, 20, 30, 40, 50}},
		// test map sorting by field value
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyInt",
			"asc",
			[]ts{{10, 10.5, "ten"}, {20, 20.5, "twenty"}, {30, 30.5, "thirty"}, {40, 40.5, "forty"}, {50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyFloat",
			"asc",
			[]ts{{10, 10.5, "ten"}, {20, 20.5, "twenty"}, {30, 30.5, "thirty"}, {40, 40.5, "forty"}, {50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyString",
			"asc",
			[]ts{{50, 50.5, "fifty"}, {40, 40.5, "forty"}, {10, 10.5, "ten"}, {30, 30.5, "thirty"}, {20, 20.5, "twenty"}},
		},
		// test sort desc
		{[]string{"class1", "class2", "class3"}, "value", "desc", []string{"class3", "class2", "class1"}},
		{[]string{"class3", "class1", "class2"}, "value", "desc", []string{"class3", "class2", "class1"}},
		// test sort by struct's method
		{
			[]TstX{{A: "i", B: "j"}, {A: "e", B: "f"}, {A: "c", B: "d"}, {A: "g", B: "h"}, {A: "a", B: "b"}},
			"TstRv",
			"asc",
			[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		{
			[]*TstX{{A: "i", B: "j"}, {A: "e", B: "f"}, {A: "c", B: "d"}, {A: "g", B: "h"}, {A: "a", B: "b"}},
			"TstRp",
			"asc",
			[]*TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		// Lower case Params, slice
		{
			[]TstParams{{params: maps.Params{"color": "indigo"}}, {params: maps.Params{"color": "blue"}}, {params: maps.Params{"color": "green"}}},
			".Params.COLOR",
			"asc",
			[]TstParams{{params: maps.Params{"color": "blue"}}, {params: maps.Params{"color": "green"}}, {params: maps.Params{"color": "indigo"}}},
		},
		// Lower case Params, map
		{
			map[string]TstParams{"1": {params: maps.Params{"color": "indigo"}}, "2": {params: maps.Params{"color": "blue"}}, "3": {params: maps.Params{"color": "green"}}},
			".Params.CoLoR",
			"asc",
			[]TstParams{{params: maps.Params{"color": "blue"}}, {params: maps.Params{"color": "green"}}, {params: maps.Params{"color": "indigo"}}},
		},
		// test map sorting by struct's method
		{
			map[string]TstX{"1": {A: "i", B: "j"}, "2": {A: "e", B: "f"}, "3": {A: "c", B: "d"}, "4": {A: "g", B: "h"}, "5": {A: "a", B: "b"}},
			"TstRv",
			"asc",
			[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		{
			map[string]*TstX{"1": {A: "i", B: "j"}, "2": {A: "e", B: "f"}, "3": {A: "c", B: "d"}, "4": {A: "g", B: "h"}, "5": {A: "a", B: "b"}},
			"TstRp",
			"asc",
			[]*TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}, {A: "g", B: "h"}, {A: "i", B: "j"}},
		},
		// test sort by dot chaining key argument
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			".foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.TstRv",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]*TstX{{"foo": &TstX{A: "e", B: "f"}}, {"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}},
			"foo.TstRp",
			"asc",
			[]map[string]*TstX{{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}}},
		},
		{
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "e", B: "f"}}}, {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.A",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		{
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "e", B: "f"}}}, {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.TstRv",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		// test map sorting by dot chaining key argument
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			".foo.A",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.TstRv",
			"asc",
			[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]*TstX{"1": {"foo": &TstX{A: "e", B: "f"}}, "2": {"foo": &TstX{A: "a", B: "b"}}, "3": {"foo": &TstX{A: "c", B: "d"}}},
			"foo.TstRp",
			"asc",
			[]map[string]*TstX{{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}}},
		},
		{
			map[string]map[string]mid{"1": {"foo": mid{Tst: TstX{A: "e", B: "f"}}}, "2": {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, "3": {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.A",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		{
			map[string]map[string]mid{"1": {"foo": mid{Tst: TstX{A: "e", B: "f"}}}, "2": {"foo": mid{Tst: TstX{A: "a", B: "b"}}}, "3": {"foo": mid{Tst: TstX{A: "c", B: "d"}}}},
			"foo.Tst.TstRv",
			"asc",
			[]map[string]mid{{"foo": mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": mid{Tst: TstX{A: "e", B: "f"}}}},
		},
		// interface slice with missing elements
		{
			[]interface{}{
				map[interface{}]interface{}{"Title": "Foo", "Weight": 10},
				map[interface{}]interface{}{"Title": "Bar"},
				map[interface{}]interface{}{"Title": "Zap", "Weight": 5},
			},
			"Weight",
			"asc",
			[]interface{}{
				map[interface{}]interface{}{"Title": "Bar"},
				map[interface{}]interface{}{"Title": "Zap", "Weight": 5},
				map[interface{}]interface{}{"Title": "Foo", "Weight": 10},
			},
		},
		// test error cases
		{(*[]TstX)(nil), nil, "asc", false},
		{TstX{A: "a", B: "b"}, nil, "asc", false},
		{
			[]map[string]TstX{{"foo": TstX{A: "e", B: "f"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}},
			"foo.NotAvailable",
			"asc",
			false,
		},
		{
			map[string]map[string]TstX{"1": {"foo": TstX{A: "e", B: "f"}}, "2": {"foo": TstX{A: "a", B: "b"}}, "3": {"foo": TstX{A: "c", B: "d"}}},
			"foo.NotAvailable",
			"asc",
			false,
		},
		{nil, nil, "asc", false},
	} {

		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			var result interface{}
			var err error
			if test.sortByField == nil {
				result, err = ns.Sort(test.seq)
			} else {
				result, err = ns.Sort(test.seq, test.sortByField, test.sortAsc)
			}

			if b, ok := test.expect.(bool); ok && !b {
				if err == nil {
					t.Fatal("Sort didn't return an expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("failed: %s", err)
				}
				if !reflect.DeepEqual(result, test.expect) {
					t.Fatalf("Sort called on sequence: %#v | sortByField: `%v` | got\n%#v but expected\n%#v", test.seq, test.sortByField, result, test.expect)
				}
			}
		})

	}
}
