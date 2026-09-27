// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"fmt"
	"path/filepath"

	"github.com/bep/overlayfs"
	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/spf13/afero"
	"golang.org/x/tools/txtar"
)

// TxtarFilename is the name of the txtar archive in the project root.
const TxtarFilename = "hugo.txtar"

// TxtarFs overlays the files in a txtar archive on top of base.
// Files in base take precedence, and all writes go to base.
type TxtarFs struct {
	*overlayfs.OverlayFs

	base     afero.Fs
	mem      afero.Fs
	filename string
	dir      string

	// Content hashes of the files currently in the archive, keyed by absolute filename.
	hashes map[string]uint64
}

// NewTxtarFsIfExists wraps base in a TxtarFs if workingDir contains a hugo.txtar file.
// It returns base unchanged if not.
func NewTxtarFsIfExists(base afero.Fs, workingDir string) (afero.Fs, error) {
	filename := filepath.Join(workingDir, TxtarFilename)
	if _, err := base.Stat(filename); err != nil {
		if herrors.IsNotExist(err) {
			return base, nil
		}
		return nil, err
	}

	mem := afero.NewMemMapFs()
	fs := &TxtarFs{
		OverlayFs: overlayfs.New(overlayfs.Options{Fss: []afero.Fs{base, mem}, FirstWritable: true}),
		base:      base,
		mem:       mem,
		filename:  filename,
		dir:       filepath.Clean(workingDir),
		hashes:    make(map[string]uint64),
	}

	if _, err := fs.Reload(); err != nil {
		return nil, err
	}

	return fs, nil
}

// Filename returns the absolute filename of the archive.
func (fs *TxtarFs) Filename() string {
	return fs.filename
}

// Reload re-reads the archive and applies any changes to the in-memory layer.
// It returns one event per added, changed or removed file.
// A missing archive is treated as empty.
func (fs *TxtarFs) Reload() ([]fsnotify.Event, error) {
	var arc *txtar.Archive
	fi, err := fs.base.Stat(fs.filename)
	if err == nil {
		var b []byte
		b, err = afero.ReadFile(fs.base, fs.filename)
		arc = txtar.Parse(b)
	}
	if err != nil {
		if !herrors.IsNotExist(err) {
			return nil, err
		}
		arc = &txtar.Archive{}
	}

	type write struct {
		name string
		data []byte
		op   fsnotify.Op
	}
	var writes []write
	hashes := make(map[string]uint64, len(arc.Files))

	for _, f := range arc.Files {
		name := filepath.FromSlash(f.Name)
		if !filepath.IsLocal(name) {
			return nil, fmt.Errorf("%s: invalid filename %q", fs.filename, f.Name)
		}
		name = filepath.Join(fs.dir, name)
		h := xxhash.Sum64(f.Data)
		hashes[name] = h
		old, seen := fs.hashes[name]
		if seen && old == h {
			continue
		}
		op := fsnotify.Write
		if !seen {
			op = fsnotify.Create
		}
		writes = append(writes, write{name, f.Data, op})
	}

	var events []fsnotify.Event

	// Remove before write to handle files replaced by directories and vice versa.
	for name := range fs.hashes {
		if _, ok := hashes[name]; ok {
			continue
		}
		if err := fs.remove(name); err != nil {
			return nil, err
		}
		events = append(events, fsnotify.Event{Name: name, Op: fsnotify.Remove})
	}

	for _, w := range writes {
		if err := fs.mem.MkdirAll(filepath.Dir(w.name), 0o777); err != nil {
			return nil, err
		}
		if err := afero.WriteFile(fs.mem, w.name, w.data, 0o666); err != nil {
			return nil, err
		}
		if fi != nil {
			fs.mem.Chtimes(w.name, fi.ModTime(), fi.ModTime())
		}
		events = append(events, fsnotify.Event{Name: w.name, Op: w.op})
	}

	fs.hashes = hashes

	return events, nil
}

// remove removes name and any directories left empty by it, up to the project root.
func (fs *TxtarFs) remove(name string) error {
	if err := fs.mem.Remove(name); err != nil {
		return err
	}
	for dir := filepath.Dir(name); dir != fs.dir; dir = filepath.Dir(dir) {
		if empty, _ := afero.IsEmpty(fs.mem, dir); !empty {
			break
		}
		if err := fs.mem.Remove(dir); err != nil {
			return err
		}
	}
	return nil
}
