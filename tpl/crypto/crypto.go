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

// MD5 hashes the v and returns its MD5 checksum.
func (ns *Namespace) MD5(v any) (string, error) {
	conv, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA1 hashes v and returns its SHA1 checksum.
func (ns *Namespace) SHA1(v any) (string, error) {
	conv, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA256 hashes v and returns its SHA256 checksum.
func (ns *Namespace) SHA256(v any) (string, error) {
	conv, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// Hash returns the hex-encoded checksum of v using the given algorithm; one of
// md5, sha1, sha256 (the default), sha384 or sha512.
//
// The supported algorithms match those used for the Subresource Integrity (SRI)
// hash in .Data.Integrity on fingerprinted resources, so an SRI hash can be
// constructed by combining this with encoding.HexDecode and encoding.Base64Encode.
func (ns *Namespace) Hash(args ...any) (string, error) {
	var algo, v any
	switch len(args) {
	case 1:
		algo, v = "sha256", args[0]
	case 2:
		algo, v = args[0], args[1]
	default:
		return "", fmt.Errorf("crypto.Hash: expected 1 or 2 arguments, got %d", len(args))
	}

	conv, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	algoS, err := cast.ToStringE(algo)
	if err != nil {
		return "", err
	}

	h, err := newHash(algoS)
	if err != nil {
		return "", err
	}

	if _, err := h.Write([]byte(conv)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func newHash(algo string) (hash.Hash, error) {
	switch algo {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha384":
		return sha512.New384(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("crypto.Hash: %q is not a supported hash algorithm", algo)
	}
}

// HMAC returns a cryptographic hash that uses a key to sign a message.
func (ns *Namespace) HMAC(h any, k any, m any, e ...any) (string, error) {
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

	encoding := "hex"
	if len(e) > 0 && e[0] != nil {
		encoding, err = cast.ToStringE(e[0])
		if err != nil {
			return "", err
		}
	}

	switch encoding {
	case "binary":
		return string(mac.Sum(nil)[:]), nil
	case "hex":
		return hex.EncodeToString(mac.Sum(nil)[:]), nil
	default:
		return "", fmt.Errorf("%q is not a supported encoding method", encoding)
	}
}
