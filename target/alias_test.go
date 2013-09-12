package target

import (
	"testing"
)

func TestHTMLRedirectAlias(t *testing.T) {
	var o Translator
	o = new(HTMLRedirectAlias)

	tests := []struct {
		value    string
		expected string
	}{
		{"", ""},
		{"alias 1", "alias-1"},
		{"alias 2/", "alias-2/index.html"},
		{"alias 3.html", "alias-3.html"},
	}

	for _, test := range tests {
		path, err := o.Translate(test.value)
		if err != nil {
			t.Fatalf("Translate returned an error: %s", err)
		}

		if path != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, path)
		}
	}
}
