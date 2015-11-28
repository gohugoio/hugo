package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

var genmandir string
var genmanCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages for the Hugo CLI",
	Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		header := &cobra.GenManHeader{
			Section: "1",
			Manual:  "Hugo Manual",
			Source:  fmt.Sprintf("Hugo %s", helpers.HugoVersion()),
		}
		if !strings.HasSuffix(genmandir, helpers.FilePathSeparator) {
			genmandir += helpers.FilePathSeparator
		}
		if found, _ := helpers.Exists(genmandir, hugofs.OsFs); !found {
			jww.FEEDBACK.Println("Directory", genmandir, "does not exist, creating...")
			hugofs.OsFs.MkdirAll(genmandir, 0777)
		}
		cmd.Root().DisableAutoGenTag = true

		jww.FEEDBACK.Println("Generating Hugo man pages in", genmandir, "...")
		cmd.Root().GenManTree(header, genmandir)

		jww.FEEDBACK.Println("Done.")
	},
}

func init() {
	genmanCmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")
}
