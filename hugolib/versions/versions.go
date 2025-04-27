// Copyright 2025 The Hugo Authors. All rights reserved.
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

package versions

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
)

type VersionConfig struct {
	// Whether this version is the default version.
	// This will be by default rendered in the root.
	// There can only be one default version.
	Default bool

	// The weight of the version.
	// Used to determine the order of the versions.
	// If zero, we use the version name to sort.
	// TODO1 sort by semantic version.
	Weight int
}

type Version struct {
	Name string
	VersionConfig
}

type Versions struct {
	versionConfigs map[string]VersionConfig
	Sorted         []Version
}

func (r Versions) IndexDefault() int {
	for i, version := range r.Sorted {
		if version.Default {
			return i
		}
	}
	panic("no default version found")
}

// IndexMatch returns the index of the first version that matches the given Glob pattern.
func (r Versions) IndexMatch(pattern string) (int, error) {
	g, err := glob.GetGlob(pattern)
	if err != nil {
		return 0, err
	}
	for i, version := range r.Sorted {
		if g.Match(version.Name) {
			return i, nil
		}
	}
	return -1, nil
}

func (r *Versions) init() error {
	if len(r.versionConfigs) == 0 {
		// Add a default version.
		r.versionConfigs[""] = VersionConfig{Default: true}
	}

	var defaultSeen int
	for k, v := range r.versionConfigs {
		if k == "" {
			return errors.New("version name cannot be empty")
		}
		if err := paths.ValidateIdentifier(k); err != nil {
			return fmt.Errorf("version name %q is invalid: %s", k, err)
		}

		if v.Default {
			defaultSeen++
		}

		if defaultSeen > 1 {
			return errors.New("only one version can be the default version")
		}
		r.Sorted = append(r.Sorted, Version{Name: k, VersionConfig: v})
	}

	// Sort by weight if set, then by name.
	sort.SliceStable(r.Sorted, func(i, j int) bool {
		ri, rj := r.Sorted[i], r.Sorted[j]
		if ri.Weight == rj.Weight {
			v1, v2 := hugo.VersionString(ri.Name), hugo.VersionString(rj.Name)
			return v1.Compare(v2) < 0
		}
		if rj.Weight == 0 {
			return true
		}
		if ri.Weight == 0 {
			return false
		}
		return ri.Weight < rj.Weight
	})

	if defaultSeen == 0 {
		// If no default version is set, we set the first one.
		first := r.Sorted[0]
		first.Default = true
		r.versionConfigs[first.Name] = first.VersionConfig
		r.Sorted[0] = first
	}

	return nil
}

func (r Versions) Has(version string) bool {
	_, found := r.versionConfigs[version]
	return found
}

func DecodeConfig(m map[string]any) (*config.ConfigNamespace[map[string]VersionConfig, Versions], error) {
	return config.DecodeNamespace[map[string]VersionConfig](m, func(in any) (Versions, any, error) {
		var versions Versions
		var conf map[string]VersionConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return versions, nil, err
		}
		versions.versionConfigs = conf
		if err := versions.init(); err != nil {
			return versions, nil, err
		}
		return versions, nil, nil
	})
}
