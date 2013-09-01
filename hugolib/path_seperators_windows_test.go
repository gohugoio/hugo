package hugolib

import (
	"testing"
)

const (
	win_base = "c:\\a\\windows\\path\\layout"
	win_path = "c:\\a\\windows\\path\\layout\\sub1\\index.html"
)

func TestTemplatePathSeperator(t *testing.T) {
	tmpl := new(GoHtmlTemplate)
	if name := tmpl.generateTemplateNameFrom(win_base, win_path); name != "sub1/index.html" {
		t.Fatalf("Template name incorrect.  Expected: %s, Got: %s", "sub1/index.html", name)
	}
}
