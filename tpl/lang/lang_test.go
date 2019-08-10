package lang

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
)

func TestNumFormat(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{})

	cases := []struct {
		prec  int
		n     float64
		runes string
		delim string

		want string
	}{
		{2, -12345.6789, "", "", "-12,345.68"},
		{2, -12345.6789, "- . ,", "", "-12,345.68"},
		{2, -12345.1234, "- . ,", "", "-12,345.12"},

		{2, 12345.6789, "- . ,", "", "12,345.68"},
		{0, 12345.6789, "- . ,", "", "12,346"},
		{11, -12345.6789, "- . ,", "", "-12,345.67890000000"},

		{3, -12345.6789, "- ,", "", "-12345,679"},
		{6, -12345.6789, "- , .", "", "-12.345,678900"},

		{3, -12345.6789, "-|,| ", "|", "-12 345,679"},
		{6, -12345.6789, "-|,| ", "|", "-12 345,678900"},

		// Arabic, ar_AE
		{6, -12345.6789, "‏- ٫ ٬", "", "‏-12٬345٫678900"},
		{6, -12345.6789, "‏-|٫| ", "|", "‏-12 345٫678900"},
	}

	for _, cas := range cases {
		var s string
		var err error

		if len(cas.runes) == 0 {
			s, err = ns.NumFmt(cas.prec, cas.n)
		} else {
			if cas.delim == "" {
				s, err = ns.NumFmt(cas.prec, cas.n, cas.runes)
			} else {
				s, err = ns.NumFmt(cas.prec, cas.n, cas.runes, cas.delim)
			}
		}

		c.Assert(err, qt.IsNil)
		c.Assert(s, qt.Equals, cas.want)
	}
}
