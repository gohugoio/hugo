package helpers

import (
	"bytes"
	"html"

	"github.com/miekg/mmark"
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
// Enabling Hugo to customise the rendering experience
type HugoHtmlRenderer struct {
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
