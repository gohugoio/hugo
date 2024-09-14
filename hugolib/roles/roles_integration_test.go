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

package roles_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// TODO1 role hierarchy.
func TestRoles(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
[roles]
  [roles.guest]
  	default = true
	weight = 100
  [roles.member]
	weight = 200
-- content/memberonly.md --
---
title: "Member Only"
roles: ["member"]
---
Member content.
-- content/public.md --
---
title: "Public"
---
Users with no (blank) role will see this.
-- layouts/_default/single.html --
{{ .Title }}|{{ .Content }}|
`

	b := hugolib.Test(t, files)

	b.AssertPublishDir("memberonly")
}
