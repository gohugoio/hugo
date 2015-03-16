package hugofs

import (
	//"fmt"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PathCleaningOsFs is a Fs implementation that uses functions provided by the os package,
// with all paths run through filepath.Clean. This makes it suitable for cross-platform usage.
//
// For details in any method, check the documentation of the os package
// (http://golang.org/pkg/os/).
type PathCleaningOsFs struct{}

func (PathCleaningOsFs) Name() string { return "PathCleaningOsFs" }

func (PathCleaningOsFs) Create(name string) (afero.File, error) {
	return os.Create(name)
}

func (PathCleaningOsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (PathCleaningOsFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (PathCleaningOsFs) Open(name string) (afero.File, error) {
	return os.Open(name)
}

func (PathCleaningOsFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (PathCleaningOsFs) Remove(name string) error {
	return os.Remove(name)
}

func (PathCleaningOsFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (PathCleaningOsFs) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (PathCleaningOsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (PathCleaningOsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (PathCleaningOsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func clean(path string) string {
	s := path
	first := s[0]
	if !os.IsPathSeparator(first) {
		s = strings.Replace(s, string(first), string(os.PathSeparator), -1)
	}

	return filepath.Clean(s)
}
