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

package config

import (
	"encoding/json"

	"github.com/gohugoio/hugo/common/hashing"
)

func DecodeNamespace[S, C any](configSource any, buildConfig func(any) (C, any, error)) (*ConfigNamespace[S, C], error) {
	// Calculate the hash of the input (not including any defaults applied later).
	// This allows us to introduce new config options without breaking the hash.
	h := hashing.HashString(configSource)

	// Build the config
	c, ext, err := buildConfig(configSource)
	if err != nil {
		return nil, err
	}

	if ext == nil {
		ext = configSource
	}

	if ext == nil {
		panic("ext is nil")
	}

	ns := &ConfigNamespace[S, C]{
		SourceStructure: ext,
		SourceHash:      h,
		Config:          c,
	}

	return ns, nil
}

// ConfigNamespace holds a Hugo configuration namespace.
// The construct looks a little odd, but it's built to make the configuration elements
// both self-documenting and contained in a common structure.
type ConfigNamespace[S, C any] struct {
	// SourceStructure represents the source configuration with any defaults applied.
	// This is used for documentation and printing of the configuration setup to the user.
	SourceStructure any

	// SourceHash is a hash of the source configuration before any defaults gets applied.
	SourceHash string

	// Config is the final configuration as used by Hugo.
	Config C
}

// MarshalJSON marshals the source structure.
func (ns *ConfigNamespace[S, C]) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.SourceStructure)
}

// Signature returns the signature of the source structure.
// Note that this is for documentation purposes only and SourceStructure may not always be cast to S (it's usually just a map).
func (ns *ConfigNamespace[S, C]) Signature() S {
	var s S
	return s
}
