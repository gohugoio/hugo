// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRootMappingFsRealName(t *testing.T) {
	assert := require.New(t)
	fs := afero.NewMemMapFs()

	rfs, err := NewRootMappingFs(fs, "f1", "f1t", "f2", "f2t")
	assert.NoError(err)

	assert.Equal(filepath.FromSlash("f1t/foo/file.txt"), rfs.realName(filepath.Join("f1", "foo", "file.txt")))

}

func TestRootMappingFsDirnames(t *testing.T) {
	assert := require.New(t)
	fs := afero.NewMemMapFs()

	testfile := "myfile.txt"
	assert.NoError(fs.Mkdir("f1t", 0755))
	assert.NoError(fs.Mkdir("f2t", 0755))
	assert.NoError(fs.Mkdir("f3t", 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("f2t", testfile), []byte("some content"), 0755))

	rfs, err := NewRootMappingFs(fs, "bf1", "f1t", "cf2", "f2t", "af3", "f3t")
	assert.NoError(err)

	fif, err := rfs.Stat(filepath.Join("cf2", testfile))
	assert.NoError(err)
	assert.Equal("myfile.txt", fif.Name())
	assert.Equal(filepath.FromSlash("f2t/myfile.txt"), fif.(RealFilenameInfo).RealFilename())

	root, err := rfs.Open(filepathSeparator)
	assert.NoError(err)

	dirnames, err := root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{"bf1", "cf2", "af3"}, dirnames)

}

func TestRootMappingFsOs(t *testing.T) {
	assert := require.New(t)
	fs := afero.NewOsFs()

	d, err := ioutil.TempDir("", "hugo-root-mapping")
	assert.NoError(err)
	defer func() {
		os.RemoveAll(d)
	}()

	testfile := "myfile.txt"
	assert.NoError(fs.Mkdir(filepath.Join(d, "f1t"), 0755))
	assert.NoError(fs.Mkdir(filepath.Join(d, "f2t"), 0755))
	assert.NoError(fs.Mkdir(filepath.Join(d, "f3t"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join(d, "f2t", testfile), []byte("some content"), 0755))

	rfs, err := NewRootMappingFs(fs, "bf1", filepath.Join(d, "f1t"), "cf2", filepath.Join(d, "f2t"), "af3", filepath.Join(d, "f3t"))
	assert.NoError(err)

	fif, err := rfs.Stat(filepath.Join("cf2", testfile))
	assert.NoError(err)
	assert.Equal("myfile.txt", fif.Name())

	root, err := rfs.Open(filepathSeparator)
	assert.NoError(err)

	dirnames, err := root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{"bf1", "cf2", "af3"}, dirnames)

}
