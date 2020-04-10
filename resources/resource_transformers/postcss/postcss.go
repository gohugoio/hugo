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
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/pkg/errors"

	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

const importIdentifier = "@import"

var cssSyntaxErrorRe = regexp.MustCompile(`> (\d+) \|`)

var shouldImportRe = regexp.MustCompile(`^@import ["'].*["'];?\s*(/\*.*\*/)?$`)

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

func DecodeOptions(m map[string]interface{}) (opts Options, err error) {
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
func (c *Client) Process(res resources.ResourceTransformer, options Options) (resource.Resource, error) {
	return res.Transform(&postcssTransformation{rs: c.rs, options: options})
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
		args = append(args, "--use", opts.Use)
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
	options Options
	rs      *resources.Spec
}

func (t *postcssTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("postcss", t.options)
}

// Transform shells out to postcss-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g postcss-cli
// npm install -g autoprefixer
func (t *postcssTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {

	const localPostCSSPath = "node_modules/.bin/"
	const binaryName = "postcss"

	// Try first in the project's node_modules.
	csiBinPath := filepath.Join(t.rs.WorkingDir, localPostCSSPath, binaryName)

	binary := csiBinPath

	if _, err := exec.LookPath(binary); err != nil {
		// Try PATH
		binary = binaryName
		if _, err := exec.LookPath(binary); err != nil {
			// This may be on a CI server etc. Will fall back to pre-built assets.
			return herrors.ErrFeatureNotAvailable
		}
	}

	var configFile string
	logger := t.rs.Logger

	if t.options.Config != "" {
		configFile = t.options.Config
	} else {
		configFile = "postcss.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an abolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		// We resolve this against the virtual Work filesystem, to allow
		// this config file to live in one of the themes if needed.
		fi, err := t.rs.BaseFs.Work.Stat(configFile)
		if err != nil {
			if t.options.Config != "" {
				// Only fail if the user specificed config file is not found.
				return errors.Wrapf(err, "postcss config %q not found:", configFile)
			}
			configFile = ""
		} else {
			configFile = fi.(hugofs.FileMetaInfo).Meta().Filename()
		}
	}

	var cmdArgs []string

	if configFile != "" {
		logger.INFO.Println("postcss: use config file", configFile)
		cmdArgs = []string{"--config", configFile}
	}

	if optArgs := t.options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, optArgs...)
	}

	cmd := exec.Command(binary, cmdArgs...)

	var errBuf bytes.Buffer

	cmd.Stdout = ctx.To
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	// TODO(bep) somehow generalize this to other external helpers that may need this.
	env := os.Environ()
	config.SetEnvVars(&env, "HUGO_ENVIRONMENT", t.rs.Cfg.GetString("environment"))
	cmd.Env = env

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	src := ctx.From

	imp := newImportResolver(
		ctx.From,
		ctx.InPath,
		t.rs.Assets.Fs, t.rs.Logger,
	)

	if t.options.InlineImports {
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

	contentSeen map[string]bool
	linemap     map[int]fileOffset
	fs          afero.Fs
	logger      *loggers.Logger
}

func newImportResolver(r io.Reader, inPath string, fs afero.Fs, logger *loggers.Logger) *importResolver {
	return &importResolver{
		r:      r,
		inPath: inPath,
		fs:     fs, logger: logger,
		linemap: make(map[int]fileOffset), contentSeen: make(map[string]bool),
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
		line = strings.TrimSpace(line)

		if !imp.shouldImport(line) {
			trackLine(i, offset, line)
		} else {
			i--
			path := strings.Trim(strings.TrimPrefix(line, importIdentifier), " \"';")
			filename := filepath.Join(basePath, path)
			importContent, hash := imp.contentHash(filename)
			if importContent == nil {
				trackLine(i, offset, "ERROR")
				imp.logger.WARN.Printf("postcss: Failed to resolve CSS @import in %q for path %q", inPath, filename)
				continue
			}

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

	content, err := ioutil.ReadAll(imp.r)
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
//     @import "navigation.css";
// This is not:
//     @import url("navigation.css");
//     @import "mobstyle.css" screen and (max-width: 768px);
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
	inErr := errors.New(strings.TrimSpace(output))

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
	realFilename := fi.(hugofs.FileMetaInfo).Meta().Filename()

	ferr := herrors.NewFileError("css", -1, file.Offset+1, 1, inErr)

	werr, ok := herrors.WithFileContextForFile(ferr, realFilename, file.Filename, imp.fs, herrors.SimpleLineMatcher)

	if !ok {
		return ferr
	}

	return werr
}
