---
aliases:
- /doc/installing/
date: 2013-07-01
menu:
  main:
    parent: getting started
next: /overview/usage
prev: /overview/quickstart
title: Installing Hugo
weight: 20
---

Hugo is written in [Go][] with support for multiple platforms.

The latest release can be found at [Hugo Releases](https://github.com/spf13/hugo/releases).
We currently provide pre-built binaries for
<i class="fa fa-windows"></i>&nbsp;Windows,
<i class="fa fa-linux"></i>&nbsp;Linux,
<i class="fa freebsd-19px"></i>&nbsp;FreeBSD
and <i class="fa fa-apple"></i>&nbsp;OS&nbsp;X (Darwin)
for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including DragonFly BSD, OpenBSD, Plan&nbsp;9 and Solaris.  See http://golang.org/doc/install/source for the full set of supported combinations of target operating systems and compilation architectures.

## Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [Hugo Releases](https://github.com/spf13/hugo/releases).
Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally, you should install it somewhere in your `PATH` for easy use.
`/usr/local/bin` is the most probable location.

On OS&nbsp;X, if you have [Homebrew](http://brew.sh/), installation is even
easier: just run `brew install hugo`.

### Installing Pygments (optional)

The Hugo executable has one *optional* external dependency for source code highlighting (Pygments).

If you want to have source code highlighting using the [highlight shortcode](/extras/highlighting/),
you need to install the Python-based Pygments program. The procedure is outlined on the [Pygments home page](http://pygments.org/).

## Upgrading Hugo

Upgrading Hugo is as easy as downloading and replacing the executable you’ve
placed in your `PATH`.


## Installing from source

### Prerequisite tools for downloading and building source code

* [Git](http://git-scm.com/)
* [Mercurial](http://mercurial.selenic.com/)
* [Go][] 1.3+ (Go 1.4+ on Windows, see Go [Issue #8090](https://code.google.com/p/go/issues/detail?id=8090))

### Get directly from GitHub

    $ export GOPATH=$HOME/go
    $ go get -v github.com/spf13/hugo

`go get` will then fetch Hugo and all its dependent libraries to your
`$GOPATH/src` directory, and compile everything into the final `hugo`
(or `hugo.exe`) executable, which you will find sitting in the
`$GOPATH/bin/` directory, all ready to go!

You may run `go get` with the `-u` option to update Hugo's dependencies:

    $ go get -u -v github.com/spf13/hugo

## Contributing

Please see the [contributing guide](/doc/contributing/).

[Go]: http://golang.org/
