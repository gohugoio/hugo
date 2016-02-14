package api

import (
	"github.com/spf13/hugo/commands"
	"github.com/spf13/viper"
)

func Run(flags []string) {
	commands.Execute(flags)
}

func Reset() {
	commands.ClearSite()
	viper.Reset()
}
