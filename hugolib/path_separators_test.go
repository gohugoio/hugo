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
	"path/filepath"
	"strings"
	"testing"
)

var SIMPLE_PAGE_YAML = `---
contenttype: ""
---
Sample Text
`

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), filepath.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in NewPageFrom")
	}
	if p.Section() != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestNewPageWithFilePath(t *testing.T) {
	toCheck := []struct {
		input   string
		section string
		layout  []string
	}{
		{filepath.Join("sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{filepath.Join("content", "foobar.html"), "", L("page/single.html", "_default/single.html")},
		{filepath.Join("content", "sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{filepath.Join("content", "dub", "sub", "foobar.html"), "dub", L("dub/single.html", "_default/single.html")},
	}

	for i, el := range toCheck {
		p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), el.input)
		if err != nil {
			t.Errorf("[%d] Reading from SIMPLE_PAGE_YAML resulted in an error: %s", i, err)
		}
		if p.Section() != el.section {
			t.Errorf("[%d] Section incorrect page %s. got %s but expected %s", i, el.input, p.Section(), el.section)
		}

		for _, y := range el.layout {
			el.layout = append(el.layout, "theme/"+y)
		}

		if !listEqual(p.layouts(), el.layout) {
			t.Errorf("[%d] Layout incorrect. got '%s' but expected '%s'", i, p.layouts(), el.layout)
		}
	}
}
