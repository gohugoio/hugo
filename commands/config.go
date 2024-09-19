// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config/allconfig"
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

	format string
	lang   string

	commands []simplecobra.Commander
}

func (c *configCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *configCommand) Name() string {
	return "config"
}

func (c *configCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	conf, err := c.r.ConfigFromProvider(configKey{counter: c.r.configVersionID.Load()}, flagsToCfg(cd, nil))
	if err != nil {
		return err
	}
	var config *allconfig.Config
	if c.lang != "" {
		var found bool
		config, found = conf.configs.LanguageConfigMap[c.lang]
		if !found {
			return fmt.Errorf("language %q not found", c.lang)
		}
	} else {
		config = conf.configs.LanguageConfigSlice[0]
	}

	var buf bytes.Buffer
	dec := json.NewEncoder(&buf)
	dec.SetIndent("", "  ")
	dec.SetEscapeHTML(false)

	if err := dec.Encode(parser.ReplacingJSONMarshaller{Value: config, KeysToLower: true, OmitEmpty: true}); err != nil {
		return err
	}

	format := strings.ToLower(c.format)

	switch format {
	case "json":
		os.Stdout.Write(buf.Bytes())
	default:
		// Decode the JSON to a map[string]interface{} and then unmarshal it again to the correct format.
		var m map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
			return err
		}
		maps.ConvertFloat64WithNoDecimalsToInt(m)
		switch format {
		case "yaml":
			return parser.InterfaceToConfig(m, metadecoders.YAML, os.Stdout)
		case "toml":
			return parser.InterfaceToConfig(m, metadecoders.TOML, os.Stdout)
		default:
			return fmt.Errorf("unsupported format: %q", format)
		}
	}

	return nil
}

func (c *configCommand) Init(cd *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	cmd := cd.CobraCommand
	cmd.Short = "Print the site configuration"
	cmd.Long = `Print the site configuration, both default and custom settings.`
	cmd.Flags().StringVar(&c.format, "format", "toml", "preferred file format (toml, yaml or json)")
	_ = cmd.RegisterFlagCompletionFunc("format", cobra.FixedCompletions([]string{"toml", "yaml", "json"}, cobra.ShellCompDirectiveNoFileComp))
	cmd.Flags().StringVar(&c.lang, "lang", "", "the language to display config for. Defaults to the first language defined.")
	_ = cmd.RegisterFlagCompletionFunc("lang", cobra.NoFileCompletions)
	applyLocalFlagsBuildConfig(cmd, c.r)

	return nil
}

func (c *configCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
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
	r         *rootCommand
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
	conf, err := r.ConfigFromProvider(configKey{counter: c.r.configVersionID.Load()}, flagsToCfg(cd, nil))
	if err != nil {
		return err
	}

	for _, m := range conf.configs.Modules {
		if err := parser.InterfaceToConfig(&configModMounts{m: m, verbose: r.isVerbose()}, metadecoders.JSON, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func (c *configMountsCommand) Init(cd *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	cmd := cd.CobraCommand
	cmd.Short = "Print the configured file mounts"
	cmd.ValidArgsFunction = cobra.NoFileCompletions
	applyLocalFlagsBuildConfig(cmd, c.r)
	return nil
}

func (c *configMountsCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.configCmd = cd.Parent.Command.(*configCommand)
	return nil
}
