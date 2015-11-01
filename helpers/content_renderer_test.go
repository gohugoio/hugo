package helpers

import (
	"bytes"
	"github.com/spf13/viper"
	"testing"
)

// Renders a codeblock using Blackfriday
func render(input string) string {
	ctx := &RenderingContext{}
	render := GetHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html")
	return buf.String()
}

// Renders a codeblock using Mmark
func renderWithMmark(input string) string {
	ctx := &RenderingContext{}
	render := GetMmarkHtmlRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html", []byte(""), false, false)
	return buf.String()
}

func TestCodeFence(t *testing.T) {

	if !HasPygments() {
		t.Skip("Skipping Pygments test as Pygments is not installed or available.")
		return
	}

	type test struct {
		enabled         bool
		input, expected string
	}
	data := []test{
		{true, "<html></html>", "<div class=\"highlight\"><pre><code class=\"language-html\" data-lang=\"html\"><span class=\"nt\">&lt;html&gt;&lt;/html&gt;</span>\n</code></pre></div>\n"},
		{false, "<html></html>", "<pre><code class=\"language-html\">&lt;html&gt;&lt;/html&gt;</code></pre>\n"},
	}

	viper.Reset()
	defer viper.Reset()

	viper.Set("PygmentsStyle", "monokai")
	viper.Set("PygmentsUseClasses", true)

	for i, d := range data {
		viper.Set("PygmentsCodeFences", d.enabled)

		result := render(d.input)
		if result != d.expected {
			t.Errorf("Test %d failed. BlackFriday enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}

		result = renderWithMmark(d.input)
		if result != d.expected {
			t.Errorf("Test %d failed. Mmark enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}
	}
}
