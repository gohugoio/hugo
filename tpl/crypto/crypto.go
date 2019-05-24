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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

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
