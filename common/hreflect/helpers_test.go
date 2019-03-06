// Copyright 2019 The Hugo Authors. All rights reserved.
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

package hreflect

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsTruthful(t *testing.T) {
	assert := require.New(t)

	assert.True(IsTruthful(true))
	assert.False(IsTruthful(false))
	assert.True(IsTruthful(time.Now()))
	assert.False(IsTruthful(time.Time{}))
}

func BenchmarkIsTruthFul(b *testing.B) {
	v := reflect.ValueOf("Hugo")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !IsTruthfulValue(v) {
			b.Fatal("not truthful")
		}
	}
}
