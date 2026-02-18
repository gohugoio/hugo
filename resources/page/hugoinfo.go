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

package page

import (
	"fmt"
	"html/template"

	"github.com/gohugoio/hugo/common/hstore"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/version"
)

var _ hstore.StoreProvider = (*HugoInfo)(nil)

type hugoInfoProviders struct {
	HugoInfoHugoSitesProvider
}

// HugoInfo contains information about the current Hugo environment.
type HugoInfo struct {
	CommitHash string
	BuildDate  string

	// Version of Go that the Hugo binary was built with.
	GoVersion string

	// Avoid exporting the embedded providers directly, as they may be used by the template layer.
	hugoInfoProviders

	// Context gives access to some of the context scoped variables.
	Context hugo.Context

	opts  HugoInfoOptions
	store *hstore.Scratch
}

// Version returns the current version as a comparable version string.
func (i HugoInfo) Version() version.VersionString {
	return hugo.CurrentVersion.Version()
}

// Generator returns a Hugo meta generator HTML tag.
func (i HugoInfo) Generator() template.HTML {
	return template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s">`, hugo.CurrentVersion.String()))
}

// Environment returns the build environment.
// Defaults are "production" (hugo) and "development" (hugo server).
// This can also be set by the user.
// It can be any string, but it will be all lower case.
func (i HugoInfo) Environment() string {
	return i.opts.Conf.Environment()
}

// IsDevelopment reports whether the current running environment is "development".
func (i HugoInfo) IsDevelopment() bool {
	return i.Environment() == hugo.EnvironmentDevelopment
}

// IsProduction reports whether the current running environment is "production".
func (i HugoInfo) IsProduction() bool {
	return i.Environment() == hugo.EnvironmentProduction
}

// IsServer reports whether the built-in server is running.
func (i HugoInfo) IsServer() bool {
	return i.opts.Conf.Running()
}

// IsExtended reports whether the Hugo binary is the extended version.
func (i HugoInfo) IsExtended() bool {
	return hugo.IsExtended
}

// WorkingDir returns the project working directory.
func (i HugoInfo) WorkingDir() string {
	return i.opts.Conf.WorkingDir()
}

// Deps gets a list of dependencies for this Hugo build.
func (i HugoInfo) Deps() []*hugo.Dependency {
	return i.opts.Deps
}

func (i HugoInfo) Store() *hstore.Scratch {
	return i.store
}

// Deprecated: Use hugo.IsMultihost instead.
func (i HugoInfo) IsMultiHost() bool {
	hugo.Deprecate("hugo.IsMultiHost", "Use hugo.IsMultihost instead.", "v0.124.0")
	return i.opts.Conf.IsMultihost()
}

// IsMultihost reports whether each configured language has a unique baseURL.
func (i HugoInfo) IsMultihost() bool {
	return i.opts.Conf.IsMultihost()
}

// IsMultilingual reports whether there are two or more configured languages.
func (i HugoInfo) IsMultilingual() bool {
	return i.opts.Conf.IsMultilingual()
}

// HugoInfoConfigProvider represents the config options that are relevant for HugoInfo.
type HugoInfoConfigProvider interface {
	Environment() string
	Running() bool
	WorkingDir() string
	IsMultihost() bool
	IsMultilingual() bool
}

// HugoInfoHugoSitesProvider provides what HugoInfo needs from hugolib.HugoSites.
type HugoInfoHugoSitesProvider interface {
	SitesProvider
	DataProvider
}

// HugoInfoOptions defines the providers required to initialize HugoInfo.
type HugoInfoOptions struct {
	Conf                      HugoInfoConfigProvider
	HugoInfoHugoSitesProvider HugoInfoHugoSitesProvider

	Deps       []*hugo.Dependency
	CommitHash string
	BuildDate  string
	GoVersion  string
}

// NewHugoInfo creates a new Hugo Info object.
func NewHugoInfo(opts HugoInfoOptions) HugoInfo {
	if opts.Conf == nil {
		panic("config provider not set")
	}
	if opts.Conf.Environment() == "" {
		panic("environment not set")
	}
	if opts.HugoInfoHugoSitesProvider == nil {
		opts.HugoInfoHugoSitesProvider = nopHugoSitesProvider{}
	}

	return HugoInfo{
		CommitHash: opts.CommitHash,
		BuildDate:  opts.BuildDate,
		GoVersion:  opts.GoVersion,

		hugoInfoProviders: hugoInfoProviders{HugoInfoHugoSitesProvider: opts.HugoInfoHugoSitesProvider},

		opts:  opts,
		store: hstore.NewScratch(),
	}
}
