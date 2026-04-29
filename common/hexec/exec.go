// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hexec

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/security"
)

var WithDir = func(dir string) func(c *commandeer) {
	return func(c *commandeer) {
		c.dir = dir
	}
}

var WithContext = func(ctx context.Context) func(c *commandeer) {
	return func(c *commandeer) {
		c.ctx = ctx
	}
}

var WithStdout = func(w io.Writer) func(c *commandeer) {
	return func(c *commandeer) {
		c.stdout = w
	}
}

var WithStderr = func(w io.Writer) func(c *commandeer) {
	return func(c *commandeer) {
		c.stderr = w
	}
}

var WithStdin = func(r io.Reader) func(c *commandeer) {
	return func(c *commandeer) {
		c.stdin = r
	}
}

var WithEnviron = func(env []string) func(c *commandeer) {
	return func(c *commandeer) {
		setOrAppend := func(s string) {
			k1, _ := config.SplitEnvVar(s)
			var found bool
			for i, v := range c.env {
				k2, _ := config.SplitEnvVar(v)
				if k1 == k2 {
					found = true
					c.env[i] = s
				}
			}

			if !found {
				c.env = append(c.env, s)
			}
		}

		for _, s := range env {
			setOrAppend(s)
		}
	}
}

// New creates a new Exec using the provided security config.
func New(cfg security.Config, workingDir string, log loggers.Logger) *Exec {
	var baseEnviron []string
	for _, v := range os.Environ() {
		k, _ := config.SplitEnvVar(v)
		if cfg.Exec.OsEnv.Accept(k) {
			baseEnviron = append(baseEnviron, v)
		}
	}

	return &Exec{
		sc:              cfg,
		workingDir:      workingDir,
		infol:           log.InfoCommand("exec"),
		baseEnviron:     baseEnviron,
		nodeRunnerCache: hmaps.NewCache[string, func(arg ...any) (Runner, error)](),
	}
}

// IsNotFound reports whether this is an error about a binary not found.
func IsNotFound(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}

// Exec enforces a security policy for commands run via os/exec.
type Exec struct {
	sc         security.Config
	workingDir string
	infol      logg.LevelLogger

	// os.Environ filtered by the Exec.OsEnviron whitelist filter.
	baseEnviron []string

	// Additional absolute paths to allow reading from in the Node.js permission model.
	nodeReadPaths []string

	nodeRunnerCache *hmaps.Cache[string, func(arg ...any) (Runner, error)]
}

// SetNodeReadPaths sets additional absolute paths to allow reading from
// in the Node.js permission model (e.g. Hugo module cache directories).
func (e *Exec) SetNodeReadPaths(paths []string) {
	e.nodeReadPaths = paths
}

func (e *Exec) New(name string, arg ...any) (Runner, error) {
	return e.new(name, "", arg...)
}

// New will fail if name is not allowed according to the configured security policy.
// Else a configured Runner will be returned ready to be Run.
func (e *Exec) new(name string, fullyQualifiedName string, arg ...any) (Runner, error) {
	if err := e.sc.CheckAllowedExec(name); err != nil {
		return nil, err
	}

	env := make([]string, len(e.baseEnviron))
	copy(env, e.baseEnviron)

	cm := &commandeer{
		name:               name,
		fullyQualifiedName: fullyQualifiedName,
		env:                env,
	}

	return cm.command(arg...)
}

type binaryLocation int

func (b binaryLocation) String() string {
	switch b {
	case binaryLocationNodeModules:
		return "node_modules/.bin"
	case binaryLocationPath:
		return "PATH"
	}
	return "unknown"
}

const (
	binaryLocationNodeModules binaryLocation = iota + 1
	binaryLocationPath
)

// Npx finds and runs a Node.js tool. The binary is located first in
// WORKINGDIR/node_modules/.bin, then in PATH. The tool is always invoked via
// "node [--permission <flags>] <script> <args>"; the --permission flags are
// added when the Node.js permission model is enabled.
func (e *Exec) Npx(name string, arg ...any) (Runner, error) {
	if err := e.sc.CheckAllowedExec(name); err != nil {
		return nil, err
	}
	if err := e.sc.CheckAllowedExec("node"); err != nil {
		// Legacy path: We replaced npx with node in v0.161.0, and anyone using these tools with a custom security.exec.allow list
		// would get an error when upgrading. To avoid this, check for npx as well.
		if err2 := e.sc.CheckAllowedExec("npx"); err2 != nil {
			return nil, err
		}
	}

	newRunner, err := e.nodeRunnerCache.GetOrCreate(name, func() (func(...any) (Runner, error), error) {
		var resolvedBin string
		var loc binaryLocation

		nodeBinFilename := filepath.Join(e.workingDir, nodeModulesBinPath, name)
		if p, err := exec.LookPath(nodeBinFilename); err == nil {
			resolvedBin = p
			loc = binaryLocationNodeModules
		} else if p, err := exec.LookPath(name); err == nil {
			resolvedBin = p
			loc = binaryLocationPath
		} else {
			return nil, &NotFoundError{name: name, method: "in PATH"}
		}

		scriptPath := resolveNodeBin(resolvedBin)

		e.infol.WithFields(logg.Fields{
			logg.Field{Name: "location", Value: loc},
			logg.Field{Name: "bin", Value: resolvedBin},
			logg.Field{Name: "script", Value: scriptPath},
		}).Logf("resolve %q", name)

		if scriptPath == "" {
			return nil, fmt.Errorf("binary %q is not a Node.js script", name)
		}

		return func(arg2 ...any) (Runner, error) {
			return e.newNode(name, scriptPath, arg2...)
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return newRunner(arg...)
}

// newNode runs a Node.js script via "node [--permission <flags>] <scriptPath> <args>".
func (e *Exec) newNode(name, scriptPath string, arg ...any) (Runner, error) {
	var allArgs []any
	for _, pa := range e.nodePermissionArgs(name, scriptPath) {
		allArgs = append(allArgs, pa)
	}
	allArgs = append(allArgs, scriptPath)
	allArgs = append(allArgs, arg...)
	// When the script lives outside the working dir (a globally installed
	// tool), point NODE_PATH at the script's node_modules ancestor so Node's
	// resolver (and tools that honor it, e.g. tailwindcss v4) can locate the
	// tool's sibling packages. tailwindcss v4's CSS resolver treats NODE_PATH
	// as a single path, not a list, so we don't concatenate with the local
	// path here. For local installs the caller's NODE_PATH (set by
	// hugo.GetExecEnviron to <workDir>/node_modules) already covers the need.
	localNM := filepath.Join(e.workingDir, "node_modules")
	if p := nodeScriptReadPath(scriptPath); p != "" && p != localNM {
		allArgs = append(allArgs, WithEnviron([]string{"NODE_PATH=" + p}))
	}

	return e.New("node", allArgs...)
}

// nodePermissionArgs builds the Node.js --permission flags from the security config.
func (e *Exec) nodePermissionArgs(name, scriptPath string) []string {
	perms := e.sc.Node.Permissions
	if !perms.IsEnabled() {
		return nil
	}

	args := []string{"--permission"}

	for _, p := range e.resolveNodePermPaths(perms.AllowRead) {
		args = append(args, "--allow-fs-read="+p)
	}
	for _, p := range e.nodeReadPaths {
		args = append(args, "--allow-fs-read="+p)
	}
	if p := nodeScriptReadPath(scriptPath); p != "" {
		args = append(args, "--allow-fs-read="+p)
	}

	for _, p := range e.resolveNodePermPaths(perms.AllowWrite) {
		args = append(args, "--allow-fs-write="+p)
	}

	var silenceSecurityWarnings bool
	if slices.Contains(perms.AllowAddons, name) {
		silenceSecurityWarnings = true
		args = append(args, "--allow-addons")
	}

	if slices.Contains(perms.AllowChildProcess, name) {
		silenceSecurityWarnings = true
		args = append(args, "--allow-child-process")
	}

	if slices.Contains(perms.AllowWorker, name) {
		silenceSecurityWarnings = true
		args = append(args, "--allow-worker")
	}

	if silenceSecurityWarnings {
		// There are no more fine grained way to do this, see https://github.com/nodejs/node/issues/59818
		// If the process is configured to allow either workers or addons, Node will print warnings that's not very helpful.
		args = append(args, "--disable-warning=SecurityWarning")
	}

	return args
}

// resolveNodePermPaths resolves relative paths against the working directory.
func (e *Exec) resolveNodePermPaths(paths []string) []string {
	resolved := make([]string, len(paths))
	for i, p := range paths {
		switch {
		case p == "*":
			resolved[i] = "*"
		case filepath.IsAbs(p):
			resolved[i] = p
		default:
			resolved[i] = filepath.Join(e.workingDir, p)
		}
	}
	return resolved
}

const nodeModulesBinPath = "node_modules/.bin"

// nodeScriptReadPath returns a path to add to the Node.js read allow-list so
// a script can load its dependencies. For scripts inside a node_modules tree
// it returns the nearest ancestor "node_modules" directory, so both nested
// and hoisted deps are reachable. Otherwise the script's own directory.
func nodeScriptReadPath(scriptPath string) string {
	if scriptPath == "" {
		return ""
	}
	dir := filepath.Dir(scriptPath)
	for {
		if filepath.Base(dir) == "node_modules" {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Dir(scriptPath)
		}
		dir = parent
	}
}

// resolveNodeBin resolves a binary path to the underlying Node.js script.
// Returns the path to the JS entry point, or "" if the binary is not a Node script.
func resolveNodeBin(path string) string {
	// 1. If the file is a symlink, resolve it (macOS/Linux npm creates symlinks in node_modules/.bin).
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			if hasJSExtension(resolved) || isNodeScript(resolved) {
				return resolved
			}
		}
		return ""
	}
	// 2. Check if the file itself is a Node script (e.g. globally installed with #!/usr/bin/env node).
	if isNodeScript(path) {
		return path
	}
	// 3. Try extracting JS entry point from an npm wrapper script (.cmd or shell).
	return extractNodeEntryPoint(path)
}

// nodeEntryPointRe matches a relative path in npm-generated wrapper scripts.
// The entry may be a .js/.mjs/.cjs file or an extensionless Node shebang
// script (e.g. postcss-cli 7's bin/postcss). Local installs reference the
// entry via "..", global installs via "node_modules" (notably on Windows,
// where npm does not symlink global binaries).
// Examples:
//
//	Local shell:  "$basedir/../postcss-cli/index.js"
//	Local cmd:    "%dp0%\..\postcss-cli\index.js"
//	Scoped:       "$basedir/../@babel/cli/bin/babel.js"
//	No ext:       "%dp0%\..\postcss-cli\bin\postcss"
//	Global shell: "$basedir/node_modules/postcss-cli/index.js"
//	Global cmd:   "%dp0%\node_modules\postcss-cli\index.js"
var nodeEntryPointRe = regexp.MustCompile(`[/\\]((?:\.\.|node_modules)[/\\][\w@][\w@./\\-]*)`)

// extractNodeEntryPoint reads an npm wrapper script and extracts the Node
// entry point path, validating that it's a JS file or a Node shebang script.
func extractNodeEntryPoint(wrapperPath string) string {
	data, err := os.ReadFile(wrapperPath)
	if err != nil {
		return ""
	}
	m := nodeEntryPointRe.FindSubmatch(data)
	if m == nil {
		return ""
	}
	// Normalize backslashes from Windows .cmd wrappers.
	relPath := strings.ReplaceAll(string(m[1]), "\\", "/")
	resolved := filepath.Join(filepath.Dir(wrapperPath), relPath)
	if _, err := os.Stat(resolved); err != nil {
		return ""
	}
	if !hasJSExtension(resolved) && !isNodeScript(resolved) {
		return ""
	}
	return resolved
}

func hasJSExtension(path string) bool {
	switch filepath.Ext(path) {
	case ".js", ".mjs", ".cjs":
		return true
	}
	return false
}

// isNodeScript reports whether the file at path has a Node.js shebang.
func isNodeScript(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false
	}
	return strings.HasPrefix(line, "#!") && strings.Contains(line, "node")
}

// Sec returns the security policies this Exec is configured with.
func (e *Exec) Sec() security.Config {
	return e.sc
}

type NotFoundError struct {
	name   string
	method string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("binary with name %q not found %s", e.name, e.method)
}

// Runner wraps a *os.Cmd.
type Runner interface {
	Run() error
	StdinPipe() (io.WriteCloser, error)
}

type cmdWrapper struct {
	name string
	c    *exec.Cmd

	outerr *bytes.Buffer
}

var notFoundRe = regexp.MustCompile(`(?s)not found:|could not determine executable`)

func (c *cmdWrapper) Run() error {
	err := c.c.Run()
	if err == nil {
		return nil
	}
	if notFoundRe.MatchString(c.outerr.String()) {
		return &NotFoundError{name: c.name, method: "in PATH"}
	}
	return fmt.Errorf("failed to execute binary %q with args %v: %s", c.name, c.c.Args[1:], c.outerr.String())
}

func (c *cmdWrapper) StdinPipe() (io.WriteCloser, error) {
	return c.c.StdinPipe()
}

type commandeer struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	dir    string
	ctx    context.Context

	name               string
	fullyQualifiedName string
	env                []string
}

func (c *commandeer) command(arg ...any) (*cmdWrapper, error) {
	if c == nil {
		return nil, nil
	}

	var args []string
	for _, a := range arg {
		switch v := a.(type) {
		case string:
			args = append(args, v)
		case func(*commandeer):
			v(c)
		default:
			return nil, fmt.Errorf("invalid argument to command: %T", a)
		}
	}

	var bin string
	if c.fullyQualifiedName != "" {
		bin = c.fullyQualifiedName
	} else {
		var err error
		bin, err = exec.LookPath(c.name)
		if err != nil {
			return nil, &NotFoundError{
				name:   c.name,
				method: "in PATH",
			}
		}
	}

	outerr := &bytes.Buffer{}
	if c.stderr == nil {
		c.stderr = outerr
	} else {
		c.stderr = io.MultiWriter(c.stderr, outerr)
	}

	var cmd *exec.Cmd

	if c.ctx != nil {
		cmd = exec.CommandContext(c.ctx, bin, args...)
	} else {
		cmd = exec.Command(bin, args...)
	}

	cmd.Stdin = c.stdin
	cmd.Stderr = c.stderr
	cmd.Stdout = c.stdout
	cmd.Env = c.env
	cmd.Dir = c.dir

	return &cmdWrapper{outerr: outerr, c: cmd, name: c.name}, nil
}

// InPath reports whether binaryName is in $PATH.
func InPath(binaryName string) bool {
	if strings.Contains(binaryName, "/") {
		panic("binary name should not contain any slash")
	}
	_, err := exec.LookPath(binaryName)
	return err == nil
}

// LookPath finds the path to binaryName in $PATH.
// Returns "" if not found.
func LookPath(binaryName string) string {
	if strings.Contains(binaryName, "/") {
		panic("binary name should not contain any slash")
	}
	s, err := exec.LookPath(binaryName)
	if err != nil {
		return ""
	}
	return s
}
