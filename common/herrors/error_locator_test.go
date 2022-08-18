// Copyright 2022 The Hugo Authors. All rights reserved.
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

// Package herrors contains common Hugo errors and error related utilities.
package herrors

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestErrorLocator(t *testing.T) {
	c := qt.New(t)

	lineMatcher := func(m LineMatcher) int {
		if strings.Contains(m.Line, "THEONE") {
			return 1
		}
		return -1
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

	location := locateErrorInString(lines, lineMatcher)
	pos := location.Position
	c.Assert(location.Lines, qt.DeepEquals, []string{"LINE 3", "LINE 4", "This is THEONE", "LINE 6", "LINE 7"})

	c.Assert(pos.LineNumber, qt.Equals, 5)
	c.Assert(location.LinesPos, qt.Equals, 2)

	locate := func(s string, m LineMatcherFn) *ErrorContext {
		ctx := locateErrorInString(s, m)
		return ctx
	}

	c.Assert(locate(`This is THEONE`, lineMatcher).Lines, qt.DeepEquals, []string{"This is THEONE"})

	location = locateErrorInString(`L1
This is THEONE
L2
`, lineMatcher)
	pos = location.Position
	c.Assert(pos.LineNumber, qt.Equals, 2)
	c.Assert(location.LinesPos, qt.Equals, 1)
	c.Assert(location.Lines, qt.DeepEquals, []string{"L1", "This is THEONE", "L2", ""})

	location = locate(`This is THEONE
L2
`, lineMatcher)
	c.Assert(location.LinesPos, qt.Equals, 0)
	c.Assert(location.Lines, qt.DeepEquals, []string{"This is THEONE", "L2", ""})

	location = locate(`L1
This THEONE
`, lineMatcher)
	c.Assert(location.Lines, qt.DeepEquals, []string{"L1", "This THEONE", ""})
	c.Assert(location.LinesPos, qt.Equals, 1)

	location = locate(`L1
L2
This THEONE
`, lineMatcher)
	c.Assert(location.Lines, qt.DeepEquals, []string{"L1", "L2", "This THEONE", ""})
	c.Assert(location.LinesPos, qt.Equals, 2)

	location = locateErrorInString("NO MATCH", lineMatcher)
	pos = location.Position
	c.Assert(pos.LineNumber, qt.Equals, 0)
	c.Assert(location.LinesPos, qt.Equals, -1)
	c.Assert(len(location.Lines), qt.Equals, 0)

	lineMatcher = func(m LineMatcher) int {
		if m.LineNumber == 6 {
			return 1
		}
		return -1
	}

	location = locateErrorInString(`A
B
C
D
E
F
G
H
I
J`, lineMatcher)
	pos = location.Position

	c.Assert(location.Lines, qt.DeepEquals, []string{"D", "E", "F", "G", "H"})
	c.Assert(pos.LineNumber, qt.Equals, 6)
	c.Assert(location.LinesPos, qt.Equals, 2)

	// Test match EOF
	lineMatcher = func(m LineMatcher) int {
		if m.LineNumber == 4 {
			return 1
		}
		return -1
	}

	location = locateErrorInString(`A
B
C
`, lineMatcher)

	pos = location.Position

	c.Assert(location.Lines, qt.DeepEquals, []string{"B", "C", ""})
	c.Assert(pos.LineNumber, qt.Equals, 4)
	c.Assert(location.LinesPos, qt.Equals, 2)

	offsetMatcher := func(m LineMatcher) int {
		if m.Offset == 1 {
			return 1
		}
		return -1
	}

	location = locateErrorInString(`A
B
C
D
E`, offsetMatcher)

	pos = location.Position

	c.Assert(location.Lines, qt.DeepEquals, []string{"A", "B", "C", "D"})
	c.Assert(pos.LineNumber, qt.Equals, 2)
	c.Assert(location.LinesPos, qt.Equals, 1)
}
