package segments

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

func TestCompileSegments(t *testing.T) {
	c := qt.New(t)
	dims := sitesmatrix.NewTestingDimensions([]string{"no", "en"}, []string{"v1.0.0", "v1.2.0", "v2.0.0"}, []string{"guest", "member"})
	sitesNotNo := sitesmatrix.Sites{
		Matrix: sitesmatrix.StringSlices{
			Languages: []string{"! no"},
		},
	}
	if sitesNotNo.IsZero() {
	}
	segments := Segments{
		builder: &segmentsBuilder{
			configuredDimensions: dims,
			segmentsToRender:     []string{"docs"},
			segmentCfg: map[string]SegmentConfig{
				"docs": {
					Rules: []SegmentMatcherRules{
						{
							Sites: sitesNotNo,
							Path:  []string{"! /bar/**", "{/docs,/docs/**}"},
						},
						{
							Kind: []string{"! {home,term}", "section"},
						},
						{
							Kind: []string{"page"},
							Path: []string{"/regularpages/**"},
						},
					},
				},
			},
		},
	}
	c.Assert(segments.compile(), qt.IsNil)

	include := segments.IncludeSegment
	no := &sitesmatrix.Vector{0, 0, 0}
	en := &sitesmatrix.Vector{1, 0, 0}

	// c.Assert(filter((SegmentMatcherQuery{Output: "rss"}), qt.Equals, true)
	c.Assert(include(SegmentMatcherQuery{Kind: "section"}), qt.Equals, true)

	c.Assert(include(SegmentMatcherQuery{Kind: "term"}), qt.Equals, false)
	c.Assert(include(SegmentMatcherQuery{Path: "/docs", Site: no}), qt.Equals, false)
	c.Assert(include(SegmentMatcherQuery{Path: "/docs", Site: en}), qt.Equals, true)
	c.Assert(include(SegmentMatcherQuery{Site: en}), qt.Equals, true)
	c.Assert(include(SegmentMatcherQuery{Kind: "page"}), qt.Equals, true)
}
