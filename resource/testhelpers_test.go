package resource

import (
	"path/filepath"
	"testing"

	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
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
	cfg.Set("resourceDir", "/res")

	imagingCfg := map[string]interface{}{
		"resampleFilter": "linear",
		"quality":        68,
		"anchor":         "left",
	}

	cfg.Set("imaging", imagingCfg)

	fs := hugofs.NewMem(cfg)

	s, err := helpers.NewPathSpec(fs, cfg)

	assert.NoError(err)

	spec, err := NewSpec(s, media.DefaultTypes)
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

	contentDir := "base"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", contentDir)
	cfg.Set("resourceDir", filepath.Join(workDir, "res"))

	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}

	s, err := helpers.NewPathSpec(fs, cfg)

	assert.NoError(err)

	spec, err := NewSpec(s, media.DefaultTypes)
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

func fetchResourceForSpec(spec *Spec, assert *require.Assertions, name string) Resource {
	src, err := os.Open("testdata/" + name)
	assert.NoError(err)

	workingDir := spec.Cfg.GetString("workingDir")
	f := filepath.Join(workingDir, name)

	out, err := spec.Fs.Source.Create(f)
	assert.NoError(err)
	_, err = io.Copy(out, src)
	out.Close()
	src.Close()
	assert.NoError(err)

	factory := func(s string) string {
		return path.Join("/a", s)
	}

	r, err := spec.NewResourceFromFilename(factory, "/public", f, name)
	assert.NoError(err)

	return r
}
func assertFileCache(assert *require.Assertions, fs *hugofs.Fs, filename string, width, height int) {
	f, err := fs.Source.Open(filepath.Join("/res/_gen/images", filename))
	assert.NoError(err)
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	assert.NoError(err)

	assert.Equal(width, config.Width)
	assert.Equal(height, config.Height)
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}
