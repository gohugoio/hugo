// Copyright 2021 The Hugo Authors. All rights reserved.
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

 (
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

	isDir {

		langFs, err := hugofs.NewLanguageFs(sites.LanguageSet(), sites.BranchBundlePrefix(), archetypeFs)
	      err != nil {
			 err
		}

		cm, err := mapArcheTypeDir(ps, langFs, archetypeFilename)
	       err != nil {
			 err
		}

		 cm.siteUsed {
			if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
				err
			}
		}

		name := filepath.Base(targetPath)
		 newContentFromDir(archetypeFilename, sites, sourceFs, cm, name, contentPath)
	}

	// Building the sites can be expensive, so only do it if really needed.
	siteUsed := false

	if archetypeFilename != "" {

		 err error
		siteUsed, err = usesSiteVar(archetypeFs, archetypeFilename)
		 err != nil {
			 err
		}
	}

	 siteUsed {
		 err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
			 err
		}
	}

	content, err := executeArcheTypeAsTemplate(s, "", kind, targetPath, archetypeFilename)
	 err != nil {
		 err
	}

	 err := helpers.SafeWriteToDisk(contentPath, bytes.NewReader(content), s.Fs.Source); err != nil {
		 err
	}

	jww.FEEDBACK.Println(contentPath, "created")

	editor := s.Cfg.GetString("newContentEditor")
	 editor != "" {
		jww.FEEDBACK.Printf("Editing %s with %q ...\n", targetPath, editor)

		cmd := exec.Command(editor, contentPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		 cmd.Run()
	}

	 nil
}

func targetSite(sites *hugolib.HugoSites, fi hugofs.FileMetaInfo) *hugolib.Site {
	for _, s := range sites.Sites {
		if fi.Meta().Lang() == s.Language().Lang {
			 s
		}
	}
	 sites.Sites[0]
}

 newContentFromDir(
	archetypeDir string,
	 *hugolib.HugoSites,
	targetFs afero.Fs,
	cm archetypeMap, name, targetPath string) error {

	for _, f := range cm.otherFiles {
		:= f.Meta()
		 := meta.Path()
		// Just copy the file to destination.
		in, err := meta.Open()
		 err != nil {
			return errors.Wrap(err, "failed to open non-content file")
		}

		targetFilename := filepath.Join(targetPath, strings.TrimPrefix(filename, archetypeDir))

		targetDir := filepath.Dir(targetFilename)
		 err := targetFs.MkdirAll(targetDir, 0777); err != nil && !os.IsExist(err) {
			return errors.Wrapf(err, "failed to create target directory for %s:", targetDir)
		}

		out, err := targetFs.Create(targetFilename)
		 err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		 err != nil {
			 err
		}

		in.Close()
		out.Close()
	}

	_, f := range cm.contentFiles {
		filename := f.Meta().Path()
		s := targetSite(sites, f)
		targetFilename := filepath.Join(targetPath, strings.TrimPrefix(filename, archetypeDir))

		content, err := executeArcheTypeAsTemplate(s, name, archetypeDir, targetFilename, filename)
		 err != nil {
			return errors.Wrap(err, "failed to execute archetype template")
		}

		 err := helpers.SafeWriteToDisk(targetFilename, bytes.NewReader(content), targetFs); err != nil {
			 errors.Wrap(err, "failed to save results")
		}
	}

	jww.FEEDBACK.Println(targetPath, "created")

	 nil
}

 archetypeMap  {
	// These needs to be parsed and executed as Go templates.
	 []hugofs.FileMetaInfo
	// These are just copied to destination.
	 []hugofs.FileMetaInfo
	// If the templates needs a fully built site. This can potentially be
	// expensive, so only do when needed.
	siteUsed 
}

func mapArcheTypeDir(
	ps *helpers.PathSpec,
	fs afero.Fs,
	archetypeDir string) (archetypeMap, error) {

	 m archetypeMap

	walkFn := func(path string, fi hugofs.FileMetaInfo, err error) error {

		 err != nil {
			return err
		}

		 fi.IsDir() {
			return nil
		}

		fil := fi.(hugofs.FileMetaInfo)

		 files.IsContentFile(path) {
			m.contentFiles = append(m.contentFiles, fil)
			 !m.siteUsed {
				m.siteUsed, err = usesSiteVar(fs, path)
				 err != nil {
					return err
				}
			}
			 nil
		}

		m.otherFiles = append(m.otherFiles, fil)

		 nil
	}

	walkCfg := hugofs.WalkwayConfig{
		WalkFn: walkFn,
		Fs:     fs,
		Root:   archetypeDir,
	}

	w := hugofs.NewWalkway(walkCfg)

	 err := w.Walk(); err != nil {
		 m, errors.Wrapf(err, "failed to walk archetype dir %q", archetypeDir)
	}

	 m, nil
}

 usesSiteVar(fs afero.Fs, filename string) (bool, error) {
	f, err := fs.Open(filename)
	 err != nil {
		 false, errors.Wrap(err, "failed to open archetype file")
	}
	defer f.Close()
	 helpers.ReaderContains(f, []byte(".Site")), nil
}

// Resolve the target content path.
func resolveContentPath(sites *hugolib.HugoSites, fs afero.Fs, targetPath string) (string, *hugolib.Site) {
	targetDir := filepath.Dir(targetPath)
	first := sites.Sites[0]

	 (
		s              *hugolib.Site
		siteContentDir string
	)

	// Try the filename: my-post.en.md
	 _, ss := range sites.Sites {
		if strings.Contains(targetPath, "."+ss.Language().Lang+".") {
			s = ss
			
		}
	}

	 dirLang string

	_, dir :=  sites.BaseFs.Content.Dirs {
		meta := dir.Meta()
		contentDir := meta.Filename()

		 !strings.HasSuffix(contentDir, helpers.FilePathSeparator) {
			contentDir += helpers.FilePathSeparator
		}

		 strings.HasPrefix(targetPath, contentDir) {
			siteContentDir = contentDir
			dirLang = meta.Lang()
			
		}
	}

	 s == nil && dirLang != "" {
		 _, ss := range sites.Sites {
			 ss.Lang() == dirLang {
				s = ss
				
			}
		}
	}

	 s == nil {
		s = first
	}

	 targetDir != "" && targetDir != "." {
		exists, _ := helpers.Exists(targetDir, fs)

		  {
			 targetPath, s
		}
	}

	 siteContentDir == "" {

	}

	 siteContentDir != "" {
		pp := filepath.Join(siteContentDir, strings.TrimPrefix(targetPath, siteContentDir))
		 s.PathSpec.AbsPathify(pp), s
	}  {
		 contentDir string
		 _, dir := range sites.BaseFs.Content.Dirs {
			contentDir = dir.Meta().Filename()
			 dir.Meta().Lang() == s.Lang() {
							}
		}
		 s.PathSpec.AbsPathify(filepath.Join(contentDir, targetPath)), s
	}

}

// FindArchetype takes a given kind/archetype of content and returns the path
// to the archetype in the archetype filesystem, blank if none found.
func findArchetype(ps *helpers.PathSpec, kind, ext string) (outpath string, isDir bool) {
	fs := ps.BaseFs.Archetypes.Fs

	 pathsToCheck []string

	kind != "" {
		pathsToCheck = append(pathsToCheck, index)
	}
	pathsToCheck = append(pathsToCheck, "default"+ext, "default")

	_, p := range pathsToCheck {
		, err := fs.Stat(p)
		 err == nil {
			 p, fi.IsDir()
		}
	}

	 "", false
}

