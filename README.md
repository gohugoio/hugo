# Hugo
A Fast and Flexible Static Site Generator built with love by [spf13](http://spf13.com) 
and [friends](http://github.com/spf13/hugo/graphs/contributors) in Go.

[![Build Status](https://travis-ci.org/spf13/hugo.png)](https://travis-ci.org/spf13/hugo)
[![wercker status](https://app.wercker.com/status/1a0de7d703ce3b80527f00f675e1eb32 "wercker status")](https://app.wercker.com/project/bykey/1a0de7d703ce3b80527f00f675e1eb32)
[![Build status](https://ci.appveyor.com/api/projects/status/n2mo912b8s2505e8/branch/master?svg=true)](https://ci.appveyor.com/project/spf13/hugo/branch/master)

## Overview

Hugo is a static site generator written in Go. It is optimized for
speed, easy use and configurability. Hugo takes a directory with content and
templates and renders them into a full HTML website.

Hugo makes use of Markdown files with front matter for meta data.

A typical website of moderate size can be
rendered in a fraction of a second. A good rule of thumb is that Hugo
takes around 1 millisecond for each piece of content.

It is written to work well with any
kind of website including blogs, tumbles and docs.

**Complete documentation is available at [Hugo Documentation](http://gohugo.io).**

# Getting Started

## Installing Hugo

Hugo is written in Go with support for Windows, Linux, FreeBSD and OS X.

The latest release can be found at [hugo releases](https://github.com/spf13/hugo/releases).
We currently build for Windows, Linux, FreeBSD and OS X for x64
and i386 architectures.

### Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [Hugo Releases](https://github.com/spf13/hugo/releases).
Once downloaded, it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally, you should install it somewhere in your path for easy use. `/usr/local/bin`
is the most probable location.

*The Hugo executable has no external dependencies.*

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

Because Go expects all of your libraries to be found in either $GOROOT or $GOPATH,
it's helpful to symlink the project to one of the following paths:

 * `ln -s /path/to/your/hugo $GOPATH/src/github.com/spf13/hugo`
 * `ln -s /path/to/your/hugo $GOROOT/src/pkg/github.com/spf13/hugo`

#### Get directly from GitHub:

If you only want to build from source, it's even easier.

    go get -v github.com/spf13/hugo

#### Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

##### Adding compile information to Hugo

When Hugo is built using the above steps, the `version` sub-command will include the `mdate` of the Hugo executable.  Instead, it is possible to have the `version` sub-command return information about the git commit used and time of compilation using `build` flags.

To do this, replace the `go build` command with the following *(replace `/path/to/hugo` with the actual path)*:

    go build -ldflags "-X /path/to/hugo/commands.commitHash `git rev-parse --short HEAD 2>/dev/null` -X github.com/spf13/hugo/commands.buildDate `date +%FT%T`"  

This will result in hugo version output that looks similar to:

    Hugo Static Site Generator v0.13-DEV buildDate: 2014-10-16T09:59:55Z
    Hugo Static Site Generator v0.13-DEV-24BBFE7 buildDate: 2014-10-16T10:00:55Z

The format of the date is configurable via the `Params.DateFormat` setting.  `DateFormat` is a string value representing the Go time layout that should be used to format the date output. If `Params.DateFormat` is not set, `time.RFC3339` will be used as the default format.See [time documentation](http://golang.org/pkg/time/#pkg-constants) for more information.

Configuration setting using config.yaml as example:

    Params:
       DateFormat: "2006-01-02"

Will result in:

    Hugo Static Site Generator v0.13-DEV buildDate: 2014-10-16
    Hugo Static Site Generator v0.13-DEV-24BBFE7 buildDate: 2014-10-16

#### Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugolib
    go run main.go

#### Contribution Guidelines

We welcome your contributions.  To make the process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes.  We encourage pull requests to discuss code changes.
* When you're ready to create a pull request, be sure to:
     * Have test cases for the new code.  If you have questions about how to do it, please ask in your pull request.
     * Run `go fmt`
     * Squash your commits into a single commit.  `git rebase -i`.  It's okay to force update your pull request.  
     * Make sure `go test ./...` passes, and go build completes.  Our Travis CI loop will catch most things that are missing.  The exception: Windows.  We run on Windows from time to time, but if you have access, please check on a Windows machine too.

**Complete documentation is available at [Hugo Documentation](http://gohugo.io).**

[![Analytics](https://ga-beacon.appspot.com/UA-7131036-6/hugo/readme)](https://github.com/igrigorik/ga-beacon)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/spf13/hugo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
