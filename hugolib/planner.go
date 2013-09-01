package hugolib

import (
	"fmt"
	"io"
)

func (s *Site) ShowPlan(out io.Writer) (err error) {
	if len(s.Files) <= 0 {
		fmt.Fprintf(out, "No source files provided.\n")
	}

	for _, file := range s.Files {
		fmt.Fprintf(out, "%s\n", file)
		if s.Target == nil {
			fmt.Fprintf(out, " *implicit* => %s\n", "!no target specified!")
			continue
		}
	}
	return
}
