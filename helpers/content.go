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

var SummaryLength = 70
var SummaryDivider = []byte("<!--more-->")

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

func StripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

func GetHtmlRenderer(defaultFlags int, footnoteref string) blackfriday.Renderer {
	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("FootnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("FootnoteReturnLinkContents"),
	}

	if len(footnoteref) != 0 {
		renderParameters.FootnoteAnchorPrefix = footnoteref + ":" +
			renderParameters.FootnoteAnchorPrefix
	}

	htmlFlags := defaultFlags
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	return blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters)
}

func GetMarkdownExtensions() int {
	return 0 | blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES | blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK | blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS | blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_HEADER_IDS
}

func MarkdownRender(content []byte, footnoteref string) []byte {
	return blackfriday.Markdown(content, GetHtmlRenderer(0, footnoteref),
		GetMarkdownExtensions())
}

func MarkdownRenderWithTOC(content []byte, footnoteref string) []byte {
	return blackfriday.Markdown(content,
		GetHtmlRenderer(blackfriday.HTML_TOC, footnoteref),
		GetMarkdownExtensions())
}

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
	correctNav := bytes.Index(content[startOfTOC:peekEnd], []byte(`#toc_0`))
	if correctNav < 0 { // no match found
		return content, toc
	}
	lengthOfTOC := bytes.Index(content[startOfTOC:], last) + len(last)
	endOfTOC := startOfTOC + lengthOfTOC

	newcontent = append(content[:startOfTOC], content[endOfTOC:]...)
	toc = append(replacement, origContent[startOfTOC+len(first):endOfTOC]...)
	return
}

func RenderBytesWithTOC(content []byte, pagefmt string, footnoteref string) []byte {
	switch pagefmt {
	default:
		return MarkdownRenderWithTOC(content, footnoteref)
	case "markdown":
		return MarkdownRenderWithTOC(content, footnoteref)
	case "rst":
		return []byte(GetRstContent(content))
	}
}

func RenderBytes(content []byte, pagefmt string, footnoteref string) []byte {
	switch pagefmt {
	default:
		return MarkdownRender(content, footnoteref)
	case "markdown":
		return MarkdownRender(content, footnoteref)
	case "rst":
		return []byte(GetRstContent(content))
	}
}

func TotalWords(s string) int {
	return len(strings.Fields(s))
}

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, f := range strings.Fields(s) {
		m[f] += 1
	}

	return m
}

func RemoveSummaryDivider(content []byte) []byte {
	return bytes.Replace(content, SummaryDivider, []byte(""), -1)
}

func TruncateWords(s string, max int) string {
	words := strings.Fields(s)
	if max > len(words) {
		return strings.Join(words, " ")
	}

	return strings.Join(words[:max], " ")
}

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

func GetRstContent(content []byte) string {
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	cmd := exec.Command("rst2html.py", "--leave-comments")
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
