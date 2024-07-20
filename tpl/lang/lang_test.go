package lang

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	translators "github.com/gohugoio/localescompressed"
)

func TestNumFmt(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(&deps.Deps{}, nil)

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

		{2, 927.675, "- .", "", "927.68"},
		{2, 1927.675, "- .", "", "1927.68"},
		{2, 2927.675, "- .", "", "2927.68"},

		{3, -12345.6789, "- ,", "", "-12345,679"},
		{6, -12345.6789, "- , .", "", "-12.345,678900"},

		{3, -12345.6789, "-|,| ", "|", "-12 345,679"},
		{6, -12345.6789, "-|,| ", "|", "-12 345,678900"},

		// Arabic, ar_AE
		{6, -12345.6789, "\u200f- ٫ ٬", "", "\u200f-12٬345٫678900"},
		{6, -12345.6789, "\u200f-|٫| ", "|", "\u200f-12 345٫678900"},
	}

	for _, cas := range cases {
		var s string
		var err error

		if len(cas.runes) == 0 {
			s, err = ns.FormatNumberCustom(cas.prec, cas.n)
		} else {
			if cas.delim == "" {
				s, err = ns.FormatNumberCustom(cas.prec, cas.n, cas.runes)
			} else {
				s, err = ns.FormatNumberCustom(cas.prec, cas.n, cas.runes, cas.delim)
			}
		}

		c.Assert(err, qt.IsNil)
		c.Assert(s, qt.Equals, cas.want)
	}
}

func TestFormatNumbers(t *testing.T) {
	c := qt.New(t)

	nsNn := New(&deps.Deps{}, translators.GetTranslator("nn"))
	nsEn := New(&deps.Deps{}, translators.GetTranslator("en"))
	pi := 3.14159265359

	c.Run("FormatNumber", func(c *qt.C) {
		c.Parallel()
		got, err := nsNn.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3,142")

		got, err = nsEn.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3.142")
	})

	c.Run("FormatPercent", func(c *qt.C) {
		c.Parallel()
		got, err := nsEn.FormatPercent(3, 67.33333)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "67.333%")
	})

	c.Run("FormatCurrency", func(c *qt.C) {
		c.Parallel()
		got, err := nsEn.FormatCurrency(2, "USD", 20000)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "$20,000.00")
	})

	c.Run("FormatAccounting", func(c *qt.C) {
		c.Parallel()
		got, err := nsEn.FormatAccounting(2, "USD", 20000)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "$20,000.00")
	})
}

// Issue 9446
func TestLanguageKeyFormat(t *testing.T) {
	c := qt.New(t)

	nsUnderscoreUpper := New(&deps.Deps{}, translators.GetTranslator("es_ES"))
	nsUnderscoreLower := New(&deps.Deps{}, translators.GetTranslator("es_es"))
	nsHyphenUpper := New(&deps.Deps{}, translators.GetTranslator("es-ES"))
	nsHyphenLower := New(&deps.Deps{}, translators.GetTranslator("es-es"))
	pi := 3.14159265359

	c.Run("FormatNumber", func(c *qt.C) {
		c.Parallel()
		got, err := nsUnderscoreUpper.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3,142")

		got, err = nsUnderscoreLower.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3,142")

		got, err = nsHyphenUpper.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3,142")

		got, err = nsHyphenLower.FormatNumber(3, pi)
		c.Assert(err, qt.IsNil)
		c.Assert(got, qt.Equals, "3,142")
	})
}
