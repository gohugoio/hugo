package target

import (
	"fmt"
	"io"
	"path"
)

type Publisher interface {
	Publish(string, io.Reader) error
}

type Translator interface {
	Translate(string) (string, error)
}

type Filesystem struct {
	UglyUrls         bool
	DefaultExtension string
}

func (fs *Filesystem) Translate(src string) (dest string, err error) {
	if fs.UglyUrls {
		return src, nil
	}

	dir, file := path.Split(src)
	ext := fs.extension(path.Ext(file))
	name := filename(file)

	return path.Join(dir, name, fmt.Sprintf("index%s", ext)), nil
}

func (fs *Filesystem) extension(ext string) string {
	if ext != "" {
		return ext
	}

	if fs.DefaultExtension != "" {
		return fs.DefaultExtension
	}

	return ".html"
}

func filename(f string) string {
	ext := path.Ext(f)
	if ext == "" {
		return f
	}

	return f[:len(f)-len(ext)]
}
