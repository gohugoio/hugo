// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package crypto provides template functions for cryptographic operations.
package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/spf13/cast"
)

// New returns a new instance of the crypto-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "crypto" namespace.
type Namespace struct{}

// MD5 hashes the given input and returns its MD5 checksum.
func (ns *Namespace) MD5(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA1 hashes the given input and returns its SHA1 checksum.
func (ns *Namespace) SHA1(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA256 hashes the given input and returns its SHA256 checksum.
func (ns *Namespace) SHA256(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// HMAC returns a cryptographic hash that uses a key to sign a message.
func (ns *Namespace) HMAC(h interface{}, k interface{}, m interface{}) (string, error) {
	ha, err := cast.ToStringE(h)
	if err != nil {
		return "", err
	}

	var hash func() hash.Hash
	switch ha {
	case "md5":
		hash = md5.New
	case "sha1":
		hash = sha1.New
	case "sha256":
		hash = sha256.New
	case "sha512":
		hash = sha512.New
	default:
		return "", fmt.Errorf("hmac: %s is not a supported hash function", ha)
	}

	msg, err := cast.ToStringE(m)
	if err != nil {
		return "", err
	}

	key, err := cast.ToStringE(k)
	if err != nil {
		return "", err
	}

	mac := hmac.New(hash, []byte(key))
	_, err = mac.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)[:]), nil
}
