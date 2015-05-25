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

	"github.com/miekg/mmark"
	"github.com/russross/blackfriday"
	bp "github.com/spf13/hugo/bufferpool"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"strings"
	"sync"
)

// Length of the summary that Hugo extracts from a content.
var SummaryLength = 70

// Custom divider <!--more--> let's user define where summarization ends.
var SummaryDivider = []byte("<!--more-->")

// Blackfriday holds configuration values for Blackfriday rendering.
type Blackfriday struct {
	AngledQuotes   bool
	Fractions      bool
	PlainIDAnchors bool
	Extensions     []string
	ExtensionsMask []string
}

// NewBlackfriday creates a new Blackfriday with some sane defaults.
func NewBlackfriday() *Blackfriday {
	return &Blackfriday{
		AngledQuotes:   false,
		Fractions:      true,
		PlainIDAnchors: false,
	}
}

var blackfridayExtensionMap = map[string]int{
	"noIntraEmphasis":        blackfriday.EXTENSION_NO_INTRA_EMPHASIS,
	"tables":                 blackfriday.EXTENSION_TABLES,
	"fencedCode":             blackfriday.EXTENSION_FENCED_CODE,
	"autolink":               blackfriday.EXTENSION_AUTOLINK,
	"strikethrough":          blackfriday.EXTENSION_STRIKETHROUGH,
	"laxHtmlBlocks":          blackfriday.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           blackfriday.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          blackfriday.EXTENSION_HARD_LINE_BREAK,
	"tabSizeEight":           blackfriday.EXTENSION_TAB_SIZE_EIGHT,
	"footnotes":              blackfriday.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              blackfriday.EXTENSION_HEADER_IDS,
	"titleblock":             blackfriday.EXTENSION_TITLEBLOCK,
	"autoHeaderIds":          blackfriday.EXTENSION_AUTO_HEADER_IDS,
}

var stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")

var mmarkExtensionMap = map[string]int{
	"tables":                 mmark.EXTENSION_TABLES,
	"fencedCode":             mmark.EXTENSION_FENCED_CODE,
	"autolink":               mmark.EXTENSION_AUTOLINK,
	"laxHtmlBlocks":          mmark.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           mmark.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          mmark.EXTENSION_HARD_LINE_BREAK,
	"footnotes":              mmark.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              mmark.EXTENSION_HEADER_IDS,
	"autoHeaderIds":          mmark.EXTENSION_AUTO_HEADER_IDS,
}

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)

	// Walk through the string removing all tags
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

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
	return b.String()
}

// StripEmptyNav strips out empty <nav> tags from content.
func StripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

// BytesToHTML converts bytes to type template.HTML.
func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

// GetHtmlRenderer creates a new Renderer with the given configuration.
func GetHTMLRenderer(defaultFlags int, ctx *RenderingContext) blackfriday.Renderer {
	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("FootnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("FootnoteReturnLinkContents"),
	}

	b := len(ctx.DocumentID) != 0

	if b && !ctx.getConfig().PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		renderParameters.HeaderIDSuffix = ":" + ctx.DocumentID
	}

	htmlFlags := defaultFlags
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	if ctx.getConfig().AngledQuotes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	}

	if ctx.getConfig().Fractions {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	}

	return blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters)
}

func getMarkdownExtensions(ctx *RenderingContext) int {
	flags := 0 | blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES | blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK | blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS | blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_HEADER_IDS | blackfriday.EXTENSION_AUTO_HEADER_IDS
	for _, extension := range ctx.getConfig().Extensions {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	for _, extension := range ctx.getConfig().ExtensionsMask {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags &= ^flag
		}
	}
	return flags
}

func markdownRender(ctx *RenderingContext) []byte {
	return blackfriday.Markdown(ctx.Content, GetHTMLRenderer(0, ctx),
		getMarkdownExtensions(ctx))
}

func markdownRenderWithTOC(ctx *RenderingContext) []byte {
	return blackfriday.Markdown(ctx.Content,
		GetHTMLRenderer(blackfriday.HTML_TOC, ctx),
		getMarkdownExtensions(ctx))
}

// mmark
func GetMmarkHtmlRenderer(defaultFlags int, ctx *RenderingContext) mmark.Renderer {
	renderParameters := mmark.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("FootnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("FootnoteReturnLinkContents"),
	}

	b := len(ctx.DocumentID) != 0

	if b && !ctx.getConfig().PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		// renderParameters.HeaderIDSuffix = ":" + ctx.DocumentId
	}

	htmlFlags := defaultFlags
	htmlFlags |= mmark.HTML_FOOTNOTE_RETURN_LINKS

	return mmark.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters)
}

func GetMmarkExtensions(ctx *RenderingContext) int {
	flags := 0
	flags |= mmark.EXTENSION_TABLES
	flags |= mmark.EXTENSION_FENCED_CODE
	flags |= mmark.EXTENSION_AUTOLINK
	flags |= mmark.EXTENSION_SPACE_HEADERS
	flags |= mmark.EXTENSION_CITATION
	flags |= mmark.EXTENSION_TITLEBLOCK_TOML
	flags |= mmark.EXTENSION_HEADER_IDS
	flags |= mmark.EXTENSION_AUTO_HEADER_IDS
	flags |= mmark.EXTENSION_UNIQUE_HEADER_IDS
	flags |= mmark.EXTENSION_FOOTNOTES
	flags |= mmark.EXTENSION_SHORT_REF
	flags |= mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	flags |= mmark.EXTENSION_INCLUDE

	for _, extension := range ctx.getConfig().Extensions {
		if flag, ok := mmarkExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	return flags
}

func MmarkRender(ctx *RenderingContext) []byte {
	return mmark.Parse(ctx.Content, GetMmarkHtmlRenderer(0, ctx),
		GetMmarkExtensions(ctx)).Bytes()
}

func MmarkRenderWithTOC(ctx *RenderingContext) []byte {
	return mmark.Parse(ctx.Content,
		GetMmarkHtmlRenderer(0, ctx),
		GetMmarkExtensions(ctx)).Bytes()
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

// RenderingContext holds contextual information, like content and configuration,
// for a given content renderin.g
type RenderingContext struct {
	Content    []byte
	PageFmt    string
	DocumentID string
	Config     *Blackfriday
	configInit sync.Once
}

func (c *RenderingContext) getConfig() *Blackfriday {
	c.configInit.Do(func() {
		if c.Config == nil {
			c.Config = NewBlackfriday()
		}
	})
	return c.Config
}

// RenderBytesWithTOC renders a []byte with table of contents included.
func RenderBytesWithTOC(ctx *RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return markdownRenderWithTOC(ctx)
	case "markdown":
		return markdownRenderWithTOC(ctx)
	case "asciidoc":
		return []byte(GetAsciidocContent(ctx.Content))
	case "mmark":
		return MmarkRenderWithTOC(ctx)
	case "rst":
		return []byte(GetRstContent(ctx.Content))
	}
}

// RenderBytes renders a []byte.
func RenderBytes(ctx *RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return markdownRender(ctx)
	case "markdown":
		return markdownRender(ctx)
	case "asciidoc":
		return []byte(GetAsciidocContent(ctx.Content))
	case "mmark":
		return MmarkRender(ctx)
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
		m[f]++
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
// and returns entire sentences from content, delimited by the int
// and whether it's truncated or not.
func TruncateWordsToWholeSentence(words []string, max int) (string, bool) {
	if max >= len(words) {
		return strings.Join(words, " "), false
	}

	for counter, word := range words[max:] {
		if strings.HasSuffix(word, ".") ||
			strings.HasSuffix(word, "?") ||
			strings.HasSuffix(word, ".\"") ||
			strings.HasSuffix(word, "!") {
			upper := max + counter + 1
			return strings.Join(words[:upper], " "), (upper < len(words))
		}
	}

	return strings.Join(words[:max], " "), true
}

// GetAsciidocContent calls asciidoctor or asciidoc as an external helper
// to convert AsciiDoc content to HTML.
func GetAsciidocContent(content []byte) string {
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	path, err := exec.LookPath("asciidoctor")
	if err != nil {
		path, err = exec.LookPath("asciidoc")
		if err != nil {
			jww.ERROR.Println("asciidoctor / asciidoc not found in $PATH: Please install.\n",
				"                 Leaving AsciiDoc content unrendered.")
			return (string(content))
		}
	}

	jww.INFO.Println("Rendering with", path, "...")
	cmd := exec.Command(path, "--safe", "-")
	cmd.Stdin = bytes.NewReader(cleanContent)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		jww.ERROR.Println(err)
	}

	asciidocLines := strings.Split(out.String(), "\n")
	for i, line := range asciidocLines {
		if strings.HasPrefix(line, "<body") {
			asciidocLines = (asciidocLines[i+1 : len(asciidocLines)-3])
		}
	}
	return strings.Join(asciidocLines, "\n")
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
			return (string(content))
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
