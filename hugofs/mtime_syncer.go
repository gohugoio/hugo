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

// MtimeSyncer provides mtime-based file synchronization, replacing the
// spf13/fsync dependency for Hugo's static file handling.
//
// BACKGROUND:
// Traditional file synchronization compares files by content to detect
// changes. While accurate, this has O(file_size) complexity per file,
// making it slow for large static directories (many images, videos, etc.).
//
// APPROACH:
// MtimeSyncer uses modification times (mtime) and file sizes to decide
// whether a file needs copying. This is O(1) per file and significantly
// faster for incremental builds where most files haven't changed.
//
// ALGORITHM:
// For each source file:
//  1. If destination does not exist, copy.
//  2. If sizes differ, copy (different size means different content).
//  3. If source mtime is newer than destination mtime, copy.
//  4. Otherwise skip — destination is up to date.
//
// TRADE-OFFS:
//   - May miss a change if a file is modified without updating its mtime
//     (uncommon outside of programmatic manipulation).
//   - Destination mtime must be preserved after each copy; see syncStats.

package hugofs

import (
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// MtimeSyncer syncs files from source to destination using mtime+size
// comparison instead of content hashing. Replaces spf13/fsync.
type MtimeSyncer struct {
	SrcFs        afero.Fs
	DestFs       afero.Fs
	NoTimes      bool
	NoChmod      bool
	ChmodFilter  func(dst, src os.FileInfo) bool
	Delete       bool
	DeleteFilter func(f FileInfo) bool
}

// FileInfo is the interface for file info used by DeleteFilter.
type FileInfo interface {
	Name() string
	IsDir() bool
}

// Sync syncs srcRoot into dstRoot. Recursive traversal correctly handles
// Hugo's union/overlay module mounts.
func (s *MtimeSyncer) Sync(dstRoot, srcRoot string) error {
	if _, err := s.SrcFs.Stat(srcRoot); err != nil {
		return err
	}

	return s.syncRecursive(dstRoot, srcRoot)
}

func (s *MtimeSyncer) syncRecursive(dst, src string) error {
	sstat, err := s.SrcFs.Stat(src)
	if os.IsNotExist(err) {
		return nil // deleted between directory listing and stat
	}
	if err != nil {
		return err
	}

	dstat, err := s.DestFs.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !sstat.IsDir() {
		if dstat != nil && dstat.IsDir() {
			if err := s.DestFs.RemoveAll(dst); err != nil {
				return err
			}
		}
		if s.needsCopy(dstat, sstat) {
			if err := s.copyFile(dst, src, sstat); err != nil {
				return err
			}
		}
		return s.syncStats(dst, src)
	}

	if dstat == nil {
		if err := s.DestFs.MkdirAll(dst, 0o755); err != nil {
			return err
		}
	} else if !dstat.IsDir() {
		if err := s.DestFs.Remove(dst); err != nil {
			return err
		}
		if err := s.DestFs.MkdirAll(dst, 0o755); err != nil {
			return err
		}
	}

	srcEntries := make(map[string]bool)
	err = s.withDirEntries(s.SrcFs, src, func(fi FileInfo) error {
		srcEntries[fi.Name()] = true
		return s.syncRecursive(filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name()))
	})
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if s.Delete {
		err = s.withDirEntries(s.DestFs, dst, func(fi FileInfo) error {
			if !srcEntries[fi.Name()] {
				if s.DeleteFilter != nil && s.DeleteFilter(fi) {
					return nil
				}
				return s.DestFs.RemoveAll(filepath.Join(dst, fi.Name()))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return s.syncStats(dst, src)
}

func (s *MtimeSyncer) needsCopy(dstat, sstat os.FileInfo) bool {
	if dstat == nil {
		return true
	}
	return dstat.Size() != sstat.Size() || sstat.ModTime().After(dstat.ModTime())
}

func (s *MtimeSyncer) copyFile(dst, src string, sstat os.FileInfo) error {
	if err := s.DestFs.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	srcFile, err := s.SrcFs.Open(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.DestFs.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (s *MtimeSyncer) syncStats(dst, src string) error {
	dstat, err1 := s.DestFs.Stat(dst)
	sstat, err2 := s.SrcFs.Stat(src)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		return nil
	}
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	if !s.NoChmod {
		if s.ChmodFilter == nil || !s.ChmodFilter(dstat, sstat) {
			if dstat.Mode().Perm() != sstat.Mode().Perm() {
				if err := s.DestFs.Chmod(dst, sstat.Mode().Perm()); err != nil {
					return err
				}
			}
		}
	}

	// Must sync mtime so future needsCopy comparisons stay accurate.
	if !s.NoTimes {
		if !dstat.ModTime().Equal(sstat.ModTime()) {
			if err := s.DestFs.Chtimes(dst, sstat.ModTime(), sstat.ModTime()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *MtimeSyncer) withDirEntries(fs afero.Fs, path string, fn func(FileInfo) error) error {
	f, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if rdf, ok := f.(iofs.ReadDirFile); ok {
		entries, err := rdf.ReadDir(-1)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := fn(entry); err != nil {
				return err
			}
		}
		return nil
	}

	fis, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if err := fn(fi); err != nil {
			return err
		}
	}
	return nil
}
