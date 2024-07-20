---
title: Exif
description: Applicable to JPEG, PNG, TIFF, and WebP images, returns an EXIF object containing image metadata.
categories: []
keywords: []
action:
  related: []
  returnType: exif.ExifInfo
  signatures: [RESOURCE.Exif]
toc: true
---

Applicable to JPEG, PNG, TIFF, and WebP images, the `Exif` method on an image `Resource` object returns an [EXIF] object containing image metadata.

## Methods

Date
: (`time.Time`) Returns the image creation date/time. Format with the [`time.Format`] function.

Lat
: (`float64`) Returns the GPS latitude in degrees.

Long
: (`float64`) Returns the GPS longitude in degrees.

Tags
: (`exif.Tags`) Returns a collection of the available EXIF tags for this image. You may include or exclude specific tags from this collection in the [site configuration].

## Examples

To list the creation date, location, and EXIF tags:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Exif }}
    <p>Date: {{ .Date }}</p>
    <p>Lat/Long: {{ .Lat }}/{{ .Long }}</p>
    {{ with .Tags }}
      <p>Tags</p>
      <table>
        <thead>
          <tr><th>Tag</th><th>Value</th></tr>
        </thead>
        <tbody>
          {{ range $k, $v := . }}
          <tr><td>{{ $k }}</td><td>{{ $v }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    {{ end }}
  {{ end }}
{{ end }}
```

To list specific values:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Exif }}
    <ul>
      {{ with .Date }}<li>Date: {{ .Format "January 02, 2006" }}</li>{{ end }}
      {{ with .Tags.ApertureValue }}<li>Aperture: {{ lang.FormatNumber 2 . }}</li>{{ end }}
      {{ with .Tags.BrightnessValue }}<li>Brightness: {{ lang.FormatNumber 2 . }}</li>{{ end }}
      {{ with .Tags.ExposureTime }}<li>Exposure Time: {{ . }}</li>{{ end }}
      {{ with .Tags.FNumber }}<li>F Number: {{ . }}</li>{{ end }}
      {{ with .Tags.FocalLength }}<li>Focal Length: {{ . }}</li>{{ end }}
      {{ with .Tags.ISOSpeedRatings }}<li>ISO Speed Ratings: {{ . }}</li>{{ end }}
      {{ with .Tags.LensModel }}<li>Lens Model: {{ . }}</li>{{ end }}
    </ul>
  {{ end }}
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

[exif]: https://en.wikipedia.org/wiki/Exif
[site configuration]: /content-management/image-processing/#exif-data
[`time.Format`]: /functions/time/format/
