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

package herrors

import (
	"fmt"
	"strings"
	"testing"

	"errors"

	"github.com/gohugoio/hugo/common/text"

	qt "github.com/frankban/quicktest"
)

func TestNewFileError(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	fe := NewFileErrorFromName(errors.New("bar"), "foo.html")
	c.Assert(fe.Error(), qt.Equals, `"foo.html:1:1": bar`)

	lines := ""
	for i := 1; i <= 100; i++ {
		lines += fmt.Sprintf("line %d\n", i)
	}

	fe.UpdatePosition(text.Position{LineNumber: 32, ColumnNumber: 2})
	c.Assert(fe.Error(), qt.Equals, `"foo.html:32:2": bar`)
	fe.UpdatePosition(text.Position{LineNumber: 0, ColumnNumber: 0, Offset: 212})
	fe.UpdateContent(strings.NewReader(lines), nil)
	c.Assert(fe.Error(), qt.Equals, `"foo.html:32:0": bar`)
	errorContext := fe.ErrorContext()
	c.Assert(errorContext, qt.IsNotNil)
	c.Assert(errorContext.Lines, qt.DeepEquals, []string{"line 30", "line 31", "line 32", "line 33", "line 34"})
	c.Assert(errorContext.LinesPos, qt.Equals, 2)
	c.Assert(errorContext.ChromaLexer, qt.Equals, "go-html-template")

}

func TestNewFileErrorExtractFromMessage(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for i, test := range []struct {
		in           error
		offset       int
		lineNumber   int
		columnNumber int
	}{
		{errors.New("no line number for you"), 0, 1, 1},
		{errors.New(`template: _default/single.html:4:15: executing "_default/single.html" at <.Titles>: can't evaluate field Titles in type *hugolib.PageOutput`), 0, 4, 15},
		{errors.New("parse failed: template: _default/bundle-resource-meta.html:11: unexpected in operand"), 0, 11, 1},
		{errors.New(`failed:: template: _default/bundle-resource-meta.html:2:7: executing "main" at <.Titles>`), 0, 2, 7},
		{errors.New(`failed to load translations: (6, 7): was expecting token =, but got "g" instead`), 0, 6, 7},
		{errors.New(`execute of template failed: template: index.html:2:5: executing "index.html" at <partial "foo.html" .>: error calling partial: "/layouts/partials/foo.html:3:6": execute of template failed: template: partials/foo.html:3:6: executing "partials/foo.html" at <.ThisDoesNotExist>: can't evaluate field ThisDoesNotExist in type *hugolib.pageStat`), 0, 2, 5},
	} {

		got := NewFileErrorFromName(test.in, "test.txt")

		errMsg := qt.Commentf("[%d][%T]", i, got)

		pos := got.Position()
		c.Assert(pos.LineNumber, qt.Equals, test.lineNumber, errMsg)
		c.Assert(pos.ColumnNumber, qt.Equals, test.columnNumber, errMsg)
		c.Assert(errors.Unwrap(got), qt.Not(qt.IsNil))
	}
}
