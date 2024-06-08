// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"os"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/spf13/afero"
)

var (
	_ afero.Fs            = (*hasBytesFs)(nil)
	_ FilesystemUnwrapper = (*hasBytesFs)(nil)
)

type hasBytesFs struct {
	afero.Fs
	shouldCheck      func(name string) bool
	hasBytesCallback func(name string, match []byte)
	patterns         [][]byte
}

func NewHasBytesReceiver(delegate afero.Fs, shouldCheck func(name string) bool, hasBytesCallback func(name string, match []byte), patterns ...[]byte) afero.Fs {
	return &hasBytesFs{Fs: delegate, shouldCheck: shouldCheck, hasBytesCallback: hasBytesCallback, patterns: patterns}
}

func (fs *hasBytesFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

func (fs *hasBytesFs) Create(name string) (afero.File, error) {
	f, err := fs.Fs.Create(name)
	if err == nil {
		f = fs.wrapFile(f)
	}
	return f, err
}

func (fs *hasBytesFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil && isWrite(flag) {
		f = fs.wrapFile(f)
	}
	return f, err
}

func (fs *hasBytesFs) wrapFile(f afero.File) afero.File {
	if !fs.shouldCheck(f.Name()) {
		return f
	}
	patterns := make([]*hugio.HasBytesPattern, len(fs.patterns))
	for i, p := range fs.patterns {
		patterns[i] = &hugio.HasBytesPattern{Pattern: p}
	}

	return &hasBytesFile{
		File: f,
		hbw: &hugio.HasBytesWriter{
			Patterns: patterns,
		},
		hasBytesCallback: fs.hasBytesCallback,
	}
}

func (fs *hasBytesFs) Name() string {
	return "hasBytesFs"
}

type hasBytesFile struct {
	hasBytesCallback func(name string, match []byte)
	hbw              *hugio.HasBytesWriter
	afero.File
}

func (h *hasBytesFile) Write(p []byte) (n int, err error) {
	n, err = h.File.Write(p)
	if err != nil {
		return
	}
	return h.hbw.Write(p)
}

func (h *hasBytesFile) Close() error {
	for _, p := range h.hbw.Patterns {
		if p.Match {
			h.hasBytesCallback(h.Name(), p.Pattern)
		}
	}
	return h.File.Close()
}
