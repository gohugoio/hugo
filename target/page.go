package target

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
)

type PagePublisher interface {
	Translator
	Publish(string, template.HTML) error
}

type PagePub struct {
	UglyUrls         bool
	DefaultExtension string
	PublishDir       string
}

func (pp *PagePub) Publish(path string, r io.Reader) (err error) {

	translated, err := pp.Translate(path)
	if err != nil {
		return
	}

	return helpers.WriteToDisk(translated, r, hugofs.DestinationFS)
}

func (pp *PagePub) Translate(src string) (dest string, err error) {
	if src == helpers.FilePathSeparator {
		if pp.PublishDir != "" {
			return filepath.Join(pp.PublishDir, "index.html"), nil
		}
		return "index.html", nil
	}

	dir, file := filepath.Split(src)
	ext := pp.extension(filepath.Ext(file))
	name := filename(file)
	if pp.PublishDir != "" {
		dir = filepath.Join(pp.PublishDir, dir)
	}

	if pp.UglyUrls || file == "index.html" || file == "404.html" {
		return filepath.Join(dir, fmt.Sprintf("%s%s", name, ext)), nil
	}

	return filepath.Join(dir, name, fmt.Sprintf("index%s", ext)), nil
}

func (pp *PagePub) extension(ext string) string {
	switch ext {
	case ".md", ".rst": // TODO make this list configurable.  page.go has the list of markup types.
		return ".html"
	}

	if ext != "" {
		return ext
	}

	if pp.DefaultExtension != "" {
		return pp.DefaultExtension
	}

	return ".html"
}
