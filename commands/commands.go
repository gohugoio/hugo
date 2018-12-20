// Copyright 2017 The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/hugolib/paths"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/nitro"
)

type commandsBuilder struct {
	hugoBuilderCommon

	commands []cmder
}

func newCommandsBuilder() *commandsBuilder {
	return &commandsBuilder{}
}

func (b *commandsBuilder) addCommands(commands ...cmder) *commandsBuilder {
	b.commands = append(b.commands, commands...)
	return b
}

func (b *commandsBuilder) addAll() *commandsBuilder {
	b.addCommands(
		b.newServerCmd(),
		newVersionCmd(),
		newEnvCmd(),
		newConfigCmd(),
		newCheckCmd(),
		newConvertCmd(),
		b.newNewCmd(),
		newListCmd(),
		newImportCmd(),
		newGenCmd(),
		createReleaser(),
	)

	return b
}

func (b *commandsBuilder) build() *hugoCmd {
	h := b.newHugoCmd()
	addCommands(h.getCommand(), b.commands...)
	return h
}

func addCommands(root *cobra.Command, commands ...cmder) {
	for _, command := range commands {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		root.AddCommand(cmd)
	}
}

type baseCmd struct {
	cmd *cobra.Command
}

var _ commandsBuilderGetter = (*baseBuilderCmd)(nil)

// Used in tests.
type commandsBuilderGetter interface {
	getCommandsBuilder() *commandsBuilder
}
type baseBuilderCmd struct {
	*baseCmd
	*commandsBuilder
}

func (b *baseBuilderCmd) getCommandsBuilder() *commandsBuilder {
	return b.commandsBuilder
}

func (c *baseCmd) getCommand() *cobra.Command {
	return c.cmd
}

func newBaseCmd(cmd *cobra.Command) *baseCmd {
	return &baseCmd{cmd: cmd}
}

func (b *commandsBuilder) newBuilderCmd(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{commandsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bcmd.hugoBuilderCommon.handleFlags(cmd)
	return bcmd
}

func (c *baseCmd) flagsToConfig(cfg config.Provider) {
	initializeFlags(c.cmd, cfg)
}

type hugoCmd struct {
	*baseBuilderCmd

	// Need to get the sites once built.
	c *commandeer
}

var _ cmder = (*nilCommand)(nil)

type nilCommand struct {
}

func (c *nilCommand) getCommand() *cobra.Command {
	return nil
}

func (c *nilCommand) flagsToConfig(cfg config.Provider) {

}

func (b *commandsBuilder) newHugoCmd() *hugoCmd {
	cc := &hugoCmd{}

	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
		Use:   "hugo",
		Short: "hugo builds your site",
		Long: `hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgInit := func(c *commandeer) error {
				if cc.buildWatch {
					c.Set("disableLiveReload", true)
				}
				return nil
			}

			c, err := initializeConfig(true, cc.buildWatch, &cc.hugoBuilderCommon, cc, cfgInit)
			if err != nil {
				return err
			}
			cc.c = c

			return c.build()
		},
	})

	cc.cmd.PersistentFlags().StringVar(&cc.cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	cc.cmd.PersistentFlags().StringVar(&cc.cfgDir, "configDir", "config", "config dir")
	cc.cmd.PersistentFlags().BoolVar(&cc.quiet, "quiet", false, "build in quiet mode")

	// Set bash-completion
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	_ = cc.cmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

	cc.cmd.PersistentFlags().BoolVarP(&cc.verbose, "verbose", "v", false, "verbose output")
	cc.cmd.PersistentFlags().BoolVarP(&cc.debug, "debug", "", false, "debug output")
	cc.cmd.PersistentFlags().BoolVar(&cc.logging, "log", false, "enable Logging")
	cc.cmd.PersistentFlags().StringVar(&cc.logFile, "logFile", "", "log File path (if set, logging enabled automatically)")
	cc.cmd.PersistentFlags().BoolVar(&cc.verboseLog, "verboseLog", false, "verbose logging")

	cc.cmd.Flags().BoolVarP(&cc.buildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")

	cc.cmd.Flags().Bool("renderToMemory", false, "render to memory (only useful for benchmark testing)")

	// Set bash-completion
	_ = cc.cmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})

	cc.cmd.SetGlobalNormalizationFunc(helpers.NormalizeHugoFlags)
	cc.cmd.SilenceUsage = true

	return cc
}

type hugoBuilderCommon struct {
	source      string
	baseURL     string
	environment string

	buildWatch bool

	gc bool

	// TODO(bep) var vs string
	logging    bool
	verbose    bool
	verboseLog bool
	debug      bool
	quiet      bool

	cfgFile string
	cfgDir  string
	logFile string
}

func (cc *hugoBuilderCommon) getConfigDir(baseDir string) string {
	if cc.cfgDir != "" {
		return paths.AbsPathify(baseDir, cc.cfgDir)
	}

	if v, found := os.LookupEnv("HUGO_CONFIGDIR"); found {
		return paths.AbsPathify(baseDir, v)
	}

	return paths.AbsPathify(baseDir, "config")
}

func (cc *hugoBuilderCommon) getEnvironment(isServer bool) string {
	if cc.environment != "" {
		return cc.environment
	}

	if v, found := os.LookupEnv("HUGO_ENVIRONMENT"); found {
		return v
	}

	if isServer {
		return hugo.EnvironmentDevelopment
	}

	return hugo.EnvironmentProduction
}

func (cc *hugoBuilderCommon) handleFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("cleanDestinationDir", false, "remove files from destination not found in static directories")
	cmd.Flags().BoolP("buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolP("buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().BoolP("buildExpired", "E", false, "include expired content")
	cmd.Flags().StringVarP(&cc.source, "source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().StringVarP(&cc.environment, "environment", "e", "", "build environment")
	cmd.Flags().StringP("contentDir", "c", "", "filesystem path to content directory")
	cmd.Flags().StringP("layoutDir", "l", "", "filesystem path to layout directory")
	cmd.Flags().StringP("cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	cmd.Flags().BoolP("ignoreCache", "", false, "ignores the cache directory")
	cmd.Flags().StringP("destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringP("theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().StringP("themesDir", "", "", "filesystem path to themes directory")
	cmd.Flags().StringVarP(&cc.baseURL, "baseURL", "b", "", "hostname (and path) to the root, e.g. http://spf13.com/")
	cmd.Flags().Bool("enableGitInfo", false, "add Git revision, date and author info to the pages")
	cmd.Flags().BoolVar(&cc.gc, "gc", false, "enable to run some cleanup tasks (remove unused cache files) after the build")

	cmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	cmd.Flags().Bool("templateMetrics", false, "display metrics about template executions")
	cmd.Flags().Bool("templateMetricsHints", false, "calculate some improvement hints when combined with --templateMetrics")
	cmd.Flags().BoolP("forceSyncStatic", "", false, "copy all files when static is changed.")
	cmd.Flags().BoolP("noTimes", "", false, "don't sync modification time of files")
	cmd.Flags().BoolP("noChmod", "", false, "don't sync permission mode of files")
	cmd.Flags().BoolP("i18n-warnings", "", false, "print missing translations")

	cmd.Flags().StringSlice("disableKinds", []string{}, "disable different kind of pages (home, RSS etc.)")

	cmd.Flags().Bool("minify", false, "minify any supported output format (HTML, XML etc.)")

	// Set bash-completion.
	// Each flag must first be defined before using the SetAnnotation() call.
	_ = cmd.Flags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("cacheDir", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})
}

func checkErr(logger *loggers.Logger, err error, s ...string) {
	if err == nil {
		return
	}
	if len(s) == 0 {
		logger.CRITICAL.Println(err)
		return
	}
	for _, message := range s {
		logger.ERROR.Println(message)
	}
	logger.ERROR.Println(err)
}
