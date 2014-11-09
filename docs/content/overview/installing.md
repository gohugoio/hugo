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

Hugo is written in Go with support for Windows, Linux, FreeBSD and OS&nbsp;X.

The latest release can be found at [Hugo Releases](https://github.com/spf13/hugo/releases).
We currently build for Windows, Linux, FreeBSD and OS&nbsp;X for x64
and i386 architectures.

## Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [Hugo Releases](https://github.com/spf13/hugo/releases).
Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally you should install it somewhere in your path for easy use. `/usr/local/bin`
is the most probable location.

If you have [Homebrew](http://brew.sh), installation is even easier.  Just run
`brew install hugo`.

### Installing Pygments (optional)

The Hugo executable has one *optional* external dependency for source code highlighting (Pygments).

If you want to have source code highlighting using the [highlight shortcode](/extras/highlighting),
you need to install the Python-based Pygments program. The procedure is outlined on the [Pygments home page](http://pygments.org).

## Upgrading Hugo

Upgrading Hugo is as easy as downloading and replacing the executable you’ve
placed in your path.


## Installing from source

### Dependencies

* Git
* Go 1.1+
* Mercurial
* Bazaar

### Get directly from GitHub:

    go get -v github.com/spf13/hugo

### Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

## Contributing

Please see the [contributing guide](/doc/contributing).
