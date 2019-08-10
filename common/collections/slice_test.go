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
	"testing"

	qt "github.com/frankban/quicktest"
)

var _ Slicer = (*tstSlicer)(nil)
var _ Slicer = (*tstSlicerIn1)(nil)
var _ Slicer = (*tstSlicerIn2)(nil)
var _ testSlicerInterface = (*tstSlicerIn1)(nil)
var _ testSlicerInterface = (*tstSlicerIn1)(nil)

type testSlicerInterface interface {
	Name() string
}

type testSlicerInterfaces []testSlicerInterface

type tstSlicerIn1 struct {
	TheName string
}

type tstSlicerIn2 struct {
	TheName string
}

type tstSlicer struct {
	TheName string
}

func (p *tstSlicerIn1) Slice(in interface{}) (interface{}, error) {
	items := in.([]interface{})
	result := make(testSlicerInterfaces, len(items))
	for i, v := range items {
		switch vv := v.(type) {
		case testSlicerInterface:
			result[i] = vv
		default:
			return nil, errors.New("invalid type")
		}

	}
	return result, nil
}

func (p *tstSlicerIn2) Slice(in interface{}) (interface{}, error) {
	items := in.([]interface{})
	result := make(testSlicerInterfaces, len(items))
	for i, v := range items {
		switch vv := v.(type) {
		case testSlicerInterface:
			result[i] = vv
		default:
			return nil, errors.New("invalid type")
		}
	}
	return result, nil
}

func (p *tstSlicerIn1) Name() string {
	return p.TheName
}

func (p *tstSlicerIn2) Name() string {
	return p.TheName
}

func (p *tstSlicer) Slice(in interface{}) (interface{}, error) {
	items := in.([]interface{})
	result := make(tstSlicers, len(items))
	for i, v := range items {
		switch vv := v.(type) {
		case *tstSlicer:
			result[i] = vv
		default:
			return nil, errors.New("invalid type")
		}
	}
	return result, nil
}

type tstSlicers []*tstSlicer

func TestSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		args     []interface{}
		expected interface{}
	}{
		{[]interface{}{"a", "b"}, []string{"a", "b"}},
		{[]interface{}{&tstSlicer{"a"}, &tstSlicer{"b"}}, tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}}},
		{[]interface{}{&tstSlicer{"a"}, "b"}, []interface{}{&tstSlicer{"a"}, "b"}},
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{nil}, []interface{}{nil}},
		{[]interface{}{5, "b"}, []interface{}{5, "b"}},
		{[]interface{}{&tstSlicerIn1{"a"}, &tstSlicerIn2{"b"}}, testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn2{"b"}}},
		{[]interface{}{&tstSlicerIn1{"a"}, &tstSlicer{"b"}}, []interface{}{&tstSlicerIn1{"a"}, &tstSlicer{"b"}}},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.args)

		result := Slice(test.args...)

		c.Assert(test.expected, qt.DeepEquals, result, errMsg)
	}

}
