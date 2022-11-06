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

package postcss

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cast"

	"errors"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

const importIdentifier = "@import"

var (
	cssSyntaxErrorRe = regexp.MustCompile(`> (\d+) \|`)
	shouldImportRe   = regexp.MustCompile(`^@import ["'].*["'];?\s*(/\*.*\*/)?$`)
)

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

func decodeOptions(m map[string]any) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)

	if !opts.NoMap {
		// There was for a long time a discrepancy between documentation and
		// implementation for the noMap property, so we need to support both
		// camel and snake case.
		opts.NoMap = cast.ToBool(m["no-map"])
	}

	return
}

// Client is the client used to do PostCSS transformations.
type Client struct {
	rs *resources.Spec
}

// Process transforms the given Resource with the PostCSS processor.
func (c *Client) Process(res resources.ResourceTransformer, options map[string]any) (resource.Resource, error) {
	return res.Transform(&postcssTransformation{rs: c.rs, optionsm: options})
}

// Some of the options from https://github.com/postcss/postcss-cli
type Options struct {

	// Set a custom path to look for a config file.
	Config string

	NoMap bool // Disable the default inline sourcemaps

	// Enable inlining of @import statements.
	// Does so recursively, but currently once only per file;
	// that is, it's not possible to import the same file in
	// different scopes (root, media query...)
	// Note that this import routine does not care about the CSS spec,
	// so you can have @import anywhere in the file.
	InlineImports bool

	// When InlineImports is enabled, we fail the build if an import cannot be resolved.
	// You can enable this to allow the build to continue and leave the import statement in place.
	// Note that the inline importer does not process url location or imports with media queries,
	// so those will be left as-is even without enabling this option.
	SkipInlineImportsNotFound bool

	// Options for when not using a config file
	Use         string // List of postcss plugins to use
	Parser      string //  Custom postcss parser
	Stringifier string // Custom postcss stringifier
	Syntax      string // Custom postcss syntax
}

func (opts Options) toArgs() []string {
	var args []string
	if opts.NoMap {
		args = append(args, "--no-map")
	}
	if opts.Use != "" {
		args = append(args, "--use")
		args = append(args, strings.Fields(opts.Use)...)
	}
	if opts.Parser != "" {
		args = append(args, "--parser", opts.Parser)
	}
	if opts.Stringifier != "" {
		args = append(args, "--stringifier", opts.Stringifier)
	}
	if opts.Syntax != "" {
		args = append(args, "--syntax", opts.Syntax)
	}
	return args
}

type postcssTransformation struct {
	optionsm map[string]any
	rs       *resources.Spec
}

func (t *postcssTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("postcss", t.optionsm)
}

// Transform shells out to postcss-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g postcss-cli
// npm install -g autoprefixer
func (t *postcssTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const binaryName = "postcss"

	ex := t.rs.ExecHelper

	var configFile string
	logger := t.rs.Logger

	var options Options
	if t.optionsm != nil {
		var err error
		options, err = decodeOptions(t.optionsm)
		if err != nil {
			return err
		}
	}

	if options.Config != "" {
		configFile = options.Config
	} else {
		configFile = "postcss.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an absolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		configFile = t.rs.BaseFs.ResolveJSConfigFile(configFile)
		if configFile == "" && options.Config != "" {
			// Only fail if the user specified config file is not found.
			return fmt.Errorf("postcss config %q not found:", options.Config)
		}
	}

	var cmdArgs []any

	if configFile != "" {
		logger.Infoln("postcss: use config file", configFile)
		cmdArgs = []any{"--config", configFile}
	}

	if optArgs := options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, collections.StringSliceToInterfaceSlice(optArgs)...)
	}

	var errBuf bytes.Buffer
	infoW := loggers.LoggerToWriterWithPrefix(logger.Info(), "postcss")

	stderr := io.MultiWriter(infoW, &errBuf)
	cmdArgs = append(cmdArgs, hexec.WithStderr(stderr))
	cmdArgs = append(cmdArgs, hexec.WithStdout(ctx.To))
	cmdArgs = append(cmdArgs, hexec.WithEnviron(hugo.GetExecEnviron(t.rs.WorkingDir, t.rs.Cfg, t.rs.BaseFs.Assets.Fs)))

	cmd, err := ex.Npx(binaryName, cmdArgs...)
	if err != nil {
		if hexec.IsNotFound(err) {
			// This may be on a CI server etc. Will fall back to pre-built assets.
			return herrors.ErrFeatureNotAvailable
		}
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	src := ctx.From

	imp := newImportResolver(
		ctx.From,
		ctx.InPath,
		options,
		t.rs.Assets.Fs, t.rs.Logger,
	)

	if options.InlineImports {
		var err error
		src, err = imp.resolve()
		if err != nil {
			return err
		}
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, src)
	}()

	err = cmd.Run()
	if err != nil {
		if hexec.IsNotFound(err) {
			return herrors.ErrFeatureNotAvailable
		}
		return imp.toFileError(errBuf.String())
	}

	return nil
}

type fileOffset struct {
	Filename string
	Offset   int
}

type importResolver struct {
	r      io.Reader
	inPath string
	opts   Options

	contentSeen map[string]bool
	linemap     map[int]fileOffset
	fs          afero.Fs
	logger      loggers.Logger
}

func newImportResolver(r io.Reader, inPath string, opts Options, fs afero.Fs, logger loggers.Logger) *importResolver {
	return &importResolver{
		r:      r,
		inPath: inPath,
		fs:     fs, logger: logger,
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
	inPath string) (int, string, error) {
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
	const importIdentifier = "@import"

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

	return shouldImportRe.MatchString(s)
}

func (imp *importResolver) toFileError(output string) error {
	output = strings.TrimSpace(loggers.RemoveANSIColours(output))
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

	//return herrors.NewFileErrorFromFile(inErr, file.Filename, realFilename, hugofs.Os, herrors.SimpleLineMatcher)

}
