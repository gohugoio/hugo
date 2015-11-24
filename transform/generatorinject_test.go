package transform

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/hugo/helpers"
)

func TestGeneratorInject(t *testing.T) {
	out := new(bytes.Buffer)
	in := helpers.StringToReader("</head>")

	tr := NewChain(GeneratorInject)
	tr.Apply(out, in, []byte("path"))

	if !strings.HasSuffix(string(out.Bytes()), "</head>") {
		t.Errorf("Expected suffix \"</head>\" got %s", string(out.Bytes()))
	}

	if !strings.HasPrefix(string(out.Bytes()), "<meta name=\"generator\" content=\"Hugo") {
		t.Errorf("Expected prefix %q got %s", "<meta name=\"generator\" content=\"Hugo", string(out.Bytes()))
	}
}
