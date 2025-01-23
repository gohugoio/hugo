---
title: Details
description: Insert an HTML details element into your content using the details shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{< new-in 0.140.0 >}}

{{% note %}}
To override Hugo's embedded `details` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl details %}}
{{% /note %}}

## Example

With this markup:

```text
{{</* details summary="See the details" */>}}
This is a **bold** word.
{{</* /details */>}}
```

Hugo renders this HTML:

```html
<details>
  <summary>See the details</summary>
  <p>This is a <strong>bold</strong> word.</p>
</details>
```

Which looks like this in your browser:

{{< details summary="See the details" >}}
This is a **bold** word.
{{< /details >}}

## Parameters

summary
: (`string`) The content of the child `summary` element rendered from Markdown to HTML. Default is `Details`.

open
: (`bool`) Whether to initially display the content of the `details` element. Default is `false`.

class
: (`string`) The `class` attribute of the `details` element.

name
: (`string`) The `name` attribute of the `details` element.

title
: (`string`) The `title` attribute of the `details` element.

## Styling

Use CSS to style the `details` element, the `summary` element, and the content itself.

```css
/* target the details element */
details { }

/* target the summary element */
details > summary { }

/* target the children of the summary element */
details > summary > * { }

/* target the content */
details > :not(summary) { }
```
