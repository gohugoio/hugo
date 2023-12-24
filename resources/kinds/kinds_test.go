// Copyright 2024 The Hugo Authors. All rights reserved.
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

package kinds

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestKind(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	// Add tests for these constants to make sure they don't change
	c.Assert(KindPage, qt.Equals, "page")
	c.Assert(KindHome, qt.Equals, "home")
	c.Assert(KindSection, qt.Equals, "section")
	c.Assert(KindTaxonomy, qt.Equals, "taxonomy")
	c.Assert(KindTerm, qt.Equals, "term")

	c.Assert(GetKindMain("TAXONOMYTERM"), qt.Equals, KindTaxonomy)
	c.Assert(GetKindMain("Taxonomy"), qt.Equals, KindTaxonomy)
	c.Assert(GetKindMain("Page"), qt.Equals, KindPage)
	c.Assert(GetKindMain("Home"), qt.Equals, KindHome)
	c.Assert(GetKindMain("SEction"), qt.Equals, KindSection)

	c.Assert(GetKindAny("Page"), qt.Equals, KindPage)
	c.Assert(GetKindAny("Robotstxt"), qt.Equals, KindRobotsTXT)
}
