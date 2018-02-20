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

package hugofs

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

type testHashReceiver struct {
	sum  string
	name string
}

func (t *testHashReceiver) OnFileClose(name, md5hash string) {
	t.name = name
	t.sum = md5hash
}

func TestHashingFs(t *testing.T) {
	assert := require.New(t)

	fs := afero.NewMemMapFs()
	observer := &testHashReceiver{}
	ofs := NewHashingFs(fs, observer)

	f, err := ofs.Create("hashme")
	assert.NoError(err)
	_, err = f.Write([]byte("content"))
	assert.NoError(err)
	assert.NoError(f.Close())
	assert.Equal("9a0364b9e99bb480dd25e1f0284c8555", observer.sum)
	assert.Equal("hashme", observer.name)

	f, err = ofs.Create("nowrites")
	assert.NoError(err)
	assert.NoError(f.Close())
	assert.Equal("d41d8cd98f00b204e9800998ecf8427e", observer.sum)

}
