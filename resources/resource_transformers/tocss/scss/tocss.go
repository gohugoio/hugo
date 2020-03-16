// Copyright 2018 The Hugo Authors. All rights reserved.
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

// +build extended

package scss

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/bep/golibsass/libsass"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/pkg/errors"
)

// Used in tests. This feature requires Hugo to be built with the extended tag.
func Supports() bool {
	return true
}

func (t *toCSSTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.CSSType

	var outName string
	if t.options.from.TargetPath != "" {
		ctx.OutPath = t.options.from.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".css")
	}

	outName = path.Base(ctx.OutPath)

	options := t.options
	baseDir := path.Dir(ctx.SourcePath)
	options.to.IncludePaths = t.c.sfs.RealDirs(baseDir)

	// Append any workDir relative include paths
	for _, ip := range options.from.IncludePaths {
		info, err := t.c.workFs.Stat(filepath.Clean(ip))
		if err == nil {
			filename := info.(hugofs.FileMetaInfo).Meta().Filename()
			options.to.IncludePaths = append(options.to.IncludePaths, filename)
		}
	}

	// To allow for overrides of SCSS files anywhere in the project/theme hierarchy, we need
	// to help libsass revolve the filename by looking in the composite filesystem first.
	// We add the entry directories for both project and themes to the include paths list, but
	// that only work for overrides on the top level.
	options.to.ImportResolver = func(url string, prev string) (newUrl string, body string, resolved bool) {
		// We get URL paths from LibSASS, but we need file paths.
		url = filepath.FromSlash(url)
		prev = filepath.FromSlash(prev)

		var basePath string
		urlDir := filepath.Dir(url)
		var prevDir string

		if prev == "stdin" {
			prevDir = baseDir
		} else {
			prevDir = t.c.sfs.MakePathRelative(filepath.Dir(prev))

			if prevDir == "" {
				// Not a member of this filesystem. Let LibSASS handle it.
				return "", "", false
			}
		}

		basePath = filepath.Join(prevDir, urlDir)
		name := filepath.Base(url)

		// Libsass throws an error in cases where you have several possible candidates.
		// We make this simpler and pick the first match.
		var namePatterns []string
		if strings.Contains(name, ".") {
			namePatterns = []string{"_%s", "%s"}
		} else if strings.HasPrefix(name, "_") {
			namePatterns = []string{"_%s.scss", "_%s.sass"}
		} else {
			namePatterns = []string{"_%s.scss", "%s.scss", "_%s.sass", "%s.sass"}
		}

		name = strings.TrimPrefix(name, "_")

		for _, namePattern := range namePatterns {
			filenameToCheck := filepath.Join(basePath, fmt.Sprintf(namePattern, name))
			fi, err := t.c.sfs.Fs.Stat(filenameToCheck)
			if err == nil {
				if fim, ok := fi.(hugofs.FileMetaInfo); ok {
					return fim.Meta().Filename(), "", true
				}
			}
		}

		// Not found, let LibSASS handle it
		return "", "", false
	}

	if ctx.InMediaType.SubType == media.SASSType.SubType {
		options.to.SassSyntax = true
	}

	if options.from.EnableSourceMap {

		options.to.SourceMapOptions.Filename = outName + ".map"
		options.to.SourceMapOptions.Root = t.c.rs.WorkingDir

		// Setting this to the relative input filename will get the source map
		// more correct for the main entry path (main.scss typically), but
		// it will mess up the import mappings. As a workaround, we do a replacement
		// in the source map itself (see below).
		//options.InputPath = inputPath
		options.to.SourceMapOptions.OutputPath = outName
		options.to.SourceMapOptions.Contents = true
		options.to.SourceMapOptions.OmitURL = false
		options.to.SourceMapOptions.EnableEmbedded = false
	}

	res, err := t.c.toCSS(options.to, ctx.To, ctx.From)
	if err != nil {
		return err
	}

	if options.from.EnableSourceMap && res.SourceMapContent != "" {
		sourcePath := t.c.sfs.RealFilename(ctx.SourcePath)

		if strings.HasPrefix(sourcePath, t.c.rs.WorkingDir) {
			sourcePath = strings.TrimPrefix(sourcePath, t.c.rs.WorkingDir+helpers.FilePathSeparator)
		}

		// This needs to be Unix-style slashes, even on Windows.
		// See https://github.com/gohugoio/hugo/issues/4968
		sourcePath = filepath.ToSlash(sourcePath)

		// This is a workaround for what looks like a bug in Libsass. But
		// getting this resolution correct in tools like Chrome Workspaces
		// is important enough to go this extra mile.
		mapContent := strings.Replace(res.SourceMapContent, `stdin",`, fmt.Sprintf("%s\",", sourcePath), 1)

		return ctx.PublishSourceMap(mapContent)
	}
	return nil
}

func (c *Client) toCSS(options libsass.Options, dst io.Writer, src io.Reader) (libsass.Result, error) {
	var res libsass.Result

	transpiler, err := libsass.New(options)
	if err != nil {
		return res, err
	}

	in := helpers.ReaderToString(src)

	// See https://github.com/gohugoio/hugo/issues/7059
	// We need to preserver the regular CSS imports. This is by far
	// a perfect solution, and only works for the main entry file, but
	// that should cover many use cases, e.g. using SCSS as a preprocessor
	// for Tailwind.
	var importsReplaced bool
	in, importsReplaced = replaceRegularImportsIn(in)

	res, err = transpiler.Execute(in)
	if err != nil {
		return res, errors.Wrap(err, "SCSS processing failed")
	}

	out := res.CSS
	if importsReplaced {
		out = replaceRegularImportsOut(out)
	}

	_, err = io.WriteString(dst, out)

	return res, err
}
