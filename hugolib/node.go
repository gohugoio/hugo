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
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/hugo/helpers"

	"github.com/spf13/cast"
)

type Node struct {
	RSSLink template.HTML
	Site    *SiteInfo `json:"-"`
	//	layout      string
	Data        map[string]interface{}
	Title       string
	Description string
	Keywords    []string
	Params      map[string]interface{}
	Date        time.Time
	Lastmod     time.Time
	Sitemap     Sitemap
	URLPath
	IsHome        bool
	paginator     *Pager
	paginatorInit sync.Once
	scratch       *Scratch

	language *Language
	lang     string // TODO(bep) multilingo

	translations     Nodes
	translationsInit sync.Once
}

// The Nodes type is temporary until we get https://github.com/spf13/hugo/issues/2297 fixed.
type Nodes []*Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].language.Weight < n[j].language.Weight
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n *Node) Now() time.Time {
	return time.Now()
}

func (n *Node) HasMenuCurrent(menuID string, inme *MenuEntry) bool {
	if inme.HasChildren() {
		me := MenuEntry{Name: n.Title, URL: n.URL()}

		for _, child := range inme.Children {
			if me.IsSameResource(child) {
				return true
			}
			if n.HasMenuCurrent(menuID, child) {
				return true
			}
		}
	}

	return false
}

func (n *Node) IsMenuCurrent(menuID string, inme *MenuEntry) bool {

	me := MenuEntry{Name: n.Title, URL: n.Site.createNodeMenuEntryURL(n.URL())}

	if !me.IsSameResource(inme) {
		return false
	}

	// this resource may be included in several menus
	// search for it to make sure that it is in the menu with the given menuId
	if menu, ok := (*n.Site.Menus)[menuID]; ok {
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

// Param is a convenience method to do lookups in Site's Params map.
//
// This method is also implemented on Page.
func (n *Node) Param(key interface{}) (interface{}, error) {
	keyStr, err := cast.ToStringE(key)
	if err != nil {
		return nil, err
	}
	return n.Site.Params[keyStr], err
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

// TODO(bep) multilingo some of these are now hidden. Consider unexport.
type URLPath struct {
	URL       string
	Permalink string
	Slug      string
	Section   string
}

func (n *Node) URL() string {
	return n.addMultilingualWebPrefix(n.URLPath.URL)
}

func (n *Node) Permalink() string {
	return permalink(n.URL())
}

// Scratch returns the writable context associated with this Node.
func (n *Node) Scratch() *Scratch {
	if n.scratch == nil {
		n.scratch = newScratch()
	}
	return n.scratch
}

// TODO(bep) multilingo consolidate. See Page.
func (n *Node) Language() *Language {
	return n.language
}

func (n *Node) Lang() string {
	if n.Language() != nil {
		return n.Language().Lang
	}
	return n.lang
}

// AllTranslations returns all translations, including the current Node.
// Note that this and the one below is kind of a temporary hack before #2297 is solved.
func (n *Node) AllTranslations() Nodes {
	n.initTranslations()
	return n.translations
}

// Translations returns the translations excluding the current Node.
func (n *Node) Translations() Nodes {
	n.initTranslations()
	translations := make(Nodes, 0)

	for _, t := range n.translations {

		if t != n {
			translations = append(translations, t)
		}
	}

	return translations
}

func (n *Node) initTranslations() {
	n.translationsInit.Do(func() {
		if n.translations != nil {
			return
		}
		n.translations = make(Nodes, 0)
		for _, l := range n.Site.Languages {
			if l == n.language {
				n.translations = append(n.translations, n)
				continue
			}

			translation := *n
			translation.language = l
			translation.translations = n.translations
			n.translations = append(n.translations, &translation)
		}

		sort.Sort(n.translations)
	})
}

func (n *Node) addMultilingualWebPrefix(outfile string) string {

	if helpers.IsAbsURL(outfile) {
		return outfile
	}

	hadSlashSuffix := strings.HasSuffix(outfile, "/")

	lang := n.Lang()
	if lang == "" || !n.Site.Multilingual {
		return outfile
	}
	outfile = "/" + path.Join(lang, outfile)
	if hadSlashSuffix {
		outfile += "/"
	}
	return outfile
}

func (n *Node) addMultilingualFilesystemPrefix(outfile string) string {
	lang := n.Lang()
	if lang == "" || !n.Site.Multilingual {
		return outfile
	}
	return string(filepath.Separator) + filepath.Join(lang, outfile)
}
