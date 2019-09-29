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

	c.Assert(Item{Val: []byte("3.14")}.ValTyped(), qt.Equals, float64(3.14))
	c.Assert(Item{Val: []byte(".14")}.ValTyped(), qt.Equals, float64(.14))
	c.Assert(Item{Val: []byte("314")}.ValTyped(), qt.Equals, 314)
	c.Assert(Item{Val: []byte("314x")}.ValTyped(), qt.Equals, "314x")
	c.Assert(Item{Val: []byte("314 ")}.ValTyped(), qt.Equals, "314 ")
	c.Assert(Item{Val: []byte("314"), isString: true}.ValTyped(), qt.Equals, "314")
	c.Assert(Item{Val: []byte("true")}.ValTyped(), qt.Equals, true)
	c.Assert(Item{Val: []byte("false")}.ValTyped(), qt.Equals, false)
	c.Assert(Item{Val: []byte("trues")}.ValTyped(), qt.Equals, "trues")

}
