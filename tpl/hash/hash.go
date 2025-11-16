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

// Package hash provides non-cryptographic hash functions for template use.
package hash

import (
	"context"
	"hash/fnv"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/spf13/cast"
)

// New returns a new instance of the hash-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "hash" namespace.
type Namespace struct{}

// FNV32a hashes v using fnv32a algorithm.
func (ns *Namespace) FNV32a(v any) (int, error) {
	conv, err := cast.ToStringE(v)
	if err != nil {
		return 0, err
	}
	algorithm := fnv.New32a()
	algorithm.Write([]byte(conv))
	return int(algorithm.Sum32()), nil
}

// XxHash returns the xxHash of the input string.
func (ns *Namespace) XxHash(v any) (string, error) {
	conv, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	return hashing.XxHashFromStringHexEncoded(conv), nil
}

const name = "hash"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.XxHash,
			[]string{"xxhash"},
			[][2]string{
				{`{{ hash.XxHash "The quick brown fox jumps over the lazy dog" }}`, `0b242d361fda71bc`},
			},
		)

		ns.AddMethodMapping(ctx.FNV32a,
			nil,
			[][2]string{
				{`{{ hash.FNV32a "Hugo Rocks!!" }}`, `1515779328`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
