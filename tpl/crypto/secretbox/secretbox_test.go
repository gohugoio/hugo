// Copyright 2020 The Hugo Authors. All rights reserved.
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

package secretbox

import (
	"encoding/hex"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// TestSecretBox is a smoketest adapted from the upstream Go implementation to
// confirm that our implementation matches the Go and C NaCl implementations.
func TestSecretBox(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	var (
		key     [32]byte
		nonce   [24]byte
		message [64]byte
	)
	for i := range key[:] {
		key[i] = 1
	}
	for i := range nonce[:] {
		nonce[i] = 2
	}
	for i := range message[:] {
		message[i] = 3
	}

	box, err := ns.Seal(key[:], message[:], nonce[:])
	c.Assert(err, qt.IsNil)

	// expected was generated using the C implementation of NaCl.
	expected, _ := hex.DecodeString("8442bc313f4626f1359e3b50122b6ce6fe66ddfe7d39d14e637eb4fd5b45beadab55198df6ab5368439792a23c87db70acb6156dc5ef957ac04f6276cf6093b84be77ff0849cc33e34b7254d5a8f65ad")

	c.Check(box[len(nonce):], qt.Equals, string(expected))
}

func TestSecretBoxSealOpen(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for i, test := range []struct {
		key     interface{}
		message interface{}
		isErr   bool // use hex encoding
	}{
		{strings.Repeat("1", 32), strings.Repeat("1234567890", 4), false},
		{[]byte{154, 211, 108, 22, 133, 234, 137, 218, 38, 238, 17, 2, 161, 77, 180, 178}, strings.Repeat("0987654321", 5), false},
		// Errors
		{"foo", t, true},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.key)

		result, err := ns.Seal(test.key, test.message)

		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)

		dec, err := ns.Open(test.key, result)
		c.Assert(err, qt.IsNil, errMsg)
		c.Check(dec, qt.Equals, test.message, errMsg)
	}
}
