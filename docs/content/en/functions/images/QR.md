---
title: images.QR
description: Encodes text to a QR code using the given error correction level, returning an image resource.
keywords: []
action:
  aliases: []
  related: []
  returnType: images.ImageResource
  signatures: ['images.QR TEXT [OPTIONS]']
toc: true
---

{{< new-in 0.141.0 >}}

The `images.QR` function encodes text to a [QR code] using the given error correction level. The generated image will always be at least 232x232 pixels, with each QR code module represented by 8 image pixels.

[QR code]: https://en.wikipedia.org/wiki/QR_code

The size of the generated image is variable and depends on two factors:

- Length of the encoded text: Longer text requires a larger image to encode all the information.
- Error correction level: Higher error correction levels lead to a slightly larger image size to ensure better readability even if the code is damaged.

See the [resizing](#resizing) section below if you wish to decrease the image size. Always test the rendered QR code after resizing, both on-screen and in print.

## Options

text
: (`string`) The text to encode.

level
: (`string`) The error correction level to use when encoding the text, one of `low`, `medium`, `quartile`, or `high`. The default value is sufficient for most applications.

Error correction level|Redundancy
:--|:--|:--
low|20%
medium (default)|38%
quartile|55%
high|65%

targetDir
: (`string`) The subdirectory within the [`publishDir`] where Hugo will place the generated image. If empty or not provided, the image is placed directly in the `publishDir` root. Hugo automatically creates the directory if it doesn't exist.

[`publishDir`]: /getting-started/configuration/#publishdir

## Examples

To render a QR code with the default error correction level:

```go-html-template
{{ $opts := dict "text" "https://gohugo.io" }}
{{ with images.QR $opts }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{< qr text="https://gohugo.io" class="w3" >}}

To render a QR code with a "high" error correction level and publish it to a "codes" directory within your `publishDir`:

```go-html-template
{{ $text := "https://gohugo.io" }}
{{ $opts := dict "text" $text "level" "high" "targetDir" "codes" }}
{{ with images.QR $opts }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{< qr text="https://gohugo.io" level=high targetDir=qr >}}

## Resizing

Use the [`Resize`] method to scale the generated image as needed.

[`Resize`]: /methods/resource/resize/

Due to the variable size of the generated image, do not specify a fixed width when resizing. Instead, calculate the new width by multiplying the original width by a scale factor. The scale factor must be a multiple of 0.125 to maintain an even number of pixels per QR code module.

```go-html-template
{{ $scaleFactor := 0.375 }}
{{ $opts := dict "text" "https://gohugo.io" }}
{{ with images.QR $opts }}
  {{ $width := .Width | mul $scaleFactor | math.Round | int }}
  {{ with .Resize (printf "%dx" $width) }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< qr text="https://gohugo.io" scaleFactor=0.375 >}}

As you decrease the size of a QR code, the maximum distance at which it can be reliably scanned by a device also decreases.

{{% note %}}
Always test the rendered QR code after resizing, both on-screen and in print.
{{% /note %}}
