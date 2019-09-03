// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package modules provides a client that can be used to manage Hugo Components,
// what's referred to as Hugo Modules. Hugo Modules is built on top of Go Modules,
// but also supports vendoring and components stored directly in the themes dir.
package modules

import (
	"github.com/gohugoio/hugo/config"
)

var _ Module = (*moduleAdapter)(nil)

type Module interface {

	// Optional config read from the configFilename above.
	Cfg() config.Provider

	// The decoded module config and mounts.
	Config() Config

	// Optional configuration filename (e.g. "/themes/mytheme/config.json").
	// This will be added to the special configuration watch list when in
	// server mode.
	ConfigFilename() string

	// Directory holding files for this module.
	Dir() string

	// This module is disabled.
	Disabled() bool

	// Returns whether this is a Go Module.
	IsGoMod() bool

	// Any directory remappings.
	Mounts() []Mount

	// In the dependency tree, this is the first module that defines this module
	// as a dependency.
	Owner() Module

	// Returns the path to this module.
	// This will either be the module path, e.g. "github.com/gohugoio/myshortcodes",
	// or the path below your /theme folder, e.g. "mytheme".
	Path() string

	// Replaced by this module.
	Replace() Module

	// Returns whether Dir points below the _vendor dir.
	Vendor() bool

	// The module version.
	Version() string

	// Whether this module's dir is a watch candidate.
	Watch() bool
}

type Modules []Module

type moduleAdapter struct {
	path       string
	dir        string
	version    string
	vendor     bool
	disabled   bool
	projectMod bool
	owner      Module

	mounts []Mount

	configFilename string
	cfg            config.Provider
	config         Config

	// Set if a Go module.
	gomod *goModule
}

func (m *moduleAdapter) Cfg() config.Provider {
	return m.cfg
}

func (m *moduleAdapter) Config() Config {
	return m.config
}

func (m *moduleAdapter) ConfigFilename() string {
	return m.configFilename
}

func (m *moduleAdapter) Dir() string {
	// This may point to the _vendor dir.
	if !m.IsGoMod() || m.dir != "" {
		return m.dir
	}
	return m.gomod.Dir
}

func (m *moduleAdapter) Disabled() bool {
	return m.disabled
}

func (m *moduleAdapter) IsGoMod() bool {
	return m.gomod != nil
}

func (m *moduleAdapter) Mounts() []Mount {
	return m.mounts
}

func (m *moduleAdapter) Owner() Module {
	return m.owner
}

func (m *moduleAdapter) Path() string {
	if !m.IsGoMod() || m.path != "" {
		return m.path
	}
	return m.gomod.Path
}

func (m *moduleAdapter) Replace() Module {
	if m.IsGoMod() && !m.Vendor() && m.gomod.Replace != nil {
		return &moduleAdapter{
			gomod: m.gomod.Replace,
			owner: m.owner,
		}
	}
	return nil
}

func (m *moduleAdapter) Vendor() bool {
	return m.vendor
}

func (m *moduleAdapter) Version() string {
	if !m.IsGoMod() || m.version != "" {
		return m.version
	}
	return m.gomod.Version
}

func (m *moduleAdapter) Watch() bool {
	if m.Owner() == nil {
		// Main project
		return true
	}

	if !m.IsGoMod() {
		// Module inside /themes
		return true
	}

	if m.Replace() != nil {
		// Version is not set when replaced by a local folder.
		return m.Replace().Version() == ""
	}

	return false
}
