// Copyright 2026 The Hugo Authors. All rights reserved.
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

package commands

import (
	"errors"
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/internal/warpc"
)

func TestErrStateShouldRecover(t *testing.T) {
	c := qt.New(t)

	broken := fmt.Errorf("render failed: %w", warpc.ErrBroken)
	other := errors.New("template error")

	var e hugoBuilderErrState

	steps := []struct {
		err  error
		want bool
	}{
		// timeout, ok, ok, timeout.
		{broken, true},
		{nil, false},
		{nil, false},
		{broken, true},
		// The recovery rebuild fails the same way: don't loop.
		{broken, false},
		// But recover again on the next, unrelated failure.
		{broken, true},
		// Any other build result resets.
		{other, false},
		{broken, true},
	}

	for i, step := range steps {
		c.Assert(e.shouldRecover(step.err), qt.Equals, step.want, qt.Commentf("step %d", i))
	}
}
