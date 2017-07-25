// Copyright 2016 The Hugo Authors. All rights reserved.
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

// TODO(bep) remove this file when the reworked tests in menu_test.go is done.
// NOTE: Do not add more tests to this file!

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/deps"

	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gohugoio/hugo/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	confMenu1 = `
[[menu.main]]
    name = "Go Home"
    url = "/"
	weight = 1
	pre = "<div>"
	post = "</div>"
[[menu.main]]
    name = "Blog"
    url = "/posts"
[[menu.main]]
    name = "ext"
    url = "http://gohugo.io"
	identifier = "ext"
[[menu.main]]
    name = "ext2"
    url = "http://foo.local/Zoo/foo"
	identifier = "ext2"
[[menu.grandparent]]
	name = "grandparent"
	url = "/grandparent"
	identifier = "grandparentId"
[[menu.grandparent]]
	name = "parent"
	url = "/parent"
	identifier = "parentId"
	parent = "grandparentId"
[[menu.grandparent]]
	name = "Go Home3"
    url = "/"
	identifier = "grandchildId"
	parent = "parentId"
[[menu.tax]]
	name = "Tax1"
    url = "/two/key/"
	identifier="1"
[[menu.tax]]
	name = "Tax2"
    url = "/two/key/"
	identifier="2"
[[menu.tax]]
	name = "Tax RSS"
    url = "/two/key.xml"
	identifier="xml"
[[menu.hash]]
   name = "Tax With #"
   url = "/resource#anchor"
   identifier="hash"
[[menu.unicode]]
   name = "Unicode Russian"
   identifier = "unicode-russian"
   url = "/новости-проекта"` // Russian => "news-project"
)

var menuPage1 = []byte(`+++
title = "One"
weight = 1
[menu]
	[menu.p_one]
+++
Front Matter with Menu Pages`)

var menuPage2 = []byte(`+++
title = "Two"
weight = 2
[menu]
	[menu.p_one]
	[menu.p_two]
		identifier = "Two"

+++
Front Matter with Menu Pages`)

var menuPage3 = []byte(`+++
title = "Three"
weight = 3
[menu]
	[menu.p_two]
		Name = "Three"
		Parent = "Two"
+++
Front Matter with Menu Pages`)

var menuPage4 = []byte(`+++
title = "Four"
weight = 4
[menu]
	[menu.p_two]
		Name = "Four"
		Parent = "Three"
+++
Front Matter with Menu Pages`)

var menuPageSources = []source.ByteSource{
	{Name: filepath.FromSlash("sect/doc1.md"), Content: menuPage1},
	{Name: filepath.FromSlash("sect/doc2.md"), Content: menuPage2},
	{Name: filepath.FromSlash("sect/doc3.md"), Content: menuPage3},
}

var menuPageSectionsSources = []source.ByteSource{
	{Name: filepath.FromSlash("first/doc1.md"), Content: menuPage1},
	{Name: filepath.FromSlash("first/doc2.md"), Content: menuPage2},
	{Name: filepath.FromSlash("second-section/doc3.md"), Content: menuPage3},
	{Name: filepath.FromSlash("Fish and Chips/doc4.md"), Content: menuPage4},
}

func tstCreateMenuPageWithNameTOML(title, menu, name string) []byte {
	return []byte(fmt.Sprintf(`+++
title = "%s"
weight = 1
[menu]
	[menu.%s]
		name = "%s"
+++
Front Matter with Menu with Name`, title, menu, name))
}

func tstCreateMenuPageWithIdentifierTOML(title, menu, identifier string) []byte {
	return []byte(fmt.Sprintf(`+++
title = "%s"
weight = 1
[menu]
	[menu.%s]
		identifier = "%s"
		name = "somename"
+++
Front Matter with Menu with Identifier`, title, menu, identifier))
}

func tstCreateMenuPageWithNameYAML(title, menu, name string) []byte {
	return []byte(fmt.Sprintf(`---
title: "%s"
weight: 1
menu:
    %s:
      name: "%s"
---
Front Matter with Menu with Name`, title, menu, name))
}

func tstCreateMenuPageWithIdentifierYAML(title, menu, identifier string) []byte {
	return []byte(fmt.Sprintf(`---
title: "%s"
weight: 1
menu:
    %s:
      identifier: "%s"
      name: "somename"
---
Front Matter with Menu with Identifier`, title, menu, identifier))
}

// Issue 817 - identifier should trump everything
func TestPageMenuWithIdentifier(t *testing.T) {
	t.Parallel()
	toml := []source.ByteSource{
		{Name: "sect/doc1.md", Content: tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i1")},
		{Name: "sect/doc2.md", Content: tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i2")},
		{Name: "sect/doc3.md", Content: tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i2")}, // duplicate
	}

	yaml := []source.ByteSource{
		{Name: "sect/doc1.md", Content: tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i1")},
		{Name: "sect/doc2.md", Content: tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i2")},
		{Name: "sect/doc3.md", Content: tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i2")}, // duplicate
	}

	doTestPageMenuWithIdentifier(t, toml)
	doTestPageMenuWithIdentifier(t, yaml)

}

func doTestPageMenuWithIdentifier(t *testing.T, menuPageSources []source.ByteSource) {

	s := setupMenuTests(t, menuPageSources)

	assert.Equal(t, 3, len(s.RegularPages), "Not enough pages")

	me1 := findTestMenuEntryByID(s, "m1", "i1")
	me2 := findTestMenuEntryByID(s, "m1", "i2")

	require.NotNil(t, me1)
	require.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.URL, "doc1"), me1.URL)
	assert.True(t, strings.Contains(me2.URL, "doc2") || strings.Contains(me2.URL, "doc3"), me2.URL)

}

// Issue 817 contd - name should be second identifier in
func TestPageMenuWithDuplicateName(t *testing.T) {
	t.Parallel()
	toml := []source.ByteSource{
		{Name: "sect/doc1.md", Content: tstCreateMenuPageWithNameTOML("t1", "m1", "n1")},
		{Name: "sect/doc2.md", Content: tstCreateMenuPageWithNameTOML("t1", "m1", "n2")},
		{Name: "sect/doc3.md", Content: tstCreateMenuPageWithNameTOML("t1", "m1", "n2")}, // duplicate
	}

	yaml := []source.ByteSource{
		{Name: "sect/doc1.md", Content: tstCreateMenuPageWithNameYAML("t1", "m1", "n1")},
		{Name: "sect/doc2.md", Content: tstCreateMenuPageWithNameYAML("t1", "m1", "n2")},
		{Name: "sect/doc3.md", Content: tstCreateMenuPageWithNameYAML("t1", "m1", "n2")}, // duplicate
	}

	doTestPageMenuWithDuplicateName(t, toml)
	doTestPageMenuWithDuplicateName(t, yaml)

}

func doTestPageMenuWithDuplicateName(t *testing.T, menuPageSources []source.ByteSource) {

	s := setupMenuTests(t, menuPageSources)

	assert.Equal(t, 3, len(s.RegularPages), "Not enough pages")

	me1 := findTestMenuEntryByName(s, "m1", "n1")
	me2 := findTestMenuEntryByName(s, "m1", "n2")

	require.NotNil(t, me1)
	require.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.URL, "doc1"), me1.URL)
	assert.True(t, strings.Contains(me2.URL, "doc2") || strings.Contains(me2.URL, "doc3"), me2.URL)

}

func TestPageMenu(t *testing.T) {
	t.Parallel()
	s := setupMenuTests(t, menuPageSources)

	if len(s.RegularPages) != 3 {
		t.Fatalf("Posts not created, expected 3 got %d", len(s.RegularPages))
	}

	first := s.RegularPages[0]
	second := s.RegularPages[1]
	third := s.RegularPages[2]

	pOne := findTestMenuEntryByName(s, "p_one", "One")
	pTwo := findTestMenuEntryByID(s, "p_two", "Two")

	for i, this := range []struct {
		menu           string
		page           *Page
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"p_one", first, pOne, true, false},
		{"p_one", first, pTwo, false, false},
		{"p_one", second, pTwo, false, false},
		{"p_two", second, pTwo, true, false},
		{"p_two", third, pTwo, false, true},
		{"p_one", third, pTwo, false, false},
	} {

		if i != 4 {
			continue
		}

		isMenuCurrent := this.page.IsMenuCurrent(this.menu, this.menuItem)
		hasMenuCurrent := this.page.HasMenuCurrent(this.menu, this.menuItem)

		if isMenuCurrent != this.isMenuCurrent {
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v", i, isMenuCurrent)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			t.Errorf("[%d] Wrong result for menuItem %v for HasMenuCurrent: %v", i, this.menuItem, hasMenuCurrent)
		}

	}

}

func TestMenuURL(t *testing.T) {
	t.Parallel()
	s := setupMenuTests(t, menuPageSources)

	for i, this := range []struct {
		me          *MenuEntry
		expectedURL string
	}{
		// issue #888
		{findTestMenuEntryByID(s, "hash", "hash"), "/Zoo/resource#anchor"},
		// issue #1774
		{findTestMenuEntryByID(s, "main", "ext"), "http://gohugo.io"},
		{findTestMenuEntryByID(s, "main", "ext2"), "http://foo.local/Zoo/foo"},
	} {

		if this.me == nil {
			t.Errorf("[%d] MenuEntry not found", i)
			continue
		}

		if this.me.URL != this.expectedURL {
			t.Errorf("[%d] Got URL %s expected %s", i, this.me.URL, this.expectedURL)
		}

	}

}

// Issue #1934
func TestYAMLMenuWithMultipleEntries(t *testing.T) {
	t.Parallel()
	ps1 := []byte(`---
title: "Yaml 1"
weight: 5
menu: ["p_one", "p_two"]
---
Yaml Front Matter with Menu Pages`)

	ps2 := []byte(`---
title: "Yaml 2"
weight: 5
menu:
    p_three:
    p_four:
---
Yaml Front Matter with Menu Pages`)

	s := setupMenuTests(t, []source.ByteSource{
		{Name: filepath.FromSlash("sect/yaml1.md"), Content: ps1},
		{Name: filepath.FromSlash("sect/yaml2.md"), Content: ps2}})

	p1 := s.RegularPages[0]
	assert.Len(t, p1.Menus(), 2, "List YAML")
	p2 := s.RegularPages[1]
	assert.Len(t, p2.Menus(), 2, "Map YAML")

}

// issue #719
func TestMenuWithUnicodeURLs(t *testing.T) {
	t.Parallel()
	for _, canonifyURLs := range []bool{true, false} {
		doTestMenuWithUnicodeURLs(t, canonifyURLs)
	}
}

func doTestMenuWithUnicodeURLs(t *testing.T, canonifyURLs bool) {

	s := setupMenuTests(t, menuPageSources, "canonifyURLs", canonifyURLs)

	unicodeRussian := findTestMenuEntryByID(s, "unicode", "unicode-russian")

	expected := "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0"

	if !canonifyURLs {
		expected = "/Zoo" + expected
	}

	assert.Equal(t, expected, unicodeRussian.URL)
}

// Issue #1114
func TestSectionPagesMenu2(t *testing.T) {
	t.Parallel()
	doTestSectionPagesMenu(true, t)
	doTestSectionPagesMenu(false, t)
}

func doTestSectionPagesMenu(canonifyURLs bool, t *testing.T) {

	s := setupMenuTests(t, menuPageSectionsSources,
		"sectionPagesMenu", "spm",
		"canonifyURLs", canonifyURLs,
	)

	sects := s.getPage(KindHome).Sections()

	require.Equal(t, 3, len(sects))

	firstSectionPages := s.getPage(KindSection, "first").Pages
	require.Equal(t, 2, len(firstSectionPages))
	secondSectionPages := s.getPage(KindSection, "second-section").Pages
	require.Equal(t, 1, len(secondSectionPages))
	fishySectionPages := s.getPage(KindSection, "Fish and Chips").Pages
	require.Equal(t, 1, len(fishySectionPages))

	nodeFirst := s.getPage(KindSection, "first")
	require.NotNil(t, nodeFirst)
	nodeSecond := s.getPage(KindSection, "second-section")
	require.NotNil(t, nodeSecond)
	nodeFishy := s.getPage(KindSection, "Fish and Chips")
	require.Equal(t, "Fish and Chips", nodeFishy.sections[0])

	firstSectionMenuEntry := findTestMenuEntryByID(s, "spm", "first")
	secondSectionMenuEntry := findTestMenuEntryByID(s, "spm", "second-section")
	fishySectionMenuEntry := findTestMenuEntryByID(s, "spm", "Fish and Chips")

	require.NotNil(t, firstSectionMenuEntry)
	require.NotNil(t, secondSectionMenuEntry)
	require.NotNil(t, nodeFirst)
	require.NotNil(t, nodeSecond)
	require.NotNil(t, fishySectionMenuEntry)
	require.NotNil(t, nodeFishy)

	require.True(t, nodeFirst.IsMenuCurrent("spm", firstSectionMenuEntry))
	require.False(t, nodeFirst.IsMenuCurrent("spm", secondSectionMenuEntry))
	require.False(t, nodeFirst.IsMenuCurrent("spm", fishySectionMenuEntry))
	require.True(t, nodeFishy.IsMenuCurrent("spm", fishySectionMenuEntry))
	require.Equal(t, "Fish and Chips", fishySectionMenuEntry.Name)

	for _, p := range firstSectionPages {
		require.True(t, p.HasMenuCurrent("spm", firstSectionMenuEntry))
		require.False(t, p.HasMenuCurrent("spm", secondSectionMenuEntry))
	}

	for _, p := range secondSectionPages {
		require.False(t, p.HasMenuCurrent("spm", firstSectionMenuEntry))
		require.True(t, p.HasMenuCurrent("spm", secondSectionMenuEntry))
	}

	for _, p := range fishySectionPages {
		require.False(t, p.HasMenuCurrent("spm", firstSectionMenuEntry))
		require.False(t, p.HasMenuCurrent("spm", secondSectionMenuEntry))
		require.True(t, p.HasMenuCurrent("spm", fishySectionMenuEntry))
	}
}

func TestMenuLimit(t *testing.T) {
	t.Parallel()
	s := setupMenuTests(t, menuPageSources)
	m := *s.Menus["main"]

	// main menu has 4 entries
	firstTwo := m.Limit(2)
	assert.Equal(t, 2, len(firstTwo))
	for i := 0; i < 2; i++ {
		assert.Equal(t, m[i], firstTwo[i])
	}
	assert.Equal(t, m, m.Limit(4))
	assert.Equal(t, m, m.Limit(5))
}

func TestMenuSortByN(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		sortFunc   func(p Menu) Menu
		assertFunc func(p Menu) bool
	}{
		{(Menu).Sort, func(p Menu) bool { return p[0].Weight == 1 && p[1].Name == "nx" && p[2].Identifier == "ib" }},
		{(Menu).ByWeight, func(p Menu) bool { return p[0].Weight == 1 && p[1].Name == "nx" && p[2].Identifier == "ib" }},
		{(Menu).ByName, func(p Menu) bool { return p[0].Name == "na" }},
		{(Menu).Reverse, func(p Menu) bool { return p[0].Identifier == "ib" && p[len(p)-1].Identifier == "ia" }},
	} {
		menu := Menu{&MenuEntry{Weight: 3, Name: "nb", Identifier: "ia"},
			&MenuEntry{Weight: 1, Name: "na", Identifier: "ic"},
			&MenuEntry{Weight: 1, Name: "nx", Identifier: "ic"},
			&MenuEntry{Weight: 2, Name: "nb", Identifier: "ix"},
			&MenuEntry{Weight: 2, Name: "nb", Identifier: "ib"}}

		sorted := this.sortFunc(menu)

		if !this.assertFunc(sorted) {
			t.Errorf("[%d] sort error", i)
		}
	}

}

func TestHomeNodeMenu(t *testing.T) {
	t.Parallel()
	s := setupMenuTests(t, menuPageSources,
		"canonifyURLs", true,
		"uglyURLs", false,
	)

	home := s.getPage(KindHome)
	homeMenuEntry := &MenuEntry{Name: home.Title, URL: home.URL()}

	for i, this := range []struct {
		menu           string
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"main", homeMenuEntry, true, false},
		{"doesnotexist", homeMenuEntry, false, false},
		{"main", &MenuEntry{Name: "Somewhere else", URL: "/somewhereelse"}, false, false},
		{"grandparent", findTestMenuEntryByID(s, "grandparent", "grandparentId"), false, true},
		{"grandparent", findTestMenuEntryByID(s, "grandparent", "parentId"), false, true},
		{"grandparent", findTestMenuEntryByID(s, "grandparent", "grandchildId"), true, false},
	} {

		isMenuCurrent := home.IsMenuCurrent(this.menu, this.menuItem)
		hasMenuCurrent := home.HasMenuCurrent(this.menu, this.menuItem)

		if isMenuCurrent != this.isMenuCurrent {
			fmt.Println("isMenuCurrent", isMenuCurrent)
			fmt.Printf("this: %#v\n", this)
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v for %q", i, isMenuCurrent, this.menuItem)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			fmt.Println("hasMenuCurrent", hasMenuCurrent)
			fmt.Printf("this: %#v\n", this)
			t.Errorf("[%d] Wrong result for menu %q menuItem %v for HasMenuCurrent: %v", i, this.menu, this.menuItem, hasMenuCurrent)
		}
	}
}

func TestHopefullyUniqueID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "i", (&MenuEntry{Identifier: "i", URL: "u", Name: "n"}).hopefullyUniqueID())
	assert.Equal(t, "u", (&MenuEntry{Identifier: "", URL: "u", Name: "n"}).hopefullyUniqueID())
	assert.Equal(t, "n", (&MenuEntry{Identifier: "", URL: "", Name: "n"}).hopefullyUniqueID())
}

func TestAddMenuEntryChild(t *testing.T) {
	t.Parallel()
	root := &MenuEntry{Weight: 1}
	root.addChild(&MenuEntry{Weight: 2})
	root.addChild(&MenuEntry{Weight: 1})
	assert.Equal(t, 2, len(root.Children))
	assert.Equal(t, 1, root.Children[0].Weight)
}

var testMenuIdentityMatcher = func(me *MenuEntry, id string) bool { return me.Identifier == id }
var testMenuNameMatcher = func(me *MenuEntry, id string) bool { return me.Name == id }

func findTestMenuEntryByID(s *Site, mn string, id string) *MenuEntry {
	return findTestMenuEntry(s, mn, id, testMenuIdentityMatcher)
}
func findTestMenuEntryByName(s *Site, mn string, id string) *MenuEntry {
	return findTestMenuEntry(s, mn, id, testMenuNameMatcher)
}

func findTestMenuEntry(s *Site, mn string, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	var found *MenuEntry
	if menu, ok := s.Menus[mn]; ok {
		for _, me := range *menu {

			if matcher(me, id) {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menu %s with id/name %s", mn, id))
				}
				found = me
			}

			descendant := findDescendantTestMenuEntry(me, id, matcher)
			if descendant != nil {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menu %s with id/name %s", mn, id))
				}
				found = descendant
			}
		}
	}
	return found
}

func findDescendantTestMenuEntry(parent *MenuEntry, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	var found *MenuEntry
	if parent.HasChildren() {
		for _, child := range parent.Children {

			if matcher(child, id) {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menuitem %s with id/name %s", parent.KeyName(), id))
				}
				found = child
			}

			descendant := findDescendantTestMenuEntry(child, id, matcher)
			if descendant != nil {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menuitem %s with id/name %s", parent.KeyName(), id))
				}
				found = descendant
			}
		}
	}
	return found
}

func setupMenuTests(t *testing.T, pageSources []source.ByteSource, configKeyValues ...interface{}) *Site {

	var (
		cfg, fs = newTestCfg()
	)

	menus, err := tomlToMap(confMenu1)
	require.NoError(t, err)

	cfg.Set("menu", menus["menu"])
	cfg.Set("baseURL", "http://foo.local/Zoo/")

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}

	for _, src := range pageSources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

	}

	return buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

}

func tomlToMap(s string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	_, err := toml.Decode(s, &data)
	return data, err
}
