package hugolib

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var deepHierarchyMenu [][]byte = [][]byte{
	[]byte(`+++
type = "post"
title = "0"
[menu.mnuDeep]
	identifier = "id005"
	weight = 5
    name = "0"

+++
01 content`),

	[]byte(`+++
type = "post"
title = "1"
[menu.mnuDeep]
	identifier = "id01"
	weight = 10
    name = "1_name"

+++
1 Content`),

	[]byte(`+++
type= "post"
title = "11"
[menu.mnuDeep]
	parent = "id01"
	identifier = "id011"
	weight = 22
    name = "11"

+++
11 content`),

	[]byte(`+++
type= "post"
title = "111"
[menu.mnuDeep]
	parent = "id011"
	identifier = "id0111"
	weight = 23
    name = "111"

+++
111 content`),

	[]byte(`+++
type= "post"
title = "112"
[menu.mnuDeep]
	parent = "id011"
	identifier = "id0112"
	weight = 24
    name = "112"

+++
112 content`),

	[]byte(`+++
type= "post"
title = "2"
[menu.mnuDeep]
	identifier = "id009"
	weight = 11
    name = "2"

+++
20 content`),

	[]byte(`+++
type= "post"
title = "21"
weight = 100
[menu.mnuDeep]
	parent = "id009"
	identifier = "id02"
	weight = 20
    name = "21_name"

+++
21 content
`),

	[]byte(`+++
type= "post"
title = "211"
[menu.mnuDeep]
	parent = "id02"
	identifier = "id0211"
	weight = 22
    name = "211"

+++
211 content`),

	[]byte(`+++
type= "post"
title = "2111"
[menu.mnuDeep]
	parent = "id0211"
	identifier = "id02111"
	weight = 22
    name = "2111"

+++
2111 content`),

	[]byte(`+++
type= "post"
title = "212"
weight = 500
[menu.mnuDeep]
	parent = "id02"
	identifier = "id022"
	weight = 24
    name = "212"

+++
212 content`),

	[]byte(`+++
type= "post"
title = "Third"
[menu.mnuDeep]
	identifier = "id03"
	weight = 30
    name = "ThirdArticle"

+++
ThirdArticle`),

	[]byte(`+++
type= "post"
title = "Fourth"
[menu.mnuDeep]
	identifier = "id04"
	weight = 40
    name = "FourthArticle"

+++
FourthArticle`),
}

var MENU_DEEP_HIERA_SOURCES = []source.ByteSource{}

func init() {

	for i := 0; i < len(deepHierarchyMenu); i++ {
		bs := source.ByteSource{}
		bs.Content = deepHierarchyMenu[i]
		bs.Name = filepath.FromSlash(fmt.Sprintf("deep-hiera-sect/doc%02v.md", i))
		MENU_DEEP_HIERA_SOURCES = append(MENU_DEEP_HIERA_SOURCES, bs)
	}

}

func EnableAsNeeded_TestWriteOutAsFiles(t *testing.T) {
	td, err := ioutil.TempDir("", "menu_hiera_testfiles")
	if err != nil {
		log.Printf("Could not create dir for testfiles %v", err)
		return
	}
	for i := 0; i < len(deepHierarchyMenu); i++ {
		pth := filepath.FromSlash(fmt.Sprintf("%v/doc%02v.md", td, i))
		err := ioutil.WriteFile(pth, deepHierarchyMenu[i], 0644)
		if err != nil {
			log.Printf("Could not write testfile %v - %v", pth, err)
			return
		}
	}
}

func TestTraverseMenu(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	s := setupMenuTests(t, MENU_DEEP_HIERA_SOURCES)

	assert.Equal(t, len(s.Pages), 12, "Not enough pages")

	(*s.Menus["mnuDeep"]).ByWeight()
	// (*s.Menus["mnuDeep"]).ByName()

	assert.Equal(t, len(*s.Menus["mnuDeep"]), 5, "Level 1 of mnuDeep should have x entries")

	me := (*s.Menus["mnuDeep"])[0]

	for i, this := range []struct {
		name  string
		isNil bool
	}{
		{"0", false},
		{"1_name", false},
		{"11", false},
		{"111", false},
		{"112", false},
		{"2", false},
		{"21_name", false},
		{"211", false},
		{"2111", false},
		{"212", false},
		{"ThirdArticle", false},
		{"FourthArticle", false},
		{"no successor", true},
	} {

		assert.Equal(t, this.isNil, me == nil, fmt.Sprintf("next menu entry must be nil; iter %v", i))
		if me != nil {
			assert.Equal(t, this.name, me.Name, fmt.Sprintf("want another menu entry as next; iter %v", i))
			me = (*s.Menus["mnuDeep"]).Next(*me)
		}

	}

	me = (*s.Menus["mnuDeep"])[len(*s.Menus["mnuDeep"])-1]

	for i, this := range []struct {
		name  string
		isNil bool
	}{

		{"FourthArticle", false},
		{"ThirdArticle", false},
		{"212", false},
		{"2111", false},
		{"211", false},
		{"21_name", false},
		{"2", false},
		{"112", false},
		{"111", false},
		{"11", false},
		{"1_name", false},
		{"0", false},
		{"no precursor", true},
	} {

		assert.Equal(t, this.isNil, me == nil, fmt.Sprintf("prev menu entry must be nil; iter %v", i))
		if me != nil {
			assert.Equal(t, this.name, me.Name, fmt.Sprintf("want another menu entry as prev; iter %v", i))
			me = (*s.Menus["mnuDeep"]).Prev(*me)
		}

	}

	me = (*s.Menus["mnuDeep"]).searchRecursive("id02111")

	for i, this := range []struct {
		name  string
		isNil bool
	}{

		{"2111", false},
		{"211", false},
		{"21_name", false},
		{"2", false},
		{"no parent", true},
	} {

		assert.Equal(t, this.isNil, me == nil, fmt.Sprintf("up menu entry must be nil; iter %v", i))
		if me != nil {
			assert.Equal(t, this.name, me.Name, fmt.Sprintf("want another menu entry as up; iter %v", i))
			me = (*s.Menus["mnuDeep"]).Up(*me)
		}

	}

}
