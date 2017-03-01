// Copyright 2015 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHugoVersion(t *testing.T) {
	assert.Equal(t, "0.15-DEV", hugoVersion(0.15, 0, "-DEV"))
	assert.Equal(t, "0.17", hugoVersionNoSuffix(0.16+0.01, 0))
	assert.Equal(t, "0.20", hugoVersionNoSuffix(0.20, 0))
	assert.Equal(t, "0.15.2-DEV", hugoVersion(0.15, 2, "-DEV"))
	assert.Equal(t, "0.17.3", hugoVersionNoSuffix(0.16+0.01, 3))
}

func TestCompareVersions(t *testing.T) {
	require.Equal(t, 0, compareVersions(0.20, 0, 0.20))
	require.Equal(t, 0, compareVersions(0.20, 0, float32(0.20)))
	require.Equal(t, 0, compareVersions(0.20, 0, float64(0.20)))
	require.Equal(t, 1, compareVersions(0.19, 1, 0.20))
	require.Equal(t, 1, compareVersions(0.19, 3, "0.20.2"))
	require.Equal(t, -1, compareVersions(0.19, 1, 0.01))
	require.Equal(t, 1, compareVersions(0, 1, 3))
	require.Equal(t, 1, compareVersions(0, 1, int32(3)))
	require.Equal(t, 1, compareVersions(0, 1, int64(3)))
	require.Equal(t, 0, compareVersions(0.20, 0, "0.20"))
	require.Equal(t, 0, compareVersions(0.20, 1, "0.20.1"))
	require.Equal(t, -1, compareVersions(0.20, 1, "0.20"))
	require.Equal(t, 1, compareVersions(0.20, 0, "0.20.1"))
	require.Equal(t, 1, compareVersions(0.20, 1, "0.20.2"))
	require.Equal(t, 1, compareVersions(0.21, 1, "0.22.1"))
}
