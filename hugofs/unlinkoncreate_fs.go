// Copyright 2026 The Hugo Authors. All rights reserved.
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

	"github.com/spf13/afero"
)

var _ FilesystemUnwrapper = (*unlinkOnCreateFs)(nil)

// NewUnlinkOnCreateFs creates a new filesystem that removes a hard linked file
// before truncating it, so the other links are left untouched.
// fs must be an OS file system.
func NewUnlinkOnCreateFs(fs afero.Fs) afero.Fs {
	return &unlinkOnCreateFs{Fs: fs}
}

type unlinkOnCreateFs struct {
	afero.Fs
}

func (fs *unlinkOnCreateFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

func (fs *unlinkOnCreateFs) Create(name string) (afero.File, error) {
	fs.unlinkIfHardlinked(name)
	return fs.Fs.Create(name)
}

func (fs *unlinkOnCreateFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if flag&os.O_TRUNC != 0 {
		fs.unlinkIfHardlinked(name)
	}
	return fs.Fs.OpenFile(name, flag, perm)
}

func (fs *unlinkOnCreateFs) unlinkIfHardlinked(name string) {
	fi, err := fs.Fs.Stat(name)
	if err != nil || fi.IsDir() || !hasMultipleLinks(name, fi) {
		return
	}
	_ = fs.Fs.Remove(name)
}
