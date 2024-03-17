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

package pagesfromdata

import "testing"

func BenchmarkHash(b *testing.B) {
	m := map[string]any{
		"foo":         "bar",
		"bar":         "foo",
		"stringSlice": []any{"a", "b", "c"},
		"intSlice":    []any{1, 2, 3},
		"largeText":   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit.",
	}

	bs := BuildState{}

	for i := 0; i < b.N; i++ {
		bs.hash(m)
	}
}
