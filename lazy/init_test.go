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

func TestInit(t *testing.T) {
	assert := require.New(t)

	var result string

	bigOrSmall := func() int {
		if rand.Intn(10) < 3 {
			return 10000 + rand.Intn(100000)
		}
		return 1 + rand.Intn(50)
	}

	f1 := func(name string) func() (interface{}, error) {
		return func() (interface{}, error) {
			result += name + "|"
			size := bigOrSmall()
			_ = strings.Repeat("Hugo Rocks! ", size)
			return name, nil
		}
	}

	f2 := func() func() (interface{}, error) {
		return func() (interface{}, error) {
			size := bigOrSmall()
			_ = strings.Repeat("Hugo Rocks! ", size)
			return size, nil
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
			if rand.Intn(10) < 5 {
				_, err = root.Do()
				assert.NoError(err)
			}

			// Add a new branch on the fly.
			if rand.Intn(10) > 5 {
				branch := branch1_2.Branch(f2())
				init := branch.Add(f2())
				_, err = init.Do()
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
