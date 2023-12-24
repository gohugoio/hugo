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

package rungroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Group is a group of workers that can be used to enqueue work and wait for
// them to finish.
type Group[T any] interface {
	Enqueue(T) error
	Wait() error
}

type runGroup[T any] struct {
	ctx context.Context
	g   *errgroup.Group
	ch  chan T
}

// Config is the configuration for a new Group.
type Config[T any] struct {
	NumWorkers int
	Handle     func(context.Context, T) error
}

// Run creates a new Group with the given configuration.
func Run[T any](ctx context.Context, cfg Config[T]) Group[T] {
	if cfg.NumWorkers <= 0 {
		cfg.NumWorkers = 1
	}
	if cfg.Handle == nil {
		panic("Handle must be set")
	}

	g, ctx := errgroup.WithContext(ctx)
	// Buffered for performance.
	ch := make(chan T, cfg.NumWorkers)

	for i := 0; i < cfg.NumWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case v, ok := <-ch:
					if !ok {
						return nil
					}
					if err := cfg.Handle(ctx, v); err != nil {
						return err
					}
				}
			}
		})
	}

	return &runGroup[T]{
		ctx: ctx,
		g:   g,
		ch:  ch,
	}
}

// Enqueue enqueues a new item to be handled by the workers.
func (r *runGroup[T]) Enqueue(t T) error {
	select {
	case <-r.ctx.Done():
		return nil
	case r.ch <- t:
	}
	return nil
}

// Wait waits for all workers to finish and returns the first error.
func (r *runGroup[T]) Wait() error {
	close(r.ch)
	return r.g.Wait()
}
