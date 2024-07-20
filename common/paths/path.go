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
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

// FilePathSeparator as defined by os.Separator.
const (
	FilePathSeparator = string(filepath.Separator)
	slash             = "/"
)

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

// AddTrailingSlash adds a trailing Unix styled slash (/) if not already
// there.
func AddTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

// AddLeadingSlash adds a leading Unix styled slash (/) if not already
// there.
func AddLeadingSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// AddTrailingAndLeadingSlash adds a leading and trailing Unix styled slash (/) if not already
// there.
func AddLeadingAndTrailingSlash(path string) string {
	return AddTrailingSlash(AddLeadingSlash(path))
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

// CommonDirPath returns the common directory of the given paths.
func CommonDirPath(path1, path2 string) string {
	if path1 == "" || path2 == "" {
		return ""
	}

	hadLeadingSlash := strings.HasPrefix(path1, "/") || strings.HasPrefix(path2, "/")

	path1 = TrimLeading(path1)
	path2 = TrimLeading(path2)

	p1 := strings.Split(path1, "/")
	p2 := strings.Split(path2, "/")

	var common []string

	for i := 0; i < len(p1) && i < len(p2); i++ {
		if p1[i] == p2[i] {
			common = append(common, p1[i])
		} else {
			break
		}
	}

	s := strings.Join(common, "/")

	if hadLeadingSlash && s != "" {
		s = "/" + s
	}

	return s
}

// Sanitize sanitizes string to be used in Hugo's file paths and URLs, allowing only
// a predefined set of special Unicode characters.
//
// Spaces will be replaced with a single hyphen.
//
// This function is the core function used to normalize paths in Hugo.
//
// Note that this is the first common step for URL/path sanitation,
// the final URL/path may end up looking differently  if the user has stricter rules defined (e.g. removePathAccents=true).
func Sanitize(s string) string {
	var willChange bool
	for i, r := range s {
		willChange = !isAllowedPathCharacter(s, i, r)
		if willChange {
			break
		}
	}

	if !willChange {
		// Prevent allocation when nothing changes.
		return s
	}

	target := make([]rune, 0, len(s))
	var (
		prependHyphen bool
		wasHyphen     bool
	)

	for i, r := range s {
		isAllowed := isAllowedPathCharacter(s, i, r)

		if isAllowed {
			// track explicit hyphen in input; no need to add a new hyphen if
			// we just saw one.
			wasHyphen = r == '-'

			if prependHyphen {
				// if currently have a hyphen, don't prepend an extra one
				if !wasHyphen {
					target = append(target, '-')
				}
				prependHyphen = false
			}
			target = append(target, r)
		} else if len(target) > 0 && !wasHyphen && unicode.IsSpace(r) {
			prependHyphen = true
		}
	}

	return string(target)
}

func isAllowedPathCharacter(s string, i int, r rune) bool {
	if r == ' ' {
		return false
	}
	// Check for the most likely first (faster).
	isAllowed := unicode.IsLetter(r) || unicode.IsDigit(r)
	isAllowed = isAllowed || r == '.' || r == '/' || r == '\\' || r == '_' || r == '#' || r == '+' || r == '~' || r == '-' || r == '@'
	isAllowed = isAllowed || unicode.IsMark(r)
	isAllowed = isAllowed || (r == '%' && i+2 < len(s) && ishex(s[i+1]) && ishex(s[i+2]))
	return isAllowed
}

// From https://golang.org/src/net/url/url.go
func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

var slashFunc = func(r rune) bool {
	return r == '/'
}

// Dir behaves like path.Dir without the path.Clean step.
//
//	The returned path ends in a slash only if it is the root "/".
func Dir(s string) string {
	dir, _ := path.Split(s)
	if len(dir) > 1 && dir[len(dir)-1] == '/' {
		return dir[:len(dir)-1]
	}
	return dir
}

// FieldsSlash cuts s into fields separated with '/'.
func FieldsSlash(s string) []string {
	f := strings.FieldsFunc(s, slashFunc)
	return f
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

// PathEscape escapes unicode letters in pth.
// Use URLEscape to escape full URLs including scheme, query etc.
// This is slightly faster for the common case.
// Note, there is a url.PathEscape function, but that also
// escapes /.
func PathEscape(pth string) string {
	u, err := url.Parse(pth)
	if err != nil {
		panic(err)
	}
	return u.EscapedPath()
}

// ToSlashTrimLeading is just a filepath.ToSlash with an added / prefix trimmer.
func ToSlashTrimLeading(s string) string {
	return TrimLeading(filepath.ToSlash(s))
}

// TrimLeading trims the leading slash from the given string.
func TrimLeading(s string) string {
	return strings.TrimPrefix(s, "/")
}

// ToSlashTrimTrailing is just a filepath.ToSlash with an added / suffix trimmer.
func ToSlashTrimTrailing(s string) string {
	return TrimTrailing(filepath.ToSlash(s))
}

// TrimTrailing trims the trailing slash from the given string.
func TrimTrailing(s string) string {
	return strings.TrimSuffix(s, "/")
}

// ToSlashTrim trims any leading and trailing slashes from the given string and converts it to a forward slash separated path.
func ToSlashTrim(s string) string {
	return strings.Trim(filepath.ToSlash(s), "/")
}

// ToSlashPreserveLeading converts the path given to a forward slash separated path
// and preserves the leading slash if present trimming any trailing slash.
func ToSlashPreserveLeading(s string) string {
	return "/" + strings.Trim(filepath.ToSlash(s), "/")
}

// IsSameFilePath checks if s1 and s2 are the same file path.
func IsSameFilePath(s1, s2 string) bool {
	return path.Clean(ToSlashTrim(s1)) == path.Clean(ToSlashTrim(s2))
}
