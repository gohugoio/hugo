package helpers

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGuessType(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect string
	}{
		{"md", "markdown"},
		{"markdown", "markdown"},
		{"mdown", "markdown"},
		{"asciidoc", "asciidoc"},
		{"ad", "asciidoc"},
		{"rst", "rst"},
		{"html", "html"},
		{"htm", "html"},
		{"excel", "unknown"},
	} {
		result := GuessType(this.in)
		if result != this.expect {
			t.Errorf("[%d] GuessType guessed wrong, expected %s, got %s", i, this.expect, result)
		}
	}
}

func TestBytesToReader(t *testing.T) {
	asBytes := ReaderToBytes(strings.NewReader("Hello World!"))
	asReader := BytesToReader(asBytes)
	assert.Equal(t, []byte("Hello World!"), asBytes)
	assert.Equal(t, asBytes, ReaderToBytes(asReader))
}

func TestStringToReader(t *testing.T) {
	asString := ReaderToString(strings.NewReader("Hello World!"))
	assert.Equal(t, "Hello World!", asString)
	asReader := StringToReader(asString)
	assert.Equal(t, asString, ReaderToString(asReader))
}

func TestFindAvailablePort(t *testing.T) {
	addr, err := FindAvailablePort()
	assert.Nil(t, err)
	assert.NotNil(t, addr)
	assert.True(t, addr.Port > 0)
}

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
