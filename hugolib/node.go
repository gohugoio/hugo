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
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/spf13/hugo/helpers"
)

// TODO(bep) np add String()
type NodeType int

const (
	NodePage NodeType = iota

	// Temporary state.
	NodeUnknown

	// The rest are node types; home page, sections etc.
	NodeHome
	NodeSection
	NodeTaxonomy
	NodeTaxonomyTerms
)

func (p NodeType) String() string {
	switch p {
	case NodePage:
		return "page"
	case NodeHome:
		return "home page"
	case NodeSection:
		return "section list"
	case NodeTaxonomy:
		return "taxonomy list"
	case NodeTaxonomyTerms:
		return "taxonomy terms"
	case NodeUnknown:
		return "unknown"
	default:
		return "invalid value"
	}
}

func (p NodeType) IsNode() bool {
	return p >= NodeHome
}

type Node struct {
	NodeType NodeType

	// a natural key that should be unique for this site
	// for the home page this will typically be "home", but it can anything
	// as long as it is the same for repeated builds.
	nodeID string

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
	paginator     *Pager
	paginatorInit sync.Once
	scratch       *Scratch

	language     *helpers.Language
	languageInit sync.Once
	lang         string

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
// This method is also implemented on Page and SiteInfo.
func (n *Node) Param(key interface{}) (interface{}, error) {
	return n.Site.Param(key)
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
	return n.NodeType.IsNode()
}

func (n *Node) IsHome() bool {
	return n.NodeType == NodeHome
}

func (n *Node) IsPage() bool {
	return n.NodeType == NodePage
}

func (n *Node) Ref(ref string) (string, error) {
	return n.Site.Ref(ref, nil)
}

func (n *Node) RelRef(ref string) (string, error) {
	return n.Site.RelRef(ref, nil)
}

type URLPath struct {
	URL       string
	Permalink string
	Slug      string
	Section   string
}

func (n *Node) URL() string {
	return n.addLangPathPrefix(n.URLPath.URL)
}

func (n *Node) Permalink() string {
	return n.Site.permalink(n.URL())
}

// Scratch returns the writable context associated with this Node.
func (n *Node) Scratch() *Scratch {
	if n.scratch == nil {
		n.scratch = newScratch()
	}
	return n.scratch
}

func (n *Node) Language() *helpers.Language {
	n.initLanguage()
	return n.language
}

func (n *Node) Lang() string {
	// When set, Language can be different from lang in the case where there is a
	// content file (doc.sv.md) with language indicator, but there is no language
	// config for that language. Then the language will fall back on the site default.
	if n.Language() != nil {
		return n.Language().Lang
	}
	return n.lang
}

func (p *Page) isTranslation(candidate *Page) bool {
	if p == candidate || p.NodeType != candidate.NodeType {
		return false
	}

	if p.lang != candidate.lang || p.language != p.language {
		return false
	}

	if p.NodeType == NodePage || p.NodeType == NodeUnknown {
		panic("Node type not currently supported for this op")
	}

	// At this point, we know that this is a traditional Node (hoe page, section, taxonomy)
	// It represents the same node, but different language, if the sections is the same.
	if len(p.sections) != len(candidate.sections) {
		return false
	}

	for i := 0; i < len(p.sections); i++ {
		if p.sections[i] != candidate.sections[i] {
			return false
		}
	}

	return true

}

func (n *Node) shouldAddLanguagePrefix() bool {
	if !n.Site.IsMultiLingual() {
		return false
	}

	if n.Lang() == "" {
		return false
	}

	if !n.Site.defaultContentLanguageInSubdir && n.Lang() == n.Site.multilingual.DefaultLang.Lang {
		return false
	}

	return true
}

func (n *Node) initLanguage() {
	n.languageInit.Do(func() {
		if n.language != nil {
			return
		}
		pageLang := n.lang
		ml := n.Site.multilingual
		if ml == nil {
			panic("Multilanguage not set")
		}
		if pageLang == "" {
			n.language = ml.DefaultLang
			return
		}

		language := ml.Language(pageLang)

		if language == nil {
			// It can be a file named stefano.chiodino.md.
			jww.WARN.Printf("Page language (if it is that) not found in multilang setup: %s.", pageLang)
			language = ml.DefaultLang
		}

		n.language = language
	})
}

func (n *Node) LanguagePrefix() string {
	return n.Site.LanguagePrefix
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

// IsTranslated returns whether this node is translated to
// other language(s).
func (n *Node) IsTranslated() bool {
	n.initTranslations()
	return len(n.translations) > 1
}

func (n *Node) initTranslations() {
	n.translationsInit.Do(func() {
		n.translations = n.Site.owner.getNodes(n.nodeID)
	})
}

func (n *Node) addLangPathPrefix(outfile string) string {
	return n.addLangPathPrefixIfFlagSet(outfile, n.shouldAddLanguagePrefix())
}

func (n *Node) addLangPathPrefixIfFlagSet(outfile string, should bool) string {
	if helpers.IsAbsURL(outfile) {
		return outfile
	}

	if !should {
		return outfile
	}

	hadSlashSuffix := strings.HasSuffix(outfile, "/")

	outfile = "/" + path.Join(n.Lang(), outfile)
	if hadSlashSuffix {
		outfile += "/"
	}
	return outfile
}

func (n *Node) addLangFilepathPrefix(outfile string) string {
	if outfile == "" {
		outfile = helpers.FilePathSeparator
	}
	if !n.shouldAddLanguagePrefix() {
		return outfile
	}
	return helpers.FilePathSeparator + filepath.Join(n.Lang(), outfile)
}

func sectionsFromFilename(filename string) []string {
	dir, _ := filepath.Split(filename)
	dir = strings.TrimSuffix(dir, helpers.FilePathSeparator)
	sections := strings.Split(dir, helpers.FilePathSeparator)
	return sections
}

// TODO(bep) np node identificator
func nodeTypeFromFilename(filename string) NodeType {
	if !strings.Contains(filename, "_index") {
		return NodePage
	}

	if strings.HasPrefix(filename, "_index") {
		return NodeHome
	}

	// We don't know enough yet to determine the type.
	return NodeUnknown
}

func (p *Page) setNodeTypeVars(s *Site) {
	// TODO(bep) np taxonomies etc.
	if p.NodeType == NodeUnknown {
		// This is either a taxonomy list, taxonomy term or a section
		nodeType := s.nodeTypeFromSections(p.sections)

		if nodeType == NodeUnknown {
			panic(fmt.Sprintf("Unable to determine node type from %q", p.sections))
		}

		p.NodeType = nodeType
	}
	// TODO(bep) np node URL
	// Set Node URL
	switch p.NodeType {
	case NodeHome:
		p.URLPath.URL = ""
	case NodeSection:
		p.URLPath.URL = p.sections[0]
	case NodeTaxonomy:
		p.URLPath.URL = path.Join(p.sections...)
	case NodeTaxonomyTerms:
		p.URLPath.URL = path.Join(p.sections...)
	}

	p.site = s

}
