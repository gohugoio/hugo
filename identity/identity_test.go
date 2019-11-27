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

package identity

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIdentityManager(t *testing.T) {
	c := qt.New(t)

	id1 := testIdentity{name: "id1"}
	im := NewManager(id1)

	c.Assert(im.Search(id1).GetIdentity(), qt.Equals, id1)
	c.Assert(im.Search(testIdentity{name: "notfound"}), qt.Equals, nil)
}

type testIdentity struct {
	name string
}

func (id testIdentity) GetIdentity() Identity {
	return id
}

func (id testIdentity) Name() string {
	return id.name
}
