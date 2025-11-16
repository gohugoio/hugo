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

To use Dart Sass with Hugo on a CI/CD platform like GitHub Pages, GitLab Pages, or Netlify, you typically must modify your build workflow to install Dart Sass before the Hugo site build begins. This is because these platforms don't have Dart Sass pre-installed, and Hugo needs it to process your Sass files.

There's one key exception where you can skip this step: you have committed your `resources` directory to your repository. This is only possible if:

- You have not changed Hugo's default asset cache location.
- You have not set [`useResourceCacheWhen`] to never in your sites configuration.

By committing the `resources` directory, you're providing the pre-built CSS files directly to your CI/CD service, so it doesn't need to run the Sass compilation itself.

For examples of how to install Dart Sass in a production environment, see the following workflow files:

- [GitHub Pages]
- [GitLab Pages]
- [Netlify]

[`publishDir`]: /configuration/all/#publishdir
[`useResourceCacheWhen`]: /configuration/build/#useresourcecachewhen
[brew.sh]: https://brew.sh/
[chocolatey.org]: https://community.chocolatey.org/packages/sass
[dart sass]: https://sass-lang.com/dart-sass
[GitHub Pages]: /host-and-deploy/host-on-github-pages/#step-7
[GitLab Pages]: /host-and-deploy/host-on-gitlab-pages/#configure-gitlab-cicd
[libsass]: https://sass-lang.com/libsass
[Netlify]: /host-and-deploy/host-on-netlify/#configuration-file
[prebuilt binaries]: https://github.com/sass/dart-sass/releases/latest
[scoop.sh]: https://scoop.sh/#/apps?q=sass
[snap package]: /installation/linux/#snap
[snapcraft.io]: https://snapcraft.io/dart-sass
