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

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func prepareSymlinks(t *testing.T) (string, func()) {
	assert := require.New(t)

	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-symlink-test")
	assert.NoError(err)
	wd, _ := os.Getwd()

	blogDir := filepath.Join(workDir, "blog")
	blogSubDir := filepath.Join(blogDir, "sub")
	assert.NoError(os.MkdirAll(blogSubDir, 0777))
	blogFile1 := filepath.Join(blogDir, "a.txt")
	blogFile2 := filepath.Join(blogSubDir, "b.txt")
	afero.WriteFile(Os, filepath.Join(blogFile1), []byte("content1"), 0777)
	afero.WriteFile(Os, filepath.Join(blogFile2), []byte("content2"), 0777)
	os.Chdir(workDir)
	assert.NoError(os.Symlink("blog", "symlinkdedir"))
	os.Chdir(blogDir)
	assert.NoError(os.Symlink("sub", "symsub"))
	assert.NoError(os.Symlink("a.txt", "symlinkdedfile.txt"))

	return workDir, func() {
		clean()
		os.Chdir(wd)
	}
}

func TestNoSymlinkFs(t *testing.T) {
	if skipSymlink() {
		t.Skip("Skip; os.Symlink needs administrator rights on Windows")
	}
	assert := require.New(t)
	workDir, clean := prepareSymlinks(t)
	defer clean()

	blogDir := filepath.Join(workDir, "blog")
	blogFile1 := filepath.Join(blogDir, "a.txt")

	logger := loggers.NewWarningLogger()

	for _, bfs := range []afero.Fs{NewBaseFileDecorator(Os), Os} {
		for _, allowFiles := range []bool{false, true} {
			logger.WarnCounter.Reset()
			fs := NewNoSymlinkFs(bfs, logger, allowFiles)
			ls := fs.(afero.Lstater)
			symlinkedDir := filepath.Join(workDir, "symlinkdedir")
			symlinkedFilename := "symlinkdedfile.txt"
			symlinkedFile := filepath.Join(blogDir, symlinkedFilename)

			assertFileErr := func(err error) {
				if allowFiles {
					assert.NoError(err)
				} else {
					assert.Equal(ErrPermissionSymlink, err)
				}
			}

			assertFileStat := func(name string, fi os.FileInfo, err error) {
				t.Helper()
				assertFileErr(err)
				if err == nil {
					assert.NotNil(fi)
					assert.Equal(name, fi.Name())
				}
			}

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
				_, err := stat(symlinkedDir)
				assert.Equal(ErrPermissionSymlink, err)
				fi, err := stat(symlinkedFile)
				assertFileStat(symlinkedFilename, fi, err)

				fi, err = stat(filepath.Join(workDir, "blog"))
				assert.NoError(err)
				assert.NotNil(fi)

				fi, err = stat(blogFile1)
				assert.NoError(err)
				assert.NotNil(fi)
			}

			// Check Open
			_, err := fs.Open(symlinkedDir)
			assert.Equal(ErrPermissionSymlink, err)
			_, err = fs.OpenFile(symlinkedDir, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			assert.Equal(ErrPermissionSymlink, err)
			_, err = fs.OpenFile(symlinkedFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			assertFileErr(err)
			_, err = fs.Open(symlinkedFile)
			assertFileErr(err)
			f, err := fs.Open(blogDir)
			assert.NoError(err)
			f.Close()
			f, err = fs.Open(blogFile1)
			assert.NoError(err)
			f.Close()

			// Check readdir
			f, err = fs.Open(workDir)
			assert.NoError(err)
			// There is at least one unsported symlink inside workDir
			_, err = f.Readdir(-1)
			f.Close()
			assert.Equal(uint64(1), logger.WarnCounter.Count())

		}
	}

}
