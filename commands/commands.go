// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"fmt"
	"os"
	"time"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	hpaths "github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cobra"
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
		b.newConfigCmd(),
		b.newDeployCmd(),
		b.newConvertCmd(),
		b.newNewCmd(),
		b.newListCmd(),
		newImportCmd(),
		newGenCmd(),
		createReleaser(),
		b.newModCmd(),
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

func (b *commandsBuilder) newBuilderBasicCmd(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{commandsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bcmd.hugoBuilderCommon.handleCommonBuilderFlags(cmd)
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

type nilCommand struct{}

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

Complete documentation is available at https://gohugo.io/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer cc.timeTrack(time.Now(), "Total")
			cfgInit := func(c *commandeer) error {
				if cc.buildWatch {
					c.Set("disableLiveReload", true)
				}
				return nil
			}

			// prevent cobra printing error so it can be handled here (before the timeTrack prints)
			cmd.SilenceErrors = true

			c, err := initializeConfig(true, true, cc.buildWatch, &cc.hugoBuilderCommon, cc, cfgInit)
			if err != nil {
				cmd.PrintErrln("Error:", err.Error())
				return err
			}
			cc.c = c

			err = c.build()
			if err != nil {
				cmd.PrintErrln("Error:", err.Error())
			}
			return err
		},
	})

	cc.cmd.PersistentFlags().StringVar(&cc.cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	cc.cmd.PersistentFlags().StringVar(&cc.cfgDir, "configDir", "config", "config dir")
	cc.cmd.PersistentFlags().BoolVar(&cc.quiet, "quiet", false, "build in quiet mode")

	// Set bash-completion
	_ = cc.cmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, config.ValidConfigFileExtensions)

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
	poll       string
	clock      string

	gc bool

	// Profile flags (for debugging of performance problems)
	cpuprofile   string
	memprofile   string
	mutexprofile string
	traceprofile string
	printm       bool

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

func (cc *hugoBuilderCommon) timeTrack(start time.Time, name string) {
	if cc.quiet {
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("%s in %v ms\n", name, int(1000*elapsed.Seconds()))
}

func (cc *hugoBuilderCommon) getConfigDir(baseDir string) string {
	if cc.cfgDir != "" {
		return hpaths.AbsPathify(baseDir, cc.cfgDir)
	}

	if v, found := os.LookupEnv("HUGO_CONFIGDIR"); found {
		return hpaths.AbsPathify(baseDir, v)
	}

	return hpaths.AbsPathify(baseDir, "config")
}

func (cc *hugoBuilderCommon) getEnvironment(isServer bool) string {
	if cc.environment != "" {
		return cc.environment
	}

	if v, found := os.LookupEnv("HUGO_ENVIRONMENT"); found {
		return v
	}

	//  Used by Netlify and Forestry
	if v, found := os.LookupEnv("HUGO_ENV"); found {
		return v
	}

	if isServer {
		return hugo.EnvironmentDevelopment
	}

	return hugo.EnvironmentProduction
}

func (cc *hugoBuilderCommon) handleCommonBuilderFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cc.source, "source", "s", "", "filesystem path to read files relative from")
	cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	cmd.PersistentFlags().StringVarP(&cc.environment, "environment", "e", "", "build environment")
	cmd.PersistentFlags().StringP("themesDir", "", "", "filesystem path to themes directory")
	cmd.PersistentFlags().StringP("ignoreVendorPaths", "", "", "ignores any _vendor for module paths matching the given Glob pattern")
	cmd.PersistentFlags().StringVar(&cc.clock, "clock", "", "set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00")
}

func (cc *hugoBuilderCommon) handleFlags(cmd *cobra.Command) {
	cc.handleCommonBuilderFlags(cmd)
	cmd.Flags().Bool("cleanDestinationDir", false, "remove files from destination not found in static directories")
	cmd.Flags().BoolP("buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolP("buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().BoolP("buildExpired", "E", false, "include expired content")
	cmd.Flags().StringP("contentDir", "c", "", "filesystem path to content directory")
	cmd.Flags().StringP("layoutDir", "l", "", "filesystem path to layout directory")
	cmd.Flags().StringP("cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	cmd.Flags().BoolP("ignoreCache", "", false, "ignores the cache directory")
	cmd.Flags().StringP("destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringSliceP("theme", "t", []string{}, "themes to use (located in /themes/THEMENAME/)")
	cmd.Flags().StringVarP(&cc.baseURL, "baseURL", "b", "", "hostname (and path) to the root, e.g. https://spf13.com/")
	cmd.Flags().Bool("enableGitInfo", false, "add Git revision, date, author, and CODEOWNERS info to the pages")
	cmd.Flags().BoolVar(&cc.gc, "gc", false, "enable to run some cleanup tasks (remove unused cache files) after the build")
	cmd.Flags().StringVar(&cc.poll, "poll", "", "set this to a poll interval, e.g --poll 700ms, to use a poll based approach to watch for file system changes")
	cmd.Flags().BoolVar(&loggers.PanicOnWarning, "panicOnWarning", false, "panic on first WARNING log")
	cmd.Flags().Bool("templateMetrics", false, "display metrics about template executions")
	cmd.Flags().Bool("templateMetricsHints", false, "calculate some improvement hints when combined with --templateMetrics")
	cmd.Flags().BoolP("forceSyncStatic", "", false, "copy all files when static is changed.")
	cmd.Flags().BoolP("noTimes", "", false, "don't sync modification time of files")
	cmd.Flags().BoolP("noChmod", "", false, "don't sync permission mode of files")
	cmd.Flags().BoolP("noBuildLock", "", false, "don't create .hugo_build.lock file")
	cmd.Flags().BoolP("printI18nWarnings", "", false, "print missing translations")
	cmd.Flags().BoolP("printPathWarnings", "", false, "print warnings on duplicate target paths etc.")
	cmd.Flags().BoolP("printUnusedTemplates", "", false, "print warnings on unused templates.")
	cmd.Flags().StringVarP(&cc.cpuprofile, "profile-cpu", "", "", "write cpu profile to `file`")
	cmd.Flags().StringVarP(&cc.memprofile, "profile-mem", "", "", "write memory profile to `file`")
	cmd.Flags().BoolVarP(&cc.printm, "printMemoryUsage", "", false, "print memory usage to screen at intervals")
	cmd.Flags().StringVarP(&cc.mutexprofile, "profile-mutex", "", "", "write Mutex profile to `file`")
	cmd.Flags().StringVarP(&cc.traceprofile, "trace", "", "", "write trace to `file` (not useful in general)")

	// Hide these for now.
	cmd.Flags().MarkHidden("profile-cpu")
	cmd.Flags().MarkHidden("profile-mem")
	cmd.Flags().MarkHidden("profile-mutex")

	cmd.Flags().StringSlice("disableKinds", []string{}, "disable different kind of pages (home, RSS etc.)")

	cmd.Flags().Bool("minify", false, "minify any supported output format (HTML, XML etc.)")

	// Set bash-completion.
	// Each flag must first be defined before using the SetAnnotation() call.
	_ = cmd.Flags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("cacheDir", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{})
	_ = cmd.Flags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})
}

func checkErr(logger loggers.Logger, err error, s ...string) {
	if err == nil {
		return
	}
	for _, message := range s {
		logger.Errorln(message)
	}
	logger.Errorln(err)
}
