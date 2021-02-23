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

package internal

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

type testStruct struct {
	Name string
	V1   int64
	V2   int32
	V3   int
	V4   uint64
}

func TestResourceTransformationKey(t *testing.T) {
	// We really need this key to be portable across OSes.
	key := NewResourceTransformationKey("testing",
		testStruct{Name: "test", V1: int64(10), V2: int32(20), V3: 30, V4: uint64(40)})
	c := qt.New(t)
	c.Assert(key.Value(), qt.Equals, "testing_518996646957295636")
}
