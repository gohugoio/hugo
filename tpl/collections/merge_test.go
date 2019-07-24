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
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/parser"

	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	ns := New(&deps.Deps{})

	simpleMap := map[string]interface{}{"a": 1, "b": 2}

	for i, test := range []struct {
		name   string
		dst    interface{}
		src    interface{}
		expect interface{}
		isErr  bool
	}{
		{
			"basic",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 42, "c": 3},
			map[string]interface{}{"a": 1, "b": 2, "c": 3}, false,
		},
		{
			"basic case insensitive",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"A": 42, "c": 3},
			map[string]interface{}{"a": 1, "b": 2, "c": 3}, false,
		},
		{
			"nested",
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"d": 1, "e": 2}},
			map[string]interface{}{"a": 42, "c": 3, "b": map[string]interface{}{"d": 55, "e": 66, "f": 3}},
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"d": 1, "e": 2, "f": 3}, "c": 3}, false,
		},
		{"src nil", simpleMap, nil, simpleMap, false},
		// Error cases.
		{"dst not a map", "not a map", nil, nil, true},
		{"src not a map", simpleMap, "not a map", nil, true},
		{"diferent map typs", simpleMap, map[int]interface{}{32: "a"}, nil, true},
		{"all nil", nil, nil, nil, true},
	} {

		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			errMsg := fmt.Sprintf("[%d] %v", i, test)

			assert := require.New(t)

			srcStr, dstStr := fmt.Sprint(test.src), fmt.Sprint(test.dst)

			result, err := ns.Merge(test.src, test.dst)

			if test.isErr {
				assert.Error(err, errMsg)
				return
			}

			assert.NoError(err, errMsg)
			assert.Equal(test.expect, result, errMsg)

			// map sort in fmt was fixed in go 1.12.
			if !strings.HasPrefix(runtime.Version(), "go1.11") {
				// Verify that the original maps are preserved.
				assert.Equal(srcStr, fmt.Sprint(test.src))
				assert.Equal(dstStr, fmt.Sprint(test.dst))
			}
		})
	}
}

func TestMergeDataFormats(t *testing.T) {
	assert := require.New(t)
	ns := New(&deps.Deps{})

	toml1 := `
V1 = "v1_1"

[V2s]
V21 = "v21_1"

`

	toml2 := `
V1 = "v1_2"
V2 = "v2_2"

[V2s]
V21 = "v21_2"
V22 = "v22_2"

`

	meta1, err := metadecoders.Default.UnmarshalToMap([]byte(toml1), metadecoders.TOML)
	assert.NoError(err)
	meta2, err := metadecoders.Default.UnmarshalToMap([]byte(toml2), metadecoders.TOML)
	assert.NoError(err)

	for _, format := range []metadecoders.Format{metadecoders.JSON, metadecoders.YAML, metadecoders.TOML} {

		var dataStr1, dataStr2 bytes.Buffer
		err = parser.InterfaceToConfig(meta1, format, &dataStr1)
		assert.NoError(err)
		err = parser.InterfaceToConfig(meta2, format, &dataStr2)
		assert.NoError(err)

		dst, err := metadecoders.Default.UnmarshalToMap(dataStr1.Bytes(), format)
		assert.NoError(err)
		src, err := metadecoders.Default.UnmarshalToMap(dataStr2.Bytes(), format)
		assert.NoError(err)

		merged, err := ns.Merge(src, dst)
		assert.NoError(err)

		assert.Equal(map[string]interface{}{"V1": "v1_1", "V2": "v2_2", "V2s": map[string]interface{}{"V21": "v21_1", "V22": "v22_2"}}, merged)
	}
}

func TestCaseInsensitiveMapLookup(t *testing.T) {
	assert := require.New(t)

	m1 := reflect.ValueOf(map[string]interface{}{
		"a": 1,
		"B": 2,
	})

	m2 := reflect.ValueOf(map[int]interface{}{
		1: 1,
		2: 2,
	})

	var found bool

	a, found := caseInsensitiveLookup(m1, reflect.ValueOf("A"))
	assert.True(found)
	assert.Equal(1, a.Interface())

	b, found := caseInsensitiveLookup(m1, reflect.ValueOf("b"))
	assert.True(found)
	assert.Equal(2, b.Interface())

	two, found := caseInsensitiveLookup(m2, reflect.ValueOf(2))
	assert.True(found)
	assert.Equal(2, two.Interface())
}
