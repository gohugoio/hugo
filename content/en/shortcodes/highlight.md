---
title: Highlight
description: Insert syntax-highlighted code into your content using the highlight shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{% note %}}
To override Hugo's embedded `highlight` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl highlight %}}
{{% /note %}}

{{% note %}}
With the Markdown [content format], the `highlight` shortcode is rarely needed because, by default, Hugo automatically applies syntax highlighting to fenced code blocks.

The primary use case for the `highlight` shortcode in Markdown is to apply syntax highlighting to inline code snippets.

[content format]: /content-management/formats/
{{% /note %}}

The `highlight` shortcode calls the [`transform.Highlight`] function which uses the [Chroma] syntax highlighter, supporting over 200 languages with more than 40 [available styles].

[chroma]: https://github.com/alecthomas/chroma
[available styles]: https://xyproto.github.io/splash/docs/
[`transform.Highlight`]: /functions/transform/highlight/

## Arguments

The `highlight` shortcode takes three arguments.

```text
{{</* highlight LANG OPTIONS */>}}
CODE
{{</* /highlight */>}}
```

CODE
: (`string`) The code to highlight.

LANG
: (`string`) The language of the code to highlight. Choose from one of the [supported languages]. This value is case-insensitive.

OPTIONS
: (`string`) Zero or more space-separated key-value pairs wrapped in quotation marks. Set default values for each option in your [site configuration]. The key names are case-insensitive.

[site configuration]: /getting-started/configuration-markup#highlight
[supported languages]: /content-management/syntax-highlighting/#list-of-chroma-highlighting-languages

## Example

```text
{{</* highlight go "linenos=inline, hl_Lines=3 6-8, style=emacs" */>}}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
{{</* /highlight */>}}
```

Hugo renders this to:

{{< highlight go "linenos=inline, hl_Lines=3 6-8, noClasses=true" >}}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
            fmt.Println("Value of i:", i)
    }
}
{{< /highlight >}}

You can also use the `highlight` shortcode for inline code snippets:

```text
This is some {{</* highlight go "hl_inline=true" */>}}fmt.Println("inline"){{</* /highlight */>}} code.
```

Hugo renders this to:

This is some {{< highlight go "hl_inline=true, noClasses=true" >}}fmt.Println("inline"){{< /highlight >}} code.

Given the verbosity of the example above, if you need to frequently highlight inline code snippets, create your own shortcode using a shorter name with preset options.

{{< code file=layouts/shortcodes/hl.html >}}
{{ $code := .Inner | strings.TrimSpace }}
{{ $lang := or (.Get 0) "go"  }}
{{ $opts := dict "hl_inline" true "noClasses" true }}
{{ transform.Highlight $code $lang $opts }}
{{< /code >}}

```text
This is some {{</* hl */>}}fmt.Println("inline"){{</* /hl */>}} code.
```

Hugo renders this to:

This is some {{< hl >}}fmt.Println("inline"){{< /hl >}} code.

## Options

Pass the options when calling the shortcode. You can set their default values in your [site configuration].

{{% include "functions/_common/highlighting-options" %}}
