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
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
)

// Linker creates hard links in the publish dir when possible.
type Linker struct {
	// Source devices where hard linking is not possible.
	disabled sync.Map
}

// Link creates dst in dstFs as a hard link to the OS file src, replacing any existing dst.
// It returns false if the link could not be created and dst must be copied instead.
func (l *Linker) Link(src string, dstFs afero.Fs, dst string) bool {
	dstFilename, ok := RealFilename(dstFs, dst)
	if !ok {
		return false
	}
	fi, err := os.Lstat(src)
	// Symlinks must be resolved by copying, and read-only files (e.g. from the module cache) must be published writable.
	if err != nil || !fi.Mode().IsRegular() || fi.Mode().Perm()&0o200 == 0 {
		return false
	}
	dev := deviceID(fi)
	if _, found := l.disabled.Load(dev); found {
		return false
	}
	_ = dstFs.Remove(dst)
	err = os.Link(src, dstFilename)
	if errors.Is(err, iofs.ErrNotExist) {
		if dstFs.MkdirAll(filepath.Dir(dst), 0o777) == nil {
			err = os.Link(src, dstFilename)
		}
	}
	if err == nil {
		return true
	}
	if isLinkUnsupported(err) {
		l.disabled.Store(dev, true)
	}
	return false
}
