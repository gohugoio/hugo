// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package os provides template functions for interacting with the operating
// system.
package os

import (
	"errors"
	"fmt"
	_os "os"
	"path/filepath"

	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// New returns a new instance of the os-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	var readFileFs, workFs afero.Fs

	// The docshelper script does not have or need all the dependencies set up.
	if d.PathSpec != nil {
		readFileFs = overlayfs.New(overlayfs.Options{
			Fss: []afero.Fs{
				d.PathSpec.BaseFs.Work,
				d.PathSpec.BaseFs.Content.Fs,
			},
		})
		// See #9599
		workFs = d.PathSpec.BaseFs.WorkDir
	}

	return &Namespace{
		readFileFs: readFileFs,
		workFs:     workFs,
		deps:       d,
	}
}

// Namespace provides template functions for the "os" namespace.
type Namespace struct {
	readFileFs afero.Fs
	workFs     afero.Fs
	deps       *deps.Deps
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func (ns *Namespace) Getenv(key any) (string, error) {
	skey, err := cast.ToStringE(key)
	if err != nil {
		return "", nil
	}

	if err = ns.deps.ExecHelper.Sec().CheckAllowedGetEnv(skey); err != nil {
		return "", err
	}

	return _os.Getenv(skey), nil
}

// readFile reads the file named by filename in the given filesystem
// and returns the contents as a string.
func readFile(fs afero.Fs, filename string) (string, error) {
	filename = filepath.Clean(filename)
	if filename == "" || filename == "." || filename == string(_os.PathSeparator) {
		return "", errors.New("invalid filename")
	}

	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ReadFile reads the file named by filename relative to the configured WorkingDir.
// It returns the contents as a string.
// There is an upper size limit set at 1 megabytes.
func (ns *Namespace) ReadFile(i any) (string, error) {
	s, err := cast.ToStringE(i)
	if err != nil {
		return "", err
	}

	if ns.deps.PathSpec != nil {
		s = ns.deps.PathSpec.RelPathify(s)
	}

	s, err = readFile(ns.readFileFs, s)
	if err != nil && herrors.IsNotExist(err) {
		return "", nil
	}
	return s, err
}

// ReadDir lists the directory contents relative to the configured WorkingDir.
func (ns *Namespace) ReadDir(i any) ([]_os.FileInfo, error) {
	path, err := cast.ToStringE(i)
	if err != nil {
		return nil, err
	}

	list, err := afero.ReadDir(ns.workFs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %s", path, err)
	}

	return list, nil
}

// FileExists checks whether a file exists under the given path.
func (ns *Namespace) FileExists(i any) (bool, error) {
	path, err := cast.ToStringE(i)
	if err != nil {
		return false, err
	}

	if path == "" {
		return false, errors.New("fileExists needs a path to a file")
	}

	status, err := afero.Exists(ns.readFileFs, path)
	if err != nil {
		return false, err
	}

	return status, nil
}

// Stat returns the os.FileInfo structure describing file.
func (ns *Namespace) Stat(i any) (_os.FileInfo, error) {
	path, err := cast.ToStringE(i)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return nil, errors.New("fileStat needs a path to a file")
	}

	r, err := ns.readFileFs.Stat(path)
	if err != nil {
		return nil, err
	}

	return r, nil
}
