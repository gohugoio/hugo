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
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var _ cmder = (*configCmd)(nil)

type configCmd struct {
	hugoBuilderCommon
	*baseCmd
}

func newConfigCmd() *configCmd {
	cc := &configCmd{}
	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "config",
		Short: "Print the site configuration",
		Long:  `Print the site configuration, both default and custom settings.`,
		RunE:  cc.printConfig,
	})

	cc.cmd.Flags().StringVarP(&cc.source, "source", "s", "", "filesystem path to read files relative from")

	return cc
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
			jww.FEEDBACK.Printf("%s%s\"%+v\"\n", k, separator, allSettings[k])
		} else {
			jww.FEEDBACK.Printf("%s%s%+v\n", k, separator, allSettings[k])
		}
	}

	return nil
}
