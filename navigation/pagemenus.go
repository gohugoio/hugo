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
	"github.com/gohugoio/hugo/common/maps"

	"github.com/pkg/errors"
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

func PageMenusFromPage(p Page) (PageMenus, error) {
	params := p.Params()

	ms, ok := params["menus"]
	if !ok {
		ms, ok = params["menu"]
	}

	pm := PageMenus{}

	if !ok {
		return nil, nil
	}

	me := MenuEntry{Page: p, Name: p.LinkTitle(), Weight: p.Weight()}

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

	// Could be a structured menu entry
	menus, err := maps.ToStringMapE(ms)
	if err != nil {
		return pm, errors.Wrapf(err, "unable to process menus for %q", p.LinkTitle())
	}

	for name, menu := range menus {
		menuEntry := MenuEntry{Page: p, Name: p.LinkTitle(), Weight: p.Weight(), Menu: name}
		if menu != nil {
			ime, err := maps.ToStringMapE(menu)
			if err != nil {
				return pm, errors.Wrapf(err, "unable to process menus for %q", p.LinkTitle())
			}

			menuEntry.MarshallMap(ime)
		}
		pm[name] = &menuEntry
	}

	return pm, nil

}

func NewMenuQueryProvider(
	setionPagesMenu string,
	pagem PageMenusGetter,
	sitem MenusGetter,
	p Page) MenuQueryProvider {

	return &pageMenus{
		p:               p,
		pagem:           pagem,
		sitem:           sitem,
		setionPagesMenu: setionPagesMenu,
	}
}

type pageMenus struct {
	pagem           PageMenusGetter
	sitem           MenusGetter
	setionPagesMenu string
	p               Page
}

func (pm *pageMenus) HasMenuCurrent(menuID string, me *MenuEntry) bool {

	// page is labeled as "shadow-member" of the menu with the same identifier as the section
	if pm.setionPagesMenu != "" {
		section := pm.p.Section()

		if section != "" && pm.setionPagesMenu == menuID && section == me.Identifier {
			return true
		}
	}

	if !me.HasChildren() {
		return false
	}

	menus := pm.pagem.Menus()

	if m, ok := menus[menuID]; ok {

		for _, child := range me.Children {
			if child.IsEqual(m) {
				return true
			}
			if pm.HasMenuCurrent(menuID, child) {
				return true
			}
		}
	}

	if pm.p == nil || pm.p.IsPage() {
		return false
	}

	// The following logic is kept from back when Hugo had both Page and Node types.
	// TODO(bep) consolidate / clean
	nme := MenuEntry{Page: pm.p, Name: pm.p.LinkTitle()}

	for _, child := range me.Children {
		if nme.IsSameResource(child) {
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
		if me.IsEqual(inme) {
			return true
		}
	}

	if pm.p == nil || pm.p.IsPage() {
		return false
	}

	// The following logic is kept from back when Hugo had both Page and Node types.
	// TODO(bep) consolidate / clean
	me := MenuEntry{Page: pm.p, Name: pm.p.LinkTitle()}

	if !me.IsSameResource(inme) {
		return false
	}

	// this resource may be included in several menus
	// search for it to make sure that it is in the menu with the given menuId
	if menu, ok := pm.sitem.Menus()[menuID]; ok {
		for _, menuEntry := range menu {
			if menuEntry.IsSameResource(inme) {
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
			if child.IsSameResource(inme) {
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
