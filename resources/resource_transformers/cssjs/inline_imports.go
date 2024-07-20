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

package cssjs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/spf13/afero"
)

const importIdentifier = "@import"

var (
	cssSyntaxErrorRe = regexp.MustCompile(`> (\d+) \|`)
	shouldImportRe   = regexp.MustCompile(`^@import ["'](.*?)["'];?\s*(/\*.*\*/)?$`)
)

type fileOffset struct {
	Filename string
	Offset   int
}

type importResolver struct {
	r      io.Reader
	inPath string
	opts   InlineImports

	contentSeen       map[string]bool
	dependencyManager identity.Manager
	linemap           map[int]fileOffset
	fs                afero.Fs
	logger            loggers.Logger
}

func newImportResolver(r io.Reader, inPath string, opts InlineImports, fs afero.Fs, logger loggers.Logger, dependencyManager identity.Manager) *importResolver {
	return &importResolver{
		r:                 r,
		dependencyManager: dependencyManager,
		inPath:            inPath,
		fs:                fs, logger: logger,
		linemap: make(map[int]fileOffset), contentSeen: make(map[string]bool),
		opts: opts,
	}
}

func (imp *importResolver) contentHash(filename string) ([]byte, string) {
	b, err := afero.ReadFile(imp.fs, filename)
	if err != nil {
		return nil, ""
	}
	h := sha256.New()
	h.Write(b)
	return b, hex.EncodeToString(h.Sum(nil))
}

func (imp *importResolver) importRecursive(
	lineNum int,
	content string,
	inPath string,
) (int, string, error) {
	basePath := path.Dir(inPath)

	var replacements []string
	lines := strings.Split(content, "\n")

	trackLine := func(i, offset int, line string) {
		// TODO(bep) this is not very efficient.
		imp.linemap[i+lineNum] = fileOffset{Filename: inPath, Offset: offset}
	}

	i := 0
	for offset, line := range lines {
		i++
		lineTrimmed := strings.TrimSpace(line)
		column := strings.Index(line, lineTrimmed)
		line = lineTrimmed

		if !imp.shouldImport(line) {
			trackLine(i, offset, line)
		} else {
			path := strings.Trim(strings.TrimPrefix(line, importIdentifier), " \"';")
			filename := filepath.Join(basePath, path)
			imp.dependencyManager.AddIdentity(identity.CleanStringIdentity(filename))
			importContent, hash := imp.contentHash(filename)

			if importContent == nil {
				if imp.opts.SkipInlineImportsNotFound {
					trackLine(i, offset, line)
					continue
				}
				pos := text.Position{
					Filename:     inPath,
					LineNumber:   offset + 1,
					ColumnNumber: column + 1,
				}
				return 0, "", herrors.NewFileErrorFromFileInPos(fmt.Errorf("failed to resolve CSS @import \"%s\"", filename), pos, imp.fs, nil)
			}

			i--

			if imp.contentSeen[hash] {
				i++
				// Just replace the line with an empty string.
				replacements = append(replacements, []string{line, ""}...)
				trackLine(i, offset, "IMPORT")
				continue
			}

			imp.contentSeen[hash] = true

			// Handle recursive imports.
			l, nested, err := imp.importRecursive(i+lineNum, string(importContent), filepath.ToSlash(filename))
			if err != nil {
				return 0, "", err
			}

			trackLine(i, offset, line)

			i += l

			importContent = []byte(nested)

			replacements = append(replacements, []string{line, string(importContent)}...)
		}
	}

	if len(replacements) > 0 {
		repl := strings.NewReplacer(replacements...)
		content = repl.Replace(content)
	}

	return i, content, nil
}

func (imp *importResolver) resolve() (io.Reader, error) {
	content, err := io.ReadAll(imp.r)
	if err != nil {
		return nil, err
	}

	contents := string(content)

	_, newContent, err := imp.importRecursive(0, contents, imp.inPath)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(newContent), nil
}

// See https://www.w3schools.com/cssref/pr_import_rule.asp
// We currently only support simple file imports, no urls, no media queries.
// So this is OK:
//
//	@import "navigation.css";
//
// This is not:
//
//	@import url("navigation.css");
//	@import "mobstyle.css" screen and (max-width: 768px);
func (imp *importResolver) shouldImport(s string) bool {
	if !strings.HasPrefix(s, importIdentifier) {
		return false
	}
	if strings.Contains(s, "url(") {
		return false
	}

	m := shouldImportRe.FindStringSubmatch(s)
	if m == nil {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if tailwindImportExclude(m[1]) {
		return false
	}

	return true
}

func (imp *importResolver) toFileError(output string) error {
	inErr := errors.New(output)

	match := cssSyntaxErrorRe.FindStringSubmatch(output)
	if match == nil {
		return inErr
	}

	lineNum, err := strconv.Atoi(match[1])
	if err != nil {
		return inErr
	}

	file, ok := imp.linemap[lineNum]
	if !ok {
		return inErr
	}

	fi, err := imp.fs.Stat(file.Filename)
	if err != nil {
		return inErr
	}

	meta := fi.(hugofs.FileMetaInfo).Meta()
	realFilename := meta.Filename
	f, err := meta.Open()
	if err != nil {
		return inErr
	}
	defer f.Close()

	ferr := herrors.NewFileErrorFromName(inErr, realFilename)
	pos := ferr.Position()
	pos.LineNumber = file.Offset + 1
	return ferr.UpdatePosition(pos).UpdateContent(f, nil)

	// return herrors.NewFileErrorFromFile(inErr, file.Filename, realFilename, hugofs.Os, herrors.SimpleLineMatcher)
}
