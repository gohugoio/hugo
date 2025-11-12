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

package doctree

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

var _ MutableTrees = MutableTrees{}

const (
	LockTypeNone LockType = iota
	LockTypeRead
	LockTypeWrite
)

// AddEventListener adds an event listener to the tree.
// Note that the handler func may not add listeners.
func (ctx *WalkContext[T]) AddEventListener(event, path string, handler func(*Event[T])) {
	ctx.eventHandlersMu.Lock()
	defer ctx.eventHandlersMu.Unlock()

	if ctx.eventHandlers == nil {
		ctx.eventHandlers = make(eventHandlers[T])
	}
	if ctx.eventHandlers[event] == nil {
		ctx.eventHandlers[event] = make([]func(*Event[T]), 0)
	}

	// We want to match all above the path, so we need to exclude any similar named siblings.
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	ctx.eventHandlers[event] = append(
		ctx.eventHandlers[event], func(e *Event[T]) {
			// Propagate events up the tree only.
			if e.Path != path && strings.HasPrefix(e.Path, path) {
				handler(e)
			}
		},
	)
}

func (ctx *WalkContext[T]) Data() *SimpleThreadSafeTree[any] {
	ctx.dataInit.Do(func() {
		ctx.data = NewSimpleThreadSafeTree[any]()
	})
	return ctx.data
}

func (ctx *WalkContext[T]) initDataRaw() {
	ctx.dataRawInit.Do(func() {
		ctx.dataRaw = maps.NewCache[sitesmatrix.Vector, *SimpleThreadSafeTree[any]]()
	})
}

func (ctx *WalkContext[T]) DataRaw(vec sitesmatrix.Vector) *SimpleThreadSafeTree[any] {
	ctx.initDataRaw()
	v, _ := ctx.dataRaw.GetOrCreate(vec, func() (*SimpleThreadSafeTree[any], error) {
		return NewSimpleThreadSafeTree[any](), nil
	})
	return v
}

func (ctx *WalkContext[T]) DataRawForEeach() iter.Seq2[sitesmatrix.Vector, *SimpleThreadSafeTree[any]] {
	ctx.initDataRaw()
	return func(yield func(vec sitesmatrix.Vector, data *SimpleThreadSafeTree[any]) bool) {
		ctx.dataRaw.ForEeach(func(vec sitesmatrix.Vector, data *SimpleThreadSafeTree[any]) bool {
			return yield(vec, data)
		})
	}
}

// SendEvent sends an event up the tree.
func (ctx *WalkContext[T]) SendEvent(event *Event[T]) {
	ctx.eventMu.Lock()
	defer ctx.eventMu.Unlock()
	ctx.events = append(ctx.events, event)
}

// StopPropagation stops the propagation of the event.
func (e *Event[T]) StopPropagation() {
	e.stopPropagation = true
}

// ValidateKey returns an error if the key is not valid.
func ValidateKey(key string) error {
	if key == "" {
		// Root node.
		return nil
	}

	if len(key) < 2 {
		return fmt.Errorf("too short key: %q", key)
	}

	if key[0] != '/' {
		return fmt.Errorf("key must start with '/': %q", key)
	}

	if key[len(key)-1] == '/' {
		return fmt.Errorf("key must not end with '/': %q", key)
	}

	return nil
}

// Event is used to communicate events in the tree.
type Event[T any] struct {
	Name            string
	Path            string
	Source          T
	stopPropagation bool
}

type LockType int

// MutableTree is a tree that can be modified.
type MutableTree interface {
	DeleteRaw(key string)
	DeletePrefix(prefix string) int
	DeletePrefixRaw(prefix string) int
	Lock(writable bool) (commit func())
	CanLock() bool // Used for troubleshooting only.
}

// WalkableTree is a tree that can be walked.
type WalkableTree[T any] interface {
	WalkPrefixRaw(prefix string, walker func(key string, value T) bool)
}

var _ WalkableTree[any] = (*WalkableTrees[any])(nil)

type WalkableTrees[T any] []WalkableTree[T]

func (t WalkableTrees[T]) WalkPrefixRaw(prefix string, walker func(key string, value T) bool) {
	for _, tree := range t {
		tree.WalkPrefixRaw(prefix, walker)
	}
}

var _ MutableTree = MutableTrees(nil)

type MutableTrees []MutableTree

func (t MutableTrees) DeleteRaw(key string) {
	for _, tree := range t {
		tree.DeleteRaw(key)
	}
}

func (t MutableTrees) DeletePrefix(prefix string) int {
	var count int
	for _, tree := range t {
		count += tree.DeletePrefix(prefix)
	}
	return count
}

func (t MutableTrees) DeletePrefixRaw(prefix string) int {
	var count int
	for _, tree := range t {
		count += tree.DeletePrefixRaw(prefix)
	}
	return count
}

func (t MutableTrees) Lock(writable bool) (commit func()) {
	commits := make([]func(), len(t))
	for i, tree := range t {
		commits[i] = tree.Lock(writable)
	}
	return func() {
		for _, commit := range commits {
			commit()
		}
	}
}

func (t MutableTrees) CanLock() bool {
	for _, tree := range t {
		if !tree.CanLock() {
			return false
		}
	}
	return true
}

// WalkContext is passed to the Walk callback.
type WalkContext[T any] struct {
	data     *SimpleThreadSafeTree[any]
	dataInit sync.Once

	dataRaw     *maps.Cache[sitesmatrix.Vector, *SimpleThreadSafeTree[any]]
	dataRawInit sync.Once

	eventHandlersMu sync.Mutex
	eventHandlers   eventHandlers[T]
	eventMu         sync.Mutex
	events          []*Event[T]

	hooksPost1Init sync.Once
	hooksPost1     *collections.Stack[func() error]

	hooksPost2Init sync.Once
	hooksPost2     *collections.Stack[func() error]
}

type eventHandlers[T any] map[string][]func(*Event[T])

func cleanKey(key string) string {
	if key == "/" {
		// The path to the home page is logically "/",
		// but for technical reasons, it's stored as "".
		// This allows us to treat the home page as a section,
		// and a prefix search for "/" will return the home page's descendants.
		return ""
	}
	return key
}

func (ctx *WalkContext[T]) HandleEvents() error {
	ctx.eventHandlersMu.Lock()
	defer ctx.eventHandlersMu.Unlock()

	for len(ctx.events) > 0 {
		event := ctx.events[0]
		ctx.events = ctx.events[1:]

		// Loop the event handlers in reverse order so
		// that events created by the handlers themselves will
		// be picked up further up the tree.
		for i := len(ctx.eventHandlers[event.Name]) - 1; i >= 0; i-- {
			ctx.eventHandlers[event.Name][i](event)
			if event.stopPropagation {
				break
			}
		}
	}
	return nil
}

func (ctx *WalkContext[T]) HooksPost1() *collections.Stack[func() error] {
	ctx.hooksPost1Init.Do(func() {
		ctx.hooksPost1 = collections.NewStack[func() error]()
	})
	return ctx.hooksPost1
}

func (ctx *WalkContext[T]) HooksPost2() *collections.Stack[func() error] {
	ctx.hooksPost2Init.Do(func() {
		ctx.hooksPost2 = collections.NewStack[func() error]()
	})
	return ctx.hooksPost2
}

func (ctx *WalkContext[T]) HandleHooks1AndEventsAndHooks2() error {
	for _, hook := range ctx.HooksPost1().All() {
		if err := hook(); err != nil {
			return err
		}
	}

	if err := ctx.HandleEvents(); err != nil {
		return err
	}

	for _, hook := range ctx.HooksPost2().All() {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}

func mustValidateKey(key string) string {
	if err := ValidateKey(key); err != nil {
		panic(err)
	}
	return key
}
