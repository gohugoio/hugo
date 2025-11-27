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

package hugolib

import (
	"strings"
	"testing"
)

func TestLegacyConfigDotToml(t *testing.T) {
	const filesTemplate = `
-- config.toml --
title = "My Site"
-- layouts/home.html --
Site: {{ .Site.Title }}

  `
	t.Run("In root", func(t *testing.T) {
		t.Parallel()
		files := filesTemplate
		b := Test(t, files)
		b.AssertFileContent("public/index.html", "Site: My Site")
	})

	t.Run("In config dir", func(t *testing.T) {
		t.Parallel()
		files := strings.Replace(filesTemplate, "-- config.toml --", "-- config/_default/config.toml --", 1)
		b := Test(t, files)
		b.AssertFileContent("public/index.html", "Site: My Site")
	})
}
