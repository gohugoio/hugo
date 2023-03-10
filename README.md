<a href="https://gohugo.io/"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/static/images/hugo-logo-wide.svg?sanitize=true" alt="Hugo" width="565"></a>

A Fast and Flexible Static Site Generator built with love by [bep](https://github.com/bep), [spf13](https://spf13.com/) and [friends](https://github.com/gohugoio/hugo/graphs/contributors) in [Go](https://go.dev/).

[Website](https://gohugo.io) |
[Forum](https://discourse.gohugo.io) |
[Documentation](https://gohugo.io/getting-started/) |
[Installation Guide](https://gohugo.io/getting-started/installing/) |
[Contribution Guide](CONTRIBUTING.md) |
[Twitter](https://twitter.com/gohugoio)

[![GoDoc](https://godoc.org/github.com/gohugoio/hugo?status.svg)](https://godoc.org/github.com/gohugoio/hugo)
[![Tests on Linux, MacOS and Windows](https://github.com/gohugoio/hugo/workflows/Test/badge.svg)](https://github.com/gohugoio/hugo/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohugoio/hugo)](https://goreportcard.com/report/github.com/gohugoio/hugo)

* [Overview](#overview)
* [Banner Sponsors](#banner-sponsors)
* [Supported Architectures](#supported-architectures)
* [Choose How to Install](#choose-how-to-install)
   * [Install Hugo as Your Site Generator (Binary Install)](#install-hugo-as-your-site-generator-binary-install)
   * [Build and Install the Binary from Source (Using the Go toolchain)](#build-and-install-the-binary-from-source-using-the-go-toolchain)
* [The Hugo Documentation](#the-hugo-documentation)
* [Contributing to Hugo](#contributing-code-to-hugo)
* [Dependencies](#dependencies)

## Overview

Hugo is a static HTML and CSS website generator written in [Go](https://go.dev/).
It is optimized for speed, ease of use, and configurability.
Hugo takes a directory with content and templates and renders them into a full HTML website.

Hugo relies on Markdown files with front matter for metadata, and you can run Hugo from any directory.
This works well for shared hosts and other systems where you don’t have a privileged account.

Hugo renders a typical website of moderate size in a fraction of a second.
A good rule of thumb is that each piece of content renders in around 1 millisecond.

Hugo is designed to work well for any kind of website including blogs, tumbles, and docs.

## Banner Sponsors
<p>&nbsp;</p>
<p float="left">
  <a href="https://www.linode.com/?utm_campaign=hugosponsor&utm_medium=banner&utm_source=hugogithub" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/assets/images/sponsors/linode-logo_standard_light_medium.png" width="200" alt="Linode"></a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="https://buttercms.com/hugo-cms/?utm_campaign=sponsorship&utm_medium=banner&utm_source=hugogithub" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/assets/images/sponsors/butter-dark.svg?sanitize=true" width="280" alt="ButterCMS"></a>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
<a href="https://www.gravitykit.com/?ref=532&campaign=hugo&utm_campaign=hugosponsor&utm_medium=banner&utm_source=hugogithub" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/assets/images/sponsors/graitykit-dark.svg?sanitize=true" width="160" alt="Gravity Kit"></a>
</p>
<p>&nbsp;</p>

## Supported Architectures

Currently, we provide pre-built Hugo binaries for Windows, Linux, FreeBSD, NetBSD, DragonFly BSD, OpenBSD, macOS (Darwin), and [Android](https://gist.github.com/bep/a0d8a26cf6b4f8bc992729b8e50b480b) for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including Plan 9 and Solaris.

**Complete documentation is available at [Hugo Documentation](https://gohugo.io/getting-started/).**

## Choose How to Install

If you want to use Hugo as your site generator, simply install the Hugo binaries.

To contribute to the Hugo source code or documentation, you should [fork the Hugo GitHub project](https://github.com/gohugoio/hugo#fork-destination-box) and clone it to your local machine.

Finally, you can install the Hugo source code with `go`, build the binaries yourself, and run Hugo that way.
Building the binaries is an easy task for an experienced `go` getter.

### Install Hugo as Your Site Generator (Binary Install)

Use the [installation instructions in the Hugo documentation](https://gohugo.io/getting-started/installing/).

### Build and Install the Binary from Source (Using the Go toolchain)

#### Prerequisite Tools

* [Go (we test it with the last 2 major versions; but note that Hugo 0.95.0 only builds with >= Go 1.18.)](https://golang.org/dl/)

#### Fetch from GitHub

To fetch, build and install from the Github source:

```bash
go install github.com/gohugoio/hugo@latest
```

If you want to compile with Sass/SCSS support use `--tags extended` and make sure `CGO_ENABLED=1` is set in your go environment. If you don't want to have CGO enabled, you may use the following command to temporarily enable CGO only for hugo compilation:

```bash
CGO_ENABLED=1 go install --tags extended github.com/gohugoio/hugo@latest
```

## The Hugo Documentation

The Hugo documentation now lives in its own repository, see https://github.com/gohugoio/hugoDocs. But we do keep a version of that documentation as a `git subtree` in this repository. To build the sub folder `/docs` as a Hugo site, you need to clone this repo:

```bash
git clone git@github.com:gohugoio/hugo.git
```
## Contributing code to Hugo

For a complete guide to contributing to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

We welcome contributions to Hugo of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Hugo community and maintainers are [very active](https://github.com/gohugoio/hugo/pulse/monthly) and helpful, and the project benefits greatly from this activity.

## Asking Support Questions

We have an active [discussion forum](https://discourse.gohugo.io) where users and developers can ask questions.
Please don't use the GitHub issue tracker to ask questions.

## Reporting Issues

If you believe you have found a defect in Hugo or its documentation, use
the GitHub issue tracker to report the problem to the Hugo maintainers.
If you're not sure if it's a bug or not, start by asking in the [discussion forum](https://discourse.gohugo.io).
When reporting the issue, please provide the version of Hugo in use (`hugo version`).

## Dependencies

Hugo stands on the shoulder of many great open source libraries.

If you run `hugo env -v` you will get a complete and up to date list.

In Hugo 0.111.2 that list is, in lexical order:

```
cloud.google.com/go/compute="v1.6.1"
cloud.google.com/go/iam="v0.3.0"
cloud.google.com/go/storage="v1.22.0"
cloud.google.com/go="v0.101.0"
github.com/Azure/azure-pipeline-go="v0.2.3"
github.com/Azure/azure-storage-blob-go="v0.14.0"
github.com/Azure/go-autorest/autorest/adal="v0.9.15"
github.com/Azure/go-autorest/autorest/date="v0.3.0"
github.com/Azure/go-autorest/autorest="v0.11.20"
github.com/Azure/go-autorest/logger="v0.2.1"
github.com/Azure/go-autorest/tracing="v0.6.0"
github.com/BurntSushi/locker="v0.0.0-20171006230638-a6e239ea1c69"
github.com/PuerkitoBio/purell="v1.1.1"
github.com/PuerkitoBio/urlesc="v0.0.0-20170810143723-de5bf2ad4578"
github.com/alecthomas/chroma/v2="v2.5.0"
github.com/armon/go-radix="v1.0.0"
github.com/aws/aws-sdk-go-v2/config="v1.7.0"
github.com/aws/aws-sdk-go-v2/credentials="v1.4.0"
github.com/aws/aws-sdk-go-v2/feature/ec2/imds="v1.5.0"
github.com/aws/aws-sdk-go-v2/internal/ini="v1.2.2"
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url="v1.3.0"
github.com/aws/aws-sdk-go-v2/service/sso="v1.4.0"
github.com/aws/aws-sdk-go-v2/service/sts="v1.7.0"
github.com/aws/aws-sdk-go-v2="v1.9.0"
github.com/aws/aws-sdk-go="v1.43.5"
github.com/aws/smithy-go="v1.8.0"
github.com/bep/clock="v0.3.0"
github.com/bep/debounce="v1.2.0"
github.com/bep/gitmap="v1.1.2"
github.com/bep/goat="v0.5.0"
github.com/bep/godartsass="v0.16.0"
github.com/bep/golibsass="v1.1.0"
github.com/bep/gowebp="v0.2.0"
github.com/bep/lazycache="v0.2.0"
github.com/bep/overlayfs="v0.6.0"
github.com/bep/tmc="v0.5.1"
github.com/clbanning/mxj/v2="v2.5.7"
github.com/cli/safeexec="v1.0.0"
github.com/cpuguy83/go-md2man/v2="v2.0.2"
github.com/disintegration/gift="v1.2.1"
github.com/dlclark/regexp2="v1.7.0"
github.com/dustin/go-humanize="v1.0.0"
github.com/evanw/esbuild="v0.17.0"
github.com/frankban/quicktest="v1.14.4"
github.com/fsnotify/fsnotify="v1.6.0"
github.com/getkin/kin-openapi="v0.110.0"
github.com/ghodss/yaml="v1.0.0"
github.com/go-openapi/jsonpointer="v0.19.5"
github.com/go-openapi/swag="v0.19.5"
github.com/gobuffalo/flect="v0.3.0"
github.com/gobwas/glob="v0.2.3"
github.com/gohugoio/go-i18n/v2="v2.1.3-0.20210430103248-4c28c89f8013"
github.com/gohugoio/locales="v0.14.0"
github.com/gohugoio/localescompressed="v1.0.1"
github.com/golang-jwt/jwt/v4="v4.0.0"
github.com/golang/groupcache="v0.0.0-20210331224755-41bb18bfe9da"
github.com/golang/protobuf="v1.5.2"
github.com/google/go-cmp="v0.5.9"
github.com/google/uuid="v1.3.0"
github.com/google/wire="v0.5.0"
github.com/googleapis/gax-go/v2="v2.3.0"
github.com/googleapis/go-type-adapters="v1.0.0"
github.com/gorilla/websocket="v1.5.0"
github.com/hairyhenderson/go-codeowners="v0.2.3-0.20201026200250-cdc7c0759690"
github.com/hashicorp/golang-lru/v2="v2.0.1"
github.com/invopop/yaml="v0.1.0"
github.com/jdkato/prose="v1.2.1"
github.com/jmespath/go-jmespath="v0.4.0"
github.com/kr/pretty="v0.3.1"
github.com/kr/text="v0.2.0"
github.com/kyokomi/emoji/v2="v2.2.11"
github.com/mailru/easyjson="v0.0.0-20190626092158-b2ccc519800e"
github.com/marekm4/color-extractor="v1.2.0"
github.com/mattn/go-ieproxy="v0.0.1"
github.com/mattn/go-isatty="v0.0.17"
github.com/mattn/go-runewidth="v0.0.9"
github.com/mitchellh/hashstructure="v1.1.0"
github.com/mitchellh/mapstructure="v1.5.0"
github.com/mohae/deepcopy="v0.0.0-20170929034955-c48cc78d4826"
github.com/muesli/smartcrop="v0.3.0"
github.com/niklasfasching/go-org="v1.6.5"
github.com/olekukonko/tablewriter="v0.0.5"
github.com/pelletier/go-toml/v2="v2.0.6"
github.com/rogpeppe/go-internal="v1.9.0"
github.com/russross/blackfriday/v2="v2.1.0"
github.com/rwcarlsen/goexif="v0.0.0-20190401172101-9e8deecbddbd"
github.com/sanity-io/litter="v1.5.5"
github.com/sass/libsass="3.6.5"
github.com/spf13/afero="v1.9.3"
github.com/spf13/cast="v1.5.0"
github.com/spf13/cobra="v1.6.1"
github.com/spf13/fsync="v0.9.0"
github.com/spf13/jwalterweatherman="v1.1.0"
github.com/spf13/pflag="v1.0.5"
github.com/tdewolff/minify/v2="v2.12.4"
github.com/tdewolff/parse/v2="v2.6.5"
github.com/webmproject/libwebp="v1.2.4"
github.com/yuin/goldmark="v1.5.4"
go.opencensus.io="v0.24.0"
go.uber.org/atomic="v1.10.0"
gocloud.dev="v0.24.0"
golang.org/x/crypto="v0.3.0"
golang.org/x/exp="v0.0.0-20221031165847-c99f073a8326"
golang.org/x/image="v0.5.0"
golang.org/x/net="v0.7.0"
golang.org/x/oauth2="v0.2.0"
golang.org/x/sync="v0.1.0"
golang.org/x/sys="v0.5.0"
golang.org/x/text="v0.7.0"
golang.org/x/tools="v0.4.0"
golang.org/x/xerrors="v0.0.0-20220907171357-04be3eba64a2"
google.golang.org/api="v0.76.0"
google.golang.org/genproto="v0.0.0-20220426171045-31bebdecfb46"
google.golang.org/grpc="v1.46.0"
google.golang.org/protobuf="v1.28.1"
gopkg.in/yaml.v2="v2.4.0"
gopkg.in/yaml.v3="v3.0.1"
```
