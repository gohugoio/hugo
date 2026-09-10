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

func BenchmarkCompare(b *testing.B) {
	files := `
-- hugo.toml --
`
	bb := hugolib.Test(b, files)

	ns := bb.H.TemplateStore.GetTemplateFuncsNamespace("compare").(*compare.Namespace)

	b.ResetTimer()

	run := func(name string, fn func(v1, v2 any)) {
		b.Run(name, func(b *testing.B) {
			b.Run("mixed int and float", func(b *testing.B) {
				for b.Loop() {
					fn(42, 42.0)
				}
			})
			b.Run("mixed float and int", func(b *testing.B) {
				for b.Loop() {
					fn(42.0, 42)
				}
			})
			b.Run("mixed string and int", func(b *testing.B) {
				for b.Loop() {
					fn("42", 42)
				}
			})
			b.Run("mixed string and float", func(b *testing.B) {
				for b.Loop() {
					fn("42", 42.0)
				}
			})
			b.Run("mixed float and string", func(b *testing.B) {
				for b.Loop() {
					fn(42.0, "42")
				}
			})
			// End of mixed type benchmark
			// start of same type benchmark
			b.Run("same int", func(b *testing.B) {
				for b.Loop() {
					fn(42, 42)
				}
			})
			b.Run("same float", func(b *testing.B) {
				for b.Loop() {
					fn(42.0, 42.0)
				}
			})
			b.Run("same string", func(b *testing.B) {
				for b.Loop() {
					fn("foo", "foo")
				}
			})
			b.Run("same bool", func(b *testing.B) {
				for b.Loop() {
					fn(true, true)
				}
			})
			b.Run("same nil", func(b *testing.B) {
				for b.Loop() {
					fn(nil, nil)
				}
			})
			b.Run("same complex", func(b *testing.B) {
				for b.Loop() {
					fn(complex(1, 2), complex(1, 2))
				}
			})
			// End of same type benchmark
			// start of edge case benchmark
			b.Run("empty slice", func(b *testing.B) {
				for b.Loop() {
					fn([]int{}, []int{})
				}
			})
			b.Run("empty map", func(b *testing.B) {
				for b.Loop() {
					fn(map[string]int{}, map[string]int{})
				}
			})
			// End of edge case benchmark
			// start of nil benchmark
			b.Run("nil and non-nil", func(b *testing.B) {
				for b.Loop() {
					fn(nil, 42)
				}
			})
			b.Run("non-nil and nil", func(b *testing.B) {
				for b.Loop() {
					fn(42, nil)
				}
			})
			// End of nil benchmark
			// start of large number benchmark
			b.Run("large numbers", func(b *testing.B) {
				for b.Loop() {
					fn(1e18, 1e18)
				}
			})
			// End of large number benchmark
			// start of small number benchmark
			b.Run("small numbers", func(b *testing.B) {
				for b.Loop() {
					fn(1e-18, 1e-18)
				}
			})
			// End of small number benchmark
		})
	}

	run("Eq", func(v1, v2 any) {
		ns.Eq(v1, v2)
	})

	run("Gt", func(v1, v2 any) {
		ns.Gt(v1, v2)
	})
}
