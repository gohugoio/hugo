package helpers

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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
	inputCopy := input

	SliceToLower(input)

	for i, e := range input {
		if e != inputCopy[i] {
			t.Errorf("TestSliceToLowerNonDestructive failed. Expected element #%d of input slice to be %s. Found %s.", i, inputCopy[i], input[i])
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

func TestDoArithmetic(t *testing.T) {
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		op     rune
		expect interface{}
	}{
		{3, 2, '+', int64(5)},
		{3, 2, '-', int64(1)},
		{3, 2, '*', int64(6)},
		{3, 2, '/', int64(1)},
		{3.0, 2, '+', float64(5)},
		{3.0, 2, '-', float64(1)},
		{3.0, 2, '*', float64(6)},
		{3.0, 2, '/', float64(1.5)},
		{3, 2.0, '+', float64(5)},
		{3, 2.0, '-', float64(1)},
		{3, 2.0, '*', float64(6)},
		{3, 2.0, '/', float64(1.5)},
		{3.0, 2.0, '+', float64(5)},
		{3.0, 2.0, '-', float64(1)},
		{3.0, 2.0, '*', float64(6)},
		{3.0, 2.0, '/', float64(1.5)},
		{uint(3), uint(2), '+', uint64(5)},
		{uint(3), uint(2), '-', uint64(1)},
		{uint(3), uint(2), '*', uint64(6)},
		{uint(3), uint(2), '/', uint64(1)},
		{uint(3), 2, '+', uint64(5)},
		{uint(3), 2, '-', uint64(1)},
		{uint(3), 2, '*', uint64(6)},
		{uint(3), 2, '/', uint64(1)},
		{3, uint(2), '+', uint64(5)},
		{3, uint(2), '-', uint64(1)},
		{3, uint(2), '*', uint64(6)},
		{3, uint(2), '/', uint64(1)},
		{uint(3), -2, '+', int64(1)},
		{uint(3), -2, '-', int64(5)},
		{uint(3), -2, '*', int64(-6)},
		{uint(3), -2, '/', int64(-1)},
		{-3, uint(2), '+', int64(-1)},
		{-3, uint(2), '-', int64(-5)},
		{-3, uint(2), '*', int64(-6)},
		{-3, uint(2), '/', int64(-1)},
		{uint(3), 2.0, '+', float64(5)},
		{uint(3), 2.0, '-', float64(1)},
		{uint(3), 2.0, '*', float64(6)},
		{uint(3), 2.0, '/', float64(1.5)},
		{3.0, uint(2), '+', float64(5)},
		{3.0, uint(2), '-', float64(1)},
		{3.0, uint(2), '*', float64(6)},
		{3.0, uint(2), '/', float64(1.5)},
		{0, 0, '+', 0},
		{0, 0, '-', 0},
		{0, 0, '*', 0},
		{"foo", "bar", '+', "foobar"},
		{3, 0, '/', false},
		{3.0, 0, '/', false},
		{3, 0.0, '/', false},
		{uint(3), uint(0), '/', false},
		{3, uint(0), '/', false},
		{-3, uint(0), '/', false},
		{uint(3), 0, '/', false},
		{3.0, uint(0), '/', false},
		{uint(3), 0.0, '/', false},
		{3, "foo", '+', false},
		{3.0, "foo", '+', false},
		{uint(3), "foo", '+', false},
		{"foo", 3, '+', false},
		{"foo", "bar", '-', false},
		{3, 2, '%', false},
	} {
		result, err := DoArithmetic(this.a, this.b, this.op)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] doArithmetic didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] doArithmetic got %v but expected %v", i, result, this.expect)
			}
		}
	}
}
