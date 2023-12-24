// Copyright 2024 The Hugo Authors. All rights reserved.
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

package source_test

import (
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"golang.org/x/text/unicode/norm"
)

func TestUnicodeNorm(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Darwin OS")
	}
	t.Parallel()
	files := `
-- hugo.toml --
-- content/å.md --
-- content/é.md --
-- content/å/å.md --
-- content/é/é.md --
-- layouts/_default/single.html --
Title: {{ .Title }}|File: {{ .File.Path}}
`
	b := hugolib.Test(t, files, hugolib.TestOptWithNFDOnDarwin())

	for _, p := range b.H.Sites[0].RegularPages() {
		f := p.File()
		b.Assert(norm.NFC.IsNormalString(f.Path()), qt.IsTrue)
		b.Assert(norm.NFC.IsNormalString(f.Dir()), qt.IsTrue)
		b.Assert(norm.NFC.IsNormalString(f.Filename()), qt.IsTrue)
		b.Assert(norm.NFC.IsNormalString(f.BaseFileName()), qt.IsTrue)
	}
}
