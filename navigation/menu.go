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

// Package navigation provides the menu functionality.
package navigation

import (
	"html/template"
	"sort"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/compare"
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"

	"github.com/spf13/cast"
)

var smc = newMenuCache()

// MenuEntry represents a menu item defined in either Page front matter
// or in the site config.
type MenuEntry struct {
	// The menu entry configuration.
	MenuConfig

	// The menu containing this menu entry.
	Menu string

	// The URL value from front matter / config.
	ConfiguredURL string

	// The Page connected to this menu entry.
	Page Page

	// Child entries.
	Children Menu
}

func (m *MenuEntry) URL() string {
	// Check page first.
	// In Hugo 0.86.0 we added `pageRef`,
	// a way to connect menu items in site config to pages.
	// This means that you now can have both a Page
	// and a configured URL.
	// Having the configured URL as a fallback if the Page isn't found
	// is obviously more useful, especially in multilingual sites.
	if !types.IsNil(m.Page) {
		return m.Page.RelPermalink()
	}

	return m.ConfiguredURL
}

// SetPageValues sets the Page and URL values for this menu entry.
func SetPageValues(m *MenuEntry, p Page) {
	m.Page = p
	if m.MenuConfig.Name == "" {
		m.MenuConfig.Name = p.LinkTitle()
	}
	if m.MenuConfig.Title == "" {
		m.MenuConfig.Title = p.Title()
	}
	if m.MenuConfig.Weight == 0 {
		m.MenuConfig.Weight = p.Weight()
	}
}

// A narrow version of page.Page.
type Page interface {
	LinkTitle() string
	Title() string
	RelPermalink() string
	Path() string
	Section() string
	Weight() int
	IsPage() bool
	IsSection() bool
	IsAncestor(other any) bool
	Params() maps.Params
}

// Menu is a collection of menu entries.
type Menu []*MenuEntry

// Menus is a dictionary of menus.
type Menus map[string]Menu

// PageMenus is a dictionary of menus defined in the Pages.
type PageMenus map[string]*MenuEntry

// HasChildren returns whether this menu item has any children.
func (m *MenuEntry) HasChildren() bool {
	return m.Children != nil
}

// KeyName returns the key used to identify this menu entry.
func (m *MenuEntry) KeyName() string {
	if m.Identifier != "" {
		return m.Identifier
	}
	return m.Name
}

func (m *MenuEntry) hopefullyUniqueID() string {
	if m.Identifier != "" {
		return m.Identifier
	} else if m.URL() != "" {
		return m.URL()
	} else {
		return m.Name
	}
}

// isEqual returns whether the two menu entries represents the same menu entry.
func (m *MenuEntry) isEqual(inme *MenuEntry) bool {
	return m.hopefullyUniqueID() == inme.hopefullyUniqueID() && m.Parent == inme.Parent
}

// isSameResource returns whether the two menu entries points to the same
// resource (URL).
func (m *MenuEntry) isSameResource(inme *MenuEntry) bool {
	if m.isSamePage(inme.Page) {
		return m.Page == inme.Page
	}
	murl, inmeurl := m.URL(), inme.URL()
	return murl != "" && inmeurl != "" && murl == inmeurl
}

func (m *MenuEntry) isSamePage(p Page) bool {
	if !types.IsNil(m.Page) && !types.IsNil(p) {
		return m.Page == p
	}
	return false
}

// MenuConfig holds the configuration for a menu.
type MenuConfig struct {
	Identifier string
	Parent     string
	Name       string
	Pre        template.HTML
	Post       template.HTML
	URL        string
	PageRef    string
	Weight     int
	Title      string
	// User defined params.
	Params maps.Params
}

// For internal use.

// This is for internal use only.
func (m Menu) Add(me *MenuEntry) Menu {
	m = append(m, me)
	// TODO(bep)
	m.Sort()
	return m
}

/*
 * Implementation of a custom sorter for Menu
 */

// A type to implement the sort interface for Menu
type menuSorter struct {
	menu Menu
	by   menuEntryBy
}

// Closure used in the Sort.Less method.
type menuEntryBy func(m1, m2 *MenuEntry) bool

func (by menuEntryBy) Sort(menu Menu) {
	ms := &menuSorter{
		menu: menu,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ms)
}

var defaultMenuEntrySort = func(m1, m2 *MenuEntry) bool {
	if m1.Weight == m2.Weight {
		c := compare.Strings(m1.Name, m2.Name)
		if c == 0 {
			return m1.Identifier < m2.Identifier
		}
		return c < 0
	}

	if m2.Weight == 0 {
		return true
	}

	if m1.Weight == 0 {
		return false
	}

	return m1.Weight < m2.Weight
}

func (ms *menuSorter) Len() int      { return len(ms.menu) }
func (ms *menuSorter) Swap(i, j int) { ms.menu[i], ms.menu[j] = ms.menu[j], ms.menu[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ms *menuSorter) Less(i, j int) bool { return ms.by(ms.menu[i], ms.menu[j]) }

// Sort sorts the menu by weight, name and then by identifier.
func (m Menu) Sort() Menu {
	menuEntryBy(defaultMenuEntrySort).Sort(m)
	return m
}

// Limit limits the returned menu to n entries.
func (m Menu) Limit(n int) Menu {
	if len(m) > n {
		return m[0:n]
	}
	return m
}

// ByWeight sorts the menu by the weight defined in the menu configuration.
func (m Menu) ByWeight() Menu {
	const key = "menuSort.ByWeight"
	menus, _ := smc.get(key, menuEntryBy(defaultMenuEntrySort).Sort, m)

	return menus
}

// ByName sorts the menu by the name defined in the menu configuration.
func (m Menu) ByName() Menu {
	const key = "menuSort.ByName"
	title := func(m1, m2 *MenuEntry) bool {
		return compare.LessStrings(m1.Name, m2.Name)
	}

	menus, _ := smc.get(key, menuEntryBy(title).Sort, m)

	return menus
}

// Reverse reverses the order of the menu entries.
func (m Menu) Reverse() Menu {
	const key = "menuSort.Reverse"
	reverseFunc := func(menu Menu) {
		for i, j := 0, len(menu)-1; i < j; i, j = i+1, j-1 {
			menu[i], menu[j] = menu[j], menu[i]
		}
	}
	menus, _ := smc.get(key, reverseFunc, m)

	return menus
}

// Clone clones the menu entries.
// This is for internal use only.
func (m Menu) Clone() Menu {
	return append(Menu(nil), m...)
}

func DecodeConfig(in any) (*config.ConfigNamespace[map[string]MenuConfig, Menus], error) {
	buildConfig := func(in any) (Menus, any, error) {
		ret := Menus{}

		if in == nil {
			return ret, map[string]any{}, nil
		}

		menus, err := maps.ToStringMapE(in)
		if err != nil {
			return ret, nil, err
		}
		menus = maps.CleanConfigStringMap(menus)

		for name, menu := range menus {
			m, err := cast.ToSliceE(menu)
			if err != nil {
				return ret, nil, err
			} else {
				for _, entry := range m {
					var menuConfig MenuConfig
					if err := mapstructure.WeakDecode(entry, &menuConfig); err != nil {
						return ret, nil, err
					}
					maps.PrepareParams(menuConfig.Params)
					menuEntry := MenuEntry{
						Menu:       name,
						MenuConfig: menuConfig,
					}
					menuEntry.ConfiguredURL = menuEntry.MenuConfig.URL

					if ret[name] == nil {
						ret[name] = Menu{}
					}
					ret[name] = ret[name].Add(&menuEntry)
				}
			}
		}

		return ret, menus, nil
	}

	return config.DecodeNamespace[map[string]MenuConfig](in, buildConfig)
}
