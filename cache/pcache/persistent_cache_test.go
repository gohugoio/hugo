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

package pcache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPersistentCache(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	cache, _ := createCache(t)

	vID := NewVersionedID("1", "ABC")

	created := false
	create := func() (Identifier, error) {
		top := &testObject{
			VersionedID: vID,
			Data:        "hi",
		}

		created = true

		return top, nil
	}

	to1v, err := cache.GetOrCreate(vID, create)
	assert.NoError(err)
	to1 := to1v.(*testObject)

	assert.Equal("ABC", to1._ID())
	assert.True(created)

	created = false
	to2, err := cache.GetOrCreate(vID, create)
	assert.NoError(err)
	assert.Equal(to1, to2)
	assert.False(created)

	assert.NoError(cache.Close())

	cache2, cleanup := createCacheFrom(t, cache)
	defer cleanup()

	to3, err := cache2.GetOrCreate(vID, create)
	assert.NoError(err)
	assert.False(created)
	assert.Equal(to1, to3)
}

func createCacheFrom(t *testing.T, from *persistentCache) (*persistentCache, func()) {
	c := New(from.filename, &testObject{})
	err := c.Open()
	if err != nil {
		t.Fatal(err)
	}

	return c.(*persistentCache), func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		os.RemoveAll(filepath.Dir(from.filename))
	}
}

func createCache(t *testing.T) (*persistentCache, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "hugodbcache")
	if err != nil {
		t.Fatal(err)
	}

	c := New(filepath.Join(dir, "hugocache.json"), &testObject{})
	err = c.Open()
	if err != nil {
		t.Fatal(err)
	}

	return c.(*persistentCache), func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		os.RemoveAll(dir)
	}
}
