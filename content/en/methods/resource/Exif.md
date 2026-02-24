---
title: Exif
description: Applicable to JPEG, PNG, TIFF, and WebP images, returns an object containing Exif metadata.
categories: []
keywords: ['metadata']
params:
  functions_and_methods:
    returnType: meta.ExifInfo
    signatures: [RESOURCE.Exif]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Applicable to JPEG, PNG, TIFF, and WebP images, the `Exif` method on an image `Resource` object returns an object containing [Exif][Exif_Definition] metadata.

To extract [Exif][Exif_Definition], [IPTC][IPTC_Definition], and [XMP][XMP_Definition] metadata, use the [`Meta`] method instead.

> [!note]
> Metadata is not preserved during image transformation. Use this method with the _original_ image resource to extract metadata from JPEG, PNG, TIFF, and WebP images.

## Methods

### Date

(`time.Time`) Returns the image creation date/time. Format with the [`time.Format`] function.

### Lat

(`float64`) Returns the GPS latitude in degrees from Exif metadata.

### Long

(`float64`) Returns the GPS longitude in degrees from Exif metadata.

### Tags

(`meta.Tags`) Returns a collection of available Exif fields for this image. Availability is determined by the [`includeFields`][] and [`excludeFields`][] settings in your project configuration.

## Examples

To list the creation date, latitude, and longitude:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Exif }}
    <pre>
      {{ printf "%-25s %v\n" "Date" .Date }}
      {{ printf "%-25s %v\n" "Latitude" .Lat }}
      {{ printf "%-25s %v\n" "Longitude" .Long }}
    </pre>
  {{ end }}
{{ end }}
```

To list the available Exif fields:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Exif }}
    <pre>
      {{ range $k, $v := .Tags -}}
        {{ printf "%-25s %v\n" $k $v }}
      {{ end }}
    </pre>
  {{ end }}
{{ end }}
```

To list specific Exif fields:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Exif }}
    <pre>
      {{ with .Tags.ApertureValue }}{{ printf "%-25s %v\n" "ApertureValue" . }}{{ end }}
      {{ with .Tags.BrightnessValue }}{{ printf "%-25s %v\n" "BrightnessValue" . }}{{ end }}
    </pre>
  {{ end }}
{{ end }}
```

[`excludeFields`]: /configuration/imaging/#excludefields
[`includeFields`]: /configuration/imaging/#includefields
[`Meta`]: /methods/resource/meta/
[`time.Format`]: /functions/time/format/
[Exif_Definition]: https://en.wikipedia.org/wiki/Exif
[IPTC_Definition]: https://en.wikipedia.org/wiki/IPTC_Information_Interchange_Model
[XMP_Definition]: https://en.wikipedia.org/wiki/Extensible_Metadata_Platform
