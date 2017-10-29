// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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

func getDep() error {
	return sh.Run(goexe, "get", "-u", "github.com/golang/dep/cmd/dep")
}

// Install Go Dep and sync Hugo's vendored dependencies
func Vendor() error {
	mg.Deps(getDep)
	return sh.Run("dep", "ensure")
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
	if strings.Contains(runtime.Version(), "1.8") {
		// Go 1.8 doesn't play along with go test ./... and /vendor.
		// We could fix that, but that would take time.
		fmt.Printf("Skip Check on %s\n", runtime.Version())
		return
	}
	mg.Deps(Test386, Fmt, Vet)
	// don't run two tests in parallel, they saturate the CPUs anyway, and running two
	// causes memory issues in CI.
	mg.Deps(TestRace)
}

// Run tests in 32-bit mode
func Test386() error {
	return sh.RunWith(map[string]string{"GOARCH": "386"}, goexe, "test", "./...")
}

// Run tests
func Test() error {
	mg.Deps(getDep)
	return sh.Run(goexe, "test", "./...")
}

// Run tests with race detector
func TestRace() error {
	mg.Deps(getDep)
	return sh.Run(goexe, "test", "-race", "./...")
}

// Run gofmt linter
func Fmt() error {
	pkgs, err := hugoPackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Println(s)
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
	mg.Deps(getDep)
	s, err := sh.Output(goexe, "list", "./...")
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

//  Run go vet linter
func Vet() error {
	mg.Deps(getDep)
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running govendor: %v", err)
	}
	return nil
}

// Generate test coverage report
func TestCoverHTML() error {
	mg.Deps(getDep)
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
