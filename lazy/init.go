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
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// New creates a new empty Init.
func New() *Init {
	return &Init{}
}

// Init holds a graph of lazily initialized dependencies.
type Init struct {
	// Used mainly for testing.
	initCount uint64

	mu sync.Mutex

	prev     *Init
	children []*Init

	init onceMore
	out  any
	err  error
	f    func(context.Context) (any, error)
}

// Add adds a func as a new child dependency.
func (ini *Init) Add(initFn func(context.Context) (any, error)) *Init {
	if ini == nil {
		ini = New()
	}
	return ini.add(false, initFn)
}

// InitCount gets the number of this this Init has been initialized.
func (ini *Init) InitCount() int {
	i := atomic.LoadUint64(&ini.initCount)
	return int(i)
}

// AddWithTimeout is same as Add, but with a timeout that aborts initialization.
func (ini *Init) AddWithTimeout(timeout time.Duration, f func(ctx context.Context) (any, error)) *Init {
	return ini.Add(func(ctx context.Context) (any, error) {
		return ini.withTimeout(ctx, timeout, f)
	})
}

// Branch creates a new dependency branch based on an existing and adds
// the given dependency as a child.
func (ini *Init) Branch(initFn func(context.Context) (any, error)) *Init {
	if ini == nil {
		ini = New()
	}
	return ini.add(true, initFn)
}

// BranchWithTimeout is same as Branch, but with a timeout.
func (ini *Init) BranchWithTimeout(timeout time.Duration, f func(ctx context.Context) (any, error)) *Init {
	return ini.Branch(func(ctx context.Context) (any, error) {
		return ini.withTimeout(ctx, timeout, f)
	})
}

// Do initializes the entire dependency graph.
func (ini *Init) Do(ctx context.Context) (any, error) {
	if ini == nil {
		panic("init is nil")
	}

	ini.init.Do(func() {
		atomic.AddUint64(&ini.initCount, 1)
		prev := ini.prev
		if prev != nil {
			// A branch. Initialize the ancestors.
			if prev.shouldInitialize() {
				_, err := prev.Do(ctx)
				if err != nil {
					ini.err = err
					return
				}
			} else if prev.inProgress() {
				// Concurrent initialization. The following init func
				// may depend on earlier state, so wait.
				prev.wait()
			}
		}

		if ini.f != nil {
			ini.out, ini.err = ini.f(ctx)
		}

		for _, child := range ini.children {
			if child.shouldInitialize() {
				_, err := child.Do(ctx)
				if err != nil {
					ini.err = err
					return
				}
			}
		}
	})

	ini.wait()

	return ini.out, ini.err
}

// TODO(bep) investigate if we can use sync.Cond for this.
func (ini *Init) wait() {
	var counter time.Duration
	for !ini.init.Done() {
		counter += 10
		if counter > 600000000 {
			panic("BUG: timed out in lazy init")
		}
		time.Sleep(counter * time.Microsecond)
	}
}

func (ini *Init) inProgress() bool {
	return ini != nil && ini.init.InProgress()
}

func (ini *Init) shouldInitialize() bool {
	return !(ini == nil || ini.init.Done() || ini.init.InProgress())
}

// Reset resets the current and all its dependencies.
func (ini *Init) Reset() {
	mu := ini.init.ResetWithLock()
	ini.err = nil
	defer mu.Unlock()
	for _, d := range ini.children {
		d.Reset()
	}
}

func (ini *Init) add(branch bool, initFn func(context.Context) (any, error)) *Init {
	ini.mu.Lock()
	defer ini.mu.Unlock()

	if branch {
		return &Init{
			f:    initFn,
			prev: ini,
		}
	}

	ini.checkDone()
	ini.children = append(ini.children, &Init{
		f: initFn,
	})

	return ini
}

func (ini *Init) checkDone() {
	if ini.init.Done() {
		panic("init cannot be added to after it has run")
	}
}

func (ini *Init) withTimeout(ctx context.Context, timeout time.Duration, f func(ctx context.Context) (any, error)) (any, error) {
	// Create a new context with a timeout not connected to the incoming context.
	waitCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := make(chan verr, 1)

	go func() {
		v, err := f(ctx)
		select {
		case <-waitCtx.Done():
			return
		default:
			c <- verr{v: v, err: err}
		}
	}()

	select {
	case <-waitCtx.Done():
		//lint:ignore ST1005 end user message.
		return nil, errors.New("timed out initializing value. You may have a circular loop in a shortcode, or your site may have resources that take longer to build than the `timeout` limit in your Hugo config file.")
	case ve := <-c:
		return ve.v, ve.err
	}
}

type verr struct {
	v   any
	err error
}
