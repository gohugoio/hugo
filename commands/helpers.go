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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/pflag"
)

const (
	ansiEsc    = "\u001B"
	clearLine  = "\r\033[K"
	hideCursor = ansiEsc + "[?25l"
	showCursor = ansiEsc + "[?25h"
)

func newUserError(a ...any) *simplecobra.CommandError {
	return &simplecobra.CommandError{Err: errors.New(fmt.Sprint(a...))}
}

func setValueFromFlag(flags *pflag.FlagSet, key string, cfg config.Provider, targetKey string, force bool) {
	key = strings.TrimSpace(key)
	if (force && flags.Lookup(key) != nil) || flags.Changed(key) {
		f := flags.Lookup(key)
		configKey := key
		if targetKey != "" {
			configKey = targetKey
		}
		// Gotta love this API.
		switch f.Value.Type() {
		case "bool":
			bv, _ := flags.GetBool(key)
			cfg.Set(configKey, bv)
		case "string":
			cfg.Set(configKey, f.Value.String())
		case "stringSlice":
			bv, _ := flags.GetStringSlice(key)
			cfg.Set(configKey, bv)
		case "int":
			iv, _ := flags.GetInt(key)
			cfg.Set(configKey, iv)
		default:
			panic(fmt.Sprintf("update switch with %s", f.Value.Type()))
		}

	}
}

func flagsToCfg(cd *simplecobra.Commandeer, cfg config.Provider) config.Provider {
	return flagsToCfgWithAdditionalConfigBase(cd, cfg, "")
}

func flagsToCfgWithAdditionalConfigBase(cd *simplecobra.Commandeer, cfg config.Provider, additionalConfigBase string) config.Provider {
	if cfg == nil {
		cfg = config.New()
	}

	// Flags with a different name in the config.
	keyMap := map[string]string{
		"minify":      "minifyOutput",
		"destination": "publishDir",
		"editor":      "newContentEditor",
	}

	// Flags that we for some reason don't want to expose in the site config.
	internalKeySet := map[string]bool{
		"quiet":          true,
		"verbose":        true,
		"watch":          true,
		"liveReloadPort": true,
		"renderToMemory": true,
		"clock":          true,
	}

	cmd := cd.CobraCommand
	flags := cmd.Flags()

	flags.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			targetKey := f.Name
			if internalKeySet[targetKey] {
				targetKey = "internal." + targetKey
			} else if mapped, ok := keyMap[targetKey]; ok {
				targetKey = mapped
			}
			setValueFromFlag(flags, f.Name, cfg, targetKey, false)
			if additionalConfigBase != "" {
				setValueFromFlag(flags, f.Name, cfg, additionalConfigBase+"."+targetKey, true)
			}
		}
	})

	return cfg
}

func mkdir(x ...string) {
	p := filepath.Join(x...)
	err := os.MkdirAll(p, 0o777) // before umask
	if err != nil {
		log.Fatal(err)
	}
}
