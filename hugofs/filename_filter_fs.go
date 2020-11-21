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

// FilenameFilterFs filters files (not directories) by pred function.
type FilenameFilterFs struct {
	pred   func(string) bool
	source afero.Fs
}

func NewFilenameFilterFs(source afero.Fs, pred func(string) bool) afero.Fs {
	return &FilenameFilterFs{source: source, pred: pred}
}

type FilenameFilterFile struct {
	f    afero.File
	pred func(string) bool
}

func (fn *FilenameFilterFs) matchesName(name string) error {
	if fn.pred(name) {
		return nil
	}
	return syscall.ENOENT
}

func (fn *FilenameFilterFs) dirOrMatches(name string) error {
	dir, err := afero.IsDir(fn.source, name)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	return fn.matchesName(name)
}

func (fn *FilenameFilterFs) Chtimes(name string, a, m time.Time) error {
	if err := fn.dirOrMatches(name); err != nil {
		return err
	}
	return fn.source.Chtimes(name, a, m)
}

func (fn *FilenameFilterFs) Chmod(name string, mode os.FileMode) error {
	if err := fn.dirOrMatches(name); err != nil {
		return err
	}
	return fn.source.Chmod(name, mode)
}

func (fn *FilenameFilterFs) Chown(name string, uid, gid int) error {
	if err := fn.dirOrMatches(name); err != nil {
		return err
	}
	return fn.source.Chown(name, uid, gid)
}

func (fn *FilenameFilterFs) Name() string {
	return "FilenameFilterFs"
}

func (fn *FilenameFilterFs) Stat(name string) (os.FileInfo, error) {
	if err := fn.dirOrMatches(name); err != nil {
		return nil, err
	}
	return fn.source.Stat(name)
}

func (fn *FilenameFilterFs) Rename(oldname, newname string) error {
	dir, err := afero.IsDir(fn.source, oldname)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	if err := fn.matchesName(oldname); err != nil {
		return err
	}
	if err := fn.matchesName(newname); err != nil {
		return err
	}
	return fn.source.Rename(oldname, newname)
}

func (fn *FilenameFilterFs) RemoveAll(p string) error {
	dir, err := afero.IsDir(fn.source, p)
	if err != nil {
		return err
	}
	if !dir {
		if err := fn.matchesName(p); err != nil {
			return err
		}
	}
	return fn.source.RemoveAll(p)
}

func (fn *FilenameFilterFs) Remove(name string) error {
	if err := fn.dirOrMatches(name); err != nil {
		return err
	}
	return fn.source.Remove(name)
}

func (fn *FilenameFilterFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if err := fn.dirOrMatches(name); err != nil {
		return nil, err
	}
	return fn.source.OpenFile(name, flag, perm)
}

func (fn *FilenameFilterFs) Open(name string) (afero.File, error) {
	dir, err := afero.IsDir(fn.source, name)
	if err != nil {
		return nil, err
	}
	if !dir {
		if err := fn.matchesName(name); err != nil {
			return nil, err
		}
	}
	f, err := fn.source.Open(name)
	if err != nil {
		return nil, err
	}
	return &FilenameFilterFile{f: f, pred: fn.pred}, nil
}

func (fn *FilenameFilterFs) Mkdir(n string, p os.FileMode) error {
	return fn.source.Mkdir(n, p)
}

func (fn *FilenameFilterFs) MkdirAll(n string, p os.FileMode) error {
	return fn.source.MkdirAll(n, p)
}

func (fn *FilenameFilterFs) Create(name string) (afero.File, error) {
	if err := fn.matchesName(name); err != nil {
		return nil, err
	}
	return fn.source.Create(name)
}

func (f *FilenameFilterFile) Close() error {
	return f.f.Close()
}

func (f *FilenameFilterFile) Read(s []byte) (int, error) {
	return f.f.Read(s)
}

func (f *FilenameFilterFile) ReadAt(s []byte, o int64) (int, error) {
	return f.f.ReadAt(s, o)
}

func (f *FilenameFilterFile) Seek(o int64, w int) (int64, error) {
	return f.f.Seek(o, w)
}

func (f *FilenameFilterFile) Write(s []byte) (int, error) {
	return f.f.Write(s)
}

func (f *FilenameFilterFile) WriteAt(s []byte, o int64) (int, error) {
	return f.f.WriteAt(s, o)
}

func (f *FilenameFilterFile) Name() string {
	return f.f.Name()
}

func (f *FilenameFilterFile) Readdir(c int) (fi []os.FileInfo, err error) {
	var rfi []os.FileInfo
	rfi, err = f.f.Readdir(c)
	if err != nil {
		return nil, err
	}
	for _, i := range rfi {
		if i.IsDir() || f.pred(i.Name()) {
			fi = append(fi, i)
		}
	}
	return fi, nil
}

func (f *FilenameFilterFile) Readdirnames(c int) (n []string, err error) {
	fi, err := f.Readdir(c)
	if err != nil {
		return nil, err
	}
	for _, s := range fi {
		n = append(n, s.Name())
	}
	return n, nil
}

func (f *FilenameFilterFile) Stat() (os.FileInfo, error) {
	return f.f.Stat()
}

func (f *FilenameFilterFile) Sync() error {
	return f.f.Sync()
}

func (f *FilenameFilterFile) Truncate(s int64) error {
	return f.f.Truncate(s)
}

func (f *FilenameFilterFile) WriteString(s string) (int, error) {
	return f.f.WriteString(s)
}
