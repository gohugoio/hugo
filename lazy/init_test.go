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

package lazy

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
	bigOrSmall = func() int {
		if rnd.Intn(10) < 5 {
			return 10000 + rnd.Intn(100000)
		}
		return 1 + rnd.Intn(50)
	}
)

func doWork() {
	doWorkOfSize(bigOrSmall())
}

func doWorkOfSize(size int) {
	_ = strings.Repeat("Hugo Rocks! ", size)
}

func TestInit(t *testing.T) {
	assert := require.New(t)

	var result string

	f1 := func(name string) func() (interface{}, error) {
		return func() (interface{}, error) {
			result += name + "|"
			doWork()
			return name, nil
		}
	}

	f2 := func() func() (interface{}, error) {
		return func() (interface{}, error) {
			doWork()
			return nil, nil
		}
	}

	root := New()

	root.Add(f1("root(1)"))
	root.Add(f1("root(2)"))

	branch1 := root.Branch(f1("branch_1"))
	branch1.Add(f1("branch_1_1"))
	branch1_2 := branch1.Add(f1("branch_1_2"))
	branch1_2_1 := branch1_2.Add(f1("branch_1_2_1"))

	var wg sync.WaitGroup

	// Add some concurrency and randomness to verify thread safety and
	// init order.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var err error
			if rnd.Intn(10) < 5 {
				_, err = root.Do()
				assert.NoError(err)
			}

			// Add a new branch on the fly.
			if rnd.Intn(10) > 5 {
				branch := branch1_2.Branch(f2())
				_, err = branch.Do()
				assert.NoError(err)
			} else {
				_, err = branch1_2_1.Do()
				assert.NoError(err)
			}
			_, err = branch1_2.Do()
			assert.NoError(err)
		}(i)

		wg.Wait()

		assert.Equal("root(1)|root(2)|branch_1|branch_1_1|branch_1_2|branch_1_2_1|", result)

	}
}

func TestInitAddWithTimeout(t *testing.T) {
	assert := require.New(t)

	init := New().AddWithTimeout(100*time.Millisecond, func(ctx context.Context) (interface{}, error) {
		return nil, nil
	})

	_, err := init.Do()

	assert.NoError(err)
}

func TestInitAddWithTimeoutTimeout(t *testing.T) {
	assert := require.New(t)

	init := New().AddWithTimeout(100*time.Millisecond, func(ctx context.Context) (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}
		t.Fatal("slept")
		return nil, nil
	})

	_, err := init.Do()

	assert.Error(err)

	assert.Contains(err.Error(), "timed out")

	time.Sleep(1 * time.Second)
}

func TestInitAddWithTimeoutError(t *testing.T) {
	assert := require.New(t)

	init := New().AddWithTimeout(100*time.Millisecond, func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("failed")
	})

	_, err := init.Do()

	assert.Error(err)
}

type T struct {
	sync.Mutex
	V1 string
	V2 string
}

func (t *T) Add1(v string) {
	t.Lock()
	t.V1 += v
	t.Unlock()
}

func (t *T) Add2(v string) {
	t.Lock()
	t.V2 += v
	t.Unlock()
}

// https://github.com/gohugoio/hugo/issues/5901
func TestInitBranchOrder(t *testing.T) {
	assert := require.New(t)

	base := New()

	work := func(size int, f func()) func() (interface{}, error) {
		return func() (interface{}, error) {
			doWorkOfSize(size)
			if f != nil {
				f()
			}

			return nil, nil
		}
	}

	state := &T{}

	base = base.Add(work(10000, func() {
		state.Add1("A")
	}))

	inits := make([]*Init, 2)
	for i := range inits {
		inits[i] = base.Branch(work(i+1*100, func() {
			// V1 is A
			ab := state.V1 + "B"
			state.Add2(ab)
		}))
	}

	var wg sync.WaitGroup

	for _, v := range inits {
		v := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := v.Do()
			assert.NoError(err)
		}()
	}

	wg.Wait()

	assert.Equal("ABAB", state.V2)
}
