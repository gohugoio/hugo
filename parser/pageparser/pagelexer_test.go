// Copyright 2018 The Hugo Authors. All rights reserved.
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

package pageparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinPositiveIndex(t *testing.T) {
	assert := require.New(t)
	assert.Equal(1, minPositiveIndex(4, 1, 2, 3))
	assert.Equal(2, minPositiveIndex(4, 0, -2, 2, 5))
	assert.Equal(-1, minPositiveIndex())
	assert.Equal(-1, minPositiveIndex(-2, -3))

}
