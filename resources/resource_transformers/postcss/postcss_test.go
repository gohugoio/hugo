// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postcss

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Issue 6166
func TestDecodeOptions(t *testing.T) {
	assert := require.New(t)
	opts1, err := DecodeOptions(map[string]interface{}{
		"no-map": true,
	})

	assert.NoError(err)
	assert.True(opts1.NoMap)

	opts2, err := DecodeOptions(map[string]interface{}{
		"noMap": true,
	})

	assert.NoError(err)
	assert.True(opts2.NoMap)

}
