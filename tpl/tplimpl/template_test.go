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
package tplimpl

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNeedsBaseTemplate(t *testing.T) {
	c := qt.New(t)

	c.Assert(needsBaseTemplate(`{{ define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`{{define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`{{-  define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`{{-define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`
	
	{{-define "main" }}
	
	`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`    {{ define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`
	{{ define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`  A  {{ define "main" }}`), qt.Equals, false)
	c.Assert(needsBaseTemplate(`  {{ printf "foo" }}`), qt.Equals, false)
	c.Assert(needsBaseTemplate(`{{/* comment */}}    {{ define "main" }}`), qt.Equals, true)
	c.Assert(needsBaseTemplate(`     {{/* comment */}}  A  {{ define "main" }}`), qt.Equals, false)
}
