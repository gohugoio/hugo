package target

import (
	"io"
	"path/filepath"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
)

type Publisher interface {
	Publish(string, io.Reader) error
}

type Translator interface {
	Translate(string) (string, error)
}

// TODO(bep) consider other ways to solve this.
type OptionalTranslator interface {
	TranslateRelative(string) (string, error)
}

type Output interface {
	Publisher
	Translator
}

type Filesystem struct {
	PublishDir string
}

func (fs *Filesystem) Publish(path string, r io.Reader) (err error) {
	translated, err := fs.Translate(path)
	if err != nil {
		return
	}

	return helpers.WriteToDisk(translated, r, hugofs.DestinationFS)
}

func (fs *Filesystem) Translate(src string) (dest string, err error) {
	return filepath.Join(fs.PublishDir, src), nil
}

func (fs *Filesystem) extension(ext string) string {
	return ext
}

func filename(f string) string {
	ext := filepath.Ext(f)
	if ext == "" {
		return f
	}

	return f[:len(f)-len(ext)]
}
