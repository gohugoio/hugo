package parser

import (
	"testing"
)

func TestFormatToLeadRune(t *testing.T) {
	for i, this := range []struct {
		kind   string
		expect rune
	}{
		{"yaml", '-'},
		{"yml", '-'},
		{"toml", '+'},
		{"json", '{'},
		{"js", '{'},
		{"unknown", '+'},
	} {
		result := FormatToLeadRune(this.kind)

		if result != this.expect {
			t.Errorf("[%d] got %q but expected %q", i, result, this.expect)
		}
	}
}
