package helpers

import (
	"bytes"
	"html"

	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
type HugoHtmlRenderer struct {
	blackfriday.Renderer
}

func (renderer *HugoHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if viper.GetBool("PygmentsCodeFences") {
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, ""))
	} else {
		renderer.Renderer.BlockCode(out, text, lang)
	}
}
