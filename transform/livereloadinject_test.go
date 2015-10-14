package transform

import (
	"bytes"
	"github.com/spf13/hugo/helpers"
	"testing"
)

func TestLiveReloadInject(t *testing.T) {
	out := new(bytes.Buffer)
	in := helpers.StringToReader("</body>")

	tr := NewChain(LiveReloadInject)
	tr.Apply(out, in, []byte("path"))

	expected := `<script data-no-instant>document.write('<script src="/livereload.js?mindelay=10"></' + 'script>')</script></body>`
	if string(out.Bytes()) != expected {
		t.Errorf("Expected %s got %s", expected, string(out.Bytes()))
	}
}
