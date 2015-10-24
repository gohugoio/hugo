// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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

type Menu []*MenuEntry
type Menus map[string]*Menu
type PageMenus map[string]*MenuEntry

func (me *MenuEntry) AddChild(child *MenuEntry) {
	me.Children = append(me.Children, child)
	me.Children.Sort()
}

func (me *MenuEntry) HasChildren() bool {
	return me.Children != nil
}

func (me *MenuEntry) KeyName() string {
	if me.Identifier != "" {
		return me.Identifier
	}
	return me.Name
}

func (me *MenuEntry) hopefullyUniqueID() string {
	if me.Identifier != "" {
		return me.Identifier
	} else if me.URL != "" {
		return me.URL
	} else {
		return me.Name
	}
}

func (me *MenuEntry) IsEqual(inme *MenuEntry) bool {
	return me.hopefullyUniqueID() == inme.hopefullyUniqueID() && me.Parent == inme.Parent
}

func (me *MenuEntry) IsSameResource(inme *MenuEntry) bool {
	return me.URL != "" && inme.URL != "" && me.URL == inme.URL
}

func (me *MenuEntry) MarshallMap(ime map[string]interface{}) {
	for k, v := range ime {
		loki := strings.ToLower(k)
		switch loki {
		case "url":
			me.URL = cast.ToString(v)
		case "weight":
			me.Weight = cast.ToInt(v)
		case "name":
			me.Name = cast.ToString(v)
		case "pre":
			me.Pre = template.HTML(cast.ToString(v))
		case "post":
			me.Post = template.HTML(cast.ToString(v))
		case "identifier":
			me.Identifier = cast.ToString(v)
		case "parent":
			me.Parent = cast.ToString(v)
		}
	}
}

func (m Menu) Add(me *MenuEntry) Menu {
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
type MenuSorter struct {
	menu Menu
	by   MenuEntryBy
}

// Closure used in the Sort.Less method.
type MenuEntryBy func(m1, m2 *MenuEntry) bool

func (by MenuEntryBy) Sort(menu Menu) {
	ms := &MenuSorter{
		menu: menu,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ms)
}

var defaultMenuEntrySort = func(m1, m2 *MenuEntry) bool {
	if m1.Weight == m2.Weight {
		return m1.Name < m2.Name
	}
	return m1.Weight < m2.Weight
}

func (ms *MenuSorter) Len() int      { return len(ms.menu) }
func (ms *MenuSorter) Swap(i, j int) { ms.menu[i], ms.menu[j] = ms.menu[j], ms.menu[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ms *MenuSorter) Less(i, j int) bool { return ms.by(ms.menu[i], ms.menu[j]) }

func (m Menu) Sort() {
	MenuEntryBy(defaultMenuEntrySort).Sort(m)
}

func (m Menu) Limit(n int) Menu {
	if len(m) < n {
		return m[0:n]
	}
	return m
}

func (m Menu) ByWeight() Menu {
	MenuEntryBy(defaultMenuEntrySort).Sort(m)
	return m
}

func (m Menu) ByName() Menu {
	title := func(m1, m2 *MenuEntry) bool {
		return m1.Name < m2.Name
	}

	MenuEntryBy(title).Sort(m)
	return m
}

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
//
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
	prev, _ := m.sameLevelPrevNext(&me)
	return prev
}

// See PrevSameLevel
func (m Menu) NextSameLevel(me MenuEntry) *MenuEntry {
	_, next := m.sameLevelPrevNext(&me)
	return next
}
