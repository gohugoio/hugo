package commands

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var autocompleteTarget string

// bash for now (zsh and others will come)
var autocompleteType string

var genautocompleteCmd = &cobra.Command{
	Use:   "genautocomplete",
	Short: "Generate shell autocompletion script for Hugo",
	Long: `Generates a shell autocompletion script for Hugo.
	
	NOTE: The current version supports Bash only. This should work for *nix systems with Bash installed.
	
	By default the file is written directly to /etc/bash_completion.d for convenience and the command may need superuser rights, e.g:
	
	sudo hugo genautocomplete
	
	Add --completionfile=/path/to/file flag to set alternative file-path and name.
	
	Logout and in again to reload the completion scripts or just source them in directly:
	
	. /etc/bash_completion
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if autocompleteType != "bash" {
			jww.FATAL.Fatalln("Only Bash is supported for now")
		}
		err := cmd.Root().GenBashCompletionFile(autocompleteTarget)
		if err != nil {
			jww.FATAL.Fatalln("Failed to generate shell completion file:", err)
		}
	},
}

func init() {
	genautocompleteCmd.PersistentFlags().StringVarP(&autocompleteTarget, "completionfile", "", "/etc/bash_completion.d/hugo.sh", "Autocompletion file")
	genautocompleteCmd.PersistentFlags().StringVarP(&autocompleteType, "type", "", "bash", "Autocompletion type (currently only bash supported)")

}
