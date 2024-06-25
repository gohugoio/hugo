---
title: ToCSS
linkTitle: Transpile Sass to CSS
description: Transpile Sass to CSS.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    returnType: resource.Resource
    weight: 30
weight: 30
action:
  aliases: [toCSS]
  returnType: resource.Resource
  signatures: ['css.Sass [OPTIONS] RESOURCE']
toc: true
aliases: [/hugo-pipes/transform-to-css/]
---

## Usage

Transpile Sass to CSS using the LibSass transpiler included in Hugo's extended edition, or [install Dart Sass](#dart-sass) to use the latest features of the Sass language.

```go-html-template
{{ $opts := dict "transpiler" "libsass" "targetPath" "css/style.css" }}
{{ with resources.Get "sass/main.scss" | toCSS $opts | minify | fingerprint }}
  <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
{{ end }}
```

Sass has two forms of syntax: [SCSS] and [indented]. Hugo supports both.

[scss]: https://sass-lang.com/documentation/syntax#scss
[indented]: https://sass-lang.com/documentation/syntax#the-indented-syntax

## Options

transpiler
: (`string`) The transpiler to use, either `libsass` (default) or `dartsass`. Hugo's extended edition includes the LibSass transpiler. To use the Dart Sass transpiler, see the [installation instructions](#dart-sass) below.

targetPath
: (`string`) If not set, the transformed resource's target path will be the original path of the asset file with its extension replaced by `.css`.

vars {{< new-in 0.109.0 >}}
: (`map`) A map of key-value pairs that will be available in the `hugo:vars` namespace. Useful for [initializing Sass variables from Hugo templates](https://discourse.gohugo.io/t/42053/).

```scss
// LibSass
@import "hugo:vars";

// Dart Sass
@use "hugo:vars" as v;
```

outputStyle
: (`string`) Output styles available to LibSass include `nested` (default), `expanded`, `compact`, and `compressed`. Output styles available to Dart Sass include `expanded` (default) and `compressed`.

precision
: (`int`) Precision of floating point math. Not applicable to Dart Sass.

enableSourceMap
: (`bool`) If `true`, generates a source map.

sourceMapIncludeSources {{< new-in 0.108.0 >}}
: (`bool`) If `true`, embeds sources in the generated source map. Not applicable to LibSass.

includePaths
: (`slice`) A slice of paths, relative to the project root, that the transpiler will use when resolving `@use` and `@import` statements.

```go-html-template
{{ $opts := dict
  "transpiler" "dartsass"
  "targetPath" "css/style.css"
  "vars" site.Params.styles
  "enableSourceMap" (not hugo.IsProduction)
  "includePaths" (slice "node_modules/bootstrap/scss")
}}
{{ with resources.Get "sass/main.scss" | toCSS $opts | minify | fingerprint }}
  <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
{{ end }}
```

## Dart Sass

The extended version of Hugo includes [LibSass] to transpile Sass to CSS. In 2020, the Sass team deprecated LibSass in favor of [Dart Sass].

Use the latest features of the Sass language by installing Dart Sass in your development and production environments.

### Installation overview

Dart Sass is compatible with Hugo v0.114.0 and later.

If you have been using Embedded Dart Sass[^1] with Hugo v0.113.0 and earlier, uninstall Embedded Dart Sass, then install Dart Sass. If you have installed both, Hugo will use Dart Sass.

If you install Hugo as a [Snap package] there is no need to install Dart Sass. The Hugo Snap package includes Dart Sass.

[^1]: In 2023, the Sass team deprecated Embedded Dart Sass in favor of Dart Sass.

### Installing in a development environment

When you install Dart Sass somewhere in your PATH, Hugo will find it.

OS|Package manager|Site|Installation
:--|:--|:--|:--
Linux|Homebrew|[brew.sh]|`brew install sass/sass/sass`
Linux|Snap|[snapcraft.io]|`sudo snap install dart-sass`
macOS|Homebrew|[brew.sh]|`brew install sass/sass/sass`
Windows|Chocolatey|[chocolatey.org]|`choco install sass`
Windows|Scoop|[scoop.sh]|`scoop install sass`

You may also install [prebuilt binaries] for Linux, macOS, and Windows.

Run `hugo env` to list the active transpilers.

### Installing in a production environment

For [CI/CD] deployments (e.g., GitHub Pages, GitLab Pages, Netlify, etc.) you must edit the workflow to install Dart Sass before Hugo builds the site[^2]. Some providers allow you to use one of the package managers above, or you can download and extract one of the prebuilt binaries.

[^2]: You do not have to do this if (a) you have not modified the assets cache location, and (b) you have not set `useResourceCacheWhen` to `never` in your [site configuration], and (c) you add and commit your resources directory to your repository.

#### GitHub Pages

To install Dart Sass for your builds on GitHub Pages, add this step to the GitHub Pages workflow file:

```yaml
- name: Install Dart Sass
  run: sudo snap install dart-sass
```

If you are using GitHub Pages for the first time with your repository, GitHub provides a [starter workflow] for Hugo that includes Dart Sass. This is the simplest way to get started.

#### GitLab Pages

To install Dart Sass for your builds on GitLab Pages, the `.gitlab-ci.yml` file should look something like this:

```yaml
variables:
  HUGO_VERSION: 0.128.0
  DART_SASS_VERSION: 1.77.5
  GIT_DEPTH: 0
  GIT_STRATEGY: clone
  GIT_SUBMODULE_STRATEGY: recursive
  TZ: America/Los_Angeles
image:
  name: golang:1.20-buster
pages:
  script:
    # Install Dart Sass
    - curl -LJO https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz
    - tar -xf dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz
    - cp -r dart-sass/* /usr/local/bin
    - rm -rf dart-sass*
    # Install Hugo
    - curl -LJO https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
    - apt install -y ./hugo_extended_${HUGO_VERSION}_linux-amd64.deb
    - rm hugo_extended_${HUGO_VERSION}_linux-amd64.deb
    # Build
    - hugo --gc --minify
  artifacts:
    paths:
      - public
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

#### Netlify

To install Dart Sass for your builds on Netlify, the `netlify.toml` file should look something like this:

```toml
[build.environment]
HUGO_VERSION = "0.128.0"
DART_SASS_VERSION = "1.77.5"
TZ = "America/Los_Angeles"

[build]
publish = "public"
command = """\
  curl -LJO https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz && \
  tar -xf dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz && \
  rm dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz && \
  export PATH=/opt/build/repo/dart-sass:$PATH && \
  hugo --gc --minify \
  """
```

### Example

To transpile with Dart Sass, set `transpiler` to `dartsass` in the options map passed to `css.Sass`. For example:

```go-html-template
{{ $opts := dict "transpiler" "dartsass" "targetPath" "css/style.css" }}
{{ with resources.Get "sass/main.scss" | toCSS $opts | minify | fingerprint }}
  <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
{{ end }}
```

### Miscellaneous

If you build Hugo from source and run `mage test -v`, the test will fail if you install Dart Sass as a Snap package. This is due to the Snap package's strict confinement model.

[brew.sh]: https://brew.sh/
[chocolatey.org]: https://community.chocolatey.org/packages/sass
[ci/cd]: https://en.wikipedia.org/wiki/CI/CD
[dart sass]: https://sass-lang.com/dart-sass
[libsass]: https://sass-lang.com/libsass
[prebuilt binaries]: https://github.com/sass/dart-sass/releases/latest
[scoop.sh]: https://scoop.sh/#/apps?q=sass
[site configuration]: /getting-started/configuration/#configure-build
[snap package]: /installation/linux/#snap
[snapcraft.io]: https://snapcraft.io/dart-sass
[starter workflow]: https://github.com/actions/starter-workflows/blob/main/pages/hugo.yml
