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

package identity

import (
	"strconv"

	"github.com/mitchellh/hashstructure/v2"
)

// HashString returns a hash from the given elements.
// It will panic if the hash cannot be calculated.
// Note that this hash should be used primarily for identity, not for change detection as
// it in the more complex values (e.g. Page) will not hash the full content.
func HashString(vs ...any) string {
	hash := HashUint64(vs...)
	return strconv.FormatUint(hash, 10)
}

// HashUint64 returns a hash from the given elements.
// It will panic if the hash cannot be calculated.
// Note that this hash should be used primarily for identity, not for change detection as
// it in the more complex values (e.g. Page) will not hash the full content.
func HashUint64(vs ...any) uint64 {
	var o any
	if len(vs) == 1 {
		o = toHashable(vs[0])
	} else {
		elements := make([]any, len(vs))
		for i, e := range vs {
			elements[i] = toHashable(e)
		}
		o = elements
	}

	hash, err := hashstructure.Hash(o, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}

type keyer interface {
	Key() string
}

// For structs, hashstructure.Hash only works on the exported fields,
// so rewrite the input slice for known identity types.
func toHashable(v any) any {
	switch t := v.(type) {
	case keyer:
		return t.Key()
	case IdentityProvider:
		return t.GetIdentity()
	default:
		return v
	}
}
