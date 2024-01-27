---
title: images.Text
description: Returns an image filter that adds text to an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: ['images.Text TEXT [OPTIONS]']
toc: true
---

## Options

Although none of the options are required, at a minimum you will want to set the `size` to be some reasonable percentage of the image height.

color
: (`string`) The font color, either a 3-digit or 6-digit hexadecimal color code. Default is `#ffffff` (white).

font
: (`resource.Resource`) The font can be a [global resource], a [page resource], or a [remote resource]. Default is [Go Regular], a proportional sans-serif TrueType font.

[Go Regular]: https://go.dev/blog/go-fonts#sans-serif

linespacing
: (`int`) The number of pixels between each line. For a line height of 1.4, set the `linespacing` to 0.4 multiplied by the `size`. Default is `2`.

size
: (`int`) The font size in pixels. Default is `20`.

x
: (`int`) The horizontal offset, in pixels, relative to the left of the image. Default is `10`.

y
: (`int`) The vertical offset, in pixels, relative to the top of the image. Default is `10`.

[global resource]: /getting-started/glossary/#global-resource
[page resource]: /getting-started/glossary/#page-resource
[remote resource]: /getting-started/glossary/#remote-resource

## Usage

Capture the font as a resource:

```go-html-template
{{ $font := "" }}
{{ $path := "https://github.com/google/fonts/raw/main/ofl/lato/Lato-Regular.ttf" }}
{{ with resources.GetRemote $path }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ $font = . }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get resource %q" $path }}
{{ end }}
```

Create the options map:

```go-html-template
{{ $opts := dict
  "color" "#fbfaf5"
  "font" $font
  "linespacing" 8
  "size" 40
  "x" 25
  "y" 190
}}
```

Set the text:

```go-html-template
{{ $text := "Zion National Park" }}
```

Create the filter:

```go-html-template
{{ $filter := images.Text $text $opts }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Text"
  filterArgs="Zion National Park,25,190,40,1.2,#fbfaf5"
  example=true
>}}
