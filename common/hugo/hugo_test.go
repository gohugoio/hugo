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

package hugo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHugoInfo(t *testing.T) {
	assert := require.New(t)

	hugoInfo := NewInfo("")

	assert.Equal(CurrentVersion.Version(), hugoInfo.Version())
	assert.IsType(VersionString(""), hugoInfo.Version())
	assert.Equal(commitHash, hugoInfo.CommitHash)
	assert.Equal(buildDate, hugoInfo.BuildDate)
	assert.Equal("production", hugoInfo.Environment)
	assert.Contains(hugoInfo.Generator(), fmt.Sprintf("Hugo %s", hugoInfo.Version()))

}
