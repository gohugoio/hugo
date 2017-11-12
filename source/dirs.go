// Copyright 2017 The Hugo Authors. All rights reserved.
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

package source

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

// Dirs holds the source directories for a given build.
// In case where there are more than one of a kind, the order matters:
// It will be used to construct a union filesystem, so the right-most directory
// will "win" on duplicates. Typically, the theme version will be the first.
type Dirs struct {
	logger   *jww.Notepad
	pathSpec *helpers.PathSpec

	staticDirs    []string
	AbsStaticDirs []string

	publishDir string
}

// NewDirs creates a new dirs with the given configuration and filesystem.
func NewDirs(fs *hugofs.Fs, cfg config.Provider, logger *jww.Notepad) (*Dirs, error) {
	ps, err := helpers.NewPathSpec(fs, cfg)
	if err != nil {
		return nil, err
	}

	d := &Dirs{pathSpec: ps, logger: logger}

	return d, d.init(cfg)

}

func (d *Dirs) init(cfg config.Provider) error {

	var (
		statics []string
	)

	if d.pathSpec.Theme() != "" {
		statics = append(statics, filepath.Join(d.pathSpec.ThemesDir(), d.pathSpec.Theme(), "static"))
	}

	_, isLanguage := cfg.(*helpers.Language)
	languages, hasLanguages := cfg.Get("languagesSorted").(helpers.Languages)

	if !isLanguage && !hasLanguages {
		return errors.New("missing languagesSorted in config")
	}

	if !isLanguage {
		// Merge all the static dirs.
		for _, l := range languages {
			addend, err := d.staticDirsFor(l)
			if err != nil {
				return err
			}

			statics = append(statics, addend...)
		}
	} else {
		addend, err := d.staticDirsFor(cfg)
		if err != nil {
			return err
		}

		statics = append(statics, addend...)
	}

	d.staticDirs = removeDuplicatesKeepRight(statics)
	d.AbsStaticDirs = make([]string, len(d.staticDirs))
	for i, di := range d.staticDirs {
		d.AbsStaticDirs[i] = d.pathSpec.AbsPathify(di) + helpers.FilePathSeparator
	}

	d.publishDir = d.pathSpec.AbsPathify(cfg.GetString("publishDir")) + helpers.FilePathSeparator

	return nil
}

func (d *Dirs) staticDirsFor(cfg config.Provider) ([]string, error) {
	var statics []string
	ps, err := helpers.NewPathSpec(d.pathSpec.Fs, cfg)
	if err != nil {
		return statics, err
	}

	statics = append(statics, ps.StaticDirs()...)

	return statics, nil
}

// CreateStaticFs will create a union filesystem with the static paths configured.
// Any missing directories will be logged as warnings.
func (d *Dirs) CreateStaticFs() (afero.Fs, error) {
	var (
		source   = d.pathSpec.Fs.Source
		absPaths []string
	)

	for _, staticDir := range d.AbsStaticDirs {
		if _, err := source.Stat(staticDir); os.IsNotExist(err) {
			d.logger.WARN.Printf("Unable to find Static Directory: %s", staticDir)
		} else {
			absPaths = append(absPaths, staticDir)
		}

	}

	if len(absPaths) == 0 {
		return nil, nil
	}

	return d.createOverlayFs(absPaths), nil

}

// IsStatic returns whether the given filename is located in one of the static
// source dirs.
func (d *Dirs) IsStatic(filename string) bool {
	for _, absPath := range d.AbsStaticDirs {
		if strings.HasPrefix(filename, absPath) {
			return true
		}
	}
	return false
}

// MakeStaticPathRelative creates a relative path from the given filename.
// It will return an empty string if the filename is not a member of dirs.
func (d *Dirs) MakeStaticPathRelative(filename string) string {
	for _, currentPath := range d.AbsStaticDirs {
		if strings.HasPrefix(filename, currentPath) {
			return strings.TrimPrefix(filename, currentPath)
		}
	}

	return ""

}

func (d *Dirs) createOverlayFs(absPaths []string) afero.Fs {
	source := d.pathSpec.Fs.Source

	if len(absPaths) == 1 {
		return afero.NewReadOnlyFs(afero.NewBasePathFs(source, absPaths[0]))
	}

	base := afero.NewReadOnlyFs(afero.NewBasePathFs(source, absPaths[0]))
	overlay := d.createOverlayFs(absPaths[1:])

	return afero.NewCopyOnWriteFs(base, overlay)
}

func removeDuplicatesKeepRight(in []string) []string {
	seen := make(map[string]bool)
	var out []string
	for i := len(in) - 1; i >= 0; i-- {
		v := in[i]
		if seen[v] {
			continue
		}
		out = append([]string{v}, out...)
		seen[v] = true
	}

	return out
}
