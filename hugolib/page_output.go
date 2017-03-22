// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"html/template"
	"strings"
	"sync"

	"github.com/spf13/hugo/output"
)

// PageOutput represents one of potentially many output formats of a given
// Page.
type PageOutput struct {
	*Page

	// Pagination
	paginator     *Pager
	paginatorInit sync.Once

	// Keep this to create URL/path variations, i.e. paginators.
	targetPathDescriptor targetPathDescriptor

	outputFormat output.Format
}

func (p *PageOutput) targetPath(addends ...string) (string, error) {
	tp, err := p.createTargetPath(p.outputFormat, addends...)
	if err != nil {
		return "", err
	}
	return tp, nil

}

func newPageOutput(p *Page, createCopy bool, f output.Format) (*PageOutput, error) {
	// For tests
	// TODO(bep) output get rid of this
	if p.targetPathDescriptorPrototype == nil {
		if err := p.initTargetPathDescriptor(); err != nil {
			return nil, err
		}
		if err := p.initURLs(); err != nil {
			return nil, err
		}
	}

	if createCopy {
		p = p.copy()
	}

	td, err := p.createTargetPathDescriptor(f)

	if err != nil {
		return nil, err
	}

	return &PageOutput{
		Page:                 p,
		outputFormat:         f,
		targetPathDescriptor: td,
	}, nil
}

// copy creates a copy of this PageOutput with the lazy sync.Once vars reset
// so they will be evaluated again, for word count calculations etc.
func (p *PageOutput) copy() *PageOutput {
	c, err := newPageOutput(p.Page, true, p.outputFormat)
	if err != nil {
		panic(err)
	}
	return c
}

func (p *PageOutput) layouts(layouts ...string) []string {
	// TODO(bep) output the logic here needs to be redone.
	if len(layouts) == 0 && len(p.layoutsCalculated) > 0 {
		return p.layoutsCalculated
	}

	layoutOverride := ""
	if len(layouts) > 0 {
		layoutOverride = layouts[0]
	}

	return p.s.layoutHandler.For(
		p.layoutDescriptor,
		layoutOverride,
		p.outputFormat)
}

func (p *PageOutput) Render(layout ...string) template.HTML {
	l := p.layouts(layout...)
	return p.s.Tmpl.ExecuteTemplateToHTML(p, l...)
}

// TODO(bep) output
func (p *Page) Render(layout ...string) template.HTML {
	return p.mainPageOutput.Render(layout...)
}

// OutputFormats holds a list of the relevant output formats for a given resource.
type OutputFormats []*OutputFormat

// And OutputFormat links to a representation of a resource.
type OutputFormat struct {
	// Rel constains a value that can be used to construct a rel link.
	// This is value is fetched from the output format definition.
	// Note that for pages with only one output format,
	// this method will always return "canonical".
	// TODO(bep) output -- the above may not be correct for CSS etc. Figure out a way around that.
	// TODO(bep) output -- re the above, maybe add a "alternate" filter to AlternativeOutputFormats.
	// As an example, the AMP output format will, by default, return "amphtml".
	//
	// See:
	// https://www.ampproject.org/docs/guides/deploy/discovery
	//
	// Most other output formats will have "alternate" as value for this.
	Rel string

	// It may be tempting to export this, but let us hold on to that horse for a while.
	f output.Format
	p *Page
}

// Name returns this OutputFormat's name, i.e. HTML, AMP, JSON etc.
func (o OutputFormat) Name() string {
	return o.f.Name
}

// TODO(bep) outputs consider just save this wrapper on Page.
// OutputFormats gives the output formats for this Page.
func (p *Page) OutputFormats() OutputFormats {
	var o OutputFormats
	isCanonical := len(p.outputFormats) == 1
	for _, f := range p.outputFormats {
		rel := f.Rel
		if isCanonical {
			rel = "canonical"
		}
		o = append(o, &OutputFormat{Rel: rel, f: f, p: p})
	}
	return o
}

// OutputFormats gives the alternative output formats for this PageOutput.
func (p *PageOutput) AlternativeOutputFormats() (OutputFormats, error) {
	var o OutputFormats
	for _, of := range p.OutputFormats() {
		if of.f == p.outputFormat {
			continue
		}
		o = append(o, of)
	}
	return o, nil
}

// AlternativeOutputFormats is only available on the top level rendering
// entry point, and not inside range loops on the Page collections.
// This method is just here to inform users of that restriction.
func (p *Page) AlternativeOutputFormats() (OutputFormats, error) {
	return nil, fmt.Errorf("AlternativeOutputFormats only available from the top level template context for page %q", p.Path())
}

// Get gets a OutputFormat given its name, i.e. json, html etc.
// It returns nil if not found.
func (o OutputFormats) Get(name string) *OutputFormat {
	name = strings.ToLower(name)
	for _, f := range o {
		if strings.ToLower(f.f.Name) == name {
			return f
		}
	}
	return nil
}

// Permalink returns the absolute permalink to this output format.
func (o *OutputFormat) Permalink() string {
	rel := o.p.createRelativePermalinkForOutputFormat(o.f)
	return o.p.s.permalink(rel)
}

// Permalink returns the relative permalink to this output format.
func (o *OutputFormat) RelPermalink() string {
	rel := o.p.createRelativePermalinkForOutputFormat(o.f)
	return o.p.s.PathSpec.PrependBasePath(rel)
}
