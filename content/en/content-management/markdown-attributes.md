---
title: Markdown attributes
description: Use Markdown attributes to add HTML attributes when rendering Markdown to HTML.
categories: [content management]
keywords: [goldmark,markdown]
menu:
  docs:
    parent: content-management
    weight: 240
weight: 240
toc: true
---

## Overview

Hugo supports Markdown attributes on images and block elements including blockquotes, fenced code blocks, headings, horizontal rules, lists, paragraphs, and tables.

For example:

```text
This is a paragraph.
{class="foo bar" id="baz"}
```

With `class` and `id` you can use shorthand notation:

```text
This is a paragraph.
{.foo .bar #baz}
```

Hugo renders both of these to:

```html
<p class="foo bar" id="baz">This is a paragraph.</p>
```

## Block elements

Update your site configuration to enable Markdown attributes for block-level elements.

{{< code-toggle file=hugo >}}
[markup.goldmark.parser.attribute]
title = true # default is true
block = true # default is false
{{< /code-toggle >}}


## Standalone images

By default, when the [Goldmark] Markdown renderer encounters a standalone image element (no other elements or text on the same line), it wraps the image element within a paragraph element per the [CommonMark specification].

[CommonMark specification]: https://spec.commonmark.org/current/
[Goldmark]: https://github.com/yuin/goldmark

If you were to place an attribute list beneath an image element, Hugo would apply the attributes to the surrounding paragraph, not the image.

To apply attributes to a standalone image element, you must disable the default wrapping behavior:

{{< code-toggle file=hugo >}}
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = false # default is true
{{< /code-toggle >}}

## Usage

You may add [global HTML attributes], or HTML attributes specific to the current element type. Consistent with its content security model, Hugo removes HTML event attributes such as `onclick` and `onmouseover`.

[global HTML attributes]: https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes

The attribute list consists of one or more key-value pairs, separated by spaces or commas, wrapped by braces. You must quote string values that contain spaces. Unlike HTML, boolean attributes must have both key and value.

For example:

```text
> This is a blockquote.
{class="foo bar" hidden=hidden}
```

Hugo renders this to:

```html
<blockquote class="foo bar" hidden="hidden">
  <p>This is a blockquote.</p>
</blockquote>
```

In most cases, place the attribute list beneath the markup element. For headings and fenced code blocks, place the attribute list on the right.

Element|Position of attribute list
:--|:--
blockquote | bottom
fenced code block | right
heading | right
horizontal rule | bottom
image | bottom
list  | bottom
paragraph | bottom
table | bottom

For example:

````text
## Section 1 {class=foo}

```bash {class=foo linenos=inline}
declare a=1
echo "${a}"
```

This is a paragraph.
{class=foo}
````

As shown above, the attribute list for fenced code blocks is not limited to HTML attributes. You can also configure syntax highlighting by passing one or more of [these options](/functions/transform/highlight/#options).
