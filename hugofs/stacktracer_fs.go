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

package hugofs

import (
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/gohugoio/hugo/common/types"

	"github.com/spf13/afero"
)

var (
	//  Make sure we don't accidentally use this in the real Hugo.
	_ types.DevMarker     = (*stacktracerFs)(nil)
	_ FilesystemUnwrapper = (*stacktracerFs)(nil)
)

// NewStacktracerFs wraps the given fs printing stack traces for file creates
// matching the given regexp pattern.
func NewStacktracerFs(fs afero.Fs, pattern string) afero.Fs {
	return &stacktracerFs{Fs: fs, re: regexp.MustCompile(pattern)}
}

// stacktracerFs can be used in hard-to-debug development situations where
// you get some input you don't understand where comes from.
type stacktracerFs struct {
	afero.Fs

	// Will print stacktrace for every file creates matching this pattern.
	re *regexp.Regexp
}

func (fs *stacktracerFs) DevOnly() {
}

func (fs *stacktracerFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

func (fs *stacktracerFs) onCreate(filename string) {
	if fs.re.MatchString(filename) {
		trace := make([]byte, 1500)
		runtime.Stack(trace, true)
		fmt.Printf("\n===========\n%q:\n%s\n", filename, trace)
	}
}

func (fs *stacktracerFs) Create(name string) (afero.File, error) {
	f, err := fs.Fs.Create(name)
	if err == nil {
		fs.onCreate(name)
	}
	return f, err
}

func (fs *stacktracerFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil && isWrite(flag) {
		fs.onCreate(name)
	}
	return f, err
}
