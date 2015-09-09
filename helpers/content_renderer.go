package helpers

import (
	"bytes"
	"html"

	"github.com/miekg/mmark"
	"github.com/russross/blackfriday"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type LinkResolverFunc func(ref string) (string, error)
type FileResolverFunc func(ref string) (string, error)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
// Enabling Hugo to customise the rendering experience
type HugoHtmlRenderer struct {
	FileResolver FileResolverFunc
	LinkResolver LinkResolverFunc
	blackfriday.Renderer
}

func (renderer *HugoHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if viper.GetBool("PygmentsCodeFences") {
		opts := viper.GetString("PygmentsOptions")
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, opts))
	} else {
		renderer.Renderer.BlockCode(out, text, lang)
	}
}

func (renderer *HugoHtmlRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if renderer.LinkResolver == nil || bytes.HasPrefix(link, []byte("{@{@HUGOSHORTCODE")) {
		// Use the blackfriday built in Link handler
		renderer.Renderer.Link(out, link, title, content)
	} else {
		newLink, err := renderer.LinkResolver(string(link))
		if err != nil {
			newLink = string(link)
			jww.ERROR.Printf("LinkResolver: %s", err)
		}
		renderer.Renderer.Link(out, []byte(newLink), title, content)
	}
}
func (renderer *HugoHtmlRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if renderer.FileResolver == nil || bytes.HasPrefix(link, []byte("{@{@HUGOSHORTCODE")) {
		// Use the blackfriday built in Image handler
		renderer.Renderer.Image(out, link, title, alt)
	} else {
		newLink, err := renderer.FileResolver(string(link))
		if err != nil {
			newLink = string(link)
			jww.ERROR.Printf("FileResolver: %s", err)
		}
		renderer.Renderer.Image(out, []byte(newLink), title, alt)
	}
}

// Wraps a mmark.Renderer, typically a mmark.html
// Enabling Hugo to customise the rendering experience
type HugoMmarkHtmlRenderer struct {
	mmark.Renderer
}

func (renderer *HugoMmarkHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure bool, callouts bool) {
	if viper.GetBool("PygmentsCodeFences") {
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, ""))
	} else {
		renderer.Renderer.BlockCode(out, text, lang, caption, subfigure, callouts)
	}
}
