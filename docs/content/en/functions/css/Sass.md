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
aliases: [/functions/resources/tocss/]
---

Transpile Sass to CSS using the LibSass transpiler included in Hugo's extended and extended/deploy editions, or [install Dart Sass](#dart-sass) to use the latest features of the Sass language.

> [!warning]
> The embedded LibSass transpiler was deprecated in [v0.153.0][] and will be removed in a future release. Use the Dart Sass transpiler instead by setting the `transpiler` option to `dartsass` as shown in the examples below.

Sass has two forms of syntax: [SCSS][] and [indented][]. Hugo supports both.

## Options

enableSourceMap
: (`bool`) Whether to generate a source map. Default is `false`.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "enableSourceMap" true 
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

includePaths
: (`slice`) A slice of paths, relative to the project root, that the transpiler will use when resolving `@use` and `@import` statements.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "includePaths" (slice "node_modules/bootstrap/scss")
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

outputStyle
: (`string`) The output style of the resulting CSS. With LibSass, one of `nested` (default), `expanded`, `compact`, or `compressed`. With Dart Sass, either `expanded` (default) or `compressed`.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "outputStyle" "compressed"
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

precision
: (`int`) The precision of floating point math. Applicable to LibSass. Default is `8`.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "precision" 10 
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

silenceDeprecations
: {{< new-in 0.139.0 />}}
: (`slice`) A slice of deprecation IDs to silence. IDs are enclosed in brackets within Dart Sass warning messages (e.g., `import` in `WARN Dart Sass: DEPRECATED [import]`). Applicable to Dart Sass.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "silenceDeprecations" (slice "import")
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

silenceDependencyDeprecations
: {{< new-in 0.146.0 />}}
: (`bool`) Whether to silence deprecation warnings from dependencies, where a dependency is considered any file transitively imported through a load path. This does not apply to `@warn` or `@debug` rules. Default is `false`.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "silenceDependencyDeprecations" true
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

sourceMapIncludeSources
: (`bool`) Whether to embed sources in the generated source map. Applicable to Dart Sass. Default is `false`.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "enableSourceMap" true "sourceMapIncludeSources" true
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

targetPath
: (`string`) The publish path for the transformed resource, relative to the [`publishDir`][]. If unset, the target path defaults to the asset's original path with a `.css` extension.

  ```go-html-template
  {{ $opts := dict
    "transpiler" "dartsass"
    "targetPath" "css/bundle.css"
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

transpiler
: (`string`) The transpiler to use, either `libsass` or `dartsass`. Hugo's extended and extended/deploy editions include the LibSass transpiler. To use the Dart Sass transpiler, see the [installation instructions](#dart-sass). Default is `libsass`.

  ```go-html-template
  {{ $opts := dict "transpiler" "dartsass" }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

vars
: (`map`) A map of key-value pairs used to generate Sass variables. The `css.Sass` function injects these variables into the stylesheet when it encounters the `hugo:vars` internal identifier within a `@use` or `@import` statement.

  ```go-html-template
  {{ $vars := dict
    "font-family" "\"Times New Roman\", Times, serif"
    "font-size" "24px" 
    "primary-color" "blue"
  }}
  {{ $opts := dict 
    "transpiler" "dartsass"
    "vars" $vars
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

  In the example above, using the identifier in your stylesheet allows you to access the values as Sass variables in the `hugo:vars` namespace:

  ```scss
  @use 'hugo:vars' as v;

  .element {
    color: v.$primary-color;
    font-family: v.$font-family;
    font-size: v.$font-size;
  }
  ```

  The above produces output equivalent to:

  ```css
  .element {
    color: blue;
    font-family: "Times New Roman", Times, serif;
    font-size: 24px;
  }
  ```

  {{< new-in 0.161.0 />}}

  The map may optionally contain nested maps. Each nested map is exposed as a separate `hugo:vars/<name>` namespace, where `<name>` is the key of the nested map (lowercased). Top-level scalar values and nested maps are independent. A top-level `@use 'hugo:vars'` only includes scalar values, while `@use 'hugo:vars/<name>'` only includes the scalars from the named nested map.

  ```go-html-template
  {{ $vars := dict
    "font-family" "\"Times New Roman\", Times, serif"
    "font-size" "24px"
    "primary-color" "blue"
    "mobile" (dict 
      "font-size" "12px"
      "primary-color" "red"
    )
  }}
  {{ $opts := dict 
    "transpiler" "dartsass"
    "vars" $vars
  }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

  In the stylesheet, reference each nested namespace with a separate `@use` statement. Assign an alias to access the variables from that namespace:

  ```scss
  @use 'hugo:vars' as v;
  @use 'hugo:vars/mobile' as mobile;

  body {
    color: v.$primary-color;
    font-family: v.$font-family;
    font-size: v.$font-size;
  }

  @media (max-width: 650px) {
    body {
      color: mobile.$primary-color;
      font-size: mobile.$font-size;
    }
  }
  ```

  The above produces output equivalent to:

  ```css
  body {
    color: blue;
    font-family: "Times New Roman", Times, serif;
    font-size: 24px;
  }

  @media (max-width: 650px) {
    body {
      color: red;
      font-size: 12px;
    }
  }
  ```

  The `vars` option is useful for setting Sass variables within your project configuration.

  {{< code-toggle file=hugo >}}
  [params.theme.style]
  font-family = '"Times New Roman", Times, serif'
  font-size = '24px'
  primary-color = 'blue'

  [params.theme.style.mobile]
  font-size = '12px'
  primary-color = 'red'
  {{< /code-toggle >}}

  ```go-html-template
  {{ $opts := dict 
    "transpiler" "dartsass"
    "vars" site.Params.theme.style }}
  {{ $r := resources.Get "sass/main.scss" | css.Sass $opts }}
  ```

  When passing a `vars` map to the `css.Sass` function, Hugo detects common typed CSS values such as `24px` or `#FF0000` using regular expression matching. If necessary, you can bypass automatic type inference by using the [`css.Quoted`][] or [`css.Unquoted`][] function to explicitly indicate a value's type.

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
  {{ with . | css.Sass $opts }}
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

Hugo's extended and extended/deploy editions include [LibSass][] to transpile Sass to CSS. In 2020, the Sass team deprecated LibSass in favor of [Dart Sass][].

Use the latest features of the Sass language by installing Dart Sass in your development and production environments.

### Installation overview

Dart Sass is compatible with Hugo v0.114.0 and later.

If you have been using Embedded Dart Sass[^1] with Hugo v0.113.0 and earlier, uninstall Embedded Dart Sass, then install Dart Sass. If you have installed both, Hugo will use Dart Sass.

If you install Hugo as a [Snap package][] there is no need to install Dart Sass. The Hugo Snap package includes Dart Sass.

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

You may also install [prebuilt binaries][] for Linux, macOS, and Windows. You must install the prebuilt binary outside of your project directory and ensure its path is included in your system's PATH environment variable.

Run `hugo env` to list the active transpilers.

> [!note]
> If you build Hugo from source and run `mage test -v`, the test will fail if you install Dart Sass as a Snap package. This is due to the Snap package's strict confinement model.

### Installing in a production environment

To use Dart Sass with Hugo on a [CI/CD](g) platform, you typically must modify your build workflow to install Dart Sass before the Hugo site build begins. This is because these platforms don't have Dart Sass pre-installed, and Hugo needs it to process your Sass files.

There's one key exception where you can skip this step: you have committed your `resources` directory to your repository. This is only possible if:

- You have not changed Hugo's default asset cache location.
- You have not set [`useResourceCacheWhen`][] to never in your project configuration.

By committing the `resources` directory, you're providing the pre-built CSS files directly to your CI/CD platform, so it doesn't need to run the Sass compilation itself.

For examples of how to install Dart Sass in a production environment, see these hosting guides:

- [Cloudflare][]
- [GitHub Pages][]
- [GitLab Pages][]
- [Netlify][]
- [Render][]
- [SourceHut][]
- [Vercel][]

[`css.Quoted`]: /functions/css/quoted/
[`css.Unquoted`]: /functions/css/unquoted/
[`publishDir`]: /configuration/all/#publishdir
[`useResourceCacheWhen`]: /configuration/build/#useresourcecachewhen
[brew.sh]: https://brew.sh/
[chocolatey.org]: https://community.chocolatey.org/packages/sass
[Cloudflare]: /host-and-deploy/host-on-cloudflare/
[Dart Sass]: https://sass-lang.com/dart-sass/
[GitHub Pages]: /host-and-deploy/host-on-github-pages/
[GitLab Pages]: /host-and-deploy/host-on-gitlab-pages/
[indented]: https://sass-lang.com/documentation/syntax#the-indented-syntax
[LibSass]: https://sass-lang.com/libsass
[Netlify]: /host-and-deploy/host-on-netlify/
[prebuilt binaries]: https://github.com/sass/dart-sass/releases/latest
[Render]: /host-and-deploy/host-on-render/
[scoop.sh]: https://scoop.sh/#/apps?q=sass
[SCSS]: https://sass-lang.com/documentation/syntax#scss
[Snap package]: https://snapcraft.io/hugo
[snapcraft.io]: https://snapcraft.io/dart-sass
[SourceHut]: /host-and-deploy/host-on-sourcehut-pages/
[v0.153.0]: https://github.com/gohugoio/hugo/releases/tag/v0.153.0
[Vercel]: /host-and-deploy/host-on-vercel/
