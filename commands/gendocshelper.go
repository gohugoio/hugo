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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/docshelper"
	"github.com/spf13/cobra"
)

var (
	_ cmder = (*genDocsHelper)(nil)
)

type genDocsHelper struct {
	target string
	*baseCmd
}

func createGenDocsHelper() *genDocsHelper {
	g := &genDocsHelper{
		baseCmd: newBaseCmd(&cobra.Command{
			Use:    "docshelper",
			Short:  "Generate some data files for the Hugo docs.",
			Hidden: true,
		}),
	}

	g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return g.generate()
	}

	g.cmd.PersistentFlags().StringVarP(&g.target, "dir", "", "docs/data", "data dir")

	return g
}

func (g *genDocsHelper) generate() error {
	fmt.Println("Generate docs data to", g.target)

	targetFile := filepath.Join(g.target, "docs.json")

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(docshelper.GetDocProvider()); err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil

}
