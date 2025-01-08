---
title: images.QR
description: Encodes the given text into a QR code using the specified options, returning an image resource.
keywords: []
action:
  aliases: []
  related: []
  returnType: images.ImageResource
  signatures: ['images.QR TEXT OPTIONS']
toc: true
math: true
---

{{< new-in 0.141.0 >}}

The `images.QR` function encodes the given text into a [QR code] using the specified options, returning an image resource. The size of the generated image depends on three factors:

- Data length: Longer text necessitates a larger image to accommodate the increased information density.
- Error correction level: Higher error correction levels enhance the QR code's resistance to damage, but this typically results in a slightly larger image size to maintain readability.
- Pixels per module: The number of image pixels assigned to each individual module (the smallest unit of the QR code) directly impacts the overall image size. A higher pixel count per module leads to a larger, higher-resolution image.

Although the default option values are sufficient for most applications, you should test the rendered QR code both on-screen and in print.

[QR code]: https://en.wikipedia.org/wiki/QR_code

## Options

level
: (`string`) The error correction level to use when encoding the text, one of `low`, `medium`, `quartile`, or `high`. Default is `medium`.

Error correction level|Redundancy
:--|:--|:--
low|20%
medium|38%
quartile|55%
high|65%

scale
: (`int`) The number of image pixels per QR code module. Must be greater than or equal to `2`. Default is `4`.

targetDir
: (`string`) The subdirectory within the [`publishDir`] where Hugo will place the generated image. Use Unix-style slashes (`/`) to separarate path segments. If empty or not provided, the image is placed directly in the `publishDir` root. Hugo automatically creates the necessary subdirectories if they don't exist.

[`publishDir`]: /getting-started/configuration/#publishdir

## Examples

To create a QR code using the default values for `level` and `scale`:

```go-html-template
{{ $text := "https://gohugo.io" }}
{{ with images.QR $text }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{< qr text="https://gohugo.io" class="qrcode" />}}

Specify `level`, `scale`, and `targetDir` as needed to achieve the desired result:

```go-html-template
{{ $text := "https://gohugo.io" }}
{{ $opts := dict 
  "level" "high" 
  "scale" 3
  "targetDir" "codes"
}}
{{ with images.QR $text $opts }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{< qr text="https://gohugo.io" level="high" scale=3 targetDir="codes" class="qrcode" />}}

## Scale

As you decrease the size of a QR code, the maximum distance at which it can be reliably scanned by a device also decreases.

In the example above, we set the `scale` to `2`, resulting in a QR code where each module consists of 2x2 pixels. While this might be sufficient for on-screen display, it's likely to be problematic when printed at 600 dpi.

\[ \frac{2\:px}{module} \times \frac{1\:inch}{600\:px} \times \frac{25.4\:mm}{1\:inch} = \frac{0.085\:mm}{module} \]

This module size is half of the commonly recommended minimum of 0.170 mm.\
If the QR code will be printed, use the default `scale` value of `4` pixels per module.

Avoid using Hugo's image processing methods to resize QR codes. Resizing can introduce blurring due to anti-aliasing when a QR code module occupies a fractional number of pixels.

{{% note %}}
Always test the rendered QR code both on-screen and in print.
{{% /note %}}

## Shortcode

Call the `qr` shortcode to insert a QR code into your content.

Use the self-closing syntax to pass the text as an argument:

```text
{{</* qr text="https://gohugo.io" /*/>}}
```

Or insert the text between the opening and closing tags:

```text
{{</* qr */>}}
https://gohugo.io
{{</* /qr */>}}
```

The `qr` shortcode accepts several arguments including `level` and `scale`. See the [related documentation] for details.

[related documentation]: /content-management/shortcodes/#qr
