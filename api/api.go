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
	b.flags = preppendFlag(b.flags, key, value.(string))
}

// NewSite generates a new site
func NewSite(path string, force bool, format string) {
	cmd := &cobra.Command{}

	if &force == nil {
		force = false
	}

	cmd.Flags().Bool("force", true, "")

	if &format == nil {
		format = "toml"
	}

	cmd.Flags().String("format", format, "")
	commands.NewSite(cmd, []string{path})
}

// NewContent is used to create new contents
type NewContent struct {
	cmd  *cobra.Command
	path string
}

// Set adds a value to the flags array
func (n *NewContent) Set(key string, value interface{}) {
	if key == "path" {
		n.path = key
		return
	}

	n.cmd.Flags().Set(key, value.(string))
}

// Make generates a new content
func (n *NewContent) Make(path string) {
	commands.NewSite(n.cmd, []string{n.path})
}

// Reset resets the current website and sets all settings to default. It can
// be useful if your application needs to run Hugo in different paths
func Reset() {
	commands.ClearSite()
	viper.Reset()
}

func preppendFlag(flags []string, key string, value string) []string {
	// If the key doesn't begin with "-"
	if key[0] != '-' {
		key = "--" + key
	}

	return append([]string{key, value}, flags...)
}
