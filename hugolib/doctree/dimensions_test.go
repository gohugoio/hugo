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

package doctree

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDimensionFlag(t *testing.T) {
	c := qt.New(t)

	var zero DimensionFlag
	var d DimensionFlag
	var o DimensionFlag = 1
	var p DimensionFlag = 12

	c.Assert(d.Has(o), qt.Equals, false)
	d = d.Set(o)
	c.Assert(d.Has(o), qt.Equals, true)
	c.Assert(d.Has(d), qt.Equals, true)
	c.Assert(func() { zero.Index() }, qt.PanicMatches, "dimension flag not set")
	c.Assert(DimensionLanguage.Index(), qt.Equals, 0)
	c.Assert(p.Index(), qt.Equals, 11)
}
