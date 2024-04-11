// Copyright 2024 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/spf13/afero"
)

// MakePath takes a string with any characters and replace it
// so the string could be used in a path.
// It does so by creating a Unicode-sanitized string, with the spaces replaced,
// whilst preserving the original casing of the string.
// E.g. Social Media -> Social-Media
func (p *PathSpec) MakePath(s string) string {
	s = paths.Sanitize(s)
	if p.Cfg.RemovePathAccents() {
		s = text.RemoveAccentsString(s)
	}
	return s
}

// MakePathsSanitized applies MakePathSanitized on every item in the slice
func (p *PathSpec) MakePathsSanitized(paths []string) {
	for i, path := range paths {
		paths[i] = p.MakePathSanitized(path)
	}
}

// MakePathSanitized creates a Unicode-sanitized string, with the spaces replaced
func (p *PathSpec) MakePathSanitized(s string) string {
	if p.Cfg.DisablePathToLower() {
		return p.MakePath(s)
	}
	return strings.ToLower(p.MakePath(s))
}

// MakeTitle converts the path given to a suitable title, trimming whitespace
// and replacing hyphens with whitespace.
func MakeTitle(inpath string) string {
	return strings.Replace(strings.TrimSpace(inpath), "-", " ", -1)
}

// MakeTitleInPath converts the path given to a suitable title, trimming whitespace
func MakePathRelative(inPath string, possibleDirectories ...string) (string, error) {
	for _, currentPath := range possibleDirectories {
		if strings.HasPrefix(inPath, currentPath) {
			return strings.TrimPrefix(inPath, currentPath), nil
		}
	}
	return inPath, errors.New("can't extract relative path, unknown prefix")
}

// Should be good enough for Hugo.
var isFileRe = regexp.MustCompile(`.*\..{1,6}$`)

// GetDottedRelativePath expects a relative path starting after the content directory.
// It returns a relative path with dots ("..") navigating up the path structure.
func GetDottedRelativePath(inPath string) string {
	inPath = path.Clean(filepath.ToSlash(inPath))

	if inPath == "." {
		return "./"
	}

	if !isFileRe.MatchString(inPath) && !strings.HasSuffix(inPath, "/") {
		inPath += "/"
	}

	if !strings.HasPrefix(inPath, "/") {
		inPath = "/" + inPath
	}

	dir, _ := path.Split(inPath)

	sectionCount := strings.Count(dir, "/")

	if sectionCount == 0 || dir == "/" {
		return "./"
	}

	var dottedPath string

	for i := 1; i < sectionCount; i++ {
		dottedPath += "../"
	}

	return dottedPath
}

type NamedSlice struct {
	Name  string
	Slice []string
}

func (n NamedSlice) String() string {
	if len(n.Slice) == 0 {
		return n.Name
	}
	return fmt.Sprintf("%s%s{%s}", n.Name, FilePathSeparator, strings.Join(n.Slice, ","))
}

func ExtractAndGroupRootPaths(paths []string) []NamedSlice {
	if len(paths) == 0 {
		return nil
	}

	pathsCopy := make([]string, len(paths))
	hadSlashPrefix := strings.HasPrefix(paths[0], FilePathSeparator)

	for i, p := range paths {
		pathsCopy[i] = strings.Trim(filepath.ToSlash(p), "/")
	}

	sort.Strings(pathsCopy)

	pathsParts := make([][]string, len(pathsCopy))

	for i, p := range pathsCopy {
		pathsParts[i] = strings.Split(p, "/")
	}

	var groups [][]string

	for i, p1 := range pathsParts {
		c1 := -1

		for j, p2 := range pathsParts {
			if i == j {
				continue
			}

			c2 := -1

			for i, v := range p1 {
				if i >= len(p2) {
					break
				}
				if v != p2[i] {
					break
				}

				c2 = i
			}

			if c1 == -1 || (c2 != -1 && c2 < c1) {
				c1 = c2
			}
		}

		if c1 != -1 {
			groups = append(groups, p1[:c1+1])
		} else {
			groups = append(groups, p1)
		}
	}

	groupsStr := make([]string, len(groups))
	for i, g := range groups {
		groupsStr[i] = strings.Join(g, "/")
	}

	groupsStr = UniqueStringsSorted(groupsStr)

	var result []NamedSlice

	for _, g := range groupsStr {
		name := filepath.FromSlash(g)
		if hadSlashPrefix {
			name = FilePathSeparator + name
		}
		ns := NamedSlice{Name: name}
		for _, p := range pathsCopy {
			if !strings.HasPrefix(p, g) {
				continue
			}

			p = strings.TrimPrefix(p, g)
			if p != "" {
				ns.Slice = append(ns.Slice, p)
			}
		}

		ns.Slice = UniqueStrings(ExtractRootPaths(ns.Slice))

		result = append(result, ns)
	}

	return result
}

// ExtractRootPaths extracts the root paths from the supplied list of paths.
// The resulting root path will not contain any file separators, but there
// may be duplicates.
// So "/content/section/" becomes "content"
func ExtractRootPaths(paths []string) []string {
	r := make([]string, len(paths))
	for i, p := range paths {
		root := filepath.ToSlash(p)
		sections := strings.Split(root, "/")
		for _, section := range sections {
			if section != "" {
				root = section
				break
			}
		}
		r[i] = root
	}
	return r
}

// FindCWD returns the current working directory from where the Hugo
// executable is run.
func FindCWD() (string, error) {
	serverFile, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", fmt.Errorf("can't get absolute path for executable: %v", err)
	}

	path := filepath.Dir(serverFile)
	realFile, err := filepath.EvalSymlinks(serverFile)
	if err != nil {
		if _, err = os.Stat(serverFile + ".exe"); err == nil {
			realFile = filepath.Clean(serverFile + ".exe")
		}
	}

	if err == nil && realFile != serverFile {
		path = filepath.Dir(realFile)
	}

	return path, nil
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root.
func Walk(fs afero.Fs, root string, walker hugofs.WalkFunc) error {
	if _, isOs := fs.(*afero.OsFs); isOs {
		fs = hugofs.NewBaseFileDecorator(fs)
	}
	w := hugofs.NewWalkway(hugofs.WalkwayConfig{
		Fs:     fs,
		Root:   root,
		WalkFn: walker,
	})

	return w.Walk()
}

// SafeWriteToDisk is the same as WriteToDisk
// but it also checks to see if file/directory already exists.
func SafeWriteToDisk(inpath string, r io.Reader, fs afero.Fs) (err error) {
	return afero.SafeWriteReader(fs, inpath, r)
}

// WriteToDisk writes content to disk.
func WriteToDisk(inpath string, r io.Reader, fs afero.Fs) (err error) {
	return afero.WriteReader(fs, inpath, r)
}

// OpenFilesForWriting opens all the given filenames for writing.
func OpenFilesForWriting(fs afero.Fs, filenames ...string) (io.WriteCloser, error) {
	var writeClosers []io.WriteCloser
	for _, filename := range filenames {
		f, err := OpenFileForWriting(fs, filename)
		if err != nil {
			for _, wc := range writeClosers {
				wc.Close()
			}
			return nil, err
		}
		writeClosers = append(writeClosers, f)
	}

	return hugio.NewMultiWriteCloser(writeClosers...), nil
}

// OpenFileForWriting opens or creates the given file. If the target directory
// does not exist, it gets created.
func OpenFileForWriting(fs afero.Fs, filename string) (afero.File, error) {
	filename = filepath.Clean(filename)
	// Create will truncate if file already exists.
	// os.Create will create any new files with mode 0666 (before umask).
	f, err := fs.Create(filename)
	if err != nil {
		if !herrors.IsNotExist(err) {
			return nil, err
		}
		if err = fs.MkdirAll(filepath.Dir(filename), 0o777); err != nil { //  before umask
			return nil, err
		}
		f, err = fs.Create(filename)
	}

	return f, err
}

// GetCacheDir returns a cache dir from the given filesystem and config.
// The dir will be created if it does not exist.
func GetCacheDir(fs afero.Fs, cacheDir string) (string, error) {
	cacheDir = cacheDirDefault(cacheDir)

	if cacheDir != "" {
		exists, err := DirExists(cacheDir, fs)
		if err != nil {
			return "", err
		}
		if !exists {
			err := fs.MkdirAll(cacheDir, 0o777) // Before umask
			if err != nil {
				return "", fmt.Errorf("failed to create cache dir: %w", err)
			}
		}
		return cacheDir, nil
	}

	const hugoCacheBase = "hugo_cache"

	// Avoid filling up the home dir with Hugo cache dirs from development.
	if !htesting.IsTest {
		userCacheDir, err := os.UserCacheDir()
		if err == nil {
			cacheDir := filepath.Join(userCacheDir, hugoCacheBase)
			if err := fs.Mkdir(cacheDir, 0o777); err == nil || os.IsExist(err) {
				return cacheDir, nil
			}
		}
	}

	// Fall back to a cache in /tmp.
	userName := os.Getenv("USER")
	if userName != "" {
		return GetTempDir(hugoCacheBase+"_"+userName, fs), nil
	} else {
		return GetTempDir(hugoCacheBase, fs), nil
	}
}

func cacheDirDefault(cacheDir string) string {
	// Always use the cacheDir config if set.
	if len(cacheDir) > 1 {
		return addTrailingFileSeparator(cacheDir)
	}

	// See Issue #8714.
	// Turns out that Cloudflare also sets NETLIFY=true in its build environment,
	// but all of these 3 should not give any false positives.
	if os.Getenv("NETLIFY") == "true" && os.Getenv("PULL_REQUEST") != "" && os.Getenv("DEPLOY_PRIME_URL") != "" {
		// Netlify's cache behavior is not documented, the currently best example
		// is this project:
		// https://github.com/philhawksworth/content-shards/blob/master/gulpfile.js
		return "/opt/build/cache/hugo_cache/"
	}

	// This will fall back to an hugo_cache folder in either os.UserCacheDir or the tmp dir, which should work fine for most CI
	// providers. See this for a working CircleCI setup:
	// https://github.com/bep/hugo-sass-test/blob/6c3960a8f4b90e8938228688bc49bdcdd6b2d99e/.circleci/config.yml
	// If not, they can set the HUGO_CACHEDIR environment variable or cacheDir config key.
	return ""
}

func addTrailingFileSeparator(s string) string {
	if !strings.HasSuffix(s, FilePathSeparator) {
		s = s + FilePathSeparator
	}
	return s
}

// GetTempDir returns a temporary directory with the given sub path.
func GetTempDir(subPath string, fs afero.Fs) string {
	return afero.GetTempDir(fs, subPath)
}

// DirExists checks if a path exists and is a directory.
func DirExists(path string, fs afero.Fs) (bool, error) {
	return afero.DirExists(fs, path)
}

// IsDir checks if a given path is a directory.
func IsDir(path string, fs afero.Fs) (bool, error) {
	return afero.IsDir(fs, path)
}

// IsEmpty checks if a given path is empty, meaning it doesn't contain any regular files.
func IsEmpty(path string, fs afero.Fs) (bool, error) {
	var hasFile bool
	err := afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		hasFile = true
		return filepath.SkipDir
	})
	return !hasFile, err
}

// Exists checks if a file or directory exists.
func Exists(path string, fs afero.Fs) (bool, error) {
	return afero.Exists(fs, path)
}
