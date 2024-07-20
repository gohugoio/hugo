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

package collections

import "sync"

// Stack is a simple LIFO stack that is safe for concurrent use.
type Stack[T any] struct {
	items []T
	zero  T
	mu    sync.RWMutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		return s.zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return s.zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

func (s *Stack[T]) Drain() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.items
	s.items = nil
	return items
}

func (s *Stack[T]) DrainMatching(predicate func(T) bool) []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	var items []T
	for i := len(s.items) - 1; i >= 0; i-- {
		if predicate(s.items[i]) {
			items = append(items, s.items[i])
			s.items = append(s.items[:i], s.items[i+1:]...)
		}
	}
	return items
}
