// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hugofs_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestMountRestrictTheme(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
theme = "mytheme"
[[module.mounts]]
source = '../file2.txt'
target = 'assets/file2.txt'
-- themes/mytheme/hugo.toml --
[[module.mounts]]
source = '../../file1.txt'
target = 'assets/file1.txt'
-- file1.txt --
file1
-- file2.txt --
file2
-- layouts/all.html --
All.
`
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "mount source must be a local path for modules/themes")
}
