// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hsync

import (
	"context"
	"sync"
	"sync/atomic"
)

// OnceMore is similar to sync.Once.
//
// Additional features are:
// * it can be reset, so the action can be repeated if needed
// * it has methods to check if it's done or in progress
type OnceMore struct {
	_    doNotCopy
	done atomic.Bool
	mu   sync.Mutex
}

func (t *OnceMore) Do(f func()) {
	if t.Done() {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Double check
	if t.Done() {
		return
	}

	defer t.done.Store(true)
	f()
}

func (t *OnceMore) Done() bool {
	return t.done.Load()
}

func (t *OnceMore) Reset() {
	t.mu.Lock()
	t.done.Store(false)
	t.mu.Unlock()
}

type ValueResetter[T any] struct {
	reset func()
	f     func(context.Context) T
}

func (v *ValueResetter[T]) Value(ctx context.Context) T {
	return v.f(ctx)
}

func (v *ValueResetter[T]) Reset() {
	v.reset()
}

// OnceMoreValue returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func OnceMoreValue[T any](f func(context.Context) T) ValueResetter[T] {
	v := struct {
		f      func(context.Context) T
		once   OnceMore
		ok     bool
		p      any
		result T
	}{
		f: f,
	}
	ff := func(ctx context.Context) T {
		v.once.Do(func() {
			v.ok = false
			defer func() {
				v.p = recover()
				if !v.ok {
					panic(v.p)
				}
			}()
			v.result = v.f(ctx)
			v.ok = true
		})
		if !v.ok {
			panic(v.p)
		}
		return v.result
	}

	return ValueResetter[T]{
		reset: v.once.Reset,
		f:     ff,
	}
}

type FuncResetter struct {
	f     func(context.Context) error
	reset func()
}

func (v *FuncResetter) Do(ctx context.Context) error {
	return v.f(ctx)
}

func (v *FuncResetter) Reset() {
	v.reset()
}

func OnceMoreFunc(f func(context.Context) error) FuncResetter {
	v := struct {
		f    func(context.Context) error
		once OnceMore
		ok   bool
		err  error
		p    any
	}{
		f: f,
	}
	ff := func(ctx context.Context) error {
		v.once.Do(func() {
			v.ok = false
			defer func() {
				v.p = recover()
				if !v.ok {
					panic(v.p)
				}
			}()
			v.err = v.f(ctx)
			v.ok = true
		})
		if !v.ok {
			panic(v.p)
		}
		return v.err
	}

	return FuncResetter{
		f:     ff,
		reset: v.once.Reset,
	}
}

type doNotCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*doNotCopy) Lock()   {}
func (*doNotCopy) Unlock() {}
