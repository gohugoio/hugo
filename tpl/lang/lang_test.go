package lang

import (
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumFormat(t *testing.T) {
	t.Parallel()

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

	for i, c := range cases {
		errMsg := fmt.Sprintf("[%d] %v", i, c)

		var s string
		var err error

		if len(c.runes) == 0 {
			s, err = ns.NumFmt(c.prec, c.n)
		} else {
			if c.delim == "" {
				s, err = ns.NumFmt(c.prec, c.n, c.runes)
			} else {
				s, err = ns.NumFmt(c.prec, c.n, c.runes, c.delim)
			}
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, c.want, s, errMsg)
	}
}
