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

package scss

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestReplaceRegularCSSImports(t *testing.T) {
	c := qt.New(t)

	scssWithImport := `
	
@import "moo";
@import "regular.css";
@import "moo";
@import "another.css";
@import "foo.scss";

/* foo */`

	scssWithoutImport := `
@import "moo";
/* foo */`

	res, replaced := replaceRegularImportsIn(scssWithImport)
	c.Assert(replaced, qt.Equals, true)
	c.Assert(res, qt.Equals, "\n\t\n@import \"moo\";\n/* HUGO_IMPORT_START regular.css HUGO_IMPORT_END */\n@import \"moo\";\n/* HUGO_IMPORT_START another.css HUGO_IMPORT_END */\n@import \"foo.scss\";\n\n/* foo */")

	res2, replaced2 := replaceRegularImportsIn(scssWithoutImport)
	c.Assert(replaced2, qt.Equals, false)
	c.Assert(res2, qt.Equals, scssWithoutImport)

	reverted := replaceRegularImportsOut(res)
	c.Assert(reverted, qt.Equals, scssWithImport)

}
