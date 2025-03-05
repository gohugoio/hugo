---
title: Syntax highlighting
description: Add syntax highlighting to code examples.
categories: []
keywords: []
aliases: [/extras/highlighting/,/extras/highlight/,/tools/syntax-highlighting/]
---

Hugo provides several methods to add syntax highlighting to code examples:

- Use the [`transform.Highlight`] function within your templates
- Use the [`highlight`] shortcode with any [content format](g)
- Use fenced code blocks with the Markdown content format

[`transform.Highlight`]: /functions/transform/highlight/
[`highlight`]: /shortcodes/highlight/

## Fenced code blocks

In its default configuration, Hugo highlights code examples within fenced code blocks, following this form:

{{< code file=content/example.md lang=text >}}
```LANG [OPTIONS]
CODE
```
{{< /code >}}

CODE
: The code to highlight.

LANG
: The language of the code to highlight. Choose from one of the [supported languages]. This value is case-insensitive.

OPTIONS
: One or more space-separated or comma-separated key-value pairs wrapped in braces. Set default values for each option in your [site configuration]. The key names are case-insensitive.

[supported languages]: #languages
[site configuration]: /configuration/markup/#highlight

For example, with this Markdown:

{{< code file=content/example.md lang=text >}}
```go {linenos=inline hl_lines=[3,"6-8"] style=emacs}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
```
{{< /code >}}

Hugo renders this:

```go {linenos=inline, hl_lines=[3, "6-8"], style=emacs}
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        fmt.Println("Value of i:", i)
    }
}
```

## Options

{{% include "_common/syntax-highlighting-options.md" %}}

## Escaping

When documenting shortcode usage, escape the tag delimiters:

{{< code file=content/example.md lang=text >}}
```text {linenos=inline}
{{</*/* shortcode-1 */*/>}}

{{%/*/* shortcode-2 */*/%}}
```
{{< /code >}}

Hugo renders this to:

```text {linenos=inline}
{{</* shortcode-1 */>}}

{{%/* shortcode-2 */%}}
```

## Languages

These are the supported languages. Use one of the identifiers, not the language name, when specifying a language for:

- The [`transform.Highlight`] function
- The [`highlight`] shortcode
- Fenced code blocks

{{< chroma-lexers >}}
