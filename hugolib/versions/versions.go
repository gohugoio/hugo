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
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/mitchellh/mapstructure"
)

var _ sitesmatrix.DimensionInfo = (*versionWrapper)(nil)

type VersionConfig struct {
	// The weight of the version.
	// Used to determine the order of the versions.
	// If zero, we use the version name to sort.
	Weight int
}

type Version interface {
	Name() string
}

type Versions []Version

type versionWrapper struct {
	v VersionInternal
}

func (v versionWrapper) Name() string {
	return v.v.Name
}

func (v versionWrapper) IsDefault() bool {
	return v.v.Default
}

func NewVersion(v VersionInternal) Version {
	return versionWrapper{v: v}
}

var _ Version = (*versionWrapper)(nil)

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

func (r VersionsInternal) Len() int {
	return len(r.Sorted)
}

func (r VersionsInternal) IndexDefault() int {
	for i, version := range r.Sorted {
		if version.Default {
			return i
		}
	}
	panic("no default version found")
}

func (r VersionsInternal) ResolveName(i int) string {
	if i < 0 || i >= len(r.Sorted) {
		panic(fmt.Sprintf("index %d out of range for versions", i))
	}
	return r.Sorted[i].Name
}

func (r VersionsInternal) ResolveIndex(name string) int {
	for i, version := range r.Sorted {
		if version.Name == name {
			return i
		}
	}
	panic(fmt.Sprintf("no version found for name %q", name))
}

// IndexMatch returns an iterator for the versions that match the filter.
func (r VersionsInternal) IndexMatch(match predicate.P[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, version := range r.Sorted {
			if match(version.Name) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

// ForEachIndex returns an iterator for the indices of the versions.
func (r VersionsInternal) ForEachIndex() iter.Seq[int] {
	return func(yield func(i int) bool) {
		for i := range r.Sorted {
			if !yield(i) {
				return
			}
		}
	}
}

const defaultContentVersionFallback = "v1.0.0"

func (r *VersionsInternal) init(defaultContentVersion string) error {
	if r.versionConfigs == nil {
		r.versionConfigs = make(map[string]VersionConfig)
	}
	defaultContentVersionProvided := defaultContentVersion != ""
	if len(r.versionConfigs) == 0 {
		if defaultContentVersion == "" {
			defaultContentVersion = defaultContentVersionFallback
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
		if defaultContentVersionProvided {
			return fmt.Errorf("the configured defaultContentVersion %q does not exist", defaultContentVersion)
		}
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
		return versions, versions.versionConfigs, nil
	})
}
