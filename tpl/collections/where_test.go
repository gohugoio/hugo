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
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/maps"
)

func TestWhere(t *testing.T) {
	t.Parallel()

	ns := newNs()

	type Mid struct {
		Tst TstX
	}

	d1 := time.Now()
	d2 := d1.Add(1 * time.Hour)
	d3 := d2.Add(1 * time.Hour)
	d4 := d3.Add(1 * time.Hour)
	d5 := d4.Add(1 * time.Hour)
	d6 := d5.Add(1 * time.Hour)

	type testt struct {
		seq    any
		key    any
		op     string
		match  any
		expect any
	}

	createTestVariants := func(test testt) []testt {
		testVariants := []testt{test}
		if islice := ToTstXIs(test.seq); islice != nil {
			variant := test
			variant.seq = islice
			expect := ToTstXIs(test.expect)
			if expect != nil {
				variant.expect = expect
			}
			testVariants = append(testVariants, variant)
		}

		return testVariants
	}

	for i, test := range []testt{
		{
			seq: []map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []map[int]string{
				{1: "a", 2: "m"},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4,
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.0,
			expect: []map[string]float64{{"a": 3, "b": 4}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.0, op: "!=",
			expect: []map[string]float64{{"a": 1, "b": 2}, {"a": 5, "x": 4}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.0, op: "<",
			expect: []map[string]float64{{"a": 1, "b": 2}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4, op: "<",
			expect: []map[string]float64{{"a": 1, "b": 2}},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.0, op: "<",
			expect: []map[string]int{{"a": 1, "b": 2}},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.2, op: "<",
			expect: []map[string]int{{"a": 1, "b": 2}, {"a": 3, "b": 4}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4.0, op: "<=",
			expect: []map[string]float64{{"a": 1, "b": 2}, {"a": 3, "b": 4}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 3}, {"a": 5, "x": 4},
			},
			key: "b", match: 2.0, op: ">",
			expect: []map[string]float64{{"a": 3, "b": 3}},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 3}, {"a": 5, "x": 4},
			},
			key: "b", match: 2.0, op: ">=",
			expect: []map[string]float64{{"a": 1, "b": 2}, {"a": 3, "b": 3}},
		},
		// Issue #8353
		// String type mismatch.
		{
			seq: []map[string]any{
				{"a": "1", "b": "2"}, {"a": "3", "b": template.HTML("4")}, {"a": "5", "x": "4"},
			},
			key: "b", match: "4",
			expect: []map[string]any{
				{"a": "3", "b": template.HTML("4")},
			},
		},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []TstX{
				{A: "e", B: "f"},
			},
		},
		{
			seq: []*map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []*map[int]string{
				{1: "a", 2: "m"},
			},
		},
		// Case insensitive maps.Params
		// Slice of structs
		{
			seq: []TstParams{{params: maps.Params{"i": 0, "color": "indigo"}}, {params: maps.Params{"i": 1, "color": "blue"}}, {params: maps.Params{"i": 2, "color": "green"}}, {params: maps.Params{"i": 3, "color": "blue"}}},
			key: ".Params.COLOR", match: "blue",
			expect: []TstParams{{params: maps.Params{"i": 1, "color": "blue"}}, {params: maps.Params{"i": 3, "color": "blue"}}},
		},
		{
			seq: []TstParams{{params: maps.Params{"nested": map[string]any{"color": "indigo"}}}, {params: maps.Params{"nested": map[string]any{"color": "blue"}}}},
			key: ".Params.NEsTED.COLOR", match: "blue",
			expect: []TstParams{{params: maps.Params{"nested": map[string]any{"color": "blue"}}}},
		},
		{
			seq: []TstParams{{params: maps.Params{"i": 0, "color": "indigo"}}, {params: maps.Params{"i": 1, "color": "blue"}}, {params: maps.Params{"i": 2, "color": "green"}}, {params: maps.Params{"i": 3, "color": "blue"}}},
			key: ".Params", match: "blue",
			expect: []TstParams{},
		},
		// Slice of maps
		{
			seq: []maps.Params{
				{"a": "a1", "b": "b1"}, {"a": "a2", "b": "b2"},
			},
			key: "B", match: "b2",
			expect: []maps.Params{
				{"a": "a2", "b": "b2"},
			},
		},
		{
			seq: []maps.Params{
				{
					"a": map[string]any{
						"b": "b1",
					},
				},
				{
					"a": map[string]any{
						"b": "b2",
					},
				},
			},
			key: "A.B", match: "b2",
			expect: []maps.Params{
				{
					"a": map[string]any{
						"b": "b2",
					},
				},
			},
		},
		{
			seq: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []*TstX{
				{A: "e", B: "f"},
			},
		},
		{
			seq: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRp", match: "rc",
			expect: []*TstX{
				{A: "c", B: "d"},
			},
		},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRv", match: "rc",
			expect: []TstX{
				{A: "e", B: "c"},
			},
		},
		{
			seq: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []map[string]TstX{
				{"baz": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: ".foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRv", match: "rd",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []map[string]*TstX{
				{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRp", match: "rc",
			expect: []map[string]*TstX{
				{"foo": &TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []TstXIHolder{
				{&TstX{A: "a", B: "b"}}, {&TstX{A: "c", B: "d"}}, {&TstX{A: "e", B: "f"}},
			},
			key: "XI.TstRp", match: "rc",
			expect: []TstXIHolder{
				{&TstX{A: "c", B: "d"}},
			},
		},
		{
			seq: []TstXIHolder{
				{&TstX{A: "a", B: "b"}}, {&TstX{A: "c", B: "d"}}, {&TstX{A: "e", B: "f"}},
			},
			key: "XI.A", match: "e",
			expect: []TstXIHolder{
				{&TstX{A: "e", B: "f"}},
			},
		},
		{
			seq: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.B", match: "d",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			seq: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRv", match: "rd",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			seq: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": &Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": &Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRp", match: "rc",
			expect: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: 3,
			expect: []map[string]int{
				{"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: 3.0,
			expect: []map[string]float64{
				{"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
		},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "!=", match: "f",
			expect: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: []int{3, 4, 5},
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: []float64{3, 4, 5},
			expect: []map[string]float64{
				{"a": 3, "b": 4},
			},
		},
		{
			seq: []map[string][]string{
				{"a": []string{"A", "B", "C"}, "b": []string{"D", "E", "F"}}, {"a": []string{"G", "H", "I"}, "b": []string{"J", "K", "L"}}, {"a": []string{"M", "N", "O"}, "b": []string{"P", "Q", "R"}},
			},
			key: "b", op: "intersect", match: []string{"D", "P", "Q"},
			expect: []map[string][]string{
				{"a": []string{"A", "B", "C"}, "b": []string{"D", "E", "F"}}, {"a": []string{"M", "N", "O"}, "b": []string{"P", "Q", "R"}},
			},
		},
		{
			seq: []map[string][]int{
				{"a": []int{1, 2, 3}, "b": []int{4, 5, 6}}, {"a": []int{7, 8, 9}, "b": []int{10, 11, 12}}, {"a": []int{13, 14, 15}, "b": []int{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int{4, 10, 12},
			expect: []map[string][]int{
				{"a": []int{1, 2, 3}, "b": []int{4, 5, 6}}, {"a": []int{7, 8, 9}, "b": []int{10, 11, 12}},
			},
		},
		{
			seq: []map[string][]int8{
				{"a": []int8{1, 2, 3}, "b": []int8{4, 5, 6}}, {"a": []int8{7, 8, 9}, "b": []int8{10, 11, 12}}, {"a": []int8{13, 14, 15}, "b": []int8{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int8{4, 10, 12},
			expect: []map[string][]int8{
				{"a": []int8{1, 2, 3}, "b": []int8{4, 5, 6}}, {"a": []int8{7, 8, 9}, "b": []int8{10, 11, 12}},
			},
		},
		{
			seq: []map[string][]int16{
				{"a": []int16{1, 2, 3}, "b": []int16{4, 5, 6}}, {"a": []int16{7, 8, 9}, "b": []int16{10, 11, 12}}, {"a": []int16{13, 14, 15}, "b": []int16{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int16{4, 10, 12},
			expect: []map[string][]int16{
				{"a": []int16{1, 2, 3}, "b": []int16{4, 5, 6}}, {"a": []int16{7, 8, 9}, "b": []int16{10, 11, 12}},
			},
		},
		{
			seq: []map[string][]int32{
				{"a": []int32{1, 2, 3}, "b": []int32{4, 5, 6}}, {"a": []int32{7, 8, 9}, "b": []int32{10, 11, 12}}, {"a": []int32{13, 14, 15}, "b": []int32{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int32{4, 10, 12},
			expect: []map[string][]int32{
				{"a": []int32{1, 2, 3}, "b": []int32{4, 5, 6}}, {"a": []int32{7, 8, 9}, "b": []int32{10, 11, 12}},
			},
		},
		{
			seq: []map[string][]int64{
				{"a": []int64{1, 2, 3}, "b": []int64{4, 5, 6}}, {"a": []int64{7, 8, 9}, "b": []int64{10, 11, 12}}, {"a": []int64{13, 14, 15}, "b": []int64{16, 17, 18}},
			},
			key: "b", op: "intersect", match: []int64{4, 10, 12},
			expect: []map[string][]int64{
				{"a": []int64{1, 2, 3}, "b": []int64{4, 5, 6}}, {"a": []int64{7, 8, 9}, "b": []int64{10, 11, 12}},
			},
		},
		{
			seq: []map[string][]float32{
				{"a": []float32{1.0, 2.0, 3.0}, "b": []float32{4.0, 5.0, 6.0}}, {"a": []float32{7.0, 8.0, 9.0}, "b": []float32{10.0, 11.0, 12.0}}, {"a": []float32{13.0, 14.0, 15.0}, "b": []float32{16.0, 17.0, 18.0}},
			},
			key: "b", op: "intersect", match: []float32{4, 10, 12},
			expect: []map[string][]float32{
				{"a": []float32{1.0, 2.0, 3.0}, "b": []float32{4.0, 5.0, 6.0}}, {"a": []float32{7.0, 8.0, 9.0}, "b": []float32{10.0, 11.0, 12.0}},
			},
		},
		{
			seq: []map[string][]float64{
				{"a": []float64{1.0, 2.0, 3.0}, "b": []float64{4.0, 5.0, 6.0}}, {"a": []float64{7.0, 8.0, 9.0}, "b": []float64{10.0, 11.0, 12.0}}, {"a": []float64{13.0, 14.0, 15.0}, "b": []float64{16.0, 17.0, 18.0}},
			},
			key: "b", op: "intersect", match: []float64{4, 10, 12},
			expect: []map[string][]float64{
				{"a": []float64{1.0, 2.0, 3.0}, "b": []float64{4.0, 5.0, 6.0}}, {"a": []float64{7.0, 8.0, 9.0}, "b": []float64{10.0, 11.0, 12.0}},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: ns.Slice(3, 4, 5),
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: ns.Slice(3.0, 4.0, 5.0),
			expect: []map[string]float64{
				{"a": 3, "b": 4},
			},
		},
		{
			seq: []map[string]time.Time{
				{"a": d1, "b": d2}, {"a": d3, "b": d4}, {"a": d5, "b": d6},
			},
			key: "b", op: "in", match: ns.Slice(d3, d4, d5),
			expect: []map[string]time.Time{
				{"a": d3, "b": d4},
			},
		},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "not in", match: []string{"c", "d", "e"},
			expect: []TstX{
				{A: "a", B: "b"}, {A: "e", B: "f"},
			},
		},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "not in", match: ns.Slice("c", t, "d", "e"),
			expect: []TstX{
				{A: "a", B: "b"}, {A: "e", B: "f"},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "", match: nil,
			expect: []map[string]int{
				{"a": 3},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "!=", match: nil,
			expect: []map[string]int{
				{"a": 1, "b": 2}, {"a": 5, "b": 6},
			},
		},
		{
			seq: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: nil,
			expect: []map[string]int{},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "", match: nil,
			expect: []map[string]float64{
				{"a": 3},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "!=", match: nil,
			expect: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 5, "b": 6},
			},
		},
		{
			seq: []map[string]float64{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: nil,
			expect: []map[string]float64{},
		},
		{
			seq: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: "", match: true,
			expect: []map[string]bool{
				{"c": true, "b": true},
			},
		},
		{
			seq: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: "!=", match: true,
			expect: []map[string]bool{
				{"a": true, "b": false}, {"d": true, "b": false},
			},
		},
		{
			seq: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b", op: ">", match: false,
			expect: []map[string]bool{},
		},
		{
			seq: []map[string]bool{
				{"a": true, "b": false}, {"c": true, "b": true}, {"d": true, "b": false},
			},
			key: "b.z", match: false,
			expect: []map[string]bool{},
		},
		{seq: (*[]TstX)(nil), key: "A", match: "a", expect: false},
		{seq: TstX{A: "a", B: "b"}, key: "A", match: "a", expect: false},
		{seq: []map[string]*TstX{{"foo": nil}}, key: "foo.B", match: "d", expect: []map[string]*TstX{}},
		{seq: []map[string]*TstX{{"foo": nil}}, key: "foo.B.Z", match: "d", expect: []map[string]*TstX{}},
		{
			seq: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "op", match: "f",
			expect: false,
		},
		{
			seq: map[string]any{
				"foo": []any{map[any]any{"a": 1, "b": 2}},
				"bar": []any{map[any]any{"a": 3, "b": 4}},
				"zap": []any{map[any]any{"a": 5, "b": 6}},
			},
			key: "b", op: "in", match: ns.Slice(3, 4, 5),
			expect: map[string]any{
				"bar": []any{map[any]any{"a": 3, "b": 4}},
			},
		},
		{
			seq: map[string]any{
				"foo": []any{map[any]any{"a": 1, "b": 2}},
				"bar": []any{map[any]any{"a": 3, "b": 4}},
				"zap": []any{map[any]any{"a": 5, "b": 6}},
			},
			key: "b", op: ">", match: 3,
			expect: map[string]any{
				"bar": []any{map[any]any{"a": 3, "b": 4}},
				"zap": []any{map[any]any{"a": 5, "b": 6}},
			},
		},
		{
			seq: map[string]any{
				"foo": []any{maps.Params{"a": 1, "b": 2}},
				"bar": []any{maps.Params{"a": 3, "b": 4}},
				"zap": []any{maps.Params{"a": 5, "b": 6}},
			},
			key: "B", op: ">", match: 3,
			expect: map[string]any{
				"bar": []any{maps.Params{"a": 3, "b": 4}},
				"zap": []any{maps.Params{"a": 5, "b": 6}},
			},
		},
	} {

		testVariants := createTestVariants(test)
		for j, test := range testVariants {
			name := fmt.Sprintf("%d/%d %T %s %s", i, j, test.seq, test.op, test.key)
			name = strings.ReplaceAll(name, "[]", "slice-of-")
			t.Run(name, func(t *testing.T) {
				var results any
				var err error

				if len(test.op) > 0 {
					results, err = ns.Where(context.Background(), test.seq, test.key, test.op, test.match)
				} else {
					results, err = ns.Where(context.Background(), test.seq, test.key, test.match)
				}
				if b, ok := test.expect.(bool); ok && !b {
					if err == nil {
						t.Fatalf("[%d] Where didn't return an expected error", i)
					}
				} else {
					if err != nil {
						t.Fatalf("[%d] failed: %s", i, err)
					}
					if !reflect.DeepEqual(results, test.expect) {
						t.Fatalf("Where clause matching %v with %v in seq %v (%T),\ngot\n%v (%T) but expected\n%v (%T)", test.key, test.match, test.seq, test.seq, results, results, test.expect, test.expect)
					}
				}
			})
		}
	}

	var err error
	_, err = ns.Where(context.Background(), map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1)
	if err == nil {
		t.Errorf("Where called with none string op value didn't return an expected error")
	}

	_, err = ns.Where(context.Background(), map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1, 2)
	if err == nil {
		t.Errorf("Where called with more than two variable arguments didn't return an expected error")
	}

	_, err = ns.Where(context.Background(), map[string]int{"a": 1, "b": 2}, "a")
	if err == nil {
		t.Errorf("Where called with no variable arguments didn't return an expected error")
	}
}

func TestCheckCondition(t *testing.T) {
	t.Parallel()

	ns := newNs()

	type expect struct {
		result  bool
		isError bool
	}

	for i, test := range []struct {
		value reflect.Value
		match reflect.Value
		op    string
		expect
	}{
		{reflect.ValueOf(123), reflect.ValueOf(123), "", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("foo"), "", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"",
			expect{true, false},
		},
		{reflect.ValueOf(true), reflect.ValueOf(true), "", expect{true, false}},
		{reflect.ValueOf(nil), reflect.ValueOf(nil), "", expect{true, false}},
		{reflect.ValueOf(123), reflect.ValueOf(456), "!=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), "!=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			"!=",
			expect{true, false},
		},
		{reflect.ValueOf(true), reflect.ValueOf(false), "!=", expect{true, false}},
		{reflect.ValueOf(123), reflect.ValueOf(nil), "!=", expect{true, false}},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">=",
			expect{true, false},
		},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<=", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<=",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{123, 45, 678}), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"foo", "bar", "baz"}), "in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.June, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"in",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{45, 678}), "not in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"bar", "baz"}), "not in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.February, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.March, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"not in",
			expect{true, false},
		},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar-foo-baz"), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar--baz"), "not in", expect{true, false}},
		{reflect.Value{}, reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.Value{}, "", expect{false, false}},
		{reflect.ValueOf((*TstX)(nil)), reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf((*TstX)(nil)), "", expect{false, false}},
		{reflect.ValueOf(true), reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf(true), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf(map[int]string{}), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]int{1, 2}), "", expect{false, false}},
		{reflect.ValueOf((*TstX)(nil)), reflect.ValueOf((*TstX)(nil)), ">", expect{false, false}},
		{reflect.ValueOf(true), reflect.ValueOf(false), ">", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf([]int{}), "in", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf(123), "op", expect{false, true}},

		// Issue #3718
		{reflect.ValueOf([]any{"a"}), reflect.ValueOf([]string{"a", "b"}), "intersect", expect{true, false}},
		{reflect.ValueOf([]string{"a"}), reflect.ValueOf([]any{"a", "b"}), "intersect", expect{true, false}},
		{reflect.ValueOf([]any{1, 2}), reflect.ValueOf([]int{1}), "intersect", expect{true, false}},
		{reflect.ValueOf([]int{1}), reflect.ValueOf([]any{1, 2}), "intersect", expect{true, false}},
	} {
		result, err := ns.checkCondition(test.value, test.match, test.op)
		if test.expect.isError {
			if err == nil {
				t.Errorf("[%d] checkCondition didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result != test.expect.result {
				t.Errorf("[%d] check condition %v %s %v, got %v but expected %v", i, test.value, test.op, test.match, result, test.expect.result)
			}
		}
	}
}

func TestEvaluateSubElem(t *testing.T) {
	t.Parallel()
	tstx := TstX{A: "foo", B: "bar"}
	var inner struct {
		S fmt.Stringer
	}
	inner.S = tstx
	interfaceValue := reflect.ValueOf(&inner).Elem().Field(0)

	for i, test := range []struct {
		value  reflect.Value
		key    string
		expect any
	}{
		{reflect.ValueOf(tstx), "A", "foo"},
		{reflect.ValueOf(&tstx), "TstRp", "rfoo"},
		{reflect.ValueOf(tstx), "TstRv", "rbar"},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1, "foo"},
		{reflect.ValueOf(map[string]string{"key1": "foo", "key2": "bar"}), "key1", "foo"},
		{interfaceValue, "String", "A: foo, B: bar"},
		{reflect.Value{}, "foo", false},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1.2, false},
		{reflect.ValueOf(tstx), "unexported", false},
		{reflect.ValueOf(tstx), "unexportedMethod", false},
		{reflect.ValueOf(tstx), "MethodWithArg", false},
		{reflect.ValueOf(tstx), "MethodReturnNothing", false},
		{reflect.ValueOf(tstx), "MethodReturnErrorOnly", false},
		{reflect.ValueOf(tstx), "MethodReturnTwoValues", false},
		{reflect.ValueOf(tstx), "MethodReturnValueWithError", false},
		{reflect.ValueOf((*TstX)(nil)), "A", false},
		{reflect.ValueOf(tstx), "C", false},
		{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), "1", false},
		{reflect.ValueOf([]string{"foo", "bar"}), "1", false},
	} {
		result, err := evaluateSubElem(reflect.ValueOf(context.Background()), test.value, test.key)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] evaluateSubElem didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result.Kind() != reflect.String || result.String() != test.expect {
				t.Errorf("[%d] evaluateSubElem with %v got %v but expected %v", i, test.key, result, test.expect)
			}
		}
	}
}

func BenchmarkWhereOps(b *testing.B) {
	ns := newNs()
	var seq []map[string]string
	ctx := context.Background()
	for i := 0; i < 500; i++ {
		seq = append(seq, map[string]string{"foo": "bar"})
	}
	for i := 0; i < 500; i++ {
		seq = append(seq, map[string]string{"foo": "baz"})
	}
	// Shuffle the sequence.
	for i := range seq {
		j := rand.Intn(i + 1)
		seq[i], seq[j] = seq[j], seq[i]
	}
	// results, err = ns.Where(context.Background(), test.seq, test.key, test.op, test.match)
	runOps := func(b *testing.B, op, match string) {
		_, err := ns.Where(ctx, seq, "foo", op, match)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Run("eq", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runOps(b, "eq", "bar")
		}
	})

	b.Run("ne", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runOps(b, "ne", "baz")
		}
	})

	b.Run("like", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runOps(b, "like", "^bar")
		}
	})
}
