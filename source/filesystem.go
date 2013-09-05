package source

import (
	"io"
	"os"
	"path/filepath"
)

type Input interface {
	Files() []*File
}

type File struct {
	Name     string
	Contents io.Reader
}

type Filesystem struct {
	files      []*File
	Base       string
	AvoidPaths []string
}

func (f *Filesystem) Files() []*File {
	f.captureFiles()
	return f.files
}

func (f *Filesystem) add(name string, reader io.Reader) {
	f.files = append(f.files, &File{Name: name, Contents: reader})
}

func (f *Filesystem) captureFiles() {

	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if fi.IsDir() {
			if f.avoid(path) {
				return filepath.SkipDir
			}
			return nil
		} else {
			if ignoreDotFile(path) {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			f.add(path, file)
			return nil
		}
	}

	filepath.Walk(f.Base, walker)
}

func (f *Filesystem) avoid(path string) bool {
	for _, avoid := range f.AvoidPaths {
		if avoid == path {
			return true
		}
	}
	return false
}

func ignoreDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}
