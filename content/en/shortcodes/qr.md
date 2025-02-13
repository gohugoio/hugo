---
title: QR
description: Insert a QR code into your content using the qr shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{< new-in 0.141.0 />}}

{{% note %}}
To override Hugo's embedded `qr` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl qr %}}
{{% /note %}}

The `qr` shortcode encodes the given text into a [QR code] using the specified options and renders the resulting image.

Internally this shortcode calls the `images.QR` function. Please read the [related documentation] for implementation details and guidance.

[QR code]: https://en.wikipedia.org/wiki/QR_code
[related documentation]: /functions/images/qr/

## Examples

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

Both of the above produce this image:

{{< qr text="https://gohugo.io" class="qrcode" targetDir="images/qr" />}}

To create a QR code for a phone number:

```text
{{</* qr text="tel:+12065550101" /*/>}}
```

{{< qr text="tel:+12065550101" class="qrcode" targetDir="images/qr" />}}

To create a QR code containing contact information in the [vCard] format:

[vCard]: https://en.wikipedia.org/wiki/VCard

```text
{{</* qr level="low" scale=2 alt="QR code of vCard for John Smith" */>}}
BEGIN:VCARD
VERSION:2.1
N;CHARSET=UTF-8:Smith;John;R.;Dr.;PhD
FN;CHARSET=UTF-8:Dr. John R. Smith, PhD.
ORG;CHARSET=UTF-8:ABC Widgets
TITLE;CHARSET=UTF-8:Vice President Engineering
TEL;TYPE=WORK:+12065550101
EMAIL;TYPE=WORK:jsmith@example.org
END:VCARD
{{</* /qr */>}}
```

{{< qr level="low" scale=2 alt="QR code of vCard for John Smith" class="qrcode" targetDir="images/qr" >}}
BEGIN:VCARD
VERSION:2.1
N;CHARSET=UTF-8:Smith;John;R.;Dr.;PhD
FN;CHARSET=UTF-8:Dr. John R. Smith, PhD.
ORG;CHARSET=UTF-8:ABC Widgets
TITLE;CHARSET=UTF-8:Vice President Engineering
TEL;TYPE=WORK:+12065550101
EMAIL;TYPE=WORK:jsmith@example.org
END:VCARD
{{< /qr >}}

## Parameters

text
: (`string`) The text to encode, falling back to the text between the opening and closing shortcode tags.

level
: (`string`) The error correction level to use when encoding the text, one of `low`, `medium`, `quartile`, or `high`. Default is `medium`.

scale
: (`int`) The number of image pixels per QR code module. Must be greater than or equal to 2. Default is `4`.

targetDir
: (`string`) The subdirectory within the [`publishDir`] where Hugo will place the generated image.

[`publishDir`]: /getting-started/configuration/#publishdir

alt
: (`string`) The `alt` attribute of the `img` element.

class
: (`string`) The `class` attribute of the `img` element.

id
: (`string`) The `id` attribute of the `img` element.

title
: (`string`) The `title` attribute of the `img` element.
