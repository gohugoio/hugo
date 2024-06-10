// Copyright 2019 The Hugo Authors. All rights reserved.
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

package navigation

import (
	"fmt"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/mitchellh/mapstructure"

	"github.com/spf13/cast"
)

type PageMenusProvider interface {
	PageMenusGetter
	MenuQueryProvider
}

type PageMenusGetter interface {
	Menus() PageMenus
}

type MenusGetter interface {
	Menus() Menus
}

type MenuQueryProvider interface {
	HasMenuCurrent(menuID string, me *MenuEntry) bool
	IsMenuCurrent(menuID string, inme *MenuEntry) bool
}

func PageMenusFromPage(ms any, p Page) (PageMenus, error) {
	if ms == nil {
		return nil, nil
	}
	pm := PageMenus{}
	me := MenuEntry{}
	SetPageValues(&me, p)

	// Could be the name of the menu to attach it to
	mname, err := cast.ToStringE(ms)

	if err == nil {
		me.Menu = mname
		pm[mname] = &me
		return pm, nil
	}

	// Could be a slice of strings
	mnames, err := cast.ToStringSliceE(ms)

	if err == nil {
		for _, mname := range mnames {
			me.Menu = mname
			pm[mname] = &me
		}
		return pm, nil
	}

	wrapErr := func(err error) error {
		return fmt.Errorf("unable to process menus for page %q: %w", p.Path(), err)
	}

	// Could be a structured menu entry
	menus, err := maps.ToStringMapE(ms)
	if err != nil {
		return pm, wrapErr(err)
	}

	for name, menu := range menus {
		menuEntry := MenuEntry{Menu: name}
		if menu != nil {
			ime, err := maps.ToStringMapE(menu)
			if err != nil {
				return pm, wrapErr(err)
			}
			if err := mapstructure.WeakDecode(ime, &menuEntry.MenuConfig); err != nil {
				return pm, err
			}
		}
		SetPageValues(&menuEntry, p)
		pm[name] = &menuEntry
	}

	return pm, nil
}

func NewMenuQueryProvider(
	pagem PageMenusGetter,
	sitem MenusGetter,
	p Page,
) MenuQueryProvider {
	return &pageMenus{
		p:     p,
		pagem: pagem,
		sitem: sitem,
	}
}

type pageMenus struct {
	pagem PageMenusGetter
	sitem MenusGetter
	p     Page
}

func (pm *pageMenus) HasMenuCurrent(menuID string, me *MenuEntry) bool {
	if !types.IsNil(me.Page) && me.Page.IsSection() {
		if ok := me.Page.IsAncestor(pm.p); ok {
			return true
		}
	}

	if !me.HasChildren() {
		return false
	}

	menus := pm.pagem.Menus()

	if m, ok := menus[menuID]; ok {
		for _, child := range me.Children {
			if child.isEqual(m) {
				return true
			}
			if pm.HasMenuCurrent(menuID, child) {
				return true
			}
		}
	}

	if pm.p == nil {
		return false
	}

	for _, child := range me.Children {
		if child.isSamePage(pm.p) {
			return true
		}

		if pm.HasMenuCurrent(menuID, child) {
			return true
		}
	}

	return false
}

func (pm *pageMenus) IsMenuCurrent(menuID string, inme *MenuEntry) bool {
	menus := pm.pagem.Menus()

	if me, ok := menus[menuID]; ok {
		if me.isEqual(inme) {
			return true
		}
	}

	if pm.p == nil {
		return false
	}

	if !inme.isSamePage(pm.p) {
		return false
	}

	// This resource may be included in several menus.
	// Search for it to make sure that it is in the menu with the given menuId.
	if menu, ok := pm.sitem.Menus()[menuID]; ok {
		for _, menuEntry := range menu {
			if menuEntry.isSameResource(inme) {
				return true
			}

			descendantFound := pm.isSameAsDescendantMenu(inme, menuEntry)
			if descendantFound {
				return descendantFound
			}

		}
	}

	return false
}

func (pm *pageMenus) isSameAsDescendantMenu(inme *MenuEntry, parent *MenuEntry) bool {
	if parent.HasChildren() {
		for _, child := range parent.Children {
			if child.isSameResource(inme) {
				return true
			}
			descendantFound := pm.isSameAsDescendantMenu(inme, child)
			if descendantFound {
				return descendantFound
			}
		}
	}
	return false
}

var NopPageMenus = new(nopPageMenus)

type nopPageMenus int

func (m nopPageMenus) Menus() PageMenus {
	return PageMenus{}
}

func (m nopPageMenus) HasMenuCurrent(menuID string, me *MenuEntry) bool {
	return false
}

func (m nopPageMenus) IsMenuCurrent(menuID string, inme *MenuEntry) bool {
	return false
}
