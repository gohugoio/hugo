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

// Package pandoc converts content to HTML using Pandoc as an external helper.
package pandoc

import (
	"bytes"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/internal"
	"golang.org/x/net/html"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("pandoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &pandocConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil
}

type pandocResult struct {
	converter.Result
	toc tableofcontents.Root
}

func (r pandocResult) TableOfContents() tableofcontents.Root {
	return r.toc
}

type pandocConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func (c *pandocConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	content, toc, err := c.extractTOC(c.getPandocContent(ctx.Src, c.ctx))
	if err != nil {
		return nil, err
	}
	return pandocResult{
		Result: converter.Bytes(content),
		toc:    toc,
	}, nil
}

func (c *pandocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func (c *pandocConverter) getPandocContent(src []byte, ctx converter.DocumentContext) ([]byte, error) {
	logger := c.cfg.Logger
	binaryName := getPandocBinaryName()
	if binaryName == "" {
		logger.Println("pandoc not found in $PATH: Please install.\n",
			"                 Leaving pandoc content unrendered.")
		return src, nil
	}
	args := []string{"--mathjax", "--toc", "--template", "toc", "-s", "--quiet"}
	return internal.ExternallyRenderContent(c.cfg, ctx, src, binaryName, args)
}

const pandocBinary = "pandoc"

func getPandocBinaryName() string {
	if hexec.InPath(pandocBinary) {
		return pandocBinary
	}
	return ""
}

// extractTOC extracts the toc from the given src html.
// It returns the html without the TOC, and the TOC data
func (a *pandocConverter) extractTOC(src []byte) ([]byte, tableofcontents.Root, error) {
	var buf bytes.Buffer
	buf.Write(src)
	node, err := html.Parse(&buf)
	if err != nil {
		return nil, tableofcontents.Root{}, err
	}
	var (
		f       func(*html.Node) bool
		toc     tableofcontents.Root
		toVisit []*html.Node
	)
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "nav" && attr(n, "id") == "TOC" {
			toc = parseTOC(n)
			if !a.cfg.MarkupConfig.Pandoc.PreserveTOC {
				n.Parent.RemoveChild(n)
			}
			return true
		}
		if n.FirstChild != nil {
			toVisit = append(toVisit, n.FirstChild)
		}
		if n.NextSibling != nil && f(n.NextSibling) {
			return true
		}
		for len(toVisit) > 0 {
			nv := toVisit[0]
			toVisit = toVisit[1:]
			if f(nv) {
				return true
			}
		}
		return false
	}
	f(node)
	if err != nil {
		return nil, tableofcontents.Root{}, err
	}
	buf.Reset()
	err = html.Render(&buf, node)
	if err != nil {
		return nil, tableofcontents.Root{}, err
	}
	// ltrim <html><head></head><body> and rtrim </body></html> which are added by html.Render
	res := buf.Bytes()[25:]
	res = res[:len(res)-14]
	return res, toc, nil
}

// parseTOC returns a TOC root from the given toc Node
func parseTOC(doc *html.Node) tableofcontents.Root {
	var (
		toc tableofcontents.Root
		f   func(*html.Node, int, int)
	)
	f = func(n *html.Node, row, level int) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "ul":
				if level == 0 {
					row++
				}
				level++
				f(n.FirstChild, row, level)
			case "li":
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.ElementNode || c.Data != "a" {
						continue
					}
					href := attr(c, "href")[1:]
					toc.AddAt(tableofcontents.Heading{
						Text: nodeContent(c),
						ID:   href,
					}, row, level)
				}
				f(n.FirstChild, row, level)
			}
		}
		if n.NextSibling != nil {
			f(n.NextSibling, row, level)
		}
	}
	f(doc.FirstChild, -1, 0)
	return toc
}

func attr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func nodeContent(node *html.Node) string {
	var buf bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}

// Supports returns whether Pandoc is installed on this computer.
func Supports() bool {
	hasBin := getPandocBinaryName() != ""
	if htesting.SupportsAll() {
		if !hasBin {
			panic("pandoc not installed")
		}
		return true
	}
	return hasBin
}
