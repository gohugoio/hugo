<img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/static/images/hugo-logo-wide.svg?sanitize=true" alt="Hugo" width="565">

A Fast and Flexible Static Site Generator built with love by [bep](https://github.com/bep), [spf13](http://spf13.com/) and [friends](https://github.com/gohugoio/hugo/graphs/contributors) in [Go][].

[Website](https://gohugo.io) |
[Forum](https://discourse.gohugo.io) |
[Documentation](https://gohugo.io/getting-started/) |
[Installation Guide](https://gohugo.io/getting-started/installing/) |
[Contribution Guide](CONTRIBUTING.md) |
[Twitter](https://twitter.com/gohugoio)

[![GoDoc](https://godoc.org/github.com/gohugoio/hugo?status.svg)](https://godoc.org/github.com/gohugoio/hugo)
[![Linux and macOS Build Status](https://api.travis-ci.org/gohugoio/hugo.svg?branch=master&label=Windows+and+Linux+and+macOS+build "Windows, Linux and macOS Build Status")](https://travis-ci.org/gohugoio/hugo)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohugoio/hugo)](https://goreportcard.com/report/github.com/gohugoio/hugo)

## Overview

Hugo is a static HTML and CSS website generator written in [Go][].
It is optimized for speed, ease of use, and configurability.
Hugo takes a directory with content and templates and renders them into a full HTML website.

Hugo relies on Markdown files with front matter for metadata, and you can run Hugo from any directory.
This works well for shared hosts and other systems where you don’t have a privileged account.

Hugo renders a typical website of moderate size in a fraction of a second.
A good rule of thumb is that each piece of content renders in around 1 millisecond.

Hugo is designed to work well for any kind of website including blogs, tumbles, and docs.

#### Supported Architectures

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD, DragonFly BSD, Open BSD, macOS (Darwin), and [Android](https://gist.github.com/bep/a0d8a26cf6b4f8bc992729b8e50b480b) for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including Plan 9 and Solaris.

**Complete documentation is available at [Hugo Documentation](https://gohugo.io/getting-started/).**

## Choose How to Install

If you want to use Hugo as your site generator, simply install the Hugo binaries.
The Hugo binaries have no external dependencies.

To contribute to the Hugo source code or documentation, you should [fork the Hugo GitHub project](https://github.com/gohugoio/hugo#fork-destination-box) and clone it to your local machine.

Finally, you can install the Hugo source code with `go`, build the binaries yourself, and run Hugo that way.
Building the binaries is an easy task for an experienced `go` getter.

### Install Hugo as Your Site Generator (Binary Install)

Use the [installation instructions in the Hugo documentation](https://gohugo.io/getting-started/installing/).

### Build and Install the Binaries from Source (Advanced Install)

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (we test it with the last 2 major versions)](https://golang.org/dl/)

#### Fetch from GitHub

Since Hugo 0.48, Hugo uses the Go Modules support built into Go 1.11 to build. The easiest is to clone Hugo in a directory outside of `GOPATH`, as in the following example:

```bash
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/gohugoio/hugo.git
cd hugo
go install
```

**If you are a Windows user, substitute the `$HOME` environment variable above with `%USERPROFILE%`.**

If you want to compile with Sass/SCSS support use `--tags extended` and make sure `CGO_ENABLED=1` is set in your go environment. If you don't want to have CGO enabled, you may use the following command to temporarily enable CGO only for hugo compilation:

```bash
CGO_ENABLED=1 go install --tags extended
```

## The Hugo Documentation

The Hugo documentation now lives in its own repository, see https://github.com/gohugoio/hugoDocs. But we do keep a version of that documentation as a `git subtree` in this repository. To build the sub folder `/docs` as a Hugo site, you need to clone this repo:

```bash
git clone git@github.com:gohugoio/hugo.git
```
## Basic Usage

Hugo’s CLI is fully featured but simple to use, even for those who have very limited experience working from the command line.

The following is a description of the most common commands you will use while developing your Hugo project. See the Command Line Reference for a comprehensive view of Hugo’s CLI.

Test Installation 
Once you have installed Hugo, make sure it is in your PATH. You can test that Hugo has been installed correctly via the help command:

hugo help
The output you see in your console should be similar to the following:

hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at https://gohugo.io/.

Usage:
  hugo [flags]
  hugo [command]

Available Commands:
  check       Contains some verification checks
  config      Print the site configuration
  convert     Convert your content to different formats
  env         Print Hugo version and environment info
  gen         A collection of several useful generators.
  help        Help about any command
  import      Import your site from others.
  list        Listing out various types of content
  new         Create new content for your site
  server      A high performance webserver
  version     Print the version number of Hugo

Flags:
  -b, --baseURL string         hostname (and path) to the root, e.g. https://spf13.com/
  -D, --buildDrafts            include content marked as draft
  -E, --buildExpired           include expired content
  -F, --buildFuture            include content with publishdate in the future
      --cacheDir string        filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/
      --cleanDestinationDir    remove files from destination not found in static directories
      --config string          config file (default is path/config.yaml|json|toml)
      --configDir string       config dir (default "config")
  -c, --contentDir string      filesystem path to content directory
      --debug                  debug output
  -d, --destination string     filesystem path to write files to
      --disableKinds strings   disable different kind of pages (home, RSS etc.)
      --enableGitInfo          add Git revision, date and author info to the pages
  -e, --environment string     build environment
      --forceSyncStatic        copy all files when static is changed.
      --gc                     enable to run some cleanup tasks (remove unused cache files) after the build
  -h, --help                   help for hugo
      --i18n-warnings          print missing translations
      --ignoreCache            ignores the cache directory
  -l, --layoutDir string       filesystem path to layout directory
      --log                    enable Logging
      --logFile string         log File path (if set, logging enabled automatically)
      --minify                 minify any supported output format (HTML, XML etc.)
      --noChmod                don't sync permission mode of files
      --noTimes                don't sync modification time of files
      --path-warnings          print warnings on duplicate target paths etc.
      --quiet                  build in quiet mode
      --renderToMemory         render to memory (only useful for benchmark testing)
  -s, --source string          filesystem path to read files relative from
      --templateMetrics        display metrics about template executions
      --templateMetricsHints   calculate some improvement hints when combined with --templateMetrics
  -t, --theme strings          themes to use (located in /themes/THEMENAME/)
      --themesDir string       filesystem path to themes directory
      --trace file             write trace to file (not useful in general)
  -v, --verbose                verbose output
      --verboseLog             verbose logging
  -w, --watch                  watch filesystem for changes and recreate as needed

Use "hugo [command] --help" for more information about a command.
The hugo Command 
The most common usage is probably to run hugo with your current directory being the input directory.

This generates your website to the public/ directory by default, although you can customize the output directory in your site configuration by changing the publishDir field.

The command hugo renders your site into public/ dir and is ready to be deployed to your web server:

hugo
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 90 ms
Draft, Future, and Expired Content 
Hugo allows you to set draft, publishdate, and even expirydate in your content’s front matter. By default, Hugo will not publish:

Content with a future publishdate value
Content with draft: true status
Content with a past expirydate value
All three of these can be overridden during both local development and deployment by adding the following flags to hugo and hugo server, respectively, or by changing the boolean values assigned to the fields of the same name (without --) in your configuration:

--buildFuture
--buildDrafts
--buildExpired
LiveReload 
Hugo comes with LiveReload built in. There are no additional packages to install. A common way to use Hugo while developing a site is to have Hugo run a server with the hugo server command and watch for changes:

hugo server
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
Watching for changes in /Users/yourname/sites/yourhugosite/{data,content,layouts,static}
Serving pages from /Users/yourname/sites/yourhugosite/public
Web Server is available at http://localhost:1313/
Press Ctrl+C to stop
This will run a fully functioning web server while simultaneously watching your file system for additions, deletions, or changes within the following areas of your project organization:

/static/*
/content/*
/data/*
/i18n/*
/layouts/*
/themes/<CURRENT-THEME>/*
config
Whenever you make changes, Hugo will simultaneously rebuild the site and continue to serve content. As soon as the build is finished, LiveReload tells the browser to silently reload the page.

Most Hugo builds are so fast that you may not notice the change unless looking directly at the site in your browser. This means that keeping the site open on a second monitor (or another half of your current monitor) allows you to see the most up-to-date version of your website without the need to leave your text editor.

Hugo injects the LiveReload <script> before the closing </body> in your templates and will therefore not work if this tag is not present..

Disable LiveReload 
LiveReload works by injecting JavaScript into the pages Hugo generates. The script creates a connection from the browser’s web socket client to the Hugo web socket server.

LiveReload is awesome for development. However, some Hugo users may use hugo server in production to instantly display updated content. The following methods make it easy to disable LiveReload:

hugo server --watch=false
Or…

hugo server --disableLiveReload
The latter flag can be omitted by adding the following key-value to your config.toml or config.yml file, respectively:

disableLiveReload = true
disableLiveReload: true
Deploy Your Website 
After running hugo server for local web development, you need to do a final hugo run without the server part of the command to rebuild your site. You may then deploy your site by copying the public/ directory to your production web server.

Since Hugo generates a static website, your site can be hosted anywhere using any web server. See Hosting and Deployment for methods for hosting and automating deployments contributed by the Hugo community.

Running hugo does not remove generated files before building. This means that you should delete your public/ directory (or the publish directory you specified via flag or configuration file) before running the hugo command. If you do not remove these files, you run the risk of the wrong files (e.g., drafts or future posts) being left in the generated site.

See Also
```
```
## Contributing to Hugo

For a complete guide to contributing to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

We welcome contributions to Hugo of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Hugo community and maintainers are [very active](https://github.com/gohugoio/hugo/pulse/monthly) and helpful, and the project benefits greatly from this activity.

### Asking Support Questions

We have an active [discussion forum](https://discourse.gohugo.io) where users and developers can ask questions.
Please don't use the GitHub issue tracker to ask questions.

### Reporting Issues

If you believe you have found a defect in Hugo or its documentation, use
the GitHub issue tracker to report the problem to the Hugo maintainers.
If you're not sure if it's a bug or not, start by asking in the [discussion forum](https://discourse.gohugo.io).
When reporting the issue, please provide the version of Hugo in use (`hugo version`).

### Submitting Patches

The Hugo project welcomes all contributors and contributions regardless of skill or experience level.
If you are interested in helping with the project, we will help you with your contribution.
Hugo is a very active project with many contributions happening daily.

We want to create the best possible product for our users and the best contribution experience for our developers,
we have a set of guidelines which ensure that all contributions are acceptable.
The guidelines are not intended as a filter or barrier to participation.
If you are unfamiliar with the contribution process, the Hugo team will help you and teach you how to bring your contribution in accordance with the guidelines.

For a complete guide to contributing code to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

[![Analytics](https://ga-beacon.appspot.com/UA-7131036-6/hugo/readme)](https://github.com/igrigorik/ga-beacon)

[Go]: https://golang.org/
[Hugo Documentation]: https://gohugo.io/overview/introduction/

## Dependencies

Hugo stands on the shoulder of many great open source libraries, in lexical order:

 | Dependency  | License |
 | :------------- | :------------- |
 | [github.com/alecthomas/chroma](https://github.com/alecthomas/chroma) | MIT License |
 | [github.com/armon/go-radix](https://github.com/armon/go-radix) | MIT License |
 | [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) | Apache License 2.0 |
 | [github.com/bep/debounce](https://github.com/bep/debounce) | MIT License |
 | [github.com/bep/gitmap](https://github.com/bep/gitmap) | MIT License |
 | [github.com/bep/golibsass](https://github.com/bep/golibsass) | MIT License |
 | [github.com/bep/tmc](https://github.com/bep/tmc) | MIT License |
 | [github.com/BurntSushi/locker](https://github.com/BurntSushi/locker) | The Unlicense |
 | [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) | MIT License |
 | [github.com/cpuguy83/go-md2man](https://github.com/cpuguy83/go-md2man) | MIT License |
 | [github.com/danwakefield/fnmatch](https://github.com/danwakefield/fnmatch) | BSD 2-Clause "Simplified" License |
 | [github.com/disintegration/gift](https://github.com/disintegration/gift) | MIT License |
 | [github.com/dustin/go-humanize](https://github.com/dustin/go-humanize) | MIT License |
 | [github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/gobwas/glob](https://github.com/gobwas/glob) | MIT License |
 | [github.com/gorilla/websocket](https://github.com/gorilla/websocket) | BSD 2-Clause "Simplified" License |
 | [github.com/hashicorp/golang-lru](https://github.com/hashicorp/golang-lru) | Mozilla Public License 2.0 |
 | [github.com/hashicorp/hcl](https://github.com/hashicorp/hcl) | Mozilla Public License 2.0 |
 | [github.com/jdkato/prose](https://github.com/jdkato/prose) | MIT License |
 | [github.com/kr/pretty](https://github.com/kr/pretty) | MIT License |
 | [github.com/kyokomi/emoji](https://github.com/kyokomi/emoji) | MIT License |
 | [github.com/magiconair/properties](https://github.com/magiconair/properties) | BSD 2-Clause "Simplified" License |
 | [github.com/markbates/inflect](https://github.com/markbates/inflect) | MIT License |
 | [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty) | MIT License |
 | [github.com/mattn/go-runewidth](https://github.com/mattn/go-runewidth) | MIT License |
 | [github.com/miekg/mmark](https://github.com/miekg/mmark) | Simplified BSD License |
 | [github.com/mitchellh/hashstructure](https://github.com/mitchellh/hashstructure) | MIT License |
 | [github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) | MIT License |
 | [github.com/muesli/smartcrop](https://github.com/muesli/smartcrop) | MIT License |
 | [github.com/nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n) | MIT License |
 | [github.com/niklasfasching/go-org](https://github.com/niklasfasching/go-org) | MIT License |
 | [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) | MIT License |
 | [github.com/pelletier/go-toml](https://github.com/pelletier/go-toml) | MIT License |
 | [github.com/pkg/errors](https://github.com/pkg/errors) | BSD 2-Clause "Simplified" License |
 | [github.com/PuerkitoBio/purell](https://github.com/PuerkitoBio/purell) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/PuerkitoBio/urlesc](https://github.com/PuerkitoBio/urlesc) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/rogpeppe/go-internal](https://github.com/rogpeppe/go-internal) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/russross/blackfriday](https://github.com/russross/blackfriday)  | Simplified BSD License |
 | [github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif) | BSD 2-Clause "Simplified" License |
 | [github.com/spf13/afero](https://github.com/spf13/afero) | Apache License 2.0 |
 | [github.com/spf13/cast](https://github.com/spf13/cast) | MIT License |
 | [github.com/spf13/cobra](https://github.com/spf13/cobra) | Apache License 2.0 |
 | [github.com/spf13/fsync](https://github.com/spf13/fsync) | MIT License |
 | [github.com/spf13/jwalterweatherman](https://github.com/spf13/jwalterweatherman) | MIT License |
 | [github.com/spf13/pflag](https://github.com/spf13/pflag) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/spf13/viper](https://github.com/spf13/viper) | MIT License |
 | [github.com/tdewolff/minify](https://github.com/tdewolff/minify) | MIT License |
 | [github.com/tdewolff/parse](https://github.com/tdewolff/parse) | MIT License |
 | [github.com/yuin/goldmark](https://github.com/yuin/goldmark) | MIT License |
 | [github.com/yuin/goldmark-highlighting](https://github.com/yuin/goldmark-highlighting) | MIT License |
 | [go.opencensus.io](https://go.opencensus.io) | Apache License 2.0 |
 | [go.uber.org/atomic](https://go.uber.org/atomic) | MIT License |
 | [gocloud.dev](https://gocloud.dev) | Apache License 2.0 |
 | [golang.org/x/image](https://golang.org/x/image) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/net](https://golang.org/x/net) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/oauth2](https://golang.org/x/oauth2) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/sync](https://golang.org/x/sync) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/sys](https://golang.org/x/sys) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/text](https://golang.org/x/text) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/xerrors](https://golang.org/x/xerrors) | BSD 3-Clause "New" or "Revised" License |
 | [google.golang.org/api](https://google.golang.org/api) | BSD 3-Clause "New" or "Revised" License |
 | [google.golang.org/genproto](https://google.golang.org/genproto) | Apache License 2.0 |
 | [gopkg.in/ini.v1](https://gopkg.in/ini.v1) | Apache License 2.0 |
 | [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2) | Apache License 2.0 |
