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

package integrity

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHashFromAlgo(t *testing.T) {

	for _, algo := range []struct {
		name string
		bits int
	}{
		{"md5", 128},
		{"sha256", 256},
		{"sha384", 384},
		{"sha512", 512},
		{"shaman", -1},
	} {

		t.Run(algo.name, func(t *testing.T) {
			c := qt.New(t)
			h, err := newHash(algo.name)
			if algo.bits > 0 {
				c.Assert(err, qt.IsNil)
				c.Assert(h.Size(), qt.Equals, algo.bits/8)
			} else {
				c.Assert(err, qt.Not(qt.IsNil))
				c.Assert(err.Error(), qt.Contains, "use either md5, sha256, sha384 or sha512")
			}

		})
	}
}
