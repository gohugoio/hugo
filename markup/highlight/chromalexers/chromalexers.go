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

package chromalexers

import (
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

type lexersMap struct {
	lexers map[string]chroma.Lexer
	mu     sync.RWMutex
}

var lexerCache = &lexersMap{lexers: make(map[string]chroma.Lexer)}

// Get returns a lexer for the given language name, nil if not found.
// This is just a wrapper around chromalexers.Get that caches the result.
// Reasoning for this is that chromalexers.Get is slow in the case where the lexer is not found,
// which is a common case in Hugo.
func Get(name string) chroma.Lexer {
	lexerCache.mu.RLock()
	lexer, found := lexerCache.lexers[name]
	lexerCache.mu.RUnlock()

	if found {
		return lexer
	}

	lexer = lexers.Get(name)

	lexerCache.mu.Lock()
	lexerCache.lexers[name] = lexer
	lexerCache.mu.Unlock()

	return lexer
}
