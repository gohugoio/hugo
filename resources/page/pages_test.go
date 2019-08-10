// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestProbablyEq(t *testing.T) {

	p1, p2, p3 := &testPage{title: "p1"}, &testPage{title: "p2"}, &testPage{title: "p3"}
	pages12 := Pages{p1, p2}
	pages21 := Pages{p2, p1}
	pages123 := Pages{p1, p2, p3}

	t.Run("Pages", func(t *testing.T) {
		c := qt.New(t)

		c.Assert(pages12.ProbablyEq(pages12), qt.Equals, true)
		c.Assert(pages123.ProbablyEq(pages12), qt.Equals, false)
		c.Assert(pages12.ProbablyEq(pages21), qt.Equals, false)
	})

	t.Run("PageGroup", func(t *testing.T) {
		c := qt.New(t)

		c.Assert(PageGroup{Key: "a", Pages: pages12}.ProbablyEq(PageGroup{Key: "a", Pages: pages12}), qt.Equals, true)
		c.Assert(PageGroup{Key: "a", Pages: pages12}.ProbablyEq(PageGroup{Key: "b", Pages: pages12}), qt.Equals, false)

	})

	t.Run("PagesGroup", func(t *testing.T) {
		c := qt.New(t)

		pg1, pg2 := PageGroup{Key: "a", Pages: pages12}, PageGroup{Key: "b", Pages: pages123}

		c.Assert(PagesGroup{pg1, pg2}.ProbablyEq(PagesGroup{pg1, pg2}), qt.Equals, true)
		c.Assert(PagesGroup{pg1, pg2}.ProbablyEq(PagesGroup{pg2, pg1}), qt.Equals, false)

	})

}

func TestToPages(t *testing.T) {
	c := qt.New(t)

	p1, p2 := &testPage{title: "p1"}, &testPage{title: "p2"}
	pages12 := Pages{p1, p2}

	mustToPages := func(in interface{}) Pages {
		p, err := ToPages(in)
		c.Assert(err, qt.IsNil)
		return p
	}

	c.Assert(mustToPages(nil), eq, Pages{})
	c.Assert(mustToPages(pages12), eq, pages12)
	c.Assert(mustToPages([]Page{p1, p2}), eq, pages12)
	c.Assert(mustToPages([]interface{}{p1, p2}), eq, pages12)

	_, err := ToPages("not a page")
	c.Assert(err, qt.Not(qt.IsNil))
}
