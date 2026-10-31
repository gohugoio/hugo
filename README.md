[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[bep]: https://github.com/bep
[bugs]: https://github.com/gohugoio/hugo/issues?q=is%3Aopen+is%3Aissue+label%3ABug
[contributing]: CONTRIBUTING.md
[create a proposal]: https://github.com/gohugoio/hugo/issues/new?labels=Proposal%2C+NeedsTriage&template=feature_request.md
[details]: https://gohugo.io/docs/guides/host-and-deploy/deploy-with-hugo-deploy/
[documentation repository]: https://github.com/gohugoio/hugoDocs
[documentation]: https://gohugo.io/docs/
[dragonfly bsd, freebsd, netbsd, and openbsd]: https://gohugo.io/docs/installation/bsd/
[features]: https://gohugo.io/about/features/
[forum]: https://discourse.gohugo.io
[friends]: https://github.com/gohugoio/hugo/graphs/contributors
[hugo modules]: https://gohugo.io/docs/concepts/modules/
[installation]: https://gohugo.io/docs/installation/
[issue queue]: https://github.com/gohugoio/hugo/issues
[linux]: https://gohugo.io/docs/installation/linux/
[macos]: https://gohugo.io/docs/installation/macos/
[prebuilt binary]: https://github.com/gohugoio/hugo/releases/latest
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
[spf13]: https://github.com/spf13
[static site generator]: https://en.wikipedia.org/wiki/Static_site_generator
[support]: https://discourse.gohugo.io
[themes]: https://themes.gohugo.io/
[website]: https://gohugo.io
[windows]: https://gohugo.io/docs/installation/windows/

<a href="https://gohugo.io/"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/static/images/hugo-logo-wide.svg?sanitize=true" alt="Hugo" width="565"></a>

A fast and flexible static site generator built with love by [bep][], [spf13][], and [friends][] in Go.

---

[![GoDoc](https://godoc.org/github.com/gohugoio/hugo?status.svg)](https://godoc.org/github.com/gohugoio/hugo)
[![Tests on Linux, MacOS and Windows](https://github.com/gohugoio/hugo/workflows/Test/badge.svg)](https://github.com/gohugoio/hugo/actions?query=workflow%3ATest)

[Website][] | [Installation][] | [Documentation][] | [Support][] | [Contributing][] | <a rel="me" href="https://fosstodon.org/@gohugoio">Mastodon</a>

## Overview

Hugo is a [static site generator][] written in Go, optimized for speed and designed for flexibility. With its advanced templating system and fast asset pipelines, Hugo renders a complete site in seconds, often less.

Due to its flexible framework, multilingual support, and powerful taxonomy system, Hugo is widely used to create:

- Corporate, government, nonprofit, education, news, event, and project sites
- Documentation sites
- Image portfolios
- Landing pages
- Business, professional, and personal blogs
- Resumes and CVs

Use Hugo's embedded web server during development to instantly see changes to content, structure, behavior, and presentation. Then deploy the site to your host, or push changes to your Git provider for automated builds and deployment.

Hugo's fast asset pipelines include:

- CSS Processing &ndash; Bundle, transform, minify, create source maps, perform SRI hashing, and integrate with PostCSS.
- Image processing &ndash; Convert, resize, crop, rotate, adjust colors, apply filters, overlay text and images, and extract metadata
- JavaScript bundling &ndash; Transpile TypeScript and JSX to JavaScript, bundle, tree shake, minify, create source maps, and perform SRI hashing.
- Sass processing &ndash; Transpile Sass to CSS, bundle, tree shake, minify, create source maps, perform SRI hashing, and integrate with PostCSS
- Tailwind CSS processing &ndash; Compile Tailwind CSS utility classes into standard CSS, bundle, tree shake, optimize, minify, perform SRI hashing, and integrate with PostCSS

And with [Hugo Modules][], you can share content, assets, data, translations, themes, templates, and configuration with other projects via public or private Git repositories.

See the [features][] section of the documentation for a comprehensive summary of Hugo's capabilities.

## Sponsors

<p>&nbsp;</p>
<p float="left">
  <a href="https://www.jetbrains.com/go/?utm_source=OSS&utm_medium=referral&utm_campaign=hugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/goland.svg" width="200" alt="The complete IDE crafted for professional Go developers."></a>
  &nbsp;&nbsp;&nbsp;
    <a href="https://cloudcannon.com/hugo-cms/?utm_campaign=HugoSponsorship&utm_source=sponsor&utm_content=gohugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/hugoDocs/master/assets/images/sponsors/cloudcannon-cms-logo.svg" width="200" alt="CloudCannon"></a>
</p>

## Editions

Hugo is available in two editions. Use the standard edition unless you need additional features.

Feature|standard|deploy
:--|:-:|:-:
Core features|:heavy_check_mark:|:heavy_check_mark:
Direct cloud deployment (1)|:x:|:heavy_check_mark:

(1) Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details][].

## Installation

Install Hugo from a [prebuilt binary][], package manager, or package repository. Please see the installation instructions for your operating system:

- [macOS][]
- [Linux][]
- [Windows][]
- [DragonFly BSD, FreeBSD, NetBSD, and OpenBSD][]

## Build from source

To build Hugo from source you must install:

1. [Git][]
1. [Go][] version 1.27.0 or later

### Standard edition

To build and install the standard edition:

```sh
CGO_ENABLED=0 go install github.com/gohugoio/hugo@latest
```

### Deploy edition

To build and install the deploy edition:

```sh
CGO_ENABLED=0 go install -tags withdeploy github.com/gohugoio/hugo@latest
```

## Documentation

Hugo's [documentation][] includes installation instructions, a quick start guide, conceptual explanations, reference information, and examples.

Please submit documentation issues and pull requests to the [documentation repository][].

## Support

Please **do not use the issue queue** for questions or troubleshooting. Unless you are certain that your issue is a software defect, use the [forum][].

Hugo's [forum][] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help][] before asking your first question.

## Contributing

You can contribute to the Hugo project by:

- Answering questions on the [forum][]
- Improving the [documentation][]
- Monitoring the [issue queue][]
- Creating or improving [themes][]
- Squashing [bugs][]

Please submit documentation issues and pull requests to the [documentation repository][].

If you have an idea for an enhancement or new feature, create a new topic on the [forum][] in the "Feature" category. This will help you to:

- Determine if the capability already exists
- Measure interest
- Refine the concept

If there is sufficient interest, [create a proposal][]. Do not submit a pull request until the project lead accepts the proposal.

For a complete guide to contributing to Hugo, see the [Contribution Guide](CONTRIBUTING.md).

## License

For the Hugo source code, see [LICENSE](/LICENSE).

We also bundle some libraries in binary/WASM form:

- [libwebp](https://github.com/webmproject/libwebp), [BSD-3-Clause license](https://github.com/webmproject/libwebp?tab=BSD-3-Clause-1-ov-file#readme)
- [Katex](https://github.com/KaTeX/KaTeX), [MIT license](https://github.com/KaTeX/KaTeX?tab=MIT-1-ov-file#readme)
- [QuickJS](https://github.com/bellard/quickjs?tab=License-1-ov-file#readme), [License](https://github.com/bellard/quickjs?tab=License-1-ov-file#readme)
