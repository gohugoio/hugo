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

package commands

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

// Issue #1133
func TestNewContentPathSectionWithForwardSlashes(t *testing.T) {
	c := qt.New(t)
	p, s := newContentPathSection(nil, "/post/new.md")
	c.Assert(p, qt.Equals, filepath.FromSlash("/post/new.md"))
	c.Assert(s, qt.Equals, "post")
}
