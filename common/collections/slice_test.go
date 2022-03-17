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

var (
	_ Slicer              = (*tstSlicer)(nil)
	_ Slicer              = (*tstSlicerIn1)(nil)
	_ Slicer              = (*tstSlicerIn2)(nil)
	_ testSlicerInterface = (*tstSlicerIn1)(nil)
	_ testSlicerInterface = (*tstSlicerIn1)(nil)
)

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

func (p *tstSlicerIn1) Slice(in any) (any, error) {
	items := in.([]any)
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

func (p *tstSlicerIn2) Slice(in any) (any, error) {
	items := in.([]any)
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

func (p *tstSlicer) Slice(in any) (any, error) {
	items := in.([]any)
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
		args     []any
		expected any
	}{
		{[]any{"a", "b"}, []string{"a", "b"}},
		{[]any{&tstSlicer{"a"}, &tstSlicer{"b"}}, tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}}},
		{[]any{&tstSlicer{"a"}, "b"}, []any{&tstSlicer{"a"}, "b"}},
		{[]any{}, []any{}},
		{[]any{nil}, []any{nil}},
		{[]any{5, "b"}, []any{5, "b"}},
		{[]any{&tstSlicerIn1{"a"}, &tstSlicerIn2{"b"}}, testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn2{"b"}}},
		{[]any{&tstSlicerIn1{"a"}, &tstSlicer{"b"}}, []any{&tstSlicerIn1{"a"}, &tstSlicer{"b"}}},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.args)

		result := Slice(test.args...)

		c.Assert(test.expected, qt.DeepEquals, result, errMsg)
	}
}
