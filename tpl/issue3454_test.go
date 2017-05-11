package tpl

import (
	"bytes"
	"html/template"
	"testing"
)

func TestDisallowedEscaper(t *testing.T) {
	data := map[string]string{
		"html": "<h1>Hi!</h1>",
	}

	tpl := `{{ .html | print }}`

	var buf bytes.Buffer
	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}

	if buf.String() != `&lt;h1&gt;Hi!&lt;/h1&gt;` {
		t.Fatalf("Got %q", buf.String())
	}

}
