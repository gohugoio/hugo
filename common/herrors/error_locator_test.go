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

// Package errors contains common Hugo errors and error related utilities.
package herrors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorLocator(t *testing.T) {
	assert := require.New(t)

	lineMatcher := func(le FileError, lineno int, line string) bool {
		return strings.Contains(line, "THEONE")
	}

	lines := `LINE 1
LINE 2
LINE 3
LINE 4
This is THEONE
LINE 6
LINE 7
LINE 8
`

	location := locateErrorInString(nil, lines, lineMatcher)
	assert.Equal([]string{"LINE 3", "LINE 4", "This is THEONE", "LINE 6", "LINE 7"}, location.Lines)

	assert.Equal(5, location.LineNumber)
	assert.Equal(2, location.Pos)

	assert.Equal([]string{"This is THEONE"}, locateErrorInString(nil, `This is THEONE`, lineMatcher).Lines)

	location = locateErrorInString(nil, `L1
This is THEONE
L2
`, lineMatcher)
	assert.Equal(1, location.Pos)
	assert.Equal([]string{"L1", "This is THEONE", "L2"}, location.Lines)

	location = locateErrorInString(nil, `This is THEONE
L2
`, lineMatcher)
	assert.Equal(0, location.Pos)
	assert.Equal([]string{"This is THEONE", "L2"}, location.Lines)

	location = locateErrorInString(nil, `L1
This THEONE
`, lineMatcher)
	assert.Equal([]string{"L1", "This THEONE"}, location.Lines)
	assert.Equal(1, location.Pos)

	location = locateErrorInString(nil, `L1
L2
This THEONE
`, lineMatcher)
	assert.Equal([]string{"L1", "L2", "This THEONE"}, location.Lines)
	assert.Equal(2, location.Pos)

	location = locateErrorInString(nil, "NO MATCH", lineMatcher)
	assert.Equal(-1, location.LineNumber)
	assert.Equal(-1, location.Pos)
	assert.Equal(0, len(location.Lines))

	lineMatcher = func(le FileError, lineno int, line string) bool {
		return lineno == 6
	}
	location = locateErrorInString(nil, `A
B
C
D
E
F
G
H
I
J`, lineMatcher)

	assert.Equal([]string{"D", "E", "F", "G", "H"}, location.Lines)
	assert.Equal(6, location.LineNumber)
	assert.Equal(2, location.Pos)

	// Test match EOF
	lineMatcher = func(le FileError, lineno int, line string) bool {
		return lineno == 4
	}

	location = locateErrorInString(nil, `A
B
C
`, lineMatcher)

	assert.Equal([]string{"B", "C", ""}, location.Lines)
	assert.Equal(4, location.LineNumber)
	assert.Equal(2, location.Pos)

}
