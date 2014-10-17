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
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

var sanitizeRegexp = regexp.MustCompile("[^a-zA-Z0-9./_-]")

// Take a string with any characters and replace it so the string could be used in a path.
// MakePath creates a Unicode sanitized string, with the spaces replaced, whilst
// preserving the original casing of the string.
// E.g. Social Media -> Social-Media
func MakePath(s string) string {
	return UnicodeSanitize(strings.Replace(strings.TrimSpace(s), " ", "-", -1))
}

// MakePathToLowerr creates a Unicode santized string, with the spaces replaced,
// and transformed to lower case.
// E.g. Social Media -> social-media
func MakePathToLower(s string) string {
	return UnicodeSanitize(strings.ToLower(strings.Replace(strings.TrimSpace(s), " ", "-", -1)))
}

func MakeTitle(inpath string) string {
	return strings.Replace(strings.TrimSpace(inpath), "-", " ", -1)
}

func Sanitize(s string) string {
	return sanitizeRegexp.ReplaceAllString(s, "")
}

func UnicodeSanitize(s string) string {
	source := []rune(s)
	target := make([]rune, 0, len(source))

	for _, r := range source {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' || r == '_' || r == '-' {
			target = append(target, r)
		}
	}

	return string(target)
}

func ReplaceExtension(path string, newExt string) string {
	f, _ := FileAndExt(path)
	return f + "." + newExt
}

// Check if Exists && is Directory
func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func IsEmpty(path string) (bool, error) {
	if b, _ := Exists(path); !b {
		return false, fmt.Errorf("%q path does not exist", path)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}
		list, err := f.Readdir(-1)
		f.Close()
		return len(list) == 0, nil
	} else {
		return fi.Size() == 0, nil
	}
}

// Check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
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

	return filepath.Clean(filepath.Join(viper.GetString("WorkingDir"), inPath))
}

func MakeStaticPathRelative(inPath string) (string, error) {
	staticDir := AbsPathify(viper.GetString("StaticDir"))
	themeStaticDir := AbsPathify("themes/"+viper.GetString("theme")) + "/static/"

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

func Filename(in string) (name string) {
	name, _ = FileAndExt(in)
	return
}

func FileAndExt(in string) (name string, ext string) {
	ext = path.Ext(in)
	base := path.Base(in)

	if strings.Contains(base, ".") {
		name = base[:strings.LastIndex(base, ".")]
	} else {
		name = in
	}

	return
}

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
	name = filepath.ToSlash(name)
	return name, nil
}

// Given a source path, determine the section
func GuessSection(in string) string {
	parts := strings.Split(in, "/")

	if len(parts) == 0 {
		return ""
	}

	// trim filename
	if !strings.HasSuffix(in, "/") {
		parts = parts[:len(parts)-1]
	}

	if len(parts) == 0 {
		return ""
	}

	// if first directory is "content", return second directory
	section := ""

	if parts[0] == "content" && len(parts) > 1 {
		section = parts[1]
	} else {
		section = parts[0]
	}

	if section == "." {
		return ""
	}

	return section
}

func PathPrep(ugly bool, in string) string {
	if ugly {
		return Uglify(in)
	} else {
		return PrettifyPath(in)
	}
}

// /section/name.html -> /section/name/index.html
// /section/name/  -> /section/name/index.html
// /section/name/index.html -> /section/name/index.html
func PrettifyPath(in string) string {
	if path.Ext(in) == "" {
		// /section/name/  -> /section/name/index.html
		if len(in) < 2 {
			return "/"
		}
		return path.Join(path.Clean(in), "index.html")
	} else {
		name, ext := FileAndExt(in)
		if name == "index" {
			// /section/name/index.html -> /section/name/index.html
			return path.Clean(in)
		} else {
			// /section/name.html -> /section/name/index.html
			return path.Join(path.Dir(in), name, "index"+ext)
		}
	}
}

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

func SafeWriteToDisk(inpath string, r io.Reader) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			return
		}
	}

	exists, err := Exists(inpath)
	if err != nil {
		return
	}
	if exists {
		return fmt.Errorf("%v already exists", inpath)
	}

	file, err := os.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}

func WriteToDisk(inpath string, r io.Reader) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}
