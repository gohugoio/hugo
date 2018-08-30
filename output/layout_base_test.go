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
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLayoutBase(t *testing.T) {

	var (
		workingDir     = "/sites/mysite/"
		layoutPath1    = "_default/single.html"
		layoutPathAmp  = "_default/single.amp.html"
		layoutPathJSON = "_default/single.json"
	)

	for _, this := range []struct {
		name                 string
		d                    TemplateLookupDescriptor
		needsBase            bool
		basePathMatchStrings string
		expect               TemplateNames
	}{
		{"No base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPath1}, false, "",
			TemplateNames{
				Name:            "_default/single.html",
				OverlayFilename: "_default/single.html",
			}},
		{"Base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPath1}, true, "",
			TemplateNames{
				Name:            "_default/single.html",
				OverlayFilename: "_default/single.html",
				MasterFilename:  "_default/single-baseof.html",
			}},
		// Issue #3893
		{"Base Lang, Default Base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "_default/list.en.html"}, true, "_default/baseof.html",
			TemplateNames{
				Name:            "_default/list.en.html",
				OverlayFilename: "_default/list.en.html",
				MasterFilename:  "_default/baseof.html",
			}},
		{"Base Lang, Lang Base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "_default/list.en.html"}, true, "_default/baseof.html|_default/baseof.en.html",
			TemplateNames{
				Name:            "_default/list.en.html",
				OverlayFilename: "_default/list.en.html",
				MasterFilename:  "_default/baseof.en.html",
			}},
		// Issue #3856
		{"Base Taxonomy Term", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "taxonomy/tag.terms.html"}, true, "_default/baseof.html",
			TemplateNames{
				Name:            "taxonomy/tag.terms.html",
				OverlayFilename: "taxonomy/tag.terms.html",
				MasterFilename:  "_default/baseof.html",
			}},

		{"Partial", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "partials/menu.html"}, true,
			"mytheme/layouts/_default/baseof.html",
			TemplateNames{
				Name:            "partials/menu.html",
				OverlayFilename: "partials/menu.html",
			}},
		{"Partial in subfolder", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "/partials/sub/menu.html"}, true,
			"_default/baseof.html",
			TemplateNames{
				Name:            "partials/sub/menu.html",
				OverlayFilename: "/partials/sub/menu.html",
			}},
		{"Shortcode in subfolder", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: "shortcodes/sub/menu.html"}, true,
			"_default/baseof.html",
			TemplateNames{
				Name:            "shortcodes/sub/menu.html",
				OverlayFilename: "shortcodes/sub/menu.html",
			}},
		{"AMP, no base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPathAmp}, false, "",
			TemplateNames{
				Name:            "_default/single.amp.html",
				OverlayFilename: "_default/single.amp.html",
			}},
		{"JSON, no base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPathJSON}, false, "",
			TemplateNames{
				Name:            "_default/single.json",
				OverlayFilename: "_default/single.json",
			}},
		{"AMP with base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPathAmp}, true, "single-baseof.html|single-baseof.amp.html",
			TemplateNames{
				Name:            "_default/single.amp.html",
				OverlayFilename: "_default/single.amp.html",
				MasterFilename:  "_default/single-baseof.amp.html",
			}},
		{"AMP with no AMP base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPathAmp}, true, "single-baseof.html",
			TemplateNames{
				Name:            "_default/single.amp.html",
				OverlayFilename: "_default/single.amp.html",
				MasterFilename:  "_default/single-baseof.html",
			}},

		{"JSON with base", TemplateLookupDescriptor{WorkingDir: workingDir, RelPath: layoutPathJSON}, true, "single-baseof.json",
			TemplateNames{
				Name:            "_default/single.json",
				OverlayFilename: "_default/single.json",
				MasterFilename:  "_default/single-baseof.json",
			}},
	} {
		t.Run(this.name, func(t *testing.T) {

			this.basePathMatchStrings = filepath.FromSlash(this.basePathMatchStrings)

			fileExists := func(filename string) (bool, error) {
				stringsToMatch := strings.Split(this.basePathMatchStrings, "|")
				for _, s := range stringsToMatch {
					if strings.Contains(filename, s) {
						return true, nil
					}

				}
				return false, nil
			}

			needsBase := func(filename string, subslices [][]byte) (bool, error) {
				return this.needsBase, nil
			}

			this.d.OutputFormats = Formats{AMPFormat, HTMLFormat, RSSFormat, JSONFormat}
			this.d.WorkingDir = filepath.FromSlash(this.d.WorkingDir)
			this.d.RelPath = filepath.FromSlash(this.d.RelPath)
			this.d.ContainsAny = needsBase
			this.d.FileExists = fileExists

			this.expect.MasterFilename = filepath.FromSlash(this.expect.MasterFilename)
			this.expect.OverlayFilename = filepath.FromSlash(this.expect.OverlayFilename)

			if strings.Contains(this.d.RelPath, "json") {
				// currently the only plain text templates in this test.
				this.expect.Name = "_text/" + this.expect.Name
			}

			id, err := CreateTemplateNames(this.d)

			require.NoError(t, err)
			require.Equal(t, this.expect, id, this.name)

		})
	}

}
