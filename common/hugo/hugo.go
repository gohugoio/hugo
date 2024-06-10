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

	godartsassv1 "github.com/bep/godartsass"
	"github.com/bep/logg"
	"github.com/mitchellh/mapstructure"

	"github.com/bep/godartsass/v2"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/spf13/afero"

	iofs "io/fs"

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

// HugoInfo contains information about the current Hugo environment
type HugoInfo struct {
	CommitHash string
	BuildDate  string

	// The build environment.
	// Defaults are "production" (hugo) and "development" (hugo server).
	// This can also be set by the user.
	// It can be any string, but it will be all lower case.
	Environment string

	// version of go that the Hugo binary was built with
	GoVersion string

	conf ConfigProvider
	deps []*Dependency
}

// Version returns the current version as a comparable version string.
func (i HugoInfo) Version() VersionString {
	return CurrentVersion.Version()
}

// Generator a Hugo meta generator HTML tag.
func (i HugoInfo) Generator() template.HTML {
	return template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s">`, CurrentVersion.String()))
}

// IsDevelopment reports whether the current running environment is "development".
func (i HugoInfo) IsDevelopment() bool {
	return i.Environment == EnvironmentDevelopment
}

// IsProduction reports whether the current running environment is "production".
func (i HugoInfo) IsProduction() bool {
	return i.Environment == EnvironmentProduction
}

// IsServer reports whether the built-in server is running.
func (i HugoInfo) IsServer() bool {
	return i.conf.Running()
}

// IsExtended reports whether the Hugo binary is the extended version.
func (i HugoInfo) IsExtended() bool {
	return IsExtended
}

// WorkingDir returns the project working directory.
func (i HugoInfo) WorkingDir() string {
	return i.conf.WorkingDir()
}

// Deps gets a list of dependencies for this Hugo build.
func (i HugoInfo) Deps() []*Dependency {
	return i.deps
}

// Deprecated: Use hugo.IsMultihost instead.
func (i HugoInfo) IsMultiHost() bool {
	Deprecate("hugo.IsMultiHost", "Use hugo.IsMultihost instead.", "v0.124.0")
	return i.conf.IsMultihost()
}

// IsMultihost reports whether each configured language has a unique baseURL.
func (i HugoInfo) IsMultihost() bool {
	return i.conf.IsMultihost()
}

// IsMultilingual reports whether there are two or more configured languages.
func (i HugoInfo) IsMultilingual() bool {
	return i.conf.IsMultilingual()
}

// ConfigProvider represents the config options that are relevant for HugoInfo.
type ConfigProvider interface {
	Environment() string
	Running() bool
	WorkingDir() string
	IsMultihost() bool
	IsMultilingual() bool
}

// NewInfo creates a new Hugo Info object.
func NewInfo(conf ConfigProvider, deps []*Dependency) HugoInfo {
	if conf.Environment() == "" {
		panic("environment not set")
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

	return HugoInfo{
		CommitHash:  commitHash,
		BuildDate:   buildDate,
		Environment: conf.Environment(),
		conf:        conf,
		deps:        deps,
		GoVersion:   goVersion,
	}
}

// GetExecEnviron creates and gets the common os/exec environment used in the
// external programs we interact with via os/exec, e.g. postcss.
func GetExecEnviron(workDir string, cfg config.AllProvider, fs afero.Fs) []string {
	var env []string
	nodepath := filepath.Join(workDir, "node_modules")
	if np := os.Getenv("NODE_PATH"); np != "" {
		nodepath = workDir + string(os.PathListSeparator) + np
	}
	config.SetEnvVars(&env, "NODE_PATH", nodepath)
	config.SetEnvVars(&env, "PWD", workDir)
	config.SetEnvVars(&env, "HUGO_ENVIRONMENT", cfg.Environment())
	config.SetEnvVars(&env, "HUGO_ENV", cfg.Environment())
	config.SetEnvVars(&env, "HUGO_PUBLISHDIR", filepath.Join(workDir, cfg.BaseConfig().PublishDir))

	if fs != nil {
		var fis []iofs.DirEntry
		d, err := fs.Open(files.FolderJSConfig)
		if err == nil {
			fis, err = d.(iofs.ReadDirFile).ReadDir(-1)
		}

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

var (
	bInfo     *buildInfo
	bInfoInit sync.Once
)

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
			formatDep("github.com/webmproject/libwebp", "v1.3.2"),
		)
	}

	if dartSass := dartSassVersion(); dartSass.ProtocolVersion != "" {
		dartSassPath := "github.com/sass/dart-sass-embedded"
		if IsDartSassV2() {
			dartSassPath = "github.com/sass/dart-sass"
		}
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
	if DartSassBinaryName == "" {
		return godartsass.DartSassVersion{}
	}
	if IsDartSassV2() {
		v, _ := godartsass.Version(DartSassBinaryName)
		return v
	}

	v, _ := godartsassv1.Version(DartSassBinaryName)
	var vv godartsass.DartSassVersion
	mapstructure.WeakDecode(v, &vv)
	return vv
}

// DartSassBinaryName is the name of the Dart Sass binary to use.
// TODO(beop) find a better place for this.
var DartSassBinaryName string

func init() {
	DartSassBinaryName = os.Getenv("DART_SASS_BINARY")
	if DartSassBinaryName == "" {
		for _, name := range dartSassBinaryNamesV2 {
			if hexec.InPath(name) {
				DartSassBinaryName = name
				break
			}
		}
		if DartSassBinaryName == "" {
			if hexec.InPath(dartSassBinaryNameV1) {
				DartSassBinaryName = dartSassBinaryNameV1
			}
		}
	}
}

var (
	dartSassBinaryNameV1  = "dart-sass-embedded"
	dartSassBinaryNamesV2 = []string{"dart-sass", "sass"}
)

func IsDartSassV2() bool {
	return !strings.Contains(DartSassBinaryName, "embedded")
}

// Deprecate informs about a deprecation starting at the given version.
//
// A deprecation typically needs a simple change in the template, but doing so will make the template incompatible with older versions.
// Theme maintainers generally want
// 1. No warnings or errors in the console when building a Hugo site.
// 2. Their theme to work for at least the last few Hugo versions.
func Deprecate(item, alternative string, version string) {
	level := deprecationLogLevelFromVersion(version)
	DeprecateLevel(item, alternative, version, level)
}

// DeprecateLevel informs about a deprecation logging at the given level.
func DeprecateLevel(item, alternative, version string, level logg.Level) {
	var msg string
	if level == logg.LevelError {
		msg = fmt.Sprintf("%s was deprecated in Hugo %s and will be removed in Hugo %s. %s", item, version, CurrentVersion.Next().ReleaseVersion(), alternative)
	} else {
		msg = fmt.Sprintf("%s was deprecated in Hugo %s and will be removed in a future release. %s", item, version, alternative)
	}

	loggers.Log().Logger().WithLevel(level).WithField(loggers.FieldNameCmd, "deprecated").Logf(msg)
}

// We ususally do about one minor version a month.
// We want people to run at least the current and previous version without any warnings.
// We want people who don't update Hugo that often to see the warnings and errors before we remove the feature.
func deprecationLogLevelFromVersion(ver string) logg.Level {
	from := MustParseVersion(ver)
	to := CurrentVersion
	minorDiff := to.Minor - from.Minor
	switch {
	case minorDiff >= 12:
		// Start failing the build after about a year.
		return logg.LevelError
	case minorDiff >= 6:
		// Start printing warnings after about six months.
		return logg.LevelWarn
	default:
		return logg.LevelInfo
	}
}
