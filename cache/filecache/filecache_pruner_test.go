// Copyright 2018 The Hugo Authors. All rights reserved.
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

package filecache

import (
	"fmt"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"

	"github.com/stretchr/testify/require"
)

func TestPrune(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	configStr := `
resourceDir = "myresources"
[caches]
[caches.getjson]
maxAge = "200ms"
dir = "/cache/c"

`

	cfg, err := config.FromConfigString(configStr, "toml")
	assert.NoError(err)
	fs := hugofs.NewMem(cfg)
	p, err := paths.New(fs, cfg)
	assert.NoError(err)

	caches, err := NewCachesFromPaths(p)
	assert.NoError(err)

	jsonCache := caches.GetJSONCache()
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("i%d", i)
		jsonCache.GetOrCreateBytes(id, func() ([]byte, error) {
			return []byte("abc"), nil
		})
		if i == 4 {
			// This will expire the first 5
			time.Sleep(201 * time.Millisecond)
		}
	}

	count, err := caches.Prune()
	assert.NoError(err)
	assert.Equal(5, count)

	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("i%d", i)
		v := jsonCache.getString(id)
		if i < 5 {
			assert.Equal("", v, id)
		} else {
			assert.Equal("abc", v, id)
		}
	}

	caches, err = NewCachesFromPaths(p)
	assert.NoError(err)
	jsonCache = caches.GetJSONCache()
	// Touch one and then prune.
	jsonCache.GetOrCreateBytes("i5", func() ([]byte, error) {
		return []byte("abc"), nil
	})

	count, err = caches.Prune()
	assert.NoError(err)
	assert.Equal(4, count)

	// Now only the i5 should be left.
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("i%d", i)
		v := jsonCache.getString(id)
		if i != 5 {
			assert.Equal("", v, id)
		} else {
			assert.Equal("abc", v, id)
		}
	}

}
