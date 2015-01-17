package hugolib

import (
	"github.com/spf13/hugo/tpl"
	"testing"
)

const (
	win_base = "c:\\a\\windows\\path\\layout"
	win_path = "c:\\a\\windows\\path\\layout\\sub1\\index.html"
)

func TestTemplatePathSeparator(t *testing.T) {
	tmpl := new(tpl.GoHtmlTemplate)
	if name := tmpl.GenerateTemplateNameFrom(win_base, win_path); name != "sub1/index.html" {
		t.Fatalf("Template name incorrect.  Expected: %s, Got: %s", "sub1/index.html", name)
	}
}
