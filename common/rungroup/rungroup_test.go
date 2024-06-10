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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNew(t *testing.T) {
	c := qt.New(t)

	var result int
	adder := func(ctx context.Context, i int) error {
		result += i
		return nil
	}

	g := Run[int](
		context.Background(),
		Config[int]{
			Handle: adder,
		},
	)

	c.Assert(g, qt.IsNotNil)
	g.Enqueue(32)
	g.Enqueue(33)
	c.Assert(g.Wait(), qt.IsNil)
	c.Assert(result, qt.Equals, 65)
}
