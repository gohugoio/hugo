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

package paths

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// filepathPathBridge is a bridge for common functionality in filepath vs path
type filepathPathBridge interface {
	Base(in string) string
	Clean(in string) string
	Dir(in string) string
	Ext(in string) string
	Join(elem ...string) string
	Separator() string
}

type filepathBridge struct{}

func (filepathBridge) Base(in string) string {
	return filepath.Base(in)
}

func (filepathBridge) Clean(in string) string {
	return filepath.Clean(in)
}

func (filepathBridge) Dir(in string) string {
	return filepath.Dir(in)
}

func (filepathBridge) Ext(in string) string {
	return filepath.Ext(in)
}

func (filepathBridge) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (filepathBridge) Separator() string {
	return FilePathSeparator
}

var fpb filepathBridge

// AbsPathify creates an absolute path if given a working dir and a relative path.
// If already absolute, the path is just cleaned.
func AbsPathify(workingDir, inPath string) string {
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}
	return filepath.Join(workingDir, inPath)
}

// MakeTitle converts the path given to a suitable title, trimming whitespace
// and replacing hyphens with whitespace.
func MakeTitle(inpath string) string {
	return strings.Replace(strings.TrimSpace(inpath), "-", " ", -1)
}

// ReplaceExtension takes a path and an extension, strips the old extension
// and returns the path with the new extension.
func ReplaceExtension(path string, newExt string) string {
	f, _ := fileAndExt(path, fpb)
	return f + "." + newExt
}

func makePathRelative(inPath string, possibleDirectories ...string) (string, error) {
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

	dir, _ := filepath.Split(inPath)

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

// ExtNoDelimiter takes a path and returns the extension, excluding the delimiter, i.e. "md".
func ExtNoDelimiter(in string) string {
	return strings.TrimPrefix(Ext(in), ".")
}

// Ext takes a path and returns the extension, including the delimiter, i.e. ".md".
func Ext(in string) string {
	_, ext := fileAndExt(in, fpb)
	return ext
}

// PathAndExt is the same as FileAndExt, but it uses the path package.
func PathAndExt(in string) (string, string) {
	return fileAndExt(in, pb)
}

// FileAndExt takes a path and returns the file and extension separated,
// the extension including the delimiter, i.e. ".md".
func FileAndExt(in string) (string, string) {
	return fileAndExt(in, fpb)
}

// FileAndExtNoDelimiter takes a path and returns the file and extension separated,
// the extension excluding the delimiter, e.g "md".
func FileAndExtNoDelimiter(in string) (string, string) {
	file, ext := fileAndExt(in, fpb)
	return file, strings.TrimPrefix(ext, ".")
}

// Filename takes a file path, strips out the extension,
// and returns the name of the file.
func Filename(in string) (name string) {
	name, _ = fileAndExt(in, fpb)
	return
}

// PathNoExt takes a path, strips out the extension,
// and returns the name of the file.
func PathNoExt(in string) string {
	return strings.TrimSuffix(in, path.Ext(in))
}

// FileAndExt returns the filename and any extension of a file path as
// two separate strings.
//
// If the path, in, contains a directory name ending in a slash,
// then both name and ext will be empty strings.
//
// If the path, in, is either the current directory, the parent
// directory or the root directory, or an empty string,
// then both name and ext will be empty strings.
//
// If the path, in, represents the path of a file without an extension,
// then name will be the name of the file and ext will be an empty string.
//
// If the path, in, represents a filename with an extension,
// then name will be the filename minus any extension - including the dot
// and ext will contain the extension - minus the dot.
func fileAndExt(in string, b filepathPathBridge) (name string, ext string) {
	ext = b.Ext(in)
	base := b.Base(in)

	return extractFilename(in, ext, base, b.Separator()), ext
}

func extractFilename(in, ext, base, pathSeparator string) (name string) {
	// No file name cases. These are defined as:
	// 1. any "in" path that ends in a pathSeparator
	// 2. any "base" consisting of just an pathSeparator
	// 3. any "base" consisting of just an empty string
	// 4. any "base" consisting of just the current directory i.e. "."
	// 5. any "base" consisting of just the parent directory i.e. ".."
	if (strings.LastIndex(in, pathSeparator) == len(in)-1) || base == "" || base == "." || base == ".." || base == pathSeparator {
		name = "" // there is NO filename
	} else if ext != "" { // there was an Extension
		// return the filename minus the extension (and the ".")
		name = base[:strings.LastIndex(base, ".")]
	} else {
		// no extension case so just return base, which will
		// be the filename
		name = base
	}
	return
}

// GetRelativePath returns the relative path of a given path.
func GetRelativePath(path, base string) (final string, err error) {
	if filepath.IsAbs(path) && base == "" {
		return "", errors.New("source: missing base directory")
	}
	name := filepath.Clean(path)
	base = filepath.Clean(base)

	name, err = filepath.Rel(base, name)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(filepath.FromSlash(path), FilePathSeparator) && !strings.HasSuffix(name, FilePathSeparator) {
		name += FilePathSeparator
	}
	return name, nil
}

func prettifyPath(in string, b filepathPathBridge) string {
	if filepath.Ext(in) == "" {
		// /section/name/  -> /section/name/index.html
		if len(in) < 2 {
			return b.Separator()
		}
		return b.Join(in, "index.html")
	}
	name, ext := fileAndExt(in, b)
	if name == "index" {
		// /section/name/index.html -> /section/name/index.html
		return b.Clean(in)
	}
	// /section/name.html -> /section/name/index.html
	return b.Join(b.Dir(in), name, "index"+ext)
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

// DirFile holds the result from path.Split.
type DirFile struct {
	Dir  string
	File string
}

// Used in test.
func (df DirFile) String() string {
	return fmt.Sprintf("%s|%s", df.Dir, df.File)
}
