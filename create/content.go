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
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	jww "github.com/spf13/jwalterweatherman"
)

// NewContent creates a new content file in the content directory based upon the
// given kind, which is used to lookup an archetype.
func NewContent(s *hugolib.Site, kind, name string) (err error) {
	jww.INFO.Println("attempting to create ", name, "of", kind)

	location := FindArchetype(s, kind)

	var by []byte

	if location != "" {
		by, err = afero.ReadFile(s.Fs.Source, location)
		if err != nil {
			jww.ERROR.Println(err)
		}
	}
	if location == "" || err != nil {
		by = []byte("+++\ndraft = true \n+++\n")
	}

	psr, err := parser.ReadFrom(bytes.NewReader(by))
	if err != nil {
		return err
	}

	metadata, err := createMetadata(psr, name)
	if err != nil {
		jww.ERROR.Printf("Error processing archetype file %s: %s\n", location, err)
		return err
	}

	page, err := s.NewPage(name)
	if err != nil {
		return err
	}

	if err = page.SetSourceMetaData(metadata, parser.FormatToLeadRune(s.Cfg.GetString("metaDataFormat"))); err != nil {
		return
	}

	page.SetSourceContent(psr.Content())

	contentPath := s.PathSpec.AbsPathify(filepath.Join(s.Cfg.GetString("contentDir"), name))

	if err = page.SafeSaveSourceAs(contentPath); err != nil {
		return
	}
	jww.FEEDBACK.Println(contentPath, "created")

	editor := s.Cfg.GetString("newContentEditor")
	if editor != "" {
		jww.FEEDBACK.Printf("Editing %s with %q ...\n", name, editor)

		cmd := exec.Command(editor, contentPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}

	return nil
}

// createMetadata generates Metadata for a new page based upon the metadata
// found in an archetype.
func createMetadata(archetype parser.Page, name string) (map[string]interface{}, error) {
	archMetadata, err := archetype.Metadata()
	if err != nil {
		return nil, err
	}

	metadata, err := cast.ToStringMapE(archMetadata)
	if err != nil {
		return nil, err
	}

	var date time.Time

	for k, v := range metadata {
		if v == "" {
			continue
		}
		lk := strings.ToLower(k)
		switch lk {
		case "date":
			date, err = cast.ToTimeE(v)
			if err != nil {
				return nil, err
			}
		case "title":
			// Use the archetype title as is
			metadata[lk] = v
		}
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	if date.IsZero() {
		date = time.Now()
	}

	if _, ok := metadata["title"]; !ok {
		metadata["title"] = helpers.MakeTitle(helpers.Filename(name))
	}

	metadata["date"] = date.Format(time.RFC3339)

	return metadata, nil
}

// FindArchetype takes a given kind/archetype of content and returns an output
// path for that archetype.  If no archetype is found, an empty string is
// returned.
func FindArchetype(s *hugolib.Site, kind string) (outpath string) {
	search := []string{s.PathSpec.AbsPathify(s.Cfg.GetString("archetypeDir"))}

	if s.Cfg.GetString("theme") != "" {
		themeDir := filepath.Join(s.PathSpec.AbsPathify(s.Cfg.GetString("themesDir")+"/"+s.Cfg.GetString("theme")), "/archetypes/")
		if _, err := s.Fs.Source.Stat(themeDir); os.IsNotExist(err) {
			jww.ERROR.Printf("Unable to find archetypes directory for theme %q at %q", s.Cfg.GetString("theme"), themeDir)
		} else {
			search = append(search, themeDir)
		}
	}

	for _, x := range search {
		// If the new content isn't in a subdirectory, kind == "".
		// Therefore it should be excluded otherwise `is a directory`
		// error will occur. github.com/spf13/hugo/issues/411
		var pathsToCheck []string

		if kind == "" {
			pathsToCheck = []string{"default.md", "default"}
		} else {
			pathsToCheck = []string{kind + ".md", kind, "default.md", "default"}
		}
		for _, p := range pathsToCheck {
			curpath := filepath.Join(x, p)
			jww.DEBUG.Println("checking", curpath, "for archetypes")
			if exists, _ := helpers.Exists(curpath, s.Fs.Source); exists {
				jww.INFO.Println("curpath: " + curpath)
				return curpath
			}
		}
	}

	return ""
}
