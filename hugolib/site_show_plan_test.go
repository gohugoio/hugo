package hugolib

import (
	"bytes"
	"testing"
)

func checkShowPlanExpected(t *testing.T, expected, got string) {
	if got != expected {
		t.Errorf("ShowPlan expected:\n%q\ngot\n%q", expected, got)
	}
}

func TestDegenerateNoFiles(t *testing.T) {
	s := new(Site)
	out := new(bytes.Buffer)
	if err := s.ShowPlan(out); err != nil {
		t.Errorf("ShowPlan unexpectedly returned an error: %s", err)
	}
	expected := "No source files provided.\n"
	got := out.String()
	checkShowPlanExpected(t, expected, got)
}

func TestDegenerateNoTarget(t *testing.T) {
	s := new(Site)
	s.Files = append(s.Files, "foo/bar/file.md")
	out := new(bytes.Buffer)
	if err := s.ShowPlan(out); err != nil {
		t.Errorf("ShowPlan unexpectedly returned an error: %s", err)
	}

	expected := "foo/bar/file.md\n *implicit* => !no target specified!\n"
	checkShowPlanExpected(t, expected, out.String())
}
