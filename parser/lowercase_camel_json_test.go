package parser

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestReplacingJSONMarshaler(t *testing.T) {
	c := qt.New(t)

	m := map[string]any{
		"foo":        "bar",
		"baz":        42,
		"zeroInt1":   0,
		"zeroInt2":   0,
		"zeroFloat":  0.0,
		"zeroString": "",
		"zeroBool":   false,
		"nil":        nil,
	}

	marshaler := ReplacingJSONMarshaler{
		Value:       m,
		KeysToLower: true,
		OmitEmpty:   true,
	}

	b, err := marshaler.MarshalJSON()
	c.Assert(err, qt.IsNil)

	c.Assert(string(b), qt.Equals, `{"baz":42,"foo":"bar"}`)
}
