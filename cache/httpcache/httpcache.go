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

package httpcache

import (
	"encoding/json"
	"time"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

// DefaultConfig holds the default configuration for the HTTP cache.
var DefaultConfig = Config{
	Cache: Cache{
		For: GlobMatcher{
			Excludes: []string{"**"},
		},
	},
	Polls: []PollConfig{
		{
			For: GlobMatcher{
				Includes: []string{"**"},
			},
			Disable: true,
		},
	},
}

// Config holds the configuration for the HTTP cache.
type Config struct {
	// Configures the HTTP cache behaviour (RFC 9111).
	// When this is not enabled for a resource, Hugo will go straight to the file cache.
	Cache Cache

	// Polls holds a list of configurations for polling remote resources to detect changes in watch mode.
	// This can be disabled for some resources, typically if they are known to not change.
	Polls []PollConfig
}

type Cache struct {
	// Enable HTTP cache behaviour (RFC 9111) for these rsources.
	For GlobMatcher
}

func (c *Config) Compile() (ConfigCompiled, error) {
	var cc ConfigCompiled

	p, err := c.Cache.For.CompilePredicate()
	if err != nil {
		return cc, err
	}

	cc.For = p

	for _, pc := range c.Polls {

		p, err := pc.For.CompilePredicate()
		if err != nil {
			return cc, err
		}

		cc.PollConfigs = append(cc.PollConfigs, PollConfigCompiled{
			For:    p,
			Config: pc,
		})
	}

	return cc, nil
}

// PollConfig holds the configuration for polling remote resources to detect changes in watch mode.
type PollConfig struct {
	// What remote resources to apply this configuration to.
	For GlobMatcher

	// Disable polling for this configuration.
	Disable bool

	// Low is the lower bound for the polling interval.
	// This is the starting point when the resource has recently changed,
	// if that resource stops changing, the polling interval will gradually increase towards High.
	Low time.Duration

	// High is the upper bound for the polling interval.
	// This is the interval used when the resource is stable.
	High time.Duration
}

func (c PollConfig) MarshalJSON() (b []byte, err error) {
	// Marshal the durations as strings.
	type Alias PollConfig
	return json.Marshal(&struct {
		Low  string
		High string
		Alias
	}{
		Low:   c.Low.String(),
		High:  c.High.String(),
		Alias: (Alias)(c),
	})
}

type GlobMatcher struct {
	// Excludes holds a list of glob patterns that will be excluded.
	Excludes []string

	// Includes holds a list of glob patterns that will be included.
	Includes []string
}

type ConfigCompiled struct {
	For         predicate.P[string]
	PollConfigs []PollConfigCompiled
}

func (c *ConfigCompiled) PollConfigFor(s string) PollConfigCompiled {
	for _, pc := range c.PollConfigs {
		if pc.For(s) {
			return pc
		}
	}
	return PollConfigCompiled{}
}

func (c *ConfigCompiled) IsPollingDisabled() bool {
	for _, pc := range c.PollConfigs {
		if !pc.Config.Disable {
			return false
		}
	}
	return true
}

type PollConfigCompiled struct {
	For    predicate.P[string]
	Config PollConfig
}

func (p PollConfigCompiled) IsZero() bool {
	return p.For == nil
}

func (gm *GlobMatcher) CompilePredicate() (func(string) bool, error) {
	var p predicate.P[string]
	for _, include := range gm.Includes {
		g, err := glob.Compile(include, '/')
		if err != nil {
			return nil, err
		}
		fn := func(s string) bool {
			return g.Match(s)
		}
		p = p.Or(fn)
	}

	for _, exclude := range gm.Excludes {
		g, err := glob.Compile(exclude, '/')
		if err != nil {
			return nil, err
		}
		fn := func(s string) bool {
			return !g.Match(s)
		}
		p = p.And(fn)
	}

	return p, nil
}

func DecodeConfig(bcfg config.BaseConfig, m map[string]any) (Config, error) {
	if len(m) == 0 {
		return DefaultConfig, nil
	}

	var c Config

	dc := &mapstructure.DecoderConfig{
		Result:           &c,
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return c, err
	}

	if err := decoder.Decode(m); err != nil {
		return c, err
	}

	return c, nil
}
