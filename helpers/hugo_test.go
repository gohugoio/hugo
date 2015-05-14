package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHugoVersion(t *testing.T) {
	assert.Equal(t, "0.15-DEV", hugoVersion(0.15, "-DEV"))
	assert.Equal(t, "0.17", hugoVersionNoSuffix(0.16+0.01))
}
