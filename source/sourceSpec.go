// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/gohugoio/hugo/langs"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
)

// SourceSpec abstracts language-specific file creation.
// TODO(bep) rename to Spec
type SourceSpec struct {
	*helpers.PathSpec

	SourceFs afero.Fs

	// This is set if the ignoreFiles config is set.
	ignoreFilesRe []*regexp.Regexp

	Languages              map[string]interface{}
	DefaultContentLanguage string
	DisabledLanguages      map[string]bool
}

// NewSourceSpec initializes SourceSpec using languages the given filesystem and PathSpec.
func NewSourceSpec(ps *helpers.PathSpec, fs afero.Fs) *SourceSpec {
	cfg := ps.Cfg
	defaultLang := cfg.GetString("defaultContentLanguage")
	languages := cfg.GetStringMap("languages")

	disabledLangsSet := make(map[string]bool)

	for _, disabledLang := range cfg.GetStringSlice("disableLanguages") {
		disabledLangsSet[disabledLang] = true
	}

	if len(languages) == 0 {
		l := langs.NewDefaultLanguage(cfg)
		languages[l.Lang] = l
		defaultLang = l.Lang
	}

	ignoreFiles := cast.ToStringSlice(cfg.Get("ignoreFiles"))
	var regexps []*regexp.Regexp
	if len(ignoreFiles) > 0 {
		for _, ignorePattern := range ignoreFiles {
			re, err := regexp.Compile(ignorePattern)
			if err != nil {
				helpers.DistinctErrorLog.Printf("Invalid regexp %q in ignoreFiles: %s", ignorePattern, err)
			} else {
				regexps = append(regexps, re)
			}

		}
	}

	return &SourceSpec{ignoreFilesRe: regexps, PathSpec: ps, SourceFs: fs, Languages: languages, DefaultContentLanguage: defaultLang, DisabledLanguages: disabledLangsSet}

}

// IgnoreFile returns whether a given file should be ignored.
func (s *SourceSpec) IgnoreFile(filename string) bool {
	if filename == "" {
		if _, ok := s.SourceFs.(*afero.OsFs); ok {
			return true
		}
		return false
	}

	base := filepath.Base(filename)

	if len(base) > 0 {
		first := base[0]
		last := base[len(base)-1]
		if first == '.' ||
			first == '#' ||
			last == '~' {
			return true
		}
	}

	if len(s.ignoreFilesRe) == 0 {
		return false
	}

	for _, re := range s.ignoreFilesRe {
		if re.MatchString(filename) {
			return true
		}
	}

	if runtime.GOOS == "windows" {
		// Also check the forward slash variant if different.
		unixFilename := filepath.ToSlash(filename)
		if unixFilename != filename {
			for _, re := range s.ignoreFilesRe {
				if re.MatchString(unixFilename) {
					return true
				}
			}
		}
	}

	return false
}

// IsRegularSourceFile returns whether filename represents a regular file in the
// source filesystem.
func (s *SourceSpec) IsRegularSourceFile(filename string) (bool, error) {
	fi, err := helpers.LstatIfPossible(s.SourceFs, filename)
	if err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, nil
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err := filepath.EvalSymlinks(filename)
		if err != nil {
			return false, err
		}

		fi, err = helpers.LstatIfPossible(s.SourceFs, link)
		if err != nil {
			return false, err
		}

		if fi.IsDir() {
			return false, nil
		}
	}

	return true, nil
}
