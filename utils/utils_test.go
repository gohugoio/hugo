package utils

import (
	"testing"
)

func TestCutUsageMessage(t *testing.T) {
	tests := []struct {
		message    string
		cutMessage string
	}{
		{"", ""},
		{" Usage of hugo: \n  -b, --baseUrl=...", ""},
		{"Some error Usage of hugo: \n", "Some error"},
		{"Usage of hugo: \n -b --baseU", ""},
		{"CRITICAL error for usage of hugo ", "CRITICAL error for usage of hugo"},
		{"Invalid short flag a in -abcde", "Invalid short flag a in -abcde"},
	}

	for _, test := range tests {
		message := cutUsageMessage(test.message)
		if message != test.cutMessage {
			t.Errorf("Expected %#v, got %#v", test.cutMessage, message)
		}
	}
}
