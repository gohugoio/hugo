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

package hugolib

import (
	"io"
	"path/filepath"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/output"
	jww "github.com/spf13/jwalterweatherman"
)

// We may find some abstractions/interface(s) here once we star with
// "Multiple Output Formats".
type siteWriter struct {
	langDir      string
	publishDir   string
	relativeURLs bool
	uglyURLs     bool
	allowRoot    bool // For aliases

	fs *hugofs.Fs

	log *jww.Notepad
}

func (w siteWriter) targetPathPage(f output.Format, src string) (string, error) {
	dir, err := w.baseTargetPathPage(f, src)
	if err != nil {
		return "", err
	}
	if w.publishDir != "" {
		dir = filepath.Join(w.publishDir, dir)
	}
	return dir, nil
}

func (w siteWriter) baseTargetPathPage(f output.Format, src string) (string, error) {
	if src == helpers.FilePathSeparator {
		return "index.html", nil
	}

	// The anatomy of a target path:
	// langDir
	// BaseName
	// Suffix
	// ROOT?
	// dir
	// name

	dir, file := filepath.Split(src)
	isRoot := dir == ""
	ext := extension(filepath.Ext(file))
	name := filename(file)

	if w.langDir != "" && dir == helpers.FilePathSeparator && name == w.langDir {
		return filepath.Join(dir, name, "index"+ext), nil
	}

	if w.uglyURLs || file == "index.html" || (isRoot && file == "404.html") {
		return filepath.Join(dir, name+ext), nil
	}

	dir = filepath.Join(dir, name, "index"+ext)

	return dir, nil

}

func (w siteWriter) targetPathFile(src string) (string, error) {
	return filepath.Join(w.publishDir, filepath.FromSlash(src)), nil
}

func extension(ext string) string {
	switch ext {
	case ".md", ".rst":
		return ".html"
	}

	if ext != "" {
		return ext
	}

	return ".html"
}

func filename(f string) string {
	ext := filepath.Ext(f)
	if ext == "" {
		return f
	}

	return f[:len(f)-len(ext)]
}

func (w siteWriter) writeDestPage(f output.Format, path string, reader io.Reader) error {
	w.log.DEBUG.Println("creating page:", path)
	path, _ = w.targetPathFile(path)
	// TODO(bep) output remove this file ... targetPath, err := w.targetPathPage(tp, path)

	return w.publish(path, reader)
}

func (w siteWriter) writeDestFile(path string, r io.Reader) (err error) {
	w.log.DEBUG.Println("creating file:", path)
	targetPath, err := w.targetPathFile(path)
	if err != nil {
		return err
	}
	return w.publish(targetPath, r)
}

func (w siteWriter) publish(path string, r io.Reader) (err error) {

	return helpers.WriteToDisk(path, r, w.fs.Destination)
}
