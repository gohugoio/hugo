package resources

import (
	"path/filepath"
	"testing"

	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func newTestResourceSpec(assert *require.Assertions) *Spec {
	return newTestResourceSpecForBaseURL(assert, "https://example.com/")
}

func newTestResourceSpecForBaseURL(assert *require.Assertions, baseURL string) *Spec {
	cfg := viper.New()
	cfg.Set("baseURL", baseURL)
	cfg.Set("resourceDir", "resources")
	cfg.Set("contentDir", "content")
	cfg.Set("dataDir", "data")
	cfg.Set("i18nDir", "i18n")
	cfg.Set("layoutDir", "layouts")
	cfg.Set("assetDir", "assets")
	cfg.Set("archetypeDir", "archetypes")
	cfg.Set("publishDir", "public")

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

func newTestResourceOsFs(assert *require.Assertions) *Spec {
	cfg := viper.New()
	cfg.Set("baseURL", "https://example.com")

	workDir, err := ioutil.TempDir("", "hugores")

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}

	cfg.Set("workingDir", workDir)
	cfg.Set("resourceDir", "resources")
	cfg.Set("contentDir", "content")
	cfg.Set("dataDir", "data")
	cfg.Set("i18nDir", "i18n")
	cfg.Set("layoutDir", "layouts")
	cfg.Set("assetDir", "assets")
	cfg.Set("archetypeDir", "archetypes")
	cfg.Set("publishDir", "public")

	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}

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

	out, err := helpers.OpenFileForWriting(spec.BaseFs.Content.Fs, name)
	assert.NoError(err)
	_, err = io.Copy(out, src)
	out.Close()
	src.Close()
	assert.NoError(err)

	factory := func(s string) string {
		return path.Join("/a", s)
	}

	r, err := spec.New(ResourceSourceDescriptor{TargetPathBuilder: factory, SourceFilename: name})
	assert.NoError(err)

	return r.(resource.ContentResource)
}

func assertImageFile(assert *require.Assertions, fs afero.Fs, filename string, width, height int) {
	f, err := fs.Open(filename)
	if err != nil {
		printFs(fs, "", os.Stdout)
	}
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

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			s := path
			if lang, ok := info.(hugofs.LanguageAnnouncer); ok {
				s = s + "\t" + lang.Lang()
			}
			if fp, ok := info.(hugofs.FilePather); ok {
				s += "\tFilename: " + fp.Filename() + "\tBase: " + fp.BaseDir()
			}
			fmt.Fprintln(w, "    ", s)
		}
		return nil
	})
}
