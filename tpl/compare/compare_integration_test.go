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

package compare_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/compare"
)

// See issue 15322.
func TestEqIntFloat(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
eq: {{ eq 1 1.0 }}|{{ eq 1.0 1 }}|{{ eq 1 1.5 }}|{{ eq 1 2 1.0 }}
ne: {{ ne 1 1.0 }}|{{ ne 1 1.5 }}
lt: {{ lt 1 1.0 }}|{{ lt 1 1.5 }}
le: {{ le 1 1.0 }}|{{ le 1 0.5 }}
gt: {{ gt 1 1.0 }}|{{ gt 1 0.5 }}
ge: {{ ge 1 1.0 }}|{{ ge 1 1.5 }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent(
		"public/index.html",
		"eq: true|true|false|true",
		"ne: false|true",
		"lt: false|true",
		"le: true|false",
		"gt: false|true",
		"ge: true|false",
	)
}

//

func BenchmarkEq(b *testing.B) {
	files := `
-- hugo.toml --
`
	bb := hugolib.Test(b, files)

	ns := bb.H.TemplateStore.GetTemplateFuncsNamespace("compare").(*compare.Namespace)

	b.ResetTimer()

	b.Run("eq int", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(42, 42)
		}
	})

	b.Run("eq float", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(42.0, 42.0)
		}
	})

	b.Run("eq string", func(b *testing.B) {
		for b.Loop() {
			ns.Eq("foo", "foo")
		}
	})

	b.Run("eq float32 and float64", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(float32(42.0), 42.0)
		}
	})

	b.Run("eq int and float", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(42, 42.0)
		}
	})

	b.Run("eq float and int", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(42.0, 42)
		}
	})

	b.Run("eq int and uint64", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(42, uint64(42))
		}
	})

	b.Run("eq uint32 and uint64", func(b *testing.B) {
		for b.Loop() {
			ns.Eq(uint32(42), uint64(42))
		}
	})
}
