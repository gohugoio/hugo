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
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

type testHashReceiver struct {
	name string
	sum  uint64
}

func (t *testHashReceiver) OnFileClose(name string, checksum uint64) {
	t.name = name
	t.sum = checksum
}

func TestHashingFs(t *testing.T) {
	c := qt.New(t)

	fs := afero.NewMemMapFs()
	observer := &testHashReceiver{}
	ofs := NewHashingFs(fs, observer)

	f, err := ofs.Create("hashme")
	c.Assert(err, qt.IsNil)
	_, err = f.Write([]byte("content"))
	c.Assert(err, qt.IsNil)
	c.Assert(f.Close(), qt.IsNil)
	c.Assert(observer.sum, qt.Equals, uint64(7807861979271768572))
	c.Assert(observer.name, qt.Equals, "hashme")

	f, err = ofs.Create("nowrites")
	c.Assert(err, qt.IsNil)
	c.Assert(f.Close(), qt.IsNil)
	c.Assert(observer.sum, qt.Equals, uint64(17241709254077376921))
}

func BenchmarkHashingFs(b *testing.B) {
	fs := afero.NewMemMapFs()
	observer := &testHashReceiver{}
	ofs := NewHashingFs(fs, observer)
	content := []byte(strings.Repeat("lorem ipsum ", 1000))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f, err := ofs.Create(fmt.Sprintf("file%d", i))
		if err != nil {
			b.Fatal(err)
		}
		_, err = f.Write(content)
		if err != nil {
			b.Fatal(err)
		}
		if err := f.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
