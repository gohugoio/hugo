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

package jsconfig

import (
	"path/filepath"
	"sort"
	"sync"
)

// Builder builds a jsconfig.json file that, currently, is used only to assist
// IntelliSense in editors.
type Builder struct {
	sourceRootsMu sync.RWMutex
	sourceRoots   map[string]bool
}

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return &Builder{sourceRoots: make(map[string]bool)}
}

// Build builds a new Config with paths relative to dir.
// This method is thread safe.
func (b *Builder) Build(dir string) *Config {
	b.sourceRootsMu.RLock()
	defer b.sourceRootsMu.RUnlock()

	if len(b.sourceRoots) == 0 {
		return nil
	}
	conf := newJSConfig()

	var roots []string
	for root := range b.sourceRoots {
		rel, err := filepath.Rel(dir, filepath.Join(root, "*"))
		if err == nil {
			roots = append(roots, rel)
		}
	}
	sort.Strings(roots)
	conf.CompilerOptions.Paths["*"] = roots

	return conf
}

// AddSourceRoot adds a new source root.
// This method is thread safe.
func (b *Builder) AddSourceRoot(root string) {
	b.sourceRootsMu.RLock()
	found := b.sourceRoots[root]
	b.sourceRootsMu.RUnlock()

	if found {
		return
	}

	b.sourceRootsMu.Lock()
	b.sourceRoots[root] = true
	b.sourceRootsMu.Unlock()
}

// CompilerOptions holds compilerOptions for jsonconfig.json.
type CompilerOptions struct {
	BaseURL string              `json:"baseUrl"`
	Paths   map[string][]string `json:"paths"`
}

// Config holds the data for jsconfig.json.
type Config struct {
	CompilerOptions CompilerOptions `json:"compilerOptions"`
}

func newJSConfig() *Config {
	return &Config{
		CompilerOptions: CompilerOptions{
			BaseURL: ".",
			Paths:   make(map[string][]string),
		},
	}
}
