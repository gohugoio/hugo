// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package types

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestKeyValues(t *testing.T) {
	c := qt.New(t)

	kv := NewKeyValuesStrings("key", "a1", "a2")

	c.Assert(kv.KeyString(), qt.Equals, "key")
	c.Assert(kv.Values, qt.DeepEquals, []any{"a1", "a2"})
}

func TestLowHigh(t *testing.T) {
	c := qt.New(t)

	lh := LowHigh[string]{
		Low:  2,
		High: 10,
	}

	s := "abcdefghijklmnopqrstuvwxyz"
	c.Assert(lh.IsZero(), qt.IsFalse)
	c.Assert(lh.Value(s), qt.Equals, "cdefghij")

	lhb := LowHigh[[]byte]{
		Low:  2,
		High: 10,
	}

	sb := []byte(s)
	c.Assert(lhb.IsZero(), qt.IsFalse)
	c.Assert(lhb.Value(sb), qt.DeepEquals, []byte("cdefghij"))
}
