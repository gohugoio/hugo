// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hsync

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestOnceMoreValue(t *testing.T) {
	c := qt.New(t)

	var counter int
	f := func(context.Context) int {
		counter++
		return counter
	}

	omf := OnceMoreValue(f)
	for range 10 {
		c.Assert(omf.Value(context.Background()), qt.Equals, 1)
	}
	omf.Reset()
	for range 10 {
		c.Assert(omf.Value(context.Background()), qt.Equals, 2)
	}
}

func TestOnceMoreFunc(t *testing.T) {
	c := qt.New(t)

	var counter int
	f := func(context.Context) error {
		counter++
		return nil
	}

	omf := OnceMoreFunc(f)
	for range 10 {
		c.Assert(omf.Do(context.Background()), qt.IsNil)
		c.Assert(counter, qt.Equals, 1)
	}
	omf.Reset()
	for range 10 {
		c.Assert(omf.Do(context.Background()), qt.IsNil)
		c.Assert(counter, qt.Equals, 2)
	}
}

func BenchmarkOnceMoreValue(b *testing.B) {
	var counter int
	f := func(context.Context) int {
		counter++
		return counter
	}

	for b.Loop() {
		omf := OnceMoreValue(f)
		for range 10 {
			omf.Value(context.Background())
		}
		omf.Reset()
		for range 10 {
			omf.Value(context.Background())
		}
	}
}
