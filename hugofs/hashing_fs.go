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
	"crypto/md5"
	"encoding/hex"
	"hash"
	"os"

	"github.com/spf13/afero"
)

var (
	_ afero.Fs = (*md5HashingFs)(nil)
)

// FileHashReceiver will receive the filename an the content's MD5 sum on file close.
type FileHashReceiver interface {
	OnFileClose(name, md5sum string)
}

type md5HashingFs struct {
	afero.Fs
	hashReceiver FileHashReceiver
}

// NewHashingFs creates a new filesystem that will receive MD5 checksums of
// any written file content on Close. Note that this is probably not a good
// idea for "full build" situations, but when doing fast render mode, the amount
// of files published is low, and it would be really nice to know exactly which
// of these files where actually changed.
// Note that this will only work for file operations that use the io.Writer
// to write content to file, but that is fine for the "publish content" use case.
func NewHashingFs(delegate afero.Fs, hashReceiver FileHashReceiver) afero.Fs {
	return &md5HashingFs{Fs: delegate, hashReceiver: hashReceiver}
}

func (fs *md5HashingFs) Create(name string) (afero.File, error) {
	f, err := fs.Fs.Create(name)
	if err == nil {
		f = fs.wrapFile(f)
	}
	return f, err
}

func (fs *md5HashingFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil && isWrite(flag) {
		f = fs.wrapFile(f)
	}
	return f, err
}

func (fs *md5HashingFs) wrapFile(f afero.File) afero.File {
	return &hashingFile{File: f, h: md5.New(), hashReceiver: fs.hashReceiver}
}

func isWrite(flag int) bool {
	return flag&os.O_RDWR != 0 || flag&os.O_WRONLY != 0
}

func (fs *md5HashingFs) Name() string {
	return "md5HashingFs"
}

type hashingFile struct {
	hashReceiver FileHashReceiver
	h            hash.Hash
	afero.File
}

func (h *hashingFile) Write(p []byte) (n int, err error) {
	n, err = h.File.Write(p)
	if err != nil {
		return
	}
	return h.h.Write(p)
}

func (h *hashingFile) Close() error {
	sum := hex.EncodeToString(h.h.Sum(nil))
	h.hashReceiver.OnFileClose(h.Name(), sum)
	return h.File.Close()
}
