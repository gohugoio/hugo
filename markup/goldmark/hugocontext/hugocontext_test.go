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

package hugocontext

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestWrap(t *testing.T) {
	c := qt.New(t)

	b := []byte("test")

	c.Assert(Wrap(b, 42), qt.Equals, "{{__hugo_ctx pid=42}}\ntest{{__hugo_ctx/}}\n")
}

func BenchmarkWrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Wrap([]byte("test"), 42)
	}
}
