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

// See issue 14855.
func TestReplacingJSONMarshallerOmitEmptySubMaps(t *testing.T) {
	c := qt.New(t)

	m := map[string]any{
		"keep": "yes",
		"target": map[string]any{
			"path": "/x",
			"sites": map[string]any{
				"matrix": map[string]any{
					"languages": []any{},
					"versions":  []any{},
				},
				"complements": map[string]any{
					"languages": []any{},
				},
			},
		},
		"all_zero_branch": map[string]any{
			"inner": map[string]any{
				"empty": "",
				"deep": map[string]any{
					"zero": 0,
				},
			},
		},
		"keep_slice_of_maps": []any{
			map[string]any{
				"name":  "a",
				"empty": "",
				"sub":   map[string]any{"x": 0},
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

	s := string(b)
	c.Assert(s, qt.Contains, `"keep":"yes"`)
	c.Assert(s, qt.Contains, `"path":"/x"`)
	c.Assert(s, qt.Not(qt.Contains), `"sites"`)
	c.Assert(s, qt.Not(qt.Contains), `"matrix"`)
	c.Assert(s, qt.Not(qt.Contains), `"complements"`)
	c.Assert(s, qt.Not(qt.Contains), `"all_zero_branch"`)
	c.Assert(s, qt.Not(qt.Contains), `"inner"`)
	c.Assert(s, qt.Not(qt.Contains), `"deep"`)
	c.Assert(s, qt.Contains, `"name":"a"`)
	c.Assert(s, qt.Not(qt.Contains), `"sub"`)
}
