package hugolib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
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

var MENU_PAGE_SOURCES = []source.ByteSource{
	{"sect/doc1.md", MENU_PAGE_1},
	{"sect/doc2.md", MENU_PAGE_2},
	{"sect/doc3.md", MENU_PAGE_3},
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

	ts := setupMenuTests(t, menuPageSources)
	defer resetMenuTestState(ts)

	assert.Equal(t, 3, len(ts.site.Pages), "Not enough pages")

	me1 := ts.findTestMenuEntryByID("m1", "i1")
	me2 := ts.findTestMenuEntryByID("m1", "i2")

	assert.NotNil(t, me1)
	assert.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.Url, "doc1"))
	assert.True(t, strings.Contains(me2.Url, "doc2"))

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
	ts := setupMenuTests(t, menuPageSources)
	defer resetMenuTestState(ts)

	assert.Equal(t, 3, len(ts.site.Pages), "Not enough pages")

	me1 := ts.findTestMenuEntryByName("m1", "n1")
	me2 := ts.findTestMenuEntryByName("m1", "n2")

	assert.NotNil(t, me1)
	assert.NotNil(t, me2)

	assert.True(t, strings.Contains(me1.Url, "doc1"))
	assert.True(t, strings.Contains(me2.Url, "doc2"))

}

func TestPageMenu(t *testing.T) {
	ts := setupMenuTests(t, MENU_PAGE_SOURCES)
	defer resetMenuTestState(ts)

	if len(ts.site.Pages) != 3 {
		t.Fatalf("Posts not created, expected 3 got %d", len(ts.site.Pages))
	}

	first := ts.site.Pages[0]
	second := ts.site.Pages[1]
	third := ts.site.Pages[2]

	pOne := ts.findTestMenuEntryByName("p_one", "One")
	pTwo := ts.findTestMenuEntryByID("p_two", "Two")

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
	ts := setupMenuTests(t, MENU_PAGE_SOURCES)
	defer resetMenuTestState(ts)

	me := ts.findTestMenuEntryByID("hash", "hash")

	assert.NotNil(t, me)

	assert.Equal(t, "/Zoo/resource/#anchor", me.Url)
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
	viper.Set("CanonifyURLs", canonifyURLs)
	viper.Set("UglyURLs", uglyURLs)

	ts := setupMenuTests(t, MENU_PAGE_SOURCES)
	defer resetMenuTestState(ts)

	unicodeRussian := ts.findTestMenuEntryByID("unicode", "unicode-russian")

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

	assert.Equal(t, expected, unicodeRussian.Url, "uglyURLs[%t]", uglyURLs)
}

func TestTaxonomyNodeMenu(t *testing.T) {
	viper.Set("CanonifyURLs", true)
	ts := setupMenuTests(t, MENU_PAGE_SOURCES)
	defer resetMenuTestState(ts)

	for i, this := range []struct {
		menu           string
		taxInfo        taxRenderInfo
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			ts.findTestMenuEntryByID("tax", "1"), true, false},
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			ts.findTestMenuEntryByID("tax", "2"), true, false},
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			&MenuEntry{Name: "Somewhere else", Url: "/somewhereelse"}, false, false},
	} {

		n, _ := ts.site.newTaxonomyNode(this.taxInfo)

		isMenuCurrent := n.IsMenuCurrent(this.menu, this.menuItem)
		hasMenuCurrent := n.HasMenuCurrent(this.menu, this.menuItem)

		if isMenuCurrent != this.isMenuCurrent {
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v", i, isMenuCurrent)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			t.Errorf("[%d] Wrong result for menuItem %v for HasMenuCurrent: %v", i, this.menuItem, hasMenuCurrent)
		}

	}

	menuEntryXML := ts.findTestMenuEntryByID("tax", "xml")

	if strings.HasSuffix(menuEntryXML.Url, "/") {
		t.Error("RSS menu item should not be padded with trailing slash")
	}
}

func TestHomeNodeMenu(t *testing.T) {
	ts := setupMenuTests(t, MENU_PAGE_SOURCES)
	defer resetMenuTestState(ts)

	home := ts.site.newHomeNode()
	homeMenuEntry := &MenuEntry{Name: home.Title, Url: home.Url}

	for i, this := range []struct {
		menu           string
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"main", homeMenuEntry, true, false},
		{"doesnotexist", homeMenuEntry, false, false},
		{"main", &MenuEntry{Name: "Somewhere else", Url: "/somewhereelse"}, false, false},
		{"grandparent", ts.findTestMenuEntryByID("grandparent", "grandparentId"), false, false},
		{"grandparent", ts.findTestMenuEntryByID("grandparent", "parentId"), false, true},
		{"grandparent", ts.findTestMenuEntryByID("grandparent", "grandchildId"), true, false},
	} {

		isMenuCurrent := home.IsMenuCurrent(this.menu, this.menuItem)
		hasMenuCurrent := home.HasMenuCurrent(this.menu, this.menuItem)

		if isMenuCurrent != this.isMenuCurrent {
			t.Errorf("[%d] Wrong result from IsMenuCurrent: %v", i, isMenuCurrent)
		}

		if hasMenuCurrent != this.hasMenuCurrent {
			t.Errorf("[%d] Wrong result for menuItem %v for HasMenuCurrent: %v", i, this.menuItem, hasMenuCurrent)
		}
	}
}

var testMenuIdentityMatcher = func(me *MenuEntry, id string) bool { return me.Identifier == id }
var testMenuNameMatcher = func(me *MenuEntry, id string) bool { return me.Name == id }

func (ts testMenuState) findTestMenuEntryByID(mn string, id string) *MenuEntry {
	return ts.findTestMenuEntry(mn, id, testMenuIdentityMatcher)
}
func (ts testMenuState) findTestMenuEntryByName(mn string, id string) *MenuEntry {
	return ts.findTestMenuEntry(mn, id, testMenuNameMatcher)
}

func (ts testMenuState) findTestMenuEntry(mn string, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	var found *MenuEntry
	if menu, ok := ts.site.Menus[mn]; ok {
		for _, me := range *menu {

			if matcher(me, id) {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menu %s with id/name %s", mn, id))
				}
				found = me
			}

			descendant := ts.findDescendantTestMenuEntry(me, id, matcher)
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

func (ts testMenuState) findDescendantTestMenuEntry(parent *MenuEntry, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	var found *MenuEntry
	if parent.HasChildren() {
		for _, child := range parent.Children {

			if matcher(child, id) {
				if found != nil {
					panic(fmt.Sprintf("Duplicate menu entry in menuitem %s with id/name %s", parent.KeyName(), id))
				}
				found = child
			}

			descendant := ts.findDescendantTestMenuEntry(child, id, matcher)
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

func getTestMenuState(s *Site, t *testing.T) *testMenuState {
	menuState := &testMenuState{site: s, oldBaseURL: viper.Get("baseurl"), oldMenu: viper.Get("menu")}

	menus, err := tomlToMap(CONF_MENU1)

	if err != nil {
		t.Fatalf("Unable to Read menus: %v", err)
	}

	viper.Set("menu", menus["menu"])
	viper.Set("baseurl", "http://foo.local/Zoo/")

	return menuState
}

func setupMenuTests(t *testing.T, pageSources []source.ByteSource) *testMenuState {
	s := createTestSite(pageSources)
	testState := getTestMenuState(s, t)
	testSiteSetup(s, t)

	return testState
}

func resetMenuTestState(state *testMenuState) {
	viper.Set("menu", state.oldMenu)
	viper.Set("baseurl", state.oldBaseURL)
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
	s.Shortcodes = make(map[string]ShortcodeFunc)

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
