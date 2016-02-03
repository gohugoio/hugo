// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	UglyURLs         bool
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
	dir, err := pp.TranslateRelative(src)
	if err != nil {
		return dir, err
	}
	if pp.PublishDir != "" {
		dir = filepath.Join(pp.PublishDir, dir)
	}
	return dir, nil
}

func (pp *PagePub) TranslateRelative(src string) (dest string, err error) {
	if src == helpers.FilePathSeparator {
		return "index.html", nil
	}

	dir, file := filepath.Split(src)
	isRoot := dir == ""
	ext := pp.extension(filepath.Ext(file))
	name := filename(file)

	if pp.UglyURLs || file == "index.html" || (isRoot && file == "404.html") {
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
