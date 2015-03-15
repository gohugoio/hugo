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
	"sync"
	"time"
)

type Node struct {
	RSSLink template.HTML
	Site    *SiteInfo
	//	layout      string
	Data        map[string]interface{}
	Title       string
	Description string
	Keywords    []string
	Params      map[string]interface{}
	Date        time.Time
	Sitemap     Sitemap
	UrlPath
	paginator     *pager
	paginatorInit sync.Once
	scratch       *Scratch
}

func (n *Node) Now() time.Time {
	return time.Now()
}

func (n *Node) HasMenuCurrent(menuId string, inme *MenuEntry) bool {
	if inme.HasChildren() {
		me := MenuEntry{Name: n.Title, Url: n.Url}

		for _, child := range inme.Children {
			if me.IsSameResource(child) {
				return true
			}
		}
	}

	return false
}

func (n *Node) IsMenuCurrent(menuId string, inme *MenuEntry) bool {

	me := MenuEntry{Name: n.Title, Url: n.Url}
	if !me.IsSameResource(inme) {
		return false
	}

	// this resource may be included in several menus
	// search for it to make sure that it is in the menu with the given menuId
	if menu, ok := (*n.Site.Menus)[menuId]; ok {
		for _, menuEntry := range *menu {
			if menuEntry.IsSameResource(inme) {
				return true
			}

			descendantFound := n.isSameAsDescendantMenu(inme, menuEntry)
			if descendantFound {
				return descendantFound
			}

		}
	}

	return false
}

func (n *Node) Hugo() *HugoInfo {
	return hugoInfo
}

func (n *Node) isSameAsDescendantMenu(inme *MenuEntry, parent *MenuEntry) bool {
	if parent.HasChildren() {
		for _, child := range parent.Children {
			if child.IsSameResource(inme) {
				return true
			}
			descendantFound := n.isSameAsDescendantMenu(inme, child)
			if descendantFound {
				return descendantFound
			}
		}
	}
	return false
}

func (n *Node) RSSlink() template.HTML {
	return n.RSSLink
}

func (n *Node) IsNode() bool {
	return true
}

func (n *Node) IsPage() bool {
	return !n.IsNode()
}

func (n *Node) Ref(ref string) (string, error) {
	return n.Site.Ref(ref, nil)
}

func (n *Node) RelRef(ref string) (string, error) {
	return n.Site.RelRef(ref, nil)
}

type UrlPath struct {
	Url       string
	Permalink template.HTML
	Slug      string
	Section   string
}

// Scratch returns the writable context associated with this Node.
func (n *Node) Scratch() *Scratch {
	if n.scratch == nil {
		n.scratch = newScratch()
	}
	return n.scratch
}
