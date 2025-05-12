[bep]: https://github.com/bep
[bugs]: https://github.com/gohugoio/hugo/issues?q=is%3Aopen+is%3Aissue+label%3ABug
[contributing]: CONTRIBUTING.md
[create a proposal]: https://github.com/gohugoio/hugo/issues/new?labels=Proposal%2C+NeedsTriage&template=feature_request.md
[documentation repository]: https://github.com/gohugoio/hugoDocs
[documentation]: https://gohugo.io/documentation
[dragonfly bsd, freebsd, netbsd, and openbsd]: https://gohugo.io/installation/bsd
[features]: https://gohugo.io/about/features/
[forum]: https://discourse.gohugo.io
[friends]: https://github.com/gohugoio/hugo/graphs/contributors
[go]: https://go.dev/
[hugo modules]: https://gohugo.io/hugo-modules/
[installation]: https://gohugo.io/installation
[issue queue]: https://github.com/gohugoio/hugo/issues
[linux]: https://gohugo.io/installation/linux
[macos]: https://gohugo.io/installation/macos
[prebuilt binary]: https://github.com/gohugoio/hugo/releases/latest
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
[spf13]: https://github.com/spf13
[static site generator]: https://en.wikipedia.org/wiki/Static_site_generator
[support]: https://discourse.gohugo.io
[themes]: https://themes.gohugo.io/
[website]: https://gohugo.io
[windows]: https://gohugo.io/installation/windows

<a href="https://gohugo.io/"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/static/images/hugo-logo-wide.svg?sanitize=true" alt="Hugo" width="565"></a>

A fast and flexible static site generator built with love by [bep], [spf13], and [friends] in [Go].

---

[![GoDoc](https://godoc.org/github.com/gohugoio/hugo?status.svg)](https://godoc.org/github.com/gohugoio/hugo)
[![Tests on Linux, MacOS and Windows](https://github.com/gohugoio/hugo/workflows/Test/badge.svg)](https://github.com/gohugoio/hugo/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohugoio/hugo)](https://goreportcard.com/report/github.com/gohugoio/hugo)

[Website] | [Installation] | [Documentation] | [Support] | [Contributing] | <a rel="me" href="https://fosstodon.org/@gohugoio">Mastodon</a>

## Overview

Hugo is a [static site generator] written in [Go], optimized for speed and designed for flexibility. With its advanced templating system and fast asset pipelines, Hugo renders a complete site in seconds, often less.

Due to its flexible framework, multilingual support, and powerful taxonomy system, Hugo is widely used to create:

- Corporate, government, nonprofit, education, news, event, and project sites
- Documentation sites
- Image portfolios
- Landing pages
- Business, professional, and personal blogs
- Resumes and CVs

Use Hugo's embedded web server during development to instantly see changes to content, structure, behavior, and presentation. Then deploy the site to your host, or push changes to your Git provider for automated builds and deployment.

Hugo's fast asset pipelines include:

- Image processing &ndash; Convert, resize, crop, rotate, adjust colors, apply filters, overlay text and images, and extract EXIF data
- JavaScript bundling &ndash; Transpile TypeScript and JSX to JavaScript, bundle, tree shake, minify, create source maps, and perform SRI hashing.
- Sass processing &ndash; Transpile Sass to CSS, bundle, tree shake, minify, create source maps, perform SRI hashing, and integrate with PostCSS
- Tailwind CSS processing &ndash; Compile Tailwind CSS utility classes into standard CSS, bundle, tree shake, optimize, minify, perform SRI hashing, and integrate with PostCSS

And with [Hugo Modules], you can share content, assets, data, translations, themes, templates, and configuration with other projects via public or private Git repositories.

See the [features] section of the documentation for a comprehensive summary of Hugo's capabilities.

## Sponsors

<p>&nbsp;</p>
<p float="left">
  <a href="https://www.linode.com/?utm_campaign=hugosponsor&utm_medium=banner&utm_source=hugogithub" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/linode-logo_standard_light_medium.png" width="200" alt="Linode"></a>
&nbsp;&nbsp;&nbsp;
  <a href="https://www.jetbrains.com/go/?utm_source=OSS&utm_medium=referral&utm_campaign=hugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/goland.svg" width="200" alt="The complete IDE crafted for professional Go developers."></a>
  &nbsp;&nbsp;&nbsp;
  <a href="https://pinme.eth.limo/?s=hugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/logo-pinme.svg" width="200" alt="PinMe."></a>
</p>

## Editions

Hugo is available in three editions: standard, extended, and extended/deploy. While the standard edition provides core functionality, the extended and extended/deploy editions offer advanced features.

Feature|extended edition|extended/deploy edition
:--|:-:|:-:
Encode to the WebP format when [processing images]. You can decode WebP images with any edition.|:heavy_check_mark:|:heavy_check_mark:
[Transpile Sass to CSS] using the embedded LibSass transpiler. You can use the [Dart Sass] transpiler with any edition.|:heavy_check_mark:|:heavy_check_mark:
Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details].|:x:|:heavy_check_mark:

[dart sass]: https://gohugo.io/functions/css/sass/#dart-sass
[processing images]: https://gohugo.io/content-management/image-processing/
[transpile sass to css]: https://gohugo.io/functions/css/sass/
[details]: https://gohugo.io/hosting-and-deployment/hugo-deploy/

Unless your specific deployment needs require the extended/deploy edition, we recommend the extended edition.

## Installation

Install Hugo from a [prebuilt binary], package manager, or package repository. Please see the installation instructions for your operating system:

- [macOS]
- [Linux]
- [Windows]
- [DragonFly BSD, FreeBSD, NetBSD, and OpenBSD]

## Build from source

Prerequisites to build Hugo from source:

- Standard edition: Go 1.23.0 or later
- Extended edition: Go 1.23.0 or later, and GCC
- Extended/deploy edition: Go 1.23.0 or later, and GCC

Build the standard edition:

```text
go install github.com/gohugoio/hugo@latest
```

Build the extended edition:

```text
CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
```

Build the extended/deploy edition:

```text
CGO_ENABLED=1 go install -tags extended,withdeploy github.com/gohugoio/hugo@latest
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=gohugoio/hugo&type=Timeline)](https://star-history.com/#gohugoio/hugo&Timeline)

## Documentation

Hugo's [documentation] includes installation instructions, a quick start guide, conceptual explanations, reference information, and examples.

Please submit documentation issues and pull requests to the [documentation repository].

## Support

Please **do not use the issue queue** for questions or troubleshooting. Unless you are certain that your issue is a software defect, use the [forum].

Hugo’s [forum] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

## Contributing

You can contribute to the Hugo project by:

- Answering questions on the [forum]
- Improving the [documentation]
- Monitoring the [issue queue]
- Creating or improving [themes]
- Squashing [bugs]

Please submit documentation issues and pull requests to the [documentation repository].

If you have an idea for an enhancement or new feature, create a new topic on the [forum] in the "Feature" category. This will help you to:

- Determine if the capability already exists
- Measure interest
- Refine the concept

If there is sufficient interest, [create a proposal]. Do not submit a pull request until the project lead accepts the proposal.

For a complete guide to contributing to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

## Dependencies

Hugo stands on the shoulders of great open source libraries. Run `hugo env --logLevel info` to display a list of dependencies.

<details>
<summary>See current dependencies</summary>

```text
github.com/BurntSushi/locker="v0.0.0-20171006230638-a6e239ea1c69"
github.com/PuerkitoBio/goquery="v1.10.1"
github.com/alecthomas/chroma/v2="v2.15.0"
github.com/andybalholm/cascadia="v1.3.3"
github.com/armon/go-radix="v1.0.1-0.20221118154546-54df44f2176c"
github.com/bep/clocks="v0.5.0"
github.com/bep/debounce="v1.2.0"
github.com/bep/gitmap="v1.6.0"
github.com/bep/goat="v0.5.0"
github.com/bep/godartsass/v2="v2.3.2"
github.com/bep/golibsass="v1.2.0"
github.com/bep/gowebp="v0.3.0"
github.com/bep/imagemeta="v0.8.4"
github.com/bep/lazycache="v0.7.0"
github.com/bep/logg="v0.4.0"
github.com/bep/mclib="v1.20400.20402"
github.com/bep/overlayfs="v0.9.2"
github.com/bep/simplecobra="v0.5.0"
github.com/bep/tmc="v0.5.1"
github.com/cespare/xxhash/v2="v2.3.0"
github.com/clbanning/mxj/v2="v2.7.0"
github.com/cpuguy83/go-md2man/v2="v2.0.4"
github.com/disintegration/gift="v1.2.1"
github.com/dlclark/regexp2="v1.11.5"
github.com/dop251/goja="v0.0.0-20250125213203-5ef83b82af17"
github.com/evanw/esbuild="v0.24.2"
github.com/fatih/color="v1.18.0"
github.com/frankban/quicktest="v1.14.6"
github.com/fsnotify/fsnotify="v1.8.0"
github.com/getkin/kin-openapi="v0.129.0"
github.com/ghodss/yaml="v1.0.0"
github.com/go-openapi/jsonpointer="v0.21.0"
github.com/go-openapi/swag="v0.23.0"
github.com/go-sourcemap/sourcemap="v2.1.4+incompatible"
github.com/gobuffalo/flect="v1.0.3"
github.com/gobwas/glob="v0.2.3"
github.com/gohugoio/go-i18n/v2="v2.1.3-0.20230805085216-e63c13218d0e"
github.com/gohugoio/hashstructure="v0.5.0"
github.com/gohugoio/httpcache="v0.7.0"
github.com/gohugoio/hugo-goldmark-extensions/extras="v0.2.0"
github.com/gohugoio/hugo-goldmark-extensions/passthrough="v0.3.0"
github.com/gohugoio/locales="v0.14.0"
github.com/gohugoio/localescompressed="v1.0.1"
github.com/golang/freetype="v0.0.0-20170609003504-e2365dfdc4a0"
github.com/google/go-cmp="v0.6.0"
github.com/google/pprof="v0.0.0-20250208200701-d0013a598941"
github.com/gorilla/websocket="v1.5.3"
github.com/hairyhenderson/go-codeowners="v0.7.0"
github.com/hashicorp/golang-lru/v2="v2.0.7"
github.com/jdkato/prose="v1.2.1"
github.com/josharian/intern="v1.0.0"
github.com/kr/pretty="v0.3.1"
github.com/kr/text="v0.2.0"
github.com/kyokomi/emoji/v2="v2.2.13"
github.com/lucasb-eyer/go-colorful="v1.2.0"
github.com/mailru/easyjson="v0.7.7"
github.com/makeworld-the-better-one/dither/v2="v2.4.0"
github.com/marekm4/color-extractor="v1.2.1"
github.com/mattn/go-colorable="v0.1.13"
github.com/mattn/go-isatty="v0.0.20"
github.com/mattn/go-runewidth="v0.0.9"
github.com/mazznoer/csscolorparser="v0.1.5"
github.com/mitchellh/mapstructure="v1.5.1-0.20231216201459-8508981c8b6c"
github.com/mohae/deepcopy="v0.0.0-20170929034955-c48cc78d4826"
github.com/muesli/smartcrop="v0.3.0"
github.com/niklasfasching/go-org="v1.7.0"
github.com/oasdiff/yaml3="v0.0.0-20241210130736-a94c01f36349"
github.com/oasdiff/yaml="v0.0.0-20241210131133-6b86fb107d80"
github.com/olekukonko/tablewriter="v0.0.5"
github.com/pbnjay/memory="v0.0.0-20210728143218-7b4eea64cf58"
github.com/pelletier/go-toml/v2="v2.2.3"
github.com/perimeterx/marshmallow="v1.1.5"
github.com/pkg/browser="v0.0.0-20240102092130-5ac0b6a4141c"
github.com/pkg/errors="v0.9.1"
github.com/rivo/uniseg="v0.4.7"
github.com/rogpeppe/go-internal="v1.13.1"
github.com/russross/blackfriday/v2="v2.1.0"
github.com/sass/libsass="3.6.6"
github.com/spf13/afero="v1.11.0"
github.com/spf13/cast="v1.7.1"
github.com/spf13/cobra="v1.8.1"
github.com/spf13/fsync="v0.10.1"
github.com/spf13/pflag="v1.0.6"
github.com/tdewolff/minify/v2="v2.20.37"
github.com/tdewolff/parse/v2="v2.7.15"
github.com/tetratelabs/wazero="v1.8.2"
github.com/webmproject/libwebp="v1.3.2"
github.com/yuin/goldmark-emoji="v1.0.4"
github.com/yuin/goldmark="v1.7.8"
go.uber.org/automaxprocs="v1.5.3"
golang.org/x/crypto="v0.33.0"
golang.org/x/exp="v0.0.0-20250210185358-939b2ce775ac"
golang.org/x/image="v0.24.0"
golang.org/x/mod="v0.23.0"
golang.org/x/net="v0.35.0"
golang.org/x/sync="v0.11.0"
golang.org/x/sys="v0.30.0"
golang.org/x/text="v0.22.0"
golang.org/x/tools="v0.30.0"
golang.org/x/xerrors="v0.0.0-20240903120638-7835f813f4da"
gonum.org/v1/plot="v0.15.0"
google.golang.org/protobuf="v1.36.5"
gopkg.in/yaml.v2="v2.4.0"
gopkg.in/yaml.v3="v3.0.1"
oss.terrastruct.com/d2="v0.6.9"
oss.terrastruct.com/util-go="v0.0.0-20241005222610-44c011a04896"
rsc.io/qr="v0.2.0"
software.sslmate.com/src/go-pkcs12="v0.2.0"
```
</details>
