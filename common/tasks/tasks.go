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

package tasks

import (
	"sync"
	"time"
)

// RunEvery runs a function at intervals defined by the function itself.
// Functions can be added and removed while running.
type RunEvery struct {
	// Any error returned from the function will be passed to this function.
	HandleError func(string, error)

	// If set, the function will be run immediately.
	RunImmediately bool

	// The named functions to run.
	funcs map[string]*Func

	mu      sync.Mutex
	started bool
	closed  bool
	quit    chan struct{}
}

type Func struct {
	// The shortest interval between each run.
	IntervalLow time.Duration

	// The longest interval between each run.
	IntervalHigh time.Duration

	// The function to run.
	F func(interval time.Duration) (time.Duration, error)

	interval time.Duration
	last     time.Time
}

func (r *RunEvery) Start() error {
	if r.started {
		return nil
	}

	r.started = true
	r.quit = make(chan struct{})

	go func() {
		if r.RunImmediately {
			r.run()
		}
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-r.quit:
				return
			case <-ticker.C:
				r.run()
			}
		}
	}()

	return nil
}

// Close stops the RunEvery from running.
func (r *RunEvery) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	if r.quit != nil {
		close(r.quit)
	}
	return nil
}

// Add adds a function to the RunEvery.
func (r *RunEvery) Add(name string, f Func) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.funcs == nil {
		r.funcs = make(map[string]*Func)
	}
	if f.IntervalLow == 0 {
		f.IntervalLow = 500 * time.Millisecond
	}
	if f.IntervalHigh <= f.IntervalLow {
		f.IntervalHigh = 20 * time.Second
	}

	start := f.IntervalHigh / 3
	if start < f.IntervalLow {
		start = f.IntervalLow
	}
	f.interval = start
	f.last = time.Now()

	r.funcs[name] = &f
}

// Remove removes a function from the RunEvery.
func (r *RunEvery) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.funcs, name)
}

// Has returns whether the RunEvery has a function with the given name.
func (r *RunEvery) Has(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, found := r.funcs[name]
	return found
}

func (r *RunEvery) run() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, f := range r.funcs {
		if time.Now().Before(f.last.Add(f.interval)) {
			continue
		}
		f.last = time.Now()
		interval, err := f.F(f.interval)
		if err != nil && r.HandleError != nil {
			r.HandleError(name, err)
		}

		if interval < f.IntervalLow {
			interval = f.IntervalLow
		}

		if interval > f.IntervalHigh {
			interval = f.IntervalHigh
		}
		f.interval = interval
	}
}
