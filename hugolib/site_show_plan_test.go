package hugolib

import (
	"bytes"
	"testing"
)

func TestDegenerateNoTarget(t *testing.T) {
	s := new(Site)
	out := new(bytes.Buffer)
	if err := s.ShowPlan(out); err != nil {
		t.Errorf("ShowPlan unexpectedly returned an error: %s", err)
	}
}
