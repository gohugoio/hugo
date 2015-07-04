package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHugoVersion(t *testing.T) {
	assert.Equal(t, "0.15-DEV", hugoVersion(0.15, "-DEV"))
	assert.Equal(t, "0.17", hugoVersionNoSuffix(0.16+0.01))
}
