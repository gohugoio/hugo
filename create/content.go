// Copyright 2019 The Hugo Authors. All rights reserved.
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

	"github.com/pkg/errors"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

// NewContent creates a new content file in the content directory based upon the
// given kind, which is used to lookup an archetype.
func NewContent(
	sites *hugolib.HugoSites, kind, targetPath string) error {
	targetPath = filepath.Clean(targetPath)
	ext := helpers.Ext(targetPath)
	ps := sites.PathSpec
	archetypeFs := ps.BaseFs.SourceFilesystems.Archetypes.Fs
	sourceFs := ps.Fs.Source

	jww.INFO.Printf("attempting to create %q of %q of ext %q", targetPath, kind, ext)

	archetypeFilename, isDir := findArchetype(ps, kind, ext)
	contentPath, s := resolveContentPath(sites, sourceFs, targetPath)

	if isDir {

		langFs, err := hugofs.NewLanguageFs(sites.LanguageSet(), archetypeFs)
		if err != nil {
			return err
		}

		cm, err := mapArcheTypeDir(ps, langFs, archetypeFilename)
		if err != nil {
			return err
		}

		if cm.siteUsed {
			if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
				return err
			}
		}

		name := filepath.Base(targetPath)
		return newContentFromDir(archetypeFilename, sites, sourceFs, cm, name, contentPath)
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

	if siteUsed {
		if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
			return err
		}
	}

	content, err := executeArcheTypeAsTemplate(s, "", kind, targetPath, archetypeFilename)
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

func targetSite(sites *hugolib.HugoSites, fi hugofs.FileMetaInfo) *hugolib.Site {
	for _, s := range sites.Sites {
		if fi.Meta().Lang() == s.Language().Lang {
			return s
		}
	}
	return sites.Sites[0]
}

func newContentFromDir(
	archetypeDir string,
	sites *hugolib.HugoSites,
	targetFs afero.Fs,
	cm archetypeMap, name, targetPath string) error {

	for _, f := range cm.otherFiles {
		meta := f.Meta()
		filename := meta.Path()
		// Just copy the file to destination.
		in, err := meta.Open()
		if err != nil {
			return errors.Wrap(err, "failed to open non-content file")
		}

		targetFilename := filepath.Join(targetPath, strings.TrimPrefix(filename, archetypeDir))

		targetDir := filepath.Dir(targetFilename)
		if err := targetFs.MkdirAll(targetDir, 0777); err != nil && !os.IsExist(err) {
			return errors.Wrapf(err, "failed to create target directory for %s:", targetDir)
		}

		out, err := targetFs.Create(targetFilename)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}

		in.Close()
		out.Close()
	}

	for _, f := range cm.contentFiles {
		filename := f.Meta().Path()
		s := targetSite(sites, f)
		targetFilename := filepath.Join(targetPath, strings.TrimPrefix(filename, archetypeDir))

		content, err := executeArcheTypeAsTemplate(s, name, archetypeDir, targetFilename, filename)
		if err != nil {
			return errors.Wrap(err, "failed to execute archetype template")
		}

		if err := helpers.SafeWriteToDisk(targetFilename, bytes.NewReader(content), targetFs); err != nil {
			return errors.Wrap(err, "failed to save results")
		}
	}

	jww.FEEDBACK.Println(targetPath, "created")

	return nil
}

type archetypeMap struct {
	// These needs to be parsed and executed as Go templates.
	contentFiles []hugofs.FileMetaInfo
	// These are just copied to destination.
	otherFiles []hugofs.FileMetaInfo
	// If the templates needs a fully built site. This can potentially be
	// expensive, so only do when needed.
	siteUsed bool
}

func mapArcheTypeDir(
	ps *helpers.PathSpec,
	fs afero.Fs,
	archetypeDir string) (archetypeMap, error) {

	var m archetypeMap

	walkFn := func(path string, fi hugofs.FileMetaInfo, err error) error {

		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		fil := fi.(hugofs.FileMetaInfo)

		if files.IsContentFile(path) {
			m.contentFiles = append(m.contentFiles, fil)
			if !m.siteUsed {
				m.siteUsed, err = usesSiteVar(fs, path)
				if err != nil {
					return err
				}
			}
			return nil
		}

		m.otherFiles = append(m.otherFiles, fil)

		return nil
	}

	walkCfg := hugofs.WalkwayConfig{
		WalkFn: walkFn,
		Fs:     fs,
		Root:   archetypeDir,
	}

	w := hugofs.NewWalkway(walkCfg)

	if err := w.Walk(); err != nil {
		return m, errors.Wrapf(err, "failed to walk archetype dir %q", archetypeDir)
	}

	return m, nil
}

func usesSiteVar(fs afero.Fs, filename string) (bool, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return false, errors.Wrap(err, "failed to open archetype file")
	}
	defer f.Close()
	return helpers.ReaderContains(f, []byte(".Site")), nil
}

// Resolve the target content path.
func resolveContentPath(sites *hugolib.HugoSites, fs afero.Fs, targetPath string) (string, *hugolib.Site) {
	targetDir := filepath.Dir(targetPath)
	first := sites.Sites[0]

	var (
		s              *hugolib.Site
		siteContentDir string
	)

	// Try the filename: my-post.en.md
	for _, ss := range sites.Sites {
		if strings.Contains(targetPath, "."+ss.Language().Lang+".") {
			s = ss
			break
		}
	}

	var dirLang string

	for _, dir := range sites.BaseFs.Content.Dirs {
		meta := dir.Meta()
		contentDir := meta.Filename()

		if !strings.HasSuffix(contentDir, helpers.FilePathSeparator) {
			contentDir += helpers.FilePathSeparator
		}

		if strings.HasPrefix(targetPath, contentDir) {
			siteContentDir = contentDir
			dirLang = meta.Lang()
			break
		}
	}

	if s == nil && dirLang != "" {
		for _, ss := range sites.Sites {
			if ss.Lang() == dirLang {
				s = ss
				break
			}
		}
	}

	if s == nil {
		s = first
	}

	if targetDir != "" && targetDir != "." {
		exists, _ := helpers.Exists(targetDir, fs)

		if exists {
			return targetPath, s
		}
	}

	if siteContentDir == "" {

	}

	if siteContentDir != "" {
		pp := filepath.Join(siteContentDir, strings.TrimPrefix(targetPath, siteContentDir))
		return s.PathSpec.AbsPathify(pp), s
	} else {
		var contentDir string
		for _, dir := range sites.BaseFs.Content.Dirs {
			contentDir = dir.Meta().Filename()
			if dir.Meta().Lang() == s.Lang() {
				break
			}
		}
		return s.PathSpec.AbsPathify(filepath.Join(contentDir, targetPath)), s
	}

}

// FindArchetype takes a given kind/archetype of content and returns the path
// to the archetype in the archetype filesystem, blank if none found.
func findArchetype(ps *helpers.PathSpec, kind, ext string) (outpath string, isDir bool) {
	fs := ps.BaseFs.Archetypes.Fs

	var pathsToCheck []string

	if kind != "" {
		pathsToCheck = append(pathsToCheck, kind+ext)
	}
	pathsToCheck = append(pathsToCheck, "default"+ext, "default")

	for _, p := range pathsToCheck {
		fi, err := fs.Stat(p)
		if err == nil {
			return p, fi.IsDir()
		}
	}

	return "", false
}
