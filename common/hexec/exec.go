// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/bep/logg"
	"github.com/cli/safeexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
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
		sc:                cfg,
		workingDir:        workingDir,
		infol:             log.InfoCommand("exec"),
		baseEnviron:       baseEnviron,
		newNPXRunnerCache: maps.NewCache[string, func(arg ...any) (Runner, error)](),
	}
}

// IsNotFound reports whether this is an error about a binary not found.
func IsNotFound(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}

// SafeCommand is a wrapper around os/exec Command which uses a LookPath
// implementation that does not search in current directory before looking in PATH.
// See https://github.com/cli/safeexec and the linked issues.
func SafeCommand(name string, arg ...string) (*exec.Cmd, error) {
	bin, err := safeexec.LookPath(name)
	if err != nil {
		return nil, err
	}

	return exec.Command(bin, arg...), nil
}

// Exec enforces a security policy for commands run via os/exec.
type Exec struct {
	sc         security.Config
	workingDir string
	infol      logg.LevelLogger

	// os.Environ filtered by the Exec.OsEnviron whitelist filter.
	baseEnviron []string

	newNPXRunnerCache *maps.Cache[string, func(arg ...any) (Runner, error)]
	npxInit           sync.Once
	npxAvailable      bool
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
	case binaryLocationNpx:
		return "npx"
	case binaryLocationPath:
		return "PATH"
	}
	return "unknown"
}

const (
	binaryLocationNodeModules binaryLocation = iota + 1
	binaryLocationNpx
	binaryLocationPath
)

// Npx will in order:
// 1. Try fo find the binary in the WORKINGDIR/node_modules/.bin directory.
// 2. If not found, and npx is available, run npx --no-install <name> <args>.
// 3. Fall back to the PATH.
// If name is "tailwindcss", we will try the PATH as the second option.
func (e *Exec) Npx(name string, arg ...any) (Runner, error) {
	if err := e.sc.CheckAllowedExec(name); err != nil {
		return nil, err
	}

	newRunner, err := e.newNPXRunnerCache.GetOrCreate(name, func() (func(...any) (Runner, error), error) {
		type tryFunc func() func(...any) (Runner, error)
		tryFuncs := map[binaryLocation]tryFunc{
			binaryLocationNodeModules: func() func(...any) (Runner, error) {
				nodeBinFilename := filepath.Join(e.workingDir, nodeModulesBinPath, name)
				_, err := safeexec.LookPath(nodeBinFilename)
				if err != nil {
					return nil
				}
				return func(arg2 ...any) (Runner, error) {
					return e.new(name, nodeBinFilename, arg2...)
				}
			},
			binaryLocationNpx: func() func(...any) (Runner, error) {
				e.checkNpx()
				if !e.npxAvailable {
					return nil
				}
				return func(arg2 ...any) (Runner, error) {
					return e.npx(name, arg2...)
				}
			},
			binaryLocationPath: func() func(...any) (Runner, error) {
				if _, err := safeexec.LookPath(name); err != nil {
					return nil
				}
				return func(arg2 ...any) (Runner, error) {
					return e.New(name, arg2...)
				}
			},
		}

		locations := []binaryLocation{binaryLocationNodeModules, binaryLocationNpx, binaryLocationPath}
		if name == "tailwindcss" {
			// See https://github.com/gohugoio/hugo/issues/13221#issuecomment-2574801253
			locations = []binaryLocation{binaryLocationNodeModules, binaryLocationPath, binaryLocationNpx}
		}
		for _, loc := range locations {
			if f := tryFuncs[loc](); f != nil {
				e.infol.Logf("resolve %q using %s", name, loc)
				return f, nil
			}
		}
		return nil, &NotFoundError{name: name, method: fmt.Sprintf("in %s", locations[len(locations)-1])}
	})
	if err != nil {
		return nil, err
	}

	return newRunner(arg...)
}

const (
	npxNoInstall       = "--no-install"
	npxBinary          = "npx"
	nodeModulesBinPath = "node_modules/.bin"
)

func (e *Exec) checkNpx() {
	e.npxInit.Do(func() {
		e.npxAvailable = InPath(npxBinary)
	})
}

// npx is a convenience method to create a Runner running npx --no-install <name> <args.
func (e *Exec) npx(name string, arg ...any) (Runner, error) {
	arg = append(arg[:0], append([]any{npxNoInstall, name}, arg[0:]...)...)
	return e.New(npxBinary, arg...)
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
	name := c.name
	method := "in PATH"
	if name == npxBinary {
		name = c.c.Args[2]
		method = "using npx"
	}
	if notFoundRe.MatchString(c.outerr.String()) {
		return &NotFoundError{name: name, method: method}
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
		bin, err = safeexec.LookPath(c.name)
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
	_, err := safeexec.LookPath(binaryName)
	return err == nil
}

// LookPath finds the path to binaryName in $PATH.
// Returns "" if not found.
func LookPath(binaryName string) string {
	if strings.Contains(binaryName, "/") {
		panic("binary name should not contain any slash")
	}
	s, err := safeexec.LookPath(binaryName)
	if err != nil {
		return ""
	}
	return s
}
