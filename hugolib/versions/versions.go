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
	"iter"
	"sort"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/version"
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

type VersionConfig struct {
	// The weight of the version.
	// Used to determine the order of the versions.
	// If zero, we use the version name to sort.
	// TODO1 sort by semantic version.
	Weight int
}

// Site is a sub set if page.Site to avoid circular dependencies.
type Site interface {
	Title() string
}
type Version interface {
	Name() string
	Site() Site
}

type Versions []Version

type VersionSite struct {
	v VersionInternal
	s Site
}

func (v VersionSite) Name() string {
	return v.v.Name
}

func (v VersionSite) Site() Site {
	return v.s
}

func NewVersionSite(v VersionInternal, s Site) Version {
	return VersionSite{v: v, s: s}
}

var _ Version = (*VersionSite)(nil)

type VersionInternal struct {
	// Name of the version.
	// This is the key from the config.
	Name string
	// Whether this version is the default version.
	// This will be by default rendered in the root.
	// There can only be one default version.
	Default bool

	VersionConfig
}

type VersionsInternal struct {
	versionConfigs map[string]VersionConfig

	Sorted []VersionInternal
}

func (r VersionsInternal) IndexDefault() int {
	for i, version := range r.Sorted {
		if version.Default {
			return i
		}
	}
	panic("no default version found")
}

// IndexMatch returns an iterator for the versions that match the filter.
func (r VersionsInternal) IndexMatch(filter predicate.Filter[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, version := range r.Sorted {
			if !filter.ShouldExcludeFine(version.Name) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

const dfaultContentVersionFallback = "v1"

func (r *VersionsInternal) init(defaultContentVersion string) error {
	if len(r.versionConfigs) == 0 {
		if defaultContentVersion == "" {
			defaultContentVersion = dfaultContentVersionFallback
		}
		// Add a default version.
		r.versionConfigs[defaultContentVersion] = VersionConfig{}
	}

	var defaultSeen bool
	for k, v := range r.versionConfigs {
		if k == "" {
			return errors.New("version name cannot be empty")
		}
		if err := paths.ValidateIdentifier(k); err != nil {
			return fmt.Errorf("version name %q is invalid: %s", k, err)
		}

		var isDefault bool
		if k == defaultContentVersion {
			isDefault = true
			defaultSeen = true
		}

		r.Sorted = append(r.Sorted, VersionInternal{Name: k, Default: isDefault, VersionConfig: v})
	}

	// Sort by weight if set, then semver descending.
	sort.SliceStable(r.Sorted, func(i, j int) bool {
		ri, rj := r.Sorted[i], r.Sorted[j]
		if ri.Weight == rj.Weight {
			v1, v2 := version.MustParseVersion(ri.Name), version.MustParseVersion(rj.Name)
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

	if !defaultSeen {
		// If no default version is set, we set the first one.
		first := r.Sorted[0]
		first.Default = true
		r.versionConfigs[first.Name] = first.VersionConfig
		r.Sorted[0] = first
	}

	return nil
}

func (r VersionsInternal) Has(version string) bool {
	_, found := r.versionConfigs[version]
	return found
}

func DecodeConfig(defaultContentVersion string, m map[string]any) (*config.ConfigNamespace[map[string]VersionConfig, VersionsInternal], error) {
	return config.DecodeNamespace[map[string]VersionConfig](m, func(in any) (VersionsInternal, any, error) {
		var versions VersionsInternal
		var conf map[string]VersionConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return versions, nil, err
		}
		versions.versionConfigs = conf
		if err := versions.init(defaultContentVersion); err != nil {
			return versions, nil, err
		}
		return versions, nil, nil
	})
}
