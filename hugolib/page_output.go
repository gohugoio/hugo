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

	"github.com/spf13/hugo/media"

	"github.com/spf13/hugo/helpers"
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
	// TODO(bep) This is only needed for tests and we should get rid of it.
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
func (p *PageOutput) copyWithFormat(f output.Format) (*PageOutput, error) {
	c, err := newPageOutput(p.Page, true, f)
	if err != nil {
		return nil, err
	}
	c.paginator = p.paginator
	return c, nil
}

func (p *PageOutput) copy() (*PageOutput, error) {
	return p.copyWithFormat(p.outputFormat)
}

func (p *PageOutput) layouts(layouts ...string) ([]string, error) {
	if len(layouts) == 0 && p.selfLayout != "" {
		return []string{p.selfLayout}, nil
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
	if !p.checkRender() {
		return ""
	}

	l, err := p.layouts(layout...)
	if err != nil {
		helpers.DistinctErrorLog.Printf("in .Render: Failed to resolve layout %q for page %q", layout, p.pathOrTitle())
		return ""
	}

	for _, layout := range l {
		templ := p.s.Tmpl.Lookup(layout)
		if templ == nil {
			// This is legacy from when we had only one output format and
			// HTML templates only. Some have references to layouts without suffix.
			// We default to good old HTML.
			templ = p.s.Tmpl.Lookup(layout + ".html")
		}
		if templ != nil {
			res, err := templ.ExecuteToString(p)
			if err != nil {
				helpers.DistinctErrorLog.Printf("in .Render: Failed to execute template %q for page %q", layout, p.pathOrTitle())
				return template.HTML("")
			}
			return template.HTML(res)
		}
	}

	return ""

}

func (p *Page) Render(layout ...string) template.HTML {
	if !p.checkRender() {
		return ""
	}

	p.pageOutputInit.Do(func() {
		if p.mainPageOutput != nil {
			return
		}
		// If Render is called in a range loop, the page output isn't available.
		// So, create one.
		outFormat := p.outputFormats[0]
		pageOutput, err := newPageOutput(p, true, outFormat)

		if err != nil {
			p.s.Log.ERROR.Printf("Failed to create output page for type %q for page %q: %s", outFormat.Name, p.pathOrTitle(), err)
			return
		}

		p.mainPageOutput = pageOutput

	})

	return p.mainPageOutput.Render(layout...)
}

// We may fix this in the future, but the layout handling in Render isn't built
// for list pages.
func (p *Page) checkRender() bool {
	if p.Kind != KindPage {
		helpers.DistinctWarnLog.Printf(".Render only available for regular pages, not for of kind %q. You probably meant .Site.RegularPages and not.Site.Pages.", p.Kind)
		return false
	}
	return true
}

// OutputFormats holds a list of the relevant output formats for a given resource.
type OutputFormats []*OutputFormat

// And OutputFormat links to a representation of a resource.
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

	// It may be tempting to export this, but let us hold on to that horse for a while.
	f output.Format

	p *Page
}

// Name returns this OutputFormat's name, i.e. HTML, AMP, JSON etc.
func (o OutputFormat) Name() string {
	return o.f.Name
}

// MediaType returns this OutputFormat's MediaType (MIME type).
func (o OutputFormat) MediaType() media.Type {
	return o.f.MediaType
}

// OutputFormats gives the output formats for this Page.
func (p *Page) OutputFormats() OutputFormats {
	var o OutputFormats
	for _, f := range p.outputFormats {
		o = append(o, newOutputFormat(p, f))
	}
	return o
}

func newOutputFormat(p *Page, f output.Format) *OutputFormat {
	rel := f.Rel
	isCanonical := len(p.outputFormats) == 1
	if isCanonical {
		rel = "canonical"
	}
	return &OutputFormat{Rel: rel, f: f, p: p}
}

// OutputFormats gives the alternative output formats for this PageOutput.
// Note that we use the term "alternative" and not "alternate" here, as it
// does not necessarily replace the other format, it is an alternative representation.
func (p *PageOutput) AlternativeOutputFormats() (OutputFormats, error) {
	var o OutputFormats
	for _, of := range p.OutputFormats() {
		if of.f.NotAlternative || of.f == p.outputFormat {
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
	for _, f := range o {
		if strings.EqualFold(f.f.Name, name) {
			return f
		}
	}
	return nil
}

// Permalink returns the absolute permalink to this output format.
func (o *OutputFormat) Permalink() string {
	rel := o.p.createRelativePermalinkForOutputFormat(o.f)
	perm, _ := o.p.s.permalinkForOutputFormat(rel, o.f)
	return perm
}

// Permalink returns the relative permalink to this output format.
func (o *OutputFormat) RelPermalink() string {
	rel := o.p.createRelativePermalinkForOutputFormat(o.f)
	return o.p.s.PathSpec.PrependBasePath(rel)
}
