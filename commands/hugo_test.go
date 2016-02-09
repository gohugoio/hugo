package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUserErrorCommandNewUserError(t *testing.T) {
	expectedUserError := commandError{"test\n", true}
	actualUserError := newUserError("test")
	assert.Equal(t, expectedUserError, actualUserError)
}

func TestSystemErrorCommandNewSystemError(t *testing.T) {
	expectedSystemError := commandError{"test\n", false}
	actualSystemError := newSystemError("test")
	assert.Equal(t, expectedSystemError, actualSystemError)
}

func TestUserErrorCommandNewFormattedUserError(t *testing.T) {
	expectedUserError := commandError{"test:12", true}
	actualUserError := newUserErrorF("test:%d", 12)
	assert.Equal(t, expectedUserError, actualUserError)
}

func TestSystemErrorCommandNewFormattedSystemError(t *testing.T) {
	expectedSystemError := commandError{"test:12", false}
	actualSystemError := newSystemErrorF("test:%d", 12)
	assert.Equal(t, expectedSystemError, actualSystemError)
	if actualSystemError != expectedSystemError {
		t.Errorf("Actual (%v) returned command did not equal expected (%v) command.", actualSystemError, expectedSystemError)
	}
}

func TestUserErrorReturnError(t *testing.T) {
	cError := commandError{"test error", false}
	expectedError := "test error"
	actualError := cError.Error()
	assert.Equal(t, expectedError, actualError)
}

func TestIsUserError(t *testing.T) {
	var errors = []struct {
		err error
	}{
		{fmt.Errorf("argument")},
		{fmt.Errorf("flag")},
		{fmt.Errorf("shorthand")},
		{commandError{"test", true}},
	}

	for _, v := range errors {
		if !isUserError(v.err) {
			t.Errorf("Error (%v) should have been flagged as a user error.", v)
		}
	}
}

func TestExecuteAddedCommands(t *testing.T) {
	var defaultCommands = []struct {
		cmd *cobra.Command
	}{
		{serverCmd},
		{versionCmd},
		{configCmd},
		{checkCmd},
		{benchmarkCmd},
		{convertCmd},
		{newCmd},
		{listCmd},
		{undraftCmd},
		{importCmd},
		{genCmd},
	}

	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.OsFs)
	hugofs.SourceFs.MkdirAll(basepath, os.ModePerm|os.ModeDir)
	paths := []string{
		filepath.Join(basepath, "public"),
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
	}

	for _, path := range paths {
		hugofs.SourceFs.MkdirAll(path, os.ModePerm|os.ModeDir)
	}

	file, err := hugofs.SourceFs.Create(configPath)
	if err != nil {
		t.Error(err.Error())
	}
	err = hugofs.SourceFs.Chmod(configPath, os.ModeTemporary|os.ModePerm)
	if err != nil {
		t.Error(err.Error())
	}
	basicConfig := `
        contentdir = "` + string(paths[2]) + `"
        layoutdir = "` + string(paths[1]) + `"
        publishdir = "` + string(paths[0]) + `"
        baseurl = "http://localhost"
        canonifyurls = true`

	_, err = file.WriteString(basicConfig)
	if err != nil {
		t.Error(err.Error())
	}
	cfgFile = configPath

	// Workaround for go test -coverprofile=cover.out since viper parses command-line
	// arguments.
	if len(os.Args) > 1 {
		os.Args = os.Args[:1]
	}
	Execute()
	commands := HugoCmd.Commands()

	for _, dc := range defaultCommands {
		if !containesCommand(dc.cmd, commands) {
			t.Errorf("Default command (%v) was not contained in the list of default commands. List: %#v", dc.cmd, commands)
		}
	}

}

func TestInitializeConfigSetsWorkingDirToSourceIfSet(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, os.ModePerm|os.ModeDir)
	file, _ := hugofs.SourceFs.Create(configPath)
	hugofs.SourceFs.Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`

	file.WriteString(basicConfig)
	cfgFile = configPath
	hugoCmdV = HugoCmd
	source = "testdir123"

	InitializeConfig()

	assert.Equal(t, true, strings.Contains(viper.GetString("WorkingDir"), source), "WorkingDir config did not include test directory.")
}

func TestInitializeConfigDefaultValues(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.OsFs)
	hugofs.SourceFs.MkdirAll(basepath, os.ModePerm|os.ModeDir)
	file, _ := hugofs.SourceFs.Create(configPath)
	hugofs.SourceFs.Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`

	file.WriteString(basicConfig)
	cfgFile = configPath
	hugoCmdV = HugoCmd

	var defaultValues = []struct {
		key   string
		value interface{}
	}{
		{"cleanDestinationDir", false},
		{"Watch", false},
		{"MetaDataFormat", "toml"},
		{"DisableRSS", false},
		{"DisableSitemap", false},
		{"DisableRobotsTXT", false},
		{"ContentDir", "content"},
		{"LayoutDir", "layouts"},
		{"StaticDir", "static"},
		{"ArchetypeDir", "archetypes"},
		{"PublishDir", "public"},
		{"DataDir", "data"},
		{"ThemesDir", "themes"},
		{"DefaultLayout", "post"},
		{"BuildDrafts", false},
		{"BuildFuture", false},
		{"UglyURLs", false},
		{"Verbose", false},
		{"IgnoreCache", false},
		{"CanonifyURLs", false},
		{"RelativeURLs", false},
		{"RemovePathAccents", false},
		{"Taxonomies", map[string]string{"tag": "tags", "category": "categories"}},
		{"Permalinks", make(hugolib.PermalinkOverrides, 0)},
		{"Sitemap", hugolib.Sitemap{Priority: -1, Filename: "sitemap.xml"}},
		{"DefaultExtension", "html"},
		{"PygmentsStyle", "monokai"},
		{"PygmentsUseClasses", false},
		{"PygmentsCodeFences", false},
		{"PygmentsOptions", ""},
		{"DisableLiveReload", false},
		{"PluralizeListTitles", true},
		{"PreserveTaxonomyNames", false},
		{"ForceSyncStatic", false},
		{"FootnoteAnchorPrefix", ""},
		{"FootnoteReturnLinkContents", ""},
		{"NewContentEditor", ""},
		{"Paginate", 10},
		{"PaginatePath", "page"},
		{"Blackfriday", helpers.NewBlackfriday()},
		{"RSSUri", "index.xml"},
		{"SectionPagesMenu", ""},
		{"DisablePathToLower", false},
		{"HasCJKLanguage", false},
	}

	InitializeConfig()
	for _, v := range defaultValues {
		assert.Equal(t, v.value, viper.Get(v.key), "Expected value key(%s) did not equal default configuration value.", v.key)
	}
}

func TestInitializeConfigCommandCanOverwriteDefaults(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, os.ModePerm|os.ModeDir)
	file, _ := hugofs.SourceFs.Create(configPath)
	hugofs.SourceFs.Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`

	file.WriteString(basicConfig)
	cfgFile = configPath
	hugoCmdV = HugoCmd

	var newFlagValues = []struct {
		key   string
		value string
	}{
		{"cleanDestinationDir", "true"},
		{"buildFuture", "true"},
		{"buildDrafts", "true"},
		{"uglyURLs", "true"},
		{"canonifyURLs", "true"},
		{"disableRSS", "true"},
		{"disableSitemap", "true"},
		{"disableRobotsTXT", "true"},
		{"preserveTaxonomyNames", "true"},
		{"ignoreCache", "true"},
		{"forceSyncStatic", "true"},
		{"pluralizeListTitles", "true"},
		{"noTimes", "true"},
	}

	for _, v := range newFlagValues {
		hugoCmdV.Flags().Set(v.key, v.value)
	}

	InitializeConfig(hugoCmdV)
	for _, v := range newFlagValues {
		assert.Equal(t, v.value, viper.GetString(v.key), "Expected value key(%s) did not equal default configuration value.", v.key)
	}
}

func TestInitializeConfig(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, os.ModePerm|os.ModeDir)
	hugofs.SourceFs.MkdirAll(basepath+"/testcache", os.ModePerm|os.ModeDir)
	file, _ := hugofs.SourceFs.Create(configPath)
	hugofs.SourceFs.Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`

	file.WriteString(basicConfig)
	cfgFile = configPath
	hugoCmdV = HugoCmd
	baseURL = "http://localhost"
	theme = "testtheme"
	destination = "testpublish"
	source = "testsource"
	cacheDir = "testcache"
	InitializeConfig()
	assert.Equal(t, baseURL, viper.GetString("BaseURL"), "BaseUrl was not set.")
	assert.Equal(t, theme, viper.GetString("theme"), "Theme was not set.")
	assert.Equal(t, destination, viper.GetString("PublishDir"), "Publish Dir was not set.")
	assert.Equal(t, strings.Contains(viper.GetString("WorkingDir"), source), true, "Source did not contain test directory.")
	assert.Equal(t, cacheDir, viper.GetString("CacheDir"), "CacheDir was not set.")
}

func containesCommand(cmd *cobra.Command, arr []*cobra.Command) bool {
	for _, c := range arr {
		if c == cmd {
			return true
		}
	}

	return false
}
