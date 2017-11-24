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
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
)

type commandeer struct {
	*deps.DepsCfg
	pathSpec    *helpers.PathSpec
	visitedURLs *types.EvictingStringQueue

	serverPorts []int

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
	return c.pathSpec
}

func (c *commandeer) languages() helpers.Languages {
	return c.Cfg.Get("languagesSorted").(helpers.Languages)
}

func (c *commandeer) initFs(fs *hugofs.Fs) error {
	c.DepsCfg.Fs = fs
	ps, err := helpers.NewPathSpec(fs, c.Cfg)
	if err != nil {
		return err
	}
	c.pathSpec = ps
	return nil
}

func newCommandeer(cfg *deps.DepsCfg) (*commandeer, error) {
	l := cfg.Language
	if l == nil {
		l = helpers.NewDefaultLanguage(cfg.Cfg)
	}
	ps, err := helpers.NewPathSpec(cfg.Fs, l)
	if err != nil {
		return nil, err
	}

	return &commandeer{DepsCfg: cfg, pathSpec: ps, visitedURLs: types.NewEvictingStringQueue(10)}, nil
}
