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

// Package para implements parallel execution helpers.
package para

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Workers configures a task executor with the most number of tasks to be executed in parallel.
type Workers struct {
	sem chan struct{}
}

// Runner wraps the lifecycle methods of a new task set.
//
// Run will block until a worker is available or the context is cancelled,
// and then run the given func in a new goroutine.
// Wait will wait for all the running goroutines to finish.
type Runner interface {
	Run(func() error)
	Wait() error
}

type errGroupRunner struct {
	*errgroup.Group
	w   *Workers
	ctx context.Context
}

func (g *errGroupRunner) Run(fn func() error) {
	select {
	case g.w.sem <- struct{}{}:
	case <-g.ctx.Done():
		return
	}

	g.Go(func() error {
		err := fn()
		<-g.w.sem
		return err
	})
}

// New creates a new Workers with the given number of workers.
func New(numWorkers int) *Workers {
	return &Workers{
		sem: make(chan struct{}, numWorkers),
	}
}

// Start starts a new Runner.
func (w *Workers) Start(ctx context.Context) (Runner, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &errGroupRunner{
		Group: g,
		ctx:   ctx,
		w:     w,
	}, ctx
}
