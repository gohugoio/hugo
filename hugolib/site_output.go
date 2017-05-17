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

package hugolib

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/output"
)

func createSiteOutputFormats(allFormats output.Formats, cfg config.Provider) (map[string]output.Formats, error) {
	if !cfg.IsSet("outputs") {
		return createDefaultOutputFormats(allFormats, cfg)
	}

	outFormats := make(map[string]output.Formats)

	outputs := cfg.GetStringMap("outputs")

	if outputs == nil || len(outputs) == 0 {
		return outFormats, nil
	}

	for k, v := range outputs {
		var formats output.Formats
		vals := cast.ToStringSlice(v)
		for _, format := range vals {
			f, found := allFormats.GetByName(format)
			if !found {
				return nil, fmt.Errorf("Failed to resolve output format %q from site config", format)
			}
			formats = append(formats, f)
		}

		if len(formats) > 0 {
			outFormats[k] = formats
		}
	}

	// Make sure every kind has at least one output format
	for _, kind := range allKinds {
		if _, found := outFormats[kind]; !found {
			outFormats[kind] = output.Formats{output.HTMLFormat}
		}
	}

	return outFormats, nil

}

func createDefaultOutputFormats(allFormats output.Formats, cfg config.Provider) (map[string]output.Formats, error) {
	outFormats := make(map[string]output.Formats)
	rssOut, _ := allFormats.GetByName(output.RSSFormat.Name)
	htmlOut, _ := allFormats.GetByName(output.HTMLFormat.Name)

	for _, kind := range allKinds {
		var formats output.Formats
		// All have HTML
		formats = append(formats, htmlOut)

		// All but page have RSS
		if kind != KindPage {

			rssBase := cfg.GetString("rssURI")
			if rssBase == "" || rssBase == "index.xml" {
				rssBase = rssOut.BaseName
			} else {
				// Remove in Hugo 0.22.
				helpers.Deprecated("Site config", "rssURI", "Set baseName in outputFormats.RSS", false)
				// RSS has now a well defined media type, so strip any suffix provided
				rssBase = strings.TrimSuffix(rssBase, path.Ext(rssBase))
			}

			rssOut.BaseName = rssBase
			formats = append(formats, rssOut)

		}

		outFormats[kind] = formats
	}

	return outFormats, nil
}
