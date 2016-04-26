![Hugo](https://raw.githubusercontent.com/spf13/hugo/master/docs/static/img/hugo-logo.png)

A Fast and Flexible Static Site Generator built with love by [spf13](http://spf13.com/) and [friends](https://github.com/spf13/hugo/graphs/contributors) in [Go][].

[Website](https://gohugo.io) |
[Forum](https://discuss.gohugo.io) |
[Dev Chat](https://gitter.im/spf13/hugo) |
[Documentation](https://gohugo.io/overview/introduction/) |
[Installation Guide](https://gohugo.io/overview/installing/) |
[Twitter](http://twitter.com/spf13)

[![GoDoc](https://godoc.org/github.com/spf13/hugo?status.svg)](https://godoc.org/github.com/spf13/hugo)
[![Linux and OS X Build Status](https://api.travis-ci.org/spf13/hugo.svg?branch=master&label=Linux+and+OS+X+build "Linux and OS X Build Status")](https://travis-ci.org/spf13/hugo)
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/n2mo912b8s2505e8/branch/master?svg=true&label=Windows+build "Windows Build Status")](https://ci.appveyor.com/project/spf13/hugo/branch/master)
[![Dev chat at https://gitter.im/spf13/hugo](https://img.shields.io/badge/gitter-dev_chat-46bc99.svg)](https://gitter.im/spf13/hugo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/spf13/hugo)](https://goreportcard.com/report/github.com/spf13/hugo)

## Overview

Hugo is a static HTML and CSS website generator written in [Go][].
It is optimized for speed, easy use and configurability.
Hugo takes a directory with content and templates and renders them into a full HTML website.

Hugo relies on Markdown files with front matter for meta data.
And you can run Hugo from any directory.
This works well for shared hosts and other systems where you don’t have a privileged account.

Hugo renders a typical website of moderate size in a fraction of a second.
A good rule of thumb is that each piece of content renders in around 1 millisecond.

Hugo is designed to work well for any kind of website including blogs, tumbles and docs.

#### Supported Architectures

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD and OS&nbsp;X (Darwin) for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including DragonFly BSD, OpenBSD, Plan&nbsp;9 and Solaris.

**Complete documentation is available at [Hugo Documentation][].**

## Choose How to Install

If you want to use Hugo as your site generator, simply install the Hugo binaries.
The Hugo binaries have no external dependencies.

To contribute to the Hugo source code or documentation, you should [fork the Hugo GitHub project](https://github.com/spf13/hugo#fork-destination-box) and clone it to your local machine.

Finally, you can install the Hugo source code with `go`, build the binaries yourself, and run Hugo that way.
Building the binaries is an easy task for an experienced `go` getter.

### Install Hugo as Your Site Generator (Binary Install)

Use the [installation instructions in the Hugo documentation](https://gohugo.io/overview/installing/).

### Clone the Hugo Project (Contributor)

1. Make sure your local environment has the following software installed:

    * [Git](https://git-scm.com/)
    * [Go][] 1.5+

2. [Fork the Hugo project on GitHub](https://github.com/spf13/hugo).

3. Clone your fork:

        git clone https://github.com/YOURNAME/hugo

4. Change into the `hugo` directory:

        cd hugo

5. Install the Hugo project’s package dependencies:

        go get -u -v github.com/spf13/hugo

6. Install the test dependencies (needed if you want to run tests):

        go get -v -t -d ./...

7. Use a symbolic link to add your locally cloned Hugo repository to your `$GOPATH`, assuming you prefer doing development work outside of `$GOPATH`:

    ```bash
    rm -rf "$GOPATH/src/github.com/spf13/hugo"
    ln -s `pwd` "$GOPATH/src/github.com/spf13/hugo"
    ```

    Go expects all of your libraries to be found in`$GOPATH`.

You can also find a [detailed guide](https://www.gohugo.io/tutorials/how-to-contribute-to-hugo/) in our documentation.

### Build and Install the Binaries from Source (Advanced Install)

Add Hugo and its package dependencies to your go `src` directory.

    go get -v github.com/spf13/hugo

Once the `get` completes, you should find your new `hugo` (or `hugo.exe`) executable sitting inside `$GOPATH/bin/`.

To update Hugo’s dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/spf13/hugo

## Contribute to Hugo

We welcome contributions to Hugo of any kind including documentation, themes, organization, tutorials, blog posts, bug reports, issues, feature requests, feature implementation, pull requests, answering questions on the forum, helping to manage issues, etc.

The Hugo community and maintainers are very active and helpful and the project benefits greatly from this activity.

[![Throughput Graph](https://graphs.waffle.io/spf13/hugo/throughput.svg)](https://waffle.io/spf13/hugo/metrics)

If you have any questions about how to contribute or what to contribute, please ask on the [forum](https://discuss.gohugo.io).

## Code Contribution Guideline

We welcome your contributions.
To make the process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes. We encourage pull requests to discuss code changes.
* When you’re ready to create a pull request, be sure to:
     * Sign the [CLA](https://cla-assistant.io/spf13/hugo)
     * Have test cases for the new code. If you have questions about how to do it, please ask in your pull request.
     * Run `go fmt`
     * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request.
     * **Write a good commit message.** This [blog article](http://chris.beams.io/posts/git-commit) is a good resource for learning how to write good commit messages, the most important part being that each commit message should have a title/subject in imperative mood starting with a capital letter and no trailing period: *"Return error on wrong use of the Paginator"*, **NOT** *"returning some error."* Also, if your commit references one or more GitHub issues, always end your commit message body with *See #1234* or *Fixes #1234*. Replace *1234* with the GitHub issue ID. The last example will close the issue when the commit is merged into *master*. Sometimes it makes sense to prefix the commit message with the packagename (or docs folder) all lowercased ending with a colon. That is fine, but the rest of the rules above apply. So it is "tpl: Add emojify template func", not "tpl: add emojify template func.", and "docs: Document emoji", not "doc: document emoji."
     * Make sure `go test ./...` passes, and `go build` completes. Our [Travis CI loop](https://travis-ci.org/spf13/hugo) (Linux and OS&nbsp;X) and [AppVeyor](https://ci.appveyor.com/project/spf13/hugo/branch/master) (Windows) will catch most things that are missing.

### Build Hugo with Your Changes

```bash
cd /path/to/hugo
go build -o hugo main.go
mv hugo /usr/local/bin/
```

### Add Compile Information to Hugo

To add compile information to Hugo, replace the `go build` command with the following *(replace `/path/to/hugo` with the actual path)*:

    go build -ldflags "-X /path/to/hugo/hugolib.CommitHash=`git rev-parse --short HEAD 2>/dev/null` -X github.com/spf13/hugo/hugolib.BuildDate=`date +%FT%T%z`"

This will result in `hugo version` output that looks similar to:

    Hugo Static Site Generator v0.13-DEV-8042E77 buildDate: 2014-12-25T03:25:57-07:00

Alternatively, just run `make` &mdash; all the “magic” above is already in the `Makefile`.  :wink:

### Run Hugo

```bash
cd /path/to/hugo
go install github.com/spf13/hugo/hugolib
go run main.go
```

**Complete documentation is available at [Hugo Documentation][].**

[![Analytics](https://ga-beacon.appspot.com/UA-7131036-6/hugo/readme)](https://github.com/igrigorik/ga-beacon)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/spf13/hugo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

[Go]: https://golang.org/
[Hugo Documentation]: https://gohugo.io/overview/introduction/
