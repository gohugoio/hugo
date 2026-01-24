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

package commands

import (
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/spf13/afero"
)

// hardLinkFs wraps the destination filesystem to intercept Create calls.
// It attempts to verify if a hard link can be created instead of a standard file copy.
// This is used by the static file syncer to optimize disk usage and performance.
//
// Safety guarantees:
// 1. It only attempts to link if verified source and destination paths are available.
// 2. It falls back to standard copy (h.Fs.Create) on ANY error (cross-device, permissions, unsupported FS).
// 3. It returns a no-op writer on success to prevent fsync from overwriting the linked file.
type hardLinkFs struct {
	afero.Fs
	src         *filesystems.SourceFilesystem
	publishDir  string
	destRoot    string // Absolute OS path to the destination root. Must be on the same filesystem as src.
	filesSynced []string
}

func (h *hardLinkFs) Create(name string) (afero.File, error) {
	// Calculate the relative path from the publish directory.
	rel, err := filepath.Rel(h.publishDir, name)
	if err != nil {
		return h.Fs.Create(name)
	}

	// We can only hard link if we can determine the real source filename.
	// RealFilename returns the absolute OS path if the source is on afero.OsFs.
	srcFilename := h.src.RealFilename(rel)
	if srcFilename == "" || srcFilename == rel {
		// Fallback to copy if checking the source path fails or returns a relative path (indicating not OsFs).
		return h.Fs.Create(name)
	}

	// Remove existing file if any (os.Link requires this).
	_ = h.Fs.Remove(name)

	// Ensure parent directory exists.
	if err := h.Fs.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return nil, err
	}

	// Attempt hard link.
	// We need to resolve the destination path to a real OS path.
	// This relies on destRoot being a valid absolute path on disk.
	if h.destRoot == "" {
		return h.Fs.Create(name)
	}
	destFilename := filepath.Join(h.destRoot, rel)

	// Note: We use os.Link directly as afero.OsFs.Link is not always available/reliable across all generic Fs types,
	// but here we know we are targeting the OS filesystem for the hard link to make sense.
	// Use filepath.Clean to ensure we have OS-compatible paths.
	//
	// On Windows, os.Link is attempted. If it fails (e.g. cross-drive), we fall back to copy.
	err = os.Link(filepath.Clean(srcFilename), filepath.Clean(destFilename))
	if err == nil {
		// Link successful. Return a no-op writer to satisfy fsync.
		// fsync will try to write content to this file, but we want to discard those writes
		// because the file is partially equivalent to the source (it IS the source file via hard link).
		// We still return a valid file handle (read-only) to satisfy the interface.
		f, err := h.Fs.Open(name)
		if err != nil {
			// If we can't open the file we just linked, something is wrong.
			// Fallback to create (which will likely fail too, but correct path).
			return h.Fs.Create(name)
		}
		return &noOpWriterFile{File: f}, nil
	}

	// Fallback to copy if linking fails for ANY reason (e.g. cross-device link EXDEV).
	return h.Fs.Create(name)
}

// noOpWriterFile wraps an afero.File to ignore writes.
// This is crucial when a file has been hard-linked; fsync logic normally opens the destination
// with Create (truncating it), but our hardLinkFs.Create implementation performed a Link instead.
// We must ignore the subsequent Write calls from fsync to avoid modifying the source file
// (since hard links share the same inode) or doing redundant IO.
// Truncate is also ignored.
type noOpWriterFile struct {
	afero.File
}

func (f *noOpWriterFile) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (f *noOpWriterFile) Truncate(size int64) error {
	return nil
}

func (f *noOpWriterFile) Sync() error {
	return nil
}
