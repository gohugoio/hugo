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

package roles

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
)

type RoleConfig struct {
	// Whether this role is the default role.
	// This will be rendered in the root.
	// There can only be one default role.
	Default bool

	// The weight of the role.
	// Used to determine the order of the roles.
	// If zero, we use the role name.
	Weight int
}

type Role struct {
	Name string
	RoleConfig
}

type Roles struct {
	roleConfigs map[string]RoleConfig
	Sorted      []Role
}

func (r Roles) IndexDefault() int {
	for i, role := range r.Sorted {
		if role.Default {
			return i
		}
	}
	panic("no default role found")
}

// IndexMatch returns the index of the first role that matches the given Glob pattern.
func (r Roles) IndexMatch(pattern string) (int, error) {
	g, err := glob.GetGlob(pattern)
	if err != nil {
		return 0, err
	}
	for i, role := range r.Sorted {
		if g.Match(role.Name) {
			return i, nil
		}
	}
	return -1, nil
}

func (r *Roles) init() error {
	if len(r.roleConfigs) == 0 {
		// Add a default role.
		r.roleConfigs["guest"] = RoleConfig{Default: true}
	}

	var defaultSeen int
	for k, v := range r.roleConfigs {
		if k == "" {
			return errors.New("role name cannot be empty")
		}

		if err := paths.ValidateIdentifier(k); err != nil {
			// TODO1 config keys gets auto lowercased, so this will (almost) never happen.
			// TODO1: Page.cloneForRole(role)
			// TODO1: Tree store: linked list for dimension nodes.
			return fmt.Errorf("role name %q is invalid: %s", k, err)
		}

		if v.Default {
			defaultSeen++
		}

		if defaultSeen > 1 {
			return errors.New("only one role can be the default role")
		}

		r.Sorted = append(r.Sorted, Role{Name: k, RoleConfig: v})
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

	if defaultSeen == 0 {
		// If no default role is set, we set the first one.
		first := r.Sorted[0]
		first.Default = true
		r.roleConfigs[first.Name] = first.RoleConfig
		r.Sorted[0] = first
	}

	return nil
}

func (r Roles) Has(role string) bool {
	_, found := r.roleConfigs[role]
	return found
}

func DecodeConfig(m map[string]any) (*config.ConfigNamespace[map[string]RoleConfig, Roles], error) {
	return config.DecodeNamespace[map[string]RoleConfig](m, func(in any) (Roles, any, error) {
		var roles Roles
		var conf map[string]RoleConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return roles, nil, err
		}
		roles.roleConfigs = conf
		if err := roles.init(); err != nil {
			return roles, nil, err
		}
		return roles, nil, nil
	})
}
