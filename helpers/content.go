// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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
	"os/exec"

	"github.com/russross/blackfriday"
	"github.com/spf13/viper"

	jww "github.com/spf13/jwalterweatherman"

	"strings"
)

// Length of the summary that Hugo extracts from a content.
var SummaryLength = 70

// Custom divider <!--more--> let's user define where summarization ends.
var SummaryDivider = []byte("<!--more-->")

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {
	output := ""

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "</p>", "\n", -1)
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "<br />", "\n", -1) // <br /> is the xhtml line break tag

		// Walk through the string removing all tags
		b := new(bytes.Buffer)
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}
	return output
}

// StripEmptyNav strips out empty <nav> tags from content.
func StripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

// BytesToHTML converts bytes to type template.HTML.
func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

func GetHtmlRenderer(defaultFlags int, ctx RenderingContext) blackfriday.Renderer {
	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("FootnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("FootnoteReturnLinkContents"),
	}

	b := len(ctx.DocumentId) != 0

	if m, ok := ctx.ConfigFlags["plainIdAnchors"]; b && ((ok && !m) || !ok) {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentId + ":" + renderParameters.FootnoteAnchorPrefix
		renderParameters.HeaderIDSuffix = ":" + ctx.DocumentId
	}

	htmlFlags := defaultFlags
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS

	angledQuotes := false
	smartFractions := true
	latexDashes := true
	footnoteReturnLink := true

	if m, ok := ctx.ConfigFlags["angledQuotes"]; ok {
		angledQuotes = m
	}
	if m, ok := ctx.ConfigFlags["smartFractions"]; ok {
		smartFractions = m
	}
	if m, ok := ctx.ConfigFlags["latexDashes"]; ok {
		latexDashes = m
	}
	if m, ok := ctx.ConfigFlags["footnoteReturnLink"]; ok {
		footnoteReturnLink = m
	}

	if angledQuotes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	}

	if smartFractions {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	}

	if latexDashes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	}

	if footnoteReturnLink {
		htmlFlags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS
	}

	return blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters)
}

func GetMarkdownExtensions() int {
	return 0 | blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES | blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK | blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS | blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_HEADER_IDS | blackfriday.EXTENSION_AUTO_HEADER_IDS
}

func MarkdownRender(ctx RenderingContext) []byte {
	return blackfriday.Markdown(ctx.Content, GetHtmlRenderer(0, ctx),
		GetMarkdownExtensions())
}

func MarkdownRenderWithTOC(ctx RenderingContext) []byte {
	return blackfriday.Markdown(ctx.Content,
		GetHtmlRenderer(blackfriday.HTML_TOC, ctx),
		GetMarkdownExtensions())
}

// ExtractTOC extracts Table of Contents from content.
func ExtractTOC(content []byte) (newcontent []byte, toc []byte) {
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
		return StripEmptyNav(content), toc
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

type RenderingContext struct {
	Content     []byte
	PageFmt     string
	DocumentId  string
	ConfigFlags map[string]bool
}

func RenderBytesWithTOC(ctx RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return MarkdownRenderWithTOC(ctx)
	case "markdown":
		return MarkdownRenderWithTOC(ctx)
	case "rst":
		return []byte(GetRstContent(ctx.Content))
	}
}

func RenderBytes(ctx RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return MarkdownRender(ctx)
	case "markdown":
		return MarkdownRender(ctx)
	case "rst":
		return []byte(GetRstContent(ctx.Content))
	}
}

// TotalWords returns an int of the total number of words in a given content.
func TotalWords(s string) int {
	return len(strings.Fields(s))
}

// WordCount takes content and returns a map of words and count of each word.
func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, f := range strings.Fields(s) {
		m[f] += 1
	}

	return m
}

// RemoveSummaryDivider removes summary-divider <!--more--> from content.
func RemoveSummaryDivider(content []byte) []byte {
	return bytes.Replace(content, SummaryDivider, []byte(""), -1)
}

// TruncateWords takes content and an int and shortens down the number
// of words in the content down to the number of int.
func TruncateWords(s string, max int) string {
	words := strings.Fields(s)
	if max > len(words) {
		return strings.Join(words, " ")
	}

	return strings.Join(words[:max], " ")
}

// TruncateWordsToWholeSentence takes content and an int
// and returns entire sentences from content, delimited by the int.
func TruncateWordsToWholeSentence(s string, max int) string {
	words := strings.Fields(s)
	if max > len(words) {
		return strings.Join(words, " ")
	}

	for counter, word := range words[max:] {
		if strings.HasSuffix(word, ".") ||
			strings.HasSuffix(word, "?") ||
			strings.HasSuffix(word, ".\"") ||
			strings.HasSuffix(word, "!") {
			return strings.Join(words[:max+counter+1], " ")
		}
	}

	return strings.Join(words[:max], " ")
}

// GetRstContent calls the Python script rst2html as an external helper
// to convert reStructuredText content to HTML.
func GetRstContent(content []byte) string {
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	path, err := exec.LookPath("rst2html")
	if err != nil {
		path, err = exec.LookPath("rst2html.py")
		if err != nil {
			jww.ERROR.Println("rst2html / rst2html.py not found in $PATH: Please install.\n",
			                  "                 Leaving reStructuredText content unrendered.")
			return(string(content))
		}
	}

	cmd := exec.Command(path, "--leave-comments")
	cmd.Stdin = bytes.NewReader(cleanContent)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		jww.ERROR.Println(err)
	}

	rstLines := strings.Split(out.String(), "\n")
	for i, line := range rstLines {
		if strings.HasPrefix(line, "<body>") {
			rstLines = (rstLines[i+1 : len(rstLines)-3])
		}
	}
	return strings.Join(rstLines, "\n")
}
