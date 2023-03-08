// Copyright 2023 The Hugo Authors. All rights reserved.
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

package modules_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestHugoModFileV5(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[module]
[[module.imports]]
path = "github.com/bep/hugo-mod-with-hugodotmod/v5"
-- layouts/_default/index.html --
// This comes from the module imported.
{{ $foo := resources.Get "foo.txt" }}
Foo: {{ with $foo }}{{ .Content }}{{ end }}|
-- hugo.mod --
module github.com/gohugoio/testmod

go 1.20

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Verbose:     true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Foo: bar")

}

func TestHugoModWithGoModFileV5(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[module]
[[module.imports]]
path = "github.com/bep/hugo-mod-with-hugodotmod/v5"
-- layouts/_default/index.html --
// This comes from the module imported.
{{ $foo := resources.Get "foo.txt" }}
Foo: {{ with $foo }}{{ .Content }}{{ end }}|
-- go.mod --
module github.com/gohugoio/testmod

go 1.20

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Verbose:     true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Foo: bar")

}

func TestHugoModFileV6(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[module]
[[module.imports]]
path = "github.com/bep/hugo-mod-with-hugodotmod/v6"
-- layouts/_default/index.html --
// This comes from the module imported.
{{ $foo := resources.Get "foo.txt" }}
Foo: {{ with $foo }}{{ .Content }}{{ end }}|
-- hugo.mod --
module github.com/gohugoio/testmod

go 1.20

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Verbose:     true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Foo: baz")

}

// same as above but no /v6, just directly to the repo
func TestHugoModFileV6WithoutPackage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[module]
[[module.imports]]
path = "github.com/bep/hugo-mod-with-hugodotmod"
-- layouts/_default/index.html --
// This comes from the module imported.
{{ $foo := resources.Get "foo.txt" }}
Foo: {{ with $foo }}{{ .Content }}{{ end }}|
-- hugo.mod --
module github.com/gohugoio/testmod

go 1.20

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Verbose:     true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Foo: baz")

}
