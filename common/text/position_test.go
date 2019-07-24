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

package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPositionStringFormatter(t *testing.T) {
	assert := require.New(t)

	pos := Position{Filename: "/my/file.txt", LineNumber: 12, ColumnNumber: 13, Offset: 14}

	assert.Equal("/my/file.txt|13|12", createPositionStringFormatter(":file|:col|:line")(pos))
	assert.Equal("13|/my/file.txt|12", createPositionStringFormatter(":col|:file|:line")(pos))
	assert.Equal("好:13", createPositionStringFormatter("好::col")(pos))
	assert.Equal("\"/my/file.txt:12:13\"", createPositionStringFormatter("")(pos))
	assert.Equal("\"/my/file.txt:12:13\"", pos.String())
}
