package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var genmanCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages for the Hugo CLI",
	Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		genmandir := "man/"
		cmd.Root().DisableAutoGenTag = true
		header := &cobra.GenManHeader{
			Section: "1",
			Manual:  "Hugo Manual",
			Source:  fmt.Sprintf("Hugo %s", helpers.HugoVersion()),
		}
		jww.FEEDBACK.Println("Generating Hugo man pages in", genmandir, "...")
		cmd.Root().GenManTree(header, genmandir)
		jww.FEEDBACK.Println("Done.")
	},
}
