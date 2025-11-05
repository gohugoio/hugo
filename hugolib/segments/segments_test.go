package segments

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

var (
	testDims    = sitesmatrix.NewTestingDimensions([]string{"no", "en"}, []string{"v1.0.0", "v1.2.0", "v2.0.0"}, []string{"guest", "member"})
	no          = sitesmatrix.Vector{0, 0, 0}
	en          = sitesmatrix.Vector{1, 0, 0}
	sitesNoStar = sitesmatrix.Sites{
		Matrix: sitesmatrix.StringSlices{
			Languages: []string{"n*"},
		},
	}
	sitesNo = sitesmatrix.Sites{
		Matrix: sitesmatrix.StringSlices{
			Languages: []string{"no"},
		},
	}

	sitesNotNo = sitesmatrix.Sites{
		Matrix: sitesmatrix.StringSlices{
			Languages: []string{"! no"},
		},
	}
)

var compileSegments = func(c *qt.C, includes ...SegmentMatcherFields) predicate.P[SegmentQuery] {
	segments := Segments{
		builder: &segmentsBuilder{
			logger:               loggers.NewDefault(),
			configuredDimensions: testDims,
			segmentsToRender:     []string{"docs"},
			segmentCfg: map[string]SegmentConfig{
				"docs": {
					Includes: includes,
				},
			},
		},
	}
	c.Assert(segments.compile(), qt.IsNil)
	return segments.SegmentFilter.ShouldExcludeFine
}

func TestCompileSegments(t *testing.T) {
	c := qt.New(t)

	c.Run("variants1", func(c *qt.C) {
		fields := []SegmentMatcherFields{
			{
				Sites:  sitesNoStar,
				Output: []string{"rss"},
			},
		}

		shouldExclude := compileSegments(c, fields...)

		check := func() {
			c.Assert(shouldExclude, qt.IsNotNil)
			c.Assert(shouldExclude(SegmentQuery{Site: no}), qt.Equals, true)
			c.Assert(shouldExclude(SegmentQuery{Site: no, Kind: "page"}), qt.Equals, true)
			c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "rss"}), qt.Equals, false)
			c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "html"}), qt.Equals, true)
			c.Assert(shouldExclude(SegmentQuery{Kind: "page"}), qt.Equals, true)
			c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "rss", Kind: "page"}), qt.Equals, false)
		}

		check()

		fields = []SegmentMatcherFields{
			{
				Path: []string{"/blog/**"},
			},
			{
				Sites:  sitesNoStar,
				Output: []string{"rss"},
			},
		}

		shouldExclude = compileSegments(c, fields...)
		check()
		c.Assert(shouldExclude(SegmentQuery{Path: "/blog/foo"}), qt.Equals, false)
	})

	c.Run("variants2", func(c *qt.C) {
		fields := []SegmentMatcherFields{
			{
				Path: []string{"/docs/**"},
			},
			{
				Sites:  sitesNo,
				Output: []string{"rss", "json"},
			},
		}

		shouldExclude := compileSegments(c, fields...)
		c.Assert(shouldExclude, qt.IsNotNil)
		c.Assert(shouldExclude(SegmentQuery{Site: no}), qt.Equals, true)
		c.Assert(shouldExclude(SegmentQuery{Kind: "page"}), qt.Equals, true)
		c.Assert(shouldExclude(SegmentQuery{Kind: "page", Path: "/blog/foo"}), qt.Equals, true)
		c.Assert(shouldExclude(SegmentQuery{Site: en}), qt.Equals, true)
		c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "rss"}), qt.Equals, false)
		c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "json"}), qt.Equals, false)
		c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "html"}), qt.Equals, true)
		c.Assert(shouldExclude(SegmentQuery{Kind: "page", Path: "/docs/foo"}), qt.Equals, false)
	})
}

func TestCompileSegmentsNegate(t *testing.T) {
	c := qt.New(t)

	fields := []SegmentMatcherFields{
		{
			Sites:  sitesNotNo,
			Output: []string{"! r**", "rem", "html"},
		},
	}

	shouldExclude := compileSegments(c, fields...)
	c.Assert(shouldExclude, qt.IsNotNil)
	c.Assert(shouldExclude(SegmentQuery{Site: no, Output: "html"}), qt.Equals, true)
	c.Assert(shouldExclude(SegmentQuery{Site: en, Output: "html"}), qt.Equals, false)
	c.Assert(shouldExclude(SegmentQuery{Site: en, Output: "rss"}), qt.Equals, true)
	c.Assert(shouldExclude(SegmentQuery{Site: en, Output: "rem"}), qt.Equals, true)
}

func BenchmarkSegmentsMatch(b *testing.B) {
	c := qt.New(b)
	fields := []SegmentMatcherFields{
		{
			Path: []string{"/docs/**"},
		},
		{
			Sites:  sitesNo,
			Output: []string{"rss"},
		},
	}

	match := compileSegments(c, fields...)

	for b.Loop() {
		match(SegmentQuery{Site: no, Output: "rss"})
	}
}
