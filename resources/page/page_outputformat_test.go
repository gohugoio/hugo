package page

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/output"
)

func TestOutputFormat_String(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	format := OutputFormat{Rel: "alternate", Format: output.HTMLFormat, permalink: "https://example.com/"}
	expected := `<link rel="alternate" type="text/html" href="https://example.com/">`
	c.Assert(format.String(), qt.Equals, expected)
	formats := OutputFormats{format, format}
	c.Assert(formats.String(), qt.Equals, expected+"\n"+expected+"\n")
}
