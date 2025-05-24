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
	"github.com/mitchellh/mapstructure"
)

type RoleConfig struct {
	// The weight of the role.
	// Used to determine the order of the roles.
	// If zero, we use the role name.
	Weight int
}

// Site is a sub set if page.Site to avoid circular dependencies.
type Site interface {
	Title() string
}

type Role interface {
	Name() string
	Site() Site
}

type Roles []Role

var _ Role = (*RoleSite)(nil)

func NewRoleSite(r RoleInternal, s Site) Role {
	return RoleSite{r: r, s: s}
}

type RoleSite struct {
	r RoleInternal
	s Site
}

func (r RoleSite) Name() string {
	return r.r.Name
}

func (r RoleSite) Site() Site {
	return r.s
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

func (r RolesInternal) IndexDefault() int {
	for i, role := range r.Sorted {
		if role.Default {
			return i
		}
	}
	panic("no default role found")
}

// IndexMatch returns an iterator for the roles that match the filter.
func (r RolesInternal) IndexMatch(filter predicate.Filter[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, role := range r.Sorted {
			if !filter.ShouldExcludeFine(role.Name) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

func (r *RolesInternal) init(defaultContentRole string) error {
	if len(r.roleConfigs) == 0 {
		// Add a default role.
		r.roleConfigs[""] = RoleConfig{}
	}

	var defaultSeen bool
	for k, v := range r.roleConfigs {
		if k == "" {
			return errors.New("role name cannot be empty")
		}

		if err := paths.ValidateIdentifier(k); err != nil {
			// TODO1 config keys gets auto lowercased, so this will (almost) never happen.
			// TODO1: Tree store: linked list for dimension nodes.
			return fmt.Errorf("role name %q is invalid: %s", k, err)
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
		// If no default role is set, we set the first one.
		first := r.Sorted[0]
		first.Default = true
		r.roleConfigs[first.Name] = first.RoleConfig
		r.Sorted[0] = first
	}

	return nil
}

func (r RolesInternal) Has(role string) bool {
	_, found := r.roleConfigs[role]
	return found
}

func DecodeConfig(defaultContentRole string, m map[string]any) (*config.ConfigNamespace[map[string]RoleConfig, RolesInternal], error) {
	return config.DecodeNamespace[map[string]RoleConfig](m, func(in any) (RolesInternal, any, error) {
		var roles RolesInternal
		var conf map[string]RoleConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return roles, nil, err
		}
		roles.roleConfigs = conf
		if err := roles.init(defaultContentRole); err != nil {
			return roles, nil, err
		}
		return roles, nil, nil
	})
}
