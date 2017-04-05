![Hugo](https://raw.githubusercontent.com/spf13/hugo/master/docs/static/img/hugo-logo.png)

A Fast and Flexible Static Site Generator built with love by [spf13](http://spf13.com/) and [friends](https://github.com/spf13/hugo/graphs/contributors) in [Go][].

[Website](https://gohugo.io) |
[Forum](https://discuss.gohugo.io) |
[Developer Chat (no support)](https://gitter.im/spf13/hugo) |
[Documentation](https://gohugo.io/overview/introduction/) |
[Installation Guide](https://gohugo.io/overview/installing/) |
[Contribution Guide](CONTRIBUTING.md) |
[Twitter](http://twitter.com/gohugoio)

[![GoDoc](https://godoc.org/github.com/spf13/hugo?status.svg)](https://godoc.org/github.com/spf13/hugo)
[![Linux and macOS Build Status](https://api.travis-ci.org/spf13/hugo.svg?branch=master&label=Linux+and+macOS+build "Linux and macOS Build Status")](https://travis-ci.org/spf13/hugo)
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/n2mo912b8s2505e8/branch/master?svg=true&label=Windows+build "Windows Build Status")](https://ci.appveyor.com/project/spf13/hugo/branch/master)
[![Dev chat at https://gitter.im/spf13/hugo](https://img.shields.io/badge/gitter-developer_chat-46bc99.svg)](https://gitter.im/spf13/hugo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
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

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD and macOS (Darwin) and [Android](https://gist.github.com/bep/a0d8a26cf6b4f8bc992729b8e50b480b) for x64, i386 and ARM architectures.

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

### Build and Install the Binaries from Source (Advanced Install)

Add Hugo and its package dependencies to your go `src` directory.

    go get -v github.com/spf13/hugo

Once the `get` completes, you should find your new `hugo` (or `hugo.exe`) executable sitting inside `$GOPATH/bin/`.

To update Hugo’s dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/spf13/hugo

## Contributing to Hugo

For a complete guide to contributing to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

We welcome contributions to Hugo of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Hugo community and maintainers are [very active](https://github.com/spf13/hugo/pulse/monthly) and helpful, and the project benefits greatly from this activity.

### Asking Support Questions

We have an active [discussion forum](http://discuss.gohugo.io) where users and developers can ask questions.
Please don't use the GitHub issue tracker to ask questions.

### Reporting Issues

If you believe you have found a defect in Hugo or its documentation, use
the GitHub issue tracker to report the problem to the Hugo maintainers.
If you're not sure if it's a bug or not, start by asking in the [discussion forum](http://discuss.gohugo.io).
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
