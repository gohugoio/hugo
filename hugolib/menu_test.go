package hugolib

import (
	"fmt"
	"strings"
	"testing"

	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kr/pretty"
	"github.com/spf13/afero"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	CONF_MENU1 = `
[[menu.main]]
    name = "Go Home"
    url = "/"
	weight = 1
	pre = "<div>"
	post = "</div>"
[[menu.main]]
    name = "Blog"
    url = "/posts"
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
    url = "/two/key"
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

var MENU_PAGE_1 = []byte(`+++
title = "One"
[menu]
	[menu.p_one]
weight = 1
+++
Front Matter with Menu Pages`)

var MENU_PAGE_2 = []byte(`+++
title = "Two"
weight = 2
[menu]
	[menu.p_one]
	[menu.p_two]
		identifier = "Two"

+++
Front Matter with Menu Pages`)

var MENU_PAGE_3 = []byte(`+++
title = "Three"
weight = 3
[menu]
	[menu.p_two]
		Name = "Three"
		Parent = "Two"
+++
Front Matter with Menu Pages`)

var MENU_PAGE_4 = []byte(`+++
title = "Four"
weight = 4
[menu]
	[menu.p_two]
		Name = "Four"
		Parent = "Three"
+++
Front Matter with Menu Pages`)

var MENU_PAGE_SOURCES = []source.ByteSource{
	{filepath.FromSlash("sect/doc1.md"), MENU_PAGE_1},
	{filepath.FromSlash("sect/doc2.md"), MENU_PAGE_2},
	{filepath.FromSlash("sect/doc3.md"), MENU_PAGE_3},
}

var MENU_PAGE_SECTIONS_SOURCES = []source.ByteSource{
	{filepath.FromSlash("first/doc1.md"), MENU_PAGE_1},
	{filepath.FromSlash("first/doc2.md"), MENU_PAGE_2},
	{filepath.FromSlash("second-section/doc3.md"), MENU_PAGE_3},
	{filepath.FromSlash("Fish and Chips/doc4.md"), MENU_PAGE_4},
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

type testMenuState struct {
	site       *Site
	oldMenu    interface{}
	oldBaseURL interface{}
}

// Issue 817 - identifier should trump everything
func TestPageMenuWithIdentifier(t *testing.T) {

	toml := []source.ByteSource{
		{"sect/doc1.md", tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i1")},
		{"sect/doc2.md", tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i2")},
		{"sect/doc3.md", tstCreateMenuPageWithIdentifierTOML("t1", "m1", "i2")}, // duplicate
	}

	yaml := []source.ByteSource{
		{"sect/doc1.md", tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i1")},
		{"sect/doc2.md", tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i2")},
		{"sect/doc3.md", tstCreateMenuPageWithIdentifierYAML("t1", "m1", "i2")}, // duplicate
	}

	doTestPageMenuWithIdentifier(t, toml)
	doTestPageMenuWithIdentifier(t, yaml)

}

func doTestPageMenuWithIdentifier(t *testing.T, menuPageSources []source.ByteSource) {

	viper.Reset()
	defer viper.Reset()

	s := setupMenuTests(t, menuPageSources)

	assert.Equal(t, 3, len(s.Pages), "Not enough pages")

	me1 := findTestMenuEntryByID(s, "m1", "i1")
	me2 := findTestMenuEntryByID(s, "m1", "i2")

	assert.NotNil(t, me1)
	assert.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.URL, "doc1"), me1.URL)
	assert.True(t, strings.Contains(me2.URL, "doc2") || strings.Contains(me2.URL, "doc3"), me2.URL)

}

// Issue 817 contd - name should be second identifier in
func TestPageMenuWithDuplicateName(t *testing.T) {

	toml := []source.ByteSource{
		{"sect/doc1.md", tstCreateMenuPageWithNameTOML("t1", "m1", "n1")},
		{"sect/doc2.md", tstCreateMenuPageWithNameTOML("t1", "m1", "n2")},
		{"sect/doc3.md", tstCreateMenuPageWithNameTOML("t1", "m1", "n2")}, // duplicate
	}

	yaml := []source.ByteSource{
		{"sect/doc1.md", tstCreateMenuPageWithNameYAML("t1", "m1", "n1")},
		{"sect/doc2.md", tstCreateMenuPageWithNameYAML("t1", "m1", "n2")},
		{"sect/doc3.md", tstCreateMenuPageWithNameYAML("t1", "m1", "n2")}, // duplicate
	}

	doTestPageMenuWithDuplicateName(t, toml)
	doTestPageMenuWithDuplicateName(t, yaml)

}

func doTestPageMenuWithDuplicateName(t *testing.T, menuPageSources []source.ByteSource) {
	viper.Reset()
	defer viper.Reset()

	s := setupMenuTests(t, menuPageSources)

	assert.Equal(t, 3, len(s.Pages), "Not enough pages")

	me1 := findTestMenuEntryByName(s, "m1", "n1")
	me2 := findTestMenuEntryByName(s, "m1", "n2")

	assert.NotNil(t, me1)
	assert.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.URL, "doc1"), me1.URL)
	assert.True(t, strings.Contains(me2.URL, "doc2") || strings.Contains(me2.URL, "doc3"), me2.URL)

}

func TestPageMenu(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	s := setupMenuTests(t, MENU_PAGE_SOURCES)

	if len(s.Pages) != 3 {
		t.Fatalf("Posts not created, expected 3 got %d", len(s.Pages))
	}

	first := s.Pages[0]
	second := s.Pages[1]
	third := s.Pages[2]

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

// issue #888
func TestMenuWithHashInURL(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	s := setupMenuTests(t, MENU_PAGE_SOURCES)

	me := findTestMenuEntryByID(s, "hash", "hash")

	assert.NotNil(t, me)

	assert.Equal(t, "/Zoo/resource/#anchor", me.URL)
}

// issue #719
func TestMenuWithUnicodeURLs(t *testing.T) {

	for _, uglyURLs := range []bool{true, false} {
		for _, canonifyURLs := range []bool{true, false} {
			doTestMenuWithUnicodeURLs(t, canonifyURLs, uglyURLs)
		}
	}
}

func doTestMenuWithUnicodeURLs(t *testing.T, canonifyURLs, uglyURLs bool) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("CanonifyURLs", canonifyURLs)
	viper.Set("UglyURLs", uglyURLs)

	s := setupMenuTests(t, MENU_PAGE_SOURCES)

	unicodeRussian := findTestMenuEntryByID(s, "unicode", "unicode-russian")

	expectedBase := "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0"

	if !canonifyURLs {
		expectedBase = "/Zoo" + expectedBase
	}

	var expected string
	if uglyURLs {
		expected = expectedBase + ".html"
	} else {
		expected = expectedBase + "/"
	}

	assert.Equal(t, expected, unicodeRussian.URL, "uglyURLs[%t]", uglyURLs)
}

// Issue #1114
func TestSectionPagesMenu(t *testing.T) {

	doTestSectionPagesMenu(true, t)
	doTestSectionPagesMenu(false, t)
}

func doTestSectionPagesMenu(canonifyUrls bool, t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("SectionPagesMenu", "spm")

	viper.Set("CanonifyURLs", canonifyUrls)
	s := setupMenuTests(t, MENU_PAGE_SECTIONS_SOURCES)

	assert.Equal(t, 3, len(s.Sections))

	firstSectionPages := s.Sections["first"]
	assert.Equal(t, 2, len(firstSectionPages))
	secondSectionPages := s.Sections["second-section"]
	assert.Equal(t, 1, len(secondSectionPages))
	fishySectionPages := s.Sections["fish-and-chips"]
	assert.Equal(t, 1, len(fishySectionPages))

	nodeFirst := s.newSectionListNode("First", "first", firstSectionPages)
	nodeSecond := s.newSectionListNode("Second Section", "second-section", secondSectionPages)
	nodeFishy := s.newSectionListNode("Fish and Chips", "fish-and-chips", fishySectionPages)
	firstSectionMenuEntry := findTestMenuEntryByID(s, "spm", "first")
	secondSectionMenuEntry := findTestMenuEntryByID(s, "spm", "second-section")
	fishySectionMenuEntry := findTestMenuEntryByID(s, "spm", "Fish and Chips")

	assert.NotNil(t, firstSectionMenuEntry)
	assert.NotNil(t, secondSectionMenuEntry)
	assert.NotNil(t, nodeFirst)
	assert.NotNil(t, nodeSecond)
	assert.NotNil(t, fishySectionMenuEntry)
	assert.NotNil(t, nodeFishy)

	assert.True(t, nodeFirst.IsMenuCurrent("spm", firstSectionMenuEntry))
	assert.False(t, nodeFirst.IsMenuCurrent("spm", secondSectionMenuEntry))
	assert.False(t, nodeFirst.IsMenuCurrent("spm", fishySectionMenuEntry))
	assert.True(t, nodeFishy.IsMenuCurrent("spm", fishySectionMenuEntry))
	assert.Equal(t, "Fish and Chips", fishySectionMenuEntry.Name)

	for _, p := range firstSectionPages {
		assert.True(t, p.Page.HasMenuCurrent("spm", firstSectionMenuEntry))
		assert.False(t, p.Page.HasMenuCurrent("spm", secondSectionMenuEntry))
	}

	for _, p := range secondSectionPages {
		assert.False(t, p.Page.HasMenuCurrent("spm", firstSectionMenuEntry))
		assert.True(t, p.Page.HasMenuCurrent("spm", secondSectionMenuEntry))
	}

	for _, p := range fishySectionPages {
		assert.False(t, p.Page.HasMenuCurrent("spm", firstSectionMenuEntry))
		assert.False(t, p.Page.HasMenuCurrent("spm", secondSectionMenuEntry))
		assert.True(t, p.Page.HasMenuCurrent("spm", fishySectionMenuEntry))
	}
}

func TestTaxonomyNodeMenu(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("CanonifyURLs", true)
	s := setupMenuTests(t, MENU_PAGE_SOURCES)

	for i, this := range []struct {
		menu           string
		taxInfo        taxRenderInfo
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			findTestMenuEntryByID(s, "tax", "1"), true, false},
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			findTestMenuEntryByID(s, "tax", "2"), true, false},
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			&MenuEntry{Name: "Somewhere else", URL: "/somewhereelse"}, false, false},
	} {

		n, _ := s.newTaxonomyNode(this.taxInfo)

		isMenuCurrent := n.IsMenuCurrent(this.menu, this.menuItem)
		hasMenuCurrent := n.HasMenuCurrent(this.menu, this.menuItem)

		if isMenuCurrent != this.isMenuCurrent {
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v", i, isMenuCurrent)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			t.Errorf("[%d] Wrong result for menuItem %v for HasMenuCurrent: %v", i, this.menuItem, hasMenuCurrent)
		}

	}

	menuEntryXML := findTestMenuEntryByID(s, "tax", "xml")

	if strings.HasSuffix(menuEntryXML.URL, "/") {
		t.Error("RSS menu item should not be padded with trailing slash")
	}
}

func TestHomeNodeMenu(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("CanonifyURLs", true)
	viper.Set("UglyURLs", true)

	s := setupMenuTests(t, MENU_PAGE_SOURCES)

	home := s.newHomeNode()
	homeMenuEntry := &MenuEntry{Name: home.Title, URL: home.URL}

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
			pretty.Println("this:", this)
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v for %q", i, isMenuCurrent, this.menu)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			fmt.Println("hasMenuCurrent", hasMenuCurrent)
			pretty.Println("this:", this)
			t.Errorf("[%d] Wrong result for menu %q menuItem %v for HasMenuCurrent: %v", i, this.menu, this.menuItem, hasMenuCurrent)
		}
	}
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

func setupTestMenuState(s *Site, t *testing.T) {
	menus, err := tomlToMap(CONF_MENU1)

	if err != nil {
		t.Fatalf("Unable to Read menus: %v", err)
	}

	viper.Set("menu", menus["menu"])
	viper.Set("baseurl", "http://foo.local/Zoo/")
}

func setupMenuTests(t *testing.T, pageSources []source.ByteSource) *Site {
	s := createTestSite(pageSources)
	setupTestMenuState(s, t)
	testSiteSetup(s, t)

	return s
}

func createTestSite(pageSources []source.ByteSource) *Site {
	hugofs.DestinationFS = new(afero.MemMapFs)

	s := &Site{
		Source: &source.InMemorySource{ByteSource: pageSources},
	}
	return s
}

func testSiteSetup(s *Site, t *testing.T) {
	s.Menus = Menus{}
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}
}

func tomlToMap(s string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	if _, err := toml.Decode(s, &data); err != nil {
		return nil, err
	}

	return data, nil

}
