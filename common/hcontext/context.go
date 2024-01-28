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

package hcontext

import "context"

// ContextDispatcher is a generic interface for setting and getting values from a context.
type ContextDispatcher[T any] interface {
	Set(ctx context.Context, value T) context.Context
	Get(ctx context.Context) T
}

// NewContextDispatcher creates a new ContextDispatcher with the given key.
func NewContextDispatcher[T any, R comparable](key R) ContextDispatcher[T] {
	return keyInContext[T, R]{
		id: key,
	}
}

type keyInContext[T any, R comparable] struct {
	zero T
	id   R
}

func (f keyInContext[T, R]) Get(ctx context.Context) T {
	v := ctx.Value(f.id)
	if v == nil {
		return f.zero
	}
	return v.(T)
}

func (f keyInContext[T, R]) Set(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, f.id, value)
}
