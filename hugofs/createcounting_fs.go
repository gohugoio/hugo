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
	"sort"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

// Reseter is implemented by some of the stateful filesystems.
type Reseter interface {
	Reset()
}

// DuplicatesReporter reports about duplicate filenames.
type DuplicatesReporter interface {
	ReportDuplicates() string
}

func NewCreateCountingFs(fs afero.Fs) afero.Fs {
	return &createCountingFs{Fs: fs, fileCount: make(map[string]int)}
}

// ReportDuplicates reports filenames written more than once.
func (c *createCountingFs) ReportDuplicates() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var dupes []string

	for k, v := range c.fileCount {
		if v > 1 {
			dupes = append(dupes, fmt.Sprintf("%s (%d)", k, v))
		}
	}

	if len(dupes) == 0 {
		return ""
	}

	sort.Strings(dupes)

	return strings.Join(dupes, ", ")
}

// createCountingFs counts filenames of created files or files opened
// for writing.
type createCountingFs struct {
	afero.Fs

	mu        sync.Mutex
	fileCount map[string]int
}

func (c *createCountingFs) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.fileCount = make(map[string]int)
}

func (fs *createCountingFs) onCreate(filename string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.fileCount[filename] = fs.fileCount[filename] + 1
}

func (fs *createCountingFs) Create(name string) (afero.File, error) {
	f, err := fs.Fs.Create(name)
	if err == nil {
		fs.onCreate(name)
	}
	return f, err
}

func (fs *createCountingFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil && isWrite(flag) {
		fs.onCreate(name)
	}
	return f, err
}
