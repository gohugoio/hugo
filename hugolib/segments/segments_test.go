package segments

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCompileSegments(t *testing.T) {
	c := qt.New(t)

	c.Run("excludes", func(c *qt.C) {
		fields := []SegmentMatcherFields{
			{
				Lang:   "n*",
				Output: "rss",
			},
		}

		match, err := compileSegments(fields)
		c.Assert(err, qt.IsNil)

		check := func() {
			c.Assert(match, qt.IsNotNil)
			c.Assert(match(SegmentMatcherFields{Lang: "no"}), qt.Equals, false)
			c.Assert(match(SegmentMatcherFields{Lang: "no", Kind: "page"}), qt.Equals, false)
			c.Assert(match(SegmentMatcherFields{Lang: "no", Output: "rss"}), qt.Equals, true)
			c.Assert(match(SegmentMatcherFields{Lang: "no", Output: "html"}), qt.Equals, false)
			c.Assert(match(SegmentMatcherFields{Kind: "page"}), qt.Equals, false)
			c.Assert(match(SegmentMatcherFields{Lang: "no", Output: "rss", Kind: "page"}), qt.Equals, true)
		}

		check()

		fields = []SegmentMatcherFields{
			{
				Path: "/blog/**",
			},
			{
				Lang:   "n*",
				Output: "rss",
			},
		}

		match, err = compileSegments(fields)
		c.Assert(err, qt.IsNil)
		check()
		c.Assert(match(SegmentMatcherFields{Path: "/blog/foo"}), qt.Equals, true)
	})

	c.Run("includes", func(c *qt.C) {
		fields := []SegmentMatcherFields{
			{
				Path: "/docs/**",
			},
			{
				Lang:   "no",
				Output: "rss",
			},
		}

		match, err := compileSegments(fields)
		c.Assert(err, qt.IsNil)
		c.Assert(match, qt.IsNotNil)
		c.Assert(match(SegmentMatcherFields{Lang: "no"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Kind: "page"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Kind: "page", Path: "/blog/foo"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Lang: "en"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Lang: "no", Output: "rss"}), qt.Equals, true)
		c.Assert(match(SegmentMatcherFields{Lang: "no", Output: "html"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Kind: "page", Path: "/docs/foo"}), qt.Equals, true)
	})

	c.Run("includes variant1", func(c *qt.C) {
		c.Skip()

		fields := []SegmentMatcherFields{
			{
				Kind: "home",
			},
			{
				Path: "{/docs,/docs/**}",
			},
		}

		match, err := compileSegments(fields)
		c.Assert(err, qt.IsNil)
		c.Assert(match, qt.IsNotNil)
		c.Assert(match(SegmentMatcherFields{Path: "/blog/foo"}), qt.Equals, false)
		c.Assert(match(SegmentMatcherFields{Kind: "page", Path: "/docs/foo"}), qt.Equals, true)
		c.Assert(match(SegmentMatcherFields{Kind: "home", Path: "/"}), qt.Equals, true)
	})
}

func BenchmarkSegmentsMatch(b *testing.B) {
	fields := []SegmentMatcherFields{
		{
			Path: "/docs/**",
		},
		{
			Lang:   "no",
			Output: "rss",
		},
	}

	match, err := compileSegments(fields)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		match(SegmentMatcherFields{Lang: "no", Output: "rss"})
	}
}
