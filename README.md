![Hugo](https://raw.githubusercontent.com/spf13/hugo/master/docs/static/img/hugo-logo.png)

A Fast and Flexible Static Site Generator built with love by [spf13](http://spf13.com/) and [friends](https://github.com/spf13/hugo/graphs/contributors) in [Go][].

[Website](http://gohugo.io) |
[Forum](https://discuss.gohugo.io) |
[Chat](https://gitter.im/spf13/hugo) |
[Documentation](http://gohugo.io/overview/introduction/) |
[Installation Guide](http://gohugo.io/overview/installing/) |
[Twitter](http://twitter.com/spf13)

[![Build Status](https://travis-ci.org/spf13/hugo.png)](https://travis-ci.org/spf13/hugo) [![wercker status](https://app.wercker.com/status/1a0de7d703ce3b80527f00f675e1eb32 "wercker status")](https://app.wercker.com/project/bykey/1a0de7d703ce3b80527f00f675e1eb32) [![Build status](https://ci.appveyor.com/api/projects/status/n2mo912b8s2505e8/branch/master?svg=true)](https://ci.appveyor.com/project/spf13/hugo/branch/master) [![Join the chat at https://gitter.im/spf13/hugo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/spf13/hugo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Overview

Hugo is a static site generator written in [Go][]. It is optimized for speed, easy use and configurability. Hugo takes a directory with content and templates and renders them into a full HTML website.

Hugo relies on Markdown files with front matter for meta data. And you can run Hugo from any directory. This works well for shared hosts and other systems where you don’t have a privileged account.

Hugo renders a typical website of moderate size in a fraction of a second. A good rule of thumb is that each piece of content renders in around 1 millisecond.

Hugo is meant to work well for any kind of website including blogs, tumbles and docs.

#### Supported Architectures

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD and OS&nbsp;X (Darwin) for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including DragonFly BSD, OpenBSD, Plan&nbsp;9 and Solaris.

**Complete documentation is available at [Hugo Documentation](http://gohugo.io/).**

## Choose How to Install

If you want to use Hugo as your site generator, simply install the Hugo binaries. The Hugo binaries have no external dependencies.

To contribute to the Hugo source code or documentation, you should fork the Hugo GitHub project and clone it to your local machine.

Finally, you can install the Hugo source code with `go`, build the binaries yourself, and run Hugo that way. Building the binaries is an easy task for an experienced `go` getter.

### Install Hugo as Your Site Generator (Binary Install)

Use the [installation instructions in the Hugo documentation](http://gohugo.io/overview/installing/).

### Clone the Hugo Project (Contributor)

1. Make sure your local environment has the following software installed:

    * [Git](http://git-scm.com/)
    * [Mercurial](http://mercurial.selenic.com/)
    * [Go][] 1.3+ (Go 1.4+ on Windows, see Go [Issue #8090](https://code.google.com/p/go/issues/detail?id=8090))

2. Fork the [Hugo project on GitHub](https://github.com/spf13/hugo).

3. Clone your fork:

        git clone https://github.com/YOURNAME/hugo

4. Change into the `hugo` directory:

        cd hugo

5. Install the Hugo project’s package dependencies:

        go get -u -v github.com/spf13/hugo

6. Use a symbolic link to add your locally cloned Hugo repository to your `$GOPATH`, assuming you prefer doing development work outside of `$GOPATH`:

        rm -rf "$GOPATH/src/github.com/spf13/hugo"
        ln -s `pwd` "$GOPATH/src/github.com/spf13/hugo"

    Go expects all of your libraries to be found in`$GOPATH`.

### Build and Install the Binaries from Source (Advanced Install)

Add Hugo and its package dependencies to your go `src` directory.

    go get -v github.com/spf13/hugo

Once the `get` completes, you should find your new `hugo` (or `hugo.exe`) executable sitting inside `$GOPATH/bin/`.

To update Hugo’s dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/spf13/hugo

## Contributing to Hugo

We welcome contributions to Hugo of any kind including documentation, themes, organization, tutorials, blog posts, documentation, bug reports, issues, feature requests, feature implementation, pull requests, answering questions on the forum, helping to manage issues, etc. The Hugo community and maintainers are very active and helpful and the project benefits greatly from this activity.

[![Throughput Graph](https://graphs.waffle.io/spf13/hugo/throughput.svg)](https://waffle.io/spf13/hugo/metrics)

If you have any questions about how to contribute or what to contribute please ask on the [forum](http://discuss.gohugo.io)


## Code Contribution Guide

Contributors should build Hugo and test their changes before submitting a code change.

### Building Hugo with Your Changes

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

### Adding compile information to Hugo

When Hugo is built using the above steps, the `version` sub-command will include the `mdate` of the Hugo executable, similar to the following:

    Hugo Static Site Generator v0.13-DEV buildDate: 2014-12-24T04:46:03-07:00

Instead, it is possible to have the `version` sub-command return information about the git commit used and time of compilation using `build` flags.

To do this, replace the `go build` command with the following *(replace `/path/to/hugo` with the actual path)*:

    go build -ldflags "-X /path/to/hugo/hugolib.CommitHash `git rev-parse --short HEAD 2>/dev/null` -X github.com/spf13/hugo/hugolib.BuildDate `date +%FT%T%z`"

This will result in `hugo version` output that looks similar to:

    Hugo Static Site Generator v0.13-DEV-8042E77 buildDate: 2014-12-25T03:25:57-07:00

The format of the date is configurable via the `Params.DateFormat` setting. `DateFormat` is a string value representing the Go time layout that should be used to format the date output. If `Params.DateFormat` is not set, `time.RFC3339` will be used as the default format. See Go’s ["time" package documentation](http://golang.org/pkg/time/#pkg-constants) for more information.

Configuration setting using config.yaml as example:

    Params:
       DateFormat: "2006-01-02"

Will result in:

    Hugo Static Site Generator v0.13-DEV buildDate: 2014-10-16
    Hugo Static Site Generator v0.13-DEV-24BBFE7 buildDate: 2014-10-16

Alternatively, just run `make` &mdash; all the “magic” above is already in the `Makefile`.  :wink:

### Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugolib
    go run main.go

## Contribution Guidelines

We welcome your contributions. To make the process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes. We encourage pull requests to discuss code changes.
* When you’re ready to create a pull request, be sure to:
     * Have test cases for the new code. If you have questions about how to do it, please ask in your pull request.
     * Run `go fmt`
     * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request.
     * Make sure `go test ./...` passes, and `go build` completes. Our [Travis CI loop](https://travis-ci.org/spf13/hugo) will catch most things that are missing. The exception: Windows. We run on Windows from time to time, but if you have access, please check on a Windows machine too.

**Complete documentation is available at [Hugo Documentation](http://gohugo.io/).**

[![Analytics](https://ga-beacon.appspot.com/UA-7131036-6/hugo/readme)](https://github.com/igrigorik/ga-beacon)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/spf13/hugo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

[Go]: http://golang.org/
