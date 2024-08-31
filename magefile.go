//go:build mage
// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/codegen"
	"github.com/gohugoio/hugo/resources/page/page_generate"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName  = "github.com/gohugoio/hugo"
	noGitLdflags = "-X github.com/gohugoio/hugo/common/hugo.vendorInfo=mage"
)

var ldflags = noGitLdflags

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

func runWith(env map[string]string, cmd string, inArgs ...any) error {
	s := argsToStrings(inArgs...)
	return sh.RunWith(env, cmd, s...)
}

// Build hugo binary
func Hugo() error {
	return runWith(flagEnv(), goexe, "build", "-ldflags", ldflags, buildFlags(), "-tags", buildTags(), packageName)
}

// Build hugo binary with race detector enabled
func HugoRace() error {
	return runWith(flagEnv(), goexe, "build", "-race", "-ldflags", ldflags, buildFlags(), "-tags", buildTags(), packageName)
}

// Install hugo binary
func Install() error {
	return runWith(flagEnv(), goexe, "install", "-ldflags", ldflags, buildFlags(), "-tags", buildTags(), packageName)
}

// Uninstall hugo binary
func Uninstall() error {
	return sh.Run(goexe, "clean", "-i", packageName)
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

// Generate autogen packages
func Generate() error {
	generatorPackages := []string{
		"livereload/gen",
	}

	for _, pkg := range generatorPackages {
		if err := runWith(flagEnv(), goexe, "generate", path.Join(packageName, pkg)); err != nil {
			return err
		}
	}

	dir, _ := os.Getwd()
	c := codegen.NewInspector(dir)

	if err := page_generate.Generate(c); err != nil {
		return err
	}

	goFmtPatterns := []string{
		// TODO(bep) check: stat ./resources/page/*autogen*: no such file or directory
		"./resources/page/page_marshaljson.autogen.go",
		"./resources/page/page_wrappers.autogen.go",
		"./resources/page/zero_file.autogen.go",
	}

	for _, pattern := range goFmtPatterns {
		if err := sh.Run("gofmt", "-w", filepath.FromSlash(pattern)); err != nil {
			return err
		}
	}

	return nil
}

// Generate docs helper
func GenDocsHelper() error {
	return runCmd(flagEnv(), goexe, "run", "-tags", buildTags(), "main.go", "gen", "docshelper")
}

// Build hugo without git info
func HugoNoGitInfo() error {
	ldflags = noGitLdflags
	return Hugo()
}

var docker = sh.RunCmd("docker")

// Build hugo Docker container
func Docker() error {
	if err := docker("build", "-t", "hugo", "."); err != nil {
		return err
	}
	// yes ignore errors here
	docker("rm", "-f", "hugo-build")
	if err := docker("run", "--name", "hugo-build", "hugo ls /go/bin"); err != nil {
		return err
	}
	if err := docker("cp", "hugo-build:/go/bin/hugo", "."); err != nil {
		return err
	}
	return docker("rm", "hugo-build")
}

// Run tests and linters
func Check() {
	if runtime.GOARCH == "amd64" && runtime.GOOS != "darwin" {
		mg.Deps(Test386)
	} else {
		fmt.Printf("Skip Test386 on %s and/or %s\n", runtime.GOARCH, runtime.GOOS)
	}

	if isCi() && isDarwin() {
		// Skip on macOS in CI (disk space issues)
	} else {
		mg.Deps(Fmt, Vet)
	}

	// don't run two tests in parallel, they saturate the CPUs anyway, and running two
	// causes memory issues in CI.
	mg.Deps(TestRace)
}

func testGoFlags() string {
	if isCI() {
		return ""
	}

	return "-timeout=1m"
}

// Run tests in 32-bit mode
// Note that we don't run with the extended tag. Currently not supported in 32 bit.
func Test386() error {
	env := map[string]string{"GOARCH": "386", "GOFLAGS": testGoFlags()}
	return runCmd(env, goexe, "test", "-p", "2", "./...")
}

// Run tests
func Test() error {
	env := map[string]string{"GOFLAGS": testGoFlags()}
	return runCmd(env, goexe, "test", "-p", "2", "./...", "-tags", buildTags())
}

// Run tests with race detector
func TestRace() error {
	env := map[string]string{"GOFLAGS": testGoFlags()}
	return runCmd(env, goexe, "test", "-p", "2", "-race", "./...", "-tags", buildTags())
}

// Run gofmt linter
func Fmt() error {
	if !isGoLatest() && !isUnix() {
		return nil
	}
	s, err := sh.Output("./check_gofmt.sh")
	if err != nil {
		fmt.Println(s)
		return fmt.Errorf("gofmt needs to be run: %s", err)
	}

	return nil
}

var (
	pkgPrefixLen = len("github.com/gohugoio/hugo")
	pkgs         []string
	pkgsInit     sync.Once
)

func hugoPackages() ([]string, error) {
	var err error
	pkgsInit.Do(func() {
		var s string
		s, err = sh.Output(goexe, "list", "./...")
		if err != nil {
			return
		}
		pkgs = strings.Split(s, "\n")
		for i := range pkgs {
			pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
		}
	})
	return pkgs, err
}

// Run golint linter
func Lint() error {
	pkgs, err := hugoPackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		// We don't actually want to fail this target if we find golint errors,
		// so we don't pass -set_exit_status, but we still print out any failures.
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", pkg); err != nil {
			fmt.Printf("ERROR: running go lint on %q: %v\n", pkg, err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	return nil
}

func isCi() bool {
	return os.Getenv("CI") != ""
}

func isDarwin() bool {
	return runtime.GOOS == "darwin"
}

// Run go vet linter
func Vet() error {
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}

// Generate test coverage report
func TestCoverHTML() error {
	const (
		coverAll = "coverage-all.out"
		cover    = "coverage.out"
	)
	f, err := os.Create(coverAll)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write([]byte("mode: count")); err != nil {
		return err
	}
	pkgs, err := hugoPackages()
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if err := sh.Run(goexe, "test", "-coverprofile="+cover, "-covermode=count", pkg); err != nil {
			return err
		}
		b, err := os.ReadFile(cover)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		idx := bytes.Index(b, []byte{'\n'})
		b = b[idx+1:]
		if _, err := f.Write(b); err != nil {
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return sh.Run(goexe, "tool", "cover", "-html="+coverAll)
}

func runCmd(env map[string]string, cmd string, args ...any) error {
	if mg.Verbose() {
		return runWith(env, cmd, args...)
	}
	output, err := sh.OutputWith(env, cmd, argsToStrings(args...)...)
	if err != nil {
		fmt.Fprint(os.Stderr, output)
	}

	return err
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.21")
}

func isUnix() bool {
	return runtime.GOOS != "windows"
}

func isCI() bool {
	return os.Getenv("CI") != ""
}

func buildFlags() []string {
	if runtime.GOOS == "windows" {
		return []string{"-buildmode", "exe"}
	}
	return nil
}

func buildTags() string {
	// To build the extended Hugo SCSS/SASS enabled version, build with
	// HUGO_BUILD_TAGS=extended mage install etc.
	// To build without `hugo deploy` for smaller binary, use HUGO_BUILD_TAGS=nodeploy
	if envtags := os.Getenv("HUGO_BUILD_TAGS"); envtags != "" {
		return envtags
	}
	return "none"
}

func argsToStrings(v ...any) []string {
	var args []string
	for _, arg := range v {
		switch v := arg.(type) {
		case string:
			if v != "" {
				args = append(args, v)
			}
		case []string:
			if v != nil {
				args = append(args, v...)
			}
		default:
			panic("invalid type")
		}
	}

	return args
}
