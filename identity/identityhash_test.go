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

package identity

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHashString(t *testing.T) {
	c := qt.New(t)

	c.Assert(HashString("a", "b"), qt.Equals, "2712570657419664240")
	c.Assert(HashString("ab"), qt.Equals, "590647783936702392")

	var vals []any = []any{"a", "b", tstKeyer{"c"}}

	c.Assert(HashString(vals...), qt.Equals, "12599484872364427450")
	c.Assert(vals[2], qt.Equals, tstKeyer{"c"})
}

type tstKeyer struct {
	key string
}

func (t tstKeyer) Key() string {
	return t.key
}

func (t tstKeyer) String() string {
	return "key: " + t.key
}
