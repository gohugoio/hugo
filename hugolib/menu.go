// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"html/template"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

// MenuEntry represents a menu item defined in either Page front matter
// or in the site config.
type MenuEntry struct {
	URL        string
	Name       string
	Menu       string
	Identifier string
	Pre        template.HTML
	Post       template.HTML
	Weight     int
	Parent     string
	Children   Menu
}

// Menu is a collection of menu entries.
type Menu []*MenuEntry

// Menus is a dictionary of menus.
type Menus map[string]*Menu

// PageMenus is a dictionary of menus defined in the Pages.
type PageMenus map[string]*MenuEntry

// addChild adds a new child to this menu entry.
// The default sort order will then be applied.
func (m *MenuEntry) addChild(child *MenuEntry) {
	m.Children = append(m.Children, child)
	m.Children.Sort()
}

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
	} else if m.URL != "" {
		return m.URL
	} else {
		return m.Name
	}
}

// IsEqual returns whether the two menu entries represents the same menu entry.
func (m *MenuEntry) IsEqual(inme *MenuEntry) bool {
	return m.hopefullyUniqueID() == inme.hopefullyUniqueID() && m.Parent == inme.Parent
}

// IsSameResource returns whether the two menu entries points to the same
// resource (URL).
func (m *MenuEntry) IsSameResource(inme *MenuEntry) bool {
	return m.URL != "" && inme.URL != "" && m.URL == inme.URL
}

func (m *MenuEntry) marshallMap(ime map[string]interface{}) {
	for k, v := range ime {
		loki := strings.ToLower(k)
		switch loki {
		case "url":
			m.URL = cast.ToString(v)
		case "weight":
			m.Weight = cast.ToInt(v)
		case "name":
			m.Name = cast.ToString(v)
		case "pre":
			m.Pre = template.HTML(cast.ToString(v))
		case "post":
			m.Post = template.HTML(cast.ToString(v))
		case "identifier":
			m.Identifier = cast.ToString(v)
		case "parent":
			m.Parent = cast.ToString(v)
		}
	}
}

func (m Menu) add(me *MenuEntry) Menu {
	app := func(slice Menu, x ...*MenuEntry) Menu {
		n := len(slice) + len(x)
		if n > cap(slice) {
			size := cap(slice) * 2
			if size < n {
				size = n
			}
			new := make(Menu, size)
			copy(new, slice)
			slice = new
		}
		slice = slice[0:n]
		copy(slice[n-len(x):], x)
		return slice
	}

	m = app(m, me)
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
		if m1.Name == m2.Name {
			return m1.Identifier < m2.Identifier
		}
		return m1.Name < m2.Name
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
	menuEntryBy(defaultMenuEntrySort).Sort(m)
	return m
}

// ByName sorts the menu by the name defined in the menu configuration.
func (m Menu) ByName() Menu {
	title := func(m1, m2 *MenuEntry) bool {
		return m1.Name < m2.Name
	}

	menuEntryBy(title).Sort(m)
	return m
}

// Reverse reverses the order of the menu entries.
func (m Menu) Reverse() Menu {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}

	return m
}
