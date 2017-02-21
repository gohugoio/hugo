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

// Minimal implementation of recursive searching for
// a menu entry inside a hierarchy of nested menus.
// Useful for finding a parent by identifier,
// since menu entries have no "uplink" reference.
func (m Menu) searchRecursive(argIdentifier string) *MenuEntry {
	for _, mel := range m {
		if mel.Identifier == argIdentifier {
			return mel
		}
		if mel.HasChildren() {
			if found := mel.Children.searchRecursive(argIdentifier); found != nil {
				return found
			}
		}
	}
	return nil
}

// Returns precursor and successor for a given slice of
// MenuEntries with a particular MenuEntry "me" as starting point.
// Prev/next are returned from the *same* level only.
// Otherwise, nil is returned.
func (m Menu) sameLevelPrevNext(me *MenuEntry) (*MenuEntry, *MenuEntry) {
	iPrev, iNext := -1, -1
	for i, mel := range m {
		if mel.Identifier == me.Identifier {
			iPrev = i - 1
			iNext = i + 1
			break
		}
	}
	var prev, next *MenuEntry
	if iPrev >= 0 && iPrev < len(m) {
		prev = m[iPrev]
	}
	if iNext >= 0 && iNext < len(m) {
		next = m[iNext]
	}
	return prev, next
}

// When searching for prev, we have to look into deeper levels.
// Taking same level prev, and searching for its last child,
// returning it, unless it has itself a latest child... recurse.
//
// We could also draw a closure over a "global previous" menu entry,
// but this would cost O(n), whereas our lookup function costs O(somelog(n)).
func (me *MenuEntry) closestRecentDeeper() *MenuEntry {
	if me == nil {
		return nil
	}
	if me.HasChildren() {
		maxIdx := len(me.Children) - 1
		crd := (me.Children[maxIdx]).closestRecentDeeper()
		if crd != nil {
			return crd
		}
		return me.Children[maxIdx]
	}
	return nil
}

// Pulling it all together - we traverse a menu tree, searching for a
// specific entry "me". When we find "me", we look for its
// predecessor and successor in *multilevel* terms, meaning
// first looking for prev/next on deeper levels, then on same level,
// finally on upper level.
// The resulting prev/next menu menu entries are suitable for
// traversing a menu hierarchy exactly like a word-processing software would:
// 1.
// 1.1
// 1.2
// 2.1
// 2.1.1
// 2.2
func (m Menu) hasMenuCurrentPrevNext(me *MenuEntry) (bool, *MenuEntry, *MenuEntry) {
	for _, mel := range m {
		if mel.IsEqual(me) {
			prev, next := m.sameLevelPrevNext(me)
			if mel.HasChildren() {
				next = mel.Children[0] // deeper level next overrides same level next
			}
			if crd := prev.closestRecentDeeper(); crd != nil {
				prev = crd // deeper level prev overrides same level prev
			}
			return true, prev, next
		}
		if mel.HasChildren() {
			found, prev, next := mel.Children.hasMenuCurrentPrevNext(me)
			if found {
				if prev == nil {
					prev = mel // deeper levels yielded no prev => upper preceding node becomes prev
				}
				if next == nil {
					_, nextUp := m.sameLevelPrevNext(mel) // looped menu entry
					next = nextUp                         // deeper levels yielded no next => upper node-next becomes next
				}
				return true, prev, next
			}
		}
	}
	return false, nil, nil
}

// Exported methods for template usage
//     {{ $mePrev :=  $myMenu.Prev $pageMenuEntry }}
func (m Menu) Prev(me MenuEntry) *MenuEntry {
	ok, prev, _ := m.hasMenuCurrentPrevNext(&me)
	if ok {
		return prev
	}
	return nil
}
func (m Menu) Next(me MenuEntry) *MenuEntry {
	ok, _, next := m.hasMenuCurrentPrevNext(&me)
	if ok {
		return next
	}
	return nil
}

// Apart from traversing the complete menu tree,
// some websites might want a "prev chapter" "level up" "next chapter"
// navigation. Level down is already provided by MenuEntry.Children.
// But "Up" is almost impossible to code with template functions.
func (m Menu) Up(me MenuEntry) *MenuEntry {
	if me.Parent == "" {
		return nil
	}
	return m.searchRecursive(me.Parent)
}

// Next menu entry is also template-accessible
// by ranging over the menu entries.
// But following func is much more concise for template usage.
func (m Menu) PrevSameLevel(me MenuEntry) *MenuEntry {
	if me.Parent == "" {
		prev, _ := m.sameLevelPrevNext(&me)
		return prev
	}
	mePar := m.searchRecursive(me.Parent)
	if mePar == nil || mePar.Children == nil {
		return nil
	}
	prev, _ := mePar.Children.sameLevelPrevNext(&me)
	return prev
}

// See PrevSameLevel
func (m Menu) NextSameLevel(me MenuEntry) *MenuEntry {
	if me.Parent == "" {
		_, next := m.sameLevelPrevNext(&me)
		return next
	}
	mePar := m.searchRecursive(me.Parent)
	if mePar == nil || mePar.Children == nil {
		return nil
	}
	_, next := mePar.Children.sameLevelPrevNext(&me)
	return next
}
