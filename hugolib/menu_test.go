package hugolib

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
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
	identifier="xml"`
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
		Identity = "Two"

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

type testMenuState struct {
	site       *Site
	oldMenu    interface{}
	oldBaseUrl interface{}
}

func TestPageMenu(t *testing.T) {
	ts := setupMenuTests(t)
	defer resetMenuTestState(ts)

	if len(ts.site.Pages) != 3 {
		t.Fatalf("Posts not created, expected 3 got %d", len(ts.site.Pages))
	}

	first := ts.site.Pages[0]
	second := ts.site.Pages[1]
	third := ts.site.Pages[2]

	pOne := ts.findTestMenuEntryByName("p_one", "One")
	pTwo := ts.findTestMenuEntryByName("p_two", "Two")

	for i, this := range []struct {
		menu           string
		page           *Page
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"p_one", first, pOne, true, false},
		{"p_one", first, pTwo, false, false},
		{"p_one", second, pTwo, true, false},
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

func TestTaxonomyNodeMenu(t *testing.T) {
	ts := setupMenuTests(t)
	defer resetMenuTestState(ts)

	for i, this := range []struct {
		menu           string
		taxInfo        taxRenderInfo
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			ts.findTestMenuEntryById("tax", "1"), true, false},
		{"tax", taxRenderInfo{key: "key", singular: "one", plural: "two"},
			ts.findTestMenuEntryById("tax", "2"), true, false},
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

	menuEntryXml := ts.findTestMenuEntryById("tax", "xml")

	if strings.HasSuffix(menuEntryXml.Url, "/") {
		t.Error("RSS menu item should not be padded with trailing slash")
	}
}

func TestHomeNodeMenu(t *testing.T) {
	ts := setupMenuTests(t)
	defer resetMenuTestState(ts)

	home := ts.site.newHomeNode()
	homeMenuEntry := &MenuEntry{Name: home.Title, Url: string(home.Permalink)}

	for i, this := range []struct {
		menu           string
		menuItem       *MenuEntry
		isMenuCurrent  bool
		hasMenuCurrent bool
	}{
		{"main", homeMenuEntry, true, false},
		{"doesnotexist", homeMenuEntry, false, false},
		{"main", &MenuEntry{Name: "Somewhere else", Url: "/somewhereelse"}, false, false},
		{"grandparent", ts.findTestMenuEntryById("grandparent", "grandparentId"), false, false},
		{"grandparent", ts.findTestMenuEntryById("grandparent", "parentId"), false, true},
		{"grandparent", ts.findTestMenuEntryById("grandparent", "grandchildId"), true, false},
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

func (ts testMenuState) findTestMenuEntryById(mn string, id string) *MenuEntry {
	return ts.findTestMenuEntry(mn, id, testMenuIdentityMatcher)
}
func (ts testMenuState) findTestMenuEntryByName(mn string, id string) *MenuEntry {
	return ts.findTestMenuEntry(mn, id, testMenuNameMatcher)
}

func (ts testMenuState) findTestMenuEntry(mn string, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	if menu, ok := ts.site.Menus[mn]; ok {
		for _, me := range *menu {

			if matcher(me, id) {
				return me
			}

			descendant := ts.findDescendantTestMenuEntry(me, id, matcher)
			if descendant != nil {
				return descendant
			}
		}
	}
	return nil
}

func (ts testMenuState) findDescendantTestMenuEntry(parent *MenuEntry, id string, matcher func(me *MenuEntry, id string) bool) *MenuEntry {
	if parent.HasChildren() {
		for _, child := range parent.Children {

			if matcher(child, id) {
				return child
			}

			descendant := ts.findDescendantTestMenuEntry(child, id, matcher)
			if descendant != nil {
				return descendant
			}
		}
	}
	return nil
}

func getTestMenuState(s *Site, t *testing.T) *testMenuState {
	menuState := &testMenuState{site: s, oldBaseUrl: viper.Get("baseurl"), oldMenu: viper.Get("menu")}

	menus, err := tomlToMap(CONF_MENU1)

	if err != nil {
		t.Fatalf("Unable to Read menus: %v", err)
	}

	viper.Set("menu", menus["menu"])
	viper.Set("baseurl", "http://foo.local/zoo/")

	return menuState
}

func setupMenuTests(t *testing.T) *testMenuState {
	s := createTestSite()
	testState := getTestMenuState(s, t)
	testSiteSetup(s, t)

	return testState
}

func resetMenuTestState(state *testMenuState) {
	viper.Set("menu", state.oldMenu)
	viper.Set("baseurl", state.oldBaseUrl)
}

func createTestSite() *Site {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}

	s := &Site{
		Target: target,
		Source: &source.InMemorySource{ByteSource: MENU_PAGE_SOURCES},
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
	var data map[string]interface{} = make(map[string]interface{})
	if _, err := toml.Decode(s, &data); err != nil {
		return nil, err
	}

	return data, nil

}
