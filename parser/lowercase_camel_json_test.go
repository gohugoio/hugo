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

func TestReplacingJSONMarshallerOmitEmptyMapsIssue14855(t *testing.T) {
	c := qt.New(t)

	m := map[string]any{
		"dropBranch": map[string]any{
			"inner": map[string]any{
				"emptyString": "",
				"emptySlice":  []any{},
				"nested": map[string]any{
					"zero": 0,
				},
			},
		},
		"keep": "yes",
		"slice": []any{
			map[string]any{
				"name": "a",
				"sub": map[string]any{
					"zero": 0,
				},
			},
		},
		"target": map[string]any{
			"path": "/{photo,photo/**}",
			"sites": map[string]any{
				"complements": map[string]any{
					"languages": []any{},
				},
				"matrix": map[string]any{
					"languages": []any{},
					"versions":  []any{},
				},
			},
		},
	}

	marshaller := ReplacingJSONMarshaller{
		Value:       m,
		KeysToLower: true,
		OmitEmpty:   true,
	}

	b, err := marshaller.MarshalJSON()
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, `{"keep":"yes","slice":[{"name":"a"}],"target":{"path":"/{photo,photo/**}"}}`)

	marshaller.OmitEmpty = false
	b, err = marshaller.MarshalJSON()
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Contains, `"sites"`)
	c.Assert(string(b), qt.Contains, `"matrix"`)
	c.Assert(string(b), qt.Contains, `"dropbranch"`)
}
