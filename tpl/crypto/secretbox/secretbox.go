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

// Package secretbox provides template functions for NaCl secretbox operations.
package secretbox

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/spf13/cast"
	"golang.org/x/crypto/nacl/secretbox"
)

// New returns a new instance of the secretbox-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "secretbox" namespace.
type Namespace struct{}

// Seal encrypts small messages using the NaCl secretbox framework.  The key and
// message parameters are required.  The 24-byte nonce is prepended to the
// resulting output.
//
// In normal use, the optional nonce parameter should not be used so that a
// random 24-byte nonce will be generated for each encryption.
func (ns *Namespace) Seal(key, message interface{}, nonce ...interface{}) (string, error) {
	// Key
	keyString, err := cast.ToStringE(key)
	if err != nil {
		return "", err
	}

	var secretKey [32]byte
	copy(secretKey[:], keyString)

	// Message
	messageString, err := cast.ToStringE(message)
	if err != nil {
		return "", err
	}

	// Nonce - must be unique for each call
	var nonceBytes [24]byte
	if len(nonce) == 0 {
		if _, err = io.ReadFull(rand.Reader, nonceBytes[:]); err != nil {
			return "", err
		}
	} else {
		// Message
		var nonceString string

		nonceString, err = cast.ToStringE(nonce[0])
		if err != nil {
			return "", err
		}

		copy(nonceBytes[:], nonceString)
	}

	res, err := secretbox.Seal(nonceBytes[:], []byte(messageString), &nonceBytes, &secretKey), nil
	return string(res), err
}

// Open decrypts and authenticates small messages using the NaCl secretbox
// framework. This implementation expects to find the 24-byte nonce at the
// beginning of the box content.
func (ns *Namespace) Open(key, box interface{}) (string, error) {
	// Key
	keyString, err := cast.ToStringE(key)
	if err != nil {
		return "", err
	}

	var secretKey [32]byte
	copy(secretKey[:], keyString)

	// Box
	boxString, err := cast.ToStringE(box)
	if err != nil {
		return "", err
	}

	if len(boxString) < 32 {
		return "", errors.New("invalid box parameter to secretbox.Open")
	}

	// Nonce - prepended to the encrypted payload
	var nonce [24]byte
	copy(nonce[:], boxString[:24])

	decrypted, ok := secretbox.Open(nil, []byte(boxString[24:]), &nonce, &secretKey)
	if !ok {
		return "", nil
	}

	return string(decrypted), nil
}
