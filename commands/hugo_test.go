package commands

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/hugolib"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func clearConfigFile() {
	viper.SetConfigFile("")
	theme = ""
	cfgFile = ""
	destination = ""
	source = ""
	baseURL = ""
}

func TestUserErrorCommandNewUserError(t *testing.T) {
	defer viper.Reset()
	expectedUserError := commandError{"test\n", true}
	actualUserError := newUserError("test")
	assert.Equal(t, expectedUserError, actualUserError)
}

func TestSystemErrorCommandNewSystemError(t *testing.T) {
	defer viper.Reset()
	expectedSystemError := commandError{"test\n", false}
	actualSystemError := newSystemError("test")
	assert.Equal(t, expectedSystemError, actualSystemError)
}

func TestUserErrorCommandNewFormattedUserError(t *testing.T) {
	defer viper.Reset()
	expectedUserError := commandError{"test:%d 12\n", true}
	actualUserError := newUserError("test:%d", 12)
	assert.Equal(t, expectedUserError, actualUserError)
}

func TestSystemErrorCommandNewFormattedSystemError(t *testing.T) {
	defer viper.Reset()
	expectedSystemError := commandError{"test:12", false}
	actualSystemError := newSystemErrorF("test:%d", 12)
	assert.Equal(t, expectedSystemError, actualSystemError)
	if actualSystemError != expectedSystemError {
		t.Errorf("Actual (%v) returned command did not equal expected (%v) command.", actualSystemError, expectedSystemError)
	}
}

func TestUserErrorReturnError(t *testing.T) {
	defer viper.Reset()
	cError := commandError{"test error", false}
	expectedError := "test error"
	actualError := cError.Error()
	assert.Equal(t, expectedError, actualError)
}

func TestIsUserError(t *testing.T) {
	defer viper.Reset()
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
	defer viper.Reset()
	defer clearConfigFile()
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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.Source().MkdirAll(basepath, os.ModePerm|os.ModeDir)
	paths := []string{
		filepath.Join(basepath, "public"),
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
	}

	for _, path := range paths {
		hugofs.Source().MkdirAll(path, os.ModePerm|os.ModeDir)
	}

	file, err := hugofs.Source().Create(configPath)
	if err != nil {
		t.Error(err.Error())
	}
	err = hugofs.Source().Chmod(configPath, os.ModeTemporary|os.ModePerm)
	if err != nil {
		t.Error(err.Error())
	}

	basicConfig := `
        contentdir = ` + strconv.Quote(paths[2]) + `
        layoutdir = ` + strconv.Quote(paths[1]) + `
        publishdir = ` + strconv.Quote(paths[0]) + `
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
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	hugofs.InitMemFs()
	hugofs.SetDestination(afero.NewMemMapFs())
	hugofs.SetSource(afero.NewMemMapFs())
	viper.SetFs(hugofs.Source())
	err := doNewSite(basepath, true)
	assert.Nil(t, err)
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.Source().Create(configPath)
	hugofs.Source().Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`
	afero.WriteFile(hugofs.Source(), configPath, []byte(basicConfig), os.ModeTemporary|os.ModePerm)
	source = basepath

	InitializeConfig()

	assert.Equal(t, true, strings.Contains(viper.GetString("WorkingDir"), source), "WorkingDir config did not include test directory.")
}

func TestInitializeConfigDefaultValues(t *testing.T) {
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	hugofs.InitMemFs()
	hugofs.SetDestination(afero.NewMemMapFs())
	hugofs.SetSource(afero.NewMemMapFs())
	viper.SetFs(hugofs.Source())
	err := doNewSite(basepath, true)
	assert.Nil(t, err)
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.Source().Create(configPath)
	hugofs.Source().Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`
	afero.WriteFile(hugofs.Source(), configPath, []byte(basicConfig), os.ModeTemporary|os.ModePerm)
	source = basepath
	cfgFile = configPath
	hugoCmdV = HugoCmd
	baseURL = "http://localhost"
	destination = "public"
	cacheDir = "testcache"
	theme = "testtheme"
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

func TestInitializeConfig(t *testing.T) {
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	hugofs.InitMemFs()
	hugofs.SetDestination(afero.NewMemMapFs())
	hugofs.SetSource(afero.NewMemMapFs())
	viper.SetFs(hugofs.Source())
	err := doNewSite(basepath, true)
	assert.Nil(t, err)
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.Source().Create(configPath)
	hugofs.Source().Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `baseurl = "http://localhost"`
	afero.WriteFile(hugofs.Source(), configPath, []byte(basicConfig), os.ModeTemporary|os.ModePerm)
	source = basepath
	cfgFile = configPath
	hugoCmdV = HugoCmd
	baseURL = "http://localhost"
	destination = "testpublish"
	cacheDir = "testcache"
	theme = "testtheme"
	InitializeConfig()
	assert.Equal(t, baseURL, viper.GetString("BaseURL"), "BaseUrl was not set.")
	assert.Equal(t, theme, viper.GetString("theme"), "Theme was not set.")
	assert.Equal(t, destination, viper.GetString("PublishDir"), "Publish Dir was not set.")
	assert.Equal(t, strings.Contains(viper.GetString("WorkingDir"), source), true, "Source did not contain test directory.")
	assert.Equal(t, cacheDir, viper.GetString("CacheDir"), "CacheDir was not set.")
}

func TestInitializeConfigWithErrorMissingConfigFile(t *testing.T) {
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	hugofs.InitMemFs()
	hugofs.SetDestination(afero.NewMemMapFs())
	hugofs.SetSource(afero.NewMemMapFs())
	viper.SetFs(hugofs.Source())
	err := doNewSite(basepath, true)
	assert.Nil(t, err)
	hugoCmdV = HugoCmd
	err = InitializeConfig()
	if err == nil {
		t.Fail()
	}
}

func TestInitializeConfigBaseURLNotSet(t *testing.T) {
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	configPath := filepath.Join(basepath, "config.toml")
	logFilePath := filepath.Join(basepath, "mylog.log")
	hugofs.Source().MkdirAll(basepath, os.ModePerm|os.ModeDir)
	hugofs.Source().MkdirAll(basepath+"/testcache", os.ModePerm|os.ModeDir)
	file, _ := hugofs.Source().Create(configPath)
	hugofs.Source().Chmod(configPath, os.ModeTemporary|os.ModePerm)
	basicConfig := `base = "http://localhost"`
	file.WriteString(basicConfig)
	cfgFile = configPath
	hugoCmdV = HugoCmd
	source = basepath
	destination = "testpublish"
	source = "testsource"
	cacheDir = "testcache"
	jww.SetLogFile(logFilePath)
	InitializeConfig()
	content, _ := ioutil.ReadFile(logFilePath)
	if !strings.Contains(string(content), "No 'baseurl' set in configuration or as a flag.") {
		t.Fail()
	}
}

func TestThemeVsHugoVersionMismatch(t *testing.T) {
	defer viper.Reset()
	defer clearConfigFile()
	defer helpers.ResetConfigProvider()
	defer hugofs.InitDefaultFs()
	defer hugoCmdV.ResetCommands()
	defer hugoCmdV.ResetFlags()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blogSuffix := strconv.Itoa(r.Intn(math.MaxInt8))
	basepath := filepath.Join(os.TempDir(), "blog"+blogSuffix)
	// This is needed because GetThemeDir is actually OS based. But then switches to Hugofs...
	hugofs.Source().MkdirAll(filepath.Join(basepath, "themes", "testtheme"), os.ModePerm|os.ModeDir)
	themeTomlPath := filepath.Join(basepath, "themes", "testtheme", "theme.toml")
	themeTomlFile, _ := hugofs.Source().Create(themeTomlPath)
	hugofs.Source().Chmod(themeTomlPath, os.ModeTemporary|os.ModePerm)
	themeToml := `min_version = 1.9`
	themeTomlFile.Write([]byte(themeToml))
	theme = "testtheme"
	themeDir := filepath.Join(basepath, "themes")
	viper.Set("theme", theme)
	viper.Set("themesDir", themeDir)
	mismatch, reqVersion := isThemeVsHugoVersionMismatch()
	assert.Equal(t, true, mismatch, "There should be a mismatch.")
	assert.Equal(t, "1.9", reqVersion, "Required version should be 1.9 (Set in the theme toml file.)")
}

func containesCommand(cmd *cobra.Command, arr []*cobra.Command) bool {
	for _, c := range arr {
		if c == cmd {
			return true
		}
	}

	return false
}
