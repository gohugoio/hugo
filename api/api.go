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

func NewBuild() *build {
	return new(build)
}

// Build is the type that handles Building methods
type build struct {
	flags []string
}

// Run regenerates the website
func (b *build) Run() error {
	_, err := commands.Execute(b.flags)
	return err
}

// Set adds a value to the flags array
func (b *build) Set(key string, value interface{}) {
	// If the key doesn't begin with "-"
	if key[0] != '-' {
		key = "--" + key
	}

	b.flags = append([]string{key, value.(string)}, b.flags...)
}

func NewSite() *newSite {
	n := new(newSite)
	*n.cmd = *commands.NewSiteCmd
	return n
}

// NewSite creates new sites
type newSite struct {
	cmd  *cobra.Command
	path string
}

// Set adds a value to the flags array
func (n *newSite) Set(key string, value interface{}) {
	if key == "path" {
		n.path = key
		return
	}

	n.cmd.Flags().Set(key, value.(string))
}

// Make generates a new site
func (n *newSite) Make() error {
	return n.cmd.RunE(n.cmd, []string{n.path})
}

func NewContent() *newContent {
	n := new(newContent)
	*n.cmd = *commands.NewCmd
	return n
}

// NewContent is used to create new contents
type newContent struct {
	cmd  *cobra.Command
	path string
}

// Set adds a value to the flags array
func (n *newContent) Set(key string, value interface{}) {
	if key == "path" {
		n.path = key
		return
	}

	n.cmd.Flags().Set(key, value.(string))
}

// Make generates a new content
func (n *newContent) Make() error {
	return n.cmd.RunE(n.cmd, []string{n.path})
}
