// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"fmt"
	"html/template"
	"os/exec"
	"unicode"
	"unicode/utf8"

	"github.com/chaseadamsio/goorgeous"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/config"
	"github.com/miekg/mmark"
	"github.com/mitchellh/mapstructure"
	"github.com/russross/blackfriday"
	jww "github.com/spf13/jwalterweatherman"

	"strings"
)

// SummaryDivider denotes where content summarization should end. The default is "<!--more-->".
var SummaryDivider = []byte("<!--more-->")

// ContentSpec provides functionality to render markdown content.
type ContentSpec struct {
	BlackFriday                *BlackFriday
	footnoteAnchorPrefix       string
	footnoteReturnLinkContents string
	// SummaryLength is the length of the summary that Hugo extracts from a content.
	summaryLength int

	BuildFuture  bool
	BuildExpired bool
	BuildDrafts  bool

	Highlight            func(code, lang, optsStr string) (string, error)
	defatultPygmentsOpts map[string]string

	cfg config.Provider
}

// NewContentSpec returns a ContentSpec initialized
// with the appropriate fields from the given config.Provider.
func NewContentSpec(cfg config.Provider) (*ContentSpec, error) {
	bf := newBlackfriday(cfg.GetStringMap("blackfriday"))
	spec := &ContentSpec{
		BlackFriday:                bf,
		footnoteAnchorPrefix:       cfg.GetString("footnoteAnchorPrefix"),
		footnoteReturnLinkContents: cfg.GetString("footnoteReturnLinkContents"),
		summaryLength:              cfg.GetInt("summaryLength"),
		BuildFuture:                cfg.GetBool("buildFuture"),
		BuildExpired:               cfg.GetBool("buildExpired"),
		BuildDrafts:                cfg.GetBool("buildDrafts"),

		cfg: cfg,
	}

	// Highlighting setup
	options, err := parseDefaultPygmentsOpts(cfg)
	if err != nil {
		return nil, err
	}
	spec.defatultPygmentsOpts = options

	// Use the Pygmentize on path if present
	useClassic := false
	h := newHiglighters(spec)

	if cfg.GetBool("pygmentsUseClassic") {
		if !hasPygments() {
			jww.WARN.Println("Highlighting with pygmentsUseClassic set requires Pygments to be installed and in the path")
		} else {
			useClassic = true
		}
	}

	if useClassic {
		spec.Highlight = h.pygmentsHighlight
	} else {
		spec.Highlight = h.chromaHighlight
	}

	return spec, nil
}

// BlackFriday holds configuration values for BlackFriday rendering.
type BlackFriday struct {
	Smartypants           bool
	SmartypantsQuotesNBSP bool
	AngledQuotes          bool
	Fractions             bool
	HrefTargetBlank       bool
	SmartDashes           bool
	LatexDashes           bool
	TaskLists             bool
	PlainIDAnchors        bool
	Extensions            []string
	ExtensionsMask        []string
}

// NewBlackfriday creates a new Blackfriday filled with site config or some sane defaults.
func newBlackfriday(config map[string]interface{}) *BlackFriday {
	defaultParam := map[string]interface{}{
		"smartypants":           true,
		"angledQuotes":          false,
		"smartypantsQuotesNBSP": false,
		"fractions":             true,
		"hrefTargetBlank":       false,
		"smartDashes":           true,
		"latexDashes":           true,
		"plainIDAnchors":        true,
		"taskLists":             true,
	}

	ToLowerMap(defaultParam)

	siteConfig := make(map[string]interface{})

	for k, v := range defaultParam {
		siteConfig[k] = v
	}

	if config != nil {
		for k, v := range config {
			siteConfig[k] = v
		}
	}

	combinedConfig := &BlackFriday{}
	if err := mapstructure.Decode(siteConfig, combinedConfig); err != nil {
		jww.FATAL.Printf("Failed to get site rendering config\n%s", err.Error())
	}

	return combinedConfig
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
	"backslashLineBreak":     blackfriday.EXTENSION_BACKSLASH_LINE_BREAK,
	"definitionLists":        blackfriday.EXTENSION_DEFINITION_LISTS,
	"joinLines":              blackfriday.EXTENSION_JOIN_LINES,
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
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}

		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}

		wasSpace = isSpace

	}
	return b.String()
}

// stripEmptyNav strips out empty <nav> tags from content.
func stripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

// BytesToHTML converts bytes to type template.HTML.
func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

// getHTMLRenderer creates a new Blackfriday HTML Renderer with the given configuration.
func (c *ContentSpec) getHTMLRenderer(defaultFlags int, ctx *RenderingContext) blackfriday.Renderer {
	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       c.footnoteAnchorPrefix,
		FootnoteReturnLinkContents: c.footnoteReturnLinkContents,
	}

	b := len(ctx.DocumentID) != 0

	if ctx.Config == nil {
		panic(fmt.Sprintf("RenderingContext of %q doesn't have a config", ctx.DocumentID))
	}

	if b && !ctx.Config.PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		renderParameters.HeaderIDSuffix = ":" + ctx.DocumentID
	}

	htmlFlags := defaultFlags
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	if ctx.Config.Smartypants {
		htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	}

	if ctx.Config.SmartypantsQuotesNBSP {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_QUOTES_NBSP
	}

	if ctx.Config.AngledQuotes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	}

	if ctx.Config.Fractions {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	}

	if ctx.Config.HrefTargetBlank {
		htmlFlags |= blackfriday.HTML_HREF_TARGET_BLANK
	}

	if ctx.Config.SmartDashes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_DASHES
	}

	if ctx.Config.LatexDashes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	}

	return &HugoHTMLRenderer{
		cs:               c,
		RenderingContext: ctx,
		Renderer:         blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters),
	}
}

func getMarkdownExtensions(ctx *RenderingContext) int {
	// Default Blackfriday common extensions
	commonExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	// Extra Blackfriday extensions that Hugo enables by default
	flags := commonExtensions |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_FOOTNOTES

	if ctx.Config == nil {
		panic(fmt.Sprintf("RenderingContext of %q doesn't have a config", ctx.DocumentID))
	}

	for _, extension := range ctx.Config.Extensions {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	for _, extension := range ctx.Config.ExtensionsMask {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags &= ^flag
		}
	}
	return flags
}

func (c ContentSpec) markdownRender(ctx *RenderingContext) []byte {
	if ctx.RenderTOC {
		return blackfriday.Markdown(ctx.Content,
			c.getHTMLRenderer(blackfriday.HTML_TOC, ctx),
			getMarkdownExtensions(ctx))
	}
	return blackfriday.Markdown(ctx.Content, c.getHTMLRenderer(0, ctx),
		getMarkdownExtensions(ctx))
}

// getMmarkHTMLRenderer creates a new mmark HTML Renderer with the given configuration.
func (c *ContentSpec) getMmarkHTMLRenderer(defaultFlags int, ctx *RenderingContext) mmark.Renderer {
	renderParameters := mmark.HtmlRendererParameters{
		FootnoteAnchorPrefix:       c.footnoteAnchorPrefix,
		FootnoteReturnLinkContents: c.footnoteReturnLinkContents,
	}

	b := len(ctx.DocumentID) != 0

	if ctx.Config == nil {
		panic(fmt.Sprintf("RenderingContext of %q doesn't have a config", ctx.DocumentID))
	}

	if b && !ctx.Config.PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		// renderParameters.HeaderIDSuffix = ":" + ctx.DocumentId
	}

	htmlFlags := defaultFlags
	htmlFlags |= mmark.HTML_FOOTNOTE_RETURN_LINKS

	return &HugoMmarkHTMLRenderer{
		cs:       c,
		Renderer: mmark.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters),
		Cfg:      c.cfg,
	}
}

func getMmarkExtensions(ctx *RenderingContext) int {
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

	if ctx.Config == nil {
		panic(fmt.Sprintf("RenderingContext of %q doesn't have a config", ctx.DocumentID))
	}

	for _, extension := range ctx.Config.Extensions {
		if flag, ok := mmarkExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	return flags
}

func (c ContentSpec) mmarkRender(ctx *RenderingContext) []byte {
	return mmark.Parse(ctx.Content, c.getMmarkHTMLRenderer(0, ctx),
		getMmarkExtensions(ctx)).Bytes()
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

// RenderingContext holds contextual information, like content and configuration,
// for a given content rendering.
// By creating you must set the Config, otherwise it will panic.
type RenderingContext struct {
	Content      []byte
	PageFmt      string
	DocumentID   string
	DocumentName string
	Config       *BlackFriday
	RenderTOC    bool
	Cfg          config.Provider
}

// RenderBytes renders a []byte.
func (c ContentSpec) RenderBytes(ctx *RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return c.markdownRender(ctx)
	case "markdown":
		return c.markdownRender(ctx)
	case "asciidoc":
		return getAsciidocContent(ctx)
	case "mmark":
		return c.mmarkRender(ctx)
	case "rst":
		return getRstContent(ctx)
	case "org":
		return orgRender(ctx, c)
	case "pandoc":
		return getPandocContent(ctx)
	}
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

// Old implementation only kept for benchmark comparison.
// TODO(bep) remove
func totalWordsOld(s string) int {
	return len(strings.Fields(s))
}

// TruncateWordsByRune truncates words by runes.
func (c *ContentSpec) TruncateWordsByRune(words []string) (string, bool) {
	count := 0
	for index, word := range words {
		if count >= c.summaryLength {
			return strings.Join(words[:index], " "), true
		}
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			count++
		} else if count+runeCount < c.summaryLength {
			count += runeCount
		} else {
			for ri := range word {
				if count >= c.summaryLength {
					truncatedWords := append(words[:index], word[:ri])
					return strings.Join(truncatedWords, " "), true
				}
				count++
			}
		}
	}

	return strings.Join(words, " "), false
}

// TruncateWordsToWholeSentence takes content and truncates to whole sentence
// limited by max number of words. It also returns whether it is truncated.
func (c *ContentSpec) TruncateWordsToWholeSentence(s string) (string, bool) {
	var (
		wordCount     = 0
		lastWordIndex = -1
	)

	for i, r := range s {
		if unicode.IsSpace(r) {
			wordCount++
			lastWordIndex = i

			if wordCount >= c.summaryLength {
				break
			}

		}
	}

	if lastWordIndex == -1 {
		return s, false
	}

	endIndex := -1

	for j, r := range s[lastWordIndex:] {
		if isEndOfSentence(r) {
			endIndex = j + lastWordIndex + utf8.RuneLen(r)
			break
		}
	}

	if endIndex == -1 {
		return s, false
	}

	return strings.TrimSpace(s[:endIndex]), endIndex < len(s)
}

func isEndOfSentence(r rune) bool {
	return r == '.' || r == '?' || r == '!' || r == '"' || r == '\n'
}

// Kept only for benchmark.
func (c *ContentSpec) truncateWordsToWholeSentenceOld(content string) (string, bool) {
	words := strings.Fields(content)

	if c.summaryLength >= len(words) {
		return strings.Join(words, " "), false
	}

	for counter, word := range words[c.summaryLength:] {
		if strings.HasSuffix(word, ".") ||
			strings.HasSuffix(word, "?") ||
			strings.HasSuffix(word, ".\"") ||
			strings.HasSuffix(word, "!") {
			upper := c.summaryLength + counter + 1
			return strings.Join(words[:upper], " "), (upper < len(words))
		}
	}

	return strings.Join(words[:c.summaryLength], " "), true
}

func getAsciidocExecPath() string {
	path, err := exec.LookPath("asciidoc")
	if err != nil {
		return ""
	}
	return path
}

func getAsciidoctorExecPath() string {
	path, err := exec.LookPath("asciidoctor")
	if err != nil {
		return ""
	}
	return path
}

// HasAsciidoc returns whether Asciidoc or Asciidoctor is installed on this computer.
func HasAsciidoc() bool {
	return (getAsciidoctorExecPath() != "" ||
		getAsciidocExecPath() != "")
}

// getAsciidocContent calls asciidoctor or asciidoc as an external helper
// to convert AsciiDoc content to HTML.
func getAsciidocContent(ctx *RenderingContext) []byte {
	var isAsciidoctor bool
	path := getAsciidoctorExecPath()
	if path == "" {
		path = getAsciidocExecPath()
		if path == "" {
			jww.ERROR.Println("asciidoctor / asciidoc not found in $PATH: Please install.\n",
				"                 Leaving AsciiDoc content unrendered.")
			return ctx.Content
		}
	} else {
		isAsciidoctor = true
	}

	jww.INFO.Println("Rendering", ctx.DocumentName, "with", path, "...")
	args := []string{"--no-header-footer", "--safe"}
	if isAsciidoctor {
		// asciidoctor-specific arg to show stack traces on errors
		args = append(args, "--trace")
	}
	args = append(args, "-")
	return externallyRenderContent(ctx, path, args)
}

// HasRst returns whether rst2html is installed on this computer.
func HasRst() bool {
	return getRstExecPath() != ""
}

func getRstExecPath() string {
	path, err := exec.LookPath("rst2html")
	if err != nil {
		path, err = exec.LookPath("rst2html.py")
		if err != nil {
			return ""
		}
	}
	return path
}

func getPythonExecPath() string {
	path, err := exec.LookPath("python")
	if err != nil {
		path, err = exec.LookPath("python.exe")
		if err != nil {
			return ""
		}
	}
	return path
}

// getRstContent calls the Python script rst2html as an external helper
// to convert reStructuredText content to HTML.
func getRstContent(ctx *RenderingContext) []byte {
	python := getPythonExecPath()
	path := getRstExecPath()

	if path == "" {
		jww.ERROR.Println("rst2html / rst2html.py not found in $PATH: Please install.\n",
			"                 Leaving reStructuredText content unrendered.")
		return ctx.Content

	}
	jww.INFO.Println("Rendering", ctx.DocumentName, "with", path, "...")
	args := []string{path, "--leave-comments", "--initial-header-level=2"}
	result := externallyRenderContent(ctx, python, args)
	// TODO(bep) check if rst2html has a body only option.
	bodyStart := bytes.Index(result, []byte("<body>\n"))
	if bodyStart < 0 {
		bodyStart = -7 //compensate for length
	}

	bodyEnd := bytes.Index(result, []byte("\n</body>"))
	if bodyEnd < 0 || bodyEnd >= len(result) {
		bodyEnd = len(result) - 1
		if bodyEnd < 0 {
			bodyEnd = 0
		}
	}

	return result[bodyStart+7 : bodyEnd]
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func getPandocContent(ctx *RenderingContext) []byte {
	path, err := exec.LookPath("pandoc")
	if err != nil {
		jww.ERROR.Println("pandoc not found in $PATH: Please install.\n",
			"                 Leaving pandoc content unrendered.")
		return ctx.Content
	}
	args := []string{"--mathjax"}
	return externallyRenderContent(ctx, path, args)
}

func orgRender(ctx *RenderingContext, c ContentSpec) []byte {
	content := ctx.Content
	cleanContent := bytes.Replace(content, []byte("# more"), []byte(""), 1)
	return goorgeous.Org(cleanContent,
		c.getHTMLRenderer(blackfriday.HTML_TOC, ctx))
}

func externallyRenderContent(ctx *RenderingContext, path string, args []string) []byte {
	content := ctx.Content
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	cmd := exec.Command(path, args...)
	cmd.Stdin = bytes.NewReader(cleanContent)
	var out, cmderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmderr
	err := cmd.Run()
	// Most external helpers exit w/ non-zero exit code only if severe, i.e.
	// halting errors occurred. -> log stderr output regardless of state of err
	for _, item := range strings.Split(string(cmderr.Bytes()), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			jww.ERROR.Printf("%s: %s", ctx.DocumentName, item)
		}
	}
	if err != nil {
		jww.ERROR.Printf("%s rendering %s: %v", path, ctx.DocumentName, err)
	}

	return normalizeExternalHelperLineFeeds(out.Bytes())
}
