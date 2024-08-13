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

// Package helpers implements general utility functions that work with
// and on content.  The helper functions defined here lay down the
// foundation of how Hugo works with files and filepaths, and perform
// string operations on content.
package helpers

import (
	"bytes"
	"html/template"
	"strings"
	"unicode"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/media"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/markup"

	"github.com/gohugoio/hugo/config"
)

// ContentSpec provides functionality to render markdown content.
type ContentSpec struct {
	Converters          markup.ConverterProvider
	anchorNameSanitizer converter.AnchorNameSanitizer
	Cfg                 config.AllProvider
}

// NewContentSpec returns a ContentSpec initialized
// with the appropriate fields from the given config.Provider.
func NewContentSpec(cfg config.AllProvider, logger loggers.Logger, contentFs afero.Fs, ex *hexec.Exec) (*ContentSpec, error) {
	spec := &ContentSpec{
		Cfg: cfg,
	}

	converterProvider, err := markup.NewConverterProvider(converter.ProviderConfig{
		Conf:      cfg,
		ContentFs: contentFs,
		Logger:    logger,
		Exec:      ex,
	})
	if err != nil {
		return nil, err
	}

	spec.Converters = converterProvider
	p := converterProvider.Get("markdown")
	conv, err := p.New(converter.DocumentContext{})
	if err != nil {
		return nil, err
	}
	if as, ok := conv.(converter.AnchorNameSanitizer); ok {
		spec.anchorNameSanitizer = as
	} else {
		// Use Goldmark's sanitizer
		p := converterProvider.Get("goldmark")
		conv, err := p.New(converter.DocumentContext{})
		if err != nil {
			return nil, err
		}
		spec.anchorNameSanitizer = conv.(converter.AnchorNameSanitizer)
	}

	return spec, nil
}

// stripEmptyNav strips out empty <nav> tags from content.
func stripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

// BytesToHTML converts bytes to type template.HTML.
func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

// ExtractTOC extracts Table of Contents from content.
func ExtractTOC(content []byte) (newcontent []byte, toc []byte) {
	if !bytes.Contains(content, []byte("<nav>")) {
		return content, nil
	}
	origContent := make([]byte, len(content))
	copy(origContent, content)
	first := []byte(`<nav>
<ul>`)

	last := []byte(`</ul>
</nav>`)

	replacement := []byte(`<nav id="TableOfContents">
<ul>`)

	startOfTOC := bytes.Index(content, first)

	peekEnd := len(content)
	if peekEnd > 70+startOfTOC {
		peekEnd = 70 + startOfTOC
	}

	if startOfTOC < 0 {
		return stripEmptyNav(content), toc
	}
	// Need to peek ahead to see if this nav element is actually the right one.
	correctNav := bytes.Index(content[startOfTOC:peekEnd], []byte(`<li><a href="#`))
	if correctNav < 0 { // no match found
		return content, toc
	}
	lengthOfTOC := bytes.Index(content[startOfTOC:], last) + len(last)
	endOfTOC := startOfTOC + lengthOfTOC

	newcontent = append(content[:startOfTOC], content[endOfTOC:]...)
	toc = append(replacement, origContent[startOfTOC+len(first):endOfTOC]...)
	return
}

func (c *ContentSpec) SanitizeAnchorName(s string) string {
	return c.anchorNameSanitizer.SanitizeAnchorName(s)
}

func (c *ContentSpec) ResolveMarkup(in string) string {
	in = strings.ToLower(in)

	if mediaType, found := c.Cfg.ContentTypes().(media.ContentTypes).Types().GetBestMatch(markup.ResolveMarkup(in)); found {
		return mediaType.SubType
	}

	if conv := c.Converters.Get(in); conv != nil {
		return markup.ResolveMarkup(conv.Name())
	}

	return ""
}

// TotalWords counts instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, in s.
// This is a cheaper way of word counting than the obvious len(strings.Fields(s)).
func TotalWords(s string) int {
	n := 0
	inWord := false
	for _, r := range s {
		wasInWord := inWord
		inWord = !unicode.IsSpace(r)
		if inWord && !wasInWord {
			n++
		}
	}
	return n
}

// TrimShortHTML removes the outer tags from HTML input where (a) the opening
// tag is present only once with the input, and (b) the opening and closing
// tags wrap the input after white space removal.
func (c *ContentSpec) TrimShortHTML(input []byte, markup string) []byte {
	openingTag := []byte("<p>")
	closingTag := []byte("</p>")

	if markup == media.DefaultContentTypes.AsciiDoc.SubType {
		openingTag = []byte("<div class=\"paragraph\">\n<p>")
		closingTag = []byte("</p>\n</div>")
	}

	if bytes.Count(input, openingTag) == 1 {
		input = bytes.TrimSpace(input)
		if bytes.HasPrefix(input, openingTag) && bytes.HasSuffix(input, closingTag) {
			input = bytes.TrimPrefix(input, openingTag)
			input = bytes.TrimSuffix(input, closingTag)
			input = bytes.TrimSpace(input)
		}
	}
	return input
}
