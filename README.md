![Hugo](https://raw.githubusercontent.com/gohugoio/hugoDocs/master/static/img/hugo-logo.png)

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

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD, macOS (Darwin), and [Android](https://gist.github.com/bep/a0d8a26cf6b4f8bc992729b8e50b480b) for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including DragonFly BSD, OpenBSD, Plan 9, and Solaris.

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
* [Go (at least Go 1.11)](https://golang.org/dl/)

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
	
## The Hugo Documentation

The Hugo documentation now lives in its own repository, see https://github.com/gohugoio/hugoDocs. But we do keep a version of that documentation as a `git subtree` in this repository. To build the sub folder `/docs` as a Hugo site, you need to clone this repo:

```bash
git clone git@github.com:gohugoio/hugo.git
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

Because we want to create the best possible product for our users and the best contribution experience for our developers,
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
 | [github.com/BurntSushi/locker](https://github.com/BurntSushi/locker) | The Unlicense |
 | [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) | MIT License |
 | [github.com/PuerkitoBio/purell](https://github.com/PuerkitoBio/purell) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/PuerkitoBio/urlesc](https://github.com/PuerkitoBio/urlesc) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/alecthomas/chroma](https://github.com/alecthomas/chroma) |  MIT License |
 | [github.com/bep/debounce](https://github.com/bep/debounce) |    MIT License |
 | [github.com/bep/gitmap](https://github.com/bep/gitmap) |  MIT License |
 | [github.com/bep/go-tocss](https://github.com/bep/go-tocss) | MIT License |
 | [github.com/niklasfasching/go-org](https://github.com/niklasfasching/go-org) | MIT License |
 | [github.com/cpuguy83/go-md2man](https://github.com/cpuguy83/go-md2man) | MIT License |
 | [github.com/danwakefield/fnmatch](https://github.com/danwakefield/fnmatch) | BSD 2-Clause "Simplified" License |
 | [github.com/disintegration/imaging](https://github.com/disintegration/imaging) |  MIT License |
 | [github.com/dlclark/regexp2](https://github.com/dlclark/regexp2) | MIT License |
 | [github.com/eknkc/amber](https://github.com/eknkc/amber) | MIT License |
 | [github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) |  BSD 3-Clause "New" or "Revised" License |
 | [github.com/gobwas/glob](https://github.com/gobwas/glob) | MIT License |
 | [github.com/gorilla/websocket](https://github.com/gorilla/websocket) |  BSD 2-Clause "Simplified" License |
 | [github.com/hashicorp/go-immutable-radix](https://github.com/hashicorp/go-immutable-radix) | Mozilla Public License 2.0 |
 | [github.com/hashicorp/golang-lru](https://github.com/hashicorp/golang-lru) | Mozilla Public License 2.0 |
 | [github.com/hashicorp/hcl](https://github.com/hashicorp/hcl) |   Mozilla Public License 2.0 |
 | [github.com/jdkato/prose](https://github.com/jdkato/prose) |    MIT License |
 | [github.com/kyokomi/emoji](https://github.com/kyokomi/emoji) |   MIT License |
 | [github.com/magiconair/properties](https://github.com/magiconair/properties) | BSD 2-Clause "Simplified" License |
 | [github.com/markbates/inflect](https://github.com/markbates/inflect) |  MIT License |
 | [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty) | MIT License |
 | [github.com/mattn/go-runewidth](https://github.com/mattn/go-runewidth) | MIT License |
 | [github.com/miekg/mmark](https://github.com/miekg/mmark) | Simplified BSD License |
 | [github.com/mitchellh/hashstructure](https://github.com/mitchellh/hashstructure) | MIT License |
 | [github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) | MIT License |
 | [github.com/muesli/smartcrop](https://github.com/muesli/smartcrop) |   MIT License |
 | [github.com/nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n) | MIT License |
 | [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) | MIT License |
 | [github.com/pelletier/go-toml](https://github.com/pelletier/go-toml) |  MIT License |
 | [github.com/pkg/errors](https://github.com/pkg/errors) |  BSD 2-Clause "Simplified" License |
 | [github.com/russross/blackfriday](https://github.com/russross/blackfriday)  | Simplified BSD License |
 | [github.com/shurcooL/sanitized_anchor_name](https://github.com/shurcooL/sanitized_anchor_name) | MIT License |
 | [github.com/spf13/afero](https://github.com/spf13/afero) | Apache License 2.0 |
 | [github.com/spf13/cast](https://github.com/spf13/cast) |  MIT License |
 | [github.com/spf13/cobra](https://github.com/spf13/cobra) | Apache License 2.0 |
 | [github.com/spf13/fsync](https://github.com/spf13/fsync) | MIT License |
 | [github.com/spf13/jwalterweatherman](https://github.com/spf13/jwalterweatherman) | MIT License |
 | [github.com/spf13/nitro](https://github.com/spf13/nitro) | Apache License 2.0 |
 | [github.com/spf13/pflag](https://github.com/spf13/pflag) | BSD 3-Clause "New" or "Revised" License |
 | [github.com/spf13/viper](https://github.com/spf13/viper) | MIT License |
 | [github.com/tdewolff/minify](https://github.com/tdewolff/minify) | MIT License |
 | [github.com/tdewolff/parse](https://github.com/tdewolff/parse) | MIT License |
 | [github.com/wellington/go-libsass](https://github.com/wellington/go-libsass) | Apache License 2.0 |
 | [github.com/yosssi/ace](https://github.com/yosssi/ace) |  MIT License |
 | [golang.org/x/image](https://golang.org/x/image) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/net](https://golang.org/x/net) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/sync](https://golang.org/x/sync) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/sys](https://golang.org/x/sys) | BSD 3-Clause "New" or "Revised" License |
 | [golang.org/x/text](https://golang.org/x/text) | BSD 3-Clause "New" or "Revised" License
 | [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2) |    Apache License 2.0 |

  
 
  
 
 
 
