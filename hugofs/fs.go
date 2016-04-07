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

package hugofs

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	sourceFs      afero.Fs
	destinationFs afero.Fs
	osFs          afero.Fs = &afero.OsFs{}
	workingDirFs  *afero.BasePathFs
)

// Source returns Hugo's source file system.
func Source() afero.Fs {
	return sourceFs
}

// SetSource sets Hugo's source file system
// and re-initializes dependent file systems.
func SetSource(fs afero.Fs) {
	sourceFs = fs
	initSourceDependencies()
}

// Destination returns Hugo's destionation file system.
func Destination() afero.Fs {
	return destinationFs
}

// SetDestination sets Hugo's destionation file system
func SetDestination(fs afero.Fs) {
	destinationFs = fs
}

// Os returns an OS file system.
func Os() afero.Fs {
	return osFs
}

// WorkingDir returns a read-only file system
// restricted to the project working dir.
func WorkingDir() *afero.BasePathFs {
	return workingDirFs
}

// InitDefaultFs initializes with the OS file system
// as source and destination file systems.
func InitDefaultFs() {
	InitFs(&afero.OsFs{})
}

// InitMemFs initializes with a MemMapFs as source and destination file systems.
// Useful for testing.
func InitMemFs() {
	InitFs(&afero.MemMapFs{})
}

// InitFs initializes with the given file system
// as source and destination file systems.
func InitFs(fs afero.Fs) {
	sourceFs = fs
	destinationFs = fs

	initSourceDependencies()
}

func initSourceDependencies() {
	workingDir := viper.GetString("WorkingDir")

	if workingDir != "" {
		workingDirFs = afero.NewBasePathFs(afero.NewReadOnlyFs(sourceFs), workingDir).(*afero.BasePathFs)
	}

}

func init() {
	InitDefaultFs()
}
