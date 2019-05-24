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

// Package page contains the core interfaces and types for the Page resource,
// a core component in Hugo.
package page

import (
	"strings"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
)

// OutputFormats holds a list of the relevant output formats for a given page.
type OutputFormats []OutputFormat

// OutputFormat links to a representation of a resource.
type OutputFormat struct {
	// Rel constains a value that can be used to construct a rel link.
	// This is value is fetched from the output format definition.
	// Note that for pages with only one output format,
	// this method will always return "canonical".
	// As an example, the AMP output format will, by default, return "amphtml".
	//
	// See:
	// https://www.ampproject.org/docs/guides/deploy/discovery
	//
	// Most other output formats will have "alternate" as value for this.
	Rel string

	Format output.Format

	relPermalink string
	permalink    string
}

// Name returns this OutputFormat's name, i.e. HTML, AMP, JSON etc.
func (o OutputFormat) Name() string {
	return o.Format.Name
}

// MediaType returns this OutputFormat's MediaType (MIME type).
func (o OutputFormat) MediaType() media.Type {
	return o.Format.MediaType
}

// Permalink returns the absolute permalink to this output format.
func (o OutputFormat) Permalink() string {
	return o.permalink
}

// RelPermalink returns the relative permalink to this output format.
func (o OutputFormat) RelPermalink() string {
	return o.relPermalink
}

func NewOutputFormat(relPermalink, permalink string, isCanonical bool, f output.Format) OutputFormat {
	rel := f.Rel
	if isCanonical {
		rel = "canonical"
	}
	return OutputFormat{Rel: rel, Format: f, relPermalink: relPermalink, permalink: permalink}
}

// Get gets a OutputFormat given its name, i.e. json, html etc.
// It returns nil if none found.
func (o OutputFormats) Get(name string) *OutputFormat {
	for _, f := range o {
		if strings.EqualFold(f.Format.Name, name) {
			return &f
		}
	}
	return nil
}
