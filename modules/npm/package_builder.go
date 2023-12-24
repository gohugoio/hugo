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
	"io/fs"
	"strings"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/hugofs/files"

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

func Pack(sourceFs, assetsWithDuplicatesPreservedFs afero.Fs) error {
	var b *packageBuilder

	// Have a package.hugo.json?
	fi, err := sourceFs.Stat(files.FilenamePackageHugoJSON)
	if err != nil {
		// Have a package.json?
		fi, err = sourceFs.Stat(packageJSONName)
		if err == nil {
			// Preserve the original in package.hugo.json.
			if err = hugio.CopyFile(sourceFs, packageJSONName, files.FilenamePackageHugoJSON); err != nil {
				return fmt.Errorf("npm pack: failed to copy package file: %w", err)
			}
		} else {
			// Create one.
			name := "project"
			// Use the Hugo site's folder name as the default name.
			// The owner can change it later.
			rfi, err := sourceFs.Stat("")
			if err == nil {
				name = rfi.Name()
			}
			packageJSONContent := fmt.Sprintf(packageJSONTemplate, name, "0.1.0")
			if err = afero.WriteFile(sourceFs, files.FilenamePackageHugoJSON, []byte(packageJSONContent), 0o666); err != nil {
				return err
			}
			fi, err = sourceFs.Stat(files.FilenamePackageHugoJSON)
			if err != nil {
				return err
			}
		}
	}

	meta := fi.(hugofs.FileMetaInfo).Meta()
	masterFilename := meta.Filename
	f, err := meta.Open()
	if err != nil {
		return fmt.Errorf("npm pack: failed to open package file: %w", err)
	}
	b = newPackageBuilder(meta.Module, f)
	f.Close()

	d, err := assetsWithDuplicatesPreservedFs.Open(files.FolderJSConfig)
	if err != nil {
		return nil
	}

	fis, err := d.(fs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return fmt.Errorf("npm pack: failed to read assets: %w", err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
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
			return fmt.Errorf("npm pack: failed to open package file: %w", err)
		}
		b.Add(meta.Module, f)
		f.Close()
	}

	if b.Err() != nil {
		return fmt.Errorf("npm pack: failed to build: %w", b.Err())
	}

	// Replace the dependencies in the original template with the merged set.
	b.originalPackageJSON[dependenciesKey] = b.dependencies
	b.originalPackageJSON[devDependenciesKey] = b.devDependencies
	var commentsm map[string]any
	comments, found := b.originalPackageJSON["comments"]
	if found {
		commentsm = maps.ToStringMap(comments)
	} else {
		commentsm = make(map[string]any)
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
		return fmt.Errorf("npm pack: failed to marshal JSON: %w", err)
	}

	if err := afero.WriteFile(sourceFs, packageJSONName, packageJSONData.Bytes(), 0o666); err != nil {
		return fmt.Errorf("npm pack: failed to write package.json: %w", err)
	}

	return nil
}

func newPackageBuilder(source string, first io.Reader) *packageBuilder {
	b := &packageBuilder{
		devDependencies:         make(map[string]any),
		devDependenciesComments: make(map[string]any),
		dependencies:            make(map[string]any),
		dependenciesComments:    make(map[string]any),
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
	originalPackageJSON map[string]any

	devDependencies         map[string]any
	devDependenciesComments map[string]any
	dependencies            map[string]any
	dependenciesComments    map[string]any
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

func (b *packageBuilder) addm(source string, m map[string]any) {
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

func (b *packageBuilder) unmarshal(r io.Reader) map[string]any {
	m := make(map[string]any)
	err := json.Unmarshal(helpers.ReaderToBytes(r), &m)
	if err != nil {
		b.err = err
	}
	return m
}

func (b *packageBuilder) Err() error {
	return b.err
}
