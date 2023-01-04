// Copyright 2023 The Hugo Authors. All rights reserved.
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

package commands

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/cobra"
)

// newConfigCommand creates a new config command and its subcommands.
func newConfigCommand() *configCommand {
	return &configCommand{
		commands: []simplecobra.Commander{
			&configMountsCommand{},
		},
	}

}

type configCommand struct {
	r *rootCommand

	commands []simplecobra.Commander
}

func (c *configCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *configCommand) Name() string {
	return "config"
}

func (c *configCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	conf, err := c.r.ConfigFromProvider(c.r.configVersionID.Load(), flagsToCfg(cd, nil))
	if err != nil {
		return err
	}
	config := conf.configs.Base

	// Print it as JSON.
	dec := json.NewEncoder(os.Stdout)
	dec.SetIndent("", "  ")
	dec.SetEscapeHTML(false)

	if err := dec.Encode(parser.ReplacingJSONMarshaller{Value: config, KeysToLower: true, OmitEmpty: true}); err != nil {
		return err
	}
	return nil
}

func (c *configCommand) WithCobraCommand(cmd *cobra.Command) error {
	cmd.Short = "Print the site configuration"
	cmd.Long = `Print the site configuration, both default and custom settings.`
	return nil
}

func (c *configCommand) Init(cd, runner *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	return nil
}

type configModMount struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Lang   string `json:"lang,omitempty"`
}

type configModMounts struct {
	verbose bool
	m       modules.Module
}

// MarshalJSON is for internal use only.
func (m *configModMounts) MarshalJSON() ([]byte, error) {
	var mounts []configModMount

	for _, mount := range m.m.Mounts() {
		mounts = append(mounts, configModMount{
			Source: mount.Source,
			Target: mount.Target,
			Lang:   mount.Lang,
		})
	}

	var ownerPath string
	if m.m.Owner() != nil {
		ownerPath = m.m.Owner().Path()
	}

	if m.verbose {
		config := m.m.Config()
		return json.Marshal(&struct {
			Path        string              `json:"path"`
			Version     string              `json:"version"`
			Time        time.Time           `json:"time"`
			Owner       string              `json:"owner"`
			Dir         string              `json:"dir"`
			Meta        map[string]any      `json:"meta"`
			HugoVersion modules.HugoVersion `json:"hugoVersion"`

			Mounts []configModMount `json:"mounts"`
		}{
			Path:        m.m.Path(),
			Version:     m.m.Version(),
			Time:        m.m.Time(),
			Owner:       ownerPath,
			Dir:         m.m.Dir(),
			Meta:        config.Params,
			HugoVersion: config.HugoVersion,
			Mounts:      mounts,
		})
	}

	return json.Marshal(&struct {
		Path    string           `json:"path"`
		Version string           `json:"version"`
		Time    time.Time        `json:"time"`
		Owner   string           `json:"owner"`
		Dir     string           `json:"dir"`
		Mounts  []configModMount `json:"mounts"`
	}{
		Path:    m.m.Path(),
		Version: m.m.Version(),
		Time:    m.m.Time(),
		Owner:   ownerPath,
		Dir:     m.m.Dir(),
		Mounts:  mounts,
	})

}

type configMountsCommand struct {
	configCmd *configCommand
}

func (c *configMountsCommand) Commands() []simplecobra.Commander {
	return nil
}

func (c *configMountsCommand) Name() string {
	return "mounts"
}

func (c *configMountsCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	r := c.configCmd.r
	conf, err := r.ConfigFromProvider(r.configVersionID.Load(), flagsToCfg(cd, nil))
	if err != nil {
		return err
	}

	for _, m := range conf.configs.Modules {
		if err := parser.InterfaceToConfig(&configModMounts{m: m, verbose: r.verbose}, metadecoders.JSON, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func (c *configMountsCommand) WithCobraCommand(cmd *cobra.Command) error {
	cmd.Short = "Print the configured file mounts"
	return nil
}

func (c *configMountsCommand) Init(cd, runner *simplecobra.Commandeer) error {
	c.configCmd = cd.Parent.Command.(*configCommand)
	return nil
}
