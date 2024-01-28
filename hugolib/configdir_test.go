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

package hugolib

import "testing"

func TestConfigDir(t *testing.T) {
	t.Parallel()

	files := `
-- config/_default/params.toml --
a = "acp1"
d = "dcp1"
-- config/_default/config.toml --
[params]
a = "ac1"
b = "bc1"

-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
ignoreErrors = ["error-missing-instagram-accesstoken"]
[params]
a = "a1"
b = "b1"
c = "c1"
-- layouts/index.html --
Params: {{ site.Params}}
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Params: map[a:acp1 b:bc1 c:c1 d:dcp1]


`)
}
