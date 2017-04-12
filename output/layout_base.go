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

package output

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/hugo/helpers"
)

const baseFileBase = "baseof"

var (
	aceTemplateInnerMarkers = [][]byte{[]byte("= content")}
	goTemplateInnerMarkers  = [][]byte{[]byte("{{define"), []byte("{{ define")}
)

type TemplateNames struct {
	// The name used as key in the template map. Note that this will be
	// prefixed with "_text/" if it should be parsed with text/template.
	Name string

	OverlayFilename string
	MasterFilename  string
}

type TemplateLookupDescriptor struct {
	// TemplateDir is the project or theme root of the current template.
	// This will be the same as WorkingDir for non-theme templates.
	TemplateDir string

	// The full path to the site root.
	WorkingDir string

	// Main project layout dir, defaults to "layouts"
	LayoutDir string

	// The path to the template relative the the base.
	//  I.e. shortcodes/youtube.html
	RelPath string

	// The template name prefix to look for, i.e. "theme".
	Prefix string

	// The theme dir if theme active.
	ThemeDir string

	// All the output formats in play. This is used to decide if text/template or
	// html/template.
	OutputFormats Formats

	FileExists  func(filename string) (bool, error)
	ContainsAny func(filename string, subslices [][]byte) (bool, error)
}

func CreateTemplateNames(d TemplateLookupDescriptor) (TemplateNames, error) {

	name := filepath.ToSlash(d.RelPath)

	if d.Prefix != "" {
		name = strings.Trim(d.Prefix, "/") + "/" + name
	}

	var (
		id TemplateNames

		// This is the path to the actual template in process. This may
		// be in the theme's or the project's /layouts.
		baseLayoutDir = filepath.Join(d.TemplateDir, d.LayoutDir)
		fullPath      = filepath.Join(baseLayoutDir, d.RelPath)

		// This is always the project's layout dir.
		baseWorkLayoutDir = filepath.Join(d.WorkingDir, d.LayoutDir)

		baseThemeLayoutDir string
	)

	if d.ThemeDir != "" {
		baseThemeLayoutDir = filepath.Join(d.ThemeDir, "layouts")
	}

	// The filename will have a suffix with an optional type indicator.
	// Examples:
	// index.html
	// index.amp.html
	// index.json
	filename := filepath.Base(d.RelPath)
	isPlainText := false
	outputFormat, found := d.OutputFormats.FromFilename(filename)

	if found && outputFormat.IsPlainText {
		isPlainText = true
	}

	var ext, outFormat string

	parts := strings.Split(filename, ".")
	if len(parts) > 2 {
		outFormat = parts[1]
		ext = parts[2]
	} else if len(parts) > 1 {
		ext = parts[1]
	}

	filenameNoSuffix := parts[0]

	id.OverlayFilename = fullPath
	id.Name = name

	if isPlainText {
		id.Name = "_text/" + id.Name
	}

	// Ace and Go templates may have both a base and inner template.
	pathDir := filepath.Dir(fullPath)

	if ext == "amber" || strings.HasSuffix(pathDir, "partials") || strings.HasSuffix(pathDir, "shortcodes") {
		// No base template support
		return id, nil
	}

	innerMarkers := goTemplateInnerMarkers

	var baseFilename string

	if outFormat != "" {
		baseFilename = fmt.Sprintf("%s.%s.%s", baseFileBase, outFormat, ext)
	} else {
		baseFilename = fmt.Sprintf("%s.%s", baseFileBase, ext)
	}

	if ext == "ace" {
		innerMarkers = aceTemplateInnerMarkers
	}

	// This may be a view that shouldn't have base template
	// Have to look inside it to make sure
	needsBase, err := d.ContainsAny(fullPath, innerMarkers)
	if err != nil {
		return id, err
	}

	if needsBase {
		currBaseFilename := fmt.Sprintf("%s-%s", filenameNoSuffix, baseFilename)

		templateDir := filepath.Dir(fullPath)

		// Find the base, e.g. "_default".
		baseTemplatedDir := strings.TrimPrefix(templateDir, baseLayoutDir)
		baseTemplatedDir = strings.TrimPrefix(baseTemplatedDir, helpers.FilePathSeparator)

		// Look for base template in the follwing order:
		//   1. <current-path>/<template-name>-baseof.<outputFormat>(optional).<suffix>, e.g. list-baseof.<outputFormat>(optional).<suffix>.
		//   2. <current-path>/baseof.<outputFormat>(optional).<suffix>
		//   3. _default/<template-name>-baseof.<outputFormat>(optional).<suffix>, e.g. list-baseof.<outputFormat>(optional).<suffix>.
		//   4. _default/baseof.<outputFormat>(optional).<suffix>
		// For each of the steps above, it will first look in the project, then, if theme is set,
		// in the theme's layouts folder.
		// Also note that the <current-path> may be both the project's layout folder and the theme's.
		pairsToCheck := [][]string{
			{baseTemplatedDir, currBaseFilename},
			{baseTemplatedDir, baseFilename},
			{"_default", currBaseFilename},
			{"_default", baseFilename},
		}

	Loop:
		for _, pair := range pairsToCheck {
			pathsToCheck := basePathsToCheck(pair, baseLayoutDir, baseWorkLayoutDir, baseThemeLayoutDir)

			for _, pathToCheck := range pathsToCheck {
				if ok, err := d.FileExists(pathToCheck); err == nil && ok {
					id.MasterFilename = pathToCheck
					break Loop
				}
			}
		}
	}

	return id, nil

}

func basePathsToCheck(path []string, layoutDir, workLayoutDir, themeLayoutDir string) []string {
	// workLayoutDir will always be the most specific, so start there.
	pathsToCheck := []string{filepath.Join((append([]string{workLayoutDir}, path...))...)}

	if layoutDir != "" && layoutDir != workLayoutDir {
		pathsToCheck = append(pathsToCheck, filepath.Join((append([]string{layoutDir}, path...))...))
	}

	// May have a theme
	if themeLayoutDir != "" && themeLayoutDir != layoutDir {
		pathsToCheck = append(pathsToCheck, filepath.Join((append([]string{themeLayoutDir}, path...))...))

	}

	return pathsToCheck

}
