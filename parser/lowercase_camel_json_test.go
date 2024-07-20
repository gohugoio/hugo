package parser

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestReplacingJSONMarshaller(t *testing.T) {
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

	marshaller := ReplacingJSONMarshaller{
		Value:       m,
		KeysToLower: true,
		OmitEmpty:   true,
	}

	b, err := marshaller.MarshalJSON()
	c.Assert(err, qt.IsNil)

	c.Assert(string(b), qt.Equals, `{"baz":42,"foo":"bar"}`)
}
