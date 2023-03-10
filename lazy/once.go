// Copyright 2019 The Hugo Authors. All rights reserved.
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

package lazy

import (
	"sync"
	"sync/atomic"
)

// onceMore is similar to sync.Once.
//
// Additional features are:
// * it can be reset, so the action can be repeated if needed
// * it has methods to check if it's done or in progress

type onceMore struct {
	mu   sync.Mutex
	lock uint32
	done uint32
}

func (t *onceMore) Do(f func()) {
	if atomic.LoadUint32(&t.done) == 1 {
		return
	}

	// f may call this Do and we would get a deadlock.
	locked := atomic.CompareAndSwapUint32(&t.lock, 0, 1)
	if !locked {
		return
	}
	defer atomic.StoreUint32(&t.lock, 0)

	t.mu.Lock()
	defer t.mu.Unlock()

	// Double check
	if t.done == 1 {
		return
	}
	defer atomic.StoreUint32(&t.done, 1)
	f()
}

func (t *onceMore) InProgress() bool {
	return atomic.LoadUint32(&t.lock) == 1
}

func (t *onceMore) Done() bool {
	return atomic.LoadUint32(&t.done) == 1
}

func (t *onceMore) ResetWithLock() *sync.Mutex {
	t.mu.Lock()
	defer atomic.StoreUint32(&t.done, 0)
	return &t.mu
}
