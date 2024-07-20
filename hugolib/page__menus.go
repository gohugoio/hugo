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

package hugolib

import (
	"context"
	"sync"

	"github.com/gohugoio/hugo/navigation"
)

type pageMenus struct {
	p *pageState

	q navigation.MenuQueryProvider

	pmInit sync.Once
	pm     navigation.PageMenus
}

func (p *pageMenus) HasMenuCurrent(menuID string, me *navigation.MenuEntry) bool {
	p.p.s.init.menus.Do(context.Background())
	p.init()
	return p.q.HasMenuCurrent(menuID, me)
}

func (p *pageMenus) IsMenuCurrent(menuID string, inme *navigation.MenuEntry) bool {
	p.p.s.init.menus.Do(context.Background())
	p.init()
	return p.q.IsMenuCurrent(menuID, inme)
}

func (p *pageMenus) Menus() navigation.PageMenus {
	// There is a reverse dependency here. initMenus will, once, build the
	// site menus and update any relevant page.
	p.p.s.init.menus.Do(context.Background())

	return p.menus()
}

func (p *pageMenus) menus() navigation.PageMenus {
	p.init()
	return p.pm
}

func (p *pageMenus) init() {
	p.pmInit.Do(func() {
		p.q = navigation.NewMenuQueryProvider(
			p,
			p.p.s,
			p.p,
		)

		params := p.p.Params()

		var menus any
		var ok bool

		if p.p.m.pageConfig.Menus != nil {
			menus = p.p.m.pageConfig.Menus
		} else {
			menus, ok = params["menus"]
			if !ok {
				menus = params["menu"]
			}
		}

		var err error
		p.pm, err = navigation.PageMenusFromPage(menus, p.p)
		if err != nil {
			p.p.s.Log.Errorln(p.p.wrapError(err))
		}
	})
}
