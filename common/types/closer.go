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

package types

import "sync"

type Closer interface {
	Close() error
}

type CloseAdder interface {
	Add(Closer)
}

type Closers struct {
	mu sync.Mutex
	cs []Closer
}

func (cs *Closers) Add(c Closer) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.cs = append(cs.cs, c)
}

func (cs *Closers) Close() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.cs {
		c.Close()
	}

	cs.cs = cs.cs[:0]

	return nil
}
