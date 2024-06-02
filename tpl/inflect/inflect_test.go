package inflect

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestInflect(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		fn     func(i any) (string, error)
		in     any
		expect any
	}{
		{ns.Humanize, "MyCamel", "My camel"},
		{ns.Humanize, "óbito", "Óbito"},
		{ns.Humanize, "", ""},
		{ns.Humanize, "103", "103rd"},
		{ns.Humanize, "41", "41st"},
		{ns.Humanize, 103, "103rd"},
		{ns.Humanize, int64(92), "92nd"},
		{ns.Humanize, "5.5", "5.5"},
		{ns.Humanize, t, false},
		{ns.Humanize, "this is a TEST", "This is a test"},
		{ns.Humanize, "my-first-Post", "My first post"},
		{ns.Pluralize, "cat", "cats"},
		{ns.Pluralize, "", ""},
		{ns.Pluralize, t, false},
		{ns.Singularize, "cats", "cat"},
		{ns.Singularize, "", ""},
		{ns.Singularize, t, false},
	} {

		result, err := test.fn(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestSI(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		in     any
		args   []any
		expect any
	}{
		{12, []any{}, "12"},
		{"123", []any{}, "123"},
		{1234, []any{}, "1.234 k"},
		{"2345", []any{}, "2.345 k"},
		{1234000, []any{}, "1.234 M"},
		{"2345000", []any{}, "2.345 M"},
		{"0.00000000223", []any{"M"}, "2.23 nM"},
		{"1000000", []any{"B"}, "1 MB"},
		{"2.2345e-12", []any{"F"}, "2.2345 pF"},
		{"invalid-number", []any{}, false},
	} {

		result, err := ns.SI(test.in, test.args...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}
