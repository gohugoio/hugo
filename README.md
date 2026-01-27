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

- Image processing &ndash; Convert, resize, crop, rotate, adjust colors, apply filters, overlay text and images, and extract metadata
- JavaScript bundling &ndash; Transpile TypeScript and JSX to JavaScript, bundle, tree shake, minify, create source maps, and perform SRI hashing.
- Sass processing &ndash; Transpile Sass to CSS, bundle, tree shake, minify, create source maps, perform SRI hashing, and integrate with PostCSS
- Tailwind CSS processing &ndash; Compile Tailwind CSS utility classes into standard CSS, bundle, tree shake, optimize, minify, perform SRI hashing, and integrate with PostCSS

And with [Hugo Modules], you can share content, assets, data, translations, themes, templates, and configuration with other projects via public or private Git repositories.

See the [features] section of the documentation for a comprehensive summary of Hugo's capabilities.

## Sponsors

<p>&nbsp;</p>
<p float="left">
  <a href="https://www.jetbrains.com/go/?utm_source=OSS&utm_medium=referral&utm_campaign=hugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/goland.svg" width="200" alt="The complete IDE crafted for professional Go developers."></a>
  &nbsp;&nbsp;&nbsp;
    <a href="https://cloudcannon.com/hugo-cms/?utm_campaign=HugoSponsorship&utm_source=sponsor&utm_content=gohugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/cloudcannon-cms-logo.svg" width="200" alt="CloudCannon"></a>
</p>

## Editions

Hugo offers a standard edition with core features, plus extended and extended/deploy editions with more. Use the standard edition unless you need the features below.

Feature|extended edition|extended/deploy edition
:--|:-:|:-:
[Transpile Sass to CSS] via embedded LibSass. Note that embedded LibSass was deprecated in v0.153.0 and will be removed in a future release. Use the [Dart Sass] transpiler instead, which is compatible with any edition.|:heavy_check_mark:|:heavy_check_mark:
Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details].|:x:|:heavy_check_mark:

[dart sass]: https://gohugo.io/functions/css/sass/#dart-sass
[transpile sass to css]: https://gohugo.io/functions/css/sass/
[details]: https://gohugo.io/hosting-and-deployment/hugo-deploy/

## Installation

Install Hugo from a [prebuilt binary], package manager, or package repository. Please see the installation instructions for your operating system:

- [macOS]
- [Linux]
- [Windows]
- [DragonFly BSD, FreeBSD, NetBSD, and OpenBSD]

## Build from source

Prerequisites to build Hugo from source:

- Standard edition: Go 1.24.0 or later
- Extended edition: Go 1.24.0 or later, and GCC
- Extended/deploy edition: Go 1.24.0 or later, and GCC

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

## License

For the Hugo source code, see [LICENSE](/LICENSE).

We also bundle some libraries in binary/WASM form:

* [libwebp](https://github.com/webmproject/libwebp), [BSD-3-Clause license](https://github.com/webmproject/libwebp?tab=BSD-3-Clause-1-ov-file#readme)
* [Katex](https://github.com/KaTeX/KaTeX), [MIT license](https://github.com/KaTeX/KaTeX?tab=MIT-1-ov-file#readme)
* [QuickJS](https://github.com/bellard/quickjs?tab=License-1-ov-file#readme), [License](https://github.com/bellard/quickjs?tab=License-1-ov-file#readme)

## Dependencies

Hugo stands on the shoulders of great open source libraries. Run `hugo env --logLevel info` to display a list of dependencies.

<details>
<summary>See current dependencies</summary>

```text
github.com/BurntSushi/locker="v0.0.0-20171006230638-a6e239ea1c69"
github.com/JohannesKaufmann/dom="v0.2.0"
github.com/JohannesKaufmann/html-to-markdown/v2="v2.5.0"
github.com/alecthomas/chroma/v2="v2.21.1"
github.com/aymerick/douceur="v0.2.0"
github.com/bep/clocks="v0.5.0"
github.com/bep/debounce="v1.2.0"
github.com/bep/gitmap="v1.9.0"
github.com/bep/goat="v0.5.0"
github.com/bep/godartsass/v2="v2.5.0"
github.com/bep/golibsass="v1.2.0"
github.com/bep/goportabletext="v0.1.0"
github.com/bep/helpers="v0.6.0"
github.com/bep/imagemeta="v0.12.0"
github.com/bep/lazycache="v0.8.0"
github.com/bep/logg="v0.4.0"
github.com/bep/mclib="v1.20400.20402"
github.com/bep/overlayfs="v0.10.0"
github.com/bep/simplecobra="v0.6.1"
github.com/bep/textandbinarywriter="v0.0.0-20251212174530-cd9f0732f60f"
github.com/bep/tmc="v0.5.1"
github.com/bits-and-blooms/bitset="v1.24.4"
github.com/cespare/xxhash/v2="v2.3.0"
github.com/clbanning/mxj/v2="v2.7.0"
github.com/clipperhouse/displaywidth="v0.6.0"
github.com/clipperhouse/stringish="v0.1.1"
github.com/clipperhouse/uax29/v2="v2.3.0"
github.com/cpuguy83/go-md2man/v2="v2.0.6"
github.com/disintegration/gift="v1.2.1"
github.com/dlclark/regexp2="v1.11.5"
github.com/evanw/esbuild="v0.27.2"
github.com/fatih/color="v1.18.0"
github.com/frankban/quicktest="v1.14.6"
github.com/fsnotify/fsnotify="v1.9.0"
github.com/getkin/kin-openapi="v0.133.0"
github.com/go-openapi/jsonpointer="v0.21.0"
github.com/go-openapi/swag="v0.23.0"
github.com/gobuffalo/flect="v1.0.3"
github.com/gobwas/glob="v0.2.3"
github.com/goccy/go-yaml="v1.19.1"
github.com/gohugoio/go-i18n/v2="v2.1.3-0.20251018145728-cfcc22d823c6"
github.com/gohugoio/go-radix="v1.2.0"
github.com/gohugoio/hashstructure="v0.6.0"
github.com/gohugoio/httpcache="v0.8.0"
github.com/gohugoio/hugo-goldmark-extensions/extras="v0.5.0"
github.com/gohugoio/hugo-goldmark-extensions/passthrough="v0.3.1"
github.com/gohugoio/locales="v0.14.0"
github.com/gohugoio/localescompressed="v1.0.1"
github.com/google/go-cmp="v0.7.0"
github.com/gorilla/css="v1.0.1"
github.com/gorilla/websocket="v1.5.3"
github.com/hairyhenderson/go-codeowners="v0.7.0"
github.com/hashicorp/golang-lru/v2="v2.0.7"
github.com/jdkato/prose="v1.2.1"
github.com/josharian/intern="v1.0.0"
github.com/kr/pretty="v0.3.1"
github.com/kr/text="v0.2.0"
github.com/kyokomi/emoji/v2="v2.2.13"
github.com/mailru/easyjson="v0.7.7"
github.com/makeworld-the-better-one/dither/v2="v2.4.0"
github.com/marekm4/color-extractor="v1.2.1"
github.com/mattn/go-colorable="v0.1.13"
github.com/mattn/go-isatty="v0.0.20"
github.com/mattn/go-runewidth="v0.0.19"
github.com/microcosm-cc/bluemonday="v1.0.27"
github.com/mitchellh/mapstructure="v1.5.1-0.20231216201459-8508981c8b6c"
github.com/mohae/deepcopy="v0.0.0-20170929034955-c48cc78d4826"
github.com/muesli/smartcrop="v0.3.0"
github.com/niklasfasching/go-org="v1.9.1"
github.com/oasdiff/yaml3="v0.0.0-20250309153720-d2182401db90"
github.com/oasdiff/yaml="v0.0.0-20250309154309-f31be36b4037"
github.com/olekukonko/cat="v0.0.0-20250911104152-50322a0618f6"
github.com/olekukonko/errors="v1.1.0"
github.com/olekukonko/ll="v0.1.3"
github.com/olekukonko/tablewriter="v1.1.2"
github.com/pbnjay/memory="v0.0.0-20210728143218-7b4eea64cf58"
github.com/pelletier/go-toml/v2="v2.2.4"
github.com/perimeterx/marshmallow="v1.1.5"
github.com/pkg/browser="v0.0.0-20240102092130-5ac0b6a4141c"
github.com/pkg/errors="v0.9.1"
github.com/rogpeppe/go-internal="v1.14.1"
github.com/russross/blackfriday/v2="v2.1.0"
github.com/sass/dart-sass/compiler="1.97.1"
github.com/sass/dart-sass/implementation="1.97.1"
github.com/sass/dart-sass/protocol="3.2.0"
github.com/spf13/afero="v1.15.0"
github.com/spf13/cast="v1.10.0"
github.com/spf13/cobra="v1.10.2"
github.com/spf13/fsync="v0.10.1"
github.com/spf13/pflag="v1.0.9"
github.com/tdewolff/minify/v2="v2.24.8"
github.com/tdewolff/parse/v2="v2.8.5"
github.com/tetratelabs/wazero="v1.10.1"
github.com/webmproject/libwebp="v1.6.0"
github.com/woodsbury/decimal128="v1.3.0"
github.com/yuin/goldmark-emoji="v1.0.6"
github.com/yuin/goldmark="v1.7.13"
go.uber.org/automaxprocs="v1.5.3"
go.yaml.in/yaml/v3="v3.0.4"
golang.org/x/crypto="v0.46.0"
golang.org/x/image="v0.34.0"
golang.org/x/mod="v0.31.0"
golang.org/x/net="v0.48.0"
golang.org/x/sync="v0.19.0"
golang.org/x/sys="v0.39.0"
golang.org/x/text="v0.32.0"
golang.org/x/tools="v0.40.0"
google.golang.org/protobuf="v1.36.10"
gopkg.in/yaml.v3="v3.0.1"
rsc.io/qr="v0.2.0"
software.sslmate.com/src/go-pkcs12="v0.2.0"
```
</details>
