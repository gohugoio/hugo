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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
)

func TestIgnoreDotFilesAndDirectories(t *testing.T) {
	c := qt.New(t)

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
		{"foobar/barfoo.md~", true, nil},
		{".foobar/barfoo.md~", true, nil},
		{"foobar~/barfoo.md", false, nil},
		{"foobar/bar~foo.md", false, nil},
		{"foobar/foo.md", true, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.html", false, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.md", true, []string{"foo.md$"}},
		{"foobar/foo.md", true, []string{"*", "\\.md$", "\\.boo$"}},
		{"foobar/.#content.md", true, []string{"/\\.#"}},
		{".#foobar.md", true, []string{"^\\.#"}},
	}

	for i, test := range tests {
		v := newTestConfig()
		v.Set("ignoreFiles", test.ignoreFilesRegexpes)
		fs := hugofs.NewMem(v)
		ps, err := helpers.NewPathSpec(fs, v, nil)
		c.Assert(err, qt.IsNil)

		s := NewSourceSpec(ps, fs.Source)

		if ignored := s.IgnoreFile(filepath.FromSlash(test.path)); test.ignore != ignored {
			t.Errorf("[%d] File not ignored", i)
		}
	}
}
