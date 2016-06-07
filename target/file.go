// Copyright 2016 The Hugo Authors. All rights reserved.
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

	return helpers.WriteToDisk(translated, r, hugofs.Destination())
}

func (fs *Filesystem) Translate(src string) (dest string, err error) {
	return filepath.Join(fs.PublishDir, filepath.FromSlash(src)), nil
}

func filename(f string) string {
	ext := filepath.Ext(f)
	if ext == "" {
		return f
	}

	return f[:len(f)-len(ext)]
}
