// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose       bool          // print out log statements
		List          bool          // print out a list of targets
		Help          bool          // print out help for a specific target
		Timeout       time.Duration // set a timeout to running the targets
		Args          []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}		
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}		
		d, err := time.ParseDuration(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}
		
	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}
	
	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func (i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
    var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {
		
		targets := map[string]string{
			"check": "Run tests and linters",
			"docker": "Build hugo Docker container",
			"fmt": "Run gofmt linter",
			"genDocsHelper": "Generate docs helper",
			"generate": "autogen packages",
			"hugo": "Build hugo binary",
			"hugoNoGitInfo": "Build hugo without git info",
			"hugoRace": "Build hugo binary with race detector enabled",
			"install": "hugo binary",
			"lint": "Run golint linter",
			"test": "Run tests",
			"test386": "Run tests in 32-bit mode Note that we don't run with the extended tag.",
			"testCoverHTML": "Generate test coverage report",
			"testRace": "Run tests with race detector",
			"uninstall": "hugo binary",
			"vet": "Run go vet linter",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		fmt.Println("Targets:")
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
		return err
	}

	var ctx context.Context
	var ctxCancel func()

	getContext := func() (context.Context, func()) {
		if ctx != nil {
			return ctx, ctxCancel
		}

		if args.Timeout != 0 {
			ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
		} else {
			ctx = context.Background()
			ctxCancel = func() {}
		}
		return ctx, ctxCancel
	}

	runTarget := func(fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		select {
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			cancel()
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	log.SetFlags(0)
	if !args.Verbose {
		log.SetOutput(ioutil.Discard)
	}
	logger := log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch strings.ToLower(args.Args[0]) {
			case "check":
				fmt.Println("Run tests and linters")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage check\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "docker":
				fmt.Println("Build hugo Docker container")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage docker\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "fmt":
				fmt.Println("Run gofmt linter")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage fmt\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "gendocshelper":
				fmt.Println("Generate docs helper")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage gendocshelper\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "generate":
				fmt.Println("Generate autogen packages")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage generate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "hugo":
				fmt.Println("Build hugo binary")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage hugo\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "hugonogitinfo":
				fmt.Println("Build hugo without git info")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage hugonogitinfo\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "hugorace":
				fmt.Println("Build hugo binary with race detector enabled")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage hugorace\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "install":
				fmt.Println("Install hugo binary")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage install\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "lint":
				fmt.Println("Run golint linter")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage lint\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "test":
				fmt.Println("Run tests")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage test\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "test386":
				fmt.Println("Run tests in 32-bit mode Note that we don't run with the extended tag. Currently not supported in 32 bit.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage test386\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "testcoverhtml":
				fmt.Println("Generate test coverage report")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage testcoverhtml\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "testrace":
				fmt.Println("Run tests with race detector")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage testrace\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "uninstall":
				fmt.Println("Uninstall hugo binary")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage uninstall\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "vet":
				fmt.Println("Run go vet linter")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage vet\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			
			default:
				logger.Printf("Unknown target: %q\n", args.Args[0])
				os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		if err := list(); err != nil {
			logger.Println("Error:", err)
			os.Exit(1)
		}
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch strings.ToLower(target) {
		
		}

		switch strings.ToLower(target) {
		
			case "check":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Check\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Check")
				}
				
				wrapFn := func(ctx context.Context) error {
					Check()
					return nil
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "docker":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "fmt":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Fmt\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Fmt")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Fmt()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "gendocshelper":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"GenDocsHelper\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "GenDocsHelper")
				}
				
				wrapFn := func(ctx context.Context) error {
					return GenDocsHelper()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "generate":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generate")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generate()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "hugo":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Hugo\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Hugo")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Hugo()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "hugonogitinfo":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"HugoNoGitInfo\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "HugoNoGitInfo")
				}
				
				wrapFn := func(ctx context.Context) error {
					return HugoNoGitInfo()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "hugorace":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"HugoRace\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "HugoRace")
				}
				
				wrapFn := func(ctx context.Context) error {
					return HugoRace()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "install":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Install\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Install")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Install()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "lint":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Lint\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Lint")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Lint()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "test":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "test386":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test386\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test386")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test386()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "testcoverhtml":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"TestCoverHTML\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "TestCoverHTML")
				}
				
				wrapFn := func(ctx context.Context) error {
					return TestCoverHTML()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "testrace":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"TestRace\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "TestRace")
				}
				
				wrapFn := func(ctx context.Context) error {
					return TestRace()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "uninstall":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Uninstall\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Uninstall")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Uninstall()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "vet":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Vet\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Vet")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Vet()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
		
		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}




