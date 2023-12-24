// Copyright 2024 The Hugo Authors. All rights reserved.
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

package skeletons

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/afero"
)

//go:embed all:site/*
var siteFs embed.FS

//go:embed all:theme/*
var themeFs embed.FS

// CreateTheme creates a theme skeleton.
func CreateTheme(createpath string, sourceFs afero.Fs) error {
	if exists, _ := helpers.Exists(createpath, sourceFs); exists {
		return errors.New(createpath + " already exists")
	}
	return copyFiles(createpath, sourceFs, themeFs)
}

// CreateSite creates a site skeleton.
func CreateSite(createpath string, sourceFs afero.Fs, force bool, format string) error {
	format = strings.ToLower(format)
	if exists, _ := helpers.Exists(createpath, sourceFs); exists {
		if isDir, _ := helpers.IsDir(createpath, sourceFs); !isDir {
			return errors.New(createpath + " already exists but not a directory")
		}

		isEmpty, _ := helpers.IsEmpty(createpath, sourceFs)

		switch {
		case !isEmpty && !force:
			return errors.New(createpath + " already exists and is not empty. See --force.")
		case !isEmpty && force:
			var all []string
			fs.WalkDir(siteFs, ".", func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() && path != "." {
					all = append(all, path)
				}
				return nil
			})
			all = append(all, filepath.Join(createpath, "hugo."+format))
			for _, path := range all {
				if exists, _ := helpers.Exists(path, sourceFs); exists {
					return errors.New(path + " already exists")
				}
			}
		}
	}

	err := newSiteCreateConfig(sourceFs, createpath, format)
	if err != nil {
		return err
	}

	return copyFiles(createpath, sourceFs, siteFs)
}

func copyFiles(createpath string, sourceFs afero.Fs, skeleton embed.FS) error {
	return fs.WalkDir(skeleton, ".", func(path string, d fs.DirEntry, err error) error {
		_, slug, _ := strings.Cut(path, "/")
		if d.IsDir() {
			return sourceFs.MkdirAll(filepath.Join(createpath, slug), 0o777)
		} else {
			if filepath.Base(path) != ".gitkeep" {
				data, _ := fs.ReadFile(skeleton, path)
				return helpers.WriteToDisk(filepath.Join(createpath, slug), bytes.NewReader(data), sourceFs)
			}
			return nil
		}
	})
}

func newSiteCreateConfig(fs afero.Fs, createpath string, format string) (err error) {
	in := map[string]string{
		"baseURL":      "https://example.org/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, metadecoders.FormatFromString(format), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(createpath, "hugo."+format), &buf, fs)
}
