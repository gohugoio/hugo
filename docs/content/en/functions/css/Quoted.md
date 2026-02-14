---
title: css.Quoted
description: Returns the given string, setting its data type to indicate that it must be quoted when used in CSS.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: css.QuotedString
    signatures: [css.Quoted STRING]
---

<!-- Added in v0.111.0 -->

> [!note]
> This function is only applicable to the [`vars`] option passed to the [`css.Sass`] function.

When passing a `vars` map to the `css.Sass` function, Hugo detects common typed CSS values such as `24px` or `#FF0000` using regular expression matching. If necessary, you can bypass automatic type inference by using the `css.Quoted` function to explicitly indicate that the value must be treated as a quoted string.

For example:

```scss {file="assets/sass/main.scss"}
@use "hugo:vars" as h;

ol li::after {
  content: h.$ol-li-after;
}

ul li::after {
  content: h.$ul-li-after;
}
```

```go-html-template {file="layouts/_partials/css.html"}
{{ $vars := dict
  "ol_li_after" ("6" | css.Quoted )
  "ul_li_after" ("7" | css.Quoted )
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

The Sass code is transpiled to:

```css {file="public/css/main.css"}
ol li::after {
  content: "6";
}

ul li::after {
  content: "7";
}
```

[`css.Sass`]: /functions/css/sass/
[`vars`]: /functions/css/sass/#vars
