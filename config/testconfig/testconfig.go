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

// This package should only be used for testing.
package testconfig

import (
	_ "unsafe"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
)

func GetTestConfigs(fs afero.Fs, cfg config.Provider) *allconfig.Configs {
	if fs == nil {
		fs = afero.NewMemMapFs()
	}
	if cfg == nil {
		cfg = config.New()
	}
	// Make sure that the workingDir exists.
	workingDir := cfg.GetString("workingDir")
	if workingDir != "" {
		if err := fs.MkdirAll(workingDir, 0o777); err != nil {
			panic(err)
		}
	}

	configs, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: fs, Flags: cfg, Environ: []string{"EMPTY_TEST_ENVIRONMENT"}})
	if err != nil {
		panic(err)
	}
	return configs
}

func GetTestConfig(fs afero.Fs, cfg config.Provider) config.AllProvider {
	return GetTestConfigs(fs, cfg).GetFirstLanguageConfig()
}

func GetTestDeps(fs afero.Fs, cfg config.Provider, beforeInit ...func(*deps.Deps)) *deps.Deps {
	if fs == nil {
		fs = afero.NewMemMapFs()
	}
	conf := GetTestConfig(fs, cfg)
	d := &deps.Deps{
		Conf: conf,
		Fs:   hugofs.NewFrom(fs, conf.BaseConfig()),
	}
	for _, f := range beforeInit {
		f(d)
	}
	if err := d.Init(); err != nil {
		panic(err)
	}
	return d
}

func GetTestConfigSectionFromStruct(section string, v any) config.AllProvider {
	data, err := toml.Marshal(v)
	if err != nil {
		panic(err)
	}
	p := maps.Params{
		section: config.FromTOMLConfigString(string(data)).Get(""),
	}
	cfg := config.NewFrom(p)
	return GetTestConfig(nil, cfg)
}
