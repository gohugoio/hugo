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
> This function is only applicable to the [`vars`] option passed to the [`css.Sass`] function.

When passing a `vars` map to the `css.Sass` function, Hugo detects common typed CSS values such as `24px` or `#FF0000` using regular expression matching. If necessary, you can bypass automatic type inference by using the `css.Unquoted` function to explicitly indicate that the value must not be treated as a quoted string.

For example:

```scss {file="assets/sass/main.scss"}
@use "hugo:vars" as h;

h1 {
  font-size: h.$font-size-h1;
}

h2 {
  font-size: h.$font-size-h2;
}
```

```go-html-template {file="layouts/_partials/css.html"}
{{ $vars := dict
  "font_size_h1" ("72px * 0.500" | css.Unquoted)
  "font_size_h2" ("72px * 0.375" | css.Unquoted)
}}

{{ with resources.Get "sass/main.scss" }}
  {{ $opts := dict
    "enableSourceMap" hugo.IsDevelopment
    "outputStyle" (cond hugo.IsDevelopment "expanded" "compressed")
    "targetPath" "css/main.css"
    "transpiler" "dartsass"
    "vars" $vars
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

The Sass rules are transpiled to:

```css {file="public/css/main.css"}
h1 {
  font-size: 36px;
}

h2 {
  font-size: 27px;
}
```

[`css.Sass`]: /functions/css/sass/
[`vars`]: /functions/css/sass/#vars
