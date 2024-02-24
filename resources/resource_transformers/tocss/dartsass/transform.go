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

package dartsass

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/internal/sass"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	godartsassv1 "github.com/bep/godartsass"
	"github.com/bep/godartsass/v2"
)

// Supports returns whether dart-sass-embedded is found in $PATH.
func Supports() bool {
	if htesting.SupportsAll() {
		return true
	}
	return hugo.DartSassBinaryName != ""
}

type transform struct {
	optsm map[string]any
	c     *Client
}

func (t *transform) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey(transformationName, t.optsm)
}

func (t *transform) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.Builtin.CSSType

	opts, err := decodeOptions(t.optsm)
	if err != nil {
		return err
	}

	if opts.TargetPath != "" {
		ctx.OutPath = opts.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".css")
	}

	baseDir := path.Dir(ctx.SourcePath)
	filename := dartSassStdinPrefix

	if ctx.SourcePath != "" {
		filename += t.c.sfs.RealFilename(ctx.SourcePath)
	}

	args := godartsass.Args{
		URL:          filename,
		IncludePaths: t.c.sfs.RealDirs(baseDir),
		ImportResolver: importResolver{
			baseDir:           baseDir,
			c:                 t.c,
			dependencyManager: ctx.DependencyManager,

			varsStylesheet: godartsass.Import{Content: sass.CreateVarsStyleSheet(opts.Vars)},
		},
		OutputStyle:             godartsass.ParseOutputStyle(opts.OutputStyle),
		EnableSourceMap:         opts.EnableSourceMap,
		SourceMapIncludeSources: opts.SourceMapIncludeSources,
	}

	// Append any workDir relative include paths
	for _, ip := range opts.IncludePaths {
		info, err := t.c.workFs.Stat(filepath.Clean(ip))
		if err == nil {
			filename := info.(hugofs.FileMetaInfo).Meta().Filename
			args.IncludePaths = append(args.IncludePaths, filename)
		}
	}

	if ctx.InMediaType.SubType == media.Builtin.SASSType.SubType {
		args.SourceSyntax = godartsass.SourceSyntaxSASS
	}

	res, err := t.c.toCSS(args, ctx.From)
	if err != nil {
		return err
	}

	out := res.CSS

	_, err = io.WriteString(ctx.To, out)
	if err != nil {
		return err
	}

	if opts.EnableSourceMap && res.SourceMap != "" {
		if err := ctx.PublishSourceMap(res.SourceMap); err != nil {
			return err
		}
		_, err = fmt.Fprintf(ctx.To, "\n\n/*# sourceMappingURL=%s */", path.Base(ctx.OutPath)+".map")
	}

	return err
}

type importResolver struct {
	baseDir           string
	c                 *Client
	dependencyManager identity.Manager
	varsStylesheet    godartsass.Import
}

func (t importResolver) CanonicalizeURL(url string) (string, error) {
	if url == sass.HugoVarsNamespace {
		return url, nil
	}

	filePath, isURL := paths.UrlToFilename(url)
	var prevDir string
	var pathDir string
	if isURL {
		var found bool
		prevDir, found = t.c.sfs.MakePathRelative(filepath.Dir(filePath), true)

		if !found {
			// Not a member of this filesystem, let Dart Sass handle it.
			return "", nil
		}
	} else {
		prevDir = t.baseDir
		pathDir = path.Dir(url)
	}

	basePath := filepath.Join(prevDir, pathDir)
	name := filepath.Base(filePath)

	// Pick the first match.
	var namePatterns []string
	if strings.Contains(name, ".") {
		namePatterns = []string{"_%s", "%s"}
	} else if strings.HasPrefix(name, "_") {
		namePatterns = []string{"_%s.scss", "_%s.sass", "_%s.css"}
	} else {
		namePatterns = []string{"_%s.scss", "%s.scss", "_%s.sass", "%s.sass", "_%s.css", "%s.css"}
	}

	name = strings.TrimPrefix(name, "_")

	for _, namePattern := range namePatterns {
		filenameToCheck := filepath.Join(basePath, fmt.Sprintf(namePattern, name))
		fi, err := t.c.sfs.Fs.Stat(filenameToCheck)
		if err == nil {
			if fim, ok := fi.(hugofs.FileMetaInfo); ok {
				t.dependencyManager.AddIdentity(identity.CleanStringIdentity(filenameToCheck))
				return "file://" + filepath.ToSlash(fim.Meta().Filename), nil
			}
		}
	}

	// Not found, let Dart Dass handle it
	return "", nil
}

func (t importResolver) Load(url string) (godartsass.Import, error) {
	if url == sass.HugoVarsNamespace {
		return t.varsStylesheet, nil
	}
	filename, _ := paths.UrlToFilename(url)
	b, err := afero.ReadFile(hugofs.Os, filename)

	sourceSyntax := godartsass.SourceSyntaxSCSS
	if strings.HasSuffix(filename, ".sass") {
		sourceSyntax = godartsass.SourceSyntaxSASS
	} else if strings.HasSuffix(filename, ".css") {
		sourceSyntax = godartsass.SourceSyntaxCSS
	}

	return godartsass.Import{Content: string(b), SourceSyntax: sourceSyntax}, err
}

type importResolverV1 struct {
	godartsass.ImportResolver
}

func (t importResolverV1) Load(url string) (godartsassv1.Import, error) {
	res, err := t.ImportResolver.Load(url)
	return godartsassv1.Import{Content: res.Content, SourceSyntax: godartsassv1.SourceSyntax(res.SourceSyntax)}, err
}
