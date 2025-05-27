package collections

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNewStack(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	s := NewStack[int]()

	c.Assert(s, qt.IsNotNil)
}

func TestStackBasic(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	s := NewStack[int]()

	c.Assert(s.Len(), qt.Equals, 0)

	s.Push(1)
	s.Push(2)
	s.Push(3)

	c.Assert(s.Len(), qt.Equals, 3)

	top, ok := s.Peek()
	c.Assert(ok, qt.Equals, true)
	c.Assert(top, qt.Equals, 3)

	popped, ok := s.Pop()
	c.Assert(ok, qt.Equals, true)
	c.Assert(popped, qt.Equals, 3)

	c.Assert(s.Len(), qt.Equals, 2)

	_, _ = s.Pop()
	_, _ = s.Pop()
	_, ok = s.Pop()

	c.Assert(ok, qt.Equals, false)
}

func TestStackDrain(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	s := NewStack[string]()
	s.Push("a")
	s.Push("b")

	got := s.Drain()

	c.Assert(got, qt.DeepEquals, []string{"a", "b"})
	c.Assert(s.Len(), qt.Equals, 0)
}

func TestStackDrainMatching(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)

	got := s.DrainMatching(func(v int) bool { return v%2 == 0 })

	c.Assert(got, qt.DeepEquals, []int{4, 2})
	c.Assert(s.Drain(), qt.DeepEquals, []int{1, 3})
}
