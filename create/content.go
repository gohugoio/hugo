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

// Package create provides functions to create new content.
package create

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	jww "github.com/spf13/jwalterweatherman"
)

// NewContent creates a new content file in the content directory based upon the
// given kind, which is used to lookup an archetype.
func NewContent(
	ps *helpers.PathSpec,
	siteFactory func(filename string, siteUsed bool) (*hugolib.Site, error), kind, targetPath string) error {
	ext := helpers.Ext(targetPath)

	jww.INFO.Printf("attempting to create %q of %q of ext %q", targetPath, kind, ext)

	archetypeFilename := findArchetype(ps, kind, ext)

	// Building the sites can be expensive, so only do it if really needed.
	siteUsed := false

	if archetypeFilename != "" {
		f, err := ps.Fs.Source.Open(archetypeFilename)
		if err != nil {
			return err
		}
		defer f.Close()

		if helpers.ReaderContains(f, []byte(".Site")) {
			siteUsed = true
		}
	}

	s, err := siteFactory(targetPath, siteUsed)
	if err != nil {
		return err
	}

	var content []byte

	content, err = executeArcheTypeAsTemplate(s, kind, targetPath, archetypeFilename)
	if err != nil {
		return err
	}

	contentPath := s.PathSpec.AbsPathify(filepath.Join(s.Cfg.GetString("contentDir"), targetPath))

	if err := helpers.SafeWriteToDisk(contentPath, bytes.NewReader(content), s.Fs.Source); err != nil {
		return err
	}

	jww.FEEDBACK.Println(contentPath, "created")

	editor := s.Cfg.GetString("newContentEditor")
	if editor != "" {
		jww.FEEDBACK.Printf("Editing %s with %q ...\n", targetPath, editor)

		cmd := exec.Command(editor, contentPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}

	return nil
}

// FindArchetype takes a given kind/archetype of content and returns an output
// path for that archetype.  If no archetype is found, an empty string is
// returned.
func findArchetype(ps *helpers.PathSpec, kind, ext string) (outpath string) {
	search := []string{ps.AbsPathify(ps.Cfg.GetString("archetypeDir"))}

	if ps.Cfg.GetString("theme") != "" {
		themeDir := filepath.Join(ps.AbsPathify(ps.Cfg.GetString("themesDir")+"/"+ps.Cfg.GetString("theme")), "/archetypes/")
		if _, err := ps.Fs.Source.Stat(themeDir); os.IsNotExist(err) {
			jww.ERROR.Printf("Unable to find archetypes directory for theme %q at %q", ps.Cfg.GetString("theme"), themeDir)
		} else {
			search = append(search, themeDir)
		}
	}

	for _, x := range search {
		// If the new content isn't in a subdirectory, kind == "".
		// Therefore it should be excluded otherwise `is a directory`
		// error will occur. github.com/gohugoio/hugo/issues/411
		var pathsToCheck = []string{"default"}

		if ext != "" {
			if kind != "" {
				pathsToCheck = append([]string{kind + ext, "default" + ext}, pathsToCheck...)
			} else {
				pathsToCheck = append([]string{"default" + ext}, pathsToCheck...)
			}
		}

		for _, p := range pathsToCheck {
			curpath := filepath.Join(x, p)
			jww.DEBUG.Println("checking", curpath, "for archetypes")
			if exists, _ := helpers.Exists(curpath, ps.Fs.Source); exists {
				jww.INFO.Println("curpath: " + curpath)
				return curpath
			}
		}
	}

	return ""
}
