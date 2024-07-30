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

// Package hashing provides common hashing utilities.
package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strconv"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/gohugoio/hashstructure"
	"github.com/gohugoio/hugo/identity"
)

// XXHashFromReader calculates the xxHash for the given reader.
func XXHashFromReader(r io.Reader) (uint64, int64, error) {
	h := getXxHashReadFrom()
	defer putXxHashReadFrom(h)

	size, err := io.Copy(h, r)
	if err != nil {
		return 0, 0, err
	}
	return h.Sum64(), size, nil
}

// XXHashFromString calculates the xxHash for the given string.
func XXHashFromString(s string) (uint64, error) {
	h := xxhash.New()
	h.WriteString(s)
	return h.Sum64(), nil
}

// XxHashFromStringHexEncoded calculates the xxHash for the given string
// and returns the hash as a hex encoded string.
func XxHashFromStringHexEncoded(f string) string {
	h := xxhash.New()
	h.WriteString(f)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

// MD5FromStringHexEncoded returns the MD5 hash of the given string.
func MD5FromStringHexEncoded(f string) string {
	h := md5.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum(nil))
}

// HashString returns a hash from the given elements.
// It will panic if the hash cannot be calculated.
// Note that this hash should be used primarily for identity, not for change detection as
// it in the more complex values (e.g. Page) will not hash the full content.
func HashString(vs ...any) string {
	hash := HashUint64(vs...)
	return strconv.FormatUint(hash, 10)
}

var hashOptsPool = sync.Pool{
	New: func() any {
		return &hashstructure.HashOptions{
			Hasher: xxhash.New(),
		}
	},
}

func getHashOpts() *hashstructure.HashOptions {
	return hashOptsPool.Get().(*hashstructure.HashOptions)
}

func putHashOpts(opts *hashstructure.HashOptions) {
	opts.Hasher.Reset()
	hashOptsPool.Put(opts)
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

	hashOpts := getHashOpts()
	defer putHashOpts(hashOpts)

	hash, err := hashstructure.Hash(o, hashOpts)
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
	case identity.IdentityProvider:
		return t.GetIdentity()
	default:
		return v
	}
}

type xxhashReadFrom struct {
	buff []byte
	*xxhash.Digest
}

func (x *xxhashReadFrom) ReadFrom(r io.Reader) (int64, error) {
	for {
		n, err := r.Read(x.buff)
		if n > 0 {
			x.Digest.Write(x.buff[:n])
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return int64(n), err
		}
	}
}

var xXhashReadFromPool = sync.Pool{
	New: func() any {
		return &xxhashReadFrom{Digest: xxhash.New(), buff: make([]byte, 48*1024)}
	},
}

func getXxHashReadFrom() *xxhashReadFrom {
	return xXhashReadFromPool.Get().(*xxhashReadFrom)
}

func putXxHashReadFrom(h *xxhashReadFrom) {
	h.Reset()
	xXhashReadFromPool.Put(h)
}
