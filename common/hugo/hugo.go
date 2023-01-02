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

package hugo

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bep/godartsass"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

var (
	// buildDate allows vendor-specified build date when .git/ is unavailable.
	buildDate string
	// vendorInfo contains vendor notes about the current build.
	vendorInfo string
)

// Info contains information about the current Hugo environment
type Info struct {
	CommitHash string
	BuildDate  string

	// The build environment.
	// Defaults are "production" (hugo) and "development" (hugo server).
	// This can also be set by the user.
	// It can be any string, but it will be all lower case.
	Environment string

	// version of go that the Hugo binary was built with
	GoVersion string

	deps []*Dependency
}

// Version returns the current version as a comparable version string.
func (i Info) Version() VersionString {
	return CurrentVersion.Version()
}

// Generator a Hugo meta generator HTML tag.
func (i Info) Generator() template.HTML {
	return template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s">`, CurrentVersion.String()))
}

func (i Info) IsProduction() bool {
	return i.Environment == EnvironmentProduction
}

func (i Info) IsExtended() bool {
	return IsExtended
}

// Deps gets a list of dependencies for this Hugo build.
func (i Info) Deps() []*Dependency {
	return i.deps
}

// NewInfo creates a new Hugo Info object.
func NewInfo(environment string, deps []*Dependency) Info {
	if environment == "" {
		environment = EnvironmentProduction
	}
	var (
		commitHash string
		buildDate  string
		goVersion  string
	)

	bi := getBuildInfo()
	if bi != nil {
		commitHash = bi.Revision
		buildDate = bi.RevisionTime
		goVersion = bi.GoVersion
	}

	return Info{
		CommitHash:  commitHash,
		BuildDate:   buildDate,
		Environment: environment,
		deps:        deps,
		GoVersion:   goVersion,
	}
}

// GetExecEnviron creates and gets the common os/exec environment used in the
// external programs we interact with via os/exec, e.g. postcss.
func GetExecEnviron(workDir string, cfg config.Provider, fs afero.Fs) []string {
	var env []string
	nodepath := filepath.Join(workDir, "node_modules")
	if np := os.Getenv("NODE_PATH"); np != "" {
		nodepath = workDir + string(os.PathListSeparator) + np
	}
	config.SetEnvVars(&env, "NODE_PATH", nodepath)
	config.SetEnvVars(&env, "PWD", workDir)
	config.SetEnvVars(&env, "HUGO_ENVIRONMENT", cfg.GetString("environment"))
	config.SetEnvVars(&env, "HUGO_ENV", cfg.GetString("environment"))

	config.SetEnvVars(&env, "HUGO_PUBLISHDIR", filepath.Join(workDir, cfg.GetString("publishDirOrig")))

	if fs != nil {
		fis, err := afero.ReadDir(fs, files.FolderJSConfig)
		if err == nil {
			for _, fi := range fis {
				key := fmt.Sprintf("HUGO_FILE_%s", strings.ReplaceAll(strings.ToUpper(fi.Name()), ".", "_"))
				value := fi.(hugofs.FileMetaInfo).Meta().Filename
				config.SetEnvVars(&env, key, value)
			}
		}
	}

	return env
}

type buildInfo struct {
	VersionControlSystem string
	Revision             string
	RevisionTime         string
	Modified             bool

	GoOS   string
	GoArch string

	*debug.BuildInfo
}

var bInfo *buildInfo
var bInfoInit sync.Once

func getBuildInfo() *buildInfo {
	bInfoInit.Do(func() {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}

		bInfo = &buildInfo{BuildInfo: bi}

		for _, s := range bInfo.Settings {
			switch s.Key {
			case "vcs":
				bInfo.VersionControlSystem = s.Value
			case "vcs.revision":
				bInfo.Revision = s.Value
			case "vcs.time":
				bInfo.RevisionTime = s.Value
			case "vcs.modified":
				bInfo.Modified = s.Value == "true"
			case "GOOS":
				bInfo.GoOS = s.Value
			case "GOARCH":
				bInfo.GoArch = s.Value
			}
		}

	})

	return bInfo
}

func formatDep(path, version string) string {
	return fmt.Sprintf("%s=%q", path, version)
}

// GetDependencyList returns a sorted dependency list on the format package="version".
// It includes both Go dependencies and (a manually maintained) list of C(++) dependencies.
func GetDependencyList() []string {
	var deps []string

	bi := getBuildInfo()
	if bi == nil {
		return deps
	}

	for _, dep := range bi.Deps {
		deps = append(deps, formatDep(dep.Path, dep.Version))
	}

	deps = append(deps, GetDependencyListNonGo()...)

	sort.Strings(deps)

	return deps
}

// GetDependencyListNonGo returns a list of non-Go dependencies.
func GetDependencyListNonGo() []string {
	var deps []string

	if IsExtended {
		deps = append(
			deps,
			formatDep("github.com/sass/libsass", "3.6.5"),
			formatDep("github.com/webmproject/libwebp", "v1.2.4"),
		)
	}

	if dartSass := dartSassVersion(); dartSass.ProtocolVersion != "" {
		const dartSassPath = "github.com/sass/dart-sass-embedded"
		deps = append(deps,
			formatDep(dartSassPath+"/protocol", dartSass.ProtocolVersion),
			formatDep(dartSassPath+"/compiler", dartSass.CompilerVersion),
			formatDep(dartSassPath+"/implementation", dartSass.ImplementationVersion),
		)
	}
	return deps
}

// IsRunningAsTest reports whether we are running as a test.
func IsRunningAsTest() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			return true
		}
	}
	return false
}

// Dependency is a single dependency, which can be either a Hugo Module or a local theme.
type Dependency struct {
	// Returns the path to this module.
	// This will either be the module path, e.g. "github.com/gohugoio/myshortcodes",
	// or the path below your /theme folder, e.g. "mytheme".
	Path string

	// The module version.
	Version string

	// Whether this dependency is vendored.
	Vendor bool

	// Time version was created.
	Time time.Time

	// In the dependency tree, this is the first module that defines this module
	// as a dependency.
	Owner *Dependency

	// Replaced by this dependency.
	Replace *Dependency
}

func dartSassVersion() godartsass.DartSassVersion {
	// This is also duplicated in the dartsass package.
	const dartSassEmbeddedBinaryName = "dart-sass-embedded"
	if !hexec.InPath(dartSassEmbeddedBinaryName) {
		return godartsass.DartSassVersion{}
	}
	v, _ := godartsass.Version(dartSassEmbeddedBinaryName)
	return v
}
