package source

import (
	"github.com/spf13/hugo/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileUniqueID(t *testing.T) {
	f1 := File{uniqueID: "123"}
	f2 := NewFile("a")

	assert.Equal(t, "123", f1.UniqueID())
	assert.Equal(t, "0cc175b9c0f1b6a831c399e269772661", f2.UniqueID())
}

func TestFileString(t *testing.T) {
	assert.Equal(t, "abc", NewFileWithContents("a", helpers.StringToReader("abc")).String())
	assert.Equal(t, "", NewFile("a").String())
}

func TestFileBytes(t *testing.T) {
	assert.Equal(t, []byte("abc"), NewFileWithContents("a", helpers.StringToReader("abc")).Bytes())
	assert.Equal(t, []byte(""), NewFile("a").Bytes())
}
