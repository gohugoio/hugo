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

package identity

import "sync"

// NewQuestion creates a new question with the given identity.
func NewQuestion[T any](id Identity) *Question[T] {
	return &Question[T]{
		Identity: id,
	}
}

// Answer takes a func that knows the answer.
// Note that this is a one-time operation,
// fn will not be invoked again it the question is already answered.
// Use Result to check if the question is answered.
func (q *Question[T]) Answer(fn func() T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.answered {
		return
	}

	q.fasit = fn()
	q.answered = true
}

// Result returns the fasit of the question (if answered),
// and a bool indicating if the question has been answered.
func (q *Question[T]) Result() (any, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.fasit, q.answered
}

// A Question is defined by its Identity and can be answered once.
type Question[T any] struct {
	Identity
	fasit T

	mu       sync.RWMutex
	answered bool
}
