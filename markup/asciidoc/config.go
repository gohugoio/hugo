// Copyright 2019 The Hugo Authors. All rights reserved.
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

package asciidoc

import (
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

// DecodeConfig creates a modules Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (Config, error) {
	c := Config{
		Args:           []string{"--no-header-footer", "--safe", "--trace"},
		CurrentContent: false,
	}

	if cfg == nil {
		return c, nil
	}

	asciidoctorSet := cfg.IsSet("asciidoctor")

	if asciidoctorSet {
		m := cfg.GetStringMap("asciidoctor")
		if err := mapstructure.WeakDecode(m, &c); err != nil {
			return c, err
		}
	}

	return c, nil
}

// Config holds a module config.
type Config struct {
	Args           []string
	CurrentContent bool
}
