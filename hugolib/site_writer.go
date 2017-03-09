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
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/output"
	jww "github.com/spf13/jwalterweatherman"
)

// We may find some abstractions/interface(s) here once we star with
// "Multiple Output Types".
type siteWriter struct {
	langDir      string
	publishDir   string
	relativeURLs bool
	uglyURLs     bool
	allowRoot    bool // For aliases

	fs *hugofs.Fs

	log *jww.Notepad
}

func (w siteWriter) targetPathPage(tp output.Type, src string) (string, error) {
	dir, err := w.baseTargetPathPage(tp, src)
	if err != nil {
		return "", err
	}
	if w.publishDir != "" {
		dir = filepath.Join(w.publishDir, dir)
	}
	return dir, nil
}

func (w siteWriter) baseTargetPathPage(tp output.Type, src string) (string, error) {
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

func (w siteWriter) targetPathAlias(src string) (string, error) {
	originalAlias := src
	if len(src) <= 0 {
		return "", fmt.Errorf("Alias \"\" is an empty string")
	}

	alias := filepath.Clean(src)
	components := strings.Split(alias, helpers.FilePathSeparator)

	if !w.allowRoot && alias == helpers.FilePathSeparator {
		return "", fmt.Errorf("Alias \"%s\" resolves to website root directory", originalAlias)
	}

	// Validate against directory traversal
	if components[0] == ".." {
		return "", fmt.Errorf("Alias \"%s\" traverses outside the website root directory", originalAlias)
	}

	// Handle Windows file and directory naming restrictions
	// See "Naming Files, Paths, and Namespaces" on MSDN
	// https://msdn.microsoft.com/en-us/library/aa365247%28v=VS.85%29.aspx?f=255&MSPPError=-2147217396
	msgs := []string{}
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM0", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT0", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	if strings.ContainsAny(alias, ":*?\"<>|") {
		msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains invalid characters on Windows: : * ? \" < > |", originalAlias))
	}
	for _, ch := range alias {
		if ch < ' ' {
			msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains ASCII control code (0x00 to 0x1F), invalid on Windows: : * ? \" < > |", originalAlias))
			continue
		}
	}
	for _, comp := range components {
		if strings.HasSuffix(comp, " ") || strings.HasSuffix(comp, ".") {
			msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains component with a trailing space or period, problematic on Windows", originalAlias))
		}
		for _, r := range reservedNames {
			if comp == r {
				msgs = append(msgs, fmt.Sprintf("Alias \"%s\" contains component with reserved name \"%s\" on Windows", originalAlias, r))
			}
		}
	}
	if len(msgs) > 0 {
		if runtime.GOOS == "windows" {
			for _, m := range msgs {
				w.log.ERROR.Println(m)
			}
			return "", fmt.Errorf("Cannot create \"%s\": Windows filename restriction", originalAlias)
		}
		for _, m := range msgs {
			w.log.WARN.Println(m)
		}
	}

	// Add the final touch
	alias = strings.TrimPrefix(alias, helpers.FilePathSeparator)
	if strings.HasSuffix(alias, helpers.FilePathSeparator) {
		alias = alias + "index.html"
	} else if !strings.HasSuffix(alias, ".html") {
		alias = alias + helpers.FilePathSeparator + "index.html"
	}
	if originalAlias != alias {
		w.log.INFO.Printf("Alias \"%s\" translated to \"%s\"\n", originalAlias, alias)
	}

	return filepath.Join(w.publishDir, alias), nil
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

func (w siteWriter) writeDestPage(tp output.Type, path string, reader io.Reader) error {
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
