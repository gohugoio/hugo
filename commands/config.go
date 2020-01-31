// Copyright 2015 The Hugo Authors. All rights reserved.
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
// limitations under the License.Print the version number of Hug

package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/modules"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ cmder = (*configCmd)(nil)

type configCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newConfigCmd() *configCmd {
	cc := &configCmd{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Print the site configuration",
		Long:  `Print the site configuration, both default and custom settings.`,
		RunE:  cc.printConfig,
	}

	printMountsCmd := &cobra.Command{
		Use:   "mounts",
		Short: "Print the configured file mounts",
		RunE:  cc.printMounts,
	}

	cmd.AddCommand(printMountsCmd)

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}

func (c *configCmd) printMounts(cmd *cobra.Command, args []string) error {
	cfg, err := initializeConfig(true, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return err
	}

	allModules := cfg.Cfg.Get("allmodules").(modules.Modules)

	for _, m := range allModules {
		if err := parser.InterfaceToConfig(&modMounts{m: m}, metadecoders.JSON, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func (c *configCmd) printConfig(cmd *cobra.Command, args []string) error {
	cfg, err := initializeConfig(true, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return err
	}

	allSettings := cfg.Cfg.(*viper.Viper).AllSettings()

	// We need to clean up this, but we store objects in the config that
	// isn't really interesting to the end user, so filter these.
	ignoreKeysRe := regexp.MustCompile("client|sorted|filecacheconfigs|allmodules|multilingual")

	separator := ": "

	if len(cfg.configFiles) > 0 && strings.HasSuffix(cfg.configFiles[0], ".toml") {
		separator = " = "
	}

	var keys []string
	for k := range allSettings {
		if ignoreKeysRe.MatchString(k) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		kv := reflect.ValueOf(allSettings[k])
		if kv.Kind() == reflect.String {
			fmt.Printf("%s%s\"%+v\"\n", k, separator, allSettings[k])
		} else {
			fmt.Printf("%s%s%+v\n", k, separator, allSettings[k])
		}
	}

	return nil
}

type modMounts struct {
	m modules.Module
}

type modMount struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Lang   string `json:"lang,omitempty"`
}

func (m *modMounts) MarshalJSON() ([]byte, error) {
	var mounts []modMount

	for _, mount := range m.m.Mounts() {
		mounts = append(mounts, modMount{
			Source: mount.Source,
			Target: mount.Target,
			Lang:   mount.Lang,
		})
	}

	return json.Marshal(&struct {
		Path   string     `json:"path"`
		Dir    string     `json:"dir"`
		Mounts []modMount `json:"mounts"`
	}{
		Path:   m.m.Path(),
		Dir:    m.m.Dir(),
		Mounts: mounts,
	})
}
