package resources_test

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
)

type specDescriptor struct {
	baseURL string
	c       *qt.C
	fs      afero.Fs
}

func newTestResourceSpec(desc specDescriptor) *resources.Spec {
	baseURL := desc.baseURL
	if baseURL == "" {
		baseURL = "https://example.com/"
	}

	afs := desc.fs
	if afs == nil {
		afs = afero.NewMemMapFs()
	}

	if hugofs.IsOsFs(afs) {
		panic("osFs not supported for this test")
	}

	if err := afs.MkdirAll("assets", 0o755); err != nil {
		panic(err)
	}

	cfg := config.New()
	cfg.Set("baseURL", baseURL)
	cfg.Set("publishDir", "public")

	imagingCfg := map[string]any{
		"resampleFilter": "linear",
		"quality":        68,
		"anchor":         "left",
	}

	cfg.Set("imaging", imagingCfg)
	d := testconfig.GetTestDeps(
		afs, cfg,
		func(d *deps.Deps) { d.Fs.PublishDir = hugofs.NewCreateCountingFs(d.Fs.PublishDir) },
	)

	desc.c.Cleanup(func() {
		if err := d.Close(); err != nil {
			panic(err)
		}
	})

	return d.ResourceSpec
}

func newTestResourceOsFs(c *qt.C) (*resources.Spec, string) {
	cfg := config.New()
	cfg.Set("baseURL", "https://example.com")

	workDir, err := os.MkdirTemp("", "hugores")
	c.Assert(err, qt.IsNil)
	c.Assert(workDir, qt.Not(qt.Equals), "")

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}

	cfg.Set("workingDir", workDir)

	os.MkdirAll(filepath.Join(workDir, "assets"), 0o755)

	d := testconfig.GetTestDeps(hugofs.Os, cfg)

	return d.ResourceSpec, workDir
}

func fetchSunset(c *qt.C) (*resources.Spec, images.ImageResource) {
	return fetchImage(c, "sunset.jpg")
}

func fetchImage(c *qt.C, name string) (*resources.Spec, images.ImageResource) {
	spec := newTestResourceSpec(specDescriptor{c: c})
	return spec, fetchImageForSpec(spec, c, name)
}

func fetchImageForSpec(spec *resources.Spec, c *qt.C, name string) images.ImageResource {
	r := fetchResourceForSpec(spec, c, name)
	img := r.(images.ImageResource)
	c.Assert(img, qt.Not(qt.IsNil))
	return img
}

func fetchResourceForSpec(spec *resources.Spec, c *qt.C, name string, targetPathAddends ...string) resource.ContentResource {
	b, err := os.ReadFile(filepath.FromSlash("testdata/" + name))
	c.Assert(err, qt.IsNil)
	open := hugio.NewOpenReadSeekCloser(hugio.NewReadSeekerNoOpCloserFromBytes(b))
	targetPath := name
	base := "/a/"
	r, err := spec.NewResource(resources.ResourceSourceDescriptor{
		LazyPublish:    true,
		NameNormalized: name, TargetPath: targetPath, BasePathRelPermalink: base, BasePathTargetPath: base, OpenReadSeekCloser: open,
		GroupIdentity: identity.Anonymous,
	})
	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.Not(qt.IsNil))

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
