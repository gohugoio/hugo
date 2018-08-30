// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
)

// Os points to an Os Afero file system.
var Os = &afero.OsFs{}

// Fs abstracts the file system to separate source and destination file systems
// and allows both to be mocked for testing.
type Fs struct {
	// Source is Hugo's source file system.
	Source afero.Fs

	// Destination is Hugo's destination file system.
	Destination afero.Fs

	// Os is an OS file system.
	// NOTE: Field is currently unused.
	Os afero.Fs

	// WorkingDir is a read-only file system
	// restricted to the project working dir.
	WorkingDir *afero.BasePathFs
}

// NewDefault creates a new Fs with the OS file system
// as source and destination file systems.
func NewDefault(cfg config.Provider) *Fs {
	fs := &afero.OsFs{}
	return newFs(fs, cfg)
}

// NewMem creates a new Fs with the MemMapFs
// as source and destination file systems.
// Useful for testing.
func NewMem(cfg config.Provider) *Fs {
	fs := &afero.MemMapFs{}
	return newFs(fs, cfg)
}

// NewFrom creates a new Fs based on the provided Afero Fs
// as source and destination file systems.
// Useful for testing.
func NewFrom(fs afero.Fs, cfg config.Provider) *Fs {
	return newFs(fs, cfg)
}

func newFs(base afero.Fs, cfg config.Provider) *Fs {
	return &Fs{
		Source:      base,
		Destination: base,
		Os:          &afero.OsFs{},
		WorkingDir:  getWorkingDirFs(base, cfg),
	}
}

func getWorkingDirFs(base afero.Fs, cfg config.Provider) *afero.BasePathFs {
	workingDir := cfg.GetString("workingDir")

	if workingDir != "" {
		return afero.NewBasePathFs(afero.NewReadOnlyFs(base), workingDir).(*afero.BasePathFs)
	}

	return nil
}
