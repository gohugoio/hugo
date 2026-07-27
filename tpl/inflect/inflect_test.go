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
		{ns.Humanize, "\u00f3bito", "\u00d3bito"},
		{ns.Humanize, "", ""},
		{ns.Humanize, "103", "103rd"},
		{ns.Humanize, "41", "41st"},
		{ns.Humanize, 103, "103rd"},
		{ns.Humanize, int64(92), "92nd"},
		{ns.Humanize, "5.5", "5.5"},
		{ns.Humanize, t, false},
		{ns.Humanize, "this is a TEST", "This is a test"},
		{ns.Humanize, "my-first-Post", "My first post"},
		// Issue #15126: flect incorrectly treats quote+uppercase as a camelCase
		// boundary, inserting a spurious space after the opening quote.
		{ns.Humanize, "Painting \u201eTitle\u201c", "Painting \u201etitle\u201c"}, // German „Title" (from issue)
		{ns.Humanize, "a \u201eB\u201c", "A \u201eb\u201c"},                       // German „B" uppercase
		{ns.Humanize, "a \u201eb\u201c", "A \u201eb\u201c"},                       // German „b" lowercase (baseline, no split)
		{ns.Humanize, "a \"B\"", "A \"b\""},                                       // ASCII " uppercase
		{ns.Humanize, "A \"B\"", "A \"b\""},                                       // ASCII " uppercase, capital-initial first word
		{ns.Humanize, "a \"b\"", "A \"b\""},                                       // ASCII " lowercase (baseline, no split)
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
