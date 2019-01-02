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
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	// This should be the only list of valid extensions for content files.
	contentFileExtensions = []string{
		"html", "htm",
		"mdown", "markdown", "md",
		"asciidoc", "adoc", "ad",
		"rest", "rst",
		"mmark",
		"org",
		"pandoc", "pdc"}

	contentFileExtensionsSet map[string]bool
)

func init() {
	contentFileExtensionsSet = make(map[string]bool)
	for _, ext := range contentFileExtensions {
		contentFileExtensionsSet[ext] = true
	}
}

func newHandlerChain(s *Site) contentHandler {
	c := &contentHandlers{s: s}

	contentFlow := c.parsePage(c.processFirstMatch(
		// Handles all files with a content file extension. See above.
		c.handlePageContent(),

		// Every HTML file without front matter will be passed on to this handler.
		c.handleHTMLContent(),
	))

	c.rootHandler = c.processFirstMatch(
		contentFlow,

		// Creates a file resource (image, CSS etc.) if there is a parent
		// page set on the current context.
		c.createResource(),

		// Everything that isn't handled above, will just be copied
		// to destination.
		c.copyFile(),
	)

	return c.rootHandler

}

type contentHandlers struct {
	s           *Site
	rootHandler contentHandler
}

func (c *contentHandlers) processFirstMatch(handlers ...contentHandler) func(ctx *handlerContext) handlerResult {
	return func(ctx *handlerContext) handlerResult {
		for _, h := range handlers {
			res := h(ctx)
			if res.handled || res.err != nil {
				return res
			}
		}
		return handlerResult{err: errors.New("no matching handler found")}
	}
}

type handlerContext struct {
	// These are the pages stored in Site.
	pages chan<- *Page

	doNotAddToSiteCollections bool

	currentPage *Page
	parentPage  *Page

	bundle *bundleDir

	source *fileInfo

	// Relative path to the target.
	target string
}

func (c *handlerContext) ext() string {
	if c.currentPage != nil {
		if c.currentPage.Markup != "" {
			return c.currentPage.Markup
		}
		return c.currentPage.Ext()
	}

	if c.bundle != nil {
		return c.bundle.fi.Ext()
	} else {
		return c.source.Ext()
	}
}

func (c *handlerContext) targetPath() string {
	if c.target != "" {
		return c.target
	}

	return c.source.Filename()
}

func (c *handlerContext) file() *fileInfo {
	if c.bundle != nil {
		return c.bundle.fi
	}

	return c.source
}

// Create a copy with the current context as its parent.
func (c handlerContext) childCtx(fi *fileInfo) *handlerContext {
	if c.currentPage == nil {
		panic("Need a Page to create a child context")
	}

	c.target = strings.TrimPrefix(fi.Path(), c.bundle.fi.Dir())
	c.source = fi

	c.doNotAddToSiteCollections = c.bundle != nil && c.bundle.tp != bundleBranch

	c.bundle = nil

	c.parentPage = c.currentPage
	c.currentPage = nil

	return &c
}

func (c *handlerContext) supports(exts ...string) bool {
	ext := c.ext()
	for _, s := range exts {
		if s == ext {
			return true
		}
	}

	return false
}

func (c *handlerContext) isContentFile() bool {
	return contentFileExtensionsSet[c.ext()]
}

type (
	handlerResult struct {
		err      error
		handled  bool
		resource resource.Resource
	}

	contentHandler func(ctx *handlerContext) handlerResult
)

var (
	notHandled handlerResult
)

func (c *contentHandlers) parsePage(h contentHandler) contentHandler {
	return func(ctx *handlerContext) handlerResult {
		if !ctx.isContentFile() {
			return notHandled
		}

		result := handlerResult{handled: true}
		fi := ctx.file()

		f, err := fi.Open()
		if err != nil {
			return handlerResult{err: fmt.Errorf("(%s) failed to open content file: %s", fi.Filename(), err)}
		}
		defer f.Close()

		p := c.s.newPageFromFile(fi)

		_, err = p.ReadFrom(f)
		if err != nil {
			return handlerResult{err: err}
		}

		if !p.shouldBuild() {
			if !ctx.doNotAddToSiteCollections {
				ctx.pages <- p
			}
			return result
		}

		ctx.currentPage = p

		if ctx.bundle != nil {
			// Add the bundled files
			for _, fi := range ctx.bundle.resources {
				childCtx := ctx.childCtx(fi)
				res := c.rootHandler(childCtx)
				if res.err != nil {
					return res
				}
				if res.resource != nil {
					if pageResource, ok := res.resource.(*Page); ok {
						pageResource.resourcePath = filepath.ToSlash(childCtx.target)
						pageResource.parent = p
					}
					p.Resources = append(p.Resources, res.resource)
				}
			}

			sort.SliceStable(p.Resources, func(i, j int) bool {
				if p.Resources[i].ResourceType() < p.Resources[j].ResourceType() {
					return true
				}

				p1, ok1 := p.Resources[i].(*Page)
				p2, ok2 := p.Resources[j].(*Page)

				if ok1 != ok2 {
					return ok2
				}

				if ok1 {
					return defaultPageSort(p1, p2)
				}

				return p.Resources[i].RelPermalink() < p.Resources[j].RelPermalink()
			})

			// Assign metadata from front matter if set
			if len(p.resourcesMetadata) > 0 {
				resources.AssignMetadata(p.resourcesMetadata, p.Resources...)
			}

		}

		return h(ctx)
	}
}

func (c *contentHandlers) handlePageContent() contentHandler {
	return func(ctx *handlerContext) handlerResult {
		if ctx.supports("html", "htm") {
			return notHandled
		}

		p := ctx.currentPage

		p.workContent = p.renderContent(p.workContent)

		tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.workContent)
		p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)
		p.workContent = tmpContent

		if !ctx.doNotAddToSiteCollections {
			ctx.pages <- p
		}

		return handlerResult{handled: true, resource: p}
	}
}

func (c *contentHandlers) handleHTMLContent() contentHandler {
	return func(ctx *handlerContext) handlerResult {
		if !ctx.supports("html", "htm") {
			return notHandled
		}

		p := ctx.currentPage

		if !ctx.doNotAddToSiteCollections {
			ctx.pages <- p
		}

		return handlerResult{handled: true, resource: p}
	}
}

func (c *contentHandlers) createResource() contentHandler {
	return func(ctx *handlerContext) handlerResult {
		if ctx.parentPage == nil {
			return notHandled
		}

		resource, err := c.s.ResourceSpec.New(
			resources.ResourceSourceDescriptor{
				TargetPathBuilder: ctx.parentPage.subResourceTargetPathFactory,
				SourceFile:        ctx.source,
				RelTargetFilename: ctx.target,
				URLBase:           c.s.GetURLLanguageBasePath(),
				TargetBasePaths:   []string{c.s.GetTargetLanguageBasePath()},
			})

		return handlerResult{err: err, handled: true, resource: resource}
	}
}

func (c *contentHandlers) copyFile() contentHandler {
	return func(ctx *handlerContext) handlerResult {
		f, err := c.s.BaseFs.Content.Fs.Open(ctx.source.Filename())
		if err != nil {
			err := fmt.Errorf("failed to open file in copyFile: %s", err)
			return handlerResult{err: err}
		}

		target := ctx.targetPath()

		defer f.Close()
		if err := c.s.publish(&c.s.PathSpec.ProcessingStats.Files, target, f); err != nil {
			return handlerResult{err: err}
		}

		return handlerResult{handled: true}
	}
}
