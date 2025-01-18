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

alignx
 {{< new-in 0.141.0 >}}
: (`string`) The horizontal alignment of the text relative to the horizontal offset, one of `left`, `center`, or `right`. Default is `left`.

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

Set the text and paths:

```go-html-template
{{ $text := "Zion National Park" }}
{{ $fontPath := "https://github.com/google/fonts/raw/main/ofl/lato/Lato-Regular.ttf" }}
{{ $imagePath := "images/original.jpg" }}
```

Capture the font as a resource:

```go-html-template
{{ $font := "" }}
{{ with try (resources.GetRemote $fontPath) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $font = . }}
  {{ else }}
    {{ errorf "Unable to get resource %s" $fontPath }}
  {{ end }}
{{ end }}
```

Create the filter, centering the text horizontally and vertically:

```go-html-template
{{ $r := "" }}
{{ $filter := "" }}
{{ with $r = resources.Get $imagePath }}
  {{ $opts := dict
    "alignx" "center"
    "color" "#fbfaf5"
    "font" $font
    "linespacing" 8
    "size" 60
    "x" (mul .Width 0.5 | int)
    "y" (mul .Height 0.5 | int)
  }}
  {{ $filter = images.Text $text $opts }}
{{ else }}
  {{ errorf "Unable to get resource %s" $imagePath }}
{{ end }}
```

Apply the filter using the [`images.Filter`] function:

```go-html-template
{{ with $r }}
  {{ with . | images.Filter $filter }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also apply the filter using the [`Filter`] method on a `Resource` object:

```go-html-template
{{ with $r }}
  {{ with .Filter $filter }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

[`images.Filter`]: /functions/images/filter/
[`Filter`]: /methods/resource/filter/

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Text"
  filterArgs="Zion National Park,25,190,40,1.2,#fbfaf5"
  example=true
>}}
