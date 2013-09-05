package hugolib

import (
	"fmt"
	"io"
)

func (s *Site) ShowPlan(out io.Writer) (err error) {
	if s.Source == nil || len(s.Source.Files()) <= 0 {
		fmt.Fprintf(out, "No source files provided.\n")
	}

	for _, p := range s.Pages {
		fmt.Fprintf(out, "%s\n", p.FileName)
		fmt.Fprintf(out, " canonical => ")
		if s.Target == nil {
			fmt.Fprintf(out, "%s\n", "!no target specified!")
			continue
		}

		trns, err := s.Target.Translate(p.OutFile)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "%s\n", trns)

	}
	return
}
