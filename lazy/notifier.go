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
	"runtime/debug"
	"strings"
	"sync"
)

// Notifier as a synchronization tool that is queried for just-in-time access
// to a resource. Callers use Wait to block until the resource is ready, and
// call Close to indicate that the resource is ready. Reset returns the
// resource to its locked state.
//
// Notifier must be initialized by calling NewNotifier.
type Notifier struct {
	// Channel to close when the protected resource is ready. This must only
	// be accessed via calling the currentCh method to avoid race conditions.
	ch chan struct{}
	// For locking the channel while resetting it
	mu *sync.RWMutex
}

// NewNotifier creates a Notifier with all synchronization mechanisms
// initialized.
func NewNotifier() *Notifier {
	return &Notifier{
		ch: make(chan struct{}),
		mu: &sync.RWMutex{},
	}
}

// isClosed checks whether a channel is closed. If calling from a Notifier
// method, the calling goroutine must hold and release the Notifier.mu lock.
func isClosed(ch chan struct{}) bool {
	select {
	// Already closed
	case <-ch:
		return true
	default:
		return false
	}
}

// currentCh safely returns the current channel. Because this locks and unlocks
// the mutex, callers must not perform any other locking until the channel is
// returned.
func (n *Notifier) currentCh() chan struct{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.ch
}

// Wait waits for the Notifier to be ready, i.e., for Close to be called
// somewhere
func (n *Notifier) Wait() {
	ch := n.currentCh()
	<-ch
	s := string(debug.Stack())
	if strings.Contains(s, "newWatcher") {
	}
	return
}

// Close unblocks any goroutines that called Wait
func (n *Notifier) Close() {
	ch := n.currentCh()
	n.mu.Lock()
	defer n.mu.Unlock()
	if !isClosed(ch) {
		close(ch)
	}
	return
}

// Reset returns the resource to its pre-ready state while locking
func (n *Notifier) Reset() {
	ch := n.currentCh()
	n.mu.Lock()
	// No need to reset since the channel is open
	if !isClosed(ch) {
		return
	}
	defer n.mu.Unlock()
	n.ch = make(chan struct{})
	return
}
