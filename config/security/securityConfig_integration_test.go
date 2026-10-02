// Copyright 2026 The Hugo Authors. All rights reserved.
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

package security_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/hugolib"
)

// We have similar tests about this elsewhere, but add one specific one to demonstrate that by default
// it should not be possible for a theme to losen the security defaults.
func TestNoThemeOverridesAllowed(t *testing.T) {
	files := `
-- hugo.toml --
theme = "mytheme"

-- themes/mytheme/hugo.toml --
[params]
themevar = "themeval"

[security]
allowContent = ['.*']
`

	b := hugolib.Test(t, files)

	conf := b.H.Conf

	_, themeOK := b.H.Site.Params()["themevar"]
	b.Assert(themeOK, qt.IsTrue)
	sec := conf.GetConfigSection("security").(security.Config)
	b.Assert(sec.AllowContent.Accept("text/html"), qt.IsFalse)
}
