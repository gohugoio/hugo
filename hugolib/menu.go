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
	"github.com/spf13/hugo/helpers"
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

// Url is deprecated. Will be removed in 0.15.
func (me *MenuEntry) Url() string {
	helpers.Deprecated("MenuEntry", ".Url", ".URL")
	return me.URL
}

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
