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

package reflect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ns = New()

type tstNoStringer struct{}

func TestKindIs(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		k      string
		v      interface{}
		expect interface{}
	}{
		{"invalid", nil, true},
		{"string", "foo", true},
		{"int", 1, true},
		{"int", "1", false},
		{"int16", int16(1), true},
		{"slice", []int{1}, true},
		{"map", map[int]int{1: 1}, true},
		{"map", map[int]interface{}{1: "a"}, true},
		{"struct", tstNoStringer{}, true},
		{"ptr", &tstNoStringer{}, true},
		{"float", 1.1, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := ns.KindIs(test.k, test.v)

		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestKindOf(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{nil, "invalid"},
		{"foo", "string"},
		{int16(1), "int16"},
		{[]int{1, 2}, "slice"},
		{map[int]int{1: 1}, "map"},
		{tstNoStringer{}, "struct"},
		{&tstNoStringer{}, "ptr"},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := ns.KindOf(test.v)

		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestTypeIs(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		t      string
		v      interface{}
		expect interface{}
	}{
		{"<nil>", nil, true},
		{"string", "foo", true},
		{"int16", int16(1), true},
		{"[]int", []int{1, 2}, true},
		{"map[int]int", map[int]int{1: 1}, true},
		{"map[int]interface {}", map[int]interface{}{1: "a"}, true},
		{"reflect.tstNoStringer", tstNoStringer{}, true},
		{"reflect.tstNoStringer", &tstNoStringer{}, false},
		{"*reflect.tstNoStringer", &tstNoStringer{}, true},
		{"*reflect.tstNoStringer", tstNoStringer{}, false},
		{"float", 1.1, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := ns.TypeIs(test.t, test.v)

		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestTypeIsLike(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		t      string
		v      interface{}
		expect interface{}
	}{
		{"<nil>", nil, true},
		{"string", "foo", true},
		{"reflect.tstNoStringer", tstNoStringer{}, true},
		{"reflect.tstNoStringer", &tstNoStringer{}, true},
		{"*reflect.tstNoStringer", &tstNoStringer{}, true},
		{"*reflect.tstNoStringer", tstNoStringer{}, false},
		{"float", 1.1, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := ns.TypeIsLike(test.t, test.v)

		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestTypeOf(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{nil, "<nil>"},
		{"foo", "string"},
		{int16(1), "int16"},
		{[]int{1, 2}, "[]int"},
		{map[int]int{1: 1}, "map[int]int"},
		{map[int]interface{}{1: 1}, "map[int]interface {}"},
		{tstNoStringer{}, "reflect.tstNoStringer"},
		{&tstNoStringer{}, "*reflect.tstNoStringer"},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result := ns.TypeOf(test.v)

		assert.Equal(t, test.expect, result, errMsg)
	}
}
