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

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestPrune(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	configStr := `
resourceDir = "myresources"
contentDir = "content"
dataDir = "data"
i18nDir = "i18n"
layoutDir = "layouts"
assetDir = "assets"
archeTypedir = "archetypes"

[caches]
[caches.getjson]
maxAge = "200ms"
dir = "/cache/c"
[caches.getcsv]
maxAge = "200ms"
dir = "/cache/d"
[caches.assets]
maxAge = "200ms"
dir = ":resourceDir/_gen"
[caches.images]
maxAge = "200ms"
dir = ":resourceDir/_gen"
`

	for _, name := range []string{cacheKeyGetCSV, cacheKeyGetJSON, cacheKeyAssets, cacheKeyImages} {
		msg := qt.Commentf("cache: %s", name)
		p := newPathsSpec(t, afero.NewMemMapFs(), configStr)
		caches, err := NewCaches(p)
		c.Assert(err, qt.IsNil)
		cache := caches[name]
		for i := 0; i < 10; i++ {
			id := fmt.Sprintf("i%d", i)
			cache.GetOrCreateBytes(id, func() ([]byte, error) {
				return []byte("abc"), nil
			})
			if i == 4 {
				// This will expire the first 5
				time.Sleep(201 * time.Millisecond)
			}
		}

		count, err := caches.Prune()
		c.Assert(err, qt.IsNil)
		c.Assert(count, qt.Equals, 5, msg)

		for i := 0; i < 10; i++ {
			id := fmt.Sprintf("i%d", i)
			v := cache.getString(id)
			if i < 5 {
				c.Assert(v, qt.Equals, "")
			} else {
				c.Assert(v, qt.Equals, "abc")
			}
		}

		caches, err = NewCaches(p)
		c.Assert(err, qt.IsNil)
		cache = caches[name]
		// Touch one and then prune.
		cache.GetOrCreateBytes("i5", func() ([]byte, error) {
			return []byte("abc"), nil
		})

		count, err = caches.Prune()
		c.Assert(err, qt.IsNil)
		c.Assert(count, qt.Equals, 4)

		// Now only the i5 should be left.
		for i := 0; i < 10; i++ {
			id := fmt.Sprintf("i%d", i)
			v := cache.getString(id)
			if i != 5 {
				c.Assert(v, qt.Equals, "")
			} else {
				c.Assert(v, qt.Equals, "abc")
			}
		}

	}
}
