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
	"encoding/hex"
	"io"
	"sync"

	"github.com/cespare/xxhash/v2"
)

// XXHashFromReader calculates the xxHash for the given reader.
func XXHashFromReader(r io.ReadSeeker) (uint64, int64, error) {
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
