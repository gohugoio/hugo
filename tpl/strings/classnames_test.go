package strings

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestClassNames(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// The unit tests port from classnames
	// see https://github.com/JedWatson/classnames/blob/495e24d9/tests/index.js

	t.Run("keeps object keys with truthy values", func(t *testing.T) {
		classes := ns.ClassNames(map[any]any{
			"a": true,
			"b": false,
			"c": 0,
			"d": nil,
			"e": "",
			"f": 1,
			"g": []any{},
			"":  0,
		})
		c.Assert(classes, qt.Equals, "a f")
	})

	t.Run("joins arrays of class names and ignore falsy values", func(t *testing.T) {
		classes := ns.ClassNames("a", 0, nil, true, 1, "b")
		c.Assert(classes, qt.Equals, "a 1 b")
	})

	t.Run("supports heterogeneous arguments", func(t *testing.T) {
		classes := ns.ClassNames(map[any]any{"a": true}, "b", 0)
		c.Assert(classes, qt.Equals, "a b")
	})

	t.Run("should be trimmed", func(t *testing.T) {
		classes := ns.ClassNames("", "b", map[any]any{}, "")
		c.Assert(classes, qt.Equals, "b")
	})

	t.Run("returns an empty string for an empty configuration", func(t *testing.T) {
		classes := ns.ClassNames(map[any]any{})
		c.Assert(classes, qt.Equals, "")
	})

	t.Run("returns an empty string for an empty configuration", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", "b"})
		c.Assert(classes, qt.Equals, "a b")
	})

	t.Run("joins array arguments with string arguments", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", "b"}, "c")
		c.Assert(classes, qt.Equals, "a b c")

		classes = ns.ClassNames("c", []any{"a", "b"})
		c.Assert(classes, qt.Equals, "c a b")
	})

	t.Run("handles multiple array arguments", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", "b"}, []any{"c", "d"})
		c.Assert(classes, qt.Equals, "a b c d")
	})

	t.Run("handles arrays that include falsy and true values", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", 0, nil, false, true, "b"})
		c.Assert(classes, qt.Equals, "a b")
	})

	t.Run("handles arrays that include arrays", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", []any{"b", "c"}})
		c.Assert(classes, qt.Equals, "a b c")
	})

	t.Run("handles arrays that include arrays", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", map[any]any{"b": true, "c": false}})
		c.Assert(classes, qt.Equals, "a b")
	})

	t.Run("handles deep array recursion", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", []any{"b", []any{"c", map[any]any{"d": true}}}})
		c.Assert(classes, qt.Equals, "a b c d")
	})

	t.Run("handles arrays that are empty", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", []any{}})
		c.Assert(classes, qt.Equals, "a")
	})

	t.Run("handles nested arrays that have empty nested arrays", func(t *testing.T) {
		classes := ns.ClassNames([]any{"a", []any{[]any{}}})
		c.Assert(classes, qt.Equals, "a")
	})
}
