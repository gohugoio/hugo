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
	"os"
	"strings"
	"sync"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/output"
)

// PageOutput represents one of potentially many output formats of a given
// Page.
type PageOutput struct {
	*Page

	// Pagination
	paginator     *Pager
	paginatorInit sync.Once

	// Page output specific resources
	resources     resource.Resources
	resourcesInit sync.Once

	// Keep this to create URL/path variations, i.e. paginators.
	targetPathDescriptor targetPathDescriptor

	outputFormat output.Format
}

func (p *PageOutput) targetPath(addends ...string) (string, error) {
	tp, err := p.createTargetPath(p.outputFormat, false, addends...)
	if err != nil {
		return "", err
	}
	return tp, nil
}

func newPageOutput(p *Page, createCopy, initContent bool, f output.Format) (*PageOutput, error) {
	// TODO(bep) This is only needed for tests and we should get rid of it.
	if p.targetPathDescriptorPrototype == nil {
		if err := p.initPaths(); err != nil {
			return nil, err
		}
	}

	if createCopy {
		p = p.copy(initContent)
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
func (p *PageOutput) copyWithFormat(f output.Format, initContent bool) (*PageOutput, error) {
	c, err := newPageOutput(p.Page, true, initContent, f)
	if err != nil {
		return nil, err
	}
	c.paginator = p.paginator
	return c, nil
}

func (p *PageOutput) copy() (*PageOutput, error) {
	return p.copyWithFormat(p.outputFormat, false)
}

func (p *PageOutput) layouts(layouts ...string) ([]string, error) {
	if len(layouts) == 0 && p.selfLayout != "" {
		return []string{p.selfLayout}, nil
	}

	layoutDescriptor := p.layoutDescriptor

	if len(layouts) > 0 {
		layoutDescriptor.Layout = layouts[0]
		layoutDescriptor.LayoutOverride = true
	}

	return p.s.layoutHandler.For(
		layoutDescriptor,
		p.outputFormat)
}

func (p *PageOutput) Render(layout ...string) template.HTML {
	l, err := p.layouts(layout...)
	if err != nil {
		p.s.DistinctErrorLog.Printf("in .Render: Failed to resolve layout %q for page %q", layout, p.pathOrTitle())
		return ""
	}

	for _, layout := range l {
		templ, found := p.s.Tmpl.Lookup(layout)
		if !found {
			// This is legacy from when we had only one output format and
			// HTML templates only. Some have references to layouts without suffix.
			// We default to good old HTML.
			templ, found = p.s.Tmpl.Lookup(layout + ".html")
		}
		if templ != nil {
			res, err := executeToString(templ, p)
			if err != nil {
				p.s.DistinctErrorLog.Printf("in .Render: Failed to execute template %q: %s", layout, err)
				return template.HTML("")
			}
			return template.HTML(res)
		}
	}

	return ""

}

func executeToString(templ tpl.Template, data interface{}) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := templ.Execute(b, data); err != nil {
		return "", err
	}
	return b.String(), nil

}

func (p *Page) Render(layout ...string) template.HTML {
	if p.mainPageOutput == nil {
		panic(fmt.Sprintf("programming error: no mainPageOutput for %q", p.Path()))
	}
	return p.mainPageOutput.Render(layout...)
}

// OutputFormats holds a list of the relevant output formats for a given resource.
type OutputFormats []*OutputFormat

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

// AlternativeOutputFormats gives the alternative output formats for this PageOutput.
// Note that we use the term "alternative" and not "alternate" here, as it
// does not necessarily replace the other format, it is an alternative representation.
func (p *PageOutput) AlternativeOutputFormats() (OutputFormats, error) {
	var o OutputFormats
	for _, of := range p.OutputFormats() {
		if of.f.NotAlternative || of.f.Name == p.outputFormat.Name {
			continue
		}
		o = append(o, of)
	}
	return o, nil
}

// deleteResource removes the resource from this PageOutput and the Page. They will
// always be of the same length, but may contain different elements.
func (p *PageOutput) deleteResource(i int) {
	p.resources = append(p.resources[:i], p.resources[i+1:]...)
	p.Page.Resources = append(p.Page.Resources[:i], p.Page.Resources[i+1:]...)

}

func (p *PageOutput) Resources() resource.Resources {
	p.resourcesInit.Do(func() {
		// If the current out shares the same path as the main page output, we reuse
		// the resource set. For the "amp" use case, we need to clone them with new
		// base folder.
		ff := p.outputFormats[0]
		if p.outputFormat.Path == ff.Path {
			p.resources = p.Page.Resources
			return
		}

		// Clone it with new base.
		resources := make(resource.Resources, len(p.Page.Resources))

		for i, r := range p.Page.Resources {
			if c, ok := r.(resource.Cloner); ok {
				// Clone the same resource with a new target.
				resources[i] = c.WithNewBase(p.outputFormat.Path)
			} else {
				resources[i] = r
			}
		}

		p.resources = resources
	})

	return p.resources
}

func (p *PageOutput) renderResources() error {

	for i, r := range p.Resources() {
		src, ok := r.(resource.Source)
		if !ok {
			// Pages gets rendered with the owning page.
			continue
		}

		if err := src.Publish(); err != nil {
			if os.IsNotExist(err) {
				// The resource has been deleted from the file system.
				// This should be extremely rare, but can happen on live reload in server
				// mode when the same resource is member of different page bundles.
				p.deleteResource(i)
			} else {
				p.s.Log.ERROR.Printf("Failed to publish Resource for page %q: %s", p.pathOrTitle(), err)
			}
		} else {
			p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Files)
		}
	}
	return nil
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

// RelPermalink returns the relative permalink to this output format.
func (o *OutputFormat) RelPermalink() string {
	rel := o.p.createRelativePermalinkForOutputFormat(o.f)
	return o.p.s.PathSpec.PrependBasePath(rel, false)
}
