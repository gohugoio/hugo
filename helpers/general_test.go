package helpers

import (
	"strings"
	"testing"
)

func TestInStringArrayCaseSensitive(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}
	array := []string{
		"Albert",
		"Joe",
		"Francis",
	}
	data := []test{
		{"Albert", true},
		{"ALBERT", false},
	}
	for _, in := range data {
		output := InStringArray(array, in.input)
		if output != in.expected {
			t.Errorf("TestInStringArrayCase failed. Expected %t. Got %t.", in.expected, output)
		}
	}
}

func TestSliceToLowerStable(t *testing.T) {
	input := []string{
		"New York",
		"BARCELONA",
		"COffEE",
		"FLOWer",
		"CanDY",
	}

	output := SliceToLower(input)

	for i, e := range output {
		if e != strings.ToLower(input[i]) {
			t.Errorf("Expected %s. Found %s.", strings.ToLower(input[i]), e)
		}
	}
}

func TestSliceToLowerNil(t *testing.T) {
	var input []string

	output := SliceToLower(input)

	if output != nil {
		t.Errorf("Expected nil to be returned. Had %s.", output)
	}
}

func TestSliceToLowerNonDestructive(t *testing.T) {
	input := []string{
		"New York",
		"BARCELONA",
		"COffEE",
		"FLOWer",
		"CanDY",
	}

	// This assignment actually copies the content
	// of input into a new object.
	// Otherwise, the test would not make sense...
	input_copy := input

	SliceToLower(input)

	for i, e := range input {
		if e != input_copy[i] {
			t.Errorf("TestSliceToLowerNonDestructive failed. Expected element #%d of input slice to be %s. Found %s.", i, input_copy[i], input[i])
		}
	}
}

// Just make sure there is no error for empty-like strings
func TestMd5StringEmpty(t *testing.T) {
	inputs := []string{"", " ", "   "}

	for _, in := range inputs {
		Md5String(in)
	}
}
