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

// Package hugofs provides the file systems used by Hugo.
package hugofs

import (
	"fmt"
	"os"
	"strings"

	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
)

// Os points to the (real) Os filesystem.
var Os = &afero.OsFs{}

// Fs holds the core filesystems used by Hugo.
type Fs struct {
	// Source is Hugo's source file system.
	// Note that this will always be a "plain" Afero filesystem:
	// * afero.OsFs when running in production
	// * afero.MemMapFs for many of the tests.
	Source afero.Fs

	// PublishDir is where Hugo publishes its rendered content.
	// It's mounted inside publishDir (default /public).
	PublishDir afero.Fs

	// PublishDirStatic is the file system used for static files.
	PublishDirStatic afero.Fs

	// PublishDirServer is the file system used for serving the public directory with Hugo's development server.
	// This will typically be the same as PublishDir, but not if --renderStaticToDisk is set.
	PublishDirServer afero.Fs

	// Os is an OS file system.
	// NOTE: Field is currently unused.
	Os afero.Fs

	// WorkingDirReadOnly is a read-only file system
	// restricted to the project working dir.
	WorkingDirReadOnly afero.Fs

	// WorkingDirWritable is a writable file system
	// restricted to the project working dir.
	WorkingDirWritable afero.Fs
}

func NewDefault(cfg config.Provider) *Fs {
	workingDir, publishDir := getWorkingPublishDir(cfg)
	fs := Os
	return newFs(fs, fs, workingDir, publishDir)
}

// NewFrom creates a new Fs based on the provided Afero Fs
// as source and destination file systems.
// Useful for testing.
func NewFrom(fs afero.Fs, conf config.BaseConfig) *Fs {
	return newFs(fs, fs, conf.WorkingDir, conf.PublishDir)
}

func NewFromOld(fs afero.Fs, cfg config.Provider) *Fs {
	workingDir, publishDir := getWorkingPublishDir(cfg)
	return newFs(fs, fs, workingDir, publishDir)
}

// NewFromSourceAndDestination creates a new Fs based on the provided Afero Fss
// as the source and destination file systems.
func NewFromSourceAndDestination(source, destination afero.Fs, cfg config.Provider) *Fs {
	workingDir, publishDir := getWorkingPublishDir(cfg)
	return newFs(source, destination, workingDir, publishDir)
}

func getWorkingPublishDir(cfg config.Provider) (string, string) {
	workingDir := cfg.GetString("workingDir")
	publishDir := cfg.GetString("publishDirDynamic")
	if publishDir == "" {
		publishDir = cfg.GetString("publishDir")
	}
	return workingDir, publishDir
}

func newFs(source, destination afero.Fs, workingDir, publishDir string) *Fs {
	if publishDir == "" {
		panic("publishDir is empty")
	}

	if workingDir == "." {
		workingDir = ""
	}

	// Sanity check
	if IsOsFs(source) && len(workingDir) < 2 {
		panic("workingDir is too short")
	}

	// If this does not exist, it will be created later.
	absPublishDir := paths.AbsPathify(workingDir, publishDir)

	pubFs := NewBasePathFs(destination, absPublishDir)

	return &Fs{
		Source:             source,
		PublishDir:         pubFs,
		PublishDirServer:   pubFs,
		PublishDirStatic:   pubFs,
		Os:                 &afero.OsFs{},
		WorkingDirReadOnly: getWorkingDirFsReadOnly(source, workingDir),
		WorkingDirWritable: getWorkingDirFsWritable(source, workingDir),
	}
}

func getWorkingDirFsReadOnly(base afero.Fs, workingDir string) afero.Fs {
	if workingDir == "" {
		return NewReadOnlyFs(base)
	}
	return NewBasePathFs(NewReadOnlyFs(base), workingDir)
}

func getWorkingDirFsWritable(base afero.Fs, workingDir string) afero.Fs {
	if workingDir == "" {
		return base
	}
	return NewBasePathFs(base, workingDir)
}

func isWrite(flag int) bool {
	return flag&os.O_RDWR != 0 || flag&os.O_WRONLY != 0
}

// MakeReadableAndRemoveAllModulePkgDir makes any subdir in dir readable and then
// removes the root.
// TODO(bep) move this to a more suitable place.
func MakeReadableAndRemoveAllModulePkgDir(fs afero.Fs, dir string) (int, error) {
	// Safe guard
	// Note that the base directory changed from pkg to gomod_cache in Go 1.23.
	if !strings.Contains(dir, "pkg") && !strings.Contains(dir, "gomod") {
		panic(fmt.Sprint("invalid dir:", dir))
	}

	counter := 0
	afero.Walk(fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			counter++
			fs.Chmod(path, 0o777)
		}
		return nil
	})

	return counter, fs.RemoveAll(dir)
}

// IsOsFs returns whether fs is an OsFs or if it fs wraps an OsFs.
// TODO(bep) make this more robust.
func IsOsFs(fs afero.Fs) bool {
	var isOsFs bool
	WalkFilesystems(fs, func(fs afero.Fs) bool {
		switch fs.(type) {
		case *afero.MemMapFs:
			isOsFs = false
		case *afero.OsFs:
			isOsFs = true
		}
		return isOsFs
	})
	return isOsFs
}

// FilesystemsUnwrapper returns the underlying filesystems.
type FilesystemsUnwrapper interface {
	UnwrapFilesystems() []afero.Fs
}

// FilesystemUnwrapper returns the underlying filesystem.
type FilesystemUnwrapper interface {
	UnwrapFilesystem() afero.Fs
}

// WalkFn is the walk func for WalkFilesystems.
type WalkFn func(fs afero.Fs) bool

// WalkFilesystems walks fs recursively and calls fn.
// If fn returns true, walking is stopped.
func WalkFilesystems(fs afero.Fs, fn WalkFn) bool {
	if fn(fs) {
		return true
	}

	if afs, ok := fs.(FilesystemUnwrapper); ok {
		if WalkFilesystems(afs.UnwrapFilesystem(), fn) {
			return true
		}
	} else if bfs, ok := fs.(FilesystemsUnwrapper); ok {
		for _, sf := range bfs.UnwrapFilesystems() {
			if WalkFilesystems(sf, fn) {
				return true
			}
		}
	} else if cfs, ok := fs.(overlayfs.FilesystemIterator); ok {
		for i := 0; i < cfs.NumFilesystems(); i++ {
			if WalkFilesystems(cfs.Filesystem(i), fn) {
				return true
			}
		}
	}

	return false
}

var _ FilesystemUnwrapper = (*filesystemsWrapper)(nil)

// NewBasePathFs creates a new BasePathFs.
func NewBasePathFs(source afero.Fs, path string) afero.Fs {
	return WrapFilesystem(afero.NewBasePathFs(source, path), source)
}

// NewReadOnlyFs creates a new ReadOnlyFs.
func NewReadOnlyFs(source afero.Fs) afero.Fs {
	return WrapFilesystem(afero.NewReadOnlyFs(source), source)
}

// WrapFilesystem is typically used to wrap a afero.BasePathFs to allow
// access to the underlying filesystem if needed.
func WrapFilesystem(container, content afero.Fs) afero.Fs {
	return filesystemsWrapper{Fs: container, content: content}
}

type filesystemsWrapper struct {
	afero.Fs
	content afero.Fs
}

func (w filesystemsWrapper) UnwrapFilesystem() afero.Fs {
	return w.content
}
