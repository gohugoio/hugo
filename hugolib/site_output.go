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

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
	"github.com/spf13/cast"
)

func createDefaultOutputFormats(allFormats output.Formats, cfg config.Provider) map[string]output.Formats {
	rssOut, _ := allFormats.GetByName(output.RSSFormat.Name)
	htmlOut, _ := allFormats.GetByName(output.HTMLFormat.Name)
	robotsOut, _ := allFormats.GetByName(output.RobotsTxtFormat.Name)
	sitemapOut, _ := allFormats.GetByName(output.SitemapFormat.Name)

	// TODO(bep) this mumbo jumbo is deprecated and should be removed, but there are tests that
	// depends on this, so that will have to wait.
	rssBase := cfg.GetString("rssURI")
	if rssBase == "" || rssBase == "index.xml" {
		rssBase = rssOut.BaseName
	} else {
		// Remove in Hugo 0.36.
		helpers.Deprecated("Site config", "rssURI", "Set baseName in outputFormats.RSS", true)
		// RSS has now a well defined media type, so strip any suffix provided
		rssBase = strings.TrimSuffix(rssBase, path.Ext(rssBase))
	}

	rssOut.BaseName = rssBase

	return map[string]output.Formats{
		KindPage:         output.Formats{htmlOut},
		KindHome:         output.Formats{htmlOut, rssOut},
		KindSection:      output.Formats{htmlOut, rssOut},
		KindTaxonomy:     output.Formats{htmlOut, rssOut},
		KindTaxonomyTerm: output.Formats{htmlOut, rssOut},
		// Below are for conistency. They are currently not used during rendering.
		kindRSS:       output.Formats{rssOut},
		kindSitemap:   output.Formats{sitemapOut},
		kindRobotsTXT: output.Formats{robotsOut},
		kind404:       output.Formats{htmlOut},
	}

}

func createSiteOutputFormats(allFormats output.Formats, cfg config.Provider) (map[string]output.Formats, error) {
	defaultOutputFormats := createDefaultOutputFormats(allFormats, cfg)

	if !cfg.IsSet("outputs") {
		return defaultOutputFormats, nil
	}

	outFormats := make(map[string]output.Formats)

	outputs := cfg.GetStringMap("outputs")

	if outputs == nil || len(outputs) == 0 {
		return outFormats, nil
	}

	seen := make(map[string]bool)

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

		// This effectively prevents empty outputs entries for a given Kind.
		// We need at least one.
		if len(formats) > 0 {
			seen[k] = true
			outFormats[k] = formats
		}
	}

	// Add defaults for the entries not provided by the user.
	for k, v := range defaultOutputFormats {
		if !seen[k] {
			outFormats[k] = v
		}
	}

	return outFormats, nil

}
