// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
)

type commandeer struct {
	*deps.DepsCfg
	pathSpec   *helpers.PathSpec
	configured bool
}

func (c *commandeer) Set(key string, value interface{}) {
	if c.configured {
		panic("commandeer cannot be changed")
	}
	c.Cfg.Set(key, value)
}

// PathSpec lazily creates a new PathSpec, as all the paths must
// be configured before it is created.
func (c *commandeer) PathSpec() *helpers.PathSpec {
	c.configured = true
	if c.pathSpec == nil {
		c.pathSpec = helpers.NewPathSpec(c.Fs, c.Cfg)
	}
	return c.pathSpec
}

func newCommandeer(cfg *deps.DepsCfg) *commandeer {
	return &commandeer{DepsCfg: cfg}
}
