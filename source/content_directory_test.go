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

package source

import (
	"testing"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"
)

func TestIgnoreDotFilesAndDirectories(t *testing.T) {

	tests := []struct {
		path                string
		ignore              bool
		ignoreFilesRegexpes interface{}
	}{
		{".foobar/", true, nil},
		{"foobar/.barfoo/", true, nil},
		{"barfoo.md", false, nil},
		{"foobar/barfoo.md", false, nil},
		{"foobar/.barfoo.md", true, nil},
		{".barfoo.md", true, nil},
		{".md", true, nil},
		{"", true, nil},
		{"foobar/barfoo.md~", true, nil},
		{".foobar/barfoo.md~", true, nil},
		{"foobar~/barfoo.md", false, nil},
		{"foobar/bar~foo.md", false, nil},
		{"foobar/foo.md", true, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.html", false, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.md", true, []string{"^foo"}},
		{"foobar/foo.md", false, []string{"*", "\\.md$", "\\.boo$"}},
		{"foobar/.#content.md", true, []string{"/\\.#"}},
		{".#foobar.md", true, []string{"^\\.#"}},
	}

	for _, test := range tests {

		v := viper.New()
		v.Set("ignoreFiles", test.ignoreFilesRegexpes)

		s := NewSourceSpec(v, hugofs.NewMem(v))

		if ignored := s.isNonProcessablePath(test.path); test.ignore != ignored {
			t.Errorf("File not ignored.  Expected: %t, got: %t", test.ignore, ignored)
		}
	}
}
