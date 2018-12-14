package plugins

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
)

type testPage struct {
	title string
}

func (t *testPage) Title() string {
	return t.title
}

func TestStarlightUserScriptOverride(t *testing.T) {
	cfg, err := config.FromConfigString(`
plugin_dir = "plugins"
theme = "test"
`, "toml")

	if err != nil {
		t.Fatal(err)
	}
	d := &deps.Deps{Cfg: cfg}

	userDir := "plugins"
	err = os.Mkdir(userDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(userDir)
	themeDir := filepath.Join("themes", "test", "plugins")
	err = os.MkdirAll(themeDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("themes")
	err = ioutil.WriteFile(filepath.Join(userDir, "hello.star"), []byte(`output("hello user " + name + " from: " + page.Title())`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(themeDir, "hello.star"), []byte(`output("hello theme " + name + " from: " + page.Title())`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	ns := New(d)

	page := &testPage{title: "Moby Dick"}
	i, err := ns.Starlight("hello.star", "name", "bob", "page", page)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := i.(string)
	if !ok {
		t.Fatalf("expected to get string out of script, but got %v (%T)", i, i)
	}
	expected := "hello user bob from: Moby Dick"
	if s != expected {
		t.Fatalf("expected %q but got %q", expected, s)
	}
}

func TestStarlightThemeScript(t *testing.T) {
	cfg, err := config.FromConfigString(`
plugin_dir = "plugins"
theme = "test"
`, "toml")

	if err != nil {
		t.Fatal(err)
	}
	d := &deps.Deps{Cfg: cfg}

	userDir := "plugins"
	err = os.Mkdir(userDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(userDir)
	themeDir := filepath.Join("themes", "test", "plugins")
	err = os.MkdirAll(themeDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("themes")
	err = ioutil.WriteFile(filepath.Join(themeDir, "hello.star"), []byte(`output("hello theme " + name + " from: " + page.Title())`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	ns := New(d)

	page := &testPage{title: "Moby Dick"}
	i, err := ns.Starlight("hello.star", "name", "bob", "page", page)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := i.(string)
	if !ok {
		t.Fatalf("expected to get string out of script, but got %v (%T)", i, i)
	}
	expected := "hello theme bob from: Moby Dick"
	if s != expected {
		t.Fatalf("expected %q but got %q", expected, s)
	}
}
