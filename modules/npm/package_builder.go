// Copyright 2020 The Hugo Authors. All rights reserved.
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

package npm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/helpers"
)

const (
	dependenciesKey    = "dependencies"
	devDependenciesKey = "devDependencies"

	packageJSONName = "package.json"

	packageJSONTemplate = `{
  "name": "%s",
  "version": "%s"
}`
)

func Pack(fs afero.Fs, fis []hugofs.FileMetaInfo) error {
	var b *packageBuilder

	// Have a package.hugo.json?
	fi, err := fs.Stat(files.FilenamePackageHugoJSON)
	if err != nil {
		// Have a package.json?
		fi, err = fs.Stat(packageJSONName)
		if err == nil {
			// Preserve the original in package.hugo.json.
			if err = hugio.CopyFile(fs, packageJSONName, files.FilenamePackageHugoJSON); err != nil {
				return errors.Wrap(err, "npm pack: failed to copy package file")
			}
		} else {
			// Create one.
			name := "project"
			// Use the Hugo site's folder name as the default name.
			// The owner can change it later.
			rfi, err := fs.Stat("")
			if err == nil {
				name = rfi.Name()
			}
			packageJSONContent := fmt.Sprintf(packageJSONTemplate, name, "0.1.0")
			if err = afero.WriteFile(fs, files.FilenamePackageHugoJSON, []byte(packageJSONContent), 0666); err != nil {
				return err
			}
			fi, err = fs.Stat(files.FilenamePackageHugoJSON)
			if err != nil {
				return err
			}
		}
	}

	meta := fi.(hugofs.FileMetaInfo).Meta()
	masterFilename := meta.Filename
	f, err := meta.Open()
	if err != nil {
		return errors.Wrap(err, "npm pack: failed to open package file")
	}
	b = newPackageBuilder(meta.Module, f)
	f.Close()

	for _, fi := range fis {
		if fi.IsDir() {
			// We only care about the files in the root.
			continue
		}

		if fi.Name() != files.FilenamePackageHugoJSON {
			continue
		}

		meta := fi.(hugofs.FileMetaInfo).Meta()

		if meta.Filename == masterFilename {
			continue
		}

		f, err := meta.Open()
		if err != nil {
			return errors.Wrap(err, "npm pack: failed to open package file")
		}
		b.Add(meta.Module, f)
		f.Close()
	}

	if b.Err() != nil {
		return errors.Wrap(b.Err(), "npm pack: failed to build")
	}

	// Replace the dependencies in the original template with the merged set.
	b.originalPackageJSON[dependenciesKey] = b.dependencies
	b.originalPackageJSON[devDependenciesKey] = b.devDependencies
	var commentsm map[string]interface{}
	comments, found := b.originalPackageJSON["comments"]
	if found {
		commentsm = maps.ToStringMap(comments)
	} else {
		commentsm = make(map[string]interface{})
	}
	commentsm[dependenciesKey] = b.dependenciesComments
	commentsm[devDependenciesKey] = b.devDependenciesComments
	b.originalPackageJSON["comments"] = commentsm

	// Write it out to the project package.json
	packageJSONData := new(bytes.Buffer)
	encoder := json.NewEncoder(packageJSONData)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", strings.Repeat(" ", 2))
	if err := encoder.Encode(b.originalPackageJSON); err != nil {
		return errors.Wrap(err, "npm pack: failed to marshal JSON")
	}

	if err := afero.WriteFile(fs, packageJSONName, packageJSONData.Bytes(), 0666); err != nil {
		return errors.Wrap(err, "npm pack: failed to write package.json")
	}

	return nil
}

func newPackageBuilder(source string, first io.Reader) *packageBuilder {
	b := &packageBuilder{
		devDependencies:         make(map[string]interface{}),
		devDependenciesComments: make(map[string]interface{}),
		dependencies:            make(map[string]interface{}),
		dependenciesComments:    make(map[string]interface{}),
	}

	m := b.unmarshal(first)
	if b.err != nil {
		return b
	}

	b.addm(source, m)
	b.originalPackageJSON = m

	return b
}

type packageBuilder struct {
	err error

	// The original package.hugo.json.
	originalPackageJSON map[string]interface{}

	devDependencies         map[string]interface{}
	devDependenciesComments map[string]interface{}
	dependencies            map[string]interface{}
	dependenciesComments    map[string]interface{}
}

func (b *packageBuilder) Add(source string, r io.Reader) *packageBuilder {
	if b.err != nil {
		return b
	}

	m := b.unmarshal(r)
	if b.err != nil {
		return b
	}

	b.addm(source, m)

	return b
}

func (b *packageBuilder) addm(source string, m map[string]interface{}) {
	if source == "" {
		source = "project"
	}

	// The version selection is currently very simple.
	// We may consider minimal version selection or something
	// after testing this out.
	//
	// But for now, the first version string for a given dependency wins.
	// These packages will be added by order of import (project, module1, module2...),
	// so that should at least give the project control over the situation.
	if devDeps, found := m[devDependenciesKey]; found {
		mm := maps.ToStringMapString(devDeps)
		for k, v := range mm {
			if _, added := b.devDependencies[k]; !added {
				b.devDependencies[k] = v
				b.devDependenciesComments[k] = source
			}
		}
	}

	if deps, found := m[dependenciesKey]; found {
		mm := maps.ToStringMapString(deps)
		for k, v := range mm {
			if _, added := b.dependencies[k]; !added {
				b.dependencies[k] = v
				b.dependenciesComments[k] = source
			}
		}
	}
}

func (b *packageBuilder) unmarshal(r io.Reader) map[string]interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal(helpers.ReaderToBytes(r), &m)
	if err != nil {
		b.err = err
	}
	return m
}

func (b *packageBuilder) Err() error {
	return b.err
}
