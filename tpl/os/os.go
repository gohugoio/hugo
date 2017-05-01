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

package os

import (
	"errors"
	"fmt"
	_os "os"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/deps"
)

// New returns a new instance of the os-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "os" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func (ns *Namespace) Getenv(key interface{}) (string, error) {
	skey, err := cast.ToStringE(key)
	if err != nil {
		return "", nil
	}

	return _os.Getenv(skey), nil
}

// readFile reads the file named by filename relative to the given basepath
// and returns the contents as a string.
// There is a upper size limit set at 1 megabytes.
func readFile(fs *afero.BasePathFs, filename string) (string, error) {
	if filename == "" {
		return "", errors.New("readFile needs a filename")
	}

	if info, err := fs.Stat(filename); err == nil {
		if info.Size() > 1000000 {
			return "", fmt.Errorf("File %q is too big", filename)
		}
	} else {
		return "", err
	}
	b, err := afero.ReadFile(fs, filename)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ReadFilereads the file named by filename relative to the configured
// WorkingDir.  It returns the contents as a string.  There is a upper size
// limit set at 1 megabytes.
func (ns *Namespace) ReadFile(i interface{}) (string, error) {
	s, err := cast.ToStringE(i)
	if err != nil {
		return "", err
	}

	return readFile(ns.deps.Fs.WorkingDir, s)
}

// ReadDir lists the directory contents relative to the configured WorkingDir.
func (ns *Namespace) ReadDir(i interface{}) ([]_os.FileInfo, error) {
	path, err := cast.ToStringE(i)
	if err != nil {
		return nil, err
	}

	list, err := afero.ReadDir(ns.deps.Fs.WorkingDir, path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Directory %s with error message %s", path, err)
	}

	return list, nil
}
