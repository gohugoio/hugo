// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package commands

import (
	"os"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/spf13/cobra"
)

type genChromaStyles struct {
	style string
	cmd   *cobra.Command
}

func createGenChromaStyles() *genChromaStyles {
	g := &genChromaStyles{
		cmd: &cobra.Command{
			Use:   "chromastyles",
			Short: "Generate CSS stylesheet for the Chroma code highlighter",
			Long: `Generate CSS stylesheet for the Chroma code highlighter for a given style. This stylesheet is needed if pygmentsUseClasses is enabled in config.

See https://help.farbox.com/pygments.html for preview of available styles`,
		},
	}

	g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return g.generate()
	}

	g.cmd.PersistentFlags().StringVarP(&g.style, "style", "", "default", "highlighter style (see https://help.farbox.com/pygments.html)")

	return g
}

func (g *genChromaStyles) generate() error {
	style := styles.Get(g.style)
	formatter := html.New(html.WithClasses())
	formatter.WriteCSS(os.Stdout, style)
	return nil
}
