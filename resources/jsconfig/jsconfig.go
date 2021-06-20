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
	"sync"
)

// Builder builds a jsconfig.json file that, currently, is used only to assist
// intellinsense in editors.
type Builder struct {
	sourceRootsMu sync.RWMutex
	sourceRoots   map[string]string
}

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return &Builder{sourceRoots: make(map[string]string)}
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

	paths := make(map[string][]string)
	for sourceRoot, mountRoot := range b.sourceRoots {
		rel, err := filepath.Rel(dir, filepath.Join(sourceRoot, "*"))
		if err == nil {
			globPattern := filepath.Join(mountRoot, "*")
			paths[globPattern] = append(paths[globPattern], rel)
		}
	}

	conf.CompilerOptions.Paths = paths
	return conf
}

// AddRoots adds a new source root and mount root.
// This method is thread safe.
func (b *Builder) AddRoots(sourceRoot, mountRoot string) {
	b.sourceRootsMu.RLock()
	_, found := b.sourceRoots[sourceRoot]
	b.sourceRootsMu.RUnlock()

	if found {
		return
	}

	b.sourceRootsMu.Lock()
	b.sourceRoots[sourceRoot] = mountRoot
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
