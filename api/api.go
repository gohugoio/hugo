// Copyright 2016 The Hugo Authors. All rights reserved.
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

package api

import (
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/commands"
	"github.com/spf13/viper"
)

// Reset resets the current website and sets all settings to default. It can
// be useful if your application needs to run Hugo in different paths
func Reset() {
	commands.ClearSite()
	viper.Reset()
}

// Build is the type that handles Building methods
type Build struct {
	flags []string
}

// Run regenerates the website
func (b *Build) Run() error {
	_, err := commands.Execute(b.flags)
	return err
}

// Set adds a value to the flags array
func (b *Build) Set(key string, value interface{}) {
	// If the key doesn't begin with "-"
	if key[0] != '-' {
		key = "--" + key
	}

	b.flags = append([]string{key, value.(string)}, b.flags...)
}

// NewSite creates new sites
type NewSite struct {
	path string
}

// Set adds a value to the flags array
func (n *NewSite) Set(key, value string) {
	if key == "path" {
		n.path = value
		return
	}

	commands.NewSiteCmd.Flags().Set(key, value)
}

// Make generates a new site
func (n *NewSite) Make() error {
	return commands.NewSiteCmd.RunE(commands.NewSiteCmd, []string{n.path})
}

// NewContent is used to create new contents
type NewContent struct {
	cmd  *cobra.Command
	path string
}

// Set adds a value to the flags array
func (n *NewContent) Set(key, value string) {
	if key == "path" {
		n.path = value
		return
	}

	commands.NewCmd.Flags().Set(key, value)
}

// Make generates a new content
func (n *NewContent) Make() error {
	return commands.NewCmd.RunE(commands.NewCmd, []string{n.path})
}
