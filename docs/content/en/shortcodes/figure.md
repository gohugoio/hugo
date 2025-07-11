---
title: Figure shortcode
linkTitle: Figure
description: Insert an HTML figure element into your content using the figure shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `figure` shortcode, copy the [source code] to a file with the same name in the `layouts/_shortcodes` directory.

## Example

With this markup:

```text
{{</* figure
  src="/images/examples/zion-national-park.jpg"
  alt="A photograph of Zion National Park"
  link="https://www.nps.gov/zion/index.htm"
  caption="Zion National Park"
  class="ma0 w-75"
*/>}}
```

Hugo renders this HTML:

```html
<figure class="ma0 w-75">
  <a href="https://www.nps.gov/zion/index.htm">
    <img 
      src="/images/examples/zion-national-park.jpg" 
      alt="A photograph of Zion National Park"
    >
  </a>
  <figcaption>
    <p>Zion National Park</p>
  </figcaption>
</figure>
```

Which looks like this in your browser:

{{< figure
  src="/images/examples/zion-national-park.jpg"
  alt="A photograph of Zion National Park"
  link="https://www.nps.gov/zion/index.htm"
  caption="Zion National Park"
  class="ma0 w-75"
>}}

## Arguments

src
: (`string`) The `src` attribute of the `img` element. Typically this is a [page resource](g) or a [global resource](g).

alt
: (`string`) The `alt` attribute of the `img` element.

width
: (`int`) The `width` attribute of the `img` element.

height
: (`int`) The `height` attribute of the `img` element.

loading
: (`string`) The `loading` attribute of the `img` element.

class
: (`string`) The `class` attribute of the `figure` element.

link
: (`string`) The `href` attribute of the anchor element that wraps the `img` element.

target
: (`string`) The `target` attribute of the anchor element that wraps the `img` element.

rel
: (`rel`) The `rel` attribute of the anchor element that wraps the `img` element.

title
: (`string`) Within the `figurecaption` element, the title is at the top, wrapped within an `h4` element.

caption
: (`string`) Within the `figurecaption` element, the caption is at the bottom and may contain plain text or markdown.

attr
: (`string`) Within the `figurecaption` element, the attribution appears next to the caption and may contain plain text or markdown.

attrlink
: (`string`) The `href` attribute of the anchor element that wraps the attribution.

## Image location

The `figure` shortcode resolves internal Markdown destinations by looking for a matching [page resource](g), falling back to a matching [global resource](g). Remote destinations are passed through, and the render hook will not throw an error or warning if unable to resolve a destination.

You must place global resources in the `assets` directory. If you have placed your resources in the `static` directory, and you are unable or unwilling to move them, you must mount the `static` directory to the `assets` directory by including both of these entries in your site configuration:

{{< code-toggle file=hugo >}}
[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'static'
target = 'assets'
{{< /code-toggle >}}

[source code]: {{% eturl figure %}}
