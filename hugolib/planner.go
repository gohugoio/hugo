// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
	"io"
)

// ShowPlan prints a build plan to the given Writer.
// Useful for debugging.
func (s *Site) ShowPlan(out io.Writer) (err error) {
	if s.Source == nil || len(s.Source.Files()) <= 0 {
		fmt.Fprintf(out, "No source files provided.\n")
	}

	for _, p := range s.Pages {
		fmt.Fprintf(out, "%s", p.Source.Path())
		if p.IsRenderable() {
			fmt.Fprintf(out, " (renderer: markdown)")
		} else {
			fmt.Fprintf(out, " (renderer: n/a)")
		}
		if s.Tmpl != nil {
			for _, l := range p.layouts() {
				fmt.Fprintf(out, " (layout: %s, exists: %t)", l, s.Tmpl.Lookup(l) != nil)
			}
		}
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, " canonical => ")
		if s.targets.page == nil {
			fmt.Fprintf(out, "%s\n\n", "!no target specified!")
			continue
		}

		trns, err := s.pageTarget().Translate(p.TargetPath())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s\n", trns)

		if s.targets.alias == nil {
			continue
		}

		for _, alias := range p.Aliases {
			aliasTrans, err := s.aliasTarget().Translate(alias)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, " %s => %s\n", alias, aliasTrans)
		}
		fmt.Fprintln(out)
	}
	return
}
