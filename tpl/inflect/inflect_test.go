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
		fn     func(i interface{}) (string, error)
		in     interface{}
		expect interface{}
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
