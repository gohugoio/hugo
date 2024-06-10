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

package security

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestWhitelist(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Run("none", func(c *qt.C) {
		c.Assert(MustNewWhitelist("none", "foo").Accept("foo"), qt.IsFalse)
		c.Assert(MustNewWhitelist().Accept("foo"), qt.IsFalse)
		c.Assert(MustNewWhitelist("").Accept("foo"), qt.IsFalse)
		c.Assert(MustNewWhitelist("  ", " ").Accept("foo"), qt.IsFalse)
		c.Assert(Whitelist{}.Accept("foo"), qt.IsFalse)
	})

	c.Run("One", func(c *qt.C) {
		w := MustNewWhitelist("^foo.*")
		c.Assert(w.Accept("foo"), qt.IsTrue)
		c.Assert(w.Accept("mfoo"), qt.IsFalse)
	})

	c.Run("Multiple", func(c *qt.C) {
		w := MustNewWhitelist("^foo.*", "^bar.*")
		c.Assert(w.Accept("foo"), qt.IsTrue)
		c.Assert(w.Accept("bar"), qt.IsTrue)
		c.Assert(w.Accept("mbar"), qt.IsFalse)
	})
}
