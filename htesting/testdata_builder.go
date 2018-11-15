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

package htesting

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

type testFile struct {
	name    string
	content string
}

type testdataBuilder struct {
	t          testing.TB
	fs         afero.Fs
	workingDir string

	files []testFile
}

func NewTestdataBuilder(fs afero.Fs, workingDir string, t testing.TB) *testdataBuilder {
	workingDir = filepath.Clean(workingDir)
	return &testdataBuilder{fs: fs, workingDir: workingDir, t: t}
}

func (b *testdataBuilder) Add(filename, content string) *testdataBuilder {
	b.files = append(b.files, testFile{name: filename, content: content})
	return b
}

func (b *testdataBuilder) Build() *testdataBuilder {
	for _, f := range b.files {
		if err := afero.WriteFile(b.fs, filepath.Join(b.workingDir, f.name), []byte(f.content), 0666); err != nil {
			b.t.Fatalf("failed to add %q: %s", f.name, err)
		}
	}
	return b
}

func (b testdataBuilder) WithWorkingDir(dir string) *testdataBuilder {
	b.workingDir = filepath.Clean(dir)
	b.files = make([]testFile, 0)
	return &b
}
