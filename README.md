# Hugo
A Fast and Flexible Static Site Generator built with love by [spf13](http://spf13.com) 
and [friends](http://github.com/spf13/hugo/graphs/contributors) in Go.

[![Build Status](https://travis-ci.org/spf13/hugo.png)](https://travis-ci.org/spf13/hugo)

## Overview

Hugo is a static site generator written in GoLang. It is optimized for
speed, easy use and configurability. Hugo takes a directory with content and
templates and renders them into a full html website.

Hugo makes use of markdown files with front matter for meta data.

A typical website of moderate size can be
rendered in a fraction of a second. A good rule of thumb is that Hugo
takes around 1 millisecond for each piece of content.

It is written to work well with any
kind of website including blogs, tumbles and docs.

**Complete documentation is available at [Hugo Documentation](http://hugo.spf13.com).**

# Getting Started

## Installing Hugo

Hugo is written in GoLang with support for Windows, Linux, FreeBSD and OSX.

The latest release can be found at [hugo releases](https://github.com/spf13/hugo/releases).
We currently build for Windows, Linux, FreeBSD and OS X for x64
and 386 architectures.

### Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [hugo releases](https://github.com/spf13/hugo/releases).
Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally you should install it somewhere in your path for easy use. `/usr/local/bin` 
is the most probable location.

*The Hugo executible has no external dependencies.*

### Installing from source

#### Dependencies

* Git
* Go 1.1+
* Mercurial
* Bazaar

#### Clone locally (for contributors):

    git clone https://github.com/spf13/hugo
    cd hugo
    go get

Because go expects all of your libraries to be found in either $GOROOT or $GOPATH,
it's helpful to symlink the project to one of the following paths:

 * ln -s /path/to/your/hugo $GOPATH/src/github.com/spf13/hugo
 * ln -s /path/to/your/hugo $GOROOT/src/pkg/github.com/spf13/hugo

#### Get directly from Github:

If you only want to build from source, it's even easier.

    go get github.com/spf13/hugo

#### Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

#### Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugolibs
    go run main.go

**Complete documentation is available at [Hugo Documentation](http://hugo.spf13.com).**
