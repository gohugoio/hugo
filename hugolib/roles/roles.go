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

package roles

import (
	"errors"
	"fmt"
	"iter"
	"sort"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/mitchellh/mapstructure"
)

var _ sitesmatrix.DimensionInfo = (*roleWrapper)(nil)

type RoleConfig struct {
	// The weight of the role.
	// Used to determine the order of the roles.
	// If zero, we use the role name.
	Weight int
}

type Role interface {
	Name() string
}

type Roles []Role

var _ Role = (*roleWrapper)(nil)

func NewRole(r RoleInternal) Role {
	return roleWrapper{r: r}
}

type roleWrapper struct {
	r RoleInternal
}

func (r roleWrapper) Name() string {
	return r.r.Name
}

func (r roleWrapper) IsDefault() bool {
	return r.r.Default
}

type RoleInternal struct {
	// Name is the name of the role, extracted from the key in the config.
	Name string

	// Whether this role is the default role.
	// This will be rendered in the root.
	// There is only be one default role.
	Default bool

	RoleConfig
}

type RolesInternal struct {
	roleConfigs map[string]RoleConfig
	Sorted      []RoleInternal
}

func (r RolesInternal) Len() int {
	return len(r.Sorted)
}

func (r RolesInternal) IndexDefault() int {
	for i, role := range r.Sorted {
		if role.Default {
			return i
		}
	}
	panic("no default role found")
}

func (r RolesInternal) ResolveName(i int) string {
	if i < 0 || i >= len(r.Sorted) {
		panic(fmt.Sprintf("index %d out of range for roles", i))
	}
	return r.Sorted[i].Name
}

func (r RolesInternal) ResolveIndex(name string) int {
	for i, role := range r.Sorted {
		if role.Name == name {
			return i
		}
	}
	panic(fmt.Sprintf("no role found for name %q", name))
}

// IndexMatch returns an iterator for the roles that match the filter.
func (r RolesInternal) IndexMatch(match predicate.P[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, role := range r.Sorted {
			if match(role.Name) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

// ForEachIndex returns an iterator for the indices of the roles.
func (r RolesInternal) ForEachIndex() iter.Seq[int] {
	return func(yield func(i int) bool) {
		for i := range r.Sorted {
			if !yield(i) {
				return
			}
		}
	}
}

const defaultContentRoleFallback = "guest"

func (r *RolesInternal) init(defaultContentRole string) (string, error) {
	if r.roleConfigs == nil {
		r.roleConfigs = make(map[string]RoleConfig)
	}
	defaultContentRoleProvided := defaultContentRole != ""
	if len(r.roleConfigs) == 0 {
		// Add a default role.
		if defaultContentRole == "" {
			defaultContentRole = defaultContentRoleFallback
		}
		r.roleConfigs[defaultContentRole] = RoleConfig{}
	}

	var defaultSeen bool
	for k, v := range r.roleConfigs {
		if k == "" {
			return "", errors.New("role name cannot be empty")
		}

		if err := paths.ValidateIdentifier(k); err != nil {
			return "", fmt.Errorf("role name %q is invalid: %s", k, err)
		}

		var isDefault bool
		if k == defaultContentRole {
			isDefault = true
			defaultSeen = true
		}

		r.Sorted = append(r.Sorted, RoleInternal{Name: k, Default: isDefault, RoleConfig: v})
	}

	// Sort by weight if set, then by name.
	sort.SliceStable(r.Sorted, func(i, j int) bool {
		ri, rj := r.Sorted[i], r.Sorted[j]
		if ri.Weight == rj.Weight {
			return ri.Name < rj.Name
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
		if defaultContentRoleProvided {
			return "", fmt.Errorf("the configured defaultContentRole %q does not exist", defaultContentRole)
		}
		// If no default role is set, we set the first one.
		first := r.Sorted[0]
		first.Default = true
		r.roleConfigs[first.Name] = first.RoleConfig
		r.Sorted[0] = first
		defaultContentRole = first.Name
	}

	return defaultContentRole, nil
}

func (r RolesInternal) Has(role string) bool {
	_, found := r.roleConfigs[role]
	return found
}

func DecodeConfig(defaultContentRole string, m map[string]any) (*config.ConfigNamespace[map[string]RoleConfig, RolesInternal], string, error) {
	v, err := config.DecodeNamespace[map[string]RoleConfig](m, func(in any) (RolesInternal, any, error) {
		var roles RolesInternal
		var conf map[string]RoleConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return roles, nil, err
		}
		roles.roleConfigs = conf
		var err error
		if defaultContentRole, err = roles.init(defaultContentRole); err != nil {
			return roles, nil, err
		}
		return roles, roles.roleConfigs, nil
	})

	return v, defaultContentRole, err
}
