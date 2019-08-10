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

package compare

import (
	"reflect"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hreflect"
)

func TestTruth(t *testing.T) {
	n := New(false)

	truthv, falsev := reflect.ValueOf(time.Now()), reflect.ValueOf(false)

	assertTruth := func(t *testing.T, v reflect.Value, expected bool) {
		if hreflect.IsTruthfulValue(v) != expected {
			t.Fatal("truth mismatch")
		}
	}

	t.Run("And", func(t *testing.T) {
		assertTruth(t, n.And(truthv, truthv), true)
		assertTruth(t, n.And(truthv, falsev), false)

	})

	t.Run("Or", func(t *testing.T) {
		assertTruth(t, n.Or(truthv, truthv), true)
		assertTruth(t, n.Or(falsev, truthv, falsev), true)
		assertTruth(t, n.Or(falsev, falsev), false)
	})

	t.Run("Not", func(t *testing.T) {
		c := qt.New(t)
		c.Assert(n.Not(falsev), qt.Equals, true)
		c.Assert(n.Not(truthv), qt.Equals, false)
	})

	t.Run("getIf", func(t *testing.T) {
		c := qt.New(t)
		assertTruth(t, n.getIf(reflect.ValueOf(nil)), false)
		s := reflect.ValueOf("Hugo")
		c.Assert(n.getIf(s), qt.Equals, s)
	})
}
