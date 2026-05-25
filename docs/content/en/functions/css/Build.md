---
title: css.Build
description: Bundle, transform, and minify CSS resources.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: ['css.Build [OPTIONS] RESOURCE']
---

{{< new-in 0.158.0 />}}

> [!note]
> The `css.Build` function is backed by the [`evanw/esbuild`][] package, providing a mature, high-performance foundation for bundling, transformation, and minification.

Use the `css.Build` function to:

- Recursively replace `@import` statements in CSS files with the content of the imported files
- Transform syntax for browser compatibility
- Apply vendor prefixes for browser compatibility
- Minify the bundled CSS code
- Create a source map

If an `@import` statement includes a media query, a feature query, or a cascade layer assignment, the function wraps the imported content in the corresponding `@media`, `@supports`, or `@layer` rule.

## Usage

In this example, Hugo bundles the local files referenced by `@import` statements to create and publish a single resource with inline content.

```text
assets/
└── css/
    ├── components/
    │   ├── a.css
    │   └── b.css
    └── main.css
```

```css {file="assets/css/main.css" copy=true}
@import url('https://cdn.jsdelivr.net/npm/the-new-css-reset/css/reset.min.css');

@import './components/a.css';
@import './components/b.css';

.c {color: blue; }
```

```css {file="assets/css/components/a.css" copy=true}
.a { color: red; }
```

```css {file="assets/css/components/b.css" copy=true}
.b { color: green; }
```

```go-html-template {file="layouts/_partials/css.html" copy=true}
{{ with resources.Get "css/main.css" | css.Build }}
  {{ if hugo.IsDevelopment }}
    <link rel="stylesheet" href="{{ .RelPermalink }}">
  {{ else }}
    {{ with . | fingerprint }}
      <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
    {{ end }}
  {{ end }}
{{ end }}
```

```go-html-template {file="layouts/baseof.html" copy=true}
{{ partialCached "css.html" . }}
```

The generated CSS code:

```css {file="public/css/main.css"}
@import "https://cdn.jsdelivr.net/npm/the-new-css-reset/css/reset.min.css";

.a {
  color: red;
}

.b {
  color: green;
}

.c {
  color: blue;
}
```

To minify the generated CSS code, use the [`minify`](#minify) option as described below.

## Options

The `css.Build` function takes an optional map of options to fine-tune bundling, minification, and browser compatibility.

externals
: (`[]string`) A slice of path patterns to exclude from bundling. The `@import` statements for these patterns remain as-is in the generated CSS code. See&nbsp;[details][esb_external].

  ```go-html-template
  {{ $opts := dict "externals" (slice "./exclude-these/*" "./exclude-these-too/*") }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

loaders
: (`map`) A map of file extensions to loader types. This determines how files with a given extension are processed during bundling. By default, Hugo uses the `css` loader for `.css` files and the `file` loader for all others. Common loaders include:

  - `css`: Processes the file as a CSS file
  - `dataurl`: Embeds the file as a base64-encoded data URL
  - `empty`: Excludes the file from the bundle
  - `file`: Copies the file to the output directory and rewrites the URL
  - `text`: Loads the file content as a string

  See&nbsp;[details][esb_loader].

  ```go-html-template
  {{ $opts := dict "loaders" (dict ".png" "dataurl" ".svg" "dataurl") }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

mainFields
: (`[]string`) A prioritized slice of field names in a `package.json` file that determine the CSS entry point of a Node package. The default is `["style", "main"]`. See&nbsp;[details][esb_mainfields].

  When an `@import` statement references a Node package, Hugo consults the metadata in the `package.json` file to find the stylesheet. Use this option to support packages that define a CSS entry point using non-standard fields.

  ```go-html-template
  {{ $opts := dict "mainFields" (slice "css" "style" "main") }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

minify
: (`bool`) Whether to minify the generated CSS code. Default is `false`. See&nbsp;[details][esb_minify].

  ```go-html-template
  {{ $opts := dict "minify" true }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

sourceMap
: (`string`) The type of source map to generate. One of `external`, `inline`, `linked`, or `none`. Default is `none`. See&nbsp;[details][esb_sourcemap].

  ```go-html-template
  {{ $opts := dict "sourceMap" "linked" }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

sourcesContent
: (`bool`) Whether to include the content of the source files in the source map. Default is `true`. See&nbsp;[details][esb_sourcesContent].

  ```go-html-template
  {{ $opts := dict "sourceMap" "linked" "sourcesContent" false }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

target
: (`[]string`) The target environment for the generated CSS code. This determines which syntax transformations to perform and which vendor prefixes to apply. If unset, no transformations or prefixing are performed. Each element consists of a target name and a version number. Supported targets include `chrome`, `edge`, `firefox`, `ie`, `ios`, `opera`, and `safari`. See&nbsp;[details][esb_target].

  ```go-html-template
  {{ $target := slice "chrome115" "edge115" "firefox116" "ios16.4" "opera101" "safari16.4" }}
  {{ $opts := dict "target" $target }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

  In the example above, the target environment is roughly equivalent to the [browserlist][] "baseline widely available" profile as of March 2026.

targetPath
: (`string`) The path to the generated CSS file, relative to the project's [`publishDir`][]. If unset, this defaults to the asset's original path with a `.css` extension.

  ```go-html-template
  {{ $opts := dict "targetPath" "css/styles.css" }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

vars
: {{< new-in 0.160.0 />}}
: (`map`) A map of key-value pairs used to generate CSS variables. The `css.Build` function injects these variables into the stylesheet when it encounters the `hugo:vars` internal identifier within an `@import` statement.
  
  ```go-html-template
  {{ $vars := dict
    "font-family" "\"Times New Roman\", Times, serif"
    "font-size" "24px" 
    "primary-color" "blue" 
  }}
  {{ $opts := dict "vars" $vars }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

  In the example above, using the identifier in your CSS allows you to access the values using standard CSS variable syntax.

  ```css
  @import 'hugo:vars';

  .element {
    color: var(--primary-color);
    font-family: var(--font-family);
    font-size: var(--font-size);
  }
  ```

  The above produces output equivalent to:

  ```css
  :root {
    --font-family:
      "Times New Roman",
      Times,
      serif;
    --font-size: 24px;
    --primary-color: blue;
  }

  .element {
    color: var(--primary-color);
    font-family: var(--font-family);
    font-size: var(--font-size);
  }
````

  {{< new-in 0.161.0 />}}

  The map may optionally contain nested maps. Each nested map is exposed as a separate `hugo:vars/<name>` namespace, where `<name>` is the key of the nested map (lowercased). Top-level scalar values and nested maps are independent. A top-level `@import 'hugo:vars'` only includes scalar values, while `@import 'hugo:vars/<name>'` only includes the scalars from the named nested map.

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
  {{ $opts := dict "vars" $vars }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

  Because nested imports follow the same rules as regular `@import` statements, you can attach a media query, feature query, or cascade layer assignment to a `hugo:vars/<name>` import.

  ```css
  @import 'hugo:vars';
  @import 'hugo:vars/mobile' (max-width: 650px);

  body {
    background-color: var(--primary-color);
    font-family: var(--font-family);
  }
  ```

  The above produces output equivalent to:

  ```css
  :root {
    --font-family: "Times New Roman", Times, serif;
    --font-size: 24px;
    --primary-color: blue;
  }

  @media (max-width: 650px) {
    :root {
      --font-size: 12px;
      --primary-color: red;
    }
  }

  body {
    background-color: var(--primary-color);
    font-family: var(--font-family);
  }
  ```

  The `vars` option is useful for setting CSS variables within your project configuration.

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
  {{ $opts := dict "vars" site.Params.theme.style }}
  {{ $r := resources.Get "css/main.css" | css.Build $opts }}
  ```

  When passing a `vars` map to the `css.Build` function, you can use the [`css.Quoted`][] function to explicitly indicate that a value must be treated as a quoted string, most commonly for `font-family` names or the `content` property.

  > [!note]
  > If you're using TailwindCSS and want to use the `vars` option to inject CSS variables, see [this section in the TailwindCSS documentation](./TailwindCSS.md#inject-css-variables-with-vars).

## Example

The example below uses several of the [options](#options) described above to bundle, transform, and minify CSS code.

```go-html-template {file="layouts/_partials/css.html" copy=true}
{{ with resources.Get "css/main.css" }}
  {{ $opts := dict
    "loaders" (dict ".png" "dataurl" ".svg" "dataurl")
    "minify" (cond hugo.IsDevelopment false true)
    "sourceMap" (cond hugo.IsDevelopment "linked" "none")
    "target" (slice "chrome115" "edge115" "firefox116" "ios16.4" "opera101" "safari16.4")
    "targetPath" "css/styles.css"
  }}
  {{ with . | css.Build $opts }}
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

Using the options above, Hugo does the following:

- Embeds PNG and SVG images as data URLs in the generated CSS code
- Minifies the output in production but not in development
- Generates an external source map in development but not in production
- Transforms syntax for compatibility with the targeted browser versions
- Adds vendor prefixes for compatibility with the targeted browser versions
- Publishes the generated CSS code to `css/styles.css`
- In production, adds an SRI hash and inserts a file hash into the filename

## Common patterns

The examples below cover the most frequent use cases for referencing resources within your project or within Node packages. These patterns apply to both `@import` statements and the `url()` functional notation used for images and fonts.

All resources referenced by a path, including images, fonts, and stylesheets, must reside in the `assets` directory of the [unified file system](g), or within a Node package.

### Files in the assets directory

To include a stylesheet from the `assets` directory, you can use a bare path, a relative path, or a root-relative path. When you use a bare path, Hugo searches relative to the current stylesheet, then relative to the `assets` directory.

```css {file="/assets/css/main.css"}
/* A bare path */
@import "variables.css";

/* A relative path */
@import "./theme.css";
@import "../layout.css";

/* A root-relative path */
@import "/css/grid.css";

/* A url() reference using the same resolution logic */
.logo { background: url("/images/logo.svg"); }
```

### Node packages

When referencing a Node package by name, Hugo consults the `package.json` file within that package to find the entry point.

```css {file="/assets/css/main.css"}
@import "bootstrap";
```

### Files within a package

To reference a specific file within a Node package, provide the path starting with the package name.

```css {file="/assets/css/main.css"}
@import "bootstrap/dist/css/bootstrap-grid.css";
```

[`css.Quoted`]: /functions/css/quoted/
[`evanw/esbuild`]: https://github.com/evanw/esbuild
[`publishDir`]: /configuration/all/#publishdir
[browserlist]: https://browsersl.ist
[esb_external]: https://esbuild.github.io/api/#external
[esb_loader]: https://esbuild.github.io/api/#loader
[esb_mainfields]: https://esbuild.github.io/api/#main-fields
[esb_minify]: https://esbuild.github.io/api/#minify
[esb_sourcemap]: https://esbuild.github.io/api/#sourcemap
[esb_sourcesContent]: https://esbuild.github.io/api/#sources-content
[esb_target]: https://esbuild.github.io/api/#target
