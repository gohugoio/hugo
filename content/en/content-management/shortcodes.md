---
title: Shortcodes
description: Use embedded, custom, or inline shortcodes to insert elements such as videos, images, and social media embeds into your content.
categories: []
keywords: []
aliases: [/extras/shortcodes/]
---

## Introduction

{{% glossary-term shortcode %}}

There are three types of shortcodes: embedded, custom, and inline.

## Embedded

Hugo's embedded shortcodes are pre-defined templates within the application. Refer to each shortcode's documentation for specific usage instructions and available arguments.

{{% list-pages-in-section path=/shortcodes %}}

## Custom

Create custom shortcodes to simplify and standardize content creation. For example, the following shortcode template generates an audio player using a [global resource](g):

```go-html-template {file="layouts/_shortcodes/audio.html"}
{{ with resources.Get (.Get "src") }}
  <audio controls preload="auto" src="{{ .RelPermalink }}"></audio>
{{ end }}
```

Then call the shortcode from within markup:

```text {file="content/example.md"}
{{</* audio src=/audio/test.mp3 */>}}
```

Learn more about creating shortcodes in the [shortcode templates] section.

## Inline

An inline shortcode is a shortcode template defined within content.

Hugo's security model is based on the premise that template and configuration authors are trusted, but content authors are not. This model enables generation of HTML output safe against code injection.

To conform with this security model, creating shortcode templates within content is disabled by default. If you trust your content authors, you can enable this functionality in your site's configuration:

{{< code-toggle file=hugo >}}
[security]
enableInlineShortcodes = true
{{< /code-toggle >}}

For more information see [configure security](/configuration/security).

The following example demonstrates an inline shortcode, `date.inline`, that accepts a single positional argument: a date/time [layout string].

```text {file="content/example.md"}
Today is
{{</* date.inline ":date_medium" */>}}
  {{- now | time.Format (.Get 0) -}}
{{</* /date.inline */>}}.

Today is {{</* date.inline ":date_full" /*/>}}.
```

In the example above, the inline shortcode is executed twice: once upon definition and again when subsequently called. Hugo renders this to:

```html
<p>Today is Jan 30, 2025.</p>
<p>Today is Thursday, January 30, 2025</p>
```

Inline shortcodes process their inner content within the same context as regular shortcode templates, allowing you to use any available [shortcode method].

> [!note]
> You cannot [nest](#nesting) inline shortcodes.

Learn more about creating shortcodes in the [shortcode templates] section.

## Calling

Shortcode calls involve three syntactical elements: tags, arguments, and notation.

### Tags

Some shortcodes expect content between opening and closing tags. For example, the embedded [`details`] shortcode requires an opening and closing tag:

```text
{{</* details summary="See the details" */>}}
This is a **bold** word.
{{</* /details */>}}
```

Some shortcodes do not accept content. For example, the embedded [`instagram`] shortcode requires a single _positional_ argument:

```text
{{</* instagram CxOWiQNP2MO */>}}
```

Some shortcodes optionally accept content. For example, you can call the embedded [`qr`] shortcode with content:

```text
{{</* qr */>}}
https://gohugo.io
{{</* /qr */>}}
```

Or use the self-closing syntax with a trailing slash to pass the text as an argument:

```text
{{</* qr text=https://gohugo.io /*/>}}
```

Refer to each shortcode's documentation for specific usage instructions and available arguments.

### Arguments

Shortcode arguments can be either _named_ or _positional_.

Named arguments are passed as case-sensitive key-value pairs, as seen in this example with the embedded [`figure`] shortcode. The `src` argument, for instance, is required.

```text
{{</* figure src=/images/kitten.jpg */>}}
```

Positional arguments, on the other hand, are determined by their position. The embedded `instagram` shortcode, for example, expects the first argument to be the Instagram post ID.

```text
{{</* instagram CxOWiQNP2MO */>}}
```

Shortcode arguments are space-delimited, and arguments with internal spaces must be quoted.

```text
{{</* figure src=/images/kitten.jpg alt="A white kitten" */>}}
```

Shortcodes accept [scalar](g) arguments, one of [string](g), [integer](g), [floating point](g), or [boolean](g).

```text
{{</* my-shortcode name="John Smith" age=24 married=false */>}}
```

You can optionally use multiple lines when providing several arguments to a shortcode for better readability:

```text
{{</* figure
  src=/images/kitten.jpg
  alt="A white kitten"
  caption="This is a white kitten"
  loading=lazy
*/>}}
```

Use a [raw string literal](g) if you need to pass a multiline string:

```text
{{</* myshortcode `This is some <b>HTML</b>,
and a new line with a "quoted string".` */>}}
```

Shortcodes can accept named arguments, positional arguments, or both, but you must use either named or positional arguments exclusively within a single shortcode call; mixing them is not allowed.

Refer to each shortcode's documentation for specific usage instructions and available arguments.

### Notation

Shortcodes can be called using two different notations, distinguished by their tag delimiters.

Notation|Example
:--|:--
Markdown|`{{%/* foo */%}} ## Section 1 {{%/* /foo */%}}`
Standard|`{{</* foo */>}} ## Section 2 {{</* /foo */>}}`

#### Markdown notation

Hugo processes the shortcode before the page content is rendered by the Markdown renderer. This means, for instance, that Markdown headings inside a Markdown-notation shortcode will be included when invoking the [`TableOfContents`] method on the `Page` object.

#### Standard notation

With standard notation, Hugo processes the shortcode separately, merging the output into the page content after Markdown rendering. This means, for instance, that Markdown headings inside a standard-notation shortcode will be excluded when invoking the `TableOfContents` method on the `Page` object.

By way of example, with this shortcode template:

```go-html-template {file="layouts/_shortcodes/foo.html"}
{{ .Inner }}
```

And this markdown:

```text {file="content/example.md"}
{{%/* foo */%}} ## Section 1 {{%/* /foo */%}}

{{</* foo */>}} ## Section 2 {{</* /foo */>}}
```

Hugo renders this HTML:

```html
<h2 id="heading">Section 1</h2>

## Section 2
```

In the above, "Section 1" will be included when invoking the `TableOfContents` method, while "Section 2" will not.

The shortcode author determines which notation to use. Consult each shortcode's documentation for specific usage instructions and available arguments.

## Nesting

Shortcodes (excluding [inline](#inline) shortcodes) can be nested, creating parent-child relationships. For example, a gallery shortcode might contain several image shortcodes:

```text {file="content/example.md"}
{{</* gallery class="content-gallery" */>}}
  {{</* image src="/images/a.jpg" */>}}
  {{</* image src="/images/b.jpg" */>}}
  {{</* image src="/images/c.jpg" */>}}
{{</* /gallery */>}}
```

The [shortcode templates][nesting] section provides a detailed explanation and examples.

[`details`]: /shortcodes/details
[`figure`]: /shortcodes/figure
[`instagram`]: /shortcodes/instagram
[`qr`]: /shortcodes/qr
[`TableOfContents`]: /methods/page/tableofcontents/
[layout string]: /functions/time/format/#layout-string
[nesting]: /templates/shortcode/#nesting
[shortcode method]: /templates/shortcode/#methods
[shortcode templates]: /templates/shortcode/
