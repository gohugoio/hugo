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

package paths_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestRemovePathAccents(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
removePathAccents = true
-- content/διακριτικός.md --
-- content/διακριτικός.fr.md --
-- layouts/_default/single.html --
{{ .Language.Lang }}|Single.
-- layouts/_default/list.html --
List
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/διακριτικός/index.html", "en|Single")
	b.AssertFileContent("public/fr/διακριτικος/index.html", "fr|Single")
}

func TestDisablePathToLower(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
disablePathToLower = true
-- content/MySection/MyPage.md --
-- content/MySection/MyPage.fr.md --
-- content/MySection/MyBundle/index.md --
-- content/MySection/MyBundle/index.fr.md --
-- layouts/_default/single.html --
{{ .Language.Lang }}|Single.
-- layouts/_default/list.html --
{{ .Language.Lang }}|List.
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/mysection/index.html", "en|List")
	b.AssertFileContent("public/en/mysection/mypage/index.html", "en|Single")
	b.AssertFileContent("public/fr/MySection/index.html", "fr|List")
	b.AssertFileContent("public/fr/MySection/MyPage/index.html", "fr|Single")
	b.AssertFileContent("public/en/mysection/mybundle/index.html", "en|Single")
	b.AssertFileContent("public/fr/MySection/MyBundle/index.html", "fr|Single")
}
