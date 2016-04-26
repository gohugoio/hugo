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
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// NewContent creates a new content file in the content directory based upon the
// given kind, which is used to lookup an archetype.
func NewContent(fs afero.Fs, kind, name string) (err error) {
	jww.INFO.Println("attempting to create ", name, "of", kind)

	location := FindArchetype(fs, kind)

	var by []byte

	if location != "" {
		by, err = afero.ReadFile(fs, location)
		if err != nil {
			jww.ERROR.Println(err)
		}
	}
	if location == "" || err != nil {
		by = []byte("+++\n title = \"title\"\n draft = true \n+++\n")
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

	page, err := hugolib.NewPage(name)
	if err != nil {
		return err
	}

	if err = page.SetSourceMetaData(metadata, parser.FormatToLeadRune(viper.GetString("MetaDataFormat"))); err != nil {
		return
	}

	page.SetSourceContent(psr.Content())

	if err = page.SafeSaveSourceAs(filepath.Join(viper.GetString("contentDir"), name)); err != nil {
		return
	}
	jww.FEEDBACK.Println(helpers.AbsPathify(filepath.Join(viper.GetString("contentDir"), name)), "created")

	editor := viper.GetString("NewContentEditor")

	if editor != "" {
		jww.FEEDBACK.Printf("Editing %s with %q ...\n", name, editor)

		cmd := exec.Command(editor, helpers.AbsPathify(path.Join(viper.GetString("contentDir"), name)))
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

	for k := range metadata {
		switch strings.ToLower(k) {
		case "date":
			metadata[k] = time.Now()
		case "title":
			metadata[k] = helpers.MakeTitle(helpers.Filename(name))
		}
	}

	caseimatch := func(m map[string]interface{}, key string) bool {
		for k := range m {
			if strings.ToLower(k) == strings.ToLower(key) {
				return true
			}
		}
		return false
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	if !caseimatch(metadata, "date") {
		metadata["date"] = time.Now()
	}

	if !caseimatch(metadata, "title") {
		metadata["title"] = helpers.MakeTitle(helpers.Filename(name))
	}

	if x := parser.FormatSanitize(viper.GetString("MetaDataFormat")); x == "json" || x == "yaml" || x == "toml" {
		metadata["date"] = time.Now().Format(time.RFC3339)
	}

	return metadata, nil
}

// FindArchetype takes a given kind/archetype of content and returns an output
// path for that archetype.  If no archetype is found, an empty string is
// returned.
func FindArchetype(fs afero.Fs, kind string) (outpath string) {
	search := []string{helpers.AbsPathify(viper.GetString("archetypeDir"))}

	if viper.GetString("theme") != "" {
		themeDir := filepath.Join(helpers.AbsPathify(viper.GetString("themesDir")+"/"+viper.GetString("theme")), "/archetypes/")
		if _, err := fs.Stat(themeDir); os.IsNotExist(err) {
			jww.ERROR.Printf("Unable to find archetypes directory for theme %q at %q", viper.GetString("theme"), themeDir)
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
			if exists, _ := helpers.Exists(curpath, fs); exists {
				jww.INFO.Println("curpath: " + curpath)
				return curpath
			}
		}
	}

	return ""
}
