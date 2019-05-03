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
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func TestNoSymlinkFs(t *testing.T) {
	if skipSymlink() {
		t.Skip("Skip; os.Symlink needs administrator rights on Windows")
	}
	assert := require.New(t)
	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-nosymlink")
	assert.NoError(err)
	defer clean()
	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	blogDir := filepath.Join(workDir, "blog")
	blogFile := filepath.Join(blogDir, "a.txt")
	assert.NoError(os.MkdirAll(blogDir, 0777))
	afero.WriteFile(Os, filepath.Join(blogFile), []byte("content"), 0777)
	os.Chdir(workDir)
	assert.NoError(os.Symlink("blog", "symlinkdedir"))
	os.Chdir(blogDir)
	assert.NoError(os.Symlink("a.txt", "symlinkdedfile.txt"))

	fs := NewNoSymlinkFs(Os)
	ls := fs.(afero.Lstater)
	symlinkedDir := filepath.Join(workDir, "symlinkdedir")
	symlinkedFile := filepath.Join(blogDir, "symlinkdedfile.txt")

	// Check Stat and Lstat
	for _, stat := range []func(name string) (os.FileInfo, error){
		func(name string) (os.FileInfo, error) {
			return fs.Stat(name)
		},
		func(name string) (os.FileInfo, error) {
			fi, _, err := ls.LstatIfPossible(name)
			return fi, err
		},
	} {
		_, err = stat(symlinkedDir)
		assert.Equal(ErrPermissionSymlink, err)
		_, err = stat(symlinkedFile)
		assert.Equal(ErrPermissionSymlink, err)

		fi, err := stat(filepath.Join(workDir, "blog"))
		assert.NoError(err)
		assert.NotNil(fi)

		fi, err = stat(blogFile)
		assert.NoError(err)
		assert.NotNil(fi)
	}

	// Check Open
	_, err = fs.Open(symlinkedDir)
	assert.Equal(ErrPermissionSymlink, err)
	_, err = fs.OpenFile(symlinkedDir, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	assert.Equal(ErrPermissionSymlink, err)
	_, err = fs.OpenFile(symlinkedFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	assert.Equal(ErrPermissionSymlink, err)
	_, err = fs.Open(symlinkedFile)
	assert.Equal(ErrPermissionSymlink, err)
	f, err := fs.Open(blogDir)
	assert.NoError(err)
	f.Close()
	f, err = fs.Open(blogFile)
	assert.NoError(err)
	f.Close()

	// os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

}
