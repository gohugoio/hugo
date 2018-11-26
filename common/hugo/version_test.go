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

package hugo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHugoVersion(t *testing.T) {
	assert.Equal(t, "0.15-DEV", version(0.15, 0, "-DEV"))
	assert.Equal(t, "0.15.2-DEV", version(0.15, 2, "-DEV"))

	v := Version{Number: 0.21, PatchLevel: 0, Suffix: "-DEV"}

	require.Equal(t, v.ReleaseVersion().String(), "0.21")
	require.Equal(t, "0.21-DEV", v.String())
	require.Equal(t, "0.22", v.Next().String())
	nextVersionString := v.Next().Version()
	require.Equal(t, "0.22", nextVersionString.String())
	require.True(t, nextVersionString.Eq("0.22"))
	require.False(t, nextVersionString.Eq("0.21"))
	require.True(t, nextVersionString.Eq(nextVersionString))
	require.Equal(t, "0.20.3", v.NextPatchLevel(3).String())
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
	require.Equal(t, -1, compareVersions(0.22, 0, "0.22-DEV"))
	require.Equal(t, 1, compareVersions(0.22, 0, "0.22.1-DEV"))
	require.Equal(t, 1, compareVersionsWithSuffix(0.22, 0, "-DEV", "0.22"))
	require.Equal(t, -1, compareVersionsWithSuffix(0.22, 1, "-DEV", "0.22"))
	require.Equal(t, 0, compareVersionsWithSuffix(0.22, 1, "-DEV", "0.22.1-DEV"))

}

func TestParseHugoVersion(t *testing.T) {
	require.Equal(t, "0.25", MustParseVersion("0.25").String())
	require.Equal(t, "0.25.2", MustParseVersion("0.25.2").String())
	require.Equal(t, "0.25-test", MustParseVersion("0.25-test").String())
	require.Equal(t, "0.25-DEV", MustParseVersion("0.25-DEV").String())

}
