package source

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Input interface {
	Files() []*File
}

type File struct {
	name        string
	LogicalName string
	Contents    io.Reader
	Section     string
	Dir         string
}

type Filesystem struct {
	files      []*File
	Base       string
	AvoidPaths []string
}

func (f *Filesystem) Files() []*File {
	if len(f.files) < 1 {
		f.captureFiles()
	}
	return f.files
}

var errMissingBaseDir = errors.New("source: missing base directory")

func (f *Filesystem) add(name string, reader io.Reader) (err error) {

	if name, err = f.getRelativePath(name); err != nil {
		return err
	}

	// section should be the first part of the path
	dir, logical := path.Split(name)
	parts := strings.Split(dir, "/")
	section := parts[0]

	if section == "." {
		section = ""
	}

	f.files = append(f.files, &File{
		name:        name,
		LogicalName: logical,
		Contents:    reader,
		Section:     section,
		Dir:         dir,
	})

	return
}

func (f *Filesystem) getRelativePath(name string) (final string, err error) {
	if filepath.IsAbs(name) && f.Base == "" {
		return "", errMissingBaseDir
	}
	name = filepath.Clean(name)
	base := filepath.Clean(f.Base)

	name, err = filepath.Rel(base, name)
	if err != nil {
		return "", err
	}
	name = filepath.ToSlash(name)
	return name, nil
}

func (f *Filesystem) captureFiles() {

	walker := func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if fi.IsDir() {
			if f.avoid(filePath) {
				return filepath.SkipDir
			}
			return nil
		} else {
			if ignoreDotFile(filePath) {
				return nil
			}
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			f.add(filePath, bytes.NewBuffer(data))
			return nil
		}
	}

	filepath.Walk(f.Base, walker)
}

func (f *Filesystem) avoid(filePath string) bool {
	for _, avoid := range f.AvoidPaths {
		if avoid == filePath {
			return true
		}
	}
	return false
}

func ignoreDotFile(filePath string) bool {
	base := filepath.Base(filePath)
	if base[0] == '.' {
		return true
	}

	if base[len(base)-1] == '~' {
		return true
	}

	return false
}
