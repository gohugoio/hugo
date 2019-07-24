package resources

import (
	"path/filepath"
	"testing"

	"image"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func newTestResourceSpec(assert *require.Assertions) *Spec {
	return newTestResourceSpecForBaseURL(assert, "https://example.com/")
}

func createTestCfg() *viper.Viper {
	cfg := viper.New()
	cfg.Set("resourceDir", "resources")
	cfg.Set("contentDir", "content")
	cfg.Set("dataDir", "data")
	cfg.Set("i18nDir", "i18n")
	cfg.Set("layoutDir", "layouts")
	cfg.Set("assetDir", "assets")
	cfg.Set("archetypeDir", "archetypes")
	cfg.Set("publishDir", "public")

	langs.LoadLanguageSettings(cfg, nil)
	mod, err := modules.CreateProjectModule(cfg)
	if err != nil {
		panic(err)
	}
	cfg.Set("allModules", modules.Modules{mod})

	return cfg
}

func newTestResourceSpecForBaseURL(assert *require.Assertions, baseURL string) *Spec {
	cfg := createTestCfg()
	cfg.Set("baseURL", baseURL)

	imagingCfg := map[string]interface{}{
		"resampleFilter": "linear",
		"quality":        68,
		"anchor":         "left",
	}

	cfg.Set("imaging", imagingCfg)

	fs := hugofs.NewMem(cfg)

	s, err := helpers.NewPathSpec(fs, cfg)
	assert.NoError(err)

	filecaches, err := filecache.NewCaches(s)
	assert.NoError(err)

	spec, err := NewSpec(s, filecaches, nil, output.DefaultFormats, media.DefaultTypes)
	assert.NoError(err)
	return spec
}

func newTargetPaths(link string) func() page.TargetPaths {
	return func() page.TargetPaths {
		return page.TargetPaths{
			SubResourceBaseTarget: filepath.FromSlash(link),
			SubResourceBaseLink:   link,
		}
	}
}

func newTestResourceOsFs(assert *require.Assertions) *Spec {
	cfg := createTestCfg()
	cfg.Set("baseURL", "https://example.com")

	workDir, _ := ioutil.TempDir("", "hugores")

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}

	cfg.Set("workingDir", workDir)

	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}
	fs.Source = afero.NewBasePathFs(hugofs.Os, workDir)

	s, err := helpers.NewPathSpec(fs, cfg)
	assert.NoError(err)

	filecaches, err := filecache.NewCaches(s)
	assert.NoError(err)

	spec, err := NewSpec(s, filecaches, nil, output.DefaultFormats, media.DefaultTypes)
	assert.NoError(err)
	return spec
}

func fetchSunset(assert *require.Assertions) *Image {
	return fetchImage(assert, "sunset.jpg")
}

func fetchImage(assert *require.Assertions, name string) *Image {
	spec := newTestResourceSpec(assert)
	return fetchImageForSpec(spec, assert, name)
}

func fetchImageForSpec(spec *Spec, assert *require.Assertions, name string) *Image {
	r := fetchResourceForSpec(spec, assert, name)
	assert.IsType(&Image{}, r)
	return r.(*Image)
}

func fetchResourceForSpec(spec *Spec, assert *require.Assertions, name string) resource.ContentResource {
	src, err := os.Open(filepath.FromSlash("testdata/" + name))
	assert.NoError(err)

	out, err := helpers.OpenFileForWriting(spec.Fs.Source, name)
	assert.NoError(err)
	_, err = io.Copy(out, src)
	out.Close()
	src.Close()
	assert.NoError(err)

	factory := newTargetPaths("/a")

	r, err := spec.New(ResourceSourceDescriptor{Fs: spec.Fs.Source, TargetPaths: factory, LazyPublish: true, SourceFilename: name})
	assert.NoError(err)

	return r.(resource.ContentResource)
}

func assertImageFile(assert *require.Assertions, fs afero.Fs, filename string, width, height int) {
	filename = filepath.Clean(filename)
	f, err := fs.Open(filename)
	assert.NoError(err)
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	assert.NoError(err)

	assert.Equal(width, config.Width)
	assert.Equal(height, config.Height)
}

func assertFileCache(assert *require.Assertions, fs afero.Fs, filename string, width, height int) {
	assertImageFile(assert, fs, filepath.Clean(filename), width, height)
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}
