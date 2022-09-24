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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/spf13/afero"
)

// Glob walks the fs and passes all matches to the handle func.
// The handle func can return true to signal a stop.
func Glob(fs afero.Fs, pattern string, handle func(fi FileMetaInfo) (bool, error)) error {
	pattern = glob.NormalizePathNoLower(pattern)
	if pattern == "" {
		return nil
	}
	root := glob.ResolveRootDir(pattern)
	pattern = strings.ToLower(pattern)

	g, err := glob.GetGlob(pattern)
	if err != nil {
		return err
	}

	hasSuperAsterisk := strings.Contains(pattern, "**")
	levels := strings.Count(pattern, "/")

	// Signals that we're done.
	done := errors.New("done")

	wfn := func(p string, info FileMetaInfo, err error) error {
		p = glob.NormalizePath(p)
		if info.IsDir() {
			if !hasSuperAsterisk {
				// Avoid walking to the bottom if we can avoid it.
				if p != "" && strings.Count(p, "/") >= levels {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if g.Match(p) {
			d, err := handle(info)
			if err != nil {
				return err
			}
			if d {
				return done
			}
		}

		return nil
	}

	root, err = CanonicalizeFilepath(fs, root)

	if os.IsNotExist(err) {
		// TODO logger
		return nil
	} else if err != nil {
		return err
	}

	w := NewWalkway(WalkwayConfig{
		Root:   root,
		Fs:     fs,
		WalkFn: wfn,
	})

	err = w.Walk()

	if err != done {
		return err
	}

	return nil
}

// return case-matched path from case-insensitive input
// before using this, you have to apply `NormalizePath`
func CanonicalizeFilepath(fs afero.Fs, path string) (string, error) {
	if path == "." || path == "" {
		return path, nil
	}
	paths := strings.Split(path, "/")

	var ret []string

	for _, p := range paths {
		dir := filepath.Join(ret...)
		joined := filepath.Join(dir, p)
		fi, ok, err := lstatIfPossible(fs, joined)

		if !os.IsNotExist(err) {
			return "", err
		}

		if ok {
			ret = append(ret, fi.Name())
			continue
		}

		fi, err = statWithCaseInsensitiveName(fs, dir, p)
		if err != nil {
			return "", err
		}

		ret = append(ret, fi.Name())
		// ret = append(ret, p)
	}

	return filepath.Join(ret...), nil
}

// case-insenstively search in parent dir, and return os.FileInfo
func statWithCaseInsensitiveName(fs afero.Fs, parent, name string) (os.FileInfo, error) {
	fi, err := fs.Stat(parent)

	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not directory", parent)
	}

	f, err := fs.Open(parent)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)

	if err != nil {
		return nil, err
	}

	baseLowered := strings.ToLower(name)
	for _, name := range names {
		if baseLowered == strings.ToLower(name) {
			p := filepath.Join(parent, name)
			fi, _, err := lstatIfPossible(fs, p)
			if err != nil {
				return nil, err
			}
			return fi, err
		}
	}
	return nil, os.ErrNotExist
}
