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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/paths"

	"errors"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/spf13/afero"
)

const (
	// DefaultArchetypeTemplateTemplate is the template used in 'hugo new site'
	// and the template we use as a fall back.
	DefaultArchetypeTemplateTemplate = `---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true
---

`
)

// NewContent creates a new content file in h (or a full bundle if the archetype is a directory)
// in targetPath.
func NewContent(h *hugolib.HugoSites, kind, targetPath string, force bool) error {
	if h.BaseFs.Content.Dirs == nil {
		return errors.New("no existing content directory configured for this project")
	}

	cf := hugolib.NewContentFactory(h)

	if kind == "" {
		var err error
		kind, err = cf.SectionFromFilename(targetPath)
		if err != nil {
			return err
		}
	}

	b := &contentBuilder{
		archeTypeFs: h.PathSpec.BaseFs.Archetypes.Fs,
		sourceFs:    h.PathSpec.Fs.Source,
		ps:          h.PathSpec,
		h:           h,
		cf:          cf,

		kind:       kind,
		targetPath: targetPath,
		force:      force,
	}

	ext := paths.Ext(targetPath)

	b.setArcheTypeFilenameToUse(ext)

	withBuildLock := func() (string, error) {
		unlock, err := h.BaseFs.LockBuild()
		if err != nil {
			return "", fmt.Errorf("failed to acquire a build lock: %s", err)
		}
		defer unlock()

		if b.isDir {
			return "", b.buildDir()
		}

		if ext == "" {
			return "", fmt.Errorf("failed to resolve %q to an archetype template", targetPath)
		}

		if !files.IsContentFile(b.targetPath) {
			return "", fmt.Errorf("target path %q is not a known content format", b.targetPath)
		}

		return b.buildFile()

	}

	filename, err := withBuildLock()
	if err != nil {
		return err
	}

	if filename != "" {
		return b.openInEditorIfConfigured(filename)
	}

	return nil

}

type contentBuilder struct {
	archeTypeFs afero.Fs
	sourceFs    afero.Fs

	ps *helpers.PathSpec
	h  *hugolib.HugoSites
	cf hugolib.ContentFactory

	// Builder state
	archetypeFilename string
	targetPath        string
	kind              string
	isDir             bool
	dirMap            archetypeMap
	force             bool
}

func (b *contentBuilder) buildDir() error {
	// Split the dir into content files and the rest.
	if err := b.mapArcheTypeDir(); err != nil {
		return err
	}

	var contentTargetFilenames []string
	var baseDir string

	for _, fi := range b.dirMap.contentFiles {
		targetFilename := filepath.Join(b.targetPath, strings.TrimPrefix(fi.Meta().Path, b.archetypeFilename))
		abs, err := b.cf.CreateContentPlaceHolder(targetFilename, b.force)
		if err != nil {
			return err
		}
		if baseDir == "" {
			baseDir = strings.TrimSuffix(abs, targetFilename)
		}

		contentTargetFilenames = append(contentTargetFilenames, abs)
	}

	var contentInclusionFilter *glob.FilenameFilter
	if !b.dirMap.siteUsed {
		// We don't need to build everything.
		contentInclusionFilter = glob.NewFilenameFilterForInclusionFunc(func(filename string) bool {
			filename = strings.TrimPrefix(filename, string(os.PathSeparator))
			for _, cn := range contentTargetFilenames {
				if strings.Contains(cn, filename) {
					return true
				}
			}
			return false
		})

	}

	if err := b.h.Build(hugolib.BuildCfg{NoBuildLock: true, SkipRender: true, ContentInclusionFilter: contentInclusionFilter}); err != nil {
		return err
	}

	for i, filename := range contentTargetFilenames {
		if err := b.applyArcheType(filename, b.dirMap.contentFiles[i].Meta().Path); err != nil {
			return err
		}
	}

	// Copy the rest as is.
	for _, f := range b.dirMap.otherFiles {
		meta := f.Meta()
		filename := meta.Path

		in, err := meta.Open()
		if err != nil {
			return fmt.Errorf("failed to open non-content file: %w", err)
		}

		targetFilename := filepath.Join(baseDir, b.targetPath, strings.TrimPrefix(filename, b.archetypeFilename))
		targetDir := filepath.Dir(targetFilename)

		if err := b.sourceFs.MkdirAll(targetDir, 0o777); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create target directory for %q: %w", targetDir, err)
		}

		out, err := b.sourceFs.Create(targetFilename)
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

	b.h.Log.Printf("Content dir %q created", filepath.Join(baseDir, b.targetPath))

	return nil
}

func (b *contentBuilder) buildFile() (string, error) {
	contentPlaceholderAbsFilename, err := b.cf.CreateContentPlaceHolder(b.targetPath, b.force)
	if err != nil {
		return "", err
	}

	usesSite, err := b.usesSiteVar(b.archetypeFilename)
	if err != nil {
		return "", err
	}

	var contentInclusionFilter *glob.FilenameFilter
	if !usesSite {
		// We don't need to build everything.
		contentInclusionFilter = glob.NewFilenameFilterForInclusionFunc(func(filename string) bool {
			filename = strings.TrimPrefix(filename, string(os.PathSeparator))
			return strings.Contains(contentPlaceholderAbsFilename, filename)
		})
	}

	if err := b.h.Build(hugolib.BuildCfg{NoBuildLock: true, SkipRender: true, ContentInclusionFilter: contentInclusionFilter}); err != nil {
		return "", err
	}

	if err := b.applyArcheType(contentPlaceholderAbsFilename, b.archetypeFilename); err != nil {
		return "", err
	}

	b.h.Log.Printf("Content %q created", contentPlaceholderAbsFilename)

	return contentPlaceholderAbsFilename, nil
}

func (b *contentBuilder) setArcheTypeFilenameToUse(ext string) {
	var pathsToCheck []string

	if b.kind != "" {
		pathsToCheck = append(pathsToCheck, b.kind+ext)
	}

	pathsToCheck = append(pathsToCheck, "default"+ext)

	for _, p := range pathsToCheck {
		fi, err := b.archeTypeFs.Stat(p)
		if err == nil {
			b.archetypeFilename = p
			b.isDir = fi.IsDir()
			return
		}
	}

}

func (b *contentBuilder) applyArcheType(contentFilename, archetypeFilename string) error {
	p := b.h.GetContentPage(contentFilename)
	if p == nil {
		panic(fmt.Sprintf("[BUG] no Page found for %q", contentFilename))
	}

	f, err := b.sourceFs.Create(contentFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	if archetypeFilename == "" {
		return b.cf.ApplyArchetypeTemplate(f, p, b.kind, DefaultArchetypeTemplateTemplate)
	}

	return b.cf.ApplyArchetypeFilename(f, p, b.kind, archetypeFilename)

}

func (b *contentBuilder) mapArcheTypeDir() error {
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
				m.siteUsed, err = b.usesSiteVar(path)
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
		Fs:     b.archeTypeFs,
		Root:   b.archetypeFilename,
	}

	w := hugofs.NewWalkway(walkCfg)

	if err := w.Walk(); err != nil {
		return fmt.Errorf("failed to walk archetype dir %q: %w", b.archetypeFilename, err)
	}

	b.dirMap = m

	return nil
}

func (b *contentBuilder) openInEditorIfConfigured(filename string) error {
	editor := b.h.Cfg.GetString("newContentEditor")
	if editor == "" {
		return nil
	}

	editorExec := strings.Fields(editor)[0]
	editorFlags := strings.Fields(editor)[1:]

	var args []any
	for _, editorFlag := range editorFlags {
		args = append(args, editorFlag)
	}
	args = append(
		args,
		filename,
		hexec.WithStdin(os.Stdin),
		hexec.WithStderr(os.Stderr),
		hexec.WithStdout(os.Stdout),
	)

	b.h.Log.Printf("Editing %q with %q ...\n", filename, editorExec)

	cmd, err := b.h.Deps.ExecHelper.New(editorExec, args...)
	if err != nil {
		return err
	}

	return cmd.Run()
}

func (b *contentBuilder) usesSiteVar(filename string) (bool, error) {
	if filename == "" {
		return false, nil
	}
	bb, err := afero.ReadFile(b.archeTypeFs, filename)
	if err != nil {
		return false, fmt.Errorf("failed to open archetype file: %w", err)
	}

	return bytes.Contains(bb, []byte(".Site")) || bytes.Contains(bb, []byte("site.")), nil

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
