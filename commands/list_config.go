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
	"sort"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the site configuration",
	Long:  `Print the site configuration, both default and custom settings.`,
}

func init() {
	configCmd.RunE = printConfig
}

func printConfig(cmd *cobra.Command, args []string) error {
	cfg, err := InitializeConfig(configCmd)

	if err != nil {
		return err
	}

	allSettings := cfg.Cfg.(*viper.Viper).AllSettings()

	var separator string
	if allSettings["metadataformat"] == "toml" {
		separator = " = "
	} else {
		separator = ": "
	}

	var keys []string
	for k := range allSettings {
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
