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

package crypto_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// crypto.Hash combined with encoding.HexDecode and encoding.Base64Encode should
// reproduce the SRI hash in .Data.Integrity on a fingerprinted resource.
// See issue 15072.
func TestHashIntegrity(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
{{ $content := "hello world" }}
{{ range $algo := slice "sha256" "sha384" "sha512" }}
{{ $integrity := $content | resources.FromString "data.txt" | fingerprint $algo }}
{{ $composed := printf "%s-%s" $algo ($content | crypto.Hash $algo | encoding.HexDecode | encoding.Base64Encode) }}
{{ $algo }}: {{ eq $integrity.Data.Integrity $composed }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"sha256: true",
		"sha384: true",
		"sha512: true",
	)
}
