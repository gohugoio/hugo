package hugolib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScratchAdd(t *testing.T) {
	scratch := newScratch()
	scratch.Add("int1", 10)
	scratch.Add("int1", 20)
	scratch.Add("int2", 20)

	assert.Equal(t, int64(30), scratch.Get("int1"))
	assert.Equal(t, 20, scratch.Get("int2"))

	scratch.Add("float1", float64(10.5))
	scratch.Add("float1", float64(20.1))

	assert.Equal(t, float64(30.6), scratch.Get("float1"))

	scratch.Add("string1", "Hello ")
	scratch.Add("string1", "big ")
	scratch.Add("string1", "World!")

	assert.Equal(t, "Hello big World!", scratch.Get("string1"))

	scratch.Add("scratch", scratch)
	_, err := scratch.Add("scratch", scratch)

	if err == nil {
		t.Errorf("Expected error from invalid arithmetic")
	}

}

func TestScratchSet(t *testing.T) {
	scratch := newScratch()
	scratch.Set("key", "val")
	assert.Equal(t, "val", scratch.Get("key"))
}

func TestScratchGet(t *testing.T) {
	scratch := newScratch()
	nothing := scratch.Get("nothing")
	if nothing != nil {
		t.Errorf("Should not return anything, but got %v", nothing)
	}
}
