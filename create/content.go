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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

// NewContent creates a new content file in the content directory based upon the
// given kind, which is used to lookup an archetype.
func NewContent(
	ps *helpers.PathSpec,
	siteFactory func(filename string, siteUsed bool) (*hugolib.Site, error), kind, targetPath string) error {
	ext := helpers.Ext(targetPath)
	archetypeFs := ps.BaseFs.SourceFilesystems.Archetypes.Fs

	jww.INFO.Printf("attempting to create %q of %q of ext %q", targetPath, kind, ext)

	archetypeFilename, isDir := findArchetype(ps, kind, ext)

	if isDir {
		cm, err := mapArcheTypeDir(ps, archetypeFs, archetypeFilename)
		if err != nil {
			return err
		}
		s, err := siteFactory(targetPath, cm.siteUsed)
		if err != nil {
			return err
		}
		contentPath := resolveContentPath(s, s.Fs.Source, targetPath)
		return newContentFromDir(kind, s, archetypeFs, s.Fs.Source, cm, contentPath)
	}

	// Building the sites can be expensive, so only do it if really needed.
	siteUsed := false

	if archetypeFilename != "" {
		var err error
		siteUsed, err = usesSiteVar(archetypeFs, archetypeFilename)
		if err != nil {
			return err
		}

	}

	s, err := siteFactory(targetPath, siteUsed)
	if err != nil {
		return err
	}

	contentPath := resolveContentPath(s, archetypeFs, targetPath)

	content, err := executeArcheTypeAsTemplate(s, kind, targetPath, archetypeFilename)
	if err != nil {
		return err
	}

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

func newContentFromDir(
	kind string,
	s *hugolib.Site,
	sourceFs, targetFs afero.Fs,
	cm archetypeMap, targetPath string) error {

	for _, filename := range cm.otherFiles {
		// Just copy the file to destination.
		in, err := sourceFs.Open(filename)
		if err != nil {
			return err
		}

		targetFilename := filepath.Join(targetPath, filename)
		targetDir := filepath.Dir(targetFilename)
		if err := targetFs.MkdirAll(targetDir, 0777); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create target directory for %s: %s", targetDir, err)
		}

		out, err := targetFs.Create(targetFilename)

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}

		in.Close()
		out.Close()
	}

	for _, filename := range cm.contentFiles {
		targetFilename := filepath.Join(targetPath, filename)
		content, err := executeArcheTypeAsTemplate(s, kind, targetFilename, filename)
		if err != nil {
			return err
		}

		if err := helpers.SafeWriteToDisk(targetFilename, bytes.NewReader(content), targetFs); err != nil {
			return err
		}
	}

	jww.FEEDBACK.Println(targetPath, "created")

	return nil
}

type archetypeMap struct {
	// These needs to be parsed and executed as Go templates.
	contentFiles []string
	// These are just copied to destination.
	otherFiles []string
	// If the templates needs a fully built site. This can potentially be
	// expensive, so only do when needed.
	siteUsed bool
}

func mapArcheTypeDir(
	ps *helpers.PathSpec,
	fs afero.Fs,
	archetypeDir string) (archetypeMap, error) {

	var m archetypeMap

	walkFn := func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		if hugolib.IsContentFile(filename) {
			m.contentFiles = append(m.contentFiles, filename)
			if !m.siteUsed {
				m.siteUsed, err = usesSiteVar(fs, filename)
				if err != nil {
					return err
				}
			}
			return nil
		}

		m.otherFiles = append(m.otherFiles, filename)

		return nil
	}

	if err := helpers.SymbolicWalk(fs, archetypeDir, walkFn); err != nil {
		return m, err
	}

	return m, nil
}

func usesSiteVar(fs afero.Fs, filename string) (bool, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return false, fmt.Errorf("failed to open archetype file: %s", err)
	}
	defer f.Close()
	return helpers.ReaderContains(f, []byte(".Site")), nil
}

func resolveContentPath(s *hugolib.Site, fs afero.Fs, targetPath string) string {
	// The site may have multiple content dirs, and we currently do not know which contentDir the
	// user wants to create this content in. We should improve on this, but we start by testing if the
	// provided path points to an existing dir. If so, use it as is.
	var contentPath string
	var exists bool
	targetDir := filepath.Dir(targetPath)

	if targetDir != "" && targetDir != "." {
		exists, _ = helpers.Exists(targetDir, fs)
	}

	if exists {
		contentPath = targetPath
	} else {
		contentPath = s.PathSpec.AbsPathify(filepath.Join(s.Cfg.GetString("contentDir"), targetPath))
	}

	return contentPath
}

// FindArchetype takes a given kind/archetype of content and returns the path
// to the archetype in the archetype filesystem, blank if none found.
func findArchetype(ps *helpers.PathSpec, kind, ext string) (outpath string, isDir bool) {
	fs := ps.BaseFs.Archetypes.Fs

	var pathsToCheck []string

	if kind != "" {
		pathsToCheck = append(pathsToCheck, kind+ext)
	}
	pathsToCheck = append(pathsToCheck, "default"+ext)

	for _, p := range pathsToCheck {
		fi, err := fs.Stat(p)
		if err == nil {
			return p, fi.IsDir()
		}
	}

	return "", false
}
