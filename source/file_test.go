// Copyright 2015 The Hugo Authors. All rights reserved.
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

package source

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestFileUniqueID(t *testing.T) {
	ss := newTestSourceSpec()

	f1 := File{uniqueID: "123"}
	f2 := ss.NewFile("a")

	assert.Equal(t, "123", f1.UniqueID())
	assert.Equal(t, "0cc175b9c0f1b6a831c399e269772661", f2.UniqueID())

	f3 := ss.NewFile(filepath.FromSlash("test1/index.md"))
	f4 := ss.NewFile(filepath.FromSlash("test2/index.md"))

	assert.NotEqual(t, f3.UniqueID(), f4.UniqueID())

	f5l := ss.NewFile("test3/index.md")
	f5w := ss.NewFile(filepath.FromSlash("test3/index.md"))

	assert.Equal(t, f5l.UniqueID(), f5w.UniqueID())
}

func TestFileString(t *testing.T) {
	ss := newTestSourceSpec()
	assert.Equal(t, "abc", ss.NewFileWithContents("a", strings.NewReader("abc")).String())
	assert.Equal(t, "", ss.NewFile("a").String())
}

func TestFileBytes(t *testing.T) {
	ss := newTestSourceSpec()
	assert.Equal(t, []byte("abc"), ss.NewFileWithContents("a", strings.NewReader("abc")).Bytes())
	assert.Equal(t, []byte(""), ss.NewFile("a").Bytes())
}

func newTestSourceSpec() SourceSpec {
	v := viper.New()
	return SourceSpec{Fs: hugofs.NewMem(v), Cfg: v}
}
