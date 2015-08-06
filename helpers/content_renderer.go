package helpers

import (
	"bytes"
	"html"

	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

// hugoHtmlRenderer wraps a blackfriday.Renderer, typically a blackfriday.Html
type hugoHtmlRenderer struct {
	blackfriday.Renderer
}

func (renderer *hugoHtmlRenderer) blockCode(out *bytes.Buffer, text []byte, lang string) {
	if viper.GetBool("PygmentsCodeFences") {
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, ""))
	} else {
		renderer.Renderer.BlockCode(out, text, lang)
	}
}
