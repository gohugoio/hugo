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

package types

import (
	"encoding/json"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestToStringSlicePreserveString(t *testing.T) {
	c := qt.New(t)

	c.Assert(ToStringSlicePreserveString("Hugo"), qt.DeepEquals, []string{"Hugo"})
	c.Assert(ToStringSlicePreserveString(qt.Commentf("Hugo")), qt.DeepEquals, []string{"Hugo"})
	c.Assert(ToStringSlicePreserveString([]any{"A", "B"}), qt.DeepEquals, []string{"A", "B"})
	c.Assert(ToStringSlicePreserveString([]int{1, 3}), qt.DeepEquals, []string{"1", "3"})
	c.Assert(ToStringSlicePreserveString(nil), qt.IsNil)
}

func TestToString(t *testing.T) {
	c := qt.New(t)

	c.Assert(ToString([]byte("Hugo")), qt.Equals, "Hugo")
	c.Assert(ToString(json.RawMessage("Hugo")), qt.Equals, "Hugo")
}

func TestToDuration(t *testing.T) {
	c := qt.New(t)

	c.Assert(ToDuration("200ms"), qt.Equals, 200*time.Millisecond)
	c.Assert(ToDuration("200"), qt.Equals, 200*time.Millisecond)
	c.Assert(ToDuration("4m"), qt.Equals, 4*time.Minute)
	c.Assert(ToDuration("asdfadf"), qt.Equals, time.Duration(0))
}
