---
title: css.ChromaStyles
description: Returns a CSS stylesheet for the syntax highlighter.
categories: []
keywords: [highlight]
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: ['css.ChromaStyles OPTIONS']
---

{{< new-in 0.165.0 />}}

The `css.ChromaStyles` function returns a CSS stylesheet for the syntax highlighter as a `Resource` object. This stylesheet is needed when the `noClasses` option is `false`, either as a [global default][] in your project configuration or as an override when using any of the following:

- The [`highlight`][] shortcode
- The [`transform.Highlight`][] function
- The [`transform.HighlightCodeBlock`][] function
- Highlight options in the [info string](g) of a fenced code block in Markdown

Hugo caches the result, so calling the function multiple times with the same options incurs no additional overhead.

## Options

The `css.ChromaStyles` function requires an options map. The [`targetPath`](#targetpath) is the only required option.

`classDark`
: (`string`) When [`modeSelector`](#modeselector) is `true` and [`mode`](#mode) is `dark`, the CSS class name used to scope selectors. Default is `dark`.

`classLight`
: (`string`) When [`modeSelector`](#modeselector) is `true` and [`mode`](#mode) is `light`, the CSS class name used to scope selectors. Default is `light`.

`highlightStyle`
: (`string`) The foreground and background colors for highlighted lines, such as `#fff000 bg:#000fff`. Defaults to the colors defined by the selected [`style`](#style).

`lineNumbersInlineStyle`
: (`string`) The foreground and background colors for inline line numbers, such as `#fff000 bg:#000fff`. Defaults to the colors defined by the selected [`style`](#style).

`lineNumbersTableStyle`
: (`string`) The foreground and background colors for table line numbers, such as `#fff000 bg:#000fff`. Defaults to the colors defined by the selected [`style`](#style).

`mode`
: (`string`) The color [mode][mode], either `light` or `dark`. The specified style must support the given mode. If you omit this option, Hugo uses the style's default mode.

`modeSelector`
: (`bool`) Whether to scope CSS selectors under a top-level mode class. For example, a `light` mode stylesheet scopes selectors under `.light`, producing `.light .chroma` instead of `.chroma`. Set to `true` when generating a paired light/dark stylesheet. Default is `false`.

`omitClassComments`
: (`bool`) Whether to omit CSS class comment prefixes in the generated stylesheet. Default is `false`.

`style`
: (`string`) The syntax highlighting style. Defaults to the [`style`][] value in your project configuration. See [syntax highlighting styles][] for a list of available styles.

`targetPath`
: (`string`) The target path of the resource, relative to the [`publishDir`][]. Required.

## Examples

In the first two examples, call the _partial_ template from your _base_ template using the [`partials.IncludeCached`][] function. Both examples assume this project configuration:

{{< code-toggle file=hugo copy=true >}}
[markup.highlight]
noClasses = false
style = 'github'
{{< /code-toggle >}}

### Single stylesheet

To generate and include a stylesheet using the [`style`][] value in your project configuration:

```go-html-template {file="layouts/_partials/highlight.html" copy=true}
{{ $opts := dict "targetPath" "css/highlight.css" }}
{{ with css.ChromaStyles $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

### Light and dark stylesheets

To generate and include paired light and dark stylesheets for a style that supports both [modes][mode]:

```go-html-template {file="layouts/_partials/highlight.html" copy=true}
{{ $opts := dict
  "mode" "light"
  "targetPath" "css/highlight-light.css"
}}
{{ with css.ChromaStyles $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}

{{ $opts := dict
  "mode" "dark"
  "modeSelector" true
  "targetPath" "css/highlight-dark.css"
}}
{{ with css.ChromaStyles $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

The light stylesheet is unscoped and acts as the default. The dark stylesheet's selectors are scoped under the `dark` class, so its rules take effect only when the root element has that class.

### Complete example

This example adds a light/dark/system theme switcher to your site, using the [`css.Build`][] function to bundle the generated stylesheets into your main CSS file.

Step 1
: Add this to your project configuration:

  {{< code-toggle file=hugo copy=true >}}
  [markup.highlight]
  noClasses = false
  style = 'github'
  {{< /code-toggle >}}

Step 2
: Create a CSS entry file with `@import` statements for the generated stylesheets, plus light and dark rules for the rest of the page:

  ```css {file="assets/css/main.css" copy=true}
  @import "./highlight-light.css";
  @import "./highlight-dark.css";

  html {
    background-color: #fff;
    color: #000;
    color-scheme: light;
  }

  a {
    color: #00e;
  }

  html.dark {
    background-color: #000;
    color: #fff;
    color-scheme: dark;
  }

  html.dark a {
    color: #6af;
  }
  ```

Step 3
: Create a _partial_ template to generate the stylesheets and bundle them with the CSS entry file:

  ```go-html-template {file="layouts/_partials/css.html" copy=true}
  {{ $opts := dict
    "mode" "light"
    "targetPath" "css/highlight-light.css"
  }}
  {{ $highlightLight := css.ChromaStyles $opts }}

  {{ $opts := dict
    "mode" "dark"
    "modeSelector" true
    "targetPath" "css/highlight-dark.css"
  }}
  {{ $highlightDark := css.ChromaStyles $opts }}

  {{ with resources.Get "css/main.css" }}
    {{ $opts := dict
      "importContext" (slice $highlightLight $highlightDark)
      "minify" (cond hugo.IsDevelopment false true)
      "sourceMap" (cond hugo.IsDevelopment "linked" "none")
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

  The [`importContext`][] option makes the generated stylesheets available when `css.Build` resolves the `@import` statements in the CSS entry file.

Step 4
: Create a JavaScript file to control the theme switcher, persisting manual selections in local storage and following the operating system's preference when the visitor selects the system option:

  ```js {file="assets/js/main.js" copy=true}
  const root = document.documentElement;
  const fieldset = document.getElementById('theme-switcher');
  const mq = window.matchMedia('(prefers-color-scheme: dark)');

  const applyTheme = (theme) => {
    if (theme === 'dark') {
      root.classList.add('dark');
    } else if (theme === 'light') {
      root.classList.remove('dark');
    } else {
      root.classList.toggle('dark', mq.matches);
    }
    fieldset.querySelector(`input[value="${theme ?? 'system'}"]`).checked = true;
  };

  // Keep class in sync with system preference when no manual override is set.
  mq.addEventListener('change', () => {
    if (!localStorage.getItem('theme')) {
      root.classList.toggle('dark', mq.matches);
    }
  });

  applyTheme(localStorage.getItem('theme'));

  fieldset.addEventListener('change', (e) => {
    const next = e.target.value === 'system' ? null : e.target.value;
    if (next === null) {
      localStorage.removeItem('theme');
    } else {
      localStorage.setItem('theme', next);
    }
    applyTheme(next);
  });
  ```

Step 5
: Create a _partial_ template to process the JavaScript:

  ```go-html-template {file="layouts/_partials/js.html" copy=true}
  {{ with resources.Get "js/main.js" }}
    {{ $opts := dict
      "minify" (cond hugo.IsDevelopment false true)
      "sourceMap" (cond hugo.IsDevelopment "linked" "none")
    }}
    {{ with . | js.Build $opts }}
      {{ if hugo.IsDevelopment }}
        <script defer src="{{ .RelPermalink }}"></script>
      {{ else }}
        {{ with . | fingerprint }}
          <script defer src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
        {{ end }}
      {{ end }}
    {{ end }}
  {{ end }}
  ```

Step 6
: Call both _partial_ templates from your _base_ template and add the theme switcher markup. The inline script in the `head` element applies the theme before the first paint, preventing a flash of light-themed content when a dark-theme visitor loads a page:

  ```go-html-template {file="layouts/baseof.html" copy=true}
  <!DOCTYPE html>
  <html lang="{{ site.Language.Locale }}">
    <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>{{ site.Title }}</title>
      <script>
        const theme = localStorage.getItem('theme');
        document.documentElement.classList.toggle('dark',
          theme === 'dark' ||
          (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches));
      </script>
      {{ partialCached "css.html" . }}
      {{ partialCached "js.html" . }}
    </head>
    <body>
      <header>
        <fieldset id="theme-switcher">
          <legend>Color scheme</legend>
          <label><input type="radio" name="theme" value="light"> Light</label>
          <label><input type="radio" name="theme" value="dark"> Dark</label>
          <label><input type="radio" name="theme" value="system"> System</label>
        </fieldset>
      </header>
      <main>
        {{ block "main" . }}{{ end }}
      </main>
    </body>
  </html>
  ```

Step 7
: To verify the setup, add a fenced code block to your home page:

  ````md {file="content/_index.md" copy=true}
  ```go
  func printGreeting(showGreeting bool) {
    if showGreeting {
      fmt.Println("Hello, World!")
    }
  }
  ```
  ````

[`css.Build`]: /functions/css/build/
[`highlight`]: /shortcodes/highlight/
[`importContext`]: /functions/css/build/#importcontext
[`partials.IncludeCached`]: /functions/partials/includecached/
[`publishDir`]: /configuration/all/#publishdir
[`style`]: /configuration/markup/#style
[`transform.HighlightCodeBlock`]: /functions/transform/highlightcodeblock/
[`transform.Highlight`]: /functions/transform/highlight/
[global default]: /configuration/markup/#noclasses
[mode]: /quick-reference/syntax-highlighting-styles/#modes
[syntax highlighting styles]: /quick-reference/syntax-highlighting-styles/
