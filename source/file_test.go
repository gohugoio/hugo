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
	"github.com/spf13/hugo/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileUniqueID(t *testing.T) {
	f1 := File{uniqueID: "123"}
	f2 := NewFile("a")

	assert.Equal(t, "123", f1.UniqueID())
	assert.Equal(t, "0cc175b9c0f1b6a831c399e269772661", f2.UniqueID())
}

func TestFileString(t *testing.T) {
	assert.Equal(t, "abc", NewFileWithContents("a", helpers.StringToReader("abc")).String())
	assert.Equal(t, "", NewFile("a").String())
}

func TestFileBytes(t *testing.T) {
	assert.Equal(t, []byte("abc"), NewFileWithContents("a", helpers.StringToReader("abc")).Bytes())
	assert.Equal(t, []byte(""), NewFile("a").Bytes())
}
