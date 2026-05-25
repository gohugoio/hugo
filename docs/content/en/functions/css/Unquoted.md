---
title: css.Unquoted
description: Returns the given string, setting its data type to indicate that it must not be quoted when used in CSS.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: css.UnquotedString
    signatures: [css.Unquoted STRING]
---

<!-- Added in v0.111.0 -->

> [!note]
> This function is only applicable to the `vars` option passed to the [`css.Sass`][] function.

When passing a `vars` map to the `css.Sass` function, Hugo detects common typed CSS values such as `24px` or `#FF0000` using regular expression matching. If necessary, you can bypass automatic type inference by using the `css.Unquoted` function to explicitly indicate that the value must be treated as an unquoted string.

In the example below, we use `css.Unquoted` to ensure the value for the `font-family` property is injected without quotes.

```go-html-template
{{ $vars := dict
  "font-main" ("sans-serif" | css.Unquoted)
}}

{{ $opts := dict "vars" $vars "transpiler" "dartsass" }}
{{ with resources.Get "sass/main.scss" | css.Sass $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

Using the `hugo:vars` identifier in your stylesheet:

```scss
@use "hugo:vars" as h;

body {
  font-family: h.$font-main;
}
```

The resulting CSS contains an unquoted string:

```css
body {
  font-family: sans-serif;
}
```

[`css.Sass`]: /functions/css/sass/#vars
