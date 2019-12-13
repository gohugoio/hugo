// Copyright 2018 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestExtractBaseof(t *testing.T) {
	c := qt.New(t)

	replaced := extractBaseOf(`failed: template: _default/baseof.html:37:11: executing "_default/baseof.html" at <.Parents>: can't evaluate field Parents in type *hugolib.PageOutput`)

	c.Assert(replaced, qt.Equals, "_default/baseof.html")
	c.Assert(extractBaseOf("not baseof for you"), qt.Equals, "")
	c.Assert(extractBaseOf("template: blog/baseof.html:23:11:"), qt.Equals, "blog/baseof.html")
}
