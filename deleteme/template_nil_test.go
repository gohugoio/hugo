package deleteme

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestTemplateNil(t *testing.T) {
	data := map[string]interface{}{}
	tpl := `Nil: {{ .noMatch }}`

	var buf bytes.Buffer
	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		t.Fatal(err)
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}

	result := strings.TrimSpace(buf.String())

	if result != "Nil:" {
		t.Fatal(result)
	}

}
