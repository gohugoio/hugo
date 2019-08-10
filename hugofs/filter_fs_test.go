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

package hugofs

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestLangInfoFrom(t *testing.T) {

	langs := map[string]int{
		"sv": 10,
		"en": 20,
	}

	c := qt.New(t)

	tests := []struct {
		input    string
		expected []string
	}{
		{"page.sv.md", []string{"sv", "page", "page.md"}},
		{"page.en.md", []string{"en", "page", "page.md"}},
		{"page.no.md", []string{"", "page.no", "page.no.md"}},
		{filepath.FromSlash("tc-lib-color/class-Com.Tecnick.Color.Css"), []string{"", "class-Com.Tecnick.Color", "class-Com.Tecnick.Color.Css"}},
		{filepath.FromSlash("class-Com.Tecnick.Color.sv.Css"), []string{"sv", "class-Com.Tecnick.Color", "class-Com.Tecnick.Color.Css"}},
	}

	for _, test := range tests {
		v1, v2, v3 := langInfoFrom(langs, test.input)
		c.Assert([]string{v1, v2, v3}, qt.DeepEquals, test.expected)
	}

}
