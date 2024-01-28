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

package pageparser

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestItemValTyped(t *testing.T) {
	c := qt.New(t)

	source := []byte("3.14")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, float64(3.14))
	source = []byte(".14")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, float64(0.14))
	source = []byte("314")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, 314)
	source = []byte("314")
	c.Assert(Item{low: 0, high: len(source), isString: true}.ValTyped(source), qt.Equals, "314")
	source = []byte("314x")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "314x")
	source = []byte("314 ")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "314 ")
	source = []byte("true")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, true)
	source = []byte("false")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, false)
	source = []byte("falsex")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "falsex")
	source = []byte("xfalse")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "xfalse")
	source = []byte("truex")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "truex")
	source = []byte("xtrue")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "xtrue")
}
