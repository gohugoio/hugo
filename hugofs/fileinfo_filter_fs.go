// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"syscall"
	"time"

	"github.com/spf13/afero"
)

// FileInfoFilterFs filters files (not directories) by pred function.
type FileInfoFilterFs struct {
	pred   func(os.FileInfo) bool
	source afero.Fs
}

func NewFileInfoFilterFs(source afero.Fs, pred func(os.FileInfo) bool) afero.Fs {
	return &FileInfoFilterFs{source: source, pred: pred}
}

type FileInfoFilterFile struct {
	f    afero.File
	pred func(os.FileInfo) bool
}

func (fs *FileInfoFilterFs) matches(name string) error {
	fi, err := fs.source.Stat(name)
	if err != nil {
		return err
	}

	if fs.pred(fi) {
		return nil
	}
	return syscall.ENOENT
}

func (fs *FileInfoFilterFs) dirOrMatches(name string) error {
	dir, err := afero.IsDir(fs.source, name)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	return fs.matches(name)
}

func (fs *FileInfoFilterFs) Chtimes(name string, a, m time.Time) error {
	if err := fs.dirOrMatches(name); err != nil {
		return err
	}
	return fs.source.Chtimes(name, a, m)
}

func (fs *FileInfoFilterFs) Chmod(name string, mode os.FileMode) error {
	if err := fs.dirOrMatches(name); err != nil {
		return err
	}
	return fs.source.Chmod(name, mode)
}

func (fs *FileInfoFilterFs) Chown(name string, uid, gid int) error {
	if err := fs.dirOrMatches(name); err != nil {
		return err
	}
	return fs.source.Chown(name, uid, gid)
}

func (fs *FileInfoFilterFs) Name() string {
	return "FileInfoFilterFs"
}

func (fs *FileInfoFilterFs) Stat(name string) (os.FileInfo, error) {
	if err := fs.dirOrMatches(name); err != nil {
		return nil, err
	}
	return fs.source.Stat(name)
}

func (fs *FileInfoFilterFs) Rename(oldname, newname string) error {
	dir, err := afero.IsDir(fs.source, oldname)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	if err := fs.matches(oldname); err != nil {
		return err
	}
	if err := fs.matches(newname); err != nil {
		return err
	}
	return fs.source.Rename(oldname, newname)
}

func (fs *FileInfoFilterFs) RemoveAll(p string) error {
	dir, err := afero.IsDir(fs.source, p)
	if err != nil {
		return err
	}
	if !dir {
		if err := fs.matches(p); err != nil {
			return err
		}
	}
	return fs.source.RemoveAll(p)
}

func (fs *FileInfoFilterFs) Remove(name string) error {
	if err := fs.dirOrMatches(name); err != nil {
		return err
	}
	return fs.source.Remove(name)
}

func (fs *FileInfoFilterFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if err := fs.dirOrMatches(name); err != nil {
		return nil, err
	}
	return fs.source.OpenFile(name, flag, perm)
}

func (fs *FileInfoFilterFs) Open(name string) (afero.File, error) {
	dir, err := afero.IsDir(fs.source, name)
	if err != nil {
		return nil, err
	}
	if !dir {
		if err := fs.matches(name); err != nil {
			return nil, err
		}
	}
	f, err := fs.source.Open(name)
	if err != nil {
		return nil, err
	}
	return &FileInfoFilterFile{f: f, pred: fs.pred}, nil
}

func (fs *FileInfoFilterFs) Mkdir(n string, p os.FileMode) error {
	return fs.source.Mkdir(n, p)
}

func (fs *FileInfoFilterFs) MkdirAll(n string, p os.FileMode) error {
	return fs.source.MkdirAll(n, p)
}

func (fs *FileInfoFilterFs) Create(name string) (afero.File, error) {
	if err := fs.matches(name); err != nil {
		return nil, err
	}
	return fs.source.Create(name)
}

func (f *FileInfoFilterFile) Close() error {
	return f.f.Close()
}

func (f *FileInfoFilterFile) Read(s []byte) (int, error) {
	return f.f.Read(s)
}

func (f *FileInfoFilterFile) ReadAt(s []byte, o int64) (int, error) {
	return f.f.ReadAt(s, o)
}

func (f *FileInfoFilterFile) Seek(o int64, w int) (int64, error) {
	return f.f.Seek(o, w)
}

func (f *FileInfoFilterFile) Write(s []byte) (int, error) {
	return f.f.Write(s)
}

func (f *FileInfoFilterFile) WriteAt(s []byte, o int64) (int, error) {
	return f.f.WriteAt(s, o)
}

func (f *FileInfoFilterFile) Name() string {
	return f.f.Name()
}

func (f *FileInfoFilterFile) Readdir(c int) (fi []os.FileInfo, err error) {
	var rfi []os.FileInfo
	rfi, err = f.f.Readdir(c)
	if err != nil {
		return nil, err
	}
	for _, i := range rfi {
		if i.IsDir() || f.pred(i) {
			fi = append(fi, i)
		}
	}
	return fi, nil
}

func (f *FileInfoFilterFile) Readdirnames(c int) (n []string, err error) {
	fi, err := f.Readdir(c)
	if err != nil {
		return nil, err
	}
	for _, s := range fi {
		n = append(n, s.Name())
	}
	return n, nil
}

func (f *FileInfoFilterFile) Stat() (os.FileInfo, error) {
	return f.f.Stat()
}

func (f *FileInfoFilterFile) Sync() error {
	return f.f.Sync()
}

func (f *FileInfoFilterFile) Truncate(s int64) error {
	return f.f.Truncate(s)
}

func (f *FileInfoFilterFile) WriteString(s string) (int, error) {
	return f.f.WriteString(s)
}
