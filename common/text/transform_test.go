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

package text

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestRemoveAccents(t *testing.T) {
	c := qt.New(t)

	c.Assert(string(RemoveAccents([]byte("Resumé"))), qt.Equals, "Resume")
	c.Assert(string(RemoveAccents([]byte("Hugo Rocks!"))), qt.Equals, "Hugo Rocks!")
	c.Assert(string(RemoveAccentsString("Resumé")), qt.Equals, "Resume")

}
