// +build mage

package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName  = "github.com/gohugoio/hugo"
	noGitLdflags = "-X $PACKAGE/hugolib.BuildDate=$BUILD_DATE"
)

var ldflags = "-X $PACKAGE/hugolib.CommitHash=$COMMIT_HASH -X $PACKAGE/hugolib.BuildDate=$BUILD_DATE"

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}
}

func govendor() error {
	return sh.Run(goexe, "get", "github.com/kardianos/govendor")
}

// Install govendor and sync Hugo's vendored dependencies
func Vendor() error {
	mg.Deps(govendor)
	return sh.Run("govendor", "sync", packageName)
}

// Build hugo binary
func Hugo() error {
	mg.Deps(Vendor)
	return sh.RunWith(flagEnv(), goexe, "build", "-ldflags", ldflags, packageName)
}

// Build hugo binary with race detector enabled
func HugoRace() error {
	mg.Deps(Vendor)
	return sh.RunWith(flagEnv(), goexe, "build", "-race", "-ldflags", ldflags, packageName)
}

// Install hugo binary
func Install() error {
	mg.Deps(Vendor)
	return sh.RunWith(flagEnv(), goexe, "install", "-ldflags", ldflags, packageName)
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
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
	mg.Deps(TestRace, Test386, Fmt, Vet)
}

// Run tests in 32-bit mode
func Test386() error {
	return sh.RunWith(map[string]string{"GOARCH": "386"}, "govendor", "test", "+local")
}

// Run tests
func Test() error {
	mg.Deps(govendor)
	return sh.Run("govendor", "test", "+local")
}

// Run tests with race detector
func TestRace() error {
	mg.Deps(govendor)
	return sh.Run("govendor", "test", "-race", "+local")
}

// Run gofmt linter
func Fmt() error {
	pkgs, err := hugoPackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			if err := sh.Run("gofmt", "-l", f); err != nil {
				failed = false
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}

var pkgPrefixLen = len("github.com/gohugoio/hugo")

func hugoPackages() ([]string, error) {
	mg.Deps(govendor)
	s, err := sh.Output("govendor", "list", "-no-status", "+local")
	if err != nil {
		return nil, err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
	}
	return pkgs, nil
}

// Run golint linter
func Lint() error {
	pkgs, err := hugoPackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		if _, err := sh.Exec(nil, os.Stderr, os.Stderr, "golint", "-set_exit_status", pkg); err != nil {
			failed = true
		}
	}
	if failed {
		return errors.New("golint errors!")
	}
	return nil
}

//  Run go vet linter
func Vet() error {
	mg.Deps(govendor)
	if err := sh.Run("govendor", "vet", "+local"); err != nil {
		return errors.New("go vet errors!")
	}
	return nil
}

// Generate test coverage report
func TestCoverHTML() error {
	mg.Deps(govendor)
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
		if err := sh.Run("govendor", "test", "-coverprofile="+cover, "-covermode=count", pkg); err != nil {
			return err
		}
		b, err := ioutil.ReadFile(cover)
		if err != nil {
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

// Verify that vendored packages match git HEAD
func CheckVendor() error {
	if err := sh.Run("git", "diff-index", "--quiet", "HEAD", "vendor/"); err != nil {
		// yes, ignore errors from this, not much we can do.
		sh.Exec(nil, os.Stdout, os.Stderr, "git", "diff", "vendor/")
		return errors.New("check-vendor target failed: vendored packages out of sync")
	}
	return nil
}
