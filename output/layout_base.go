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

	"github.com/gohugoio/hugo/helpers"
)

const (
	baseFileBase = "baseof"
)

var (
	aceTemplateInnerMarkers = [][]byte{[]byte("= content")}
	goTemplateInnerMarkers  = [][]byte{[]byte("{{define"), []byte("{{ define"), []byte("{{- define"), []byte("{{-define")}
)

// TemplateNames represents a template naming scheme.
type TemplateNames struct {
	// The name used as key in the template map. Note that this will be
	// prefixed with "_text/" if it should be parsed with text/template.
	Name string

	OverlayFilename string
	MasterFilename  string
}

// TemplateLookupDescriptor describes the template lookup configuration.
type TemplateLookupDescriptor struct {
	// The full path to the site root.
	WorkingDir string

	// The path to the template relative the the base.
	//  I.e. shortcodes/youtube.html
	RelPath string

	// The template name prefix to look for.
	Prefix string

	// All the output formats in play. This is used to decide if text/template or
	// html/template.
	OutputFormats Formats

	FileExists  func(filename string) (bool, error)
	ContainsAny func(filename string, subslices [][]byte) (bool, error)
}

func isShorthCodeOrPartial(name string) bool {
	return strings.HasPrefix(name, "shortcodes/") || strings.HasPrefix(name, "partials/")
}

// CreateTemplateNames returns a TemplateNames object for a given template.
func CreateTemplateNames(d TemplateLookupDescriptor) (TemplateNames, error) {

	name := filepath.ToSlash(d.RelPath)
	name = strings.TrimPrefix(name, "/")

	if d.Prefix != "" {
		name = strings.Trim(d.Prefix, "/") + "/" + name
	}

	var (
		id TemplateNames
	)

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

	id.OverlayFilename = d.RelPath
	id.Name = name

	if isPlainText {
		id.Name = "_text/" + id.Name
	}

	// Ace and Go templates may have both a base and inner template.
	if ext == "amber" || isShorthCodeOrPartial(name) {
		// No base template support
		return id, nil
	}

	pathDir := filepath.Dir(d.RelPath)

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
	needsBase, err := d.ContainsAny(d.RelPath, innerMarkers)
	if err != nil {
		return id, err
	}

	if needsBase {
		currBaseFilename := fmt.Sprintf("%s-%s", filenameNoSuffix, baseFilename)

		// Look for base template in the follwing order:
		//   1. <current-path>/<template-name>-baseof.<outputFormat>(optional).<suffix>, e.g. list-baseof.<outputFormat>(optional).<suffix>.
		//   2. <current-path>/baseof.<outputFormat>(optional).<suffix>
		//   3. _default/<template-name>-baseof.<outputFormat>(optional).<suffix>, e.g. list-baseof.<outputFormat>(optional).<suffix>.
		//   4. _default/baseof.<outputFormat>(optional).<suffix>
		//
		// The filesystem it looks in a a composite of the project and potential theme(s).
		pathsToCheck := createPathsToCheck(pathDir, baseFilename, currBaseFilename)

		// We may have language code and/or "terms" in the template name. We want the most specific,
		// but need to fall back to the baseof.html or baseof.ace if needed.
		// E.g. list-baseof.en.html and list-baseof.terms.en.html
		// See #3893, #3856.
		baseBaseFilename, currBaseBaseFilename := helpers.Filename(baseFilename), helpers.Filename(currBaseFilename)
		p1, p2 := strings.Split(baseBaseFilename, "."), strings.Split(currBaseBaseFilename, ".")
		if len(p1) > 0 && len(p1) == len(p2) {
			for i := len(p1); i > 0; i-- {
				v1, v2 := strings.Join(p1[:i], ".")+"."+ext, strings.Join(p2[:i], ".")+"."+ext
				pathsToCheck = append(pathsToCheck, createPathsToCheck(pathDir, v1, v2)...)

			}
		}

		for _, p := range pathsToCheck {
			if ok, err := d.FileExists(p); err == nil && ok {
				id.MasterFilename = p
				break
			}
		}
	}

	return id, nil

}

func createPathsToCheck(baseTemplatedDir, baseFilename, currBaseFilename string) []string {
	return []string{
		filepath.Join(baseTemplatedDir, currBaseFilename),
		filepath.Join(baseTemplatedDir, baseFilename),
		filepath.Join("_default", currBaseFilename),
		filepath.Join("_default", baseFilename),
	}
}
