// Copyright 2021 The Hugo Authors. All rights reserved.
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

package pagekinds

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestKind(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	// Add tests for these constants to make sure they don't change
	c.Assert(Page, qt.Equals, "page")
	c.Assert(Home, qt.Equals, "home")
	c.Assert(Section, qt.Equals, "section")
	c.Assert(Taxonomy, qt.Equals, "taxonomy")
	c.Assert(Term, qt.Equals, "term")

	c.Assert(Get("TAXONOMYTERM"), qt.Equals, Taxonomy)
	c.Assert(Get("Taxonomy"), qt.Equals, Taxonomy)
	c.Assert(Get("Page"), qt.Equals, Page)
	c.Assert(Get("Home"), qt.Equals, Home)
	c.Assert(Get("SEction"), qt.Equals, Section)
}
