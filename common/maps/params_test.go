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

package maps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetNestedParam(t *testing.T) {
	m := map[string]interface{}{
		"first":           1,
		"with_underscore": 2,
		"nested": map[string]interface{}{
			"color": "blue",
		},
	}

	assert := require.New(t)

	must := func(keyStr, separator string, candidates ...map[string]interface{}) interface{} {
		v, err := GetNestedParam(keyStr, separator, candidates...)
		assert.NoError(err)
		return v
	}

	assert.Equal(1, must("first", "_", m))
	assert.Equal(1, must("First", "_", m))
	assert.Equal(2, must("with_underscore", "_", m))
	assert.Equal("blue", must("nested_color", "_", m))
}
