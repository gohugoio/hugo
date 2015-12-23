// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"bytes"
	"os"
	"testing"
)

func BenchmarkParsePage(b *testing.B) {
	f, _ := os.Open("redis.cn.md")
	sample := new(bytes.Buffer)
	sample.ReadFrom(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		page, _ := NewPage("bench")
		page.ReadFrom(bytes.NewReader(sample.Bytes()))
	}
}
