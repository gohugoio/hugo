package helpers

import (
	"testing"
)

func TestStripHTML(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"<h1>strip h1 tag <h1>", "strip h1 tag "},
		{"<p> strip p tag </p>", " strip p tag \n"},
		{"</br> strip br<br>", " strip br\n"},
	}
	for i, d := range data {
		output := StripHTML(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}
