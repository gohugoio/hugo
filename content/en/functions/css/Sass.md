---
title: css.Sass
description: Transpiles Sass to CSS.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [toCSS]
    returnType: resource.Resource
    signatures: ['css.Sass [OPTIONS] RESOURCE']
---

{{< new-in 0.128.0 />}}

Transpile Sass to CSS using the LibSass transpiler included in Hugo's extended and extended/deploy editions, or [install Dart Sass](#dart-sass) to use the latest features of the Sass language.

Sass has two forms of syntax: [SCSS] and [indented]. Hugo supports both.

[scss]: https://sass-lang.com/documentation/syntax#scss
[indented]: https://sass-lang.com/documentation/syntax#the-indented-syntax

## Options

enableSourceMap
: (`bool`) Whether to generate a source map. Default is `false`.

includePaths
: (`slice`) A slice of paths, relative to the project root, that the transpiler will use when resolving `@use` and `@import` statements.

outputStyle
: (`string`) The output style of the resulting CSS. With LibSass, one of `nested` (default), `expanded`, `compact`, or `compressed`. With Dart Sass, either `expanded` (default) or `compressed`.

precision
: (`int`) The precision of floating point math. Applicable to LibSass. Default is `8`.

silenceDeprecations
: {{< new-in 0.139.0 />}}
: (`slice`) A slice of deprecation IDs to silence. IDs are enclosed in brackets within Dart Sass warning messages (e.g., `import` in `WARN Dart Sass: DEPRECATED [import]`). Applicable to Dart Sass. Default is `false`.

silenceDependencyDeprecations
: {{< new-in 0.146.0 />}}
: (`bool`) Whether to silence deprecation warnings from dependencies, where a dependency is considered any file transitively imported through a load path. This does not apply to `@warn` or `@debug` rules.Default is `false`.

sourceMapIncludeSources
: (`bool`) Whether to embed sources in the generated source map. Applicable to Dart Sass. Default is `false`.

targetPath
: (`string`) The publish path for the transformed resource, relative to the[`publishDir`]. If unset, the target path defaults to the asset's original path with a `.css` extension.

transpiler
: (`string`) The transpiler to use, either `libsass` or `dartsass`. Hugo's extended and extended/deploy editions include the LibSass transpiler. To use the Dart Sass transpiler, see the [installation instructions](#dart-sass). Default is `libsass`.

vars
: (`map`) A map of key-value pairs that will be available in the `hugo:vars` namespace. Useful for [initializing Sass variables from Hugo templates](https://discourse.gohugo.io/t/42053/).

  ```scss
  // LibSass
  @import "hugo:vars";

  // Dart Sass
  @use "hugo:vars" as v;
  ```

## Example

```go-html-template {copy=true}
{{ with resources.Get "sass/main.scss" }}
  {{ $opts := dict
    "enableSourceMap" hugo.IsDevelopment
    "outputStyle" (cond hugo.IsDevelopment "expanded" "compressed")
    "targetPath" "css/main.css"
    "transpiler" "dartsass"
    "vars" site.Params.styles
    "includePaths" (slice "node_modules/bootstrap/scss")
  }}
  {{ with . | toCSS $opts }}
    {{ if hugo.IsDevelopment }}
      <link rel="stylesheet" href="{{ .RelPermalink }}">
    {{ else }}
      {{ with . | fingerprint }}
        <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
```

## Dart Sass

Hugo's extended and extended/deploy editions include [LibSass] to transpile Sass to CSS. In 2020, the Sass team deprecated LibSass in favor of [Dart Sass].

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

You may also install [prebuilt binaries] for Linux, macOS, and Windows. You must install the prebuilt binary outside of your project directory and ensure its path is included in your system's PATH environment variable.

Run `hugo env` to list the active transpilers.

> [!note]
> If you build Hugo from source and run `mage test -v`, the test will fail if you install Dart Sass as a Snap package. This is due to the Snap package's strict confinement model.

### Installing in a production environment

For [CI/CD](g) deployments (e.g., GitHub Pages, GitLab Pages, Netlify, etc.) you must edit the workflow to install Dart Sass before Hugo builds the site[^2]. Some providers allow you to use one of the package managers above, or you can download and extract one of the prebuilt binaries.

[^2]: You do not have to do this if (a) you have not modified the assets cache location, and (b) you have not set `useResourceCacheWhen` to `never` in your [site configuration], and (c) you add and commit your `resources` directory to your repository.

#### GitHub Pages

To install Dart Sass for your builds on GitHub Pages, add this step to the GitHub Pages workflow file:

```yaml
- name: Install Dart Sass
  run: sudo snap install dart-sass
```

#### GitLab Pages

To install Dart Sass for your builds on GitLab Pages, the `.gitlab-ci.yml` file should look something like this:

```yaml
variables:
  HUGO_VERSION: 0.147.9
  DART_SASS_VERSION: 1.89.2
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
HUGO_VERSION = "0.147.9"
DART_SASS_VERSION = "1.89.2"
NODE_VERSION = "22"
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

[brew.sh]: https://brew.sh/
[chocolatey.org]: https://community.chocolatey.org/packages/sass
[dart sass]: https://sass-lang.com/dart-sass
[libsass]: https://sass-lang.com/libsass
[prebuilt binaries]: https://github.com/sass/dart-sass/releases/latest
[scoop.sh]: https://scoop.sh/#/apps?q=sass
[site configuration]: /configuration/build/
[snap package]: /installation/linux/#snap
[snapcraft.io]: https://snapcraft.io/dart-sass
[starter workflow]: https://github.com/actions/starter-workflows/blob/main/pages/hugo.yml
[`publishDir`]: /configuration/all/#publishdir
