package hugolib

import (
	"testing"
)

func TestTemplatePathSeperator(t *testing.T) {
	config := Config{
		LayoutDir: "c:\\a\\windows\\path\\layout",
		Path:      "c:\\a\\windows\\path",
	}
	s := &Site{Config: config}
	if name := s.generateTemplateNameFrom("c:\\a\\windows\\path\\layout\\sub1\\index.html"); name != "sub1/index.html" {
		t.Fatalf("Template name incorrect.  Expected: %s, Got: %s", "sub1/index.html", name)
	}
}
