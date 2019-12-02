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

package partials

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCreateKey(t *testing.T) {
	c := qt.New(t)
	m := make(map[interface{}]bool)

	create := func(name string, variants ...interface{}) partialCacheKey {
		k, err := createKey(name, variants...)
		c.Assert(err, qt.IsNil)
		m[k] = true
		return k
	}

	for i := 0; i < 123; i++ {
		c.Assert(create("a", "b"), qt.Equals, partialCacheKey{name: "a", variant: "b"})
		c.Assert(create("a", "b", "c"), qt.Equals, partialCacheKey{name: "a", variant: "9629524865311698396"})
		c.Assert(create("a", 1), qt.Equals, partialCacheKey{name: "a", variant: 1})
		c.Assert(create("a", map[string]string{"a": "av"}), qt.Equals, partialCacheKey{name: "a", variant: "4809626101226749924"})
		c.Assert(create("a", []string{"a", "b"}), qt.Equals, partialCacheKey{name: "a", variant: "2712570657419664240"})
	}

}
