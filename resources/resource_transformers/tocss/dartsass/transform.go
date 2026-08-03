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
	"context"
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
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/sass"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/bep/godartsass/v2"
)

// Prefix for canonical URLs of stylesheets resolved in the user provided import context.
// Note: This prefix must be all lower case.
const dartSassImportContextPrefix = "hugoimportcontext:"

// Supports returns whether sass, dart-sass, or dart-sass-embedded is found in $PATH.
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

	var ic resource.ResourceGetter
	if opts.ImportContext != nil {
		ic = resource.NewCachedResourceGetter(opts.ImportContext)
	}

	args := godartsass.Args{
		URL:          filename,
		IncludePaths: t.c.sfs.RealDirs(baseDir),
		ImportResolver: importResolver{
			baseDir:           baseDir,
			c:                 t.c,
			dependencyManager: ctx.DependencyManager,
			importContext:     ic,
			ctx:               ctx.Ctx,

			vars: opts.Vars,
		},
		OutputStyle:                   godartsass.ParseOutputStyle(opts.OutputStyle),
		EnableSourceMap:               opts.EnableSourceMap,
		SourceMapIncludeSources:       opts.SourceMapIncludeSources,
		SilenceDeprecations:           opts.SilenceDeprecations,
		SilenceDependencyDeprecations: opts.SilenceDependencyDeprecations,
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
		if err := ctx.PublishSourceMap([]byte(res.SourceMap)); err != nil {
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
	importContext     resource.ResourceGetter
	ctx               context.Context
	vars              map[string]any
}

func (t importResolver) CanonicalizeURL(url string) (string, error) {
	if _, ok := sass.HugoVarsSubPath(url); ok {
		return strings.ToLower(url), nil
	}

	if r := t.resolveInImportContext(url); r != nil {
		t.dependencyManager.AddIdentity(identity.FirstIdentity(r))
		return dartSassImportContextPrefix + paths.ToSlashTrimLeading(resources.InternalResourceTargetPath(r)), nil
	}

	filePath, isURL := paths.UrlStringToFilename(url)
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
	namePatterns := sassNamePatterns(name)
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

	// Not found, let Dart Sass handle it
	return "", nil
}

func sassNamePatterns(name string) []string {
	if strings.Contains(name, ".") {
		return []string{"_%s", "%s"}
	}
	if strings.HasPrefix(name, "_") {
		return []string{"_%s.scss", "_%s.sass", "_%s.css"}
	}
	return []string{
		"_%s.scss", "%s.scss",
		"_%s.sass", "%s.sass",
		"_%s.css", "%s.css",
		"%s/_index.scss", "%s/_index.sass",
		"%s/index.scss", "%s/index.sass",
	}
}

// resolveInImportContext resolves url in the user provided import context, if any,
// using the same name patterns as for the file system.
func (t importResolver) resolveInImportContext(url string) resource.Resource {
	if t.importContext == nil {
		return nil
	}
	url = strings.TrimPrefix(url, dartSassImportContextPrefix)
	if r := t.importContext.Get(url); r != nil {
		return r
	}
	dir, name := path.Split(url)
	namePatterns := sassNamePatterns(name)
	name = strings.TrimPrefix(name, "_")
	for _, namePattern := range namePatterns {
		if r := t.importContext.Get(path.Join(dir, fmt.Sprintf(namePattern, name))); r != nil {
			return r
		}
	}
	return nil
}

func (t importResolver) Load(url string) (godartsass.Import, error) {
	if subPath, ok := sass.HugoVarsSubPath(url); ok {
		return godartsass.Import{
			Content: sass.CreateVarsStyleSheet(sass.TranspilerDart, sass.ResolveVars(t.vars, subPath)),
		}, nil
	}

	if strings.HasPrefix(url, dartSassImportContextPrefix) {
		r := t.resolveInImportContext(url)
		if r == nil {
			return godartsass.Import{}, fmt.Errorf("could not find %q in the import context", url)
		}
		content, err := resources.InternalResourceSourceContent(t.ctx, r)
		return godartsass.Import{Content: content, SourceSyntax: sassSourceSyntax(url)}, err
	}

	filename, _ := paths.UrlStringToFilename(url)
	b, err := afero.ReadFile(hugofs.Os, filename)

	return godartsass.Import{Content: string(b), SourceSyntax: sassSourceSyntax(filename)}, err
}

func sassSourceSyntax(name string) godartsass.SourceSyntax {
	switch {
	case strings.HasSuffix(name, ".sass"):
		return godartsass.SourceSyntaxSASS
	case strings.HasSuffix(name, ".css"):
		return godartsass.SourceSyntaxCSS
	default:
		return godartsass.SourceSyntaxSCSS
	}
}
