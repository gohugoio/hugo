package commands

import (
	"github.com/spf13/afero"
	"github.com/spf13/hugo/hugofs"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"os"
)

// Issue #1133
func TestNewContentPathSectionWithForwardSlashes(t *testing.T) {
	p, s := newContentPathSection("/post/new.md")
	assert.Equal(t, filepath.FromSlash("/post/new.md"), p)
	assert.Equal(t, "post", s)
}

func checkNewSiteInited(basepath string, t *testing.T) {
	paths := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "config.toml"),
	}

	for _, path := range paths {
		_, err := hugofs.SourceFs.Stat(path)
		assert.Nil(t, err)
	}
}

func TestDoNewSite(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	hugofs.SourceFs = new(afero.MemMapFs)
	err := doNewSite(basepath, false)
	assert.Nil(t, err)

	checkNewSiteInited(basepath, t)
}

func TestDoNewSite_error_base_exists(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, 777)
	err := doNewSite(basepath, false)
	assert.NotNil(t, err)
}

func TestDoNewSite_force_empty_dir(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, 777)
	err := doNewSite(basepath, true)
	assert.Nil(t, err)

	checkNewSiteInited(basepath, t)
}

func TestDoNewSite_error_force_dir_inside_exists(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	contentPath := filepath.Join(basepath, "content")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(contentPath, 777)
	err := doNewSite(basepath, true)
	assert.NotNil(t, err)
}

func TestDoNewSite_error_force_config_inside_exists(t *testing.T) {
	basepath := filepath.Join(os.TempDir(), "blog")
	configPath := filepath.Join(basepath, "config.toml")
	hugofs.SourceFs = new(afero.MemMapFs)
	hugofs.SourceFs.MkdirAll(basepath, 777)
	hugofs.SourceFs.Create(configPath)
	err := doNewSite(basepath, true)
	assert.NotNil(t, err)
}
