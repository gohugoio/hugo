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

package files

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestComponentFolders(t *testing.T) {
	c := qt.New(t)

	// It's important that these are absolutely right and not changed.
	c.Assert(len(componentFoldersSet), qt.Equals, len(ComponentFolders))
	c.Assert(IsComponentFolder("archetypes"), qt.Equals, true)
	c.Assert(IsComponentFolder("layouts"), qt.Equals, true)
	c.Assert(IsComponentFolder("data"), qt.Equals, true)
	c.Assert(IsComponentFolder("i18n"), qt.Equals, true)
	c.Assert(IsComponentFolder("assets"), qt.Equals, true)
	c.Assert(IsComponentFolder("resources"), qt.Equals, false)
	c.Assert(IsComponentFolder("static"), qt.Equals, true)
	c.Assert(IsComponentFolder("content"), qt.Equals, true)
	c.Assert(IsComponentFolder("foo"), qt.Equals, false)
	c.Assert(IsComponentFolder(""), qt.Equals, false)
}
