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

package encoding

import (
	"fmt"
	"html/template"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tstNoStringer struct{}

func TestBase64Decode(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"YWJjMTIzIT8kKiYoKSctPUB+", "abc123!?$*&()'-=@~"},
		// errors
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.Base64Decode(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestBase64Encode(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"YWJjMTIzIT8kKiYoKSctPUB+", "WVdKak1USXpJVDhrS2lZb0tTY3RQVUIr"},
		// errors
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.Base64Encode(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestJsonify(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{[]string{"a", "b"}, template.HTML(`["a","b"]`)},
		{tstNoStringer{}, template.HTML("{}")},
		{nil, template.HTML("null")},
		// errors
		{math.NaN(), false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.Jsonify(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
