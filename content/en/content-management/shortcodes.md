---
title: Shortcodes
description: Shortcodes are simple snippets inside your content files calling built-in or custom templates.
categories: [content management]
keywords: [markdown,content,shortcodes]
menu:
  docs:
    parent: content-management
    weight: 100
weight: 100
toc: true
aliases: [/extras/shortcodes/]
testparam: "Hugo Rocks!"
---

## What a shortcode is

Hugo loves Markdown because of its simple content format, but there are times when Markdown falls short. Often, content authors are forced to add raw HTML (e.g., video `<iframe>`'s) to Markdown content. We think this contradicts the beautiful simplicity of Markdown's syntax.

Hugo created **shortcodes** to circumvent these limitations.

A shortcode is a simple snippet inside a content file that Hugo will render using a predefined template. Note that shortcodes will not work in template files. If you need the type of drop-in functionality that shortcodes provide but in a template, you most likely want a [partial template][partials] instead.

In addition to cleaner Markdown, shortcodes can be updated any time to reflect new classes, techniques, or standards. At the point of site generation, Hugo shortcodes will easily merge in your changes. You avoid a possibly complicated search and replace operation.

## Use shortcodes

{{< youtube 2xkNJL4gJ9E >}}

In your content files, a shortcode can be called by calling `{{%/* shortcodename arguments */%}}`. Shortcode arguments are space delimited, and arguments with internal spaces must be quoted.

The first word in the shortcode declaration is always the name of the shortcode. Arguments follow the name. Depending upon how the shortcode is defined, the arguments may be named, positional, or both, although you can't mix argument types in a single call. The format for named arguments models that of HTML with the format `name="value"`.

Some shortcodes use or require closing shortcodes. Again like HTML, the opening and closing shortcodes match (name only) with the closing declaration, which is prepended with a slash.

Here are two examples of paired shortcodes:

```go-html-template
{{%/* mdshortcode */%}}Stuff to `process` in the *center*.{{%/* /mdshortcode */%}}
```

```go-html-template
{{</* highlight go */>}} A bunch of code here {{</* /highlight */>}}
```

The examples above use two different delimiters, the difference being the `%` character in the first and the `<>` characters in the second.

### Shortcodes with raw string arguments

You can pass multiple lines as arguments to a shortcode by using raw string literals:

```go-html-template
{{</*  myshortcode `This is some <b>HTML</b>,
and a new line with a "quoted string".` */>}}
```

### Shortcodes with Markdown

Shortcodes using the `%` as the outer-most delimiter will be fully rendered when sent to the content renderer. This means that the rendered output from a shortcode can be part of the page's table of contents, footnotes, etc.

### Shortcodes without Markdown

The `<` character indicates that the shortcode's inner content does *not* need further rendering. Often shortcodes without Markdown include internal HTML:

```go-html-template
{{</* myshortcode */>}}<p>Hello <strong>World!</strong></p>{{</* /myshortcode */>}}
```

### Nested shortcodes

You can call shortcodes within other shortcodes by creating your own templates that leverage the `.Parent` method. `.Parent` allows you to check the context in which the shortcode is being called. See [Shortcode templates][sctemps].

## Embedded shortcodes

See the [shortcodes](/shortcodes/) section.

## Privacy configuration

To learn how to configure your Hugo site to meet the new EU privacy regulation, see [privacy protections].

## Create custom shortcodes

To learn more about creating custom shortcodes, see the [shortcode template documentation].

[privacy protections]: /about/privacy/
[partials]: /templates/partial/
[quickstart]: /getting-started/quick-start/
[sctemps]: /templates/shortcode/
[shortcode template documentation]: /templates/shortcode/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/
