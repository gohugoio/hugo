// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"github.com/spf13/hugo/tpl"
	"testing"
)

const (
	win_base = "c:\\a\\windows\\path\\layout"
	win_path = "c:\\a\\windows\\path\\layout\\sub1\\index.html"
)

func TestTemplatePathSeparator(t *testing.T) {
	tmpl := new(tpl.GoHTMLTemplate)
	if name := tmpl.GenerateTemplateNameFrom(win_base, win_path); name != "sub1/index.html" {
		t.Fatalf("Template name incorrect. got %s but expected %s", name, "sub1/index.html")
	}
}
