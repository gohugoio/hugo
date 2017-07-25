package inflect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInflect(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		fn     func(i interface{}) (string, error)
		in     interface{}
		expect interface{}
	}{
		{ns.Humanize, "MyCamel", "My camel"},
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := test.fn(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
