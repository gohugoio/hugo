package resources

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/htesting/hqt"

	"image"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func newTestResourceSpec(c *qt.C) *Spec {
	return newTestResourceSpecForBaseURL(c, "https://example.com/")
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

func newTestResourceSpecForBaseURL(c *qt.C, baseURL string) *Spec {
	cfg := createTestCfg()
	cfg.Set("baseURL", baseURL)

	imagingCfg := map[string]interface{}{
		"resampleFilter": "linear",
		"quality":        68,
		"anchor":         "left",
	}

	cfg.Set("imaging", imagingCfg)

	fs := hugofs.NewMem(cfg)

	s, err := helpers.NewPathSpec(fs, cfg, nil)
	c.Assert(err, qt.IsNil)

	filecaches, err := filecache.NewCaches(s)
	c.Assert(err, qt.IsNil)

	spec, err := NewSpec(s, filecaches, nil, output.DefaultFormats, media.DefaultTypes)
	c.Assert(err, qt.IsNil)
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

func newTestResourceOsFs(c *qt.C) *Spec {
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

	s, err := helpers.NewPathSpec(fs, cfg, nil)
	c.Assert(err, qt.IsNil)

	filecaches, err := filecache.NewCaches(s)
	c.Assert(err, qt.IsNil)

	spec, err := NewSpec(s, filecaches, nil, output.DefaultFormats, media.DefaultTypes)
	c.Assert(err, qt.IsNil)
	return spec

}

func fetchSunset(c *qt.C) *Image {
	return fetchImage(c, "sunset.jpg")
}

func fetchImage(c *qt.C, name string) *Image {
	spec := newTestResourceSpec(c)
	return fetchImageForSpec(spec, c, name)
}

func fetchImageForSpec(spec *Spec, c *qt.C, name string) *Image {
	r := fetchResourceForSpec(spec, c, name)
	c.Assert(r, hqt.IsSameType, &Image{})
	return r.(*Image)
}

func fetchResourceForSpec(spec *Spec, c *qt.C, name string) resource.ContentResource {
	src, err := os.Open(filepath.FromSlash("testdata/" + name))
	c.Assert(err, qt.IsNil)

	out, err := helpers.OpenFileForWriting(spec.Fs.Source, name)
	c.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	out.Close()
	src.Close()
	c.Assert(err, qt.IsNil)

	factory := newTargetPaths("/a")

	r, err := spec.New(ResourceSourceDescriptor{Fs: spec.Fs.Source, TargetPaths: factory, LazyPublish: true, SourceFilename: name})
	c.Assert(err, qt.IsNil)

	return r.(resource.ContentResource)
}

func assertImageFile(c *qt.C, fs afero.Fs, filename string, width, height int) {
	filename = filepath.Clean(filename)
	f, err := fs.Open(filename)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	c.Assert(err, qt.IsNil)

	c.Assert(config.Width, qt.Equals, width)
	c.Assert(config.Height, qt.Equals, height)
}

func assertFileCache(c *qt.C, fs afero.Fs, filename string, width, height int) {
	assertImageFile(c, fs, filepath.Clean(filename), width, height)
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}
