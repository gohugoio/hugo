// Copyright © 2013-15 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.Print the version number of Hug

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
	"sort"
)

var config = &cobra.Command{
	Use:   "config",
	Short: "Print the site configuration",
	Long:  `Print the site configuration, both default and custom settings`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		allSettings := viper.AllSettings()

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
				fmt.Printf("%s%s\"%+v\"\n", k, separator, allSettings[k])
			} else {
				fmt.Printf("%s%s%+v\n", k, separator, allSettings[k])
			}
		}
	},
}
