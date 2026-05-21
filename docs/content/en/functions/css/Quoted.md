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
> This function is only applicable to the `vars` option passed to the [`css.Build`][] or [`css.Sass`][] functions.

When passing a `vars` map to the `css.Sass` function, Hugo detects common typed CSS values such as `24px` or `#FF0000` using regular expression matching. If necessary, you can bypass automatic type inference by using the `css.Quoted` function to explicitly indicate that the value must be treated as a quoted string.

For the `css.Build` function, use `css.Quoted` to explicitly indicate that a value must be treated as a quoted string, most commonly for `font-family` names or the `content` property.

In the example below, we use `css.Quoted` to ensure the values for the `content` property are injected as strings.

```go-html-template
{{ $vars := dict
  "ol-li-after" ("6" | css.Quoted)
  "ul-li-after" ("7" | css.Quoted)
}}

{{ $opts := dict "vars" $vars "transpiler" "dartsass" }}
{{ with resources.Get "sass/main.scss" | css.Sass $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

Using the `hugo:vars` identifier in your stylesheet:

```scss
@use "hugo:vars" as h;

ol li::after {
  content: h.$ol-li-after;
}

ul li::after {
  content: h.$ul-li-after;
}
```

The resulting CSS contains quoted strings:

```css
ol li::after {
  content: "6";
}

ul li::after {
  content: "7";
}
```

[`css.Build`]: /functions/css/build/#vars
[`css.Sass`]: /functions/css/sass/#vars
