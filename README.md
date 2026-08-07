[Clang]: https://clang.llvm.org
[GCC]: https://gcc.gnu.org
[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[bep]: https://github.com/bep
[bugs]: https://github.com/gohugoio/hugo/issues?q=is%3Aopen+is%3Aissue+label%3ABug
[contributing]: CONTRIBUTING.md
[create a proposal]: https://github.com/gohugoio/hugo/issues/new?labels=Proposal%2C+NeedsTriage&template=feature_request.md
[dart sass]: https://gohugo.io/docs/reference/functions/css/sass/#dart-sass
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
[transpile sass to css]: https://gohugo.io/docs/reference/functions/css/sass/
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

Hugo is available in several editions. Use the standard edition unless you need additional features.

Feature|standard|deploy|extended|extended/deploy
:--|:-:|:-:|:-:|:-:
Core features|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Direct cloud deployment (1)|:x:|:heavy_check_mark:|:x:|:heavy_check_mark:
LibSass support (2)|:x:|:x:|:heavy_check_mark:|:heavy_check_mark:

(1) Deploy your site directly to a Google Cloud Storage bucket, an AWS S3 bucket, or an Azure Storage container. See&nbsp;[details][].

(2) [Transpile Sass to CSS][] via embedded LibSass. Note that embedded LibSass was deprecated in v0.153.0 and will be removed in a future release. Use the [Dart Sass][] transpiler instead, which is compatible with any edition.

## Installation

Install Hugo from a [prebuilt binary][], package manager, or package repository. Please see the installation instructions for your operating system:

- [macOS][]
- [Linux][]
- [Windows][]
- [DragonFly BSD, FreeBSD, NetBSD, and OpenBSD][]

## Build from source

To build Hugo from source you must install:

1. [Git][]
1. [Go][] version 1.26.0 or later

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

### Extended edition

To build and install the extended edition, first install a C compiler such as [GCC][] or [Clang][] and then run the following command.

```sh
CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
```

### Extended/deploy edition

To build and install the extended/deploy edition, first install a C compiler such as [GCC][] or [Clang][] and then run the following command.

```sh
CGO_ENABLED=1 go install -tags extended,withdeploy github.com/gohugoio/hugo@latest
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=gohugoio/hugo&type=Timeline)](https://star-history.com/#gohugoio/hugo&Timeline)

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

## Dependencies

Hugo stands on the shoulders of great open source libraries. Run `hugo env --logLevel info` to display a list of dependencies.

<details>
<summary>See current dependencies</summary>

```text
cel.dev/expr="v0.25.1"
cloud.google.com/go/auth/oauth2adapt="v0.2.8"
cloud.google.com/go/auth="v0.20.0"
cloud.google.com/go/compute/metadata="v0.9.0"
cloud.google.com/go/iam="v1.5.3"
cloud.google.com/go/monitoring="v1.24.3"
cloud.google.com/go/storage="v1.57.2"
cloud.google.com/go="v0.123.0"
github.com/Azure/azure-sdk-for-go/sdk/azcore="v1.20.0"
github.com/Azure/azure-sdk-for-go/sdk/azidentity="v1.13.1"
github.com/Azure/azure-sdk-for-go/sdk/internal="v1.11.2"
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob="v1.6.3"
github.com/Azure/go-autorest/autorest/to="v0.4.1"
github.com/AzureAD/microsoft-authentication-library-for-go="v1.6.0"
github.com/BurntSushi/locker="v0.0.0-20171006230638-a6e239ea1c69"
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp="v1.31.0"
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric="v0.54.0"
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping="v0.54.0"
github.com/JohannesKaufmann/dom="v0.3.1"
github.com/JohannesKaufmann/html-to-markdown/v2="v2.5.2"
github.com/alecthomas/chroma/v2="v2.27.0"
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream="v1.7.8"
github.com/aws/aws-sdk-go-v2/config="v1.32.2"
github.com/aws/aws-sdk-go-v2/credentials="v1.19.2"
github.com/aws/aws-sdk-go-v2/feature/ec2/imds="v1.18.14"
github.com/aws/aws-sdk-go-v2/feature/s3/manager="v1.20.12"
github.com/aws/aws-sdk-go-v2/internal/configsources="v1.4.22"
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2="v2.7.22"
github.com/aws/aws-sdk-go-v2/internal/ini="v1.8.4"
github.com/aws/aws-sdk-go-v2/internal/v4a="v1.4.22"
github.com/aws/aws-sdk-go-v2/service/cloudfront="v1.61.1"
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding="v1.13.7"
github.com/aws/aws-sdk-go-v2/service/internal/checksum="v1.9.13"
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url="v1.13.21"
github.com/aws/aws-sdk-go-v2/service/internal/s3shared="v1.19.21"
github.com/aws/aws-sdk-go-v2/service/s3="v1.97.3"
github.com/aws/aws-sdk-go-v2/service/signin="v1.0.2"
github.com/aws/aws-sdk-go-v2/service/sso="v1.30.5"
github.com/aws/aws-sdk-go-v2/service/ssooidc="v1.35.10"
github.com/aws/aws-sdk-go-v2/service/sts="v1.41.2"
github.com/aws/aws-sdk-go-v2="v1.41.6"
github.com/aws/smithy-go="v1.25.0"
github.com/aymerick/douceur="v0.2.0"
github.com/bep/clocks="v0.5.0"
github.com/bep/debounce="v1.2.0"
github.com/bep/gitmap="v1.9.0"
github.com/bep/goat="v0.5.0"
github.com/bep/godartsass/v2="v2.5.0"
github.com/bep/golibsass="v1.2.0"
github.com/bep/golocales="v0.2.0"
github.com/bep/goportabletext="v0.2.0"
github.com/bep/helpers="v0.12.0"
github.com/bep/imagemeta="v0.17.2"
github.com/bep/lazycache="v0.8.1"
github.com/bep/logg="v0.4.0"
github.com/bep/mclib="v1.20401.20400"
github.com/bep/overlayfs="v0.11.0"
github.com/bep/simplecobra="v0.7.0"
github.com/bep/textandbinarywriter="v0.1.0"
github.com/bep/tmc="v0.6.0"
github.com/bits-and-blooms/bitset="v1.24.5"
github.com/cespare/xxhash/v2="v2.3.0"
github.com/clbanning/mxj/v2="v2.7.0"
github.com/clipperhouse/displaywidth="v0.10.0"
github.com/clipperhouse/uax29/v2="v2.6.0"
github.com/cncf/xds/go="v0.0.0-20251210132809-ee656c7534f5"
github.com/cpuguy83/go-md2man/v2="v2.0.6"
github.com/dlclark/regexp2/v2="v2.2.1"
github.com/dustin/go-humanize="v1.0.1"
github.com/envoyproxy/go-control-plane/envoy="v1.36.0"
github.com/envoyproxy/protoc-gen-validate="v1.3.0"
github.com/evanw/esbuild="v0.28.1"
github.com/fatih/color="v1.18.0"
github.com/felixge/httpsnoop="v1.0.4"
github.com/frankban/quicktest="v1.14.6"
github.com/fsnotify/fsnotify="v1.9.0"
github.com/getkin/kin-openapi="v0.140.0"
github.com/go-jose/go-jose/v4="v4.1.4"
github.com/go-logr/logr="v1.4.3"
github.com/go-logr/stdr="v1.2.2"
github.com/go-openapi/jsonpointer="v0.22.5"
github.com/go-openapi/swag/jsonname="v0.25.5"
github.com/gobuffalo/flect="v1.0.3"
github.com/gobwas/glob="v0.2.3"
github.com/goccy/go-yaml="v1.19.2"
github.com/gohugoio/gift="v0.2.0"
github.com/gohugoio/go-i18n/v2="v2.1.3-0.20251018145728-cfcc22d823c6"
github.com/gohugoio/go-radix="v1.2.0"
github.com/gohugoio/hashstructure="v0.6.0"
github.com/gohugoio/httpcache="v0.8.0"
github.com/gohugoio/hugo-goldmark-extensions/extras="v0.7.0"
github.com/gohugoio/hugo-goldmark-extensions/passthrough="v0.5.0"
github.com/golang-jwt/jwt/v5="v5.3.0"
github.com/google/go-cmp="v0.7.0"
github.com/google/s2a-go="v0.1.9"
github.com/google/uuid="v1.6.0"
github.com/google/wire="v0.7.0"
github.com/googleapis/enterprise-certificate-proxy="v0.3.14"
github.com/googleapis/gax-go/v2="v2.21.0"
github.com/gorilla/css="v1.0.1"
github.com/gorilla/websocket="v1.5.3"
github.com/hairyhenderson/go-codeowners="v0.7.0"
github.com/hashicorp/golang-lru/v2="v2.0.7"
github.com/jdkato/prose="v1.2.1"
github.com/kr/pretty="v0.3.1"
github.com/kr/text="v0.2.0"
github.com/kylelemons/godebug="v1.1.0"
github.com/kyokomi/emoji/v2="v2.2.13"
github.com/makeworld-the-better-one/dither/v2="v2.4.0"
github.com/marekm4/color-extractor="v1.2.1"
github.com/mattn/go-colorable="v0.1.14"
github.com/mattn/go-isatty="v0.0.22"
github.com/mattn/go-runewidth="v0.0.19"
github.com/microcosm-cc/bluemonday="v1.0.27"
github.com/mitchellh/mapstructure="v1.5.1-0.20231216201459-8508981c8b6c"
github.com/muesli/smartcrop="v0.3.0"
github.com/niklasfasching/go-org="v1.9.1"
github.com/oasdiff/yaml3="v0.0.13"
github.com/oasdiff/yaml="v0.1.0"
github.com/olekukonko/cat="v0.0.0-20250911104152-50322a0618f6"
github.com/olekukonko/errors="v1.2.0"
github.com/olekukonko/ll="v0.1.6"
github.com/olekukonko/tablewriter="v1.1.4"
github.com/pbnjay/memory="v0.0.0-20210728143218-7b4eea64cf58"
github.com/pelletier/go-toml/v2="v2.4.3"
github.com/pkg/browser="v0.0.0-20240102092130-5ac0b6a4141c"
github.com/pkg/errors="v0.9.1"
github.com/rogpeppe/go-internal="v1.15.0"
github.com/russross/blackfriday/v2="v2.1.0"
github.com/santhosh-tekuri/jsonschema/v6="v6.0.2"
github.com/spf13/afero="v1.15.0"
github.com/spf13/cast="v1.10.0"
github.com/spf13/cobra="v1.10.2"
github.com/spf13/fsync="v0.10.1"
github.com/spf13/pflag="v1.0.10"
github.com/spiffe/go-spiffe/v2="v2.6.0"
github.com/tdewolff/minify/v2="v2.24.13"
github.com/tdewolff/parse/v2="v2.8.12"
github.com/tetratelabs/wazero="v1.12.0"
github.com/webmproject/libwebp="v1.6.0"
github.com/yuin/goldmark-emoji="v1.0.6"
github.com/yuin/goldmark="v1.8.2"
go.opentelemetry.io/auto/sdk="v1.2.1"
go.opentelemetry.io/contrib/detectors/gcp="v1.39.0"
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc="v0.67.0"
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp="v0.67.0"
go.opentelemetry.io/otel/metric="v1.43.0"
go.opentelemetry.io/otel/sdk/metric="v1.43.0"
go.opentelemetry.io/otel/sdk="v1.43.0"
go.opentelemetry.io/otel/trace="v1.43.0"
go.opentelemetry.io/otel="v1.43.0"
go.uber.org/automaxprocs="v1.5.3"
go.yaml.in/yaml/v3="v3.0.4"
gocloud.dev="v0.45.0"
golang.org/x/crypto="v0.53.0"
golang.org/x/image="v0.43.0"
golang.org/x/mod="v0.37.0"
golang.org/x/net="v0.56.0"
golang.org/x/oauth2="v0.36.0"
golang.org/x/sync="v0.21.0"
golang.org/x/sys="v0.46.0"
golang.org/x/text="v0.38.0"
golang.org/x/time="v0.15.0"
golang.org/x/tools="v0.47.0"
golang.org/x/xerrors="v0.0.0-20240903120638-7835f813f4da"
google.golang.org/api="v0.276.0"
google.golang.org/genproto/googleapis/api="v0.0.0-20260319201613-d00831a3d3e7"
google.golang.org/genproto/googleapis/rpc="v0.0.0-20260401024825-9d38bb4040a9"
google.golang.org/genproto="v0.0.0-20260319201613-d00831a3d3e7"
google.golang.org/grpc="v1.80.0"
google.golang.org/protobuf="v1.36.11"
rsc.io/qr="v0.2.0"
software.sslmate.com/src/go-pkcs12="v0.7.0"
```
</details>
