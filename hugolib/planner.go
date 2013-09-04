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
		fmt.Fprintf(out, " canonical => ")
		if s.Target == nil {
			fmt.Fprintf(out, "%s\n", "!no target specified!")
			continue
		}

		trns, err := s.Target.Translate(file)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "%s\n", trns)

	}
	return
}
