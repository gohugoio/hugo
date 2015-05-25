// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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
	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
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

type filepathBridge struct {
}

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
var sanitizeRegexp = regexp.MustCompile("[^a-zA-Z0-9./_-]")

// MakePath takes a string with any characters and replace it
// so the string could be used in a path.
// It does so by creating a Unicode-sanitized string, with the spaces replaced,
// whilst preserving the original casing of the string.
// E.g. Social Media -> Social-Media
func MakePath(s string) string {
	return UnicodeSanitize(strings.Replace(strings.TrimSpace(s), " ", "-", -1))
}

// MakePathToLower creates a Unicode-sanitized string, with the spaces replaced,
// and transformed to lower case.
// E.g. Social Media -> social-media
func MakePathToLower(s string) string {
	return strings.ToLower(MakePath(s))
}

func MakeTitle(inpath string) string {
	return strings.Replace(strings.TrimSpace(inpath), "-", " ", -1)
}

func UnicodeSanitize(s string) string {
	source := []rune(s)
	target := make([]rune, 0, len(source))

	for _, r := range source {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' || r == '_' || r == '-' || r == '#' {
			target = append(target, r)
		}
	}

	return string(target)
}

// ReplaceExtension takes a path and an extension, strips the old extension
// and returns the path with the new extension.
func ReplaceExtension(path string, newExt string) string {
	f, _ := FileAndExt(path, fpb)
	return f + "." + newExt
}

// DirExists checks if a path exists and is a directory.
func DirExists(path string, fs afero.Fs) (bool, error) {
	fi, err := fs.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir checks if a given path is a directory.
func IsDir(path string, fs afero.Fs) (bool, error) {
	fi, err := fs.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// IsEmpty checks if a given path is empty.
func IsEmpty(path string, fs afero.Fs) (bool, error) {
	if b, _ := Exists(path, fs); !b {
		return false, fmt.Errorf("%q path does not exist", path)
	}
	fi, err := fs.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		f, err := os.Open(path)
		// FIX: Resource leak - f.close() should be called here by defer or is missed
		// if the err != nil branch is taken.
		defer f.Close()
		if err != nil {
			return false, err
		}
		list, err := f.Readdir(-1)
		// f.Close() - see bug fix above
		return len(list) == 0, nil
	}
	return fi.Size() == 0, nil
}

// Check if a file contains a specified string.
func FileContains(filename string, subslice []byte, fs afero.Fs) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	return ReaderContains(f, subslice), nil
}

// Check if a file or directory exists.
func Exists(path string, fs afero.Fs) (bool, error) {
	_, err := fs.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func AbsPathify(inPath string) string {
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	// todo consider move workingDir to argument list
	return filepath.Clean(filepath.Join(viper.GetString("WorkingDir"), inPath))
}

func GetStaticDirPath() string {
	return AbsPathify(viper.GetString("StaticDir"))
}

// GetThemeStaticDirPath returns the theme's static dir path if theme is set.
// If theme is set and the static dir doesn't exist, an error is returned.
func GetThemeStaticDirPath() (string, error) {
	return getThemeDirPath("static")
}

// GetThemeStaticDirPath returns the theme's data dir path if theme is set.
// If theme is set and the data dir doesn't exist, an error is returned.
func GetThemeDataDirPath() (string, error) {
	return getThemeDirPath("data")
}

func getThemeDirPath(path string) (string, error) {
	var themeDir string
	if ThemeSet() {
		themeDir = AbsPathify("themes/"+viper.GetString("theme")) + FilePathSeparator + path
		if _, err := os.Stat(themeDir); os.IsNotExist(err) {
			return "", fmt.Errorf("Unable to find %s directory for theme %s in %s", path, viper.GetString("theme"), themeDir)
		}
	}
	return themeDir, nil
}

func GetThemesDirPath() string {
	return AbsPathify(filepath.Join("themes", viper.GetString("theme"), "static"))
}

func MakeStaticPathRelative(inPath string) (string, error) {
	staticDir := GetStaticDirPath()
	themeStaticDir := GetThemesDirPath()

	return MakePathRelative(inPath, staticDir, themeStaticDir)
}

func MakePathRelative(inPath string, possibleDirectories ...string) (string, error) {

	for _, currentPath := range possibleDirectories {
		if strings.HasPrefix(inPath, currentPath) {
			return strings.TrimPrefix(inPath, currentPath), nil
		}
	}
	return inPath, errors.New("Can't extract relative path, unknown prefix")
}

// Should be good enough for Hugo.
var isFileRe = regexp.MustCompile(".*\\..{1,6}$")

// Expects a relative path starting after the content directory.
func GetDottedRelativePath(inPath string) string {
	inPath = filepath.Clean(filepath.FromSlash(inPath))
	if inPath == "." {
		return "./"
	}
	isFile := isFileRe.MatchString(inPath)
	if !isFile {
		if !strings.HasSuffix(inPath, FilePathSeparator) {
			inPath += FilePathSeparator
		}
	}
	if !strings.HasPrefix(inPath, FilePathSeparator) {
		inPath = FilePathSeparator + inPath
	}
	dir, _ := filepath.Split(inPath)

	sectionCount := strings.Count(dir, FilePathSeparator)

	if sectionCount == 0 || dir == FilePathSeparator {
		return "./"
	}

	var dottedPath string

	for i := 1; i < sectionCount; i++ {
		dottedPath += "../"
	}

	return dottedPath
}

// Filename takes a path, strips out the extension,
// and returns the name of the file.
func Filename(in string) (name string) {
	name, _ = FileAndExt(in, fpb)
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
func FileAndExt(in string, b filepathPathBridge) (name string, ext string) {
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
		// no extension case so just return base, which willi
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
	return name, nil
}

func PaginateAliasPath(base string, page int) string {
	paginatePath := viper.GetString("paginatePath")
	uglify := viper.GetBool("UglyURLs")
	var p string
	if base != "" {
		p = filepath.FromSlash(fmt.Sprintf("/%s/%s/%d", base, paginatePath, page))
	} else {
		p = filepath.FromSlash(fmt.Sprintf("/%s/%d", paginatePath, page))
	}
	if uglify {
		p += ".html"
	}

	return p
}

// GuessSection returns the section given a source path.
// A section is the part between the root slash and the second slash
// or before the first slash.
func GuessSection(in string) string {
	parts := strings.Split(in, FilePathSeparator)
	// This will include an empty entry before and after paths with leading and trailing slashes
	// eg... /sect/one/ -> ["", "sect", "one", ""]

	// Needs to have at least a value and a slash
	if len(parts) < 2 {
		return ""
	}

	// If it doesn't have a leading slash and value and file or trailing slash, then return ""
	if parts[0] == "" && len(parts) < 3 {
		return ""
	}

	// strip leading slash
	if parts[0] == "" {
		parts = parts[1:]
	}

	// if first directory is "content", return second directory
	if parts[0] == "content" {
		if len(parts) > 2 {
			return parts[1]
		}
		return ""
	}

	return parts[0]
}

func PathPrep(ugly bool, in string) string {
	if ugly {
		return Uglify(in)
	}
	return PrettifyPath(in)
}

// Same as PrettifyURLPath() but for file paths.
//     /section/name.html       becomes /section/name/index.html
//     /section/name/           becomes /section/name/index.html
//     /section/name/index.html becomes /section/name/index.html
func PrettifyPath(in string) string {
	return PrettiyPath(in, fpb)
}

func PrettiyPath(in string, b filepathPathBridge) string {
	if filepath.Ext(in) == "" {
		// /section/name/  -> /section/name/index.html
		if len(in) < 2 {
			return b.Separator()
		}
		return b.Join(b.Clean(in), "index.html")
	}
	name, ext := FileAndExt(in, b)
	if name == "index" {
		// /section/name/index.html -> /section/name/index.html
		return b.Clean(in)
	}
	// /section/name.html -> /section/name/index.html
	return b.Join(b.Dir(in), name, "index"+ext)
}

// RemoveSubpaths takes a list of paths and removes everything that
// contains another path in the list as a prefix. Ignores any empty
// strings. Used mostly for logging.
//
// e.g. ["hello/world", "hello", "foo/bar", ""] -> ["hello", "foo/bar"]
func RemoveSubpaths(paths []string) []string {
	a := make([]string, 0)
	for _, cur := range paths {
		// ignore trivial case
		if cur == "" {
			continue
		}

		isDupe := false
		for i, old := range a {
			if strings.HasPrefix(cur, old) {
				isDupe = true
				break
			} else if strings.HasPrefix(old, cur) {
				a[i] = cur
				isDupe = true
				break
			}
		}

		if !isDupe {
			a = append(a, cur)
		}
	}

	return a
}

// FindCWD returns the current working directory from where the Hugo
// executable is run.
func FindCWD() (string, error) {
	serverFile, err := filepath.Abs(os.Args[0])

	if err != nil {
		return "", fmt.Errorf("Can't get absolute path for executable: %v", err)
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

// Same as WriteToDisk but checks to see if file/directory already exists.
func SafeWriteToDisk(inpath string, r io.Reader, fs afero.Fs) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = fs.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			return
		}
	}

	exists, err := Exists(inpath, fs)
	if err != nil {
		return
	}
	if exists {
		return fmt.Errorf("%v already exists", inpath)
	}

	file, err := fs.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}

// Writes content to disk.
func WriteToDisk(inpath string, r io.Reader, fs afero.Fs) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = fs.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			if err != os.ErrExist {
				jww.FATAL.Fatalln(err)
			}
		}
	}

	file, err := fs.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}

// GetTempDir returns the OS default temp directory with trailing slash
// if subPath is not empty then it will be created recursively with mode 777 rwx rwx rwx
func GetTempDir(subPath string, fs afero.Fs) string {
	addSlash := func(p string) string {
		if FilePathSeparator != p[len(p)-1:] {
			p = p + FilePathSeparator
		}
		return p
	}
	dir := addSlash(os.TempDir())

	if subPath != "" {
		// preserve windows backslash :-(
		if FilePathSeparator == "\\" {
			subPath = strings.Replace(subPath, "\\", "____", -1)
		}
		dir = dir + MakePath(subPath)
		if FilePathSeparator == "\\" {
			dir = strings.Replace(dir, "____", "\\", -1)
		}

		if exists, _ := Exists(dir, fs); exists {
			return addSlash(dir)
		}

		err := fs.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
		dir = addSlash(dir)
	}
	return dir
}
