package bundle

import (
	"testing"
)

func TestNothing(t *testing.T) {
	b := NewTemplate()
	if b.Lookup("alias") == nil {
		t.Fatalf("Expecting alias to be initialized with new bundle")
	}
}
